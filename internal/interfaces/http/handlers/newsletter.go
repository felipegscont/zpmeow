package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"meow/internal/application"
	"meow/internal/infrastructure/wmeow"
	"meow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	waTypes "go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// NewsletterHandler handles newsletter-related HTTP requests
type NewsletterHandler struct {
	sessionService *application.SessionApp
	wmeowService   wmeow.Service
}

// NewNewsletterHandler creates a new newsletter handler
func NewNewsletterHandler(sessionService *application.SessionApp, wmeowService wmeow.Service) *NewsletterHandler {
	return &NewsletterHandler{
		sessionService: sessionService,
		wmeowService:   wmeowService,
	}
}

// GetNewsletterMessageUpdates gets message updates from a newsletter
func (h *NewsletterHandler) GetNewsletterMessageUpdates(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Parse query parameters
	var req dto.GetNewsletterMessageUpdatesRequest
	req.JID = newsletterJID
	req.Count = 50 // Default count

	if countStr := c.Query("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
			req.Count = count
		}
	}
	if before := c.Query("before"); before != "" {
		req.Before = before
	}

	// Convert to whatsmeow params
	params, err := req.ToWhatsmeowParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Get newsletter message updates
	updates, err := h.wmeowService.GetNewsletterMessageUpdates(c.Request.Context(), sessionID, newsletterJID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to get newsletter message updates: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updates,
		"count":   len(updates),
	})
}

// MarkNewsletterViewed marks newsletter messages as viewed
func (h *NewsletterHandler) MarkNewsletterViewed(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Parse request body
	var req dto.MarkViewedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Convert server IDs to whatsmeow type
	serverIDs := make([]waTypes.MessageServerID, len(req.ServerIDs))
	for i, id := range req.ServerIDs {
		serverID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.StandardResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid server ID %s: %v", id, err),
			})
			return
		}
		serverIDs[i] = waTypes.MessageServerID(serverID)
	}

	// Mark messages as viewed
	err := h.wmeowService.NewsletterMarkViewed(c.Request.Context(), sessionID, newsletterJID, serverIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to mark messages as viewed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StandardResponse{
		Success: true,
		Message: "Messages marked as viewed successfully",
	})
}

// SendNewsletterReaction sends a reaction to a newsletter message
func (h *NewsletterHandler) SendNewsletterReaction(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Parse request body
	var req dto.SendReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.ServerID == "" || req.Reaction == "" || req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "server_id, reaction, and message_id are required",
		})
		return
	}

	// Convert server ID to int
	serverID, err := strconv.Atoi(req.ServerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid server ID: %v", err),
		})
		return
	}

	// Send reaction
	err = h.wmeowService.NewsletterSendReaction(
		c.Request.Context(),
		sessionID,
		newsletterJID,
		waTypes.MessageServerID(serverID),
		req.Reaction,
		waTypes.MessageID(req.MessageID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to send reaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StandardResponse{
		Success: true,
		Message: "Reaction sent successfully",
	})
}

// ToggleNewsletterMute toggles mute status for a newsletter
func (h *NewsletterHandler) ToggleNewsletterMute(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Parse request body
	var req dto.ToggleMuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Toggle mute
	err := h.wmeowService.NewsletterToggleMute(c.Request.Context(), sessionID, newsletterJID, req.Mute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to toggle mute: " + err.Error(),
		})
		return
	}

	action := "muted"
	if !req.Mute {
		action = "unmuted"
	}

	c.JSON(http.StatusOK, dto.StandardResponse{
		Success: true,
		Message: fmt.Sprintf("Newsletter %s successfully", action),
	})
}

// SubscribeLiveUpdates subscribes to live updates for a newsletter
func (h *NewsletterHandler) SubscribeLiveUpdates(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Subscribe to live updates
	err := h.wmeowService.NewsletterSubscribeLiveUpdates(c.Request.Context(), sessionID, newsletterJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to subscribe to live updates: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StandardResponse{
		Success: true,
		Message: "Successfully subscribed to live updates",
	})
}

