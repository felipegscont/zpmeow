package dto

import (
	"fmt"
	"time"
)

// ============================================================================
// WEBHOOK REQUEST DTOs
// ============================================================================

// RegisterWebhookRequest represents a request to register a webhook
type RegisterWebhookRequest struct {
	SessionID string   `json:"session_id" binding:"required" example:"default"`
	URL       string   `json:"url" binding:"required" example:"https://webhook.example.com/whatsapp"`
	Events    []string `json:"events" binding:"required" example:"[\"message\", \"status\", \"connection\"]"`
	Secret    string   `json:"secret,omitempty" example:"webhook_secret_key"`
}

// UpdateWebhookRequest represents a request to update a webhook
type UpdateWebhookRequest struct {
	URL    string   `json:"url,omitempty" example:"https://webhook.example.com/whatsapp"`
	Events []string `json:"events,omitempty" example:"[\"message\", \"status\"]"`
	Secret string   `json:"secret,omitempty" example:"new_webhook_secret_key"`
	Status string   `json:"status,omitempty" example:"active"`
}

// TestWebhookRequest represents a request to test a webhook
type TestWebhookRequest struct {
	EventType string                 `json:"event_type" binding:"required" example:"message"`
	TestData  map[string]interface{} `json:"test_data,omitempty"`
}

// ============================================================================
// WEBHOOK DATA STRUCTURES
// ============================================================================

// WebhookInfo represents comprehensive webhook information
type WebhookInfo struct {
	WebhookID string    `json:"webhook_id" example:"webhook_123456789"`
	SessionID string    `json:"session_id" example:"default"`
	URL       string    `json:"url" example:"https://webhook.example.com/whatsapp"`
	Events    []string  `json:"events" example:"[\"message\", \"status\", \"connection\"]"`
	Status    string    `json:"status" example:"active"`
	Secret    string    `json:"secret,omitempty" example:"webhook_secret_key"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// WebhookEventPayload represents a webhook event payload
type WebhookEventPayload struct {
	EventType string                 `json:"event_type" example:"message"`
	SessionID string                 `json:"session_id" example:"default"`
	Timestamp time.Time              `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Data      map[string]interface{} `json:"data"`
}

