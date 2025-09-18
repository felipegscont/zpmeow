package webhook

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, webhook *WebhookConfiguration) error

	GetBySessionID(ctx context.Context, sessionID string) (*WebhookConfiguration, error)

	Update(ctx context.Context, webhook *WebhookConfiguration) error

	Delete(ctx context.Context, sessionID string) error

	Exists(ctx context.Context, sessionID string) (bool, error)

	GetActive(ctx context.Context) ([]*WebhookConfiguration, error)

	GetSubscribedToEvent(ctx context.Context, eventType EventType) ([]*WebhookConfiguration, error)

	List(ctx context.Context, limit, offset int) ([]*WebhookConfiguration, int, error)
}
