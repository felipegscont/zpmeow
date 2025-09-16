package application

import (
	"context"
	"fmt"
	"time"

	"zpmeow/internal/application/commands"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/interfaces/dto"
	"zpmeow/internal/shared/validation"
)

// MessagingService implements messaging use cases following Clean Architecture
type MessagingService struct {
	wameowService interface {
		SendTextMessage(ctx context.Context, sessionID, chatJID, content string) (interface{}, error)
		SendImageMessage(ctx context.Context, sessionID, chatJID string, imageData []byte, caption string) (interface{}, error)
		SendVideoMessage(ctx context.Context, sessionID, chatJID string, videoData []byte, caption string) (interface{}, error)
		SendAudioMessage(ctx context.Context, sessionID, chatJID string, audioData []byte) (interface{}, error)
		SendDocumentMessage(ctx context.Context, sessionID, chatJID string, documentData []byte, filename, mimetype string) (interface{}, error)
		SendStickerMessage(ctx context.Context, sessionID, chatJID string, stickerData []byte) (interface{}, error)
		SendContactMessage(ctx context.Context, sessionID, chatJID, contactVCard string) (interface{}, error)
		SendLocationMessage(ctx context.Context, sessionID, chatJID string, latitude, longitude float64, name, address string) (interface{}, error)
		MarkAsRead(ctx context.Context, sessionID, chatJID, messageID string) error
		ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, reaction string) error
		EditMessage(ctx context.Context, sessionID, chatJID, messageID, newContent string) error
		DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string) error
	}
	sessionRepo session.Repository
	validator   *validation.Validator
}

// NewMessagingService creates a new MessagingService instance
func NewMessagingService(
	wameowService interface{},
	sessionRepo session.Repository,
	validator *validation.Validator,
) *MessagingService {
	return &MessagingService{
		wameowService: wameowService.(interface {
			SendTextMessage(ctx context.Context, sessionID, chatJID, content string) (interface{}, error)
			SendImageMessage(ctx context.Context, sessionID, chatJID string, imageData []byte, caption string) (interface{}, error)
			SendVideoMessage(ctx context.Context, sessionID, chatJID string, videoData []byte, caption string) (interface{}, error)
			SendAudioMessage(ctx context.Context, sessionID, chatJID string, audioData []byte) (interface{}, error)
			SendDocumentMessage(ctx context.Context, sessionID, chatJID string, documentData []byte, filename, mimetype string) (interface{}, error)
			SendStickerMessage(ctx context.Context, sessionID, chatJID string, stickerData []byte) (interface{}, error)
			SendContactMessage(ctx context.Context, sessionID, chatJID, contactVCard string) (interface{}, error)
			SendLocationMessage(ctx context.Context, sessionID, chatJID string, latitude, longitude float64, name, address string) (interface{}, error)
			MarkAsRead(ctx context.Context, sessionID, chatJID, messageID string) error
			ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, reaction string) error
			EditMessage(ctx context.Context, sessionID, chatJID, messageID, newContent string) error
			DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string) error
		}),
		sessionRepo: sessionRepo,
		validator:   validator,
	}
}

// SendText sends a text message using command pattern
func (s *MessagingService) SendText(ctx context.Context, cmd *commands.SendTextCommand) (*dto.MessageResponse, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	// 2. Validate session
	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	// 3. Send message via wameow service
	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendTextMessage(ctx, cmd.SessionID, chatJID, cmd.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to send text message: %w", err)
	}

	// 4. Build response (whatsmeow will handle events)
	response := s.buildMessageResponse(result)
	return response, nil
}

