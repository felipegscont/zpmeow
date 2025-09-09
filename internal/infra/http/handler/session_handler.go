package handler

import (
	"net/http"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/utils"

	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate request
	if !utils.IsValidSessionName(req.Name) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid session name")
		return
	}

	// Create session through service
	h.logger.Infof("Creating new session: %s", req.Name)
	sess, err := h.sessionService.CreateSession(c.Request.Context(), req.Name)
	if err != nil {
		h.logger.Errorf("Failed to create session %s: %v", req.Name, err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create session", err.Error())
		return
	}
	h.logger.Infof("Session created successfully: %s (ID: %s)", sess.Name, sess.ID)

	// Convert to response DTO
	response := session.CreateSessionResponse{
		ID:        sess.ID,
		Name:      sess.Name,
		Status:    string(sess.Status),
		CreatedAt: sess.CreatedAt,
		UpdatedAt: sess.UpdatedAt,
	}

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
	sessions, err := h.sessionService.GetAllSessions(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list sessions", err.Error())
		return
	}

	// Convert to response DTOs
	sessionResponses := make([]session.SessionInfoResponse, len(sessions))
	for i, sess := range sessions {
		sessionResponses[i] = session.SessionInfoResponse{
			ID:          sess.ID,
			Name:        sess.Name,
			WhatsAppJID: sess.WhatsAppJID,
			Status:      string(sess.Status),
			QRCode:      sess.QRCode,
			ProxyURL:    sess.ProxyURL,
			CreatedAt:   sess.CreatedAt,
			UpdatedAt:   sess.UpdatedAt,
		}
	}

	response := session.SessionListResponse{
		Sessions: sessionResponses,
		Total:    len(sessionResponses),
	}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	sess, err := h.sessionService.GetSession(c.Request.Context(), id)
	if err != nil {
		if err == session.ErrSessionNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Session not found")
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get session", err.Error())
		return
	}

	response := session.SessionInfoResponse{
		ID:          sess.ID,
		Name:        sess.Name,
		WhatsAppJID: sess.WhatsAppJID,
		Status:      string(sess.Status),
		QRCode:      sess.QRCode,
		ProxyURL:    sess.ProxyURL,
		CreatedAt:   sess.CreatedAt,
		UpdatedAt:   sess.UpdatedAt,
	}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	err := h.sessionService.DeleteSession(c.Request.Context(), id)
	if err != nil {
		if err == session.ErrSessionNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Session not found")
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete session", err.Error())
		return
	}

	utils.RespondNoContent(c)
}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	err := h.sessionService.ConnectSession(c.Request.Context(), id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to connect session", err.Error())
		return
	}

	response := session.MessageResponse{
		Message: "Connection process started. Check status and QR code endpoints.",
	}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	err := h.sessionService.DisconnectSession(c.Request.Context(), id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to logout session", err.Error())
		return
	}

	response := session.MessageResponse{
		Message: "Session logged out successfully.",
	}

	utils.RespondWithData(c, response)
}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	qrCode, err := h.sessionService.GetQRCode(c.Request.Context(), id)
	if err != nil {
		if err == session.ErrSessionNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Session not found")
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get QR code", err.Error())
		return
	}

	sess, _ := h.sessionService.GetSession(c.Request.Context(), id)
	response := session.QRCodeResponse{
		QRCode: qrCode,
		Status: string(sess.Status),
	}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req session.PairSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if !utils.IsValidPhoneNumber(req.PhoneNumber) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number")
		return
	}

	code, err := h.sessionService.PairWithPhone(c.Request.Context(), id, req.PhoneNumber)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to pair session", err.Error())
		return
	}

	response := session.PairSessionResponse{
		PairingCode: code,
	}

	utils.RespondWithData(c, response)
}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req session.ProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if !utils.IsValidProxyURL(req.ProxyURL) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid proxy URL")
		return
	}

	err := h.sessionService.SetProxy(c.Request.Context(), id, req.ProxyURL)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to set proxy", err.Error())
		return
	}

	response := session.ProxyResponse{
		ProxyURL: req.ProxyURL,
		Message:  "Proxy updated successfully.",
	}

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
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	proxyURL, err := h.sessionService.GetProxy(c.Request.Context(), id)
	if err != nil {
		if err == session.ErrSessionNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Session not found")
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get proxy", err.Error())
		return
	}

	response := session.ProxyResponse{
		ProxyURL: proxyURL,
		Message:  "Proxy configuration retrieved successfully.",
	}

	utils.RespondWithData(c, response)
}
