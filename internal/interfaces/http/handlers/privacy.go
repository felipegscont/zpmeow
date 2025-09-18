package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	waTypes "go.mau.fi/whatsmeow/types"

	"meow/internal/application"
	"meow/internal/infrastructure/wmeow"
	"meow/internal/interfaces/dto"
)

type PrivacyHandler struct {
	sessionService *application.SessionApp
	wmeowService   wmeow.Service
}

func NewPrivacyHandler(sessionService *application.SessionApp, wmeowService wmeow.Service) *PrivacyHandler {
	return &PrivacyHandler{
		sessionService: sessionService,
		wmeowService:   wmeowService,
	}
}

//	@Summary		Set multiple privacy settings
//	@Description	Set multiple privacy settings in a single request
//	@Tags			Privacy
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string								true	"Session ID"
//	@Param			request		body		dto.SetAllPrivacySettingsRequest	true	"Privacy settings request"
//	@Success		200			{object}	dto.PrivacySettingsResponse
//	@Failure		400			{object}	dto.PrivacySettingsResponse
//	@Failure		404			{object}	dto.PrivacySettingsResponse
//	@Failure		500			{object}	dto.PrivacySettingsResponse
//	@Router			/session/{sessionId}/privacy/set [put]
func (h *PrivacyHandler) SetAllPrivacySettings(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.SetAllPrivacySettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.PrivacySettingsResponse{
			Success: false,
			Error: &dto.PrivacyErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
		})
		return
	}

	if !req.HasAnySettings() {
		c.JSON(http.StatusBadRequest, dto.PrivacySettingsResponse{
			Success: false,
			Error: &dto.PrivacyErrorResponse{
				Code:    "NO_SETTINGS_PROVIDED",
				Message: "No privacy settings provided to update",
				Details: "At least one privacy setting must be provided",
			},
		})
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, dto.PrivacySettingsResponse{
			Success: false,
			Error: &dto.PrivacyErrorResponse{
				Code:    "VALIDATION_ERRORS",
				Message: "One or more privacy settings are invalid",
				Details: fmt.Sprintf("Validation errors: %v", validationErrors),
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	currentSettings, err := h.wmeowService.GetPrivacySettings(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
			Success: false,
			Error: &dto.PrivacyErrorResponse{
				Code:    "PRIVACY_FETCH_ERROR",
				Message: "Failed to fetch current privacy settings",
				Details: err.Error(),
			},
		})
		return
	}

	updatedSettings := []string{}

	if req.GroupAdd != nil {
		_, err := h.wmeowService.SetPrivacySetting(ctx, sessionID, "groupadd", *req.GroupAdd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "PRIVACY_UPDATE_ERROR",
					Message: "Failed to update group add privacy",
					Details: err.Error(),
				},
			})
			return
		}
		updatedSettings = append(updatedSettings, "groupAdd")
		currentSettings.GroupAdd = *req.GroupAdd
	}

	if req.LastSeen != nil {
		_, err := h.wmeowService.SetPrivacySetting(ctx, sessionID, "last", *req.LastSeen)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "PRIVACY_UPDATE_ERROR",
					Message: "Failed to update last seen privacy",
					Details: err.Error(),
				},
			})
			return
		}
		updatedSettings = append(updatedSettings, "lastSeen")
		currentSettings.LastSeen = *req.LastSeen
	}

	if req.Status != nil {
		_, err := h.wmeowService.SetPrivacySetting(ctx, sessionID, "status", *req.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "PRIVACY_UPDATE_ERROR",
					Message: "Failed to update status privacy",
					Details: err.Error(),
				},
			})
			return
		}
		updatedSettings = append(updatedSettings, "status")
		currentSettings.Status = *req.Status
	}

	if req.Profile != nil {
		_, err := h.wmeowService.SetPrivacySetting(ctx, sessionID, "profile", *req.Profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "PRIVACY_UPDATE_ERROR",
					Message: "Failed to update profile privacy",
					Details: err.Error(),
				},
			})
			return
		}
		updatedSettings = append(updatedSettings, "profile")
		currentSettings.Profile = *req.Profile
	}

	if req.ReadReceipts != nil {
		value := "none"
		if *req.ReadReceipts {
			value = "all"
		}
		_, err := h.wmeowService.SetPrivacySetting(ctx, sessionID, "readreceipts", value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "PRIVACY_UPDATE_ERROR",
					Message: "Failed to update read receipts privacy",
					Details: err.Error(),
				},
			})
			return
		}
		updatedSettings = append(updatedSettings, "readReceipts")
		currentSettings.ReadReceipts = value
	}

	if req.CallAdd != nil {
		_, err := h.wmeowService.SetPrivacySetting(ctx, sessionID, "calladd", *req.CallAdd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "PRIVACY_UPDATE_ERROR",
					Message: "Failed to update call privacy",
					Details: err.Error(),
				},
			})
			return
		}
		updatedSettings = append(updatedSettings, "callAdd")
		currentSettings.CallAdd = *req.CallAdd
	}

	if req.Online != nil {
		_, err := h.wmeowService.SetPrivacySetting(ctx, sessionID, "online", *req.Online)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "PRIVACY_UPDATE_ERROR",
					Message: "Failed to update online privacy",
					Details: err.Error(),
				},
			})
			return
		}
		updatedSettings = append(updatedSettings, "online")
		currentSettings.Online = *req.Online
	}

	data := &dto.PrivacySettingsData{
		GroupAdd:     currentSettings.GroupAdd,
		LastSeen:     currentSettings.LastSeen,
		Status:       currentSettings.Status,
		Profile:      currentSettings.Profile,
		ReadReceipts: currentSettings.ReadReceipts,
		CallAdd:      currentSettings.CallAdd,
		Online:       currentSettings.Online,
	}

	c.JSON(http.StatusOK, dto.PrivacySettingsResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully updated %d privacy settings: %v", len(updatedSettings), updatedSettings),
		Data:    data,
	})
}

