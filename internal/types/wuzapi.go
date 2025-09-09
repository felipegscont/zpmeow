package types

// SendTextRequest represents a text message send request
type SendTextRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Body        string      `json:"body" binding:"required" example:"Hello, World!"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendImageRequest represents an image message send request
type SendImageRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Image       string      `json:"image" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."`
	Caption     string      `json:"caption,omitempty" example:"Image caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	MimeType    string      `json:"mimeType,omitempty" example:"image/jpeg"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendAudioRequest represents an audio message send request
type SendAudioRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Audio       string      `json:"audio" binding:"required" example:"data:audio/ogg;base64,T2dnU..."`
	Caption     string      `json:"caption,omitempty" example:"Audio caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendDocumentRequest represents a document message send request
type SendDocumentRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Document    string      `json:"document" binding:"required" example:"data:application/pdf;base64,JVBERi0x..."`
	Filename    string      `json:"filename,omitempty" example:"document.pdf"`
	Caption     string      `json:"caption,omitempty" example:"Document caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendVideoRequest represents a video message send request
type SendVideoRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Video       string      `json:"video" binding:"required" example:"data:video/mp4;base64,AAAAIGZ0eXA..."`
	Caption     string      `json:"caption,omitempty" example:"Video caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendStickerRequest represents a sticker message send request
type SendStickerRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Sticker     string      `json:"sticker" binding:"required" example:"data:image/webp;base64,UklGRv4..."`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendLocationRequest represents a location message send request
type SendLocationRequest struct {
	Phone     string  `json:"phone" binding:"required" example:"+5511999999999"`
	Latitude  float64 `json:"latitude" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	Address   string  `json:"address,omitempty" example:"São Paulo, SP, Brazil"`
	ID        string  `json:"id,omitempty" example:"custom-message-id"`
}

// SendContactRequest represents a contact message send request
type SendContactRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Contact     Contact     `json:"contact" binding:"required"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendButtonsRequest represents an interactive buttons message send request
type SendButtonsRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"+5511999999999"`
	Text       string   `json:"text" binding:"required" example:"Choose an option:"`
	Buttons    []Button `json:"buttons" binding:"required"`
	Footer     string   `json:"footer,omitempty" example:"Footer text"`
	ID         string   `json:"id,omitempty" example:"custom-message-id"`
}

// SendListRequest represents an interactive list message send request
type SendListRequest struct {
	Phone      string    `json:"phone" binding:"required" example:"+5511999999999"`
	Text       string    `json:"text" binding:"required" example:"Choose from the list:"`
	ButtonText string    `json:"buttonText" binding:"required" example:"Select Option"`
	Sections   []Section `json:"sections" binding:"required"`
	Footer     string    `json:"footer,omitempty" example:"Footer text"`
	ID         string    `json:"id,omitempty" example:"custom-message-id"`
}

// SendPollRequest represents a poll message send request
type SendPollRequest struct {
	Phone           string   `json:"phone" binding:"required" example:"+5511999999999"`
	Name            string   `json:"name" binding:"required" example:"What's your favorite color?"`
	Options         []string `json:"options" binding:"required" example:"Red,Blue,Green"`
	SelectableCount int      `json:"selectableCount,omitempty" example:"1"`
	ID              string   `json:"id,omitempty" example:"custom-message-id"`
}

// ContextInfo represents message context for replies/quotes
type ContextInfo struct {
	StanzaID      string `json:"stanzaId,omitempty" example:"3EB0C431C26A1916E07A"`
	Participant   string `json:"participant,omitempty" example:"+5511888888888@s.whatsapp.net"`
	QuotedMessage any    `json:"quotedMessage,omitempty"`
}

// Contact represents a contact structure
type Contact struct {
	DisplayName string `json:"displayName" binding:"required" example:"John Doe"`
	VCard       string `json:"vcard" binding:"required" example:"BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"`
}

// Button represents an interactive button
type Button struct {
	ButtonID   string     `json:"buttonId" binding:"required" example:"btn_1"`
	ButtonText ButtonText `json:"buttonText" binding:"required"`
	Type       int        `json:"type" example:"1"`
}

// ButtonText represents button text structure
type ButtonText struct {
	DisplayText string `json:"displayText" binding:"required" example:"Option 1"`
}

