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

var supportedEventTypes = []string{
	"Message",
	"UndecryptableMessage",
	"Receipt",
	"MediaRetry",
	"MediaRetryError",
	"ReadReceipt",

	"GroupInfo",
	"JoinedGroup",
	"Picture",
	"BlocklistChange",
	"Blocklist",
	"Contact",

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

	"PrivacySettings",
	"PushNameSetting",
	"PushName",
	"UserAbout",
	"BusinessName",

	"AppState",
	"AppStateSyncComplete",
	"HistorySync",
	"OfflineSyncCompleted",
	"OfflineSyncPreview",

	"CallOffer",
	"CallAccept",
	"CallTerminate",
	"CallOfferNotice",
	"CallRelayLatency",
	"CallPreAccept",
	"CallReject",
	"CallTransport",
	"UnknownCallEvent",

	"Presence",
	"ChatPresence",

	"IdentityChange",

	"CATRefreshError",

	"NewsletterJoin",
	"NewsletterLeave",
	"NewsletterMuteChange",
	"NewsletterLiveUpdate",
	"NewsletterMessageMeta",

	"FBMessage",

	"Archive",
	"ClearChat",
	"DeleteChat",
	"DeleteForMe",
	"MarkChatAsRead",
	"Mute",
	"Pin",
	"Star",
	"UnarchiveChatsSetting",

	"LabelAssociationChat",
	"LabelAssociationMessage",
	"LabelEdit",

	"UserStatusMute",

	"All",
}

var eventTypeMap map[string]bool

func init() {
	eventTypeMap = make(map[string]bool)
	for _, eventType := range supportedEventTypes {
		eventTypeMap[eventType] = true
	}
}

func isValidEventType(eventType string) bool {
	return eventTypeMap[eventType]
}

func IsValidEventType(eventType string) bool {
	if isValidEventType(eventType) {
		return true
	}

	mapped := mapEventType(eventType)
	return isValidEventType(mapped)
}

func MapEventType(eventType string) string {
	return mapEventType(eventType)
}

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

	return eventType
}

func GetSupportedEventTypes() []string {
	return supportedEventTypes
}

type EventProcessor struct {
	sessionID      string
	sessionName    string
	webhookURL     string
	sessionRepo    session.Repository
	webhookService *webhooks.Service
	logger         logging.Logger
}

var eventHandlers = map[string]func(*EventProcessor, interface{}){
	"*events.Message": (*EventProcessor).handleMessage,
	"*events.Receipt": (*EventProcessor).handleReceipt,

	"*events.Connected":    (*EventProcessor).handleConnected,
	"*events.Disconnected": (*EventProcessor).handleDisconnected,
	"*events.LoggedOut":    (*EventProcessor).handleLoggedOut,

	"*events.QR":          (*EventProcessor).handleQR,
	"*events.PairSuccess": (*EventProcessor).handlePairSuccess,
	"*events.PairError":   (*EventProcessor).handlePairError,

	"*events.Presence":     (*EventProcessor).handlePresence,
	"*events.ChatPresence": (*EventProcessor).handleChatPresence,

	"*events.PrivacySettings": (*EventProcessor).handlePrivacySettings,
	"*events.Blocklist":       (*EventProcessor).handleBlocklist,
}

