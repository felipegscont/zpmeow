package handlers

import (
	"context"
	"net/http"

	"zpmeow/internal/application"
	"zpmeow/internal/infrastructure/wameow"
	"zpmeow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	sessionService *application.SessionService
	wameowService  *wameow.MeowService
}

// NewUserHandler creates a new user handler
func NewUserHandler(sessionService *application.SessionService, wameowService *wameow.MeowService) *UserHandler {
	return &UserHandler{
		sessionService: sessionService,
		wameowService:  wameowService,
	}
}

// CheckUser handles checking if users are on WhatsApp
// @Summary Check users on WhatsApp
// @Description Check if phone numbers are registered on WhatsApp
// @Tags User
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.CheckUserRequest true "Check user request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.UserResponse
// @Failure 500 {object} dto.UserResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/user/check [post]
func (h *UserHandler) CheckUser(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.CheckUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if len(req.Phones) == 0 {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONES",
			"At least one phone number is required",
			"",
		))
		return
	}

	// Check users via wameow service
	ctx := context.Background()
	results, err := h.wameowService.CheckUser(ctx, sessionID, req.Phones)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewUserErrorResponse(
			http.StatusInternalServerError,
			"CHECK_USER_FAILED",
			"Failed to check users",
			err.Error(),
		))
		return
	}

	// Convert wameow results to DTO format
	var checkResults []dto.UserCheckResult
	for _, result := range results {
		checkResults = append(checkResults, dto.UserCheckResult{
			Query:        result.Query,
			IsInWhatsapp: result.IsInWhatsapp,
			JID:          result.JID,
			VerifiedName: result.VerifiedName,
		})
	}

	response := dto.NewUserSuccessResponse("check_users", checkResults, nil)
	c.JSON(http.StatusOK, response)
}

// GetUserInfo handles getting user information
// @Summary Get user information
// @Description Get detailed information about users
// @Tags User
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.GetUserInfoRequest true "Get user info request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.UserResponse
// @Failure 500 {object} dto.UserResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/user/info [post]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.GetUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if len(req.Phones) == 0 {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONES",
			"At least one phone number is required",
			"",
		))
		return
	}

	// Get user info via wameow service
	ctx := context.Background()
	results, err := h.wameowService.GetUserInfo(ctx, sessionID, req.Phones)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewUserErrorResponse(
			http.StatusInternalServerError,
			"GET_USER_INFO_FAILED",
			"Failed to get user information",
			err.Error(),
		))
		return
	}

	// Convert wameow results to DTO format
	var userInfos []dto.UserInfo
	for _, result := range results {
		userInfos = append(userInfos, dto.UserInfo{
			JID:          result.JID,
			DisplayName:  result.DisplayName,
			VerifiedName: result.VerifiedName,
			Avatar:       result.Avatar,
			Status:       result.Status,
			PictureID:    result.PictureID,
			DeviceCount:  result.DeviceCount,
		})
	}

	response := dto.NewUserSuccessResponse("get_user_info", nil, userInfos)
	c.JSON(http.StatusOK, response)
}

// GetAvatar handles getting user avatar/profile picture
// @Summary Get user avatar
// @Description Get user's profile picture/avatar
// @Tags User
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.GetAvatarRequest true "Get avatar request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.UserResponse
// @Failure 500 {object} dto.UserResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/user/avatar [post]
func (h *UserHandler) GetAvatar(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.GetAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	// Get avatar via wameow service
	ctx := context.Background()
	result, err := h.wameowService.GetAvatar(ctx, sessionID, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewUserErrorResponse(
			http.StatusInternalServerError,
			"GET_AVATAR_FAILED",
			"Failed to get user avatar",
			err.Error(),
		))
		return
	}

	// Convert wameow result to DTO format
	avatarInfo := &dto.AvatarInfo{
		Phone:     result.Phone,
		JID:       result.JID,
		AvatarURL: result.AvatarURL,
		PictureID: result.PictureID,
	}

	response := dto.NewUserAvatarResponse(avatarInfo)
	c.JSON(http.StatusOK, response)
}

// SetUserPresence handles setting global user presence
// @Summary Set user presence
// @Description Set global user presence (available/unavailable)
// @Tags User
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.SetUserPresenceRequest true "Set presence request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.UserResponse
// @Failure 500 {object} dto.UserResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/user/presence [post]
func (h *UserHandler) SetUserPresence(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.SetUserPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.State == "" {
		c.JSON(http.StatusBadRequest, dto.NewUserErrorResponse(
			http.StatusBadRequest,
			"MISSING_STATE",
			"State is required",
			"Valid states: available, unavailable",
		))
		return
	}

	// Set user presence via wameow service
	ctx := context.Background()
	err := h.wameowService.SetUserPresence(ctx, sessionID, req.State)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewUserErrorResponse(
			http.StatusInternalServerError,
			"SET_PRESENCE_FAILED",
			"Failed to set user presence",
			err.Error(),
		))
		return
	}

	response := dto.NewUserSuccessResponse("set_user_presence", nil, nil)
	c.JSON(http.StatusOK, response)
}

// GetContacts handles getting user contacts
// @Summary Get contacts
// @Description Get all contacts from user's WhatsApp
// @Tags User
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} dto.ContactsResponse
// @Failure 500 {object} dto.UserResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/user/contacts [get]
func (h *UserHandler) GetContacts(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Get contacts via wameow service
	ctx := context.Background()
	results, err := h.wameowService.GetContacts(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewUserErrorResponse(
			http.StatusInternalServerError,
			"GET_CONTACTS_FAILED",
			"Failed to get contacts",
			err.Error(),
		))
		return
	}

	// Convert wameow results to DTO format
	var contacts []dto.ContactInfo
	for _, result := range results {
		contacts = append(contacts, dto.ContactInfo{
			JID:          result.JID,
			Name:         result.Name,
			Notify:       result.Notify,
			PushName:     result.PushName,
			BusinessName: result.BusinessName,
			IsBlocked:    result.IsBlocked,
			IsMuted:      result.IsMuted,
		})
	}

	response := dto.NewContactsResponse(contacts)
	c.JSON(http.StatusOK, response)
}







// GetBlockedContacts handles getting blocked contacts
func (h *UserHandler) GetBlockedContacts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get blocked contacts endpoint - implementation pending",
	})
}

// UpdateProfile handles updating user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update profile endpoint - implementation pending",
	})
}

// SetProfilePicture handles setting profile picture
func (h *UserHandler) SetProfilePicture(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Set profile picture endpoint - implementation pending",
	})
}

// RemoveProfilePicture handles removing profile picture
func (h *UserHandler) RemoveProfilePicture(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Remove profile picture endpoint - implementation pending",
	})
}

// Router compatibility methods
func (h *UserHandler) CheckUsers(c *gin.Context) { h.CheckUser(c) }
func (h *UserHandler) GetUserAvatar(c *gin.Context) { h.GetAvatar(c) }
func (h *UserHandler) GetUserStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Get user status - implementation pending"})
}
func (h *UserHandler) SetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Set status - implementation pending"})
}
func (h *UserHandler) GetPrivacySettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Get privacy settings - implementation pending"})
}
func (h *UserHandler) UpdatePrivacySettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Update privacy settings - implementation pending"})
}
func (h *UserHandler) SetPresence(c *gin.Context) { h.SetUserPresence(c) }
func (h *UserHandler) BlockUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Block user - implementation pending"})
}
func (h *UserHandler) UnblockUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Unblock user - implementation pending"})
}
