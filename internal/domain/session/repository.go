package session

import (
	"context"
	"zpmeow/internal/shared/types"
)



type SessionRepository interface {

	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	GetByName(ctx context.Context, name string) (*Session, error)
	GetAll(ctx context.Context) ([]*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error


	Exists(ctx context.Context, id string) (bool, error)
	GetByStatus(ctx context.Context, status types.Status) ([]*Session, error)
}
