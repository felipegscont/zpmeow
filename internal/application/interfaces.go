package application

import "context"

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
	SendWebhook(ctx context.Context, sessionID, eventType string, payload interface{}) error
}

type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
	Warn(msg string)
	Warnf(format string, args ...interface{})
}

// Messaging Interfaces
type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID, chatJID, content string) (interface{}, error)
	SendImageMessage(ctx context.Context, sessionID, chatJID string, imageData []byte, caption string) (interface{}, error)
	SendVideoMessage(ctx context.Context, sessionID, chatJID string, videoData []byte, caption string) (interface{}, error)
	SendAudioMessage(ctx context.Context, sessionID, chatJID string, audioData []byte) (interface{}, error)
	SendDocumentMessage(ctx context.Context, sessionID, chatJID string, documentData []byte, filename, mimetype string) (interface{}, error)
	SendStickerMessage(ctx context.Context, sessionID, chatJID string, stickerData []byte) (interface{}, error)
}

// WhatsApp Service Interface
type WhatsAppService interface {
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	IsClientConnected(sessionID string) bool
}
