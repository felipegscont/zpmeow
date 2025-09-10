package handler

import (
	"net/http"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/infra/meow"
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

// resolveSessionID resolves session ID from path parameter (handles both ID and name)
func (h *ChatHandler) resolveSessionID(c *gin.Context) (string, bool) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return "", false
	}

	// Resolve session to get the actual ID (in case sessionID is a name)
	sess, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		if err == session.ErrSessionNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Session not found")
		} else {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to resolve session", err.Error())
		}
		return "", false
	}

	return sess.ID, true
}

// SetPresence godoc
// @Summary Set chat presence (typing/recording/paused)
// @Description Set chat presence status for a specific chat
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatPresenceRequest true "Chat presence request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/presence [post]
func (h *ChatHandler) SetPresence(c *gin.Context) {
	// Resolve session ID
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	// Parse request body
	var req types.ChatPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Map user-friendly states to whatsmeow ChatPresence values
	state, media := h.mapPresenceState(req.State, req.Media)
	if state == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid state. Allowed values: typing, recording, paused")
		return
	}

	// Set chat presence through Meow service
	h.logger.Infof("Setting chat presence for %s to %s from session %s", req.Phone, state, sessionID)

	err := h.meowService.SetChatPresence(c.Request.Context(), sessionID, req.Phone, state, media)
	if err != nil {
		h.logger.Errorf("Failed to set chat presence: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to set chat presence", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Chat presence set successfully",
	})
}

// mapPresenceState maps user-friendly presence states to whatsmeow values
func (h *ChatHandler) mapPresenceState(inputState, inputMedia string) (state, media string) {
	switch inputState {
	case "typing":
		return "composing", "" // ChatPresenceMediaText is empty string
	case "recording":
		return "composing", "audio" // ChatPresenceMediaAudio
	case "paused":
		return "paused", ""
	default:
		return "", "" // Invalid state
	}
}

// MarkRead godoc
// @Summary Mark messages as read
// @Description Mark one or more messages as read in a chat
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatMarkReadRequest true "Mark read request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/markread [post]
func (h *ChatHandler) MarkRead(c *gin.Context) {
	// Resolve session ID
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	// Parse request body
	var req types.ChatMarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number
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
// @Description Send a reaction emoji to a message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatReactRequest true "React request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/react [post]
func (h *ChatHandler) React(c *gin.Context) {
	// Resolve session ID
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	// Parse request body
	var req types.ChatReactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// React to message through Meow service
	h.logger.Infof("Reacting to message %s with emoji %s for %s from session %s", req.MessageID, req.Emoji, req.Phone, sessionID)

	err := h.meowService.ReactToMessage(c.Request.Context(), sessionID, req.Phone, req.MessageID, req.Emoji)
	if err != nil {
		h.logger.Errorf("Failed to react to message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to react to message", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Message reaction sent successfully",
	})
}

// Delete godoc
// @Summary Delete a message
// @Description Delete a message from a chat
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatDeleteRequest true "Delete request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/delete [post]
func (h *ChatHandler) Delete(c *gin.Context) {
	// Resolve session ID
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	// Parse request body
	var req types.ChatDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Delete message through Meow service
	h.logger.Infof("Deleting message %s for %s (forEveryone: %v) from session %s", req.MessageID, req.Phone, req.ForEveryone, sessionID)

	err := h.meowService.DeleteMessage(c.Request.Context(), sessionID, req.Phone, req.MessageID, req.ForEveryone)
	if err != nil {
		h.logger.Errorf("Failed to delete message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete message", err.Error())
		return
	}

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
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatEditRequest true "Edit request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/edit [post]
func (h *ChatHandler) Edit(c *gin.Context) {
	// Resolve session ID
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	// Parse request body
	var req types.ChatEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Edit message through Meow service
	h.logger.Infof("Editing message %s with new text for %s from session %s", req.MessageID, req.Phone, sessionID)

	sendResp, err := h.meowService.EditMessage(c.Request.Context(), sessionID, req.Phone, req.MessageID, req.NewText)
	if err != nil {
		h.logger.Errorf("Failed to edit message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to edit message", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Message edited successfully",
		Data:    sendResp,
	})
}

// DownloadImage godoc
// @Summary Download image media
// @Description Download image media from a message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/image [post]
func (h *ChatHandler) DownloadImage(c *gin.Context) {
	h.downloadMedia(c, "image")
}

// DownloadVideo godoc
// @Summary Download video media
// @Description Download video media from a message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/video [post]
func (h *ChatHandler) DownloadVideo(c *gin.Context) {
	h.downloadMedia(c, "video")
}

// DownloadAudio godoc
// @Summary Download audio media
// @Description Download audio media from a message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/audio [post]
func (h *ChatHandler) DownloadAudio(c *gin.Context) {
	h.downloadMedia(c, "audio")
}

// DownloadDocument godoc
// @Summary Download document media
// @Description Download document media from a message
// @Tags chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.ChatDownloadRequest true "Download request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/chat/download/document [post]
func (h *ChatHandler) DownloadDocument(c *gin.Context) {
	h.downloadMedia(c, "document")
}

// downloadMedia is a helper function to download media of any type
func (h *ChatHandler) downloadMedia(c *gin.Context, mediaType string) {
	// Resolve session ID
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	// Parse request body
	var req types.ChatDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Download media through Meow service
	h.logger.Infof("Downloading %s media for message %s from session %s", mediaType, req.MessageID, sessionID)

	mediaData, mimeType, err := h.meowService.DownloadMedia(c.Request.Context(), sessionID, req.MessageID)
	if err != nil {
		h.logger.Errorf("Failed to download %s media: %v", mediaType, err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to download media", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Media downloaded successfully",
		Data: map[string]interface{}{
			"mediaType": mediaType,
			"mimeType":  mimeType,
			"size":      len(mediaData),
			"data":      mediaData, // Base64 encoded data
		},
	})
}
