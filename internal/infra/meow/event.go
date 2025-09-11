package meow

import (
	"context"
	"fmt"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/infra/webhook"

	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)


type EventHandler struct {
	sessionID       string
	logger          logger.Logger
	waLogger        waLog.Logger
	client          *MeowClient
	webhookService  *webhook.WebhookService
	sessionService  session.SessionService


	messageCount    int64
	lastMessageTime time.Time
}


func NewEventHandler(sessionID string, waLogger waLog.Logger, client *MeowClient, webhookService *webhook.WebhookService, sessionService session.SessionService) *EventHandler {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("event-handler").Sub(sessionID)

	return &EventHandler{
		sessionID:      sessionID,
		logger:         appLogger,
		waLogger:       waLogger,
		client:         client,
		webhookService: webhookService,
		sessionService: sessionService,
	}
}


// sendWebhook sends an event to the configured webhook URL if the event is subscribed
func (eh *EventHandler) sendWebhook(eventType string, data interface{}) {
	if eh.webhookService == nil || eh.sessionService == nil {
		return
	}

	// Get session configuration
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sess, err := eh.sessionService.GetSession(ctx, eh.sessionID)
	if err != nil {
		eh.logger.Debugf("Failed to get session for webhook: %v", err)
		return
	}

	// Check if webhook is configured and event is subscribed
	if !sess.HasWebhook() || !sess.IsEventSubscribed(eventType) {
		return
	}

	// Send webhook asynchronously
	eh.webhookService.SendWebhookAsync(sess.WebhookURL, eventType, eh.sessionID, data)
}

func (eh *EventHandler) HandleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		eh.handleMessage(v)
		eh.sendWebhook("message", v)
	case *events.Receipt:
		eh.handleReceipt(v)
		eh.sendWebhook("receipt", v)
	case *events.Presence:
		eh.handlePresence(v)
		eh.sendWebhook("presence", v)
	case *events.ChatPresence:
		eh.handleChatPresence(v)
		eh.sendWebhook("chat_presence", v)
	case *events.Connected:
		eh.handleConnected(v)
		eh.sendWebhook("connected", v)
	case *events.Disconnected:
		eh.handleDisconnected(v)
		eh.sendWebhook("disconnected", v)
	case *events.LoggedOut:
		eh.handleLoggedOut(v)
		eh.sendWebhook("logged_out", v)
	case *events.QR:
		eh.handleQR(v)
		eh.sendWebhook("qr", v)
	case *events.PairSuccess:
		eh.handlePairSuccess(v)
		eh.sendWebhook("pair_success", v)
	case *events.ConnectFailure:
		eh.handleConnectFailure(v)
		eh.sendWebhook("connect_failure", v)
	case *events.StreamError:
		eh.handleStreamError(v)
		eh.sendWebhook("stream_error", v)
	case *events.StreamReplaced:
		eh.handleStreamReplaced(v)
		eh.sendWebhook("stream_replaced", v)
	case *events.TemporaryBan:
		eh.handleTemporaryBan(v)
		eh.sendWebhook("temporary_ban", v)
	case *events.GroupInfo:
		eh.handleGroupInfo(v)
		eh.sendWebhook("group_info", v)
	case *events.JoinedGroup:
		eh.handleJoinedGroup(v)
		eh.sendWebhook("joined_group", v)
	case *events.Contact:
		eh.handleContact(v)
		eh.sendWebhook("contact", v)
	case *events.PushName:
		eh.handlePushName(v)
		eh.sendWebhook("push_name", v)
	case *events.BusinessName:
		eh.handleBusinessName(v)
		eh.sendWebhook("business_name", v)
	case *events.IdentityChange:
		eh.handleIdentityChange(v)
		eh.sendWebhook("identity_change", v)
	case *events.PrivacySettings:
		eh.handlePrivacySettings(v)
		eh.sendWebhook("privacy_settings", v)
	case *events.OfflineSyncPreview:
		eh.handleOfflineSyncPreview(v)
		eh.sendWebhook("offline_sync_preview", v)
	case *events.OfflineSyncCompleted:
		eh.handleOfflineSyncCompleted(v)
		eh.sendWebhook("offline_sync_completed", v)
	case *events.AppStateSyncComplete:
		eh.handleAppStateSyncComplete(v)
		eh.sendWebhook("app_state_sync_complete", v)
	case *events.HistorySync:
		eh.handleHistorySync(v)
		eh.sendWebhook("history_sync", v)
	case *events.AppState:
		eh.handleAppState(v)
		eh.sendWebhook("app_state", v)
	case *events.KeepAliveTimeout:
		eh.handleKeepAliveTimeout(v)
		eh.sendWebhook("keep_alive_timeout", v)
	case *events.KeepAliveRestored:
		eh.handleKeepAliveRestored(v)
		eh.sendWebhook("keep_alive_restored", v)
	case *events.Blocklist:
		eh.handleBlocklist(v)
		eh.sendWebhook("blocklist", v)
	default:
		eh.logger.Debugf("Session %s: Unhandled event type: %T", eh.sessionID, evt)
	}
}


