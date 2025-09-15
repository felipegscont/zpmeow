package utils

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidURL(url string) bool {
	if url == "" {
		return false
	}
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func SanitizeString(input string) string {

	cleaned := regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}

func GenerateUUID() string {
	return uuid.New().String()
}
