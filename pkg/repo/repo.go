package repo

import "context"

type ObjectRepo interface {
	Init(ctx context.Context) error

	CreateCommitting(ctx context.Context, obj Object) error
	GetByPath(ctx context.Context, bucket, key string) (*Object, error)
	GetByID(ctx context.Context, objectID string) (*Object, error)

	MarkAvailable(ctx context.Context, objectID string) error
	MarkFailed(ctx context.Context, objectID string, lastError string) error
}

type JobRepo interface {
	Init(ctx context.Context) error

	Enqueue(ctx context.Context, objectID string) (string, error)
	FetchOnePending(ctx context.Context) (*Job, error)

	MarkRunning(ctx context.Context, jobID string) error
	MarkDone(ctx context.Context, jobID string) error
	MarkFailed(ctx context.Context, jobID string, lastError string) error
}

