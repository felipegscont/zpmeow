package domain

// Re-export session domain types to avoid aliases
import (
	"zpmeow/internal/domain/session"
)

// Session domain types
type Session = session.Session
type SessionRepository = session.SessionRepository
type SessionService = session.SessionService
type WhatsAppService = session.WhatsAppService
type DomainError = session.DomainError

// Session domain errors
var (
	ErrInvalidSessionID          = session.ErrInvalidSessionID
	ErrInvalidSessionName        = session.ErrInvalidSessionName
	ErrSessionNameTooShort       = session.ErrSessionNameTooShort
	ErrSessionNameTooLong        = session.ErrSessionNameTooLong
	ErrInvalidSessionNameChar    = session.ErrInvalidSessionNameChar
	ErrInvalidSessionNameFormat  = session.ErrInvalidSessionNameFormat
	ErrReservedSessionName       = session.ErrReservedSessionName
	ErrInvalidSessionStatus      = session.ErrInvalidSessionStatus
	ErrInvalidProxyURL           = session.ErrInvalidProxyURL
	ErrSessionAlreadyExists      = session.ErrSessionAlreadyExists
	ErrSessionNotFound           = session.ErrSessionNotFound
	ErrSessionAlreadyConnected   = session.ErrSessionAlreadyConnected
	ErrSessionCannotConnect      = session.ErrSessionCannotConnect
	ErrSessionCannotDisconnect   = session.ErrSessionCannotDisconnect
	ErrSessionCannotDelete       = session.ErrSessionCannotDelete
	ErrWhatsAppServiceUnavailable = session.ErrWhatsAppServiceUnavailable
	ErrQRCodeGenerationFailed     = session.ErrQRCodeGenerationFailed
	ErrPairingFailed              = session.ErrPairingFailed
	ErrConnectionFailed           = session.ErrConnectionFailed
)

// Session domain functions
var (
	ValidateSessionName = session.ValidateSessionName
	ValidateSessionID   = session.ValidateSessionID
	ValidateProxyURL    = session.ValidateProxyURL
	ValidateSessionStatus = session.ValidateSessionStatus
	NewDomainError      = session.NewDomainError
	NewDomainErrorWithCause = session.NewDomainErrorWithCause
)
