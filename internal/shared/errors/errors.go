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



// Global instance for convenience
var Error = &ErrorUtils{}
