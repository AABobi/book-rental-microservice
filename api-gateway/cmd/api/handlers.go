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

	url := "http://auth2:91/authenticate"
	err := helpers.ReadJSON(w, r, &user)
	fmt.Println(user)
	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}

	resp, err := sendRequest(user, w, url, "POST", authServicePassword)
	if err != nil {
		helpers.ErrorJson(w, "Problem to send request", 500)
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		http.Error(w, "New request problem", 440)
	}

	err = helpers.WriteJSON(w, response, resp.StatusCode)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func DockerTest(w http.ResponseWriter, r *http.Request) {
	var response data.ResponseMessage

	url := "http://localhost:91/test"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		http.Error(w, "Get request problem", 440)
	}
	hashedPassword, err := utils.HashPassword(authServicePassword)

	if err != nil {
		helpers.ErrorJson(w, "Password hashing problem", 500)

	}
	req.Header.Set("key", hashedPassword)
	client := &http.Client{}
	resp, err := client.Do(req)
	err = json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println("TEST")
	fmt.Println(response)
	helpers.WriteJSON(w, response, 200)
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

	err = helpers.WriteJSON(w, response, resp.StatusCode)
	if err != nil {
		helpers.ErrorJson(w, "Problem with writing json", 500)
		return
	}
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

	err = helpers.WriteJSON(w, books, resp.StatusCode)
	if err != nil {
		helpers.ErrorJson(w, "Problem with writing json", 500)
		return
	}
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

	err = helpers.WriteJSON(w, books, resp.StatusCode)
	if err != nil {
		helpers.ErrorJson(w, "Problem with writing json", 500)
		return
	}
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

	err = helpers.WriteJSON(w, books, resp.StatusCode)
	if err != nil {
		helpers.ErrorJson(w, "Problem with writing json", 500)
		return
	}
	defer resp.Body.Close()
}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response data.ResponseMessage

	userId := contextUserId(r)

	err := helpers.ReadJSON(w, r, &books)
	if err != nil {
		helpers.ErrorJson(w, "Problem with reading json", 500)
		return
	}
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

	err = helpers.WriteJSON(w, response, resp.StatusCode)
	if err != nil {
		helpers.ErrorJson(w, "Problem with writing json", 500)
		return
	}
	defer resp.Body.Close()
}

func RentBook(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response data.ResponseMessage

	userId := contextUserId(r)

	err := helpers.ReadJSON(w, r, &books)
	if err != nil {
		helpers.ErrorJson(w, "Problem with reading json", 500)
		return
	}
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

	err = helpers.WriteJSON(w, response, resp.StatusCode)
	if err != nil {
		return
	}
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

			if err != nil {
				helpers.ErrorJson(w, "Problem with new request", 500)
				return nil, err
			}
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

	if err != nil {
		helpers.ErrorJson(w, "Password hashing problem", 500)
		return nil, err
	}
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
