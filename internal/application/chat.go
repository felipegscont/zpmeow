package application

import (
	"context"
	"fmt"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

type ChatApp struct {
	sessionRepo session.Repository
	validator   *validation.Validator
}

func NewChatApp(
	sessionRepo session.Repository,
	validator *validation.Validator,
) *ChatApp {
	return &ChatApp{
		sessionRepo: sessionRepo,
		validator:   validator,
	}
}

func (s *ChatApp) GetChatHistory(ctx context.Context, sessionID, phone string, limit int) ([]*dto.ChatHistoryData, error) {
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if err := s.validator.ValidatePhoneNumber(phone); err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}

	history := make([]*dto.ChatHistoryData, 0)

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

func (s *ChatApp) ListChats(ctx context.Context, sessionID string) ([]ChatInfo, error) {
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	chats := make([]ChatInfo, 0)

	chats = append(chats, ChatInfo{
		JID:           "5511999999999@s.meow.net",
		Name:          "Contact Name",
		LastMessage:   "Hello!",
		UnreadCount:   0,
		IsGroup:       false,
		LastTimestamp: 0,
		IsMuted:       false,
		IsPinned:      false,
		IsArchived:    false,
	})

	return chats, nil
}

func (s *ChatApp) MarkAsRead(ctx context.Context, sessionID string, req *dto.MarkAsReadRequest) (*dto.ChatResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if err := s.validator.ValidatePhoneNumber(req.Phone); err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	if err := s.validator.ValidateMessageIDs(req.MessageIDs); err != nil {
		return nil, fmt.Errorf("invalid message IDs: %w", err)
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	response := dto.NewChatSuccessResponse(req.Phone, "", "mark_read")

	return response, nil
}

func (s *ChatApp) SetPresence(ctx context.Context, sessionID string, req *dto.SetPresenceRequest) (*dto.ChatResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if req.Phone != "" {
		if err := s.validator.ValidatePhoneNumber(req.Phone); err != nil {
			return nil, fmt.Errorf("invalid phone number: %w", err)
		}
	}

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

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	response := dto.NewChatSuccessResponse(req.Phone, "", "set_presence")

	return response, nil
}

func (s *ChatApp) ArchiveChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	response := dto.NewChatSuccessResponse(chatJID, "", "archive_chat")

	return response, nil
}

func (s *ChatApp) UnarchiveChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	response := dto.NewChatSuccessResponse(chatJID, "", "unarchive_chat")

	return response, nil
}

func (s *ChatApp) DeleteChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	response := dto.NewChatSuccessResponse(chatJID, "", "delete_chat")

	return response, nil
}

func (s *ChatApp) MuteChat(ctx context.Context, sessionID, chatJID string, duration int) (*dto.ChatResponse, error) {
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	if duration < 0 {
		return nil, fmt.Errorf("duration cannot be negative")
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	response := dto.NewChatSuccessResponse(chatJID, "", "mute_chat")

	return response, nil
}

func (s *ChatApp) UnmuteChat(ctx context.Context, sessionID, chatJID string) (*dto.ChatResponse, error) {
	if err := s.validator.ValidateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if chatJID == "" {
		return nil, fmt.Errorf("chat JID is required")
	}

	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	response := dto.NewChatSuccessResponse(chatJID, "", "unmute_chat")

	return response, nil
}
