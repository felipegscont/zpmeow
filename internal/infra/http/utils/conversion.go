package utils

import (
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

// DTO conversion utilities
func ToCreateSessionResponse(session *domain.Session) application.CreateSessionResponse {
	return application.CreateSessionResponse{
		ID:     session.ID,
		Name:   session.Name,
		Status: string(session.Status),
	}
}

func ToQRCodeResponse(qrCode string) application.QRCodeResponse {
	return application.QRCodeResponse{
		QRCode: qrCode,
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
