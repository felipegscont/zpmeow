package response

import "time"

// BaseSessionInfo represents basic session information
type BaseSessionInfo struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"My WhatsApp Session"`
	Status    string    `json:"status" example:"disconnected"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ExtendedSessionInfo represents detailed session information
type ExtendedSessionInfo struct {
	BaseSessionInfo
	WhatsAppJID string `json:"whatsapp_jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	QRCode      string `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	ProxyURL    string `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
}

// Response type aliases
type CreateSessionResponse = BaseSessionInfo
type SessionInfoResponse = ExtendedSessionInfo

// SessionListResponse represents a list of sessions
type SessionListResponse struct {
	Sessions []SessionInfoResponse `json:"sessions"`
	Total    int                   `json:"total" example:"5"`
}

// PairSessionResponse represents the response from pairing a session
type PairSessionResponse struct {
	PairingCode string `json:"pairing_code" example:"ABCD-EFGH"`
}

// QRCodeResponse represents a QR code response
type QRCodeResponse struct {
	QRCode string `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Status string `json:"status" example:"waiting_for_scan"`
}

// ProxyResponse represents a proxy configuration response
type ProxyResponse struct {
	ProxyURL string `json:"proxy_url" example:"http://proxy.example.com:8080"`
	Message  string `json:"message" example:"Proxy set successfully"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

// SessionStatusResponse represents session status information
type SessionStatusResponse struct {
	ID        string   `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string   `json:"name" example:"my-session"`
	Connected bool     `json:"connected" example:"true"`
	LoggedIn  bool     `json:"logged_in" example:"true"`
	Status    string   `json:"status" example:"connected"`
	JID       string   `json:"jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	Webhook   string   `json:"webhook,omitempty" example:"https://api.example.com/webhook"`
	Events    []string `json:"events,omitempty" example:"message.received,session.connected"`
	ProxyURL  string   `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
}

// Response type aliases
type MessageResponse = SuccessResponse
type PingResponse = SuccessResponse
