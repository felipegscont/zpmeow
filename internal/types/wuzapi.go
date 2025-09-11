package types

import (
	"fmt"
	"go.mau.fi/whatsmeow"
)


type SendTextRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Body        string      `json:"body" binding:"required" example:"Hello, World!"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}


type SendImageRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Image       string      `json:"image" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."`
	Caption     string      `json:"caption,omitempty" example:"Image caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	MimeType    string      `json:"mimeType,omitempty" example:"image/jpeg"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}


type SendAudioRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Audio       string      `json:"audio" binding:"required" example:"data:audio/ogg;base64,T2dnU..."`
	Caption     string      `json:"caption,omitempty" example:"Audio caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}


type SendDocumentRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Document    string      `json:"document" binding:"required" example:"data:application/pdf;base64,JVBERi0x..."`
	Filename    string      `json:"filename,omitempty" example:"document.pdf"`
	Caption     string      `json:"caption,omitempty" example:"Document caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}


type SendVideoRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Video       string      `json:"video" binding:"required" example:"data:video/mp4;base64,AAAAIGZ0eXA..."`
	Caption     string      `json:"caption,omitempty" example:"Video caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}


type SendStickerRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Sticker     string      `json:"sticker" binding:"required" example:"data:image/webp;base64,UklGRv4..."`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}



type SendMediaRequest struct {
	Phone       string      `json:"phone" form:"phone" binding:"required" example:"+5511999999999"`
	MediaType   string      `json:"mediaType" form:"mediaType" binding:"required" example:"image" enums:"image,audio,document,video"`
	Media       string      `json:"media" form:"media" example:"data:image/jpeg;base64,/9j/4AAQ..." swaggertype:"string"`
	Caption     string      `json:"caption" form:"caption" example:"Media caption"`
	Filename    string      `json:"filename" form:"filename" example:"document.pdf"`
	ID          string      `json:"id" form:"id" example:"custom-message-id"`
	MimeType    string      `json:"mimeType" form:"mimeType" example:"image/jpeg"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}


type SendLocationRequest struct {
	Phone     string  `json:"phone" binding:"required" example:"+5511999999999"`
	Latitude  float64 `json:"latitude" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	Address   string  `json:"address,omitempty" example:"São Paulo, SP, Brazil"`
	ID        string  `json:"id,omitempty" example:"custom-message-id"`
}


type SendContactRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Contact     Contact     `json:"contact" binding:"required"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}


type SendButtonsRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"+5511999999999"`
	Text       string   `json:"text" binding:"required" example:"Choose an option:"`
	Buttons    []Button `json:"buttons" binding:"required"`
	Footer     string   `json:"footer,omitempty" example:"Footer text"`
	ID         string   `json:"id,omitempty" example:"custom-message-id"`
}


type SendListRequest struct {
	Phone      string    `json:"phone" binding:"required" example:"+5511999999999"`
	Text       string    `json:"text" binding:"required" example:"Choose from the list:"`
	ButtonText string    `json:"buttonText" binding:"required" example:"Select Option"`
	Sections   []Section `json:"sections" binding:"required"`
	Footer     string    `json:"footer,omitempty" example:"Footer text"`
	ID         string    `json:"id,omitempty" example:"custom-message-id"`
}


type SendPollRequest struct {
	Phone           string   `json:"phone" binding:"required" example:"+5511999999999"`
	Name            string   `json:"name" binding:"required" example:"What's your favorite color?"`
	Options         []string `json:"options" binding:"required" example:"Red,Blue,Green"`
	SelectableCount int      `json:"selectableCount,omitempty" example:"1"`
	ID              string   `json:"id,omitempty" example:"custom-message-id"`
}


type ContextInfo struct {
	StanzaID      string `json:"stanzaId,omitempty" example:"3EB0C431C26A1916E07A"`
	Participant   string `json:"participant,omitempty" example:"+5511888888888@s.whatsapp.net"`
	QuotedMessage any    `json:"quotedMessage,omitempty"`
}


type Contact struct {
	DisplayName string `json:"displayName" binding:"required" example:"John Doe"`
	VCard       string `json:"vcard" binding:"required" example:"BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"`
}


type Button struct {
	ButtonID   string     `json:"buttonId" binding:"required" example:"btn_1"`
	ButtonText ButtonText `json:"buttonText" binding:"required"`
	Type       int        `json:"type" example:"1"`
}


type ButtonText struct {
	DisplayText string `json:"displayText" binding:"required" example:"Option 1"`
}


type Section struct {
	Title string `json:"title" binding:"required" example:"Section 1"`
	Rows  []Row  `json:"rows" binding:"required"`
}


type Row struct {
	Title       string `json:"title" binding:"required" example:"Row 1"`
	Description string `json:"description,omitempty" example:"Row description"`
	RowID       string `json:"rowId" binding:"required" example:"row_1"`
}


type SendResponse struct {

	Timestamp int64 `json:"timestamp" example:"1640995200"`


	ID string `json:"id" example:"3EB0C431C26A1916E07A"`


	ServerID string `json:"serverId,omitempty" example:"wamid.HBgNNTU5OTgxNzY5NTM2FQIAERgSMzNFNzE4QzY5QzE5MjE2RTdB"`


	Sender string `json:"sender,omitempty" example:"5511999999999@s.whatsapp.net"`


	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"messageId" example:"3EB0C431C26A1916E07A"`
}

// Webhook types
type SetWebhookRequest struct {
	WebhookURL string   `json:"webhookurl" binding:"required" example:"https://example.com/webhook"`
	Events     []string `json:"events,omitempty" example:"message,status"`
}

