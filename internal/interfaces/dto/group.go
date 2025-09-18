package dto

import (
	"time"
)

type CreateGroupRequest struct {
	Name         string   `json:"name" binding:"required" example:"My Group"`
	Participants []string `json:"participants" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

type GetGroupInfoRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

type JoinGroupRequest struct {
	InviteLink string `json:"invite_link" binding:"required" example:"https://chat.meow.com/ABC123"`
}

type JoinGroupWithInviteRequest struct {
	GroupJID   string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Inviter    string `json:"inviter" binding:"required" example:"5511999999999@s.meow.net"`
	Code       string `json:"code" binding:"required" example:"ABC123DEF456"`
	Expiration int64  `json:"expiration" binding:"required" example:"1640995200"`
}

type LeaveGroupRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

type GetInviteLinkRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Reset    bool   `json:"reset,omitempty" example:"false"`
}

type GetInviteInfoRequest struct {
	InviteLink string `json:"invite_link" binding:"required" example:"https://chat.meow.com/ABC123"`
}

type GroupInviteInfoReq struct {
	GroupJID   string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Inviter    string `json:"inviter" binding:"required" example:"5511999999999@s.meow.net"`
	Code       string `json:"code" binding:"required" example:"ABC123DEF456"`
	Expiration int64  `json:"expiration" binding:"required" example:"1640995200"`
}

type UpdateParticipantsRequest struct {
	GroupJID     string   `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Action       string   `json:"action" binding:"required" example:"add"`
	Participants []string `json:"participants" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

type SetGroupNameRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Name     string `json:"name" binding:"required" example:"New Group Name"`
}

type SetGroupTopicRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Topic    string `json:"topic" binding:"required" example:"Group topic description"`
}

type SetGroupPhotoRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Photo    string `json:"photo" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD..."`
}

type RemoveGroupPhotoRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

type SetGroupAnnounceRequest struct {
	GroupJID     string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	AnnounceOnly bool   `json:"announce_only" example:"true"`
}

type SetGroupLockedRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Locked   bool   `json:"locked" example:"true"`
}

type SetGroupEphemeralRequest struct {
	GroupJID  string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Ephemeral bool   `json:"ephemeral" example:"true"`
	Duration  int    `json:"duration,omitempty" example:"86400"`
}

type GroupJoinApprovalReq struct {
	GroupJID        string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	RequireApproval bool   `json:"require_approval" binding:"required" example:"true"`
}

type GroupMemberModeReq struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Mode     string `json:"mode" binding:"required" example:"admin" enum:"all,admin"`
}

type GetGroupRequestsReq struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

