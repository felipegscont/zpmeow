package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"meow/internal/application"
	"meow/internal/infrastructure/wmeow"
	"meow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow"
)

// MessageHandler handles message-related HTTP requests (both sending and actions)
type MessageHandler struct {
	sessionService *application.SessionApp
	wmeowService   wmeow.Service
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(sessionService *application.SessionApp, wmeowService wmeow.Service) *MessageHandler {
	return &MessageHandler{
		sessionService: sessionService,
		wmeowService:   wmeowService,
	}
}

// resolveSessionID resolves session ID or name to actual session ID
func (h *MessageHandler) resolveSessionID(c *gin.Context, sessionIDOrName string) (string, error) {
	if h.sessionService == nil {
		// Fallback: assume it's already an ID
		return sessionIDOrName, nil
	}

	ctx := c.Request.Context()
	session, err := h.sessionService.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return "", err
	}

	return session.ID.Value(), nil
}

// decodeMediaData decodes base64 data URL to bytes or downloads from URL
func (h *MessageHandler) decodeMediaData(dataURL string) ([]byte, error) {
	// Check if it's a HTTP/HTTPS URL
	if strings.HasPrefix(dataURL, "http://") || strings.HasPrefix(dataURL, "https://") {
		// Download from URL
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

	// Check if it's a data URL (data:mime/type;base64,data)
	if strings.HasPrefix(dataURL, "data:") {
		// Find the comma that separates the header from the data
		commaIndex := strings.Index(dataURL, ",")
		if commaIndex == -1 {
			return nil, fmt.Errorf("invalid data URL format")
		}

		// Extract the base64 data part
		base64Data := dataURL[commaIndex+1:]

		// Decode base64
		data, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 data: %w", err)
		}

		return data, nil
	}

	// If it's neither URL nor data URL, assume it's raw base64
	data, err := base64.StdEncoding.DecodeString(dataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return data, nil
}

// ============================================================================
// MESSAGE SENDING METHODS
// ============================================================================

// SendText handles sending text messages
//
//	@Summary		Send text message
//	@Description	Send a text message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string				true	"Session ID"
//	@Param			request		body		dto.SendTextRequest	true	"Text message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/text [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Send text message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendTextMessage(ctx, sessionID, req.Phone, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_TEXT_FAILED",
			"Failed to send text message",
			err.Error(),
		))
		return
	}

	// Create standardized response
	response := dto.NewTextResponse(true, http.StatusOK, req.Phone, resp.ID, req.Body, true)
	c.JSON(http.StatusOK, response)
}

// SendMedia handles sending media messages
//
//	@Summary		Send media message
//	@Description	Send a media message (image, video, audio, document) to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendMediaRequest	true	"Media message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/media [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Decode media data
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

	// Send media message based on type using meow service
	ctx := c.Request.Context()
	var resp *whatsmeow.SendResponse
	switch req.MediaType {
	case "image":
		resp, err = h.wmeowService.SendImageMessage(ctx, sessionID, req.Phone, mediaData, req.Caption, "image/jpeg")
	case "audio":
		resp, err = h.wmeowService.SendAudioMessage(ctx, sessionID, req.Phone, mediaData, "audio/mpeg")
	case "video":
		resp, err = h.wmeowService.SendVideoMessage(ctx, sessionID, req.Phone, mediaData, req.Caption, "video/mp4")
	case "document":
		resp, err = h.wmeowService.SendDocumentMessage(ctx, sessionID, req.Phone, mediaData, "document", req.Caption, "application/octet-stream")
	case "sticker":
		resp, err = h.wmeowService.SendStickerMessage(ctx, sessionID, req.Phone, mediaData, "image/webp")
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

	// Create appropriate response based on media type
	var response *dto.MessageResponse
	switch req.MediaType {
	case "image":
		response = dto.NewImageResponse(true, http.StatusOK, req.Phone, resp.ID, "", req.Caption, true)
	case "audio":
		response = dto.NewAudioResponse(true, http.StatusOK, req.Phone, resp.ID, "", false, true)
	case "video":
		response = dto.NewVideoResponse(true, http.StatusOK, req.Phone, resp.ID, "", req.Caption, false, true)
	case "document":
		response = dto.NewDocumentResponse(true, http.StatusOK, req.Phone, resp.ID, "", "document", "application/octet-stream", true)
	case "sticker":
		response = dto.NewStickerResponse(true, http.StatusOK, req.Phone, resp.ID, "", true)
	}

	c.JSON(http.StatusOK, response)
}

