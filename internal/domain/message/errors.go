package message

import "fmt"

type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeBusiness   ErrorType = "business"
	ErrorTypeNotFound   ErrorType = "not_found"
)

type DomainError struct {
	Type    ErrorType
	Message string
	Code    string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func NewValidationError(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    "MSG_VALIDATION_ERROR",
	}
}

func NewBusinessError(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeBusiness,
		Message: message,
		Code:    "MSG_BUSINESS_ERROR",
	}
}

func NewNotFoundError(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    "MSG_NOT_FOUND",
	}
}

var (
	ErrEmptyContent     = NewValidationError("message content cannot be empty")
	ErrInvalidJID       = NewValidationError("invalid JID format")
	ErrInvalidType      = NewValidationError("invalid message type")
	ErrMissingSessionID = NewValidationError("session ID is required")
	ErrMissingChatJID   = NewValidationError("chat JID is required")
	ErrMissingMediaURL  = NewValidationError("media URL is required for media messages")
	ErrUnsupportedType  = NewBusinessError("unsupported message type")
)

func IsValidationError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeValidation
	}
	return false
}

func IsBusinessError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeBusiness
	}
	return false
}

func IsNotFoundError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeNotFound
	}
	return false
}
