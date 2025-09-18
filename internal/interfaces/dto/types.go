package dto

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// SessionID represents a session identifier with validation
type SessionID struct {
	value string
}

// NewSessionID creates a new SessionID with validation
func NewSessionID(value string) (SessionID, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return SessionID{}, fmt.Errorf("session ID cannot be empty")
	}

	if len(cleaned) < 1 || len(cleaned) > 100 {
		return SessionID{}, fmt.Errorf("session ID must be between 1 and 100 characters")
	}

	return SessionID{value: cleaned}, nil
}

// Value returns the string value of the SessionID
func (s SessionID) Value() string {
	return s.value
}

// String returns the string representation of the SessionID
func (s SessionID) String() string {
	return s.value
}

// IsEmpty checks if the SessionID is empty
func (s SessionID) IsEmpty() bool {
	return s.value == ""
}

// Equals checks if two SessionIDs are equal
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}

// SessionInfo represents session information for API responses
type SessionInfo struct {
	ID        string    `json:"id" example:"session_123"`
	Name      string    `json:"name" example:"My Session"`
	Status    string    `json:"status" example:"connected"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ID represents a generic identifier with validation
type ID string

// NewID creates a new ID with validation
func NewID(value string) (ID, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("ID cannot be empty")
	}
	if len(value) > 255 {
		return "", fmt.Errorf("ID too long (max 255 characters)")
	}
	return ID(value), nil
}

// String returns the string representation of the ID
func (id ID) String() string {
	return string(id)
}

// IsEmpty checks if the ID is empty
func (id ID) IsEmpty() bool {
	return strings.TrimSpace(string(id)) == ""
}

// IsValid checks if the ID is valid
func (id ID) IsValid() bool {
	return !id.IsEmpty() && len(string(id)) <= 255
}

// MarshalJSON implements json.Marshaler
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler
func (id *ID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	newID, err := NewID(s)
	if err != nil {
		return err
	}
	*id = newID
	return nil
}

// Timestamp represents a timestamp with utility methods
type Timestamp time.Time

// Now creates a new Timestamp with current time
func Now() Timestamp {
	return Timestamp(time.Now())
}

// NewTimestamp creates a Timestamp from time.Time
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp(t)
}

// Time returns the underlying time.Time
func (ts Timestamp) Time() time.Time {
	return time.Time(ts)
}

// Unix returns the Unix timestamp
func (ts Timestamp) Unix() int64 {
	return time.Time(ts).Unix()
}

// UnixMilli returns the Unix timestamp in milliseconds
func (ts Timestamp) UnixMilli() int64 {
	return time.Time(ts).UnixMilli()
}

// String returns the string representation
func (ts Timestamp) String() string {
	return time.Time(ts).Format(time.RFC3339)
}

// IsZero checks if the timestamp is zero
func (ts Timestamp) IsZero() bool {
	return time.Time(ts).IsZero()
}

// MarshalJSON implements json.Marshaler
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ts))
}

// UnmarshalJSON implements json.Unmarshaler
func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*ts = Timestamp(t)
	return nil
}

// ConnectionStatus represents connection status
type ConnectionStatus string

const (
	ConnectionStatusDisconnected ConnectionStatus = "disconnected"
	ConnectionStatusConnecting   ConnectionStatus = "connecting"
	ConnectionStatusConnected    ConnectionStatus = "connected"
	ConnectionStatusError        ConnectionStatus = "error"
	ConnectionStatusDeleted      ConnectionStatus = "deleted"
)

// String returns the string representation
func (s ConnectionStatus) String() string {
	return string(s)
}

// IsValid checks if the status is valid
func (s ConnectionStatus) IsValid() bool {
	switch s {
	case ConnectionStatusDisconnected, ConnectionStatusConnecting,
		ConnectionStatusConnected, ConnectionStatusError, ConnectionStatusDeleted:
		return true
	default:
		return false
	}
}

// IsConnected checks if status indicates connection
func (s ConnectionStatus) IsConnected() bool {
	return s == ConnectionStatusConnected
}

// IsError checks if status indicates error
func (s ConnectionStatus) IsError() bool {
	return s == ConnectionStatusError
}

// Phone represents a phone number with validation
type Phone string

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

// NewPhone creates a new Phone with validation
func NewPhone(value string) (Phone, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}

	if !phoneRegex.MatchString(cleaned) {
		return "", fmt.Errorf("invalid phone number format: %s", value)
	}

	return Phone(cleaned), nil
}

// String returns the string representation
func (p Phone) String() string {
	return string(p)
}

// IsValid checks if the phone number is valid
func (p Phone) IsValid() bool {
	return phoneRegex.MatchString(string(p))
}

// ToJID converts phone to WhatsApp JID format
func (p Phone) ToJID() string {
	phone := string(p)
	if strings.Contains(phone, "@") {
		return phone
	}
	return phone + "@s.whatsapp.net"
}

// JID represents a WhatsApp JID (Jabber ID)
type JID string

var jidRegex = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)

// NewJID creates a new JID with validation
func NewJID(value string) (JID, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return "", fmt.Errorf("JID cannot be empty")
	}

	if !jidRegex.MatchString(cleaned) {
		return "", fmt.Errorf("invalid JID format: %s", value)
	}

	return JID(cleaned), nil
}

// String returns the string representation
func (j JID) String() string {
	return string(j)
}

// IsValid checks if the JID is valid
func (j JID) IsValid() bool {
	return jidRegex.MatchString(string(j))
}

// IsGroup checks if JID is a group JID
func (j JID) IsGroup() bool {
	return strings.Contains(string(j), "@g.us")
}

// IsUser checks if JID is a user JID
func (j JID) IsUser() bool {
	return strings.Contains(string(j), "@s.whatsapp.net")
}

// ToPhone extracts phone number from user JID
func (j JID) ToPhone() (Phone, error) {
	if !j.IsUser() {
		return "", fmt.Errorf("JID is not a user JID")
	}

	parts := strings.Split(string(j), "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid JID format")
	}

	return NewPhone(parts[0])
}

// MessageID represents a WhatsApp message ID
type MessageID string

var messageIDRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// NewMessageID creates a new MessageID with validation
func NewMessageID(value string) (MessageID, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return "", fmt.Errorf("message ID cannot be empty")
	}

	if len(cleaned) < 10 || len(cleaned) > 100 {
		return "", fmt.Errorf("message ID length must be between 10 and 100 characters")
	}

	if !messageIDRegex.MatchString(cleaned) {
		return "", fmt.Errorf("message ID contains invalid characters")
	}

	return MessageID(cleaned), nil
}

// String returns the string representation
func (m MessageID) String() string {
	return string(m)
}

// IsValid checks if the message ID is valid
func (m MessageID) IsValid() bool {
	return len(string(m)) >= 10 && len(string(m)) <= 100 &&
		messageIDRegex.MatchString(string(m))
}

// MessageType represents different types of messages
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVideo    MessageType = "video"
	MessageTypeDocument MessageType = "document"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeContact  MessageType = "contact"
	MessageTypeLocation MessageType = "location"
	MessageTypeButton   MessageType = "button"
	MessageTypeList     MessageType = "list"
	MessageTypePoll     MessageType = "poll"
)

// String returns the string representation
func (mt MessageType) String() string {
	return string(mt)
}

// IsValid checks if the message type is valid
func (mt MessageType) IsValid() bool {
	switch mt {
	case MessageTypeText, MessageTypeImage, MessageTypeAudio, MessageTypeVideo,
		MessageTypeDocument, MessageTypeSticker, MessageTypeContact,
		MessageTypeLocation, MessageTypeButton, MessageTypeList, MessageTypePoll:
		return true
	default:
		return false
	}
}

// IsMedia checks if message type is media
func (mt MessageType) IsMedia() bool {
	switch mt {
	case MessageTypeImage, MessageTypeAudio, MessageTypeVideo,
		MessageTypeDocument, MessageTypeSticker:
		return true
	default:
		return false
	}
}

// MessageStatus represents message delivery status
type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// String returns the string representation
func (ms MessageStatus) String() string {
	return string(ms)
}

// IsValid checks if the message status is valid
func (ms MessageStatus) IsValid() bool {
	switch ms {
	case MessageStatusPending, MessageStatusSent, MessageStatusDelivered,
		MessageStatusRead, MessageStatusFailed:
		return true
	default:
		return false
	}
}

// IsFinal checks if status is final (delivered or read)
func (ms MessageStatus) IsFinal() bool {
	return ms == MessageStatusDelivered || ms == MessageStatusRead
}

// IsError checks if status indicates error
func (ms MessageStatus) IsError() bool {
	return ms == MessageStatusFailed
}

// SessionStatus represents session status
type SessionStatus string

const (
	SessionStatusActive   SessionStatus = "active"
	SessionStatusInactive SessionStatus = "inactive"
	SessionStatusPairing  SessionStatus = "pairing"
	SessionStatusError    SessionStatus = "error"
)

// String returns the string representation
func (ss SessionStatus) String() string {
	return string(ss)
}

// IsValid checks if the session status is valid
func (ss SessionStatus) IsValid() bool {
	switch ss {
	case SessionStatusActive, SessionStatusInactive, SessionStatusPairing, SessionStatusError:
		return true
	default:
		return false
	}
}

// IsActive checks if session is active
func (ss SessionStatus) IsActive() bool {
	return ss == SessionStatusActive
}

// WebhookStatus represents webhook status
type WebhookStatus string

const (
	WebhookStatusActive   WebhookStatus = "active"
	WebhookStatusInactive WebhookStatus = "inactive"
	WebhookStatusError    WebhookStatus = "error"
)

// String returns the string representation
func (ws WebhookStatus) String() string {
	return string(ws)
}

// IsValid checks if the webhook status is valid
func (ws WebhookStatus) IsValid() bool {
	switch ws {
	case WebhookStatusActive, WebhookStatusInactive, WebhookStatusError:
		return true
	default:
		return false
	}
}

// IsActive checks if webhook is active
func (ws WebhookStatus) IsActive() bool {
	return ws == WebhookStatusActive
}

// URL represents a URL with validation
type URL string

var urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

// NewURL creates a new URL with validation
func NewURL(value string) (URL, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	if len(cleaned) > 2048 {
		return "", fmt.Errorf("URL too long (max 2048 characters)")
	}

	if !urlRegex.MatchString(cleaned) {
		return "", fmt.Errorf("invalid URL format: %s", value)
	}

	return URL(cleaned), nil
}

// String returns the string representation
func (u URL) String() string {
	return string(u)
}

// IsValid checks if the URL is valid
func (u URL) IsValid() bool {
	return len(string(u)) <= 2048 && urlRegex.MatchString(string(u))
}

// IsHTTPS checks if URL uses HTTPS
func (u URL) IsHTTPS() bool {
	return strings.HasPrefix(string(u), "https://")
}

// MimeType represents a MIME type
type MimeType string

const (
	MimeTypeTextPlain      MimeType = "text/plain"
	MimeTypeImageJPEG      MimeType = "image/jpeg"
	MimeTypeImagePNG       MimeType = "image/png"
	MimeTypeImageWebP      MimeType = "image/webp"
	MimeTypeAudioMPEG      MimeType = "audio/mpeg"
	MimeTypeAudioOGG       MimeType = "audio/ogg"
	MimeTypeVideoMP4       MimeType = "video/mp4"
	MimeTypeApplicationPDF MimeType = "application/pdf"
	MimeTypeApplicationZip MimeType = "application/zip"
	MimeTypeOctetStream    MimeType = "application/octet-stream"
)

// String returns the string representation
func (mt MimeType) String() string {
	return string(mt)
}

// IsValid checks if the MIME type is valid
func (mt MimeType) IsValid() bool {
	return strings.Contains(string(mt), "/")
}

// IsImage checks if MIME type is image
func (mt MimeType) IsImage() bool {
	return strings.HasPrefix(string(mt), "image/")
}

// IsAudio checks if MIME type is audio
func (mt MimeType) IsAudio() bool {
	return strings.HasPrefix(string(mt), "audio/")
}

// IsVideo checks if MIME type is video
func (mt MimeType) IsVideo() bool {
	return strings.HasPrefix(string(mt), "video/")
}

// IsText checks if MIME type is text
func (mt MimeType) IsText() bool {
	return strings.HasPrefix(string(mt), "text/")
}

// Utility functions for common validations

// ValidateRequired checks if a string field is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateLength checks if a string is within length limits
func ValidateLength(value, fieldName string, min, max int) error {
	length := len(value)
	if length < min {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, min)
	}
	if length > max {
		return fmt.Errorf("%s must be at most %d characters long", fieldName, max)
	}
	return nil
}

// ValidateArrayLength checks if an array is within length limits
func ValidateArrayLength[T any](arr []T, fieldName string, min, max int) error {
	length := len(arr)
	if length < min {
		return fmt.Errorf("%s must have at least %d items", fieldName, min)
	}
	if length > max {
		return fmt.Errorf("%s cannot have more than %d items", fieldName, max)
	}
	return nil
}

// GroupRole represents roles in a WhatsApp group
type GroupRole string

const (
	GroupRoleParticipant GroupRole = "participant"
	GroupRoleAdmin       GroupRole = "admin"
	GroupRoleOwner       GroupRole = "owner"
)

// String returns the string representation
func (gr GroupRole) String() string {
	return string(gr)
}

// IsValid checks if the group role is valid
func (gr GroupRole) IsValid() bool {
	switch gr {
	case GroupRoleParticipant, GroupRoleAdmin, GroupRoleOwner:
		return true
	default:
		return false
	}
}

// HasAdminRights checks if role has admin rights
func (gr GroupRole) HasAdminRights() bool {
	return gr == GroupRoleAdmin || gr == GroupRoleOwner
}

// PresenceStatus represents user presence status
type PresenceStatus string

const (
	PresenceStatusAvailable   PresenceStatus = "available"
	PresenceStatusUnavailable PresenceStatus = "unavailable"
	PresenceStatusComposing   PresenceStatus = "composing"
	PresenceStatusRecording   PresenceStatus = "recording"
	PresenceStatusPaused      PresenceStatus = "paused"
)

// String returns the string representation
func (ps PresenceStatus) String() string {
	return string(ps)
}

// IsValid checks if the presence status is valid
func (ps PresenceStatus) IsValid() bool {
	switch ps {
	case PresenceStatusAvailable, PresenceStatusUnavailable,
		PresenceStatusComposing, PresenceStatusRecording, PresenceStatusPaused:
		return true
	default:
		return false
	}
}

// IsActive checks if presence indicates activity
func (ps PresenceStatus) IsActive() bool {
	return ps == PresenceStatusComposing || ps == PresenceStatusRecording
}

// EventType represents webhook event types
type EventType string

const (
	EventTypeMessage    EventType = "message"
	EventTypeReceipt    EventType = "receipt"
	EventTypePresence   EventType = "presence"
	EventTypeConnection EventType = "connection"
	EventTypeGroup      EventType = "group"
	EventTypeContact    EventType = "contact"
	EventTypeCall       EventType = "call"
	EventTypeNewsletter EventType = "newsletter"
)

// String returns the string representation
func (et EventType) String() string {
	return string(et)
}

// IsValid checks if the event type is valid
func (et EventType) IsValid() bool {
	switch et {
	case EventTypeMessage, EventTypeReceipt, EventTypePresence,
		EventTypeConnection, EventTypeGroup, EventTypeContact,
		EventTypeCall, EventTypeNewsletter:
		return true
	default:
		return false
	}
}

// PrivacySetting represents privacy setting values
type PrivacySetting string

const (
	PrivacySettingEveryone PrivacySetting = "everyone"
	PrivacySettingContacts PrivacySetting = "contacts"
	PrivacySettingNobody   PrivacySetting = "nobody"
)

// String returns the string representation
func (ps PrivacySetting) String() string {
	return string(ps)
}

// IsValid checks if the privacy setting is valid
func (ps PrivacySetting) IsValid() bool {
	switch ps {
	case PrivacySettingEveryone, PrivacySettingContacts, PrivacySettingNobody:
		return true
	default:
		return false
	}
}

// MediaQuality represents media quality settings
type MediaQuality string

const (
	MediaQualityAuto   MediaQuality = "auto"
	MediaQualityLow    MediaQuality = "low"
	MediaQualityMedium MediaQuality = "medium"
	MediaQualityHigh   MediaQuality = "high"
)

// String returns the string representation
func (mq MediaQuality) String() string {
	return string(mq)
}

// IsValid checks if the media quality is valid
func (mq MediaQuality) IsValid() bool {
	switch mq {
	case MediaQualityAuto, MediaQualityLow, MediaQualityMedium, MediaQualityHigh:
		return true
	default:
		return false
	}
}

// Note: HTTP status codes and common error codes are now defined in common.go
// Additional specific error codes for this module
const (
	ErrorCodeNotImplemented = "NOT_IMPLEMENTED"
	ErrorCodeServiceError   = "SERVICE_ERROR"
	ErrorCodeInvalidPhone   = "INVALID_PHONE"
	ErrorCodeInvalidSession = "INVALID_SESSION"
	ErrorCodeInvalidMessage = "INVALID_MESSAGE"
	ErrorCodeMessageFailed  = "MESSAGE_FAILED"
)