// ============================================================================
// MESSAGE ACTION METHODS
// ============================================================================

// MarkAsRead handles marking messages as read
//
//	@Summary		Mark messages as read
//	@Description	Mark one or more messages as read in a chat
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.MarkAsReadRequest	true	"Mark as read request"
//	@Success		200			{object}	dto.MessageActionResponse
//	@Failure		400			{object}	dto.MessageActionResponse
//	@Failure		500			{object}	dto.MessageActionResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/markread [post]
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

	// Validate required fields
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

	// Mark messages as read via meow service
	ctx := c.Request.Context()
	err := h.wmeowService.MarkAsRead(ctx, sessionID, req.Phone, req.MessageIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageActionErrorResponse(
			http.StatusInternalServerError,
			"MARK_READ_FAILED",
			"Failed to mark messages as read",
			err.Error(),
		))
		return
	}

	response := dto.NewMessageActionSuccessResponse(req.Phone, "", "mark_read")
	c.JSON(http.StatusOK, response)
}

// ReactToMessage handles reacting to a message
//
//	@Summary		React to message
//	@Description	Add or remove reaction to a message
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			messageId	path		string						true	"Message ID"
//	@Param			request		body		dto.ReactToMessageRequest	true	"React request"
//	@Success		200			{object}	dto.MessageActionResponse
//	@Failure		400			{object}	dto.MessageActionResponse
//	@Failure		500			{object}	dto.MessageActionResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/{messageId}/reactions [post]
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

	// Validate required fields
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

	// React to message via meow service
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

// DeleteMessage handles deleting a message
//
//	@Summary		Delete message
//	@Description	Delete a message for everyone or just for me
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			messageId	path		string						true	"Message ID"
//	@Param			request		body		dto.DeleteMessageRequest	true	"Delete request"
//	@Success		200			{object}	dto.MessageActionResponse
//	@Failure		400			{object}	dto.MessageActionResponse
//	@Failure		500			{object}	dto.MessageActionResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/{messageId} [delete]
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

	// Validate required fields
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

	// Delete message via meow service
	ctx := c.Request.Context()
	err := h.wmeowService.DeleteMessage(ctx, sessionID, req.Phone, req.MessageID, req.ForEveryone)
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

// EditMessage handles editing a message
//
//	@Summary		Edit message
//	@Description	Edit the text content of a message
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			messageId	path		string					true	"Message ID"
//	@Param			request		body		dto.EditMessageRequest	true	"Edit request"
//	@Success		200			{object}	dto.MessageActionResponse
//	@Failure		400			{object}	dto.MessageActionResponse
//	@Failure		500			{object}	dto.MessageActionResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/{messageId}/edit [put]
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

	// Validate required fields
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

	// Edit message via meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.EditMessage(ctx, sessionID, req.Phone, req.MessageID, req.NewText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageActionErrorResponse(
			http.StatusInternalServerError,
			"EDIT_FAILED",
			"Failed to edit message",
			err.Error(),
		))
		return
	}

	response := dto.NewMessageActionSuccessResponse(req.Phone, resp.ID, "edit")
	c.JSON(http.StatusOK, response)
}

// ============================================================================
// ADDITIONAL MESSAGE SENDING METHODS
// ============================================================================

// SendLocation handles sending location messages
//
//	@Summary		Send location message
//	@Description	Send a location message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendLocationRequest	true	"Location message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/location [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Send location message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendLocationMessage(ctx, sessionID, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_LOCATION_FAILED",
			"Failed to send location message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewLocationResponse(true, http.StatusOK, req.Phone, resp.ID, req.Latitude, req.Longitude, req.Name, "", true)
	c.JSON(http.StatusOK, messageResponse)
}

// SendContact handles sending contact messages
//
//	@Summary		Send contact message
//	@Description	Send a contact message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendContactRequest	true	"Contact message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/contact [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Send contact message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendContactMessage(ctx, sessionID, req.Phone, req.ContactName, req.ContactPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_CONTACT_FAILED",
			"Failed to send contact message",
			err.Error(),
		))
		return
	}

	// Create vCard for response
	vcard := "BEGIN:VCARD\nVERSION:3.0\nFN:" + req.ContactName + "\nTEL:" + req.ContactPhone + "\nEND:VCARD"
	messageResponse := dto.NewContactResponse(true, http.StatusOK, req.Phone, resp.ID, req.ContactName, vcard, true)
	c.JSON(http.StatusOK, messageResponse)
}

