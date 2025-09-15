package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"zpmeow/internal/application"
	"zpmeow/internal/infrastructure/wameow"
	"zpmeow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	sessionService *application.SessionService
	wameowService  *wameow.MeowService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(sessionService *application.SessionService, wameowService *wameow.MeowService) *MessageHandler {
	return &MessageHandler{
		sessionService: sessionService,
		wameowService:  wameowService,
	}
}

// resolveSessionID resolves session ID or name to actual session ID
func (h *MessageHandler) resolveSessionID(c *gin.Context, sessionIDOrName string) (string, error) {
	if h.sessionService == nil {
		// Fallback: assume it's already an ID
		return sessionIDOrName, nil
	}

	ctx := context.Background()
	session, err := h.sessionService.GetSession(ctx, sessionIDOrName)
	if err != nil {
		return "", err
	}

	return session.ID, nil
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
			return nil, fmt.Errorf("failed to download from URL: HTTP %d", resp.StatusCode)
		}

		// Read the response body
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return data, nil
	}

	// Check if it's a data URL (data:mime/type;base64,data)
	if strings.HasPrefix(dataURL, "data:") {
		// Find the comma that separates the header from the data
		commaIndex := strings.Index(dataURL, ",")
		if commaIndex == -1 {
			return nil, http.ErrNotSupported
		}

		// Extract the base64 data part
		base64Data := dataURL[commaIndex+1:]

		// Decode base64
		return base64.StdEncoding.DecodeString(base64Data)
	}

	// If it's not a data URL or HTTP URL, assume it's raw base64
	return base64.StdEncoding.DecodeString(dataURL)
}

// SendTextMessage handles sending text messages
func (h *MessageHandler) SendTextMessage(c *gin.Context) {
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
			"Invalid request format",
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

	// Send message via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendTextMessage(ctx, sessionID, req.Phone, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
			"Failed to send text message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewTextResponse(true, http.StatusOK, req.Phone, resp.ID, req.Body, true)
	c.JSON(http.StatusOK, messageResponse)
}

// SendMediaMessage handles sending media messages
func (h *MessageHandler) SendMediaMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	var req dto.SendMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request format",
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

	// TODO: Download media from URL and send via appropriate method
	// For now, return a placeholder response
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"Media message sending not yet implemented",
		"This endpoint requires media download and processing implementation",
	))
}

