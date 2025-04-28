package main

import (
	"log"
	"net/http"

	"GOLANG/controller"
	"GOLANG/db"
	"GOLANG/middleware" // Fix package name to small "middleware"

	"github.com/gorilla/mux"
)

func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome, authorized user!"))
}

func main() {
	db.Connect()

	userController := &controller.UserController{
		DB: db.DB,
	}

	router := mux.NewRouter()

	// Public routes (no AuthMiddleware)
	router.HandleFunc("/users/login", userController.Login).Methods("POST")
	router.HandleFunc("/users/logout", userController.Logout).Methods("POST")
	router.HandleFunc("/users", userController.CreateUser).Methods("POST")

	// Protected routes (with AuthMiddleware)
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUsers))).Methods("GET")
	router.Handle("/user", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUsersbyid))).Methods("GET")
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(userController.DeleteUser))).Methods("DELETE")
	router.Handle("/users/delete-all", middleware.AuthMiddleware(http.HandlerFunc(userController.DeleteAllUsers))).Methods("DELETE")
	router.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(ProtectedRoute))).Methods("GET")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