// SendImage sends an image message using command pattern
func (s *MessagingService) SendImage(ctx context.Context, cmd *commands.SendImageCommand) (*dto.MessageResponse, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendImageMessage(ctx, cmd.SessionID, chatJID, cmd.ImageData, cmd.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send image message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendVideo sends a video message using command pattern
func (s *MessagingService) SendVideo(ctx context.Context, cmd *commands.SendVideoCommand) (*dto.MessageResponse, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendVideoMessage(ctx, cmd.SessionID, chatJID, cmd.VideoData, cmd.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send video message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendAudio sends an audio message using command pattern
func (s *MessagingService) SendAudio(ctx context.Context, cmd *commands.SendAudioCommand) (*dto.MessageResponse, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendAudioMessage(ctx, cmd.SessionID, chatJID, cmd.AudioData)
	if err != nil {
		return nil, fmt.Errorf("failed to send audio message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendDocument sends a document message using command pattern
func (s *MessagingService) SendDocument(ctx context.Context, cmd *commands.SendDocumentCommand) (*dto.MessageResponse, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendDocumentMessage(ctx, cmd.SessionID, chatJID, cmd.DocumentData, cmd.Filename, cmd.Mimetype)
	if err != nil {
		return nil, fmt.Errorf("failed to send document message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendSticker sends a sticker message using command pattern
func (s *MessagingService) SendSticker(ctx context.Context, cmd *commands.SendStickerCommand) (*dto.MessageResponse, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendStickerMessage(ctx, cmd.SessionID, chatJID, cmd.StickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to send sticker message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendContact sends a contact message using command pattern
func (s *MessagingService) SendContact(ctx context.Context, cmd *commands.SendContactCommand) (*dto.MessageResponse, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendContactMessage(ctx, cmd.SessionID, chatJID, cmd.VCard)
	if err != nil {
		return nil, fmt.Errorf("failed to send contact message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendLocation sends a location message using command pattern
func (s *MessagingService) SendLocation(ctx context.Context, cmd *commands.SendLocationCommand) (*dto.MessageResponse, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	result, err := s.wameowService.SendLocationMessage(ctx, cmd.SessionID, chatJID, cmd.Latitude, cmd.Longitude, cmd.Name, cmd.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to send location message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// MarkAsRead marks a message as read using command pattern
func (s *MessagingService) MarkAsRead(ctx context.Context, cmd *commands.MarkAsReadCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	return s.wameowService.MarkAsRead(ctx, cmd.SessionID, chatJID, cmd.MessageID)
}

// ReactToMessage reacts to a message using command pattern
func (s *MessagingService) ReactToMessage(ctx context.Context, cmd *commands.ReactToMessageCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	return s.wameowService.ReactToMessage(ctx, cmd.SessionID, chatJID, cmd.MessageID, cmd.Reaction)
}

// EditMessage edits a message using command pattern
func (s *MessagingService) EditMessage(ctx context.Context, cmd *commands.EditMessageCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	return s.wameowService.EditMessage(ctx, cmd.SessionID, chatJID, cmd.MessageID, cmd.NewContent)
}

// DeleteMessage deletes a message using command pattern
func (s *MessagingService) DeleteMessage(ctx context.Context, cmd *commands.DeleteMessageCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(cmd.ChatJID)
	return s.wameowService.DeleteMessage(ctx, cmd.SessionID, chatJID, cmd.MessageID)
}

// Helper methods

func (s *MessagingService) validateSession(ctx context.Context, sessionID string) error {
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return fmt.Errorf("session is not connected")
	}

	return nil
}

func (s *MessagingService) resolveChatJID(chatJID string) string {
	// Add logic to resolve chat JID if needed
	// For now, return as-is
	return chatJID
}

func (s *MessagingService) buildMessageResponse(result interface{}) *dto.MessageResponse {
	messageID := extractMessageID(result)

	return &dto.MessageResponse{
		Success: true,
		Code:    200,
		Data: dto.MessageResponseData{
			Key: dto.MessageKey{
				ID: messageID,
			},
			Timestamp: time.Now().Unix(),
		},
	}
}

// Helper function to extract message ID from wameow service response
func extractMessageID(result interface{}) string {
	if result == nil {
		return "unknown"
	}
	return "msg_" + fmt.Sprintf("%v", result)
}
