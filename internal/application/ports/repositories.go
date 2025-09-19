package ports

import (
	"context"

	"zpmeow/internal/domain/session"
)

type SessionRepository interface {
	Create(ctx context.Context, session *session.Session) error

	CreateWithGeneratedID(ctx context.Context, session *session.Session) (string, error)

	GetByID(ctx context.Context, id string) (*session.Session, error)

	GetByName(ctx context.Context, name string) (*session.Session, error)

	GetByApiKey(ctx context.Context, apiKey string) (*session.Session, error)

	GetAll(ctx context.Context) ([]*session.Session, error)

	Update(ctx context.Context, session *session.Session) error

	Delete(ctx context.Context, id string) error

	Exists(ctx context.Context, identifier string) (bool, error)
}

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error

	Rollback() error

	SessionRepository() SessionRepository
}
