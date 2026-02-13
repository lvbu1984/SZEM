package local

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/weihaoli/szem/pkg/storage"
)

type LocalStorage struct {
	BasePath string
}

func New(base string) storage.Storage {
	return &LocalStorage{BasePath: base}
}

func (l *LocalStorage) Store(ctx context.Context, bucket, key string, r io.Reader, size int64) error {

	path := filepath.Join(l.BasePath, bucket)
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

func (l *LocalStorage) Fetch(ctx context.Context, bucket, key string) (io.ReadCloser, error) {

	fullPath := filepath.Join(l.BasePath, bucket, key)
	return os.Open(fullPath)
}

func (l *LocalStorage) Delete(ctx context.Context, bucket, key string) error {

	fullPath := filepath.Join(l.BasePath, bucket, key)
	return os.Remove(fullPath)
}

