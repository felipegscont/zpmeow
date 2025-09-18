package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"meow/internal/application"
	"meow/internal/infra/wmeow"
	"meow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	sessionService *application.SessionApp
	wmeowService   wmeow.Service
}

func NewMessageHandler(sessionService *application.SessionApp, wmeowService wmeow.Service) *MessageHandler {
	return &MessageHandler{
		sessionService: sessionService,
		wmeowService:   wmeowService,
	}
}

func (h *MessageHandler) resolveSessionID(c *gin.Context, sessionIDOrName string) (string, error) {
	if h.sessionService == nil {
		return sessionIDOrName, nil
	}

	ctx := c.Request.Context()
	session, err := h.sessionService.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return "", err
	}

	return session.ID().String(), nil
}

func (h *MessageHandler) decodeMediaData(dataURL string) ([]byte, error) {
	if strings.HasPrefix(dataURL, "http://") || strings.HasPrefix(dataURL, "https://") {
		resp, err := http.Get(dataURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download from URL: status %d", resp.StatusCode)
		}

		return io.ReadAll(resp.Body)
	}

	if strings.HasPrefix(dataURL, "data:") {
		commaIndex := strings.Index(dataURL, ",")
		if commaIndex == -1 {
			return nil, fmt.Errorf("invalid data URL format")
		}

		base64Data := dataURL[commaIndex+1:]

		data, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 data: %w", err)
		}

		return data, nil
	}

	data, err := base64.StdEncoding.DecodeString(dataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return data, nil
}

// @Summary		Send text message
// @Description	Send a text message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string				true	"Session ID"
// @Param			request		body		dto.SendTextRequest	true	"Text message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/text [post]
func (h *MessageHandler) SendText(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	err = h.wmeowService.SendTextMessage(ctx, sessionID, req.Phone, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_TEXT_FAILED",
			"Failed to send text message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	response := dto.NewTextResponse(true, http.StatusOK, req.Phone, messageID, req.Body, true)
	c.JSON(http.StatusOK, response)
}

// @Summary		Send media message
// @Description	Send a media message (image, video, audio, document) to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendMediaRequest	true	"Media message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/media [post]
func (h *MessageHandler) SendMedia(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	mediaData, err := h.decodeMediaData(req.MediaURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_MEDIA_DATA",
			"Failed to decode media data",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	// Convert media data to base64 string for now
	mediaDataStr := base64.StdEncoding.EncodeToString(mediaData)

	switch req.MediaType {
	case "image":
		err = h.wmeowService.SendImageMessage(ctx, sessionID, req.Phone, mediaDataStr)
	case "audio":
		err = h.wmeowService.SendAudioMessage(ctx, sessionID, req.Phone, mediaDataStr)
	case "video":
		err = h.wmeowService.SendVideoMessage(ctx, sessionID, req.Phone, mediaDataStr)
	case "document":
		err = h.wmeowService.SendDocumentMessage(ctx, sessionID, req.Phone, mediaDataStr)
	case "sticker":
		// For now, use image message for stickers
		err = h.wmeowService.SendImageMessage(ctx, sessionID, req.Phone, mediaDataStr)
	default:
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_MEDIA_TYPE",
			"Invalid media type",
			"Supported types: image, audio, video, document, sticker",
		))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_MEDIA_FAILED",
			"Failed to send media message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())

	var response *dto.MessageResponse
	switch req.MediaType {
	case "image":
		response = dto.NewImageResponse(true, http.StatusOK, req.Phone, messageID, "", req.Caption, true)
	case "audio":
		response = dto.NewAudioResponse(true, http.StatusOK, req.Phone, messageID, "", false, true)
	case "video":
		response = dto.NewVideoResponse(true, http.StatusOK, req.Phone, messageID, "", req.Caption, false, true)
	case "document":
		response = dto.NewDocumentResponse(true, http.StatusOK, req.Phone, messageID, "", "document", "application/octet-stream", true)
	case "sticker":
		response = dto.NewStickerResponse(true, http.StatusOK, req.Phone, messageID, "", true)
	}

	c.JSON(http.StatusOK, response)
}

// @Summary		Mark messages as read
// @Description	Mark one or more messages as read in a chat
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.MarkAsReadRequest	true	"Mark as read request"
// @Success		200			{object}	dto.MessageActionResponse
// @Failure		400			{object}	dto.MessageActionResponse
// @Failure		500			{object}	dto.MessageActionResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/markread [post]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if len(req.MessageIDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_IDS",
			"At least one message ID is required",
			"",
		))
		return
	}

	ctx := c.Request.Context()
	// Mark each message as read
	for _, messageID := range req.MessageIDs {
		if err := h.wmeowService.MarkAsRead(ctx, sessionID, req.Phone, messageID); err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewMessageActionErrorResponse(
				http.StatusInternalServerError,
				"MARK_READ_FAILED",
				"Failed to mark messages as read",
				err.Error(),
			))
			return
		}
	}

	response := dto.NewMessageActionSuccessResponse(req.Phone, "", "mark_read")
	c.JSON(http.StatusOK, response)
}

