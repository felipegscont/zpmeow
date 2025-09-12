package session

import (
	"strings"
	"time"
	"zpmeow/internal/shared/types"
)


type Session struct {
	ID          string
	Name        string
	WhatsAppJID string
	Status      types.Status
	QRCode      string
	ProxyURL    string
	WebhookURL  string
	Events      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}


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




func (s *Session) IsConnected() bool {
	return s.Status == types.StatusConnected
}


func (s *Session) IsDisconnected() bool {
	return s.Status == types.StatusDisconnected
}


func (s *Session) IsConnecting() bool {
	return s.Status == types.StatusConnecting
}


func (s *Session) HasError() bool {
	return s.Status == types.StatusError
}


func (s *Session) CanConnect() bool {
	return s.IsDisconnected() || s.HasError() || s.IsConnecting()
}


func (s *Session) HasQRCode() bool {
	return strings.TrimSpace(s.QRCode) != ""
}


func (s *Session) HasProxy() bool {
	return strings.TrimSpace(s.ProxyURL) != ""
}


func (s *Session) IsAuthenticated() bool {
	return strings.TrimSpace(s.WhatsAppJID) != ""
}




func (s *Session) SetStatus(status types.Status) {
	if !status.IsValid() {
		return // Ignore invalid status changes
	}
	s.Status = status
	s.updateTimestamp()
}


func (s *Session) SetQRCode(qrCode string) {
	s.QRCode = strings.TrimSpace(qrCode)
	s.updateTimestamp()
}


func (s *Session) SetWhatsAppJID(jid string) {
	s.WhatsAppJID = strings.TrimSpace(jid)
	s.updateTimestamp()
}


func (s *Session) SetProxyURL(proxyURL string) {
	s.ProxyURL = strings.TrimSpace(proxyURL)
	s.updateTimestamp()
}


func (s *Session) ClearQRCode() {
	s.QRCode = ""
	s.updateTimestamp()
}


func (s *Session) ClearProxy() {
	s.ProxyURL = ""
	s.updateTimestamp()
}

// Webhook methods
func (s *Session) SetWebhook(url string, events []string) {
	s.WebhookURL = strings.TrimSpace(url)
	s.Events = events
	s.updateTimestamp()
}

func (s *Session) ClearWebhook() {
	s.WebhookURL = ""
	s.Events = nil
	s.updateTimestamp()
}

func (s *Session) HasWebhook() bool {
	return strings.TrimSpace(s.WebhookURL) != ""
}

func (s *Session) IsEventSubscribed(event string) bool {
	if !s.HasWebhook() {
		return false
	}
	for _, e := range s.Events {
		if e == event {
			return true
		}
	}
	return false
}


func (s *Session) updateTimestamp() {
	s.UpdatedAt = time.Now()
}




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


func (s *Session) validateID() error {
	return ValidateSessionID(s.ID)
}


func (s *Session) validateName() error {
	return ValidateSessionName(s.Name)
}


func (s *Session) validateStatus() error {
	if !s.Status.IsValid() {
		return ErrInvalidSessionStatus
	}
	return nil
}



