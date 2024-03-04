package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	recorder := httptest.NewRecorder()

	data := map[string]interface{}{
		"key": "value",
	}

	err := WriteJSON(recorder, http.StatusOK, data, http.Header{"Custom-Header": []string{"header-value"}})
	if err != nil {
		t.Fatalf("Error writing JSON: %v", err)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', but got '%s'", contentType)
	}

	customHeaderValue := recorder.Header().Get("Custom-Header")
	if customHeaderValue != "header-value" {
		t.Errorf("Expected Custom-Header to be 'header-value', but got '%s'", customHeaderValue)
	}

	statusCode := recorder.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, statusCode)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	if !reflect.DeepEqual(responseBody, data) {
		t.Errorf("Expected response body %+v, but got %+v", data, responseBody)
	}
}

func TestReadJSON(t *testing.T) {
	jsonPayload := `{"key": "value"}`
	req, err := http.NewRequest("POST", "/example", strings.NewReader(jsonPayload))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	recorder := httptest.NewRecorder()

	var data map[string]string
	err = ReadJSON(recorder, req, &data)
	if err != nil {
		t.Fatalf("Error reading JSON: %v", err)
	}

	expectedData := map[string]string{"key": "value"}
	if !jsonEqual(data, expectedData) {
		t.Errorf("Expected data %+v, but got %+v", expectedData, data)
	}

	remainingBody := recorder.Body.String()
	if remainingBody != "" {
		t.Errorf("Expected no remaining body, but got: %s", remainingBody)
	}
}

// jsonEqual checks if two JSON-encoded strings are equal
func jsonEqual(a, b map[string]string) bool {
	aJSON, errA := json.Marshal(a)
	bJSON, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	return string(aJSON) == string(bJSON)
}
