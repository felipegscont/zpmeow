package session

import (
	"fmt"

	"zpmeow/internal/domain/common"
)

// Status represents the current status of a session
type Status string

const (
	StatusDisconnected Status = "disconnected"
	StatusConnecting   Status = "connecting"
	StatusConnected    Status = "connected"
	StatusError        Status = "error"
)

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusDisconnected, StatusConnecting, StatusConnected, StatusError:
		return true
	default:
		return false
	}
}

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
}

// Session represents a WhatsApp session aggregate root
type Session struct {
	common.AggregateRoot

	// Core attributes
	id     SessionID
	name   SessionName
	status Status

	// WhatsApp specific
	deviceJID DeviceJID
	qrCode    QRCode

	// Configuration
	proxyConfig     ProxyConfiguration
	webhookEndpoint WebhookEndpoint
	apiKey          ApiKey

	// Timestamps
	createdAt common.Timestamp
	updatedAt common.Timestamp
}

// NewSession creates a new session aggregate
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
	waJID, _ := NewDeviceJID("")
	qrCode, _ := NewQRCode("")
	proxyConfig, _ := NewProxyConfiguration("")
	webhookEndpoint, _ := NewWebhookEndpoint("")
	apiKey, _ := NewApiKey("temp-key") // Temporary key, will be replaced by application layer

	now := common.Now()

	session := &Session{
		AggregateRoot:   common.NewAggregateRoot(sessionID.ID),
		id:              sessionID,
		name:            sessionName,
		status:          StatusDisconnected,
		deviceJID:       waJID,
		qrCode:          qrCode,
		proxyConfig:     proxyConfig,
		webhookEndpoint: webhookEndpoint,
		apiKey:          apiKey,
		createdAt:       now,
		updatedAt:       now,
	}

	// Publish domain event
	event := NewSessionCreatedEvent(sessionID.Value(), sessionName.Value())
	session.AddEvent(event)

	return session, nil
}

// Getters for aggregate attributes
func (s *Session) SessionID() SessionID {
	return s.id
}

func (s *Session) Name() SessionName {
	return s.name
}

func (s *Session) Status() Status {
	return s.status
}

func (s *Session) WaJID() DeviceJID {
	return s.deviceJID
}

func (s *Session) QRCode() QRCode {
	return s.qrCode
}

func (s *Session) ProxyConfiguration() ProxyConfiguration {
	return s.proxyConfig
}

func (s *Session) WebhookEndpoint() WebhookEndpoint {
	return s.webhookEndpoint
}

func (s *Session) ApiKey() ApiKey {
	return s.apiKey
}

func (s *Session) CreatedAt() common.Timestamp {
	return s.createdAt
}

func (s *Session) UpdatedAt() common.Timestamp {
	return s.updatedAt
}

// Business logic methods
func (s *Session) IsConnected() bool {
	return s.status == StatusConnected
}

func (s *Session) IsDisconnected() bool {
	return s.status == StatusDisconnected
}

func (s *Session) IsConnecting() bool {
	return s.status == StatusConnecting
}

func (s *Session) HasError() bool {
	return s.status == StatusError
}

func (s *Session) CanConnect() bool {
	return s.IsDisconnected() || s.HasError() || s.IsConnecting()
}

func (s *Session) HasQRCode() bool {
	return !s.qrCode.IsEmpty()
}

func (s *Session) HasProxy() bool {
	return !s.proxyConfig.IsEmpty()
}

func (s *Session) IsAuthenticated() bool {
	return !s.deviceJID.IsEmpty()
}

// Connect attempts to connect the session
func (s *Session) Connect() error {
	if !s.CanConnect() {
		return fmt.Errorf("session cannot connect from current status: %s", s.status)
	}

	oldStatus := s.status
	s.status = StatusConnecting
	s.updateTimestamp()

	// Publish event if status actually changed
	if oldStatus != StatusConnecting {
		event := NewSessionConnectedEvent(s.id.Value(), s.deviceJID.Value())
		s.AddEvent(event)
	}

	return nil
}

// Disconnect disconnects the session
func (s *Session) Disconnect(reason string) error {
	if s.IsDisconnected() {
		return nil // Already disconnected
	}

	s.status = StatusDisconnected
	s.updateTimestamp()

	// Clear QR code when disconnecting
	s.qrCode, _ = NewQRCode("")

	// Publish event
	event := NewSessionDisconnectedEvent(s.id.Value(), reason)
	s.AddEvent(event)

	return nil
}

// SetConnected marks the session as connected
func (s *Session) SetConnected() error {
	if !s.IsConnecting() {
		return fmt.Errorf("session must be connecting to be marked as connected")
	}

	s.status = StatusConnected
	s.updateTimestamp()

	// Publish event
	event := NewSessionConnectedEvent(s.id.Value(), s.deviceJID.Value())
	s.AddEvent(event)

	return nil
}

