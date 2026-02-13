package s3

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/weihaoli/szem/pkg/mk20"
	"github.com/weihaoli/szem/pkg/model"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	// ===== GET =====
	if r.Method == http.MethodGet {
		handleGET(w, r)
		return
	}

	// ===== PUT only =====
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	facts := model.RequestFacts{
		Method: r.Method,
		Path:   r.URL.Path,
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

	cleanPath := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.SplitN(cleanPath, "/", 2)
	if len(parts) != 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	bucket := parts[0]
	key := parts[1]

	if r.ContentLength > MaxPutObjectSizeBytes {
		writeError(w, http.StatusBadRequest, "EntityTooLarge")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxPutObjectSizeBytes)
	defer r.Body.Close()

	var storage mk20.Storage = mk20.NewLocalClient("./data")

	err := storage.Store(context.Background(), bucket, path.Clean(key), r.Body)
	if err != nil {

		if strings.Contains(err.Error(), "http: request body too large") {
			writeError(w, http.StatusBadRequest, "EntityTooLarge")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGET(w http.ResponseWriter, r *http.Request) {

	cleanPath := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.SplitN(cleanPath, "/", 2)
	if len(parts) != 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	bucket := parts[0]
	key := parts[1]

	fullPath := path.Join("./data", bucket, key)

	http.ServeFile(w, r, fullPath)
}

