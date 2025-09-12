package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/database/models"
	"zpmeow/internal/shared/types"

	"github.com/jmoiron/sqlx"
)

// PostgresSessionRepository implements session.SessionRepository using PostgreSQL
type PostgresSessionRepository struct {
	db *sqlx.DB
}

// NewPostgresSessionRepository creates a new PostgreSQL session repository
func NewPostgresSessionRepository(db *sqlx.DB) session.SessionRepository {
	return &PostgresSessionRepository{db: db}
}

// Create implements session.SessionRepository
func (r *PostgresSessionRepository) Create(ctx context.Context, sess *session.Session) error {
	model := models.FromDomain(sess)

	query := `
		INSERT INTO sessions (id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, created_at, updated_at)
		VALUES (:id, :name, :device_jid, :status, :qr_code, :proxy_url, :webhook_url, :webhook_events, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetByID implements session.SessionRepository
func (r *PostgresSessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	var model models.SessionModel

	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, created_at, updated_at
		FROM sessions
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}

	return model.ToDomain(), nil
}

// GetByName implements session.SessionRepository
func (r *PostgresSessionRepository) GetByName(ctx context.Context, name string) (*session.Session, error) {
	var model models.SessionModel

	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, created_at, updated_at
		FROM sessions
		WHERE name = $1
	`

	err := r.db.GetContext(ctx, &model, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by name: %w", err)
	}

	return model.ToDomain(), nil
}

// GetAll implements session.SessionRepository
func (r *PostgresSessionRepository) GetAll(ctx context.Context) ([]*session.Session, error) {
	var models []models.SessionModel

	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, created_at, updated_at
		FROM sessions
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sessions: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sessions[i] = model.ToDomain()
	}

	return sessions, nil
}

// Update implements session.SessionRepository
func (r *PostgresSessionRepository) Update(ctx context.Context, sess *session.Session) error {
	sess.UpdatedAt = time.Now()
	model := models.FromDomain(sess)

	query := `
		UPDATE sessions
		SET name = :name, device_jid = :device_jid, status = :status, qr_code = :qr_code,
		    proxy_url = :proxy_url, webhook_url = :webhook_url, webhook_events = :webhook_events, updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

// Delete implements session.SessionRepository
func (r *PostgresSessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

// Exists implements session.SessionRepository
func (r *PostgresSessionRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM sessions WHERE id = $1`

	err := r.db.GetContext(ctx, &count, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if session exists: %w", err)
	}

	return count > 0, nil
}

// GetByStatus implements session.SessionRepository
func (r *PostgresSessionRepository) GetByStatus(ctx context.Context, status types.Status) ([]*session.Session, error) {
	var models []models.SessionModel

	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, created_at, updated_at
		FROM sessions
		WHERE status = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by status: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sessions[i] = model.ToDomain()
	}

	return sessions, nil
}
