package application

import (
	"auth-service/api/utils"
	"auth-service/internal/database"
	"auth-service/internal/models"

	"github.com/go-playground/validator/v10"

	"encoding/json"
	"log"
	"net/http"
)

type ApplicationLoginBody struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ApplicationLoginResponse struct {
	BaseSecretKey string `json:"base_secret"`
}

// ApplicationLogin godoc
// @Summary      Login an application
// @Description  Login an application
// @Tags         Application
// @Accept       json
// @Produce      json
// @Param        Auth-Service-Authorization header string true "Authorization key"
// @Param        application body ApplicationLoginBody true "Application login body"
// @Success      200  {object}  utils.SuccessResponse[ApplicationLoginResponse]
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /application/login [post]
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
		utils.SendResponse(w, utils.Error("Application not found", nil), http.StatusNotFound)
		tx.Rollback()
		return
	}
	if application.Password != utils.EncryptPassword(body.Password) {
		utils.SendResponse(w, utils.Error("Invalid password", nil), http.StatusUnauthorized)
		tx.Rollback()
		return
	}
	tx.Commit()

	log.Printf("Application login: %v", application)

	// Send the response
	utils.SendResponse(w, utils.Success("Application login successfully", ApplicationLoginResponse{
		BaseSecretKey: application.BaseSecretKey,
	}), http.StatusOK)
}
