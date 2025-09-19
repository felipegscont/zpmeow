package group

import (
	"context"
	"fmt"
	"strings"

	"zpmeow/internal/application/common"
	"zpmeow/internal/application/ports"
)

// CreateGroupCommand represents the command to create a new group
type CreateGroupCommand struct {
	SessionID    string
	Name         string
	Description  string
	Participants []string
}

// Validate validates the create group command
func (c CreateGroupCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.Name) == "" {
		return common.NewValidationError("name", c.Name, "group name is required")
	}

	if len(c.Name) > 100 {
		return common.NewValidationError("name", c.Name, "group name must not exceed 100 characters")
	}

	if len(c.Description) > 500 {
		return common.NewValidationError("description", c.Description, "group description must not exceed 500 characters")
	}

	if len(c.Participants) == 0 {
		return common.NewValidationError("participants", "", "at least one participant is required")
	}

	if len(c.Participants) > 256 {
		return common.NewValidationError("participants", "", "maximum 256 participants allowed")
	}

	// Validate participant JIDs
	for i, participant := range c.Participants {
		if strings.TrimSpace(participant) == "" {
			return common.NewValidationError("participants", participant, fmt.Sprintf("participant %d cannot be empty", i))
		}
	}

	return nil
}

// GroupView represents a group view model
type GroupView struct {
	JID          string
	Name         string
	Description  string
	Participants []string
	Admins       []string
	Owner        string
	CreatedAt    string
	IsAnnounce   bool
	IsLocked     bool
}

// CreateGroupResult represents the result of creating a group
type CreateGroupResult struct {
	SessionID string
	Group     GroupView
	Success   bool
}

// CreateGroupUseCase handles creating new WhatsApp groups
type CreateGroupUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewCreateGroupUseCase creates a new CreateGroupUseCase
func NewCreateGroupUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *CreateGroupUseCase {
	return &CreateGroupUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the create group use case
func (uc *CreateGroupUseCase) Handle(ctx context.Context, cmd CreateGroupCommand) (*CreateGroupResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid create group command", "error", err)
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
			fmt.Sprintf("session must be connected to create groups, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Check if session is authenticated
	if !sessionEntity.IsAuthenticated() {
		return nil, common.NewBusinessRuleError(
			"session_not_authenticated",
			"session must be authenticated to create groups",
		)
	}

	// 5. Create group via WhatsApp service
	groupJID, err := uc.whatsappService.CreateGroup(ctx, cmd.SessionID, cmd.Name, cmd.Participants)
	if err != nil {
		uc.logger.Error(ctx, "Failed to create group",
			"sessionID", cmd.SessionID,
			"groupName", cmd.Name,
			"participantCount", len(cmd.Participants),
			"error", err)
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// 6. Get group info to return complete data
	groupInfo, err := uc.whatsappService.GetGroupInfo(ctx, cmd.SessionID, groupJID)
	if err != nil {
		uc.logger.Warn(ctx, "Failed to get group info after creation",
			"sessionID", cmd.SessionID,
			"groupJID", groupJID,
			"error", err)
		// Don't fail the operation, just return basic info
		groupInfo = &ports.GroupInfo{
			JID:          groupJID,
			Name:         cmd.Name,
			Description:  cmd.Description,
			Participants: cmd.Participants,
		}
	}

	// 7. Convert to view model
	groupView := GroupView{
		JID:          groupInfo.JID,
		Name:         groupInfo.Name,
		Description:  groupInfo.Description,
		Participants: groupInfo.Participants,
		Admins:       groupInfo.Admins,
		Owner:        groupInfo.Owner,
		CreatedAt:    groupInfo.CreatedAt,
		IsAnnounce:   groupInfo.IsAnnounce,
		IsLocked:     groupInfo.IsLocked,
	}

	uc.logger.Info(ctx, "Group created successfully",
		"sessionID", cmd.SessionID,
		"groupJID", groupJID,
		"groupName", cmd.Name,
		"participantCount", len(cmd.Participants))

	// 8. Return result
	return &CreateGroupResult{
		SessionID: cmd.SessionID,
		Group:     groupView,
		Success:   true,
	}, nil
}
