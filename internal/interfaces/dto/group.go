package dto

import (
	"time"
)

// ============================================================================
// GROUP REQUEST DTOs
// ============================================================================

// CreateGroupRequest represents the request to create a new group
type CreateGroupRequest struct {
	Name         string   `json:"name" binding:"required" example:"My Group"`
	Participants []string `json:"participants" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

// GetGroupInfoRequest represents the request to get group information
type GetGroupInfoRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

// JoinGroupRequest represents the request to join a group via invite link
type JoinGroupRequest struct {
	InviteLink string `json:"invite_link" binding:"required" example:"https://chat.whatsapp.com/ABC123"`
}

// JoinGroupWithInviteRequest represents the request to join a group via specific invite
type JoinGroupWithInviteRequest struct {
	GroupJID   string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Inviter    string `json:"inviter" binding:"required" example:"5511999999999@s.whatsapp.net"`
	Code       string `json:"code" binding:"required" example:"ABC123DEF456"`
	Expiration int64  `json:"expiration" binding:"required" example:"1640995200"`
}

// LeaveGroupRequest represents the request to leave a group
type LeaveGroupRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

// GetInviteLinkRequest represents the request to get group invite link
type GetInviteLinkRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Reset    bool   `json:"reset,omitempty" example:"false"`
}

// GetInviteInfoRequest represents the request to get invite info from link
type GetInviteInfoRequest struct {
	InviteLink string `json:"invite_link" binding:"required" example:"https://chat.whatsapp.com/ABC123"`
}

// GetGroupInfoFromInviteRequest represents the request to get group info from specific invite
type GetGroupInfoFromInviteRequest struct {
	GroupJID   string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Inviter    string `json:"inviter" binding:"required" example:"5511999999999@s.whatsapp.net"`
	Code       string `json:"code" binding:"required" example:"ABC123DEF456"`
	Expiration int64  `json:"expiration" binding:"required" example:"1640995200"`
}

// UpdateParticipantsRequest represents the request to update group participants
type UpdateParticipantsRequest struct {
	GroupJID     string   `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Action       string   `json:"action" binding:"required" example:"add"`
	Participants []string `json:"participants" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

// SetGroupNameRequest represents the request to set group name
type SetGroupNameRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Name     string `json:"name" binding:"required" example:"New Group Name"`
}

// SetGroupTopicRequest represents the request to set group topic
type SetGroupTopicRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Topic    string `json:"topic" binding:"required" example:"Group topic description"`
}

// SetGroupPhotoRequest represents the request to set group photo
type SetGroupPhotoRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Photo    string `json:"photo" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD..."`
}

// RemoveGroupPhotoRequest represents the request to remove group photo
type RemoveGroupPhotoRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

// SetGroupAnnounceRequest represents the request to set group announce setting
type SetGroupAnnounceRequest struct {
	GroupJID     string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	AnnounceOnly bool   `json:"announce_only" example:"true"`
}

// SetGroupLockedRequest represents the request to set group locked setting
type SetGroupLockedRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Locked   bool   `json:"locked" example:"true"`
}

// SetGroupEphemeralRequest represents the request to set group ephemeral setting
type SetGroupEphemeralRequest struct {
	GroupJID  string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Ephemeral bool   `json:"ephemeral" example:"true"`
	Duration  int    `json:"duration,omitempty" example:"86400"`
}

// SetGroupJoinApprovalRequest represents the request to set group join approval mode
type SetGroupJoinApprovalRequest struct {
	GroupJID        string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	RequireApproval bool   `json:"require_approval" binding:"required" example:"true"`
}

// SetGroupMemberAddModeRequest represents the request to set group member add mode
type SetGroupMemberAddModeRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Mode     string `json:"mode" binding:"required" example:"admin" enum:"all,admin"`
}

// GetGroupRequestParticipantsRequest represents the request to get group request participants
type GetGroupRequestParticipantsRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
}

