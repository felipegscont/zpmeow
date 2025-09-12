package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}


type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}


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


func RespondWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}


func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}


func RespondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}


func RespondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}
