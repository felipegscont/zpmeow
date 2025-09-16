package handlers

import (
	"net/http"

	"zpmeow/internal/application"
	"zpmeow/internal/infrastructure/wameow"
	"zpmeow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

// CommunityHandler handles community-related HTTP requests
type CommunityHandler struct {
	sessionService *application.SessionService
	wameowService  *wameow.MeowService
}

// NewCommunityHandler creates a new community handler
func NewCommunityHandler(sessionService *application.SessionService, wameowService *wameow.MeowService) *CommunityHandler {
	return &CommunityHandler{
		sessionService: sessionService,
		wameowService:  wameowService,
	}
}

// resolveSessionID resolves session ID or name to actual session ID
func (h *CommunityHandler) resolveSessionID(c *gin.Context, sessionIDOrName string) (string, error) {
	// For now, just return the sessionIDOrName as-is
	// In the future, this could resolve session names to IDs
	return sessionIDOrName, nil
}

// LinkGroup handles linking a group to a community
//
//	@Summary		Link group to community
//	@Description	Link a group to a community
//	@Tags			Community
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.LinkGroupRequest	true	"Link group request"
//	@Success		200			{object}	dto.CommunityResponse
//	@Failure		400			{object}	dto.CommunityResponse
//	@Failure		500			{object}	dto.CommunityResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/community/link [post]
func (h *CommunityHandler) LinkGroup(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
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
		c.JSON(http.StatusNotFound, dto.NewCommunityErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.LinkGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Link group using wameow service
	ctx := c.Request.Context()
	err = h.wameowService.LinkGroup(ctx, sessionID, req.CommunityJID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewCommunityErrorResponse(
			http.StatusInternalServerError,
			"LINK_GROUP_FAILED",
			"Failed to link group to community",
			err.Error(),
		))
		return
	}

	response := dto.NewCommunitySuccessResponse(sessionID, req.CommunityJID, req.GroupJID, "link_group")
	c.JSON(http.StatusOK, response)
}

// UnlinkGroup handles unlinking a group from a community
//
//	@Summary		Unlink group from community
//	@Description	Unlink a group from a community
//	@Tags			Community
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.UnlinkGroupRequest	true	"Unlink group request"
//	@Success		200			{object}	dto.CommunityResponse
//	@Failure		400			{object}	dto.CommunityResponse
//	@Failure		500			{object}	dto.CommunityResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/community/unlink [post]
func (h *CommunityHandler) UnlinkGroup(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
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
		c.JSON(http.StatusNotFound, dto.NewCommunityErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.UnlinkGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Unlink group using wameow service
	ctx := c.Request.Context()
	err = h.wameowService.UnlinkGroup(ctx, sessionID, req.CommunityJID, req.GroupJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewCommunityErrorResponse(
			http.StatusInternalServerError,
			"UNLINK_GROUP_FAILED",
			"Failed to unlink group from community",
			err.Error(),
		))
		return
	}

	response := dto.NewCommunitySuccessResponse(sessionID, req.CommunityJID, req.GroupJID, "unlink_group")
	c.JSON(http.StatusOK, response)
}

// GetSubGroups handles getting subgroups of a community
//
//	@Summary		Get community subgroups
//	@Description	Get all subgroups of a community
//	@Tags			Community
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.GetSubGroupsRequest	true	"Get subgroups request"
//	@Success		200			{object}	dto.CommunitySubGroupsResponse
//	@Failure		400			{object}	dto.CommunityResponse
//	@Failure		500			{object}	dto.CommunityResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/community/subgroups [post]
func (h *CommunityHandler) GetSubGroups(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
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
		c.JSON(http.StatusNotFound, dto.NewCommunityErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetSubGroupsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Get subgroups using wameow service
	ctx := c.Request.Context()
	subGroups, err := h.wameowService.GetSubGroups(ctx, sessionID, req.CommunityJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewCommunityErrorResponse(
			http.StatusInternalServerError,
			"GET_SUBGROUPS_FAILED",
			"Failed to get subgroups",
			err.Error(),
		))
		return
	}

	response := dto.NewCommunitySubGroupsResponse(sessionID, req.CommunityJID, subGroups)
	c.JSON(http.StatusOK, response)
}

// GetLinkedGroupsParticipants handles getting participants of linked groups
//
//	@Summary		Get community participants
//	@Description	Get all participants of linked groups in a community
//	@Tags			Community
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string									true	"Session ID"
//	@Param			request		body		dto.GetLinkedGroupsParticipantsRequest	true	"Get participants request"
//	@Success		200			{object}	dto.CommunityParticipantsResponse
//	@Failure		400			{object}	dto.CommunityResponse
//	@Failure		500			{object}	dto.CommunityResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/community/participants [post]
func (h *CommunityHandler) GetLinkedGroupsParticipants(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
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
		c.JSON(http.StatusNotFound, dto.NewCommunityErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.GetLinkedGroupsParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewCommunityErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Get linked groups participants using wameow service
	ctx := c.Request.Context()
	participants, err := h.wameowService.GetLinkedGroupsParticipants(ctx, sessionID, req.CommunityJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewCommunityErrorResponse(
			http.StatusInternalServerError,
			"GET_LINKED_GROUPS_PARTICIPANTS_FAILED",
			"Failed to get linked groups participants",
			err.Error(),
		))
		return
	}

	response := dto.NewCommunityParticipantsResponse(sessionID, req.CommunityJID, participants)
	c.JSON(http.StatusOK, response)
}