// UploadNewsletterMedia uploads media for newsletter use
//
//	@Summary		Upload media for newsletter
//	@Description	Upload media files (image, video, audio, document) for use in newsletter messages. Returns MediaHandle required for sending media messages.
//	@Tags			Newsletters
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Param			file		formData	file	true	"Media file to upload"
//	@Param			media_type	formData	string	true	"Media type (image, video, audio, document)"
//	@Success		200			{object}	dto.UploadNewsletterMediaResponse	"Media uploaded successfully"
//	@Failure		400			{object}	dto.StandardResponse				"Bad request - Invalid file or parameters"
//	@Failure		404			{object}	dto.StandardResponse				"Session not found or not connected"
//	@Failure		500			{object}	dto.StandardResponse				"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/newsletter/upload [post]
func (h *NewsletterHandler) UploadNewsletterMedia(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Parse request body
	var req dto.UploadNewsletterMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Data == "" || req.MediaType == "" {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "data and media_type are required",
		})
		return
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid base64 data: " + err.Error(),
		})
		return
	}

	// Convert media type
	var mediaType whatsmeow.MediaType
	switch req.MediaType {
	case "image":
		mediaType = whatsmeow.MediaImage
	case "video":
		mediaType = whatsmeow.MediaVideo
	case "audio":
		mediaType = whatsmeow.MediaAudio
	case "document":
		mediaType = whatsmeow.MediaDocument
	default:
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid media type. Supported: image, video, audio, document",
		})
		return
	}

	// Upload media
	resp, err := h.wmeowService.UploadNewsletter(c.Request.Context(), sessionID, data, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to upload media: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Media uploaded successfully",
		"media_handle": resp.Handle,
		"url":          resp.URL,
		"direct_path":  resp.DirectPath,
	})
}

// GetNewsletterByInvite gets newsletter information by invite key
func (h *NewsletterHandler) GetNewsletterByInvite(c *gin.Context) {
	sessionID := c.Param("sessionId")
	inviteKey := c.Param("inviteKey")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.NewsletterInfoResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Validate invite key
	if inviteKey == "" {
		c.JSON(http.StatusBadRequest, dto.NewsletterInfoResponse{
			Success: false,
			Error:   "Invite key is required",
		})
		return
	}

	// Get newsletter info by invite
	info, err := h.wmeowService.GetNewsletterInfoWithInvite(c.Request.Context(), sessionID, inviteKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewsletterInfoResponse{
			Success: false,
			Error:   "Failed to get newsletter info: " + err.Error(),
		})
		return
	}

	// Convert to DTO
	result := &dto.NewsletterInfo{
		JID:         info.JID,
		Name:        info.Name,
		Description: info.Description,
		Picture:     info.Picture,
		Verified:    info.Verified,
		CreatedAt:   info.CreatedAt,
		UpdatedAt:   info.UpdatedAt,
		OwnerJID:    info.OwnerJID,
		Subscribers: info.Subscribers,
		Muted:       info.Muted,
	}

	c.JSON(http.StatusOK, dto.NewsletterInfoResponse{
		Success: true,
		Data:    result,
	})
}

