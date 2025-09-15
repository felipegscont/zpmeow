package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// MESSAGE REQUEST DTOs
// ============================================================================

// SendTextRequest represents a request to send a text message
type SendTextRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	Body  string `json:"body" binding:"required" example:"Hello, World!"`
}

// SendMediaRequest represents a request to send media via URL
type SendMediaRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MediaType string `json:"media_type" binding:"required" example:"image"`
	MediaURL  string `json:"media_url" binding:"required" example:"https://example.com/image.jpg"`
	Caption   string `json:"caption,omitempty" example:"Check out this image!"`
}

// SendLocationRequest represents a request to send location
type SendLocationRequest struct {
	Phone     string  `json:"phone" binding:"required" example:"5511999999999"`
	Latitude  float64 `json:"latitude" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	Address   string  `json:"address,omitempty" example:"São Paulo, SP, Brazil"`
}

// SendContactRequest represents a request to send contact information
type SendContactRequest struct {
	Phone        string `json:"phone" binding:"required" example:"5511999999999"`
	ContactName  string `json:"contact_name" binding:"required" example:"John Doe"`
	ContactPhone string `json:"contact_phone" binding:"required" example:"5511888888888"`
}

// SendImageRequest represents an image message request
type SendImageRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Image   string `json:"image" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."`   // Base64 data URL or HTTP URL
	Caption string `json:"caption,omitempty" example:"Check out this image!"`
}

// SendAudioRequest represents an audio message request
type SendAudioRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	Audio string `json:"audio" binding:"required" example:"data:audio/mp3;base64,SUQzBAA..."`   // Base64 data URL or HTTP URL
	PTT   bool   `json:"ptt,omitempty" example:"true"`             // Push to talk
}

// SendVideoRequest represents a video message request
type SendVideoRequest struct {
	Phone       string `json:"phone" binding:"required" example:"5511999999999"`
	Video       string `json:"video" binding:"required" example:"data:video/mp4;base64,AAAAIGZ0eXA..."`       // Base64 data URL or HTTP URL
	Caption     string `json:"caption,omitempty" example:"Check out this video!"`
	GifPlayback bool   `json:"gif_playback,omitempty" example:"false"`
}

// SendDocumentRequest represents a document message request
type SendDocumentRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Document string `json:"document" binding:"required" example:"data:application/pdf;base64,JVBERi0x..."`   // Base64 data URL or HTTP URL
	FileName string `json:"filename,omitempty" example:"document.pdf"`
	MimeType string `json:"mimetype,omitempty" example:"application/pdf"`
}

// SendStickerRequest represents a sticker message request
type SendStickerRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Sticker string `json:"sticker" binding:"required" example:"data:image/webp;base64,UklGRnoGAABXRUJQ..."`   // Base64 data URL or HTTP URL
}

// MessageStatusRequest represents a request to get message status
type MessageStatusRequest struct {
	MessageID string `json:"message_id" binding:"required" example:"msg_123456789"`
}

// ChatHistoryRequest represents a request to get chat history
type ChatHistoryRequest struct {
	Phone  string `json:"phone" binding:"required" example:"5511999999999"`
	Limit  int    `json:"limit,omitempty" example:"50"`
	Offset int    `json:"offset,omitempty" example:"0"`
}

// BulkMessageRequest represents a request to send bulk messages
type BulkMessageRequest struct {
	Recipients []string `json:"recipients" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
	Message    string   `json:"message" binding:"required" example:"Hello, everyone!"`
}

// ============================================================================
// MESSAGE DATA STRUCTURES
// ============================================================================

// MessageStatusData represents message status information
type MessageStatusData struct {
	MessageID string    `json:"message_id" example:"msg_123456789"`
	Status    string    `json:"status" example:"delivered"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}



// BulkMessageResult represents bulk message operation result
type BulkMessageResult struct {
	Phone     string `json:"phone" example:"5511999999999"`
	MessageID string `json:"message_id" example:"msg_123456789"`
	Status    string `json:"status" example:"sent"`
	Error     string `json:"error,omitempty" example:""`
}

// MessageEventData represents message event data for webhooks
type MessageEventData struct {
	MessageID   string `json:"message_id"`
	From        string `json:"from"`
	To          string `json:"to"`
	IsFromMe    bool   `json:"is_from_me"`
	IsGroup     bool   `json:"is_group"`
	Timestamp   int64  `json:"timestamp"`
	MessageType string `json:"message_type"`
	Body        string `json:"body,omitempty"`
	Caption     string `json:"caption,omitempty"`
	MediaURL    string `json:"media_url,omitempty"`
}

// ============================================================================
// LEGACY RESPONSE DATA STRUCTURES (for backward compatibility)
// ============================================================================

