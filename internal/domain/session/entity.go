package session

import (
	"time"
	"zpmeow/internal/types"
)

// Session represents a WhatsApp session entity (pure business logic)
type Session struct {
	ID          string
	Name        string
	WhatsAppJID string
	Status      types.Status
	QRCode      string
	ProxyURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewSession creates a new session with default values
func NewSession(id, name string) *Session {
	now := time.Now()
	return &Session{
		ID:        id,
		Name:      name,
		Status:    types.StatusDisconnected,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsConnected checks if the session is connected
func (s *Session) IsConnected() bool {
	return s.Status == types.StatusConnected
}

// IsDisconnected checks if the session is disconnected
func (s *Session) IsDisconnected() bool {
	return s.Status == types.StatusDisconnected
}

// CanConnect checks if the session can be connected
func (s *Session) CanConnect() bool {
	// Allow connection if disconnected, error, or if it's been connecting for too long
	return s.Status == types.StatusDisconnected ||
		   s.Status == types.StatusError ||
		   s.Status == types.StatusConnecting // Allow reconnection attempts
}

// SetStatus updates the session status and timestamp
func (s *Session) SetStatus(status types.Status) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

// SetQRCode updates the QR code and timestamp
func (s *Session) SetQRCode(qrCode string) {
	s.QRCode = qrCode
	s.UpdatedAt = time.Now()
}

// SetWhatsAppJID updates the WhatsApp JID and timestamp
func (s *Session) SetWhatsAppJID(jid string) {
	s.WhatsAppJID = jid
	s.UpdatedAt = time.Now()
}

// SetProxyURL updates the proxy URL and timestamp
func (s *Session) SetProxyURL(proxyURL string) {
	s.ProxyURL = proxyURL
	s.UpdatedAt = time.Now()
}

// ClearQRCode removes the QR code
func (s *Session) ClearQRCode() {
	s.QRCode = ""
	s.UpdatedAt = time.Now()
}

// Validate checks if the session is valid
func (s *Session) Validate() error {
	if s.ID == "" {
		return ErrInvalidSessionID
	}
	if s.Name == "" {
		return ErrInvalidSessionName
	}
	if !s.Status.IsValid() {
		return ErrInvalidSessionStatus
	}
	return nil
}

// Domain errors
var (
	ErrInvalidSessionID     = NewDomainError("invalid session ID")
	ErrInvalidSessionName   = NewDomainError("invalid session name")
	ErrInvalidSessionStatus = NewDomainError("invalid session status")
	ErrSessionNotFound      = NewDomainError("session not found")
	ErrSessionAlreadyExists = NewDomainError("session already exists")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Message string
}

func (e DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error
func NewDomainError(message string) DomainError {
	return DomainError{Message: message}
}
