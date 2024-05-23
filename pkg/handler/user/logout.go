// handler/handler.go

package user

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
)

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
