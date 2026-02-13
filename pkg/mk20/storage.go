package mk20

import (
	"context"
	"io"
)

type Storage interface {
	Store(ctx context.Context, bucket, key string, r io.Reader) error
}

