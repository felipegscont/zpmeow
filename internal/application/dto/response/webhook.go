package response

import "time"

// WebhookResponse represents the response from webhook operations
type WebhookResponse struct {
	Webhook   string   `json:"webhook" example:"https://example.com/webhook"`
	Events    []string `json:"events,omitempty" example:"message,status"`
	Active    bool     `json:"active,omitempty" example:"true"`
	Subscribe []string `json:"subscribe,omitempty" example:"message,status"`
}

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
