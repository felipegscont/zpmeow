package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra/logger"
)

type WebhookHandler struct {
	sessionService domain.SessionService
	logger         logger.Logger
}

func NewWebhookHandler(sessionService domain.SessionService) *WebhookHandler {
	return &WebhookHandler{
		sessionService: sessionService,
		logger:         logger.GetLogger().Sub("webhook-handler"),
	}
}

func (h *WebhookHandler) SetWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SetWebhook - stub implementation"})
}

func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetWebhook - stub implementation"})
}

func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateWebhook - stub implementation"})
}

func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeleteWebhook - stub implementation"})
}
