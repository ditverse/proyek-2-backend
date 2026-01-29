package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(connString string) *sql.DB {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}

	// Configure connection pool untuk production
	db.SetMaxOpenConns(25)                 // Max concurrent connections
	db.SetMaxIdleConns(5)                  // Min idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping db: %v", err)
	}

	log.Println("Koneksi Database Berhasil")
	return db
}
