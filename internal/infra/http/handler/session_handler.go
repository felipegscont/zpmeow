package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/utils"
)

// SessionHandler handles HTTP requests for session operations
type SessionHandler struct {
	sessionService session.SessionService
	logger         logger.Logger
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(sessionService session.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		logger:         logger.GetLogger().Sub("session-handler"),
	}
}

// ============================================================================
// Helper Methods
// ============================================================================

// handleDomainError handles domain errors with appropriate HTTP status codes
func (h *SessionHandler) handleDomainError(c *gin.Context, err error, defaultMessage string) {
	statusCode, message := MapDomainError(err)

	if statusCode == http.StatusInternalServerError {
		// Log internal errors and use the provided default message
		h.logger.Errorf("%s: %v", defaultMessage, err)
		utils.RespondWithError(c, statusCode, defaultMessage, err.Error())
	} else {
		// Use the mapped message for known domain errors
		utils.RespondWithError(c, statusCode, message)
	}
}

// ============================================================================
// Session Lifecycle Handlers
// ============================================================================

// CreateSession godoc
// @Summary Create a new WhatsApp session
// @Description Creates a new WhatsApp session with the provided name
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body session.CreateSessionRequest true "Session creation request"
// @Success 201 {object} session.CreateSessionResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/create [post]
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req session.CreateSessionRequest
	if !ValidateAndBindJSON(c, &req) {
		return
	}

	// Create session through service (includes validation)
	h.logger.Infof("Creating new session: %s", req.Name)
	createdSession, err := h.sessionService.CreateSession(c.Request.Context(), req.Name)
	if err != nil {
		h.handleDomainError(c, err, "Failed to create session")
		return
	}
	h.logger.Infof("Session created successfully: %s (ID: %s)", createdSession.Name, createdSession.ID)

	// Convert to response DTO
	response := ToCreateSessionResponse(createdSession)
	utils.RespondCreated(c, response)
}

// ListSessions godoc
// @Summary List all WhatsApp sessions
// @Description Retrieves a list of all WhatsApp sessions in the system
// @Tags sessions
// @Accept json
// @Produce json
// @Success 200 {object} session.SessionListResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/list [get]
func (h *SessionHandler) ListSessions(c *gin.Context) {
	allSessions, err := h.sessionService.GetAllSessions(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list sessions", err.Error())
		return
	}

	// Convert to response DTO
	response := ToSessionListResponse(allSessions)
	utils.RespondWithData(c, response)
}

// GetSessionInfo godoc
// @Summary Get session information
// @Description Retrieves detailed information about a specific WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} session.SessionInfoResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/{id}/info [get]
func (h *SessionHandler) GetSessionInfo(c *gin.Context) {
	sessionID, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	sessionInfo, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.handleDomainError(c, err, "Failed to get session")
		return
	}

	response := ToSessionInfoResponse(sessionInfo)
	utils.RespondWithData(c, response)
}

// DeleteSession godoc
// @Summary Delete a session
// @Description Deletes a WhatsApp session and logs out the client if active
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 204 "Session deleted successfully"
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/{id}/delete [delete]
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	id, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	err := h.sessionService.DeleteSession(c.Request.Context(), id)
	if err != nil {
		h.handleDomainError(c, err, "Failed to delete session")
		return
	}

	utils.RespondNoContent(c)
}

// ============================================================================
// Connection Management Handlers
// ============================================================================

// ConnectSession godoc
// @Summary Connect a session to WhatsApp
// @Description Starts the connection process for a WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 202 {object} session.MessageResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/{id}/connect [post]
func (h *SessionHandler) ConnectSession(c *gin.Context) {
	id, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	err := h.sessionService.ConnectSession(c.Request.Context(), id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to connect session", err.Error())
		return
	}

	response := ToMessageResponse("Connection process started. Check status and QR code endpoints.")
	c.JSON(http.StatusAccepted, response)
}

// LogoutSession godoc
// @Summary Logout a session from WhatsApp
// @Description Logs out a WhatsApp session and disconnects the client
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} session.MessageResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/{id}/logout [post]
func (h *SessionHandler) LogoutSession(c *gin.Context) {
	id, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	err := h.sessionService.DisconnectSession(c.Request.Context(), id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to logout session", err.Error())
		return
	}

	response := ToMessageResponse("Session logged out successfully.")
	utils.RespondWithData(c, response)
}

// ============================================================================
// Authentication Handlers
// ============================================================================

// GetSessionQR godoc
// @Summary Get QR code for session
// @Description Retrieves the QR code for a WhatsApp session to scan with mobile device
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} session.QRCodeResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /sessions/{id}/qr [get]
func (h *SessionHandler) GetSessionQR(c *gin.Context) {
	id, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	qrCode, err := h.sessionService.GetQRCode(c.Request.Context(), id)
	if err != nil {
		h.handleDomainError(c, err, "Failed to get QR code")
		return
	}

	sess, _ := h.sessionService.GetSession(c.Request.Context(), id)
	response := ToQRCodeResponse(qrCode, sess)
	utils.RespondWithData(c, response)
}

// PairSession godoc
// @Summary Pair session with phone number
// @Description Pairs a WhatsApp session with a phone number using pairing code method
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body session.PairSessionRequest true "Phone number to pair"
// @Success 200 {object} session.PairSessionResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/{id}/pair [post]
func (h *SessionHandler) PairSession(c *gin.Context) {
	id, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	var req session.PairSessionRequest
	if !ValidateAndBindJSON(c, &req) {
		return
	}

	if err := session.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	code, err := h.sessionService.PairWithPhone(c.Request.Context(), id, req.PhoneNumber)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to pair session", err.Error())
		return
	}

	response := ToPairSessionResponse(code)
	utils.RespondWithData(c, response)
}

// ============================================================================
// Configuration Handlers
// ============================================================================

// SetProxy godoc
// @Summary Set proxy for session
// @Description Sets or updates the proxy configuration for a WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body session.ProxyRequest true "Proxy configuration"
// @Success 200 {object} session.ProxyResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /sessions/{id}/proxy/set [post]
func (h *SessionHandler) SetProxy(c *gin.Context) {
	id, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	var req session.ProxyRequest
	if !ValidateAndBindJSON(c, &req) {
		return
	}

	if err := session.ValidateProxyURL(req.ProxyURL); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.sessionService.SetProxy(c.Request.Context(), id, req.ProxyURL)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to set proxy", err.Error())
		return
	}

	response := ToProxyResponse(req.ProxyURL, "Proxy updated successfully.")
	utils.RespondWithData(c, response)
}

// GetProxy godoc
// @Summary Get proxy configuration for session
// @Description Retrieves the current proxy configuration for a WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} session.ProxyResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /sessions/{id}/proxy/find [get]
func (h *SessionHandler) GetProxy(c *gin.Context) {
	id, ok := ValidateSessionIDParam(c)
	if !ok {
		return
	}

	sess, err := h.sessionService.GetSession(c.Request.Context(), id)
	if err != nil {
		h.handleDomainError(c, err, "Failed to get session")
		return
	}

	response := ToProxyResponse(sess.ProxyURL, "Proxy configuration retrieved successfully.")
	utils.RespondWithData(c, response)
}
