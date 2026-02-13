package model

import "net/http"

// request.go
//
// Phase 1 frozen responsibility:
// - Define request facts (bucket, key, headers)
//
// Forbidden in this file:
// - Any logic or behavior

type RequestFacts struct {
	Method   string
	Path     string
	Bucket   string
	ObjectKey string
	Headers  http.Header
}

