package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"meow/internal/infra/logging"
	"meow/internal/interfaces/dto"
)

// Handler interface for route registration
type Handler interface {
	RegisterRoutes(router *gin.Engine)
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	logger    logging.Logger
	validator *dto.Validator
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(moduleName string) *BaseHandler {
	return &BaseHandler{
		logger:    logging.GetLogger().Sub(moduleName),
		validator: dto.NewValidator(),
	}
}

// ValidateRequest validates a request DTO
func (h *BaseHandler) ValidateRequest(req interface{}) error {
	// Check if the request has a Validate method
	if validator, ok := req.(interface{ Validate() error }); ok {
		return validator.Validate()
	}

	// Fallback to struct validation
	return h.validator.Validate(req)
}

// BindAndValidate binds JSON request and validates it
func (h *BaseHandler) BindAndValidate(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return h.ValidateRequest(req)
}

// HTTPHandler is deprecated - use BaseHandler instead
type HTTPHandler struct {
	*BaseHandler
}

// SendSuccessResponse sends a standardized success response
func (h *BaseHandler) SendSuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	response := dto.NewSuccessResponse(statusCode, data)
	c.JSON(statusCode, response)
}

// SendActionResponse sends a response for action-based operations
func (h *BaseHandler) SendActionResponse(c *gin.Context, statusCode int, action string, data interface{}) {
	response := dto.NewActionResponse(statusCode, action, data)
	c.JSON(statusCode, response)
}

// SendErrorResponse sends a standardized error response
func (h *BaseHandler) SendErrorResponse(c *gin.Context, statusCode int, errorCode, message string, err error) {
	details := ""
	if err != nil {
		details = err.Error()
	}
	response := dto.NewErrorResponse(statusCode, errorCode, message, details)
	c.JSON(statusCode, response)
}

// Legacy method for backward compatibility
func (h *HTTPHandler) SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	h.BaseHandler.SendSuccessResponse(c, statusCode, data)
}

// Legacy method for backward compatibility
func (h *HTTPHandler) SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	h.BaseHandler.SendErrorResponse(c, statusCode, dto.ErrorCodeInternalError, message, err)
}

// Standardized error response methods using BaseHandler
func (h *BaseHandler) SendValidationErrorResponse(c *gin.Context, err error) {
	h.SendErrorResponse(c, dto.StatusBadRequest, dto.ErrorCodeValidationFailed, "Validation error", err)
}

func (h *BaseHandler) SendInternalErrorResponse(c *gin.Context, err error) {
	h.SendErrorResponse(c, dto.StatusInternalServerError, dto.ErrorCodeInternalError, "Internal server error", err)
}

func (h *BaseHandler) SendNotFoundResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, dto.StatusNotFound, dto.ErrorCodeNotFound, message, nil)
}

func (h *BaseHandler) SendUnauthorizedResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, dto.StatusUnauthorized, dto.ErrorCodeUnauthorized, message, nil)
}

func (h *BaseHandler) SendForbiddenResponse(c *gin.Context, message string) {
	h.SendErrorResponse(c, dto.StatusForbidden, dto.ErrorCodeForbidden, message, nil)
}

func (h *BaseHandler) SendConflictResponse(c *gin.Context, message string, err error) {
	h.SendErrorResponse(c, dto.StatusConflict, dto.ErrorCodeConflict, message, err)
}

// Legacy methods for backward compatibility
func (h *HTTPHandler) SendValidationErrorResponse(c *gin.Context, err error) {
	h.BaseHandler.SendValidationErrorResponse(c, err)
}

func (h *HTTPHandler) SendInternalErrorResponse(c *gin.Context, err error) {
	h.BaseHandler.SendInternalErrorResponse(c, err)
}

func (h *HTTPHandler) SendNotFoundResponse(c *gin.Context, message string) {
	h.BaseHandler.SendNotFoundResponse(c, message)
}

func (h *HTTPHandler) SendUnauthorizedResponse(c *gin.Context, message string) {
	h.BaseHandler.SendUnauthorizedResponse(c, message)
}

func (h *HTTPHandler) SendForbiddenResponse(c *gin.Context, message string) {
	h.BaseHandler.SendForbiddenResponse(c, message)
}

func (h *HTTPHandler) SendConflictResponse(c *gin.Context, message string, err error) {
	h.BaseHandler.SendConflictResponse(c, message, err)
}

// Utility methods for BaseHandler
func (h *BaseHandler) GetSessionIDFromPath(c *gin.Context) string {
	return c.Param("sessionId")
}

func (h *BaseHandler) GetQueryParam(c *gin.Context, key, defaultValue string) string {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return defaultValue
	}
	return value
}

func (h *BaseHandler) GetQueryParamInt(c *gin.Context, key string, defaultValue int) int {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return defaultValue
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

// Legacy methods for backward compatibility
func (h *HTTPHandler) GetSessionIDFromPath(c *gin.Context) string {
	return h.BaseHandler.GetSessionIDFromPath(c)
}

func (h *HTTPHandler) GetQueryParam(c *gin.Context, key, defaultValue string) string {
	return h.BaseHandler.GetQueryParam(c, key, defaultValue)
}

func (h *HTTPHandler) GetQueryParamInt(c *gin.Context, key string, defaultValue int) int {
	return h.BaseHandler.GetQueryParamInt(c, key, defaultValue)
}

// Binding methods for BaseHandler
func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.logger.Errorf("JSON binding failed: %v", err)
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}

func (h *BaseHandler) BindQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		h.logger.Errorf("Query binding failed: %v", err)
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}

func (h *BaseHandler) BindURI(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		h.logger.Errorf("URI binding failed: %v", err)
		h.SendValidationErrorResponse(c, err)
		return err
	}
	return nil
}

// Legacy methods for backward compatibility
func (h *HTTPHandler) BindJSON(c *gin.Context, obj interface{}) error {
	return h.BaseHandler.BindJSON(c, obj)
}

func (h *HTTPHandler) BindQuery(c *gin.Context, obj interface{}) error {
	return h.BaseHandler.BindQuery(c, obj)
}

func (h *HTTPHandler) BindURI(c *gin.Context, obj interface{}) error {
	return h.BaseHandler.BindURI(c, obj)
}
