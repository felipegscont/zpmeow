package application

import (
	"context"
	"fmt"

	"zpmeow/internal/application/commands"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/shared/validation"
)

// GroupService implements group use cases following Clean Architecture
type GroupService struct {
	wameowService interface {
		CreateGroup(ctx context.Context, sessionID, name string, participants []string) (interface{}, error)
		GetGroupInfo(ctx context.Context, sessionID, groupJID string) (interface{}, error)
		ListGroups(ctx context.Context, sessionID string) (interface{}, error)
		JoinGroup(ctx context.Context, sessionID, inviteLink string) (interface{}, error)
		JoinGroupWithInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (interface{}, error)
		LeaveGroup(ctx context.Context, sessionID, groupJID string) error
		GetInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error)
		GetInviteInfo(ctx context.Context, sessionID, inviteLink string) (interface{}, error)
		GetGroupInfoFromInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (interface{}, error)
		UpdateParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
		SetGroupName(ctx context.Context, sessionID, groupJID, name string) error
		SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error
		SetGroupPhoto(ctx context.Context, sessionID, groupJID string, photoData []byte) error
		RemoveGroupPhoto(ctx context.Context, sessionID, groupJID string) error
		SetGroupAnnounce(ctx context.Context, sessionID, groupJID string, announceOnly bool) error
		SetGroupLocked(ctx context.Context, sessionID, groupJID string, locked bool) error
		SetGroupEphemeral(ctx context.Context, sessionID, groupJID string, ephemeral bool, duration int) error
		SetGroupJoinApprovalMode(ctx context.Context, sessionID, groupJID string, requireApproval bool) error
		SetGroupMemberAddMode(ctx context.Context, sessionID, groupJID, mode string) error
		GetGroupRequestParticipants(ctx context.Context, sessionID, groupJID string) (interface{}, error)
		UpdateGroupRequestParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
		LinkGroup(ctx context.Context, sessionID, groupJID, parentGroupJID string) error
		UnlinkGroup(ctx context.Context, sessionID, groupJID string) error
		GetSubGroups(ctx context.Context, sessionID, parentGroupJID string) (interface{}, error)
		GetLinkedGroupsParticipants(ctx context.Context, sessionID, parentGroupJID string) (interface{}, error)
	}
	sessionRepo session.Repository
	validator   *validation.Validator
}

// NewGroupService creates a new GroupService instance
func NewGroupService(
	wameowService interface{},
	sessionRepo session.Repository,
	validator *validation.Validator,
) *GroupService {
	return &GroupService{
		wameowService: wameowService.(interface {
			CreateGroup(ctx context.Context, sessionID, name string, participants []string) (interface{}, error)
			GetGroupInfo(ctx context.Context, sessionID, groupJID string) (interface{}, error)
			ListGroups(ctx context.Context, sessionID string) (interface{}, error)
			JoinGroup(ctx context.Context, sessionID, inviteLink string) (interface{}, error)
			JoinGroupWithInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (interface{}, error)
			LeaveGroup(ctx context.Context, sessionID, groupJID string) error
			GetInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error)
			GetInviteInfo(ctx context.Context, sessionID, inviteLink string) (interface{}, error)
			GetGroupInfoFromInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (interface{}, error)
			UpdateParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
			SetGroupName(ctx context.Context, sessionID, groupJID, name string) error
			SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error
			SetGroupPhoto(ctx context.Context, sessionID, groupJID string, photoData []byte) error
			RemoveGroupPhoto(ctx context.Context, sessionID, groupJID string) error
			SetGroupAnnounce(ctx context.Context, sessionID, groupJID string, announceOnly bool) error
			SetGroupLocked(ctx context.Context, sessionID, groupJID string, locked bool) error
			SetGroupEphemeral(ctx context.Context, sessionID, groupJID string, ephemeral bool, duration int) error
			SetGroupJoinApprovalMode(ctx context.Context, sessionID, groupJID string, requireApproval bool) error
			SetGroupMemberAddMode(ctx context.Context, sessionID, groupJID, mode string) error
			GetGroupRequestParticipants(ctx context.Context, sessionID, groupJID string) (interface{}, error)
			UpdateGroupRequestParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
			LinkGroup(ctx context.Context, sessionID, groupJID, parentGroupJID string) error
			UnlinkGroup(ctx context.Context, sessionID, groupJID string) error
			GetSubGroups(ctx context.Context, sessionID, parentGroupJID string) (interface{}, error)
			GetLinkedGroupsParticipants(ctx context.Context, sessionID, parentGroupJID string) (interface{}, error)
		}),
		sessionRepo: sessionRepo,
		validator:   validator,
	}
}

// CreateGroup creates a new group using command pattern
func (s *GroupService) CreateGroup(ctx context.Context, cmd *commands.CreateGroupCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.CreateGroup(ctx, cmd.SessionID, cmd.Name, cmd.Participants)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// GetGroupInfo gets group information using command pattern
func (s *GroupService) GetGroupInfo(ctx context.Context, cmd *commands.GetGroupInfoCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	return s.wameowService.GetGroupInfo(ctx, cmd.SessionID, cmd.GroupJID)
}

// ListGroups lists all groups
func (s *GroupService) ListGroups(ctx context.Context, sessionID string) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.wameowService.ListGroups(ctx, sessionID)
}

// JoinGroup joins a group via invite link using command pattern
func (s *GroupService) JoinGroup(ctx context.Context, cmd *commands.JoinGroupCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.JoinGroup(ctx, cmd.SessionID, cmd.InviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// JoinGroupWithInvite joins a group with invite details using command pattern
func (s *GroupService) JoinGroupWithInvite(ctx context.Context, cmd *commands.JoinGroupWithInviteCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.JoinGroupWithInvite(ctx, cmd.SessionID, cmd.GroupJID, cmd.Inviter, cmd.Code, cmd.Expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to join group with invite: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// LeaveGroup leaves a group using command pattern
func (s *GroupService) LeaveGroup(ctx context.Context, cmd *commands.LeaveGroupCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	err := s.wameowService.LeaveGroup(ctx, cmd.SessionID, cmd.GroupJID)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	// whatsmeow will handle events automatically
	return nil
}

// GetInviteLink gets group invite link using command pattern
func (s *GroupService) GetInviteLink(ctx context.Context, cmd *commands.GetInviteLinkCommand) (string, error) {
	if err := cmd.Validate(); err != nil {
		return "", fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return "", err
	}

	return s.wameowService.GetInviteLink(ctx, cmd.SessionID, cmd.GroupJID, cmd.Reset)
}

// UpdateParticipants updates group participants using command pattern
func (s *GroupService) UpdateParticipants(ctx context.Context, cmd *commands.UpdateParticipantsCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	return s.wameowService.UpdateParticipants(ctx, cmd.SessionID, cmd.GroupJID, cmd.Action, cmd.Participants)
}

// SetGroupName sets group name using command pattern
func (s *GroupService) SetGroupName(ctx context.Context, cmd *commands.SetGroupNameCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	return s.wameowService.SetGroupName(ctx, cmd.SessionID, cmd.GroupJID, cmd.Name)
}

// Helper methods

func (s *GroupService) validateSession(ctx context.Context, sessionID string) error {
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return fmt.Errorf("session is not connected")
	}

	return nil
}
