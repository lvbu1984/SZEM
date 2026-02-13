package s3

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/weihaoli/szem/pkg/core"
	"github.com/weihaoli/szem/pkg/model"
	"github.com/weihaoli/szem/pkg/repo"
	"github.com/weihaoli/szem/pkg/storage"
)

func Handle(w http.ResponseWriter, r *http.Request, s storage.Storage, objects repo.ObjectRepo, jobs repo.JobRepo) {

	if r.Method == http.MethodGet {
		handleGET(w, r, s, objects)
		return
	}

	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// ===== Build RequestFacts (给 validatePUT 用) =====
	facts := model.RequestFacts{
		Method: r.Method,
	}

	_, vErr := validatePUT(r, facts)
	if vErr != nil {
		msg := vErr.Error()
		if msg == "missing Content-Length" {
			writeError(w, http.StatusLengthRequired, msg)
			return
		}
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	// ===== Body Limit =====
	if r.ContentLength > MaxPutObjectSizeBytes {
		writeError(w, http.StatusBadRequest, "EntityTooLarge")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxPutObjectSizeBytes)
	defer r.Body.Close()

	// ===== Parse path =====
	bucket, key, err := parsePath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	// ===== Stage to local (temporary only) =====
	objectID := uuid.NewString()
	stagingPath := "./data/staging/" + objectID

	if err := core.WriteToFile(stagingPath, r.Body); err != nil {
		writeError(w, http.StatusInternalServerError, "stage write failed: "+err.Error())
		return
	}

	// ===== Metadata: committing =====
	obj := repo.Object{
		ObjectID:    objectID,
		Bucket:      bucket,
		Key:         key,
		Size:        r.ContentLength,
		Status:      repo.ObjectCommitting,
		StagingPath: stagingPath,
	}
	if err := objects.CreateCommitting(r.Context(), obj); err != nil {
		_ = os.Remove(stagingPath)
		writeError(w, http.StatusInternalServerError, "db object create failed: "+err.Error())
		return
	}

	// ===== Enqueue job =====
	jobID, err := jobs.Enqueue(r.Context(), objectID)
	if err != nil {
		_ = objects.MarkFailed(r.Context(), objectID, "enqueue job failed: "+err.Error())
		writeError(w, http.StatusInternalServerError, "enqueue job failed: "+err.Error())
		return
	}

	// ===== Return 202 Accepted (explicit async) =====
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"object_id": objectID,
		"job_id":    jobID,
		"status":    "committing",
	})
}

func handleGET(w http.ResponseWriter, r *http.Request, s storage.Storage, objects repo.ObjectRepo) {

	bucket, key, err := parsePath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	// 查 metadata 状态
	obj, err := objects.GetByPath(r.Context(), bucket, key)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	switch obj.Status {
	case repo.ObjectCommitting:
		writeError(w, http.StatusConflict, "object is committing (async), try later")
		return
	case repo.ObjectFailed:
		writeError(w, http.StatusInternalServerError, "object failed: "+obj.LastError)
		return
	case repo.ObjectAvailable:
		// OK
	default:
		writeError(w, http.StatusInternalServerError, "unknown object status")
		return
	}

	reader, err := s.Fetch(r.Context(), bucket, key)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	defer reader.Close()

	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, reader)
}

func parsePath(path string) (string, string, error) {
	if len(path) < 2 {
		return "", "", http.ErrNotSupported
	}
	path = path[1:]
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", http.ErrNotSupported
	}
	return parts[0], parts[1], nil
}

