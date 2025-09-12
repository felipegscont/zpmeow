package services

import (
	"zpmeow/internal/application/dto/session"
	sessionDomain "zpmeow/internal/domain/session"
)

// SessionConverter handles conversion between domain entities and DTOs
type SessionConverter struct{}

// ToDTO converts a domain Session to DTO
func (c *SessionConverter) ToDTO(s *sessionDomain.Session) session.SessionInfoResponse {
	return session.SessionInfoResponse{
		BaseSessionInfo: session.BaseSessionInfo{
			ID:        s.ID,
			Name:      s.Name,
			Status:    string(s.Status),
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		WhatsAppJID: s.WhatsAppJID,
		QRCode:      s.QRCode,
		ProxyURL:    s.ProxyURL,
	}
}

// ToDTOBatch converts multiple domain Sessions to DTOs
func (c *SessionConverter) ToDTOBatch(sessions []*sessionDomain.Session) []session.SessionInfoResponse {
	result := make([]session.SessionInfoResponse, len(sessions))
	for i, s := range sessions {
		result[i] = c.ToDTO(s)
	}
	return result
}

// Global instance
var SessionToDTOConverter = &SessionConverter{}
