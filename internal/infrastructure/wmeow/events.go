package wmeow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"meow/internal/domain/session"
	"meow/internal/infrastructure/logging"
	"meow/internal/infrastructure/webhooks"

	waE2E "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// List of supported event types
var supportedEventTypes = []string{
	// Messages and Communication
	"Message",
	"UndecryptableMessage",
	"Receipt",
	"MediaRetry",
	"MediaRetryError",
	"ReadReceipt",

	// Groups and Contacts
	"GroupInfo",
	"JoinedGroup",
	"Picture",
	"BlocklistChange",
	"Blocklist",
	"Contact",

	// Connection and Session
	"Connected",
	"Disconnected",
	"ConnectFailure",
	"KeepAliveRestored",
	"KeepAliveTimeout",
	"LoggedOut",
	"ClientOutdated",
	"TemporaryBan",
	"StreamError",
	"StreamReplaced",
	"PairSuccess",
	"PairError",
	"QR",
	"QRScannedWithoutMultidevice",
	"ManualLoginReconnect",

	// Privacy and Settings
	"PrivacySettings",
	"PushNameSetting",
	"PushName",
	"UserAbout",
	"BusinessName",

	// Synchronization and State
	"AppState",
	"AppStateSyncComplete",
	"HistorySync",
	"OfflineSyncCompleted",
	"OfflineSyncPreview",

	// Calls
	"CallOffer",
	"CallAccept",
	"CallTerminate",
	"CallOfferNotice",
	"CallRelayLatency",
	"CallPreAccept",
	"CallReject",
	"CallTransport",
	"UnknownCallEvent",

	// Presence and Activity
	"Presence",
	"ChatPresence",

	// Identity
	"IdentityChange",

	// Erros
	"CATRefreshError",

	// Newsletter (WhatsApp Channels)
	"NewsletterJoin",
	"NewsletterLeave",
	"NewsletterMuteChange",
	"NewsletterLiveUpdate",
	"NewsletterMessageMeta",

	// Facebook/Meta Bridge
	"FBMessage",

	// Chat Management
	"Archive",
	"ClearChat",
	"DeleteChat",
	"DeleteForMe",
	"MarkChatAsRead",
	"Mute",
	"Pin",
	"Star",
	"UnarchiveChatsSetting",

	// Labels
	"LabelAssociationChat",
	"LabelAssociationMessage",
	"LabelEdit",

	// User Status
	"UserStatusMute",

	// Special - receives all events
	"All",
}

// Map for quick validation
var eventTypeMap map[string]bool

func init() {
	eventTypeMap = make(map[string]bool)
	for _, eventType := range supportedEventTypes {
		eventTypeMap[eventType] = true
	}
}

// isValidEventType validates if an event type is supported
func isValidEventType(eventType string) bool {
	return eventTypeMap[eventType]
}

// IsValidEventType validates if an event type is supported (exported version)
func IsValidEventType(eventType string) bool {
	// Check direct match first
	if isValidEventType(eventType) {
		return true
	}

	// Check mapped event types (lowercase to proper case)
	mapped := mapEventType(eventType)
	return isValidEventType(mapped)
}

// MapEventType maps common lowercase event names to proper event types (exported version)
func MapEventType(eventType string) string {
	return mapEventType(eventType)
}

// mapEventType maps common lowercase event names to proper event types
func mapEventType(eventType string) string {
	eventMap := map[string]string{
		"message":    "Message",
		"status":     "Receipt",
		"connection": "Connected",
		"call":       "CallOffer",
		"contact":    "Contact",
		"group":      "GroupInfo",
		"presence":   "Presence",
		"qr":         "QR",
		"pair":       "PairSuccess",
		"disconnect": "Disconnected",
		"error":      "ConnectFailure",
		"all":        "All",
	}

	if mapped, exists := eventMap[strings.ToLower(eventType)]; exists {
		return mapped
	}

	// Return original if no mapping found
	return eventType
}