type UpdateGroupRequestsReq struct {
	GroupJID     string   `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Action       string   `json:"action" binding:"required" example:"approve"`
	Participants []string `json:"participants" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

type GroupInfo struct {
	JID          string   `json:"jid" example:"120363025246125486@g.us"`
	Name         string   `json:"name" example:"My Group"`
	Topic        string   `json:"topic,omitempty" example:"Group topic description"`
	Participants []string `json:"participants" example:"[\"5511999999999@s.meow.net\", \"5511888888888@s.meow.net\"]"`
	Admins       []string `json:"admins" example:"[\"5511999999999@s.meow.net\"]"`
	Owner        string   `json:"owner" example:"5511999999999@s.meow.net"`
	CreatedAt    int64    `json:"created_at" example:"1640995200"`
	Size         int      `json:"size" example:"2"`
	Announce     bool     `json:"announce" example:"false"`
	Locked       bool     `json:"locked" example:"false"`
	Ephemeral    bool     `json:"ephemeral" example:"false"`
}

type GroupResponse struct {
	Success bool                `json:"success"`
	Code    int                 `json:"code"`
	Data    GroupData           `json:"data"`
	Error   *GroupErrorResponse `json:"error,omitempty"`
}

type GroupData struct {
	SessionID  string      `json:"session_id,omitempty" example:"default"`
	Action     string      `json:"action" example:"create"`
	Status     string      `json:"status" example:"success"`
	Timestamp  time.Time   `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group      *GroupInfo  `json:"group,omitempty"`
	Groups     []GroupInfo `json:"groups,omitempty"`
	InviteLink string      `json:"invite_link,omitempty" example:"https://chat.meow.com/ABC123"`
	Message    string      `json:"message,omitempty" example:"Operation completed successfully"`
}

type GroupErrorResponse struct {
	Code    string `json:"code" example:"INVALID_GROUP_JID"`
	Message string `json:"message" example:"Invalid group JID format"`
	Details string `json:"details,omitempty" example:"Group JID must end with @g.us"`
}

type CreateGroupResponse struct {
	Success bool                    `json:"success" example:"true"`
	Code    int                     `json:"code" example:"201"`
	Data    CreateGroupResponseData `json:"data"`
	Error   *GroupErrorResponse     `json:"error,omitempty"`
}

type CreateGroupResponseData struct {
	SessionID string     `json:"session_id" example:"default"`
	Action    string     `json:"action" example:"create"`
	Status    string     `json:"status" example:"success"`
	Timestamp time.Time  `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group     *GroupInfo `json:"group"`
}

type ListGroupsResponse struct {
	Success bool                   `json:"success" example:"true"`
	Code    int                    `json:"code" example:"200"`
	Data    ListGroupsResponseData `json:"data"`
	Error   *GroupErrorResponse    `json:"error,omitempty"`
}

type ListGroupsResponseData struct {
	SessionID string      `json:"session_id" example:"default"`
	Action    string      `json:"action" example:"list"`
	Status    string      `json:"status" example:"success"`
	Timestamp time.Time   `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Groups    []GroupInfo `json:"groups"`
	Total     int         `json:"total" example:"5"`
}

type GetGroupInfoResponse struct {
	Success bool                     `json:"success" example:"true"`
	Code    int                      `json:"code" example:"200"`
	Data    GetGroupInfoResponseData `json:"data"`
	Error   *GroupErrorResponse      `json:"error,omitempty"`
}

type GetGroupInfoResponseData struct {
	SessionID string     `json:"session_id" example:"default"`
	Action    string     `json:"action" example:"info"`
	Status    string     `json:"status" example:"success"`
	Timestamp time.Time  `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group     *GroupInfo `json:"group"`
}

type JoinGroupResponse struct {
	Success bool                  `json:"success" example:"true"`
	Code    int                   `json:"code" example:"200"`
	Data    JoinGroupResponseData `json:"data"`
	Error   *GroupErrorResponse   `json:"error,omitempty"`
}

type JoinGroupResponseData struct {
	SessionID string     `json:"session_id" example:"default"`
	Action    string     `json:"action" example:"join"`
	Status    string     `json:"status" example:"success"`
	Timestamp time.Time  `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group     *GroupInfo `json:"group"`
}

type GetInviteLinkResponse struct {
	Success bool                      `json:"success" example:"true"`
	Code    int                       `json:"code" example:"200"`
	Data    GetInviteLinkResponseData `json:"data"`
	Error   *GroupErrorResponse       `json:"error,omitempty"`
}

type GetInviteLinkResponseData struct {
	SessionID  string    `json:"session_id" example:"default"`
	Action     string    `json:"action" example:"invite_link"`
	Status     string    `json:"status" example:"success"`
	Timestamp  time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	GroupJID   string    `json:"group_jid" example:"120363025246125486@g.us"`
	InviteLink string    `json:"invite_link" example:"https://chat.meow.com/ABC123"`
}

func NewGroupSuccessResponse(sessionID, action string, group *GroupInfo) *GroupResponse {
	return &GroupResponse{
		Success: true,
		Code:    200,
		Data: GroupData{
			SessionID: sessionID,
			Action:    action,
			Status:    "success",
			Timestamp: time.Now(),
			Group:     group,
		},
	}
}

func NewGroupListResponse(sessionID string, groups []GroupInfo) *GroupResponse {
	return &GroupResponse{
		Success: true,
		Code:    200,
		Data: GroupData{
			SessionID: sessionID,
			Action:    "list",
			Status:    "success",
			Timestamp: time.Now(),
			Groups:    groups,
		},
	}
}

func NewGroupErrorResponse(code int, errorCode, message, details string) *GroupResponse {
	return &GroupResponse{
		Success: false,
		Code:    code,
		Data: GroupData{
			Status:    "error",
			Timestamp: time.Now(),
		},
		Error: &GroupErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
}

func NewInviteLinkResponse(sessionID, groupJID, inviteLink string) *GroupResponse {
	return &GroupResponse{
		Success: true,
		Code:    200,
		Data: GroupData{
			SessionID:  sessionID,
			Action:     "invite_link",
			Status:     "success",
			Timestamp:  time.Now(),
			InviteLink: inviteLink,
		},
	}
}

func NewGroupOperationResponse(sessionID, action, message string) *GroupResponse {
	return &GroupResponse{
		Success: true,
		Code:    200,
		Data: GroupData{
			SessionID: sessionID,
			Action:    action,
			Status:    "success",
			Timestamp: time.Now(),
			Message:   message,
		},
	}
}

func (r *CreateGroupRequest) Validate() error {
	if r.Name == "" {
		return &GroupValidationError{Field: "name", Message: "Group name is required"}
	}
	if len(r.Participants) == 0 {
		return &GroupValidationError{Field: "participants", Message: "At least one participant is required"}
	}
	for _, participant := range r.Participants {
		if !validatePhone(participant) {
			return &GroupValidationError{Field: "participants", Message: "Invalid phone number format: " + participant}
		}
	}
	return nil
}

func (r *UpdateParticipantsRequest) Validate() error {
	if r.GroupJID == "" {
		return &GroupValidationError{Field: "group_jid", Message: "Group JID is required"}
	}
	if !validateGroupJID(r.GroupJID) {
		return &GroupValidationError{Field: "group_jid", Message: "Invalid group JID format"}
	}
	validActions := []string{"add", "remove", "promote", "demote"}
	isValidAction := false
	for _, action := range validActions {
		if r.Action == action {
			isValidAction = true
			break
		}
	}
	if !isValidAction {
		return &GroupValidationError{Field: "action", Message: "Invalid action. Must be: add, remove, promote, or demote"}
	}
	if len(r.Participants) == 0 {
		return &GroupValidationError{Field: "participants", Message: "At least one participant is required"}
	}
	for _, participant := range r.Participants {
		if !validatePhone(participant) {
			return &GroupValidationError{Field: "participants", Message: "Invalid phone number format: " + participant}
		}
	}
	return nil
}

type GroupValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *GroupValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func validatePhone(phone string) bool {
	if phone == "" {
		return false
	}
	return len(phone) >= 10 && len(phone) <= 15
}

func validateGroupJID(groupJID string) bool {
	if groupJID == "" {
		return false
	}
	return len(groupJID) > 5 && groupJID[len(groupJID)-5:] == "@g.us"
}

func (r *GetGroupRequestsReq) Validate() error {
	if r.GroupJID == "" {
		return &GroupValidationError{Field: "group_jid", Message: "Group JID is required"}
	}
	if !validateGroupJID(r.GroupJID) {
		return &GroupValidationError{Field: "group_jid", Message: "Invalid group JID format"}
	}
	return nil
}

func (r *UpdateGroupRequestsReq) Validate() error {
	if r.GroupJID == "" {
		return &GroupValidationError{Field: "group_jid", Message: "Group JID is required"}
	}
	if !validateGroupJID(r.GroupJID) {
		return &GroupValidationError{Field: "group_jid", Message: "Invalid group JID format"}
	}
	validActions := []string{"approve", "reject"}
	isValidAction := false
	for _, action := range validActions {
		if r.Action == action {
			isValidAction = true
			break
		}
	}
	if !isValidAction {
		return &GroupValidationError{Field: "action", Message: "Invalid action. Must be: approve or reject"}
	}
	if len(r.Participants) == 0 {
		return &GroupValidationError{Field: "participants", Message: "At least one participant is required"}
	}
	for _, participant := range r.Participants {
		if !validatePhone(participant) {
			return &GroupValidationError{Field: "participants", Message: "Invalid phone number format: " + participant}
		}
	}
	return nil
}
