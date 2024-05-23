package application

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/utils"

	"github.com/go-playground/validator/v10"

	"encoding/json"
	"log"
	"net/http"
)

type ApplicationLoginBody struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func ApplicationLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var body ApplicationLoginBody
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
	r.Body.Close()
	// Validations - End

	// Get the Application record in the database
	tx := database.DB.Begin()
	var application models.Application
	if err := tx.Where("name = ?", body.Name).First(&application).Error; err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		tx.Rollback()
		return
	}
	if application.Password != utils.EncryptPassword(body.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		tx.Rollback()
		return
	}
	tx.Commit()

	log.Printf("Application login: %v", application)

	response := map[string]string{
		"message":         "Application login successfully",
		"base_secret_key": application.BaseSecretKey,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
