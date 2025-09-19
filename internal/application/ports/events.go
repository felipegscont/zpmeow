package ports

import (
	"context"

	"zpmeow/internal/domain/common"
)

// EventPublisher defines the contract for publishing domain events
type EventPublisher interface {
	// Publish publishes a domain event
	Publish(ctx context.Context, event common.DomainEvent) error

	// PublishBatch publishes multiple domain events
	PublishBatch(ctx context.Context, events []common.DomainEvent) error
}

// EventHandler defines the contract for handling domain events
type EventHandler interface {
	// Handle processes a domain event
	Handle(ctx context.Context, event common.DomainEvent) error

	// CanHandle checks if this handler can process the given event type
	CanHandle(eventType string) bool
}

// EventBus defines the contract for event bus operations
type EventBus interface {
	// Subscribe registers an event handler for specific event types
	Subscribe(eventType string, handler EventHandler) error

	// Unsubscribe removes an event handler
	Unsubscribe(eventType string, handler EventHandler) error

	// Publish publishes an event to all registered handlers
	Publish(ctx context.Context, event common.DomainEvent) error
}
