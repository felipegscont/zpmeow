package common

import (
	"crypto/rand"
	"encoding/hex"
)

type MessageIDGenerator struct{}

func NewMessageIDGen() *MessageIDGenerator {
	return &MessageIDGenerator{}
}

func (g *MessageIDGenerator) Generate() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (g *MessageIDGenerator) CreateBasicSendResponse(to, text, session string) interface{} {
	return nil
}

type ResponseService struct {
	*MessageIDGenerator
}

func NewResponseService() *ResponseService {
	return &ResponseService{
		MessageIDGenerator: NewMessageIDGen(),
	}
}

func (r *ResponseService) CreateErrorResponse(message string, details ...string) map[string]interface{} {
	return map[string]interface{}{
		"status": "error",
		"error":  message,
	}
}

func (r *ResponseService) CreateSuccessResponse(message string, data ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	}
}

func (r *ResponseService) CreateWebhookResponse(webhookURL string, events []string, active bool) map[string]interface{} {
	return map[string]interface{}{
		"url":    webhookURL,
		"events": events,
		"active": active,
	}
}

func (r *ResponseService) CreateGroupResponse(groupJID, action string, participants []string) map[string]interface{} {
	return map[string]interface{}{
		"group_jid":    groupJID,
		"action":       action,
		"participants": participants,
	}
}

func (r *ResponseService) CreateChatResponse(action, chatJID, messageID string) map[string]interface{} {
	return map[string]interface{}{
		"action":     action,
		"chat_jid":   chatJID,
		"message_id": messageID,
	}
}

func (r *ResponseService) CreatePresenceResponse(phone, state, media string) map[string]interface{} {
	response := map[string]interface{}{
		"phone": phone,
		"state": state,
	}

	if media != "" {
		response["media"] = media
	}

	return response
}

var DefaultResponseService = NewResponseService()
