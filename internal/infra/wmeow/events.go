package wmeow

import (
	"fmt"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logging"

	"go.mau.fi/whatsmeow/types/events"
)

// EventProcessor - Simplified event handling with map-based handlers
type EventProcessor struct {
	sessionID   string
	webhookURL  string
	sessionRepo session.Repository
	logger      logging.Logger
}

// Map of handlers organized by event category
var eventHandlers = map[string]func(*EventProcessor, interface{}){
	// Message events
	"*events.Message": (*EventProcessor).handleMessage,
	"*events.Receipt": (*EventProcessor).handleReceipt,

	// Connection events
	"*events.Connected":    (*EventProcessor).handleConnected,
	"*events.Disconnected": (*EventProcessor).handleDisconnected,
	"*events.LoggedOut":    (*EventProcessor).handleLoggedOut,

	// Authentication events
	"*events.QR":          (*EventProcessor).handleQR,
	"*events.PairSuccess": (*EventProcessor).handlePairSuccess,
	"*events.PairError":   (*EventProcessor).handlePairError,

	// Presence events
	"*events.Presence":     (*EventProcessor).handlePresence,
	"*events.ChatPresence": (*EventProcessor).handleChatPresence,
}

// NewEventProcessor creates a new simplified event processor
func NewEventProcessor(sessionID, webhookURL string, sessionRepo session.Repository) *EventProcessor {
	return &EventProcessor{
		sessionID:   sessionID,
		webhookURL:  webhookURL,
		sessionRepo: sessionRepo,
		logger:      logging.GetLogger().Sub("events").Sub(sessionID),
	}
}

// HandleEvent processes events using map-based handlers
func (ep *EventProcessor) HandleEvent(evt interface{}) {
	eventType := fmt.Sprintf("%T", evt)

	// Log all events in DEBUG for monitoring
	ep.logger.Debugf("📨 Event received: %s", eventType)
	ep.logger.Debugf("📄 Event raw: %+v", evt)

	if handler, exists := eventHandlers[eventType]; exists {
		handler(ep, evt)
	} else {
		ep.logger.Debugf("❓ Unhandled event: %s", eventType)
	}
}

// Event category helpers for better organization (removed unused methods)

// processConnectionEvents - Common processing for connection events
func (ep *EventProcessor) processConnectionEvents(eventType, status string) {
	ep.logger.Infof("Session %s %s", ep.sessionID, status)
	data := map[string]interface{}{
		"sessionId": ep.sessionID,
		"status":    status,
		"event":     eventType,
	}
	sendWebhook(ep.webhookURL, data)
}

// processAuthEvents - Common processing for authentication events
func (ep *EventProcessor) processAuthEvents(eventType string, eventData interface{}) {
	data := map[string]interface{}{
		"sessionId": ep.sessionID,
		"event":     eventType,
	}

	// Add specific data based on event type
	switch eventType {
	case "qr":
		if qr, ok := eventData.(*events.QR); ok && len(qr.Codes) > 0 {
			data["qr"] = qr.Codes
			ep.logger.Debugf("QR codes available: %v", qr.Codes)
		}
	case "pair_success":
		data["status"] = "paired"
	case "pair_error":
		data["status"] = "pair_failed"
		if pairErr, ok := eventData.(*events.PairError); ok {
			data["error"] = pairErr.Error.Error()
		}
	}

	sendWebhook(ep.webhookURL, data)
}

// EventProcessor handler methods - small and focused
func (ep *EventProcessor) handleMessage(evt interface{}) {
	msg := evt.(*events.Message)
	ep.logger.Infof("Message received from %s in session %s", msg.Info.Sender, ep.sessionID)
	data := createMessageData(msg)
	sendWebhook(ep.webhookURL, data)
}

func (ep *EventProcessor) handleConnected(evt interface{}) {
	ep.processConnectionEvents("connected", "connected")
}

func (ep *EventProcessor) handleDisconnected(evt interface{}) {
	ep.processConnectionEvents("disconnected", "disconnected")
}

func (ep *EventProcessor) handleQR(evt interface{}) {
	qr := evt.(*events.QR)
	ep.logger.Infof("QR code generated for session %s", ep.sessionID)
	ep.processAuthEvents("qr", qr)
}

func (ep *EventProcessor) handlePairSuccess(evt interface{}) {
	ep.logger.Infof("Pair success for session %s", ep.sessionID)
	ep.processAuthEvents("pair_success", evt)
}

func (ep *EventProcessor) handlePairError(evt interface{}) {
	pairError := evt.(*events.PairError)
	ep.logger.Errorf("Pair error for session %s: %v", ep.sessionID, pairError.Error)
	ep.processAuthEvents("pair_error", pairError)
}

func (ep *EventProcessor) handleLoggedOut(evt interface{}) {
	ep.processConnectionEvents("logged_out", "logged_out")
}

func (ep *EventProcessor) handleReceipt(evt interface{}) {
	receipt := evt.(*events.Receipt)
	ep.logger.Debugf("Receipt received for session %s", ep.sessionID)
	data := map[string]interface{}{
		"sessionId": ep.sessionID,
		"event":     "receipt",
		"messageId": receipt.MessageIDs,
	}
	sendWebhook(ep.webhookURL, data)
}

func (ep *EventProcessor) handlePresence(evt interface{}) {
	presence := evt.(*events.Presence)
	ep.logger.Debugf("Presence update for session %s", ep.sessionID)

	status := "available"
	if presence.Unavailable {
		status = "unavailable"
	}

	data := map[string]interface{}{
		"sessionId": ep.sessionID,
		"event":     "presence",
		"from":      presence.From.String(),
		"presence":  status,
	}
	sendWebhook(ep.webhookURL, data)
}

func (ep *EventProcessor) handleChatPresence(evt interface{}) {
	chatPresence := evt.(*events.ChatPresence)
	ep.logger.Debugf("Chat presence update for session %s", ep.sessionID)
	data := map[string]interface{}{
		"sessionId": ep.sessionID,
		"event":     "chat_presence",
		"chat":      chatPresence.MessageSource.Chat.String(),
		"presence":  string(chatPresence.State),
	}
	sendWebhook(ep.webhookURL, data)
}

// Helper functions shared across handlers (DRY principle)
func sendWebhook(url string, _ interface{}) error {
	if url == "" {
		return nil // No webhook configured
	}

	// Simple webhook implementation
	// In a real implementation, you'd want proper HTTP client with retries
	return nil
}

func createMessageData(msg *events.Message) map[string]interface{} {
	return map[string]interface{}{
		"messageId": msg.Info.ID,
		"from":      msg.Info.Sender.String(),
		"chat":      msg.Info.Chat.String(),
		"timestamp": msg.Info.Timestamp.Unix(),
		"fromMe":    msg.Info.IsFromMe,
		"isGroup":   msg.Info.IsGroup,
	}
}
