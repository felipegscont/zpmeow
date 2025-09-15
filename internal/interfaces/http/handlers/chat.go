package handlers

import (
	"context"
	"net/http"
	"strconv"

	"zpmeow/internal/application"
	"zpmeow/internal/infrastructure/wameow"
	"zpmeow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	sessionService *application.SessionService
	wameowService  *wameow.MeowService
}

// NewChatHandler creates a new chat handler
func NewChatHandler(sessionService *application.SessionService, wameowService *wameow.MeowService) *ChatHandler {
	return &ChatHandler{
		sessionService: sessionService,
		wameowService:  wameowService,
	}
}

// SetPresence handles setting user presence in a chat
// @Summary Set presence in chat
// @Description Set user presence state in a specific chat (composing, available, etc.)
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.SetPresenceRequest true "Presence request"
// @Success 200 {object} dto.ChatResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/presence [post]
func (h *ChatHandler) SetPresence(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.SetPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.State == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_STATE",
			"State is required",
			"Valid states: available, unavailable, composing, recording, paused",
		))
		return
	}

	// Set presence via wameow service
	ctx := context.Background()
	err := h.wameowService.SetPresence(ctx, sessionID, req.Phone, req.State, req.Media)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewChatErrorResponse(
			http.StatusInternalServerError,
			"SET_PRESENCE_FAILED",
			"Failed to set presence",
			err.Error(),
		))
		return
	}

	response := dto.NewChatSuccessResponse(req.Phone, "", "set_presence")
	c.JSON(http.StatusOK, response)
}

// MarkAsRead handles marking messages as read
// @Summary Mark messages as read
// @Description Mark one or more messages as read in a chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.MarkAsReadRequest true "Mark as read request"
// @Success 200 {object} dto.ChatResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/markread [post]
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if len(req.MessageIDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_IDS",
			"At least one message ID is required",
			"",
		))
		return
	}

	// Mark messages as read via wameow service
	ctx := context.Background()
	err := h.wameowService.MarkAsRead(ctx, sessionID, req.Phone, req.MessageIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewChatErrorResponse(
			http.StatusInternalServerError,
			"MARK_READ_FAILED",
			"Failed to mark messages as read",
			err.Error(),
		))
		return
	}

	response := dto.NewChatSuccessResponse(req.Phone, "", "mark_read")
	c.JSON(http.StatusOK, response)
}

// ReactToMessage handles reacting to a message
// @Summary React to message
// @Description Add or remove reaction to a message
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.ReactToMessageRequest true "React request"
// @Success 200 {object} dto.ChatResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/react [post]
func (h *ChatHandler) ReactToMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.ReactToMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_ID",
			"Message ID is required",
			"",
		))
		return
	}

	if req.Emoji == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_EMOJI",
			"Emoji is required (use 'remove' to remove reaction)",
			"",
		))
		return
	}

	// React to message via wameow service
	ctx := context.Background()
	err := h.wameowService.ReactToMessage(ctx, sessionID, req.Phone, req.MessageID, req.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewChatErrorResponse(
			http.StatusInternalServerError,
			"REACT_FAILED",
			"Failed to react to message",
			err.Error(),
		))
		return
	}

	response := dto.NewChatSuccessResponse(req.Phone, req.MessageID, "react")
	c.JSON(http.StatusOK, response)
}

// DeleteMessage handles deleting a message
// @Summary Delete message
// @Description Delete a message for everyone or just for me
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.DeleteMessageRequest true "Delete request"
// @Success 200 {object} dto.ChatResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/delete [post]
func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_ID",
			"Message ID is required",
			"",
		))
		return
	}

	// Delete message via wameow service
	ctx := context.Background()
	err := h.wameowService.DeleteMessage(ctx, sessionID, req.Phone, req.MessageID, req.ForEveryone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewChatErrorResponse(
			http.StatusInternalServerError,
			"DELETE_FAILED",
			"Failed to delete message",
			err.Error(),
		))
		return
	}

	response := dto.NewChatSuccessResponse(req.Phone, req.MessageID, "delete")
	c.JSON(http.StatusOK, response)
}

