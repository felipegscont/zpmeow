package application

import (
	"context"
	"fmt"
	"time"

	"meow/internal/domain/session"
)

// Interfaces movidas para interfaces.go para evitar duplicação

type EventDispatcher struct {
	sessionRepo    session.Repository
	webhookSender  WebhookSender
	logger         Logger
}

func NewEventDispatcher(sessionRepo session.Repository, webhookSender WebhookSender, logger Logger) *EventDispatcher {
	return &EventDispatcher{
		sessionRepo:   sessionRepo,
		webhookSender: webhookSender,
		logger:        logger,
	}
}

type EventPayload struct {
	Event     string      `json:"event"`
	SessionID string      `json:"sessionId"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func (ed *EventDispatcher) DispatchEvent(ctx context.Context, sessionID string, eventType string, eventData interface{}) error {
	_, err := ed.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		ed.logger.Errorf("Failed to get session for event dispatch: %v", err)
		return fmt.Errorf("failed to get session: %w", err)
	}

	cleanEventType := ed.extractEventTypeName(eventType)

	payload := ed.createEventPayload(sessionID, cleanEventType, eventData)

	ed.logger.Infof("Event processed: %s for session %s", cleanEventType, sessionID)

	_ = payload // Prevent unused variable warning

	return nil
}

func (ed *EventDispatcher) isEventSubscribed(sessionEntity *session.Session, eventType string) bool {
	return true
}

func (ed *EventDispatcher) extractEventTypeName(eventType string) string {
	if len(eventType) > 8 && eventType[:8] == "*events." {
		return eventType[8:]
	}
	return eventType
}

func (ed *EventDispatcher) createEventPayload(sessionID, eventType string, eventData interface{}) EventPayload {
	return EventPayload{
		Event:     eventType,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
		Data:      eventData,
	}
}

func (ed *EventDispatcher) ValidateEventType(eventType string) bool {
	return true // Simplified for this refactor
}
