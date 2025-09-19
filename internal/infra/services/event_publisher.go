package services

import (
	"context"
	"fmt"

	"zpmeow/internal/application/ports"
	"zpmeow/internal/domain/common"
)

type LoggingEventPublisher struct {
	logger ports.Logger
}

func NewLoggingEventPublisher(logger ports.Logger) *LoggingEventPublisher {
	return &LoggingEventPublisher{
		logger: logger,
	}
}

func (p *LoggingEventPublisher) Publish(ctx context.Context, event common.DomainEvent) error {
	p.logger.Info(ctx, "Domain event published",
		"eventID", event.EventID(),
		"eventType", event.EventType(),
		"aggregateID", event.AggregateID(),
		"occurredAt", event.OccurredAt(),
	)
	return nil
}

func (p *LoggingEventPublisher) PublishBatch(ctx context.Context, events []common.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	p.logger.Info(ctx, fmt.Sprintf("Publishing batch of %d domain events", len(events)))

	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.EventID(), err)
		}
	}

	p.logger.Info(ctx, fmt.Sprintf("Successfully published batch of %d domain events", len(events)))
	return nil
}
