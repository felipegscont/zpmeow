package session

import (
	"context"
)

// WhatsAppService defines the interface for WhatsApp operations
// This interface abstracts the WhatsApp SDK implementation from the domain layer
type WhatsAppService interface {
	// Session Management
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	
	// Connection Status
	IsClientConnected(sessionID string) bool
	GetClientStatus(sessionID string) Status
	
	// Authentication
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	
	// Startup Operations
	ConnectOnStartup(ctx context.Context) error
	
	// Message Operations
	SendTextMessage(ctx context.Context, sessionID, phone, text string) (*MessageResponse, error)
	SendImageMessage(ctx context.Context, sessionID, phone string, imageData []byte, caption, mimeType string) (*MessageResponse, error)
	SendAudioMessage(ctx context.Context, sessionID, phone string, audioData []byte, mimeType string) (*MessageResponse, error)
	SendVideoMessage(ctx context.Context, sessionID, phone string, videoData []byte, caption, mimeType string) (*MessageResponse, error)
	SendDocumentMessage(ctx context.Context, sessionID, phone string, documentData []byte, filename, caption, mimeType string) (*MessageResponse, error)
	SendLocationMessage(ctx context.Context, sessionID, phone string, latitude, longitude float64, name, address string) (*MessageResponse, error)
	SendContactMessage(ctx context.Context, sessionID, phone, contactName, contactJID string) (*MessageResponse, error)
	
	// Group Operations
	CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*GroupResponse, error)
	AddParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	RemoveParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	PromoteParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	DemoteParticipants(ctx context.Context, sessionID, groupJID string, participants []string) error
	LeaveGroup(ctx context.Context, sessionID, groupJID string) error
	SetGroupName(ctx context.Context, sessionID, groupJID, name string) error
	SetGroupDescription(ctx context.Context, sessionID, groupJID, description string) error
	SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error
	SetGroupPhoto(ctx context.Context, sessionID, groupJID string, photo []byte) error
	RemoveGroupPhoto(ctx context.Context, sessionID, groupJID string) error
	SetGroupAnnounce(ctx context.Context, sessionID, groupJID string, announceOnly bool) error
	SetGroupLocked(ctx context.Context, sessionID, groupJID string, locked bool) error
	SetGroupEphemeral(ctx context.Context, sessionID, groupJID string, ephemeral bool, duration int) error
	
	// User Operations
	GetUserInfo(ctx context.Context, sessionID, userJID string) (*UserInfo, error)
	GetUserProfilePicture(ctx context.Context, sessionID, userJID string) (string, error)
	BlockUser(ctx context.Context, sessionID, userJID string) error
	UnblockUser(ctx context.Context, sessionID, userJID string) error
	
	// Chat Operations
	GetChats(ctx context.Context, sessionID string) ([]*ChatInfo, error)
	GetChatHistory(ctx context.Context, sessionID, chatJID string, limit int) ([]*MessageInfo, error)
	MarkAsRead(ctx context.Context, sessionID, chatJID, messageID string) error
	ArchiveChat(ctx context.Context, sessionID, chatJID string) error
	UnarchiveChat(ctx context.Context, sessionID, chatJID string) error
	DeleteChat(ctx context.Context, sessionID, chatJID string) error
	MuteChat(ctx context.Context, sessionID, chatJID string, duration int) error
	UnmuteChat(ctx context.Context, sessionID, chatJID string) error
	PinChat(ctx context.Context, sessionID, chatJID string) error
	UnpinChat(ctx context.Context, sessionID, chatJID string) error
}

// MessageResponse represents the response from sending a message
type MessageResponse struct {
	ID        string
	Timestamp int64
	Status    string
}

// GroupResponse represents the response from group operations
type GroupResponse struct {
	GroupJID string
	Name     string
	Members  []string
}

// UserInfo represents user information
type UserInfo struct {
	JID           string
	Name          string
	ProfilePicURL string
	Status        string
	IsBlocked     bool
}

// ChatInfo represents chat information
type ChatInfo struct {
	JID           string
	Name          string
	IsGroup       bool
	LastMessage   string
	LastTimestamp int64
	UnreadCount   int
	IsMuted       bool
	IsPinned      bool
	IsArchived    bool
}

// MessageInfo represents message information
type MessageInfo struct {
	ID        string
	FromJID   string
	ToJID     string
	Content   string
	Type      string
	Timestamp int64
	IsFromMe  bool
}