func (eh *EventHandler) handleMessage(evt *events.Message) {
	
	go func() {
		defer func() {
			if r := recover(); r != nil {
				eh.logger.Errorf("Panic in message handler for session %s: %v", eh.sessionID, r)
			}
		}()

		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		eh.processMessageWithTimeout(ctx, evt)
	}()
}


func (eh *EventHandler) processMessageWithTimeout(ctx context.Context, evt *events.Message) {
	start := time.Now()

	select {
	case <-ctx.Done():
		eh.logger.Warnf("Session %s: Message processing timeout for message %s from %s",
			eh.sessionID, evt.Info.ID, evt.Info.Sender)
		return
	default:
		
		eh.messageCount++
		eh.lastMessageTime = time.Now()

		
		eh.logger.Infof("Session %s: Received message from %s (ID: %s): %s",
			eh.sessionID, evt.Info.Sender, evt.Info.ID, evt.Message.GetConversation())

		
		if eh.client != nil {
			eh.client.updateActivity()
		}

		
		

		duration := time.Since(start)
		eh.logger.Debugf("Session %s: Message processed successfully: %s (took %v)",
			eh.sessionID, evt.Info.ID, duration)

		
		if duration > 5*time.Second {
			eh.logger.Warnf("Session %s: Slow message processing: %s took %v",
				eh.sessionID, evt.Info.ID, duration)
		}
	}
}


func (eh *EventHandler) handleReceipt(evt *events.Receipt) {
	
	go func() {
		defer func() {
			if r := recover(); r != nil {
				eh.logger.Errorf("Panic in receipt handler for session %s: %v", eh.sessionID, r)
			}
		}()

		eh.logger.Debugf("Session %s: Received receipt for %s from %s",
			eh.sessionID, evt.MessageIDs, evt.SourceString())
	}()
}


func (eh *EventHandler) handlePresence(evt *events.Presence) {
	eh.logger.Debugf("Session %s: Presence update from %s: unavailable=%t",
		eh.sessionID, evt.From, evt.Unavailable)
}


func (eh *EventHandler) handleChatPresence(evt *events.ChatPresence) {
	eh.logger.Debugf("Session %s: Chat presence from %s in %s: %s",
		eh.sessionID, evt.Sender, evt.Chat, evt.State)
}


func (eh *EventHandler) handleConnected(_ *events.Connected) {
	eh.logger.Infof("Session %s: Connected to WhatsApp", eh.sessionID)

	if eh.client != nil {
		eh.client.onConnected()

		
		if eh.client.manager != nil && eh.client.client != nil && eh.client.client.Store.ID != nil {
			deviceJID := eh.client.client.Store.ID.String()
			eh.logger.Infof("Session %s: Connected event - updating device JID to %s", eh.sessionID, deviceJID)
			eh.client.manager.OnPairSuccess(eh.sessionID, deviceJID)
		}
	}
}


func (eh *EventHandler) handleDisconnected(_ *events.Disconnected) {
	eh.logger.Infof("Session %s: Disconnected from WhatsApp", eh.sessionID)

	if eh.client != nil {
		eh.client.onDisconnected()
	}
}


func (eh *EventHandler) handleLoggedOut(evt *events.LoggedOut) {
	eh.logger.Warnf("Session %s: Logged out from WhatsApp: %s", eh.sessionID, evt.Reason)

	if eh.client != nil {
		eh.client.onDisconnected()
	}
}


func (eh *EventHandler) handleQR(evt *events.QR) {
	eh.logger.Infof("Session %s: QR codes received (%d codes)", eh.sessionID, len(evt.Codes))

	
	
	if eh.client != nil && len(evt.Codes) > 0 {
		
		eh.logger.Debugf("Session %s: QR code event processed by event handler", eh.sessionID)
	}
}


func (eh *EventHandler) handlePairSuccess(evt *events.PairSuccess) {
	eh.logger.Infof("Session %s: Pairing successful with %s", eh.sessionID, evt.ID)

	if eh.client != nil {
		eh.client.onConnected()

		
		if eh.client.manager != nil {
			eh.logger.Infof("Session %s: Calling OnPairSuccess with device JID %s", eh.sessionID, evt.ID.String())
			eh.client.manager.OnPairSuccess(eh.sessionID, evt.ID.String())
		} else {
			eh.logger.Warnf("Session %s: Manager is nil, cannot update device JID", eh.sessionID)
		}
	}
}


