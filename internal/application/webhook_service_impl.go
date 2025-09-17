package application

import (
	"context"
	"fmt"

	"meow/internal/domain/session"
	"meow/internal/infrastructure/webhooks"
)

// WebhookServiceImpl implements the WebhookService interface
type WebhookServiceImpl struct {
	sessionRepo    session.Repository
	webhookService *webhooks.Service
}

// NewWebhookService creates a new webhook service implementation
func NewWebhookService(sessionRepo session.Repository, webhookService *webhooks.Service) WebhookService {
	return &WebhookServiceImpl{
		sessionRepo:    sessionRepo,
		webhookService: webhookService,
	}
}

// SetWebhook configures a webhook for a session
func (w *WebhookServiceImpl) SetWebhook(ctx context.Context, sessionID, webhookURL string, events []string) error {
	// Get session
	sessionEntity, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Set webhook on session entity
	sessionEntity.SetWebhook(webhookURL, events)

	// Save session
	err = w.sessionRepo.Update(ctx, sessionEntity)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// GetWebhook retrieves webhook configuration for a session
func (w *WebhookServiceImpl) GetWebhook(ctx context.Context, sessionID string) (*WebhookInfo, error) {
	// Get session
	sessionEntity, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if webhook is configured
	if !sessionEntity.HasWebhook() {
		return nil, nil
	}

	return &WebhookInfo{
		URL:    sessionEntity.WebhookURL,
		Events: sessionEntity.Events,
		Active: true,
	}, nil
}

// UpdateWebhook updates webhook configuration for a session
func (w *WebhookServiceImpl) UpdateWebhook(ctx context.Context, sessionID, webhookURL string, events []string, active bool) error {
	// Get session
	sessionEntity, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if active {
		// Set webhook
		sessionEntity.SetWebhook(webhookURL, events)
	} else {
		// Clear webhook
		sessionEntity.ClearWebhook()
	}

	// Save session
	err = w.sessionRepo.Update(ctx, sessionEntity)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// DeleteWebhook removes webhook configuration for a session
func (w *WebhookServiceImpl) DeleteWebhook(ctx context.Context, sessionID string) error {
	// Get session
	sessionEntity, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Clear webhook
	sessionEntity.ClearWebhook()

	// Save session
	err = w.sessionRepo.Update(ctx, sessionEntity)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// TestWebhook sends a test webhook to verify configuration
func (w *WebhookServiceImpl) TestWebhook(ctx context.Context, sessionID, message string) error {
	// Get session
	sessionEntity, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Check if webhook is configured
	if !sessionEntity.HasWebhook() {
		return fmt.Errorf("no webhook configured for session")
	}

	// Create test payload
	testPayload := map[string]interface{}{
		"type":      "test",
		"sessionId": sessionID,
		"message":   message,
		"timestamp": fmt.Sprintf("%d", ctx.Value("timestamp")),
	}

	// Send test webhook
	err = w.webhookService.SendWebhook(ctx, sessionEntity.WebhookURL, "test", sessionID, testPayload)
	if err != nil {
		return fmt.Errorf("failed to send test webhook: %w", err)
	}

	return nil
}
