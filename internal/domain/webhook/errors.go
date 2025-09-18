package webhook

import "errors"

var (
	ErrInvalidSessionID = errors.New("invalid session ID")
	ErrEmptySessionID   = errors.New("session ID cannot be empty")

	ErrInvalidWebhookURL  = errors.New("invalid webhook URL")
	ErrWebhookURLNotHTTPS = errors.New("webhook URL must use HTTPS protocol")
	ErrWebhookURLNoHost   = errors.New("webhook URL must have a valid host")

	ErrInvalidEventType   = errors.New("invalid event type")
	ErrNoEventsSubscribed = errors.New("at least one event must be subscribed")
	ErrDuplicateEventType = errors.New("duplicate event type in subscription")

	ErrWebhookNotActive     = errors.New("webhook is not active")
	ErrWebhookNotConfigured = errors.New("webhook is not configured")
	ErrWebhookAlreadyExists = errors.New("webhook already exists for this session")

	ErrCannotActivateWithoutURL    = errors.New("cannot activate webhook without URL")
	ErrCannotActivateWithoutEvents = errors.New("cannot activate webhook without events")
)

type WebhookError struct {
	Code    string
	Message string
	Cause   error
}

func (e *WebhookError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *WebhookError) Unwrap() error {
	return e.Cause
}

func NewWebhookError(code, message string, cause error) *WebhookError {
	return &WebhookError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

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

func IsWebhookError(err error) bool {
	_, ok := err.(*WebhookError)
	return ok
}

func GetWebhookErrorCode(err error) string {
	if webhookErr, ok := err.(*WebhookError); ok {
		return webhookErr.Code
	}
	return "UNKNOWN_ERROR"
}
