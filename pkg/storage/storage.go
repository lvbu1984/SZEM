package storage

import (
	"context"
	"io"
)

type Storage interface {

	// Store stores object data into backend
	Store(ctx context.Context, bucket, key string, r io.Reader, size int64) error

	// Fetch returns reader for object
	Fetch(ctx context.Context, bucket, key string) (io.ReadCloser, error)

	// Delete permanently removes object
	Delete(ctx context.Context, bucket, key string) error
}

