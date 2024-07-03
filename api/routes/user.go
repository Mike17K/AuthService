package routes

import (
	"auth-service/api/handler/user"
	"auth-service/api/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func UserRouter() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.ValidateRequestApplicationMiddleware)
		r.Post("/register", user.UserRegisterHandler)
		r.Post("/login", user.UserLoginHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyJWTMiddleware)
			r.Post("/logout", user.UserLogoutHandler)
			r.Get("/token", user.UserGetAccessTokenHandler)
		})
	})

	return r
}
