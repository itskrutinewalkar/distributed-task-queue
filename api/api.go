package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/krutinewalkar/distributed-task-queue/db"
	"github.com/krutinewalkar/distributed-task-queue/queue"
)

type QueryBody struct {
	Queue   string `json:"queue"`
	Payload string `json:"payload"`
}

func HandleEnqueue(c *gin.Context, database *sql.DB, rdb *redis.Client) {
	var body QueryBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	job, err := queue.Enqueue(database, rdb, body.Queue, body.Payload)
	if err != nil {
		log.Printf("Failed to enqueue job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job_id": job.ID, "status": job.Status})

}

func HandleGetJob(c *gin.Context, database *sql.DB) {
	jobID := c.Param("id")
	job, err := db.GetJob(database, jobID)
	if err != nil {
		log.Printf("Failed to get job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job"})
		return
	}

	c.JSON(http.StatusOK, job)
}
