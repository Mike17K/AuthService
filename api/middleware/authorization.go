package middleware

import (
	"auth-service/api/constants"
	"auth-service/api/utils"
	"auth-service/internal/database"
	"auth-service/internal/models"
	"context"
	"log"
	"net/http"
	"os"

	"fmt"
)

// middleware to check if the request has the correct AuthServiceAuthorization header with the secret key
func ServiceToServicePrivateRouteAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := r.Header.Get(constants.ServiceToServiceAuthorizationHeader)
		if secretKey == "" {
			http.Error(w, fmt.Sprintf("%s header is required for this type of request", constants.ServiceToServiceAuthorizationHeader), http.StatusUnauthorized)
			return
		}

		// Load the secret key from the environment variable
		expectedSecretKey := utils.GenerateSecret(os.Getenv(constants.SERVICE_SECRET))
		log.Println("Expected secret key: ", expectedSecretKey)
		if expectedSecretKey == "" {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if the secret key in the Auth-Service-Authorization header matches the expected secret key
		if secretKey != expectedSecretKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// middleware for validating application secret and getting the application object
func ValidateRequestApplicationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var application models.Application
		if err := database.DB.Where("base_secret_key = ?", r.Header.Get(constants.ApplicationAuthorizationHeader)).First(&application).Error; err != nil {
			http.Error(w, "Application not found", http.StatusNotFound)
			return
		}

		// Store application in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.ApplicationContextKey, application)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// middleware for JWT verification
func VerifyJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var application models.Application
		if err := database.DB.Where("base_secret_key = ?", r.Header.Get(constants.ApplicationAuthorizationHeader)).First(&application).Error; err != nil {
			http.Error(w, "Application not found", http.StatusNotFound)
			return
		}

		// Verify the JWT token
		token, err := utils.VerifyJWT(application.BaseSecretKey, r.Header.Get(constants.Authorization))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store the user in the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.UserContextKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