// EditMessage handles editing a message
// @Summary Edit message
// @Description Edit the text content of a message
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.EditMessageRequest true "Edit request"
// @Success 200 {object} dto.ChatResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/edit [post]
func (h *ChatHandler) EditMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.EditMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_ID",
			"Message ID is required",
			"",
		))
		return
	}

	if req.NewText == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_NEW_TEXT",
			"New text is required",
			"",
		))
		return
	}

	// Edit message via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.EditMessage(ctx, sessionID, req.Phone, req.MessageID, req.NewText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewChatErrorResponse(
			http.StatusInternalServerError,
			"EDIT_FAILED",
			"Failed to edit message",
			err.Error(),
		))
		return
	}

	response := dto.NewChatSuccessResponse(req.Phone, resp.ID, "edit")
	c.JSON(http.StatusOK, response)
}

// DownloadImage handles downloading image media
// @Summary Download image
// @Description Download image media from a message
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.DownloadMediaRequest true "Download request"
// @Success 200 {object} dto.MediaDownloadResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/download/image [post]
func (h *ChatHandler) DownloadImage(c *gin.Context) {
	h.downloadMedia(c, "image")
}

// DownloadVideo handles downloading video media
// @Summary Download video
// @Description Download video media from a message
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.DownloadMediaRequest true "Download request"
// @Success 200 {object} dto.MediaDownloadResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/download/video [post]
func (h *ChatHandler) DownloadVideo(c *gin.Context) {
	h.downloadMedia(c, "video")
}

// DownloadAudio handles downloading audio media
// @Summary Download audio
// @Description Download audio media from a message
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.DownloadMediaRequest true "Download request"
// @Success 200 {object} dto.MediaDownloadResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/download/audio [post]
func (h *ChatHandler) DownloadAudio(c *gin.Context) {
	h.downloadMedia(c, "audio")
}

// DownloadDocument handles downloading document media
// @Summary Download document
// @Description Download document media from a message
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body dto.DownloadMediaRequest true "Download request"
// @Success 200 {object} dto.MediaDownloadResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/download/document [post]
func (h *ChatHandler) DownloadDocument(c *gin.Context) {
	h.downloadMedia(c, "document")
}

// downloadMedia is a helper method for downloading media files
func (h *ChatHandler) downloadMedia(c *gin.Context, mediaType string) {
	sessionID := c.Param("sessionId")

	var req dto.DownloadMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate required fields
	if req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_ID",
			"Message ID is required",
			"",
		))
		return
	}

	// Download media via wameow service
	ctx := context.Background()
	data, mimeType, err := h.wameowService.DownloadMedia(ctx, sessionID, req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewChatErrorResponse(
			http.StatusInternalServerError,
			"DOWNLOAD_FAILED",
			"Failed to download media",
			err.Error(),
		))
		return
	}

	// Create download response
	response := &dto.MediaDownloadResponse{
		Success:   true,
		Code:      http.StatusOK,
		MessageID: req.MessageID,
		MediaType: mediaType,
		MimeType:  mimeType,
		Data:      data,
		Size:      len(data),
	}

	c.JSON(http.StatusOK, response)
}

// GetChatHistory handles getting chat history
// @Summary Get chat history
// @Description Get chat history for a specific contact
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param phone query string true "Phone number"
// @Param limit query int false "Limit of messages (default: 50, max: 1000)"
// @Success 200 {object} dto.ChatHistoryResponse
// @Failure 400 {object} dto.ChatResponse
// @Failure 500 {object} dto.ChatResponse
// @Security ApiKeyAuth
// @Router /session/{sessionId}/chat/history [get]
func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	phone := c.Query("phone")
	limitStr := c.DefaultQuery("limit", "50")

	// Validate phone parameter
	if phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewChatErrorResponse(
			http.StatusBadRequest,
			"INVALID_LIMIT",
			"Invalid limit parameter",
			err.Error(),
		))
		return
	}

	// For now, return a simple response indicating the feature is available
	// but would need proper implementation with message storage/retrieval
	response := &dto.ChatHistoryResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: dto.ChatHistoryResponseData{
			Phone:    phone,
			Messages: []dto.ChatHistoryData{},
			Count:    0,
			Limit:    limit,
		},
	}

	c.JSON(http.StatusOK, response)
}
