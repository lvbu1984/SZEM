package s3ingress

import "net/http"

// respond.go
//
// Phase 1 frozen responsibility:
// - Translate internal results to S3-compatible HTTP responses
// - Set status code and headers (ETag, Content-Length)
//
// Forbidden in this file:
// - Any business logic
// - Any data persistence
// - Any retry or commit logic

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(msg + "\n"))
}

func writeNotImplemented(w http.ResponseWriter) {
	writeError(w, http.StatusNotImplemented,
		"not implemented: Phase 1 only freezes architecture and S3 ingress")
}

