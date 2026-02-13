package mk20

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalClient struct {
	BasePath string
}

func NewLocalClient(base string) *LocalClient {
	return &LocalClient{BasePath: base}
}

func (c *LocalClient) Store(ctx context.Context, bucket, key string, r io.Reader) error {

	path := filepath.Join(c.BasePath, bucket)

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	fullPath := filepath.Join(path, key)

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}

