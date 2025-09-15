package common

import (
	"context"
	"time"
)

type ValidatorInterface interface {
	Validate(req interface{}) error
	ValidatePhoneNumber(phone string) error
	ValidatePhoneNumbers(phones []string) error
	ValidateMessageID(messageID string) error
	ValidateMessageIDs(messageIDs []string) error
}

type SessionManager interface {
	CreateSession(ctx context.Context, name string) (string, error)
	GetSession(ctx context.Context, idOrName string) (SessionInfo, error)
	DeleteSession(ctx context.Context, id string) error
	ListSessions(ctx context.Context, limit, offset int) ([]SessionInfo, error)
	ConnectSession(ctx context.Context, id string) error
	DisconnectSession(ctx context.Context, id string) error
	GetSessionStatus(ctx context.Context, id string) (string, error)
}

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

type WebhookManager interface {
	RegisterWebhook(ctx context.Context, sessionID, url string, events []string) (WebhookInfo, error)
	UpdateWebhook(ctx context.Context, webhookID, url string, events []string, active bool) error
	DeleteWebhook(ctx context.Context, webhookID string) error
	GetWebhook(ctx context.Context, webhookID string) (WebhookInfo, error)
	ListWebhooks(ctx context.Context, sessionID string, active *bool, limit, offset int) ([]WebhookInfo, error)
	NotifyWebhook(ctx context.Context, webhookID string, event string, data interface{}) error
}

type SessionRepository interface {
	Create(ctx context.Context, session SessionEntity) error
	GetByID(ctx context.Context, id string) (SessionEntity, error)
	GetByName(ctx context.Context, name string) (SessionEntity, error)
	Update(ctx context.Context, session SessionEntity) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]SessionEntity, error)
	Exists(ctx context.Context, id string) (bool, error)
}

type MessageRepository interface {
	Save(ctx context.Context, message MessageEntity) error
	GetByID(ctx context.Context, id string) (MessageEntity, error)
	GetBySessionID(ctx context.Context, sessionID string, limit, offset int) ([]MessageEntity, error)
	UpdateStatus(ctx context.Context, id string, status MessageStatus) error
	Delete(ctx context.Context, id string) error
}

type WebhookRepository interface {
	Create(ctx context.Context, webhook WebhookEntity) error
	GetByID(ctx context.Context, id string) (WebhookEntity, error)
	GetBySessionID(ctx context.Context, sessionID string) ([]WebhookEntity, error)
	Update(ctx context.Context, webhook WebhookEntity) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, sessionID string, active *bool, limit, offset int) ([]WebhookEntity, error)
}

type SessionService interface {
	SessionManager
	PairSession(ctx context.Context, sessionID, phoneNumber string) (string, error)
	GetQRCode(ctx context.Context, sessionID string) (string, error)
	SetProxy(ctx context.Context, sessionID, proxyURL string) error
	ClearProxy(ctx context.Context, sessionID string) error
}

type MessagingService interface {
	MessageSender
	CheckUsers(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error)
	GetUserInfo(ctx context.Context, sessionID string, phones []string) ([]UserInfo, error)
	SetPresence(ctx context.Context, sessionID string, presence string) error
}

type WebhookService interface {
	WebhookManager
	TriggerWebhook(ctx context.Context, sessionID, event string, data interface{}) error
	ValidateWebhookURL(url string) error
}

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

type SessionEntity struct {
	ID         string
	Name       string
	Status     string
	WaJID      string
	QRCode     string
	ProxyURL   string
	WebhookURL string
	Events     []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MessageResult struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	RemoteJID string    `json:"remote_jid"`
	Type      string    `json:"type"`
	Error     string    `json:"error,omitempty"`
}

type MessageEntity struct {
	ID        string
	SessionID string
	RemoteJID string
	Type      string
	Content   string
	Status    MessageStatus
	Timestamp time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

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

type UserInfo struct {
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

type WebhookEntity struct {
	ID            string
	SessionID     string
	URL           string
	Events        []string
	Active        bool
	Secret        string
	LastTriggered *time.Time
	FailureCount  int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EventPublisher interface {
	PublishSessionEvent(ctx context.Context, sessionID, event string, data interface{}) error
	PublishMessageEvent(ctx context.Context, sessionID, event string, data interface{}) error
	PublishContactEvent(ctx context.Context, sessionID, event string, data interface{}) error
	PublishCallEvent(ctx context.Context, sessionID, event string, data interface{}) error
	PublishGroupEvent(ctx context.Context, sessionID, event string, data interface{}) error
}

type EventSubscriber interface {
	Subscribe(ctx context.Context, events []string, handler EventHandler) error
	Unsubscribe(ctx context.Context, events []string) error
}

type EventHandler interface {
	Handle(ctx context.Context, event string, data interface{}) error
}

type ConfigProvider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetDuration(key string) time.Duration
}

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) ([]byte, error)
	Post(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Put(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Delete(ctx context.Context, url string, headers map[string]string) ([]byte, error)
}
