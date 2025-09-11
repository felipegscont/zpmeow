package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"zpmeow/internal/infra/logger"
)

// WebhookService handles sending webhook events
type WebhookService struct {
	client *http.Client
	logger logger.Logger
}

// NewWebhookService creates a new webhook service
func NewWebhookService() *WebhookService {
	return &WebhookService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.GetLogger().Sub("webhook-service"),
	}
}

// WebhookPayload represents the structure of webhook data sent to clients
type WebhookPayload struct {
	Event     string      `json:"event"`
	SessionID string      `json:"sessionId"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// SendWebhook sends a webhook event to the specified URL
func (w *WebhookService) SendWebhook(ctx context.Context, webhookURL, event, sessionID string, data interface{}) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	payload := WebhookPayload{
		Event:     event,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		w.logger.Errorf("Failed to marshal webhook payload: %v", err)
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		w.logger.Errorf("Failed to create webhook request: %v", err)
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "zpmeow-webhook/1.0")

	w.logger.Infof("Sending webhook to %s for event %s (session: %s)", webhookURL, event, sessionID)

	resp, err := w.client.Do(req)
	if err != nil {
		w.logger.Errorf("Failed to send webhook to %s: %v", webhookURL, err)
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		w.logger.Warnf("Webhook returned non-2xx status: %d for URL %s", resp.StatusCode, webhookURL)
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	w.logger.Infof("Successfully sent webhook to %s (status: %d)", webhookURL, resp.StatusCode)
	return nil
}

// SendWebhookWithRetry sends a webhook with retry logic
func (w *WebhookService) SendWebhookWithRetry(ctx context.Context, webhookURL, event, sessionID string, data interface{}, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, 8s...
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			w.logger.Infof("Retrying webhook in %v (attempt %d/%d)", backoff, attempt, maxRetries)
			
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := w.SendWebhook(ctx, webhookURL, event, sessionID, data)
		if err == nil {
			if attempt > 0 {
				w.logger.Infof("Webhook succeeded after %d retries", attempt)
			}
			return nil
		}

		lastErr = err
		w.logger.Warnf("Webhook attempt %d failed: %v", attempt+1, err)
	}

	w.logger.Errorf("Webhook failed after %d attempts: %v", maxRetries+1, lastErr)
	return fmt.Errorf("webhook failed after %d attempts: %w", maxRetries+1, lastErr)
}

// SendWebhookAsync sends a webhook asynchronously without blocking
func (w *WebhookService) SendWebhookAsync(webhookURL, event, sessionID string, data interface{}) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		err := w.SendWebhookWithRetry(ctx, webhookURL, event, sessionID, data, 3)
		if err != nil {
			w.logger.Errorf("Async webhook failed: %v", err)
		}
	}()
}
