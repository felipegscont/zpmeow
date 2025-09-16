package application

import (
	"context"
	"fmt"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/interfaces/dto"
	"zpmeow/internal/shared/validation"
)

// ChatService implements chat use cases following Clean Architecture
type ChatService struct {
	sessionRepo session.Repository
	validator   *validation.Validator
}

// NewChatService creates a new ChatService instance
func NewChatService(
	sessionRepo session.Repository,
	validator *validation.Validator,
) *ChatService {
	return &ChatService{
		sessionRepo: sessionRepo,
		validator:   validator,
	}
}

// GetChatHistory retrieves chat history for a specific contact
func (s *ChatService) GetChatHistory(ctx context.Context, sessionID, phone string, limit int) ([]*dto.ChatHistoryData, error) {
	// 1. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 2. Validate phone number
	if err := s.validator.ValidatePhoneNumber(phone); err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	// 3. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 4. Validate limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}

	// 5. Return mock chat history (actual implementation would query infrastructure)
	// In a real implementation, this would call the infrastructure layer to get chat history
	history := make([]*dto.ChatHistoryData, 0)

	// Mock data for demonstration
	history = append(history, &dto.ChatHistoryData{
		MessageID:   "msg_001",
		Phone:       phone,
		FromPhone:   phone,
		MessageType: "text",
		Content:     "Hello!",
		IsFromMe:    false,
	})

	return history, nil
}

// ListChats retrieves all chats for a session
func (s *ChatService) ListChats(ctx context.Context, sessionID string) ([]map[string]interface{}, error) {
	// 1. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 2. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 3. Return mock chat list (actual implementation would query infrastructure)
	chats := make([]map[string]interface{}, 0)

	// Mock data for demonstration
	chats = append(chats, map[string]interface{}{
		"jid":          "5511999999999@s.whatsapp.net",
		"name":         "Contact Name",
		"last_message": "Hello!",
		"unread_count": 0,
		"is_group":     false,
	})

	return chats, nil
}

// MarkAsRead marks messages as read in a chat
func (s *ChatService) MarkAsRead(ctx context.Context, sessionID string, req *dto.MarkAsReadRequest) (*dto.ChatResponse, error) {
	// 1. Validate input DTO
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 3. Validate phone number
	if err := s.validator.ValidatePhoneNumber(req.Phone); err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	// 4. Validate message IDs
	if err := s.validator.ValidateMessageIDs(req.MessageIDs); err != nil {
		return nil, fmt.Errorf("invalid message IDs: %w", err)
	}

	// 5. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 6. Create success response (actual implementation would call infrastructure)
	response := dto.NewChatSuccessResponse(req.Phone, "", "mark_read")

	return response, nil
}

// SetPresence sets user presence in a chat
func (s *ChatService) SetPresence(ctx context.Context, sessionID string, req *dto.SetPresenceRequest) (*dto.ChatResponse, error) {
	// 1. Validate input DTO
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 3. Validate phone number if provided
	if req.Phone != "" {
		if err := s.validator.ValidatePhoneNumber(req.Phone); err != nil {
			return nil, fmt.Errorf("invalid phone number: %w", err)
		}
	}

	// 4. Validate presence state
	validStates := []string{"available", "unavailable", "composing", "recording", "paused"}
	isValidState := false
	for _, state := range validStates {
		if req.State == state {
			isValidState = true
			break
		}
	}
	if !isValidState {
		return nil, fmt.Errorf("invalid presence state: %s", req.State)
	}

	// 5. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 6. Create success response (actual implementation would call infrastructure)
	response := dto.NewChatSuccessResponse(req.Phone, "", "set_presence")

	return response, nil
}

// ArchiveChat archives a chat
func (s *ChatService) ArchiveChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	// 1. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 2. Validate chat JID
	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	// 3. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 4. Create success response (actual implementation would call infrastructure)
	response := dto.NewChatSuccessResponse(chatJID, "", "archive_chat")

	return response, nil
}

// UnarchiveChat unarchives a chat
func (s *ChatService) UnarchiveChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	// 1. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 2. Validate chat JID
	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	// 3. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 4. Create success response (actual implementation would call infrastructure)
	response := dto.NewChatSuccessResponse(chatJID, "", "unarchive_chat")

	return response, nil
}

// DeleteChat deletes a chat
func (s *ChatService) DeleteChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	// 1. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 2. Validate chat JID
	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	// 3. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 4. Create success response (actual implementation would call infrastructure)
	response := dto.NewChatSuccessResponse(chatJID, "", "delete_chat")

	return response, nil
}

// MuteChat mutes a chat for a specified duration
func (s *ChatService) MuteChat(ctx context.Context, sessionID, chatJID string, duration int) (*dto.ChatResponse, error) {
	// 1. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 2. Validate chat JID
	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	// 3. Validate duration
	if duration < 0 {
		return nil, fmt.Errorf("duration cannot be negative")
	}

	// 4. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 5. Create success response (actual implementation would call infrastructure)
	response := dto.NewChatSuccessResponse(chatJID, "", "mute_chat")

	return response, nil
}

// UnmuteChat unmutes a chat
func (s *ChatService) UnmuteChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	// 1. Validate session ID
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 2. Validate chat JID
	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	// 3. Verify session exists and is active
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	// 4. Create success response (actual implementation would call infrastructure)
	response := dto.NewChatSuccessResponse(chatJID, "", "unmute_chat")

	return response, nil
}
