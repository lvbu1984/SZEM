package s3ingress

import (
	"net/http"
	"strings"
)

// handler.go
//
// Phase 1 frozen responsibility:
// - Provide a single HTTP handler that performs Ingress + Validate
// - Do NOT read body and do NOT call Core/MK20 yet
//
// Forbidden in this file:
// - Any Stage/Commit logic
// - Any success response for PUT

type handler struct{}

func NewHandler() http.Handler {
	return &handler{}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Minimal health probe (useful for ops, no coupling to S3 semantics)
	if r.Method == http.MethodGet && (r.URL.Path == "/health" || r.URL.Path == "/") {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
		return
	}

	facts, err := parseRequestFacts(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, vErr := validatePUT(r, facts)
	if vErr != nil {
		// Map Phase 1 errors to minimal HTTP codes
		msg := vErr.Error()
		if strings.Contains(msg, "Content-Length") {
			// 411 is more specific for missing length; otherwise 400
			if msg == "missing Content-Length" {
				writeError(w, http.StatusLengthRequired, msg)
			} else {
				writeError(w, http.StatusBadRequest, msg)
			}
			return
		}

		// Method not allowed etc.
		writeError(w, http.StatusMethodNotAllowed, msg)
		return
	}

	// Phase 1: we explicitly do NOT proceed to Stage/Commit,
	// therefore we must NOT accept body and must NOT return 200.
	// Returning 501 here is correct and prevents large uploads.
	writeNotImplemented(w)
}

