package commands

import (
	"fmt"
	"strings"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVideo    MessageType = "video"
	MessageTypeDocument MessageType = "document"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
	MessageTypeSticker  MessageType = "sticker"
)

// SendTextCommand represents a command to send a text message
type SendTextCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	Content   string `json:"content" validate:"required"`
	ReplyTo   string `json:"reply_to,omitempty"`
}

// SendMediaCommand represents a command to send a media message
type SendMediaCommand struct {
	SessionID string      `json:"session_id" validate:"required"`
	ChatJID   string      `json:"chat_jid" validate:"required"`
	MediaType MessageType `json:"media_type" validate:"required"`
	MediaURL  string      `json:"media_url,omitempty"`
	MediaData []byte      `json:"media_data,omitempty"`
	Caption   string      `json:"caption,omitempty"`
	Filename  string      `json:"filename,omitempty"`
	ReplyTo   string      `json:"reply_to,omitempty"`
}

// SendLocationCommand represents a command to send a location message
type SendLocationCommand struct {
	SessionID string  `json:"session_id" validate:"required"`
	ChatJID   string  `json:"chat_jid" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// SendContactCommand represents a command to send a contact message
type SendContactCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	VCard     string `json:"vcard" validate:"required"`
}

// SendImageCommand represents a command to send an image
type SendImageCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	ImageData []byte `json:"image_data,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	Caption   string `json:"caption,omitempty"`
	ViewOnce  bool   `json:"view_once,omitempty"`
}

// SendAudioCommand represents a command to send an audio message
type SendAudioCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	AudioData []byte `json:"audio_data,omitempty"`
	AudioURL  string `json:"audio_url,omitempty"`
	Ptt       bool   `json:"ptt,omitempty"` // Push to talk
}

// SendVideoCommand represents a command to send a video message
type SendVideoCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	VideoData []byte `json:"video_data,omitempty"`
	VideoURL  string `json:"video_url,omitempty"`
	Caption   string `json:"caption,omitempty"`
	ViewOnce  bool   `json:"view_once,omitempty"`
}

// SendDocumentCommand represents a command to send a document
type SendDocumentCommand struct {
	SessionID    string `json:"session_id" validate:"required"`
	ChatJID      string `json:"chat_jid" validate:"required"`
	DocumentData []byte `json:"document_data,omitempty"`
	DocumentURL  string `json:"document_url,omitempty"`
	Filename     string `json:"filename" validate:"required"`
	Mimetype     string `json:"mimetype,omitempty"`
}

// SendStickerCommand represents a command to send a sticker
type SendStickerCommand struct {
	SessionID   string `json:"session_id" validate:"required"`
	ChatJID     string `json:"chat_jid" validate:"required"`
	StickerData []byte `json:"sticker_data,omitempty"`
	StickerURL  string `json:"sticker_url,omitempty"`
}

// Validation methods

func (c *SendTextCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if strings.TrimSpace(c.Content) == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}

func (c *SendMediaCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.MediaType == "" {
		return fmt.Errorf("media type is required")
	}
	if c.MediaURL == "" && len(c.MediaData) == 0 {
		return fmt.Errorf("either media URL or media data is required")
	}
	return nil
}

func (c *SendLocationCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.Latitude == 0 && c.Longitude == 0 {
		return fmt.Errorf("latitude and longitude are required")
	}
	return nil
}

func (c *SendContactCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if strings.TrimSpace(c.VCard) == "" {
		return fmt.Errorf("vCard is required")
	}
	return nil
}

func (c *SendImageCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.ImageURL == "" && len(c.ImageData) == 0 {
		return fmt.Errorf("either image URL or image data is required")
	}
	return nil
}

func (c *SendAudioCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.AudioURL == "" && len(c.AudioData) == 0 {
		return fmt.Errorf("either audio URL or audio data is required")
	}
	return nil
}

func (c *SendVideoCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.VideoURL == "" && len(c.VideoData) == 0 {
		return fmt.Errorf("either video URL or video data is required")
	}
	return nil
}

func (c *SendDocumentCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.DocumentURL == "" && len(c.DocumentData) == 0 {
		return fmt.Errorf("either document URL or document data is required")
	}
	if strings.TrimSpace(c.Filename) == "" {
		return fmt.Errorf("filename is required")
	}
	return nil
}

func (c *SendStickerCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.StickerURL == "" && len(c.StickerData) == 0 {
		return fmt.Errorf("either sticker URL or sticker data is required")
	}
	return nil
}

// FormatJID formats a JID to include the WhatsApp domain if not present
func FormatJID(jid string) string {
	jid = strings.TrimSpace(jid)
	if strings.Contains(jid, "@") {
		return jid
	}
	return jid + "@s.whatsapp.net"
}

// MarkAsReadCommand represents a command to mark a message as read
type MarkAsReadCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	MessageID string `json:"message_id" validate:"required"`
}

func (c *MarkAsReadCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	return nil
}

// ReactToMessageCommand represents a command to react to a message
type ReactToMessageCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	MessageID string `json:"message_id" validate:"required"`
	Reaction  string `json:"reaction" validate:"required"`
}

func (c *ReactToMessageCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	if c.Reaction == "" {
		return fmt.Errorf("reaction is required")
	}
	return nil
}

// EditMessageCommand represents a command to edit a message
type EditMessageCommand struct {
	SessionID  string `json:"session_id" validate:"required"`
	ChatJID    string `json:"chat_jid" validate:"required"`
	MessageID  string `json:"message_id" validate:"required"`
	NewContent string `json:"new_content" validate:"required"`
}

func (c *EditMessageCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	if strings.TrimSpace(c.NewContent) == "" {
		return fmt.Errorf("new content is required")
	}
	return nil
}

// DeleteMessageCommand represents a command to delete a message
type DeleteMessageCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	ChatJID   string `json:"chat_jid" validate:"required"`
	MessageID string `json:"message_id" validate:"required"`
}

func (c *DeleteMessageCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.ChatJID == "" {
		return fmt.Errorf("chat JID is required")
	}
	if c.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	return nil
}
