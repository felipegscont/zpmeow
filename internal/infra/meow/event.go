package meow

import (
	"fmt"

	"zpmeow/internal/infra/logger"

	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// EventHandler handles WhatsApp events for a specific session
type EventHandler struct {
	sessionID string
	logger    logger.Logger
	waLogger  waLog.Logger
	client    *MeowClient
}

// NewEventHandler creates a new event handler
func NewEventHandler(sessionID string, waLogger waLog.Logger, client *MeowClient) *EventHandler {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("event-handler").Sub(sessionID)

	return &EventHandler{
		sessionID: sessionID,
		logger:    appLogger,
		waLogger:  waLogger,
		client:    client,
	}
}

// HandleEvent is the main event handler function
func (eh *EventHandler) HandleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		eh.handleMessage(v)
	case *events.Receipt:
		eh.handleReceipt(v)
	case *events.Presence:
		eh.handlePresence(v)
	case *events.ChatPresence:
		eh.handleChatPresence(v)
	case *events.Connected:
		eh.handleConnected(v)
	case *events.Disconnected:
		eh.handleDisconnected(v)
	case *events.LoggedOut:
		eh.handleLoggedOut(v)
	case *events.QR:
		eh.handleQR(v)
	case *events.PairSuccess:
		eh.handlePairSuccess(v)
	case *events.ConnectFailure:
		eh.handleConnectFailure(v)
	case *events.StreamError:
		eh.handleStreamError(v)
	case *events.StreamReplaced:
		eh.handleStreamReplaced(v)
	case *events.TemporaryBan:
		eh.handleTemporaryBan(v)
	case *events.GroupInfo:
		eh.handleGroupInfo(v)
	case *events.JoinedGroup:
		eh.handleJoinedGroup(v)
	case *events.Contact:
		eh.handleContact(v)
	case *events.PushName:
		eh.handlePushName(v)
	case *events.BusinessName:
		eh.handleBusinessName(v)
	case *events.IdentityChange:
		eh.handleIdentityChange(v)
	case *events.PrivacySettings:
		eh.handlePrivacySettings(v)
	case *events.OfflineSyncPreview:
		eh.handleOfflineSyncPreview(v)
	case *events.OfflineSyncCompleted:
		eh.handleOfflineSyncCompleted(v)
	case *events.AppStateSyncComplete:
		eh.handleAppStateSyncComplete(v)
	case *events.HistorySync:
		eh.handleHistorySync(v)
	case *events.AppState:
		eh.handleAppState(v)
	case *events.KeepAliveTimeout:
		eh.handleKeepAliveTimeout(v)
	case *events.KeepAliveRestored:
		eh.handleKeepAliveRestored(v)
	case *events.Blocklist:
		eh.handleBlocklist(v)
	default:
		eh.logger.Debugf("Session %s: Unhandled event type: %T", eh.sessionID, evt)
	}
}

// handleMessage handles incoming messages
func (eh *EventHandler) handleMessage(evt *events.Message) {
	eh.logger.Infof("Session %s: Received message from %s: %s", 
		eh.sessionID, evt.Info.Sender, evt.Message.GetConversation())
	
	// Update client activity
	if eh.client != nil {
		eh.client.updateActivity()
	}

	// Here you can add custom message processing logic
	// For example: save to database, forward to webhooks, etc.
}

// handleReceipt handles message receipts
func (eh *EventHandler) handleReceipt(evt *events.Receipt) {
	eh.logger.Debugf("Session %s: Received receipt for %s from %s", 
		eh.sessionID, evt.MessageIDs, evt.SourceString())
}

// handlePresence handles presence updates
func (eh *EventHandler) handlePresence(evt *events.Presence) {
	eh.logger.Debugf("Session %s: Presence update from %s: unavailable=%t", 
		eh.sessionID, evt.From, evt.Unavailable)
}