// resolveSessionID resolves session ID or name to actual session ID
func (h *NewsletterHandler) resolveSessionID(c *gin.Context, sessionIDOrName string) (string, error) {
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

// CreateNewsletter handles creating a newsletter
//
//	@Summary		Create newsletter
//	@Description	Create a new meow newsletter
//	@Tags			Newsletters
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			request		body		dto.CreateNewsletterRequest	true	"Newsletter creation request"
//	@Success		201			{object}	dto.CreateNewsletterResponse	"Newsletter created successfully"
//	@Failure		400			{object}	dto.CreateNewsletterResponse	"Bad request"
//	@Failure		500			{object}	dto.CreateNewsletterResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/newsletter [post]
func (h *NewsletterHandler) CreateNewsletter(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Resolve session ID and validate session exists
	resolvedSessionID, err := h.resolveSessionID(c, sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateNewsletterResponse{
			Success: false,
			Error:   "Session not found: " + err.Error(),
		})
		return
	}

	// Validate session is connected
	if !h.wmeowService.IsClientConnected(resolvedSessionID) {
		c.JSON(http.StatusBadRequest, dto.CreateNewsletterResponse{
			Success: false,
			Error:   "Session not connected",
		})
		return
	}

	// Parse request body
	var req dto.CreateNewsletterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateNewsletterResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate request
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, dto.CreateNewsletterResponse{
			Success: false,
			Error:   "Newsletter name is required",
		})
		return
	}

	// Convert to whatsmeow params
	params, err := req.ToWhatsmeowParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateNewsletterResponse{
			Success: false,
			Error:   "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Create newsletter
	resp, err := h.wmeowService.CreateNewsletter(c.Request.Context(), resolvedSessionID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateNewsletterResponse{
			Success: false,
			Error:   "Failed to create newsletter: " + err.Error(),
		})
		return
	}

	// Build response
	result := &dto.CreateNewsletterResult{
		JID:         resp.ID.String(),
		ServerID:    "",         // NewsletterMetadata doesn't have ServerID
		Timestamp:   time.Now(), // Using current time as creation timestamp
		Name:        req.Name,
		Description: req.Description,
	}

	c.JSON(http.StatusCreated, dto.CreateNewsletterResponse{
		Success: true,
		Data:    result,
	})
}

// GetNewsletter handles getting newsletter information
//
//	@Summary		Get newsletter information
//	@Description	Get information about a specific newsletter
//	@Tags			Newsletters
//	@Accept			json
//	@Produce		json
//	@Param			sessionId		path		string						true	"Session ID"
//	@Param			newsletterId	path		string						true	"Newsletter JID"
//	@Success		200				{object}	dto.NewsletterInfoResponse	"Newsletter information retrieved"
//	@Failure		400				{object}	dto.NewsletterInfoResponse	"Bad request"
//	@Failure		404				{object}	dto.NewsletterInfoResponse	"Newsletter not found"
//	@Failure		500				{object}	dto.NewsletterInfoResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/newsletter/{newsletterId} [get]
func (h *NewsletterHandler) GetNewsletter(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.NewsletterInfoResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Validate newsletter JID
	if newsletterJID == "" {
		c.JSON(http.StatusBadRequest, dto.NewsletterInfoResponse{
			Success: false,
			Error:   "Newsletter JID is required",
		})
		return
	}

	// Get newsletter info
	info, err := h.wmeowService.GetNewsletterInfo(c.Request.Context(), sessionID, newsletterJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewsletterInfoResponse{
			Success: false,
			Error:   "Failed to get newsletter info: " + err.Error(),
		})
		return
	}

	// Convert to DTO
	result := &dto.NewsletterInfo{
		JID:         info.JID,
		Name:        info.Name,
		Description: info.Description,
		Picture:     info.Picture,
		Verified:    info.Verified,
		CreatedAt:   info.CreatedAt,
		UpdatedAt:   info.UpdatedAt,
		OwnerJID:    info.OwnerJID,
		Subscribers: info.Subscribers,
		Muted:       info.Muted,
	}

	c.JSON(http.StatusOK, dto.NewsletterInfoResponse{
		Success: true,
		Data:    result,
	})
}

// UpdateNewsletter - REMOVED: Not supported by whatsmeow
// This functionality does not exist in the whatsmeow library
// Reference: Analysis of whatsmeow documentation and source code

// DeleteNewsletter - REMOVED: Not supported by whatsmeow
// This functionality does not exist in the whatsmeow library
// Reference: Analysis of whatsmeow documentation and source code