// GetSupportedEventTypes returns the list of all supported event types
func GetSupportedEventTypes() []string {
	return supportedEventTypes
}

// EventProcessor - Simplified event handling with map-based handlers
type EventProcessor struct {
	sessionID      string
	sessionName    string
	webhookURL     string
	sessionRepo    session.Repository
	webhookService *webhooks.Service
	logger         logging.Logger
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

	// Privacy events
	"*events.PrivacySettings": (*EventProcessor).handlePrivacySettings,
	"*events.Blocklist":       (*EventProcessor).handleBlocklist,
}

// NewEventProcessor creates a new simplified event processor
func NewEventProcessor(sessionID, webhookURL string, sessionRepo session.Repository, webhookService *webhooks.Service) *EventProcessor {
	// Try to get session name for better logging
	sessionName := "unknown"
	if sessionRepo != nil {
		if sessionEntity, err := sessionRepo.GetByID(context.Background(), sessionID); err == nil {
			sessionName = sessionEntity.Name.Value()
		}
	}

	return &EventProcessor{
		sessionID:      sessionID,
		sessionName:    sessionName,
		webhookURL:     webhookURL,
		sessionRepo:    sessionRepo,
		webhookService: webhookService,
		logger:         logging.GetLogger().Sub("events").Sub(sessionID),
	}
}

// HandleEvent processes events using map-based handlers
func (ep *EventProcessor) HandleEvent(evt interface{}) {
	eventType := fmt.Sprintf("%T", evt)

	// Log events in a cleaner format
	ep.logger.Debugf("Event: %s", eventType)
	ep.logger.Debugf("Data:\n%s", ep.formatEventData(evt))

	// Send webhook if configured
	ep.sendWebhook(evt, eventType)

	if handler, exists := eventHandlers[eventType]; exists {
		handler(ep, evt)
	} else {
		ep.logger.Debugf("Unhandled: %s", eventType)
	}
}

// sendWebhook sends webhook notification for events
func (ep *EventProcessor) sendWebhook(evt interface{}, eventType string) {
	if ep.webhookService == nil {
		return
	}

	// Get session to check webhook configuration
	_, err := ep.sessionRepo.GetByID(context.Background(), ep.sessionID)
	if err != nil {
		ep.logger.Errorf("Failed to get session for webhook: %v", err)
		return
	}

	// Webhook functionality temporarily disabled during refactoring
	// Webhook configuration is now handled by separate webhook aggregate
	ep.logger.Debugf("Event processed: %s for session %s (webhook disabled during refactoring)", eventType, ep.sessionID)
}

// extractEventTypeName extracts clean event type name from Go type string
func (ep *EventProcessor) extractEventTypeName(eventType string) string {
	// Convert "*events.Message" to "Message"
	if len(eventType) > 8 && eventType[:8] == "*events." {
		return eventType[8:]
	}
	return eventType
}

