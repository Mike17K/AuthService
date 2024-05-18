package middleware

import (
	"auth-service/pkg/utils"
	"net/http"
	"os"

	"fmt"
)

// middleware to check if the request has the correct Authorization header with the secret key
func ApplicationAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := r.Header.Get("Authorization")
		if secretKey == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Load the secret key from the environment variable
		fmt.Println("Expecting: ", utils.GenerateSecret(os.Getenv("SERVICE_SECRET_KEY")))
		expectedSecretKey := utils.GenerateSecret(os.Getenv("SERVICE_SECRET_KEY"))
		if expectedSecretKey == "" {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if the secret key in the Authorization header matches the expected secret key
		if secretKey != expectedSecretKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
