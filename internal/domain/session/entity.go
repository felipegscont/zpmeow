package session

import (
	"strings"
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
		Name:      strings.TrimSpace(name),
		Status:    types.StatusDisconnected,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Status query methods

// IsConnected checks if the session is connected
func (s *Session) IsConnected() bool {
	return s.Status == types.StatusConnected
}

// IsDisconnected checks if the session is disconnected
func (s *Session) IsDisconnected() bool {
	return s.Status == types.StatusDisconnected
}

// IsConnecting checks if the session is in connecting state
func (s *Session) IsConnecting() bool {
	return s.Status == types.StatusConnecting
}

// HasError checks if the session is in error state
func (s *Session) HasError() bool {
	return s.Status == types.StatusError
}

// CanConnect checks if the session can be connected
func (s *Session) CanConnect() bool {
	return s.IsDisconnected() || s.HasError() || s.IsConnecting()
}

// HasQRCode checks if the session has a QR code available
func (s *Session) HasQRCode() bool {
	return strings.TrimSpace(s.QRCode) != ""
}

// HasProxy checks if the session has a proxy configured
func (s *Session) HasProxy() bool {
	return strings.TrimSpace(s.ProxyURL) != ""
}

// IsAuthenticated checks if the session has a WhatsApp JID (is paired)
func (s *Session) IsAuthenticated() bool {
	return strings.TrimSpace(s.WhatsAppJID) != ""
}

// State modification methods

// SetStatus updates the session status and timestamp
func (s *Session) SetStatus(status types.Status) {
	if !status.IsValid() {
		return // Ignore invalid status changes
	}
	s.Status = status
	s.updateTimestamp()
}

// SetQRCode updates the QR code and timestamp
func (s *Session) SetQRCode(qrCode string) {
	s.QRCode = strings.TrimSpace(qrCode)
	s.updateTimestamp()
}

// SetWhatsAppJID updates the WhatsApp JID and timestamp
func (s *Session) SetWhatsAppJID(jid string) {
	s.WhatsAppJID = strings.TrimSpace(jid)
	s.updateTimestamp()
}

// SetProxyURL updates the proxy URL and timestamp
func (s *Session) SetProxyURL(proxyURL string) {
	s.ProxyURL = strings.TrimSpace(proxyURL)
	s.updateTimestamp()
}

// ClearQRCode removes the QR code
func (s *Session) ClearQRCode() {
	s.QRCode = ""
	s.updateTimestamp()
}

// ClearProxy removes the proxy configuration
func (s *Session) ClearProxy() {
	s.ProxyURL = ""
	s.updateTimestamp()
}

// updateTimestamp is a private helper to update the timestamp
func (s *Session) updateTimestamp() {
	s.UpdatedAt = time.Now()
}

// Validation methods

// Validate checks if the session is valid
func (s *Session) Validate() error {
	if err := s.validateID(); err != nil {
		return err
	}
	if err := s.validateName(); err != nil {
		return err
	}
	if err := s.validateStatus(); err != nil {
		return err
	}
	return nil
}

// validateID validates the session ID
func (s *Session) validateID() error {
	if strings.TrimSpace(s.ID) == "" {
		return ErrInvalidSessionID
	}
	return nil
}

// validateName validates the session name
func (s *Session) validateName() error {
	if strings.TrimSpace(s.Name) == "" {
		return ErrInvalidSessionName
	}
	return nil
}

// validateStatus validates the session status
func (s *Session) validateStatus() error {
	if !s.Status.IsValid() {
		return ErrInvalidSessionStatus
	}
	return nil
}

// Domain errors with improved messages
var (
	ErrInvalidSessionID          = NewDomainError("session ID cannot be empty")
	ErrInvalidSessionName        = NewDomainError("session name cannot be empty")
	ErrSessionNameTooShort       = NewDomainError("session name must be at least 3 characters long")
	ErrSessionNameTooLong        = NewDomainError("session name cannot exceed 50 characters")
	ErrInvalidSessionNameChar    = NewDomainError("session name can only contain letters, numbers, hyphens, and underscores")
	ErrInvalidSessionNameFormat  = NewDomainError("session name cannot start or end with hyphen or underscore")
	ErrReservedSessionName       = NewDomainError("session name is reserved and cannot be used")
	ErrInvalidSessionStatus      = NewDomainError("session status is invalid")
	ErrSessionNotFound           = NewDomainError("session not found")
	ErrSessionAlreadyExists      = NewDomainError("session already exists")
	ErrSessionAlreadyConnected   = NewDomainError("session is already connected")
	ErrSessionCannotConnect      = NewDomainError("session cannot be connected in current state")
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
