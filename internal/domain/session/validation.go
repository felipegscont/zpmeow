package session

import (
	"net/url"
	"regexp"
	"strings"
)

// ValidationService provides centralized validation for session-related inputs
type ValidationService struct {
	phoneRegex    *regexp.Regexp
	sanitizeRegex *regexp.Regexp
}

// NewValidationService creates a new validation service instance
func NewValidationService() *ValidationService {
	return &ValidationService{
		phoneRegex:    regexp.MustCompile(`\D`), // Non-digit characters
		sanitizeRegex: regexp.MustCompile(`[\x00-\x1f\x7f]`), // Control characters
	}
}

// Phone number validation

// ValidatePhoneNumber validates a phone number format
func (v *ValidationService) ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return NewDomainError("phone number cannot be empty")
	}

	// Remove all non-digit characters
	cleaned := v.phoneRegex.ReplaceAllString(phone, "")
	
	// Check if it's a valid length (10-15 digits according to E.164)
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return NewDomainError("phone number must be between 10 and 15 digits")
	}
	
	return nil
}

// Session name validation (comprehensive business rules)

// ValidateSessionName validates a session name for URL safety and business rules
func (v *ValidationService) ValidateSessionName(name string) error {
	if name == "" {
		return ErrInvalidSessionName
	}

	// Check length constraints
	if len(name) < 3 {
		return ErrSessionNameTooShort
	}
	if len(name) > 50 {
		return ErrSessionNameTooLong
	}

	// Check for URL-safe characters only
	// Allow: letters, numbers, hyphens, underscores
	// Disallow: spaces, special characters, unicode
	for i, r := range name {
		if !v.isValidSessionNameChar(r) {
			return ErrInvalidSessionNameChar
		}

		// Don't allow starting or ending with hyphen or underscore
		if i == 0 || i == len(name)-1 {
			if r == '-' || r == '_' {
				return ErrInvalidSessionNameFormat
			}
		}
	}

	// Don't allow reserved names
	if v.isReservedSessionName(name) {
		return ErrReservedSessionName
	}

	return nil
}

// isValidSessionNameChar checks if a character is valid for session names
func (v *ValidationService) isValidSessionNameChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_'
}

// isReservedSessionName checks if a name is reserved
func (v *ValidationService) isReservedSessionName(name string) bool {
	reserved := map[string]bool{
		"api": true, "admin": true, "root": true, "system": true, "config": true, "settings": true,
		"health": true, "ping": true, "status": true, "info": true, "debug": true, "test": true,
		"sessions": true, "session": true, "create": true, "list": true, "delete": true,
		"connect": true, "disconnect": true, "logout": true, "qr": true, "pair": true,
		"send": true, "chat": true, "group": true, "contact": true, "media": true, "file": true,
		"swagger": true, "docs": true, "documentation": true,
	}

	return reserved[strings.ToLower(name)]
}

// Proxy URL validation

// ValidateProxyURL validates a proxy URL format
func (v *ValidationService) ValidateProxyURL(proxyURL string) error {
	if proxyURL == "" {
		return nil // Empty is valid (no proxy)
	}

	// Parse the URL to validate format
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return NewDomainError("invalid proxy URL format")
	}

	// Check supported schemes
	supportedSchemes := map[string]bool{
		"http":   true,
		"https":  true,
		"socks5": true,
	}

	if !supportedSchemes[parsedURL.Scheme] {
		return NewDomainError("proxy URL must use http, https, or socks5 scheme")
	}

	// Ensure host is present
	if parsedURL.Host == "" {
		return NewDomainError("proxy URL must include a host")
	}

	return nil
}

// Session ID validation

// ValidateSessionID validates a session ID
func (v *ValidationService) ValidateSessionID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidSessionID
	}
	return nil
}

// ValidateSessionIDOrName validates a session ID or name
func (v *ValidationService) ValidateSessionIDOrName(idOrName string) error {
	if strings.TrimSpace(idOrName) == "" {
		return ErrInvalidSessionID // Use ID error as it's more generic
	}
	return nil
}

// String sanitization

// SanitizeString removes potentially harmful characters from a string
func (v *ValidationService) SanitizeString(input string) string {
	// Remove null bytes and control characters
	cleaned := v.sanitizeRegex.ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}

// Global validation service instance
var DefaultValidationService = NewValidationService()

// Convenience functions using the default validation service

// ValidatePhoneNumber validates a phone number using the default service
func ValidatePhoneNumber(phone string) error {
	return DefaultValidationService.ValidatePhoneNumber(phone)
}

// ValidateSessionName validates a session name using the default service
func ValidateSessionName(name string) error {
	return DefaultValidationService.ValidateSessionName(name)
}

// ValidateProxyURL validates a proxy URL using the default service
func ValidateProxyURL(proxyURL string) error {
	return DefaultValidationService.ValidateProxyURL(proxyURL)
}

// ValidateSessionID validates a session ID using the default service
func ValidateSessionID(id string) error {
	return DefaultValidationService.ValidateSessionID(id)
}

// ValidateSessionIDOrName validates a session ID or name using the default service
func ValidateSessionIDOrName(idOrName string) error {
	return DefaultValidationService.ValidateSessionIDOrName(idOrName)
}

// SanitizeString sanitizes a string using the default service
func SanitizeString(input string) string {
	return DefaultValidationService.SanitizeString(input)
}

// Legacy compatibility functions for utils package

// IsValidPhoneNumber checks if a phone number is valid (legacy compatibility)
func IsValidPhoneNumber(phone string) bool {
	return ValidatePhoneNumber(phone) == nil
}

// IsValidProxyURL checks if a proxy URL is valid (legacy compatibility)
func IsValidProxyURL(proxyURL string) bool {
	return ValidateProxyURL(proxyURL) == nil
}
