package main

import (
	"fmt"
	"net/http"

	_ "auth-service/docs"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	r := chi.NewRouter()

	fmt.Println("Swagger UI running on http://localhost:1323/swagger/index.html")
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"), //The url pointing to API definition
	))

	http.ListenAndServe(":1323", r)
}
