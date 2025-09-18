package application

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/domain/session"
)

type WebhookInfo struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Active bool     `json:"active"`
}

type WebhookService interface {
	SetWebhook(ctx context.Context, sessionID, webhookURL string, events []string) error

	GetWebhook(ctx context.Context, sessionID string) (*WebhookInfo, error)

	UpdateWebhook(ctx context.Context, sessionID, webhookURL string, events []string, active bool) error

	DeleteWebhook(ctx context.Context, sessionID string) error

	TestWebhook(ctx context.Context, sessionID, message string) error

	ValidateWebhookURL(url string) error

	ValidateEvents(events []string) error
}

type WebhookApp struct {
	sessionRepo    session.Repository
	webhookService *webhooks.Service
}

func NewWebhookApp(sessionRepo session.Repository, webhookService *webhooks.Service) WebhookService {
	return &WebhookApp{
		sessionRepo:    sessionRepo,
		webhookService: webhookService,
	}
}

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

func (w *WebhookApp) ValidateEvents(events []string) error {
	if len(events) == 0 {
		return fmt.Errorf("at least one event must be specified")
	}

	for _, event := range events {
		if strings.TrimSpace(event) == "" {
			return fmt.Errorf("event cannot be empty")
		}
	}

	return nil
}

func (w *WebhookApp) SetWebhook(ctx context.Context, sessionID, webhookURL string, events []string) error {
	if err := w.ValidateWebhookURL(webhookURL); err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	if err := w.ValidateEvents(events); err != nil {
		return fmt.Errorf("invalid events: %w", err)
	}

	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	return nil
}

func (w *WebhookApp) GetWebhook(ctx context.Context, sessionID string) (*WebhookInfo, error) {
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	return &WebhookInfo{
		URL:    "",
		Events: []string{},
		Active: false,
	}, nil
}

func (w *WebhookApp) UpdateWebhook(ctx context.Context, sessionID, webhookURL string, events []string, active bool) error {
	if active {
		if err := w.ValidateWebhookURL(webhookURL); err != nil {
			return fmt.Errorf("invalid webhook URL: %w", err)
		}

		if err := w.ValidateEvents(events); err != nil {
			return fmt.Errorf("invalid events: %w", err)
		}
	}

	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	return nil
}

func (w *WebhookApp) DeleteWebhook(ctx context.Context, sessionID string) error {
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	return nil
}

func (w *WebhookApp) TestWebhook(ctx context.Context, sessionID, message string) error {
	_, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	return fmt.Errorf("webhook testing not available - webhook aggregate not implemented")
}
