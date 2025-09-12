package handlers

import (
	"net/http"
	"zpmeow/internal/application/dto/request"
	"zpmeow/internal/application/dto/response"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra/logger"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService domain.SessionService
	logger         logger.Logger
}

func NewSessionHandler(sessionService domain.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		logger:         logger.GetLogger().Sub("session-handler"),
	}
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req request.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	session, err := h.sessionService.CreateSession(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to create session"})
		return
	}

	resp := response.CreateSessionResponse{
		ID:     session.ID,
		Name:   session.Name,
		Status: string(session.Status),
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Session ID is required"})
		return
	}

	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Session not found"})
		return
	}

	resp := response.SessionInfoResponse{
		BaseSessionInfo: response.BaseSessionInfo{
			ID:     session.ID,
			Name:   session.Name,
			Status: string(session.Status),
		},
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SessionHandler) GetAllSessions(c *gin.Context) {
	sessions, err := h.sessionService.GetAllSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to get sessions"})
		return
	}

	var sessionResponses []response.SessionInfoResponse
	for _, session := range sessions {
		sessionResponses = append(sessionResponses, response.SessionInfoResponse{
			BaseSessionInfo: response.BaseSessionInfo{
				ID:     session.ID,
				Name:   session.Name,
				Status: string(session.Status),
			},
		})
	}

	resp := response.SessionListResponse{
		Sessions: sessionResponses,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Session ID is required"})
		return
	}

	err := h.sessionService.DeleteSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Session deleted successfully",
	})
}
