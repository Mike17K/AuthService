package routes

import (
	"auth-service/api/handler/application"
	"auth-service/api/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

/**
 * The Routes for the application have prefix "/application/*"
 * Are being used to comunicate with another server in order to share keys
 */
func ApplicationRouter() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.ServiceToServicePrivateRouteAuthorization)

		r.Post("/register", application.ApplicationRegisterHandler)
		r.Post("/login", application.ApplicationLoginHandler)
	})

	return r
}
