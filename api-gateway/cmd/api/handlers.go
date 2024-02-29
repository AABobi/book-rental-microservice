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

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type Book struct {
	BookID          uint       `json:"book_id"`
	Name            string     `json:"name"`
	DateToReturn    *time.Time `json:"date_to_return"`
	UserID          uint       `json:"user_id"`
	ShippingAddress string     `json:"shipping_address"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN")
	var user data.User
	var response AuthResponse

	url := "http://localhost:91/authenticate"
	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}

	resp, err := sendRequest(user, w, url, "POST", authServicePassword)
	err = json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println("TEST")
	fmt.Println(response)
	//body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "New request problem", 440)
	}
	helpers.WriteJSON(w, response, resp.StatusCode)

	defer resp.Body.Close()
}

func SingUp(w http.ResponseWriter, r *http.Request) {
	/*fmt.Println("SIGN")
	var user data.User
	var response ResponseMessage

	url := "http://localhost:91/add-new-user"

	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}

	resp, err := sendRequest(user, w, url, "POST")

	if err != nil {
		fmt.Println("Send request problem")
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println(response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "New request problem", 440)
		return
	}
	fmt.Println("Response Status:", response)

	helpers.WriteJSON(w, response, resp.StatusCode)*/
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SIGN")
	var books []Book
	var response ResponseMessage

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
			fmt.Println(err)
			http.Error(w, "New request problem", 440)
		}
		return
	}
	fmt.Println("Response Status:", books)

	helpers.WriteJSON(w, books, resp.StatusCode)
}

func GetRentedBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("2")
	userId, _ := r.Context().Value("myKey").(uint)
	fmt.Println("2")
	var books []Book
	var response ResponseMessage

	strUserId := strconv.Itoa(int(userId))
	url := fmt.Sprintf("http://localhost:80/get-rented-books/%s", strUserId)

	resp, err := sendRequest(books, w, url, "GET", bookServicePassword)
	fmt.Println("RESP")
	if err != nil {
		fmt.Println("Send request problem")
		http.Error(w, "Send request problem", 440)
		return
	}
	fmt.Println("NewDecoder")
	err = json.NewDecoder(resp.Body).Decode(&books)
	if err != nil {
		err = json.NewDecoder(resp.Body).Decode(&response)
		fmt.Println("problem")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "New request problem", 440)
		}
		return
	}
	fmt.Println("Response Status:", books)

	helpers.WriteJSON(w, books, resp.StatusCode)
}

func GetAvailableBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetAvailableBooks")
	var books []Book
	var response ResponseMessage

	url := "http://localhost:80/get-available-books"

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
			fmt.Println(err)
			http.Error(w, "New request problem", 440)
		}
		return
	}
	fmt.Println("Response Status:", books)

	helpers.WriteJSON(w, books, resp.StatusCode)
}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response ResponseMessage
	userId, _ := r.Context().Value("myKey").(uint)
	strUserId := strconv.Itoa(int(userId))
	helpers.ReadJSON(w, r, &books)
	url := fmt.Sprintf("http://localhost:80/return-books/%s", strUserId)

	fmt.Println("RETURN")
	resp, err := sendRequest(books, w, url, "POST", bookServicePassword)

	if err != nil {
		fmt.Println("Send request problem")
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	fmt.Println("Response Status:", response)

	helpers.WriteJSON(w, response, resp.StatusCode)
}

func RentBook(w http.ResponseWriter, r *http.Request) {
	var books []Book
	var response ResponseMessage
	userId, _ := r.Context().Value("myKey").(uint)
	strUserId := strconv.Itoa(int(userId))

	helpers.ReadJSON(w, r, &books)
	url := fmt.Sprintf("http://localhost:80/rent/%s", strUserId)

	resp, err := sendRequest(books, w, url, "POST", bookServicePassword)

	if err != nil {
		fmt.Println("Send request problem")
		http.Error(w, "Send request problem", 440)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	fmt.Println("Response Status:", response)

	helpers.WriteJSON(w, response, resp.StatusCode)
}
func sendRequest(data any, w http.ResponseWriter, url string, method string, password string) (*http.Response, error) {
	fmt.Println("SENDREQ")
	fmt.Println(data)
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
			fmt.Println("GET")
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