// ListNewsletters handles listing newsletters
//
//	@Summary		List newsletters
//	@Description	Get a list of all subscribed newsletters for a session
//	@Tags			Newsletters
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Success		200			{object}	dto.NewsletterListResponse	"Newsletters retrieved successfully"
//	@Failure		400			{object}	dto.NewsletterListResponse	"Bad request"
//	@Failure		500			{object}	dto.NewsletterListResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/newsletters [get]
func (h *NewsletterHandler) ListNewsletters(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.NewsletterListResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Get subscribed newsletters
	newsletters, err := h.wmeowService.GetSubscribedNewsletters(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewsletterListResponse{
			Success: false,
			Error:   "Failed to get newsletters: " + err.Error(),
		})
		return
	}

	// Convert to DTOs
	result := make([]dto.NewsletterInfo, len(newsletters))
	for i, newsletter := range newsletters {
		result[i] = dto.NewsletterInfo{
			JID:         newsletter.JID,
			Name:        newsletter.Name,
			Description: newsletter.Description,
			Picture:     newsletter.Picture,
			Verified:    newsletter.Verified,
			CreatedAt:   newsletter.CreatedAt,
			UpdatedAt:   newsletter.UpdatedAt,
			OwnerJID:    newsletter.OwnerJID,
			Subscribers: newsletter.Subscribers,
			Muted:       newsletter.Muted,
		}
	}

	// Build response
	listResult := &dto.NewsletterList{
		Newsletters: result,
		Count:       len(result),
		Total:       len(result),
	}

	c.JSON(http.StatusOK, dto.NewsletterListResponse{
		Success: true,
		Data:    listResult,
	})
}

// SubscribeToNewsletter handles subscribing to a newsletter
//
//	@Summary		Subscribe to newsletter
//	@Description	Subscribe to a newsletter to receive updates
//	@Tags			Newsletters
//	@Accept			json
//	@Produce		json
//	@Param			sessionId		path		string				true	"Session ID"
//	@Param			newsletterId	path		string				true	"Newsletter JID"
//	@Success		200				{object}	dto.StandardResponse	"Subscribed successfully"
//	@Failure		400				{object}	dto.StandardResponse	"Bad request"
//	@Failure		404				{object}	dto.StandardResponse	"Newsletter not found"
//	@Failure		500				{object}	dto.StandardResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/newsletter/{newsletterId}/subscribe [post]
func (h *NewsletterHandler) SubscribeToNewsletter(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session
	_, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found: " + err.Error(),
		})
		return
	}

	// Validate newsletter JID
	if newsletterJID == "" {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Newsletter JID is required",
		})
		return
	}

	// Follow newsletter
	err = h.wmeowService.FollowNewsletter(c.Request.Context(), sessionID, newsletterJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to subscribe to newsletter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StandardResponse{
		Success: true,
		Message: "Successfully subscribed to newsletter",
	})
}

// UnsubscribeFromNewsletter handles unsubscribing from a newsletter
//
//	@Summary		Unsubscribe from newsletter
//	@Description	Unsubscribe from a newsletter to stop receiving updates
//	@Tags			Newsletters
//	@Accept			json
//	@Produce		json
//	@Param			sessionId		path		string				true	"Session ID"
//	@Param			newsletterId	path		string				true	"Newsletter JID"
//	@Success		200				{object}	dto.StandardResponse	"Unsubscribed successfully"
//	@Failure		400				{object}	dto.StandardResponse	"Bad request"
//	@Failure		404				{object}	dto.StandardResponse	"Newsletter not found"
//	@Failure		500				{object}	dto.StandardResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/newsletter/{newsletterId}/unsubscribe [post]
func (h *NewsletterHandler) UnsubscribeFromNewsletter(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session
	_, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found: " + err.Error(),
		})
		return
	}

	// Validate newsletter JID
	if newsletterJID == "" {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Newsletter JID is required",
		})
		return
	}

	// Unfollow newsletter
	err = h.wmeowService.UnfollowNewsletter(c.Request.Context(), sessionID, newsletterJID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to unsubscribe from newsletter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StandardResponse{
		Success: true,
		Message: "Successfully unsubscribed from newsletter",
	})
}

