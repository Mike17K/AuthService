package main

import (
	"auth-service/internal/database"
	"auth-service/pkg/routes"
	"auth-service/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
		os.Exit(1)
	}

	fmt.Println("Access key: ", os.Getenv("SERVICE_SECRET_KEY"))
	fmt.Println("Access: ", utils.GenerateSecret(os.Getenv("SERVICE_SECRET_KEY")))

	database.InitDB()

	r := chi.NewRouter()

	// Apply middleware to the entire router
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(10)) // 10 seconds max request timeout

	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "deny"))
	r.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))
	// r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=5184000; includeSubDomains")) // Uncomment this line if you are using HTTPS
	r.Use(httprate.LimitByRealIP(100, time.Minute))

	// Register your handlers
	r.Mount("/application", routes.ApplicationRouter())
	r.Mount("/user", routes.UserRouter())

	// Start the server on port from .env file
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Println("Serving on https://localhost:" + port)
	err = http.ListenAndServeTLS(":"+port, os.Getenv("SSL_CERT_PATH"), os.Getenv("SSL_KEY_PATH"), r)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
