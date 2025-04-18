package main

import (
	"log"
	"net/http"

	"GOLANG/controller"
	"GOLANG/db"

	"github.com/gorilla/mux"
)

func main() {
	db.Connect()

	userController := &controller.UserController{
		DB: db.DB, // use the DB instance from db package
	}

	router := mux.NewRouter()
	router.HandleFunc("/users", userController.GetUsers).Methods("GET")
	router.HandleFunc("/users", userController.CreateUser).Methods("POST")
	router.HandleFunc("/users", userController.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/delete-all", userController.DeleteAllUsers).Methods("DELETE")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
