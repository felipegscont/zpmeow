package application

import (
	"context"
	"fmt"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/interfaces/dto"
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

// CreateGroup creates a new group using DTO
func (s *GroupService) CreateGroup(ctx context.Context, sessionID string, req *dto.CreateGroupRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.CreateGroup(ctx, sessionID, req.Name, req.Participants)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// GetGroupInfo gets group information using DTO
func (s *GroupService) GetGroupInfo(ctx context.Context, sessionID string, req *dto.GetGroupInfoRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.wameowService.GetGroupInfo(ctx, sessionID, req.GroupJID)
}

// ListGroups lists all groups
func (s *GroupService) ListGroups(ctx context.Context, sessionID string) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.wameowService.ListGroups(ctx, sessionID)
}

// JoinGroup joins a group via invite link using DTO
func (s *GroupService) JoinGroup(ctx context.Context, sessionID string, req *dto.JoinGroupRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.JoinGroup(ctx, sessionID, req.InviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// JoinGroupWithInvite joins a group with invite details using DTO
func (s *GroupService) JoinGroupWithInvite(ctx context.Context, sessionID string, req *dto.JoinGroupWithInviteRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.JoinGroupWithInvite(ctx, sessionID, req.GroupJID, req.Inviter, req.Code, req.Expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to join group with invite: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// LeaveGroup leaves a group using DTO
func (s *GroupService) LeaveGroup(ctx context.Context, sessionID string, req *dto.LeaveGroupRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	err := s.wameowService.LeaveGroup(ctx, sessionID, req.GroupJID)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	// whatsmeow will handle events automatically
	return nil
}

// GetInviteLink gets group invite link using DTO
func (s *GroupService) GetInviteLink(ctx context.Context, sessionID string, req *dto.GetInviteLinkRequest) (string, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return "", err
	}

	return s.wameowService.GetInviteLink(ctx, sessionID, req.GroupJID, req.Reset)
}

// UpdateParticipants updates group participants using DTO
func (s *GroupService) UpdateParticipants(ctx context.Context, sessionID string, req *dto.UpdateParticipantsRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	return s.wameowService.UpdateParticipants(ctx, sessionID, req.GroupJID, req.Action, req.Participants)
}

// SetGroupName sets group name using DTO
func (s *GroupService) SetGroupName(ctx context.Context, sessionID string, req *dto.SetGroupNameRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	return s.wameowService.SetGroupName(ctx, sessionID, req.GroupJID, req.Name)
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
