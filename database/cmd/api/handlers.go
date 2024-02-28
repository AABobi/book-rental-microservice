package main

import (
	"auth-db/db"
	"auth-db/helpers"
	"fmt"
	"net/http"
)

func FindUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FIND USER11")
	var requestData db.User
	err := helpers.ReadJSON(w, r, &requestData)
	if err != nil {
		fmt.Println("ERROR")
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("FIND USER")
	fmt.Println(requestData)
	requestData.GetUser(db.DB)
	if requestData.UserID != 0 {
		fmt.Println("SENDING")
		helpers.WriteJSON(w, 200, requestData)
		return
	}
	http.Error(w, "Cannot find a user", http.StatusBadRequest)
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser db.User
	findUser := newUser
	err := helpers.ReadJSON(w, r, &newUser)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	findUser.GetUser(db.DB)
	if findUser.UserID != 0 {
		http.Error(w, "User exist", http.StatusBadRequest)
		return
	}

	newUser.AddUser(db.DB)

	newUser.GetUser(db.DB)

	if newUser.UserID != 0 {
		helpers.WriteJSON(w, 200, "User added correctly")
		return
	}

	http.Error(w, "Cannot add new user", http.StatusBadRequest)
}

// Handler only for tests
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []db.User
	db.GetAllUsers(db.DB, &users)
	helpers.WriteJSON(w, 200, users)
}
