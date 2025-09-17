package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"meow/internal/application"
	"meow/internal/domain/session"
	"meow/internal/infrastructure/logging"
	"meow/internal/infrastructure/wmeow"
	"meow/internal/interfaces/dto"
)

// ============================================================================
// SESSION HANDLER - SELF-CONTAINED WITH INTERNAL HELPERS
// ============================================================================

type SessionHandler struct {
	sessionService *application.SessionApp
	wmeowService   wmeow.Service
	logger         logging.Logger
}

func NewSessionHandler(sessionService *application.SessionApp, wmeowService wmeow.Service) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		wmeowService:   wmeowService,
		logger:         logging.GetLogger().Sub("session-handler"),
	}
}

// ============================================================================
// REMOVED DUPLICATE DTOs - NOW USING CENTRALIZED DTOs FROM dto PACKAGE
// ============================================================================

// ============================================================================
// INTERNAL HELPER FUNCTIONS
// ============================================================================

// validateSessionID checks if a session ID or name is provided
// This function accepts both UUID and session name
func (h *SessionHandler) validateSessionID(c *gin.Context) (string, bool) {
	sessionIDOrName := c.Param("id")
	if sessionIDOrName == "" {
		h.sendErrorResponse(c, http.StatusBadRequest, "SESSION_ID_REQUIRED", "Session ID or name is required", "Missing session ID or name in path")
		return "", false
	}
	return sessionIDOrName, true
}

// bindAndValidateRequest binds JSON request and validates it
func (h *SessionHandler) bindAndValidateRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		h.sendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return false
	}
	return true
}

// sendSuccessResponse sends a standardized success response using centralized DTOs
func (h *SessionHandler) sendSuccessResponse(c *gin.Context, sessionID, action string, data interface{}) {
	response := &dto.SessionResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: dto.SessionData{
			SessionID: sessionID,
			Action:    action,
			Status:    "success",
			Timestamp: time.Now(),
		},
	}

	// Set specific data based on type
	switch v := data.(type) {
	case *dto.SessionInfo:
		response.Data.Session = v
	case []dto.SessionInfo:
		response.Data.Sessions = v
	case string:
		response.Data.QRCode = v
	}

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		h.logger.Errorf("Failed to marshal response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format response"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
}

// sendErrorResponse sends a standardized error response using centralized DTOs
func (h *SessionHandler) sendErrorResponse(c *gin.Context, status int, errorCode, message, details string) {
	response := &dto.SessionResponse{
		Success: false,
		Code:    status,
		Data: dto.SessionData{
			Status:    "error",
			Timestamp: time.Now(),
		},
		Error: &dto.SessionErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		h.logger.Errorf("Failed to marshal error response: %v", err)
		c.JSON(status, gin.H{"error": "Failed to format error response"})
		return
	}

	c.Data(status, "application/json", jsonBytes)
}

// convertToSessionInfo converts domain session to SessionInfo DTO
func (h *SessionHandler) convertToSessionInfo(session *session.Session) *dto.SessionInfo {
	sessionInfo := &dto.SessionInfo{
		ID:        session.ID.Value(),
		Name:      session.Name.Value(),
		Status:    string(session.Status),
		WaJID:     session.WaJID.Value(),
		ProxyURL:  session.ProxyURL.Value(),
		ApiKey:    session.ApiKey.Value(),
		CreatedAt: session.CreatedAt,
		// Note: WebhookURL and Events removed - now handled by separate webhook aggregate
	}

	return sessionInfo
}

// logOperation logs the start of an operation
func (h *SessionHandler) logOperation(operation, details string) {
	h.logger.Infof("%s: %s", operation, details)
}

// logSuccess logs successful completion of an operation
func (h *SessionHandler) logSuccess(operation, details string) {
	h.logger.Infof("%s completed successfully: %s", operation, details)
}

// logError logs an error during operation
func (h *SessionHandler) logError(operation string, err error) {
	h.logger.Errorf("Failed to %s: %v", operation, err)
}

// GetSessions godoc
//
//	@Summary		Get all sessions
//	@Description	Retrieves a list of all meow sessions
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SessionListResponse	"Sessions retrieved successfully"
//	@Failure		500	{object}	dto.SessionResponse		"Internal server error"
//	@Router			/sessions/list [get]
//	@Security		ApiKeyAuth
func (h *SessionHandler) GetSessions(c *gin.Context) {
	h.logOperation("Getting all sessions", "")

	sessions, err := h.sessionService.GetAllSessions(c.Request.Context())
	if err != nil {
		h.logError("get all sessions", err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "GET_SESSIONS_FAILED", "Failed to get sessions", err.Error())
		return
	}

	// Convert to SessionInfo DTOs
	sessionInfos := make([]dto.SessionInfo, len(sessions))
	for i, session := range sessions {
		sessionInfos[i] = *h.convertToSessionInfo(session)
	}

	h.sendSuccessResponse(c, "", "list", sessionInfos)
	h.logSuccess("Get all sessions", fmt.Sprintf("retrieved %d sessions", len(sessions)))
}