// SendImage handles sending image messages
//
//	@Summary		Send image message
//	@Description	Send an image message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendImageRequest	true	"Image message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/image [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Decode image data
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

	// Send image message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendImageMessage(ctx, sessionID, req.Phone, imageData, req.Caption, "image/jpeg")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_IMAGE_FAILED",
			"Failed to send image message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewImageResponse(true, http.StatusOK, req.Phone, resp.ID, req.Image, req.Caption, true)
	c.JSON(http.StatusOK, messageResponse)
}

// SendAudio handles sending audio messages
//
//	@Summary		Send audio message
//	@Description	Send an audio message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendAudioRequest	true	"Audio message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/audio [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Decode audio data
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

	// Send audio message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendAudioMessage(ctx, sessionID, req.Phone, audioData, "audio/mpeg")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_AUDIO_FAILED",
			"Failed to send audio message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewAudioResponse(true, http.StatusOK, req.Phone, resp.ID, req.Audio, req.PTT, true)
	c.JSON(http.StatusOK, messageResponse)
}

// SendDocument handles sending document messages
//
//	@Summary		Send document message
//	@Description	Send a document message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendDocumentRequest	true	"Document message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/document [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Decode document data
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

	// Use provided filename and mimetype, or defaults
	filename := req.FileName
	if filename == "" {
		filename = "document"
	}
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Send document message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendDocumentMessage(ctx, sessionID, req.Phone, documentData, filename, "", mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_DOCUMENT_FAILED",
			"Failed to send document message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewDocumentResponse(true, http.StatusOK, req.Phone, resp.ID, req.Document, filename, mimeType, true)
	c.JSON(http.StatusOK, messageResponse)
}

// SendVideo handles sending video messages
//
//	@Summary		Send video message
//	@Description	Send a video message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendVideoRequest	true	"Video message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/video [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Decode video data
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

	// Send video message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendVideoMessage(ctx, sessionID, req.Phone, videoData, req.Caption, "video/mp4")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_VIDEO_FAILED",
			"Failed to send video message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewVideoResponse(true, http.StatusOK, req.Phone, resp.ID, req.Video, req.Caption, req.GifPlayback, true)
	c.JSON(http.StatusOK, messageResponse)
}

// SendSticker handles sending sticker messages
//
//	@Summary		Send sticker message
//	@Description	Send a sticker message to a meow contact
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.SendStickerRequest	true	"Sticker message request"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.MessageResponse
//	@Failure		500			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/sticker [post]
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

	// Resolve session ID or name to actual session ID
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

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Request validation failed",
			err.Error(),
		))
		return
	}

	// Decode sticker data
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

	// Send sticker message using meow service
	ctx := c.Request.Context()
	resp, err := h.wmeowService.SendStickerMessage(ctx, sessionID, req.Phone, stickerData, "image/webp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_STICKER_FAILED",
			"Failed to send sticker message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewStickerResponse(true, http.StatusOK, req.Phone, resp.ID, req.Sticker, true)
	c.JSON(http.StatusOK, messageResponse)
}

// ============================================================================
// PLACEHOLDER METHODS FOR FUTURE IMPLEMENTATION
// ============================================================================

// SendButton handles sending button messages (placeholder)
//
//	@Summary		Send button message
//	@Description	Send a button message to a meow contact (not yet implemented)
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		501			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/buttons [post]
func (h *MessageHandler) SendButton(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"Button messages not yet implemented",
		"This endpoint requires button message structure implementation",
	))
}

// SendList handles sending list messages (placeholder)
//
//	@Summary		Send list message
//	@Description	Send a list message to a meow contact (not yet implemented)
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		501			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/list [post]
func (h *MessageHandler) SendList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"List messages not yet implemented",
		"This endpoint requires list message structure implementation",
	))
}

// SendPoll handles sending poll messages (placeholder)
//
//	@Summary		Send poll message
//	@Description	Send a poll message to a meow contact (not yet implemented)
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		501			{object}	dto.MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/messages/poll [post]
func (h *MessageHandler) SendPoll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"Poll messages not yet implemented",
		"This endpoint requires poll message structure implementation",
	))
}
