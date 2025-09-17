package session

import (
	"time"
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
	ID        SessionID
	Name      SessionName
	WaJID     WaJID
	Status    Status
	QRCode    QRCode
	ProxyURL  ProxyURL
	ApiKey    ApiKey
	CreatedAt time.Time
	UpdatedAt time.Time
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

	// Create empty value objects for new session
	waJID, _ := NewWaJID("")
	qrCode, _ := NewQRCode("")
	apiKey, _ := NewApiKey("temp-key") // Temporary key, will be replaced by application layer

	now := time.Now()
	return &Session{
		ID:        sessionID,
		Name:      sessionName,
		WaJID:     waJID,
		Status:    StatusDisconnected,
		QRCode:    qrCode,
		ProxyURL:  ProxyURL{}, // Empty proxy URL
		ApiKey:    apiKey,
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
	return !s.QRCode.IsEmpty()
}

func (s *Session) HasProxy() bool {
	return !s.ProxyURL.IsEmpty()
}

func (s *Session) IsAuthenticated() bool {
	return !s.WaJID.IsEmpty()
}

func (s *Session) SetStatus(status Status) {
	if !status.IsValid() {
		return // Ignore invalid status changes
	}
	s.Status = status
	s.updateTimestamp()
}

func (s *Session) SetQRCode(qrCode string) error {
	qr, err := NewQRCode(qrCode)
	if err != nil {
		return err
	}
	s.QRCode = qr
	s.updateTimestamp()
	return nil
}

func (s *Session) SetWaJID(jid string) error {
	waJID, err := NewWaJID(jid)
	if err != nil {
		return err
	}
	s.WaJID = waJID
	s.updateTimestamp()
	return nil
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
	s.QRCode = QRCode{}
	s.updateTimestamp()
}

func (s *Session) ClearProxy() {
	s.ProxyURL = ProxyURL{} // Empty proxy URL
	s.updateTimestamp()
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
	return !s.ApiKey.IsEmpty()
}

// RegenerateApiKey should be called from application layer with generated key
func (s *Session) RegenerateApiKey(newApiKey string) error {
	apiKey, err := NewApiKey(newApiKey)
	if err != nil {
		return err
	}
	s.ApiKey = apiKey
	s.updateTimestamp()
	return nil
}

func (s *Session) SetApiKey(apiKey string) error {
	key, err := NewApiKey(apiKey)
	if err != nil {
		return err
	}
	s.ApiKey = key
	s.updateTimestamp()
	return nil
}

// generateApiKey moved to application layer to avoid external dependencies in domain

// Compatibility methods for gradual migration (will be removed in future versions)

// GetWaJIDString returns WaJID as string for compatibility
func (s *Session) GetWaJIDString() string {
	return s.WaJID.Value()
}

// GetQRCodeString returns QRCode as string for compatibility
func (s *Session) GetQRCodeString() string {
	return s.QRCode.Value()
}

// GetApiKeyString returns ApiKey as string for compatibility
func (s *Session) GetApiKeyString() string {
	return s.ApiKey.Value()
}

// SetWaJIDString sets WaJID from string for compatibility
func (s *Session) SetWaJIDString(jid string) error {
	return s.SetWaJID(jid)
}

// SetQRCodeString sets QRCode from string for compatibility
func (s *Session) SetQRCodeString(qr string) error {
	return s.SetQRCode(qr)
}

// SetApiKeyString sets ApiKey from string for compatibility
func (s *Session) SetApiKeyString(key string) error {
	return s.SetApiKey(key)
}
