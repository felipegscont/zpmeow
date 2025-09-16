package errors

import (
	"fmt"
)

type ErrorUtils struct{}

func (eu *ErrorUtils) WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

func NewValidationError(message string) error {
	return &ValidationError{message: message}
}

var Error = &ErrorUtils{}
