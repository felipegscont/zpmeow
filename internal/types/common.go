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
)


func (s Status) String() string {
	return string(s)
}


func (s Status) IsValid() bool {
	switch s {
	case StatusDisconnected, StatusConnecting, StatusConnected, StatusError:
		return true
	default:
		return false
	}
}
