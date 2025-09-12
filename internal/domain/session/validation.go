package session

import (
	"net/url"
	"regexp"
	"strings"

	"zpmeow/internal/shared/types"
)

// Domain validation rules - pure business logic only

var (
	// Session name validation
	sessionNameMinLength = 3
	sessionNameMaxLength = 50
	sessionNameRegex     = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	
	// Reserved session names
	reservedNames = map[string]bool{
		"admin":   true,
		"system":  true,
		"default": true,
		"test":    true,
	}
)

// ValidateSessionName validates session name according to domain rules
func ValidateSessionName(name string) error {
	if name == "" {
		return ErrInvalidSessionName
	}
	
	if len(name) < sessionNameMinLength {
		return ErrSessionNameTooShort
	}
	
	if len(name) > sessionNameMaxLength {
		return ErrSessionNameTooLong
	}
	
	if !sessionNameRegex.MatchString(name) {
		return ErrInvalidSessionNameChar
	}
	
	if reservedNames[strings.ToLower(name)] {
		return ErrReservedSessionName
	}
	
	return nil
}

// ValidateSessionID validates session ID format
func ValidateSessionID(id string) error {
	if id == "" {
		return ErrInvalidSessionID
	}
	
	// UUID format validation
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(id) {
		return ErrInvalidSessionID
	}
	
	return nil
}

// ValidateProxyURL validates proxy URL format
func ValidateProxyURL(proxyURL string) error {
	if proxyURL == "" {
		return nil // Empty proxy URL is valid (means no proxy)
	}
	
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return ErrInvalidProxyURL
	}
	
	if parsedURL.Scheme == "" {
		return ErrInvalidProxyURL
	}
	
	if parsedURL.Host == "" {
		return ErrInvalidProxyURL
	}
	
	return nil
}

// ValidateSessionStatus validates session status transitions
func ValidateSessionStatus(currentStatus, newStatus types.Status) error {
	// Define valid status transitions
	validTransitions := map[types.Status][]types.Status{
		types.StatusDisconnected: {types.StatusConnecting},
		types.StatusConnecting:   {types.StatusConnected, types.StatusDisconnected, types.StatusError},
		types.StatusConnected:    {types.StatusDisconnected, types.StatusError},
		types.StatusError:        {types.StatusDisconnected, types.StatusConnecting},
	}
	
	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return ErrInvalidSessionStatus
	}
	
	for _, allowed := range allowedStatuses {
		if newStatus == allowed {
			return nil
		}
	}
	
	return ErrInvalidSessionStatus
}