// UpdateGroupRequestParticipantsRequest represents the request to update group request participants
type UpdateGroupRequestParticipantsRequest struct {
	GroupJID     string   `json:"group_jid" binding:"required" example:"120363025246125486@g.us"`
	Action       string   `json:"action" binding:"required" example:"approve"`
	Participants []string `json:"participants" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

// LinkGroupRequest represents the request to link a group to a community
type LinkGroupRequest struct {
	CommunityJID string `json:"community_jid" binding:"required" example:"120363025246125486@g.us"`
	GroupJID     string `json:"group_jid" binding:"required" example:"120363025246125487@g.us"`
}

// UnlinkGroupRequest represents the request to unlink a group from a community
type UnlinkGroupRequest struct {
	CommunityJID string `json:"community_jid" binding:"required" example:"120363025246125486@g.us"`
	GroupJID     string `json:"group_jid" binding:"required" example:"120363025246125487@g.us"`
}

// GetSubGroupsRequest represents the request to get subgroups of a community
type GetSubGroupsRequest struct {
	CommunityJID string `json:"community_jid" binding:"required" example:"120363025246125486@g.us"`
}

// GetLinkedGroupsParticipantsRequest represents the request to get participants of linked groups
type GetLinkedGroupsParticipantsRequest struct {
	CommunityJID string `json:"community_jid" binding:"required" example:"120363025246125486@g.us"`
}

// ============================================================================
// GROUP DATA STRUCTURES
// ============================================================================

// GroupInfo represents group information
type GroupInfo struct {
	JID          string   `json:"jid" example:"120363025246125486@g.us"`
	Name         string   `json:"name" example:"My Group"`
	Topic        string   `json:"topic,omitempty" example:"Group topic description"`
	Participants []string `json:"participants" example:"[\"5511999999999@s.whatsapp.net\", \"5511888888888@s.whatsapp.net\"]"`
	Admins       []string `json:"admins" example:"[\"5511999999999@s.whatsapp.net\"]"`
	Owner        string   `json:"owner" example:"5511999999999@s.whatsapp.net"`
	CreatedAt    int64    `json:"created_at" example:"1640995200"`
	Size         int      `json:"size" example:"2"`
	Announce     bool     `json:"announce" example:"false"`
	Locked       bool     `json:"locked" example:"false"`
	Ephemeral    bool     `json:"ephemeral" example:"false"`
}

// ============================================================================
// GROUP RESPONSE DTOs
// ============================================================================

// GroupResponse represents the standardized response format for group operations
type GroupResponse struct {
	Success bool                 `json:"success"`
	Code    int                  `json:"code"`
	Data    GroupData            `json:"data"`
	Error   *GroupErrorResponse  `json:"error,omitempty"`
}

// GroupData contains the response data for group operations
type GroupData struct {
	SessionID string     `json:"session_id,omitempty" example:"default"`
	Action    string     `json:"action" example:"create"`
	Status    string     `json:"status" example:"success"`
	Timestamp time.Time  `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group     *GroupInfo `json:"group,omitempty"`
	Groups    []GroupInfo `json:"groups,omitempty"`
	InviteLink string    `json:"invite_link,omitempty" example:"https://chat.whatsapp.com/ABC123"`
	Message   string     `json:"message,omitempty" example:"Operation completed successfully"`
}

// GroupErrorResponse represents error information for group operations
type GroupErrorResponse struct {
	Code    string `json:"code" example:"INVALID_GROUP_JID"`
	Message string `json:"message" example:"Invalid group JID format"`
	Details string `json:"details,omitempty" example:"Group JID must end with @g.us"`
}

// ============================================================================
// SPECIFIC RESPONSE DTOs FOR SWAGGER DOCUMENTATION
// ============================================================================

// CreateGroupResponse represents the response for group creation
type CreateGroupResponse struct {
	Success bool                       `json:"success" example:"true"`
	Code    int                        `json:"code" example:"201"`
	Data    CreateGroupResponseData    `json:"data"`
	Error   *GroupErrorResponse        `json:"error,omitempty"`
}

// CreateGroupResponseData contains the data for group creation response
type CreateGroupResponseData struct {
	SessionID string     `json:"session_id" example:"default"`
	Action    string     `json:"action" example:"create"`
	Status    string     `json:"status" example:"success"`
	Timestamp time.Time  `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group     *GroupInfo `json:"group"`
}

// ListGroupsResponse represents the response for listing groups
type ListGroupsResponse struct {
	Success bool                      `json:"success" example:"true"`
	Code    int                       `json:"code" example:"200"`
	Data    ListGroupsResponseData    `json:"data"`
	Error   *GroupErrorResponse       `json:"error,omitempty"`
}

// ListGroupsResponseData contains the data for group list response
type ListGroupsResponseData struct {
	SessionID string      `json:"session_id" example:"default"`
	Action    string      `json:"action" example:"list"`
	Status    string      `json:"status" example:"success"`
	Timestamp time.Time   `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Groups    []GroupInfo `json:"groups"`
	Total     int         `json:"total" example:"5"`
}

