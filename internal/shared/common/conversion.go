package common

import (
	domainSessions "zpmeow/internal/domain/session"
	"zpmeow/internal/interfaces/dto"
)

type SessionConverter struct{}

func (c *SessionConverter) ToDTO(s *domainSessions.Session) dto.SessionInfo {
	return dto.SessionInfo{
		ID:         s.ID,
		Name:       s.Name,
		Status:     string(s.Status),
		WaJID:      s.WaJID,
		ProxyURL:   s.ProxyURL,
		WebhookURL: s.WebhookURL,
		Events:     s.Events,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func (c *SessionConverter) ToListDTO(s *domainSessions.Session) dto.SessionInfo {
	return dto.SessionInfo{
		ID:        s.ID,
		Name:      s.Name,
		Status:    string(s.Status),
		WaJID:     s.WaJID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (c *SessionConverter) ToDTOBatch(sessions []*domainSessions.Session) []dto.SessionInfo {
	result := make([]dto.SessionInfo, len(sessions))
	for i, s := range sessions {
		result[i] = c.ToListDTO(s)
	}
	return result
}

var SessionToDTOConverter = &SessionConverter{}
