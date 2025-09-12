package types

import (
	"fmt"
	"regexp"
	"strings"
)

// SessionID represents a session identifier with validation
type SessionID struct {
	value string
}

// NewSessionID creates a new SessionID with validation
func NewSessionID(value string) (SessionID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return SessionID{}, fmt.Errorf("session ID cannot be empty")
	}

	// Allow both UUID format and alphanumeric format
	if !isValidSessionIDFormat(trimmed) {
		return SessionID{}, fmt.Errorf("invalid session ID format")
	}

	return SessionID{value: trimmed}, nil
}

// MustNewSessionID creates a new SessionID, panicking on error (use for constants)
func MustNewSessionID(value string) SessionID {
	id, err := NewSessionID(value)
	if err != nil {
		panic(fmt.Sprintf("invalid session ID: %v", err))
	}
	return id
}

// String returns the string representation of the SessionID
func (s SessionID) String() string {
	return s.value
}

// IsEmpty checks if the SessionID is empty
func (s SessionID) IsEmpty() bool {
	return s.value == ""
}

// Equals checks if two SessionIDs are equal
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}

// SessionName represents a session name with validation
type SessionName struct {
	value string
}

// NewSessionName creates a new SessionName with validation
func NewSessionName(value string) (SessionName, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return SessionName{}, fmt.Errorf("session name cannot be empty")
	}

	if len(trimmed) < 3 {
		return SessionName{}, fmt.Errorf("session name must be at least 3 characters long")
	}

	if len(trimmed) > 50 {
		return SessionName{}, fmt.Errorf("session name must be at most 50 characters long")
	}

	if !isValidSessionNameFormat(trimmed) {
		return SessionName{}, fmt.Errorf("session name can only contain letters, numbers, underscores, and hyphens")
	}

	if isReservedSessionName(trimmed) {
		return SessionName{}, fmt.Errorf("'%s' is a reserved session name", trimmed)
	}

	return SessionName{value: trimmed}, nil
}

// MustNewSessionName creates a new SessionName, panicking on error
func MustNewSessionName(value string) SessionName {
	name, err := NewSessionName(value)
	if err != nil {
		panic(fmt.Sprintf("invalid session name: %v", err))
	}
	return name
}

// String returns the string representation of the SessionName
func (s SessionName) String() string {
	return s.value
}

// IsEmpty checks if the SessionName is empty
func (s SessionName) IsEmpty() bool {
	return s.value == ""
}

// Equals checks if two SessionNames are equal
func (s SessionName) Equals(other SessionName) bool {
	return s.value == other.value
}

// ProxyURL represents a proxy URL with validation
type ProxyURL struct {
	value string
}

// NewProxyURL creates a new ProxyURL with validation
func NewProxyURL(value string) (ProxyURL, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ProxyURL{}, nil // Empty proxy URL is valid (means no proxy)
	}

	if !isValidProxyURLFormat(trimmed) {
		return ProxyURL{}, fmt.Errorf("invalid proxy URL format")
	}

	return ProxyURL{value: trimmed}, nil
}

// MustNewProxyURL creates a new ProxyURL, panicking on error
func MustNewProxyURL(value string) ProxyURL {
	url, err := NewProxyURL(value)
	if err != nil {
		panic(fmt.Sprintf("invalid proxy URL: %v", err))
	}
	return url
}

// String returns the string representation of the ProxyURL
func (p ProxyURL) String() string {
	return p.value
}

// IsEmpty checks if the ProxyURL is empty
func (p ProxyURL) IsEmpty() bool {
	return p.value == ""
}

// Equals checks if two ProxyURLs are equal
func (p ProxyURL) Equals(other ProxyURL) bool {
	return p.value == other.value
}

// PhoneNumber represents a phone number with validation
type PhoneNumber struct {
	value string
}

// NewPhoneNumber creates a new PhoneNumber with validation
func NewPhoneNumber(value string) (PhoneNumber, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return PhoneNumber{}, fmt.Errorf("phone number cannot be empty")
	}

	if !isValidPhoneNumberFormat(trimmed) {
		return PhoneNumber{}, fmt.Errorf("invalid phone number format")
	}

	return PhoneNumber{value: trimmed}, nil
}

// MustNewPhoneNumber creates a new PhoneNumber, panicking on error
func MustNewPhoneNumber(value string) PhoneNumber {
	phone, err := NewPhoneNumber(value)
	if err != nil {
		panic(fmt.Sprintf("invalid phone number: %v", err))
	}
	return phone
}

// String returns the string representation of the PhoneNumber
func (p PhoneNumber) String() string {
	return p.value
}

// IsEmpty checks if the PhoneNumber is empty
func (p PhoneNumber) IsEmpty() bool {
	return p.value == ""
}

// Equals checks if two PhoneNumbers are equal
func (p PhoneNumber) Equals(other PhoneNumber) bool {
	return p.value == other.value
}

// Validation helper functions

var (
	sessionIDRegex     = regexp.MustCompile(`^[a-fA-F0-9]{32}$|^[a-zA-Z0-9_-]{3,50}$`)
	sessionNameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	proxyURLRegex      = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	phoneNumberRegex   = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	reservedNames      = map[string]bool{
		"admin":   true,
		"system":  true,
		"default": true,
		"test":    true,
		"api":     true,
		"www":     true,
	}
)

func isValidSessionIDFormat(value string) bool {
	return sessionIDRegex.MatchString(value)
}

func isValidSessionNameFormat(value string) bool {
	return sessionNameRegex.MatchString(value)
}

func isReservedSessionName(value string) bool {
	return reservedNames[strings.ToLower(value)]
}

func isValidProxyURLFormat(value string) bool {
	return proxyURLRegex.MatchString(value)
}

func isValidPhoneNumberFormat(value string) bool {
	return phoneNumberRegex.MatchString(value)
}

// Status transitions validation

// CanTransitionTo checks if a status can transition to another status
func (s Status) CanTransitionTo(target Status) bool {
	switch s {
	case StatusDisconnected:
		return target == StatusConnecting || target == StatusDeleted
	case StatusConnecting:
		return target == StatusConnected || target == StatusError || target == StatusDisconnected
	case StatusConnected:
		return target == StatusDisconnected || target == StatusError
	case StatusError:
		return target == StatusDisconnected || target == StatusConnecting
	case StatusDeleted:
		return false // Cannot transition from deleted
	default:
		return false
	}
}

// GetValidTransitions returns all valid transitions from current status
func (s Status) GetValidTransitions() []Status {
	switch s {
	case StatusDisconnected:
		return []Status{StatusConnecting, StatusDeleted}
	case StatusConnecting:
		return []Status{StatusConnected, StatusError, StatusDisconnected}
	case StatusConnected:
		return []Status{StatusDisconnected, StatusError}
	case StatusError:
		return []Status{StatusDisconnected, StatusConnecting}
	case StatusDeleted:
		return []Status{} // No transitions from deleted
	default:
		return []Status{}
	}
}
