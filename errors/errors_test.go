package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testErrorResponse(t *testing.T, message string, status int) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	ErrorResponse(w, r, status, message)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != status {
		t.Errorf("Expected status code %d, but got %d", status, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	jsonString := fmt.Sprintf(`{"error": "%s"}`, message)
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		message string
		status  int
	}{
		{name: "Valid error response", message: "An error occurred", status: http.StatusInternalServerError},
		{name: "Empty error message", message: "", status: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testErrorResponse(t, tc.message, tc.status)
		})
	}
}

func TestServerErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ServerErrorResponse(w, r, errors.New("An error occured"))
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "The server encountered a problem and could not process your request: An error occured"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestNotFoundResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	NotFoundResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "The requested resource could not be found"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestMethodNotAllowedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	MethodNotAllowedResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "The GET method is not supported for this resource"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestBadRequestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	BadRequestResponse(w, r, errors.New("Bad Request"))
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "Bad Request"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestFailedValidationResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	errors := make(map[string]string)
	errors["field1"] = "cannot be empty"
	errors["field2"] = "should be more then 8 characters"
	FailedValidationResponse(w, r, errors)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, but got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": {"field1": "cannot be empty", "field2": "should be more then 8 characters"}}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestEditConflictResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	EditConflictResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, but got %d", http.StatusConflict, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "Unable to update the record, please try again"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestRateLimitExceededResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	RateLimitExceededResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, but got %d", http.StatusTooManyRequests, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "Rate limit exceeded"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestInvalidCredentialsResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	InvalidCredentialsResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "Invalid authentication credentials"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestInvalidAuthenticationTokenResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	InvalidAuthenticationTokenResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "Invalid or missing authentication token"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestAuthenticationRequiredResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	AuthenticationRequiredResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "You must be authenticated to access this resource"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestInactiveAccountResponseResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	InactiveAccountResponse(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, but got %d", http.StatusForbidden, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(`{"error": "Your user account must be activated to access this resource"}`), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}
