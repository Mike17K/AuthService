// handler/handler.go

package handler

import (
	"encoding/json"
	"net/http"
)

// HelloHandler handles requests to the /hello endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Hello, World!"}
	json.NewEncoder(w).Encode(response)
}
