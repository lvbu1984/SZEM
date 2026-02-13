package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/weihaoli/szem/pkg/repo"
)

type JobRepo struct {
	db *sql.DB
}

func NewJobRepo(db *sql.DB) *JobRepo {
	return &JobRepo{db: db}
}

func (r *JobRepo) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS jobs (
	id TEXT PRIMARY KEY,
	object_id TEXT NOT NULL,
	status TEXT NOT NULL,
	attempts INTEGER NOT NULL DEFAULT 0,
	last_error TEXT NOT NULL DEFAULT '',
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
`)
	return err
}

func (r *JobRepo) Enqueue(ctx context.Context, objectID string) (string, error) {
	id := uuid.NewString()
	now := time.Now()

	_, err := r.db.ExecContext(ctx, `
INSERT INTO jobs (id, object_id, status, attempts, last_error, created_at, updated_at)
VALUES (?, ?, ?, 0, '', ?, ?)`,
		id, objectID, repo.JobPending, now, now,
	)
	return id, err
}

func (r *JobRepo) FetchOnePending(ctx context.Context) (*repo.Job, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, object_id, status, attempts, last_error, created_at, updated_at
FROM jobs
WHERE status=?
ORDER BY created_at ASC
LIMIT 1`,
		repo.JobPending,
	)

	var j repo.Job
	var status string
	if err := row.Scan(&j.ID, &j.ObjectID, &status, &j.Attempts, &j.LastError, &j.CreatedAt, &j.UpdatedAt); err != nil {
		return nil, err
	}
	j.Status = repo.JobStatus(status)
	return &j, nil
}

func (r *JobRepo) MarkRunning(ctx context.Context, jobID string) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE jobs
SET status=?, attempts=attempts+1, updated_at=?
WHERE id=?`,
		repo.JobRunning, time.Now(), jobID,
	)
	return err
}

func (r *JobRepo) MarkDone(ctx context.Context, jobID string) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE jobs SET status=?, last_error='', updated_at=? WHERE id=?`,
		repo.JobDone, time.Now(), jobID,
	)
	return err
}

func (r *JobRepo) MarkFailed(ctx context.Context, jobID string, lastError string) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE jobs SET status=?, last_error=?, updated_at=? WHERE id=?`,
		repo.JobFailed, lastError, time.Now(), jobID,
	)
	return err
}

