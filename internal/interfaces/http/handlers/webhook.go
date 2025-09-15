package handlers

import (
	"net/http"

	"zpmeow/internal/shared/common"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook-related HTTP requests
type WebhookHandler struct {
	sessionService common.SessionService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(sessionService common.SessionService) *WebhookHandler {
	return &WebhookHandler{
		sessionService: sessionService,
	}
}

// RegisterWebhook handles registering a webhook
func (h *WebhookHandler) RegisterWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Register webhook endpoint - implementation pending",
	})
}

// GetWebhook handles getting webhook information
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get webhook endpoint - implementation pending",
	})
}

// Router compatibility methods
func (h *WebhookHandler) SetWebhook(c *gin.Context) { h.RegisterWebhook(c) }
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Update webhook - implementation pending"})
}
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Delete webhook - implementation pending"})
}
