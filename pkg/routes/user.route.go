package routes

import (
	"auth-service/pkg/handler/user"
	"net/http"

	"github.com/go-chi/chi"
)

func UserRouter() http.Handler {
	r := chi.NewRouter()

	// Handle GET method for "/user/hello" route
	// r.Get("/hello", handler.HelloHandler)

	r.Group(func(r chi.Router) {
		r.Post("/register", user.UserRegisterHandler)
		r.Post("/login", user.UserLoginHandler)
		r.Post("/logout", user.UserLogoutHandler)
		r.Get("/token", user.UserGetAccessTokenHandler)
	})

	return r
}
