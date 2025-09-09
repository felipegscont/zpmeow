package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SessionHandler holds the whatsapp service and database connection
type SessionHandler struct {
	WhatsAppService *WhatsAppService
	DB              *sqlx.DB
}

// NewSessionHandler creates a new SessionHandler
func NewSessionHandler(waService *WhatsAppService, db *sqlx.DB) *SessionHandler {
	return &SessionHandler{
		WhatsAppService: waService,
		DB:              db,
	}
}

// CreateSession godoc
// @Summary Create a new WhatsApp session
// @Description Creates a new WhatsApp session with the provided name
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body CreateSessionRequest true "Session creation request"
// @Success 201 {object} CreateSessionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/create [post]
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	session := &Session{
		ID:     uuid.New().String(),
		Name:   req.Name,
		Status: "disconnected",
	}

	query := `INSERT INTO sessions (id, name, status) VALUES ($1, $2, $3) RETURNING created_at, updated_at`
	err := h.DB.QueryRowx(query, session.ID, session.Name, session.Status).StructScan(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create session"})
		return
	}

	response := CreateSessionResponse{
		ID:        session.ID,
		Name:      session.Name,
		Status:    session.Status,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// ListSessions godoc
// @Summary List all WhatsApp sessions
// @Description Retrieves a list of all WhatsApp sessions in the system
// @Tags sessions
// @Accept json
// @Produce json
// @Success 200 {object} SessionListResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/list [get]
func (h *SessionHandler) ListSessions(c *gin.Context) {
	var sessions []Session
	query := `SELECT id, name, COALESCE(whatsapp_jid, '') as whatsapp_jid, status, COALESCE(qr_code, '') as qr_code, COALESCE(proxy_url, '') as proxy_url, created_at, updated_at FROM sessions ORDER BY created_at DESC`
	if err := h.DB.Select(&sessions, query); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list sessions"})
		return
	}

	// Convert to response format
	sessionResponses := make([]SessionInfoResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = SessionInfoResponse{
			ID:          session.ID,
			Name:        session.Name,
			WhatsAppJID: session.WhatsAppJID,
			Status:      session.Status,
			QRCode:      session.QRCode,
			ProxyURL:    session.ProxyURL,
			CreatedAt:   session.CreatedAt,
			UpdatedAt:   session.UpdatedAt,
		}
	}

	response := SessionListResponse{
		Sessions: sessionResponses,
		Total:    len(sessionResponses),
	}

	c.JSON(http.StatusOK, response)
}

// GetSessionInfo godoc
// @Summary Get session information
// @Description Retrieves detailed information about a specific WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} SessionInfoResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/info [get]
func (h *SessionHandler) GetSessionInfo(c *gin.Context) {
	id := c.Param("id")
	var session Session
	query := `SELECT id, name, COALESCE(whatsapp_jid, '') as whatsapp_jid, status, COALESCE(qr_code, '') as qr_code, COALESCE(proxy_url, '') as proxy_url, created_at, updated_at FROM sessions WHERE id = $1`
	if err := h.DB.Get(&session, query, id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Session not found"})
		return
	}

	response := SessionInfoResponse{
		ID:          session.ID,
		Name:        session.Name,
		WhatsAppJID: session.WhatsAppJID,
		Status:      session.Status,
		QRCode:      session.QRCode,
		ProxyURL:    session.ProxyURL,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSession godoc
// @Summary Delete a session
// @Description Deletes a WhatsApp session and logs out the client if active
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 204 "Session deleted successfully"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/delete [delete]
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	id := c.Param("id")

	// First, logout the client if it's active
	_ = h.WhatsAppService.LogoutClient(id)

	query := `DELETE FROM sessions WHERE id = $1`
	result, err := h.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete session"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Session not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ConnectSession godoc
// @Summary Connect a session to WhatsApp
// @Description Starts the connection process for a WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 202 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/connect [post]
func (h *SessionHandler) ConnectSession(c *gin.Context) {
	id := c.Param("id")

	if err := h.WhatsAppService.StartClient(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, MessageResponse{
		Message: "Connection process started. Check status and QR code endpoints.",
	})
}

// LogoutSession godoc
// @Summary Logout a session from WhatsApp
// @Description Logs out a WhatsApp session and disconnects the client
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/logout [post]
func (h *SessionHandler) LogoutSession(c *gin.Context) {
	id := c.Param("id")

	if err := h.WhatsAppService.LogoutClient(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Session logged out successfully.",
	})
}

// GetSessionQR godoc
// @Summary Get QR code for session
// @Description Retrieves the QR code for a WhatsApp session to scan with mobile device
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} QRCodeResponse
// @Failure 404 {object} ErrorResponse
// @Router /sessions/{id}/qr [get]
func (h *SessionHandler) GetSessionQR(c *gin.Context) {
	id := c.Param("id")
	var session Session
	query := `SELECT COALESCE(qr_code, '') as qr_code, status FROM sessions WHERE id = $1`
	if err := h.DB.Get(&session, query, id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Session not found"})
		return
	}

	response := QRCodeResponse{
		QRCode: session.QRCode,
		Status: session.Status,
	}

	c.JSON(http.StatusOK, response)
}

// PairSession godoc
// @Summary Pair session with phone number
// @Description Pairs a WhatsApp session with a phone number using pairing code method
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body PairSessionRequest true "Phone number to pair"
// @Success 200 {object} PairSessionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/pair [post]
func (h *SessionHandler) PairSession(c *gin.Context) {
	id := c.Param("id")
	var req PairSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	code, err := h.WhatsAppService.PairPhone(id, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := PairSessionResponse{
		PairingCode: code,
	}

	c.JSON(http.StatusOK, response)
}

// SetProxy godoc
// @Summary Set proxy for session
// @Description Sets or updates the proxy configuration for a WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body ProxyRequest true "Proxy configuration"
// @Success 200 {object} ProxyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/proxy/set [post]
func (h *SessionHandler) SetProxy(c *gin.Context) {
	id := c.Param("id")
	var req ProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	query := `UPDATE sessions SET proxy_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := h.DB.Exec(query, req.ProxyURL, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to set proxy"})
		return
	}

	response := ProxyResponse{
		ProxyURL: req.ProxyURL,
		Message:  "Proxy updated successfully.",
	}

	c.JSON(http.StatusOK, response)
}

// GetProxy godoc
// @Summary Get proxy configuration for session
// @Description Retrieves the current proxy configuration for a WhatsApp session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} ProxyResponse
// @Failure 404 {object} ErrorResponse
// @Router /sessions/{id}/proxy/find [get]
func (h *SessionHandler) GetProxy(c *gin.Context) {
	id := c.Param("id")
	var session Session
	query := `SELECT COALESCE(proxy_url, '') as proxy_url FROM sessions WHERE id = $1`
	if err := h.DB.Get(&session, query, id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Session not found"})
		return
	}

	response := ProxyResponse{
		ProxyURL: session.ProxyURL,
		Message:  "Proxy configuration retrieved successfully.",
	}

	c.JSON(http.StatusOK, response)
}
