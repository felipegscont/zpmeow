package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// PingResponse represents the response for ping endpoint
type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

// Ping godoc
// @Summary Health check endpoint
// @Description Returns a simple pong message to verify the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, PingResponse{
		Message: "pong",
	})
}
