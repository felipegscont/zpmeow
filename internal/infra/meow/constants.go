package meow

import "time"


const (

	AppName    = "zpmeow"
	AppVersion = "1.0.0"
	

	DefaultUserAgent = "zpmeow/1.0.0"
	DefaultPlatform  = "Chrome (Linux)"
)


const (
	DefaultTimeout        = 30 * time.Second
	QRTimeout            = 2 * time.Minute
	MessageTimeout       = 10 * time.Second
	HistorySyncTimeout   = 30 * time.Second
	ClientCleanupTimeout = 30 * time.Minute
	ConnectionTimeout    = 60 * time.Second
	ShutdownTimeout      = 10 * time.Second
)


const (
	MaxRetries         = 3
	RetryDelay         = 1 * time.Second
	MaxRetryDelay      = 10 * time.Second
	RetryBackoffFactor = 2.0
)


const (
	KillChannelBuffer     = 1
	QRStopChannelBuffer   = 1
	EventChannelBuffer    = 100
	MessageChannelBuffer  = 50
)


const (
	MaxDBConnections     = 25
	MaxIdleDBConnections = 5
	DBConnectionLifetime = 5 * time.Minute
)


const (
	MaxImageSize    = 16 * 1024 * 1024  // 16MB
	MaxVideoSize    = 64 * 1024 * 1024  // 64MB
	MaxAudioSize    = 16 * 1024 * 1024  // 16MB
	MaxDocumentSize = 100 * 1024 * 1024 // 100MB
	
	DefaultImageMimeType    = "image/jpeg"
	DefaultAudioMimeType    = "audio/ogg; codecs=opus"
	DefaultVideoMimeType    = "video/mp4"
	DefaultDocumentMimeType = "application/octet-stream"
)


const (
	ErrClientNotConnected    = "client is not connected"
	ErrClientNotFound        = "client not found for session"
	ErrClientAlreadyAuth     = "client is already authenticated"
	ErrSessionNotConnected   = "session is not connected"
	ErrSessionNotFound       = "session not found"
	ErrQRNotAvailable       = "QR code not yet available, please wait"
	ErrInvalidJID           = "invalid JID format"
	ErrEmptySessionID       = "session ID cannot be empty"
	ErrEmptyDeviceJID       = "device JID cannot be empty"
	ErrInvalidPhoneNumber   = "invalid phone number format"
	ErrMessageTooLarge      = "message content too large"
	ErrUnsupportedMediaType = "unsupported media type"
	ErrConnectionFailed     = "connection failed"
	ErrAuthenticationFailed = "authentication failed"
	ErrTimeout              = "operation timed out"
	ErrInvalidInput         = "invalid input provided"
	ErrOperationCancelled   = "operation was cancelled"
)


const (
	MsgClientStarted     = "client started successfully"
	MsgClientStopped     = "client stopped successfully"
	MsgClientConnected   = "client connected successfully"
	MsgClientDisconnected = "client disconnected successfully"
	MsgMessageSent       = "message sent successfully"
	MsgQRGenerated       = "QR code generated successfully"
	MsgPairingSuccessful = "pairing completed successfully"
)


const (
	LogClientCreated    = "created new client for session"
	LogClientRemoved    = "removed client for session"
	LogClientConnecting = "connecting client for session"
	LogClientConnected  = "client connected for session"
	LogQRGenerated      = "QR code generated for session"
	LogPairSuccess      = "pairing successful for session"
	LogMessageReceived  = "message received for session"
	LogMessageSent      = "message sent for session"
)


const (
	EventTypeMessage      = "message"
	EventTypeReceipt      = "receipt"
	EventTypePresence     = "presence"
	EventTypeChatPresence = "chat_presence"
	EventTypeConnected    = "connected"
	EventTypeDisconnected = "disconnected"
	EventTypeQR           = "qr"
	EventTypePairSuccess  = "pair_success"
	EventTypeError        = "error"
)


const (
	StatusConnecting    = "connecting"
	StatusConnected     = "connected"
	StatusDisconnected  = "disconnected"
	StatusError         = "error"
	StatusAuthenticating = "authenticating"
	StatusReconnecting  = "reconnecting"
)


const (
	WhatsAppUserServer  = "s.whatsapp.net"
	WhatsAppGroupServer = "g.us"
	WhatsAppBroadcastServer = "broadcast"
	
	QREventCode    = "code"
	QREventTimeout = "timeout"
	QREventSuccess = "success"
)


const (
	ExtensionJPEG = ".jpg"
	ExtensionPNG  = ".png"
	ExtensionGIF  = ".gif"
	ExtensionMP4  = ".mp4"
	ExtensionOGG  = ".ogg"
	ExtensionPDF  = ".pdf"
	ExtensionTXT  = ".txt"
)


const (
	PhoneNumberPattern = `^\+?[1-9]\d{1,14}$`
	SessionIDPattern   = `^[a-zA-Z0-9_-]+$`
	JIDPattern         = `^[0-9]+@[a-z.]+$`
)


const (
	MaxConcurrentConnections = 100
	MaxMessageQueueSize      = 1000
	MaxEventHandlers         = 50
	CleanupInterval          = 5 * time.Minute
	HealthCheckInterval      = 30 * time.Second
)


const (
	EnableDebugLogging     = false
	EnableMetrics          = true
	EnableHealthChecks     = true
	EnableAutoReconnect    = true
	EnableMessageQueue     = true
	EnableEventBatching    = false
)
