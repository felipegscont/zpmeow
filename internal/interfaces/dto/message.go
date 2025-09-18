package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SendTextRequest struct {
	Phone string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	Body  string `json:"body" validate:"required,min=1,max=4096" binding:"required" example:"Hello, World!"`
}

// Validate validates the SendTextRequest
func (r *SendTextRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendMediaRequest struct {
	Phone     string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	MediaType string `json:"media_type" validate:"required,oneof=image video audio document" binding:"required" example:"image"`
	MediaURL  string `json:"media_url" validate:"required,url" binding:"required" example:"https://example.com/image.jpg"`
	Caption   string `json:"caption,omitempty" validate:"omitempty,max=1024" example:"Check out this image!"`
}

// Validate validates the SendMediaRequest
func (r *SendMediaRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendLocationRequest struct {
	Phone     string  `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" validate:"omitempty,max=100" example:"São Paulo"`
	Address   string  `json:"address,omitempty" validate:"omitempty,max=500" example:"São Paulo, SP, Brazil"`
}

// Validate validates the SendLocationRequest
func (r *SendLocationRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendContactRequest struct {
	Phone        string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	ContactName  string `json:"contact_name" validate:"required,min=1,max=100" binding:"required" example:"John Doe"`
	ContactPhone string `json:"contact_phone" validate:"required,phone_number" binding:"required" example:"5511888888888"`
}

// Validate validates the SendContactRequest
func (r *SendContactRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendImageRequest struct {
	Phone   string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	Image   string `json:"image" validate:"required,min=1" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."` // Base64 data URL or HTTP URL
	Caption string `json:"caption,omitempty" validate:"omitempty,max=1024" example:"Check out this image!"`
}

// Validate validates the SendImageRequest
func (r *SendImageRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendAudioRequest struct {
	Phone string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	Audio string `json:"audio" validate:"required,min=1" binding:"required" example:"data:audio/mp3;base64,SUQzBAA..."` // Base64 data URL or HTTP URL
	PTT   bool   `json:"ptt,omitempty" example:"true"`                                        // Push to talk
}

// Validate validates the SendAudioRequest
func (r *SendAudioRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendVideoRequest struct {
	Phone       string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	Video       string `json:"video" validate:"required,min=1" binding:"required" example:"data:video/mp4;base64,AAAAIGZ0eXA..."` // Base64 data URL or HTTP URL
	Caption     string `json:"caption,omitempty" validate:"omitempty,max=1024" example:"Check out this video!"`
	GifPlayback bool   `json:"gif_playback,omitempty" example:"false"`
}

// Validate validates the SendVideoRequest
func (r *SendVideoRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendDocumentRequest struct {
	Phone    string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	Document string `json:"document" validate:"required,min=1" binding:"required" example:"data:application/pdf;base64,JVBERi0x..."` // Base64 data URL or HTTP URL
	FileName string `json:"filename,omitempty" validate:"omitempty,max=255" example:"document.pdf"`
	MimeType string `json:"mimetype,omitempty" validate:"omitempty,max=100" example:"application/pdf"`
}

// Validate validates the SendDocumentRequest
func (r *SendDocumentRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type SendStickerRequest struct {
	Phone   string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
	Sticker string `json:"sticker" validate:"required,min=1" binding:"required" example:"data:image/webp;base64,UklGRnoGAABXRUJQ..."` // Base64 data URL or HTTP URL
}

// Validate validates the SendStickerRequest
func (r *SendStickerRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

type MessageStatusRequest struct {
	MessageID string `json:"message_id" binding:"required" example:"msg_123456789"`
}

type ChatHistoryRequest struct {
	Phone  string `json:"phone" binding:"required" example:"5511999999999"`
	Limit  int    `json:"limit,omitempty" example:"50"`
	Offset int    `json:"offset,omitempty" example:"0"`
}

type BulkMessageRequest struct {
	Recipients []string `json:"recipients" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
	Message    string   `json:"message" binding:"required" example:"Hello, everyone!"`
}

type MarkAsReadRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"5511999999999"`
	MessageIDs []string `json:"message_ids" binding:"required" example:"[\"msg_1\", \"msg_2\"]"`
}

type ReactToMessageRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"message_id" binding:"required" example:"3EB0D098B5FD4BF3BC4327"`
	Emoji     string `json:"emoji" binding:"required" example:"👍"` // Use "remove" to remove reaction
}

type DeleteMessageRequest struct {
	Phone       string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID   string `json:"message_id" binding:"required" example:"3EB0D098B5FD4BF3BC4327"`
	ForEveryone bool   `json:"for_everyone" example:"true"` // true = delete for everyone, false = delete for me
}

type EditMessageRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"message_id" binding:"required" example:"3EB0D098B5FD4BF3BC4327"`
	NewText   string `json:"new_text" binding:"required" example:"Edited message text"`
}

type MessageActionResponse struct {
	Success bool                        `json:"success"`
	Code    int                         `json:"code"`
	Data    MessageActionData           `json:"data"`
	Error   *MessageActionErrorResponse `json:"error,omitempty"`
}

type MessageActionData struct {
	Phone     string    `json:"phone" example:"5511999999999"`
	MessageID string    `json:"message_id,omitempty" example:"msg_123"`
	Action    string    `json:"action" example:"mark_read"`
	Status    string    `json:"status" example:"success"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

type MessageActionErrorResponse struct {
	Code    string `json:"code" example:"INVALID_PHONE"`
	Message string `json:"message" example:"Invalid phone number format"`
	Details string `json:"details,omitempty" example:"Phone number must include country code"`
}

type MessageStatusData struct {
	MessageID string    `json:"message_id" example:"msg_123456789"`
	Status    string    `json:"status" example:"delivered"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

type BulkMessageResult struct {
	Phone     string `json:"phone" example:"5511999999999"`
	MessageID string `json:"message_id" example:"msg_123456789"`
	Status    string `json:"status" example:"sent"`
	Error     string `json:"error,omitempty" example:""`
}

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

type SendMessageRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Message string `json:"message" binding:"required" example:"Hello, World!"`
}

type SendMessageResponseData struct {
	MessageID   string    `json:"message_id" example:"msg_123456789"`
	Phone       string    `json:"phone" example:"5511999999999"`
	MessageType string    `json:"message_type" example:"text"`
	Content     string    `json:"content" example:"Hello, World!"`
	Status      string    `json:"status" example:"sent"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	ServerID    string    `json:"server_id,omitempty" example:"server_msg_123"`
	Sender      string    `json:"sender,omitempty" example:"5511888888888@s.meow.net"`
}

type SendMessageResponse struct {
	Status  int                     `json:"status" example:"200"`
	Message string                  `json:"message" example:"Message sent successfully"`
	Data    SendMessageResponseData `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Validation error"`
	Error   string `json:"error" example:"Field validation failed"`
}

type MessageStatusResponseData struct {
	MessageID string    `json:"message_id" example:"msg_123456789"`
	Status    string    `json:"status" example:"delivered"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

type BulkMessageResponseData struct {
	Phone     string `json:"phone" example:"5511999999999"`
	MessageID string `json:"message_id" example:"msg_123456789"`
	Status    string `json:"status" example:"sent"`
	Error     string `json:"error,omitempty" example:""`
}

type MessageResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Data    MessageResponseData   `json:"data"`
	Error   *MessageErrorResponse `json:"error,omitempty"`
}

type MessageResponseData struct {
	Key       MessageKey     `json:"key"`
	Message   MessagePayload `json:"message"`
	Timestamp int64          `json:"timestamp"`
}

type MessageKey struct {
	RemoteJid string `json:"remoteJid"`
	ID        string `json:"id"`
	FromMe    bool   `json:"fromMe"`
}

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

type MessageErrorResponse struct {
	Code    string `json:"code" example:"INVALID_PHONE"`
	Message string `json:"message" example:"Invalid phone number format"`
	Details string `json:"details,omitempty" example:"Phone number must include country code"`
}

type TextMessagePayload struct {
	Text string `json:"text" example:"Hello, World!"`
}

type ImageMessagePayload struct {
	URL     string `json:"url" example:"https://example.com/image.jpg"`
	Caption string `json:"caption,omitempty" example:"Check out this image!"`
}

type AudioMessagePayload struct {
	URL string `json:"url" example:"https://example.com/audio.mp3"`
	PTT bool   `json:"ptt" example:"false"`
}

type VideoMessagePayload struct {
	URL         string `json:"url" example:"https://example.com/video.mp4"`
	Caption     string `json:"caption,omitempty" example:"Check out this video!"`
	GifPlayback bool   `json:"gifPlayback" example:"false"`
}

type DocumentMessagePayload struct {
	URL      string `json:"url" example:"https://example.com/document.pdf"`
	FileName string `json:"fileName" example:"document.pdf"`
	Mimetype string `json:"mimetype" example:"application/pdf"`
}

type StickerMessagePayload struct {
	URL string `json:"url" example:"https://example.com/sticker.webp"`
}

type ContactMessagePayload struct {
	DisplayName string `json:"displayName" example:"John Doe"`
	Vcard       string `json:"vcard" example:"BEGIN:VCARD..."`
}

type LocationMessagePayload struct {
	Latitude  float64 `json:"latitude" example:"-23.5505"`
	Longitude float64 `json:"longitude" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	URL       string  `json:"url,omitempty" example:"https://maps.google.com/..."`
}

func phoneToRemoteJid(phone string) string {
	if phone == "" {
		return ""
	}
	if strings.Contains(phone, "@s.meow.net") {
		return phone
	}
	return phone + "@s.meow.net"
}

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

func NewTextResponse(success bool, code int, phone, messageID, text string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Text: &TextMessagePayload{Text: text},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

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

func NewContactResponse(success bool, code int, phone, messageID, displayName, vcard string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Contact: &ContactMessagePayload{
			DisplayName: displayName,
			Vcard:       vcard,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

func NewImageResponse(success bool, code int, phone, messageID, url, caption string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Image: &ImageMessagePayload{
			URL:     url,
			Caption: caption,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

func NewAudioResponse(success bool, code int, phone, messageID, url string, ptt, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Audio: &AudioMessagePayload{
			URL: url,
			PTT: ptt,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

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

func NewStickerResponse(success bool, code int, phone, messageID, url string, fromMe bool) *MessageResponse {
	payload := MessagePayload{
		Sticker: &StickerMessagePayload{
			URL: url,
		},
	}
	return NewMessageResponse(success, code, phoneToRemoteJid(phone), messageID, fromMe, payload)
}

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

func (r *MessageResponse) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r *MessageResponse) ToJSONString() (string, error) {
	jsonBytes, err := r.ToJSON()
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

type MessageValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *MessageValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validateMessagePhone(phone string) bool {
	if phone == "" {
		return false
	}
	return len(phone) >= 10 && len(phone) <= 15
}

func NewMessageActionSuccessResponse(phone, messageID, action string) *MessageActionResponse {
	return &MessageActionResponse{
		Success: true,
		Code:    200,
		Data: MessageActionData{
			Phone:     phone,
			MessageID: messageID,
			Action:    action,
			Status:    "success",
			Timestamp: time.Now(),
		},
	}
}

func NewMessageActionErrorResponse(code int, errorCode, message, details string) *MessageActionResponse {
	return &MessageActionResponse{
		Success: false,
		Code:    code,
		Error: &MessageActionErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
}
