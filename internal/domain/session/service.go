package session

import (
	"context"
	"zpmeow/internal/types"

	"github.com/google/uuid"
)

// SessionService defines the business logic interface for session management
type SessionService interface {
	// Session CRUD operations
	CreateSession(ctx context.Context, name string) (*Session, error)
	GetSession(ctx context.Context, id string) (*Session, error)
	GetAllSessions(ctx context.Context) ([]*Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id string) error

	// WhatsApp operations
	ConnectSession(ctx context.Context, id string) error
	DisconnectSession(ctx context.Context, id string) error
	GetQRCode(ctx context.Context, id string) (string, error)
	PairWithPhone(ctx context.Context, id, phoneNumber string) (string, error)

	// Proxy operations
	SetProxy(ctx context.Context, id, proxyURL string) error
	GetProxy(ctx context.Context, id string) (string, error)

	// Startup operations
	ConnectOnStartup(ctx context.Context) error
}

// WhatsAppService defines the interface for WhatsApp operations
type WhatsAppService interface {
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	IsClientConnected(sessionID string) bool
	GetClientStatus(sessionID string) types.Status
	ConnectOnStartup(ctx context.Context) error
}

// SessionServiceImpl implements the SessionService interface
type SessionServiceImpl struct {
	repo            SessionRepository
	whatsappService WhatsAppService
}

// NewSessionService creates a new session service
func NewSessionService(repo SessionRepository, whatsappService WhatsAppService) SessionService {
	return &SessionServiceImpl{
		repo:            repo,
		whatsappService: whatsappService,
	}
}

// CreateSession creates a new session
func (s *SessionServiceImpl) CreateSession(ctx context.Context, name string) (*Session, error) {
	// Validate input
	if name == "" {
		return nil, ErrInvalidSessionName
	}

	// Create new session entity
	session := NewSession(generateSessionID(), name)

	// Validate entity
	if err := session.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.repo.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (s *SessionServiceImpl) GetSession(ctx context.Context, id string) (*Session, error) {
	if id == "" {
		return nil, ErrInvalidSessionID
	}

	session, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetAllSessions retrieves all sessions
func (s *SessionServiceImpl) GetAllSessions(ctx context.Context) ([]*Session, error) {
	return s.repo.FindAll(ctx)
}

// UpdateSession updates an existing session
func (s *SessionServiceImpl) UpdateSession(ctx context.Context, session *Session) error {
	if err := session.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, session)
}

// DeleteSession deletes a session
func (s *SessionServiceImpl) DeleteSession(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidSessionID
	}

	// Disconnect WhatsApp client first
	_ = s.whatsappService.StopClient(id)

	return s.repo.Delete(ctx, id)
}

// ConnectSession connects a session to WhatsApp
func (s *SessionServiceImpl) ConnectSession(ctx context.Context, id string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	// Check if client is already connected (like wuzapi does)
	if s.whatsappService.IsClientConnected(id) {
		return NewDomainError("already connected")
	}

	if !session.CanConnect() {
		return NewDomainError("session cannot be connected in current state")
	}

	// Update status to connecting
	session.SetStatus(types.StatusConnecting)
	if err := s.repo.Update(ctx, session); err != nil {
		return err
	}

	// Start WhatsApp client
	return s.whatsappService.StartClient(id)
}

// DisconnectSession disconnects a session from WhatsApp
func (s *SessionServiceImpl) DisconnectSession(ctx context.Context, id string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	// Stop WhatsApp client
	if err := s.whatsappService.StopClient(id); err != nil {
		return err
	}

	// Update status
	session.SetStatus(types.StatusDisconnected)
	return s.repo.Update(ctx, session)
}

// GetQRCode retrieves the QR code for a session
func (s *SessionServiceImpl) GetQRCode(ctx context.Context, id string) (string, error) {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return "", err
	}

	return session.QRCode, nil
}

// PairWithPhone pairs a session with a phone number
func (s *SessionServiceImpl) PairWithPhone(ctx context.Context, id, phoneNumber string) (string, error) {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return "", err
	}

	if !session.CanConnect() {
		return "", NewDomainError("session cannot be paired in current state")
	}

	return s.whatsappService.PairPhone(id, phoneNumber)
}

// SetProxy sets the proxy for a session
func (s *SessionServiceImpl) SetProxy(ctx context.Context, id, proxyURL string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	session.SetProxyURL(proxyURL)
	return s.repo.Update(ctx, session)
}

// GetProxy gets the proxy for a session
func (s *SessionServiceImpl) GetProxy(ctx context.Context, id string) (string, error) {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return "", err
	}

	return session.ProxyURL, nil
}

// ConnectOnStartup connects all previously connected sessions on startup
func (s *SessionServiceImpl) ConnectOnStartup(ctx context.Context) error {
	return s.whatsappService.ConnectOnStartup(ctx)
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return uuid.New().String()
}
