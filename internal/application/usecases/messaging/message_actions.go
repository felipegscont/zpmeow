package messaging

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// MarkAsReadCommand represents the command to mark messages as read
type MarkAsReadCommand struct {
	SessionID string
	ChatJID   string
	MessageID string
}

// Validate validates the mark as read command
func (c MarkAsReadCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if strings.TrimSpace(c.MessageID) == "" {
		return common.NewValidationError("messageID", c.MessageID, "message ID is required")
	}

	return nil
}

// ReactToMessageCommand represents the command to react to a message
type ReactToMessageCommand struct {
	SessionID string
	ChatJID   string
	MessageID string
	Emoji     string
	Remove    bool
}

// Validate validates the react to message command
func (c ReactToMessageCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if strings.TrimSpace(c.MessageID) == "" {
		return common.NewValidationError("messageID", c.MessageID, "message ID is required")
	}

	if !c.Remove && strings.TrimSpace(c.Emoji) == "" {
		return common.NewValidationError("emoji", c.Emoji, "emoji is required when not removing reaction")
	}

	return nil
}

// EditMessageCommand represents the command to edit a message
type EditMessageCommand struct {
	SessionID  string
	ChatJID    string
	MessageID  string
	NewContent string
}

// Validate validates the edit message command
func (c EditMessageCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if strings.TrimSpace(c.MessageID) == "" {
		return common.NewValidationError("messageID", c.MessageID, "message ID is required")
	}

	if strings.TrimSpace(c.NewContent) == "" {
		return common.NewValidationError("newContent", c.NewContent, "new content is required")
	}

	if len(c.NewContent) > 4096 {
		return common.NewValidationError("newContent", c.NewContent, "new content must not exceed 4096 characters")
	}

	return nil
}

// DeleteMessageCommand represents the command to delete a message
type DeleteMessageCommand struct {
	SessionID   string
	ChatJID     string
	MessageID   string
	ForEveryone bool
}

// Validate validates the delete message command
func (c DeleteMessageCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if strings.TrimSpace(c.MessageID) == "" {
		return common.NewValidationError("messageID", c.MessageID, "message ID is required")
	}

	return nil
}

// MessageActionResult represents the result of message actions
type MessageActionResult struct {
	SessionID string
	ChatJID   string
	MessageID string
	Action    string
	Success   bool
	Message   string
}

// MarkAsReadUseCase handles marking messages as read
type MarkAsReadUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewMarkAsReadUseCase creates a new MarkAsReadUseCase
func NewMarkAsReadUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *MarkAsReadUseCase {
	return &MarkAsReadUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the mark as read use case
func (uc *MarkAsReadUseCase) Handle(ctx context.Context, cmd MarkAsReadCommand) (*MessageActionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid mark as read command", "error", err)
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
			fmt.Sprintf("session must be connected to mark messages as read, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Mark message as read via WhatsApp service
	if err := uc.whatsappService.MarkAsRead(ctx, cmd.SessionID, cmd.ChatJID, cmd.MessageID); err != nil {
		uc.logger.Error(ctx, "Failed to mark message as read",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"messageID", cmd.MessageID,
			"error", err)
		return nil, fmt.Errorf("failed to mark message as read: %w", err)
	}

	uc.logger.Info(ctx, "Message marked as read successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"messageID", cmd.MessageID)

	// 5. Return result
	return &MessageActionResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		MessageID: cmd.MessageID,
		Action:    "mark_as_read",
		Success:   true,
		Message:   "Message marked as read successfully",
	}, nil
}

// ReactToMessageUseCase handles reacting to messages
type ReactToMessageUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewReactToMessageUseCase creates a new ReactToMessageUseCase
func NewReactToMessageUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *ReactToMessageUseCase {
	return &ReactToMessageUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the react to message use case
func (uc *ReactToMessageUseCase) Handle(ctx context.Context, cmd ReactToMessageCommand) (*MessageActionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid react to message command", "error", err)
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
			fmt.Sprintf("session must be connected to react to messages, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. React to message via WhatsApp service
	if err := uc.whatsappService.ReactToMessage(ctx, cmd.SessionID, cmd.ChatJID, cmd.MessageID, cmd.Emoji, cmd.Remove); err != nil {
		uc.logger.Error(ctx, "Failed to react to message",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"messageID", cmd.MessageID,
			"emoji", cmd.Emoji,
			"error", err)
		return nil, fmt.Errorf("failed to react to message: %w", err)
	}

	action := "add_reaction"
	if cmd.Remove {
		action = "remove_reaction"
	}

	uc.logger.Info(ctx, "Message reaction updated successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"messageID", cmd.MessageID,
		"action", action)

	// 5. Return result
	return &MessageActionResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		MessageID: cmd.MessageID,
		Action:    action,
		Success:   true,
		Message:   fmt.Sprintf("Message reaction %s successfully", action),
	}, nil
}

// EditMessageUseCase handles editing messages
type EditMessageUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewEditMessageUseCase creates a new EditMessageUseCase
func NewEditMessageUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *EditMessageUseCase {
	return &EditMessageUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the edit message use case
func (uc *EditMessageUseCase) Handle(ctx context.Context, cmd EditMessageCommand) (*MessageActionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid edit message command", "error", err)
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
			fmt.Sprintf("session must be connected to edit messages, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Edit message via WhatsApp service
	if err := uc.whatsappService.EditMessage(ctx, cmd.SessionID, cmd.ChatJID, cmd.MessageID, cmd.NewContent); err != nil {
		uc.logger.Error(ctx, "Failed to edit message",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"messageID", cmd.MessageID,
			"error", err)
		return nil, fmt.Errorf("failed to edit message: %w", err)
	}

	uc.logger.Info(ctx, "Message edited successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"messageID", cmd.MessageID)

	// 5. Return result
	return &MessageActionResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		MessageID: cmd.MessageID,
		Action:    "edit",
		Success:   true,
		Message:   "Message edited successfully",
	}, nil
}

// DeleteMessageUseCase handles deleting messages
type DeleteMessageUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewDeleteMessageUseCase creates a new DeleteMessageUseCase
func NewDeleteMessageUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *DeleteMessageUseCase {
	return &DeleteMessageUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the delete message use case
func (uc *DeleteMessageUseCase) Handle(ctx context.Context, cmd DeleteMessageCommand) (*MessageActionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid delete message command", "error", err)
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
			fmt.Sprintf("session must be connected to delete messages, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Delete message via WhatsApp service
	if err := uc.whatsappService.DeleteMessage(ctx, cmd.SessionID, cmd.ChatJID, cmd.MessageID, cmd.ForEveryone); err != nil {
		uc.logger.Error(ctx, "Failed to delete message",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"messageID", cmd.MessageID,
			"forEveryone", cmd.ForEveryone,
			"error", err)
		return nil, fmt.Errorf("failed to delete message: %w", err)
	}

	action := "delete_for_me"
	if cmd.ForEveryone {
		action = "delete_for_everyone"
	}

	uc.logger.Info(ctx, "Message deleted successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"messageID", cmd.MessageID,
		"action", action)

	// 5. Return result
	return &MessageActionResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		MessageID: cmd.MessageID,
		Action:    action,
		Success:   true,
		Message:   fmt.Sprintf("Message %s successfully", action),
	}, nil
}