func (eh *EventHandler) handleConnectFailure(evt *events.ConnectFailure) {
	eh.logger.Errorf("Session %s: Connection failed: %s", eh.sessionID, evt.Reason)

	if eh.client != nil {
		eh.client.onError(fmt.Errorf("connection failed: %s", evt.Reason))
	}
}


func (eh *EventHandler) handleStreamError(evt *events.StreamError) {
	eh.logger.Errorf("Session %s: Stream error: %s", eh.sessionID, evt.Code)

	if eh.client != nil {
		eh.client.onError(fmt.Errorf("stream error: %s", evt.Code))
	}
}


func (eh *EventHandler) handleStreamReplaced(_ *events.StreamReplaced) {
	eh.logger.Warnf("Session %s: Stream replaced", eh.sessionID)

	if eh.client != nil {
		eh.client.onDisconnected()
	}
}


func (eh *EventHandler) handleTemporaryBan(evt *events.TemporaryBan) {
	eh.logger.Warnf("Session %s: Temporary ban: %s (expires in %s)",
		eh.sessionID, evt.Code, evt.Expire)

	if eh.client != nil {
		eh.client.onError(fmt.Errorf("temporary ban: %s", evt.Code))
	}
}


func (eh *EventHandler) handleGroupInfo(evt *events.GroupInfo) {
	eh.logger.Debugf("Session %s: Group info update for %s", eh.sessionID, evt.JID)
}


func (eh *EventHandler) handleJoinedGroup(evt *events.JoinedGroup) {
	eh.logger.Infof("Session %s: Joined group %s", eh.sessionID, evt.JID)
}


func (eh *EventHandler) handleContact(evt *events.Contact) {
	eh.logger.Debugf("Session %s: Contact update for %s", eh.sessionID, evt.JID)
}


func (eh *EventHandler) handlePushName(evt *events.PushName) {
	eh.logger.Debugf("Session %s: Push name update for %s: %s",
		eh.sessionID, evt.JID, evt.Message.PushName)
}


func (eh *EventHandler) handleBusinessName(evt *events.BusinessName) {
	eh.logger.Debugf("Session %s: Business name update for %s: %s -> %s",
		eh.sessionID, evt.JID, evt.OldBusinessName, evt.NewBusinessName)
}


func (eh *EventHandler) handleIdentityChange(evt *events.IdentityChange) {
	eh.logger.Warnf("Session %s: Identity change for %s", eh.sessionID, evt.JID)
}


func (eh *EventHandler) handlePrivacySettings(_ *events.PrivacySettings) {
	eh.logger.Debugf("Session %s: Privacy settings updated", eh.sessionID)
}


func (eh *EventHandler) handleOfflineSyncPreview(evt *events.OfflineSyncPreview) {
	eh.logger.Infof("Session %s: Offline sync preview: %d total, %d messages",
		eh.sessionID, evt.Total, evt.Messages)
}


func (eh *EventHandler) handleOfflineSyncCompleted(_ *events.OfflineSyncCompleted) {
	eh.logger.Infof("Session %s: Offline sync completed", eh.sessionID)
}


func (eh *EventHandler) handleAppStateSyncComplete(evt *events.AppStateSyncComplete) {
	eh.logger.Debugf("Session %s: App state sync complete for %s", eh.sessionID, evt.Name)
}


func (eh *EventHandler) handleHistorySync(evt *events.HistorySync) {
	
	go func() {
		defer func() {
			if r := recover(); r != nil {
				eh.logger.Errorf("Panic in history sync handler for session %s: %v", eh.sessionID, r)
			}
		}()

		eh.logger.Infof("Session %s: History sync: %d conversations",
			eh.sessionID, len(evt.Data.Conversations))

		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			eh.logger.Warnf("Session %s: History sync processing timeout", eh.sessionID)
		default:
			
			eh.logger.Debugf("Session %s: History sync processed successfully", eh.sessionID)
		}
	}()
}


func (eh *EventHandler) handleAppState(evt *events.AppState) {
	eh.logger.Debugf("Session %s: App state update with %d indices", eh.sessionID, len(evt.Index))
}


func (eh *EventHandler) handleKeepAliveTimeout(_ *events.KeepAliveTimeout) {
	eh.logger.Warnf("Session %s: Keep alive timeout", eh.sessionID)
}


func (eh *EventHandler) handleKeepAliveRestored(_ *events.KeepAliveRestored) {
	eh.logger.Infof("Session %s: Keep alive restored", eh.sessionID)
}


func (eh *EventHandler) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"session_id":        eh.sessionID,
		"message_count":     eh.messageCount,
		"last_message_time": eh.lastMessageTime,
		"uptime":            time.Since(eh.lastMessageTime),
	}
}


func (eh *EventHandler) handleBlocklist(evt *events.Blocklist) {
	eh.logger.Debugf("Session %s: Blocklist updated: %d changes",
		eh.sessionID, len(evt.Changes))
}
