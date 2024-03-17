package middlewares

import (
	"api-gateway/data"
	"api-gateway/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorization(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()

	hashPassword, _ := utils.HashPassword("superSecretAuthDbPassword")
	req.Header.Set("key", hashPassword)

	Authorization(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rec, req)

	var response data.AuthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		return
	}

	expected := "Incorrect token"

	if response.Message != expected {
		t.Errorf("Validation doesn't work correctly")
	}
}
