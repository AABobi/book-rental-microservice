package main

import (
	"api-gateway/data"
	"api-gateway/helpers"
	"api-gateway/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const authServicePassword = "superSecretAuthDbPassword"
const bookServicePassword = "superSecretBooksPassword"

type Book struct {
	BookID          uint       `json:"book_id"`
	Name            string     `json:"name"`
	DateToReturn    *time.Time `json:"date_to_return"`
	UserID          uint       `json:"user_id"`
	ShippingAddress string     `json:"shipping_address"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user data.User
	var response data.AuthResponse

	url := "http://localhost:91/authenticate"
	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}

	resp, err := sendRequest(user, w, url, "POST", authServicePassword)
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		http.Error(w, "New request problem", 440)
	}

	helpers.WriteJSON(w, response, resp.StatusCode)
	defer resp.Body.Close()
}

func SingUp(w http.ResponseWriter, r *http.Request) {
	var user data.User
	var response data.ResponseMessage

	url := "http://localhost:91/add-new-user"
	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}

	resp, err := sendRequest(user, w, url, "POST", authServicePassword)

	if err != nil {
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		http.Error(w, "New request problem", 440)
		return
	}

	helpers.WriteJSON(w, response, resp.StatusCode)
	defer resp.Body.Close()
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response data.ResponseMessage

	url := "http://localhost:80/get-all-books"

	resp, err := sendRequest(books, w, url, "GET", bookServicePassword)

	if err != nil {
		fmt.Println("Send request problem")
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&books)
	if err != nil {
		err = json.NewDecoder(resp.Body).Decode(&response)

		if err != nil {
			http.Error(w, "New request problem", 440)
		}
		return
	}

	helpers.WriteJSON(w, books, resp.StatusCode)
	defer resp.Body.Close()
}

func GetRentedBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response data.ResponseMessage

	userId := contextUserId(r)

	url := fmt.Sprintf("http://localhost:80/get-rented-books/%s", userId)
	resp, err := sendRequest(books, w, url, "GET", bookServicePassword)

	if err != nil {
		fmt.Println("Send request problem")
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&books)
	if err != nil {
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			http.Error(w, "New request problem", 440)
		}
		return
	}

	helpers.WriteJSON(w, books, resp.StatusCode)
	defer resp.Body.Close()
}

func GetAvailableBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response data.ResponseMessage

	url := "http://localhost:80/get-available-books"

	resp, err := sendRequest(books, w, url, "GET", bookServicePassword)

	if err != nil {
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&books)
	if err != nil {
		err = json.NewDecoder(resp.Body).Decode(&response)

		if err != nil {
			http.Error(w, "New request problem", 440)
		}
		return
	}

	helpers.WriteJSON(w, books, resp.StatusCode)
	defer resp.Body.Close()
}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response data.ResponseMessage

	userId := contextUserId(r)

	helpers.ReadJSON(w, r, &books)
	url := fmt.Sprintf("http://localhost:80/return-books/%s", userId)

	resp, err := sendRequest(books, w, url, "POST", bookServicePassword)

	if err != nil {
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return
	}

	helpers.WriteJSON(w, response, resp.StatusCode)
	defer resp.Body.Close()
}

func RentBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("R")
	var books []Book
	var response data.ResponseMessage

	userId := contextUserId(r)

	helpers.ReadJSON(w, r, &books)
	url := fmt.Sprintf("http://localhost:80/rent/%s", userId)

	resp, err := sendRequest(books, w, url, "POST", bookServicePassword)

	if err != nil {
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}

	helpers.WriteJSON(w, response, resp.StatusCode)
	defer resp.Body.Close()
}
func sendRequest(data any, w http.ResponseWriter, url string, method string, password string) (*http.Response, error) {
	var req *http.Request
	var err error

	switch method {
	case "POST":
		{
			jsonData, err := json.Marshal(data)
			if err != nil {
				http.Error(w, "Post request problem", 440)
				return nil, err
			}
			req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
			break
		}
	case "GET":
		{
			req, err = http.NewRequest(method, url, nil)
			if err != nil {
				http.Error(w, "Get request problem", 440)
				return nil, err
			}
			break
		}
	default:
		break
	}

	hashedPassword, err := utils.HashPassword(password)
	req.Header.Set("key", hashedPassword)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

var contextUserId = func(r *http.Request) string {
	userId, _ := r.Context().Value("myKey").(uint)
	strUserId := strconv.Itoa(int(userId))
	return strUserId
}
