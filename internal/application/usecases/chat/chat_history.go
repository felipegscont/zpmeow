package chat

import (
	"context"
	"fmt"
	"strings"

	"zpmeow/internal/application/common"
	"zpmeow/internal/application/ports"
)

// GetChatHistoryQuery represents the query to get chat history
type GetChatHistoryQuery struct {
	SessionID string
	ChatJID   string
	Limit     int
	Offset    int
}

// Validate validates the get chat history query
func (q GetChatHistoryQuery) Validate() error {
	if strings.TrimSpace(q.SessionID) == "" {
		return common.NewValidationError("sessionID", q.SessionID, "session ID is required")
	}

	if strings.TrimSpace(q.ChatJID) == "" {
		return common.NewValidationError("chatJID", q.ChatJID, "chat JID is required")
	}

	if q.Limit < 0 {
		return common.NewValidationError("limit", q.Limit, "limit cannot be negative")
	}

	if q.Limit > 1000 {
		return common.NewValidationError("limit", q.Limit, "limit cannot exceed 1000")
	}

	if q.Offset < 0 {
		return common.NewValidationError("offset", q.Offset, "offset cannot be negative")
	}

	return nil
}

// MessageView represents a message view model
type MessageView struct {
	ID        string
	ChatJID   string
	FromJID   string
	Content   string
	Type      string
	Timestamp string
	IsFromMe  bool
	IsRead    bool
	MediaURL  string
	Caption   string
}

// GetChatHistoryResult represents the result of getting chat history
type GetChatHistoryResult struct {
	SessionID string
	ChatJID   string
	Messages  []MessageView
	Total     int
	Limit     int
	Offset    int
}

// GetChatHistoryUseCase handles getting chat message history
type GetChatHistoryUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewGetChatHistoryUseCase creates a new GetChatHistoryUseCase
func NewGetChatHistoryUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *GetChatHistoryUseCase {
	return &GetChatHistoryUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the get chat history use case
func (uc *GetChatHistoryUseCase) Handle(ctx context.Context, query GetChatHistoryQuery) (*GetChatHistoryResult, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid get chat history query", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, query.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", query.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is connected (business rule)
	if !sessionEntity.IsConnected() {
		return nil, common.NewBusinessRuleError(
			"session_not_connected",
			fmt.Sprintf("session must be connected to get chat history, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Get chat history via WhatsApp service
	messages, err := uc.whatsappService.GetChatHistory(ctx, query.SessionID, query.ChatJID, query.Limit, query.Offset)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get chat history",
			"sessionID", query.SessionID,
			"chatJID", query.ChatJID,
			"error", err)
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	// 5. Convert to view models
	messageViews := make([]MessageView, len(messages))
	for i, message := range messages {
		messageViews[i] = MessageView{
			ID:        message.ID,
			ChatJID:   message.ChatJID,
			FromJID:   message.FromJID,
			Content:   message.Content,
			Type:      message.Type,
			Timestamp: message.Timestamp,
			IsFromMe:  message.IsFromMe,
			IsRead:    message.IsRead,
			MediaURL:  message.MediaURL,
			Caption:   message.Caption,
		}
	}

	uc.logger.Debug(ctx, "Chat history retrieved successfully",
		"sessionID", query.SessionID,
		"chatJID", query.ChatJID,
		"messageCount", len(messageViews))

	// 6. Return result
	return &GetChatHistoryResult{
		SessionID: query.SessionID,
		ChatJID:   query.ChatJID,
		Messages:  messageViews,
		Total:     len(messageViews), // In a real implementation, this would be the total count
		Limit:     query.Limit,
		Offset:    query.Offset,
	}, nil
}

// SetPresenceCommand represents the command to set presence in a chat
type SetPresenceCommand struct {
	SessionID string
	ChatJID   string
	State     string // available, unavailable, composing, recording, paused
	Media     string // optional media type for recording state
}

// Validate validates the set presence command
func (c SetPresenceCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if strings.TrimSpace(c.State) == "" {
		return common.NewValidationError("state", c.State, "presence state is required")
	}

	validStates := map[string]bool{
		"available":   true,
		"unavailable": true,
		"composing":   true,
		"recording":   true,
		"paused":      true,
	}

	if !validStates[c.State] {
		return common.NewValidationError("state", c.State, "invalid presence state. Valid states: available, unavailable, composing, recording, paused")
	}

	return nil
}

// SetPresenceResult represents the result of setting presence
type SetPresenceResult struct {
	SessionID string
	ChatJID   string
	State     string
	Success   bool
}

// SetPresenceUseCase handles setting presence in chats
type SetPresenceUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewSetPresenceUseCase creates a new SetPresenceUseCase
func NewSetPresenceUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *SetPresenceUseCase {
	return &SetPresenceUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the set presence use case
func (uc *SetPresenceUseCase) Handle(ctx context.Context, cmd SetPresenceCommand) (*SetPresenceResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid set presence command", "error", err)
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
			fmt.Sprintf("session must be connected to set presence, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Set presence via WhatsApp service
	if err := uc.whatsappService.SetPresence(ctx, cmd.SessionID, cmd.ChatJID, cmd.State, cmd.Media); err != nil {
		uc.logger.Error(ctx, "Failed to set presence",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"state", cmd.State,
			"error", err)
		return nil, fmt.Errorf("failed to set presence: %w", err)
	}

	uc.logger.Info(ctx, "Presence set successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"state", cmd.State)

	// 5. Return result
	return &SetPresenceResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		State:     cmd.State,
		Success:   true,
	}, nil
}
