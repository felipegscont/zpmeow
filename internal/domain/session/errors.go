package session

import "errors"

// Domain errors - business rule violations
var (
	// Session validation errors
	ErrInvalidSessionID         = errors.New("invalid session ID")
	ErrInvalidSessionName       = errors.New("invalid session name")
	ErrSessionNameTooShort      = errors.New("session name too short")
	ErrSessionNameTooLong       = errors.New("session name too long")
	ErrInvalidSessionNameChar   = errors.New("session name contains invalid characters")
	ErrInvalidSessionNameFormat = errors.New("invalid session name format")
	ErrReservedSessionName      = errors.New("session name is reserved")
	ErrInvalidSessionStatus     = errors.New("invalid session status")
	ErrInvalidProxyURL          = errors.New("invalid proxy URL")

	// Session business rule errors
	ErrSessionAlreadyExists    = errors.New("session already exists")
	ErrSessionNotFound         = errors.New("session not found")
	ErrSessionAlreadyConnected = errors.New("session is already connected")
	ErrSessionCannotConnect    = errors.New("session cannot be connected in current state")
	ErrSessionCannotDisconnect = errors.New("session cannot be disconnected in current state")
	ErrSessionCannotDelete     = errors.New("session cannot be deleted in current state")

	// WhatsApp service errors
	ErrWhatsAppServiceUnavailable = errors.New("WhatsApp service unavailable")
	ErrQRCodeGenerationFailed     = errors.New("QR code generation failed")
	ErrPairingFailed              = errors.New("phone pairing failed")
	ErrConnectionFailed           = errors.New("connection failed")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewDomainError creates a new domain error
func NewDomainError(message string) *DomainError {
	return &DomainError{
		Message: message,
	}
}

// NewDomainErrorWithCause creates a new domain error with a cause
func NewDomainErrorWithCause(message string, cause error) *DomainError {
	return &DomainError{
		Message: message,
		Cause:   cause,
	}
}
