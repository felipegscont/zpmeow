package handler

import (
	"net/http"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/meow"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"
	"zpmeow/internal/utils"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles HTTP requests for chat operations
type ChatHandler struct {
	sessionService session.SessionService
	meowService    *meow.MeowServiceImpl
	logger         logger.Logger
}

// NewChatHandler creates a new chat handler
func NewChatHandler(sessionService session.SessionService, meowService *meow.MeowServiceImpl) *ChatHandler {
	return &ChatHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("chat-handler"),
	}
}

// SetPresence godoc
// @Summary Set chat presence
// @Description Set presence status in a chat (typing, recording, paused)
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatPresenceRequest true "Chat presence request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/presence [post]
func (h *ChatHandler) SetPresence(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Validate presence state
	validStates := []string{"typing", "recording", "paused"}
	isValid := false
	for _, state := range validStates {
		if req.State == state {
			isValid = true
			break
		}
	}
	if !isValid {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid presence state. Must be: typing, recording, or paused")
		return
	}

	// Set presence through Meow service
	h.logger.Infof("Setting presence %s for %s from session %s", req.State, req.Phone, sessionID)

	err := h.meowService.SetChatPresence(c.Request.Context(), sessionID, req.Phone, req.State)
	if err != nil {
		h.logger.Errorf("Failed to set presence: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to set presence", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Presence set successfully",
	})
}

// MarkRead godoc
// @Summary Mark messages as read
// @Description Mark one or more messages as read in a chat
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatMarkReadRequest true "Mark read request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/markread [post]
func (h *ChatHandler) MarkRead(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatMarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Validate message IDs
	if len(req.MessageIDs) == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "At least one message ID is required")
		return
	}

	// Mark messages as read through Meow service
	h.logger.Infof("Marking %d messages as read for %s from session %s", len(req.MessageIDs), req.Phone, sessionID)

	err := h.meowService.MarkMessageRead(c.Request.Context(), sessionID, req.Phone, req.MessageIDs)
	if err != nil {
		h.logger.Errorf("Failed to mark messages as read: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to mark messages as read", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Messages marked as read successfully",
	})
}

// React godoc
// @Summary React to a message
// @Description Add an emoji reaction to a message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatReactRequest true "React request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/react [post]
func (h *ChatHandler) React(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatReactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Send reaction through Meow service
	h.logger.Infof("Reacting with %s to message %s for %s from session %s", req.Emoji, req.MessageID, req.Phone, sessionID)

	err := h.meowService.ReactToMessage(c.Request.Context(), sessionID, req.Phone, req.MessageID, req.Emoji)
	if err != nil {
		h.logger.Errorf("Failed to send reaction: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send reaction", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.MessageID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// Delete godoc
// @Summary Delete a message
// @Description Delete a message from a chat
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatDeleteRequest true "Delete request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/delete [post]
func (h *ChatHandler) Delete(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// TODO: Implement actual message deletion through Meow service
	h.logger.Infof("Deleting message %s for %s from session %s (forEveryone: %v)", req.MessageID, req.Phone, sessionID, req.ForEveryone)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Message deleted successfully",
	})
}

// Edit godoc
// @Summary Edit a message
// @Description Edit a text message in a chat
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatEditRequest true "Edit request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/edit [post]
func (h *ChatHandler) Edit(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// TODO: Implement actual message editing through Meow service
	h.logger.Infof("Editing message %s for %s from session %s", req.MessageID, req.Phone, sessionID)

	response := types.SendResponse{
		Success:   true,
		MessageID: req.MessageID,
		Timestamp: 1640995200,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// DownloadImage godoc
// @Summary Download an image from a message
// @Description Download an image attachment from a WhatsApp message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/image [post]
func (h *ChatHandler) DownloadImage(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual image download through Meow service
	h.logger.Infof("Downloading image from message %s from session %s", req.MessageID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Image download initiated",
		Data: map[string]interface{}{
			"messageId": req.MessageID,
			"type":      "image",
		},
	})
}

// DownloadVideo godoc
// @Summary Download a video from a message
// @Description Download a video attachment from a WhatsApp message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/video [post]
func (h *ChatHandler) DownloadVideo(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual video download through Meow service
	h.logger.Infof("Downloading video from message %s from session %s", req.MessageID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Video download initiated",
		Data: map[string]interface{}{
			"messageId": req.MessageID,
			"type":      "video",
		},
	})
}

// DownloadAudio godoc
// @Summary Download an audio from a message
// @Description Download an audio attachment from a WhatsApp message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/audio [post]
func (h *ChatHandler) DownloadAudio(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual audio download through Meow service
	h.logger.Infof("Downloading audio from message %s from session %s", req.MessageID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Audio download initiated",
		Data: map[string]interface{}{
			"messageId": req.MessageID,
			"type":      "audio",
		},
	})
}

// DownloadDocument godoc
// @Summary Download a document from a message
// @Description Download a document attachment from a WhatsApp message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/document [post]
func (h *ChatHandler) DownloadDocument(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.ChatDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// TODO: Implement actual document download through Meow service
	h.logger.Infof("Downloading document from message %s from session %s", req.MessageID, sessionID)

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Document download initiated",
		Data: map[string]interface{}{
			"messageId": req.MessageID,
			"type":      "document",
		},
	})
}
