package session

import (
	"context"
)

// SessionManager defines the interface for session management operations
type SessionManager interface {
	CreateSession(ctx context.Context, name string) (*Session, error)
	GetSession(ctx context.Context, idOrName string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
	ListSessions(ctx context.Context, limit, offset int) ([]*Session, error)
	ConnectSession(ctx context.Context, id string) error
	DisconnectSession(ctx context.Context, id string) error
	GetSessionStatus(ctx context.Context, id string) (Status, error)
}

// MessageSender defines the interface for sending messages (domain abstraction)
type MessageSender interface {
	SendMessage(ctx context.Context, session *Session, message Message) error
	CanSendMessage(ctx context.Context, session *Session, message Message) error
}

// WebhookManager defines the interface for webhook management (domain abstraction)
type WebhookManager interface {
	RegisterWebhook(ctx context.Context, session *Session, webhook Webhook) error
	CanRegisterWebhook(ctx context.Context, session *Session, webhook Webhook) error
	ValidateWebhookConfiguration(webhook Webhook) error
}

// ApplicationSessionService combines session management with additional operations
type ApplicationSessionService interface {
	SessionManager
	PairSession(ctx context.Context, session *Session, phoneNumber string) error
	GetQRCode(ctx context.Context, session *Session) (string, error)
	SetProxy(ctx context.Context, session *Session, proxyURL string) error
	ClearProxy(ctx context.Context, session *Session) error
}

// Domain Value Objects and Entities

// Message represents a domain message concept
type Message struct {
	Content     string
	MessageType MessageType
	Recipient   string
}

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVideo    MessageType = "video"
	MessageTypeDocument MessageType = "document"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
	MessageTypePoll     MessageType = "poll"
)

// Webhook represents a domain webhook concept
type Webhook struct {
	URL    string
	Events []string
	Secret string
}
