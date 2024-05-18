// handler/handler.go

package handler

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/utils"
	"auth-service/pkg/validations"
	"strings"

	"github.com/go-playground/validator/v10"

	"encoding/json"
	"log"
	"net/http"
)

// HelloHandler handles requests to the /hello endpoint
func ApplicationRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var body validations.ApplicationRegisterBody
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

	var application = models.Application{
		ID:            utils.GenerateRandomString(),
		Name:          body.Name,
		Password:      utils.EncryptPassword(body.Password),
		Description:   body.Description,
		IsBlocked:     false,
		BaseSecretKey: utils.GenerateSecret(body.Name),
	}

	// Create the Application record in the database
	if err := database.DB.Create(&application).Error; err != nil {
		// Check for duplicate entry error
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "Application with the same name already exists", http.StatusConflict)
		} else {
			// Log other types of errors
			log.Printf("Error creating application: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Application created: %v", application)

	response := map[string]string{
		"message":         "Application created successfully",
		"base_secret_key": application.BaseSecretKey,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
