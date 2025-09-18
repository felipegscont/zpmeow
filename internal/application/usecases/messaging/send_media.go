package messaging

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// MediaType represents the type of media being sent
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeDocument MediaType = "document"
	MediaTypeSticker  MediaType = "sticker"
)

// IsValid checks if the media type is valid
func (mt MediaType) IsValid() bool {
	switch mt {
	case MediaTypeImage, MediaTypeVideo, MediaTypeAudio, MediaTypeDocument, MediaTypeSticker:
		return true
	default:
		return false
	}
}

// SendMediaMessageCommand represents the command to send a media message
type SendMediaMessageCommand struct {
	SessionID string
	ChatJID   string
	MediaType MediaType
	MediaData []byte
	MimeType  string
	Caption   string
	Filename  string
}

// Validate validates the send media message command
func (c SendMediaMessageCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if !c.MediaType.IsValid() {
		return common.NewValidationError("mediaType", c.MediaType, "invalid media type")
	}

	if len(c.MediaData) == 0 {
		return common.NewValidationError("mediaData", "", "media data is required")
	}

	if len(c.MediaData) > 100*1024*1024 { // 100MB limit
		return common.NewValidationError("mediaData", "", "media data exceeds 100MB limit")
	}

	if strings.TrimSpace(c.MimeType) == "" {
		return common.NewValidationError("mimeType", c.MimeType, "mime type is required")
	}

	if len(c.Caption) > 1024 {
		return common.NewValidationError("caption", c.Caption, "caption must not exceed 1024 characters")
	}

	return nil
}

// SendMediaMessageResult represents the result of sending a media message
type SendMediaMessageResult struct {
	SessionID string
	ChatJID   string
	MediaType string
	MessageID string
	Sent      bool
}

// SendMediaMessageUseCase handles sending media messages via WhatsApp
type SendMediaMessageUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewSendMediaMessageUseCase creates a new SendMediaMessageUseCase
func NewSendMediaMessageUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *SendMediaMessageUseCase {
	return &SendMediaMessageUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the send media message use case
func (uc *SendMediaMessageUseCase) Handle(ctx context.Context, cmd SendMediaMessageCommand) (*SendMediaMessageResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid send media message command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, cmd.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is connected (business rule)
	if !sessionEntity.IsConnected() {
		return nil, common.NewBusinessRuleError(
			"session_not_connected",
			fmt.Sprintf("session must be connected to send messages, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Check if session is authenticated
	if !sessionEntity.IsAuthenticated() {
		return nil, common.NewBusinessRuleError(
			"session_not_authenticated",
			"session must be authenticated to send messages",
		)
	}

	// 5. Prepare media message
	mediaMessage := ports.MediaMessage{
		Type:     string(cmd.MediaType),
		Data:     cmd.MediaData,
		MimeType: cmd.MimeType,
		Caption:  cmd.Caption,
		Filename: cmd.Filename,
	}

	// 6. Send media message via WhatsApp service
	if err := uc.whatsappService.SendMediaMessage(ctx, cmd.SessionID, cmd.ChatJID, mediaMessage); err != nil {
		uc.logger.Error(ctx, "Failed to send media message",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"mediaType", cmd.MediaType,
			"error", err)
		return nil, fmt.Errorf("failed to send media message: %w", err)
	}

	uc.logger.Info(ctx, "Media message sent successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"mediaType", cmd.MediaType,
		"dataSize", len(cmd.MediaData))

	// 7. Return result
	return &SendMediaMessageResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		MediaType: string(cmd.MediaType),
		MessageID: "", // Would be provided by WhatsApp service
		Sent:      true,
	}, nil
}
