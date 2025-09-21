package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidatePassword checks if password meets requirements
func ValidatePassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "Password must be at least 6 characters long"
	}
	if len(password) > 100 {
		return false, "Password must not exceed 100 characters"
	}
	return true, ""
}

// ValidateName checks if name is valid
func ValidateName(name string) (bool, string) {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return false, "Name must be at least 2 characters long"
	}
	if len(name) > 100 {
		return false, "Name must not exceed 100 characters"
	}

	// Check if name contains only valid characters
	validName := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !validName.MatchString(name) {
		return false, "Name can only contain letters and spaces"
	}

	return true, ""
}

// ValidateTaskTitle checks if task title is valid
func ValidateTaskTitle(title string) (bool, string) {
	title = strings.TrimSpace(title)
	if len(title) < 1 {
		return false, "Title is required"
	}
	if len(title) > 255 {
		return false, "Title must not exceed 255 characters"
	}
	return true, ""
}

// ValidateTaskDescription checks if task description is valid
func ValidateTaskDescription(description string) (bool, string) {
	if len(description) > 1000 {
		return false, "Description must not exceed 1000 characters"
	}
	return true, ""
}

// ValidateTaskStatus checks if task status is valid
func ValidateTaskStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
	}
	return validStatuses[status]
}