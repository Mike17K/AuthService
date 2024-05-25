package main

import (
	"auth-service/internal/database"
	"auth-service/pkg/constants"
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
		// os.Exit(1)
	}

	// print all env variables
	for _, element := range os.Environ() {
		fmt.Println(element)
	}

	// Set a secret key for the service if not already set in .env file
	if os.Getenv(constants.SERVICE_SECRET) == "" {
		fmt.Println(constants.SERVICE_SECRET, "not set in .env file, setting a random secret key")
		os.Setenv(constants.SERVICE_SECRET, utils.GenerateRandomString(32))
		fmt.Println("New secret key: ", os.Getenv(constants.SERVICE_SECRET))
		os.Exit(1)
	}

	fmt.Println("Time Generated Access Token: ", utils.GenerateSecret(os.Getenv(constants.SERVICE_SECRET)))

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
	port := os.Getenv(constants.PORT)
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Println("Serve on http://localhost:" + port)
	//err = http.ListenAndServeTLS(":"+port, os.Getenv("SSL_CERT_PATH"), os.Getenv("SSL_KEY_PATH"), r)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
