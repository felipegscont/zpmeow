package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra"
	"zpmeow/internal/infra/logger"
)

type GroupHandler struct {
	sessionService domain.SessionService
	meowService    *infra.MeowServiceImpl
	logger         logger.Logger
}

func NewGroupHandler(sessionService domain.SessionService, meowService *infra.MeowServiceImpl) *GroupHandler {
	return &GroupHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("group-handler"),
	}
}

func (h *GroupHandler) GetGroups(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetGroups - stub implementation"})
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "CreateGroup - stub implementation"})
}

func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "LeaveGroup - stub implementation"})
}

func (h *GroupHandler) GetInviteLink(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetInviteLink - stub implementation"})
}
