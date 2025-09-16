package application

import (
	"context"
	"fmt"
	"strings"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/interfaces/dto"
	"zpmeow/internal/shared/validation"

	"github.com/google/uuid"
)

type SessionService struct {
	sessionRepo    session.Repository
	sessionService session.DomainService
	validator      *validation.Validator
}

func NewSessionService(
	sessionRepo session.Repository,
	sessionService session.DomainService,
	validator *validation.Validator,
) *SessionService {
	return &SessionService{
		sessionRepo:    sessionRepo,
		sessionService: sessionService,
		validator:      validator,
	}
}

func (s *SessionService) CreateSessionWithRequest(ctx context.Context, req *dto.CreateSessionRequest) (*session.Session, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	sessionEntity := session.NewSession(uuid.New().String(), req.Name)
	sessionEntity.ProxyURL = req.ProxyURL
	sessionEntity.WebhookURL = req.WebhookURL
	if req.Events != "" {
		sessionEntity.Events = []string{req.Events}
	}

	if err := s.sessionService.ValidateSessionConfiguration(sessionEntity); err != nil {
		return nil, fmt.Errorf("session configuration validation failed: %w", err)
	}

	if err := s.sessionRepo.Create(ctx, sessionEntity); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return sessionEntity, nil
}

func (s *SessionService) GetSession(ctx context.Context, sessionIDOrName string) (*session.Session, error) {
	if sessionIDOrName == "" {
		return nil, fmt.Errorf("session ID or name is required")
	}

	// Clean the input
	sessionIDOrName = strings.TrimSpace(sessionIDOrName)

	// Check if it looks like a UUID (contains hyphens and is 36 chars)
	isUUID := len(sessionIDOrName) == 36 && strings.Contains(sessionIDOrName, "-")

	if isUUID {
		// If it looks like a UUID, try to parse it first to validate
		if _, err := uuid.Parse(sessionIDOrName); err == nil {
			// Valid UUID, try to get by ID first
			sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionIDOrName)
			if err == nil {
				return sessionEntity, nil
			}
		}
	}

	// If not a valid UUID or not found by ID, try by name first
	sessionEntity, err := s.sessionRepo.GetByName(ctx, sessionIDOrName)
	if err == nil {
		return sessionEntity, nil
	}

	// If not found by name and wasn't tried as ID yet, try by ID
	if !isUUID {
		sessionEntity, err = s.sessionRepo.GetByID(ctx, sessionIDOrName)
		if err == nil {
			return sessionEntity, nil
		}
	}

	// Not found by either method
	return nil, fmt.Errorf("session not found")
}

func (s *SessionService) GetSessionByDeviceJID(ctx context.Context, deviceJID string) (*session.Session, error) {
	if deviceJID == "" {
		return nil, fmt.Errorf("device JID is required")
	}

	return s.sessionRepo.GetByDeviceJID(ctx, deviceJID)
}

func (s *SessionService) ListSessionEntities(ctx context.Context) ([]*session.Session, error) {
	sessions, err := s.sessionRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

func (s *SessionService) GetAllSessions(ctx context.Context) ([]*session.Session, error) {
	return s.ListSessionEntities(ctx)
}

func (s *SessionService) ConnectSession(ctx context.Context, sessionIDOrName string) (*session.Session, error) {
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

	// Validate device connection consistency before updating
	if err := s.sessionService.ValidateDeviceConnection(sessionEntity, sessionEntity.WaJID); err != nil {
		return nil, fmt.Errorf("device connection validation failed: %w", err)
	}

	if err := s.sessionRepo.Update(ctx, sessionEntity); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Set a placeholder QR code - this should be handled by the infrastructure layer
	sessionEntity.QRCode = "data:image/png;base64,placeholder_qr_code"

	return sessionEntity, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionIDOrName string) error {
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

	if err := s.sessionRepo.Delete(ctx, sessionEntity.ID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *SessionService) RegenerateApiKey(ctx context.Context, sessionIDOrName string) (string, error) {
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

	// Regenerate API key using domain method
	sessionEntity.RegenerateApiKey()

	if err := s.sessionRepo.Update(ctx, sessionEntity); err != nil {
		return "", fmt.Errorf("failed to update session with new API key: %w", err)
	}

	return sessionEntity.ApiKey, nil
}
