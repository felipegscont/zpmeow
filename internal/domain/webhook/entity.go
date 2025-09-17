package webhook

import (
	"time"
)

// WebhookConfiguration represents webhook configuration as a separate aggregate
type WebhookConfiguration struct {
	SessionID SessionID
	URL       WebhookURL
	Events    []EventType
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewWebhookConfiguration creates a new webhook configuration
func NewWebhookConfiguration(sessionID, url string, events []string) (*WebhookConfiguration, error) {
	sid, err := NewSessionID(sessionID)
	if err != nil {
		return nil, err
	}

	webhookURL, err := NewWebhookURL(url)
	if err != nil {
		return nil, err
	}

	eventTypes, err := ParseEventTypes(events)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &WebhookConfiguration{
		SessionID: sid,
		URL:       webhookURL,
		Events:    eventTypes,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// IsActive returns whether the webhook is active
func (w *WebhookConfiguration) IsActive() bool {
	return w.Active && !w.URL.IsEmpty()
}

// IsEventSubscribed checks if the webhook is subscribed to a specific event
func (w *WebhookConfiguration) IsEventSubscribed(event EventType) bool {
	if !w.IsActive() {
		return false
	}

	// Check for "All" subscription
	for _, e := range w.Events {
		if e == EventTypeAll || e == event {
			return true
		}
	}

	return false
}

// UpdateURL updates the webhook URL
func (w *WebhookConfiguration) UpdateURL(url string) error {
	webhookURL, err := NewWebhookURL(url)
	if err != nil {
		return err
	}

	w.URL = webhookURL
	w.updateTimestamp()
	return nil
}

// UpdateEvents updates the subscribed events
func (w *WebhookConfiguration) UpdateEvents(events []string) error {
	eventTypes, err := ParseEventTypes(events)
	if err != nil {
		return err
	}

	w.Events = eventTypes
	w.updateTimestamp()
	return nil
}

// Activate activates the webhook
func (w *WebhookConfiguration) Activate() {
	w.Active = true
	w.updateTimestamp()
}

// Deactivate deactivates the webhook
func (w *WebhookConfiguration) Deactivate() {
	w.Active = false
	w.updateTimestamp()
}

// Note: Getter methods removed to maintain encapsulation
// Access to internal state should be through business methods only

// Validate validates the webhook configuration
func (w *WebhookConfiguration) Validate() error {
	if w.SessionID.IsEmpty() {
		return ErrInvalidSessionID
	}

	if w.Active && w.URL.IsEmpty() {
		return ErrInvalidWebhookURL
	}

	if w.Active && len(w.Events) == 0 {
		return ErrNoEventsSubscribed
	}

	return nil
}

// updateTimestamp updates the UpdatedAt timestamp
func (w *WebhookConfiguration) updateTimestamp() {
	w.UpdatedAt = time.Now()
}

// CanReceiveEvent checks if webhook can receive a specific event
func (w *WebhookConfiguration) CanReceiveEvent(event EventType) bool {
	return w.IsActive() && w.IsEventSubscribed(event)
}

// GetEventCount returns the number of subscribed events
func (w *WebhookConfiguration) GetEventCount() int {
	return len(w.Events)
}

// HasEvent checks if a specific event is in the subscription list
func (w *WebhookConfiguration) HasEvent(event EventType) bool {
	for _, e := range w.Events {
		if e == event {
			return true
		}
	}
	return false
}

// AddEvent adds an event to the subscription list
func (w *WebhookConfiguration) AddEvent(event EventType) error {
	if !event.IsValid() {
		return ErrInvalidEventType
	}

	if w.HasEvent(event) {
		return nil // Already subscribed
	}

	w.Events = append(w.Events, event)
	w.updateTimestamp()
	return nil
}

// RemoveEvent removes an event from the subscription list
func (w *WebhookConfiguration) RemoveEvent(event EventType) {
	for i, e := range w.Events {
		if e == event {
			w.Events = append(w.Events[:i], w.Events[i+1:]...)
			w.updateTimestamp()
			break
		}
	}
}

// ClearEvents removes all events from subscription
func (w *WebhookConfiguration) ClearEvents() {
	w.Events = []EventType{}
	w.updateTimestamp()
}

// Note: Clone method removed - not a domain responsibility
// Cloning should be handled by application or infrastructure layers
