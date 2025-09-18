package common

import (
	"fmt"
	"strings"
)

// TimeProvider abstracts time operations for domain purity
type TimeProvider interface {
	Now() Timestamp
}

// defaultTimeProvider uses system time
type systemTimeProvider struct{}

func (p *systemTimeProvider) Now() Timestamp {
	return Now()
}

// Default time provider
var defaultTimeProvider TimeProvider = &systemTimeProvider{}

// SetTimeProvider allows injection of time provider (useful for testing)
func SetTimeProvider(provider TimeProvider) {
	defaultTimeProvider = provider
}

// GetCurrentTime returns current time using the configured provider
func GetCurrentTime() Timestamp {
	return defaultTimeProvider.Now()
}

// URLValidator abstracts URL validation for domain purity
type URLValidator interface {
	ValidateURL(url string) error
	ValidateScheme(url string, allowedSchemes []string) error
	ExtractScheme(url string) string
	HasHost(url string) bool
}

// basicURLValidator provides basic URL validation without external dependencies
type basicURLValidator struct{}

func (v *basicURLValidator) ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.Contains(url, "://") {
		return fmt.Errorf("URL must contain scheme (://)")
	}

	return nil
}

func (v *basicURLValidator) ValidateScheme(url string, allowedSchemes []string) error {
	scheme := v.ExtractScheme(url)
	if scheme == "" {
		return fmt.Errorf("URL must have a scheme")
	}

	for _, allowed := range allowedSchemes {
		if strings.ToLower(scheme) == strings.ToLower(allowed) {
			return nil
		}
	}

	return fmt.Errorf("unsupported scheme: %s", scheme)
}

func (v *basicURLValidator) ExtractScheme(url string) string {
	parts := strings.Split(url, "://")
	if len(parts) < 2 {
		return ""
	}
	return strings.ToLower(parts[0])
}

func (v *basicURLValidator) HasHost(url string) bool {
	parts := strings.Split(url, "://")
	if len(parts) < 2 {
		return false
	}
	return parts[1] != ""
}

// Default URL validator
var defaultURLValidator URLValidator = &basicURLValidator{}

// SetURLValidator allows injection of URL validator
func SetURLValidator(validator URLValidator) {
	defaultURLValidator = validator
}

// GetURLValidator returns the configured URL validator
func GetURLValidator() URLValidator {
	return defaultURLValidator
}
