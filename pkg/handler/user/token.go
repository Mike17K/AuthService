// handler/handler.go

package user

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"encoding/json"
	"log"
	"net/http"
)

func UserGetAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var application = models.Application{}
	if err := database.DB.Where("base_secret_key = ?", r.Header.Get("Application-Secret")).First(&application).Error; err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}
	var token, err = utils.VerifyJWT(application.BaseSecretKey, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if token["token_type"] != "refresh_token" {
		http.Error(w, "Wrong Token Type", http.StatusUnauthorized)
		return
	}
	r.Body.Close()
	// Validations - End

	newAccessToken := utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
		"user_id":    token["user_id"],
		"user_type":  token["user_type"],
		"token_type": "access_token",
	}, 30*60*time.Second)

	log.Printf("New access token generated for user: %v", token["user_id"])

	response := map[string]string{
		"message":      "New access token generated successfully",
		"access_token": newAccessToken,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
