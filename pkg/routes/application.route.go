package routes

import (
	"auth-service/middleware"
	"auth-service/pkg/handler/application"
	"log"
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
		r.Use(middleware.ApplicationAuthServiceAuthorization)

		r.Post("/register", application.ApplicationRegisterHandler)
		r.Post("/login", application.ApplicationLoginHandler)
	})

	return r
}

// simple middleware to log incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
