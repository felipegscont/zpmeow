package types

import "time"

// Generic types used across the application

type ID string

type Timestamp time.Time

// Status represents the connection status of a session
type Status string

const (
	StatusDisconnected Status = "disconnected"
	StatusConnecting   Status = "connecting"
	StatusConnected    Status = "connected"
	StatusError        Status = "error"
	StatusDeleted      Status = "deleted"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	switch s {
	case StatusDisconnected, StatusConnecting, StatusConnected, StatusError, StatusDeleted:
		return true
	default:
		return false
	}
}

// SendResponse represents a generic response for send operations
type SendResponse struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Success   bool   `json:"success"`
}
