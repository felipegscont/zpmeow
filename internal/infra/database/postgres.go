package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"zpmeow/internal/config"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/types"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


type PostgresSessionRepository struct {
	db *sqlx.DB
}


func NewPostgresSessionRepository(db *sqlx.DB) session.SessionRepository {
	return &PostgresSessionRepository{db: db}
}


type sessionModel struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	WhatsAppJID   string    `db:"device_jid"`
	Status        string    `db:"status"`
	QRCode        string    `db:"qr_code"`
	ProxyURL      string    `db:"proxy_url"`
	WebhookURL    string    `db:"webhook_url"`
	WebhookEvents string    `db:"webhook_events"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}


func (m *sessionModel) toEntity() *session.Session {
	// Parse webhook events from comma-separated string
	var events []string
	if m.WebhookEvents != "" {
		events = strings.Split(m.WebhookEvents, ",")
	}

	return &session.Session{
		ID:            m.ID,
		Name:          m.Name,
		WhatsAppJID:   m.WhatsAppJID,
		Status:        types.Status(m.Status),
		QRCode:        m.QRCode,
		ProxyURL:      m.ProxyURL,
		WebhookURL:    m.WebhookURL,
		Events:        events,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}


func fromEntity(s *session.Session) *sessionModel {
	// Convert webhook events to comma-separated string
	var eventsStr string
	if len(s.Events) > 0 {
		eventsStr = strings.Join(s.Events, ",")
	}

	return &sessionModel{
		ID:            s.ID,
		Name:          s.Name,
		WhatsAppJID:   s.WhatsAppJID,
		Status:        string(s.Status),
		QRCode:        s.QRCode,
		ProxyURL:      s.ProxyURL,
		WebhookURL:    s.WebhookURL,
		WebhookEvents: eventsStr,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}


func (r *PostgresSessionRepository) Create(ctx context.Context, sess *session.Session) error {
	model := fromEntity(sess)

	query := `
		INSERT INTO sessions (id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		model.ID, model.Name, model.WhatsAppJID, model.Status,
		model.QRCode, model.ProxyURL, model.WebhookURL, model.WebhookEvents, model.CreatedAt, model.UpdatedAt)

	return err
}


func (r *PostgresSessionRepository) Save(ctx context.Context, sess *session.Session) error {
	model := fromEntity(sess)
	
	query := `
		INSERT INTO sessions (id, name, device_jid, status, qr_code, proxy_url, webhook_url, webhook_events, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			device_jid = EXCLUDED.device_jid,
			status = EXCLUDED.status,
			qr_code = EXCLUDED.qr_code,
			proxy_url = EXCLUDED.proxy_url,
			webhook_url = EXCLUDED.webhook_url,
			webhook_events = EXCLUDED.webhook_events,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		model.ID, model.Name, model.WhatsAppJID, model.Status,
		model.QRCode, model.ProxyURL, model.WebhookURL, model.WebhookEvents, model.CreatedAt, model.UpdatedAt)
	
	return err
}


func (r *PostgresSessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	var model sessionModel
	
	query := `
		SELECT id, name, COALESCE(device_jid, '') as device_jid, status,
			   COALESCE(qr_code, '') as qr_code, COALESCE(proxy_url, '') as proxy_url,
			   COALESCE(webhook_url, '') as webhook_url, COALESCE(webhook_events, '') as webhook_events,
			   created_at, updated_at
		FROM sessions WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, err
	}
	
	return model.toEntity(), nil
}


func (r *PostgresSessionRepository) GetByName(ctx context.Context, name string) (*session.Session, error) {
	var model sessionModel

	query := `
		SELECT id, name, COALESCE(device_jid, '') as device_jid, status,
			   COALESCE(qr_code, '') as qr_code, COALESCE(proxy_url, '') as proxy_url,
			   COALESCE(webhook_url, '') as webhook_url, COALESCE(webhook_events, '') as webhook_events,
			   created_at, updated_at
		FROM sessions WHERE name = $1
	`

	err := r.db.GetContext(ctx, &model, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, err
	}

	return model.toEntity(), nil
}


