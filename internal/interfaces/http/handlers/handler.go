package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(router *gin.Engine)
}

type HTTPHandler struct{}

func (h *HTTPHandler) SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"message": message,
		"data":    data,
	})
}

func (h *HTTPHandler) SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := gin.H{
		"status":  statusCode,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	c.JSON(statusCode, response)
}

func (h *HTTPHandler) SendValidationErrorResponse(c *gin.Context, err error) {
	h.SendErrorResponse(c, http.StatusBadRequest, "Validation error", err)
}

func (h *HTTPHandler) SendInternalErrorResponse(c *gin.Context, err error) {
	h.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", err)
}

func (h *HTTPHandler) SendNotFoundResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, http.StatusNotFound, message, nil)
}

func (h *HTTPHandler) SendUnauthorizedResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, http.StatusUnauthorized, message, nil)
}

func (h *HTTPHandler) SendForbiddenResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, http.StatusForbidden, message, nil)
}

func (h *HTTPHandler) SendConflictResponse(c *gin.Context, message string, err error) {
	h.SendErrorResponse(c, http.StatusConflict, message, err)
}

func (h *HTTPHandler) GetSessionIDFromPath(c *gin.Context) string {
	return c.Param("sessionId")
}

func (h *HTTPHandler) GetQueryParam(c *gin.Context, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (h *HTTPHandler) GetQueryParamInt(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	if intValue := parseIntOrDefault(value, defaultValue); intValue != defaultValue {
		return intValue
	}
	return defaultValue
}

func parseIntOrDefault(value string, defaultValue int) int {
	switch value {
	case "0":
		return 0
	case "1":
		return 1
	case "10":
		return 10
	case "25":
		return 25
	case "50":
		return 50
	case "100":
		return 100
	default:
		return defaultValue
	}
}

func (h *HTTPHandler) BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}

func (h *HTTPHandler) BindQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}

func (h *HTTPHandler) BindURI(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}
