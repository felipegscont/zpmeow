package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"zpmeow/internal/application"
	"zpmeow/internal/infrastructure/wameow"
	"zpmeow/internal/interfaces/dto"
	"zpmeow/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

// GroupHandler handles group-related HTTP requests
type GroupHandler struct {
	sessionService *application.SessionService
	wameowService  *wameow.MeowService
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(sessionService *application.SessionService, wameowService *wameow.MeowService) *GroupHandler {
	return &GroupHandler{
		sessionService: sessionService,
		wameowService:  wameowService,
	}
}

// resolveSessionID resolves session ID or name to actual session ID
func (h *GroupHandler) resolveSessionID(c *gin.Context, sessionIDOrName string) (string, error) {
	if h.sessionService == nil {
		// Fallback: assume it's already an ID
		return sessionIDOrName, nil
	}

	ctx := context.Background()
	session, err := h.sessionService.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return "", err
	}

	return session.ID, nil
}

// CreateGroup handles creating a group
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Create group using wameow service
	ctx := context.Background()
	groupInfo, err := h.wameowService.CreateGroup(ctx, sessionID, req.Name, req.Participants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"CREATE_GROUP_FAILED",
			"Failed to create group",
			err.Error(),
		))
		return
	}

	// Convert wameow.GroupInfo to dto.GroupInfo
	dtoGroupInfo := convertWameowGroupInfoToDTO(groupInfo)

	response := dto.NewGroupSuccessResponse(sessionID, "create", dtoGroupInfo)
	response.Code = http.StatusCreated
	c.JSON(http.StatusCreated, response)
}

// GetGroupInfo handles getting group information
func (h *GroupHandler) GetGroupInfo(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetGroupInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Get group info using wameow service
	ctx := context.Background()
	groupInfo, err := h.wameowService.GetGroupInfo(ctx, sessionID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"GET_GROUP_INFO_FAILED",
			"Failed to get group information",
			err.Error(),
		))
		return
	}

	// Convert wameow.GroupInfo to dto.GroupInfo
	dtoGroupInfo := convertWameowGroupInfoToDTO(groupInfo)

	response := dto.NewGroupSuccessResponse(sessionID, "info", dtoGroupInfo)
	c.JSON(http.StatusOK, response)
}

// ListGroups handles listing groups
func (h *GroupHandler) ListGroups(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	// List groups using wameow service
	ctx := context.Background()
	groups, err := h.wameowService.ListGroups(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"LIST_GROUPS_FAILED",
			"Failed to list groups",
			err.Error(),
		))
		return
	}

	// Convert wameow.GroupInfo slice to dto.GroupInfo slice
	dtoGroups := convertWameowGroupInfoSliceToDTO(groups)

	response := dto.NewGroupListResponse(sessionID, dtoGroups)
	c.JSON(http.StatusOK, response)
}

// JoinGroup handles joining a group via invite link
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.JoinGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Join group using wameow service
	ctx := context.Background()
	groupInfo, err := h.wameowService.JoinGroup(ctx, sessionID, req.InviteLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"JOIN_GROUP_FAILED",
			"Failed to join group",
			err.Error(),
		))
		return
	}

	// Convert wameow.GroupInfo to dto.GroupInfo
	dtoGroupInfo := convertWameowGroupInfoToDTO(groupInfo)

	response := dto.NewGroupSuccessResponse(sessionID, "join", dtoGroupInfo)
	c.JSON(http.StatusOK, response)
}

// JoinGroupWithInvite handles joining a group with specific invite
func (h *GroupHandler) JoinGroupWithInvite(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.JoinGroupWithInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Join group using wameow service
	ctx := context.Background()
	groupInfo, err := h.wameowService.JoinGroupWithInvite(ctx, sessionID, req.GroupJID, req.Inviter, req.Code, req.Expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"JOIN_GROUP_WITH_INVITE_FAILED",
			"Failed to join group with invite",
			err.Error(),
		))
		return
	}

	// Convert wameow.GroupInfo to dto.GroupInfo
	dtoGroupInfo := convertWameowGroupInfoToDTO(groupInfo)

	response := dto.NewGroupSuccessResponse(sessionID, "join_with_invite", dtoGroupInfo)
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.LeaveGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Leave group using wameow service
	ctx := context.Background()
	err = h.wameowService.LeaveGroup(ctx, sessionID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"LEAVE_GROUP_FAILED",
			"Failed to leave group",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "leave", "Successfully left the group")
	c.JSON(http.StatusOK, response)
}
func (h *GroupHandler) GetInviteLink(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetInviteLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Get invite link using wameow service
	ctx := context.Background()
	inviteLink, err := h.wameowService.GetInviteLink(ctx, sessionID, req.GroupJID, req.Reset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"GET_INVITE_LINK_FAILED",
			"Failed to get invite link",
			err.Error(),
		))
		return
	}

	response := dto.NewInviteLinkResponse(sessionID, req.GroupJID, inviteLink)
	c.JSON(http.StatusOK, response)
}

