package handler

import (
	"net/http"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/infra/meow"
	"zpmeow/internal/types"
	"zpmeow/internal/utils"

	"github.com/gin-gonic/gin"
	waTypes "go.mau.fi/whatsmeow/types"
)

// UserHandler handles user-related operations
type UserHandler struct {
	sessionService session.SessionService
	meowService    *meow.MeowServiceImpl
	logger         logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(sessionService session.SessionService, meowService *meow.MeowServiceImpl) *UserHandler {
	return &UserHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("user-handler"),
	}
}

// Helper function to resolve session ID from path parameter
func (h *UserHandler) resolveSessionID(c *gin.Context) (string, bool) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return "", false
	}

	// Check if session exists
	_, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Session not found", err.Error())
		return "", false
	}

	return sessionID, true
}

// @Summary Set user presence
// @Description Set global presence status for the WhatsApp session
// @Tags user
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.UserPresenceRequest true "Presence request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/user/presence [post]
func (h *UserHandler) SetPresence(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	var req types.UserPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate presence type
	var presence waTypes.Presence
	switch req.Type {
	case "available":
		presence = waTypes.PresenceAvailable
	case "unavailable":
		presence = waTypes.PresenceUnavailable
	default:
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid presence type. Allowed values: available, unavailable")
		return
	}

	h.logger.Infof("Setting global presence to %s for session %s", req.Type, sessionID)

	err := h.meowService.SetGlobalPresence(c.Request.Context(), sessionID, presence)
	if err != nil {
		h.logger.Errorf("Failed to set presence: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to set presence", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Presence set successfully",
		Data: map[string]interface{}{
			"Details": "Presence set successfully",
		},
	})
}

// @Summary Check if users are on WhatsApp
// @Description Check if phone numbers are registered on WhatsApp
// @Tags user
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.CheckUserRequest true "Check user request"
// @Success 200 {object} types.UserCheckResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/user/check [post]
func (h *UserHandler) CheckUser(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	var req types.CheckUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone numbers
	for _, phone := range req.Phone {
		if !utils.IsValidPhoneNumber(phone) {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format: "+phone)
			return
		}
	}

	h.logger.Infof("Checking %d phone numbers for session %s", len(req.Phone), sessionID)

	resp, err := h.meowService.CheckUsersOnWhatsApp(c.Request.Context(), sessionID, req.Phone)
	if err != nil {
		h.logger.Errorf("Failed to check users: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to check users", err.Error())
		return
	}

	// Convert response
	var users []types.UserCheckResult
	for _, item := range resp {
		verifiedName := ""
		if item.VerifiedName != nil {
			verifiedName = item.VerifiedName.Details.GetVerifiedName()
		}

		users = append(users, types.UserCheckResult{
			Query:        item.Query,
			IsInWhatsapp: item.IsIn,
			JID:          item.JID.String(),
			VerifiedName: verifiedName,
		})
	}

	response := types.UserCheckResponse{
		Users: users,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// @Summary Get user information
// @Description Get detailed information about WhatsApp users
// @Tags user
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.GetUserInfoRequest true "Get user info request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/user/info [post]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	var req types.GetUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone numbers
	for _, phone := range req.Phone {
		if !utils.IsValidPhoneNumber(phone) {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format: "+phone)
			return
		}
	}

	h.logger.Infof("Getting user info for %d phone numbers from session %s", len(req.Phone), sessionID)

	userInfo, err := h.meowService.GetUserInfo(c.Request.Context(), sessionID, req.Phone)
	if err != nil {
		h.logger.Errorf("Failed to get user info: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get user info", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, map[string]interface{}{
		"users": userInfo,
	})
}

// @Summary Get user avatar
// @Description Get avatar URL for a WhatsApp user
// @Tags user
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.GetAvatarRequest true "Get avatar request"
// @Success 200 {object} types.AvatarResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/user/avatar [post]
func (h *UserHandler) GetAvatar(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	var req types.GetAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	h.logger.Infof("Getting avatar for %s from session %s", req.Phone, sessionID)

	avatarInfo, err := h.meowService.GetUserAvatar(c.Request.Context(), sessionID, req.Phone, req.Preview)
	if err != nil {
		h.logger.Errorf("Failed to get avatar: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get avatar", err.Error())
		return
	}

	response := types.AvatarResponse{
		URL:       avatarInfo.URL,
		ID:        avatarInfo.ID,
		Type:      avatarInfo.Type,
		DirectURL: avatarInfo.DirectURL,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// @Summary Get contacts
// @Description Get all contacts from the WhatsApp session
// @Tags user
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/user/contacts [get]
func (h *UserHandler) GetContacts(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	h.logger.Infof("Getting contacts for session %s", sessionID)

	contacts, err := h.meowService.GetContacts(c.Request.Context(), sessionID)
	if err != nil {
		h.logger.Errorf("Failed to get contacts: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get contacts", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, contacts)
}
