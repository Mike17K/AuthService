package application

import (
	"auth-service/api/utils"
	"auth-service/internal/database"
	"auth-service/internal/models"
	"strings"

	"github.com/go-playground/validator/v10"

	"encoding/json"
	"log"
	"net/http"
)

type ApplicationRegisterBody struct {
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type ApplicationRegisterResponse struct {
	BaseSecretKey string `json:"base_secret"`
}

// ApplicationRegister godoc
// @Summary      Register an application
// @Description  Register an application
// @Tags         Application
// @Accept       json
// @Produce      json
// @Param        Auth-Service-Authorization header string true "Authorization key"
// @Param        application body ApplicationRegisterBody true "Application register body"
// @Success      200  {object}  utils.SuccessResponse[ApplicationRegisterResponse]
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /application/register [post]
func ApplicationRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	var body ApplicationRegisterBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		utils.SendResponse(w, utils.Error(err.Error(), nil), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		utils.SendResponse(w, utils.Error(err.Error(), nil), http.StatusBadRequest)
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
	tx := database.DB.Begin() // create a transaction for rollback in case of error
	if err := tx.Create(&application).Error; err != nil {
		// Check for duplicate entry error
		if strings.Contains(err.Error(), "Duplicate entry") {
			utils.SendResponse(w, utils.Error("Application with the same name already exists", nil), http.StatusConflict)
		} else {
			// Log other types of errors
			log.Printf("Error creating application: %v", err)
			utils.SendResponse(w, utils.Error("Internal server error", nil), http.StatusInternalServerError)
		}
		tx.Rollback()
		return
	}
	tx.Commit()

	log.Printf("Application created: %v", application)

	utils.SendResponse(w, utils.Success("Application created", ApplicationRegisterResponse{BaseSecretKey: application.BaseSecretKey}), http.StatusOK)
}
