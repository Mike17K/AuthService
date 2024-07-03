package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"auth-service/api/constants"
	"auth-service/internal/models"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	// Establish initial database connection
	if err := connectToDatabase(); err != nil {
		return nil, err
	}

	// Start goroutine for periodic health checks
	go periodicHealthCheck()

	return DB, nil
}

// connectToDatabase establishes the initial database connection
func connectToDatabase() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv(constants.DB_USER),
		os.Getenv(constants.DB_PASSWORD),
		os.Getenv(constants.DB_HOST),
		os.Getenv(constants.DB_PORT),
		os.Getenv(constants.DB_NAME))
	fmt.Println(dsn)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		// call reconnect() here after a delay
		delay := 2 * time.Second

		go func() {
			<-time.After(delay) // Wait for the specified duration
			reconnect()
		}()

		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Set maximum open connections
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(time.Hour)

	// AutoMigrate creates tables based on the model structs
	if !DB.HasTable(&models.User{}) {
		DB.AutoMigrate(&models.User{})
	}

	if !DB.HasTable(&models.Application{}) {
		DB.AutoMigrate(&models.Application{})
	}

	fmt.Println("Connected to the database")

	return nil
}

// periodicHealthCheck periodically checks the health of the database connection
func periodicHealthCheck() {
	ticker := time.NewTicker(1 * time.Hour) // Adjust the interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := DB.DB().Ping(); err != nil {
				fmt.Println("Error checking database connection health:", err)
				// Attempt to reconnect
				if err := reconnect(); err != nil {
					fmt.Printf("Failed to reconnect to database: %v\n", err)
				}
			}
		}
	}
}

// reconnect attempts to reconnect to the database
func reconnect() error {
	fmt.Println("Attempting to reconnect to the database...")
	// Close the existing connection
	if err := DB.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %v", err)
	}
	// Reopen the connection
	return connectToDatabase()
}
