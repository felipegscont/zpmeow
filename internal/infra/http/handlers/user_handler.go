package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra"
	"zpmeow/internal/infra/logger"
)

type UserHandler struct {
	sessionService domain.SessionService
	meowService    *infra.MeowServiceImpl
	logger         logger.Logger
}

func NewUserHandler(sessionService domain.SessionService, meowService *infra.MeowServiceImpl) *UserHandler {
	return &UserHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("user-handler"),
	}
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetUserInfo - stub implementation"})
}

func (h *UserHandler) CheckUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "CheckUser - stub implementation"})
}

func (h *UserHandler) SetPresence(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "SetPresence - stub implementation"})
}

func (h *UserHandler) GetAvatar(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetAvatar - stub implementation"})
}

func (h *UserHandler) GetContacts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetContacts - stub implementation"})
}
