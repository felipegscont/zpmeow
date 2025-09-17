package session

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDisconnected Status = "disconnected"
	StatusConnecting   Status = "connecting"
	StatusConnected    Status = "connected"
	StatusError        Status = "error"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusDisconnected, StatusConnecting, StatusConnected, StatusError:
		return true
	default:
		return false
	}
}

type Session struct {
	ID         SessionID
	Name       SessionName
	WaJID      string
	Status     Status
	QRCode     string
	ProxyURL   ProxyURL
	WebhookURL string
	Events     []string
	ApiKey     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewSession(id, name string) (*Session, error) {
	sessionID, err := NewSessionID(id)
	if err != nil {
		return nil, err
	}

	sessionName, err := NewSessionName(name)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Session{
		ID:        sessionID,
		Name:      sessionName,
		Status:    StatusDisconnected,
		ProxyURL:  ProxyURL{}, // Empty proxy URL
		ApiKey:    uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *Session) IsConnected() bool {
	return s.Status == StatusConnected
}

func (s *Session) IsDisconnected() bool {
	return s.Status == StatusDisconnected
}

func (s *Session) IsConnecting() bool {
	return s.Status == StatusConnecting
}

func (s *Session) HasError() bool {
	return s.Status == StatusError
}

func (s *Session) CanConnect() bool {
	return s.IsDisconnected() || s.HasError() || s.IsConnecting()
}

func (s *Session) HasQRCode() bool {
	return strings.TrimSpace(s.QRCode) != ""
}

func (s *Session) HasProxy() bool {
	return !s.ProxyURL.IsEmpty()
}

func (s *Session) IsAuthenticated() bool {
	return strings.TrimSpace(s.WaJID) != ""
}

func (s *Session) SetStatus(status Status) {
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

func (s *Session) SetWaJID(jid string) {
	s.WaJID = strings.TrimSpace(jid)
	s.updateTimestamp()
}

func (s *Session) SetProxyURL(proxyURL string) error {
	proxy, err := NewProxyURL(proxyURL)
	if err != nil {
		return err
	}
	s.ProxyURL = proxy
	s.updateTimestamp()
	return nil
}

func (s *Session) ClearQRCode() {
	s.QRCode = ""
	s.updateTimestamp()
}

func (s *Session) ClearProxy() {
	s.ProxyURL = ProxyURL{} // Empty proxy URL
	s.updateTimestamp()
}

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
	if s.ID.IsEmpty() {
		return ErrInvalidSessionID
	}
	return nil
}

func (s *Session) validateName() error {
	if s.Name.IsEmpty() {
		return ErrInvalidSessionName
	}
	return nil
}

func (s *Session) validateStatus() error {
	if !s.Status.IsValid() {
		return ErrInvalidSessionStatus
	}
	return nil
}

func (s *Session) HasApiKey() bool {
	return strings.TrimSpace(s.ApiKey) != ""
}

func (s *Session) RegenerateApiKey() {
	s.ApiKey = generateApiKey()
	s.updateTimestamp()
}

func (s *Session) SetApiKey(apiKey string) {
	s.ApiKey = strings.TrimSpace(apiKey)
	s.updateTimestamp()
}

func generateApiKey() string {
	return uuid.New().String()
}
