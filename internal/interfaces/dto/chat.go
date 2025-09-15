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

// ReactToMessageRequest represents a request to react to a message
type ReactToMessageRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"message_id" binding:"required" example:"msg_123"`
	Emoji     string `json:"emoji" binding:"required" example:"👍"` // Use "remove" to remove reaction
}

// DeleteMessageRequest represents a request to delete a message
type DeleteMessageRequest struct {
	Phone       string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID   string `json:"message_id" binding:"required" example:"msg_123"`
	ForEveryone bool   `json:"for_everyone" example:"true"` // true = delete for everyone, false = delete for me
}

// EditMessageRequest represents a request to edit a message
type EditMessageRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"message_id" binding:"required" example:"msg_123"`
	NewText   string `json:"new_text" binding:"required" example:"Edited message text"`
}

// DownloadMediaRequest represents a request to download media
type DownloadMediaRequest struct {
	MessageID string `json:"message_id" binding:"required" example:"msg_123"`
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

// MediaDownloadResponse represents the response for media download
type MediaDownloadResponse struct {
	Success   bool   `json:"success" example:"true"`
	Code      int    `json:"code" example:"200"`
	MessageID string `json:"message_id" example:"msg_123"`
	MediaType string `json:"media_type" example:"image"`
	MimeType  string `json:"mime_type" example:"image/jpeg"`
	Data      []byte `json:"data"` // Base64 encoded media data
	Size      int    `json:"size" example:"1024"`
}

// ChatHistoryResponse represents the response for chat history
type ChatHistoryResponse struct {
	Success bool                     `json:"success" example:"true"`
	Code    int                      `json:"code" example:"200"`
	Data    ChatHistoryResponseData  `json:"data"`
	Error   *ChatErrorResponse       `json:"error,omitempty"`
}

// ChatHistoryResponseData contains chat history data
type ChatHistoryResponseData struct {
	Phone    string            `json:"phone" example:"5511999999999"`
	Messages []ChatHistoryData `json:"messages"`
	Count    int               `json:"count" example:"10"`
	Limit    int               `json:"limit" example:"50"`
}

// ChatHistoryData represents a single message in chat history
type ChatHistoryData struct {
	MessageID   string    `json:"message_id" example:"msg_123456789"`
	Phone       string    `json:"phone" example:"5511999999999"`
	FromPhone   string    `json:"from_phone" example:"5511888888888"`
	MessageType string    `json:"message_type" example:"text"`
	Content     string    `json:"content" example:"Hello, World!"`
	MediaURL    string    `json:"media_url,omitempty" example:"https://example.com/image.jpg"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	IsFromMe    bool      `json:"is_from_me" example:"false"`
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
