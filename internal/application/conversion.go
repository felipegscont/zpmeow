package application

import (
	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}


func (c *Converter) SessionToInfo(session *session.Session) dto.SessionInfo {
	return dto.SessionInfo{
		ID:        session.ID.Value(),
		Name:      session.Name.Value(),
		Status:    string(session.Status),
		WaJID:     session.WaJID.Value(),
		ProxyURL:  session.ProxyURL.Value(),
		CreatedAt: session.CreatedAt,
	}
}

func (c *Converter) SessionsToInfoList(sessions []*session.Session) []dto.SessionInfo {
	result := make([]dto.SessionInfo, len(sessions))
	for i, s := range sessions {
		result[i] = c.SessionToInfo(s)
	}
	return result
}
