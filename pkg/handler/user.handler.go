// handler/handler.go

package handler

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/utils"
	"auth-service/pkg/validations"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"

	"encoding/json"
	"log"
	"net/http"
)

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var body validations.UserRegisterBody
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

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var body validations.UserLoginBody
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

	// Create the User record in the database
	tx := database.DB.Begin() // create a transaction for rollback in case of error
	var User = models.User{}
	if err := tx.Where("Email = ? AND Password = ?", body.Email, utils.EncryptPassword(body.Password)).First(&User).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		tx.Rollback()
		return
	}
	Refreshtoken := utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
		"user_id":    User.ID,
		"user_type":  models.SimpleUser,
		"token_type": "refresh_token",
	}, 30*24*60*60*time.Second)
	User.RefreshToken = Refreshtoken
	if err := tx.Save(&User).Error; err != nil {
		log.Printf("Error updating User: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	tx.Commit()

	log.Printf("User login: %v", User)

	response := map[string]string{
		"message": "User login successfully",
		"access_token": utils.GenerateJWT(application.BaseSecretKey, jwt.MapClaims{
			"user_id":    User.ID,
			"user_type":  User.UserType,
			"token_type": "access_token",
		}, 30*60*time.Second),
		"refresh_token": User.RefreshToken,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var application = models.Application{}
	if err := database.DB.Where("base_secret_key = ?", r.Header.Get("Application-Secret")).First(&application).Error; err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}
	var token, err = utils.VerifyJWT(application.BaseSecretKey, r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("Error verifying JWT: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if token["token_type"] != "access_token" {
		http.Error(w, "Wrong Token Type", http.StatusUnauthorized)
		return
	}
	r.Body.Close()
	// Validations - End

	// Delete the refresh_token from the User record in the database
	tx := database.DB.Begin() // create a transaction for rollback in case of error
	var User = models.User{}
	if err := tx.Where("id = ?", token["user_id"]).First(&User).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		tx.Rollback()
		return
	}
	User.RefreshToken = ""
	if err := tx.Save(&User).Error; err != nil {
		log.Printf("Error updating User: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	tx.Commit()

	log.Printf("User logout: %v", User)

	response := map[string]string{
		"message": "User logout successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

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