// Section represents a list section
type Section struct {
	Title string `json:"title" binding:"required" example:"Section 1"`
	Rows  []Row  `json:"rows" binding:"required"`
}

// Row represents a list row
type Row struct {
	Title       string `json:"title" binding:"required" example:"Row 1"`
	Description string `json:"description,omitempty" example:"Row description"`
	RowID       string `json:"rowId" binding:"required" example:"row_1"`
}

// SendResponse represents a successful send response
type SendResponse struct {
	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"messageId" example:"3EB0C431C26A1916E07A"`
	Timestamp int64  `json:"timestamp" example:"1640995200"`
}

// Chat operation types

// ChatPresenceRequest represents a chat presence request
type ChatPresenceRequest struct {
	Phone string `json:"phone" binding:"required" example:"+5511999999999"`
	State string `json:"state" binding:"required" example:"typing"` // typing, recording, paused
}

// ChatMarkReadRequest represents a mark as read request
type ChatMarkReadRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"+5511999999999"`
	MessageIDs []string `json:"messageIds" binding:"required" example:"3EB0C431C26A1916E07A,3EB0C431C26A1916E07B"`
}

// ChatReactRequest represents a message reaction request
type ChatReactRequest struct {
	Phone     string `json:"phone" binding:"required" example:"+5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	Emoji     string `json:"emoji" binding:"required" example:"👍"`
}

// ChatDeleteRequest represents a message deletion request
type ChatDeleteRequest struct {
	Phone       string `json:"phone" binding:"required" example:"+5511999999999"`
	MessageID   string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	ForEveryone bool   `json:"forEveryone,omitempty" example:"false"`
}

// ChatEditRequest represents a message edit request
type ChatEditRequest struct {
	Phone     string `json:"phone" binding:"required" example:"+5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	NewText   string `json:"newText" binding:"required" example:"Edited message"`
}

// ChatDownloadRequest represents a media download request
type ChatDownloadRequest struct {
	MessageID string `json:"messageId" binding:"required" example:"3EB0C431C26A1916E07A"`
	Phone     string `json:"phone,omitempty" example:"+5511999999999"`
}

// Group operation types

// GroupCreateRequest represents a group creation request
type GroupCreateRequest struct {
	Name         string   `json:"name" binding:"required" example:"My Group"`
	Participants []string `json:"participants" binding:"required" example:"+5511999999999,+5511888888888"`
}

// GroupJoinRequest represents a group join request
type GroupJoinRequest struct {
	InviteCode string `json:"inviteCode" binding:"required" example:"CjQKOAokMjU5NzE4NzAtNzBiYy00"`
}

// GroupLeaveRequest represents a group leave request
type GroupLeaveRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
}

// GroupUpdateParticipantsRequest represents a group participants update request
type GroupUpdateParticipantsRequest struct {
	GroupJID     string   `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Participants []string `json:"participants" binding:"required" example:"+5511999999999,+5511888888888"`
	Action       string   `json:"action" binding:"required" example:"add"` // add, remove, promote, demote
}

// GroupSetNameRequest represents a group name update request
type GroupSetNameRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Name     string `json:"name" binding:"required" example:"New Group Name"`
}

// GroupSetTopicRequest represents a group topic update request
type GroupSetTopicRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Topic    string `json:"topic" binding:"required" example:"New group description"`
}

// GroupSetPhotoRequest represents a group photo update request
type GroupSetPhotoRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Image    string `json:"image" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."`
}

// GroupRemovePhotoRequest represents a group photo removal request
type GroupRemovePhotoRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
}

// GroupSetAnnounceRequest represents a group announce mode update request
type GroupSetAnnounceRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Announce bool   `json:"announce" binding:"required" example:"true"`
}

// GroupSetLockedRequest represents a group locked mode update request
type GroupSetLockedRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Locked   bool   `json:"locked" binding:"required" example:"true"`
}

// GroupSetEphemeralRequest represents a group ephemeral messages update request
type GroupSetEphemeralRequest struct {
	GroupJID string `json:"groupJid" binding:"required" example:"120363025246125486@g.us"`
	Duration int64  `json:"duration" binding:"required" example:"86400"` // seconds
}

// GroupInviteInfoRequest represents a group invite info request
type GroupInviteInfoRequest struct {
	InviteCode string `json:"inviteCode" binding:"required" example:"CjQKOAokMjU5NzE4NzAtNzBiYy00"`
}