// Note: RevokeInviteLink and UpdateGroupSettings are not implemented yet
// These can be added later if needed
// Note: AddParticipant, RemoveParticipant, PromoteParticipant, and DemoteParticipant
// are now handled by the UpdateParticipants method with different action parameters

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// convertWameowGroupInfoToDTO converts wameow.GroupInfo to dto.GroupInfo
func convertWameowGroupInfoToDTO(groupInfo *wameow.GroupInfo) *dto.GroupInfo {
	if groupInfo == nil {
		return nil
	}

	return &dto.GroupInfo{
		JID:          groupInfo.JID,
		Name:         groupInfo.Name,
		Topic:        groupInfo.Topic,
		Participants: groupInfo.Participants,
		Admins:       groupInfo.Admins,
		Owner:        groupInfo.Owner,
		CreatedAt:    groupInfo.CreatedAt,
		Size:         groupInfo.Size,
		Announce:     groupInfo.Announce,
		Locked:       groupInfo.Locked,
		Ephemeral:    groupInfo.Ephemeral,
	}
}

// convertWameowGroupInfoSliceToDTO converts slice of wameow.GroupInfo to slice of dto.GroupInfo
func convertWameowGroupInfoSliceToDTO(groups []wameow.GroupInfo) []dto.GroupInfo {
	var dtoGroups []dto.GroupInfo
	for _, group := range groups {
		dtoGroups = append(dtoGroups, dto.GroupInfo{
			JID:          group.JID,
			Name:         group.Name,
			Topic:        group.Topic,
			Participants: group.Participants,
			Admins:       group.Admins,
			Owner:        group.Owner,
			CreatedAt:    group.CreatedAt,
			Size:         group.Size,
			Announce:     group.Announce,
			Locked:       group.Locked,
			Ephemeral:    group.Ephemeral,
		})
	}
	return dtoGroups
}
func (h *GroupHandler) GetInviteInfo(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetInviteInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Get invite info using wameow service
	ctx := context.Background()
	groupInfo, err := h.wameowService.GetInviteInfo(ctx, sessionID, req.InviteLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"GET_INVITE_INFO_FAILED",
			"Failed to get invite info",
			err.Error(),
		))
		return
	}

	// Convert wameow.GroupInfo to dto.GroupInfo
	dtoGroupInfo := convertWameowGroupInfoToDTO(groupInfo)

	response := dto.NewGroupSuccessResponse(sessionID, "invite_info", dtoGroupInfo)
	c.JSON(http.StatusOK, response)
}

