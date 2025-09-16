package session

import (
	"context"
	"time"
)

// SessionManager defines the interface for session management operations
type SessionManager interface {
	CreateSession(ctx context.Context, name string) (string, error)
	GetSession(ctx context.Context, idOrName string) (SessionInfo, error)
	DeleteSession(ctx context.Context, id string) error
	ListSessions(ctx context.Context, limit, offset int) ([]SessionInfo, error)
	ConnectSession(ctx context.Context, id string) error
	DisconnectSession(ctx context.Context, id string) error
	GetSessionStatus(ctx context.Context, id string) (string, error)
}

// MessageSender defines the interface for sending messages
type MessageSender interface {
	SendText(ctx context.Context, sessionID, phone, text string) (MessageResult, error)
	SendImage(ctx context.Context, sessionID, phone string, image []byte, caption string) (MessageResult, error)
	SendAudio(ctx context.Context, sessionID, phone string, audio []byte, caption string) (MessageResult, error)
	SendVideo(ctx context.Context, sessionID, phone string, video []byte, caption string) (MessageResult, error)
	SendDocument(ctx context.Context, sessionID, phone string, document []byte, filename, caption string) (MessageResult, error)
	SendSticker(ctx context.Context, sessionID, phone string, sticker []byte) (MessageResult, error)
	SendLocation(ctx context.Context, sessionID, phone string, latitude, longitude float64, name, address string) (MessageResult, error)
	SendContact(ctx context.Context, sessionID, phone string, contact ContactInfo) (MessageResult, error)
	SendPoll(ctx context.Context, sessionID, phone, question string, options []string, selectableCount int) (MessageResult, error)
}

// WebhookManager defines the interface for webhook management
type WebhookManager interface {
	RegisterWebhook(ctx context.Context, sessionID, url string, events []string) (WebhookInfo, error)
	UpdateWebhook(ctx context.Context, webhookID, url string, events []string, active bool) error
	DeleteWebhook(ctx context.Context, webhookID string) error
	GetWebhook(ctx context.Context, webhookID string) (WebhookInfo, error)
	ListWebhooks(ctx context.Context, sessionID string, active *bool, limit, offset int) ([]WebhookInfo, error)
	NotifyWebhook(ctx context.Context, webhookID string, event string, data interface{}) error
}

// ApplicationSessionService combines session management with additional operations
type ApplicationSessionService interface {
	SessionManager
	PairSession(ctx context.Context, sessionID, phoneNumber string) (string, error)
	GetQRCode(ctx context.Context, sessionID string) (string, error)
	SetProxy(ctx context.Context, sessionID, proxyURL string) error
	ClearProxy(ctx context.Context, sessionID string) error
}

// ApplicationMessagingService combines message sending with user operations
type ApplicationMessagingService interface {
	MessageSender
	CheckUsers(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error)
	GetUserInfo(ctx context.Context, sessionID string, phones []string) ([]ApplicationUserInfo, error)
	SetPresence(ctx context.Context, sessionID string, presence string) error
}

// ApplicationWebhookService combines webhook management with additional operations
type ApplicationWebhookService interface {
	WebhookManager
	TriggerWebhook(ctx context.Context, sessionID, event string, data interface{}) error
	ValidateWebhookURL(url string) error
}

// DTOs and Value Objects

type SessionInfo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	WaJID      string    `json:"wa_jid,omitempty"`
	QRCode     string    `json:"qr_code,omitempty"`
	ProxyURL   string    `json:"proxy_url,omitempty"`
	WebhookURL string    `json:"webhook_url,omitempty"`
	Events     []string  `json:"events,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MessageResult struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	RemoteJID string    `json:"remote_jid"`
	Type      string    `json:"type"`
	Error     string    `json:"error,omitempty"`
}

type ContactInfo struct {
	DisplayName string `json:"display_name"`
	VCard       string `json:"vcard"`
}

type UserCheckResult struct {
	Query        string `json:"query"`
	IsInWhatsapp bool   `json:"is_in_whatsapp"`
	JID          string `json:"jid"`
	VerifiedName string `json:"verified_name,omitempty"`
}

type ApplicationUserInfo struct {
	JID          string `json:"jid"`
	DisplayName  string `json:"display_name"`
	VerifiedName string `json:"verified_name,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Status       string `json:"status,omitempty"`
}

type WebhookInfo struct {
	ID            string     `json:"id"`
	SessionID     string     `json:"session_id"`
	URL           string     `json:"url"`
	Events        []string   `json:"events"`
	Active        bool       `json:"active"`
	Secret        string     `json:"secret,omitempty"`
	LastTriggered *time.Time `json:"last_triggered,omitempty"`
	FailureCount  int        `json:"failure_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
