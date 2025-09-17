package webhook

import (
	"fmt"
	"net/url"
	"strings"

	"meow/internal/shared/events"
	"meow/internal/shared/types"
)

// SessionID is imported from shared types
type SessionID = types.SessionID

// NewSessionID creates a new SessionID (re-exported for convenience)
var NewSessionID = types.NewSessionID

// WebhookURL represents a webhook URL value object
type WebhookURL struct {
	value string
}

// NewWebhookURL creates a new WebhookURL
func NewWebhookURL(value string) (WebhookURL, error) {
	if value == "" {
		return WebhookURL{}, nil // Empty URL is allowed for inactive webhooks
	}

	// Validate URL format
	parsedURL, err := url.Parse(value)
	if err != nil {
		return WebhookURL{}, fmt.Errorf("invalid URL format: %w", err)
	}

	// Ensure HTTPS for security
	if parsedURL.Scheme != "https" {
		return WebhookURL{}, fmt.Errorf("webhook URL must use HTTPS protocol")
	}

	// Ensure host is present
	if parsedURL.Host == "" {
		return WebhookURL{}, fmt.Errorf("webhook URL must have a valid host")
	}

	return WebhookURL{value: value}, nil
}

// Value returns the webhook URL value
func (w WebhookURL) Value() string {
	return w.value
}

// String returns the string representation
func (w WebhookURL) String() string {
	return w.value
}

// IsEmpty checks if the webhook URL is empty
func (w WebhookURL) IsEmpty() bool {
	return w.value == ""
}

// URL returns the parsed URL
func (w WebhookURL) URL() (*url.URL, error) {
	if w.value == "" {
		return nil, nil
	}
	return url.Parse(w.value)
}

// Host returns the host part of the URL
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

// Path returns the path part of the URL
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

// EventType is an alias to the shared event type
type EventType = events.EventType

// Re-export commonly used event types for convenience
const (
	EventTypeMessage      = events.EventTypeMessage
	EventTypeConnected    = events.EventTypeConnected
	EventTypeDisconnected = events.EventTypeDisconnected
	EventTypeAll          = events.EventTypeAll
)

// ParseEventType parses a string into an EventType
func ParseEventType(s string) (EventType, error) {
	eventType := EventType(s)
	if !eventType.IsValid() {
		return "", fmt.Errorf("invalid event type: %s", s)
	}
	return eventType, nil
}

// ParseEventTypes parses a slice of strings into EventTypes
func ParseEventTypes(events []string) ([]EventType, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("at least one event type is required")
	}

	eventTypes := make([]EventType, 0, len(events))
	seen := make(map[EventType]bool)

	for _, event := range events {
		// Normalize event name (handle case variations)
		trimmed := strings.TrimSpace(event)
		if len(trimmed) == 0 {
			continue
		}
		// Convert to title case manually (strings.Title is deprecated)
		normalizedEvent := strings.ToUpper(trimmed[:1]) + strings.ToLower(trimmed[1:])

		eventType, err := ParseEventType(normalizedEvent)
		if err != nil {
			return nil, err
		}

		// Avoid duplicates
		if !seen[eventType] {
			eventTypes = append(eventTypes, eventType)
			seen[eventType] = true
		}
	}

	return eventTypes, nil
}

// GetAllEventTypes returns all valid event types
func GetAllEventTypes() []EventType {
	return events.GetAllEventTypes()
}

// GetEventTypeNames returns all valid event type names as strings
func GetEventTypeNames() []string {
	return events.GetEventTypeNames()
}
