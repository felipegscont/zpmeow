package errors

import (
	"net/http"

	"zpmeow/internal/domain/session"
)

// ErrorMapping represents the mapping of domain errors to HTTP responses
type ErrorMapping struct {
	StatusCode int
	Message    string
}

// DomainErrorMappings maps domain errors to HTTP status codes and messages
var DomainErrorMappings = map[error]ErrorMapping{
	// Validation errors (400 Bad Request)
	session.ErrInvalidSessionID:          {http.StatusBadRequest, "Invalid session ID"},
	session.ErrInvalidSessionName:        {http.StatusBadRequest, "Invalid session name"},
	session.ErrSessionNameTooShort:       {http.StatusBadRequest, "Session name too short"},
	session.ErrSessionNameTooLong:        {http.StatusBadRequest, "Session name too long"},
	session.ErrInvalidSessionNameChar:    {http.StatusBadRequest, "Session name contains invalid characters"},
	session.ErrInvalidSessionNameFormat:  {http.StatusBadRequest, "Invalid session name format"},
	session.ErrReservedSessionName:       {http.StatusBadRequest, "Session name is reserved"},
	session.ErrInvalidSessionStatus:      {http.StatusBadRequest, "Invalid session status"},
	
	// Conflict errors (409 Conflict)
	session.ErrSessionAlreadyExists:      {http.StatusConflict, "Session already exists"},
	session.ErrSessionAlreadyConnected:   {http.StatusConflict, "Session is already connected"},
	session.ErrSessionCannotConnect:      {http.StatusConflict, "Session cannot be connected in current state"},
	
	// Not found errors (404 Not Found)
	session.ErrSessionNotFound:           {http.StatusNotFound, "Session not found"},
}

// MapDomainError maps a domain error to HTTP status code and message
func MapDomainError(err error) (statusCode int, message string) {
	if mapping, exists := DomainErrorMappings[err]; exists {
		return mapping.StatusCode, mapping.Message
	}
	
	// Default to internal server error
	return http.StatusInternalServerError, "Internal server error"
}

// IsValidationError checks if an error is a validation error (400 status)
func IsValidationError(err error) bool {
	mapping, exists := DomainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusBadRequest
}

// IsConflictError checks if an error is a conflict error (409 status)
func IsConflictError(err error) bool {
	mapping, exists := DomainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusConflict
}

// IsNotFoundError checks if an error is a not found error (404 status)
func IsNotFoundError(err error) bool {
	mapping, exists := DomainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusNotFound
}
