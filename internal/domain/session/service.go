package session

import (
	"context"
	"zpmeow/internal/shared/types"
)

// SessionService defines the interface for session business logic
type SessionService interface {
	// Session CRUD operations
	CreateSession(ctx context.Context, name string) (*Session, error)
	GetSession(ctx context.Context, idOrName string) (*Session, error)
	GetAllSessions(ctx context.Context) ([]*Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id string) error

	// Session connection operations
	ConnectSession(ctx context.Context, id string) error
	DisconnectSession(ctx context.Context, id string) error

	// Session pairing operations
	GetQRCode(ctx context.Context, id string) (string, error)
	PairWithPhone(ctx context.Context, id, phoneNumber string) (string, error)

	// Proxy operations
	SetProxy(ctx context.Context, id, proxyURL string) error
	ClearProxy(ctx context.Context, id string) error

	// Startup operations
	ConnectOnStartup(ctx context.Context) error
}

// WhatsAppService defines the interface for WhatsApp operations
type WhatsAppService interface {
	// Client lifecycle
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	IsClientConnected(sessionID string) bool
	GetClientStatus(sessionID string) types.Status
	ConnectOnStartup(ctx context.Context) error

	// Message operations
	DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error
	EditMessage(ctx context.Context, sessionID, chatJID, messageID, newText string) (*types.SendResponse, error)
	DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error)
	ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string) error
}
