package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/krutinewalkar/distributed-task-queue/api"
	"github.com/krutinewalkar/distributed-task-queue/db"
	"github.com/krutinewalkar/distributed-task-queue/worker"
)

func main() {
	// Connect to Postgres
	database := db.Connect("postgres://taskqueue:password@localhost:5432/taskqueue?sslmode=disable")
	defer database.Close()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Redis connected")

	// Start worker in background
	go worker.StartWorker(database, rdb, "1234567890")

	// Setup routes
	router := gin.Default()
	router.POST("/jobs", func(c *gin.Context) {
		api.HandleEnqueue(c, database, rdb)
	})
	router.GET("/jobs/:id", func(c *gin.Context) {
		api.HandleGetJob(c, database)
	})

	log.Println("API running on :8082")
	router.Run(":8082")
}
