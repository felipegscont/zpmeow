package session

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	sessionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type SessionID struct {
	value string
}

func NewSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, fmt.Errorf("session ID cannot be empty")
	}

	if len(value) < 1 || len(value) > 100 {
		return SessionID{}, fmt.Errorf("session ID must be between 1 and 100 characters")
	}

	return SessionID{value: value}, nil
}

func (s SessionID) Value() string {
	return s.value
}

func (s SessionID) String() string {
	return s.value
}

func (s SessionID) IsEmpty() bool {
	return s.value == ""
}

func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}

type SessionName struct {
	value string
}

func NewSessionName(value string) (SessionName, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return SessionName{}, fmt.Errorf("session name cannot be empty")
	}

	if len(trimmed) < 3 {
		return SessionName{}, fmt.Errorf("session name must be at least 3 characters long")
	}

	if len(trimmed) > 100 {
		return SessionName{}, fmt.Errorf("session name cannot exceed 100 characters")
	}

	if !sessionNameRegex.MatchString(trimmed) {
		return SessionName{}, fmt.Errorf("session name can only contain letters, numbers, hyphens, and underscores")
	}

	return SessionName{value: trimmed}, nil
}

func (s SessionName) Value() string {
	return s.value
}

func (s SessionName) String() string {
	return s.value
}

func (s SessionName) IsEmpty() bool {
	return s.value == ""
}

type ProxyURL struct {
	value string
}

func NewProxyURL(value string) (ProxyURL, error) {
	if value == "" {
		return ProxyURL{}, nil // Empty proxy URL is allowed (no proxy)
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return ProxyURL{}, fmt.Errorf("invalid proxy URL format: %w", err)
	}

	supportedSchemes := map[string]bool{
		"http":   true,
		"https":  true,
		"socks5": true,
	}

	if !supportedSchemes[parsedURL.Scheme] {
		return ProxyURL{}, fmt.Errorf("unsupported proxy scheme: %s (supported: http, https, socks5)", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return ProxyURL{}, fmt.Errorf("proxy URL must include a host")
	}

	return ProxyURL{value: value}, nil
}

func (p ProxyURL) Value() string {
	return p.value
}

func (p ProxyURL) String() string {
	return p.value
}

func (p ProxyURL) IsEmpty() bool {
	return p.value == ""
}

func (p ProxyURL) URL() (*url.URL, error) {
	if p.value == "" {
		return nil, nil
	}
	return url.Parse(p.value)
}

type SessionConfig struct {
	ID       SessionID
	Name     SessionName
	ProxyURL ProxyURL
}

func NewSessionConfig(id, name, proxyURL string) (SessionConfig, error) {
	sessionID, err := NewSessionID(id)
	if err != nil {
		return SessionConfig{}, fmt.Errorf("invalid session ID: %w", err)
	}

	sessionName, err := NewSessionName(name)
	if err != nil {
		return SessionConfig{}, fmt.Errorf("invalid session name: %w", err)
	}

	sessionProxyURL, err := NewProxyURL(proxyURL)
	if err != nil {
		return SessionConfig{}, fmt.Errorf("invalid proxy URL: %w", err)
	}

	return SessionConfig{
		ID:       sessionID,
		Name:     sessionName,
		ProxyURL: sessionProxyURL,
	}, nil
}

func (sc SessionConfig) Validate() error {
	if sc.ID.IsEmpty() {
		return fmt.Errorf("session ID is required")
	}

	if sc.Name.IsEmpty() {
		return fmt.Errorf("session name is required")
	}

	return nil
}

type WaJID struct {
	value string
}

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

func (w WaJID) Value() string {
	return w.value
}

func (w WaJID) String() string {
	return w.value
}

func (w WaJID) IsEmpty() bool {
	return w.value == ""
}

func (w WaJID) User() string {
	if w.value == "" {
		return ""
	}
	parts := strings.Split(w.value, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func (w WaJID) Server() string {
	if w.value == "" {
		return ""
	}
	parts := strings.Split(w.value, "@")
	if len(parts) > 1 {
		server := parts[1]
		if idx := strings.Index(server, "/"); idx != -1 {
			server = server[:idx]
		}
		return server
	}
	return ""
}

type QRCode struct {
	value string
}

func NewQRCode(value string) (QRCode, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return QRCode{}, nil // Empty QR code is allowed
	}

	if len(trimmed) < 10 {
		return QRCode{}, fmt.Errorf("QR code too short")
	}

	return QRCode{value: trimmed}, nil
}

func (q QRCode) Value() string {
	return q.value
}

func (q QRCode) String() string {
	return q.value
}

func (q QRCode) IsEmpty() bool {
	return q.value == ""
}

func (q QRCode) IsDataURL() bool {
	return strings.HasPrefix(q.value, "data:")
}

type ApiKey struct {
	value string
}

func NewApiKey(value string) (ApiKey, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return ApiKey{}, fmt.Errorf("API key cannot be empty")
	}

	if len(trimmed) < 10 {
		return ApiKey{}, fmt.Errorf("API key too short")
	}

	if len(trimmed) > 100 {
		return ApiKey{}, fmt.Errorf("API key too long")
	}

	return ApiKey{value: trimmed}, nil
}

func (a ApiKey) Value() string {
	return a.value
}

func (a ApiKey) String() string {
	if len(a.value) <= 8 {
		return "****"
	}
	return a.value[:4] + "****" + a.value[len(a.value)-4:]
}

func (a ApiKey) IsEmpty() bool {
	return a.value == ""
}

func (a ApiKey) Equals(other ApiKey) bool {
	return a.value == other.value
}
