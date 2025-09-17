package types

import "fmt"

// SessionID represents a session identifier (shared across domains)
type SessionID struct {
	value string
}

// NewSessionID creates a new SessionID
func NewSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, fmt.Errorf("session ID cannot be empty")
	}

	if len(value) < 1 || len(value) > 100 {
		return SessionID{}, fmt.Errorf("session ID must be between 1 and 100 characters")
	}

	return SessionID{value: value}, nil
}

// Value returns the session ID value
func (s SessionID) Value() string {
	return s.value
}

// String returns the string representation
func (s SessionID) String() string {
	return s.value
}

// IsEmpty checks if the session ID is empty
func (s SessionID) IsEmpty() bool {
	return s.value == ""
}

// Equals checks if two SessionIDs are equal
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}
