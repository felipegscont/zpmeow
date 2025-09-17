package application

import (
	"context"
)

// WebhookInfo represents webhook configuration information
type WebhookInfo struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Active bool     `json:"active"`
}

// WebhookService defines the interface for webhook operations
type WebhookService interface {
	// SetWebhook configures a webhook for a session
	SetWebhook(ctx context.Context, sessionID, webhookURL string, events []string) error
	
	// GetWebhook retrieves webhook configuration for a session
	GetWebhook(ctx context.Context, sessionID string) (*WebhookInfo, error)
	
	// UpdateWebhook updates webhook configuration for a session
	UpdateWebhook(ctx context.Context, sessionID, webhookURL string, events []string, active bool) error
	
	// DeleteWebhook removes webhook configuration for a session
	DeleteWebhook(ctx context.Context, sessionID string) error
	
	// TestWebhook sends a test webhook to verify configuration
	TestWebhook(ctx context.Context, sessionID, message string) error
}
