package session

import (
	"context"
	"zpmeow/internal/types"

	"github.com/google/uuid"
)



type SessionService interface {

	CreateSession(ctx context.Context, name string) (*Session, error)
	GetSession(ctx context.Context, idOrName string) (*Session, error)
	GetAllSessions(ctx context.Context) ([]*Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id string) error


	ConnectSession(ctx context.Context, id string) error
	DisconnectSession(ctx context.Context, id string) error


	GetQRCode(ctx context.Context, id string) (string, error)
	PairWithPhone(ctx context.Context, id, phoneNumber string) (string, error)


	SetProxy(ctx context.Context, id, proxyURL string) error
	ClearProxy(ctx context.Context, id string) error


	ConnectOnStartup(ctx context.Context) error
}



type WhatsAppService interface {

	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	IsClientConnected(sessionID string) bool
	GetClientStatus(sessionID string) types.Status
	ConnectOnStartup(ctx context.Context) error


	DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error
	EditMessage(ctx context.Context, sessionID, chatJID, messageID, newText string) (*types.SendResponse, error)
	DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error)
	ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string) error
}


type SessionServiceImpl struct {
	repo            SessionRepository
	whatsappService WhatsAppService
}


func NewSessionService(repo SessionRepository, whatsappService WhatsAppService) SessionService {
	return &SessionServiceImpl{
		repo:            repo,
		whatsappService: whatsappService,
	}
}




func (s *SessionServiceImpl) CreateSession(ctx context.Context, name string) (*Session, error) {
	if err := s.validateSessionName(name); err != nil {
		return nil, err
	}


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


func (s *SessionServiceImpl) GetSession(ctx context.Context, idOrName string) (*Session, error) {
	if err := s.validateSessionIDOrName(idOrName); err != nil {
		return nil, err
	}


	session, err := s.repo.GetByID(ctx, idOrName)
	if err == nil {
		return session, nil
	}


	if err == ErrSessionNotFound {
		session, nameErr := s.repo.GetByName(ctx, idOrName)
		if nameErr == nil {
			return session, nil
		}

		if nameErr == ErrSessionNotFound {
			return nil, ErrSessionNotFound
		}

		return nil, nameErr
	}


	return nil, err
}


func (s *SessionServiceImpl) GetAllSessions(ctx context.Context) ([]*Session, error) {
	return s.repo.GetAll(ctx)
}

func (s *SessionServiceImpl) UpdateSession(ctx context.Context, session *Session) error {
	return s.repo.Update(ctx, session)
}

func (s *SessionServiceImpl) DeleteSession(ctx context.Context, id string) error {
	if err := s.validateSessionID(id); err != nil {
		return err
	}


	_ = s.whatsappService.StopClient(id)

	return s.repo.Delete(ctx, id)
}




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


func (s *SessionServiceImpl) performConnection(ctx context.Context, session *Session) error {
	session.SetStatus(types.StatusConnecting)
	if err := s.repo.Update(ctx, session); err != nil {
		return err
	}

	return s.whatsappService.StartClient(session.ID)
}




func (s *SessionServiceImpl) GetQRCode(ctx context.Context, id string) (string, error) {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return "", err
	}

	if session.IsConnected() {
		return "", ErrSessionAlreadyConnected
	}


	if !session.HasQRCode() && session.CanConnect() {
		return s.whatsappService.GetQRCode(id)
	}

	return session.QRCode, nil
}


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




func (s *SessionServiceImpl) SetProxy(ctx context.Context, id, proxyURL string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	session.SetProxyURL(proxyURL)
	return s.repo.Update(ctx, session)
}


func (s *SessionServiceImpl) ClearProxy(ctx context.Context, id string) error {
	session, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}

	session.ClearProxy()
	return s.repo.Update(ctx, session)
}




func (s *SessionServiceImpl) ConnectOnStartup(ctx context.Context) error {
	return s.whatsappService.ConnectOnStartup(ctx)
}




func (s *SessionServiceImpl) validateSessionID(id string) error {
	return ValidateSessionID(id)
}


func (s *SessionServiceImpl) validateSessionName(name string) error {
	return ValidateSessionName(name)
}


func (s *SessionServiceImpl) validateSessionIDOrName(idOrName string) error {
	return ValidateSessionIDOrName(idOrName)
}




func generateSessionID() string {
	return uuid.New().String()
}
