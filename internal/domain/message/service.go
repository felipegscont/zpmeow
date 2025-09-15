package message

import (
	"fmt"
	"strings"
)

type Service interface {
	ValidateMessage(msg *Message) error
	BuildTextMessage(sessionID, chatJID, content string) (*Message, error)
	BuildMediaMessage(sessionID, chatJID string, messageType MessageType, mediaURL string) (*Message, error)
	FormatJID(jid string) string
}

type DomainService struct{}

func NewDomainService() Service {
	return &DomainService{}
}

func (s *DomainService) ValidateMessage(msg *Message) error {
	if msg == nil {
		return NewValidationError("message cannot be nil")
	}

	if msg.SessionID == "" {
		return NewValidationError("session ID is required")
	}

	if msg.ChatJID == "" {
		return NewValidationError("chat JID is required")
	}

	if msg.MessageType == "" {
		return NewValidationError("message type is required")
	}

	switch msg.MessageType {
	case MessageTypeText:
		if strings.TrimSpace(msg.Content) == "" {
			return NewValidationError("text content cannot be empty")
		}
	case MessageTypeImage, MessageTypeAudio, MessageTypeVideo, MessageTypeDocument:
		if msg.MediaURL == "" {
			return NewValidationError("media URL is required for media messages")
		}
	}

	return nil
}

func (s *DomainService) BuildTextMessage(sessionID, chatJID, content string) (*Message, error) {
	if sessionID == "" {
		return nil, NewValidationError("session ID is required")
	}

	if chatJID == "" {
		return nil, NewValidationError("chat JID is required")
	}

	if strings.TrimSpace(content) == "" {
		return nil, NewValidationError("text content cannot be empty")
	}

	formattedJID := s.FormatJID(chatJID)
	msg := NewMessage(sessionID, formattedJID, "", formattedJID, MessageTypeText, content)

	return msg, nil
}

func (s *DomainService) BuildMediaMessage(sessionID, chatJID string, messageType MessageType, mediaURL string) (*Message, error) {
	if sessionID == "" {
		return nil, NewValidationError("session ID is required")
	}

	if chatJID == "" {
		return nil, NewValidationError("chat JID is required")
	}

	if mediaURL == "" {
		return nil, NewValidationError("media URL is required")
	}

	switch messageType {
	case MessageTypeImage, MessageTypeAudio, MessageTypeVideo, MessageTypeDocument:
	default:
		return nil, NewValidationError(fmt.Sprintf("invalid media message type: %s", messageType))
	}

	formattedJID := s.FormatJID(chatJID)
	msg := NewMessage(sessionID, formattedJID, "", formattedJID, messageType, "")
	msg.SetMediaURL(mediaURL)

	return msg, nil
}

func (s *DomainService) FormatJID(jid string) string {
	jid = strings.TrimSpace(jid)
	
	if strings.Contains(jid, "@") {
		return jid
	}

	return jid + "@s.whatsapp.net"
}
