package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(router *gin.Engine)
}

type BaseHandler struct{}

func (h *BaseHandler) SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"message": message,
		"data":    data,
	})
}

func (h *BaseHandler) SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := gin.H{
		"status":  statusCode,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	c.JSON(statusCode, response)
}

func (h *BaseHandler) SendValidationErrorResponse(c *gin.Context, err error) {
	h.SendErrorResponse(c, http.StatusBadRequest, "Validation error", err)
}

func (h *BaseHandler) SendInternalErrorResponse(c *gin.Context, err error) {
	h.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", err)
}

func (h *BaseHandler) SendNotFoundResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, http.StatusNotFound, message, nil)
}

func (h *BaseHandler) SendUnauthorizedResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, http.StatusUnauthorized, message, nil)
}

func (h *BaseHandler) SendForbiddenResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, http.StatusForbidden, message, nil)
}

func (h *BaseHandler) SendConflictResponse(c *gin.Context, message string, err error) {
	h.SendErrorResponse(c, http.StatusConflict, message, err)
}

func (h *BaseHandler) GetSessionIDFromPath(c *gin.Context) string {
	return c.Param("sessionId")
}

func (h *BaseHandler) GetQueryParam(c *gin.Context, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (h *BaseHandler) GetQueryParamInt(c *gin.Context, key string, defaultValue int) int {
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

func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}

func (h *BaseHandler) BindQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}

func (h *BaseHandler) BindURI(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}
