package webhook

import (
	"context"
)

// Repository defines the interface for webhook persistence
type Repository interface {
	// Create creates a new webhook configuration
	Create(ctx context.Context, webhook *WebhookConfiguration) error

	// GetBySessionID retrieves webhook configuration by session ID
	GetBySessionID(ctx context.Context, sessionID string) (*WebhookConfiguration, error)

	// Update updates an existing webhook configuration
	Update(ctx context.Context, webhook *WebhookConfiguration) error

	// Delete removes webhook configuration by session ID
	Delete(ctx context.Context, sessionID string) error

	// Exists checks if webhook configuration exists for session
	Exists(ctx context.Context, sessionID string) (bool, error)

	// GetActive retrieves all active webhook configurations
	GetActive(ctx context.Context) ([]*WebhookConfiguration, error)

	// GetSubscribedToEvent retrieves webhooks subscribed to a specific event
	GetSubscribedToEvent(ctx context.Context, eventType EventType) ([]*WebhookConfiguration, error)

	// List retrieves webhook configurations with pagination
	List(ctx context.Context, limit, offset int) ([]*WebhookConfiguration, int, error)
}
