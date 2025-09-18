package application

import "context"

// Response Types (eliminar interface{})
type MessageResult struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

type GroupResult struct {
	GroupJID string   `json:"groupJid"`
	Name     string   `json:"name"`
	Members  []string `json:"members"`
}

type UserInfo struct {
	JID           string `json:"jid"`
	Name          string `json:"name"`
	ProfilePicURL string `json:"profilePicUrl"`
	Status        string `json:"status"`
	IsBlocked     bool   `json:"isBlocked"`
}

type ChatInfo struct {
	JID           string `json:"jid"`
	Name          string `json:"name"`
	IsGroup       bool   `json:"isGroup"`
	LastMessage   string `json:"lastMessage"`
	LastTimestamp int64  `json:"lastTimestamp"`
	UnreadCount   int    `json:"unreadCount"`
	IsMuted       bool   `json:"isMuted"`
	IsPinned      bool   `json:"isPinned"`
	IsArchived    bool   `json:"isArchived"`
}

type EventData struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// Group Types
type GroupInfo struct {
	JID              string   `json:"jid"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Participants     []string `json:"participants"`
	Admins           []string `json:"admins"`
	Owner            string   `json:"owner"`
	CreatedAt        int64    `json:"createdAt"`
	IsAnnounce       bool     `json:"isAnnounce"`
	IsLocked         bool     `json:"isLocked"`
	IsEphemeral      bool     `json:"isEphemeral"`
	ParticipantCount int      `json:"participantCount"`
}

type GroupList struct {
	Groups []GroupInfo `json:"groups"`
	Total  int         `json:"total"`
}

type InviteInfo struct {
	GroupJID   string `json:"groupJid"`
	GroupName  string `json:"groupName"`
	InviteCode string `json:"inviteCode"`
	Inviter    string `json:"inviter"`
	ExpiresAt  int64  `json:"expiresAt"`
	IsValid    bool   `json:"isValid"`
}

// Newsletter Types
type NewsletterInfo struct {
	JID             string `json:"jid"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SubscriberCount int    `json:"subscriberCount"`
	CreatedAt       int64  `json:"createdAt"`
	IsVerified      bool   `json:"isVerified"`
}

type NewsletterList struct {
	Newsletters []NewsletterInfo `json:"newsletters"`
	Total       int              `json:"total"`
}

type NewsletterMessage struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	MediaURL  string `json:"mediaUrl,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}

type NewsletterMessages struct {
	Messages []NewsletterMessage `json:"messages"`
	Total    int                 `json:"total"`
}

// Media Upload Result
type MediaUploadResult struct {
	URL      string `json:"url"`
	MediaID  string `json:"mediaId"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}

// Event Processing Interfaces
type EventProcessor interface {
	HandleEvent(evt interface{})
}

type EventDispatcherInterface interface {
	DispatchEvent(ctx context.Context, sessionID string, eventType string, eventData interface{}) error
	ValidateEventType(eventType string) bool
}

// Infrastructure Interfaces (Dependency Inversion)
type WebhookSender interface {
	SendWebhook(ctx context.Context, sessionID, url, eventType string, payload interface{}) error
}

type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
	Warn(msg string)
	Warnf(format string, args ...interface{})
}

// ID Generator Interface (para abstrair geração de IDs)
type IDGenerator interface {
	GenerateSessionID() string
	GenerateAPIKey() string
}

// Messaging Interfaces
type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID, chatJID, content string) (*MessageResult, error)
	SendImageMessage(ctx context.Context, sessionID, chatJID string, imageData []byte, caption string) (*MessageResult, error)
	SendVideoMessage(ctx context.Context, sessionID, chatJID string, videoData []byte, caption string) (*MessageResult, error)
	SendAudioMessage(ctx context.Context, sessionID, chatJID string, audioData []byte) (*MessageResult, error)
	SendDocumentMessage(ctx context.Context, sessionID, chatJID string, documentData []byte, filename, mimetype string) (*MessageResult, error)
	SendStickerMessage(ctx context.Context, sessionID, chatJID string, stickerData []byte) (*MessageResult, error)
}

// Extended Message Sender interface for additional operations
type ExtendedMessageSender interface {
	MessageSender
	SendContactMessage(ctx context.Context, sessionID, chatJID, contactVCard string) (*MessageResult, error)
	SendLocationMessage(ctx context.Context, sessionID, chatJID string, latitude, longitude float64, name, address string) (*MessageResult, error)
	MarkAsRead(ctx context.Context, sessionID, chatJID, messageID string) error
	ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, reaction string) error
	EditMessage(ctx context.Context, sessionID, chatJID, messageID, newContent string) error
	DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string) error
}

