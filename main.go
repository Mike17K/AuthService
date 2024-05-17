package main

import (
	"auth-service/middleware"
	"auth-service/pkg/handler"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("failed to load .env file")
		os.Exit(1)
	}

	// Create a new mux router
	mux := http.NewServeMux()

	// Register your handlers
	mux.HandleFunc("/hello", handler.HelloHandler)

	// Apply middleware to the entire router
	http.Handle("/", middleware.LoggingMiddleware(mux))

	// Start the server on port from .env file
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Println("Server started on port " + port)
	http.ListenAndServe(":"+port, nil)
}
