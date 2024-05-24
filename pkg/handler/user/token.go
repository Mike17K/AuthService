// handler/handler.go

package user

import (
	"auth-service/internal/models"
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"encoding/json"
	"log"
	"net/http"
)

type UserTokenResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func UserGetAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	token, ok := r.Context().Value(constants.UserContextKey).(jwt.MapClaims)
	if !ok || token[constants.JWTTokenTypeField] != constants.RefreshToken {
		http.Error(w, "Wrong Token Type", http.StatusUnauthorized)
		return
	}
	application, ok := r.Context().Value(constants.ApplicationContextKey).(models.Application)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	// Validations - End

	// Generate new access token
	newAccessToken := utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
		constants.JWTUserIdField:    token[constants.JWTUserIdField],
		constants.JWTUserTypeField:  token[constants.JWTUserTypeField],
		constants.JWTTokenTypeField: constants.AccessToken,
	}, 30*60*time.Second)
	log.Printf("New access token generated for user: %v", token["user_id"])

	// Generate the Response
	response := UserTokenResponse{
		Message:     "New access token generated successfully",
		AccessToken: newAccessToken,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
