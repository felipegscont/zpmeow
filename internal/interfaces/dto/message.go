package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SendTextRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	Body  string `json:"body" binding:"required" example:"Hello, World!"`
}

type SendMediaRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MediaType string `json:"media_type" binding:"required" example:"image"`
	MediaURL  string `json:"media_url" binding:"required" example:"https://example.com/image.jpg"`
	Caption   string `json:"caption,omitempty" example:"Check out this image!"`
}

type SendLocationRequest struct {
	Phone     string  `json:"phone" binding:"required" example:"5511999999999"`
	Latitude  float64 `json:"latitude" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	Address   string  `json:"address,omitempty" example:"São Paulo, SP, Brazil"`
}

type SendContactRequest struct {
	Phone        string `json:"phone" binding:"required" example:"5511999999999"`
	ContactName  string `json:"contact_name" binding:"required" example:"John Doe"`
	ContactPhone string `json:"contact_phone" binding:"required" example:"5511888888888"`
}

type SendImageRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Image   string `json:"image" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."` // Base64 data URL or HTTP URL
	Caption string `json:"caption,omitempty" example:"Check out this image!"`
}

type SendAudioRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	Audio string `json:"audio" binding:"required" example:"data:audio/mp3;base64,SUQzBAA..."` // Base64 data URL or HTTP URL
	PTT   bool   `json:"ptt,omitempty" example:"true"`                                        // Push to talk
}

type SendVideoRequest struct {
	Phone       string `json:"phone" binding:"required" example:"5511999999999"`
	Video       string `json:"video" binding:"required" example:"data:video/mp4;base64,AAAAIGZ0eXA..."` // Base64 data URL or HTTP URL
	Caption     string `json:"caption,omitempty" example:"Check out this video!"`
	GifPlayback bool   `json:"gif_playback,omitempty" example:"false"`
}

type SendDocumentRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Document string `json:"document" binding:"required" example:"data:application/pdf;base64,JVBERi0x..."` // Base64 data URL or HTTP URL
	FileName string `json:"filename,omitempty" example:"document.pdf"`
	MimeType string `json:"mimetype,omitempty" example:"application/pdf"`
}

type SendStickerRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Sticker string `json:"sticker" binding:"required" example:"data:image/webp;base64,UklGRnoGAABXRUJQ..."` // Base64 data URL or HTTP URL
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

func (r *SendTextRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Body == "" {
		return &MessageValidationError{Field: "body", Message: "Message body is required"}
	}
	return nil
}

func (r *SendMediaRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
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

func (r *SendLocationRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Latitude < -90 || r.Latitude > 90 {
		return &MessageValidationError{Field: "latitude", Message: "Latitude must be between -90 and 90"}
	}
	if r.Longitude < -180 || r.Longitude > 180 {
		return &MessageValidationError{Field: "longitude", Message: "Longitude must be between -180 and 180"}
	}
	return nil
}

func (r *SendContactRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.ContactName == "" {
		return &MessageValidationError{Field: "contact_name", Message: "Contact name is required"}
	}
	if r.ContactPhone == "" {
		return &MessageValidationError{Field: "contact_phone", Message: "Contact phone is required"}
	}
	return nil
}

func (r *SendImageRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Image == "" {
		return &MessageValidationError{Field: "image", Message: "Image data is required"}
	}
	return nil
}

func (r *SendAudioRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Audio == "" {
		return &MessageValidationError{Field: "audio", Message: "Audio data is required"}
	}
	return nil
}

func (r *SendVideoRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Video == "" {
		return &MessageValidationError{Field: "video", Message: "Video data is required"}
	}
	return nil
}

func (r *SendDocumentRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Document == "" {
		return &MessageValidationError{Field: "document", Message: "Document data is required"}
	}
	return nil
}

func (r *SendStickerRequest) Validate() error {
	if !validateMessagePhone(r.Phone) {
		return &MessageValidationError{Field: "phone", Message: "Invalid phone number format"}
	}
	if r.Sticker == "" {
		return &MessageValidationError{Field: "sticker", Message: "Sticker data is required"}
	}
	return nil
}
