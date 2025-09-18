package webhook

import (
	"fmt"
	"net/url"
	"strings"
)

// SessionID represents a session identifier for webhook
type SessionID struct {
	value string
}

// NewSessionID creates a new SessionID
func NewSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, fmt.Errorf("session ID cannot be empty")
	}

	if len(value) < 1 || len(value) > 100 {
		return SessionID{}, fmt.Errorf("session ID must be between 1 and 100 characters")
	}

	return SessionID{value: value}, nil
}

// Value returns the session ID value
func (s SessionID) Value() string {
	return s.value
}

// String returns the string representation
func (s SessionID) String() string {
	return s.value
}

// IsEmpty checks if the session ID is empty
func (s SessionID) IsEmpty() bool {
	return s.value == ""
}

type WebhookURL struct {
	value string
}

func NewWebhookURL(value string) (WebhookURL, error) {
	if value == "" {
		return WebhookURL{}, nil // Empty URL is allowed for inactive webhooks
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return WebhookURL{}, fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return WebhookURL{}, fmt.Errorf("webhook URL must use HTTPS protocol")
	}

	if parsedURL.Host == "" {
		return WebhookURL{}, fmt.Errorf("webhook URL must have a valid host")
	}

	return WebhookURL{value: value}, nil
}

func (w WebhookURL) Value() string {
	return w.value
}

func (w WebhookURL) String() string {
	return w.value
}

func (w WebhookURL) IsEmpty() bool {
	return w.value == ""
}

func (w WebhookURL) URL() (*url.URL, error) {
	if w.value == "" {
		return nil, nil
	}
	return url.Parse(w.value)
}

func (w WebhookURL) Host() string {
	if w.value == "" {
		return ""
	}

	parsedURL, err := url.Parse(w.value)
	if err != nil {
		return ""
	}

	return parsedURL.Host
}

func (w WebhookURL) Path() string {
	if w.value == "" {
		return ""
	}

	parsedURL, err := url.Parse(w.value)
	if err != nil {
		return ""
	}

	return parsedURL.Path
}

// EventType represents the type of event that can trigger a webhook
type EventType string

const (
	EventTypeMessage      EventType = "message"
	EventTypeConnected    EventType = "connected"
	EventTypeDisconnected EventType = "disconnected"
	EventTypeAll          EventType = "all"
)

// String returns the string representation of the event type
func (e EventType) String() string {
	return string(e)
}

// IsValid checks if the event type is valid
func (e EventType) IsValid() bool {
	switch e {
	case EventTypeMessage, EventTypeConnected, EventTypeDisconnected, EventTypeAll:
		return true
	default:
		return false
	}
}

func ParseEventType(s string) (EventType, error) {
	eventType := EventType(s)
	if !eventType.IsValid() {
		return "", fmt.Errorf("invalid event type: %s", s)
	}
	return eventType, nil
}

func ParseEventTypes(events []string) ([]EventType, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("at least one event type is required")
	}

	eventTypes := make([]EventType, 0, len(events))
	seen := make(map[EventType]bool)

	for _, event := range events {
		trimmed := strings.TrimSpace(event)
		if len(trimmed) == 0 {
			continue
		}
		normalizedEvent := strings.ToUpper(trimmed[:1]) + strings.ToLower(trimmed[1:])

		eventType, err := ParseEventType(normalizedEvent)
		if err != nil {
			return nil, err
		}

		if !seen[eventType] {
			eventTypes = append(eventTypes, eventType)
			seen[eventType] = true
		}
	}

	return eventTypes, nil
}

func GetAllEventTypes() []EventType {
	return []EventType{
		EventTypeMessage,
		EventTypeConnected,
		EventTypeDisconnected,
		EventTypeAll,
	}
}

func GetEventTypeNames() []string {
	return []string{
		string(EventTypeMessage),
		string(EventTypeConnected),
		string(EventTypeDisconnected),
		string(EventTypeAll),
	}
}
