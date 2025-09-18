package chat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// MuteChatCommand represents the command to mute/unmute a chat
type MuteChatCommand struct {
	SessionID string
	ChatJID   string
	Mute      bool
	Duration  time.Duration // Duration to mute (0 for permanent)
}

// Validate validates the mute chat command
func (c MuteChatCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if c.Mute && c.Duration < 0 {
		return common.NewValidationError("duration", c.Duration, "duration cannot be negative")
	}

	return nil
}

// ArchiveChatCommand represents the command to archive/unarchive a chat
type ArchiveChatCommand struct {
	SessionID string
	ChatJID   string
	Archive   bool
}

// Validate validates the archive chat command
func (c ArchiveChatCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	return nil
}

// BlockChatCommand represents the command to block/unblock a chat
type BlockChatCommand struct {
	SessionID string
	ChatJID   string
	Block     bool
}

// Validate validates the block chat command
func (c BlockChatCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	return nil
}

// ChatManagementResult represents the result of chat management operations
type ChatManagementResult struct {
	SessionID string
	ChatJID   string
	Action    string
	Success   bool
}

// MuteChatUseCase handles muting/unmuting chats
type MuteChatUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewMuteChatUseCase creates a new MuteChatUseCase
func NewMuteChatUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *MuteChatUseCase {
	return &MuteChatUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the mute chat use case
func (uc *MuteChatUseCase) Handle(ctx context.Context, cmd MuteChatCommand) (*ChatManagementResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid mute chat command", "error", err)
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
			fmt.Sprintf("session must be connected to manage chats, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Mute/unmute chat via WhatsApp service
	if err := uc.whatsappService.MuteChat(ctx, cmd.SessionID, cmd.ChatJID, cmd.Mute, cmd.Duration); err != nil {
		uc.logger.Error(ctx, "Failed to mute/unmute chat",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"mute", cmd.Mute,
			"error", err)
		return nil, fmt.Errorf("failed to mute/unmute chat: %w", err)
	}

	action := "unmute"
	if cmd.Mute {
		action = "mute"
	}

	uc.logger.Info(ctx, "Chat mute status changed successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"action", action)

	// 5. Return result
	return &ChatManagementResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		Action:    action,
		Success:   true,
	}, nil
}

// ArchiveChatUseCase handles archiving/unarchiving chats
type ArchiveChatUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewArchiveChatUseCase creates a new ArchiveChatUseCase
func NewArchiveChatUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *ArchiveChatUseCase {
	return &ArchiveChatUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the archive chat use case
func (uc *ArchiveChatUseCase) Handle(ctx context.Context, cmd ArchiveChatCommand) (*ChatManagementResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid archive chat command", "error", err)
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
			fmt.Sprintf("session must be connected to manage chats, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Archive/unarchive chat via WhatsApp service
	if err := uc.whatsappService.ArchiveChat(ctx, cmd.SessionID, cmd.ChatJID, cmd.Archive); err != nil {
		uc.logger.Error(ctx, "Failed to archive/unarchive chat",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"archive", cmd.Archive,
			"error", err)
		return nil, fmt.Errorf("failed to archive/unarchive chat: %w", err)
	}

	action := "unarchive"
	if cmd.Archive {
		action = "archive"
	}

	uc.logger.Info(ctx, "Chat archive status changed successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"action", action)

	// 5. Return result
	return &ChatManagementResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		Action:    action,
		Success:   true,
	}, nil
}
