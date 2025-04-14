package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to default postgres database
	adminDBURL := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := sql.Open("postgres", adminDBURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()

	// Check if db already exist
	dbName := os.Getenv("DB_NAME")

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking database existence: %v\n", err)
		return
	}

	if exists {
		fmt.Printf("Database %s already exists\n", dbName)
		return
	}

	// Create the database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		return
	}

	fmt.Printf("Database %s created successfully\n", dbName)
}
