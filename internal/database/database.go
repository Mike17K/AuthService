// database.go

package database

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// DB is a global variable to hold the database connection
var DB *sqlx.DB

// InitDB initializes the database connection
func InitDB() error {
	var err error

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}

	// Get database connection string from environment variables
	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	// Connect to the database
	DB, err = sqlx.Open("mysql", dbConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Set maximum idle connections
	DB.SetMaxIdleConns(10)

	// Set maximum open connections
	DB.SetMaxOpenConns(100)

	// Set the maximum amount of time a connection may be reused
	DB.SetConnMaxLifetime(time.Hour)

	// Ping the database to check if the connection is successful
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping the database: %v", err)
	}

	fmt.Println("Connected to the database")
	return nil
}
