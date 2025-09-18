package session

import (
	"meow/internal/domain/common"
)

// Event types for session domain
const (
	SessionCreatedEventType              = "session.created"
	SessionConnectedEventType            = "session.connected"
	SessionDisconnectedEventType         = "session.disconnected"
	SessionAuthenticatedEventType        = "session.authenticated"
	SessionConfigurationChangedEventType = "session.configuration_changed"
	SessionDeletedEventType              = "session.deleted"
	SessionErrorEventType                = "session.error"
)

// SessionCreatedEvent represents a session creation event
type SessionCreatedEvent struct {
	common.BaseDomainEvent
}

// NewSessionCreatedEvent creates a new session created event
func NewSessionCreatedEvent(sessionID string, sessionName string) SessionCreatedEvent {
	data := map[string]interface{}{
		"session_id":   sessionID,
		"session_name": sessionName,
	}

	return SessionCreatedEvent{
		BaseDomainEvent: common.NewBaseDomainEvent(
			SessionCreatedEventType,
			sessionID,
			data,
		),
	}
}

// SessionConnectedEvent represents a session connection event
type SessionConnectedEvent struct {
	common.BaseDomainEvent
}

// NewSessionConnectedEvent creates a new session connected event
func NewSessionConnectedEvent(sessionID string, waJID string) SessionConnectedEvent {
	data := map[string]interface{}{
		"session_id": sessionID,
		"wa_jid":     waJID,
	}

	return SessionConnectedEvent{
		BaseDomainEvent: common.NewBaseDomainEvent(
			SessionConnectedEventType,
			sessionID,
			data,
		),
	}
}

// SessionDisconnectedEvent represents a session disconnection event
type SessionDisconnectedEvent struct {
	common.BaseDomainEvent
}

// NewSessionDisconnectedEvent creates a new session disconnected event
func NewSessionDisconnectedEvent(sessionID string, reason string) SessionDisconnectedEvent {
	data := map[string]interface{}{
		"session_id": sessionID,
		"reason":     reason,
	}

	return SessionDisconnectedEvent{
		BaseDomainEvent: common.NewBaseDomainEvent(
			SessionDisconnectedEventType,
			sessionID,
			data,
		),
	}
}

// SessionAuthenticatedEvent represents a session authentication event
type SessionAuthenticatedEvent struct {
	common.BaseDomainEvent
}

// NewSessionAuthenticatedEvent creates a new session authenticated event
func NewSessionAuthenticatedEvent(sessionID string, waJID string) SessionAuthenticatedEvent {
	data := map[string]interface{}{
		"session_id": sessionID,
		"wa_jid":     waJID,
	}

	return SessionAuthenticatedEvent{
		BaseDomainEvent: common.NewBaseDomainEvent(
			SessionAuthenticatedEventType,
			sessionID,
			data,
		),
	}
}

// SessionConfigurationChangedEvent represents a session configuration change event
type SessionConfigurationChangedEvent struct {
	common.BaseDomainEvent
}

// NewSessionConfigurationChangedEvent creates a new session configuration changed event
func NewSessionConfigurationChangedEvent(sessionID string, changes map[string]interface{}) SessionConfigurationChangedEvent {
	data := map[string]interface{}{
		"session_id": sessionID,
		"changes":    changes,
	}

	return SessionConfigurationChangedEvent{
		BaseDomainEvent: common.NewBaseDomainEvent(
			SessionConfigurationChangedEventType,
			sessionID,
			data,
		),
	}
}

// SessionDeletedEvent represents a session deletion event
type SessionDeletedEvent struct {
	common.BaseDomainEvent
}

// NewSessionDeletedEvent creates a new session deleted event
func NewSessionDeletedEvent(sessionID string, sessionName string) SessionDeletedEvent {
	data := map[string]interface{}{
		"session_id":   sessionID,
		"session_name": sessionName,
	}

	return SessionDeletedEvent{
		BaseDomainEvent: common.NewBaseDomainEvent(
			SessionDeletedEventType,
			sessionID,
			data,
		),
	}
}

// SessionErrorEvent represents a session error event
type SessionErrorEvent struct {
	common.BaseDomainEvent
}

// NewSessionErrorEvent creates a new session error event
func NewSessionErrorEvent(sessionID string, errorMessage string, errorCode string) SessionErrorEvent {
	data := map[string]interface{}{
		"session_id":    sessionID,
		"error_message": errorMessage,
		"error_code":    errorCode,
	}

	return SessionErrorEvent{
		BaseDomainEvent: common.NewBaseDomainEvent(
			SessionErrorEventType,
			sessionID,
			data,
		),
	}
}
