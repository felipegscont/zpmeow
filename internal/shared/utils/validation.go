package utils

import (
	"regexp"
	"strings"
)

func IsValidPhoneNumber(phone string) bool {
	// Check if it's a group JID (ends with @g.us)
	if strings.HasSuffix(phone, "@g.us") {
		return true
	}

	// Check if it's a complete JID (contains @)
	if strings.Contains(phone, "@") {
		return true
	}

	// Remove all non-digit characters for phone number validation
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(cleaned) < 10 || len(cleaned) > 15 {
		return false
	}

	return true
}

func IsValidSessionName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) >= 1 && len(name) <= 100
}

func IsValidProxyURL(url string) bool {
	if url == "" {
		return true // Empty is valid (no proxy)
	}

	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "socks5://")
}

func SanitizeString(input string) string {

	cleaned := regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}
