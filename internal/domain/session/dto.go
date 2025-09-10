package session

import "time"

// Base structures to avoid duplication

// BaseSessionInfo contains common session fields
type BaseSessionInfo struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"My WhatsApp Session"`
	Status    string    `json:"status" example:"disconnected"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ExtendedSessionInfo includes additional optional fields
type ExtendedSessionInfo struct {
	BaseSessionInfo
	WhatsAppJID string `json:"whatsapp_jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	QRCode      string `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	ProxyURL    string `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
}

// Request DTOs

// CreateSessionRequest represents the request body for creating a new session
type CreateSessionRequest struct {
	Name string `json:"name" binding:"required" example:"My WhatsApp Session"`
}

// Response DTOs

// CreateSessionResponse represents the response for creating a new session
type CreateSessionResponse = BaseSessionInfo

// SessionInfoResponse represents the response for getting session information
type SessionInfoResponse = ExtendedSessionInfo

// SessionListResponse represents the response for listing sessions
type SessionListResponse struct {
	Sessions []SessionInfoResponse `json:"sessions"`
	Total    int                   `json:"total" example:"5"`
}

// Operation-specific DTOs

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

// Common response DTOs

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

// Convenience type aliases for backward compatibility
type MessageResponse = SuccessResponse
type PingResponse = SuccessResponse
