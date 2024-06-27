package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/windevkay/flhoutils/assert"
)

func TestGenerateUniqueId(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "Generate ID of length 20", length: 20},
		{name: "Generate ID of length 10", length: 10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateUniqueId(tc.length)
			assert.Equal(t, len(result), tc.length)
		})
	}
}

func TestReadIDParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	tests := []struct {
		name string
		args string
		want int64
		err  error
	}{
		{name: "Valid param", args: "1", want: 1, err: nil},
		{name: "Invalid param", args: "0", want: 0, err: errors.New("invalid ID parameter")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params := httprouter.Params{{Key: "id", Value: tc.args}}
			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			id, err := ReadIDParam(req)

			if tc.err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			} else {
				assert.Equal(t, id, tc.want)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name    string
		message string
		status  int
		data    Envelope
		headers http.Header
	}{
		{name: "200 status with custom headers", message: "success", status: http.StatusOK, data: Envelope{"data": "success"}, headers: http.Header{
			"X-Custom-Header": []string{"value1"},
		}},
		{name: "500 status with no custom headers", message: "error", status: http.StatusInternalServerError, data: Envelope{"data": "error"}, headers: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSON(w, tc.status, tc.data, tc.headers)
			resp := w.Result()
			checkCustomHeader(t, tc.message, resp)
			checkGeneralHeader(t, resp)
			checkStatusCode(t, tc.status, resp)
			checkResponseBody(t, tc.message, resp)
		})
	}
}

func checkCustomHeader(t *testing.T, message string, resp *http.Response) {
	if message == "success" {
		headerValue := resp.Header.Get("X-Custom-Header")
		expectedHeaderValue := "value1"
		if headerValue != expectedHeaderValue {
			t.Errorf("Expected header 'X-Custom-Header' to have value '%s', but got '%s'", expectedHeaderValue, headerValue)
		}
	}
}

func checkGeneralHeader(t *testing.T, resp *http.Response) {
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		t.Errorf("Expected header 'Content-Type' to have value '%s', but got '%s'", expectedContentType, contentType)
	}
}

func checkStatusCode(t *testing.T, expectedStatus int, resp *http.Response) {
	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status code %d, but got %d", expectedStatus, resp.StatusCode)
	}
}

func checkResponseBody(t *testing.T, message string, resp *http.Response) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	jsonString := fmt.Sprintf(`{"data": "%s"}`, message)
	var expectedResponse map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected response: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, actualResponse)
	}
}

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "Valid request body", err: nil},
		{name: "Unknown key in request body", err: errors.New(`body contains unknown key "oddKey"`)},
		{name: "Empty request body", err: errors.New("body must not be empty")},
		{name: "Badly formed JSON", err: errors.New("body contains badly-formed JSON")},
		{name: "Large request body", err: errors.New("body must not be larger than 1048576 bytes")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			checkValidRequestBody(t, tc.name, tc.err)
			checkUnknownKeyInRequestBody(t, tc.name, tc.err)
			checkEmptyRequestBody(t, tc.name, tc.err)
			checkBadlyFormedJSON(t, tc.name, tc.err)
			checkLargeRequestBody(t, tc.name, tc.err)
		})
	}
}

func checkLargeRequestBody(t *testing.T, testName string, expectedErr error) {
	if testName == "Large request body" {
		w := httptest.NewRecorder()

		jsonSnippet := `{"data":"value"},`
		repeatCount := (1_048_576 / len(jsonSnippet)) + 1
		largeJSON := "[" + strings.Repeat(jsonSnippet, repeatCount)
		largeJSON = strings.TrimRight(largeJSON, ",") + "]"

		reader := bytes.NewReader([]byte(largeJSON))
		r := httptest.NewRequest("POST", "/", reader)
		var dst struct {
			Data string `json:"data"`
		}

		err := ReadJSON(w, r, &dst)
		assert.Equal(t, err.Error(), expectedErr.Error())
	}
}

func checkEmptyRequestBody(t *testing.T, testName string, expectedErr error) {
	if testName == "Empty request body" {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		var dst struct {
			Data string `json:"data"`
		}

		err := ReadJSON(w, r, &dst)
		assert.Equal(t, err.Error(), expectedErr.Error())
	}
}

func checkValidRequestBody(t *testing.T, testName string, expectedErr error) {
	if testName == "Valid request body" {
		w := httptest.NewRecorder()
		requestBodyContent := map[string]interface{}{
			"data": "value",
		}
		requestBodyBytes, err := json.Marshal(requestBodyContent)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		requestBodyReader := bytes.NewReader(requestBodyBytes)
		r := httptest.NewRequest("POST", "/", requestBodyReader)
		var dst struct {
			Data string `json:"data"`
		}

		err = ReadJSON(w, r, &dst)
		assert.Equal(t, err, expectedErr)
	}
}

func checkBadlyFormedJSON(t *testing.T, testName string, expectedErr error) {
	if testName == "Badly formed JSON" {
		w := httptest.NewRecorder()
		badJSON := `{"data": "some value"`
		r := httptest.NewRequest("POST", "/", strings.NewReader(badJSON))
		var dst struct {
			Data string `json:"data"`
		}

		err := ReadJSON(w, r, &dst)
		assert.Equal(t, err.Error(), expectedErr.Error())
	}
}

func checkUnknownKeyInRequestBody(t *testing.T, testName string, expectedErr error) {
	if testName == "Unknown key in request body" {
		w := httptest.NewRecorder()
		requestBodyContent := map[string]interface{}{
			"oddKey": "oddValue",
		}
		requestBodyBytes, err := json.Marshal(requestBodyContent)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		requestBodyReader := bytes.NewReader(requestBodyBytes)
		r := httptest.NewRequest("POST", "/", requestBodyReader)
		var dst struct {
			Data string `json:"data"`
		}

		err = ReadJSON(w, r, &dst)
		assert.Equal(t, err.Error(), expectedErr.Error())
	}
}
