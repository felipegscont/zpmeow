package chat

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// GetChatsQuery represents the query to get chats
type GetChatsQuery struct {
	SessionID string
	Limit     int
	Offset    int
}

// Validate validates the get chats query
func (q GetChatsQuery) Validate() error {
	if strings.TrimSpace(q.SessionID) == "" {
		return common.NewValidationError("sessionID", q.SessionID, "session ID is required")
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

// ChatView represents a chat view model
type ChatView struct {
	JID           string
	Name          string
	LastMessage   string
	LastMessageAt string
	UnreadCount   int
	IsGroup       bool
	IsMuted       bool
	IsArchived    bool
	IsBlocked     bool
}

// GetChatsResult represents the result of getting chats
type GetChatsResult struct {
	SessionID string
	Chats     []ChatView
	Total     int
	Limit     int
	Offset    int
}

// GetChatsUseCase handles getting chats for a session
type GetChatsUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewGetChatsUseCase creates a new GetChatsUseCase
func NewGetChatsUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *GetChatsUseCase {
	return &GetChatsUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the get chats use case
func (uc *GetChatsUseCase) Handle(ctx context.Context, query GetChatsQuery) (*GetChatsResult, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid get chats query", "error", err)
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
			fmt.Sprintf("session must be connected to get chats, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Get chats via WhatsApp service
	chats, err := uc.whatsappService.GetChats(ctx, query.SessionID, query.Limit, query.Offset)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get chats",
			"sessionID", query.SessionID,
			"error", err)
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	// 5. Convert to view models
	chatViews := make([]ChatView, len(chats))
	for i, chat := range chats {
		chatViews[i] = ChatView{
			JID:           chat.JID,
			Name:          chat.Name,
			LastMessage:   chat.LastMessage,
			LastMessageAt: chat.LastMessageAt,
			UnreadCount:   chat.UnreadCount,
			IsGroup:       chat.IsGroup,
			IsMuted:       chat.IsMuted,
			IsArchived:    chat.IsArchived,
			IsBlocked:     chat.IsBlocked,
		}
	}

	uc.logger.Debug(ctx, "Chats retrieved successfully",
		"sessionID", query.SessionID,
		"count", len(chatViews))

	// 6. Return result
	return &GetChatsResult{
		SessionID: query.SessionID,
		Chats:     chatViews,
		Total:     len(chatViews), // In a real implementation, this would be the total count
		Limit:     query.Limit,
		Offset:    query.Offset,
	}, nil
}