func NewEventProcessor(sessionID, webhookURL string, sessionRepo session.Repository, webhookService *webhooks.Service) *EventProcessor {
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

func (ep *EventProcessor) HandleEvent(evt interface{}) {
	eventType := fmt.Sprintf("%T", evt)

	ep.logger.Debugf("Event: %s", eventType)
	ep.logger.Debugf("Data:\n%s", ep.formatEventData(evt))

	ep.sendWebhook(evt, eventType)

	if handler, exists := eventHandlers[eventType]; exists {
		handler(ep, evt)
	} else {
		ep.logger.Debugf("Unhandled: %s", eventType)
	}
}

func (ep *EventProcessor) sendWebhook(evt interface{}, eventType string) {
	if ep.webhookService == nil {
		return
	}

	_, err := ep.sessionRepo.GetByID(context.Background(), ep.sessionID)
	if err != nil {
		ep.logger.Errorf("Failed to get session for webhook: %v", err)
		return
	}

	ep.logger.Debugf("Event processed: %s for session %s (webhook disabled during refactoring)", eventType, ep.sessionID)
}

func (ep *EventProcessor) extractEventTypeName(eventType string) string {
	if len(eventType) > 8 && eventType[:8] == "*events." {
		return eventType[8:]
	}
	return eventType
}

type WebhookPayload struct {
	Event     string      `json:"event"`
	SessionID string      `json:"sessionId"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func (ep *EventProcessor) createWebhookPayload(evt interface{}, eventType string) interface{} {
	return WebhookPayload{
		Event:     eventType,
		SessionID: ep.sessionID,
		Timestamp: time.Now().Unix(),
		Data:      evt, // Raw event payload from WhatsApp Meow - will appear last due to struct field order
	}
}


func (ep *EventProcessor) isMessageEvent(eventType string) bool {
	messageEvents := []string{"*events.Message", "*events.Receipt"}
	for _, msgEvent := range messageEvents {
		if eventType == msgEvent {
			return true
		}
	}
	return false
}

func (ep *EventProcessor) isConnectionEvent(eventType string) bool {
	connectionEvents := []string{"*events.Connected", "*events.Disconnected", "*events.LoggedOut"}
	for _, connEvent := range connectionEvents {
		if eventType == connEvent {
			return true
		}
	}
	return false
}

func (ep *EventProcessor) isAuthEvent(eventType string) bool {
	authEvents := []string{"*events.QR", "*events.PairSuccess", "*events.PairError"}
	for _, authEvent := range authEvents {
		if eventType == authEvent {
			return true
		}
	}
	return false
}

func (ep *EventProcessor) processConnectionEvents(eventType, status string) {
	ep.logger.Infof("Session %s %s", ep.sessionID, status)
}

func (ep *EventProcessor) processAuthEvents(eventType string, eventData interface{}) {
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

}

func (ep *EventProcessor) handleMessage(evt interface{}) {
	msg := evt.(*events.Message)

	messageID := msg.Info.ID
	sender := msg.Info.Sender.String()
	chat := msg.Info.Chat.String()
	timestamp := msg.Info.Timestamp.Format("15:04:05")
	isGroup := msg.Info.IsGroup
	isFromMe := msg.Info.IsFromMe

	messageType := getMessageType(msg)

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

	receiptType := getReceiptTypeString(receipt.Type)
	messageCount := len(receipt.MessageIDs)
	sender := receipt.MessageSource.Sender.String()
	timestamp := receipt.Timestamp.Format("15:04:05")

	sessionInfo := fmt.Sprintf("Session:%-10s", truncateString(ep.sessionName, 10))

	if messageCount == 1 {
		ep.logger.Infof("Receipt received [%-8s] ID:%-20s Type:%-8s From:%-25s %s Time:%s",
			"RECEIPT", receipt.MessageIDs[0], receiptType, truncateString(sender, 25), sessionInfo, timestamp)
	} else {
		ep.logger.Infof("Receipt received [%-8s] Count:%-3d        Type:%-8s From:%-25s %s Time:%s",
			"RECEIPT", messageCount, receiptType, truncateString(sender, 25), sessionInfo, timestamp)
	}

}

func (ep *EventProcessor) handlePresence(evt interface{}) {
	ep.logger.Debugf("Presence update for session %s", ep.sessionID)

}

func (ep *EventProcessor) handleChatPresence(evt interface{}) {
	ep.logger.Debugf("Chat presence update for session %s", ep.sessionID)
}

func (ep *EventProcessor) handlePrivacySettings(evt interface{}) {
	privacySettings := evt.(*events.PrivacySettings)
	ep.logger.Debugf("Privacy settings changed for session %s", ep.sessionID)

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

}

func (ep *EventProcessor) handleBlocklist(evt interface{}) {
	blocklist := evt.(*events.Blocklist)
	ep.logger.Debugf("Blocklist changed for session %s", ep.sessionID)

	changes := []map[string]interface{}{}
	for _, change := range blocklist.Changes {
		changeData := map[string]interface{}{
			"jid":    change.JID.String(),
			"action": string(change.Action),
		}
		changes = append(changes, changeData)
		ep.logger.Debugf("Blocklist change: %s %s", string(change.Action), change.JID.String())
	}

}


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

func getMessageType(msg *events.Message) string {
	if msg.Message == nil {
		return "unknown"
	}

	switch {
	case msg.Message.Conversation != nil:
		return "text"
	case msg.Message.ExtendedTextMessage != nil:
		return "text"

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

	case msg.Message.LocationMessage != nil:
		return "location"
	case msg.Message.LiveLocationMessage != nil:
		return "live_loc"

	case msg.Message.ContactMessage != nil:
		return "contact"
	case msg.Message.ContactsArrayMessage != nil:
		return "contacts"

	case msg.Message.ButtonsMessage != nil:
		return "buttons"
	case msg.Message.ListMessage != nil:
		return "list"
	case msg.Message.TemplateMessage != nil:
		return "template"
	case msg.Message.InteractiveMessage != nil:
		return "interact"

	case msg.Message.PollCreationMessage != nil:
		return "poll"
	case msg.Message.PollUpdateMessage != nil:
		return "poll_upd"
	case msg.Message.ReactionMessage != nil:
		return "reaction"

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

	case msg.Message.GroupInviteMessage != nil:
		return "grp_inv"

	case msg.Message.PaymentInviteMessage != nil:
		return "payment"
	case msg.Message.OrderMessage != nil:
		return "order"

	case msg.Message.Call != nil:
		return "call"

	case msg.Message.EditedMessage != nil:
		return "edit"
	case msg.Message.KeepInChatMessage != nil:
		return "keep_chat"
	case msg.Message.PinInChatMessage != nil:
		return "pin_chat"

	case msg.Message.DeviceSentMessage != nil:
		return "dev_sent"
	case msg.Message.EncReactionMessage != nil:
		return "enc_react"

	default:
		return "other"
	}
}

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
		if receiptType == "" {
			return "delivered" // Most common default
		}
		return string(receiptType)
	}
}

func (ep *EventProcessor) formatEventData(evt interface{}) string {
	if jsonData, err := json.MarshalIndent(evt, "", "  "); err == nil {
		return string(jsonData)
	}

	return ep.formatStructuredData(evt)
}

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
		return fmt.Sprintf("%+v", evt)
	}
}

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

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
