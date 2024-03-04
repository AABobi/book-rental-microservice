package main

import (
	"auth-db/db"
	"auth-db/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("TO")
	CreateDBforTest(m)
}
func CreateDBforTest(m *testing.M) {
	var err error

	rootDir := filepath.Join(".", "..", "..")
	if err != nil {
		panic("failed to get current working directory")
	}

	dbName := filepath.Join(rootDir, "auth-db-test.db")
	os.Remove(dbName)

	db.DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	pass, _ := utils.HashPassword("password")

	users := []db.User{}

	user := db.User{
		Email:    "test@gmail.com",
		Password: pass,
	}

	users = append(users, user)
	db.DB.AutoMigrate(&db.User{})
	db.DB.Create(&users)

	os.Exit(m.Run())
}

func TestCreateNewUser(t *testing.T) {
	//pass, _ := utils.HashPassword("password")
	newUser := db.User{
		Email:    "test1@gmail.com",
		Password: "pass",
	}
	jsonString, _ := json.Marshal(newUser)

	reader := bytes.NewReader(jsonString)

	req := httptest.NewRequest(http.MethodPost, "/signup", reader)

	record := httptest.NewRecorder()

	CreateUser(record, req)

	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	expected := db.ResponseMessage{Message: "User added correctly"}
	b := string(body)

	var gotResponse db.ResponseMessage
	json.Unmarshal([]byte(b), &gotResponse)

	if expected.Message != gotResponse.Message {
		t.Errorf("Cannot create a user \n %v", string(body))
	}
}

func TestCreateNewUser_userExistInDB(t *testing.T) {
	newUser := db.User{
		Email:    "test1@gmail.com",
		Password: "pass",
	}
	jsonString, _ := json.Marshal(newUser)

	reader := bytes.NewReader(jsonString)

	req := httptest.NewRequest(http.MethodPost, "/signup", reader)

	record := httptest.NewRecorder()

	CreateUser(record, req)

	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	expected := db.ResponseMessage{Message: "User exist"}
	b := string(body)

	var gotResponse db.ResponseMessage
	json.Unmarshal([]byte(b), &gotResponse)

	if expected.Message != gotResponse.Message {
		t.Errorf("User should exist")
	}
}

func TestLogin(t *testing.T) {
	newUser := db.User{
		Email:    "test1@gmail.com",
		Password: "password",
	}
	jsonString, _ := json.Marshal(newUser)

	reader := bytes.NewReader(jsonString)

	req := httptest.NewRequest(http.MethodPost, "/login", reader)

	record := httptest.NewRecorder()

	FindUser(record, req)

	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	expected := db.AuthResponse{Message: "Authenticated"}
	b := string(body)

	var gotResponse db.ResponseMessage
	err = json.Unmarshal([]byte(b), &gotResponse)

	if err != nil || expected.Message != gotResponse.Message {
		t.Errorf("Login test failed")
	}
}

func TestAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/Auth", nil)
	user := db.User{Email: "test1@gmail.com"}
	user.GetUser(db.DB)

	record := httptest.NewRecorder()
	token, _ := utils.GenerateToken(user.Email, &user.UserID)
	req.Header.Set("Authorization", token)

	Auth(record, req)

	resp := record.Result()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	expected := db.User{UserID: 2}
	b := string(body)
	var gotResponse db.User
	err = json.Unmarshal([]byte(b), &gotResponse)

	if err != nil || expected.UserID != gotResponse.UserID {
		t.Errorf("Login test failed")
	}
}

func TestGetAllUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/Auth", nil)
	record := httptest.NewRecorder()
	GetAllUsers(record, req)
	resp := record.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	expected := 2
	var gotResponse []db.User
	b := string(body)
	err = json.Unmarshal([]byte(b), &gotResponse)

	if err != nil || len(gotResponse) != expected {
		t.Errorf("Login test failed")
	}
}
