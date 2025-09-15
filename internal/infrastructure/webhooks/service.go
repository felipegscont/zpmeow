package webhook

import (
	"context"
	"fmt"
	"time"

	"zpmeow/internal/infrastructure/logging"
	"zpmeow/internal/interfaces/dto"
)

type WebhookService struct {
	httpClient    HTTPClient
	retryStrategy *RetryStrategy
	logger        logging.Logger
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		httpClient:    NewWebhookHTTPClient(30 * time.Second),
		retryStrategy: NewRetryStrategy(DefaultRetryConfig()),
		logger:        logging.GetLogger().Sub("webhook-service"),
	}
}

func NewWebhookServiceWithConfig(timeout time.Duration, retryConfig *RetryConfig) *WebhookService {
	return &WebhookService{
		httpClient:    NewWebhookHTTPClient(timeout),
		retryStrategy: NewRetryStrategy(retryConfig),
		logger:        logging.GetLogger().Sub("webhook-service"),
	}
}

func (w *WebhookService) SendWebhook(ctx context.Context, webhookURL, event, sessionID string, data interface{}) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	var dataMap map[string]interface{}
	if data != nil {
		if dm, ok := data.(map[string]interface{}); ok {
			dataMap = dm
		} else {
			dataMap = map[string]interface{}{"data": data}
		}
	}

	payload := dto.WebhookEventPayload{
		EventType: event,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      dataMap,
	}
	w.logger.Infof("Sending webhook to %s for event %s (session: %s)", webhookURL, event, sessionID)

	err := w.httpClient.Post(ctx, webhookURL, payload, nil)
	if err != nil {
		w.logger.Errorf("Failed to send webhook to %s: %v", webhookURL, err)
		return err
	}

	w.logger.Infof("Successfully sent webhook to %s", webhookURL)
	return nil
}

func (w *WebhookService) SendWebhookWithRetry(ctx context.Context, webhookURL, event, sessionID string, data interface{}) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	var dataMap map[string]interface{}
	if data != nil {
		if dm, ok := data.(map[string]interface{}); ok {
			dataMap = dm
		} else {
			dataMap = map[string]interface{}{"data": data}
		}
	}

	payload := dto.WebhookEventPayload{
		EventType: event,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      dataMap,
	}
	operationName := fmt.Sprintf("webhook to %s for event %s", webhookURL, event)

	return w.retryStrategy.ExecuteWithRetry(ctx, func() error {
		w.logger.Debugf("Attempting to send %s (session: %s)", operationName, sessionID)
		return w.httpClient.Post(ctx, webhookURL, payload, nil)
	}, operationName)
}

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

func (w *WebhookService) SendWebhookWithHeaders(ctx context.Context, webhookURL, event, sessionID string, data interface{}, headers map[string]string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	var dataMap map[string]interface{}
	if data != nil {
		if dm, ok := data.(map[string]interface{}); ok {
			dataMap = dm
		} else {
			dataMap = map[string]interface{}{"data": data}
		}
	}

	payload := dto.WebhookEventPayload{
		EventType: event,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      dataMap,
	}
	w.logger.Infof("Sending webhook with headers to %s for event %s (session: %s)", webhookURL, event, sessionID)

	err := w.httpClient.Post(ctx, webhookURL, payload, headers)
	if err != nil {
		w.logger.Errorf("Failed to send webhook to %s: %v", webhookURL, err)
		return err
	}

	w.logger.Infof("Successfully sent webhook to %s", webhookURL)
	return nil
}

func (w *WebhookService) SendWebhookWithHeadersAndRetry(ctx context.Context, webhookURL, event, sessionID string, data interface{}, headers map[string]string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	var dataMap map[string]interface{}
	if data != nil {
		if dm, ok := data.(map[string]interface{}); ok {
			dataMap = dm
		} else {
			dataMap = map[string]interface{}{"data": data}
		}
	}

	payload := dto.WebhookEventPayload{
		EventType: event,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      dataMap,
	}
	operationName := fmt.Sprintf("webhook with headers to %s for event %s", webhookURL, event)

	return w.retryStrategy.ExecuteWithRetry(ctx, func() error {
		w.logger.Debugf("Attempting to send %s (session: %s)", operationName, sessionID)
		return w.httpClient.Post(ctx, webhookURL, payload, headers)
	}, operationName)
}