func (r *PostgresSessionRepository) GetAll(ctx context.Context) ([]*session.Session, error) {
	var models []sessionModel
	
	query := `
		SELECT id, name, COALESCE(device_jid, '') as device_jid, status,
			   COALESCE(qr_code, '') as qr_code, COALESCE(proxy_url, '') as proxy_url,
			   COALESCE(webhook_url, '') as webhook_url, COALESCE(webhook_events, '') as webhook_events,
			   created_at, updated_at
		FROM sessions ORDER BY created_at DESC
	`
	
	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, err
	}
	
	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sessions[i] = model.toEntity()
	}
	
	return sessions, nil
}


func (r *PostgresSessionRepository) Update(ctx context.Context, sess *session.Session) error {
	model := fromEntity(sess)
	
	query := `
		UPDATE sessions SET
			name = $2, device_jid = $3, status = $4,
			qr_code = $5, proxy_url = $6, webhook_url = $7, webhook_events = $8, updated_at = $9
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		model.ID, model.Name, model.WhatsAppJID, model.Status,
		model.QRCode, model.ProxyURL, model.WebhookURL, model.WebhookEvents, model.UpdatedAt)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}
	
	return nil
}


func (r *PostgresSessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}
	
	return nil
}


func (r *PostgresSessionRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM sessions WHERE id = $1`
	
	err := r.db.GetContext(ctx, &count, query, id)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}


func (r *PostgresSessionRepository) FindByName(ctx context.Context, name string) ([]*session.Session, error) {
	var models []sessionModel
	
	query := `
		SELECT id, name, COALESCE(device_jid, '') as device_jid, status,
			   COALESCE(qr_code, '') as qr_code, COALESCE(proxy_url, '') as proxy_url,
			   COALESCE(webhook_url, '') as webhook_url, COALESCE(webhook_events, '') as webhook_events,
			   created_at, updated_at
		FROM sessions WHERE name ILIKE $1 ORDER BY created_at DESC
	`
	
	err := r.db.SelectContext(ctx, &models, query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	
	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sessions[i] = model.toEntity()
	}
	
	return sessions, nil
}


func (r *PostgresSessionRepository) GetByStatus(ctx context.Context, status types.Status) ([]*session.Session, error) {
	var models []sessionModel
	
	query := `
		SELECT id, name, COALESCE(device_jid, '') as device_jid, status,
			   COALESCE(qr_code, '') as qr_code, COALESCE(proxy_url, '') as proxy_url,
			   COALESCE(webhook_url, '') as webhook_url, COALESCE(webhook_events, '') as webhook_events,
			   created_at, updated_at
		FROM sessions WHERE status = $1 ORDER BY created_at DESC
	`
	
	err := r.db.SelectContext(ctx, &models, query, string(status))
	if err != nil {
		return nil, err
	}
	
	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sessions[i] = model.toEntity()
	}
	
	return sessions, nil
}




func (r *PostgresSessionRepository) FindByID(ctx context.Context, id string) (*session.Session, error) {
	return r.GetByID(ctx, id)
}


func (r *PostgresSessionRepository) FindAll(ctx context.Context) ([]*session.Session, error) {
	return r.GetAll(ctx)
}


func (r *PostgresSessionRepository) FindByStatus(ctx context.Context, status string) ([]*session.Session, error) {
	return r.GetByStatus(ctx, types.Status(status))
}




func Connect(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword,
		cfg.DBName, cfg.DBSslMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}


	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}


func RunMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationsDir := "internal/infra/database/migrations"
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {

		if err.Error() == "Dirty database version 2. Fix and force version." {

			if forceErr := m.Force(1); forceErr != nil {
				return fmt.Errorf("failed to force migration version: %w", forceErr)
			}

			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to run migrations after force: %w", err)
			}
		} else {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return nil
}
