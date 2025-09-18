package session

import (
	"fmt"
	"strings"

	"meow/internal/domain/common"
)

// SessionID represents a unique session identifier
type SessionID struct {
	common.ID
}

// NewSessionID creates a new session ID
func NewSessionID(value string) (SessionID, error) {
	id, err := common.NewID(value)
	if err != nil {
		return SessionID{}, fmt.Errorf("invalid session ID: %w", err)
	}
	return SessionID{ID: id}, nil
}

// GenerateSessionID creates a new random session ID
func GenerateSessionID() SessionID {
	return SessionID{ID: common.GenerateID()}
}

// SessionName represents a session name with domain-specific validation
type SessionName struct {
	common.Name
}

// NewSessionName creates a new session name with validation
func NewSessionName(value string) (SessionName, error) {
	name, err := common.NewName(value, 3, 100)
	if err != nil {
		return SessionName{}, fmt.Errorf("invalid session name: %w", err)
	}
	return SessionName{Name: name}, nil
}

// ApiKey represents an API key for session authentication
type ApiKey struct {
	value string
}

// NewApiKey creates a new API key with validation
func NewApiKey(value string) (ApiKey, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return ApiKey{}, fmt.Errorf("API key cannot be empty")
	}

	if len(trimmed) < 10 {
		return ApiKey{}, fmt.Errorf("API key too short (minimum 10 characters)")
	}

	if len(trimmed) > 100 {
		return ApiKey{}, fmt.Errorf("API key too long (maximum 100 characters)")
	}

	return ApiKey{value: trimmed}, nil
}

// Value returns the API key value
func (a ApiKey) Value() string {
	return a.value
}

// String returns a masked representation for logging
func (a ApiKey) String() string {
	if len(a.value) <= 8 {
		return "****"
	}
	return a.value[:4] + "****" + a.value[len(a.value)-4:]
}

// IsEmpty checks if API key is empty
func (a ApiKey) IsEmpty() bool {
	return a.value == ""
}

// WaJID represents a WhatsApp JID (Jabber ID)
type WaJID struct {
	value string
}

// NewWaJID creates a new WhatsApp JID with validation
func NewWaJID(value string) (WaJID, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return WaJID{}, nil // Empty JID is allowed for new sessions
	}

	if !strings.Contains(trimmed, "@") {
		return WaJID{}, fmt.Errorf("invalid JID format: must contain @")
	}

	parts := strings.Split(trimmed, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return WaJID{}, fmt.Errorf("invalid JID format: user@server required")
	}

	return WaJID{value: trimmed}, nil
}

// Value returns the JID value
func (w WaJID) Value() string {
	return w.value
}

// String returns the string representation
func (w WaJID) String() string {
	return w.value
}

// IsEmpty checks if JID is empty
func (w WaJID) IsEmpty() bool {
	return w.value == ""
}

// QRCode represents a QR code for session pairing
type QRCode struct {
	value string
}

// NewQRCode creates a new QR code
func NewQRCode(value string) (QRCode, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return QRCode{}, nil // Empty QR code is allowed
	}

	// Basic validation - QR codes should be reasonable length
	if len(trimmed) > 10000 {
		return QRCode{}, fmt.Errorf("QR code too long")
	}

	return QRCode{value: trimmed}, nil
}

// Value returns the QR code value
func (q QRCode) Value() string {
	return q.value
}

// String returns the string representation
func (q QRCode) String() string {
	return q.value
}

// IsEmpty checks if QR code is empty
func (q QRCode) IsEmpty() bool {
	return q.value == ""
}

// IsDataURL checks if QR code is a data URL
func (q QRCode) IsDataURL() bool {
	return strings.HasPrefix(q.value, "data:")
}

// ProxyConfiguration represents a proxy configuration for the session
type ProxyConfiguration struct {
	value string
}

