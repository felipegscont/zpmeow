package utils

import (
	"fmt"
	"strings"

	waTypes "go.mau.fi/whatsmeow/types"
)

const (
	ErrInvalidJID = "invalid JID"
)

// JIDUtils provides utilities for working with WhatsApp JIDs
type JIDUtils struct{}

// ParseJID parses a string into a WhatsApp JID
func (ju *JIDUtils) ParseJID(arg string) (waTypes.JID, error) {
	if arg == "" {
		return waTypes.JID{}, fmt.Errorf(ErrInvalidJID + ": empty string")
	}

	// Remove leading + if present
	if arg[0] == '+' {
		arg = arg[1:]
	}

	// If no @ symbol, assume it's a phone number for default user server
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

// IsValidJID checks if a string is a valid JID
func (ju *JIDUtils) IsValidJID(jidStr string) bool {
	_, err := ju.ParseJID(jidStr)
	return err == nil
}

// Global instance for convenience
var JID = &JIDUtils{}