// @Summary		React to message
// @Description	Add or remove reaction to a message
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"
// @Param			messageId	path		string						true	"Message ID"
// @Param			request		body		dto.ReactToMessageRequest	true	"React request"
// @Success		200			{object}	dto.MessageActionResponse
// @Failure		400			{object}	dto.MessageActionResponse
// @Failure		500			{object}	dto.MessageActionResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/{messageId}/reactions [post]
func (h *MessageHandler) ReactToMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.ReactToMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_ID",
			"Message ID is required",
			"",
		))
		return
	}

	if req.Emoji == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_EMOJI",
			"Emoji is required (use 'remove' to remove reaction)",
			"",
		))
		return
	}

	ctx := c.Request.Context()
	err := h.wmeowService.ReactToMessage(ctx, sessionID, req.Phone, req.MessageID, req.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageActionErrorResponse(
			http.StatusInternalServerError,
			"REACT_FAILED",
			"Failed to react to message",
			err.Error(),
		))
		return
	}

	response := dto.NewMessageActionSuccessResponse(req.Phone, req.MessageID, "react")
	c.JSON(http.StatusOK, response)
}

// @Summary		Delete message
// @Description	Delete a message for everyone or just for me
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"
// @Param			messageId	path		string						true	"Message ID"
// @Param			request		body		dto.DeleteMessageRequest	true	"Delete request"
// @Success		200			{object}	dto.MessageActionResponse
// @Failure		400			{object}	dto.MessageActionResponse
// @Failure		500			{object}	dto.MessageActionResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/{messageId} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_ID",
			"Message ID is required",
			"",
		))
		return
	}

	ctx := c.Request.Context()
	err := h.wmeowService.DeleteMessage(ctx, sessionID, req.Phone, req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageActionErrorResponse(
			http.StatusInternalServerError,
			"DELETE_FAILED",
			"Failed to delete message",
			err.Error(),
		))
		return
	}

	response := dto.NewMessageActionSuccessResponse(req.Phone, req.MessageID, "delete")
	c.JSON(http.StatusOK, response)
}

// @Summary		Edit message
// @Description	Edit the text content of a message
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			messageId	path		string					true	"Message ID"
// @Param			request		body		dto.EditMessageRequest	true	"Edit request"
// @Success		200			{object}	dto.MessageActionResponse
// @Failure		400			{object}	dto.MessageActionResponse
// @Failure		500			{object}	dto.MessageActionResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/{messageId}/edit [put]
func (h *MessageHandler) EditMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req dto.EditMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_PHONE",
			"Phone number is required",
			"",
		))
		return
	}

	if req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_MESSAGE_ID",
			"Message ID is required",
			"",
		))
		return
	}

	if req.NewText == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageActionErrorResponse(
			http.StatusBadRequest,
			"MISSING_NEW_TEXT",
			"New text is required",
			"",
		))
		return
	}

	ctx := c.Request.Context()
	err := h.wmeowService.EditMessage(ctx, sessionID, req.Phone, req.MessageID, req.NewText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageActionErrorResponse(
			http.StatusInternalServerError,
			"EDIT_FAILED",
			"Failed to edit message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	response := dto.NewMessageActionSuccessResponse(req.Phone, messageID, "edit")
	c.JSON(http.StatusOK, response)
}

// @Summary		Send location message
// @Description	Send a location message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendLocationRequest	true	"Location message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/location [post]
func (h *MessageHandler) SendLocation(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	err = h.wmeowService.SendLocationMessage(ctx, sessionID, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_LOCATION_FAILED",
			"Failed to send location message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	messageResponse := dto.NewLocationResponse(true, http.StatusOK, req.Phone, messageID, req.Latitude, req.Longitude, req.Name, "", true)
	c.JSON(http.StatusOK, messageResponse)
}

// @Summary		Send contact message
// @Description	Send a contact message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendContactRequest	true	"Contact message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/contact [post]
func (h *MessageHandler) SendContact(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	err = h.wmeowService.SendContactMessage(ctx, sessionID, req.Phone, req.ContactName, req.ContactPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_CONTACT_FAILED",
			"Failed to send contact message",
			err.Error(),
		))
		return
	}

	vcard := "BEGIN:VCARD\nVERSION:3.0\nFN:" + req.ContactName + "\nTEL:" + req.ContactPhone + "\nEND:VCARD"
	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	messageResponse := dto.NewContactResponse(true, http.StatusOK, req.Phone, messageID, req.ContactName, vcard, true)
	c.JSON(http.StatusOK, messageResponse)
}

