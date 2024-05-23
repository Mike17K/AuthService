package user

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/utils"
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
	var application = models.Application{}
	if err := database.DB.Where("base_secret_key = ?", r.Header.Get("Application-Secret")).First(&application).Error; err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}
	r.Body.Close()
	// Validations - End

	ID := utils.GenerateRandomString()
	Refreshtoken := utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
		"user_id":    ID,
		"user_type":  models.SimpleUser,
		"token_type": "refresh_token",
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

	// Create the User record in the database
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

	response := map[string]string{
		"message": "User created successfully",
		"access_token": utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
			"user_id":    User.ID,
			"user_type":  User.UserType,
			"token_type": "access_token",
		}, 30*60*time.Second),
		"refresh_token": Refreshtoken,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
