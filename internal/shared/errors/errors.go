package errors

import (
	"fmt"
)

// ErrorUtils provides utilities for error handling
type ErrorUtils struct{}

// WrapError wraps an error with additional context
func (eu *ErrorUtils) WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// ValidationError represents a validation error
type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) error {
	return &ValidationError{message: message}
}



// Global instance for convenience
var Error = &ErrorUtils{}
