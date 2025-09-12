package session

// CreateSessionRequest represents a request to create a new session
type CreateSessionRequest struct {
	Name string `json:"name" binding:"required" example:"My WhatsApp Session"`
}

// PairSessionRequest represents a request to pair a session with a phone number
type PairSessionRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"5511999999999"`
}

// ProxyRequest represents a request to set proxy configuration
type ProxyRequest struct {
	ProxyURL string `json:"proxy_url" binding:"required" example:"http://proxy.example.com:8080"`
}
