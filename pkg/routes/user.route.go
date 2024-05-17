package routes

import (
	"auth-service/pkg/handler"
	"net/http"

	"github.com/go-chi/chi"
)

func UserRouter() http.Handler {
	r := chi.NewRouter()

	// Handle GET method for "/user/hello" route
	r.Get("/hello", handler.HelloHandler)

	return r
}