// GetGroupInfoFromInvite handles getting group info from specific invite
func (h *GroupHandler) GetGroupInfoFromInvite(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetGroupInfoFromInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Get group info from invite using wameow service
	ctx := context.Background()
	groupInfo, err := h.wameowService.GetGroupInfoFromInvite(ctx, sessionID, req.GroupJID, req.Inviter, req.Code, req.Expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"GET_GROUP_INFO_FROM_INVITE_FAILED",
			"Failed to get group info from invite",
			err.Error(),
		))
		return
	}

	// Convert wameow.GroupInfo to dto.GroupInfo
	dtoGroupInfo := convertWameowGroupInfoToDTO(groupInfo)

	response := dto.NewGroupSuccessResponse(sessionID, "invite_info_specific", dtoGroupInfo)
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) UpdateParticipants(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.UpdateParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Update participants using wameow service
	ctx := context.Background()
	err = h.wameowService.UpdateParticipants(ctx, sessionID, req.GroupJID, req.Action, req.Participants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"UPDATE_PARTICIPANTS_FAILED",
			"Failed to update participants",
			err.Error(),
		))
		return
	}

	message := "Successfully " + req.Action + "ed participants"
	response := dto.NewGroupOperationResponse(sessionID, "update_participants", message)
	c.JSON(http.StatusOK, response)
}
func (h *GroupHandler) SetName(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Set group name using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupName(ctx, sessionID, req.GroupJID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_GROUP_NAME_FAILED",
			"Failed to set group name",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_name", "Group name updated successfully")
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) SetTopic(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Set group topic using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupTopic(ctx, sessionID, req.GroupJID, req.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_GROUP_TOPIC_FAILED",
			"Failed to set group topic",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_topic", "Group topic updated successfully")
	c.JSON(http.StatusOK, response)
}
func (h *GroupHandler) SetPhoto(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Process photo data (supports URLs, base64, data URLs)
	photoData, _, err := utils.ProcessUnifiedMedia(context.Background(), req.Photo, nil, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_PHOTO",
			"Invalid photo data",
			err.Error(),
		))
		return
	}

	// Set group photo using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupPhoto(ctx, sessionID, req.GroupJID, photoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_GROUP_PHOTO_FAILED",
			"Failed to set group photo",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_photo", "Group photo updated successfully")
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) RemovePhoto(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.RemoveGroupPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Remove group photo using wameow service
	ctx := context.Background()
	err = h.wameowService.RemoveGroupPhoto(ctx, sessionID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"REMOVE_GROUP_PHOTO_FAILED",
			"Failed to remove group photo",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "remove_photo", "Group photo removed successfully")
	c.JSON(http.StatusOK, response)
}
func (h *GroupHandler) SetAnnounce(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupAnnounceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Set group announce using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupAnnounce(ctx, sessionID, req.GroupJID, req.AnnounceOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_GROUP_ANNOUNCE_FAILED",
			"Failed to set group announce setting",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_announce", "Group announce setting updated successfully")
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) SetLocked(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupLockedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Set group locked using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupLocked(ctx, sessionID, req.GroupJID, req.Locked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_GROUP_LOCKED_FAILED",
			"Failed to set group locked setting",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_locked", "Group locked setting updated successfully")
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) SetEphemeral(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupEphemeralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Set group ephemeral using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupEphemeral(ctx, sessionID, req.GroupJID, req.Ephemeral, req.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_GROUP_EPHEMERAL_FAILED",
			"Failed to set group ephemeral setting",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_ephemeral", "Group ephemeral setting updated successfully")
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) SetJoinApproval(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupJoinApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Set group join approval mode using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupJoinApprovalMode(ctx, sessionID, req.GroupJID, req.RequireApproval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_JOIN_APPROVAL_FAILED",
			"Failed to set group join approval mode",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_join_approval", "Group join approval mode updated successfully")
	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) SetMemberAddMode(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SetGroupMemberAddModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Set group member add mode using wameow service
	ctx := context.Background()
	err = h.wameowService.SetGroupMemberAddMode(ctx, sessionID, req.GroupJID, req.Mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"SET_MEMBER_ADD_MODE_FAILED",
			"Failed to set group member add mode",
			err.Error(),
		))
		return
	}

	response := dto.NewGroupOperationResponse(sessionID, "set_member_add_mode", "Group member add mode updated successfully")
	c.JSON(http.StatusOK, response)
}

// GetGroupRequestParticipants handles getting group request participants
func (h *GroupHandler) GetGroupRequestParticipants(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetGroupRequestParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Get group request participants using wameow service
	ctx := context.Background()
	participants, err := h.wameowService.GetGroupRequestParticipants(ctx, sessionID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"GET_GROUP_REQUEST_PARTICIPANTS_FAILED",
			"Failed to get group request participants",
			err.Error(),
		))
		return
	}

	// Create response with participants list
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": gin.H{
			"session_id":   sessionID,
			"action":       "get_request_participants",
			"status":       "success",
			"timestamp":    time.Now(),
			"group_jid":    req.GroupJID,
			"participants": participants,
			"total":        len(participants),
		},
	})
}

// UpdateGroupRequestParticipants handles updating group request participants
func (h *GroupHandler) UpdateGroupRequestParticipants(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.UpdateGroupRequestParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Update group request participants using wameow service
	ctx := context.Background()
	err = h.wameowService.UpdateGroupRequestParticipants(ctx, sessionID, req.GroupJID, req.Action, req.Participants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"UPDATE_GROUP_REQUEST_PARTICIPANTS_FAILED",
			"Failed to update group request participants",
			err.Error(),
		))
		return
	}

	message := fmt.Sprintf("Successfully %sed %d participants", req.Action, len(req.Participants))
	response := dto.NewGroupOperationResponse(sessionID, "update_request_participants", message)
	c.JSON(http.StatusOK, response)
}

