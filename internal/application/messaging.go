package application

import (
	"context"
	"fmt"
	"time"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

// Interface movida para interfaces.go para centralização

type MessageApp struct {
	messageSender ExtendedMessageSender
	sessionRepo   session.Repository
	validator     *validation.Validator
}

func NewMessageApp(
	messageSender ExtendedMessageSender,
	sessionRepo session.Repository,
	validator *validation.Validator,
) *MessageApp {
	return &MessageApp{
		messageSender: messageSender,
		sessionRepo:   sessionRepo,
		validator:     validator,
	}
}

func (s *MessageApp) SendText(ctx context.Context, sessionID string, req *dto.SendTextRequest) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendTextMessage(ctx, sessionID, chatJID, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to send text message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) SendImage(ctx context.Context, sessionID string, req *dto.SendImageRequest, imageData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendImageMessage(ctx, sessionID, chatJID, imageData, req.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send image message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) SendVideo(ctx context.Context, sessionID string, req *dto.SendVideoRequest, videoData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendVideoMessage(ctx, sessionID, chatJID, videoData, req.Caption)
	if err != nil {
		return nil, fmt.Errorf("failed to send video message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) SendAudio(ctx context.Context, sessionID string, req *dto.SendAudioRequest, audioData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendAudioMessage(ctx, sessionID, chatJID, audioData)
	if err != nil {
		return nil, fmt.Errorf("failed to send audio message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) SendDocument(ctx context.Context, sessionID string, req *dto.SendDocumentRequest, documentData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendDocumentMessage(ctx, sessionID, chatJID, documentData, req.FileName, req.MimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to send document message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) SendSticker(ctx context.Context, sessionID string, req *dto.SendStickerRequest, stickerData []byte) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendStickerMessage(ctx, sessionID, chatJID, stickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to send sticker message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) SendContact(ctx context.Context, sessionID string, req *dto.SendContactRequest) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	vCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL:%s\nEND:VCARD", req.ContactName, req.ContactPhone)

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendContactMessage(ctx, sessionID, chatJID, vCard)
	if err != nil {
		return nil, fmt.Errorf("failed to send contact message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) SendLocation(ctx context.Context, sessionID string, req *dto.SendLocationRequest) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	chatJID := s.resolveChatJID(req.Phone)
	result, err := s.messageSender.SendLocationMessage(ctx, sessionID, chatJID, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to send location message: %w", err)
	}

	response := s.buildMessageResponse(result)
	return response, nil
}

func (s *MessageApp) MarkAsRead(ctx context.Context, sessionID string, req *dto.MarkAsReadRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	if len(req.MessageIDs) > 0 {
		return s.messageSender.MarkAsRead(ctx, sessionID, chatJID, req.MessageIDs[0])
	}
	return fmt.Errorf("no message IDs provided")
}

func (s *MessageApp) ReactToMessage(ctx context.Context, sessionID string, req *dto.ReactToMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.messageSender.ReactToMessage(ctx, sessionID, chatJID, req.MessageID, req.Emoji)
}

func (s *MessageApp) EditMessage(ctx context.Context, sessionID string, req *dto.EditMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.messageSender.EditMessage(ctx, sessionID, chatJID, req.MessageID, req.NewText)
}

func (s *MessageApp) DeleteMessage(ctx context.Context, sessionID string, req *dto.DeleteMessageRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	chatJID := s.resolveChatJID(req.Phone)
	return s.messageSender.DeleteMessage(ctx, sessionID, chatJID, req.MessageID)
}

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
	return chatJID
}

func (s *MessageApp) buildMessageResponse(result *MessageResult) *dto.MessageResponse {
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

func extractMessageID(result *MessageResult) string {
	if result == nil {
		return "unknown"
	}
	return result.ID
}
