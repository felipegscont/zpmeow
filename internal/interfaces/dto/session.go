package dto

import (
	"time"
)


type CreateSessionRequest struct {
	Name       string `json:"name" binding:"required" example:"default"`
	ProxyURL   string `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
	WebhookURL string `json:"webhook_url,omitempty" example:"https://webhook.example.com/meow"`
	Events     string `json:"events,omitempty" example:"message,status"`
}

type PairPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"5511999999999"`
}


type SessionInfo struct {
	ID         string    `json:"id" example:"default"`
	Name       string    `json:"name" example:"default"`
	Status     string    `json:"status" example:"connected"`
	WaJID      string    `json:"wa_jid,omitempty" example:"5511999999999@s.meow.net"`
	ProxyURL   string    `json:"proxy_url,omitempty" example:"http://proxy.example.com:8080"`
	WebhookURL string    `json:"webhook_url,omitempty" example:"https://webhook.example.com/meow"`
	Events     []string  `json:"events,omitempty" example:"message,status"`
	ApiKey     string    `json:"api_key,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

type SessionConnectionInfo struct {
	QRCode      string    `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Connected   bool      `json:"connected" example:"true"`
	IsConnected bool      `json:"is_connected" example:"true"`
	LastSeen    time.Time `json:"last_seen,omitempty" example:"2023-01-01T00:00:00Z"`
	PairCode    string    `json:"pair_code,omitempty" example:"ABCD-1234"`
}


type SessionResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Data    SessionData           `json:"data"`
	Error   *SessionErrorResponse `json:"error,omitempty"`
}

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

type SessionErrorResponse struct {
	Code    string `json:"code" example:"INVALID_SESSION_ID"`
	Message string `json:"message" example:"Invalid session ID format"`
	Details string `json:"details,omitempty" example:"Session ID must be alphanumeric"`
}


type CreateSessionResponse struct {
	Success bool                  `json:"success" example:"true"`
	Code    int                   `json:"code" example:"201"`
	Data    SessionCreateData     `json:"data"`
	Error   *SessionErrorResponse `json:"error,omitempty"`
}

type SessionCreateData struct {
	SessionID string       `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Action    string       `json:"action" example:"create"`
	Status    string       `json:"status" example:"success"`
	Timestamp time.Time    `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session   *SessionInfo `json:"session"`
}

type SessionInfoResponse struct {
	Success bool                    `json:"success" example:"true"`
	Code    int                     `json:"code" example:"200"`
	Data    SessionInfoResponseData `json:"data"`
	Error   *SessionErrorResponse   `json:"error,omitempty"`
}

type SessionInfoResponseData struct {
	SessionID string       `json:"session_id" example:"default"`
	Action    string       `json:"action" example:"get"`
	Status    string       `json:"status" example:"success"`
	Timestamp time.Time    `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session   *SessionInfo `json:"session"`
}

type SessionListResponse struct {
	Success bool                    `json:"success" example:"true"`
	Code    int                     `json:"code" example:"200"`
	Data    SessionListResponseData `json:"data"`
	Error   *SessionErrorResponse   `json:"error,omitempty"`
}

type SessionListResponseData struct {
	Action    string        `json:"action" example:"list"`
	Status    string        `json:"status" example:"success"`
	Timestamp time.Time     `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Sessions  []SessionInfo `json:"sessions"`
	Total     int           `json:"total" example:"5"`
}

type ConnectSessionResponse struct {
	Success bool                  `json:"success" example:"true"`
	Code    int                   `json:"code" example:"200"`
	Data    SessionConnectData    `json:"data"`
	Error   *SessionErrorResponse `json:"error,omitempty"`
}

type SessionConnectData struct {
	SessionID  string                 `json:"session_id" example:"default"`
	Action     string                 `json:"action" example:"connect"`
	Status     string                 `json:"status" example:"success"`
	Timestamp  time.Time              `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session    *SessionInfo           `json:"session,omitempty"`
	Connection *SessionConnectionInfo `json:"connection,omitempty"`
	QRCode     string                 `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
}

type QRCodeResponse struct {
	Success bool                  `json:"success" example:"true"`
	Code    int                   `json:"code" example:"200"`
	Data    QRCodeResponseData    `json:"data"`
	Error   *SessionErrorResponse `json:"error,omitempty"`
}

type QRCodeResponseData struct {
	SessionID string    `json:"session_id" example:"default"`
	Action    string    `json:"action" example:"qr"`
	Status    string    `json:"status" example:"success"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	QRCode    string    `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
}

type PairPhoneResponse struct {
	Success bool                  `json:"success" example:"true"`
	Code    int                   `json:"code" example:"200"`
	Data    PairPhoneResponseData `json:"data"`
	Error   *SessionErrorResponse `json:"error,omitempty"`
}

type PairPhoneResponseData struct {
	SessionID   string    `json:"session_id" example:"default"`
	Action      string    `json:"action" example:"pair"`
	Status      string    `json:"status" example:"success"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	PhoneNumber string    `json:"phone_number" example:"5511999999999"`
	PairCode    string    `json:"pair_code" example:"ABCD-1234"`
}

type SessionStatusResponse struct {
	Success bool                      `json:"success" example:"true"`
	Code    int                       `json:"code" example:"200"`
	Data    SessionStatusResponseData `json:"data"`
	Error   *SessionErrorResponse     `json:"error,omitempty"`
}

type SessionStatusResponseData struct {
	SessionID     string    `json:"session_id" example:"default"`
	Action        string    `json:"action" example:"status"`
	Status        string    `json:"status" example:"success"`
	Timestamp     time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Name          string    `json:"name" example:"default"`
	SessionStatus string    `json:"session_status" example:"connected"`
	WaJID         string    `json:"wa_jid,omitempty" example:"5511999999999@s.meow.net"`
	IsConnected   bool      `json:"is_connected" example:"true"`
	ClientStatus  string    `json:"client_status" example:"connected"`
	CreatedAt     time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}


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
