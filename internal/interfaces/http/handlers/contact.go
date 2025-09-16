package handlers

import (
	"net/http"

	"zpmeow/internal/application"
	"zpmeow/internal/infrastructure/wameow"
	"zpmeow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

// ContactHandler handles contact-related HTTP requests
type ContactHandler struct {
	sessionService *application.SessionService
	wameowService  *wameow.MeowService
}

// NewContactHandler creates a new contact handler
func NewContactHandler(sessionService *application.SessionService, wameowService *wameow.MeowService) *ContactHandler {
	return &ContactHandler{
		sessionService: sessionService,
		wameowService:  wameowService,
	}
}

// CheckContact handles checking if contacts are on WhatsApp
//
//	@Summary		Check contacts on WhatsApp
//	@Description	Check if phone numbers are registered on WhatsApp
//	@Tags			Contacts
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.CheckContactRequest	true	"Check contact request"
//	@Success		200			{object}	dto.ContactResponse
//	@Failure		400			{object}	dto.ContactResponse
//	@Failure		500			{object}	dto.ContactResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/contacts/check [post]
func (h *ContactHandler) CheckContact(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.CheckContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if len(req.Phones) == 0 {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONES",
			"At least one phone number is required",
			"",
		))
		return
	}

	// Check contacts via wameow service
	ctx := c.Request.Context()
	results, err := h.wameowService.CheckUser(ctx, sessionID, req.Phones)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewContactErrorResponse(
			http.StatusInternalServerError,
			"CHECK_CONTACT_FAILED",
			"Failed to check contacts",
			err.Error(),
		))
		return
	}

	// Convert wameow results to DTO format
	var checkResults []dto.ContactCheckResult
	for _, result := range results {
		checkResults = append(checkResults, dto.ContactCheckResult{
			Query:        result.Query,
			IsInWhatsapp: result.IsInWhatsapp,
			JID:          result.JID,
			VerifiedName: result.VerifiedName,
		})
	}

	response := dto.NewContactSuccessResponse("check_contacts", checkResults, nil)
	c.JSON(http.StatusOK, response)
}

// GetContactInfo handles getting contact information
//
//	@Summary		Get contact information
//	@Description	Get detailed information about contacts
//	@Tags			Contacts
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			request		body		dto.GetContactInfoRequest	true	"Get contact info request"
//	@Success		200			{object}	dto.ContactResponse
//	@Failure		400			{object}	dto.ContactResponse
//	@Failure		500			{object}	dto.ContactResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/contacts/info [post]
func (h *ContactHandler) GetContactInfo(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.GetContactInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if len(req.Phones) == 0 {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONES",
			"At least one phone number is required",
			"",
		))
		return
	}

	// Get contact info via wameow service
	ctx := c.Request.Context()
	results, err := h.wameowService.GetUserInfo(ctx, sessionID, req.Phones)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewContactErrorResponse(
			http.StatusInternalServerError,
			"GET_CONTACT_INFO_FAILED",
			"Failed to get contact information",
			err.Error(),
		))
		return
	}

	// Convert wameow results to DTO format
	var contactInfos []dto.ContactInfo
	for _, result := range results {
		contactInfos = append(contactInfos, dto.ContactInfo{
			JID:          result.JID,
			DisplayName:  result.DisplayName,
			VerifiedName: result.VerifiedName,
			Avatar:       result.Avatar,
			Status:       result.Status,
			PictureID:    result.PictureID,
			DeviceCount:  result.DeviceCount,
		})
	}

	response := dto.NewContactSuccessResponse("get_contact_info", nil, contactInfos)
	c.JSON(http.StatusOK, response)
}

