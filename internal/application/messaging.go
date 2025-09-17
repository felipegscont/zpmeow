package application

import (
	"context"
	"fmt"
	"time"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

// MessageApp implements messaging use cases following Clean Architecture
type MessageApp struct {
	meowService interface {
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

// NewMessageApp creates a new MessageApp instance
func NewMessageApp(
	meowService interface{},
	sessionRepo session.Repository,
	validator *validation.Validator,
) *MessageApp {
	return &MessageApp{
		meowService: meowService.(interface {
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
func (s *MessageApp) SendText(ctx context.Context, sessionID string, req *dto.SendTextRequest) (*dto.MessageResponse, error) {
	// 1. Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	// 2. Validate session
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	// 3. Send message via meow service
	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendTextMessage(ctx, sessionID, chatJID, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to send text message: %w", err)
	}

	// 4. Build response (whatsmeow will handle events)
	response := s.buildMessageResponse(result)
	return response, nil
}

// SendImage sends an image message using DTO
func (s *MessageApp) SendImage(ctx context.Context, sessionID string, req *dto.SendImageRequest, imageData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendImageMessage(ctx, sessionID, chatJID, imageData, req.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send image message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendVideo sends a video message using DTO
func (s *MessageApp) SendVideo(ctx context.Context, sessionID string, req *dto.SendVideoRequest, videoData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendVideoMessage(ctx, sessionID, chatJID, videoData, req.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send video message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendAudio sends an audio message using DTO
func (s *MessageApp) SendAudio(ctx context.Context, sessionID string, req *dto.SendAudioRequest, audioData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendAudioMessage(ctx, sessionID, chatJID, audioData)
	if err != nil {
		return nil, fmt.Errorf("failed to send audio message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendDocument sends a document message using DTO
func (s *MessageApp) SendDocument(ctx context.Context, sessionID string, req *dto.SendDocumentRequest, documentData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendDocumentMessage(ctx, sessionID, chatJID, documentData, req.FileName, req.MimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to send document message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendSticker sends a sticker message using DTO
func (s *MessageApp) SendSticker(ctx context.Context, sessionID string, req *dto.SendStickerRequest, stickerData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendStickerMessage(ctx, sessionID, chatJID, stickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to send sticker message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendContact sends a contact message using DTO
func (s *MessageApp) SendContact(ctx context.Context, sessionID string, req *dto.SendContactRequest) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	// Build vCard from contact info
	vCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL:%s\nEND:VCARD", req.ContactName, req.ContactPhone)

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendContactMessage(ctx, sessionID, chatJID, vCard)
	if err != nil {
		return nil, fmt.Errorf("failed to send contact message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// SendLocation sends a location message using DTO
func (s *MessageApp) SendLocation(ctx context.Context, sessionID string, req *dto.SendLocationRequest) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.meowService.SendLocationMessage(ctx, sessionID, chatJID, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to send location message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

// MarkAsRead marks messages as read using DTO
func (s *MessageApp) MarkAsRead(ctx context.Context, sessionID string, req *dto.MarkAsReadRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	// For now, mark the first message ID as read (the DTO supports multiple but meow service expects single)
	if len(req.MessageIDs) > 0 {
		return s.meowService.MarkAsRead(ctx, sessionID, chatJID, req.MessageIDs[0])
	}
	return fmt.Errorf("no message IDs provided")
}

// ReactToMessage reacts to a message using DTO
func (s *MessageApp) ReactToMessage(ctx context.Context, sessionID string, req *dto.ReactToMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.meowService.ReactToMessage(ctx, sessionID, chatJID, req.MessageID, req.Emoji)
}

// EditMessage edits a message using DTO
func (s *MessageApp) EditMessage(ctx context.Context, sessionID string, req *dto.EditMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.meowService.EditMessage(ctx, sessionID, chatJID, req.MessageID, req.NewText)
}

// DeleteMessage deletes a message using DTO
func (s *MessageApp) DeleteMessage(ctx context.Context, sessionID string, req *dto.DeleteMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.meowService.DeleteMessage(ctx, sessionID, chatJID, req.MessageID)
}

// Helper methods

func (s *MessageApp) validateSession(ctx context.Context, sessionID string) error {
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return fmt.Errorf("session is not connected")
	}

	return nil
}

func (s *MessageApp) resolveChatJID(chatJID string) string {
	// Add logic to resolve chat JID if needed
	// For now, return as-is
	return chatJID
}

func (s *MessageApp) buildMessageResponse(result interface{}) *dto.MessageResponse {
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

// Helper function to extract message ID from meow service response
func extractMessageID(result interface{}) string {
	if result == nil {
		return "unknown"
	}
	return "msg_" + fmt.Sprintf("%v", result)
}
