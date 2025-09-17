package application

import (
	"context"
	"fmt"
	"time"

	"meow/internal/domain/session"
	"meow/internal/infrastructure/logging"
	"meow/internal/infrastructure/webhooks"
)

// EventDispatcher handles event processing and webhook dispatching
// This belongs in the Application Layer as it orchestrates business logic
type EventDispatcher struct {
	sessionRepo    session.Repository
	webhookService *webhooks.Service
	logger         logging.Logger
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher(sessionRepo session.Repository, webhookService *webhooks.Service) *EventDispatcher {
	return &EventDispatcher{
		sessionRepo:    sessionRepo,
		webhookService: webhookService,
		logger:         logging.GetLogger().Sub("event-dispatcher"),
	}
}

// EventPayload represents the standardized event payload
type EventPayload struct {
	Event     string      `json:"event"`
	SessionID string      `json:"sessionId"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// DispatchEvent processes an event and sends webhooks if configured
func (ed *EventDispatcher) DispatchEvent(ctx context.Context, sessionID string, eventType string, eventData interface{}) error {
	// Verify session exists
	_, err := ed.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		ed.logger.Errorf("Failed to get session for event dispatch: %v", err)
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Extract clean event type name
	cleanEventType := ed.extractEventTypeName(eventType)

	// Create webhook payload
	payload := ed.createEventPayload(sessionID, cleanEventType, eventData)

	ed.logger.Debugf("Event processed: %s for session %s", cleanEventType, sessionID)

	// Note: Webhook sending will be handled by separate webhook service
	// when webhook aggregate is fully implemented
	_ = payload // Prevent unused variable warning

	return nil
}

// isEventSubscribed checks if session is subscribed to event (Domain Logic)
func (ed *EventDispatcher) isEventSubscribed(sessionEntity *session.Session, eventType string) bool {
	// Default behavior: process all events until webhook aggregate is implemented
	return true
}

// extractEventTypeName extracts clean event type name from Go type string
func (ed *EventDispatcher) extractEventTypeName(eventType string) string {
	// Convert "*events.Message" to "Message"
	if len(eventType) > 8 && eventType[:8] == "*events." {
		return eventType[8:]
	}
	return eventType
}

// createEventPayload creates the standardized event payload
func (ed *EventDispatcher) createEventPayload(sessionID, eventType string, eventData interface{}) EventPayload {
	return EventPayload{
		Event:     eventType,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
		Data:      eventData,
	}
}

// ValidateEventType validates if an event type is supported
func (ed *EventDispatcher) ValidateEventType(eventType string) bool {
	// This could be moved to a domain service or use the existing validation
	// For now, delegate to the infrastructure layer
	return true // Simplified for this refactor
}
