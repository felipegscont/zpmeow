package commands

import (
	"fmt"
	"strings"
)

// CreateGroupCommand represents a command to create a group
type CreateGroupCommand struct {
	SessionID    string   `json:"session_id" validate:"required"`
	Name         string   `json:"name" validate:"required"`
	Participants []string `json:"participants" validate:"required"`
}

// JoinGroupCommand represents a command to join a group via invite link
type JoinGroupCommand struct {
	SessionID  string `json:"session_id" validate:"required"`
	InviteLink string `json:"invite_link" validate:"required"`
}

// JoinGroupWithInviteCommand represents a command to join a group with invite details
type JoinGroupWithInviteCommand struct {
	SessionID  string `json:"session_id" validate:"required"`
	GroupJID   string `json:"group_jid" validate:"required"`
	Inviter    string `json:"inviter" validate:"required"`
	Code       string `json:"code" validate:"required"`
	Expiration int64  `json:"expiration"`
}

// LeaveGroupCommand represents a command to leave a group
type LeaveGroupCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
}

// GetGroupInfoCommand represents a command to get group information
type GetGroupInfoCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
}

// GetInviteLinkCommand represents a command to get group invite link
type GetInviteLinkCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
	Reset     bool   `json:"reset,omitempty"`
}

// GetInviteInfoCommand represents a command to get invite information
type GetInviteInfoCommand struct {
	SessionID  string `json:"session_id" validate:"required"`
	InviteLink string `json:"invite_link" validate:"required"`
}

// GetGroupInfoFromInviteCommand represents a command to get group info from invite
type GetGroupInfoFromInviteCommand struct {
	SessionID  string `json:"session_id" validate:"required"`
	GroupJID   string `json:"group_jid" validate:"required"`
	Inviter    string `json:"inviter" validate:"required"`
	Code       string `json:"code" validate:"required"`
	Expiration int64  `json:"expiration"`
}

// UpdateParticipantsCommand represents a command to update group participants
type UpdateParticipantsCommand struct {
	SessionID    string   `json:"session_id" validate:"required"`
	GroupJID     string   `json:"group_jid" validate:"required"`
	Action       string   `json:"action" validate:"required"` // add, remove, promote, demote
	Participants []string `json:"participants" validate:"required"`
}

// SetGroupNameCommand represents a command to set group name
type SetGroupNameCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

// SetGroupTopicCommand represents a command to set group topic
type SetGroupTopicCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
	Topic     string `json:"topic" validate:"required"`
}

// SetGroupPhotoCommand represents a command to set group photo
type SetGroupPhotoCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
	Photo     string `json:"photo" validate:"required"` // base64, URL, or file path
}

// RemoveGroupPhotoCommand represents a command to remove group photo
type RemoveGroupPhotoCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
}

// SetGroupAnnounceCommand represents a command to set group announce mode
type SetGroupAnnounceCommand struct {
	SessionID    string `json:"session_id" validate:"required"`
	GroupJID     string `json:"group_jid" validate:"required"`
	AnnounceOnly bool   `json:"announce_only"`
}

// SetGroupLockedCommand represents a command to set group locked mode
type SetGroupLockedCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
	Locked    bool   `json:"locked"`
}

// SetGroupEphemeralCommand represents a command to set group ephemeral mode
type SetGroupEphemeralCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
	Ephemeral bool   `json:"ephemeral"`
	Duration  int    `json:"duration,omitempty"` // in seconds
}

// SetGroupJoinApprovalCommand represents a command to set group join approval mode
type SetGroupJoinApprovalCommand struct {
	SessionID       string `json:"session_id" validate:"required"`
	GroupJID        string `json:"group_jid" validate:"required"`
	RequireApproval bool   `json:"require_approval"`
}

// SetGroupMemberAddModeCommand represents a command to set group member add mode
type SetGroupMemberAddModeCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
	Mode      string `json:"mode" validate:"required"` // all_members, admins_only
}

