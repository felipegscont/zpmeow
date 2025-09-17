package webhook

import "errors"

// Domain-specific errors for webhook aggregate
var (
	// Session ID errors
	ErrInvalidSessionID = errors.New("invalid session ID")
	ErrEmptySessionID   = errors.New("session ID cannot be empty")

	// Webhook URL errors
	ErrInvalidWebhookURL  = errors.New("invalid webhook URL")
	ErrWebhookURLNotHTTPS = errors.New("webhook URL must use HTTPS protocol")
	ErrWebhookURLNoHost   = errors.New("webhook URL must have a valid host")

	// Event type errors
	ErrInvalidEventType   = errors.New("invalid event type")
	ErrNoEventsSubscribed = errors.New("at least one event must be subscribed")
	ErrDuplicateEventType = errors.New("duplicate event type in subscription")

	// Webhook configuration errors
	ErrWebhookNotActive     = errors.New("webhook is not active")
	ErrWebhookNotConfigured = errors.New("webhook is not configured")
	ErrWebhookAlreadyExists = errors.New("webhook already exists for this session")

	// Business rule violations
	ErrCannotActivateWithoutURL    = errors.New("cannot activate webhook without URL")
	ErrCannotActivateWithoutEvents = errors.New("cannot activate webhook without events")

	// Note: Session-related errors should be handled by session domain
	// Authorization errors should be handled by application layer
)

// WebhookError represents a webhook-specific error with context
type WebhookError struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *WebhookError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *WebhookError) Unwrap() error {
	return e.Cause
}

// NewWebhookError creates a new webhook error
func NewWebhookError(code, message string, cause error) *WebhookError {
	return &WebhookError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Common webhook error constructors
func NewInvalidURLError(url string, cause error) *WebhookError {
	return NewWebhookError(
		"INVALID_WEBHOOK_URL",
		"invalid webhook URL: "+url,
		cause,
	)
}

func NewInvalidEventError(event string) *WebhookError {
	return NewWebhookError(
		"INVALID_EVENT_TYPE",
		"invalid event type: "+event,
		ErrInvalidEventType,
	)
}

func NewWebhookNotFoundError(sessionID string) *WebhookError {
	return NewWebhookError(
		"WEBHOOK_NOT_FOUND",
		"webhook not found for session: "+sessionID,
		ErrWebhookNotConfigured,
	)
}

func NewWebhookNotActiveError(sessionID string) *WebhookError {
	return NewWebhookError(
		"WEBHOOK_NOT_ACTIVE",
		"webhook is not active for session: "+sessionID,
		ErrWebhookNotActive,
	)
}

// Note: Session-related error constructors removed
// Session errors should be handled by session domain

// IsWebhookError checks if an error is a webhook error
func IsWebhookError(err error) bool {
	_, ok := err.(*WebhookError)
	return ok
}

// GetWebhookErrorCode extracts the error code from a webhook error
func GetWebhookErrorCode(err error) string {
	if webhookErr, ok := err.(*WebhookError); ok {
		return webhookErr.Code
	}
	return "UNKNOWN_ERROR"
}
