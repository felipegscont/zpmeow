package application

import (
	"context"
	"fmt"
	"time"

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

// SendText sends a text message using DTO
func (s *MessagingService) SendText(ctx context.Context, sessionID string, req *dto.SendTextRequest) (*dto.MessageResponse, error) {
	// 1. Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	// 2. Validate session
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	// 3. Send message via wameow service
	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendTextMessage(ctx, sessionID, chatJID, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to send text message: %w", err)
	}

	// 4. Build response (whatsmeow will handle events)
	response := s.buildMessageResponse(result)
	return response, nil
}

// SendImage sends an image message using DTO
func (s *MessagingService) SendImage(ctx context.Context, sessionID string, req *dto.SendImageRequest, imageData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendImageMessage(ctx, sessionID, chatJID, imageData, req.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send image message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendVideo sends a video message using DTO
func (s *MessagingService) SendVideo(ctx context.Context, sessionID string, req *dto.SendVideoRequest, videoData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendVideoMessage(ctx, sessionID, chatJID, videoData, req.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send video message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendAudio sends an audio message using DTO
func (s *MessagingService) SendAudio(ctx context.Context, sessionID string, req *dto.SendAudioRequest, audioData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendAudioMessage(ctx, sessionID, chatJID, audioData)
	if err != nil {
		return nil, fmt.Errorf("failed to send audio message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendDocument sends a document message using DTO
func (s *MessagingService) SendDocument(ctx context.Context, sessionID string, req *dto.SendDocumentRequest, documentData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendDocumentMessage(ctx, sessionID, chatJID, documentData, req.FileName, req.MimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to send document message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendSticker sends a sticker message using DTO
func (s *MessagingService) SendSticker(ctx context.Context, sessionID string, req *dto.SendStickerRequest, stickerData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendStickerMessage(ctx, sessionID, chatJID, stickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to send sticker message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendContact sends a contact message using DTO
func (s *MessagingService) SendContact(ctx context.Context, sessionID string, req *dto.SendContactRequest) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	// Build vCard from contact info
	vCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL:%s\nEND:VCARD", req.ContactName, req.ContactPhone)

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendContactMessage(ctx, sessionID, chatJID, vCard)
	if err != nil {
		return nil, fmt.Errorf("failed to send contact message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendLocation sends a location message using DTO
func (s *MessagingService) SendLocation(ctx context.Context, sessionID string, req *dto.SendLocationRequest) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.wameowService.SendLocationMessage(ctx, sessionID, chatJID, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to send location message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// MarkAsRead marks messages as read using DTO
func (s *MessagingService) MarkAsRead(ctx context.Context, sessionID string, req *dto.MarkAsReadRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	// For now, mark the first message ID as read (the DTO supports multiple but wameow service expects single)
	if len(req.MessageIDs) > 0 {
		return s.wameowService.MarkAsRead(ctx, sessionID, chatJID, req.MessageIDs[0])
	}
	return fmt.Errorf("no message IDs provided")
}

// ReactToMessage reacts to a message using DTO
func (s *MessagingService) ReactToMessage(ctx context.Context, sessionID string, req *dto.ReactToMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.wameowService.ReactToMessage(ctx, sessionID, chatJID, req.MessageID, req.Emoji)
}

// EditMessage edits a message using DTO
func (s *MessagingService) EditMessage(ctx context.Context, sessionID string, req *dto.EditMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.wameowService.EditMessage(ctx, sessionID, chatJID, req.MessageID, req.NewText)
}

// DeleteMessage deletes a message using DTO
func (s *MessagingService) DeleteMessage(ctx context.Context, sessionID string, req *dto.DeleteMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.wameowService.DeleteMessage(ctx, sessionID, chatJID, req.MessageID)
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
