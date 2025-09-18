package session

import (
	"context"
)

type SessionManager interface {
	CreateSession(ctx context.Context, name string) (*Session, error)
	GetSession(ctx context.Context, idOrName string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
	ListSessions(ctx context.Context, limit, offset int) ([]*Session, error)
}
