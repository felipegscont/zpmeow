package session

import "time"

// CreateSessionRequest represents the request body for creating a new session
type CreateSessionRequest struct {
	Name string `json:"name" binding:"required" example:"My WhatsApp Session"`
}

// CreateSessionResponse represents the response for creating a new session
type CreateSessionResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"My WhatsApp Session"`
	Status    string    `json:"status" example:"disconnected"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SessionInfoResponse represents the response for getting session information
type SessionInfoResponse struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"My WhatsApp Session"`
	WhatsAppJID string    `json:"whatsapp_jid" example:"5511999999999@s.whatsapp.net"`
	Status      string    `json:"status" example:"connected"`
	QRCode      string    `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	ProxyURL    string    `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SessionListResponse represents the response for listing sessions
type SessionListResponse struct {
	Sessions []SessionInfoResponse `json:"sessions"`
	Total    int                   `json:"total" example:"5"`
}

// PairSessionRequest represents the request body for pairing a session with a phone number
type PairSessionRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"5511999999999"`
}

// PairSessionResponse represents the response for pairing a session
type PairSessionResponse struct {
	PairingCode string `json:"pairing_code" example:"ABCD-EFGH"`
}

// QRCodeResponse represents the response for getting QR code
type QRCodeResponse struct {
	QRCode string `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Status string `json:"status" example:"waiting_for_scan"`
}

// ProxyRequest represents the request body for setting proxy
type ProxyRequest struct {
	ProxyURL string `json:"proxy_url" binding:"required" example:"http://proxy.example.com:8080"`
}

// ProxyResponse represents the response for proxy operations
type ProxyResponse struct {
	ProxyURL string `json:"proxy_url" example:"http://proxy.example.com:8080"`
	Message  string `json:"message" example:"Proxy set successfully"`
}

// MessageResponse represents a generic message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

// PingResponse represents the response for ping endpoint
type PingResponse struct {
	Message string `json:"message" example:"pong"`
}
