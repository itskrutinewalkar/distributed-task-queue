package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/krutinewalkar/distributed-task-queue/db"
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

	log.Println("All systems go")
}
