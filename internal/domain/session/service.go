package session

import (
	"context"
	"zpmeow/internal/types"

	"github.com/google/uuid"
)

// SessionService defines the business logic interface for session management
// Organized by responsibility following Single Responsibility Principle
type SessionService interface {
	// Session lifecycle management
	CreateSession(ctx context.Context, name string) (*Session, error)
	GetSession(ctx context.Context, idOrName string) (*Session, error)
	GetAllSessions(ctx context.Context) ([]*Session, error)
	DeleteSession(ctx context.Context, id string) error

	// Connection management
	ConnectSession(ctx context.Context, id string) error
	DisconnectSession(ctx context.Context, id string) error

	// Authentication management
	GetQRCode(ctx context.Context, id string) (string, error)
	PairWithPhone(ctx context.Context, id, phoneNumber string) (string, error)

	// Configuration management
	SetProxy(ctx context.Context, id, proxyURL string) error
	ClearProxy(ctx context.Context, id string) error

	// System operations
	ConnectOnStartup(ctx context.Context) error
}

// WhatsAppService defines the interface for WhatsApp operations
// Separated for better testability and dependency inversion
type WhatsAppService interface {
	// Session management
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	IsClientConnected(sessionID string) bool
	GetClientStatus(sessionID string) types.Status
	ConnectOnStartup(ctx context.Context) error

	// Chat operations
	DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error
	EditMessage(ctx context.Context, sessionID, chatJID, messageID, newText string) (*types.SendResponse, error)
	DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error)
	ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string) error
}

// SessionServiceImpl implements the SessionService interface
type SessionServiceImpl struct {
	repo            SessionRepository
	whatsappService WhatsAppService
}

// NewSessionService creates a new session service with dependency injection
func NewSessionService(repo SessionRepository, whatsappService WhatsAppService) SessionService {
	return &SessionServiceImpl{
		repo:            repo,
		whatsappService: whatsappService,
	}
}

// Session lifecycle management

// CreateSession creates a new session with validation
func (s *SessionServiceImpl) CreateSession(ctx context.Context, name string) (*Session, error) {
	if err := s.validateSessionName(name); err != nil {
		return nil, err
	}

	// Check if a session with this name already exists
	_, err := s.repo.GetByName(ctx, name)
	if err == nil {
		return nil, ErrSessionAlreadyExists
	}
	if err != ErrSessionNotFound {
		return nil, err
	}

	session := NewSession(generateSessionID(), name)

	if err := session.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by ID or name with validation
func (s *SessionServiceImpl) GetSession(ctx context.Context, idOrName string) (*Session, error) {
	if err := s.validateSessionIDOrName(idOrName); err != nil {
		return nil, err
	}

	// First try to get by ID
	session, err := s.repo.GetByID(ctx, idOrName)
	if err == nil {
		return session, nil
	}

	// If not found by ID and error is "not found", try by name
	if err == ErrSessionNotFound {
		session, nameErr := s.repo.GetByName(ctx, idOrName)
		if nameErr == nil {
			return session, nil
		}
		// If both ID and name lookups failed, return the original "not found" error
		if nameErr == ErrSessionNotFound {
			return nil, ErrSessionNotFound
		}
		// If name lookup failed with a different error, return that error
		return nil, nameErr
	}

	// If ID lookup failed with a different error, return that error
	return nil, err
}

// GetAllSessions retrieves all sessions
func (s *SessionServiceImpl) GetAllSessions(ctx context.Context) ([]*Session, error) {
	return s.repo.GetAll(ctx)
}

// DeleteSession deletes a session with proper cleanup
func (s *SessionServiceImpl) DeleteSession(ctx context.Context, id string) error {
	if err := s.validateSessionID(id); err != nil {
		return err
	}

	// Ensure WhatsApp client is stopped before deletion
	_ = s.whatsappService.StopClient(id)

	return s.repo.Delete(ctx, id)
}

// Connection management

// ConnectSession connects a session to WhatsApp with proper validation
func (s *SessionServiceImpl) ConnectSession(ctx context.Context, id string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	if s.whatsappService.IsClientConnected(id) {
		return ErrSessionAlreadyConnected
	}

	if !session.CanConnect() {
		return ErrSessionCannotConnect
	}

	return s.performConnection(ctx, session)
}

// DisconnectSession disconnects a session from WhatsApp
func (s *SessionServiceImpl) DisconnectSession(ctx context.Context, id string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	if err := s.whatsappService.StopClient(id); err != nil {
		return err
	}

	session.SetStatus(types.StatusDisconnected)
	return s.repo.Update(ctx, session)
}

// performConnection handles the connection process
func (s *SessionServiceImpl) performConnection(ctx context.Context, session *Session) error {
	session.SetStatus(types.StatusConnecting)
	if err := s.repo.Update(ctx, session); err != nil {
		return err
	}

	return s.whatsappService.StartClient(session.ID)
}

// Authentication management

// GetQRCode retrieves the QR code for a session
func (s *SessionServiceImpl) GetQRCode(ctx context.Context, id string) (string, error) {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return "", err
	}

	if session.IsConnected() {
		return "", ErrSessionAlreadyConnected
	}

	// Try to get fresh QR code if session can connect and doesn't have one
	if !session.HasQRCode() && session.CanConnect() {
		return s.whatsappService.GetQRCode(id)
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
		return "", ErrSessionCannotConnect
	}

	return s.whatsappService.PairPhone(id, phoneNumber)
}

// Configuration management

// SetProxy sets the proxy for a session
func (s *SessionServiceImpl) SetProxy(ctx context.Context, id, proxyURL string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	session.SetProxyURL(proxyURL)
	return s.repo.Update(ctx, session)
}

// ClearProxy removes the proxy configuration for a session
func (s *SessionServiceImpl) ClearProxy(ctx context.Context, id string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	session.ClearProxy()
	return s.repo.Update(ctx, session)
}

// System operations

// ConnectOnStartup connects all previously connected sessions on startup
func (s *SessionServiceImpl) ConnectOnStartup(ctx context.Context) error {
	return s.whatsappService.ConnectOnStartup(ctx)
}

// Private helper methods

// validateSessionID validates a session ID using centralized validation
func (s *SessionServiceImpl) validateSessionID(id string) error {
	return ValidateSessionID(id)
}

// validateSessionName validates a session name using centralized validation
func (s *SessionServiceImpl) validateSessionName(name string) error {
	return ValidateSessionName(name)
}

// validateSessionIDOrName validates a session ID or name using centralized validation
func (s *SessionServiceImpl) validateSessionIDOrName(idOrName string) error {
	return ValidateSessionIDOrName(idOrName)
}



// generateSessionID generates a unique session ID
func generateSessionID() string {
	return uuid.New().String()
}
