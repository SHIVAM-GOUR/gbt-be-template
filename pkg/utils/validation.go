package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail checks if the email format is valid
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if password meets minimum requirements
func IsValidPassword(password string) bool {
	// At least 6 characters
	if len(password) < 6 {
		return false
	}
	return true
}

// IsValidUsername checks if username is valid
func IsValidUsername(username string) bool {
	// 3-50 characters, alphanumeric and underscore only
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// SanitizeString removes leading/trailing whitespace and converts to lowercase
func SanitizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// TruncateString truncates a string to the specified length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}
