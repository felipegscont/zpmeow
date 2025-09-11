package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/utils"
)






type ErrorMapping struct {
	StatusCode int
	Message    string
}


var domainErrorMappings = map[error]ErrorMapping{

	session.ErrInvalidSessionID:          {http.StatusBadRequest, "Invalid session ID"},
	session.ErrInvalidSessionName:        {http.StatusBadRequest, "Invalid session name"},
	session.ErrSessionNameTooShort:       {http.StatusBadRequest, "Session name too short"},
	session.ErrSessionNameTooLong:        {http.StatusBadRequest, "Session name too long"},
	session.ErrInvalidSessionNameChar:    {http.StatusBadRequest, "Session name contains invalid characters"},
	session.ErrInvalidSessionNameFormat:  {http.StatusBadRequest, "Invalid session name format"},
	session.ErrReservedSessionName:       {http.StatusBadRequest, "Session name is reserved"},
	session.ErrInvalidSessionStatus:      {http.StatusBadRequest, "Invalid session status"},
	

	session.ErrSessionAlreadyExists:      {http.StatusConflict, "Session already exists"},
	session.ErrSessionAlreadyConnected:   {http.StatusConflict, "Session is already connected"},
	session.ErrSessionCannotConnect:      {http.StatusConflict, "Session cannot be connected in current state"},
	

	session.ErrSessionNotFound:           {http.StatusNotFound, "Session not found"},
}



func MapDomainError(err error) (statusCode int, message string) {
	if mapping, exists := domainErrorMappings[err]; exists {
		return mapping.StatusCode, mapping.Message
	}
	

	return http.StatusInternalServerError, "Internal server error"
}


func IsValidationError(err error) bool {
	mapping, exists := domainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusBadRequest
}


func IsConflictError(err error) bool {
	mapping, exists := domainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusConflict
}


func IsNotFoundError(err error) bool {
	mapping, exists := domainErrorMappings[err]
	return exists && mapping.StatusCode == http.StatusNotFound
}






type SessionToDTOConverter struct{}


func NewSessionToDTOConverter() *SessionToDTOConverter {
	return &SessionToDTOConverter{}
}


func (c *SessionToDTOConverter) ToCreateSessionResponse(sess *session.Session) session.CreateSessionResponse {
	return session.CreateSessionResponse{
		ID:        sess.ID,
		Name:      sess.Name,
		Status:    string(sess.Status),
		CreatedAt: sess.CreatedAt,
		UpdatedAt: sess.UpdatedAt,
	}
}


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


func (c *SessionToDTOConverter) ToQRCodeResponse(qrCode string, sess *session.Session) session.QRCodeResponse {
	return session.QRCodeResponse{
		QRCode: qrCode,
		Status: string(sess.Status),
	}
}


func (c *SessionToDTOConverter) ToProxyResponse(proxyURL, message string) session.ProxyResponse {
	return session.ProxyResponse{
		ProxyURL: proxyURL,
		Message:  message,
	}
}


func (c *SessionToDTOConverter) ToPairSessionResponse(pairingCode string) session.PairSessionResponse {
	return session.PairSessionResponse{
		PairingCode: pairingCode,
	}
}


func (c *SessionToDTOConverter) ToMessageResponse(message string) session.MessageResponse {
	return session.MessageResponse{
		Message: message,
	}
}


var DefaultConverter = NewSessionToDTOConverter()


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







func ValidateSessionIDParam(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return "", false
	}
	return id, true
}



func ValidateAndBindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return false
	}
	return true
}
