package main

import (
	"log"
	"net/http"

	"GOLANG/controller"
	"GOLANG/db"
	"GOLANG/Middleware"

	"github.com/gorilla/mux"
)
func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome, authorized user!"))
}
func main() {
	http.HandleFunc("/protected", middleware.AuthMiddleware(ProtectedRoute))
	db.Connect()
	userController := &controller.UserController{
		DB: db.DB, // use the DB instance from db package
	}

	router := mux.NewRouter()
	router.HandleFunc("/users", userController.GetUsers).Methods("GET")
	router.HandleFunc("/users", userController.CreateUser).Methods("POST")
	router.HandleFunc("/users/login", userController.Login).Methods("POST")
	router.HandleFunc("/users", userController.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/delete-all", userController.DeleteAllUsers).Methods("DELETE")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
