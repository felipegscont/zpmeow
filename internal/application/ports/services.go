package ports

import (
	"context"
	"time"
)

type WhatsAppService interface {
	ConnectSession(ctx context.Context, sessionID string) error
	DisconnectSession(ctx context.Context, sessionID string) error
	GetSessionStatus(ctx context.Context, sessionID string) (string, error)
	PairWithPhone(ctx context.Context, sessionID, phoneNumber string) error
	GetQRCode(ctx context.Context, sessionID string) (string, error)

	SendTextMessage(ctx context.Context, sessionID, chatJID, message string) error
	SendMediaMessage(ctx context.Context, sessionID, chatJID string, media MediaMessage) error
	SendLocationMessage(ctx context.Context, sessionID, chatJID string, latitude, longitude float64, name, address string) error
	SendContactMessage(ctx context.Context, sessionID, chatJID string, contacts []ContactInfo) error
	MarkAsRead(ctx context.Context, sessionID, chatJID, messageID string) error
	ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string, remove bool) error
	EditMessage(ctx context.Context, sessionID, chatJID, messageID, newContent string) error
	DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error

	GetChats(ctx context.Context, sessionID string, limit, offset int) ([]ChatInfo, error)
	GetChatHistory(ctx context.Context, sessionID, chatJID string, limit, offset int) ([]MessageInfo, error)
	SetPresence(ctx context.Context, sessionID, chatJID, state, media string) error
	MuteChat(ctx context.Context, sessionID, chatJID string, mute bool, duration time.Duration) error
	ArchiveChat(ctx context.Context, sessionID, chatJID string, archive bool) error
	BlockContact(ctx context.Context, sessionID, contactJID string, block bool) error

	GetContacts(ctx context.Context, sessionID string, limit, offset int) ([]ContactInfo, error)
	CheckContact(ctx context.Context, sessionID, phone string) (bool, string, error)
	GetUserInfo(ctx context.Context, sessionID, userJID string) (*UserInfo, error)

	CreateGroup(ctx context.Context, sessionID, name string, participants []string) (string, error)
	JoinGroup(ctx context.Context, sessionID, inviteLink string) error
	LeaveGroup(ctx context.Context, sessionID, groupJID string) error
	GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error)
	ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error)
	AddParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	RemoveParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	GetGroupInviteLink(ctx context.Context, sessionID, groupJID string) (string, error)

	CreateNewsletter(ctx context.Context, sessionID, name, description string) (string, error)
	GetNewsletterInfo(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
	SubscribeNewsletter(ctx context.Context, sessionID, newsletterJID string) error
	UnsubscribeNewsletter(ctx context.Context, sessionID, newsletterJID string) error
}

type MediaMessage struct {
	Type     string // image, video, audio, document
	Data     []byte
	MimeType string
	Caption  string
	Filename string
}

type ContactInfo struct {
	JID          string
	Name         string
	Notify       string
	PushName     string
	BusinessName string
	Phone        string
	Organization string
	Email        string
	IsBlocked    bool
	IsMuted      bool
	IsContact    bool
	Avatar       string
}

type ChatInfo struct {
	JID           string
	Name          string
	LastMessage   string
	LastMessageAt string
	UnreadCount   int
	IsGroup       bool
	IsMuted       bool
	IsArchived    bool
	IsBlocked     bool
}

type GroupInfo struct {
	JID          string
	Name         string
	Description  string
	Participants []string
	Admins       []string
	Owner        string
	CreatedAt    string
	IsAnnounce   bool
	IsLocked     bool
}

type MessageInfo struct {
	ID        string
	ChatJID   string
	FromJID   string
	Content   string
	Type      string
	Timestamp string
	IsFromMe  bool
	IsRead    bool
	MediaURL  string
	Caption   string
}

type UserInfo struct {
	JID          string
	Name         string
	Notify       string
	PushName     string
	BusinessName string
	Phone        string
	Status       string
	Avatar       string
	IsBlocked    bool
	IsMuted      bool
	IsContact    bool
	LastSeen     string
}

type NewsletterInfo struct {
	JID             string
	Name            string
	Description     string
	SubscriberCount int
	IsSubscribed    bool
	CreatedAt       string
	UpdatedAt       string
}

type NotificationService interface {
	SendWebhook(ctx context.Context, url string, payload interface{}) error

	SendEmail(ctx context.Context, to, subject, body string) error
}

type IDGenerator interface {
	GenerateSessionID() string

	GenerateAPIKey() string
}

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...interface{})
	Info(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
}

type Validator interface {
	Validate(v interface{}) error

	ValidateField(field interface{}, tag string) error
}

type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID, chatJID, content string) error
	SendMediaMessage(ctx context.Context, sessionID, chatJID string, media MediaMessage) error
}

type ChatService interface {
	GetChats(ctx context.Context, sessionID string, limit, offset int) ([]ChatInfo, error)
	GetChatHistory(ctx context.Context, sessionID, chatJID string, limit, offset int) ([]MessageInfo, error)
	SetPresence(ctx context.Context, sessionID, chatJID, state, media string) error
}

type GroupService interface {
	CreateGroup(ctx context.Context, sessionID, name string, participants []string) (string, error)
	GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error)
	ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error)
}

type ContactService interface {
	GetContacts(ctx context.Context, sessionID string, limit, offset int) ([]ContactInfo, error)
	CheckContact(ctx context.Context, sessionID, phone string) (bool, string, error)
	GetUserInfo(ctx context.Context, sessionID, userJID string) (*UserInfo, error)
}

type NewsletterService interface {
	CreateNewsletter(ctx context.Context, sessionID, name, description string) (string, error)
	GetNewsletterInfo(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
	SubscribeNewsletter(ctx context.Context, sessionID, newsletterJID string) error
}