// SetError sets the session to error status
func (s *Session) SetError(errorMessage string) {
	s.status = StatusError
	s.updateTimestamp()

	// Publish event
	event := NewSessionErrorEvent(s.id.Value(), errorMessage, "connection_error")
	s.AddEvent(event)
}

// SetStatus sets the session status
func (s *Session) SetStatus(status Status) {
	s.status = status
	s.updateTimestamp()
}

// SetQRCode sets the QR code for the session
func (s *Session) SetQRCode(qrCode string) error {
	qr, err := NewQRCode(qrCode)
	if err != nil {
		return err
	}

	s.qrCode = qr
	s.updateTimestamp()

	// Publish configuration change event
	changes := map[string]interface{}{
		"qr_code_updated": true,
	}
	event := NewSessionConfigurationChangedEvent(s.id.Value(), changes)
	s.AddEvent(event)

	return nil
}

// Authenticate sets the WhatsApp JID and marks session as authenticated
func (s *Session) Authenticate(jid string) error {
	deviceJID, err := NewDeviceJID(jid)
	if err != nil {
		return err
	}

	s.deviceJID = deviceJID
	s.updateTimestamp()

	// Clear QR code after authentication
	s.qrCode, _ = NewQRCode("")

	// Publish authentication event
	event := NewSessionAuthenticatedEvent(s.id.Value(), jid)
	s.AddEvent(event)

	return nil
}

// SetProxyConfiguration configures the proxy
func (s *Session) SetProxyConfiguration(proxyConfig string) error {
	proxy, err := NewProxyConfiguration(proxyConfig)
	if err != nil {
		return err
	}

	s.proxyConfig = proxy
	s.updateTimestamp()

	// Publish configuration change event
	changes := map[string]interface{}{
		"proxy_configuration": proxyConfig,
	}
	event := NewSessionConfigurationChangedEvent(s.id.Value(), changes)
	s.AddEvent(event)

	return nil
}

// SetWebhookEndpoint configures the webhook endpoint
func (s *Session) SetWebhookEndpoint(webhookEndpoint string) error {
	webhook, err := NewWebhookEndpoint(webhookEndpoint)
	if err != nil {
		return err
	}

	s.webhookEndpoint = webhook
	s.updateTimestamp()

	// Publish configuration change event
	changes := map[string]interface{}{
		"webhook_endpoint": webhookEndpoint,
	}
	event := NewSessionConfigurationChangedEvent(s.id.Value(), changes)
	s.AddEvent(event)

	return nil
}

// SetApiKey sets the API key for the session
func (s *Session) SetApiKey(apiKey string) error {
	key, err := NewApiKey(apiKey)
	if err != nil {
		return err
	}

	s.apiKey = key
	s.updateTimestamp()

	return nil
}

// ClearQRCode clears the QR code
func (s *Session) ClearQRCode() {
	s.qrCode, _ = NewQRCode("")
	s.updateTimestamp()
}

// ClearProxy removes the proxy configuration
func (s *Session) ClearProxy() {
	s.proxyConfig, _ = NewProxyConfiguration("")
	s.updateTimestamp()

	// Publish configuration change event
	changes := map[string]interface{}{
		"proxy_cleared": true,
	}
	event := NewSessionConfigurationChangedEvent(s.id.Value(), changes)
	s.AddEvent(event)
}

// updateTimestamp updates the last modified timestamp
func (s *Session) updateTimestamp() {
	s.updatedAt = common.Now()
}

// Validate validates the session aggregate
func (s *Session) Validate() error {
	if s.id.IsEmpty() {
		return ErrInvalidSessionID
	}

	if s.name.IsEmpty() {
		return ErrInvalidSessionName
	}

	if !s.status.IsValid() {
		return ErrInvalidSessionStatus
	}

	return nil
}

// Delete marks the session for deletion
func (s *Session) Delete() {
	// Publish deletion event
	event := NewSessionDeletedEvent(s.id.Value(), s.name.Value())
	s.AddEvent(event)
}

// HasApiKey checks if session has an API key
func (s *Session) HasApiKey() bool {
	return !s.apiKey.IsEmpty()
}

// GetDeviceJIDString returns the WhatsApp JID as string
func (s *Session) GetDeviceJIDString() string {
	return s.deviceJID.Value()
}

// GetQRCodeString returns the QR code as string
func (s *Session) GetQRCodeString() string {
	return s.qrCode.Value()
}

// GetApiKeyString returns the API key as string
func (s *Session) GetApiKeyString() string {
	return s.apiKey.Value()
}

// HasWebhook checks if session has webhook configured
func (s *Session) HasWebhook() bool {
	return !s.webhookEndpoint.IsEmpty()
}

// GetWebhookEndpointString returns the webhook endpoint as string
func (s *Session) GetWebhookEndpointString() string {
	return s.webhookEndpoint.Value()
}
