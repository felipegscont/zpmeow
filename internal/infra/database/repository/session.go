package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	sessiondomain "meow/internal/domain/session"
	"meow/internal/infra/persistence/postgres"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) sessiondomain.Repository {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, sessionEntity *sessiondomain.Session) error {
	// ✅ CORREÇÃO: Não modificar a entidade de domínio!
	// Se o ID estiver vazio, gerar um novo ID apenas para persistência
	sessionID := sessionEntity.SessionID().Value()
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	eventsJSON := []byte("[]")

	// ✅ CORREÇÃO: Usar timestamps da entidade ou gerar apenas para persistência
	now := time.Now()
	createdAt := now
	updatedAt := now

	// Se a entidade já tem timestamps, usar os dela
	if !sessionEntity.CreatedAt().Value().IsZero() {
		createdAt = sessionEntity.CreatedAt().Value()
	}
	if !sessionEntity.UpdatedAt().Value().IsZero() {
		updatedAt = sessionEntity.UpdatedAt().Value()
	}

	query := `
		INSERT INTO sessions (id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	// ✅ CORREÇÃO: Usar métodos corretos da entidade
	isConnected := sessionEntity.Status() == sessiondomain.StatusConnected && sessionEntity.WaJID().Value() != ""

	_, err := r.db.ExecContext(ctx, query,
		sessionID, // usar o ID gerado ou existente
		sessionEntity.Name().Value(),
		sessionEntity.WaJID().Value(),
		string(sessionEntity.Status()),
		sessionEntity.QRCode().Value(),
		sessionEntity.ProxyConfiguration().Value(), // ✅ CORREÇÃO: usar método correto
		sessionEntity.WebhookEndpoint().Value(),    // ✅ CORREÇÃO: usar método correto
		string(eventsJSON),
		isConnected,
		sessionEntity.ApiKey().Value(),
		createdAt,
		updatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique_session_name") {
			return fmt.Errorf("session already exists")
		}
		return fmt.Errorf("failed to create session: %w", err)
	}

	// ✅ CORREÇÃO: Não modificar a entidade de domínio!
	// A entidade deve ser criada com o ID correto desde o início
	// Se precisar do ID gerado, isso deve ser tratado na camada de aplicação

	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*sessiondomain.Session, error) {
	var model postgres.SessionModel
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

func (r *PostgresRepo) GetByName(ctx context.Context, name string) (*sessiondomain.Session, error) {
	var model postgres.SessionModel
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

func (r *PostgresRepo) GetAll(ctx context.Context) ([]*sessiondomain.Session, error) {
	var models []postgres.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sessions: %w", err)
	}

	sessions := make([]*sessiondomain.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *PostgresRepo) Update(ctx context.Context, session *sessiondomain.Session) error {
	// ✅ CORREÇÃO: Não modificar a entidade de domínio!
	eventsJSON := []byte("[]")
	updatedAt := time.Now()

	query := `
		UPDATE sessions
		SET name = $2, device_jid = $3, status = $4, qr_code = $5, proxy_url = $6,
		    webhook_url = $7, webhook_events = $8, connected = $9, apikey = $10, updated_at = $11
		WHERE id = $1
	`

	// ✅ CORREÇÃO: Usar métodos corretos da entidade
	isConnected := session.Status() == sessiondomain.StatusConnected && session.WaJID().Value() != ""

	result, err := r.db.ExecContext(ctx, query,
		session.SessionID().Value(),
		session.Name().Value(),
		session.WaJID().Value(),
		string(session.Status()),
		session.QRCode().Value(),
		session.ProxyConfiguration().Value(),
		session.WebhookEndpoint().Value(),
		string(eventsJSON),
		isConnected,
		session.ApiKey().Value(),
		updatedAt,
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

func (r *PostgresRepo) List(ctx context.Context, limit, offset int, status string) ([]*sessiondomain.Session, int, error) {
	var models []postgres.SessionModel
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

	sessions := make([]*sessiondomain.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, 0, err
		}
		sessions[i] = session
	}

	return sessions, totalCount, nil
}

func (r *PostgresRepo) GetActive(ctx context.Context) ([]*sessiondomain.Session, error) {
	var models []postgres.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE device_jid IS NOT NULL AND device_jid != '' ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions with credentials: %w", err)
	}

	sessions := make([]*sessiondomain.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *PostgresRepo) GetInactive(ctx context.Context) ([]*sessiondomain.Session, error) {
	var models []postgres.SessionModel
	query := `
		SELECT id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, connected, apikey, created_at, updated_at
		FROM sessions WHERE status != $1 ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &models, query, string("connected"))
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive sessions: %w", err)
	}

	sessions := make([]*sessiondomain.Session, len(models))
	for i, model := range models {
		session, err := r.modelToDomain(&model)
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *PostgresRepo) GetByApiKey(ctx context.Context, apiKey string) (*sessiondomain.Session, error) {
	var model postgres.SessionModel
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

func (r *PostgresRepo) GetByDeviceJID(ctx context.Context, deviceJID string) (*sessiondomain.Session, error) {
	if deviceJID == "" {
		return nil, fmt.Errorf("device JID cannot be empty")
	}

	var model postgres.SessionModel
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

func (r *PostgresRepo) modelToDomain(model *postgres.SessionModel) (*sessiondomain.Session, error) {

	// status := sessiondomain.Status(model.Status) // não usado por enquanto

	sessionID, err := sessiondomain.NewSessionID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session ID: %w", err)
	}

	sessionName, err := sessiondomain.NewSessionName(model.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create session name: %w", err)
	}

	// Reconstruct proxy configuration
	var proxyURL sessiondomain.ProxyConfiguration
	if model.ProxyURL != "" {
		proxy, err := sessiondomain.NewProxyConfiguration(model.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create proxy configuration: %w", err)
		}
		proxyURL = proxy
	}

	// Reconstruct device JID
	var waJID sessiondomain.WaJID
	if model.DeviceJID != "" {
		jid, err := sessiondomain.NewWaJID(model.DeviceJID)
		if err != nil {
			return nil, fmt.Errorf("failed to create device JID: %w", err)
		}
		waJID = jid
	}

	// Reconstruct QR code
	var qrCode sessiondomain.QRCode
	if model.QRCode != "" {
		qr, err := sessiondomain.NewQRCode(model.QRCode)
		if err != nil {
			return nil, fmt.Errorf("failed to create QR code: %w", err)
		}
		qrCode = qr
	}

	// Reconstruct API key
	var apiKey sessiondomain.ApiKey
	if model.ApiKey != "" {
		key, err := sessiondomain.NewApiKey(model.ApiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create API key: %w", err)
		}
		apiKey = key
	}

	// Parse webhook events
	var events []string
	if model.Events != "" && model.Events != "[]" {
		if err := json.Unmarshal([]byte(model.Events), &events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal webhook events: %w", err)
		}
	}

	// Create session entity with basic information
	sessionEntity, err := sessiondomain.NewSession(sessionID.Value(), sessionName.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to create session entity: %w", err)
	}

	// Reconstruct complete state
	// Note: Status is set during NewSession, we need to update it if different
	if model.Status != string(sessiondomain.StatusDisconnected) {
		switch sessiondomain.Status(model.Status) {
		case sessiondomain.StatusConnecting:
			if err := sessionEntity.Connect(); err != nil {
				return nil, fmt.Errorf("failed to set connecting status: %w", err)
			}
		case sessiondomain.StatusConnected:
			if err := sessionEntity.Connect(); err != nil {
				return nil, fmt.Errorf("failed to set connecting status: %w", err)
			}
			if err := sessionEntity.SetConnected(); err != nil {
				return nil, fmt.Errorf("failed to set connected status: %w", err)
			}
		case sessiondomain.StatusError:
			sessionEntity.SetError("Restored from database")
		}
	}

	if !proxyURL.IsEmpty() {
		if err := sessionEntity.SetProxyConfiguration(proxyURL.Value()); err != nil {
			return nil, fmt.Errorf("failed to set proxy URL: %w", err)
		}
	}

	if !waJID.IsEmpty() {
		if err := sessionEntity.Authenticate(waJID.Value()); err != nil {
			return nil, fmt.Errorf("failed to set device JID: %w", err)
		}
	}

	if !qrCode.IsEmpty() {
		if err := sessionEntity.SetQRCode(qrCode.Value()); err != nil {
			return nil, fmt.Errorf("failed to set QR code: %w", err)
		}
	}

	if !apiKey.IsEmpty() {
		if err := sessionEntity.SetApiKey(apiKey.Value()); err != nil {
			return nil, fmt.Errorf("failed to set API key: %w", err)
		}
	}

	if model.WebhookURL != "" {
		if err := sessionEntity.SetWebhookEndpoint(model.WebhookURL); err != nil {
			return nil, fmt.Errorf("failed to set webhook URL: %w", err)
		}
	}

	// Note: Webhook events are not stored in the domain entity directly
	// They would be handled by the application layer
	_ = events // Suppress unused variable warning

	// Note: Timestamps are managed internally by the domain entity
	// We don't set them directly from the database

	return sessionEntity, nil
}

// generateAPIKey generates a new API key in the format: cPXeznTF44e1d8BQ57SJoGxtEyffdazE
func (r *PostgresRepo) generateAPIKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 32

	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
