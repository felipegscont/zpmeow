package request

// ContextInfo represents message context information
type ContextInfo struct {
	StanzaID      string `json:"stanzaId,omitempty" example:"3EB0C431C26A1916E07A"`
	Participant   string `json:"participant,omitempty" example:"+5511888888888@s.whatsapp.net"`
	QuotedMessage any    `json:"quotedMessage,omitempty"`
}

// SendTextRequest represents a text message request
type SendTextRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Body        string      `json:"body" binding:"required" example:"Hello, World!"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendImageRequest represents an image message request
type SendImageRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Image       string      `json:"image" binding:"required" example:"data:image/jpeg;base64,/9j/4AAQ..."`
	Caption     string      `json:"caption,omitempty" example:"Image caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	MimeType    string      `json:"mimeType,omitempty" example:"image/jpeg"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendAudioRequest represents an audio message request
type SendAudioRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Audio       string      `json:"audio" binding:"required" example:"data:audio/ogg;base64,T2dnU..."`
	Caption     string      `json:"caption,omitempty" example:"Audio caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	MimeType    string      `json:"mimeType,omitempty" example:"audio/ogg; codecs=opus"`
	PTT         bool        `json:"ptt,omitempty" example:"false"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendMediaRequest represents a generic media message request
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

// SendLocationRequest represents a location message request
type SendLocationRequest struct {
	Phone     string  `json:"phone" binding:"required" example:"+5511999999999"`
	Latitude  float64 `json:"latitude" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	Address   string  `json:"address,omitempty" example:"São Paulo, SP, Brazil"`
	ID        string  `json:"id,omitempty" example:"custom-message-id"`
}

// SendContactRequest represents a contact message request
type SendContactRequest struct {
	Phone   string  `json:"phone" binding:"required" example:"+5511999999999"`
	Contact Contact `json:"contact" binding:"required"`
	ID      string  `json:"id,omitempty" example:"custom-message-id"`
}

// SendDocumentRequest represents a document message request
type SendDocumentRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Document    string      `json:"document" binding:"required" example:"data:application/pdf;base64,JVBERi0x..."`
	Filename    string      `json:"filename,omitempty" example:"document.pdf"`
	Caption     string      `json:"caption,omitempty" example:"Document caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendVideoRequest represents a video message request
type SendVideoRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Video       string      `json:"video" binding:"required" example:"data:video/mp4;base64,AAAAIGZ0eXA..."`
	Caption     string      `json:"caption,omitempty" example:"Video caption"`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendStickerRequest represents a sticker message request
type SendStickerRequest struct {
	Phone       string      `json:"phone" binding:"required" example:"+5511999999999"`
	Sticker     string      `json:"sticker" binding:"required" example:"data:image/webp;base64,UklGRiQAAABXRUJQ..."`
	ID          string      `json:"id,omitempty" example:"custom-message-id"`
	ContextInfo ContextInfo `json:"contextInfo,omitempty"`
}

// SendPollRequest represents a poll message request
type SendPollRequest struct {
	Phone           string   `json:"phone" binding:"required" example:"+5511999999999"`
	Name            string   `json:"name" binding:"required" example:"What's your favorite color?"`
	Options         []string `json:"options" binding:"required" example:"Red,Blue,Green"`
	SelectableCount int      `json:"selectableCount,omitempty" example:"1"`
	ID              string   `json:"id,omitempty" example:"custom-message-id"`
}

// Contact represents a contact for contact messages
type Contact struct {
	DisplayName string `json:"displayName" binding:"required" example:"John Doe"`
	VCard       string `json:"vcard" binding:"required" example:"BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"`
}

// Button represents a button for interactive messages
type Button struct {
	ButtonID   string     `json:"buttonId" binding:"required" example:"btn_1"`
	ButtonText ButtonText `json:"buttonText" binding:"required"`
	Type       int        `json:"type" example:"1"`
}

// ButtonText represents button text
type ButtonText struct {
	DisplayText string `json:"displayText" binding:"required" example:"Option 1"`
}

// Section represents a section for list messages
type Section struct {
	Title string `json:"title" binding:"required" example:"Section 1"`
	Rows  []Row  `json:"rows" binding:"required"`
}

// Row represents a row in a list section
type Row struct {
	Title       string `json:"title" binding:"required" example:"Row 1"`
	Description string `json:"description,omitempty" example:"Row description"`
	RowID       string `json:"rowId" binding:"required" example:"row_1"`
}

// SendButtonRequest represents a button message request
type SendButtonRequest struct {
	Phone   string   `json:"phone" binding:"required" example:"+5511999999999"`
	Text    string   `json:"text" binding:"required" example:"Choose an option:"`
	Buttons []Button `json:"buttons" binding:"required"`
	Footer  string   `json:"footer,omitempty" example:"Footer text"`
	ID      string   `json:"id,omitempty" example:"custom-message-id"`
}

// SendListRequest represents a list message request
type SendListRequest struct {
	Phone       string    `json:"phone" binding:"required" example:"+5511999999999"`
	Text        string    `json:"text" binding:"required" example:"Choose from the list:"`
	ButtonText  string    `json:"buttonText" binding:"required" example:"View Options"`
	Sections    []Section `json:"sections" binding:"required"`
	Footer      string    `json:"footer,omitempty" example:"Footer text"`
	Title       string    `json:"title,omitempty" example:"List Title"`
	ID          string    `json:"id,omitempty" example:"custom-message-id"`
}

// Group-related requests
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
