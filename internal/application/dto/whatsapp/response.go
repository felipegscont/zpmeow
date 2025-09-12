package whatsapp

import (
	"fmt"

	"go.mau.fi/whatsmeow"
)

// SendResponse represents the response from sending a WhatsApp message
type SendResponse struct {
	// Unix timestamp when the message was sent
	Timestamp int64 `json:"timestamp" example:"1640995200"`

	// Message ID assigned by WhatsApp
	ID string `json:"id" example:"3EB0C431C26A1916E07A"`

	// Server-assigned message ID (optional)
	ServerID string `json:"serverId,omitempty" example:"wamid.HBgNNTU5OTgxNzY5NTM2FQIAERgSMzNFNzE4QzY5QzE5MjE2RTdB"`

	// Sender JID (optional)
	Sender string `json:"sender,omitempty" example:"5511999999999@s.whatsapp.net"`

	// Success status
	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"messageId" example:"3EB0C431C26A1916E07A"`
}

// NewSendResponseFromWhatsmeow creates a SendResponse from whatsmeow.SendResponse
func NewSendResponseFromWhatsmeow(resp *whatsmeow.SendResponse, requestID string) SendResponse {
	return SendResponse{
		Timestamp: resp.Timestamp.Unix(),
		ID:        string(resp.ID),
		ServerID: func() string {
			if resp.ServerID == 0 {
				return ""
			}
			return fmt.Sprintf("%d", resp.ServerID)
		}(),
		Sender:  resp.Sender.String(),
		Success: true,
		MessageID: func() string {
			if requestID != "" {
				return requestID
			}
			return string(resp.ID)
		}(),
	}
}

// UserPresenceRequest represents a user presence update request
type UserPresenceRequest struct {
	Type string `json:"type" binding:"required" example:"available"` // available, unavailable
}

// CheckUserRequest represents a request to check if users are on WhatsApp
type CheckUserRequest struct {
	Phone []string `json:"phone" binding:"required" example:"+5511999999999,+5511888888888"`
}

// GetUserInfoRequest represents a request to get user information
type GetUserInfoRequest struct {
	Phone []string `json:"phone" binding:"required" example:"+5511999999999,+5511888888888"`
}

// ChatPresenceRequest represents a chat presence update request
type ChatPresenceRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	State string `json:"state" binding:"required" example:"composing"` // composing, paused
	Media string `json:"media" example:"audio"`                        // text, audio (optional, only for composing)
}

// ChatMarkReadRequest represents a request to mark messages as read
type ChatMarkReadRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"5511999999999"`
	MessageIDs []string `json:"messageIds" binding:"required" example:"3EB0C431C26A1916E07A,3EB0C431C26A1916E07B"`
}

// ChatReactRequest represents a request to react to a message
type ChatReactRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	Emoji     string `json:"emoji" binding:"required" example:"👍"`
}

// ChatDeleteRequest represents a request to delete a message
type ChatDeleteRequest struct {
	Phone       string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID   string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	ForEveryone bool   `json:"forEveryone,omitempty" example:"false"`
}

// ChatEditRequest represents a request to edit a message
type ChatEditRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	NewText   string `json:"newText" binding:"required" example:"Updated message text"`
}

// ChatDownloadRequest represents a request to download media
type ChatDownloadRequest struct {
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
}