// GetSession godoc
//
//	@Summary		Get session information
//	@Description	Retrieves detailed information about a specific session by ID or name
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"Session ID or Name"
//	@Success		200	{object}	dto.SessionInfoResponse	"Session information retrieved successfully"
//	@Failure		400	{object}	dto.SessionResponse		"Invalid session ID or name"
//	@Failure		404	{object}	dto.SessionResponse		"Session not found"
//	@Failure		500	{object}	dto.SessionResponse		"Internal server error"
//	@Router			/sessions/{id}/info [get]
//	@Security		ApiKeyAuth
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID, ok := h.validateSessionID(c)
	if !ok {
		return
	}

	h.logOperation("Getting session", sessionID)

	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.logError("get session "+sessionID, err)
		h.sendErrorResponse(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found", err.Error())
		return
	}

	sessionInfo := h.convertToSessionInfo(session)
	h.sendSuccessResponse(c, sessionID, "get", sessionInfo)
	h.logSuccess("Get session", sessionID)
}

// CreateSession godoc
//
//	@Summary		Create a new meow session
//	@Description	Creates a new meow session and starts the client
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateSessionRequest	true	"Session creation request"
//	@Success		201		{object}	dto.CreateSessionResponse	"Session created successfully"
//	@Failure		400		{object}	dto.SessionResponse			"Invalid request body"
//	@Failure		500		{object}	dto.SessionResponse			"Internal server error"
//	@Router			/sessions/create [post]
//	@Security		ApiKeyAuth
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req dto.CreateSessionRequest
	if !h.bindAndValidateRequest(c, &req) {
		return
	}

	h.logOperation("Creating session", "name: "+req.Name)

	// Create session via application service
	session, err := h.sessionService.CreateSessionWithRequest(c.Request.Context(), &req)
	if err != nil {
		h.logError("create session", err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "CREATE_SESSION_FAILED", "Failed to create session", err.Error())
		return
	}

	sessionInfo := h.convertToSessionInfo(session)

	// Send 201 Created status for session creation
	response := &dto.CreateSessionResponse{
		Success: true,
		Code:    http.StatusCreated,
		Data: dto.SessionCreateData{
			SessionID: session.ID.Value(),
			Action:    "create",
			Status:    "success",
			Timestamp: time.Now(),
			Session:   sessionInfo,
		},
	}

	c.JSON(http.StatusCreated, response)
	h.logSuccess("Create session", session.ID.Value())
}

// DeleteSession godoc
//
//	@Summary		Delete a session
//	@Description	Deletes a meow session and stops the client
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Session ID"
//	@Success		200	{object}	dto.SessionResponse	"Session deleted successfully"
//	@Failure		400	{object}	dto.SessionResponse	"Invalid session ID"
//	@Failure		404	{object}	dto.SessionResponse	"Session not found"
//	@Failure		500	{object}	dto.SessionResponse	"Internal server error"
//	@Router			/sessions/{id}/delete [delete]
//	@Security		ApiKeyAuth
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID, ok := h.validateSessionID(c)
	if !ok {
		return
	}

	h.logOperation("Deleting session", sessionID)

	// Get session first to get the actual ID
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.logError("get session "+sessionID+" for deletion", err)
		h.sendErrorResponse(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found", err.Error())
		return
	}

	// Stop the meow client first
	if err := h.wmeowService.StopClient(session.ID.Value()); err != nil {
		h.logger.Errorf("Failed to stop client for session %s: %v", session.ID.Value(), err)
		// Continue with deletion even if stop fails
	}

	// Delete the session
	if err := h.sessionService.DeleteSession(c.Request.Context(), sessionID); err != nil {
		h.logError("delete session "+sessionID, err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "DELETE_SESSION_FAILED", "Failed to delete session", err.Error())
		return
	}

	h.sendSuccessResponse(c, sessionID, "delete", nil)
	h.logSuccess("Delete session", sessionID)
}

// REMOVED DUPLICATE GetAllSessions - Using GetSessions instead

