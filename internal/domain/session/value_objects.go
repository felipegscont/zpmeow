package session

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
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

	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(trimmed) {
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
