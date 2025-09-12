package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPClient interface for webhook HTTP operations
type HTTPClient interface {
	Post(ctx context.Context, url string, payload interface{}, headers map[string]string) error
}

// WebhookHTTPClient implements HTTPClient for webhook operations
type WebhookHTTPClient struct {
	client  *http.Client
	timeout time.Duration
}

// NewWebhookHTTPClient creates a new webhook HTTP client
func NewWebhookHTTPClient(timeout time.Duration) *WebhookHTTPClient {
	return &WebhookHTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// Post sends a POST request with JSON payload
func (c *WebhookHTTPClient) Post(ctx context.Context, url string, payload interface{}, headers map[string]string) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook request failed with status %d", resp.StatusCode)
	}

	return nil
}