// SendMessageRequest represents legacy send message request
type SendMessageRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Message string `json:"message" binding:"required" example:"Hello, World!"`
}

// SendMessageResponseData represents legacy message response data
type SendMessageResponseData struct {
	MessageID   string    `json:"message_id" example:"msg_123456789"`
	Phone       string    `json:"phone" example:"5511999999999"`
	MessageType string    `json:"message_type" example:"text"`
	Content     string    `json:"content" example:"Hello, World!"`
	Status      string    `json:"status" example:"sent"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	ServerID    string    `json:"server_id,omitempty" example:"server_msg_123"`
	Sender      string    `json:"sender,omitempty" example:"5511888888888@s.whatsapp.net"`
}

// SendMessageResponse represents legacy send message response
type SendMessageResponse struct {
	Status  int                     `json:"status" example:"200"`
	Message string                  `json:"message" example:"Message sent successfully"`
	Data    SendMessageResponseData `json:"data"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Validation error"`
	Error   string `json:"error" example:"Field validation failed"`
}



// MessageStatusResponseData represents legacy message status response data
type MessageStatusResponseData struct {
	MessageID string    `json:"message_id" example:"msg_123456789"`
	Status    string    `json:"status" example:"delivered"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

// BulkMessageResponseData represents legacy bulk message response data (single result)
type BulkMessageResponseData struct {
	Phone     string `json:"phone" example:"5511999999999"`
	MessageID string `json:"message_id" example:"msg_123456789"`
	Status    string `json:"status" example:"sent"`
	Error     string `json:"error,omitempty" example:""`
}

// ============================================================================
// STANDARD MESSAGE RESPONSE FORMAT
// ============================================================================

// MessageResponse represents the standardized message API response format
type MessageResponse struct {
	Success bool                    `json:"success"`
	Code    int                     `json:"code"`
	Data    MessageResponseData     `json:"data"`
	Error   *MessageErrorResponse   `json:"error,omitempty"`
}

// MessageResponseData contains the response data
type MessageResponseData struct {
	Key       MessageKey     `json:"key"`
	Message   MessagePayload `json:"message"`
	Timestamp int64          `json:"timestamp"`
}

// MessageKey identifies the message
type MessageKey struct {
	RemoteJid string `json:"remoteJid"`
	ID        string `json:"id"`
	FromMe    bool   `json:"fromMe"`
}

// MessagePayload contains the message content based on type
type MessagePayload struct {
	Text     *TextMessagePayload     `json:"text,omitempty"`
	Image    *ImageMessagePayload    `json:"image,omitempty"`
	Audio    *AudioMessagePayload    `json:"audio,omitempty"`
	Video    *VideoMessagePayload    `json:"video,omitempty"`
	Document *DocumentMessagePayload `json:"document,omitempty"`
	Sticker  *StickerMessagePayload  `json:"sticker,omitempty"`
	Contact  *ContactMessagePayload  `json:"contact,omitempty"`
	Location *LocationMessagePayload `json:"location,omitempty"`
}

// MessageErrorResponse represents error information
type MessageErrorResponse struct {
	Code    string `json:"code" example:"INVALID_PHONE"`
	Message string `json:"message" example:"Invalid phone number format"`
	Details string `json:"details,omitempty" example:"Phone number must include country code"`
}

// TextMessagePayload represents text message content
type TextMessagePayload struct {
	Text string `json:"text" example:"Hello, World!"`
}

// ImageMessagePayload represents image message content
type ImageMessagePayload struct {
	URL     string `json:"url" example:"https://example.com/image.jpg"`
	Caption string `json:"caption,omitempty" example:"Check out this image!"`
}

// AudioMessagePayload represents audio message content
type AudioMessagePayload struct {
	URL string `json:"url" example:"https://example.com/audio.mp3"`
	PTT bool   `json:"ptt" example:"false"`
}

// VideoMessagePayload represents video message content
type VideoMessagePayload struct {
	URL         string `json:"url" example:"https://example.com/video.mp4"`
	Caption     string `json:"caption,omitempty" example:"Check out this video!"`
	GifPlayback bool   `json:"gifPlayback" example:"false"`
}

// DocumentMessagePayload represents document message content
type DocumentMessagePayload struct {
	URL      string `json:"url" example:"https://example.com/document.pdf"`
	FileName string `json:"fileName" example:"document.pdf"`
	Mimetype string `json:"mimetype" example:"application/pdf"`
}

// StickerMessagePayload represents sticker message content
type StickerMessagePayload struct {
	URL string `json:"url" example:"https://example.com/sticker.webp"`
}

// ContactMessagePayload represents contact message content
type ContactMessagePayload struct {
	DisplayName string `json:"displayName" example:"John Doe"`
	Vcard       string `json:"vcard" example:"BEGIN:VCARD..."`
}

// LocationMessagePayload represents location message content
type LocationMessagePayload struct {
	Latitude  float64 `json:"latitude" example:"-23.5505"`
	Longitude float64 `json:"longitude" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	URL       string  `json:"url,omitempty" example:"https://maps.google.com/..."`
}