// WhatsApp Service Interface
type WhatsAppService interface {
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	IsClientConnected(sessionID string) bool
	ConnectOnStartup(ctx context.Context) error
}

// Webhook Service Interface (movida de webhook.go)
type WebhookService interface {
	SetWebhook(ctx context.Context, sessionID, webhookURL string, events []string) error
	GetWebhook(ctx context.Context, sessionID string) (*WebhookInfo, error)
	UpdateWebhook(ctx context.Context, sessionID, webhookURL string, events []string, active bool) error
	DeleteWebhook(ctx context.Context, sessionID string) error
	TestWebhook(ctx context.Context, sessionID, message string) error
	ValidateWebhookURL(url string) error
	ValidateEvents(events []string) error
}

// Webhook Info Type (movida de webhook.go)
type WebhookInfo struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Active bool     `json:"active"`
}

// Complete Service Interfaces - covering all existing functionality
type GroupService interface {
	// Group management
	CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*GroupInfo, error)
	GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error)
	ListGroups(ctx context.Context, sessionID string) (*GroupList, error)
	JoinGroup(ctx context.Context, sessionID, inviteLink string) (*GroupInfo, error)
	JoinGroupWithInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error)
	LeaveGroup(ctx context.Context, sessionID, groupJID string) error

	// Group invites
	GetInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error)
	GetInviteInfo(ctx context.Context, sessionID, inviteLink string) (*InviteInfo, error)
	GetGroupInfoFromInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error)

	// Group participants
	UpdateParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
	GetGroupRequestParticipants(ctx context.Context, sessionID, groupJID string) ([]string, error)
	UpdateGroupRequestParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error

	// Group settings
	SetGroupName(ctx context.Context, sessionID, groupJID, name string) error
	SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error
	SetGroupPhoto(ctx context.Context, sessionID, groupJID string, photoData []byte) error
	RemoveGroupPhoto(ctx context.Context, sessionID, groupJID string) error
	SetGroupAnnounce(ctx context.Context, sessionID, groupJID string, announceOnly bool) error
	SetGroupLocked(ctx context.Context, sessionID, groupJID string, locked bool) error
	SetGroupEphemeral(ctx context.Context, sessionID, groupJID string, ephemeral bool, duration int) error
	SetGroupJoinApprovalMode(ctx context.Context, sessionID, groupJID string, requireApproval bool) error
	SetGroupMemberAddMode(ctx context.Context, sessionID, groupJID, mode string) error

	// Group linking
	LinkGroup(ctx context.Context, sessionID, groupJID, parentGroupJID string) error
	UnlinkGroup(ctx context.Context, sessionID, groupJID string) error
	GetSubGroups(ctx context.Context, sessionID, parentGroupJID string) (*GroupList, error)
	GetLinkedGroupsParticipants(ctx context.Context, sessionID, parentGroupJID string) ([]string, error)
}

type NewsletterService interface {
	// Newsletter management
	CreateNewsletter(ctx context.Context, sessionID, name, description string) (*NewsletterInfo, error)
	GetNewsletter(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
	ListNewsletters(ctx context.Context, sessionID string) (*NewsletterList, error)

	// Newsletter subscription
	SubscribeToNewsletter(ctx context.Context, sessionID, newsletterJID string) error
	UnsubscribeFromNewsletter(ctx context.Context, sessionID, newsletterJID string) error

	// Newsletter messaging
	SendNewsletterMessage(ctx context.Context, sessionID, newsletterJID, content string) (*MessageResult, error)
	GetNewsletterMessages(ctx context.Context, sessionID, newsletterJID string, limit int) (*NewsletterMessages, error)
	GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterJID string) (*NewsletterMessages, error)
	MarkNewsletterViewed(ctx context.Context, sessionID, newsletterJID, messageID string) error
	SendNewsletterReaction(ctx context.Context, sessionID, newsletterJID, messageID, reaction string) error

	// Newsletter settings
	ToggleNewsletterMute(ctx context.Context, sessionID, newsletterJID string, mute bool) error
	SubscribeLiveUpdates(ctx context.Context, sessionID, newsletterJID string) error

	// Newsletter media
	UploadNewsletterMedia(ctx context.Context, sessionID string, mediaData []byte, mediaType string) (*MediaUploadResult, error)
	GetNewsletterByInvite(ctx context.Context, sessionID, inviteKey string) (*NewsletterInfo, error)
}
