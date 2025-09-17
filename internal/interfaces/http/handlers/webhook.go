package handlers

import (
	"net/http"
	"time"

	"meow/internal/application"
	"meow/internal/infrastructure/wmeow"
	"meow/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook-related HTTP requests
type WebhookHandler struct {
	sessionService *application.SessionApp
	webhookApp     application.WebhookService
}

// resolveSessionID resolves session ID or name to actual session ID
func (h *WebhookHandler) resolveSessionID(c *gin.Context, sessionIDOrName string) (string, error) {
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

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(sessionService *application.SessionApp, webhookApp application.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		sessionService: sessionService,
		webhookApp:     webhookApp,
	}
}

// SetWebhook handles setting/configuring a webhook
//
//	@Summary		Set webhook
//	@Description	Set a webhook URL to receive meow events
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string						true	"Session ID"
//	@Param			request		body		dto.RegisterWebhookRequest	true	"Set webhook request"
//	@Success		201			{object}	dto.StandardWebhookCreateResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhook [post]
func (h *WebhookHandler) SetWebhook(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusNotFound,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "SESSION_NOT_FOUND",
				Message: "Session not found: " + err.Error(),
			},
		})
		return
	}

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

	// Validate and map event types using wmeow service
	validEvents := make([]string, 0)
	for _, event := range req.Events {
		if wmeow.IsValidEventType(event) {
			// Map to proper event type if needed
			mappedEvent := wmeow.MapEventType(event)
			validEvents = append(validEvents, mappedEvent)
		}
	}

	// If no valid events, return error
	if len(validEvents) == 0 {
		c.JSON(http.StatusBadRequest, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "INVALID_EVENTS",
				Message: "No valid events provided. Valid events include: message, status, connection, call, contact, group, presence, qr, pair, disconnect, error, all",
			},
		})
		return
	}

	// Set webhook for session using application service (includes validation)
	err = h.webhookApp.SetWebhook(c.Request.Context(), sessionID, req.URL, validEvents)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "Failed to register webhook: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.StandardWebhookCreateResponse{
		Status:  http.StatusCreated,
		Message: "Webhook registered successfully",
		Data: dto.StandardWebhookData{
			CreatedAt: time.Now(),
			Events:    validEvents,
			SessionID: sessionIDOrName,
			Status:    "active",
			URL:       req.URL,
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
//	@Success		200			{object}	dto.StandardWebhookResponse
//	@Failure		400			{object}	dto.WebhookResponse
//	@Failure		404			{object}	dto.WebhookResponse
//	@Failure		500			{object}	dto.WebhookResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhook [get]
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	sessionIDOrName := c.Param("sessionId")

	// Resolve session ID or name to actual session ID
	sessionID, err := h.resolveSessionID(c, sessionIDOrName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.WebhookResponse{
			Success: false,
			Code:    http.StatusNotFound,
			Data:    dto.WebhookResponseData{},
			Error: &dto.WebhookErrorResponse{
				Code:    "SESSION_NOT_FOUND",
				Message: "Session not found: " + err.Error(),
			},
		})
		return
	}

	webhook, err := h.webhookApp.GetWebhook(c.Request.Context(), sessionID)
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

	c.JSON(http.StatusOK, dto.StandardWebhookResponse{
		Status:  http.StatusOK,
		Message: "Webhook retrieved successfully",
		Data: dto.StandardWebhookData{
			CreatedAt: time.Now(),
			Events:    webhook.Events,
			SessionID: sessionIDOrName,
			Status:    "active",
			URL:       webhook.URL,
		},
	})
}

// ListEvents handles listing supported webhook event types
//
//	@Summary		List supported events
//	@Description	Get list of all supported webhook event types
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SupportedEventsResponse
//	@Security		ApiKeyAuth
//	@Router			/session/{sessionId}/webhooks/events [get]
func (h *WebhookHandler) ListEvents(c *gin.Context) {
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
