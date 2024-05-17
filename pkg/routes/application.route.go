package routes

import (
	"auth-service/pkg/handler"
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

	//r.Group https://artursiarohau.medium.com/go-chi-rate-limiter-useful-examples-8277dc4d4ff5

	// Handle GET method for "/application/hello" route
	r.With(LoggingMiddleware).Post("/hello", handler.HelloHandler)
	r.With(LoggingMiddleware).Get("/hello", handler.HelloHandler)

	return r
}

// simple middleware to log incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
