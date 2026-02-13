package s3ingress

import (
	"errors"
	"net/http"
	"strings"

	"github.com/your-org/szem/internal/model"
)

// ingress.go
//
// Phase 1 frozen responsibility:
// - Receive HTTP requests
// - Parse request line (method / path)
// - Extract bucket and object-key
//
// Forbidden in this file:
// - Reading request body
// - Any storage or metadata operation
// - Any Core or MK20 logic

var (
	ErrInvalidPath = errors.New("invalid S3 path, expected /{bucket}/{object-key}")
)

func parseRequestFacts(r *http.Request) (model.RequestFacts, error) {
	f := model.RequestFacts{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: r.Header,
	}

	// Path-style only: /{bucket}/{object-key}
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return model.RequestFacts{}, ErrInvalidPath
	}

	f.Bucket = parts[0]
	f.ObjectKey = parts[1]
	return f, nil
}

