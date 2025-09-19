package application

import (
	"context"
	"fmt"

	"zpmeow/internal/application/ports"
	"zpmeow/internal/domain/session"
)

type SessionApp struct {
	sessionRepo   session.Repository
	domainService session.Service
}

func NewSessionApp(sessionRepo session.Repository, domainService session.Service) *SessionApp {
	return &SessionApp{
		sessionRepo:   sessionRepo,
		domainService: domainService,
	}
}

func (s *SessionApp) GetSession(ctx context.Context, sessionIDOrName string) (*session.Session, error) {
	sess, err := s.sessionRepo.GetByID(ctx, sessionIDOrName)
	if err == nil {
		return sess, nil
	}

	return s.sessionRepo.GetByName(ctx, sessionIDOrName)
}

func (s *SessionApp) GetAllSessions(ctx context.Context) ([]*session.Session, error) {
	return s.sessionRepo.GetAll(ctx)
}

func (s *SessionApp) CreateSessionWithRequest(ctx context.Context, req CreateSessionRequest) (*session.Session, error) {
	sess, err := session.NewSession("", req.Name)
	if err != nil {
		return nil, err
	}

	generatedID, err := s.sessionRepo.CreateWithGeneratedID(ctx, sess)
	if err != nil {
		return nil, err
	}

	return s.sessionRepo.GetByID(ctx, generatedID)
}

func (s *SessionApp) DeleteSession(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *SessionApp) GetSessionByDeviceJID(ctx context.Context, deviceJID string) (*session.Session, error) {
	return nil, fmt.Errorf("GetSessionByDeviceJID not implemented")
}

type CreateSessionRequest struct {
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
}

type WebhookApp struct {
	sessionRepo session.Repository
}

func NewWebhookApp(sessionRepo session.Repository) *WebhookApp {
	return &WebhookApp{
		sessionRepo: sessionRepo,
	}
}

func (w *WebhookApp) SetWebhook(ctx context.Context, sessionID, webhookURL string, events []string) error {
	sess, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if err := sess.SetWebhookEndpoint(webhookURL); err != nil {
		return err
	}

	return w.sessionRepo.Update(ctx, sess)
}

func (w *WebhookApp) GetWebhook(ctx context.Context, sessionID string) (string, []string, error) {
	sess, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", nil, err
	}

	webhookURL := sess.GetWebhookEndpointString()
	events := []string{}

	return webhookURL, events, nil
}

func (w *WebhookApp) ListEvents(ctx context.Context) ([]string, error) {
	events := []string{
		"message",
		"message_ack",
		"message_revoke",
		"presence",
		"chat_presence",
		"connected",
		"disconnected",
		"qr",
		"group_join",
		"group_leave",
	}
	return events, nil
}

type MessageApp struct {
	sessionRepo   session.Repository
	messageSender ports.MessageSender
}

func NewMessageApp(sessionRepo session.Repository, messageSender ports.MessageSender) *MessageApp {
	return &MessageApp{
		sessionRepo:   sessionRepo,
		messageSender: messageSender,
	}
}

type ChatApp struct {
	sessionRepo session.Repository
	chatService ports.ChatService
}

func NewChatApp(sessionRepo session.Repository, chatService ports.ChatService) *ChatApp {
	return &ChatApp{
		sessionRepo: sessionRepo,
		chatService: chatService,
	}
}

type GroupApp struct {
	sessionRepo  session.Repository
	groupService ports.GroupService
}

func NewGroupApp(sessionRepo session.Repository, groupService ports.GroupService) *GroupApp {
	return &GroupApp{
		sessionRepo:  sessionRepo,
		groupService: groupService,
	}
}

type ContactApp struct {
	sessionRepo    session.Repository
	contactService ports.ContactService
}

func NewContactApp(sessionRepo session.Repository, contactService ports.ContactService) *ContactApp {
	return &ContactApp{
		sessionRepo:    sessionRepo,
		contactService: contactService,
	}
}

type NewsletterApp struct {
	sessionRepo       session.Repository
	newsletterService ports.NewsletterService
}

func NewNewsletterApp(sessionRepo session.Repository, newsletterService ports.NewsletterService) *NewsletterApp {
	return &NewsletterApp{
		sessionRepo:       sessionRepo,
		newsletterService: newsletterService,
	}
}
