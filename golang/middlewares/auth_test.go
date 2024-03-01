package middlewares

import (
	"book-rental/mock"
	"book-rental/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	mock.CreateDBforTest(m)
}

func TestAuthorization(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get-available-books", nil)
	hashPassword, _ := utils.HashPassword("superSecretBooksPassword")
	req.Header.Set("key", hashPassword)
	w := httptest.NewRecorder()

	Authorization(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthorizationUnhappy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get-available-books", nil)
	hashPassword, _ := utils.HashPassword("superSecretBooksPasswor")
	req.Header.Set("key", hashPassword)
	w := httptest.NewRecorder()

	Authorization(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
