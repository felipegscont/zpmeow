package session

import "context"

// SessionRepository defines the contract for session persistence
type SessionRepository interface {
	// Save creates or updates a session
	Save(ctx context.Context, session *Session) error
	
	// FindByID retrieves a session by its ID
	FindByID(ctx context.Context, id string) (*Session, error)
	
	// FindAll retrieves all sessions
	FindAll(ctx context.Context) ([]*Session, error)
	
	// Update updates an existing session
	Update(ctx context.Context, session *Session) error
	
	// Delete removes a session by ID
	Delete(ctx context.Context, id string) error
	
	// Exists checks if a session exists by ID
	Exists(ctx context.Context, id string) (bool, error)
	
	// FindByName retrieves sessions by name (for search functionality)
	FindByName(ctx context.Context, name string) ([]*Session, error)
	
	// FindByStatus retrieves sessions by status
	FindByStatus(ctx context.Context, status string) ([]*Session, error)
}
