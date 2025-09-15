package application

import (
	"zpmeow/internal/domain/session"
	"zpmeow/internal/interfaces/dto"
)

// ConversionService handles conversions between domain entities and DTOs
type ConversionService struct{}

// NewConversionService creates a new conversion service
func NewConversionService() *ConversionService {
	return &ConversionService{}
}

// Session Conversions

// SessionToInfo converts a session entity to session info DTO
func (c *ConversionService) SessionToInfo(session *session.Session) dto.SessionInfo {
	return dto.SessionInfo{
		ID:         session.ID,
		Name:       session.Name,
		Status:     string(session.Status),
		WaJID:      session.WaJID,
		ProxyURL:   session.ProxyURL,
		WebhookURL: session.WebhookURL,
		Events:     session.Events,
		CreatedAt:  session.CreatedAt,
		UpdatedAt:  session.UpdatedAt,
	}
}

// SessionsToInfoList converts a list of session entities to session info DTOs
func (c *ConversionService) SessionsToInfoList(sessions []*session.Session) []dto.SessionInfo {
	result := make([]dto.SessionInfo, len(sessions))
	for i, s := range sessions {
		result[i] = c.SessionToInfo(s)
	}
	return result
}
