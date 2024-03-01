package main

import (
	"api-gateway/data"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var authToken string
var testEmail string

func TestSingUp(t *testing.T) {
	testEmail = randomEmail(10)
	user := data.User{Email: testEmail, Password: "testpass"}
	jsonString, _ := json.Marshal(user)

	r := bytes.NewReader(jsonString)
	req, err := http.NewRequest("POST", "/add-new-user", r)

	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a mock HTTP server
	/*mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the response from the authentication service
		responseJSON := `{"message": "User added correctly"}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))

	defer mockServer.Close()*/

	recorder := httptest.NewRecorder()
	SingUp(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	var response data.ResponseMessage
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error decoding JSON: %v", err)
	}

	// Check the expected values in the response
	expectedResponse := data.ResponseMessage{Message: "User added correctly"}
	if response.Message != expectedResponse.Message {
		t.Errorf("Expected response %+v, but got %+v", expectedResponse, response)
	}
}

func TestLogin(t *testing.T) {
	// Create a sample user for the request body

	user := data.User{Email: testEmail, Password: "testpass"}
	jsonString, _ := json.Marshal(user)

	r := bytes.NewReader(jsonString)

	req, err := http.NewRequest("POST", "/login", r)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authResponseJSON := `{"token": "some_token", "message": "Authentication successful"}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(authResponseJSON))
	}))

	defer mockServer.Close()
	recorder := httptest.NewRecorder()
	Login(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	var response data.AuthResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error decoding JSON: %v", err)
	}

	// Check the expected values in the response
	expectedResponse := data.AuthResponse{Token: "some_token", Message: "Authenticated"}
	authToken = response.Token
	if response.Message != expectedResponse.Message {
		t.Errorf("Expected response %+v, but got %+v", expectedResponse, response)
	}
}

func TestGetAvailableBooks(t *testing.T) {
	req := httptest.NewRequest("GET", "/get-available-books", nil)
	rec := httptest.NewRecorder()

	req.Header.Set("key", "superSecretBooksPassword")
	GetAvailableBooks(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected response 200, but got %+v", rec.Code)
	}
}

func TestGetAllBooks(t *testing.T) {
	req := httptest.NewRequest("GET", "/get", nil)
	rec := httptest.NewRecorder()

	req.Header.Set("key", "superSecretBooksPassword")
	GetAllBooks(rec, req)

	var response []data.Book
	json.Unmarshal(rec.Body.Bytes(), &response)
	fmt.Println(len(response))
	if len(response) == 0 {
		t.Errorf("Slice is empty")
	}
}

func TestRentBook(t *testing.T) {
	_contextUserId := contextUserId

	defer func() {
		contextUserId = _contextUserId
	}()

	contextUserId = func(*http.Request) string {
		return "2"
	}

	jsonString := `[{
        "book_id": 1,
        "name": "The Alchemist",
        "date_to_return": null,
        "user_id": 0,
        "shipping_address": ""
    }]`
	req, _ := http.NewRequest("POST", "/rent-book", strings.NewReader(jsonString))
	req.Header.Set("key", "superSecretBooksPassword")

	rec := httptest.NewRecorder()

	RentBook(rec, req)

	var response data.ResponseMessage
	err := json.Unmarshal(rec.Body.Bytes(), &response)

	if rec.Code != 200 || err != nil {
		t.Errorf("Expected response 200, but got %+v", rec.Code)
	}

	expected := data.ResponseMessage{Message: "The rental has been made correctly"}

	if expected.Message != response.Message {
		t.Errorf("Expected %s got %s", expected.Message, response.Message)
	}
}

func TestGetRentedBooks(t *testing.T) {
	_contextUserId := contextUserId

	defer func() {
		contextUserId = _contextUserId
	}()

	contextUserId = func(*http.Request) string {
		return "2"
	}

	req := httptest.NewRequest("GET", "/get-rented-books", nil)
	rec := httptest.NewRecorder()

	GetRentedBooks(rec, req)

	var response []data.Book
	json.Unmarshal(rec.Body.Bytes(), &response)
	fmt.Println(len(response))
	fmt.Println("HALO")
	if len(response) == 0 {
		t.Errorf("Incorrect slice lenght")
	}

}

func TestReturnBook(t *testing.T) {
	_contextUserId := contextUserId

	defer func() {
		contextUserId = _contextUserId
	}()

	contextUserId = func(*http.Request) string {
		return "2"
	}

	jsonString := `[{
        "book_id": 1,
        "name": "The Alchemist",
        "date_to_return": null,
        "user_id": 0,
        "shipping_address": ""
    }]`
	req, _ := http.NewRequest("POST", "/return-books", strings.NewReader(jsonString))
	req.Header.Set("key", "superSecretBooksPassword")

	rec := httptest.NewRecorder()
	ReturnBook(rec, req)
	if rec.Code != 200 {
		t.Errorf("Expected response 200, but got %+v", rec.Code)
	}

	var response data.ResponseMessage
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error decoding JSON: %v", err)
	}

	// Check the expected values in the response
	expectedResponse := data.ResponseMessage{Message: "Success: Book is updated correctly"}

	if response.Message != expectedResponse.Message {
		t.Errorf("Expected response %+v, but got %+v", expectedResponse, response)
	}
}

func randomEmail(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
