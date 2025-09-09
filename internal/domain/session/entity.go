package session

import "time"

// Session represents a WhatsApp session in the system.
// We are using struct tags to map database columns to the struct fields.
// The `db` tag is for sqlx, and the `json` tag is for gin to serialize the struct to JSON.
type Session struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	WhatsAppJID string     `db:"whatsapp_jid" json:"whatsapp_jid"`
	Status      string     `db:"status" json:"status"`
	QRCode      string     `db:"qr_code" json:"qr_code,omitempty"`
	ProxyURL    string     `db:"proxy_url" json:"proxy_url,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