// @Summary		Send image message
// @Description	Send an image message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendImageRequest	true	"Image message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/image [post]
func (h *MessageHandler) SendImage(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	imageData, err := h.decodeMediaData(req.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_IMAGE_DATA",
			"Failed to decode image data",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	// Convert image data to base64 string
	imageDataStr := base64.StdEncoding.EncodeToString(imageData)
	err = h.wmeowService.SendImageMessage(ctx, sessionID, req.Phone, imageDataStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_IMAGE_FAILED",
			"Failed to send image message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	messageResponse := dto.NewImageResponse(true, http.StatusOK, req.Phone, messageID, req.Image, req.Caption, true)
	c.JSON(http.StatusOK, messageResponse)
}

// @Summary		Send audio message
// @Description	Send an audio message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendAudioRequest	true	"Audio message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/audio [post]
func (h *MessageHandler) SendAudio(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendAudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	audioData, err := h.decodeMediaData(req.Audio)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_AUDIO_DATA",
			"Failed to decode audio data",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	// Convert audio data to base64 string
	audioDataStr := base64.StdEncoding.EncodeToString(audioData)
	err = h.wmeowService.SendAudioMessage(ctx, sessionID, req.Phone, audioDataStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_AUDIO_FAILED",
			"Failed to send audio message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	messageResponse := dto.NewAudioResponse(true, http.StatusOK, req.Phone, messageID, req.Audio, req.PTT, true)
	c.JSON(http.StatusOK, messageResponse)
}

// @Summary		Send document message
// @Description	Send a document message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendDocumentRequest	true	"Document message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/document [post]
func (h *MessageHandler) SendDocument(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	documentData, err := h.decodeMediaData(req.Document)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_DOCUMENT_DATA",
			"Failed to decode document data",
			err.Error(),
		))
		return
	}

	filename := req.FileName
	if filename == "" {
		filename = "document"
	}
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	ctx := c.Request.Context()
	// Convert document data to base64 string
	documentDataStr := base64.StdEncoding.EncodeToString(documentData)
	err = h.wmeowService.SendDocumentMessage(ctx, sessionID, req.Phone, documentDataStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_DOCUMENT_FAILED",
			"Failed to send document message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	messageResponse := dto.NewDocumentResponse(true, http.StatusOK, req.Phone, messageID, req.Document, filename, mimeType, true)
	c.JSON(http.StatusOK, messageResponse)
}

// @Summary		Send video message
// @Description	Send a video message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendVideoRequest	true	"Video message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/video [post]
func (h *MessageHandler) SendVideo(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	videoData, err := h.decodeMediaData(req.Video)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_VIDEO_DATA",
			"Failed to decode video data",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	// Convert video data to base64 string
	videoDataStr := base64.StdEncoding.EncodeToString(videoData)
	err = h.wmeowService.SendVideoMessage(ctx, sessionID, req.Phone, videoDataStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_VIDEO_FAILED",
			"Failed to send video message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	messageResponse := dto.NewVideoResponse(true, http.StatusOK, req.Phone, messageID, req.Video, req.Caption, req.GifPlayback, true)
	c.JSON(http.StatusOK, messageResponse)
}

// @Summary		Send sticker message
// @Description	Send a sticker message to a meow contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendStickerRequest	true	"Sticker message request"
// @Success		200			{object}	dto.MessageResponse
// @Failure		400			{object}	dto.MessageResponse
// @Failure		500			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/sticker [post]
func (h *MessageHandler) SendSticker(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")
	if sessionIDOrName == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewMessageErrorResponse(
			http.StatusNotFound,
			"SESSION_NOT_FOUND",
			"Session not found",
			err.Error(),
		))
		return
	}

	var req dto.SendStickerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	stickerData, err := h.decodeMediaData(req.Sticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_STICKER_DATA",
			"Failed to decode sticker data",
			err.Error(),
		))
		return
	}

	ctx := c.Request.Context()
	// Convert sticker data to base64 string and use image message for now
	stickerDataStr := base64.StdEncoding.EncodeToString(stickerData)
	err = h.wmeowService.SendImageMessage(ctx, sessionID, req.Phone, stickerDataStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_STICKER_FAILED",
			"Failed to send sticker message",
			err.Error(),
		))
		return
	}

	// Generate a mock message ID for now
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	messageResponse := dto.NewStickerResponse(true, http.StatusOK, req.Phone, messageID, req.Sticker, true)
	c.JSON(http.StatusOK, messageResponse)
}

// @Summary		Send button message
// @Description	Send a button message to a meow contact (not yet implemented)
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string	true	"Session ID"
// @Success		501			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/buttons [post]
func (h *MessageHandler) SendButton(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"Button messages not yet implemented",
		"This endpoint requires button message structure implementation",
	))
}

// @Summary		Send list message
// @Description	Send a list message to a meow contact (not yet implemented)
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string	true	"Session ID"
// @Success		501			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/list [post]
func (h *MessageHandler) SendList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"List messages not yet implemented",
		"This endpoint requires list message structure implementation",
	))
}

// @Summary		Send poll message
// @Description	Send a poll message to a meow contact (not yet implemented)
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string	true	"Session ID"
// @Success		501			{object}	dto.MessageResponse
// @Security		ApiKeyAuth
// @Router			/session/{sessionId}/messages/poll [post]
func (h *MessageHandler) SendPoll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"Poll messages not yet implemented",
		"This endpoint requires poll message structure implementation",
	))
}
