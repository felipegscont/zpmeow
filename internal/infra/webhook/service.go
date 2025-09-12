package webhook

import (
	"context"
	"fmt"
	"time"

	"zpmeow/internal/application/dto/response"
	"zpmeow/internal/infra/logger"
)

// WebhookService handles sending webhook events
type WebhookService struct {
	httpClient HTTPClient
	logger     logger.Logger
}

// NewWebhookService creates a new webhook service
func NewWebhookService() *WebhookService {
	return &WebhookService{
		httpClient: NewWebhookHTTPClient(30 * time.Second),
		logger:     logger.GetLogger().Sub("webhook-service"),
	}
}

// SendWebhook sends a webhook event to the specified URL
func (w *WebhookService) SendWebhook(ctx context.Context, webhookURL, event, sessionID string, data interface{}) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	payload := response.NewWebhookPayload(event, sessionID, data)
	w.logger.Infof("Sending webhook to %s for event %s (session: %s)", webhookURL, event, sessionID)

	err := w.httpClient.Post(ctx, webhookURL, payload, nil)
	if err != nil {
		w.logger.Errorf("Failed to send webhook to %s: %v", webhookURL, err)
		return err
	}

	w.logger.Infof("Successfully sent webhook to %s", webhookURL)
	return nil
}

// SendWebhookWithRetry sends a webhook with retry logic
func (w *WebhookService) SendWebhookWithRetry(ctx context.Context, webhookURL, event, sessionID string, data interface{}) error {
	// For now, just call SendWebhook directly without retry
	// TODO: Implement proper retry strategy
	return w.SendWebhook(ctx, webhookURL, event, sessionID, data)
}

// SendWebhookAsync sends a webhook asynchronously without blocking
func (w *WebhookService) SendWebhookAsync(webhookURL, event, sessionID string, data interface{}) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		err := w.SendWebhookWithRetry(ctx, webhookURL, event, sessionID, data)
		if err != nil {
			w.logger.Errorf("Async webhook failed: %v", err)
		}
	}()
}
