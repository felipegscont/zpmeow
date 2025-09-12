package request

// CreateSessionRequest represents a request to create a new session
type CreateSessionRequest struct {
	Name string `json:"name" binding:"required" example:"My WhatsApp Session"`
}

// PairSessionRequest represents a request to pair a session with a phone number
type PairSessionRequest struct {
	SessionID   string `json:"session_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	PhoneNumber string `json:"phone_number" binding:"required" example:"5511999999999"`
}

// ProxyRequest represents a request to set proxy configuration
type ProxyRequest struct {
	ProxyURL string `json:"proxy_url" binding:"required" example:"http://proxy.example.com:8080"`
}

// GetSessionRequest represents a request to get session information
type GetSessionRequest struct {
	IDOrName string `json:"id_or_name" binding:"required" example:"my-session"`
}

// DeleteSessionRequest represents a request to delete a session
type DeleteSessionRequest struct {
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// ConnectSessionRequest represents a request to connect a session
type ConnectSessionRequest struct {
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// DisconnectSessionRequest represents a request to disconnect a session
type DisconnectSessionRequest struct {
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// GetQRCodeRequest represents a request to get QR code
type GetQRCodeRequest struct {
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// SetProxyRequest represents a request to set proxy
type SetProxyRequest struct {
	ID       string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProxyURL string `json:"proxy_url" binding:"required" example:"http://proxy.example.com:8080"`
}

// ClearProxyRequest represents a request to clear proxy
type ClearProxyRequest struct {
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// GetSessionStatusRequest represents a request to get session status
type GetSessionStatusRequest struct {
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// LogoutSessionRequest represents a request to logout session
type LogoutSessionRequest struct {
	ID string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}
