package handler

import (
	"net/http"
	"strings"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/infra/meow"
	"zpmeow/internal/types"
	"zpmeow/internal/utils"

	"github.com/gin-gonic/gin"
)

// GroupHandler handles HTTP requests for group operations
type GroupHandler struct {
	sessionService session.SessionService
	meowService    *meow.MeowServiceImpl
	logger         logger.Logger
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(sessionService session.SessionService, meowService *meow.MeowServiceImpl) *GroupHandler {
	return &GroupHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("group-handler"),
	}
}

// CreateGroup godoc
// @Summary Create a new WhatsApp group
// @Description Create a new WhatsApp group with specified participants
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupCreateRequest true "Group creation request"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/create [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate participants
	if len(req.Participants) == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "At least one participant is required")
		return
	}

	// Validate phone numbers
	for _, phone := range req.Participants {
		if !utils.IsValidPhoneNumber(phone) {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format: "+phone)
			return
		}
	}

	// Create group through Meow service
	h.logger.Infof("Creating group '%s' with %d participants from session %s", req.Name, len(req.Participants), sessionID)

	groupInfo, err := h.meowService.CreateGroup(c.Request.Context(), sessionID, req.Name, req.Participants)
	if err != nil {
		h.logger.Errorf("Failed to create group: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create group", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, utils.SuccessResponse{
		Success: true,
		Message: "Group created successfully",
		Data: map[string]interface{}{
			"groupJid": groupInfo.JID.String(),
			"name":     groupInfo.Name,
		},
	})
}

// ListGroups godoc
// @Summary List all groups
// @Description Get a list of all groups the session is part of
// @Tags group
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/list [get]
func (h *GroupHandler) ListGroups(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	// List groups through Meow service
	h.logger.Infof("Listing groups for session %s", sessionID)

	groups, err := h.meowService.ListGroups(c.Request.Context(), sessionID)
	if err != nil {
		h.logger.Errorf("Failed to list groups: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list groups", err.Error())
		return
	}

	// Convert to response format
	groupList := make([]map[string]interface{}, len(groups))
	for i, group := range groups {
		groupList[i] = map[string]interface{}{
			"jid":          group.JID.String(),
			"name":         group.Name,
			"topic":        group.Topic,
			"participants": len(group.Participants),
			"owner":        group.OwnerJID.String(),
			"announce":     group.IsAnnounce,
			"locked":       group.IsLocked,
			"ephemeral":    group.IsEphemeral,
		}
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Groups retrieved successfully",
		Data: map[string]interface{}{
			"groups": groupList,
			"count":  len(groupList),
		},
	})
}

// GetGroupInfo godoc
// @Summary Get group information
// @Description Get detailed information about a specific group
// @Tags group
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param groupJid query string true "Group JID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/info [get]
func (h *GroupHandler) GetGroupInfo(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	groupJID := c.Query("groupJid")
	if groupJID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Group JID is required")
		return
	}

	// Get group info through Meow service
	h.logger.Infof("Getting info for group %s from session %s", groupJID, sessionID)

	groupInfo, err := h.meowService.GetGroupInfo(c.Request.Context(), sessionID, groupJID)
	if err != nil {
		h.logger.Errorf("Failed to get group info: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get group info", err.Error())
		return
	}

	// Convert participants to response format
	participants := make([]map[string]interface{}, len(groupInfo.Participants))
	for i, participant := range groupInfo.Participants {
		participants[i] = map[string]interface{}{
			"jid":        participant.JID.String(),
			"admin":      participant.IsAdmin,
			"superAdmin": participant.IsSuperAdmin,
		}
	}

	// Convert to response format
	responseData := map[string]interface{}{
		"jid":          groupInfo.JID.String(),
		"name":         groupInfo.Name,
		"topic":        groupInfo.Topic,
		"owner":        groupInfo.OwnerJID.String(),
		"participants": participants,
		"announce":     groupInfo.IsAnnounce,
		"locked":       groupInfo.IsLocked,
		"ephemeral":    groupInfo.IsEphemeral,
		"createdAt":    groupInfo.GroupCreated.Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group info retrieved successfully",
		Data:    responseData,
	})
}

