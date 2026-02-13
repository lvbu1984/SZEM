package s3

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/weihaoli/szem/pkg/model"
)

// validate.go
//
// Phase 1 frozen responsibility:
// - Validate S3 request semantics
// - Handle Expect: 100-continue (only as a header check in Phase 1)
// - Reject invalid requests early
//
// Forbidden in this file:
// - Reading request body
// - Creating any persistent state
// - Any MK20 interaction

var (
	ErrMethodNotAllowed   = errors.New("only PUT is supported in Phase 1")
	ErrContentLengthMiss  = errors.New("missing Content-Length")
	ErrContentLengthBad   = errors.New("invalid Content-Length")
	ErrContentLengthNeg   = errors.New("negative Content-Length")
)

type ValidateResult struct {
	ContentLength int64
	ContentType   string
	Expect100     bool
}

func validatePUT(r *http.Request, facts model.RequestFacts) (ValidateResult, error) {
	if facts.Method != http.MethodPut {
		return ValidateResult{}, ErrMethodNotAllowed
	}

	cl := r.Header.Get("Content-Length")
	if strings.TrimSpace(cl) == "" {
		return ValidateResult{}, ErrContentLengthMiss
	}

	n, err := strconv.ParseInt(cl, 10, 64)
	if err != nil {
		return ValidateResult{}, ErrContentLengthBad
	}
	if n < 0 {
		return ValidateResult{}, ErrContentLengthNeg
	}

	ct := r.Header.Get("Content-Type")
	expect := strings.EqualFold(strings.TrimSpace(r.Header.Get("Expect")), "100-continue")

	return ValidateResult{
		ContentLength: n,
		ContentType:   ct,
		Expect100:     expect,
	}, nil
}

