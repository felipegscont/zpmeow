package message

import (
	"time"
)

type Message struct {
	ID          string
	SessionID   string
	ChatJID     string
	FromJID     string
	ToJID       string
	MessageType MessageType
	Content     string
	MediaURL    string
	Timestamp   time.Time
	Status      MessageStatus
}

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVideo    MessageType = "video"
	MessageTypeDocument MessageType = "document"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
)

type MessageStatus string

const (
	MessageStatusPending MessageStatus = "pending"
	MessageStatusSent    MessageStatus = "sent"
	MessageStatusFailed  MessageStatus = "failed"
)

func NewMessage(sessionID, chatJID, fromJID, toJID string, messageType MessageType, content string) *Message {
	return &Message{
		SessionID:   sessionID,
		ChatJID:     chatJID,
		FromJID:     fromJID,
		ToJID:       toJID,
		MessageType: messageType,
		Content:     content,
		Timestamp:   time.Now(),
		Status:      MessageStatusPending,
	}
}

func (m *Message) IsValid() bool {
	return m.SessionID != "" && m.ChatJID != "" && m.MessageType != ""
}

func (m *Message) SetMediaURL(url string) {
	m.MediaURL = url
}

func (m *Message) MarkAsSent() {
	m.Status = MessageStatusSent
}

func (m *Message) MarkAsFailed() {
	m.Status = MessageStatusFailed
}
