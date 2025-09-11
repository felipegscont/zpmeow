package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type HealthHandler struct{}


func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}


type PingResponse struct {
	Message string `json:"message" example:"pong"`
}


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
