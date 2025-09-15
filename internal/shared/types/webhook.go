package types

import (
	"time"
)

type EventType string

const (
	EventTypeMessage         EventType = "message"
	EventTypeMessageAck      EventType = "message.ack"
	EventTypeMessageRevoked  EventType = "message.revoked"
	
	EventTypeSessionConnected    EventType = "session.connected"
	EventTypeSessionDisconnected EventType = "session.disconnected"
	EventTypeSessionQR           EventType = "session.qr"
	
	EventTypeContactUpdate   EventType = "contact.update"
	EventTypeContactPresence EventType = "contact.presence"
	
	EventTypeGroupJoin       EventType = "group.join"
	EventTypeGroupLeave      EventType = "group.leave"
	EventTypeGroupUpdate     EventType = "group.update"
	
	EventTypeCallOffer       EventType = "call.offer"
	EventTypeCallAccept      EventType = "call.accept"
	EventTypeCallTerminate   EventType = "call.terminate"
)

type Webhook struct {
	ID          string
	SessionID   string
	URL         string
	Events      []EventType
	Secret      string
	Active      bool
	RetryCount  int
	MaxRetries  int
	Timeout     time.Duration
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WebhookEvent struct {
	ID        string
	WebhookID string
	SessionID string
	Type      EventType
	Data      interface{}
	Timestamp time.Time
	Attempts  int
	LastError string
	CreatedAt time.Time
}

type WebhookPayload struct {
	Event     string      `json:"event"`
	SessionID string      `json:"session_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}
