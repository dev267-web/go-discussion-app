// db/postgres.go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// InitPostgres opens a connection to PostgreSQL using environment variables.
// It returns a *sql.DB that’s ready for queries.
func InitPostgres(ctx context.Context) (*sql.DB, error) {
	// 1) Read environment variables (provide sensible defaults if missing)
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	// Note: if DB_PASSWORD is empty, make sure your Postgres user allows no-password login
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "discussion_app"
	}
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	// 2) Construct the DSN string
	//    e.g.: "host=localhost port=5432 user=postgres password=secret dbname=discussion_app sslmode=disable"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode,
	)

	// 3) Open a database handle
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open Postgres connection: %w", err)
	}

	// 4) Set reasonable connection pool settings
	//    (these numbers can be tuned later as needed)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	// 5) Verify with a context-backed ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("unable to ping Postgres: %w", err)
	}

	return db, nil
}
