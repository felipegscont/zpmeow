package ports

import (
	"context"
	"time"
)

// WhatsAppService defines the contract for WhatsApp operations
// This abstracts the external WhatsApp integration
type WhatsAppService interface {
	// Session Management
	ConnectSession(ctx context.Context, sessionID string) error
	DisconnectSession(ctx context.Context, sessionID string) error
	GetSessionStatus(ctx context.Context, sessionID string) (string, error)
	PairWithPhone(ctx context.Context, sessionID, phoneNumber string) error
	GetQRCode(ctx context.Context, sessionID string) (string, error)

	// Messaging
	SendTextMessage(ctx context.Context, sessionID, chatJID, message string) error
	SendMediaMessage(ctx context.Context, sessionID, chatJID string, media MediaMessage) error
	SendLocationMessage(ctx context.Context, sessionID, chatJID string, latitude, longitude float64, name, address string) error
	SendContactMessage(ctx context.Context, sessionID, chatJID string, contacts []ContactInfo) error
	MarkAsRead(ctx context.Context, sessionID, chatJID, messageID string) error
	ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string, remove bool) error
	EditMessage(ctx context.Context, sessionID, chatJID, messageID, newContent string) error
	DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error

	// Chat Management
	GetChats(ctx context.Context, sessionID string, limit, offset int) ([]ChatInfo, error)
	GetChatHistory(ctx context.Context, sessionID, chatJID string, limit, offset int) ([]MessageInfo, error)
	SetPresence(ctx context.Context, sessionID, chatJID, state, media string) error
	MuteChat(ctx context.Context, sessionID, chatJID string, mute bool, duration time.Duration) error
	ArchiveChat(ctx context.Context, sessionID, chatJID string, archive bool) error
	BlockContact(ctx context.Context, sessionID, contactJID string, block bool) error

	// Contact Management
	GetContacts(ctx context.Context, sessionID string, limit, offset int) ([]ContactInfo, error)
	CheckContact(ctx context.Context, sessionID, phone string) (bool, string, error)
	GetUserInfo(ctx context.Context, sessionID, userJID string) (*UserInfo, error)

	// Group Management
	CreateGroup(ctx context.Context, sessionID, name string, participants []string) (string, error)
	JoinGroup(ctx context.Context, sessionID, inviteLink string) error
	LeaveGroup(ctx context.Context, sessionID, groupJID string) error
	GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error)
	ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error)
	AddParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	RemoveParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	GetGroupInviteLink(ctx context.Context, sessionID, groupJID string) (string, error)

	// Newsletter Management
	CreateNewsletter(ctx context.Context, sessionID, name, description string) (string, error)
	GetNewsletterInfo(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
	SubscribeNewsletter(ctx context.Context, sessionID, newsletterJID string) error
	UnsubscribeNewsletter(ctx context.Context, sessionID, newsletterJID string) error
}

// MediaMessage represents a media message to be sent
type MediaMessage struct {
	Type     string // image, video, audio, document
	Data     []byte
	MimeType string
	Caption  string
	Filename string
}

// ContactInfo represents contact information
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

// ChatInfo represents chat information
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

// GroupInfo represents group information
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

// MessageInfo represents message information
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

// UserInfo represents detailed user information
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

// NewsletterInfo represents newsletter information
type NewsletterInfo struct {
	JID             string
	Name            string
	Description     string
	SubscriberCount int
	IsSubscribed    bool
	CreatedAt       string
	UpdatedAt       string
}

// NotificationService defines the contract for sending notifications
type NotificationService interface {
	// SendWebhook sends a webhook notification
	SendWebhook(ctx context.Context, url string, payload interface{}) error

	// SendEmail sends an email notification
	SendEmail(ctx context.Context, to, subject, body string) error
}

// IDGenerator defines the contract for generating unique identifiers
type IDGenerator interface {
	// GenerateSessionID generates a unique session identifier
	GenerateSessionID() string

	// GenerateAPIKey generates a unique API key
	GenerateAPIKey() string
}

// Logger defines the contract for logging
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...interface{})
	Info(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
}

// Validator defines the contract for input validation
type Validator interface {
	// Validate validates a struct and returns validation errors
	Validate(v interface{}) error

	// ValidateField validates a single field
	ValidateField(field interface{}, tag string) error
}

// MessageSender defines the contract for sending messages
type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID, chatJID, content string) error
	SendMediaMessage(ctx context.Context, sessionID, chatJID string, media MediaMessage) error
}

// ChatService defines the contract for chat operations
type ChatService interface {
	GetChats(ctx context.Context, sessionID string, limit, offset int) ([]ChatInfo, error)
	GetChatHistory(ctx context.Context, sessionID, chatJID string, limit, offset int) ([]MessageInfo, error)
	SetPresence(ctx context.Context, sessionID, chatJID, state, media string) error
}

// GroupService defines the contract for group operations
type GroupService interface {
	CreateGroup(ctx context.Context, sessionID, name string, participants []string) (string, error)
	GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error)
	ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error)
}

// ContactService defines the contract for contact operations
type ContactService interface {
	GetContacts(ctx context.Context, sessionID string, limit, offset int) ([]ContactInfo, error)
	CheckContact(ctx context.Context, sessionID, phone string) (bool, string, error)
	GetUserInfo(ctx context.Context, sessionID, userJID string) (*UserInfo, error)
}

// NewsletterService defines the contract for newsletter operations
type NewsletterService interface {
	CreateNewsletter(ctx context.Context, sessionID, name, description string) (string, error)
	GetNewsletterInfo(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
	SubscribeNewsletter(ctx context.Context, sessionID, newsletterJID string) error
}
