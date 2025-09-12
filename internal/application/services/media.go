package services

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"zpmeow/internal/shared/patterns"
	"zpmeow/internal/shared/utils"
)

// ImageStrategy handles image media processing
type ImageStrategy struct{}

func (s *ImageStrategy) ValidateMedia(data []byte) error {
	return utils.ValidateMediaSize(data, "image")
}

func (s *ImageStrategy) ProcessMedia(ctx context.Context, data []byte, filename string) ([]byte, string, error) {
	// For images, we can return the data as-is with normalized MIME type
	mimeType := utils.NormalizeMimeType("image/jpeg", "image")
	return data, mimeType, nil
}

func (s *ImageStrategy) SendMessage(ctx context.Context, service patterns.MediaSender, sessionID, phone string, data []byte, caption, filename, mimeType string) (*whatsmeow.SendResponse, error) {
	return service.SendImageMessage(ctx, sessionID, phone, data, caption, mimeType)
}

// AudioStrategy handles audio media processing
type AudioStrategy struct{}

func (s *AudioStrategy) ValidateMedia(data []byte) error {
	return utils.ValidateMediaSize(data, "audio")
}

func (s *AudioStrategy) ProcessMedia(ctx context.Context, data []byte, filename string) ([]byte, string, error) {
	// For audio, normalize to WhatsApp's preferred format
	mimeType := utils.NormalizeMimeType("audio/ogg", "audio")
	return data, mimeType, nil
}

func (s *AudioStrategy) SendMessage(ctx context.Context, service patterns.MediaSender, sessionID, phone string, data []byte, caption, filename, mimeType string) (*whatsmeow.SendResponse, error) {
	return service.SendAudioMessage(ctx, sessionID, phone, data, mimeType)
}

// VideoStrategy handles video media processing
type VideoStrategy struct{}

func (s *VideoStrategy) ValidateMedia(data []byte) error {
	return utils.ValidateMediaSize(data, "video")
}

func (s *VideoStrategy) ProcessMedia(ctx context.Context, data []byte, filename string) ([]byte, string, error) {
	// For video, normalize to MP4
	mimeType := utils.NormalizeMimeType("video/mp4", "video")
	return data, mimeType, nil
}

func (s *VideoStrategy) SendMessage(ctx context.Context, service patterns.MediaSender, sessionID, phone string, data []byte, caption, filename, mimeType string) (*whatsmeow.SendResponse, error) {
	return service.SendVideoMessage(ctx, sessionID, phone, data, caption, mimeType)
}

// DocumentStrategy handles document media processing
type DocumentStrategy struct{}

func (s *DocumentStrategy) ValidateMedia(data []byte) error {
	return utils.ValidateMediaSize(data, "document")
}

func (s *DocumentStrategy) ProcessMedia(ctx context.Context, data []byte, filename string) ([]byte, string, error) {
	// For documents, preserve original MIME type
	mimeType := utils.NormalizeMimeType("application/octet-stream", "document")
	return data, mimeType, nil
}

func (s *DocumentStrategy) SendMessage(ctx context.Context, service patterns.MediaSender, sessionID, phone string, data []byte, caption, filename, mimeType string) (*whatsmeow.SendResponse, error) {
	return service.SendDocumentMessage(ctx, sessionID, phone, data, filename, caption, mimeType)
}

// StickerStrategy handles sticker media processing
type StickerStrategy struct{}

func (s *StickerStrategy) ValidateMedia(data []byte) error {
	return utils.ValidateMediaSize(data, "sticker")
}

func (s *StickerStrategy) ProcessMedia(ctx context.Context, data []byte, filename string) ([]byte, string, error) {
	// For stickers, normalize to WebP
	mimeType := utils.NormalizeMimeType("image/webp", "sticker")
	return data, mimeType, nil
}

func (s *StickerStrategy) SendMessage(ctx context.Context, service patterns.MediaSender, sessionID, phone string, data []byte, caption, filename, mimeType string) (*whatsmeow.SendResponse, error) {
	// Stickers don't have captions and use image message method
	return service.SendImageMessage(ctx, sessionID, phone, data, "", mimeType)
}

// MediaStrategyFactoryImpl implements the MediaStrategyFactory
type MediaStrategyFactoryImpl struct{}

// CreateStrategy creates a strategy for the given media type
func (f *MediaStrategyFactoryImpl) CreateStrategy(mediaType string) patterns.MediaStrategy {
	switch mediaType {
	case "image":
		return &ImageStrategy{}
	case "audio":
		return &AudioStrategy{}
	case "video":
		return &VideoStrategy{}
	case "document":
		return &DocumentStrategy{}
	case "sticker":
		return &StickerStrategy{}
	default:
		return nil
	}
}

// GetSupportedMediaTypes returns the list of supported media types
func (f *MediaStrategyFactoryImpl) GetSupportedMediaTypes() []string {
	return []string{"image", "audio", "video", "document", "sticker"}
}

// IsMediaTypeSupported checks if a media type is supported
func (f *MediaStrategyFactoryImpl) IsMediaTypeSupported(mediaType string) bool {
	for _, supported := range f.GetSupportedMediaTypes() {
		if mediaType == supported {
			return true
		}
	}
	return false
}

// ValidateMediaType validates that a media type is supported
func (f *MediaStrategyFactoryImpl) ValidateMediaType(mediaType string) error {
	if !f.IsMediaTypeSupported(mediaType) {
		return fmt.Errorf("unsupported media type: %s. Supported types: %v", mediaType, f.GetSupportedMediaTypes())
	}
	return nil
}

// Global instance
var MediaStrategyFactory = &MediaStrategyFactoryImpl{}
