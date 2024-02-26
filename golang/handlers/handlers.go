package handlers

import (
	"book-rental/data"
	"book-rental/db"
	"book-rental/helpers"
	"book-rental/utils"
	"fmt"

	"net/http"
)

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func GetAvailableBooks(w http.ResponseWriter, r *http.Request) {
	condition := "user_id IS 0"
	books := data.FindBooks(db.DB, condition)
	_ = helpers.WriteJSON(w, http.StatusOK, books)
}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var user data.User

	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Failed to read json", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		http.Error(w, "Hash error", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	err = data.CreateNewUser(db.DB, &user)

	if err != nil {
		http.Error(w, "User exist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	message := "Success: User created"

	_, err = fmt.Fprintln(w, message)

	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user data.User

	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Failed to read json", http.StatusInternalServerError)
		return
	}

	passwordCorrect := user.CheckCredentials(db.DB)

	if !passwordCorrect {
		http.Error(w, "Wrong credentials", http.StatusBadRequest)
		return
	}

	userId := user.FindUser(db.DB)

	token, err := utils.GenerateToken(user.Email, &userId)

	if err != nil {
		http.Error(w, "Generate token problem", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Email: user.Email,
		Token: token,
	}

	helpers.WriteJSON(w, 200, response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := data.GetAllUsers(db.DB)
	helpers.WriteJSON(w, 200, users)
}

func GetRentedBooks(w http.ResponseWriter, r *http.Request) {
	var books []data.Book

	userId, ok := userIdFromContext("myKey", r)

	if !ok {
		http.Error(w, "Failed to retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	data.GetRentedBooksFromDb(db.DB, userId, &books)

	err := helpers.WriteJSON(w, 200, books)

	if err != nil {
		http.Error(w, "Cannot rent a book", http.StatusInternalServerError)
		return
	}
}

func RentBook(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	helpers.ReadJSON(w, r, &books)
	userId, ok := userIdFromContext("myKey", r)
	if !ok {
		http.Error(w, "Failed to retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	err := data.RentBook(db.DB, userId, books)

	if err != nil {
		http.Error(w, "Failed to rent a book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = fmt.Fprint(w, "The rental has been made correctly")

	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	data.FindAllBooks(db.DB, &books)
	helpers.WriteJSON(w, 200, books)
}

func ReturnTheBook(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	helpers.ReadJSON(w, r, &books)

	userId, ok := userIdFromContext("myKey", r)

	if !ok {
		http.Error(w, "Failed to retrieve user ID from context", http.StatusInternalServerError)
		return
	}
	err := data.RemoveUserIdFromBooks(db.DB, *userId, &books)

	if err != nil {
		http.Error(w, "The book cannot be updated", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	message := "Success: Book is updated correctly"

	_, err = fmt.Fprintln(w, message)

	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

var userIdFromContext = func(key string, r *http.Request) (*uint, bool) {
	userId, ok := r.Context().Value("myKey").(*uint)
	return userId, ok
}
