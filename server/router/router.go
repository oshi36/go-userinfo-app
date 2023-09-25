package router

import (
	"htmlgo/server/middleware"

	"github.com/gorilla/mux"
)

// Captial means exporting the method
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/new", middleware.NewUser).Methods("POST")
	router.HandleFunc("/", middleware.GetUsers).Methods("POST", "GET")
	router.HandleFunc("/api/users/{id}", middleware.User).Methods("GET", "DELETE")

	return router
}
