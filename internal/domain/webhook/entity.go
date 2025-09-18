package webhook

import (
	"time"
)

type WebhookConfiguration struct {
	SessionID SessionID
	URL       WebhookURL
	Events    []EventType
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

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

func (w *WebhookConfiguration) IsActive() bool {
	return w.Active && !w.URL.IsEmpty()
}

func (w *WebhookConfiguration) IsEventSubscribed(event EventType) bool {
	if !w.IsActive() {
		return false
	}

	for _, e := range w.Events {
		if e == EventTypeAll || e == event {
			return true
		}
	}

	return false
}

func (w *WebhookConfiguration) UpdateURL(url string) error {
	webhookURL, err := NewWebhookURL(url)
	if err != nil {
		return err
	}

	w.URL = webhookURL
	w.updateTimestamp()
	return nil
}

func (w *WebhookConfiguration) UpdateEvents(events []string) error {
	eventTypes, err := ParseEventTypes(events)
	if err != nil {
		return err
	}

	w.Events = eventTypes
	w.updateTimestamp()
	return nil
}

func (w *WebhookConfiguration) Activate() {
	w.Active = true
	w.updateTimestamp()
}

func (w *WebhookConfiguration) Deactivate() {
	w.Active = false
	w.updateTimestamp()
}


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

func (w *WebhookConfiguration) updateTimestamp() {
	w.UpdatedAt = time.Now()
}

func (w *WebhookConfiguration) CanReceiveEvent(event EventType) bool {
	return w.IsActive() && w.IsEventSubscribed(event)
}

func (w *WebhookConfiguration) GetEventCount() int {
	return len(w.Events)
}

func (w *WebhookConfiguration) HasEvent(event EventType) bool {
	for _, e := range w.Events {
		if e == event {
			return true
		}
	}
	return false
}

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

func (w *WebhookConfiguration) RemoveEvent(event EventType) {
	for i, e := range w.Events {
		if e == event {
			w.Events = append(w.Events[:i], w.Events[i+1:]...)
			w.updateTimestamp()
			break
		}
	}
}

func (w *WebhookConfiguration) ClearEvents() {
	w.Events = []EventType{}
	w.updateTimestamp()
}

