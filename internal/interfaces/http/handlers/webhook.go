package handlers

import (
	"net/http"

	"zpmeow/internal/domain/session"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook-related HTTP requests
type WebhookHandler struct {
	sessionService session.ApplicationSessionService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(sessionService session.ApplicationSessionService) *WebhookHandler {
	return &WebhookHandler{
		sessionService: sessionService,
	}
}

// RegisterWebhook handles registering a webhook
//
//	@Summary		Register webhook
//	@Description	Register a webhook URL to receive WhatsApp events
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			request		body		dto.RegisterWebhookRequest	true	"Register webhook request"
//	@Success		201			{object}	dto.RegisterWebhookResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhook/register [post]
func (h *WebhookHandler) RegisterWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Register webhook endpoint - implementation pending",
	})
}

// GetWebhook handles getting webhook information
//
//	@Summary		Get webhook information
//	@Description	Get information about registered webhooks for a session
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		200			{object}	dto.GetWebhookResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		404			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhook [get]
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get webhook endpoint - implementation pending",
	})
}

// Router compatibility methods
func (h *WebhookHandler) SetWebhook(c *gin.Context) { h.RegisterWebhook(c) }

// UpdateWebhook handles updating webhook configuration
//
//	@Summary		Update webhook
//	@Description	Update webhook URL, events, or status for a session
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			request		body		dto.UpdateWebhookRequest	true	"Update webhook request"
//	@Success		200			{object}	dto.UpdateWebhookResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		404			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhook [put]
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Update webhook - implementation pending"})
}

// DeleteWebhook handles deleting a webhook
//
//	@Summary		Delete webhook
//	@Description	Delete/unregister a webhook for a session
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		200			{object}	dto.DeleteWebhookResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		404			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhook [delete]
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Delete webhook - implementation pending"})
}

// ListWebhooks handles listing all webhooks
//
//	@Summary		List webhooks
//	@Description	Get a list of all registered webhooks for a session
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string	true	"Session ID"
//	@Success		200			{object}	dto.ListWebhooksResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhooks [get]
func (h *WebhookHandler) ListWebhooks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "List webhooks - implementation pending"})
}

// TestWebhook handles testing a webhook
//
//	@Summary		Test webhook
//	@Description	Send a test event to a webhook to verify it's working
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"
//	@Param			request		body		dto.TestWebhookRequest	true	"Test webhook request"
//	@Success		200			{object}	dto.TestWebhookResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		404			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhook/test [post]
func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Test webhook - implementation pending"})
}
