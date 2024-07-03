package user

import (
	"auth-service/api/constants"
	"auth-service/api/utils"
	"auth-service/internal/database"
	"auth-service/internal/models"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"

	"encoding/json"
	"log"
	"net/http"
)

type UserRegisterBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var body UserRegisterBody
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
	application, ok := r.Context().Value(constants.ApplicationContextKey).(models.Application)
	if !ok {
		http.Error(w, "Application not found in context", http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	// Validations - End

	// Create the User object
	ID := utils.GenerateRandomString()
	Refreshtoken := utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
		constants.JWTUserIdField:    ID,
		constants.JWTUserTypeField:  models.SimpleUser,
		constants.JWTTokenTypeField: constants.RefreshToken,
	}, 30*24*60*60*time.Second)
	var User = models.User{
		ID:                              ID,
		Name:                            body.Name,
		Password:                        utils.EncryptPassword(body.Password),
		Email:                           body.Email,
		EmailVerified:                   false,
		AuthProvider:                    models.Email,
		TwoFactor:                       false,
		TwoFactorValidated:              false,
		TwoFactorInitialSecret:          utils.GenerateRandomString(),
		IsBlocked:                       false,
		OpenToChangePasswordWithTokenID: "",
		ApplicationID:                   application.ID,
		UserType:                        models.SimpleUser,
		RefreshToken:                    Refreshtoken,
	}

	// Save the User object to the database
	tx := database.DB.Begin() // create a transaction for rollback in case of error
	if err := tx.Create(&User).Error; err != nil {
		// Check for duplicate entry error
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "User with the same email already exists", http.StatusConflict)
		} else {
			// Log other types of errors
			log.Printf("Error creating User: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		tx.Rollback()
		return
	}
	tx.Commit()
	log.Printf("User created: %v", User)
	// Database save successful

	// Send the response
	response := UserRegisterResponse{
		Message: "User created successfully",
		AccessToken: utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
			constants.JWTUserIdField:    User.ID,
			constants.JWTUserTypeField:  User.UserType,
			constants.JWTTokenTypeField: constants.AccessToken,
		}, 30*60*time.Second),
		RefreshToken: Refreshtoken,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
