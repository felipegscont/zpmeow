package dto

import (
	"time"
)

// Session-related request types

// CreateSessionRequest represents a request to create a new session
type CreateSessionRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=50" binding:"required" example:"default"`
	WebhookURL string `json:"webhook_url,omitempty" validate:"omitempty,webhook_url" example:"https://webhook.example.com/whatsapp"`
	Events     string `json:"events,omitempty" validate:"omitempty,max=500" example:"message,status"`
	ProxyURL   string `json:"proxy_url,omitempty" validate:"omitempty,url" example:"http://proxy.example.com:8080"`
}

// Validate validates the CreateSessionRequest
func (r *CreateSessionRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

// PairPhoneRequest represents a request to pair a phone number
type PairPhoneRequest struct {
	Phone string `json:"phone" validate:"required,phone_number" binding:"required" example:"5511999999999"`
}

// Validate validates the PairPhoneRequest
func (r *PairPhoneRequest) Validate() error {
	validator := NewValidator()
	return validator.Validate(r)
}

// Session-related response types

// SessionResponse represents a generic session response
type SessionResponse struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Data      SessionData `json:"data"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SessionData represents session response data
type SessionData struct {
	SessionID string        `json:"session_id" example:"session_123"`
	Action    string        `json:"action" example:"create"`
	Status    string        `json:"status" example:"success"`
	Timestamp time.Time     `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session   *SessionInfo  `json:"session,omitempty"`
	Sessions  []SessionInfo `json:"sessions,omitempty"`
	QRCode    string        `json:"qr_code,omitempty"`
}

// CreateSessionResponse represents a response to session creation
type CreateSessionResponse struct {
	Success   bool              `json:"success"`
	Code      int               `json:"code"`
	Data      SessionCreateData `json:"data"`
	Error     *ErrorInfo        `json:"error,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// SessionCreateData represents session creation response data
type SessionCreateData struct {
	SessionID string       `json:"session_id" example:"session_123"`
	Action    string       `json:"action" example:"create"`
	Status    string       `json:"status" example:"success"`
	Timestamp time.Time    `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session   *SessionInfo `json:"session"`
}

// ConnectSessionResponse represents a response to session connection
type ConnectSessionResponse struct {
	Success   bool               `json:"success"`
	Code      int                `json:"code"`
	Data      SessionConnectData `json:"data"`
	Error     *ErrorInfo         `json:"error,omitempty"`
	Timestamp time.Time          `json:"timestamp"`
}

// SessionConnectData represents session connection response data
type SessionConnectData struct {
	SessionID  string                 `json:"session_id" example:"session_123"`
	Action     string                 `json:"action" example:"connect"`
	Status     string                 `json:"status" example:"success"`
	Timestamp  time.Time              `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Session    *SessionInfo           `json:"session"`
	Connection *SessionConnectionInfo `json:"connection"`
	QRCode     string                 `json:"qr_code,omitempty"`
}

// SessionConnectionInfo represents session connection information
type SessionConnectionInfo struct {
	QRCode      string `json:"qr_code,omitempty"`
	Connected   bool   `json:"connected"`
	IsConnected bool   `json:"is_connected"`
}

// PairPhoneResponse represents a response to phone pairing
type PairPhoneResponse struct {
	Success   bool                  `json:"success"`
	Code      int                   `json:"code"`
	Data      PairPhoneResponseData `json:"data"`
	Error     *ErrorInfo            `json:"error,omitempty"`
	Timestamp time.Time             `json:"timestamp"`
}

// PairPhoneResponseData represents phone pairing response data
type PairPhoneResponseData struct {
	SessionID string    `json:"session_id" example:"session_123"`
	Action    string    `json:"action" example:"pair"`
	Status    string    `json:"status" example:"success"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Phone     string    `json:"phone" example:"5511999999999"`
	Code      string    `json:"code,omitempty" example:"123456"`
}

