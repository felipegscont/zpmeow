package handler

import (
	"net/http"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"
	"zpmeow/internal/utils"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook-related operations
type WebhookHandler struct {
	sessionService session.SessionService
	logger         logger.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(sessionService session.SessionService) *WebhookHandler {
	return &WebhookHandler{
		sessionService: sessionService,
		logger:         logger.GetLogger().Sub("webhook-handler"),
	}
}

// Helper function to resolve session ID from path parameter
func (h *WebhookHandler) resolveSessionID(c *gin.Context) (string, bool) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return "", false
	}

	// Check if session exists
	_, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Session not found", err.Error())
		return "", false
	}

	return sessionID, true
}

// Supported event types for webhooks
var supportedEventTypes = []string{
	"message",
	"status",
	"presence",
	"typing",
	"recording",
	"paused",
	"group",
	"call",
	"receipt",
	"reaction",
	"edit",
	"delete",
}

// Helper function to check if an event type is supported
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// @Summary Set webhook configuration
// @Description Configure webhook URL and event subscriptions for the session
// @Tags webhook
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.SetWebhookRequest true "Webhook configuration"
// @Success 200 {object} types.WebhookResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/webhook [post]
func (h *WebhookHandler) SetWebhook(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}
	var req types.SetWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate webhook URL
	if req.WebhookURL == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Webhook URL is required")
		return
	}

	// Validate and filter events
	var validEvents []string
	for _, event := range req.Events {
		if !contains(supportedEventTypes, event) {
			h.logger.Warnf("Event type '%s' is not supported and will be discarded", event)
			continue
		}
		validEvents = append(validEvents, event)
	}

	h.logger.Infof("Setting webhook URL: %s with events: %v for session %s", req.WebhookURL, validEvents, sessionID)

	// Get current session
	sess, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get session", err.Error())
		return
	}

	// Update webhook configuration
	sess.SetWebhook(req.WebhookURL, validEvents)

	// Save session
	err = h.sessionService.UpdateSession(c.Request.Context(), sess)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save webhook configuration", err.Error())
		return
	}

	response := types.WebhookResponse{
		Webhook: req.WebhookURL,
		Events:  validEvents,
		Active:  true,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// @Summary Get webhook configuration
// @Description Retrieve current webhook configuration for the session
// @Tags webhook
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Success 200 {object} types.WebhookResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/webhook [get]
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	h.logger.Infof("Getting webhook configuration for session %s", sessionID)

	// Get current session
	sess, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get session", err.Error())
		return
	}

	response := types.WebhookResponse{
		Webhook:   sess.WebhookURL,
		Subscribe: sess.Events,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// @Summary Update webhook configuration
// @Description Update existing webhook configuration
// @Tags webhook
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Param request body types.UpdateWebhookRequest true "Webhook update configuration"
// @Success 200 {object} types.WebhookResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/webhook [put]
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}
	var req types.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate webhook URL
	if req.WebhookURL == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Webhook URL is required")
		return
	}

	// Validate and filter events
	var validEvents []string
	for _, event := range req.Events {
		if !contains(supportedEventTypes, event) {
			h.logger.Warnf("Event type '%s' is not supported and will be discarded", event)
			continue
		}
		validEvents = append(validEvents, event)
	}

	h.logger.Infof("Updating webhook URL: %s with events: %v, active: %t for session %s", req.WebhookURL, validEvents, req.Active, sessionID)

	// Get current session
	sess, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get session", err.Error())
		return
	}

	// Update webhook configuration
	if req.Active {
		sess.SetWebhook(req.WebhookURL, validEvents)
	} else {
		sess.ClearWebhook()
	}

	// Save session
	err = h.sessionService.UpdateSession(c.Request.Context(), sess)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update webhook configuration", err.Error())
		return
	}

	response := types.WebhookResponse{
		Webhook: req.WebhookURL,
		Events:  validEvents,
		Active:  req.Active,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// @Summary Delete webhook configuration
// @Description Remove webhook configuration for the session
// @Tags webhook
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID or Name"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/webhook [delete]
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	sessionID, ok := h.resolveSessionID(c)
	if !ok {
		return
	}

	h.logger.Infof("Deleting webhook configuration for session %s", sessionID)

	// Get current session
	sess, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get session", err.Error())
		return
	}

	// Clear webhook configuration
	sess.ClearWebhook()

	// Save session
	err = h.sessionService.UpdateSession(c.Request.Context(), sess)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete webhook configuration", err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Webhook and events deleted successfully",
		Data: map[string]interface{}{
			"Details": "Webhook and events deleted successfully",
		},
	})
}
