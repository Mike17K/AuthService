// handler/handler.go

package user

import (
	"auth-service/api/constants"
	"auth-service/api/utils"
	"auth-service/internal/database"
	"auth-service/internal/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"

	"encoding/json"
	"log"
	"net/http"
)

type UserLoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var body UserLoginBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get data from context
	application, ok := r.Context().Value(constants.ApplicationContextKey).(models.Application)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Validations - End

	// Create the User record in the database
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var User models.User
	if err := tx.Where("Email = ? AND Password = ?", body.Email, utils.EncryptPassword(body.Password)).First(&User).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		tx.Rollback()
		return
	}
	refreshToken := utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
		constants.JWTUserIdField:    User.ID,
		constants.JWTUserTypeField:  models.SimpleUser,
		constants.JWTTokenTypeField: constants.RefreshToken,
	}, 30*24*60*60*time.Second)
	User.RefreshToken = refreshToken
	if err := tx.Save(&User).Error; err != nil {
		log.Printf("Error updating User: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	tx.Commit()

	log.Printf("User login: %v", User)

	accessToken := utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
		constants.JWTUserIdField:    User.ID,
		constants.JWTUserTypeField:  User.UserType,
		constants.JWTTokenTypeField: constants.AccessToken,
	}, 30*60*time.Second)

	response := UserLoginResponse{
		Message:      "User login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
