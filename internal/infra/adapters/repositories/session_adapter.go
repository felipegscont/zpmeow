package repositories

import (
	"context"
	"fmt"

	"meow/internal/application/ports"
	"meow/internal/domain/session"
	"meow/internal/infra/database/repository"

	"github.com/jmoiron/sqlx"
)

// SessionRepositoryAdapter adapts the existing domain repository to implement application ports
// This is a temporary solution to bridge the gap between domain and application interfaces
type SessionRepositoryAdapter struct {
	domainRepo session.Repository
	logger     ports.Logger
}

// NewSessionRepositoryAdapter creates a new session repository adapter
func NewSessionRepositoryAdapter(db *sqlx.DB, logger ports.Logger) ports.SessionRepository {
	// Use the existing domain repository implementation
	domainRepo := repository.NewPostgresRepo(db)

	return &SessionRepositoryAdapter{
		domainRepo: domainRepo,
		logger:     logger,
	}
}

// Create implements ports.SessionRepository
func (a *SessionRepositoryAdapter) Create(ctx context.Context, sessionEntity *session.Session) error {
	if a.logger != nil {
		a.logger.Debug(ctx, "Creating session via adapter", "sessionID", sessionEntity.SessionID().Value())
	}

	// Delegate to domain repository
	err := a.domainRepo.Create(ctx, sessionEntity)
	if err != nil {
		if a.logger != nil {
			a.logger.Error(ctx, "Failed to create session", "error", err)
		}
		return fmt.Errorf("failed to create session: %w", err)
	}

	if a.logger != nil {
		a.logger.Debug(ctx, "Session created successfully", "sessionID", sessionEntity.SessionID().Value())
	}
	return nil
}

// GetByID implements ports.SessionRepository
func (a *SessionRepositoryAdapter) GetByID(ctx context.Context, id string) (*session.Session, error) {
	if a.logger != nil {
		a.logger.Debug(ctx, "Getting session by ID", "sessionID", id)
	}

	sessionEntity, err := a.domainRepo.GetByID(ctx, id)
	if err != nil {
		if a.logger != nil {
			a.logger.Error(ctx, "Failed to get session by ID", "sessionID", id, "error", err)
		}
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}

	return sessionEntity, nil
}

// GetByName implements ports.SessionRepository
func (a *SessionRepositoryAdapter) GetByName(ctx context.Context, name string) (*session.Session, error) {
	if a.logger != nil {
		a.logger.Debug(ctx, "Getting session by name", "name", name)
	}

	sessionEntity, err := a.domainRepo.GetByName(ctx, name)
	if err != nil {
		if a.logger != nil {
			a.logger.Error(ctx, "Failed to get session by name", "name", name, "error", err)
		}
		return nil, fmt.Errorf("failed to get session by name: %w", err)
	}

	return sessionEntity, nil
}

// GetByApiKey implements ports.SessionRepository
func (a *SessionRepositoryAdapter) GetByApiKey(ctx context.Context, apiKey string) (*session.Session, error) {
	if a.logger != nil {
		a.logger.Debug(ctx, "Getting session by API key")
	}

	sessionEntity, err := a.domainRepo.GetByApiKey(ctx, apiKey)
	if err != nil {
		if a.logger != nil {
			a.logger.Error(ctx, "Failed to get session by API key", "error", err)
		}
		return nil, fmt.Errorf("failed to get session by API key: %w", err)
	}

	return sessionEntity, nil
}

// GetAll implements ports.SessionRepository
func (a *SessionRepositoryAdapter) GetAll(ctx context.Context) ([]*session.Session, error) {
	if a.logger != nil {
		a.logger.Debug(ctx, "Getting all sessions")
	}

	sessions, err := a.domainRepo.GetAll(ctx)
	if err != nil {
		if a.logger != nil {
			a.logger.Error(ctx, "Failed to get all sessions", "error", err)
		}
		return nil, fmt.Errorf("failed to get all sessions: %w", err)
	}

	if a.logger != nil {
		a.logger.Debug(ctx, "Retrieved all sessions", "count", len(sessions))
	}
	return sessions, nil
}

// Update implements ports.SessionRepository
func (a *SessionRepositoryAdapter) Update(ctx context.Context, sessionEntity *session.Session) error {
	if a.logger != nil {
		a.logger.Debug(ctx, "Updating session", "sessionID", sessionEntity.SessionID().Value())
	}

	err := a.domainRepo.Update(ctx, sessionEntity)
	if err != nil {
		if a.logger != nil {
			a.logger.Error(ctx, "Failed to update session", "sessionID", sessionEntity.SessionID().Value(), "error", err)
		}
		return fmt.Errorf("failed to update session: %w", err)
	}

	if a.logger != nil {
		a.logger.Debug(ctx, "Session updated successfully", "sessionID", sessionEntity.SessionID().Value())
	}
	return nil
}

// Delete implements ports.SessionRepository
func (a *SessionRepositoryAdapter) Delete(ctx context.Context, id string) error {
	if a.logger != nil {
		a.logger.Debug(ctx, "Deleting session", "sessionID", id)
	}

	err := a.domainRepo.Delete(ctx, id)
	if err != nil {
		if a.logger != nil {
			a.logger.Error(ctx, "Failed to delete session", "sessionID", id, "error", err)
		}
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if a.logger != nil {
		a.logger.Debug(ctx, "Session deleted successfully", "sessionID", id)
	}
	return nil
}

// Exists implements ports.SessionRepository
// This method was missing in the domain repository, so we implement it here
func (a *SessionRepositoryAdapter) Exists(ctx context.Context, identifier string) (bool, error) {
	if a.logger != nil {
		a.logger.Debug(ctx, "Checking if session exists", "identifier", identifier)
	}

	// Try to get by ID first
	_, err := a.domainRepo.GetByID(ctx, identifier)
	if err == nil {
		return true, nil
	}

	// If not found by ID, try by name
	_, err = a.domainRepo.GetByName(ctx, identifier)
	if err == nil {
		return true, nil
	}

	// If both fail, session doesn't exist
	return false, nil
}
