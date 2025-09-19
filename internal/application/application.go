package application

import (
	"context"
	"fmt"

	"zpmeow/internal/application/ports"
	"zpmeow/internal/domain/session"
)

// SessionApp represents the session application service
type SessionApp struct {
	sessionRepo   session.Repository
	domainService session.Service
}

// NewSessionApp creates a new session application service
func NewSessionApp(sessionRepo session.Repository, domainService session.Service) *SessionApp {
	return &SessionApp{
		sessionRepo:   sessionRepo,
		domainService: domainService,
	}
}

// GetSession retrieves a session by ID or name
func (s *SessionApp) GetSession(ctx context.Context, sessionIDOrName string) (*session.Session, error) {
	// Try to get by ID first
	sess, err := s.sessionRepo.GetByID(ctx, sessionIDOrName)
	if err == nil {
		return sess, nil
	}

	// If not found by ID, try by name
	return s.sessionRepo.GetByName(ctx, sessionIDOrName)
}

// GetAllSessions retrieves all sessions
func (s *SessionApp) GetAllSessions(ctx context.Context) ([]*session.Session, error) {
	return s.sessionRepo.GetAll(ctx)
}

// CreateSessionWithRequest creates a new session with the given request
func (s *SessionApp) CreateSessionWithRequest(ctx context.Context, req CreateSessionRequest) (*session.Session, error) {
	// Create new session with empty ID (PostgreSQL will generate it)
	sess, err := session.NewSession("", req.Name)
	if err != nil {
		return nil, err
	}

	// Save to repository and get the generated ID
	generatedID, err := s.sessionRepo.CreateWithGeneratedID(ctx, sess)
	if err != nil {
		return nil, err
	}

	// Get the session back from the database with the generated ID
	return s.sessionRepo.GetByID(ctx, generatedID)
}

// DeleteSession deletes a session
func (s *SessionApp) DeleteSession(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

// GetSessionByDeviceJID retrieves a session by device JID
func (s *SessionApp) GetSessionByDeviceJID(ctx context.Context, deviceJID string) (*session.Session, error) {
	// This would need to be implemented in the repository
	// For now, return an error
	return nil, fmt.Errorf("GetSessionByDeviceJID not implemented")
}

// CreateSessionRequest represents a request to create a session
type CreateSessionRequest struct {
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
}

// WebhookApp represents the webhook application service
type WebhookApp struct {
	sessionRepo session.Repository
}

// NewWebhookApp creates a new webhook application service
func NewWebhookApp(sessionRepo session.Repository) *WebhookApp {
	return &WebhookApp{
		sessionRepo: sessionRepo,
	}
}

// SetWebhook sets webhook configuration for a session
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

// GetWebhook gets webhook configuration for a session
func (w *WebhookApp) GetWebhook(ctx context.Context, sessionID string) (string, []string, error) {
	sess, err := w.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", nil, err
	}

	webhookURL := sess.GetWebhookEndpointString()
	// For now, return empty events list as it's not stored in domain
	events := []string{}

	return webhookURL, events, nil
}

// ListEvents returns available webhook events
func (w *WebhookApp) ListEvents(ctx context.Context) ([]string, error) {
	// Return standard webhook events
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

// MessageApp represents the message application service
type MessageApp struct {
	sessionRepo   session.Repository
	messageSender ports.MessageSender
}

// NewMessageApp creates a new message application service
func NewMessageApp(sessionRepo session.Repository, messageSender ports.MessageSender) *MessageApp {
	return &MessageApp{
		sessionRepo:   sessionRepo,
		messageSender: messageSender,
	}
}

// ChatApp represents the chat application service
type ChatApp struct {
	sessionRepo session.Repository
	chatService ports.ChatService
}

// NewChatApp creates a new chat application service
func NewChatApp(sessionRepo session.Repository, chatService ports.ChatService) *ChatApp {
	return &ChatApp{
		sessionRepo: sessionRepo,
		chatService: chatService,
	}
}

// GroupApp represents the group application service
type GroupApp struct {
	sessionRepo  session.Repository
	groupService ports.GroupService
}

// NewGroupApp creates a new group application service
func NewGroupApp(sessionRepo session.Repository, groupService ports.GroupService) *GroupApp {
	return &GroupApp{
		sessionRepo:  sessionRepo,
		groupService: groupService,
	}
}

// ContactApp represents the contact application service
type ContactApp struct {
	sessionRepo    session.Repository
	contactService ports.ContactService
}

// NewContactApp creates a new contact application service
func NewContactApp(sessionRepo session.Repository, contactService ports.ContactService) *ContactApp {
	return &ContactApp{
		sessionRepo:    sessionRepo,
		contactService: contactService,
	}
}

// NewsletterApp represents the newsletter application service
type NewsletterApp struct {
	sessionRepo       session.Repository
	newsletterService ports.NewsletterService
}

// NewNewsletterApp creates a new newsletter application service
func NewNewsletterApp(sessionRepo session.Repository, newsletterService ports.NewsletterService) *NewsletterApp {
	return &NewsletterApp{
		sessionRepo:       sessionRepo,
		newsletterService: newsletterService,
	}
}
