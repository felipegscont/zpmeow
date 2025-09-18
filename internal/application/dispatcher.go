package application

import (
	"context"
	"fmt"
	"time"

	"meow/internal/domain/session"
)

// Interfaces movidas para interfaces.go para evitar duplicação

type EventDispatcher struct {
	sessionRepo   session.Repository
	webhookSender WebhookSender
	logger        Logger
}

func NewEventDispatcher(sessionRepo session.Repository, webhookSender WebhookSender, logger Logger) *EventDispatcher {
	return &EventDispatcher{
		sessionRepo:   sessionRepo,
		webhookSender: webhookSender,
		logger:        logger,
	}
}

func (ed *EventDispatcher) DispatchEvent(ctx context.Context, sessionID string, eventType string, eventData interface{}) error {
	_, err := ed.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		ed.logger.Errorf("Failed to get session for event dispatch: %v", err)
		return fmt.Errorf("failed to get session: %w", err)
	}

	cleanEventType := ed.extractEventTypeName(eventType)

	payload := ed.createEventPayload(sessionID, cleanEventType, eventData)

	// Enviar webhook se configurado
	if err := ed.webhookSender.SendWebhook(ctx, sessionID, "", cleanEventType, payload); err != nil {
		ed.logger.Errorf("Failed to send webhook for session %s: %v", sessionID, err)
		// Não retornar erro para não bloquear o fluxo principal
	}

	ed.logger.Infof("Event processed: %s for session %s", cleanEventType, sessionID)

	return nil
}

func (ed *EventDispatcher) extractEventTypeName(eventType string) string {
	if len(eventType) > 8 && eventType[:8] == "*events." {
		return eventType[8:]
	}
	return eventType
}

func (ed *EventDispatcher) createEventPayload(sessionID, eventType string, eventData interface{}) EventData {
	return EventData{
		Type:      eventType,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
		Payload:   map[string]interface{}{"data": eventData},
	}
}

func (ed *EventDispatcher) ValidateEventType(eventType string) bool {
	validEvents := map[string]bool{
		"message.received":     true,
		"message.sent":         true,
		"session.connected":    true,
		"session.disconnected": true,
		"qr.generated":         true,
		"session.error":        true,
		"session.connecting":   true,
	}

	return validEvents[eventType]
}
