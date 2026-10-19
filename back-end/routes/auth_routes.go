package routes

import (
	"back-end/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/auth/register", controllers.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", controllers.Login).Methods("POST")

	r.HandleFunc("/api/auth/logout", controllers.Logout).Methods("POST")

	return r
}