// ConnectSession godoc
//
//	@Summary		Connect a session to meow
//	@Description	Initiates connection to meow and returns QR code if needed. Accepts session ID or name.
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Session ID or Name"
//	@Success		200	{object}	dto.ConnectSessionResponse	"Connection initiated successfully"
//	@Failure		400	{object}	dto.SessionResponse			"Invalid session ID or name"
//	@Failure		404	{object}	dto.SessionResponse			"Session not found"
//	@Failure		500	{object}	dto.SessionResponse			"Internal server error"
//	@Router			/sessions/{id}/connect [post]
//	@Security		ApiKeyAuth
func (h *SessionHandler) ConnectSession(c *gin.Context) {
	sessionID, ok := h.validateSessionID(c)
	if !ok {
		return
	}

	h.logOperation("Connecting session", sessionID)

	// Check if session exists
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.logError("get session "+sessionID+" for connection", err)
		h.sendErrorResponse(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found", err.Error())
		return
	}

	// Validate that no other session is using the same device (if this session has a device)
	if !session.WaJID.IsEmpty() {
		// Check if another session is already using this device
		existingSession, err := h.sessionService.GetSessionByDeviceJID(c.Request.Context(), session.WaJID.Value())
		if err == nil && existingSession.ID != session.ID {
			h.sendErrorResponse(c, http.StatusConflict, "DEVICE_ALREADY_IN_USE",
				fmt.Sprintf("Device %s is already in use by session %s (%s)", session.WaJID.Value(), existingSession.ID.Value(), existingSession.Name.Value()),
				"Each meow device can only be used by one session at a time")
			return
		}
	}

	// Start the client if not already started
	if err := h.wmeowService.StartClient(session.ID.Value()); err != nil {
		h.logError("start client for session "+session.ID.Value(), err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "START_CLIENT_FAILED", "Failed to start client", err.Error())
		return
	}

	// Get QR code if not connected
	var qrCode string
	isConnected := h.wmeowService.IsClientConnected(session.ID.Value())
	if !isConnected {
		qrCode, err = h.wmeowService.GetQRCode(session.ID.Value())
		if err != nil {
			h.logger.Errorf("Failed to get QR code for session %s: %v", session.ID.Value(), err)
		}
	}

	// Create connection info
	connectionInfo := &dto.SessionConnectionInfo{
		QRCode:      qrCode,
		Connected:   isConnected,
		IsConnected: isConnected,
	}

	// Create standardized response
	response := &dto.ConnectSessionResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: dto.SessionConnectData{
			SessionID:  session.ID.Value(),
			Action:     "connect",
			Status:     "success",
			Timestamp:  time.Now(),
			Session:    h.convertToSessionInfo(session),
			Connection: connectionInfo,
			QRCode:     qrCode,
		},
	}

	c.JSON(http.StatusOK, response)
	h.logSuccess("Connect session", sessionID)
}

// DisconnectSession godoc
//
//	@Summary		Disconnect a session from meow
//	@Description	Disconnects the session from meow without deleting it
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Session ID"
//	@Success		200	{object}	dto.SessionResponse	"Session disconnected successfully"
//	@Failure		400	{object}	dto.SessionResponse	"Invalid session ID"
//	@Failure		404	{object}	dto.SessionResponse	"Session not found"
//	@Failure		500	{object}	dto.SessionResponse	"Internal server error"
//	@Router			/sessions/{id}/disconnect [post]
//	@Security		ApiKeyAuth
func (h *SessionHandler) DisconnectSession(c *gin.Context) {
	sessionID, ok := h.validateSessionID(c)
	if !ok {
		return
	}

	h.logOperation("Disconnecting session", sessionID)

	// Check if session exists
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.logError("get session "+sessionID+" for disconnection", err)
		h.sendErrorResponse(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found", err.Error())
		return
	}

	// Stop the meow client
	if err := h.wmeowService.StopClient(session.ID.Value()); err != nil {
		h.logError("stop client for session "+session.ID.Value(), err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "STOP_CLIENT_FAILED", "Failed to disconnect session", err.Error())
		return
	}

	h.sendSuccessResponse(c, sessionID, "disconnect", nil)
	h.logSuccess("Disconnect session", sessionID)
}

// PairPhone godoc
//
//	@Summary		Pair phone number with session
//	@Description	Pairs a phone number with the session for meow connection
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Session ID"
//	@Param			request	body		dto.PairPhoneRequest	true	"Phone pairing request"
//	@Success		200		{object}	dto.PairPhoneResponse	"Phone pairing initiated successfully"
//	@Failure		400		{object}	dto.SessionResponse		"Invalid request body or session ID"
//	@Failure		404		{object}	dto.SessionResponse		"Session not found"
//	@Failure		500		{object}	dto.SessionResponse		"Internal server error"
//	@Router			/sessions/{id}/pair [post]
//	@Security		ApiKeyAuth
func (h *SessionHandler) PairPhone(c *gin.Context) {
	sessionID, ok := h.validateSessionID(c)
	if !ok {
		return
	}

	var req dto.PairPhoneRequest
	if !h.bindAndValidateRequest(c, &req) {
		return
	}

	h.logOperation("Pairing phone for session", fmt.Sprintf("session: %s, phone: %s", sessionID, req.PhoneNumber))

	// Check if session exists
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.logError("get session "+sessionID+" for phone pairing", err)
		h.sendErrorResponse(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found", err.Error())
		return
	}

	// Pair phone via MeowService
	pairCode, err := h.wmeowService.PairPhone(session.ID.Value(), req.PhoneNumber)
	if err != nil {
		h.logError("pair phone for session "+session.ID.Value(), err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "PHONE_PAIRING_FAILED", "Failed to pair phone", err.Error())
		return
	}

	// Create standardized response
	response := &dto.PairPhoneResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: dto.PairPhoneResponseData{
			SessionID:   sessionID,
			Action:      "pair",
			Status:      "success",
			Timestamp:   time.Now(),
			PhoneNumber: req.PhoneNumber,
			PairCode:    pairCode,
		},
	}

	c.JSON(http.StatusOK, response)
	h.logSuccess("Pair phone", fmt.Sprintf("session: %s, phone: %s", sessionID, req.PhoneNumber))
}

