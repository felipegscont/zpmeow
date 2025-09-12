package request

// SetWebhookRequest represents a request to set webhook configuration
type SetWebhookRequest struct {
	WebhookURL string   `json:"webhookurl" binding:"required" example:"https://example.com/webhook"`
	Events     []string `json:"events,omitempty" example:"message,status"`
}

// UpdateWebhookRequest represents a request to update webhook configuration
type UpdateWebhookRequest struct {
	WebhookURL string   `json:"webhook" binding:"required" example:"https://example.com/webhook"`
	Events     []string `json:"events,omitempty" example:"message,status"`
	Active     bool     `json:"active" example:"true"`
}
