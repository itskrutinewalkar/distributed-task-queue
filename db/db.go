package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open db: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to db:", err)
	}
	log.Println("Postgres connected")
	return db
}