// SessionStatusResponse represents a response to session status request
type SessionStatusResponse struct {
	Success   bool                      `json:"success"`
	Code      int                       `json:"code"`
	Data      SessionStatusResponseData `json:"data"`
	Error     *ErrorInfo                `json:"error,omitempty"`
	Timestamp time.Time                 `json:"timestamp"`
}

// SessionStatusResponseData represents session status response data
type SessionStatusResponseData struct {
	SessionID     string    `json:"session_id" example:"session_123"`
	Action        string    `json:"action" example:"status"`
	Status        string    `json:"status" example:"success"`
	Timestamp     time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Name          string    `json:"name" example:"My Session"`
	SessionStatus string    `json:"session_status" example:"connected"`
	WaJID         string    `json:"wa_jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	IsConnected   bool      `json:"is_connected"`
	ClientStatus  string    `json:"client_status" example:"connected"`
	CreatedAt     time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SessionErrorResponse represents an error response for session operations
type SessionErrorResponse struct {
	Success   bool       `json:"success"`
	Code      int        `json:"code"`
	Error     *ErrorInfo `json:"error"`
	Timestamp time.Time  `json:"timestamp"`
}

// SessionListResponse represents a response containing multiple sessions
type SessionListResponse struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Data      SessionData `json:"data"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Response constructors

// NewSessionSuccessResponse creates a standardized success response for session operations
func NewSessionSuccessResponse(sessionID, action string, data interface{}) *SessionResponse {
	response := &SessionResponse{
		Success: true,
		Code:    200,
		Data: SessionData{
			SessionID: sessionID,
			Action:    action,
			Status:    "success",
			Timestamp: time.Now(),
		},
		Timestamp: time.Now(),
	}

	switch v := data.(type) {
	case *SessionInfo:
		response.Data.Session = v
	case []SessionInfo:
		response.Data.Sessions = v
	case string:
		response.Data.QRCode = v
	}

	return response
}

// NewSessionErrorResponse creates a standardized error response for session operations
func NewSessionErrorResponse(code int, errorCode, message, details string) *SessionResponse {
	return &SessionResponse{
		Success: false,
		Code:    code,
		Data: SessionData{
			Status:    "error",
			Timestamp: time.Now(),
		},
		Error: &ErrorInfo{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

// NewCreateSessionSuccessResponse creates a success response for session creation
func NewCreateSessionSuccessResponse(sessionInfo *SessionInfo) *CreateSessionResponse {
	return &CreateSessionResponse{
		Success: true,
		Code:    201,
		Data: SessionCreateData{
			SessionID: sessionInfo.ID,
			Action:    "create",
			Status:    "success",
			Timestamp: time.Now(),
			Session:   sessionInfo,
		},
		Timestamp: time.Now(),
	}
}

// NewConnectSessionSuccessResponse creates a success response for session connection
func NewConnectSessionSuccessResponse(sessionInfo *SessionInfo, connectionInfo *SessionConnectionInfo, qrCode string) *ConnectSessionResponse {
	return &ConnectSessionResponse{
		Success: true,
		Code:    200,
		Data: SessionConnectData{
			SessionID:  sessionInfo.ID,
			Action:     "connect",
			Status:     "success",
			Timestamp:  time.Now(),
			Session:    sessionInfo,
			Connection: connectionInfo,
			QRCode:     qrCode,
		},
		Timestamp: time.Now(),
	}
}

// NewPairPhoneSuccessResponse creates a success response for phone pairing
func NewPairPhoneSuccessResponse(sessionID, phone, code string) *PairPhoneResponse {
	return &PairPhoneResponse{
		Success: true,
		Code:    200,
		Data: PairPhoneResponseData{
			SessionID: sessionID,
			Action:    "pair",
			Status:    "success",
			Timestamp: time.Now(),
			Phone:     phone,
			Code:      code,
		},
		Timestamp: time.Now(),
	}
}

// NewSessionStatusSuccessResponse creates a success response for session status
func NewSessionStatusSuccessResponse(data SessionStatusResponseData) *SessionStatusResponse {
	return &SessionStatusResponse{
		Success:   true,
		Code:      200,
		Data:      data,
		Timestamp: time.Now(),
	}
}
