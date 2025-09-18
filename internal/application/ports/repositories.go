package ports

import (
	"context"

	"meow/internal/domain/session"
)

// SessionRepository defines the contract for session persistence
// This is a port that will be implemented by the infrastructure layer
type SessionRepository interface {
	// Create persists a new session
	Create(ctx context.Context, session *session.Session) error

	// GetByID retrieves a session by its ID
	GetByID(ctx context.Context, id string) (*session.Session, error)

	// GetByName retrieves a session by its name
	GetByName(ctx context.Context, name string) (*session.Session, error)

	// GetByApiKey retrieves a session by its API key
	GetByApiKey(ctx context.Context, apiKey string) (*session.Session, error)

	// GetAll retrieves all sessions
	GetAll(ctx context.Context) ([]*session.Session, error)

	// Update updates an existing session
	Update(ctx context.Context, session *session.Session) error

	// Delete removes a session
	Delete(ctx context.Context, id string) error

	// Exists checks if a session exists by ID or name
	Exists(ctx context.Context, identifier string) (bool, error)
}

// UnitOfWork defines transaction management
type UnitOfWork interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (Transaction, error)
}

// Transaction represents a database transaction
type Transaction interface {
	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error

	// SessionRepository returns the session repository within this transaction
	SessionRepository() SessionRepository
}
