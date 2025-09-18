package session

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
	"meow/internal/domain/session"
)

// GetSessionQuery represents the query to get a session
type GetSessionQuery struct {
	SessionID string
	Name      string
	ApiKey    string
}

// Validate validates the get session query
func (q GetSessionQuery) Validate() error {
	// At least one identifier must be provided
	if strings.TrimSpace(q.SessionID) == "" &&
		strings.TrimSpace(q.Name) == "" &&
		strings.TrimSpace(q.ApiKey) == "" {
		return common.NewValidationError("identifier", "", "at least one identifier (sessionID, name, or apiKey) is required")
	}

	return nil
}

// SessionView represents the view model for a session
type SessionView struct {
	SessionID          string
	Name               string
	Status             string
	WaJID              string
	QRCode             string
	ProxyConfiguration string
	WebhookEndpoint    string
	HasProxy           bool
	HasWebhook         bool
	IsConnected        bool
	IsAuthenticated    bool
	CreatedAt          string
	UpdatedAt          string
}

// GetSessionUseCase handles getting session information
type GetSessionUseCase struct {
	sessionRepo ports.SessionRepository
	logger      ports.Logger
}

// NewGetSessionUseCase creates a new GetSessionUseCase
func NewGetSessionUseCase(
	sessionRepo ports.SessionRepository,
	logger ports.Logger,
) *GetSessionUseCase {
	return &GetSessionUseCase{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

// Handle executes the get session use case
func (uc *GetSessionUseCase) Handle(ctx context.Context, query GetSessionQuery) (*SessionView, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid get session query", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var sessionEntity *session.Session
	var err error

	// 2. Get session by the provided identifier
	switch {
	case query.SessionID != "":
		sessionEntity, err = uc.sessionRepo.GetByID(ctx, query.SessionID)
		if err != nil {
			uc.logger.Error(ctx, "Failed to get session by ID", "sessionID", query.SessionID, "error", err)
			return nil, fmt.Errorf("failed to get session by ID: %w", err)
		}

	case query.Name != "":
		sessionEntity, err = uc.sessionRepo.GetByName(ctx, query.Name)
		if err != nil {
			uc.logger.Error(ctx, "Failed to get session by name", "name", query.Name, "error", err)
			return nil, fmt.Errorf("failed to get session by name: %w", err)
		}

	case query.ApiKey != "":
		sessionEntity, err = uc.sessionRepo.GetByApiKey(ctx, query.ApiKey)
		if err != nil {
			uc.logger.Error(ctx, "Failed to get session by API key", "error", err)
			return nil, fmt.Errorf("failed to get session by API key: %w", err)
		}
	}

	// 3. Convert domain entity to view model
	view := &SessionView{
		SessionID:          sessionEntity.SessionID().Value(),
		Name:               sessionEntity.Name().Value(),
		Status:             sessionEntity.Status().String(),
		WaJID:              sessionEntity.WaJID().Value(),
		QRCode:             sessionEntity.QRCode().Value(),
		ProxyConfiguration: sessionEntity.ProxyConfiguration().Value(),
		WebhookEndpoint:    sessionEntity.WebhookEndpoint().Value(),
		HasProxy:           sessionEntity.HasProxy(),
		HasWebhook:         sessionEntity.HasWebhook(),
		IsConnected:        sessionEntity.IsConnected(),
		IsAuthenticated:    sessionEntity.IsAuthenticated(),
		CreatedAt:          sessionEntity.CreatedAt().Value().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          sessionEntity.UpdatedAt().Value().Format("2006-01-02T15:04:05Z07:00"),
	}

	uc.logger.Debug(ctx, "Session retrieved successfully", "sessionID", view.SessionID, "name", view.Name)

	return view, nil
}

// GetAllSessionsQuery represents the query to get all sessions
type GetAllSessionsQuery struct {
	// Future: Add filtering, pagination, sorting parameters
}

// Validate validates the get all sessions query
func (q GetAllSessionsQuery) Validate() error {
	// No validation needed for now
	return nil
}

// GetAllSessionsUseCase handles getting all sessions
type GetAllSessionsUseCase struct {
	sessionRepo ports.SessionRepository
	logger      ports.Logger
}

// NewGetAllSessionsUseCase creates a new GetAllSessionsUseCase
func NewGetAllSessionsUseCase(
	sessionRepo ports.SessionRepository,
	logger ports.Logger,
) *GetAllSessionsUseCase {
	return &GetAllSessionsUseCase{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

// Handle executes the get all sessions use case
func (uc *GetAllSessionsUseCase) Handle(ctx context.Context, query GetAllSessionsQuery) ([]*SessionView, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid get all sessions query", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get all sessions from repository
	sessions, err := uc.sessionRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get all sessions", "error", err)
		return nil, fmt.Errorf("failed to get all sessions: %w", err)
	}

	// 3. Convert domain entities to view models
	views := make([]*SessionView, len(sessions))
	for i, sessionEntity := range sessions {
		views[i] = &SessionView{
			SessionID:          sessionEntity.SessionID().Value(),
			Name:               sessionEntity.Name().Value(),
			Status:             sessionEntity.Status().String(),
			WaJID:              sessionEntity.WaJID().Value(),
			QRCode:             sessionEntity.QRCode().Value(),
			ProxyConfiguration: sessionEntity.ProxyConfiguration().Value(),
			WebhookEndpoint:    sessionEntity.WebhookEndpoint().Value(),
			HasProxy:           sessionEntity.HasProxy(),
			HasWebhook:         sessionEntity.HasWebhook(),
			IsConnected:        sessionEntity.IsConnected(),
			IsAuthenticated:    sessionEntity.IsAuthenticated(),
			CreatedAt:          sessionEntity.CreatedAt().Value().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:          sessionEntity.UpdatedAt().Value().Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	uc.logger.Debug(ctx, "All sessions retrieved successfully", "count", len(views))

	return views, nil
}
