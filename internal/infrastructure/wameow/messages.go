package wameow

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
)

// Pure functions for message sending - no structs needed
// Following DRY principle with shared helper functions

// Common message sending helpers to reduce duplication

// sendMessageToJID - Helper function to send any message to a JID
func sendMessageToJID(client *whatsmeow.Client, to string, message *waProto.Message) (*whatsmeow.SendResponse, error) {
	jid, err := parsePhoneToJID(to)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendMessage(context.Background(), jid, message)
	return &resp, err
}

// createMediaMessage - Helper function to create media message with uploaded content
func createMediaMessage(client *whatsmeow.Client, data []byte, mediaType whatsmeow.MediaType) (*whatsmeow.UploadResponse, error) {
	return uploadMedia(client, data, mediaType)
}

// validateMessageInput - Helper function to validate common message inputs
func validateMessageInput(client *whatsmeow.Client, to string) error {
	if client == nil {
		return fmt.Errorf("client cannot be nil")
	}
	if to == "" {
		return fmt.Errorf("recipient cannot be empty")
	}
	return nil
}

// SendTextMessage sends a text message using pure function approach
func SendTextMessage(client *whatsmeow.Client, to, text string) (*whatsmeow.SendResponse, error) {
	if err := validateMessageInput(client, to); err != nil {
		return nil, err
	}

	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	message := &waProto.Message{
		Conversation: &text,
	}

	return sendMessageToJID(client, to, message)
}

// SendImageMessage sends an image message using pure function approach
func SendImageMessage(client *whatsmeow.Client, to string, data []byte, caption string) (*whatsmeow.SendResponse, error) {
	if err := validateMessageInput(client, to); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("image data cannot be empty")
	}

	uploaded, err := createMediaMessage(client, data, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	mimeType := "image/jpeg"
	message := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:       &caption,
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
			Mimetype:      &mimeType,
		},
	}

	return sendMessageToJID(client, to, message)
}

// SendAudioMessage sends an audio message using pure function approach
func SendAudioMessage(client *whatsmeow.Client, to string, data []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	if err := validateMessageInput(client, to); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("audio data cannot be empty")
	}

	uploaded, err := createMediaMessage(client, data, whatsmeow.MediaAudio)
	if err != nil {
		return nil, err
	}

	message := &waProto.Message{
		AudioMessage: &waProto.AudioMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
		},
	}

	return sendMessageToJID(client, to, message)
}

// SendVideoMessage sends a video message using pure function approach
func SendVideoMessage(client *whatsmeow.Client, to string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	if err := validateMessageInput(client, to); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("video data cannot be empty")
	}

	uploaded, err := createMediaMessage(client, data, whatsmeow.MediaVideo)
	if err != nil {
		return nil, err
	}

	message := &waProto.Message{
		VideoMessage: &waProto.VideoMessage{
			Caption:       &caption,
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
		},
	}

	return sendMessageToJID(client, to, message)
}

// SendDocumentMessage sends a document message using pure function approach
func SendDocumentMessage(client *whatsmeow.Client, to string, data []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	jid, err := parsePhoneToJID(to)
	if err != nil {
		return nil, err
	}

	uploaded, err := uploadMedia(client, data, whatsmeow.MediaDocument)
	if err != nil {
		return nil, err
	}

	message := &waProto.Message{
		DocumentMessage: &waProto.DocumentMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
			FileName:      &filename,
			Caption:       &caption,
		},
	}

	resp, err := client.SendMessage(context.Background(), jid, message)
	return &resp, err
}

// SendStickerMessage sends a sticker message using pure function approach
func SendStickerMessage(client *whatsmeow.Client, to string, data []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	jid, err := parsePhoneToJID(to)
	if err != nil {
		return nil, err
	}

	uploaded, err := uploadMedia(client, data, whatsmeow.MediaImage) // Stickers use image media type
	if err != nil {
		return nil, err
	}

	message := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
		},
	}

	resp, err := client.SendMessage(context.Background(), jid, message)
	return &resp, err
}

// SendContactMessage sends a contact message using pure function approach
func SendContactMessage(client *whatsmeow.Client, to, contactName, contactPhone string) (*whatsmeow.SendResponse, error) {
	jid, err := parsePhoneToJID(to)
	if err != nil {
		return nil, err
	}

	vcard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL;type=CELL;type=VOICE;waid=%s:+%s\nEND:VCARD", contactName, contactPhone, contactPhone)

	message := &waProto.Message{
		ContactMessage: &waProto.ContactMessage{
			DisplayName: &contactName,
			Vcard:       &vcard,
		},
	}

	resp, err := client.SendMessage(context.Background(), jid, message)
	return &resp, err
}

// SendLocationMessage sends a location message using pure function approach
func SendLocationMessage(client *whatsmeow.Client, to string, latitude, longitude float64, name, address string) (*whatsmeow.SendResponse, error) {
	jid, err := parsePhoneToJID(to)
	if err != nil {
		return nil, err
	}

	message := &waProto.Message{
		LocationMessage: &waProto.LocationMessage{
			DegreesLatitude:  &latitude,
			DegreesLongitude: &longitude,
			Name:             &name,
			Address:          &address,
		},
	}

	resp, err := client.SendMessage(context.Background(), jid, message)
	return &resp, err
}

// Additional message functions can be implemented here as pure functions
// Following the same pattern as above

// All remaining MessageSender methods removed - implement as pure functions when needed

// Helper functions shared across message types (DRY principle)
func parsePhoneToJID(phone string) (waTypes.JID, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return waTypes.EmptyJID, fmt.Errorf("phone number cannot be empty")
	}

	// Remove leading + if present
	if phone[0] == '+' {
		phone = phone[1:]
	}

	// Extract only digits
	var digits strings.Builder
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	formattedPhone := digits.String()

	// Validate phone number
	if formattedPhone == "" {
		return waTypes.EmptyJID, fmt.Errorf("phone number cannot be empty")
	}

	if len(formattedPhone) < 7 || len(formattedPhone) > 15 {
		return waTypes.EmptyJID, fmt.Errorf("phone number must be between 7 and 15 digits")
	}

	if formattedPhone[0] == '0' {
		return waTypes.EmptyJID, fmt.Errorf("phone number should not start with 0")
	}

	return waTypes.NewJID(formattedPhone, waTypes.DefaultUserServer), nil
}

func uploadMedia(client *whatsmeow.Client, data []byte, mediaType whatsmeow.MediaType) (*whatsmeow.UploadResponse, error) {
	resp, err := client.Upload(context.Background(), data, mediaType)
	return &resp, err
}
