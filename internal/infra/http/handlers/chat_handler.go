package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/domain"
	"zpmeow/internal/infra"
	httpUtils "zpmeow/internal/infra/http/utils"
	"zpmeow/internal/infra/logger"
)

type ChatHandler struct {
	sessionService domain.SessionService
	meowService    *infra.MeowServiceImpl
	logger         logger.Logger
}

func NewChatHandler(sessionService domain.SessionService, meowService *infra.MeowServiceImpl) *ChatHandler {
	return &ChatHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("chat-handler"),
	}
}

func (h *ChatHandler) SetPresence(c *gin.Context) {
	// Validate session
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req struct {
		Phone string `json:"phone" binding:"required"`
		State string `json:"state" binding:"required"`
		Media string `json:"media,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	h.logger.Infof("SetPresence request: phone=%s state=%s media=%s session=%s", req.Phone, req.State, req.Media, sessionID)
	c.JSON(http.StatusOK, gin.H{"message": "SetPresence - stub implementation"})
}

func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "MarkAsRead - stub implementation"})
}

func (h *ChatHandler) ReactToMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ReactToMessage - stub implementation"})
}

func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeleteMessage - stub implementation"})
}

func (h *ChatHandler) EditMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "EditMessage - stub implementation"})
}

func (h *ChatHandler) DownloadMedia(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DownloadMedia - stub implementation"})
}

// Aliases for router compatibility
func (h *ChatHandler) MarkRead(c *gin.Context) {
	h.MarkAsRead(c)
}

func (h *ChatHandler) React(c *gin.Context) {
	h.ReactToMessage(c)
}

func (h *ChatHandler) Delete(c *gin.Context) {
	h.DeleteMessage(c)
}

func (h *ChatHandler) Edit(c *gin.Context) {
	h.EditMessage(c)
}

// Specific download methods
func (h *ChatHandler) DownloadImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DownloadImage - stub implementation"})
}

func (h *ChatHandler) DownloadVideo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DownloadVideo - stub implementation"})
}

func (h *ChatHandler) DownloadAudio(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DownloadAudio - stub implementation"})
}

func (h *ChatHandler) DownloadDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DownloadDocument - stub implementation"})
}
