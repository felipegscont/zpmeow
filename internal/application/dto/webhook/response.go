package webhook

// WebhookResponse represents the response from webhook operations
type WebhookResponse struct {
	Webhook   string   `json:"webhook" example:"https://example.com/webhook"`
	Events    []string `json:"events,omitempty" example:"message,status"`
	Active    bool     `json:"active,omitempty" example:"true"`
	Subscribe []string `json:"subscribe,omitempty" example:"message,status"`
}
