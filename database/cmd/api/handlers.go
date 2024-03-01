package main

import (
	"auth-db/db"
	"auth-db/helpers"
	"auth-db/utils"
	"fmt"
	"net/http"
)

var response = db.ResponseMessage{
	Message: "EMPTY MESSAGE",
}

func Auth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TEST")
	token := r.Header.Get("Authorization")
	id, err := utils.VerifyToken(token)

	if err != nil {
		fmt.Println("ERROR")
		authResponse := db.AuthResponse{Message: "Incorrect token", Token: token}
		helpers.WriteJSON(w, 440, authResponse)
		return
	}

	user := db.User{
		UserID: *id,
	}

	helpers.WriteJSON(w, 200, user)
}

func FindUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TEST")
	var requestData db.User

	err := helpers.ReadJSON(w, r, &requestData)
	if err != nil {
		response.Message = "Error decoding JSON"
		helpers.WriteJSON(w, 400, response)
		return
	}

	requestData.GetUser(db.DB)
	if requestData.UserID != 0 {
		token, err := utils.GenerateToken(requestData.Email, &requestData.UserID)

		if err != nil {
			response.Message = "Cannot generate token"
			helpers.WriteJSON(w, 400, response)
			return
		}

		responseAuth := db.AuthResponse{
			Message: "Authenticated",
			Token:   token,
		}

		helpers.WriteJSON(w, 200, responseAuth)
		return
	}
	response.Message = "Cannot find a user"
	helpers.WriteJSON(w, 400, response)
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser db.User
	err := helpers.ReadJSON(w, r, &newUser)

	if err != nil {
		response.Message = "Error decoding JSON"
		helpers.WriteJSON(w, 400, response)
		return
	}
	findUser := newUser
	findUser.GetUser(db.DB)

	if findUser.UserID != 0 {
		response.Message = "User exist"
		helpers.WriteJSON(w, 400, response)
		return
	}

	newUser.AddUser(db.DB)
	newUser.GetUser(db.DB)

	if newUser.UserID != 0 {
		response.Message = "User added correctly"
		helpers.WriteJSON(w, 200, response)
		return
	}

	response.Message = "Cannot add new user"
	helpers.WriteJSON(w, 400, response)
}

// Handler only for tests
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []db.User
	db.GetAllUsers(db.DB, &users)
	helpers.WriteJSON(w, 200, users)
}
