package application

import (
	"context"
	"fmt"
	"time"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

type SessionApp struct {
	sessionRepo       session.Repository
	sessionService    session.Service
	identifierService session.IdentifierService
	validator         *validation.Validator
	idGenerator       IDGenerator
}

type DefaultIDGenerator struct{}

func (d DefaultIDGenerator) GenerateSessionID() string {
	return "session-" + generateRandomID()
}

func (d DefaultIDGenerator) GenerateAPIKey() string {
	return "api-" + generateRandomID()
}

func generateRandomID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func NewSessionApp(
	sessionRepo session.Repository,
	sessionService session.Service,
	validator *validation.Validator,
) *SessionApp {
	return &SessionApp{
		sessionRepo:       sessionRepo,
		sessionService:    sessionService,
		identifierService: session.NewIdentifierService(),
		validator:         validator,
		idGenerator:       DefaultIDGenerator{},
	}
}

func NewSessionAppWithIDGenerator(
	sessionRepo session.Repository,
	sessionService session.Service,
	validator *validation.Validator,
	idGenerator IDGenerator,
) *SessionApp {
	return &SessionApp{
		sessionRepo:       sessionRepo,
		sessionService:    sessionService,
		identifierService: session.NewIdentifierService(),
		validator:         validator,
		idGenerator:       idGenerator,
	}
}

func (s *SessionApp) CreateSessionWithRequest(ctx context.Context, req *dto.CreateSessionRequest) (*session.Session, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	sessionEntity, err := session.NewSession(s.idGenerator.GenerateSessionID(), req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create session entity: %w", err)
	}

	if req.ProxyURL != "" {
		if err := sessionEntity.SetProxyURL(req.ProxyURL); err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
	}

	if err := s.sessionService.ValidateSessionConfiguration(sessionEntity); err != nil {
		return nil, fmt.Errorf("session configuration validation failed: %w", err)
	}

	if err := s.sessionRepo.Create(ctx, sessionEntity); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return sessionEntity, nil
}

func (s *SessionApp) GetSession(ctx context.Context, sessionIDOrName string) (*session.Session, error) {
	identifierInfo, err := s.identifierService.ResolveIdentifier(sessionIDOrName)
	if err != nil {
		return nil, fmt.Errorf("invalid session identifier: %w", err)
	}

	switch identifierInfo.Type {
	case session.IdentifierTypeUUID:
		sessionEntity, err := s.sessionRepo.GetByID(ctx, identifierInfo.Normalized)
		if err == nil {
			return sessionEntity, nil
		}
		sessionEntity, err = s.sessionRepo.GetByName(ctx, identifierInfo.Normalized)
		if err == nil {
			return sessionEntity, nil
		}
	case session.IdentifierTypeName:
		sessionEntity, err := s.sessionRepo.GetByName(ctx, identifierInfo.Normalized)
		if err == nil {
			return sessionEntity, nil
		}
		sessionEntity, err = s.sessionRepo.GetByID(ctx, identifierInfo.Normalized)
		if err == nil {
			return sessionEntity, nil
		}
	default:
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierInfo.Type)
	}

	return nil, fmt.Errorf("session not found: %s", sessionIDOrName)
}

func (s *SessionApp) GetSessionDTO(ctx context.Context, sessionIDOrName string) (*dto.SessionInfo, error) {
	sessionEntity, err := s.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return nil, err
	}

	return s.convertSessionToDTO(sessionEntity), nil
}

func (s *SessionApp) convertSessionToDTO(sessionEntity *session.Session) *dto.SessionInfo {
	return &dto.SessionInfo{
		ID:        sessionEntity.ID.Value(),
		Name:      sessionEntity.Name.Value(),
		Status:    string(sessionEntity.Status),
		WaJID:     sessionEntity.WaJID.Value(),
		ProxyURL:  sessionEntity.ProxyURL.Value(),
		ApiKey:    sessionEntity.ApiKey.Value(),
		CreatedAt: sessionEntity.CreatedAt,
	}
}

func (s *SessionApp) GetSessionByDeviceJID(ctx context.Context, deviceJID string) (*session.Session, error) {
	if deviceJID == "" {
		return nil, fmt.Errorf("device JID is required")
	}

	return s.sessionRepo.GetByDeviceJID(ctx, deviceJID)
}

func (s *SessionApp) GetAllSessions(ctx context.Context) ([]*session.Session, error) {
	sessions, err := s.sessionRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

func (s *SessionApp) ConnectSession(ctx context.Context, sessionIDOrName string) (*session.Session, error) {
	if sessionIDOrName == "" {
		return nil, fmt.Errorf("session ID or name is required")
	}

	sessionEntity, err := s.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if !s.sessionService.CanConnect(sessionEntity) {
		return nil, fmt.Errorf("session cannot be connected in current status: %s", sessionEntity.Status)
	}

	sessionEntity.Status = session.StatusConnecting

	if err := s.sessionService.ValidateDeviceConnection(sessionEntity, sessionEntity.WaJID.Value()); err != nil {
		return nil, fmt.Errorf("device connection validation failed: %w", err)
	}

	if err := s.sessionRepo.Update(ctx, sessionEntity); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	err = sessionEntity.SetQRCode("data:image/png;base64,placeholder_qr_code")
	if err != nil {
		return nil, fmt.Errorf("failed to set QR code: %w", err)
	}

	return sessionEntity, nil
}

func (s *SessionApp) DeleteSession(ctx context.Context, sessionIDOrName string) error {
	if sessionIDOrName == "" {
		return fmt.Errorf("session ID or name is required")
	}

	sessionEntity, err := s.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !s.sessionService.CanDelete(sessionEntity) {
		return fmt.Errorf("session cannot be deleted in current status: %s", sessionEntity.Status)
	}

	if err := s.sessionRepo.Delete(ctx, sessionEntity.ID.Value()); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *SessionApp) RegenerateApiKey(ctx context.Context, sessionIDOrName string) (string, error) {
	if sessionIDOrName == "" {
		return "", fmt.Errorf("session ID or name is required")
	}

	sessionEntity, err := s.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return "", fmt.Errorf("session not found")
	}

	if !s.sessionService.CanRegenerateApiKey(sessionEntity) {
		return "", fmt.Errorf("cannot regenerate API key for connected session")
	}

	newApiKey := s.idGenerator.GenerateAPIKey()
	err = sessionEntity.RegenerateApiKey(newApiKey)
	if err != nil {
		return "", fmt.Errorf("failed to regenerate API key: %w", err)
	}

	if err := s.sessionRepo.Update(ctx, sessionEntity); err != nil {
		return "", fmt.Errorf("failed to update session with new API key: %w", err)
	}

	return sessionEntity.ApiKey.Value(), nil
}