// JoinGroup godoc
// @Summary Join a group via invite link
// @Description Join a WhatsApp group using an invite code
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupJoinRequest true "Group join request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/join [post]
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupJoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Join group through Meow service
	h.logger.Infof("Joining group with invite code %s from session %s", req.InviteCode, sessionID)

	groupInfo, err := h.meowService.JoinGroupWithLink(c.Request.Context(), sessionID, req.InviteCode)
	if err != nil {
		h.logger.Errorf("Failed to join group: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to join group", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Successfully joined group",
		Data: map[string]interface{}{
			"groupJid": groupInfo.JID.String(),
			"name":     groupInfo.Name,
		},
	})
}

// LeaveGroup godoc
// @Summary Leave a group
// @Description Leave a WhatsApp group
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupLeaveRequest true "Group leave request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/leave [post]
func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Leave group through Meow service
	h.logger.Infof("Leaving group %s from session %s", req.GroupJID, sessionID)

	err := h.meowService.LeaveGroup(c.Request.Context(), sessionID, req.GroupJID)
	if err != nil {
		h.logger.Errorf("Failed to leave group: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to leave group", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Successfully left group",
	})
}

// GetInviteLink godoc
// @Summary Get group invite link
// @Description Get the invite link for a group
// @Tags group
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param groupJid query string true "Group JID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/invitelink [get]
func (h *GroupHandler) GetInviteLink(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	groupJID := c.Query("groupJid")
	if groupJID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Group JID is required")
		return
	}

	// Get invite link through Meow service
	h.logger.Infof("Getting invite link for group %s from session %s", groupJID, sessionID)

	inviteLink, err := h.meowService.GetGroupInviteLink(c.Request.Context(), sessionID, groupJID, false)
	if err != nil {
		h.logger.Errorf("Failed to get invite link: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get invite link", err.Error())
		return
	}

	// Extract invite code from link
	inviteCode := ""
	if strings.Contains(inviteLink, "chat.whatsapp.com/") {
		parts := strings.Split(inviteLink, "/")
		if len(parts) > 0 {
			inviteCode = parts[len(parts)-1]
		}
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Invite link retrieved successfully",
		Data: map[string]interface{}{
			"inviteLink": inviteLink,
			"inviteCode": inviteCode,
		},
	})
}

// GetInviteInfo godoc
// @Summary Get invite information
// @Description Get information about a group invite without joining
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupInviteInfoRequest true "Invite info request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/inviteinfo [post]
func (h *GroupHandler) GetInviteInfo(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupInviteInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual invite info retrieval through Meow service
	h.logger.Infof("Getting invite info for code %s from session %s", req.InviteCode, sessionID)

	// Mock response
	inviteInfo := map[string]interface{}{
		"groupJid":     "120363025246125486@g.us",
		"groupName":    "Test Group",
		"groupTopic":   "This is a test group",
		"participants": 5,
		"inviteCode":   req.InviteCode,
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Invite info retrieved successfully",
		Data:    inviteInfo,
	})
}

// UpdateParticipants godoc
// @Summary Update group participants
// @Description Add, remove, promote, or demote group participants
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupUpdateParticipantsRequest true "Update participants request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/participants/update [post]
func (h *GroupHandler) UpdateParticipants(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupUpdateParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate action
	validActions := []string{"add", "remove", "promote", "demote"}
	isValid := false
	for _, action := range validActions {
		if req.Action == action {
			isValid = true
			break
		}
	}
	if !isValid {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid action. Must be: add, remove, promote, or demote")
		return
	}

	// Validate participants
	if len(req.Participants) == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "At least one participant is required")
		return
	}

	// Validate phone numbers
	for _, phone := range req.Participants {
		if !utils.IsValidPhoneNumber(phone) {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format: "+phone)
			return
		}
	}

	// Update participants through Meow service
	h.logger.Infof("Updating participants in group %s with action %s from session %s", req.GroupJID, req.Action, sessionID)

	err := h.meowService.UpdateGroupParticipants(c.Request.Context(), sessionID, req.GroupJID, req.Participants, req.Action)
	if err != nil {
		h.logger.Errorf("Failed to update participants: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update participants", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Participants updated successfully",
		Data: map[string]interface{}{
			"action":       req.Action,
			"participants": req.Participants,
		},
	})
}

