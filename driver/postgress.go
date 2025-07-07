package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("Waiting for the database to be ready...")
	time.Sleep(5 * time.Second) // Wait for 5 seconds before connecting

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Database connection established successfully")
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Failed to close the database connection: %v", err)
	}
}
