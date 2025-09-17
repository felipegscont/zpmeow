package session

import (
	"context"
)

// SessionManager defines the interface for session management operations
type SessionManager interface {
	CreateSession(ctx context.Context, name string) (*Session, error)
	GetSession(ctx context.Context, idOrName string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
	ListSessions(ctx context.Context, limit, offset int) ([]*Session, error)
	// Note: Technical operations moved to application layer
	// Domain interfaces should only contain business concepts
}

// Note: Technical interfaces like MessageSender, WebhookManager moved to application layer
// Domain interfaces should only contain pure business concepts
