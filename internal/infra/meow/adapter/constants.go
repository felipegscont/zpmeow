package adapter

import "errors"

// Validation patterns
const (
	PhoneNumberPattern = `^\+?[1-9]\d{1,14}$`
	SessionIDPattern   = `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	JIDPattern         = `^[0-9]+@[a-z.]+$`
)

// Error constants
var (
	ErrEmptySessionID = errors.New("session ID cannot be empty")
	ErrInvalidJID     = errors.New("invalid JID format")
)

// Timeout constants
const (
	DefaultTimeout = 30 // seconds
)

// MeowClient placeholder type (will be properly defined in core)
type MeowClient interface {
	// Placeholder interface
}
