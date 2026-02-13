package repo

import "time"

type ObjectStatus string

const (
	ObjectCommitting ObjectStatus = "committing"
	ObjectAvailable  ObjectStatus = "available"
	ObjectFailed     ObjectStatus = "failed"
)

type Object struct {
	ObjectID     string
	Bucket       string
	Key          string
	Size         int64
	Status       ObjectStatus
	StagingPath  string
	LastError    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type JobStatus string

const (
	JobPending JobStatus = "pending"
	JobRunning JobStatus = "running"
	JobDone    JobStatus = "done"
	JobFailed  JobStatus = "failed"
)

type Job struct {
	ID        string
	ObjectID  string
	Status    JobStatus
	Attempts  int
	LastError string
	CreatedAt time.Time
	UpdatedAt time.Time
}

