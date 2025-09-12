package whatsapp

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
	Phone       string `json:"phone" binding:"required" example:"+5511999999999"`
	DisplayName string `json:"displayName" binding:"required" example:"John Doe"`
	Vcard       string `json:"vcard" binding:"required" example:"BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"`
	ID          string `json:"id,omitempty" example:"custom-message-id"`
}

// SendPollRequest represents a poll message request
type SendPollRequest struct {
	Phone           string   `json:"phone" binding:"required" example:"+5511999999999"`
	Name            string   `json:"name" binding:"required" example:"What's your favorite color?"`
	Options         []string `json:"options" binding:"required" example:"Red,Blue,Green"`
	SelectableCount int      `json:"selectableCount,omitempty" example:"1"`
	ID              string   `json:"id,omitempty" example:"custom-message-id"`
}
