// handler/handler.go

package user

import (
	"auth-service/api/constants"
	"auth-service/internal/database"
	"auth-service/internal/models"
	"encoding/json"
	"log"
	"net/http"

	_ "auth-service/api/utils"

	"github.com/golang-jwt/jwt/v4"
)

type UserLogoutResponse struct {
	Message string `json:"message"`
}

// UserLogout godoc
// @Summary      Logout an user
// @Description  Logout an user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Application-Secret header string true "Application base secret key"
// @Param        Authorization header string true "Bearer access token"
// @Success      200  {object}  utils.SuccessResponse[UserLogoutResponse]
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /user/logout [post]
func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Validations - Start
	// Get data from context
	token, ok := r.Context().Value(constants.UserContextKey).(jwt.MapClaims)
	if !ok || token[constants.JWTTokenTypeField] != constants.AccessToken {
		http.Error(w, "Wrong Token Type", http.StatusUnauthorized)
		return
	}
	r.Body.Close()
	// Validations - End

	// Update the user's refresh token in the database
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var user = models.User{}
	if err := tx.Where("id = ?", token[constants.JWTUserIdField]).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		tx.Rollback()
		return
	}
	user.RefreshToken = ""
	if err := tx.Save(&user).Error; err != nil {
		log.Printf("Error updating User: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	tx.Commit()
	log.Printf("User logout: %v", user)
	// Database update successful

	// Send the response
	response := UserLogoutResponse{
		Message: "User logged out successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
