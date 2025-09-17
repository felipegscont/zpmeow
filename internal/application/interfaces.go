package application

import "context"

// EventProcessor defines the interface for processing WhatsApp events
// This belongs in Application layer as it defines use case capabilities
type EventProcessor interface {
	// HandleEvent processes a WhatsApp event
	HandleEvent(evt interface{})
}

// EventDispatcherInterface defines the interface for dispatching events to webhooks
// This belongs in Application layer as it orchestrates business logic
type EventDispatcherInterface interface {
	// DispatchEvent processes an event and sends webhooks if configured
	DispatchEvent(ctx context.Context, sessionID string, eventType string, eventData interface{}) error

	// ValidateEventType validates if an event type is supported
	ValidateEventType(eventType string) bool
}
