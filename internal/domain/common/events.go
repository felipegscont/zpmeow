package common

import (
	"time"
)

// DomainEvent represents a domain event that occurred
type DomainEvent interface {
	// EventID returns the unique identifier of the event
	EventID() string

	// EventType returns the type of the event
	EventType() string

	// AggregateID returns the ID of the aggregate that generated the event
	AggregateID() string

	// OccurredAt returns when the event occurred
	OccurredAt() time.Time

	// EventData returns the event-specific data
	EventData() interface{}
}

// BaseDomainEvent provides common functionality for domain events
type BaseDomainEvent struct {
	eventID     string
	eventType   string
	aggregateID string
	occurredAt  time.Time
	data        interface{}
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType, aggregateID string, data interface{}) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     GenerateID().Value(),
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
		data:        data,
	}
}

// EventID returns the unique identifier of the event
func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

// EventType returns the type of the event
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the ID of the aggregate that generated the event
func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventData returns the event-specific data
func (e BaseDomainEvent) EventData() interface{} {
	return e.data
}

// AggregateRoot represents the base for all aggregate roots
type AggregateRoot struct {
	id     ID
	events []DomainEvent
}

// NewAggregateRoot creates a new aggregate root
func NewAggregateRoot(id ID) AggregateRoot {
	return AggregateRoot{
		id:     id,
		events: make([]DomainEvent, 0),
	}
}

// ID returns the aggregate root ID
func (ar *AggregateRoot) ID() ID {
	return ar.id
}

// AddEvent adds a domain event to the aggregate
func (ar *AggregateRoot) AddEvent(event DomainEvent) {
	ar.events = append(ar.events, event)
}

// GetEvents returns all uncommitted events
func (ar *AggregateRoot) GetEvents() []DomainEvent {
	return ar.events
}

// ClearEvents clears all events (typically called after persistence)
func (ar *AggregateRoot) ClearEvents() {
	ar.events = make([]DomainEvent, 0)
}

// HasEvents checks if there are uncommitted events
func (ar *AggregateRoot) HasEvents() bool {
	return len(ar.events) > 0
}