type UpdateWebhookRequest struct {
	WebhookURL string   `json:"webhook" binding:"required" example:"https://example.com/webhook"`
	Events     []string `json:"events,omitempty" example:"message,status"`
	Active     bool     `json:"active" example:"true"`
}

type WebhookResponse struct {
	Webhook   string   `json:"webhook" example:"https://example.com/webhook"`
	Events    []string `json:"events,omitempty" example:"message,status"`
	Active    bool     `json:"active,omitempty" example:"true"`
	Subscribe []string `json:"subscribe,omitempty" example:"message,status"`
}

// User types
type UserPresenceRequest struct {
	Type string `json:"type" binding:"required" example:"available"` // available, unavailable
}

type CheckUserRequest struct {
	Phone []string `json:"phone" binding:"required" example:"+5511999999999,+5511888888888"`
}

type GetUserInfoRequest struct {
	Phone []string `json:"phone" binding:"required" example:"+5511999999999,+5511888888888"`
}

type GetAvatarRequest struct {
	Phone   string `json:"phone" binding:"required" example:"+5511999999999"`
	Preview bool   `json:"preview,omitempty" example:"false"`
}

type UserCheckResult struct {
	Query        string `json:"query" example:"+5511999999999"`
	IsInWhatsapp bool   `json:"isInWhatsapp" example:"true"`
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	VerifiedName string `json:"verifiedName,omitempty" example:"John Doe"`
}

type UserCheckResponse struct {
	Users []UserCheckResult `json:"users"`
}

type AvatarResponse struct {
	URL       string `json:"url,omitempty" example:"https://example.com/avatar.jpg"`
	ID        string `json:"id,omitempty" example:"avatar_id_123"`
	Type      string `json:"type,omitempty" example:"image"`
	DirectURL string `json:"directUrl,omitempty" example:"https://example.com/direct_avatar.jpg"`
}

// Newsletter types
type NewsletterResponse struct {
	Newsletter []NewsletterMetadata `json:"newsletter"`
}

type NewsletterMetadata struct {
	ID          string `json:"id" example:"newsletter_123"`
	Name        string `json:"name" example:"Tech News"`
	Description string `json:"description,omitempty" example:"Latest technology updates"`
	Handle      string `json:"handle,omitempty" example:"@technews"`
	Picture     string `json:"picture,omitempty" example:"https://example.com/newsletter.jpg"`
	Preview     string `json:"preview,omitempty" example:"https://example.com/preview.jpg"`
	Reaction    string `json:"reaction,omitempty" example:"👍"`
	Verified    bool   `json:"verified,omitempty" example:"true"`
}


func NewSendResponseFromWhatsmeow(resp *whatsmeow.SendResponse, requestID string) SendResponse {
	return SendResponse{
		Timestamp: resp.Timestamp.Unix(),
		ID:        string(resp.ID),
		ServerID:  func() string {
			if resp.ServerID == 0 {
				return ""
			}
			return fmt.Sprintf("%d", resp.ServerID)
		}(),
		Sender:    resp.Sender.String(),
		Success:   true,
		MessageID: func() string {
			if requestID != "" {
				return requestID
			}
			return string(resp.ID)
		}(),
	}
}




type ChatPresenceRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	State string `json:"state" binding:"required" example:"composing"` // composing, paused
	Media string `json:"media" example:"audio"`                        // text, audio (optional, only for composing)
}


type ChatMarkReadRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"+5511999999999"`
	MessageIDs []string `json:"messageIds" binding:"required" example:"3EB0C431C26A1916E07A,3EB0C431C26A1916E07B"`
}


type ChatReactRequest struct {
	Phone     string `json:"phone" binding:"required" example:"+5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	Emoji     string `json:"emoji" binding:"required" example:"👍"`
}


type ChatDeleteRequest struct {
	Phone       string `json:"phone" binding:"required" example:"+5511999999999"`
	MessageID   string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	ForEveryone bool   `json:"forEveryone,omitempty" example:"false"`
}


type ChatEditRequest struct {
	Phone     string `json:"phone" binding:"required" example:"+5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	NewText   string `json:"newText" binding:"required" example:"Edited message"`
}


type ChatDownloadRequest struct {
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	Phone     string `json:"phone,omitempty" example:"+5511999999999"`
}




type GroupCreateRequest struct {
	Name         string   `json:"name" binding:"required" example:"My Group"`
	Participants []string `json:"participants" binding:"required" example:"+5511999999999,+5511888888888"`
}


type GroupJoinRequest struct {
	InviteCode string `json:"inviteCode" binding:"required" example:"CjQKOAokMjU5NzE4NzAtNzBiYy00"`
}


type GroupLeaveRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
}


type GroupUpdateParticipantsRequest struct {
	GroupJID     string   `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Participants []string `json:"participants" binding:"required" example:"+5511999999999,+5511888888888"`
	Action       string   `json:"action" binding:"required" example:"add"` // add, remove, promote, demote
}


type GroupSetNameRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Name     string `json:"name" binding:"required" example:"New Group Name"`
}


type GroupSetTopicRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Topic    string `json:"topic" binding:"required" example:"New group description"`
}


type GroupSetPhotoRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Image    string `json:"image" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."`
}


type GroupRemovePhotoRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
}


type GroupSetAnnounceRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Announce bool   `json:"announce" binding:"required" example:"true"`
}


type GroupSetLockedRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Locked   bool   `json:"locked" binding:"required" example:"true"`
}


type GroupSetEphemeralRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Duration int64  `json:"duration" binding:"required" example:"86400"` // seconds
}


type GroupInviteInfoRequest struct {
	InviteCode string `json:"inviteCode" binding:"required" example:"CjQKOAokMjU5NzE4NzAtNzBiYy00"`
}
