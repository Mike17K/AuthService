package utils

import (
	"auth-service/api/constants"

	"encoding/json"
	"net/http"
	"strings"
)

func GetTokenFromHeader(r *http.Request) string {
	const bearerPrefix = "Bearer "
	authHeader := r.Header.Get(constants.Authorization)
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return strings.TrimPrefix(authHeader, bearerPrefix)
	}
	return authHeader
}

func DecodeBody(r *http.Request, body interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(body); err != nil {
		return err
	}
	return nil
}

func SendResponse(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func Error(message string, err error) ErrorResponse {
	res := ErrorResponse{
		Success: false,
		Message: message,
	}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}

type SuccessResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func Success[T any](message string, data T) SuccessResponse[T] {
	return SuccessResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}
