package meow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"zpmeow/internal/types"

	waTypes "go.mau.fi/whatsmeow/types"
)

// Note: Constants moved to constants.go file

// JIDUtils provides utility functions for JID operations
type JIDUtils struct{}

// ParseJID parses a phone number or JID string into a WhatsApp JID
func (ju *JIDUtils) ParseJID(arg string) (waTypes.JID, error) {
	if arg == "" {
		return waTypes.JID{}, fmt.Errorf(ErrInvalidJID + ": empty string")
	}

	// Remove leading + if present
	if arg[0] == '+' {
		arg = arg[1:]
	}

	// If no @ symbol, assume it's a phone number
	if !strings.ContainsRune(arg, '@') {
		return waTypes.NewJID(arg, waTypes.DefaultUserServer), nil
	}

	// Parse as full JID
	recipient, err := waTypes.ParseJID(arg)
	if err != nil {
		return waTypes.JID{}, fmt.Errorf(ErrInvalidJID + ": %w", err)
	}

	if recipient.User == "" {
		return waTypes.JID{}, fmt.Errorf(ErrInvalidJID + ": no user specified")
	}

	return recipient, nil
}

// IsValidJID checks if a JID string is valid
func (ju *JIDUtils) IsValidJID(jidStr string) bool {
	_, err := ju.ParseJID(jidStr)
	return err == nil
}

// Note: Validation utilities moved to validation.go file

// ErrorUtils provides error handling utilities
type ErrorUtils struct{}

// WrapError wraps an error with additional context
func (eu *ErrorUtils) WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// IsRetryableError determines if an error should trigger a retry
func (eu *ErrorUtils) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	retryableErrors := []string{
		"connection refused",
		"timeout",
		"temporary failure",
		"network unreachable",
	}
	
	for _, retryable := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), retryable) {
			return true
		}
	}
	
	return false
}

// ContextUtils provides context-related utilities
type ContextUtils struct{}

// WithTimeout creates a context with the default timeout
func (cu *ContextUtils) WithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, DefaultTimeout)
}

// WithCustomTimeout creates a context with a custom timeout
func (cu *ContextUtils) WithCustomTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, timeout)
}

// StatusUtils provides status-related utilities
type StatusUtils struct{}

// IsConnectedStatus checks if the status indicates a connected state
func (su *StatusUtils) IsConnectedStatus(status types.Status) bool {
	return status == types.StatusConnected
}

// IsDisconnectedStatus checks if the status indicates a disconnected state
func (su *StatusUtils) IsDisconnectedStatus(status types.Status) bool {
	return status == types.StatusDisconnected
}

// IsErrorStatus checks if the status indicates an error state
func (su *StatusUtils) IsErrorStatus(status types.Status) bool {
	return status == types.StatusError
}

// ChannelUtils provides channel operation utilities
type ChannelUtils struct{}

// SafeChannelSend sends to a channel with timeout to prevent blocking
func (cu *ChannelUtils) SafeChannelSend(ch chan bool, value bool, timeout time.Duration) bool {
	select {
	case ch <- value:
		return true
	case <-time.After(timeout):
		return false
	default:
		return false
	}
}

// SafeChannelClose closes a channel safely
func (cu *ChannelUtils) SafeChannelClose(ch chan bool) {
	select {
	case <-ch:
		// Channel already closed or has value
	default:
		close(ch)
	}
}

// Global utility instances
var (
	JID     = &JIDUtils{}
	Error   = &ErrorUtils{}
	Context = &ContextUtils{}
	Status  = &StatusUtils{}
	Channel = &ChannelUtils{}

	// Use the global validator from validation.go
	Validation = DefaultValidator
)
