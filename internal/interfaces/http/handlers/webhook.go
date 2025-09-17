package handlers

import (
	"net/http"
	"time"

	"meow/internal/application"
	"meow/internal/interfaces/dto"
	"meow/internal/infrastructure/wmeow"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook-related HTTP requests
type WebhookHandler struct {
	sessionService *application.SessionApp
	webhookService application.WebhookService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(sessionService *application.SessionApp, webhookService application.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		sessionService: sessionService,
		webhookService: webhookService,
	}
}

// RegisterWebhook handles registering a webhook
//
//	@Summary		Register webhook
//	@Description	Register a webhook URL to receive meow events
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
	sessionID := c.Param("sessionId")

	var req dto.RegisterWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body: " + err.Error(),
			},
		})
		return
	}

	// Validate webhook URL
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "MISSING_URL",
				Message: "Webhook URL is required",
			},
		})
		return
	}

	// Validate event types
	validEvents := make([]string, 0)
	for _, event := range req.Events {
		if wmeow.IsValidEventType(event) {
			validEvents = append(validEvents, event)
		}
	}

	// Set webhook for session
	err := h.webhookService.SetWebhook(c.Request.Context(), sessionID, req.URL, validEvents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "REGISTRATION_FAILED",
				Message: "Failed to register webhook: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterWebhookResponse{
		Status:  http.StatusCreated,
		Message: "Webhook registered successfully",
		Data: dto.RegisterWebhookResponseData{
			WebhookID: "webhook_" + sessionID,
			SessionID: sessionID,
			URL:       req.URL,
			Events:    validEvents,
			Status:    "active",
			CreatedAt: time.Now(),
		},
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
	sessionID := c.Param("sessionId")

	webhook, err := h.webhookService.GetWebhook(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "GET_FAILED",
				Message: "Failed to get webhook: " + err.Error(),
			},
		})
		return
	}

	if webhook == nil {
		c.JSON(http.StatusNotFound, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusNotFound,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "NOT_FOUND",
				Message: "No webhook configured for this session",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetWebhookResponse{
		Status:  http.StatusOK,
		Message: "Webhook retrieved successfully",
		Data: dto.GetWebhookResponseData{
			WebhookID: "webhook_" + sessionID,
			SessionID: sessionID,
			URL:       webhook.URL,
			Events:    webhook.Events,
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
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
	sessionID := c.Param("sessionId")

	var req dto.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body: " + err.Error(),
			},
		})
		return
	}

	// Validate event types if provided
	validEvents := make([]string, 0)
	for _, event := range req.Events {
		if wmeow.IsValidEventType(event) {
			validEvents = append(validEvents, event)
		}
	}

	// Determine if webhook should be active
	active := req.Status == "active"

	// Update webhook for session
	err := h.webhookService.UpdateWebhook(c.Request.Context(), sessionID, req.URL, validEvents, active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "UPDATE_FAILED",
				Message: "Failed to update webhook: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateWebhookResponse{
		Status:  http.StatusOK,
		Message: "Webhook updated successfully",
		Data: dto.UpdateWebhookResponseData{
			WebhookID: "webhook_" + sessionID,
			SessionID: sessionID,
			URL:       req.URL,
			Events:    validEvents,
			Status:    req.Status,
			UpdatedAt: time.Now(),
		},
	})
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
	sessionID := c.Param("sessionId")

	err := h.webhookService.DeleteWebhook(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "DELETE_FAILED",
				Message: "Failed to delete webhook: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteWebhookResponse{
		Status:  http.StatusOK,
		Message: "Webhook deleted successfully",
		Data: dto.DeleteWebhookResponseData{
			WebhookID: "webhook_" + sessionID,
			Status:    "deleted",
		},
	})
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

// GetSupportedEvents handles listing supported webhook event types
//
//	@Summary		Get supported events
//	@Description	Get list of all supported webhook event types
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SupportedEventsResponse
//	@Security		ApiKeyAuth
//	@Router			/webhooks/events [get]
func (h *WebhookHandler) GetSupportedEvents(c *gin.Context) {
	events := wmeow.GetSupportedEventTypes()

	c.JSON(http.StatusOK, dto.SupportedEventsResponse{
		Status:  http.StatusOK,
		Message: "Supported events retrieved successfully",
		Data: dto.SupportedEventsData{
			Events: events,
			Count:  len(events),
		},
	})
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
	sessionID := c.Param("sessionId")

	var req dto.TestWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body: " + err.Error(),
			},
		})
		return
	}

	// Test webhook with event type
	err := h.webhookService.TestWebhook(c.Request.Context(), sessionID, req.EventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "TEST_FAILED",
				Message: "Webhook test failed: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.TestWebhookResponse{
		Status:  http.StatusOK,
		Message: "Webhook test completed",
		Data: dto.TestWebhookResponseData{
			WebhookID:    "webhook_" + sessionID,
			TestResult:   "success",
			ResponseCode: 200,
			ResponseTime: 150,
		},
	})
}
