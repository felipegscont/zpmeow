package common

import (
	"errors"
	"fmt"
)

// Application layer specific errors
var (
	// Validation errors
	ErrInvalidInput     = errors.New("invalid input")
	ErrValidationFailed = errors.New("validation failed")
	ErrMissingRequired  = errors.New("missing required field")

	// Business rule errors
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrSessionNotConnected  = errors.New("session not connected")
	ErrSessionInUse         = errors.New("session is in use")

	// Operation errors
	ErrOperationFailed        = errors.New("operation failed")
	ErrConcurrentModification = errors.New("concurrent modification detected")
	ErrResourceUnavailable    = errors.New("resource unavailable")
)

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// BusinessRuleError represents a business rule violation
type BusinessRuleError struct {
	Rule    string
	Message string
}

func (e BusinessRuleError) Error() string {
	return fmt.Sprintf("business rule violation '%s': %s", e.Rule, e.Message)
}

// NewBusinessRuleError creates a new business rule error
func NewBusinessRuleError(rule, message string) *BusinessRuleError {
	return &BusinessRuleError{
		Rule:    rule,
		Message: message,
	}
}

// ApplicationError represents a general application error
type ApplicationError struct {
	Code    string
	Message string
	Cause   error
}

func (e ApplicationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e ApplicationError) Unwrap() error {
	return e.Cause
}

// NewApplicationError creates a new application error
func NewApplicationError(code, message string, cause error) *ApplicationError {
	return &ApplicationError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsBusinessRuleError checks if an error is a business rule error
func IsBusinessRuleError(err error) bool {
	var businessErr *BusinessRuleError
	return errors.As(err, &businessErr)
}

// IsApplicationError checks if an error is an application error
func IsApplicationError(err error) bool {
	var appErr *ApplicationError
	return errors.As(err, &appErr)
}
