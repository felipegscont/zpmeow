package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/application"
	"zpmeow/internal/domain"
	"zpmeow/internal/shared"
)

// Error handling utilities
func MapDomainError(err error) (statusCode int, message string) {
	return shared.MapDomainError(err)
}

func IsValidationError(err error) bool {
	return shared.IsValidationError(err)
}

func IsConflictError(err error) bool {
	return shared.IsConflictError(err)
}

func IsNotFoundError(err error) bool {
	return shared.IsNotFoundError(err)
}

// Use the converter from application services
var DefaultConverter = application.SessionToDTOConverter

func ToSessionInfoResponse(sess *domain.Session) application.SessionInfoResponse {
	return DefaultConverter.ToDTO(sess)
}

func ToSessionListResponse(sessions []*domain.Session) application.SessionListResponse {
	sessionResponses := DefaultConverter.ToDTOBatch(sessions)
	return application.SessionListResponse{
		Sessions: sessionResponses,
		Total:    len(sessionResponses),
	}
}

func ToCreateSessionResponse(sess *domain.Session) application.CreateSessionResponse {
	dto := DefaultConverter.ToDTO(sess)
	return application.CreateSessionResponse{
		ID:        dto.ID,
		Name:      dto.Name,
		Status:    dto.Status,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func ToQRCodeResponse(qrCode string, sess *domain.Session) application.QRCodeResponse {
	return application.QRCodeResponse{
		QRCode: qrCode,
		Status: string(sess.Status),
	}
}

func ToProxyResponse(proxyURL, message string) application.ProxyResponse {
	return application.ProxyResponse{
		ProxyURL: proxyURL,
		Message:  message,
	}
}

func ToPairSessionResponse(pairingCode string) application.PairSessionResponse {
	return application.PairSessionResponse{
		PairingCode: pairingCode,
	}
}

func ToMessageResponse(message string) application.MessageResponse {
	return application.MessageResponse{
		Message: message,
	}
}

// HTTP response utilities
func RespondWithError(c *gin.Context, statusCode int, message string, details ...string) {
	response := gin.H{"error": message}
	if len(details) > 0 {
		response["details"] = details[0]
	}
	c.JSON(statusCode, response)
}

func RespondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// Validation utilities
func ValidateSessionIDParam(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return "", false
	}
	return id, true
}

func ValidateAndBindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return false
	}
	return true
}