// GetAvatar handles getting contact avatar/profile picture
//
//	@Summary		Get contact avatar
//	@Description	Get contact's profile picture/avatar
//	@Tags			Contacts
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.GetAvatarRequest	true	"Get avatar request"
//	@Success		200			{object}	dto.ContactResponse
//	@Failure		400			{object}	dto.ContactResponse
//	@Failure		500			{object}	dto.ContactResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/contacts/avatar [post]
func (h *ContactHandler) GetAvatar(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.GetAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	// Get avatar via wameow service
	ctx := c.Request.Context()
	result, err := h.wameowService.GetAvatar(ctx, sessionID, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewContactErrorResponse(
			http.StatusInternalServerError,
			"GET_AVATAR_FAILED",
			"Failed to get contact avatar",
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

	response := dto.NewContactAvatarResponse(avatarInfo)
	c.JSON(http.StatusOK, response)
}

// SetPresence handles setting global contact presence
//
//	@Summary		Set contact presence
//	@Description	Set global contact presence (available/unavailable)
//	@Tags			Contacts
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			request		body		dto.SetContactPresenceRequest	true	"Set presence request"
//	@Success		200			{object}	dto.ContactResponse
//	@Failure		400			{object}	dto.ContactResponse
//	@Failure		500			{object}	dto.ContactResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/contacts/presence [post]
func (h *ContactHandler) SetPresence(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.SetContactPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.State == "" {
		c.JSON(http.StatusBadRequest, dto.NewContactErrorResponse(
			http.StatusBadRequest,
			"MISSING_STATE",
			"State is required",
			"Valid states: available, unavailable",
		))
		return
	}

	// Set contact presence via wameow service
	ctx := c.Request.Context()
	err := h.wameowService.SetUserPresence(ctx, sessionID, req.State)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewContactErrorResponse(
			http.StatusInternalServerError,
			"SET_PRESENCE_FAILED",
			"Failed to set contact presence",
			err.Error(),
		))
		return
	}

	response := dto.NewContactSuccessResponse("set_contact_presence", nil, nil)
	c.JSON(http.StatusOK, response)
}

// GetContacts handles getting contact list
//
//	@Summary		Get contacts
//	@Description	Get all contacts from contact's WhatsApp
//	@Tags			Contacts
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		200			{object}	dto.ContactsResponse
//	@Failure		500			{object}	dto.ContactResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/contacts/list [get]
func (h *ContactHandler) GetContacts(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Get contacts via wameow service
	ctx := c.Request.Context()
	results, err := h.wameowService.GetContacts(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewContactErrorResponse(
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
func (h *ContactHandler) GetBlockedContacts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get blocked contacts endpoint - implementation pending",
	})
}

// UpdateProfile handles updating contact profile
func (h *ContactHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update profile endpoint - implementation pending",
	})
}

// SetProfilePicture handles setting profile picture
func (h *ContactHandler) SetProfilePicture(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Set profile picture endpoint - implementation pending",
	})
}

// RemoveProfilePicture handles removing profile picture
func (h *ContactHandler) RemoveProfilePicture(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Remove profile picture endpoint - implementation pending",
	})
}

// Router compatibility methods
func (h *ContactHandler) CheckUser(c *gin.Context)     { h.CheckContact(c) }
func (h *ContactHandler) GetUserInfo(c *gin.Context)   { h.GetContactInfo(c) }
func (h *ContactHandler) CheckUsers(c *gin.Context)    { h.CheckContact(c) }
func (h *ContactHandler) GetUserAvatar(c *gin.Context) { h.GetAvatar(c) }
func (h *ContactHandler) GetUserStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Get contact status - implementation pending"})
}
func (h *ContactHandler) SetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Set status - implementation pending"})
}
func (h *ContactHandler) GetPrivacySettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Get privacy settings - implementation pending"})
}
func (h *ContactHandler) UpdatePrivacySettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Update privacy settings - implementation pending"})
}
func (h *ContactHandler) SetUserPresence(c *gin.Context) { h.SetPresence(c) }
func (h *ContactHandler) BlockUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Block contact - implementation pending"})
}
func (h *ContactHandler) UnblockUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Unblock contact - implementation pending"})
}
