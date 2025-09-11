package session

import "time"




type BaseSessionInfo struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"My WhatsApp Session"`
	Status    string    `json:"status" example:"disconnected"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}


type ExtendedSessionInfo struct {
	BaseSessionInfo
	WhatsAppJID string `json:"whatsapp_jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	QRCode      string `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	ProxyURL    string `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
}




type CreateSessionRequest struct {
	Name string `json:"name" binding:"required" example:"My WhatsApp Session"`
}




type CreateSessionResponse = BaseSessionInfo


type SessionInfoResponse = ExtendedSessionInfo


type SessionListResponse struct {
	Sessions []SessionInfoResponse `json:"sessions"`
	Total    int                   `json:"total" example:"5"`
}




type PairSessionRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"5511999999999"`
}


type PairSessionResponse struct {
	PairingCode string `json:"pairing_code" example:"ABCD-EFGH"`
}


type QRCodeResponse struct {
	QRCode string `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Status string `json:"status" example:"waiting_for_scan"`
}


type ProxyRequest struct {
	ProxyURL string `json:"proxy_url" binding:"required" example:"http://proxy.example.com:8080"`
}


type ProxyResponse struct {
	ProxyURL string `json:"proxy_url" example:"http://proxy.example.com:8080"`
	Message  string `json:"message" example:"Proxy set successfully"`
}




type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}


type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}


type MessageResponse = SuccessResponse
type PingResponse = SuccessResponse
