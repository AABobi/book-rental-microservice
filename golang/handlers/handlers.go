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

var response = data.ResponseMessage{
	Message: "EMPTY MESSAGE",
}

func GetAvailableBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GETEVAIL")
	condition := "user_id IS 0"
	books := data.FindBooks(db.DB, condition)
	_ = helpers.WriteJSON(w, 201, books)
}

func GetRentedBooks(w http.ResponseWriter, r *http.Request) {

	var books []data.Book

	id := pathId(r)
	uintValue, err := strconv.ParseUint(id, 10, 64)
	userId := uint(uintValue)

	data.GetRentedBooksFromDb(db.DB, &userId, &books)

	err = helpers.WriteJSON(w, 201, books)

	if err != nil {
		http.Error(w, "Cannot rent a book", http.StatusInternalServerError)
		return
	}
}

func RentBook(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	helpers.ReadJSON(w, r, &books)
	id := pathId(r)

	fmt.Println("book")
	fmt.Println(books)
	uintValue, err := strconv.ParseUint(id, 10, 64)
	userId := uint(uintValue)

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
	helpers.WriteJSON(w, 200, books)
}

func ReturnTheBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ReturnTheBook RENTED")
	var books []data.Book
	helpers.ReadJSON(w, r, &books)

	id := pathId(r)
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

var pathId = func(r *http.Request) string {
	id := chi.URLParam(r, "id")
	return id
}
