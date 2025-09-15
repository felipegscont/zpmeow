package dto

import (
	"time"
)

// ============================================================================
// SESSION REQUEST DTOs
// ============================================================================

// CreateSessionRequest represents the request to create a new session
type CreateSessionRequest struct {
	Name       string `json:"name" binding:"required" example:"default"`
	ProxyURL   string `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
	WebhookURL string `json:"webhook_url,omitempty" example:"https://webhook.example.com/whatsapp"`
	Events     string `json:"events,omitempty" example:"message,status"`
}

// PairPhoneRequest represents the request to pair a phone number
type PairPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"5511999999999"`
}

// ============================================================================
// SESSION DATA STRUCTURES
// ============================================================================

// SessionInfo represents session information
type SessionInfo struct {
	ID         string    `json:"id" example:"default"`
	Name       string    `json:"name" example:"default"`
	Status     string    `json:"status" example:"connected"`
	WaJID      string    `json:"wa_jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	ProxyURL   string    `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
	WebhookURL string    `json:"webhook_url,omitempty" example:"https://webhook.example.com/whatsapp"`
	Events     []string  `json:"events,omitempty" example:"message,status"`
	ApiKey     string    `json:"api_key,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SessionConnectionInfo represents session connection information
type SessionConnectionInfo struct {
	QRCode      string    `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Connected   bool      `json:"connected" example:"true"`
	IsConnected bool      `json:"is_connected" example:"true"`
	LastSeen    time.Time `json:"last_seen,omitempty" example:"2023-01-01T00:00:00Z"`
	PairCode    string    `json:"pair_code,omitempty" example:"ABCD-1234"`
}

// ============================================================================
// SESSION RESPONSE DTOs
// ============================================================================

// SessionResponse represents the standardized response format for session operations
type SessionResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Data    SessionData           `json:"data"`
	Error   *SessionErrorResponse `json:"error,omitempty"`
}

// SessionData contains the response data for session operations
type SessionData struct {
	SessionID  string                 `json:"session_id,omitempty" example:"default"`
	Action     string                 `json:"action" example:"create"`
	Status     string                 `json:"status" example:"success"`
	Timestamp  time.Time              `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session    *SessionInfo           `json:"session,omitempty"`
	Connection *SessionConnectionInfo `json:"connection,omitempty"`
	Sessions   []SessionInfo          `json:"sessions,omitempty"`
	QRCode     string                 `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	PairCode   string                 `json:"pair_code,omitempty" example:"ABCD-1234"`
}

// SessionErrorResponse represents error information for session operations
type SessionErrorResponse struct {
	Code    string `json:"code" example:"INVALID_SESSION_ID"`
	Message string `json:"message" example:"Invalid session ID format"`
	Details string `json:"details,omitempty" example:"Session ID must be alphanumeric"`
}

// ============================================================================
// SPECIFIC RESPONSE DTOs FOR SWAGGER DOCUMENTATION
// ============================================================================

// CreateSessionResponse represents the response for session creation
type CreateSessionResponse struct {
	Success bool                      `json:"success" example:"true"`
	Code    int                       `json:"code" example:"201"`
	Data    CreateSessionResponseData `json:"data"`
	Error   *SessionErrorResponse     `json:"error,omitempty"`
}

// CreateSessionResponseData contains the data for session creation response
type CreateSessionResponseData struct {
	SessionID string       `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Action    string       `json:"action" example:"create"`
	Status    string       `json:"status" example:"success"`
	Timestamp time.Time    `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session   *SessionInfo `json:"session"`
}

// SessionInfoResponse represents the response for getting session information
type SessionInfoResponse struct {
	Success bool                        `json:"success" example:"true"`
	Code    int                         `json:"code" example:"200"`
	Data    SessionInfoResponseData     `json:"data"`
	Error   *SessionErrorResponse       `json:"error,omitempty"`
}