// GetGroupInfoResponse represents the response for getting group information
type GetGroupInfoResponse struct {
	Success bool                        `json:"success" example:"true"`
	Code    int                         `json:"code" example:"200"`
	Data    GetGroupInfoResponseData    `json:"data"`
	Error   *GroupErrorResponse         `json:"error,omitempty"`
}

// GetGroupInfoResponseData contains the data for group info response
type GetGroupInfoResponseData struct {
	SessionID string     `json:"session_id" example:"default"`
	Action    string     `json:"action" example:"info"`
	Status    string     `json:"status" example:"success"`
	Timestamp time.Time  `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group     *GroupInfo `json:"group"`
}

// JoinGroupResponse represents the response for joining a group
type JoinGroupResponse struct {
	Success bool                     `json:"success" example:"true"`
	Code    int                      `json:"code" example:"200"`
	Data    JoinGroupResponseData    `json:"data"`
	Error   *GroupErrorResponse      `json:"error,omitempty"`
}

// JoinGroupResponseData contains the data for join group response
type JoinGroupResponseData struct {
	SessionID string     `json:"session_id" example:"default"`
	Action    string     `json:"action" example:"join"`
	Status    string     `json:"status" example:"success"`
	Timestamp time.Time  `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Group     *GroupInfo `json:"group"`
}

// GetInviteLinkResponse represents the response for getting invite link
type GetInviteLinkResponse struct {
	Success bool                         `json:"success" example:"true"`
	Code    int                          `json:"code" example:"200"`
	Data    GetInviteLinkResponseData    `json:"data"`
	Error   *GroupErrorResponse          `json:"error,omitempty"`
}

