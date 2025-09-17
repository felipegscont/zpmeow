package wmeow

import (
	"context"
	"fmt"

	"meow/internal/application"
	"meow/internal/infrastructure/logging"

	"go.mau.fi/whatsmeow/types/events"
)

// EventProcessorRefactored - Focused only on event processing
// Webhook dispatching is delegated to the Application Layer
type EventProcessorRefactored struct {
	sessionID       string
	sessionName     string
	eventDispatcher *application.EventDispatcher
	logger          logging.Logger
}

// NewEventProcessorRefactored creates a new refactored event processor
func NewEventProcessorRefactored(sessionID, sessionName string, eventDispatcher *application.EventDispatcher) *EventProcessorRefactored {
	return &EventProcessorRefactored{
		sessionID:       sessionID,
		sessionName:     sessionName,
		eventDispatcher: eventDispatcher,
		logger:          logging.GetLogger().Sub("events").Sub(sessionID),
	}
}

// HandleEvent processes events and delegates webhook dispatching
func (ep *EventProcessorRefactored) HandleEvent(evt interface{}) {
	eventType := fmt.Sprintf("%T", evt)

	// Log events in a cleaner format
	ep.logger.Debugf("Event: %s", eventType)
	ep.logger.Debugf("Data:\n%s", ep.formatEventData(evt))

	// Delegate webhook dispatching to Application Layer
	if ep.eventDispatcher != nil {
		ctx := context.Background()
		err := ep.eventDispatcher.DispatchEvent(ctx, ep.sessionID, eventType, evt)
		if err != nil {
			ep.logger.Errorf("Failed to dispatch event: %v", err)
		}
	}

	// Process event with specific handlers
	if handler, exists := refactoredEventHandlers[eventType]; exists {
		handler(ep, evt)
	} else {
		ep.logger.Debugf("Unhandled: %s", eventType)
	}
}

// formatEventData formats event data for logging (Infrastructure concern)
func (ep *EventProcessorRefactored) formatEventData(evt interface{}) string {
	// Implementation for formatting event data
	// This is purely infrastructure/logging concern
	return fmt.Sprintf("%+v", evt)
}

// Event handlers - focused only on logging and infrastructure concerns
func (ep *EventProcessorRefactored) handleMessage(evt interface{}) {
	msg := evt.(*events.Message)

	// Extract message info for logging
	messageID := msg.Info.ID
	messageType := getMessageType(msg)
	sender := msg.Info.Sender.String()
	timestamp := msg.Info.Timestamp.Format("15:04:05")
	sessionInfo := fmt.Sprintf("Session:%-10s", ep.sessionName)

	// Log message info
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

	// Extract receipt info for logging
	messageCount := len(receipt.MessageIDs)
	receiptType := string(receipt.Type)
	sender := receipt.MessageSource.Sender.String()
	timestamp := receipt.Timestamp.Format("15:04:05")
	sessionInfo := fmt.Sprintf("Session:%-10s", ep.sessionName)

	// Log receipt info
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

// Helper functions for logging
func truncateStringRefactored(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Map of refactored handlers
var refactoredEventHandlers = map[string]func(*EventProcessorRefactored, interface{}){
	// Message events
	"*events.Message": (*EventProcessorRefactored).handleMessage,
	"*events.Receipt": (*EventProcessorRefactored).handleReceipt,

	// Connection events
	"*events.Connected":    (*EventProcessorRefactored).handleConnected,
	"*events.Disconnected": (*EventProcessorRefactored).handleDisconnected,
	"*events.LoggedOut":    (*EventProcessorRefactored).handleLoggedOut,

	// Authentication events
	"*events.QR":          (*EventProcessorRefactored).handleQR,
	"*events.PairSuccess": (*EventProcessorRefactored).handlePairSuccess,
	"*events.PairError":   (*EventProcessorRefactored).handlePairError,

	// Presence events
	"*events.Presence":     (*EventProcessorRefactored).handlePresence,
	"*events.ChatPresence": (*EventProcessorRefactored).handleChatPresence,

	// Privacy events
	"*events.PrivacySettings": (*EventProcessorRefactored).handlePrivacySettings,
	"*events.Blocklist":       (*EventProcessorRefactored).handleBlocklist,
}
