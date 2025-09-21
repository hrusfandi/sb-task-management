package utils

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"Valid email", "test@example.com", true},
		{"Valid email with subdomain", "user@mail.example.com", true},
		{"Valid email with plus", "user+tag@example.com", true},
		{"Valid email with dots", "first.last@example.com", true},
		{"Invalid - no @", "testexample.com", false},
		{"Invalid - no domain", "test@", false},
		{"Invalid - no local part", "@example.com", false},
		{"Invalid - spaces", "test @example.com", false},
		{"Invalid - double @", "test@@example.com", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantOk   bool
		wantMsg  string
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantOk:   true,
			wantMsg:  "",
		},
		{
			name:     "Minimum length",
			password: "123456",
			wantOk:   true,
			wantMsg:  "",
		},
		{
			name:     "Too short",
			password: "12345",
			wantOk:   false,
			wantMsg:  "Password must be at least 6 characters long",
		},
		{
			name:     "Too long",
			password: strings.Repeat("a", 101),
			wantOk:   false,
			wantMsg:  "Password must not exceed 100 characters",
		},
		{
			name:     "Empty password",
			password: "",
			wantOk:   false,
			wantMsg:  "Password must be at least 6 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := ValidatePassword(tt.password)
			if ok != tt.wantOk {
				t.Errorf("ValidatePassword() ok = %v, want %v", ok, tt.wantOk)
			}
			if msg != tt.wantMsg {
				t.Errorf("ValidatePassword() msg = %v, want %v", msg, tt.wantMsg)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOk  bool
		wantMsg string
	}{
		{
			name:    "Valid name",
			input:   "John Doe",
			wantOk:  true,
			wantMsg: "",
		},
		{
			name:    "Single name",
			input:   "John",
			wantOk:  true,
			wantMsg: "",
		},
		{
			name:    "Minimum length",
			input:   "Jo",
			wantOk:  true,
			wantMsg: "",
		},
		{
			name:    "Too short",
			input:   "J",
			wantOk:  false,
			wantMsg: "Name must be at least 2 characters long",
		},
		{
			name:    "Too long",
			input:   strings.Repeat("a", 101),
			wantOk:  false,
			wantMsg: "Name must not exceed 100 characters",
		},
		{
			name:    "With numbers",
			input:   "John123",
			wantOk:  false,
			wantMsg: "Name can only contain letters and spaces",
		},
		{
			name:    "With special characters",
			input:   "John@Doe",
			wantOk:  false,
			wantMsg: "Name can only contain letters and spaces",
		},
		{
			name:    "Empty name",
			input:   "",
			wantOk:  false,
			wantMsg: "Name must be at least 2 characters long",
		},
		{
			name:    "Only spaces",
			input:   "   ",
			wantOk:  false,
			wantMsg: "Name must be at least 2 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := ValidateName(tt.input)
			if ok != tt.wantOk {
				t.Errorf("ValidateName(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if msg != tt.wantMsg {
				t.Errorf("ValidateName(%q) msg = %v, want %v", tt.input, msg, tt.wantMsg)
			}
		})
	}
}

func TestValidateTaskTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantOk  bool
		wantMsg string
	}{
		{
			name:    "Valid title",
			title:   "Complete project documentation",
			wantOk:  true,
			wantMsg: "",
		},
		{
			name:    "Single character",
			title:   "A",
			wantOk:  true,
			wantMsg: "",
		},
		{
			name:    "Maximum length",
			title:   strings.Repeat("a", 255),
			wantOk:  true,
			wantMsg: "",
		},
		{
			name:    "Too long",
			title:   strings.Repeat("a", 256),
			wantOk:  false,
			wantMsg: "Title must not exceed 255 characters",
		},
		{
			name:    "Empty title",
			title:   "",
			wantOk:  false,
			wantMsg: "Title is required",
		},
		{
			name:    "Only spaces",
			title:   "   ",
			wantOk:  false,
			wantMsg: "Title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := ValidateTaskTitle(tt.title)
			if ok != tt.wantOk {
				t.Errorf("ValidateTaskTitle() ok = %v, want %v", ok, tt.wantOk)
			}
			if msg != tt.wantMsg {
				t.Errorf("ValidateTaskTitle() msg = %v, want %v", msg, tt.wantMsg)
			}
		})
	}
}

func TestValidateTaskDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantOk      bool
		wantMsg     string
	}{
		{
			name:        "Valid description",
			description: "This is a task description",
			wantOk:      true,
			wantMsg:     "",
		},
		{
			name:        "Empty description",
			description: "",
			wantOk:      true,
			wantMsg:     "",
		},
		{
			name:        "Maximum length",
			description: strings.Repeat("a", 1000),
			wantOk:      true,
			wantMsg:     "",
		},
		{
			name:        "Too long",
			description: strings.Repeat("a", 1001),
			wantOk:      false,
			wantMsg:     "Description must not exceed 1000 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := ValidateTaskDescription(tt.description)
			if ok != tt.wantOk {
				t.Errorf("ValidateTaskDescription() ok = %v, want %v", ok, tt.wantOk)
			}
			if msg != tt.wantMsg {
				t.Errorf("ValidateTaskDescription() msg = %v, want %v", msg, tt.wantMsg)
			}
		})
	}
}

func TestValidateTaskStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"Valid - pending", "pending", true},
		{"Valid - in_progress", "in_progress", true},
		{"Valid - completed", "completed", true},
		{"Invalid - empty", "", false},
		{"Invalid - unknown status", "cancelled", false},
		{"Invalid - uppercase", "PENDING", false},
		{"Invalid - with space", "in progress", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateTaskStatus(tt.status); got != tt.want {
				t.Errorf("ValidateTaskStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}