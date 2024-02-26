package middlewares

import (
	"book-rental/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	mock.CreateDBforTest(m)
}

func TestAuthorization(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get-available-books", nil)

	w := httptest.NewRecorder()

	Authorization(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	expectedErrorMessage := "Authorization header is missing\n"
	if w.Body.String() != expectedErrorMessage {
		t.Errorf("Expe:'%s':, got :'%s':", expectedErrorMessage, w.Body.String())
	}
}
