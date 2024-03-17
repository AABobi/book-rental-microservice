package handlers

import (
	"book-rental/data"
	"book-rental/mock"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func TestMain(m *testing.M) {
	mock.CreateDBforTest(m)
}

func TestGetAvailableBooks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get-available-books", nil)

	record := httptest.NewRecorder()

	GetAvailableBooks(record, req)

	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	b := string(body)
	var books []data.Book
	err = json.Unmarshal([]byte(b), &books)
	if err != nil {
		t.Errorf("Problem with unmarshal")
	}

	if len(books) != 37 {
		t.Errorf("Unexpeted body return")
	}
}

func TestRentBook(t *testing.T) {
	_pathId := pathId

	defer func() {
		pathId = _pathId
	}()

	pathId = func(*http.Request) string {
		return "1"
	}

	book := []data.Book{
		{Name: "The Catcher in the Rye", ShippingAddress: "test"},
	}

	jsonString, _ := json.Marshal(book)

	reader := bytes.NewReader(jsonString)

	req := httptest.NewRequest(http.MethodPost, "/rent", reader)

	record := httptest.NewRecorder()

	RentBook(record, req)
	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var expected = data.ResponseMessage{Message: `The rental has been made correctly`}
	var gotResponse data.ResponseMessage
	err = json.Unmarshal(body, &gotResponse)

	if gotResponse.Message != expected.Message && err == nil {
		t.Errorf("Rent a book test failed")
	}
}

func TestGetRentedBooks(t *testing.T) {
	_pathId := pathId

	defer func() {
		pathId = _pathId
	}()

	pathId = func(*http.Request) string {
		return "1"
	}

	req := httptest.NewRequest(http.MethodPost, "/get-rented-books", nil)

	record := httptest.NewRecorder()

	GetRentedBooks(record, req)
	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	expectedBook := []data.Book{
		{Name: "The Catcher in the Rye", ShippingAddress: "test"},
	}

	b := string(body)
	var books []data.Book
	err = json.Unmarshal([]byte(b), &books)
	if err != nil {
		t.Errorf("Unmarshal problem")
		return
	}

	if books[0].Name != expectedBook[0].Name {
		t.Errorf("Get rented books test failed")
	}
}

func TestReturnTheBook(t *testing.T) {
	_pathId := pathId

	defer func() {
		pathId = _pathId
	}()

	pathId = func(*http.Request) string {
		return "1"
	}

	book := []data.Book{
		{Name: "The Catcher in the Rye", ShippingAddress: "test"},
	}

	jsonString, _ := json.Marshal(book)

	reader := bytes.NewReader(jsonString)

	req := httptest.NewRequest(http.MethodPost, "/return-books", reader)

	record := httptest.NewRecorder()

	ReturnTheBook(record, req)
	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var expected = data.ResponseMessage{Message: `Success: Book is updated correctly`}
	var gotResponse data.ResponseMessage
	err = json.Unmarshal(body, &gotResponse)

	if gotResponse.Message != expected.Message && err == nil {
		t.Errorf("Rent a book test failed")
	}
}