// SessionInfoResponseData contains the data for session info response
type SessionInfoResponseData struct {
	SessionID string       `json:"session_id" example:"default"`
	Action    string       `json:"action" example:"get"`
	Status    string       `json:"status" example:"success"`
	Timestamp time.Time    `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session   *SessionInfo `json:"session"`
}

// SessionListResponse represents the response for listing sessions
type SessionListResponse struct {
	Success bool                       `json:"success" example:"true"`
	Code    int                        `json:"code" example:"200"`
	Data    SessionListResponseData    `json:"data"`
	Error   *SessionErrorResponse      `json:"error,omitempty"`
}

// SessionListResponseData contains the data for session list response
type SessionListResponseData struct {
	Action    string        `json:"action" example:"list"`
	Status    string        `json:"status" example:"success"`
	Timestamp time.Time     `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Sessions  []SessionInfo `json:"sessions"`
	Total     int           `json:"total" example:"5"`
}

// ConnectSessionResponse represents the response for session connection
type ConnectSessionResponse struct {
	Success bool                           `json:"success" example:"true"`
	Code    int                            `json:"code" example:"200"`
	Data    ConnectSessionResponseData     `json:"data"`
	Error   *SessionErrorResponse          `json:"error,omitempty"`
}

// ConnectSessionResponseData contains the data for session connection response
type ConnectSessionResponseData struct {
	SessionID  string                 `json:"session_id" example:"default"`
	Action     string                 `json:"action" example:"connect"`
	Status     string                 `json:"status" example:"success"`
	Timestamp  time.Time              `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session    *SessionInfo           `json:"session,omitempty"`
	Connection *SessionConnectionInfo `json:"connection,omitempty"`
	QRCode     string                 `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
}

// QRCodeResponse represents the response for QR code retrieval
type QRCodeResponse struct {
	Success bool                  `json:"success" example:"true"`
	Code    int                   `json:"code" example:"200"`
	Data    QRCodeResponseData    `json:"data"`
	Error   *SessionErrorResponse `json:"error,omitempty"`
}

// QRCodeResponseData contains the data for QR code response
type QRCodeResponseData struct {
	SessionID string    `json:"session_id" example:"default"`
	Action    string    `json:"action" example:"qr"`
	Status    string    `json:"status" example:"success"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	QRCode    string    `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
}

// PairPhoneResponse represents the response for phone pairing
type PairPhoneResponse struct {
	Success bool                     `json:"success" example:"true"`
	Code    int                      `json:"code" example:"200"`
	Data    PairPhoneResponseData    `json:"data"`
	Error   *SessionErrorResponse    `json:"error,omitempty"`
}

// PairPhoneResponseData contains the data for phone pairing response
type PairPhoneResponseData struct {
	SessionID   string    `json:"session_id" example:"default"`
	Action      string    `json:"action" example:"pair"`
	Status      string    `json:"status" example:"success"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	PhoneNumber string    `json:"phone_number" example:"5511999999999"`
	PairCode    string    `json:"pair_code" example:"ABCD-1234"`
}

// SessionStatusResponse represents the response for session status
type SessionStatusResponse struct {
	Success bool                         `json:"success" example:"true"`
	Code    int                          `json:"code" example:"200"`
	Data    SessionStatusResponseData    `json:"data"`
	Error   *SessionErrorResponse        `json:"error,omitempty"`
}

// SessionStatusResponseData contains the data for session status response
type SessionStatusResponseData struct {
	SessionID    string    `json:"session_id" example:"default"`
	Action       string    `json:"action" example:"status"`
	Status       string    `json:"status" example:"success"`
	Timestamp    time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Name         string    `json:"name" example:"default"`
	SessionStatus string   `json:"session_status" example:"connected"`
	WaJID        string    `json:"wa_jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	IsConnected  bool      `json:"is_connected" example:"true"`
	ClientStatus string    `json:"client_status" example:"connected"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// NewSessionSuccessResponse creates a successful session operation response
func NewSessionSuccessResponse(sessionID, action string, session *SessionInfo) *SessionResponse {
	return &SessionResponse{
		Success: true,
		Code:    200,
		Data: SessionData{
			SessionID: sessionID,
			Action:    action,
			Status:    "success",
			Timestamp: time.Now(),
			Session:   session,
		},
	}
}

// NewSessionErrorResponse creates an error response for session operations
func NewSessionErrorResponse(code int, errorCode, message, details string) *SessionResponse {
	return &SessionResponse{
		Success: false,
		Code:    code,
		Data: SessionData{
			Status:    "error",
			Timestamp: time.Now(),
		},
		Error: &SessionErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
}
