package webhook

import (
	"fmt"
	"net/url"
	"strings"

	"meow/internal/shared/events"
	"meow/internal/shared/types"
)

type SessionID = types.SessionID

var NewSessionID = types.NewSessionID

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

type EventType = events.EventType

const (
	EventTypeMessage      = events.EventTypeMessage
	EventTypeConnected    = events.EventTypeConnected
	EventTypeDisconnected = events.EventTypeDisconnected
	EventTypeAll          = events.EventTypeAll
)

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
	return events.GetAllEventTypes()
}

func GetEventTypeNames() []string {
	return events.GetEventTypeNames()
}
