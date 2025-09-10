package session

import (
	"context"
	"zpmeow/internal/types"
)

// SessionRepository defines the contract for session persistence
// Following KISS principle - only essential operations
type SessionRepository interface {
	// Core CRUD operations
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	GetAll(ctx context.Context) ([]*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error

	// Essential queries
	Exists(ctx context.Context, id string) (bool, error)
	GetByStatus(ctx context.Context, status types.Status) ([]*Session, error)
}