// SendNewsletterMessage handles sending a newsletter message
// ✅ FULLY IMPLEMENTED - Uses SendMessage + MediaHandle approach based on whatsmeow issues #481, #697, #498
//
//	@Summary		Send newsletter message
//	@Description	Send a text or media message to newsletter subscribers. For media messages, MediaHandle is required.
//	@Tags			Newsletters
//	@Accept			json
//	@Produce		json
//	@Param			sessionId		path		string					true	"Session ID"
//	@Param			newsletterId	path		string					true	"Newsletter ID (format: {id}@newsletter)"
//	@Param			request			body		dto.SendNewsletterMessageRequest	true	"Message content (text or media)"
//	@Success		200				{object}	dto.SendNewsletterMessageResponse	"Message sent successfully"
//	@Failure		400				{object}	dto.SendNewsletterMessageResponse	"Bad request - Invalid parameters or missing MediaHandle for media"
//	@Failure		404				{object}	dto.SendNewsletterMessageResponse	"Newsletter not found or session not connected"
//	@Failure		500				{object}	dto.SendNewsletterMessageResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/newsletter/{newsletterId}/send [post]
func (h *NewsletterHandler) SendNewsletterMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.SendNewsletterMessageResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Parse request body
	var req dto.SendNewsletterMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.SendNewsletterMessageResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate that we have either content or media
	if req.Content == "" && req.MediaHandle == "" {
		c.JSON(http.StatusBadRequest, dto.SendNewsletterMessageResponse{
			Success: false,
			Error:   "Either content or media_handle is required",
		})
		return
	}

	// Build message based on content type
	var message *waE2E.Message

	if req.MediaHandle != "" {
		// Media message - determine type based on media_type
		switch req.MediaType {
		case "image":
			message = &waE2E.Message{
				ImageMessage: &waE2E.ImageMessage{
					Caption: proto.String(req.Caption),
				},
			}
		case "video":
			message = &waE2E.Message{
				VideoMessage: &waE2E.VideoMessage{
					Caption: proto.String(req.Caption),
				},
			}
		case "audio":
			message = &waE2E.Message{
				AudioMessage: &waE2E.AudioMessage{},
			}
		case "document":
			message = &waE2E.Message{
				DocumentMessage: &waE2E.DocumentMessage{
					Caption: proto.String(req.Caption),
				},
			}
		default:
			// Default to image if media_type not specified
			message = &waE2E.Message{
				ImageMessage: &waE2E.ImageMessage{
					Caption: proto.String(req.Caption),
				},
			}
		}
	} else {
		// Text message
		message = &waE2E.Message{
			Conversation: proto.String(req.Content),
		}
	}

	// Send message to newsletter using the new implementation
	resp, err := h.wmeowService.SendNewsletterMessage(c.Request.Context(), sessionID, newsletterJID, message, req.MediaHandle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.SendNewsletterMessageResponse{
			Success: false,
			Error:   "Failed to send newsletter message: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SendNewsletterMessageResponse{
		Success:   true,
		MessageID: string(resp.ID),
		ServerID:  fmt.Sprintf("%d", resp.ServerID),
		Timestamp: resp.Timestamp.Format(time.RFC3339),
	})
}

// GetNewsletterMessages gets messages from a newsletter
func (h *NewsletterHandler) GetNewsletterMessages(c *gin.Context) {
	sessionID := c.Param("sessionId")
	newsletterJID := c.Param("newsletterId")

	// Validate session exists and is connected
	if !h.wmeowService.IsClientConnected(sessionID) {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Session not found or not connected",
		})
		return
	}

	// Parse query parameters
	var req dto.GetNewsletterMessagesRequest
	req.JID = newsletterJID
	req.Count = 50 // Default count

	if countStr := c.Query("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
			req.Count = count
		}
	}
	if before := c.Query("before"); before != "" {
		req.Before = before
	}

	// Convert to whatsmeow params
	params, err := req.ToWhatsmeowParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.StandardResponse{
			Success: false,
			Error:   "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Get newsletter messages
	messages, err := h.wmeowService.GetNewsletterMessages(c.Request.Context(), sessionID, newsletterJID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.StandardResponse{
			Success: false,
			Error:   "Failed to get newsletter messages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    messages,
		"count":   len(messages),
	})
}

// GetNewsletterSubscribers - REMOVED: Not supported by whatsmeow
// This functionality does not exist in the whatsmeow library
// Reference: Analysis of whatsmeow documentation and source code

// GetNewsletterMetrics - REMOVED: Not supported by whatsmeow
// This functionality does not exist in the whatsmeow library
// Reference: Analysis of whatsmeow documentation and source code
