package application

import (
	"context"
	"fmt"

	"zpmeow/internal/domain/message"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/shared/validation"
	"zpmeow/internal/interfaces/dto"
)

// MessageService implements message use cases following Clean Architecture
type MessageService struct {
	messageService message.Service
	sessionRepo    session.Repository
	validator      *validation.Validator
}

// NewMessageService creates a new MessageService instance
func NewMessageService(
	messageService message.Service,
	sessionRepo session.Repository,
	validator *validation.Validator,
) *MessageService {
	return &MessageService{
		messageService: messageService,
		sessionRepo:    sessionRepo,
		validator:      validator,
	}
}

// SendText sends a text message
func (s *MessageService) SendText(ctx context.Context, sessionID string, req *dto.SendTextRequest) (*dto.MessageResponse, error) {
	// 1. Validate input DTO
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 3. Build message using domain service
	chatJID := s.messageService.FormatJID(req.Phone)
	msg, err := s.messageService.BuildTextMessage(sessionID, chatJID, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to build message: %w", err)
	}

	// 4. Validate message using domain rules
	if err := s.messageService.ValidateMessage(msg); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// 5. Create response (actual sending would be handled by infrastructure layer)
	response := dto.NewTextResponse(true, 200, req.Phone, "msg_"+sessionID+"_"+req.Phone, req.Body, true)
	
	return response, nil
}

// SendMedia sends a media message
func (s *MessageService) SendMedia(ctx context.Context, sessionID string, req *dto.SendMediaRequest) (*dto.MessageResponse, error) {
	// 1. Validate input DTO
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 3. Convert media type to domain type
	var messageType message.MessageType
	switch req.MediaType {
	case "image":
		messageType = message.MessageTypeImage
	case "audio":
		messageType = message.MessageTypeAudio
	case "video":
		messageType = message.MessageTypeVideo
	case "document":
		messageType = message.MessageTypeDocument
	default:
		return nil, fmt.Errorf("unsupported media type: %s", req.MediaType)
	}

	// 4. Build message using domain service
	chatJID := s.messageService.FormatJID(req.Phone)
	msg, err := s.messageService.BuildMediaMessage(sessionID, chatJID, messageType, req.MediaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build media message: %w", err)
	}

	// 5. Validate message using domain rules
	if err := s.messageService.ValidateMessage(msg); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// 6. Create response based on media type
	var response *dto.MessageResponse
	messageID := "msg_" + sessionID + "_" + req.Phone

	switch messageType {
	case message.MessageTypeImage:
		payload := dto.MessagePayload{
			Image: &dto.ImageMessagePayload{
				URL:     req.MediaURL,
				Caption: req.Caption,
			},
		}
		response = dto.NewMessageResponse(true, 200, chatJID, messageID, true, payload)
	case message.MessageTypeAudio:
		payload := dto.MessagePayload{
			Audio: &dto.AudioMessagePayload{
				URL: req.MediaURL,
				PTT: false,
			},
		}
		response = dto.NewMessageResponse(true, 200, chatJID, messageID, true, payload)
	case message.MessageTypeVideo:
		payload := dto.MessagePayload{
			Video: &dto.VideoMessagePayload{
				URL:         req.MediaURL,
				Caption:     req.Caption,
				GifPlayback: false,
			},
		}
		response = dto.NewMessageResponse(true, 200, chatJID, messageID, true, payload)
	case message.MessageTypeDocument:
		payload := dto.MessagePayload{
			Document: &dto.DocumentMessagePayload{
				URL:      req.MediaURL,
				FileName: "document",
				Mimetype: "application/octet-stream",
			},
		}
		response = dto.NewMessageResponse(true, 200, chatJID, messageID, true, payload)
	}

	return response, nil
}

// SendLocation sends a location message
func (s *MessageService) SendLocation(ctx context.Context, sessionID string, req *dto.SendLocationRequest) (*dto.MessageResponse, error) {
	// 1. Validate input DTO
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 3. Validate coordinates
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("invalid latitude: must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("invalid longitude: must be between -180 and 180")
	}

	// 4. Build message using domain service
	chatJID := s.messageService.FormatJID(req.Phone)
	msg, err := s.messageService.BuildTextMessage(sessionID, chatJID, fmt.Sprintf("Location: %s", req.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to build location message: %w", err)
	}

	// 5. Validate message using domain rules
	if err := s.messageService.ValidateMessage(msg); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// 6. Create location response
	messageID := "msg_" + sessionID + "_" + req.Phone
	response := dto.NewLocationResponse(true, 200, req.Phone, messageID, req.Latitude, req.Longitude, req.Name, "", true)

	return response, nil
}

// SendContact sends a contact message
func (s *MessageService) SendContact(ctx context.Context, sessionID string, req *dto.SendContactRequest) (*dto.MessageResponse, error) {
	// 1. Validate input DTO
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 3. Build message using domain service
	chatJID := s.messageService.FormatJID(req.Phone)
	msg, err := s.messageService.BuildTextMessage(sessionID, chatJID, fmt.Sprintf("Contact: %s", req.ContactName))
	if err != nil {
		return nil, fmt.Errorf("failed to build contact message: %w", err)
	}

	// 4. Validate message using domain rules
	if err := s.messageService.ValidateMessage(msg); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// 5. Create vCard format
	vcard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL:%s\nEND:VCARD", req.ContactName, req.ContactPhone)

	// 6. Create contact response
	messageID := "msg_" + sessionID + "_" + req.Phone
	response := dto.NewContactResponse(true, 200, req.Phone, messageID, req.ContactName, vcard, true)

	return response, nil
}

// GetMessageStatus gets the status of a message
func (s *MessageService) GetMessageStatus(ctx context.Context, sessionID string, req *dto.MessageStatusRequest) (*dto.MessageStatusData, error) {
	// 1. Validate input DTO
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Verify session exists
	_, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// 3. Return mock status (actual implementation would query infrastructure)
	return &dto.MessageStatusData{
		MessageID: req.MessageID,
		Status:    "delivered",
	}, nil
}