// handleChatPresence handles chat presence (typing indicators)
func (eh *EventHandler) handleChatPresence(evt *events.ChatPresence) {
	eh.logger.Debugf("Session %s: Chat presence from %s in %s: %s", 
		eh.sessionID, evt.Sender, evt.Chat, evt.State)
}

// handleConnected handles successful connection
func (eh *EventHandler) handleConnected(evt *events.Connected) {
	eh.logger.Infof("Session %s: Connected to WhatsApp", eh.sessionID)

	if eh.client != nil {
		eh.client.onConnected()

		// Also update device JID on connection if available and not already set
		if eh.client.manager != nil && eh.client.client != nil && eh.client.client.Store.ID != nil {
			deviceJID := eh.client.client.Store.ID.String()
			eh.logger.Infof("Session %s: Connected event - updating device JID to %s", eh.sessionID, deviceJID)
			eh.client.manager.OnPairSuccess(eh.sessionID, deviceJID)
		}
	}
}

// handleDisconnected handles disconnection
func (eh *EventHandler) handleDisconnected(evt *events.Disconnected) {
	eh.logger.Infof("Session %s: Disconnected from WhatsApp", eh.sessionID)
	
	if eh.client != nil {
		eh.client.onDisconnected()
	}
}

// handleLoggedOut handles logout events
func (eh *EventHandler) handleLoggedOut(evt *events.LoggedOut) {
	eh.logger.Warnf("Session %s: Logged out from WhatsApp: %s", eh.sessionID, evt.Reason)
	
	if eh.client != nil {
		eh.client.onDisconnected()
	}
}

// handleQR handles QR code events
func (eh *EventHandler) handleQR(evt *events.QR) {
	eh.logger.Infof("Session %s: QR codes received (%d codes)", eh.sessionID, len(evt.Codes))

	// Note: QR code display and storage is handled by the QR channel loop in client.go
	// This avoids duplication of QR code printing
	if eh.client != nil && len(evt.Codes) > 0 {
		// Just log the event, the QR channel loop will handle display and storage
		eh.logger.Debugf("Session %s: QR code event processed by event handler", eh.sessionID)
	}
}

// handlePairSuccess handles successful pairing
func (eh *EventHandler) handlePairSuccess(evt *events.PairSuccess) {
	eh.logger.Infof("Session %s: Pairing successful with %s", eh.sessionID, evt.ID)

	if eh.client != nil {
		eh.client.onConnected()

		// Notify manager about successful pairing with device JID
		if eh.client.manager != nil {
			eh.logger.Infof("Session %s: Calling OnPairSuccess with device JID %s", eh.sessionID, evt.ID.String())
			eh.client.manager.OnPairSuccess(eh.sessionID, evt.ID.String())
		} else {
			eh.logger.Warnf("Session %s: Manager is nil, cannot update device JID", eh.sessionID)
		}
	}
}

// handleConnectFailure handles connection failures
func (eh *EventHandler) handleConnectFailure(evt *events.ConnectFailure) {
	eh.logger.Errorf("Session %s: Connection failed: %s", eh.sessionID, evt.Reason)
	
	if eh.client != nil {
		eh.client.onError(fmt.Errorf("connection failed: %s", evt.Reason))
	}
}

// handleStreamError handles stream errors
func (eh *EventHandler) handleStreamError(evt *events.StreamError) {
	eh.logger.Errorf("Session %s: Stream error: %s", eh.sessionID, evt.Code)
	
	if eh.client != nil {
		eh.client.onError(fmt.Errorf("stream error: %s", evt.Code))
	}
}

// handleStreamReplaced handles stream replacement
func (eh *EventHandler) handleStreamReplaced(evt *events.StreamReplaced) {
	eh.logger.Warnf("Session %s: Stream replaced", eh.sessionID)
	
	if eh.client != nil {
		eh.client.onDisconnected()
	}
}

