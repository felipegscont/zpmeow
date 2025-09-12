package types

import "time"

type ID string

type Timestamp time.Time

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

// SendResponse represents a response from sending a message
type SendResponse struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Success   bool   `json:"success"`
}
