package group

import (
	"context"
	"fmt"
	"strings"

	"zpmeow/internal/application/common"
	"zpmeow/internal/application/ports"
)

// JoinGroupCommand represents the command to join a group via invite link
type JoinGroupCommand struct {
	SessionID  string
	InviteLink string
}

// Validate validates the join group command
func (c JoinGroupCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.InviteLink) == "" {
		return common.NewValidationError("inviteLink", c.InviteLink, "invite link is required")
	}

	if !strings.Contains(c.InviteLink, "chat.whatsapp.com") {
		return common.NewValidationError("inviteLink", c.InviteLink, "invalid WhatsApp invite link format")
	}

	return nil
}

// LeaveGroupCommand represents the command to leave a group
type LeaveGroupCommand struct {
	SessionID string
	GroupJID  string
}

// Validate validates the leave group command
func (c LeaveGroupCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.GroupJID) == "" {
		return common.NewValidationError("groupJID", c.GroupJID, "group JID is required")
	}

	return nil
}

// GroupManagementResult represents the result of group management operations
type GroupManagementResult struct {
	SessionID string
	GroupJID  string
	Action    string
	Success   bool
	Message   string
}

// JoinGroupUseCase handles joining groups via invite links
type JoinGroupUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewJoinGroupUseCase creates a new JoinGroupUseCase
func NewJoinGroupUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *JoinGroupUseCase {
	return &JoinGroupUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the join group use case
func (uc *JoinGroupUseCase) Handle(ctx context.Context, cmd JoinGroupCommand) (*GroupManagementResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid join group command", "error", err)
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
			fmt.Sprintf("session must be connected to join groups, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Join group via WhatsApp service
	if err := uc.whatsappService.JoinGroup(ctx, cmd.SessionID, cmd.InviteLink); err != nil {
		uc.logger.Error(ctx, "Failed to join group",
			"sessionID", cmd.SessionID,
			"inviteLink", cmd.InviteLink,
			"error", err)
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	uc.logger.Info(ctx, "Successfully joined group",
		"sessionID", cmd.SessionID,
		"inviteLink", cmd.InviteLink)

	// 5. Return result
	return &GroupManagementResult{
		SessionID: cmd.SessionID,
		GroupJID:  "", // Group JID would be extracted from invite link
		Action:    "join",
		Success:   true,
		Message:   "Successfully joined group",
	}, nil
}

// LeaveGroupUseCase handles leaving groups
type LeaveGroupUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewLeaveGroupUseCase creates a new LeaveGroupUseCase
func NewLeaveGroupUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *LeaveGroupUseCase {
	return &LeaveGroupUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the leave group use case
func (uc *LeaveGroupUseCase) Handle(ctx context.Context, cmd LeaveGroupCommand) (*GroupManagementResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid leave group command", "error", err)
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
			fmt.Sprintf("session must be connected to leave groups, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Leave group via WhatsApp service
	if err := uc.whatsappService.LeaveGroup(ctx, cmd.SessionID, cmd.GroupJID); err != nil {
		uc.logger.Error(ctx, "Failed to leave group",
			"sessionID", cmd.SessionID,
			"groupJID", cmd.GroupJID,
			"error", err)
		return nil, fmt.Errorf("failed to leave group: %w", err)
	}

	uc.logger.Info(ctx, "Successfully left group",
		"sessionID", cmd.SessionID,
		"groupJID", cmd.GroupJID)

	// 5. Return result
	return &GroupManagementResult{
		SessionID: cmd.SessionID,
		GroupJID:  cmd.GroupJID,
		Action:    "leave",
		Success:   true,
		Message:   "Successfully left group",
	}, nil
}
