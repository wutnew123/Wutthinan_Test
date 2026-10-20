package main

import (
	"log"
	"net/http"

	"back-end/config"
	"back-end/routes"

	"github.com/gorilla/handlers"
)

func main() {
	config.ConnectDB()

	r := routes.SetupRoutes()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:4200",
		}),
		handlers.AllowedMethods([]string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		}),
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"Accept",
			"Origin",
		}),
	)

	log.Println("Server is running on :8080")

	log.Fatal(
		http.ListenAndServe(
			":8080",
			corsHandler(r),
		),
	)
}