// ============================================================================
// INTERNAL HELPER FUNCTIONS
// ============================================================================

// phoneToRemoteJid converts a phone number to WhatsApp JID format
func phoneToRemoteJid(phone string) string {
	if phone == "" {
		return ""
	}
	// If already contains @s.whatsapp.net, return as is
	if strings.Contains(phone, "@s.whatsapp.net") {
		return phone
	}
	return phone + "@s.whatsapp.net"
}

// NewMessageResponse creates a new standardized message response
func NewMessageResponse(success bool, code int, remoteJid, messageID string, fromMe bool, messagePayload MessagePayload) *MessageResponse {
	return &MessageResponse{
		Success: success,
		Code:    code,
		Data: MessageResponseData{
			Key: MessageKey{
				RemoteJid: remoteJid,
				ID:        messageID,
				FromMe:    fromMe,
			},
			Message:   messagePayload,
			Timestamp: time.Now().Unix(),
		},
	}
}

// NewTextResponse creates a response for text message
func NewTextResponse(success bool, code int, phone, messageID, text string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Text: &TextMessagePayload{Text: text},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewLocationResponse creates a response for location message
func NewLocationResponse(success bool, code int, phone, messageID string, latitude, longitude float64, name, url string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Location: &LocationMessagePayload{
			Latitude:  latitude,
			Longitude: longitude,
			Name:      name,
			URL:       url,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewContactResponse creates a response for contact message
func NewContactResponse(success bool, code int, phone, messageID, displayName, vcard string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Contact: &ContactMessagePayload{
			DisplayName: displayName,
			Vcard:       vcard,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewImageResponse creates a response for image message
func NewImageResponse(success bool, code int, phone, messageID, url, caption string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Image: &ImageMessagePayload{
			URL:     url,
			Caption: caption,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewAudioResponse creates a response for audio message
func NewAudioResponse(success bool, code int, phone, messageID, url string, ptt, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Audio: &AudioMessagePayload{
			URL: url,
			PTT: ptt,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewVideoResponse creates a response for video message
func NewVideoResponse(success bool, code int, phone, messageID, url, caption string, gifPlayback, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Video: &VideoMessagePayload{
			URL:         url,
			Caption:     caption,
			GifPlayback: gifPlayback,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewDocumentResponse creates a response for document message
func NewDocumentResponse(success bool, code int, phone, messageID, url, fileName, mimeType string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Document: &DocumentMessagePayload{
			URL:      url,
			FileName: fileName,
			Mimetype: mimeType,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewStickerResponse creates a response for sticker message
func NewStickerResponse(success bool, code int, phone, messageID, url string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Sticker: &StickerMessagePayload{
			URL: url,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

// NewMessageErrorResponse creates a standardized error response
func NewMessageErrorResponse(code int, errorCode, message, details string) *MessageResponse {
	return &MessageResponse{
		Success: false,
		Code:    code,
		Data: MessageResponseData{
			Timestamp: time.Now().Unix(),
		},
		Error: &MessageErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
}

// ToJSON converts the response to pretty-formatted JSON
func (r *MessageResponse) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// ToJSONString converts the response to pretty-formatted JSON string
func (r *MessageResponse) ToJSONString() (string, error) {
	jsonBytes, err := r.ToJSON()
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// MessageValidationError represents a validation error for message requests
type MessageValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *MessageValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// validateMessagePhone checks if a phone number is valid for messages
func validateMessagePhone(phone string) bool {
	if phone == "" {
		return false
	}
	// Basic validation - should start with country code
	return len(phone) >= 10 && len(phone) <= 15
}

// Validate validates a SendTextRequest
func (r *SendTextRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Body == "" {
		return &MessageValidationError{Field: "body", Message: "Message body is required"}
	}
	return nil
}

// Validate validates a SendMediaRequest
func (r *SendMediaRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	// Use basic media type validation
	validTypes := []string{"image", "audio", "video", "document", "sticker"}
	isValid := false
	for _, validType := range validTypes {
		if r.MediaType == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return &MessageValidationError{Field: "media_type", Message: "Invalid media type"}
	}
	if r.MediaURL == "" {
		return &MessageValidationError{Field: "media_url", Message: "Media URL is required"}
	}
	return nil
}
