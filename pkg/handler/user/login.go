// handler/handler.go

package user

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/utils"
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
