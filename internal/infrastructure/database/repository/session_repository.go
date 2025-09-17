package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"meow/internal/domain/session"
	"meow/internal/infrastructure/database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) session.Repository {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, sessionEntity *session.Session) error {
	if sessionEntity.ID.IsEmpty() {
		newID, err := session.NewSessionID(uuid.New().String())
		if err != nil {
			return fmt.Errorf("failed to create session ID: %w", err)
		}
		sessionEntity.ID = newID
	}

	eventsJSON, err := json.Marshal(sessionEntity.Events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	now := time.Now()
	sessionEntity.CreatedAt = now
	sessionEntity.UpdatedAt = now

	query := `
		INSERT INTO sessions (id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	// Ensure consistency: connected field should only be true if status is "connected" AND device_jid is not empty
	isConnected := string(sessionEntity.Status) == "connected" && sessionEntity.WaJID != ""

	_, err = r.db.ExecContext(ctx, query,
		sessionEntity.ID.Value(),
		sessionEntity.Name.Value(),
		sessionEntity.WaJID,
		string(sessionEntity.Status),
		sessionEntity.QRCode,
		sessionEntity.ProxyURL.Value(),
		sessionEntity.WebhookURL,
		string(eventsJSON),
		isConnected, // connected field with validation
		sessionEntity.ApiKey,
		sessionEntity.CreatedAt,
		sessionEntity.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique_session_name") {
			return fmt.Errorf("session already exists")
		}
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*session.Session, error) {
	var model database.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE id = $1
	`

	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}

	return r.modelToDomain(&model)
}

func (r *PostgresRepo) GetByName(ctx context.Context, name string) (*session.Session, error) {
	var model database.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE name = $1
	`

	err := r.db.GetContext(ctx, &model, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session by name: %w", err)
	}

	return r.modelToDomain(&model)
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]*session.Session, error) {
	var models []database.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sessions: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *PostgresRepo) Update(ctx context.Context, session *session.Session) error {
	eventsJSON, err := json.Marshal(session.Events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	session.UpdatedAt = time.Now()

	query := `
		UPDATE sessions
		SET name = $2, device_jid = $3, status = $4, qr_code = $5, proxy_url = $6,
		    webhook_url = $7, webhook_events = $8, connected = $9, apikey = $10, updated_at = $11
		WHERE id = $1
	`

	// Ensure consistency: connected field should only be true if status is "connected" AND device_jid is not empty
	isConnected := string(session.Status) == "connected" && session.WaJID != ""

	result, err := r.db.ExecContext(ctx, query,
		session.ID.Value(),
		session.Name.Value(),
		session.WaJID,
		string(session.Status),
		session.QRCode,
		session.ProxyURL.Value(),
		session.WebhookURL,
		string(eventsJSON),
		isConnected, // connected field with validation
		session.ApiKey,
		session.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique_session_name") {
			return fmt.Errorf("session already exists")
		}
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id string) error {
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
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *PostgresRepo) List(ctx context.Context, limit, offset int, status string) ([]*session.Session, int, error) {
	var models []database.SessionModel
	var totalCount int

	countQuery := `SELECT COUNT(*) FROM sessions`
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		countQuery += fmt.Sprintf(" WHERE status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions
	`

	if status != "" {
		query += fmt.Sprintf(" WHERE status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	err = r.db.SelectContext(ctx, &models, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sessions: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, 0, err
		}
		sessions[i] = session
	}

	return sessions, totalCount, nil
}

func (r *PostgresRepo) GetActive(ctx context.Context) ([]*session.Session, error) {
	var models []database.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE device_jid IS NOT NULL AND device_jid != '' ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions with credentials: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *PostgresRepo) GetInactive(ctx context.Context) ([]*session.Session, error) {
	var models []database.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE status != $1 ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query, string("connected"))
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive sessions: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *PostgresRepo) GetByApiKey(ctx context.Context, apiKey string) (*session.Session, error) {
	var model database.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE apikey = $1
	`

	err := r.db.GetContext(ctx, &model, query, apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session by API key: %w", err)
	}

	return r.modelToDomain(&model)
}

// GetByDeviceJID gets a session by device JID
func (r *PostgresRepo) GetByDeviceJID(ctx context.Context, deviceJID string) (*session.Session, error) {
	if deviceJID == "" {
		return nil, fmt.Errorf("device JID cannot be empty")
	}

	var model database.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE device_jid = $1
	`

	err := r.db.GetContext(ctx, &model, query, deviceJID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session by device JID: %w", err)
	}

	return r.modelToDomain(&model)
}

// ValidateDeviceUniqueness validates that a device JID is not already in use by another session
func (r *PostgresRepo) ValidateDeviceUniqueness(ctx context.Context, sessionID, deviceJID string) error {
	if deviceJID == "" {
		return nil // Empty device JID is allowed (not connected yet)
	}

	var count int
	query := `
		SELECT COUNT(*) FROM sessions
		WHERE device_jid = $1 AND id != $2
	`

	err := r.db.GetContext(ctx, &count, query, deviceJID, sessionID)
	if err != nil {
		return fmt.Errorf("failed to validate device uniqueness: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("device JID %s is already in use by another session", deviceJID)
	}

	return nil
}

func (r *PostgresRepo) modelToDomain(model *database.SessionModel) (*session.Session, error) {
	var events []string
	if model.Events != "" && model.Events != "null" {
		err := json.Unmarshal([]byte(model.Events), &events)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}
	}

	status := session.Status(model.Status)

	// Create value objects
	sessionID, err := session.NewSessionID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session ID: %w", err)
	}

	sessionName, err := session.NewSessionName(model.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create session name: %w", err)
	}

	proxyURL, err := session.NewProxyURL(model.ProxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy URL: %w", err)
	}

	return &session.Session{
		ID:         sessionID,
		Name:       sessionName,
		WaJID:      model.DeviceJID,
		Status:     status,
		QRCode:     model.QRCode,
		ProxyURL:   proxyURL,
		WebhookURL: model.WebhookURL,
		Events:     events,
		ApiKey:     model.ApiKey,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}, nil
}