// NewProxyConfiguration creates a new proxy configuration with domain validation
func NewProxyConfiguration(value string) (ProxyConfiguration, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return ProxyConfiguration{}, nil // Empty proxy is allowed (no proxy)
	}

	// Domain validation - basic format checks without external dependencies
	if len(trimmed) < 7 { // Minimum: "a://b:1"
		return ProxyConfiguration{}, fmt.Errorf("proxy configuration too short")
	}

	if len(trimmed) > 500 {
		return ProxyConfiguration{}, fmt.Errorf("proxy configuration too long")
	}

	// Basic scheme validation without net/url dependency
	if !strings.Contains(trimmed, "://") {
		return ProxyConfiguration{}, fmt.Errorf("proxy configuration must contain scheme (://)")
	}

	parts := strings.Split(trimmed, "://")
	if len(parts) != 2 {
		return ProxyConfiguration{}, fmt.Errorf("invalid proxy configuration format")
	}

	scheme := strings.ToLower(parts[0])
	if scheme != "http" && scheme != "https" && scheme != "socks5" {
		return ProxyConfiguration{}, fmt.Errorf("unsupported proxy scheme: %s (supported: http, https, socks5)", scheme)
	}

	if parts[1] == "" {
		return ProxyConfiguration{}, fmt.Errorf("proxy configuration must have host")
	}

	return ProxyConfiguration{value: trimmed}, nil
}

// Value returns the proxy configuration value
func (p ProxyConfiguration) Value() string {
	return p.value
}

// String returns the string representation
func (p ProxyConfiguration) String() string {
	return p.value
}

// IsEmpty checks if proxy configuration is empty
func (p ProxyConfiguration) IsEmpty() bool {
	return p.value == ""
}

// Scheme returns the proxy scheme (http, https, socks5)
func (p ProxyConfiguration) Scheme() string {
	if p.value == "" {
		return ""
	}
	parts := strings.Split(p.value, "://")
	if len(parts) < 2 {
		return ""
	}
	return strings.ToLower(parts[0])
}

// WebhookEndpoint represents a webhook endpoint configuration
type WebhookEndpoint struct {
	value string
}

// NewWebhookEndpoint creates a new webhook endpoint with domain validation
func NewWebhookEndpoint(value string) (WebhookEndpoint, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return WebhookEndpoint{}, nil // Empty webhook endpoint is allowed
	}

	// Domain validation without external dependencies
	if len(trimmed) < 10 { // Minimum: "http://a.b"
		return WebhookEndpoint{}, fmt.Errorf("webhook endpoint too short")
	}

	if len(trimmed) > 2000 {
		return WebhookEndpoint{}, fmt.Errorf("webhook endpoint too long")
	}

	// Basic scheme validation
	if !strings.Contains(trimmed, "://") {
		return WebhookEndpoint{}, fmt.Errorf("webhook endpoint must contain scheme (://)")
	}

	parts := strings.Split(trimmed, "://")
	if len(parts) != 2 {
		return WebhookEndpoint{}, fmt.Errorf("invalid webhook endpoint format")
	}

	scheme := strings.ToLower(parts[0])
	if scheme != "http" && scheme != "https" {
		return WebhookEndpoint{}, fmt.Errorf("webhook endpoint must use http or https scheme")
	}

	if parts[1] == "" {
		return WebhookEndpoint{}, fmt.Errorf("webhook endpoint must have host")
	}

	return WebhookEndpoint{value: trimmed}, nil
}

// Value returns the webhook endpoint value
func (w WebhookEndpoint) Value() string {
	return w.value
}

// String returns the string representation
func (w WebhookEndpoint) String() string {
	return w.value
}

// IsEmpty checks if webhook endpoint is empty
func (w WebhookEndpoint) IsEmpty() bool {
	return w.value == ""
}

// IsSecure checks if the webhook uses HTTPS
func (w WebhookEndpoint) IsSecure() bool {
	return strings.HasPrefix(strings.ToLower(w.value), "https://")
}
