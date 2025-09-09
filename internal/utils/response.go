package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, message string, details ...string) {
	response := ErrorResponse{
		Error: message,
		Code:  statusCode,
	}
	
	if len(details) > 0 {
		response.Details = details[0]
	}
	
	c.JSON(statusCode, response)
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, message string, data ...interface{}) {
	response := SuccessResponse{
		Success: true,
		Message: message,
	}

	if len(data) > 0 {
		response.Data = data[0]
	}

	c.JSON(http.StatusOK, response)
}

// RespondWithData sends data with a 200 status
func RespondWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RespondCreated sends data with a 201 status
func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// RespondNoContent sends a 204 status with no body
func RespondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RespondWithJSON sends a JSON response with the specified status code
func RespondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}
