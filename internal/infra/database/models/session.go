package models

import (
	"strings"
	"time"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/shared/types"
)

// SessionModel represents the database model for sessions
type SessionModel struct {
	ID            string    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Status        string    `db:"status" json:"status"`
	WhatsAppJID   string    `db:"device_jid" json:"whatsapp_jid"`
	QRCode        string    `db:"qr_code" json:"qr_code"`
	ProxyURL      string    `db:"proxy_url" json:"proxy_url"`
	WebhookURL    string    `db:"webhook_url" json:"webhook_url"`
	WebhookEvents string    `db:"webhook_events" json:"webhook_events"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// ToDomain converts database model to domain entity
func (m *SessionModel) ToDomain() *session.Session {
	// Parse webhook events from comma-separated string
	var events []string
	if m.WebhookEvents != "" {
		events = strings.Split(m.WebhookEvents, ",")
		// Trim whitespace from each event
		for i, event := range events {
			events[i] = strings.TrimSpace(event)
		}
	}

	return &session.Session{
		ID:          m.ID,
		Name:        m.Name,
		Status:      types.Status(m.Status),
		WhatsAppJID: m.WhatsAppJID,
		QRCode:      m.QRCode,
		ProxyURL:    m.ProxyURL,
		WebhookURL:  m.WebhookURL,
		Events:      events,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromDomain converts domain entity to database model
func FromDomain(s *session.Session) *SessionModel {
	// Convert events slice to comma-separated string
	var eventsStr string
	if len(s.Events) > 0 {
		eventsStr = strings.Join(s.Events, ",")
	}

	return &SessionModel{
		ID:            s.ID,
		Name:          s.Name,
		Status:        string(s.Status),
		WhatsAppJID:   s.WhatsAppJID,
		QRCode:        s.QRCode,
		ProxyURL:      s.ProxyURL,
		WebhookURL:    s.WebhookURL,
		WebhookEvents: eventsStr,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

// FromDomain converts domain entity to database model
func (m *SessionModel) FromDomain(s *session.Session) {
	m.ID = s.ID
	m.Name = s.Name
	m.Status = string(s.Status)
	m.WhatsAppJID = s.WhatsAppJID
	m.QRCode = s.QRCode
	m.ProxyURL = s.ProxyURL
	m.CreatedAt = s.CreatedAt
	m.UpdatedAt = s.UpdatedAt
}
