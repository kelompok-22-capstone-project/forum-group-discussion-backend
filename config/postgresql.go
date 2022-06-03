package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func NewPostgreSQLDatabase() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("string conversion failed: %w", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var psqlInfo string
	if os.Getenv("ENV") == "production" {
		psqlInfo = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s",
			host, port, user, password, dbName,
		)
	} else {
		psqlInfo = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbName,
		)
	}

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		defer func(db *sql.DB) {
			if err := db.Close(); err != nil {
				log.Fatalf("failed to close the connection: %s\n", err.Error())
			}
		}(db)
		return nil, fmt.Errorf("can't sent ping to database: %w", err)
	}

	return db, nil
}
