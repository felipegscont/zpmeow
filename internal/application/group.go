package application

import (
	"context"
	"fmt"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

type GroupApp struct {
	meowService interface {
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

func NewGroupApp(
	meowService interface{},
	sessionRepo session.Repository,
	validator *validation.Validator,
) *GroupApp {
	return &GroupApp{
		meowService: meowService.(interface {
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

func (s *GroupApp) CreateGroup(ctx context.Context, sessionID string, req *dto.CreateGroupRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.meowService.CreateGroup(ctx, sessionID, req.Name, req.Participants)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return result, nil
}

func (s *GroupApp) GetGroupInfo(ctx context.Context, sessionID string, req *dto.GetGroupInfoRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.meowService.GetGroupInfo(ctx, sessionID, req.GroupJID)
}

func (s *GroupApp) ListGroups(ctx context.Context, sessionID string) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.meowService.ListGroups(ctx, sessionID)
}

func (s *GroupApp) JoinGroup(ctx context.Context, sessionID string, req *dto.JoinGroupRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.meowService.JoinGroup(ctx, sessionID, req.InviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	return result, nil
}

func (s *GroupApp) JoinGroupWithInvite(ctx context.Context, sessionID string, req *dto.JoinGroupWithInviteRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.meowService.JoinGroupWithInvite(ctx, sessionID, req.GroupJID, req.Inviter, req.Code, req.Expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to join group with invite: %w", err)
	}

	return result, nil
}

func (s *GroupApp) LeaveGroup(ctx context.Context, sessionID string, req *dto.LeaveGroupRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	err := s.meowService.LeaveGroup(ctx, sessionID, req.GroupJID)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	return nil
}

func (s *GroupApp) GetInviteLink(ctx context.Context, sessionID string, req *dto.GetInviteLinkRequest) (string, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return "", err
	}

	return s.meowService.GetInviteLink(ctx, sessionID, req.GroupJID, req.Reset)
}

func (s *GroupApp) UpdateParticipants(ctx context.Context, sessionID string, req *dto.UpdateParticipantsRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	return s.meowService.UpdateParticipants(ctx, sessionID, req.GroupJID, req.Action, req.Participants)
}

func (s *GroupApp) SetGroupName(ctx context.Context, sessionID string, req *dto.SetGroupNameRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	return s.meowService.SetGroupName(ctx, sessionID, req.GroupJID, req.Name)
}


func (s *GroupApp) validateSession(ctx context.Context, sessionID string) error {
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return fmt.Errorf("session is not connected")
	}

	return nil
}