//	@Summary		Get blocked contacts list
//	@Description	Get the list of blocked contacts (blocklist)
//	@Tags			Privacy
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		200			{object}	dto.BlocklistResponse
//	@Failure		400			{object}	dto.BlocklistResponse
//	@Failure		404			{object}	dto.BlocklistResponse
//	@Failure		500			{object}	dto.BlocklistResponse
//	@Router			/session/{sessionId}/privacy/blocklist [get]
func (h *PrivacyHandler) GetBlocklist(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if err := validateSessionID(sessionID); err != nil {
		c.JSON(http.StatusBadRequest, dto.BlocklistResponse{
			Success: false,
			Error:   createValidationError("session ID", err.Error()),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	blocklist, err := h.wmeowService.GetBlocklist(ctx, sessionID)
	if err != nil {
		if strings.Contains(err.Error(), "client not found") {
			c.JSON(http.StatusNotFound, dto.BlocklistResponse{
				Success: false,
				Error:   createNotFoundError("Session"),
			})
		} else if strings.Contains(err.Error(), "not connected") {
			c.JSON(http.StatusServiceUnavailable, dto.BlocklistResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "SESSION_NOT_CONNECTED",
					Message: "Session is not connected",
					Details: "Please ensure the meow session is connected before performing this operation",
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.BlocklistResponse{
				Success: false,
				Error:   createInternalError("fetch blocklist", err.Error()),
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.BlocklistResponse{
		Success: true,
		Message: fmt.Sprintf("Blocklist retrieved successfully (%d blocked contacts)", len(blocklist)),
		Data:    blocklist,
	})
}

//	@Summary		Update blocklist (block/unblock contact)
//	@Description	Block or unblock a contact
//	@Tags			Privacy
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			request		body		dto.UpdateBlocklistRequest	true	"Update blocklist request"
//	@Success		200			{object}	dto.BlocklistResponse
//	@Failure		400			{object}	dto.BlocklistResponse
//	@Failure		404			{object}	dto.BlocklistResponse
//	@Failure		500			{object}	dto.BlocklistResponse
//	@Router			/session/{sessionId}/privacy/blocklist [put]
func (h *PrivacyHandler) UpdateBlocklist(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if err := validateSessionID(sessionID); err != nil {
		c.JSON(http.StatusBadRequest, dto.BlocklistResponse{
			Success: false,
			Error:   createValidationError("session ID", err.Error()),
		})
		return
	}

	var req dto.UpdateBlocklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BlocklistResponse{
			Success: false,
			Error:   createValidationError("request body", err.Error()),
		})
		return
	}

	if err := validateJID(req.JID); err != nil {
		c.JSON(http.StatusBadRequest, dto.BlocklistResponse{
			Success: false,
			Error:   createValidationError("JID", err.Error()),
		})
		return
	}

	if err := validateBlocklistAction(req.Action); err != nil {
		c.JSON(http.StatusBadRequest, dto.BlocklistResponse{
			Success: false,
			Error:   createValidationError("action", err.Error()),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	updatedBlocklist, err := h.wmeowService.UpdateBlocklist(ctx, sessionID, req.JID, req.Action)
	if err != nil {
		if strings.Contains(err.Error(), "client not found") {
			c.JSON(http.StatusNotFound, dto.BlocklistResponse{
				Success: false,
				Error:   createNotFoundError("Session"),
			})
		} else if strings.Contains(err.Error(), "not connected") {
			c.JSON(http.StatusServiceUnavailable, dto.BlocklistResponse{
				Success: false,
				Error: &dto.PrivacyErrorResponse{
					Code:    "SESSION_NOT_CONNECTED",
					Message: "Session is not connected",
					Details: "Please ensure the meow session is connected before performing this operation",
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.BlocklistResponse{
				Success: false,
				Error:   createInternalError(fmt.Sprintf("%s contact", req.Action), err.Error()),
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.BlocklistResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully %sed contact %s (%d contacts in blocklist)", req.Action, req.JID, len(updatedBlocklist)),
		Data:    updatedBlocklist,
	})
}


func validateJID(jidStr string) error {
	if jidStr == "" {
		return fmt.Errorf("JID cannot be empty")
	}

	if !strings.Contains(jidStr, "@") {
		return fmt.Errorf("JID must contain '@' symbol")
	}

	_, err := waTypes.ParseJID(jidStr)
	if err != nil {
		return fmt.Errorf("invalid JID format: %w", err)
	}

	return nil
}

func validateSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	if len(sessionID) != 36 {
		return fmt.Errorf("session ID must be a valid UUID")
	}

	return nil
}

func validateBlocklistAction(action string) error {
	if action == "" {
		return fmt.Errorf("action cannot be empty")
	}

	if action != "block" && action != "unblock" {
		return fmt.Errorf("action must be 'block' or 'unblock', got '%s'", action)
	}

	return nil
}

func createValidationError(field, message string) *dto.PrivacyErrorResponse {
	return &dto.PrivacyErrorResponse{
		Code:    "VALIDATION_ERROR",
		Message: fmt.Sprintf("Invalid %s", field),
		Details: message,
	}
}

func createInternalError(operation, details string) *dto.PrivacyErrorResponse {
	return &dto.PrivacyErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: fmt.Sprintf("Failed to %s", operation),
		Details: details,
	}
}

func createNotFoundError(resource string) *dto.PrivacyErrorResponse {
	return &dto.PrivacyErrorResponse{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Details: fmt.Sprintf("The requested %s could not be found", strings.ToLower(resource)),
	}
}

//	@Summary		Find specific privacy settings
//	@Description	Get specific privacy settings or all settings if none specified
//	@Tags			Privacy
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string							true	"Session ID"
//	@Param			request		body		dto.FindPrivacySettingsRequest	false	"Find privacy settings request (optional)"
//	@Success		200			{object}	dto.PrivacySettingsResponse
//	@Failure		400			{object}	dto.PrivacySettingsResponse
//	@Failure		404			{object}	dto.PrivacySettingsResponse
//	@Failure		500			{object}	dto.PrivacySettingsResponse
//	@Router			/session/{sessionId}/privacy/find [post]
func (h *PrivacyHandler) FindPrivacySettings(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.FindPrivacySettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Settings = []string{}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	allSettings, err := h.wmeowService.GetPrivacySettings(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.PrivacySettingsResponse{
			Success: false,
			Error: &dto.PrivacyErrorResponse{
				Code:    "PRIVACY_FETCH_ERROR",
				Message: "Failed to fetch privacy settings",
				Details: err.Error(),
			},
		})
		return
	}

	if len(req.Settings) == 0 {
		data := &dto.PrivacySettingsData{
			GroupAdd:     allSettings.GroupAdd,
			LastSeen:     allSettings.LastSeen,
			Status:       allSettings.Status,
			Profile:      allSettings.Profile,
			ReadReceipts: allSettings.ReadReceipts,
			CallAdd:      allSettings.CallAdd,
			Online:       allSettings.Online,
		}

		c.JSON(http.StatusOK, dto.PrivacySettingsResponse{
			Success: true,
			Message: "All privacy settings retrieved successfully",
			Data:    data,
		})
		return
	}

	filteredData := &dto.PrivacySettingsData{}
	requestedSettings := []string{}

	for _, setting := range req.Settings {
		switch setting {
		case "groupAdd":
			filteredData.GroupAdd = allSettings.GroupAdd
			requestedSettings = append(requestedSettings, "groupAdd")
		case "lastSeen":
			filteredData.LastSeen = allSettings.LastSeen
			requestedSettings = append(requestedSettings, "lastSeen")
		case "status":
			filteredData.Status = allSettings.Status
			requestedSettings = append(requestedSettings, "status")
		case "profile":
			filteredData.Profile = allSettings.Profile
			requestedSettings = append(requestedSettings, "profile")
		case "readReceipts":
			filteredData.ReadReceipts = allSettings.ReadReceipts
			requestedSettings = append(requestedSettings, "readReceipts")
		case "callAdd":
			filteredData.CallAdd = allSettings.CallAdd
			requestedSettings = append(requestedSettings, "callAdd")
		case "online":
			filteredData.Online = allSettings.Online
			requestedSettings = append(requestedSettings, "online")
		}
	}

	c.JSON(http.StatusOK, dto.PrivacySettingsResponse{
		Success: true,
		Message: fmt.Sprintf("Requested privacy settings retrieved successfully: %v", requestedSettings),
		Data:    filteredData,
	})
}
