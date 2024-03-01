package authenticate

import (
	"auth-db/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/authorization", nil)
	hashPassword, _ := utils.HashPassword("superSecretAuthDbPassword")
	req.Header.Set("key", hashPassword)
	w := httptest.NewRecorder()

	Authenticate(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthenticateUnhappy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/authorization", nil)
	hashPassword, _ := utils.HashPassword("superSecretAuthDbPasswor")
	req.Header.Set("key", hashPassword)
	w := httptest.NewRecorder()

	Authenticate(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
