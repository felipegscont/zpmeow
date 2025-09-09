package utils

import (
	"regexp"
	"strings"
)

// IsValidPhoneNumber validates a phone number format
func IsValidPhoneNumber(phone string) bool {
	// Remove all non-digit characters
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	// Check if it's a valid length (10-15 digits)
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return false
	}
	
	return true
}

// IsValidSessionName validates a session name
func IsValidSessionName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) >= 1 && len(name) <= 100
}

// IsValidProxyURL validates a proxy URL format
func IsValidProxyURL(url string) bool {
	if url == "" {
		return true // Empty is valid (no proxy)
	}
	
	// Basic URL validation
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "socks5://")
}

// SanitizeString removes potentially harmful characters from a string
func SanitizeString(input string) string {
	// Remove null bytes and control characters
	cleaned := regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}
