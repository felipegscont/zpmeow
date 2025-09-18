package wmeow

import (
	"context"
	"fmt"

	"meow/internal/application"
	"meow/internal/infrastructure/logging"

	"go.mau.fi/whatsmeow/types/events"
)

type EventProcessorRefactored struct {
	sessionID       string
	sessionName     string
	eventDispatcher *application.EventDispatcher
	logger          logging.Logger
}

func NewEventProcessorRefactored(sessionID, sessionName string, eventDispatcher *application.EventDispatcher) *EventProcessorRefactored {
	return &EventProcessorRefactored{
		sessionID:       sessionID,
		sessionName:     sessionName,
		eventDispatcher: eventDispatcher,
		logger:          logging.GetLogger().Sub("events").Sub(sessionID),
	}
}

func (ep *EventProcessorRefactored) HandleEvent(evt interface{}) {
	eventType := fmt.Sprintf("%T", evt)

	ep.logger.Debugf("Event: %s", eventType)
	ep.logger.Debugf("Data:\n%s", ep.formatEventData(evt))

	if ep.eventDispatcher != nil {
		ctx := context.Background()
		err := ep.eventDispatcher.DispatchEvent(ctx, ep.sessionID, eventType, evt)
		if err != nil {
			ep.logger.Errorf("Failed to dispatch event: %v", err)
		}
	}

	if handler, exists := refactoredEventHandlers[eventType]; exists {
		handler(ep, evt)
	} else {
		ep.logger.Debugf("Unhandled: %s", eventType)
	}
}

func (ep *EventProcessorRefactored) formatEventData(evt interface{}) string {
	return fmt.Sprintf("%+v", evt)
}

func (ep *EventProcessorRefactored) handleMessage(evt interface{}) {
	msg := evt.(*events.Message)

	messageID := msg.Info.ID
	messageType := getMessageType(msg)
	sender := msg.Info.Sender.String()
	timestamp := msg.Info.Timestamp.Format("15:04:05")
	sessionInfo := fmt.Sprintf("Session:%-10s", ep.sessionName)

	if msg.Info.IsGroup {
		groupName := truncateStringRefactored(msg.Info.Chat.String(), 25)
		ep.logger.Infof("Message received [%-8s] ID:%-20s Type:%-8s From:%-25s Group:%-25s %s Time:%s",
			"GROUP", messageID, messageType, truncateStringRefactored(sender, 25), groupName, sessionInfo, timestamp)
	} else {
		ep.logger.Infof("Message received [%-8s] ID:%-20s Type:%-8s From:%-25s %s Time:%s",
			"DIRECT", messageID, messageType, truncateStringRefactored(sender, 25), sessionInfo, timestamp)
	}
}

func (ep *EventProcessorRefactored) handleConnected(evt interface{}) {
	ep.logger.Infof("Session %s connected", ep.sessionID)
}

func (ep *EventProcessorRefactored) handleDisconnected(evt interface{}) {
	ep.logger.Infof("Session %s disconnected", ep.sessionID)
}

func (ep *EventProcessorRefactored) handleLoggedOut(evt interface{}) {
	ep.logger.Infof("Session %s logged out", ep.sessionID)
}

func (ep *EventProcessorRefactored) handleQR(evt interface{}) {
	ep.logger.Debugf("QR code event for session %s", ep.sessionID)
}

func (ep *EventProcessorRefactored) handlePairSuccess(evt interface{}) {
	ep.logger.Infof("Pairing successful for session %s", ep.sessionID)
}

func (ep *EventProcessorRefactored) handlePairError(evt interface{}) {
	if pairErr, ok := evt.(*events.PairError); ok {
		ep.logger.Errorf("Pairing error for session %s: %v", ep.sessionID, pairErr.Error)
	}
}

func (ep *EventProcessorRefactored) handleReceipt(evt interface{}) {
	receipt := evt.(*events.Receipt)

	messageCount := len(receipt.MessageIDs)
	receiptType := string(receipt.Type)
	sender := receipt.MessageSource.Sender.String()
	timestamp := receipt.Timestamp.Format("15:04:05")
	sessionInfo := fmt.Sprintf("Session:%-10s", ep.sessionName)

	if receipt.MessageSource.IsGroup {
		groupName := truncateStringRefactored(receipt.MessageSource.Chat.String(), 25)
		ep.logger.Infof("Receipt received [%-8s] Count:%-3d        Type:%-8s From:%-25s Group:%-25s %s Time:%s",
			"GROUP", messageCount, receiptType, truncateStringRefactored(sender, 25), groupName, sessionInfo, timestamp)
	} else {
		ep.logger.Infof("Receipt received [%-8s] Count:%-3d        Type:%-8s From:%-25s %s Time:%s",
			"RECEIPT", messageCount, receiptType, truncateStringRefactored(sender, 25), sessionInfo, timestamp)
	}
}

func (ep *EventProcessorRefactored) handlePresence(evt interface{}) {
	ep.logger.Debugf("Presence update for session %s", ep.sessionID)
}

func (ep *EventProcessorRefactored) handleChatPresence(evt interface{}) {
	ep.logger.Debugf("Chat presence update for session %s", ep.sessionID)
}

func (ep *EventProcessorRefactored) handlePrivacySettings(evt interface{}) {
	ep.logger.Debugf("Privacy settings changed for session %s", ep.sessionID)
}

func (ep *EventProcessorRefactored) handleBlocklist(evt interface{}) {
	ep.logger.Debugf("Blocklist changed for session %s", ep.sessionID)
}

func truncateStringRefactored(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

var refactoredEventHandlers = map[string]func(*EventProcessorRefactored, interface{}){
	"*events.Message": (*EventProcessorRefactored).handleMessage,
	"*events.Receipt": (*EventProcessorRefactored).handleReceipt,

	"*events.Connected":    (*EventProcessorRefactored).handleConnected,
	"*events.Disconnected": (*EventProcessorRefactored).handleDisconnected,
	"*events.LoggedOut":    (*EventProcessorRefactored).handleLoggedOut,

	"*events.QR":          (*EventProcessorRefactored).handleQR,
	"*events.PairSuccess": (*EventProcessorRefactored).handlePairSuccess,
	"*events.PairError":   (*EventProcessorRefactored).handlePairError,

	"*events.Presence":     (*EventProcessorRefactored).handlePresence,
	"*events.ChatPresence": (*EventProcessorRefactored).handleChatPresence,

	"*events.PrivacySettings": (*EventProcessorRefactored).handlePrivacySettings,
	"*events.Blocklist":       (*EventProcessorRefactored).handleBlocklist,
}
