package queue

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/krutinewalkar/distributed-task-queue/db"
	"github.com/krutinewalkar/distributed-task-queue/models"
)

func Enqueue(database *sql.DB, rdb *redis.Client, queueName string, payload string) (models.Job, error) {
	job := models.Job{
		ID:         uuid.New().String(),
		Queue:      queueName,
		Payload:    payload,
		Status:     models.StatusPending,
		Attempts:   0,
		MaxRetries: 3,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := SaveJob(database, rdb, job)
	if err != nil {
		return models.Job{}, err
	}
	return job, nil
}

//save to database and add to Redis queue
func SaveJob(database *sql.DB, rdb *redis.Client, job models.Job) error {
	// Save to database, call InsertJob() full record written, status is pending
	err := db.InsertJob(database, job)
	if err != nil {
		log.Printf("Failed to save job to database: %v", err)
		return err
	}

	// Add to Redis queue, only job ID is added to Redis list
	// LPUSH queue:default <job-id> — job is now in line
	// LPUSH adds the job to the beginning of the list 
	// LPUSH job-A  →  [job-A]
	// LPUSH job-B  →  [job-B] [job-A]
	// LPUSH job-C  →  [job-C] [job-B] [job-A]
	//worker will do RPOP hence job is processed in the order they were added
	//process starvation is eliminated

	err = rdb.LPush(context.Background(), "queue:"+job.Queue, job.ID).Err()
	if err != nil {
		log.Printf("Failed to add job to Redis queue: %v", err)
		return err
	}

	return nil
}