// handleTemporaryBan handles temporary bans
func (eh *EventHandler) handleTemporaryBan(evt *events.TemporaryBan) {
	eh.logger.Warnf("Session %s: Temporary ban: %s (expires in %s)", 
		eh.sessionID, evt.Code, evt.Expire)
	
	if eh.client != nil {
		eh.client.onError(fmt.Errorf("temporary ban: %s", evt.Code))
	}
}

// handleGroupInfo handles group information updates
func (eh *EventHandler) handleGroupInfo(evt *events.GroupInfo) {
	eh.logger.Debugf("Session %s: Group info update for %s", eh.sessionID, evt.JID)
}

// handleJoinedGroup handles group join events
func (eh *EventHandler) handleJoinedGroup(evt *events.JoinedGroup) {
	eh.logger.Infof("Session %s: Joined group %s", eh.sessionID, evt.JID)
}

// handleContact handles contact updates
func (eh *EventHandler) handleContact(evt *events.Contact) {
	eh.logger.Debugf("Session %s: Contact update for %s", eh.sessionID, evt.JID)
}

// handlePushName handles push name updates
func (eh *EventHandler) handlePushName(evt *events.PushName) {
	eh.logger.Debugf("Session %s: Push name update for %s: %s", 
		eh.sessionID, evt.JID, evt.Message.PushName)
}

// handleBusinessName handles business name updates
func (eh *EventHandler) handleBusinessName(evt *events.BusinessName) {
	eh.logger.Debugf("Session %s: Business name update for %s: %s -> %s",
		eh.sessionID, evt.JID, evt.OldBusinessName, evt.NewBusinessName)
}

// handleIdentityChange handles identity changes
func (eh *EventHandler) handleIdentityChange(evt *events.IdentityChange) {
	eh.logger.Warnf("Session %s: Identity change for %s", eh.sessionID, evt.JID)
}

// handlePrivacySettings handles privacy settings updates
func (eh *EventHandler) handlePrivacySettings(evt *events.PrivacySettings) {
	eh.logger.Debugf("Session %s: Privacy settings updated", eh.sessionID)
}

// handleOfflineSyncPreview handles offline sync preview
func (eh *EventHandler) handleOfflineSyncPreview(evt *events.OfflineSyncPreview) {
	eh.logger.Infof("Session %s: Offline sync preview: %d total, %d messages",
		eh.sessionID, evt.Total, evt.Messages)
}

// handleOfflineSyncCompleted handles offline sync completion
func (eh *EventHandler) handleOfflineSyncCompleted(evt *events.OfflineSyncCompleted) {
	eh.logger.Infof("Session %s: Offline sync completed", eh.sessionID)
}

// handleAppStateSyncComplete handles app state sync completion
func (eh *EventHandler) handleAppStateSyncComplete(evt *events.AppStateSyncComplete) {
	eh.logger.Debugf("Session %s: App state sync complete for %s", eh.sessionID, evt.Name)
}

// handleHistorySync handles history sync events
func (eh *EventHandler) handleHistorySync(evt *events.HistorySync) {
	eh.logger.Infof("Session %s: History sync: %d conversations", 
		eh.sessionID, len(evt.Data.Conversations))
}

// handleAppState handles app state events
func (eh *EventHandler) handleAppState(evt *events.AppState) {
	eh.logger.Debugf("Session %s: App state update with %d indices", eh.sessionID, len(evt.Index))
}

// handleKeepAliveTimeout handles keep alive timeouts
func (eh *EventHandler) handleKeepAliveTimeout(evt *events.KeepAliveTimeout) {
	eh.logger.Warnf("Session %s: Keep alive timeout", eh.sessionID)
}

// handleKeepAliveRestored handles keep alive restoration
func (eh *EventHandler) handleKeepAliveRestored(evt *events.KeepAliveRestored) {
	eh.logger.Infof("Session %s: Keep alive restored", eh.sessionID)
}

// handleBlocklist handles blocklist updates
func (eh *EventHandler) handleBlocklist(evt *events.Blocklist) {
	eh.logger.Debugf("Session %s: Blocklist updated: %d changes",
		eh.sessionID, len(evt.Changes))
}
