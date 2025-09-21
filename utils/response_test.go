package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		wantStatus int
	}{
		{
			name:       "Success response",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Error response",
			statusCode: http.StatusBadRequest,
			data:       map[string]string{"error": "bad request"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Nil data",
			statusCode: http.StatusNoContent,
			data:       nil,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondJSON(w, tt.statusCode, tt.data)

			if w.Code != tt.wantStatus {
				t.Errorf("RespondJSON() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("RespondJSON() Content-Type = %v, want application/json", w.Header().Get("Content-Type"))
			}

			if tt.data != nil {
				var response interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("RespondJSON() failed to decode response: %v", err)
				}
			}
		})
	}
}

func TestRespondSuccess(t *testing.T) {
	tests := []struct {
		name    string
		message string
		data    interface{}
	}{
		{
			name:    "With data",
			message: "Operation successful",
			data:    map[string]string{"id": "123"},
		},
		{
			name:    "Without data",
			message: "Operation successful",
			data:    nil,
		},
		{
			name:    "With array data",
			message: "List fetched",
			data:    []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondSuccess(w, tt.message, tt.data)

			if w.Code != http.StatusOK {
				t.Errorf("RespondSuccess() status = %v, want %v", w.Code, http.StatusOK)
			}

			var response Response
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if !response.Success {
				t.Error("RespondSuccess() Success field should be true")
			}

			if response.Message != tt.message {
				t.Errorf("RespondSuccess() Message = %v, want %v", response.Message, tt.message)
			}

			if response.Error != "" {
				t.Error("RespondSuccess() Error field should be empty")
			}
		})
	}
}

func TestRespondError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "Bad request",
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
		},
		{
			name:       "Unauthorized",
			statusCode: http.StatusUnauthorized,
			message:    "Authentication required",
		},
		{
			name:       "Internal error",
			statusCode: http.StatusInternalServerError,
			message:    "Server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondError(w, tt.statusCode, tt.message)

			if w.Code != tt.statusCode {
				t.Errorf("RespondError() status = %v, want %v", w.Code, tt.statusCode)
			}

			var response Response
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Success {
				t.Error("RespondError() Success field should be false")
			}

			if response.Error != tt.message {
				t.Errorf("RespondError() Error = %v, want %v", response.Error, tt.message)
			}

			if response.Message != "" {
				t.Error("RespondError() Message field should be empty")
			}

			if response.Data != nil {
				t.Error("RespondError() Data field should be nil")
			}
		})
	}
}

func TestRespondCreated(t *testing.T) {
	tests := []struct {
		name    string
		message string
		data    interface{}
	}{
		{
			name:    "New resource created",
			message: "User created successfully",
			data:    map[string]interface{}{"id": 1, "name": "John"},
		},
		{
			name:    "Created without data",
			message: "Resource created",
			data:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondCreated(w, tt.message, tt.data)

			if w.Code != http.StatusCreated {
				t.Errorf("RespondCreated() status = %v, want %v", w.Code, http.StatusCreated)
			}

			var response Response
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if !response.Success {
				t.Error("RespondCreated() Success field should be true")
			}

			if response.Message != tt.message {
				t.Errorf("RespondCreated() Message = %v, want %v", response.Message, tt.message)
			}

			if response.Error != "" {
				t.Error("RespondCreated() Error field should be empty")
			}
		})
	}
}

func TestResponseStructure(t *testing.T) {
	// Test that Response struct properly marshals to JSON
	response := Response{
		Success: true,
		Message: "Test message",
		Data:    map[string]string{"key": "value"},
		Error:   "",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal Response: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Response: %v", err)
	}

	// Check that fields are present in JSON
	if _, ok := decoded["success"]; !ok {
		t.Error("Response JSON missing 'success' field")
	}

	// Check omitempty works
	if decoded["error"] != nil && decoded["error"] != "" {
		t.Error("Empty error field should be omitted or empty in JSON")
	}
}