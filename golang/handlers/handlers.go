package handlers

import (
	"book-rental/data"
	"book-rental/db"
	"book-rental/helpers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

var response = ResponseMessage{
	Message: "EMPTY MESSAGE",
}

func GetAvailableBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GETAVAIL")
	condition := "user_id IS 0"
	books := data.FindBooks(db.DB, condition)
	_ = helpers.WriteJSON(w, http.StatusOK, books)
}

func GetRentedBooks(w http.ResponseWriter, r *http.Request) {
	var books []data.Book

	//userId, ok := userIdFromContext("myKey", r)
	id := chi.URLParam(r, "id")
	uintValue, err := strconv.ParseUint(id, 10, 64)
	userId := uint(uintValue)
	/*if !ok {
		http.Error(w, "Failed to retrieve user ID from context", http.StatusInternalServerError)
		return
	}*/

	data.GetRentedBooksFromDb(db.DB, &userId, &books)

	fmt.Println(books)
	err = helpers.WriteJSON(w, 200, books)

	if err != nil {
		http.Error(w, "Cannot rent a book", http.StatusInternalServerError)
		return
	}
}

func RentBook(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	helpers.ReadJSON(w, r, &books)
	//userId, ok := userIdFromContext("myKey", r)
	fmt.Println("RENT BOOK")
	fmt.Println(books)
	id := chi.URLParam(r, "id")
	uintValue, err := strconv.ParseUint(id, 10, 64)
	userId := uint(uintValue)
	/*	if !ok {
		http.Error(w, "Failed to retrieve user ID from context", http.StatusInternalServerError)
		return
	}*/
	err = data.RentBook(db.DB, &userId, books)

	if err != nil {
		response.Message = "Failed to rent a book"
		_ = helpers.WriteJSON(w, 440, response)
		return
	}

	response.Message = "The rental has been made correctly"
	_ = helpers.WriteJSON(w, 200, response)

	if err != nil {
		response.Message = "Failed to write response"
		_ = helpers.WriteJSON(w, 200, response)
		return
	}
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	data.FindAllBooks(db.DB, &books)
	fmt.Println(len(books))
	helpers.WriteJSON(w, 200, books)
}

func ReturnTheBook(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	helpers.ReadJSON(w, r, &books)
	fmt.Println("ReturnTheBook")
	/*userId, ok := userIdFromContext("myKey", r)

	if !ok {
		http.Error(w, "Failed to retrieve user ID from context", http.StatusInternalServerError)
		return
	}*/
	id := chi.URLParam(r, "id")
	uintValue, err := strconv.ParseUint(id, 10, 64)
	userId := uint(uintValue)
	err = data.RemoveUserIdFromBooks(db.DB, userId, &books)

	if err != nil {
		response.Message = "The book cannot be updated"
		_ = helpers.WriteJSON(w, 440, response)
		return
	}

	w.WriteHeader(http.StatusOK)

	response.Message = "Success: Book is updated correctly"
	_ = helpers.WriteJSON(w, 200, response)

	if err != nil {
		response.Message = "Failed to write response"
		_ = helpers.WriteJSON(w, 440, response)
		return
	}
}

var userIdFromContext = func(key string, r *http.Request) (*uint, bool) {
	userId, ok := r.Context().Value("myKey").(*uint)
	return userId, ok
}
