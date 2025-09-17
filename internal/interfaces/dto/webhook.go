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
	URL    string   `json:"url" binding:"required" example:"https://webhook.example.com/meow"`
	Events []string `json:"events" binding:"required" example:"Message,Receipt,Connected"`
	Secret string   `json:"secret,omitempty" example:"webhook_secret_key"`
}

// UpdateWebhookRequest represents a request to update a webhook
type UpdateWebhookRequest struct {
	URL    string   `json:"url,omitempty" example:"https://webhook.example.com/meow"`
	Events []string `json:"events,omitempty" example:"Message,Receipt"`
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
	URL       string    `json:"url" example:"https://webhook.example.com/meow"`
	Events    []string  `json:"events" example:"Message,Receipt,Connected"`
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
	ChatJID     string    `json:"chat_jid" example:"5511999999999@s.meow.net"`
	FromJID     string    `json:"from_jid" example:"5511888888888@s.meow.net"`
	MessageType string    `json:"message_type" example:"text"`
	Content     string    `json:"content" example:"Hello, World!"`
	MediaURL    string    `json:"media_url,omitempty" example:"https://example.com/image.jpg"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	IsFromMe    bool      `json:"is_from_me" example:"false"`
}

// StatusWebhookData represents status webhook event data
type StatusWebhookData struct {
	MessageID string    `json:"message_id" example:"msg_123456789"`
	ChatJID   string    `json:"chat_jid" example:"5511999999999@s.meow.net"`
	Status    string    `json:"status" example:"delivered"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

// ConnectionWebhookData represents connection webhook event data
type ConnectionWebhookData struct {
	Status    string    `json:"status" example:"connected"`
	JID       string    `json:"jid,omitempty" example:"5511999999999@s.meow.net"`
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

// SupportedEventsResponse represents supported events response
type SupportedEventsResponse struct {
	Status  int                 `json:"status" example:"200"`
	Message string              `json:"message" example:"Supported events retrieved successfully"`
	Data    SupportedEventsData `json:"data"`
}

// SupportedEventsData represents supported events data
type SupportedEventsData struct {
	Events []string `json:"events" example:"[\"Message\", \"Receipt\", \"Connected\"]"`
	Count  int      `json:"count" example:"65"`
}

// ============================================================================
// NEW STANDARDIZED RESPONSE STRUCTURES
// ============================================================================

// StandardWebhookResponse represents the new standardized webhook response format
type StandardWebhookResponse struct {
	Data    StandardWebhookData `json:"data"`
	Message string              `json:"message" example:"Webhook retrieved successfully"`
	Status  int                 `json:"status" example:"200"`
}

// StandardWebhookData represents the data structure for standardized webhook responses
type StandardWebhookData struct {
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	Events    []string  `json:"events" example:"Message,Receipt,Connected"`
	SessionID string    `json:"sessionID" example:"default"`
	Status    string    `json:"status" example:"active"`
	URL       string    `json:"url" example:"https://webhook.example.com/meow"`
}

// StandardWebhookListResponse represents the standardized list webhooks response
type StandardWebhookListResponse struct {
	Data    []StandardWebhookData `json:"data"`
	Message string                `json:"message" example:"Webhooks retrieved successfully"`
	Status  int                   `json:"status" example:"200"`
}

// StandardWebhookCreateResponse represents the standardized create webhook response
type StandardWebhookCreateResponse struct {
	Data    StandardWebhookData `json:"data"`
	Message string              `json:"message" example:"Webhook created successfully"`
	Status  int                 `json:"status" example:"201"`
}

// StandardWebhookUpdateResponse represents the standardized update webhook response
type StandardWebhookUpdateResponse struct {
	Data    StandardWebhookData `json:"data"`
	Message string              `json:"message" example:"Webhook updated successfully"`
	Status  int                 `json:"status" example:"200"`
}

// StandardWebhookDeleteResponse represents the standardized delete webhook response
type StandardWebhookDeleteResponse struct {
	Data    StandardWebhookDeleteData `json:"data"`
	Message string                    `json:"message" example:"Webhook deleted successfully"`
	Status  int                       `json:"status" example:"200"`
}

// StandardWebhookDeleteData represents the data for delete webhook response
type StandardWebhookDeleteData struct {
	SessionID string `json:"sessionID" example:"default"`
	Status    string `json:"status" example:"deleted"`
}

// Note: Validation functions moved to internal/application/webhook.go
// DTOs should only contain data structures, not business logic