// SetName godoc
// @Summary Set group name
// @Description Update the name of a group
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupSetNameRequest true "Set name request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/name/set [post]
func (h *GroupHandler) SetName(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupSetNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Set group name through Meow service
	h.logger.Infof("Setting name '%s' for group %s from session %s", req.Name, req.GroupJID, sessionID)

	err := h.meowService.SetGroupName(c.Request.Context(), sessionID, req.GroupJID, req.Name)
	if err != nil {
		h.logger.Errorf("Failed to set group name: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to set group name", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group name updated successfully",
		Data: map[string]interface{}{
			"groupJid": req.GroupJID,
			"name":     req.Name,
		},
	})
}

// SetTopic godoc
// @Summary Set group topic/description
// @Description Update the topic/description of a group
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupSetTopicRequest true "Set topic request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/topic/set [post]
func (h *GroupHandler) SetTopic(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupSetTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Set group topic through Meow service
	h.logger.Infof("Setting topic for group %s from session %s", req.GroupJID, sessionID)

	err := h.meowService.SetGroupTopic(c.Request.Context(), sessionID, req.GroupJID, req.Topic)
	if err != nil {
		h.logger.Errorf("Failed to set group topic: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to set group topic", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group topic updated successfully",
		Data: map[string]interface{}{
			"groupJid": req.GroupJID,
			"topic":    req.Topic,
		},
	})
}

// SetPhoto godoc
// @Summary Set group photo
// @Description Update the photo of a group
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupSetPhotoRequest true "Set photo request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/photo/set [post]
func (h *GroupHandler) SetPhoto(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupSetPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual group photo setting through Meow service
	h.logger.Infof("Setting photo for group %s from session %s", req.GroupJID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group photo updated successfully",
		Data: map[string]interface{}{
			"groupJid": req.GroupJID,
		},
	})
}

// RemovePhoto godoc
// @Summary Remove group photo
// @Description Remove the photo of a group
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupRemovePhotoRequest true "Remove photo request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/photo/remove [post]
func (h *GroupHandler) RemovePhoto(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupRemovePhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual group photo removal through Meow service
	h.logger.Infof("Removing photo for group %s from session %s", req.GroupJID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group photo removed successfully",
		Data: map[string]interface{}{
			"groupJid": req.GroupJID,
		},
	})
}

// SetAnnounce godoc
// @Summary Set group announce mode
// @Description Set whether only admins can send messages
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupSetAnnounceRequest true "Set announce request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/announce/set [post]
func (h *GroupHandler) SetAnnounce(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupSetAnnounceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual group announce setting through Meow service
	h.logger.Infof("Setting announce mode %v for group %s from session %s", req.Announce, req.GroupJID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group announce mode updated successfully",
		Data: map[string]interface{}{
			"groupJid": req.GroupJID,
			"announce": req.Announce,
		},
	})
}

// SetLocked godoc
// @Summary Set group locked mode
// @Description Set whether only admins can edit group info
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupSetLockedRequest true "Set locked request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/locked/set [post]
func (h *GroupHandler) SetLocked(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupSetLockedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual group locked setting through Meow service
	h.logger.Infof("Setting locked mode %v for group %s from session %s", req.Locked, req.GroupJID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group locked mode updated successfully",
		Data: map[string]interface{}{
			"groupJid": req.GroupJID,
			"locked":   req.Locked,
		},
	})
}

// SetEphemeral godoc
// @Summary Set group ephemeral messages
// @Description Set the duration for disappearing messages in a group
// @Tags group
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.GroupSetEphemeralRequest true "Set ephemeral request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/group/ephemeral/set [post]
func (h *GroupHandler) SetEphemeral(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.GroupSetEphemeralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate duration (0 = disabled, common values: 86400 = 1 day, 604800 = 7 days, 7776000 = 90 days)
	validDurations := []int64{0, 86400, 604800, 7776000}
	isValid := false
	for _, duration := range validDurations {
		if req.Duration == duration {
			isValid = true
			break
		}
	}
	if !isValid {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid duration. Must be 0 (disabled), 86400 (1 day), 604800 (7 days), or 7776000 (90 days)")
		return
	}

	// TODO: Implement actual group ephemeral setting through Meow service
	h.logger.Infof("Setting ephemeral duration %d for group %s from session %s", req.Duration, req.GroupJID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Group ephemeral messages updated successfully",
		Data: map[string]interface{}{
			"groupJid": req.GroupJID,
			"duration": req.Duration,
		},
	})
}
