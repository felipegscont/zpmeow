package session

import "context"

// SessionService defines domain business rules for session management
// This is a pure domain service focused on business logic only
type SessionService interface {
	// Business rules for session lifecycle
	CanCreateSession(name string) error
	CanDeleteSession(session *Session) error
	CanConnectSession(session *Session) error
	CanDisconnectSession(session *Session) error

	// Business validation rules
	ValidateSessionName(name string) error
	ValidateSessionConfiguration(session *Session) error

	// Business state transitions
	TransitionToConnecting(session *Session) error
	TransitionToConnected(session *Session) error
	TransitionToDisconnected(session *Session) error
	TransitionToError(session *Session, reason string) error
}

// SessionManager interface is defined in interfaces.go to avoid duplication

// ConnectionManager defines domain capabilities for connection management
// Pure domain interface - no technical details
type ConnectionManager interface {
	// Domain connection concepts
	EstablishConnection(ctx context.Context, sessionID string) error
	TerminateConnection(ctx context.Context, sessionID string) error
	CheckConnectionStatus(ctx context.Context, sessionID string) (Status, error)

	// Authentication from domain perspective
	InitiateAuthentication(ctx context.Context, sessionID string) (string, error) // Returns QR or pairing code
	CompleteAuthentication(ctx context.Context, sessionID, credential string) error
}

// CommunicationManager defines domain capabilities for messaging
// Focused on business concepts, not technical implementation
type CommunicationManager interface {
	// Domain messaging concepts
	SendMessage(ctx context.Context, sessionID, recipient, content string) error
	ReceiveMessage(ctx context.Context, sessionID string) error

	// Group management from domain perspective
	CreateCommunicationGroup(ctx context.Context, sessionID, name string, members []string) error
	ManageGroupMembership(ctx context.Context, sessionID, groupID string, members []string, action string) error
}

// SessionDomainService implements business rules for sessions
type SessionDomainService struct{}

// NewSessionDomainService creates a new session domain service
func NewSessionDomainService() SessionService {
	return &SessionDomainService{}
}

// CanCreateSession validates if a session can be created (business rule)
func (s *SessionDomainService) CanCreateSession(name string) error {
	if name == "" {
		return ErrInvalidSessionName
	}

	// Business rule: Session name must be unique and valid
	_, err := NewSessionName(name)
	if err != nil {
		return err
	}

	return nil
}

// CanDeleteSession validates if a session can be deleted (business rule)
func (s *SessionDomainService) CanDeleteSession(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	// Business rule: Cannot delete connected sessions
	if session.Status == StatusConnected {
		return ErrCannotDeleteConnectedSession
	}

	return nil
}

// CanConnectSession validates if a session can connect (business rule)
func (s *SessionDomainService) CanConnectSession(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	// Business rule: Can only connect disconnected or error sessions
	if session.Status == StatusConnected || session.Status == StatusConnecting {
		return ErrSessionAlreadyConnected
	}

	return nil
}

// CanDisconnectSession validates if a session can disconnect (business rule)
func (s *SessionDomainService) CanDisconnectSession(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	// Business rule: Can only disconnect connected sessions
	if session.Status != StatusConnected {
		return ErrSessionNotConnected
	}

	return nil
}

// ValidateSessionName validates session name according to business rules
func (s *SessionDomainService) ValidateSessionName(name string) error {
	_, err := NewSessionName(name)
	return err
}

// ValidateSessionConfiguration validates session configuration
func (s *SessionDomainService) ValidateSessionConfiguration(session *Session) error {
	if session == nil {
		return ErrInvalidSession
	}

	// Validate all session components
	if session.ID.IsEmpty() {
		return ErrInvalidSessionID
	}

	if session.Name.IsEmpty() {
		return ErrInvalidSessionName
	}

	return nil
}

// TransitionToConnecting transitions session to connecting state
func (s *SessionDomainService) TransitionToConnecting(session *Session) error {
	if err := s.CanConnectSession(session); err != nil {
		return err
	}

	session.SetStatus(StatusConnecting)
	return nil
}

// TransitionToConnected transitions session to connected state
func (s *SessionDomainService) TransitionToConnected(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	session.SetStatus(StatusConnected)
	return nil
}

// TransitionToDisconnected transitions session to disconnected state
func (s *SessionDomainService) TransitionToDisconnected(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	session.SetStatus(StatusDisconnected)
	return nil
}

// TransitionToError transitions session to error state
func (s *SessionDomainService) TransitionToError(session *Session, reason string) error {
	if session == nil {
		return ErrSessionNotFound
	}

	session.SetStatus(StatusError)
	// In a real implementation, you might want to store the error reason
	return nil
}