// LinkGroup handles linking a group to a community
func (h *GroupHandler) LinkGroup(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.LinkGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Link group using wameow service
	ctx := context.Background()
	err = h.wameowService.LinkGroup(ctx, sessionID, req.CommunityJID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"LINK_GROUP_FAILED",
			"Failed to link group to community",
			err.Error(),
		))
		return
	}

	message := fmt.Sprintf("Successfully linked group %s to community %s", req.GroupJID, req.CommunityJID)
	response := dto.NewGroupOperationResponse(sessionID, "link_group", message)
	c.JSON(http.StatusOK, response)
}

// UnlinkGroup handles unlinking a group from a community
func (h *GroupHandler) UnlinkGroup(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.UnlinkGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Unlink group using wameow service
	ctx := context.Background()
	err = h.wameowService.UnlinkGroup(ctx, sessionID, req.CommunityJID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"UNLINK_GROUP_FAILED",
			"Failed to unlink group from community",
			err.Error(),
		))
		return
	}

	message := fmt.Sprintf("Successfully unlinked group %s from community %s", req.GroupJID, req.CommunityJID)
	response := dto.NewGroupOperationResponse(sessionID, "unlink_group", message)
	c.JSON(http.StatusOK, response)
}

// GetSubGroups handles getting subgroups of a community
func (h *GroupHandler) GetSubGroups(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetSubGroupsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Get subgroups using wameow service
	ctx := context.Background()
	subGroups, err := h.wameowService.GetSubGroups(ctx, sessionID, req.CommunityJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"GET_SUBGROUPS_FAILED",
			"Failed to get subgroups",
			err.Error(),
		))
		return
	}

	// Create response with subgroups list
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": gin.H{
			"session_id":    sessionID,
			"action":        "get_subgroups",
			"status":        "success",
			"timestamp":     time.Now(),
			"community_jid": req.CommunityJID,
			"subgroups":     subGroups,
			"total":         len(subGroups),
		},
	})
}

// GetLinkedGroupsParticipants handles getting participants of linked groups
func (h *GroupHandler) GetLinkedGroupsParticipants(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewGroupErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetLinkedGroupsParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGroupErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Get linked groups participants using wameow service
	ctx := context.Background()
	participants, err := h.wameowService.GetLinkedGroupsParticipants(ctx, sessionID, req.CommunityJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGroupErrorResponse(
			http.StatusInternalServerError,
			"GET_LINKED_GROUPS_PARTICIPANTS_FAILED",
			"Failed to get linked groups participants",
			err.Error(),
		))
		return
	}

	// Create response with participants list
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": gin.H{
			"session_id":    sessionID,
			"action":        "get_linked_groups_participants",
			"status":        "success",
			"timestamp":     time.Now(),
			"community_jid": req.CommunityJID,
			"participants":  participants,
			"total":         len(participants),
		},
	})
}

// ============================================================================
// HELPER FUNCTIONS FOR MEDIA HANDLING
// ============================================================================

// decodeMediaData decodes base64 data URL to bytes or downloads from URL
func (h *GroupHandler) decodeMediaData(dataURL string) ([]byte, error) {
	// Check if it's a HTTP/HTTPS URL
	if strings.HasPrefix(dataURL, "http://") || strings.HasPrefix(dataURL, "https://") {
		// Download from URL
		resp, err := http.Get(dataURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download from URL: HTTP %d", resp.StatusCode)
		}

		// Read the response body
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return data, nil
	}

	// Handle data URL format: data:image/jpeg;base64,<base64-data>
	if strings.HasPrefix(dataURL, "data:") {
		// Find the comma that separates the header from the data
		commaIndex := strings.Index(dataURL, ",")
		if commaIndex == -1 {
			return nil, fmt.Errorf("invalid data URL format")
		}

		// Extract the base64 data part
		base64Data := dataURL[commaIndex+1:]

		// Decode base64
		data, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 data: %w", err)
		}

		return data, nil
	}

	// Assume it's raw base64 data
	data, err := base64.StdEncoding.DecodeString(dataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return data, nil
}
