package meow

import (
	"context"
	"errors"
	"fmt"

	"zpmeow/internal/types"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	waTypes "go.mau.fi/whatsmeow/types"
)

// MessageBuilder provides utilities for building WhatsApp messages
type MessageBuilder struct{}

// BuildTextMessage creates a text message
func (mb *MessageBuilder) BuildTextMessage(text string, contextInfo *waE2E.ContextInfo) *waE2E.Message {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &text,
		},
	}

	if contextInfo != nil {
		msg.ExtendedTextMessage.ContextInfo = contextInfo
	}

	return msg
}

// BuildLocationMessage creates a location message
func (mb *MessageBuilder) BuildLocationMessage(latitude, longitude float64, name, address string) *waE2E.Message {
	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  &latitude,
			DegreesLongitude: &longitude,
		},
	}

	if name != "" {
		msg.LocationMessage.Name = &name
	}
	if address != "" {
		msg.LocationMessage.Address = &address
	}

	return msg
}

// BuildContactMessage creates a contact message
func (mb *MessageBuilder) BuildContactMessage(displayName, vcard string) *waE2E.Message {
	return &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: &displayName,
			Vcard:       &vcard,
		},
	}
}

// MediaMessageParams holds parameters for media messages
type MediaMessageParams struct {
	UploadResponse whatsmeow.UploadResponse
	Caption        string
	MimeType       string
	FileName       string
}

// BuildImageMessage creates an image message
func (mb *MessageBuilder) BuildImageMessage(params MediaMessageParams) *waE2E.Message {
	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           &params.UploadResponse.URL,
			DirectPath:    &params.UploadResponse.DirectPath,
			MediaKey:      params.UploadResponse.MediaKey,
			Mimetype:      &params.MimeType,
			FileEncSHA256: params.UploadResponse.FileEncSHA256,
			FileSHA256:    params.UploadResponse.FileSHA256,
			FileLength:    &params.UploadResponse.FileLength,
		},
	}

	if params.Caption != "" {
		msg.ImageMessage.Caption = &params.Caption
	}

	return msg
}

// BuildAudioMessage creates an audio message
func (mb *MessageBuilder) BuildAudioMessage(params MediaMessageParams, isPTT bool) *waE2E.Message {
	mimeType := params.MimeType
	if mimeType == "" {
		mimeType = "audio/ogg; codecs=opus"
	}

	return &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           &params.UploadResponse.URL,
			DirectPath:    &params.UploadResponse.DirectPath,
			MediaKey:      params.UploadResponse.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: params.UploadResponse.FileEncSHA256,
			FileSHA256:    params.UploadResponse.FileSHA256,
			FileLength:    &params.UploadResponse.FileLength,
			PTT:           &isPTT,
		},
	}
}

// BuildDocumentMessage creates a document message
func (mb *MessageBuilder) BuildDocumentMessage(params MediaMessageParams) *waE2E.Message {
	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           &params.UploadResponse.URL,
			DirectPath:    &params.UploadResponse.DirectPath,
			MediaKey:      params.UploadResponse.MediaKey,
			Mimetype:      &params.MimeType,
			FileEncSHA256: params.UploadResponse.FileEncSHA256,
			FileSHA256:    params.UploadResponse.FileSHA256,
			FileLength:    &params.UploadResponse.FileLength,
			FileName:      &params.FileName,
		},
	}

	if params.Caption != "" {
		msg.DocumentMessage.Caption = &params.Caption
	}

	return msg
}

// BuildVideoMessage creates a video message
func (mb *MessageBuilder) BuildVideoMessage(params MediaMessageParams) *waE2E.Message {
	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           &params.UploadResponse.URL,
			DirectPath:    &params.UploadResponse.DirectPath,
			MediaKey:      params.UploadResponse.MediaKey,
			Mimetype:      &params.MimeType,
			FileEncSHA256: params.UploadResponse.FileEncSHA256,
			FileSHA256:    params.UploadResponse.FileSHA256,
			FileLength:    &params.UploadResponse.FileLength,
		},
	}

	if params.Caption != "" {
		msg.VideoMessage.Caption = &params.Caption
	}

	return msg
}

