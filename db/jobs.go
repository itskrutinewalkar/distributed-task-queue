package db

import (
	"database/sql"
	"time"

	"github.com/krutinewalkar/distributed-task-queue/models"
)

func InsertJob(db *sql.DB, job models.Job) error {
	query := `
        INSERT INTO jobs (id, queue, payload, status, attempts, max_retries, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := db.Exec(query,
		job.ID,
		job.Queue,
		job.Payload,
		job.Status,
		job.Attempts,
		job.MaxRetries,
		job.CreatedAt,
		job.UpdatedAt,
	)
	return err
}

func UpdateJobStatus(db *sql.DB, id string, status models.Status, errMsg string) error {
	query := `
        UPDATE jobs
        SET status = $1, error = $2, updated_at = $3
        WHERE id = $4
    `
	var nullErr sql.NullString

	if errMsg != "" {
        nullErr = sql.NullString{String: errMsg, Valid: true}
    }

	_, err := db.Exec(query, status, nullErr, time.Now(), id)
	return err
}

func GetJob(db *sql.DB, id string) (models.Job, error) {
	query := `SELECT id, queue, payload, status, attempts, max_retries, error, created_at, updated_at FROM jobs WHERE id = $1`
	row := db.QueryRow(query, id)

	var job models.Job
	err := row.Scan(
		&job.ID,
		&job.Queue,
		&job.Payload,
		&job.Status,
		&job.Attempts,
		&job.MaxRetries,
		&job.Error,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	return job, err
}
