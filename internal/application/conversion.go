package application

import (
	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
)

// Converter handles conversions between domain entities and DTOs
type Converter struct{}

// NewConverter creates a new conversion service
func NewConverter() *Converter {
	return &Converter{}
}

// Session Conversions

// SessionToInfo converts a session entity to session info DTO
func (c *Converter) SessionToInfo(session *session.Session) dto.SessionInfo {
	return dto.SessionInfo{
		ID:         session.ID.Value(),
		Name:       session.Name.Value(),
		Status:     string(session.Status),
		WaJID:      session.WaJID,
		ProxyURL:   session.ProxyURL.Value(),
		WebhookURL: session.WebhookURL,
		Events:     session.Events,
		CreatedAt:  session.CreatedAt,
		UpdatedAt:  session.UpdatedAt,
	}
}

// SessionsToInfoList converts a list of session entities to session info DTOs
func (c *Converter) SessionsToInfoList(sessions []*session.Session) []dto.SessionInfo {
	result := make([]dto.SessionInfo, len(sessions))
	for i, s := range sessions {
		result[i] = c.SessionToInfo(s)
	}
	return result
}
