package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/weihaoli/szem/pkg/repo"
)

type ObjectRepo struct {
	db *sql.DB
}

func NewObjectRepo(db *sql.DB) *ObjectRepo {
	return &ObjectRepo{db: db}
}

func (r *ObjectRepo) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS objects (
	object_id TEXT PRIMARY KEY,
	bucket TEXT NOT NULL,
	key TEXT NOT NULL,
	size INTEGER NOT NULL,
	status TEXT NOT NULL,
	staging_path TEXT NOT NULL,
	last_error TEXT NOT NULL DEFAULT '',
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	UNIQUE(bucket, key)
);`)
	return err
}

func (r *ObjectRepo) CreateCommitting(ctx context.Context, obj repo.Object) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
INSERT INTO objects (object_id, bucket, key, size, status, staging_path, last_error, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, '', ?, ?)
ON CONFLICT(bucket, key) DO UPDATE SET
	object_id=excluded.object_id,
	size=excluded.size,
	status=excluded.status,
	staging_path=excluded.staging_path,
	last_error='',
	updated_at=excluded.updated_at
;`,
		obj.ObjectID, obj.Bucket, obj.Key, obj.Size, repo.ObjectCommitting, obj.StagingPath, now, now,
	)
	return err
}

func (r *ObjectRepo) GetByPath(ctx context.Context, bucket, key string) (*repo.Object, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT object_id, bucket, key, size, status, staging_path, last_error, created_at, updated_at
FROM objects WHERE bucket=? AND key=?`,
		bucket, key,
	)

	var o repo.Object
	var status string
	if err := row.Scan(&o.ObjectID, &o.Bucket, &o.Key, &o.Size, &status, &o.StagingPath, &o.LastError, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, err
	}
	o.Status = repo.ObjectStatus(status)
	return &o, nil
}

func (r *ObjectRepo) GetByID(ctx context.Context, objectID string) (*repo.Object, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT object_id, bucket, key, size, status, staging_path, last_error, created_at, updated_at
FROM objects WHERE object_id=?`,
		objectID,
	)

	var o repo.Object
	var status string
	if err := row.Scan(&o.ObjectID, &o.Bucket, &o.Key, &o.Size, &status, &o.StagingPath, &o.LastError, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, err
	}
	o.Status = repo.ObjectStatus(status)
	return &o, nil
}

func (r *ObjectRepo) MarkAvailable(ctx context.Context, objectID string) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE objects SET status=?, last_error='', updated_at=? WHERE object_id=?`,
		repo.ObjectAvailable, time.Now(), objectID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("object not found")
	}
	return nil
}

func (r *ObjectRepo) MarkFailed(ctx context.Context, objectID string, lastError string) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE objects SET status=?, last_error=?, updated_at=? WHERE object_id=?`,
		repo.ObjectFailed, lastError, time.Now(), objectID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("object not found")
	}
	return nil
}

