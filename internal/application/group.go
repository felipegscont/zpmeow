package application

import (
	"context"
	"fmt"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

type GroupApp struct {
	groupService GroupService
	sessionRepo  session.Repository
	validator    *validation.Validator
}

func NewGroupApp(
	groupService GroupService,
	sessionRepo session.Repository,
	validator *validation.Validator,
) *GroupApp {
	return &GroupApp{
		groupService: groupService,
		sessionRepo:  sessionRepo,
		validator:    validator,
	}
}

func (s *GroupApp) CreateGroup(ctx context.Context, sessionID string, req *dto.CreateGroupRequest) (*GroupInfo, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.groupService.CreateGroup(ctx, sessionID, req.Name, req.Participants)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return result, nil
}

func (s *GroupApp) GetGroupInfo(ctx context.Context, sessionID string, req *dto.GetGroupInfoRequest) (*GroupInfo, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.groupService.GetGroupInfo(ctx, sessionID, req.GroupJID)
}

func (s *GroupApp) ListGroups(ctx context.Context, sessionID string) (*GroupList, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.groupService.ListGroups(ctx, sessionID)
}

func (s *GroupApp) JoinGroup(ctx context.Context, sessionID string, req *dto.JoinGroupRequest) (*GroupInfo, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.groupService.JoinGroup(ctx, sessionID, req.InviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	return result, nil
}

func (s *GroupApp) JoinGroupWithInvite(ctx context.Context, sessionID string, req *dto.JoinGroupWithInviteRequest) (*GroupInfo, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.groupService.JoinGroupWithInvite(ctx, sessionID, req.GroupJID, req.Inviter, req.Code, req.Expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to join group with invite: %w", err)
	}

	return result, nil
}

func (s *GroupApp) LeaveGroup(ctx context.Context, sessionID string, req *dto.LeaveGroupRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	err := s.groupService.LeaveGroup(ctx, sessionID, req.GroupJID)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	return nil
}

func (s *GroupApp) GetInviteLink(ctx context.Context, sessionID string, req *dto.GetInviteLinkRequest) (string, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return "", err
	}

	return s.groupService.GetInviteLink(ctx, sessionID, req.GroupJID, req.Reset)
}

func (s *GroupApp) UpdateParticipants(ctx context.Context, sessionID string, req *dto.UpdateParticipantsRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	return s.groupService.UpdateParticipants(ctx, sessionID, req.GroupJID, req.Action, req.Participants)
}

func (s *GroupApp) SetGroupName(ctx context.Context, sessionID string, req *dto.SetGroupNameRequest) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	return s.groupService.SetGroupName(ctx, sessionID, req.GroupJID, req.Name)
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
