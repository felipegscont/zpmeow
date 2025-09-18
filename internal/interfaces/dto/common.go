package dto

import (
	"time"
)

// StandardResponse represents the common response structure used across all endpoints
type StandardResponse struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorInfo represents error information in responses
type ErrorInfo struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request parameters"`
	Details string `json:"details,omitempty" example:"Additional error details"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int `json:"page" example:"1"`
	PageSize   int `json:"page_size" example:"10"`
	Total      int `json:"total" example:"100"`
	TotalPages int `json:"total_pages" example:"10"`
}

// ActionData represents common action response data
type ActionData struct {
	Action    string    `json:"action" example:"create_session"`
	Status    string    `json:"status" example:"success"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

// NewSuccessResponse creates a standardized success response
func NewSuccessResponse(code int, data interface{}) *StandardResponse {
	return &StandardResponse{
		Success:   true,
		Code:      code,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse creates a standardized error response
func NewErrorResponse(code int, errorCode, message, details string) *StandardResponse {
	return &StandardResponse{
		Success: false,
		Code:    code,
		Data:    nil,
		Error: &ErrorInfo{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

// NewActionResponse creates a response for action-based operations
func NewActionResponse(code int, action string, data interface{}) *StandardResponse {
	actionData := ActionData{
		Action:    action,
		Status:    "success",
		Timestamp: time.Now(),
	}

	// If data is provided, merge it with action data
	responseData := map[string]interface{}{
		"action":    actionData.Action,
		"status":    actionData.Status,
		"timestamp": actionData.Timestamp,
	}

	if data != nil {
		responseData["result"] = data
	}

	return &StandardResponse{
		Success:   true,
		Code:      code,
		Data:      responseData,
		Timestamp: time.Now(),
	}
}

// Common HTTP status codes as constants
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusAccepted            = 202
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusInternalServerError = 500
)

// Common error codes
const (
	ErrorCodeInvalidRequest     = "INVALID_REQUEST"
	ErrorCodeValidationFailed   = "VALIDATION_FAILED"
	ErrorCodeUnauthorized       = "UNAUTHORIZED"
	ErrorCodeForbidden          = "FORBIDDEN"
	ErrorCodeNotFound           = "NOT_FOUND"
	ErrorCodeConflict           = "CONFLICT"
	ErrorCodeInternalError      = "INTERNAL_ERROR"
	ErrorCodeSessionNotFound    = "SESSION_NOT_FOUND"
	ErrorCodeSessionInactive    = "SESSION_INACTIVE"
	ErrorCodeInvalidPhoneNumber = "INVALID_PHONE_NUMBER"
	ErrorCodeInvalidMessageID   = "INVALID_MESSAGE_ID"
	ErrorCodeInvalidURL         = "INVALID_URL"
)
