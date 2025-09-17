package application

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/domain/session"
	"meow/internal/infrastructure/webhooks"
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

	// ValidateWebhookURL validates webhook URL format
	ValidateWebhookURL(url string) error

	// ValidateEvents validates webhook events
	ValidateEvents(events []string) error
}

// WebhookApp implements the WebhookService interface following the application layer pattern
type WebhookApp struct {
	sessionRepo    session.Repository
	webhookService *webhooks.Service
}

// NewWebhookApp creates a new webhook application service
func NewWebhookApp(sessionRepo session.Repository, webhookService *webhooks.Service) WebhookService {
	return &WebhookApp{
		sessionRepo:    sessionRepo,
		webhookService: webhookService,
	}
}

// ValidateWebhookURL validates webhook URL format
func (w *WebhookApp) ValidateWebhookURL(url string) error {
	if url == "" {
		return fmt.Errorf("webhook URL is required")
	}

	url = strings.TrimSpace(url)
	if len(url) < 8 {
		return fmt.Errorf("webhook URL is too short")
	}

	if !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("webhook URL must use HTTPS")
	}

	return nil
}

// ValidateEvents validates webhook events
func (w *WebhookApp) ValidateEvents(events []string) error {
	if len(events) == 0 {
		return fmt.Errorf("at least one event must be specified")
	}

	// Note: Event validation logic will be moved here from handlers
	// For now, basic validation
	for _, event := range events {
		if strings.TrimSpace(event) == "" {
			return fmt.Errorf("event cannot be empty")
		}
	}

	return nil
}

// SetWebhook configures a webhook for a session
func (w *WebhookApp) SetWebhook(ctx context.Context, sessionID, webhookURL string, events []string) error {
	// Validate inputs
	if err := w.ValidateWebhookURL(webhookURL); err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	if err := w.ValidateEvents(events); err != nil {
		return fmt.Errorf("invalid events: %w", err)
	}

	// Verify session exists
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Webhook configuration is now handled by separate webhook aggregate
	// For now, just validate the inputs and return success
	return nil
}

// GetWebhook retrieves webhook configuration for a session
func (w *WebhookApp) GetWebhook(ctx context.Context, sessionID string) (*WebhookInfo, error) {
	// Verify session exists
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Webhook configuration is now handled by separate webhook aggregate
	return &WebhookInfo{
		URL:    "",
		Events: []string{},
		Active: false,
	}, nil
}

// UpdateWebhook updates webhook configuration for a session
func (w *WebhookApp) UpdateWebhook(ctx context.Context, sessionID, webhookURL string, events []string, active bool) error {
	// Validate inputs if webhook is being activated
	if active {
		if err := w.ValidateWebhookURL(webhookURL); err != nil {
			return fmt.Errorf("invalid webhook URL: %w", err)
		}

		if err := w.ValidateEvents(events); err != nil {
			return fmt.Errorf("invalid events: %w", err)
		}
	}

	// Verify session exists
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Webhook configuration is now handled by separate webhook aggregate
	return nil
}

// DeleteWebhook removes webhook configuration for a session
func (w *WebhookApp) DeleteWebhook(ctx context.Context, sessionID string) error {
	// Verify session exists
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Webhook configuration is now handled by separate webhook aggregate
	return nil
}

// TestWebhook sends a test webhook to verify configuration
func (w *WebhookApp) TestWebhook(ctx context.Context, sessionID, message string) error {
	// Verify session exists
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Webhook testing is now handled by separate webhook aggregate
	return fmt.Errorf("webhook testing not available - webhook aggregate not implemented")
}
