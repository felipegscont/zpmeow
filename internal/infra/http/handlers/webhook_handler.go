package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/application"
	"zpmeow/internal/application/services"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra"
	"zpmeow/internal/infra/logger"
)

type WebhookHandler struct {
	sessionService    domain.SessionService
	validationService *services.ApplicationValidationService
	logger            logger.Logger
}

func NewWebhookHandler(sessionService domain.SessionService) *WebhookHandler {
	return &WebhookHandler{
		sessionService:    sessionService,
		validationService: services.DefaultValidationService,
		logger:            logger.GetLogger().Sub("webhook-handler"),
	}
}

// GetSupportedEvents returns all supported webhook events
func (h *WebhookHandler) GetSupportedEvents(c *gin.Context) {
	supportedEvents := []string{
		// Message Events
		infra.EventMessage,
		infra.EventUndecryptableMessage,
		infra.EventReceipt,
		infra.EventMediaRetry,
		infra.EventReadReceipt,

		// Connection Events
		infra.EventConnected,
		infra.EventDisconnected,
		infra.EventConnectFailure,
		infra.EventKeepAliveRestored,
		infra.EventKeepAliveTimeout,
		infra.EventLoggedOut,
		infra.EventClientOutdated,
		infra.EventTemporaryBan,
		infra.EventStreamError,
		infra.EventStreamReplaced,
		infra.EventPairSuccess,
		infra.EventPairError,
		infra.EventQR,
		infra.EventQRScannedWithoutMultidevice,

		// Group Events
		infra.EventGroupInfo,
		infra.EventJoinedGroup,
		infra.EventPicture,
		infra.EventBlocklistChange,
		infra.EventBlocklist,

		// Call Events
		infra.EventCallOffer,
		infra.EventCallAccept,
		infra.EventCallTerminate,
		infra.EventCallOfferNotice,
		infra.EventCallRelayLatency,

		// Presence Events
		infra.EventPresence,
		infra.EventChatPresence,

		// Privacy and Settings
		infra.EventPrivacySettings,
		infra.EventPushNameSetting,
		infra.EventUserAbout,

		// Sync Events
		infra.EventAppState,
		infra.EventAppStateSyncComplete,
		infra.EventHistorySync,
		infra.EventOfflineSyncCompleted,
		infra.EventOfflineSyncPreview,

		// Newsletter Events
		infra.EventNewsletterJoin,
		infra.EventNewsletterLeave,
		infra.EventNewsletterMuteChange,
		infra.EventNewsletterLiveUpdate,

		// Other Events
		infra.EventIdentityChange,
		infra.EventCATRefreshError,
		infra.EventFBMessage,

		// Special
		infra.EventAll,
	}

	c.JSON(http.StatusOK, gin.H{
		"events": supportedEvents,
		"count":  len(supportedEvents),
	})
}

