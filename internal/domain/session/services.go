package session

import "context"

type SessionService interface {
	CanCreateSession(name string) error
	CanDeleteSession(session *Session) error
	CanConnectSession(session *Session) error
	CanDisconnectSession(session *Session) error

	ValidateSessionName(name string) error
	ValidateSessionConfiguration(session *Session) error

	TransitionToConnecting(session *Session) error
	TransitionToConnected(session *Session) error
	TransitionToDisconnected(session *Session) error
	TransitionToError(session *Session, reason string) error
}

type ConnectionManager interface {
	EstablishConnection(ctx context.Context, sessionID string) error
	TerminateConnection(ctx context.Context, sessionID string) error
	CheckConnectionStatus(ctx context.Context, sessionID string) (Status, error)

	InitiateAuthentication(ctx context.Context, sessionID string) (string, error) // Returns QR or pairing code
	CompleteAuthentication(ctx context.Context, sessionID, credential string) error
}

type CommunicationManager interface {
	SendMessage(ctx context.Context, sessionID, recipient, content string) error
	ReceiveMessage(ctx context.Context, sessionID string) error

	CreateCommunicationGroup(ctx context.Context, sessionID, name string, members []string) error
	ManageGroupMembership(ctx context.Context, sessionID, groupID string, members []string, action string) error
}

type SessionDomainService struct{}

func NewSessionDomainService() SessionService {
	return &SessionDomainService{}
}

func (s *SessionDomainService) CanCreateSession(name string) error {
	if name == "" {
		return ErrInvalidSessionName
	}

	_, err := NewSessionName(name)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionDomainService) CanDeleteSession(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	if session.Status == StatusConnected {
		return ErrCannotDeleteConnectedSession
	}

	return nil
}

func (s *SessionDomainService) CanConnectSession(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	if session.Status == StatusConnected || session.Status == StatusConnecting {
		return ErrSessionAlreadyConnected
	}

	return nil
}

func (s *SessionDomainService) CanDisconnectSession(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	if session.Status != StatusConnected {
		return ErrSessionNotConnected
	}

	return nil
}

func (s *SessionDomainService) ValidateSessionName(name string) error {
	_, err := NewSessionName(name)
	return err
}

func (s *SessionDomainService) ValidateSessionConfiguration(session *Session) error {
	if session == nil {
		return ErrInvalidSession
	}

	if session.ID.IsEmpty() {
		return ErrInvalidSessionID
	}

	if session.Name.IsEmpty() {
		return ErrInvalidSessionName
	}

	return nil
}

func (s *SessionDomainService) TransitionToConnecting(session *Session) error {
	if err := s.CanConnectSession(session); err != nil {
		return err
	}

	session.SetStatus(StatusConnecting)
	return nil
}

func (s *SessionDomainService) TransitionToConnected(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	session.SetStatus(StatusConnected)
	return nil
}

func (s *SessionDomainService) TransitionToDisconnected(session *Session) error {
	if session == nil {
		return ErrSessionNotFound
	}

	session.SetStatus(StatusDisconnected)
	return nil
}

func (s *SessionDomainService) TransitionToError(session *Session, reason string) error {
	if session == nil {
		return ErrSessionNotFound
	}

	session.SetStatus(StatusError)
	return nil
}