// WebhookPayload represents the webhook payload structure with ordered fields
type WebhookPayload struct {
	Event     string      `json:"event"`
	SessionID string      `json:"sessionId"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// createWebhookPayload creates the webhook payload for different event types
func (ep *EventProcessor) createWebhookPayload(evt interface{}, eventType string) interface{} {
	// Use struct to ensure consistent field ordering in JSON
	return WebhookPayload{
		Event:     eventType,
		SessionID: ep.sessionID,
		Timestamp: time.Now().Unix(),
		Data:      evt, // Raw event payload from WhatsApp Meow - will appear last due to struct field order
	}
}

// Event category helpers for better organization

// isMessageEvent checks if event is message-related
func (ep *EventProcessor) isMessageEvent(eventType string) bool {
	messageEvents := []string{"*events.Message", "*events.Receipt"}
	for _, msgEvent := range messageEvents {
		if eventType == msgEvent {
			return true
		}
	}
	return false
}

// isConnectionEvent checks if event is connection-related
func (ep *EventProcessor) isConnectionEvent(eventType string) bool {
	connectionEvents := []string{"*events.Connected", "*events.Disconnected", "*events.LoggedOut"}
	for _, connEvent := range connectionEvents {
		if eventType == connEvent {
			return true
		}
	}
	return false
}

// isAuthEvent checks if event is authentication-related
func (ep *EventProcessor) isAuthEvent(eventType string) bool {
	authEvents := []string{"*events.QR", "*events.PairSuccess", "*events.PairError"}
	for _, authEvent := range authEvents {
		if eventType == authEvent {
			return true
		}
	}
	return false
}

// processConnectionEvents - Common processing for connection events
func (ep *EventProcessor) processConnectionEvents(eventType, status string) {
	ep.logger.Infof("Session %s %s", ep.sessionID, status)
	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
}

// processAuthEvents - Common processing for authentication events
func (ep *EventProcessor) processAuthEvents(eventType string, eventData interface{}) {
	// Add specific logging based on event type
	switch eventType {
	case "qr":
		if qr, ok := eventData.(*events.QR); ok && len(qr.Codes) > 0 {
			ep.logger.Debugf("QR codes available: %v", qr.Codes)
		}
	case "pair_success":
		ep.logger.Debugf("Pairing successful for session %s", ep.sessionID)
	case "pair_error":
		if pairErr, ok := eventData.(*events.PairError); ok {
			ep.logger.Errorf("Pairing error for session %s: %v", ep.sessionID, pairErr.Error)
		}
	}

	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
}

// EventProcessor handler methods - small and focused
func (ep *EventProcessor) handleMessage(evt interface{}) {
	msg := evt.(*events.Message)

	// Extract message details for comprehensive logging
	messageID := msg.Info.ID
	sender := msg.Info.Sender.String()
	chat := msg.Info.Chat.String()
	timestamp := msg.Info.Timestamp.Format("15:04:05")
	isGroup := msg.Info.IsGroup
	isFromMe := msg.Info.IsFromMe

	// Determine message type
	messageType := getMessageType(msg)

	// Create comprehensive log message with session name
	sessionInfo := fmt.Sprintf("Session:%-10s", truncateString(ep.sessionName, 10))

	if isFromMe {
		ep.logger.Infof("Message sent     [%-8s] ID:%-20s Type:%-8s Chat:%-25s %s Time:%s",
			"OUTGOING", messageID, messageType, truncateString(chat, 25), sessionInfo, timestamp)
	} else if isGroup {
		ep.logger.Infof("Message received [%-8s] ID:%-20s Type:%-8s From:%-25s Group:%-20s %s Time:%s",
			"GROUP", messageID, messageType, truncateString(sender, 25), truncateString(chat, 20), sessionInfo, timestamp)
	} else {
		ep.logger.Infof("Message received [%-8s] ID:%-20s Type:%-8s From:%-25s %s Time:%s",
			"DIRECT", messageID, messageType, truncateString(sender, 25), sessionInfo, timestamp)
	}

	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
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

	// Extract receipt details
	receiptType := getReceiptTypeString(receipt.Type)
	messageCount := len(receipt.MessageIDs)
	sender := receipt.MessageSource.Sender.String()
	timestamp := receipt.Timestamp.Format("15:04:05")

	// Log receipt with detailed information including session name
	sessionInfo := fmt.Sprintf("Session:%-10s", truncateString(ep.sessionName, 10))

	if messageCount == 1 {
		ep.logger.Infof("Receipt received [%-8s] ID:%-20s Type:%-8s From:%-25s %s Time:%s",
			"RECEIPT", receipt.MessageIDs[0], receiptType, truncateString(sender, 25), sessionInfo, timestamp)
	} else {
		ep.logger.Infof("Receipt received [%-8s] Count:%-3d        Type:%-8s From:%-25s %s Time:%s",
			"RECEIPT", messageCount, receiptType, truncateString(sender, 25), sessionInfo, timestamp)
	}

	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
}

func (ep *EventProcessor) handlePresence(evt interface{}) {
	ep.logger.Debugf("Presence update for session %s", ep.sessionID)

	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
}

func (ep *EventProcessor) handleChatPresence(evt interface{}) {
	ep.logger.Debugf("Chat presence update for session %s", ep.sessionID)
	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
}

func (ep *EventProcessor) handlePrivacySettings(evt interface{}) {
	privacySettings := evt.(*events.PrivacySettings)
	ep.logger.Debugf("Privacy settings changed for session %s", ep.sessionID)

	// Log detailed changes
	changes := []string{}
	if privacySettings.GroupAddChanged {
		changes = append(changes, fmt.Sprintf("GroupAdd: %s", string(privacySettings.NewSettings.GroupAdd)))
	}
	if privacySettings.LastSeenChanged {
		changes = append(changes, fmt.Sprintf("LastSeen: %s", string(privacySettings.NewSettings.LastSeen)))
	}
	if privacySettings.StatusChanged {
		changes = append(changes, fmt.Sprintf("Status: %s", string(privacySettings.NewSettings.Status)))
	}
	if privacySettings.ProfileChanged {
		changes = append(changes, fmt.Sprintf("Profile: %s", string(privacySettings.NewSettings.Profile)))
	}
	if privacySettings.ReadReceiptsChanged {
		changes = append(changes, fmt.Sprintf("ReadReceipts: %s", string(privacySettings.NewSettings.ReadReceipts)))
	}
	if privacySettings.OnlineChanged {
		changes = append(changes, fmt.Sprintf("Online: %s", string(privacySettings.NewSettings.Online)))
	}
	if privacySettings.CallAddChanged {
		changes = append(changes, fmt.Sprintf("CallAdd: %s", string(privacySettings.NewSettings.CallAdd)))
	}

	ep.logger.Infof("🔒 Privacy changes: %v", changes)

	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
}

func (ep *EventProcessor) handleBlocklist(evt interface{}) {
	blocklist := evt.(*events.Blocklist)
	ep.logger.Debugf("Blocklist changed for session %s", ep.sessionID)

	// Log detailed changes
	changes := []map[string]interface{}{}
	for _, change := range blocklist.Changes {
		changeData := map[string]interface{}{
			"jid":    change.JID.String(),
			"action": string(change.Action),
		}
		changes = append(changes, changeData)
		ep.logger.Debugf("Blocklist change: %s %s", string(change.Action), change.JID.String())
	}

	// Note: Webhook is already sent by HandleEvent method with complete event data
	// No need to send duplicate webhook here
}

// Helper functions shared across handlers (DRY principle)
// Note: The old sendWebhook stub function has been removed.
// All webhook calls now use ep.sendWebhook() which properly sends webhooks.

func createMessageData(msg *events.Message) map[string]interface{} {
	return map[string]interface{}{
		"messageId": msg.Info.ID,
		"from":      msg.Info.Sender.String(),
		"chat":      msg.Info.Chat.String(),
		"timestamp": msg.Info.Timestamp.Unix(),
		"fromMe":    msg.Info.IsFromMe,
		"isGroup":   msg.Info.IsGroup,
		"type":      getMessageType(msg),
	}
}

// getMessageType determines the type of message from the message content
func getMessageType(msg *events.Message) string {
	if msg.Message == nil {
		return "unknown"
	}

	switch {
	// Text messages
	case msg.Message.Conversation != nil:
		return "text"
	case msg.Message.ExtendedTextMessage != nil:
		return "text"

	// Media messages
	case msg.Message.ImageMessage != nil:
		return "image"
	case msg.Message.VideoMessage != nil:
		return "video"
	case msg.Message.AudioMessage != nil:
		return "audio"
	case msg.Message.DocumentMessage != nil:
		return "document"
	case msg.Message.StickerMessage != nil:
		return "sticker"

	// Location messages
	case msg.Message.LocationMessage != nil:
		return "location"
	case msg.Message.LiveLocationMessage != nil:
		return "live_loc"

	// Contact messages
	case msg.Message.ContactMessage != nil:
		return "contact"
	case msg.Message.ContactsArrayMessage != nil:
		return "contacts"

	// Interactive messages
	case msg.Message.ButtonsMessage != nil:
		return "buttons"
	case msg.Message.ListMessage != nil:
		return "list"
	case msg.Message.TemplateMessage != nil:
		return "template"
	case msg.Message.InteractiveMessage != nil:
		return "interact"

	// Poll and reactions
	case msg.Message.PollCreationMessage != nil:
		return "poll"
	case msg.Message.PollUpdateMessage != nil:
		return "poll_upd"
	case msg.Message.ReactionMessage != nil:
		return "reaction"

	// Protocol and system messages
	case msg.Message.ProtocolMessage != nil:
		return getProtocolMessageType(msg.Message.ProtocolMessage)
	case msg.Message.EphemeralMessage != nil:
		return "ephemeral"
	case msg.Message.ViewOnceMessage != nil:
		return "view_once"
	case msg.Message.ViewOnceMessageV2 != nil:
		return "view_once"
	case msg.Message.ViewOnceMessageV2Extension != nil:
		return "view_once"

	// Group messages
	case msg.Message.GroupInviteMessage != nil:
		return "grp_inv"

	// Payment and order
	case msg.Message.PaymentInviteMessage != nil:
		return "payment"
	case msg.Message.OrderMessage != nil:
		return "order"

	// Call messages
	case msg.Message.Call != nil:
		return "call"

	// Special message types
	case msg.Message.EditedMessage != nil:
		return "edit"
	case msg.Message.KeepInChatMessage != nil:
		return "keep_chat"
	case msg.Message.PinInChatMessage != nil:
		return "pin_chat"

	// Encrypted/wrapped messages
	case msg.Message.DeviceSentMessage != nil:
		return "dev_sent"
	case msg.Message.EncReactionMessage != nil:
		return "enc_react"

	default:
		return "other"
	}
}

// getProtocolMessageType determines the specific type of protocol message
func getProtocolMessageType(protocolMsg *waE2E.ProtocolMessage) string {
	if protocolMsg.Type == nil {
		return "protocol"
	}

	switch *protocolMsg.Type {
	case waE2E.ProtocolMessage_REVOKE:
		return "revoke"
	case waE2E.ProtocolMessage_EPHEMERAL_SETTING:
		return "ephemeral"
	case waE2E.ProtocolMessage_EPHEMERAL_SYNC_RESPONSE:
		return "eph_sync"
	case waE2E.ProtocolMessage_HISTORY_SYNC_NOTIFICATION:
		return "hist_sync"
	case waE2E.ProtocolMessage_APP_STATE_SYNC_KEY_SHARE:
		return "key_share"
	case waE2E.ProtocolMessage_APP_STATE_SYNC_KEY_REQUEST:
		return "key_req"
	case waE2E.ProtocolMessage_MSG_FANOUT_BACKFILL_REQUEST:
		return "backfill"
	case waE2E.ProtocolMessage_INITIAL_SECURITY_NOTIFICATION_SETTING_SYNC:
		return "sec_sync"
	case waE2E.ProtocolMessage_APP_STATE_FATAL_EXCEPTION_NOTIFICATION:
		return "fatal_exc"
	default:
		return "protocol"
	}
}

// getReceiptTypeString converts receipt type to a readable string
func getReceiptTypeString(receiptType types.ReceiptType) string {
	switch receiptType {
	case types.ReceiptTypeDelivered:
		return "delivered"
	case types.ReceiptTypeRead:
		return "read"
	case types.ReceiptTypePlayed:
		return "played"
	case types.ReceiptTypeReadSelf:
		return "read_self"
	case types.ReceiptTypeSender:
		return "sender"
	case types.ReceiptTypeRetry:
		return "retry"
	default:
		// If empty or unknown, try to infer from context
		if receiptType == "" {
			return "delivered" // Most common default
		}
		return string(receiptType)
	}
}

// formatEventData formats event data in a more readable JSON format
func (ep *EventProcessor) formatEventData(evt interface{}) string {
	// Try to marshal as JSON with indentation for better readability
	if jsonData, err := json.MarshalIndent(evt, "", "  "); err == nil {
		return string(jsonData)
	}

	// Fallback to structured format if JSON marshaling fails
	return ep.formatStructuredData(evt)
}

// formatStructuredData creates a structured representation of the event
func (ep *EventProcessor) formatStructuredData(evt interface{}) string {
	switch e := evt.(type) {
	case *events.Message:
		return ep.formatMessageEvent(e)
	case *events.Receipt:
		return ep.formatReceiptEvent(e)
	case *events.Connected:
		return "{\n  \"status\": \"connected\"\n}"
	case *events.OfflineSyncCompleted:
		return fmt.Sprintf("{\n  \"count\": %d\n}", e.Count)
	default:
		// Fallback to default formatting
		return fmt.Sprintf("%+v", evt)
	}
}

// formatMessageEvent formats message events in a clean structure
func (ep *EventProcessor) formatMessageEvent(msg *events.Message) string {
	msgType := getMessageType(msg)
	text := ep.extractMessageText(msg)

	return fmt.Sprintf(`{
  "id": "%s",
  "type": "%s",
  "from": "%s",
  "chat": "%s",
  "timestamp": "%s",
  "text": "%s",
  "isFromMe": %t,
  "isGroup": %t
}`,
		msg.Info.ID,
		msgType,
		msg.Info.MessageSource.Sender.String(),
		msg.Info.MessageSource.Chat.String(),
		msg.Info.Timestamp.Format("15:04:05"),
		text,
		msg.Info.MessageSource.IsFromMe,
		msg.Info.MessageSource.IsGroup,
	)
}

// formatReceiptEvent formats receipt events in a clean structure
func (ep *EventProcessor) formatReceiptEvent(receipt *events.Receipt) string {
	receiptType := getReceiptTypeString(receipt.Type)

	return fmt.Sprintf(`{
  "messageIds": ["%s"],
  "type": "%s",
  "from": "%s",
  "chat": "%s",
  "timestamp": "%s"
}`,
		receipt.MessageIDs[0],
		receiptType,
		receipt.MessageSource.Sender.String(),
		receipt.MessageSource.Chat.String(),
		receipt.Timestamp.Format("15:04:05"),
	)
}

// extractMessageText extracts text content from different message types
func (ep *EventProcessor) extractMessageText(msg *events.Message) string {
	if msg.Message == nil {
		return ""
	}

	switch {
	case msg.Message.Conversation != nil:
		return *msg.Message.Conversation
	case msg.Message.ExtendedTextMessage != nil:
		return *msg.Message.ExtendedTextMessage.Text
	case msg.Message.ImageMessage != nil:
		if msg.Message.ImageMessage.Caption != nil {
			return *msg.Message.ImageMessage.Caption
		}
		return "[Image]"
	case msg.Message.VideoMessage != nil:
		if msg.Message.VideoMessage.Caption != nil {
			return *msg.Message.VideoMessage.Caption
		}
		return "[Video]"
	case msg.Message.AudioMessage != nil:
		return "[Audio]"
	case msg.Message.DocumentMessage != nil:
		if msg.Message.DocumentMessage.FileName != nil {
			return fmt.Sprintf("[Document: %s]", *msg.Message.DocumentMessage.FileName)
		}
		return "[Document]"
	case msg.Message.StickerMessage != nil:
		return "[Sticker]"
	case msg.Message.LocationMessage != nil:
		return "[Location]"
	case msg.Message.ContactMessage != nil:
		return "[Contact]"
	default:
		return "[Other]"
	}
}

// truncateString truncates a string to the specified length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
