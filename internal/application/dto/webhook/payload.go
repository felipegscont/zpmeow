package webhook

import "time"

// WebhookPayload represents the structure of webhook data sent to clients
type WebhookPayload struct {
	Event     string      `json:"event" example:"message"`
	SessionID string      `json:"sessionId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp int64       `json:"timestamp" example:"1640995200"`
	Data      interface{} `json:"data"`
}

// NewWebhookPayload creates a new webhook payload
func NewWebhookPayload(event, sessionID string, data interface{}) *WebhookPayload {
	return &WebhookPayload{
		Event:     event,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}
