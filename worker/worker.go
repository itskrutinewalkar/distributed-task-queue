package worker

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/krutinewalkar/distributed-task-queue/db"
	"github.com/krutinewalkar/distributed-task-queue/models"
)

func StartWorker(database *sql.DB, rdb *redis.Client, queueName string) {
	log.Printf("Worker started for queue: %s", queueName)

	for {
		timeout := 5 * time.Second
		res, err := rdb.BRPop(context.Background(), timeout, "queue:"+queueName).Result()

		if err != nil {
			if err == redis.Nil {
				log.Printf("No jobs in queue: %s, waiting...", queueName)
				continue
			}
			log.Printf("Error fetching job from Redis: %v", err)
			continue
		}

		if len(res) < 2 {
			log.Printf("Invalid response from Redis: %v", res)
			continue
		}

		jobID := res[1]
		log.Printf("Fetched job ID: %s from queue: %s", jobID, queueName)

		job, err := db.GetJob(database, jobID)
		if err != nil {
			log.Printf("Error fetching job from database: %v", err)
			continue
		}

		// Mark as running
		err = db.UpdateJobStatus(database, jobID, models.StatusRunning, "")
		if err != nil {
			log.Printf("Failed to mark job %s as running: %v", job.ID, err)
			continue
		}

		// Process the job
		processErr := processJob(job)

		if processErr != nil {
			log.Printf("Error processing job %s: %v", job.ID, processErr)
			job.Attempts++

			if job.Attempts >= job.MaxRetries {
				// Out of retries — mark as failed
				job.Status = models.StatusFailed
				log.Printf("Job %s failed after %d attempts", job.ID, job.Attempts)
			} else {
				// Re-enqueue at the back of the queue
				job.Status = models.StatusPending
				log.Printf("Retrying job %s, attempt %d", job.ID, job.Attempts)
				err = rdb.RPush(context.Background(), "queue:"+job.Queue, job.ID).Err()
				if err != nil {
					log.Printf("Failed to re-enqueue job %s: %v", job.ID, err)
				}
			}

			err = db.UpdateJobStatus(database, jobID, job.Status, processErr.Error())
			if err != nil {
				log.Printf("Failed to update job %s in database: %v", job.ID, err)
			}
			continue
		}

		// Mark as completed
		err = db.UpdateJobStatus(database, jobID, models.StatusCompleted, "")
		if err != nil {
			log.Printf("Failed to update job %s in database: %v", job.ID, err)
		}
		log.Printf("Job %s completed successfully", job.ID)
	}
}

func processJob(job models.Job) error {
	time.Sleep(2 * time.Second)
	log.Printf("Job %s processed with payload: %s", job.ID, job.Payload)
	return nil
}
