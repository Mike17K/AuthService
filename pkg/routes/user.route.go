package routes

import (
	"auth-service/middleware"
	"auth-service/pkg/handler"
	"net/http"

	"github.com/go-chi/chi"
)

func UserRouter() http.Handler {
	r := chi.NewRouter()

	// Handle GET method for "/user/hello" route
	// r.Get("/hello", handler.HelloHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.ApplicationAuthServiceAuthorization)

		r.Post("/register", handler.UserRegisterHandler)
		r.Post("/login", handler.UserLoginHandler)
		r.Post("/logout", handler.UserLogoutHandler)
		r.Get("/token", handler.UserGetAccessTokenHandler)
	})

	return r
}
