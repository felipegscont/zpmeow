package application

import "context"

type EventProcessor interface {
	HandleEvent(evt interface{})
}

type EventDispatcherInterface interface {
	DispatchEvent(ctx context.Context, sessionID string, eventType string, eventData interface{}) error

	ValidateEventType(eventType string) bool
}