// Router compatibility methods
func (h *MessageHandler) SendText(c *gin.Context) { h.SendTextMessage(c) }
func (h *MessageHandler) SendMedia(c *gin.Context) { h.SendMediaMessage(c) }
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
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Basic validation
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Phone number is required",
			"Phone field cannot be empty",
		))
		return
	}

	// Send location via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendLocationMessage(ctx, sessionID, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
			"Failed to send location message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewLocationResponse(true, http.StatusOK, req.Phone, resp.ID, req.Latitude, req.Longitude, req.Name, "", true)
	c.JSON(http.StatusOK, messageResponse)
}
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
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Basic validation
	if req.Phone == "" || req.ContactName == "" || req.ContactPhone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Required fields missing",
			"Phone, contact_name, and contact_phone are required",
		))
		return
	}

	// Send contact via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendContactMessage(ctx, sessionID, req.Phone, req.ContactName, req.ContactPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
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
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Basic validation
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Phone number is required",
			"Phone field cannot be empty",
		))
		return
	}

	if req.Image == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Image is required",
			"Image field cannot be empty",
		))
		return
	}

	// Convert base64 to bytes
	imageData, err := h.decodeMediaData(req.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_MEDIA",
			"Invalid image data",
			err.Error(),
		))
		return
	}

	// Send image via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendImageMessage(ctx, sessionID, req.Phone, imageData, req.Caption, "image/jpeg")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
			"Failed to send image message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewImageResponse(true, http.StatusOK, req.Phone, resp.ID, req.Image, req.Caption, true)
	c.JSON(http.StatusOK, messageResponse)
}

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
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Basic validation
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Phone number is required",
			"Phone field cannot be empty",
		))
		return
	}

	if req.Audio == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Audio is required",
			"Audio field cannot be empty",
		))
		return
	}

	// Convert base64 to bytes (simplified - in production you'd handle URLs too)
	audioData, err := h.decodeMediaData(req.Audio)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_MEDIA",
			"Invalid audio data",
			err.Error(),
		))
		return
	}

	// Send audio via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendAudioMessage(ctx, sessionID, req.Phone, audioData, "audio/mpeg")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
			"Failed to send audio message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewAudioResponse(true, http.StatusOK, req.Phone, resp.ID, req.Audio, req.PTT, true)
	c.JSON(http.StatusOK, messageResponse)
}

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
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Basic validation
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Phone number is required",
			"Phone field cannot be empty",
		))
		return
	}

	if req.Document == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Document is required",
			"Document field cannot be empty",
		))
		return
	}

	// Convert base64 to bytes
	documentData, err := h.decodeMediaData(req.Document)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_MEDIA",
			"Invalid document data",
			err.Error(),
		))
		return
	}

	// Default filename if not provided
	filename := req.FileName
	if filename == "" {
		filename = "document"
	}

	// Default mime type if not provided
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Send document via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendDocumentMessage(ctx, sessionID, req.Phone, documentData, filename, "", mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
			"Failed to send document message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewDocumentResponse(true, http.StatusOK, req.Phone, resp.ID, req.Document, filename, mimeType, true)
	c.JSON(http.StatusOK, messageResponse)
}

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
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Basic validation
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Phone number is required",
			"Phone field cannot be empty",
		))
		return
	}

	if req.Video == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Video is required",
			"Video field cannot be empty",
		))
		return
	}

	// Convert base64 to bytes
	videoData, err := h.decodeMediaData(req.Video)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_MEDIA",
			"Invalid video data",
			err.Error(),
		))
		return
	}

	// Send video via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendVideoMessage(ctx, sessionID, req.Phone, videoData, req.Caption, "video/mp4")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
			"Failed to send video message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewVideoResponse(true, http.StatusOK, req.Phone, resp.ID, req.Video, req.Caption, req.GifPlayback, true)
	c.JSON(http.StatusOK, messageResponse)
}

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
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Basic validation
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Phone number is required",
			"Phone field cannot be empty",
		))
		return
	}

	if req.Sticker == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Sticker is required",
			"Sticker field cannot be empty",
		))
		return
	}

	// Convert base64 to bytes
	stickerData, err := h.decodeMediaData(req.Sticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"INVALID_MEDIA",
			"Invalid sticker data",
			err.Error(),
		))
		return
	}

	// Send sticker via wameow service
	ctx := context.Background()
	resp, err := h.wameowService.SendStickerMessage(ctx, sessionID, req.Phone, stickerData, "image/webp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageErrorResponse(
			http.StatusInternalServerError,
			"SEND_FAILED",
			"Failed to send sticker message",
			err.Error(),
		))
		return
	}

	// Create successful response
	messageResponse := dto.NewStickerResponse(true, http.StatusOK, req.Phone, resp.ID, req.Sticker, true)
	c.JSON(http.StatusOK, messageResponse)
}


func (h *MessageHandler) SendButton(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Button messages require complex structure - not implemented yet
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"Button messages not yet implemented",
		"This endpoint requires button message structure implementation",
	))
}

func (h *MessageHandler) SendList(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// List messages require complex structure - not implemented yet
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"List messages not yet implemented",
		"This endpoint requires list message structure implementation",
	))
}

func (h *MessageHandler) SendPoll(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, dto.NewMessageErrorResponse(
			http.StatusBadRequest,
			"MISSING_SESSION_ID",
			"Session ID is required",
			"Session ID must be provided in the URL path",
		))
		return
	}

	// Poll messages require complex structure - not implemented yet
	c.JSON(http.StatusNotImplemented, dto.NewMessageErrorResponse(
		http.StatusNotImplemented,
		"NOT_IMPLEMENTED",
		"Poll messages not yet implemented",
		"This endpoint requires poll message structure implementation",
	))
}

// Alias for SendHandler compatibility
type SendHandler = MessageHandler

// NewSendHandler creates a new send handler (alias for MessageHandler)
func NewSendHandler(sessionService *application.SessionService, wameowService *wameow.MeowService) *SendHandler {
	return NewMessageHandler(sessionService, wameowService)
}
