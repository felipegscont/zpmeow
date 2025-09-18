package types

import "fmt"

type SessionID struct {
	value string
}

func NewSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, fmt.Errorf("session ID cannot be empty")
	}

	if len(value) < 1 || len(value) > 100 {
		return SessionID{}, fmt.Errorf("session ID must be between 1 and 100 characters")
	}

	return SessionID{value: value}, nil
}

func (s SessionID) Value() string {
	return s.value
}

func (s SessionID) String() string {
	return s.value
}

func (s SessionID) IsEmpty() bool {
	return s.value == ""
}

func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}
