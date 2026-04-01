package models

import (
	"database/sql"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Job struct {
	ID         string         `json:"id"` //struct tags for JSON serialization/deserialization to/from JSON format
	Queue      string         `json:"queue"`
	Payload    string         `json:"payload"`
	Status     Status         `json:"status"`
	Attempts   int            `json:"attempts"`
	MaxRetries int            `json:"max_retries"`
	Error      sql.NullString `json:"error,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
