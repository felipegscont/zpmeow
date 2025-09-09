package types

import "time"

// Common types used across the application

// ID represents a unique identifier
type ID string

// Timestamp represents a timestamp
type Timestamp time.Time

// Status represents a status string
type Status string

// Common status values
const (
	StatusDisconnected Status = "disconnected"
	StatusConnecting   Status = "connecting"
	StatusConnected    Status = "connected"
	StatusError        Status = "error"
)

// String returns the string representation of Status
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusDisconnected, StatusConnecting, StatusConnected, StatusError:
		return true
	default:
		return false
	}
}