// BuildStickerMessage creates a sticker message
func (mb *MessageBuilder) BuildStickerMessage(params MediaMessageParams) *waE2E.Message {
	return &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           &params.UploadResponse.URL,
			DirectPath:    &params.UploadResponse.DirectPath,
			MediaKey:      params.UploadResponse.MediaKey,
			Mimetype:      &params.MimeType,
			FileEncSHA256: params.UploadResponse.FileEncSHA256,
			FileSHA256:    params.UploadResponse.FileSHA256,
			FileLength:    &params.UploadResponse.FileLength,
		},
	}
}

// BuildPollMessage creates a poll message
func (mb *MessageBuilder) BuildPollMessage(name string, options []string, selectableCount int) *waE2E.Message {
	pollOptions := make([]*waE2E.PollCreationMessage_Option, len(options))
	for i, option := range options {
		pollOptions[i] = &waE2E.PollCreationMessage_Option{
			OptionName: &option,
		}
	}

	if selectableCount <= 0 {
		selectableCount = 1
	}

	return &waE2E.Message{
		PollCreationMessage: &waE2E.PollCreationMessage{
			Name:    &name,
			Options: pollOptions,
		},
	}
}



// BuildButtonsMessage creates a buttons message (fallback to text with numbered options)
func (mb *MessageBuilder) BuildButtonsMessage(text string, buttons []types.Button, footer string) *waE2E.Message {
	// WhatsApp deprecated interactive buttons, so we fallback to text with numbered options
	buttonText := text + "\n\n"
	for i, button := range buttons {
		buttonText += fmt.Sprintf("%d. %s\n", i+1, button.ButtonText.DisplayText)
	}
	if footer != "" {
		buttonText += "\n" + footer
	}

	return &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &buttonText,
		},
	}
}

// BuildListMessage creates a list message (fallback to text with sections)
func (mb *MessageBuilder) BuildListMessage(text, buttonText string, sections []types.Section, footer string) *waE2E.Message {
	// WhatsApp deprecated interactive lists, so we fallback to text with sections
	listText := text + "\n\n"
	for _, section := range sections {
		if section.Title != "" {
			listText += "📋 " + section.Title + "\n"
		}
		for i, row := range section.Rows {
			listText += fmt.Sprintf("%d. %s", i+1, row.Title)
			if row.Description != "" {
				listText += " - " + row.Description
			}
			listText += "\n"
		}
		listText += "\n"
	}
	if footer != "" {
		listText += footer
	}

	return &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &listText,
		},
	}
}

// MediaUploader handles media upload operations
type MediaUploader struct {
	client *whatsmeow.Client
}

// NewMediaUploader creates a new media uploader
func NewMediaUploader(client *whatsmeow.Client) *MediaUploader {
	return &MediaUploader{client: client}
}

// UploadMedia uploads media data and returns upload response
func (mu *MediaUploader) UploadMedia(ctx context.Context, data []byte, mediaType whatsmeow.MediaType) (whatsmeow.UploadResponse, error) {
	if mu.client == nil {
		return whatsmeow.UploadResponse{}, errors.New(ErrClientNotFound)
	}

	uploaded, err := mu.client.Upload(ctx, data, mediaType)
	if err != nil {
		return whatsmeow.UploadResponse{}, Error.WrapError(err, "failed to upload media")
	}

	return uploaded, nil
}

// MessageSender handles message sending operations
type MessageSender struct {
	client *whatsmeow.Client
}

// NewMessageSender creates a new message sender
func NewMessageSender(client *whatsmeow.Client) *MessageSender {
	return &MessageSender{client: client}
}

// SendMessage sends a message to the specified recipient and returns the whatsmeow SendResponse
func (ms *MessageSender) SendMessage(ctx context.Context, to waTypes.JID, message *waE2E.Message) (*whatsmeow.SendResponse, error) {
	if ms.client == nil {
		return nil, errors.New(ErrClientNotFound)
	}

	// Skip connection check to avoid potential deadlock - validation is done at service level

	resp, err := ms.client.SendMessage(ctx, to, message)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send message")
	}

	return &resp, nil
}

// Global message utilities
var (
	MsgBuilder = &MessageBuilder{}
)
