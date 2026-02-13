package pdp

import (
	"context"
	"errors"
	"io"

	"github.com/weihaoli/szem/pkg/storage"
)

type PDPStorage struct {
	Endpoint string
	APIKey   string
}

func New(endpoint, apiKey string) storage.Storage {
	return &PDPStorage{
		Endpoint: endpoint,
		APIKey:   apiKey,
	}
}

func (p *PDPStorage) Store(ctx context.Context, bucket, key string, r io.Reader, size int64) error {
	return errors.New("PDP Store not implemented yet")
}

func (p *PDPStorage) Fetch(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	return nil, errors.New("PDP Fetch not implemented yet")
}

func (p *PDPStorage) Delete(ctx context.Context, bucket, key string) error {
	return errors.New("PDP Delete not implemented yet")
}