// MessageWebhookData represents message webhook event data
type MessageWebhookData struct {
	MessageID   string    `json:"message_id" example:"msg_123456789"`
	ChatJID     string    `json:"chat_jid" example:"5511999999999@s.whatsapp.net"`
	FromJID     string    `json:"from_jid" example:"5511888888888@s.whatsapp.net"`
	MessageType string    `json:"message_type" example:"text"`
	Content     string    `json:"content" example:"Hello, World!"`
	MediaURL    string    `json:"media_url,omitempty" example:"https://example.com/image.jpg"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	IsFromMe    bool      `json:"is_from_me" example:"false"`
}

// StatusWebhookData represents status webhook event data
type StatusWebhookData struct {
	MessageID string    `json:"message_id" example:"msg_123456789"`
	ChatJID   string    `json:"chat_jid" example:"5511999999999@s.whatsapp.net"`
	Status    string    `json:"status" example:"delivered"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

// ConnectionWebhookData represents connection webhook event data
type ConnectionWebhookData struct {
	Status    string    `json:"status" example:"connected"`
	JID       string    `json:"jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

// TestWebhookResult represents webhook test result
type TestWebhookResult struct {
	WebhookID    string `json:"webhook_id" example:"webhook_123456789"`
	TestResult   string `json:"test_result" example:"success"`
	ResponseCode int    `json:"response_code" example:"200"`
	ResponseTime int64  `json:"response_time_ms" example:"150"`
	Error        string `json:"error,omitempty" example:""`
}

// ============================================================================
// WEBHOOK RESPONSE DTOs
// ============================================================================

// WebhookResponse represents the standardized response format for webhook operations
type WebhookResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Data    WebhookResponseData   `json:"data"`
	Error   *WebhookErrorResponse `json:"error,omitempty"`
}

// WebhookResponseData contains the response data for webhook operations
type WebhookResponseData struct {
	WebhookID  string             `json:"webhook_id,omitempty" example:"webhook_123456789"`
	Action     string             `json:"action" example:"register"`
	Status     string             `json:"status" example:"success"`
	Timestamp  time.Time          `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Webhook    *WebhookInfo       `json:"webhook,omitempty"`
	Webhooks   []WebhookInfo      `json:"webhooks,omitempty"`
	TestResult *TestWebhookResult `json:"test_result,omitempty"`
}

// WebhookErrorResponse represents error information for webhook operations
type WebhookErrorResponse struct {
	Code    string `json:"code" example:"INVALID_WEBHOOK_URL"`
	Message string `json:"message" example:"Invalid webhook URL format"`
	Details string `json:"details,omitempty" example:"Webhook URL must be a valid HTTPS URL"`
}

// WebhookValidationError represents a validation error for webhook requests
type WebhookValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *WebhookValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ============================================================================
// LEGACY RESPONSE STRUCTURES (for backward compatibility)
// ============================================================================

// RegisterWebhookResponse represents legacy register webhook response
type RegisterWebhookResponse struct {
	Status  int                         `json:"status" example:"201"`
	Message string                      `json:"message" example:"Webhook registered successfully"`
	Data    RegisterWebhookResponseData `json:"data"`
}

// RegisterWebhookResponseData represents legacy register webhook response data
type RegisterWebhookResponseData struct {
	WebhookID string    `json:"webhook_id" example:"webhook_123456789"`
	SessionID string    `json:"session_id" example:"default"`
	URL       string    `json:"url" example:"https://webhook.example.com/whatsapp"`
	Events    []string  `json:"events" example:"[\"message\", \"status\", \"connection\"]"`
	Status    string    `json:"status" example:"active"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// GetWebhookResponse represents legacy get webhook response
type GetWebhookResponse struct {
	Status  int                    `json:"status" example:"200"`
	Message string                 `json:"message" example:"Webhook retrieved successfully"`
	Data    GetWebhookResponseData `json:"data"`
}

// GetWebhookResponseData represents legacy get webhook response data
type GetWebhookResponseData struct {
	WebhookID string    `json:"webhook_id" example:"webhook_123456789"`
	SessionID string    `json:"session_id" example:"default"`
	URL       string    `json:"url" example:"https://webhook.example.com/whatsapp"`
	Events    []string  `json:"events" example:"[\"message\", \"status\", \"connection\"]"`
	Status    string    `json:"status" example:"active"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// UpdateWebhookResponse represents legacy update webhook response
type UpdateWebhookResponse struct {
	Status  int                       `json:"status" example:"200"`
	Message string                    `json:"message" example:"Webhook updated successfully"`
	Data    UpdateWebhookResponseData `json:"data"`
}

// UpdateWebhookResponseData represents legacy update webhook response data
type UpdateWebhookResponseData struct {
	WebhookID string    `json:"webhook_id" example:"webhook_123456789"`
	SessionID string    `json:"session_id" example:"default"`
	URL       string    `json:"url" example:"https://webhook.example.com/whatsapp"`
	Events    []string  `json:"events" example:"[\"message\", \"status\", \"connection\"]"`
	Status    string    `json:"status" example:"active"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ListWebhooksResponse represents legacy list webhooks response
type ListWebhooksResponse struct {
	Status  int                        `json:"status" example:"200"`
	Message string                     `json:"message" example:"Webhooks retrieved successfully"`
	Data    []ListWebhooksResponseData `json:"data"`
}

// ListWebhooksResponseData represents legacy list webhooks response data (single item)
type ListWebhooksResponseData struct {
	WebhookID string    `json:"webhook_id" example:"webhook_123456789"`
	SessionID string    `json:"session_id" example:"default"`
	URL       string    `json:"url" example:"https://webhook.example.com/whatsapp"`
	Events    []string  `json:"events" example:"[\"message\", \"status\", \"connection\"]"`
	Status    string    `json:"status" example:"active"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// DeleteWebhookResponse represents legacy delete webhook response
type DeleteWebhookResponse struct {
	Status  int                       `json:"status" example:"200"`
	Message string                    `json:"message" example:"Webhook deleted successfully"`
	Data    DeleteWebhookResponseData `json:"data"`
}

// DeleteWebhookResponseData represents legacy delete webhook response data
type DeleteWebhookResponseData struct {
	WebhookID string `json:"webhook_id" example:"webhook_123456789"`
	Status    string `json:"status" example:"deleted"`
}

// TestWebhookResponse represents legacy test webhook response
type TestWebhookResponse struct {
	Status  int                     `json:"status" example:"200"`
	Message string                  `json:"message" example:"Webhook test completed"`
	Data    TestWebhookResponseData `json:"data"`
}

// TestWebhookResponseData represents legacy test webhook response data
type TestWebhookResponseData struct {
	WebhookID    string `json:"webhook_id" example:"webhook_123456789"`
	TestResult   string `json:"test_result" example:"success"`
	ResponseCode int    `json:"response_code" example:"200"`
	ResponseTime int64  `json:"response_time_ms" example:"150"`
	Error        string `json:"error,omitempty" example:""`
}

// ============================================================================
// INTERNAL HELPER FUNCTIONS
// ============================================================================

// validateWebhookURL checks if a URL is valid (basic validation)
func validateWebhookURL(url string) bool {
	return url != "" && len(url) > 8 && url[:8] == "https://"
}

// validateWebhookSessionID checks if a session ID is valid
func validateWebhookSessionID(sessionID string) bool {
	return sessionID != "" && len(sessionID) >= 1
}

// validateEventType checks if an event type is valid
func validateEventType(eventType string) bool {
	validEvents := []string{"message", "status", "connection", "call", "contact"}
	for _, validEvent := range validEvents {
		if eventType == validEvent {
			return true
		}
	}
	return false
}

// validateEvents checks if all events in the slice are valid
func validateEvents(events []string) bool {
	if len(events) == 0 {
		return false
	}
	for _, event := range events {
		if !validateEventType(event) {
			return false
		}
	}
	return true
}

// validateWebhookStatus checks if a webhook status is valid
func validateWebhookStatus(status string) bool {
	validStatuses := []string{"active", "inactive", "paused"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}
