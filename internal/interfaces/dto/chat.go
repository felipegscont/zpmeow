package dto

import (
	"time"
)

// ============================================================================
// CHAT REQUEST DTOs
// ============================================================================

// SetPresenceRequest represents a request to set user presence in a chat
type SetPresenceRequest struct {
	Phone string `json:"phone,omitempty" example:"5511999999999"`
	State string `json:"state" binding:"required" example:"available"` // available, unavailable, composing, recording, paused
	Media string `json:"media,omitempty" example:""`                   // Optional media type for composing state
}

// MarkAsReadRequest represents a request to mark messages as read
type MarkAsReadRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"5511999999999"`
	MessageIDs []string `json:"message_ids" binding:"required" example:"[\"msg_1\", \"msg_2\"]"`
}

// ============================================================================
// CHAT RESPONSE DTOs
// ============================================================================

// ChatResponse represents the standardized response format for chat operations
type ChatResponse struct {
	Success bool                `json:"success"`
	Code    int                 `json:"code"`
	Data    ChatData            `json:"data"`
	Error   *ChatErrorResponse  `json:"error,omitempty"`
}

// ChatData contains the response data for chat operations
type ChatData struct {
	Phone     string    `json:"phone" example:"5511999999999"`
	MessageID string    `json:"message_id,omitempty" example:"msg_123"`
	Action    string    `json:"action" example:"mark_read"`
	Status    string    `json:"status" example:"success"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

// ChatErrorResponse represents error information for chat operations
type ChatErrorResponse struct {
	Code    string `json:"code" example:"INVALID_PHONE"`
	Message string `json:"message" example:"Invalid phone number format"`
	Details string `json:"details,omitempty" example:"Phone number must include country code"`
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// NewChatSuccessResponse creates a successful chat operation response
func NewChatSuccessResponse(phone, messageID, action string) *ChatResponse {
	return &ChatResponse{
		Success: true,
		Code:    200,
		Data: ChatData{
			Phone:     phone,
			MessageID: messageID,
			Action:    action,
			Status:    "success",
			Timestamp: time.Now(),
		},
	}
}

// NewChatErrorResponse creates an error response for chat operations
func NewChatErrorResponse(code int, errorCode, message, details string) *ChatResponse {
	return &ChatResponse{
		Success: false,
		Code:    code,
		Data: ChatData{
			Status:    "error",
			Timestamp: time.Now(),
		},
		Error: &ChatErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
}
