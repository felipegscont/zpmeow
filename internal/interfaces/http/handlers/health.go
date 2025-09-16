package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/infrastructure/logging"
)

// ============================================================================
// HEALTH HANDLER - SELF-CONTAINED WITH INTERNAL HELPERS
// ============================================================================

type HealthHandler struct {
	logger logging.Logger
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		logger: logging.GetLogger().Sub("health-handler"),
	}
}

// ============================================================================
// HEALTH RESPONSE DTOs
// ============================================================================

// HealthStandardResponse represents the standardized response format for health operations
type HealthStandardResponse struct {
	Success bool                 `json:"success"`
	Code    int                  `json:"code"`
	Data    HealthData           `json:"data"`
	Error   *HealthErrorResponse `json:"error,omitempty"`
}

// HealthData contains the response data for health operations
type HealthData struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Service is healthy"`
	Version string `json:"version,omitempty" example:"1.0.0"`
	Service string `json:"service" example:"zpmeow"`
}

// HealthErrorResponse represents error information for health operations
type HealthErrorResponse struct {
	Code    string `json:"code" example:"HEALTH_CHECK_FAILED"`
	Message string `json:"message" example:"Health check failed"`
	Details string `json:"details,omitempty" example:"Service is not responding"`
}

// ============================================================================
// INTERNAL HELPER FUNCTIONS
// ============================================================================

// sendSuccessResponse sends a standardized success response
func (h *HealthHandler) sendSuccessResponse(c *gin.Context, status, message, version string) {
	response := &HealthStandardResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: HealthData{
			Status:  status,
			Message: message,
			Version: version,
			Service: "zpmeow",
		},
	}

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		h.logger.Errorf("Failed to marshal health response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format response"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
}

// Health godoc
//
//	@Summary		Health check endpoint
//	@Description	Returns the health status of the service using standardized response format
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthStandardResponse	"Service is healthy"
//	@Router			/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	h.logger.Infof("Health check requested")

	// Perform basic health checks here if needed
	// For now, we'll assume the service is healthy if we can respond

	h.sendSuccessResponse(c, "ok", "Service is healthy", "1.0.0")
	h.logger.Infof("Health check completed successfully")
}

// Ping godoc
//
//	@Summary		Ping endpoint
//	@Description	Simple ping endpoint that returns pong using standardized response format
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthStandardResponse	"Pong response"
//	@Router			/ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	h.logger.Infof("Ping requested")

	h.sendSuccessResponse(c, "ok", "pong", "")
	h.logger.Infof("Ping completed successfully")
}