// GetGroupRequestParticipantsCommand represents a command to get group request participants
type GetGroupRequestParticipantsCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
}

// UpdateGroupRequestParticipantsCommand represents a command to update group request participants
type UpdateGroupRequestParticipantsCommand struct {
	SessionID    string   `json:"session_id" validate:"required"`
	GroupJID     string   `json:"group_jid" validate:"required"`
	Action       string   `json:"action" validate:"required"` // approve, reject
	Participants []string `json:"participants" validate:"required"`
}

// Community commands

// LinkGroupCommand represents a command to link a group to a community
type LinkGroupCommand struct {
	SessionID      string `json:"session_id" validate:"required"`
	GroupJID       string `json:"group_jid" validate:"required"`
	ParentGroupJID string `json:"parent_group_jid" validate:"required"`
}

// UnlinkGroupCommand represents a command to unlink a group from a community
type UnlinkGroupCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	GroupJID  string `json:"group_jid" validate:"required"`
}

// GetSubGroupsCommand represents a command to get sub groups of a community
type GetSubGroupsCommand struct {
	SessionID      string `json:"session_id" validate:"required"`
	ParentGroupJID string `json:"parent_group_jid" validate:"required"`
}

// GetLinkedGroupsParticipantsCommand represents a command to get linked groups participants
type GetLinkedGroupsParticipantsCommand struct {
	SessionID      string `json:"session_id" validate:"required"`
	ParentGroupJID string `json:"parent_group_jid" validate:"required"`
}

// Validation methods

func (c *CreateGroupCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("group name is required")
	}
	if len(c.Participants) == 0 {
		return fmt.Errorf("at least one participant is required")
	}
	return nil
}

func (c *JoinGroupCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.InviteLink == "" {
		return fmt.Errorf("invite link is required")
	}
	return nil
}

func (c *JoinGroupWithInviteCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	if c.Inviter == "" {
		return fmt.Errorf("inviter is required")
	}
	if c.Code == "" {
		return fmt.Errorf("code is required")
	}
	return nil
}

func (c *LeaveGroupCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	return nil
}

func (c *GetGroupInfoCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	return nil
}

func (c *UpdateParticipantsCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	if c.Action == "" {
		return fmt.Errorf("action is required")
	}
	if len(c.Participants) == 0 {
		return fmt.Errorf("at least one participant is required")
	}

	// Validate action
	validActions := []string{"add", "remove", "promote", "demote"}
	isValidAction := false
	for _, validAction := range validActions {
		if c.Action == validAction {
			isValidAction = true
			break
		}
	}
	if !isValidAction {
		return fmt.Errorf("invalid action: %s. Valid actions are: %v", c.Action, validActions)
	}

	return nil
}

func (c *SetGroupNameCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("group name is required")
	}
	return nil
}

func (c *SetGroupMemberAddModeCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	if c.Mode == "" {
		return fmt.Errorf("mode is required")
	}

	// Validate mode
	validModes := []string{"all_members", "admins_only"}
	isValidMode := false
	for _, validMode := range validModes {
		if c.Mode == validMode {
			isValidMode = true
			break
		}
	}
	if !isValidMode {
		return fmt.Errorf("invalid mode: %s. Valid modes are: %v", c.Mode, validModes)
	}

	return nil
}

func (c *GetInviteLinkCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	return nil
}

func (c *LinkGroupCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	if c.ParentGroupJID == "" {
		return fmt.Errorf("parent group JID is required")
	}
	return nil
}

func (c *UnlinkGroupCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.GroupJID == "" {
		return fmt.Errorf("group JID is required")
	}
	return nil
}

func (c *GetSubGroupsCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ParentGroupJID == "" {
		return fmt.Errorf("parent group JID is required")
	}
	return nil
}

func (c *GetLinkedGroupsParticipantsCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ParentGroupJID == "" {
		return fmt.Errorf("parent group JID is required")
	}
	return nil
}
