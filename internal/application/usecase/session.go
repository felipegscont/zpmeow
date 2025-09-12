package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"zpmeow/internal/domain"
	"zpmeow/internal/shared/types"
)

// SessionServiceImpl implements domain.SessionService
type SessionServiceImpl struct {
	repo            domain.SessionRepository
	whatsappService domain.WhatsAppService
}

// NewSessionService creates a new session service
func NewSessionService(repo domain.SessionRepository, whatsappService domain.WhatsAppService) domain.SessionService {
	return &SessionServiceImpl{
		repo:            repo,
		whatsappService: whatsappService,
	}
}

// CreateSession creates a new session with validation
func (s *SessionServiceImpl) CreateSession(ctx context.Context, name string) (*domain.Session, error) {
	// Validate session name using domain validation
	if err := domain.ValidateSessionName(name); err != nil {
		return nil, err
	}

	// Check if session with same name already exists
	existing, err := s.repo.GetByName(ctx, name)
	if err == nil && existing != nil {
		return nil, domain.ErrSessionAlreadyExists
	}

	// Generate unique ID
	sessionID := generateSessionID()

	// Create new session entity
	session := &domain.Session{
		ID:        sessionID,
		Name:      name,
		Status:    types.StatusDisconnected,
		CreatedAt: getCurrentTime(),
		UpdatedAt: getCurrentTime(),
	}

	// Save to repository
	if err := s.repo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *SessionServiceImpl) GetSession(ctx context.Context, idOrName string) (*domain.Session, error) {
	// Validate input
	if idOrName == "" {
		return nil, domain.ErrInvalidSessionID
	}

	// Try to get by ID first
	session, err := s.repo.GetByID(ctx, idOrName)
	if err == nil && session != nil {
		return session, nil
	}

	// If not found by ID, try by name
	session, err = s.repo.GetByName(ctx, idOrName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return nil, domain.ErrSessionNotFound
	}

	return session, nil
}

func (s *SessionServiceImpl) GetAllSessions(ctx context.Context) ([]*domain.Session, error) {
	sessions, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sessions: %w", err)
	}

	return sessions, nil
}

func (s *SessionServiceImpl) UpdateSession(ctx context.Context, session *domain.Session) error {
	// Validate session
	if session == nil {
		return domain.ErrInvalidSessionID
	}

	if err := domain.ValidateSessionName(session.Name); err != nil {
		return err
	}

	// Update timestamp
	session.UpdatedAt = getCurrentTime()

	// Save to repository
	if err := s.repo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (s *SessionServiceImpl) DeleteSession(ctx context.Context, id string) error {
	// Validate input
	if err := domain.ValidateSessionID(id); err != nil {
		return err
	}

	// Check if session exists
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return domain.ErrSessionNotFound
	}

	// Check if session can be deleted (business rule)
	if session.Status == types.StatusConnected {
		return domain.ErrSessionCannotDelete
	}

	// Delete from repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *SessionServiceImpl) ConnectSession(ctx context.Context, id string) error {
	// Validate input
	if err := domain.ValidateSessionID(id); err != nil {
		return err
	}

	// Get session
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return domain.ErrSessionNotFound
	}

	// Check if already connected
	if session.Status == types.StatusConnected {
		return domain.ErrSessionAlreadyConnected
	}

	// Use WhatsApp service to connect if available
	if s.whatsappService != nil {
		if err := s.whatsappService.StartClient(id); err != nil {
			return fmt.Errorf("failed to start WhatsApp client: %w", err)
		}
	}

	// Update session status
	session.Status = types.StatusConnecting
	session.UpdatedAt = getCurrentTime()

	if err := s.repo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	return nil
}

func (s *SessionServiceImpl) DisconnectSession(ctx context.Context, id string) error {
	// Validate input
	if err := domain.ValidateSessionID(id); err != nil {
		return err
	}

	// Get session
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return domain.ErrSessionNotFound
	}

	// Use WhatsApp service to disconnect if available
	if s.whatsappService != nil {
		if err := s.whatsappService.StopClient(id); err != nil {
			return fmt.Errorf("failed to stop WhatsApp client: %w", err)
		}
	}

	// Update session status
	session.Status = types.StatusDisconnected
	session.UpdatedAt = getCurrentTime()

	if err := s.repo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	return nil
}

func (s *SessionServiceImpl) GetQRCode(ctx context.Context, id string) (string, error) {
	// Validate input
	if err := domain.ValidateSessionID(id); err != nil {
		return "", err
	}

	// Check if session exists
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return "", domain.ErrSessionNotFound
	}

	// Use WhatsApp service to get QR code if available
	if s.whatsappService != nil {
		qrCode, err := s.whatsappService.GetQRCode(id)
		if err != nil {
			return "", fmt.Errorf("failed to get QR code: %w", err)
		}
		return qrCode, nil
	}

	return "", domain.ErrWhatsAppServiceUnavailable
}

func (s *SessionServiceImpl) PairWithPhone(ctx context.Context, id, phoneNumber string) (string, error) {
	// Validate input
	if err := domain.ValidateSessionID(id); err != nil {
		return "", err
	}

	if phoneNumber == "" {
		return "", domain.NewDomainError("phone number is required")
	}

	// Check if session exists
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return "", domain.ErrSessionNotFound
	}

	// Use WhatsApp service to pair with phone if available
	if s.whatsappService != nil {
		pairingCode, err := s.whatsappService.PairPhone(id, phoneNumber)
		if err != nil {
			return "", fmt.Errorf("failed to pair with phone: %w", err)
		}
		return pairingCode, nil
	}

	return "", domain.ErrWhatsAppServiceUnavailable
}

func (s *SessionServiceImpl) SetProxy(ctx context.Context, id, proxyURL string) error {
	// Validate input
	if err := domain.ValidateSessionID(id); err != nil {
		return err
	}

	if err := domain.ValidateProxyURL(proxyURL); err != nil {
		return err
	}

	// Get session
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return domain.ErrSessionNotFound
	}

	// Update session with proxy URL
	session.ProxyURL = proxyURL
	session.UpdatedAt = getCurrentTime()

	if err := s.repo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session proxy: %w", err)
	}

	return nil
}

func (s *SessionServiceImpl) ClearProxy(ctx context.Context, id string) error {
	// Validate input
	if err := domain.ValidateSessionID(id); err != nil {
		return err
	}

	// Get session
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return domain.ErrSessionNotFound
	}

	// Clear proxy URL
	session.ProxyURL = ""
	session.UpdatedAt = getCurrentTime()

	if err := s.repo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to clear session proxy: %w", err)
	}

	return nil
}

func (s *SessionServiceImpl) ConnectOnStartup(ctx context.Context) error {
	// Get all sessions that should be connected on startup
	sessions, err := s.repo.GetByStatus(ctx, types.StatusConnected)
	if err != nil {
		return fmt.Errorf("failed to get connected sessions: %w", err)
	}

	// Use WhatsApp service to connect on startup if available
	if s.whatsappService != nil {
		if err := s.whatsappService.ConnectOnStartup(ctx); err != nil {
			return fmt.Errorf("failed to connect sessions on startup: %w", err)
		}
	}

	// Log the number of sessions that should be reconnected
	if len(sessions) > 0 {
		// Note: In a real implementation, you might want to use a proper logger here
		fmt.Printf("Found %d sessions to reconnect on startup\n", len(sessions))
	}

	return nil
}

// Helper functions

// generateSessionID generates a unique session ID
func generateSessionID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("session_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// getCurrentTime returns current time (useful for testing)
func getCurrentTime() time.Time {
	return time.Now()
}
