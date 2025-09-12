package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra"
	"zpmeow/internal/infra/logger"
)

type SendHandler struct {
	sessionService domain.SessionService
	meowService    *infra.MeowServiceImpl
	logger         logger.Logger
}

func NewSendHandler(sessionService domain.SessionService, meowService *infra.MeowServiceImpl) *SendHandler {
	return &SendHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("send-handler"),
	}
}

func (h *SendHandler) SendText(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendText - stub implementation"})
}

func (h *SendHandler) SendMedia(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendMedia - stub implementation"})
}

func (h *SendHandler) SendImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendImage - stub implementation"})
}

func (h *SendHandler) SendAudio(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendAudio - stub implementation"})
}

func (h *SendHandler) SendDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendDocument - stub implementation"})
}

func (h *SendHandler) SendVideo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendVideo - stub implementation"})
}

func (h *SendHandler) SendLocation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendLocation - stub implementation"})
}

func (h *SendHandler) SendContact(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendContact - stub implementation"})
}

func (h *SendHandler) SendPoll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendPoll - stub implementation"})
}

func (h *SendHandler) SendSticker(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendSticker - stub implementation"})
}

func (h *SendHandler) SendButtons(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendButtons - stub implementation"})
}

func (h *SendHandler) SendList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SendList - stub implementation"})
}
