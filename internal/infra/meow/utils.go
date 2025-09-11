package meow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"zpmeow/internal/types"

	waTypes "go.mau.fi/whatsmeow/types"
)




type JIDUtils struct{}


func (ju *JIDUtils) ParseJID(arg string) (waTypes.JID, error) {
	if arg == "" {
		return waTypes.JID{}, fmt.Errorf(ErrInvalidJID + ": empty string")
	}


	if arg[0] == '+' {
		arg = arg[1:]
	}


	if !strings.ContainsRune(arg, '@') {
		return waTypes.NewJID(arg, waTypes.DefaultUserServer), nil
	}


	recipient, err := waTypes.ParseJID(arg)
	if err != nil {
		return waTypes.JID{}, fmt.Errorf(ErrInvalidJID + ": %w", err)
	}

	if recipient.User == "" {
		return waTypes.JID{}, fmt.Errorf(ErrInvalidJID + ": no user specified")
	}

	return recipient, nil
}


func (ju *JIDUtils) IsValidJID(jidStr string) bool {
	_, err := ju.ParseJID(jidStr)
	return err == nil
}




type ErrorUtils struct{}


func (eu *ErrorUtils) WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}


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


type ContextUtils struct{}


func (cu *ContextUtils) WithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, DefaultTimeout)
}


func (cu *ContextUtils) WithCustomTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, timeout)
}


type StatusUtils struct{}


func (su *StatusUtils) IsConnectedStatus(status types.Status) bool {
	return status == types.StatusConnected
}


func (su *StatusUtils) IsDisconnectedStatus(status types.Status) bool {
	return status == types.StatusDisconnected
}


func (su *StatusUtils) IsErrorStatus(status types.Status) bool {
	return status == types.StatusError
}


type ChannelUtils struct{}


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


func (cu *ChannelUtils) SafeChannelClose(ch chan bool) {
	select {
	case <-ch:

	default:
		close(ch)
	}
}


var (
	JID     = &JIDUtils{}
	Error   = &ErrorUtils{}
	Context = &ContextUtils{}
	Status  = &StatusUtils{}
	Channel = &ChannelUtils{}


	Validation = DefaultValidator
)
