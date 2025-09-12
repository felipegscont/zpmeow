package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"zpmeow/internal/application/dto/response"
)

// ResponseService handles the creation of standardized responses
type ResponseService struct{}

// generateMessageID generates a unique message ID
func (r *ResponseService) generateMessageID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CreateSendResponse creates a standardized send response
func (r *ResponseService) CreateSendResponse() response.SendResponse {
	return response.SendResponse{
		Timestamp: time.Now().Unix(),
		ID:        r.generateMessageID(),
	}
}

// CreateSendResponseWithDetails creates a send response with additional details
func (r *ResponseService) CreateSendResponseWithDetails(serverID, sender string) response.SendResponse {
	return response.SendResponse{
		Timestamp: time.Now().Unix(),
		ID:        r.generateMessageID(),
		ServerID:  serverID,
		Sender:    sender,
	}
}

// CreateErrorResponse creates a standardized error response
func (r *ResponseService) CreateErrorResponse(message string, details ...string) map[string]interface{} {
	response := map[string]interface{}{
		"error": message,
	}
	
	if len(details) > 0 {
		response["details"] = details[0]
	}
	
	return response
}

// CreateSuccessResponse creates a standardized success response
func (r *ResponseService) CreateSuccessResponse(message string, data ...interface{}) map[string]interface{} {
	response := map[string]interface{}{
		"message": message,
	}
	
	if len(data) > 0 {
		response["data"] = data[0]
	}
	
	return response
}

// CreateWebhookResponse creates a webhook configuration response
func (r *ResponseService) CreateWebhookResponse(webhookURL string, events []string, active bool) map[string]interface{} {
	return map[string]interface{}{
		"webhook": webhookURL,
		"events":  events,
		"active":  active,
	}
}

// CreatePresenceResponse creates a presence update response
func (r *ResponseService) CreatePresenceResponse(phone, state, media string) map[string]interface{} {
	response := map[string]interface{}{
		"phone": phone,
		"state": state,
	}
	
	if media != "" {
		response["media"] = media
	}
	
	return response
}

// CreateUserInfoResponse creates a user info response
func (r *ResponseService) CreateUserInfoResponse(users []response.UserCheckResult) map[string]interface{} {
	return map[string]interface{}{
		"users": users,
		"count": len(users),
	}
}

// CreateAvatarResponse creates an avatar response
func (r *ResponseService) CreateAvatarResponse(url, id, avatarType, directURL string) response.AvatarResponse {
	return response.AvatarResponse{
		URL:       url,
		ID:        id,
		Type:      avatarType,
		DirectURL: directURL,
	}
}

// CreateGroupResponse creates a group operation response
func (r *ResponseService) CreateGroupResponse(groupJID, action string, participants []string) map[string]interface{} {
	response := map[string]interface{}{
		"groupJid": groupJID,
		"action":   action,
	}
	
	if len(participants) > 0 {
		response["participants"] = participants
	}
	
	return response
}

// CreateNewsletterResponse creates a newsletter response
func (r *ResponseService) CreateNewsletterResponse(newsletters []response.NewsletterMetadata) response.NewsletterResponse {
	return response.NewsletterResponse{
		Newsletter: newsletters,
	}
}

// CreateChatResponse creates a chat operation response
func (r *ResponseService) CreateChatResponse(action, chatJID, messageID string) map[string]interface{} {
	response := map[string]interface{}{
		"action": action,
	}
	
	if chatJID != "" {
		response["chatJid"] = chatJID
	}
	
	if messageID != "" {
		response["messageId"] = messageID
	}
	
	return response
}

// CreateDownloadResponse creates a media download response
func (r *ResponseService) CreateDownloadResponse(data []byte, mimeType, filename string) map[string]interface{} {
	response := map[string]interface{}{
		"data":     data,
		"mimeType": mimeType,
		"size":     len(data),
	}
	
	if filename != "" {
		response["filename"] = filename
	}
	
	return response
}

// CreateValidationErrorResponse creates a validation error response with field details
func (r *ResponseService) CreateValidationErrorResponse(field, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": "Validation failed",
		"field": field,
		"message": message,
	}
}

// CreatePaginatedResponse creates a paginated response
func (r *ResponseService) CreatePaginatedResponse(data interface{}, total, page, limit int) map[string]interface{} {
	return map[string]interface{}{
		"data":  data,
		"total": total,
		"page":  page,
		"limit": limit,
	}
}

// Global instance
var DefaultResponseService = &ResponseService{}
