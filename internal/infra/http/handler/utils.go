package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/utils"
)

// ============================================================================
// Error Mapping
// ============================================================================

// ErrorMapping defines the mapping between domain errors and HTTP status codes
type ErrorMapping struct {
	StatusCode int
	Message    string
}

// domainErrorMappings maps domain errors to HTTP responses
var domainErrorMappings = map[error]ErrorMapping{
	// Session validation errors (400 Bad Request)
	session.ErrInvalidSessionID:          {http.StatusBadRequest, "Invalid session ID"},
	session.ErrInvalidSessionName:        {http.StatusBadRequest, "Invalid session name"},
	session.ErrSessionNameTooShort:       {http.StatusBadRequest, "Session name too short"},
	session.ErrSessionNameTooLong:        {http.StatusBadRequest, "Session name too long"},
	session.ErrInvalidSessionNameChar:    {http.StatusBadRequest, "Session name contains invalid characters"},
	session.ErrInvalidSessionNameFormat:  {http.StatusBadRequest, "Invalid session name format"},
	session.ErrReservedSessionName:       {http.StatusBadRequest, "Session name is reserved"},
	session.ErrInvalidSessionStatus:      {http.StatusBadRequest, "Invalid session status"},
	
	// Session state errors (409 Conflict)
	session.ErrSessionAlreadyExists:      {http.StatusConflict, "Session already exists"},
	session.ErrSessionAlreadyConnected:   {http.StatusConflict, "Session is already connected"},
	session.ErrSessionCannotConnect:      {http.StatusConflict, "Session cannot be connected in current state"},
	
	// Session not found errors (404 Not Found)
	session.ErrSessionNotFound:           {http.StatusNotFound, "Session not found"},
}

// MapDomainError maps a domain error to an HTTP status code and message
// Returns the mapped status code and message, or 500 Internal Server Error for unmapped errors
func MapDomainError(err error) (statusCode int, message string) {
	if mapping, exists := domainErrorMappings[err]; exists {
		return mapping.StatusCode, mapping.Message
	}
	
	// Default to internal server error for unmapped errors
	return http.StatusInternalServerError, "Internal server error"
}

// IsValidationError checks if an error is a validation-related error
func IsValidationError(err error) bool {
	mapping, exists := domainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusBadRequest
}

// IsConflictError checks if an error is a conflict-related error
func IsConflictError(err error) bool {
	mapping, exists := domainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusConflict
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	mapping, exists := domainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusNotFound
}

// ============================================================================
// DTO Conversion
// ============================================================================

// SessionToDTOConverter provides centralized conversion functions for Session entities to DTOs
type SessionToDTOConverter struct{}

// NewSessionToDTOConverter creates a new converter instance
func NewSessionToDTOConverter() *SessionToDTOConverter {
	return &SessionToDTOConverter{}
}

// ToCreateSessionResponse converts a Session entity to CreateSessionResponse DTO
func (c *SessionToDTOConverter) ToCreateSessionResponse(sess *session.Session) session.CreateSessionResponse {
	return session.CreateSessionResponse{
		ID:        sess.ID,
		Name:      sess.Name,
		Status:    string(sess.Status),
		CreatedAt: sess.CreatedAt,
		UpdatedAt: sess.UpdatedAt,
	}
}

// ToSessionInfoResponse converts a Session entity to SessionInfoResponse DTO
func (c *SessionToDTOConverter) ToSessionInfoResponse(sess *session.Session) session.SessionInfoResponse {
	return session.SessionInfoResponse{
		BaseSessionInfo: session.BaseSessionInfo{
			ID:        sess.ID,
			Name:      sess.Name,
			Status:    string(sess.Status),
			CreatedAt: sess.CreatedAt,
			UpdatedAt: sess.UpdatedAt,
		},
		WhatsAppJID: sess.WhatsAppJID,
		QRCode:      sess.QRCode,
		ProxyURL:    sess.ProxyURL,
	}
}

// ToSessionListResponse converts a slice of Session entities to SessionListResponse DTO
func (c *SessionToDTOConverter) ToSessionListResponse(sessions []*session.Session) session.SessionListResponse {
	sessionResponses := make([]session.SessionInfoResponse, len(sessions))
	for i, sess := range sessions {
		sessionResponses[i] = c.ToSessionInfoResponse(sess)
	}

	return session.SessionListResponse{
		Sessions: sessionResponses,
		Total:    len(sessionResponses),
	}
}

// ToQRCodeResponse converts session data to QRCodeResponse DTO
func (c *SessionToDTOConverter) ToQRCodeResponse(qrCode string, sess *session.Session) session.QRCodeResponse {
	return session.QRCodeResponse{
		QRCode: qrCode,
		Status: string(sess.Status),
	}
}

// ToProxyResponse converts proxy data to ProxyResponse DTO
func (c *SessionToDTOConverter) ToProxyResponse(proxyURL, message string) session.ProxyResponse {
	return session.ProxyResponse{
		ProxyURL: proxyURL,
		Message:  message,
	}
}

// ToPairSessionResponse converts pairing code to PairSessionResponse DTO
func (c *SessionToDTOConverter) ToPairSessionResponse(pairingCode string) session.PairSessionResponse {
	return session.PairSessionResponse{
		PairingCode: pairingCode,
	}
}

// ToMessageResponse creates a generic message response
func (c *SessionToDTOConverter) ToMessageResponse(message string) session.MessageResponse {
	return session.MessageResponse{
		Message: message,
	}
}

// Global converter instance for convenience
var DefaultConverter = NewSessionToDTOConverter()

// Convenience functions using the default converter
func ToCreateSessionResponse(sess *session.Session) session.CreateSessionResponse {
	return DefaultConverter.ToCreateSessionResponse(sess)
}

func ToSessionInfoResponse(sess *session.Session) session.SessionInfoResponse {
	return DefaultConverter.ToSessionInfoResponse(sess)
}

func ToSessionListResponse(sessions []*session.Session) session.SessionListResponse {
	return DefaultConverter.ToSessionListResponse(sessions)
}

func ToQRCodeResponse(qrCode string, sess *session.Session) session.QRCodeResponse {
	return DefaultConverter.ToQRCodeResponse(qrCode, sess)
}

func ToProxyResponse(proxyURL, message string) session.ProxyResponse {
	return DefaultConverter.ToProxyResponse(proxyURL, message)
}

func ToPairSessionResponse(pairingCode string) session.PairSessionResponse {
	return DefaultConverter.ToPairSessionResponse(pairingCode)
}

func ToMessageResponse(message string) session.MessageResponse {
	return DefaultConverter.ToMessageResponse(message)
}

// ============================================================================
// Parameter Validation Helpers
// ============================================================================

// ValidateSessionIDParam validates and extracts session ID from URL parameter
// Returns the session ID if valid, or responds with error and returns empty string
func ValidateSessionIDParam(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return "", false
	}
	return id, true
}

// ValidateAndBindJSON validates and binds JSON request body
// Returns true if successful, false if validation failed (response already sent)
func ValidateAndBindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return false
	}
	return true
}
