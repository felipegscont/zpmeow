package patterns

import (
	"context"

	"go.mau.fi/whatsmeow"
)

// MediaStrategy defines the interface for different media type handling strategies
type MediaStrategy interface {
	// ValidateMedia validates the media data for this specific type
	ValidateMedia(data []byte) error

	// ProcessMedia processes the media data and returns the processed data and mime type
	ProcessMedia(ctx context.Context, data []byte, filename string) ([]byte, string, error)

	// SendMessage sends the media message using the WhatsApp service
	SendMessage(ctx context.Context, service MediaSender, sessionID, phone string, data []byte, caption, filename, mimeType string) (*whatsmeow.SendResponse, error)
}

// MediaSender interface for sending different types of media
type MediaSender interface {
	SendImageMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error)
	SendAudioMessage(ctx context.Context, sessionID, phone string, data []byte, mimeType string) (*whatsmeow.SendResponse, error)
	SendVideoMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error)
	SendDocumentMessage(ctx context.Context, sessionID, phone string, data []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error)
}

// MediaStrategyFactory creates media strategies based on media type
type MediaStrategyFactory struct{}

// CreateStrategy creates a strategy for the given media type
func (f *MediaStrategyFactory) CreateStrategy(mediaType string) MediaStrategy {
	// Implementation will be provided by application layer
	return nil
}