// GetSessionStatus godoc
//
//	@Summary		Get session status
//	@Description	Retrieves the current status and connection information of a session
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Session ID"
//	@Success		200	{object}	dto.SessionStatusResponse	"Session status retrieved successfully"
//	@Failure		400	{object}	dto.SessionResponse			"Invalid session ID"
//	@Failure		404	{object}	dto.SessionResponse			"Session not found"
//	@Failure		500	{object}	dto.SessionResponse			"Internal server error"
//	@Router			/sessions/{id}/status [get]
//	@Security		ApiKeyAuth
func (h *SessionHandler) GetSessionStatus(c *gin.Context) {
	sessionID, ok := h.validateSessionID(c)
	if !ok {
		return
	}

	h.logOperation("Getting status for session", sessionID)

	// Check if session exists
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.logError("get session "+sessionID+" for status", err)
		h.sendErrorResponse(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found", err.Error())
		return
	}

	// Get client status from MeowService
	clientStatus := h.wmeowService.GetClientStatus(session.ID.Value())
	isConnected := h.wmeowService.IsClientConnected(session.ID.Value())

	// Create standardized response
	response := &dto.SessionStatusResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: dto.SessionStatusResponseData{
			SessionID:     sessionID,
			Action:        "status",
			Status:        "success",
			Timestamp:     time.Now(),
			Name:          session.Name.Value(),
			SessionStatus: string(session.Status),
			WaJID:         session.WaJID.Value(),
			IsConnected:   isConnected,
			ClientStatus:  string(clientStatus), // Convert types.Status to string
			CreatedAt:     session.CreatedAt,
			UpdatedAt:     session.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
	h.logSuccess("Get session status", sessionID)
}

// UpdateSessionWebhook godoc
//
//	@Summary		Update session webhook URL
//	@Description	Updates the webhook URL and events for a session
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Session ID"
//	@Param			request	body		dto.UpdateWebhookRequest	true	"Webhook update request"
//	@Success		200		{object}	dto.SessionResponse			"Webhook updated successfully"
//	@Failure		400		{object}	dto.SessionResponse			"Invalid request"
//	@Failure		404		{object}	dto.SessionResponse			"Session not found"
//	@Failure		500		{object}	dto.SessionResponse			"Internal server error"
//	@Router			/sessions/{id}/webhook [put]
//	@Security		ApiKeyAuth
func (h *SessionHandler) UpdateSessionWebhook(c *gin.Context) {
	sessionID, ok := h.validateSessionID(c)
	if !ok {
		return
	}

	var req dto.UpdateWebhookRequest
	if !h.bindAndValidateRequest(c, &req) {
		return
	}

	h.logOperation("Updating webhook for session", fmt.Sprintf("session: %s, url: %s", sessionID, req.URL))

	// Check if session exists
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.logError("get session "+sessionID+" for webhook update", err)
		h.sendErrorResponse(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found", err.Error())
		return
	}

	// Update webhook in wmeowService
	if err := h.wmeowService.UpdateSessionWebhook(session.ID.Value(), req.URL); err != nil {
		h.logError("update webhook for session "+session.ID.Value(), err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "WEBHOOK_UPDATE_FAILED", "Failed to update webhook", err.Error())
		return
	}

	// Update events subscription
	if len(req.Events) > 0 {
		if err := h.wmeowService.UpdateSessionSubscriptions(session.ID.Value(), req.Events); err != nil {
			h.logError("update events subscription for session "+session.ID.Value(), err)
			h.sendErrorResponse(c, http.StatusInternalServerError, "EVENTS_UPDATE_FAILED", "Failed to update events subscription", err.Error())
			return
		}
	}

	h.sendSuccessResponse(c, sessionID, "webhook_update", nil)
	h.logSuccess("Update session webhook", sessionID)
}

// REMOVED: ValidateSessionsIntegrity endpoint - not needed