// GetInviteLinkResponseData contains the data for invite link response
type GetInviteLinkResponseData struct {
	SessionID  string    `json:"session_id" example:"default"`
	Action     string    `json:"action" example:"invite_link"`
	Status     string    `json:"status" example:"success"`
	Timestamp  time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	GroupJID   string    `json:"group_jid" example:"120363025246125486@g.us"`
	InviteLink string    `json:"invite_link" example:"https://chat.whatsapp.com/ABC123"`
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// NewGroupSuccessResponse creates a successful group operation response
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

// NewGroupListResponse creates a successful group list response
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

// NewGroupErrorResponse creates an error response for group operations
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

// NewInviteLinkResponse creates a successful invite link response
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

// NewGroupOperationResponse creates a response for simple group operations
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

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

// Validate validates a CreateGroupRequest
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

// Validate validates an UpdateParticipantsRequest
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

// GroupValidationError represents a validation error for group requests
type GroupValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *GroupValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// validatePhone checks if a phone number is valid
func validatePhone(phone string) bool {
	if phone == "" {
		return false
	}
	// Basic validation - should start with country code
	return len(phone) >= 10 && len(phone) <= 15
}

// validateGroupJID checks if a group JID is valid
func validateGroupJID(groupJID string) bool {
	if groupJID == "" {
		return false
	}
	// Group JIDs should end with @g.us
	return len(groupJID) > 5 && groupJID[len(groupJID)-5:] == "@g.us"
}

// Validate validates a GetGroupRequestParticipantsRequest
func (r *GetGroupRequestParticipantsRequest) Validate() error {
	if r.GroupJID == "" {
		return &GroupValidationError{Field: "group_jid", Message: "Group JID is required"}
	}
	if !validateGroupJID(r.GroupJID) {
		return &GroupValidationError{Field: "group_jid", Message: "Invalid group JID format"}
	}
	return nil
}

// Validate validates an UpdateGroupRequestParticipantsRequest
func (r *UpdateGroupRequestParticipantsRequest) Validate() error {
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

// Validate validates a LinkGroupRequest
func (r *LinkGroupRequest) Validate() error {
	if r.CommunityJID == "" {
		return &GroupValidationError{Field: "community_jid", Message: "Community JID is required"}
	}
	if !validateGroupJID(r.CommunityJID) {
		return &GroupValidationError{Field: "community_jid", Message: "Invalid community JID format"}
	}
	if r.GroupJID == "" {
		return &GroupValidationError{Field: "group_jid", Message: "Group JID is required"}
	}
	if !validateGroupJID(r.GroupJID) {
		return &GroupValidationError{Field: "group_jid", Message: "Invalid group JID format"}
	}
	return nil
}

// Validate validates an UnlinkGroupRequest
func (r *UnlinkGroupRequest) Validate() error {
	if r.CommunityJID == "" {
		return &GroupValidationError{Field: "community_jid", Message: "Community JID is required"}
	}
	if !validateGroupJID(r.CommunityJID) {
		return &GroupValidationError{Field: "community_jid", Message: "Invalid community JID format"}
	}
	if r.GroupJID == "" {
		return &GroupValidationError{Field: "group_jid", Message: "Group JID is required"}
	}
	if !validateGroupJID(r.GroupJID) {
		return &GroupValidationError{Field: "group_jid", Message: "Invalid group JID format"}
	}
	return nil
}

// Validate validates a GetSubGroupsRequest
func (r *GetSubGroupsRequest) Validate() error {
	if r.CommunityJID == "" {
		return &GroupValidationError{Field: "community_jid", Message: "Community JID is required"}
	}
	if !validateGroupJID(r.CommunityJID) {
		return &GroupValidationError{Field: "community_jid", Message: "Invalid community JID format"}
	}
	return nil
}

// Validate validates a GetLinkedGroupsParticipantsRequest
func (r *GetLinkedGroupsParticipantsRequest) Validate() error {
	if r.CommunityJID == "" {
		return &GroupValidationError{Field: "community_jid", Message: "Community JID is required"}
	}
	if !validateGroupJID(r.CommunityJID) {
		return &GroupValidationError{Field: "community_jid", Message: "Invalid community JID format"}
	}
	return nil
}

// Newsletter/Channel DTOs

// CreateNewsletterRequest represents the request to create a newsletter/channel
type CreateNewsletterRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100" example:"My Channel"`
	Description string `json:"description" binding:"max=500" example:"Channel description"`
	Picture     string `json:"picture,omitempty" example:"https://example.com/image.jpg"`
}

// CreateNewsletterResponse represents the response for creating a newsletter
type CreateNewsletterResponse struct {
	JID         string `json:"jid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

// GetNewsletterInfoRequest represents the request to get newsletter info
type GetNewsletterInfoRequest struct {
	NewsletterJID string `json:"newsletter_jid" binding:"required" example:"120363197100429171@newsletter"`
}

// NewsletterInfo represents newsletter information
type NewsletterInfo struct {
	JID           string `json:"jid"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Picture       string `json:"picture,omitempty"`
	FollowerCount int64  `json:"follower_count"`
	CreatedAt     int64  `json:"created_at"`
	Verified      bool   `json:"verified"`
	Muted         bool   `json:"muted"`
}

// GetNewsletterInfoResponse represents the response for getting newsletter info
type GetNewsletterInfoResponse struct {
	Newsletter NewsletterInfo `json:"newsletter"`
}

// FollowNewsletterRequest represents the request to follow a newsletter
type FollowNewsletterRequest struct {
	NewsletterJID string `json:"newsletter_jid" binding:"required" example:"120363197100429171@newsletter"`
}

// UnfollowNewsletterRequest represents the request to unfollow a newsletter
type UnfollowNewsletterRequest struct {
	NewsletterJID string `json:"newsletter_jid" binding:"required" example:"120363197100429171@newsletter"`
}

// GetSubscribedNewslettersResponse represents the response for getting subscribed newsletters
type GetSubscribedNewslettersResponse struct {
	Newsletters []NewsletterInfo `json:"newsletters"`
}

// NewsletterReactionRequest represents the request to react to a newsletter message
type NewsletterReactionRequest struct {
	NewsletterJID string `json:"newsletter_jid" binding:"required" example:"120363197100429171@newsletter"`
	ServerID      string `json:"server_id" binding:"required" example:"msg_server_id"`
	Reaction      string `json:"reaction" binding:"required" example:"👍"`
	MessageID     string `json:"message_id,omitempty" example:"3EB0123456789ABCDEF"`
}

// NewsletterMarkViewedRequest represents the request to mark newsletter messages as viewed
type NewsletterMarkViewedRequest struct {
	NewsletterJID string   `json:"newsletter_jid" binding:"required" example:"120363197100429171@newsletter"`
	ServerIDs     []string `json:"server_ids" binding:"required" example:"[\"msg_id_1\",\"msg_id_2\"]"`
}

// NewsletterToggleMuteRequest represents the request to toggle newsletter mute
type NewsletterToggleMuteRequest struct {
	NewsletterJID string `json:"newsletter_jid" binding:"required" example:"120363197100429171@newsletter"`
	Mute          bool   `json:"mute" example:"true"`
}
