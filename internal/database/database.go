// database.go

package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/joho/godotenv"
)

// Define a struct representing a database model
type User struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);"`
	Age  uint
}

// DB is a global variable to hold the database connection
var db *gorm.DB

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	var err error

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the database")

	defer db.Close()

	db.DB().SetMaxIdleConns(10)

	// Set maximum open connections
	db.DB().SetMaxOpenConns(100)

	// Set the maximum amount of time a connection may be reused
	db.DB().SetConnMaxLifetime(time.Hour)

	// AutoMigrate creates tables based on the model structs
	if !db.HasTable(&User{}) {
		// AutoMigrate creates tables based on the model structs
		db.AutoMigrate(&User{})
	}

	return db, nil
}
