package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/krutinewalkar/distributed-task-queue/db"
	"github.com/krutinewalkar/distributed-task-queue/queue"
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

	// Test Enqueue
	job, err := queue.Enqueue(database, rdb, "default", "Hello, World!")
	if err != nil {
		log.Fatal("Failed to enqueue job:", err)
	}
	log.Printf("Job enqueued: %v", job)

	// Start worker
	worker.StartWorker(database, rdb, "default")
}