// validateEvents validates that all provided events are supported
func (h *WebhookHandler) validateEvents(events []string) []string {
	if len(events) == 0 {
		return nil // Empty means all events
	}

	supportedEventsMap := map[string]bool{
		infra.EventMessage:                       true,
		infra.EventUndecryptableMessage:          true,
		infra.EventReceipt:                       true,
		infra.EventMediaRetry:                    true,
		infra.EventReadReceipt:                   true,
		infra.EventConnected:                     true,
		infra.EventDisconnected:                  true,
		infra.EventConnectFailure:                true,
		infra.EventKeepAliveRestored:             true,
		infra.EventKeepAliveTimeout:              true,
		infra.EventLoggedOut:                     true,
		infra.EventClientOutdated:                true,
		infra.EventTemporaryBan:                  true,
		infra.EventStreamError:                   true,
		infra.EventStreamReplaced:                true,
		infra.EventPairSuccess:                   true,
		infra.EventPairError:                     true,
		infra.EventQR:                            true,
		infra.EventQRScannedWithoutMultidevice:   true,
		infra.EventGroupInfo:                     true,
		infra.EventJoinedGroup:                   true,
		infra.EventPicture:                       true,
		infra.EventBlocklistChange:               true,
		infra.EventBlocklist:                     true,
		infra.EventCallOffer:                     true,
		infra.EventCallAccept:                    true,
		infra.EventCallTerminate:                 true,
		infra.EventCallOfferNotice:               true,
		infra.EventCallRelayLatency:              true,
		infra.EventPresence:                      true,
		infra.EventChatPresence:                  true,
		infra.EventPrivacySettings:               true,
		infra.EventPushNameSetting:               true,
		infra.EventUserAbout:                     true,
		infra.EventAppState:                      true,
		infra.EventAppStateSyncComplete:          true,
		infra.EventHistorySync:                   true,
		infra.EventOfflineSyncCompleted:          true,
		infra.EventOfflineSyncPreview:            true,
		infra.EventNewsletterJoin:                true,
		infra.EventNewsletterLeave:               true,
		infra.EventNewsletterMuteChange:          true,
		infra.EventNewsletterLiveUpdate:          true,
		infra.EventIdentityChange:                true,
		infra.EventCATRefreshError:               true,
		infra.EventFBMessage:                     true,
		infra.EventAll:                           true,
	}

	var invalidEvents []string
	for _, event := range events {
		if !supportedEventsMap[event] {
			invalidEvents = append(invalidEvents, event)
		}
	}

	return invalidEvents
}

func (h *WebhookHandler) SetWebhook(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	var req application.SetWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate webhook URL
	if err := h.validationService.ValidateWebhookURL(req.WebhookURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate events
	if invalidEvents := h.validateEvents(req.Events); len(invalidEvents) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Invalid events provided",
			"invalid_events": invalidEvents,
		})
		return
	}

	// TODO: Implement actual webhook setting logic
	h.logger.Infof("Setting webhook for session %s: %s (events: %v)", sessionID, req.WebhookURL, req.Events)

	response := application.WebhookResponse{
		Webhook: req.WebhookURL,
		Events:  req.Events,
		Active:  true,
	}

	c.JSON(http.StatusOK, response)
}

func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// TODO: Implement actual webhook retrieval logic
	h.logger.Infof("Getting webhook for session %s", sessionID)

	// For now, return a placeholder response
	response := application.WebhookResponse{
		Webhook: "",
		Events:  []string{},
		Active:  false,
	}

	c.JSON(http.StatusOK, response)
}

func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	var req application.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate webhook URL
	if err := h.validationService.ValidateWebhookURL(req.WebhookURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate events
	if invalidEvents := h.validateEvents(req.Events); len(invalidEvents) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Invalid events provided",
			"invalid_events": invalidEvents,
		})
		return
	}

	// TODO: Implement actual webhook update logic
	h.logger.Infof("Updating webhook for session %s: %s (events: %v, active: %v)",
		sessionID, req.WebhookURL, req.Events, req.Active)

	response := application.WebhookResponse{
		Webhook: req.WebhookURL,
		Events:  req.Events,
		Active:  req.Active,
	}

	c.JSON(http.StatusOK, response)
}

func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// TODO: Implement actual webhook deletion logic
	h.logger.Infof("Deleting webhook for session %s", sessionID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook deleted successfully",
		"session": sessionID,
	})
}

// normalizeEvents normalizes event names and removes duplicates
func (h *WebhookHandler) normalizeEvents(events []string) []string {
	if len(events) == 0 {
		return events
	}

	seen := make(map[string]bool)
	var normalized []string

	for _, event := range events {
		// Trim whitespace and normalize case
		event = strings.TrimSpace(event)
		if event == "" {
			continue
		}

		// Avoid duplicates
		if !seen[event] {
			seen[event] = true
			normalized = append(normalized, event)
		}
	}

	return normalized
}
