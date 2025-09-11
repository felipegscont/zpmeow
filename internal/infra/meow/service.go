package meow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/infra/webhook"
	"zpmeow/internal/types"

	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)


type MediaDownloadInfo struct {
	MediaType     string `json:"mediaType"`     // image, video, audio, document
	URL           string `json:"url"`
	DirectPath    string `json:"directPath"`
	MediaKey      []byte `json:"mediaKey"`
	MimeType      string `json:"mimeType"`
	FileEncSHA256 []byte `json:"fileEncSHA256"`
	FileSHA256    []byte `json:"fileSHA256"`
	FileLength    uint64 `json:"fileLength"`
}




func parseJID(arg string) (waTypes.JID, bool) {
	// First try to parse as a complete JID (for groups like 120363313346913103@g.us)
	if strings.Contains(arg, "@") {
		jid, err := waTypes.ParseJID(arg)
		if err == nil {
			return jid, true
		}
	}

	// If it's a phone number, convert to WhatsApp JID
	// Remove any non-numeric characters except +
	phone := strings.ReplaceAll(arg, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Remove + prefix if present
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}

	// Validate that it's all digits
	for _, char := range phone {
		if char < '0' || char > '9' {
			return waTypes.JID{}, false
		}
	}

	// Create WhatsApp user JID
	jid := waTypes.NewJID(phone, waTypes.DefaultUserServer)
	return jid, true
}


func (m *MeowServiceImpl) validateAndGetClient(sessionID, to string) (*MeowClient, waTypes.JID, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, waTypes.JID{}, fmt.Errorf("client not found for session %s", sessionID)
	}

	jid, ok := parseJID(to)
	if !ok {
		return nil, waTypes.JID{}, fmt.Errorf("invalid JID %s", to)
	}

	return client, jid, nil
}


type MeowServiceImpl struct {
	clientManager *ClientManager
	logger        logger.Logger
	waLogger      waLog.Logger
}


func NewMeowService(db *sqlx.DB, container *sqlstore.Container, waLogger waLog.Logger, sessionService session.SessionService) session.WhatsAppService {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("meow-service")

	// Create webhook service
	webhookService := webhook.NewWebhookService()

	clientManager := NewClientManager(db, container, waLogger, webhookService, sessionService)

	service := &MeowServiceImpl{
		clientManager: clientManager,
		logger:        appLogger,
		waLogger:      waLogger,
	}

	return service
}


func (m *MeowServiceImpl) StartClient(sessionID string) error {

	if err := Validation.ValidateSessionID(sessionID); err != nil {
		return Error.WrapError(err, "start client validation failed")
	}

	ctx, cancel := Context.WithTimeout(context.Background())
	defer cancel()

	m.logger.Infof("Starting client for session %s", sessionID)

	err := m.clientManager.StartClient(ctx, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to start client for session %s: %v", sessionID, err)
		return Error.WrapError(err, "failed to start client")
	}

	m.logger.Infof("Client started successfully for session %s", sessionID)
	return nil
}


func (m *MeowServiceImpl) StopClient(sessionID string) error {

	if err := Validation.ValidateSessionID(sessionID); err != nil {
		return Error.WrapError(err, "stop client validation failed")
	}

	m.logger.Infof("Stopping client for session %s", sessionID)

	err := m.clientManager.StopClient(sessionID)
	if err != nil {
		m.logger.Errorf("Failed to stop client for session %s: %v", sessionID, err)
		return Error.WrapError(err, "failed to stop client")
	}

	m.logger.Infof("Client stopped successfully for session %s", sessionID)
	return nil
}


func (m *MeowServiceImpl) LogoutClient(sessionID string) error {
	ctx := context.Background()
	
	m.logger.Infof("Logging out client for session %s", sessionID)
	
	err := m.clientManager.LogoutClient(ctx, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to logout client for session %s: %v", sessionID, err)
		return fmt.Errorf("failed to logout client: %w", err)
	}

	m.logger.Infof("Client logged out successfully for session %s", sessionID)
	return nil
}


func (m *MeowServiceImpl) GetQRCode(sessionID string) (string, error) {
	ctx := context.Background()
	
	m.logger.Debugf("Getting QR code for session %s", sessionID)
	
	qrCode, err := m.clientManager.GetQRCode(ctx, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to get QR code for session %s: %v", sessionID, err)
		return "", fmt.Errorf("failed to get QR code: %w", err)
	}

	if qrCode == "" {
		return "", fmt.Errorf("QR code not available yet")
	}

	m.logger.Debugf("QR code retrieved for session %s", sessionID)
	return qrCode, nil
}


func (m *MeowServiceImpl) PairPhone(sessionID, phoneNumber string) (string, error) {
	ctx := context.Background()
	
	m.logger.Infof("Pairing phone %s with session %s", phoneNumber, sessionID)
	
	pairingCode, err := m.clientManager.PairPhone(ctx, sessionID, phoneNumber)
	if err != nil {
		m.logger.Errorf("Failed to pair phone for session %s: %v", sessionID, err)
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	m.logger.Infof("Phone pairing initiated for session %s", sessionID)
	return pairingCode, nil
}


func (m *MeowServiceImpl) IsClientConnected(sessionID string) bool {
	connected := m.clientManager.IsClientConnected(sessionID)
	m.logger.Debugf("Client connection status for session %s: %t", sessionID, connected)
	return connected
}


func (m *MeowServiceImpl) GetClientStatus(sessionID string) types.Status {
	status := m.clientManager.GetClientStatus(sessionID)
	m.logger.Debugf("Client status for session %s: %s", sessionID, status)
	return status
}


func (m *MeowServiceImpl) CheckAndUpdateDeviceJID(sessionID string) {
	m.logger.Infof("Checking and updating device JID for session %s", sessionID)
	m.clientManager.CheckAndUpdateDeviceJID(sessionID)
}


func (m *MeowServiceImpl) GetClient(sessionID string) (*MeowClient, bool) {
	return m.clientManager.GetClient(sessionID)
}


func (m *MeowServiceImpl) GetAllClients() map[string]*MeowClient {
	return m.clientManager.GetAllClients()
}


func (m *MeowServiceImpl) GetStats() map[string]interface{} {
	stats := m.clientManager.GetStats()
	m.logger.Debugf("Service stats: %+v", stats)
	return stats
}


func (m *MeowServiceImpl) Cleanup() {
	m.logger.Infof("Performing service cleanup")
	m.clientManager.Cleanup()
}


func (m *MeowServiceImpl) Shutdown(ctx context.Context) error {
	m.logger.Infof("Shutting down Meow service")
	return m.clientManager.Shutdown(ctx)
}




func (m *MeowServiceImpl) IsHealthy() bool {
	
	return m.clientManager != nil
}


func (m *MeowServiceImpl) GetHealthStatus() map[string]interface{} {
	stats := m.GetStats()
	
	return map[string]interface{}{
		"healthy":     m.IsHealthy(),
		"stats":       stats,
		"service":     "meow-whatsapp",
		"version":     "1.0.0",
	}
}




func (m *MeowServiceImpl) SendMessage(ctx context.Context, sessionID, to, message string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	

	
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendMessage(ctx, jid, message)
}


func (m *MeowServiceImpl) SendTextMessage(ctx context.Context, sessionID, to, text string, contextInfo *waE2E.ContextInfo) (*whatsmeow.SendResponse, error) {
	m.logger.Infof("DEBUG: SendTextMessage called - sessionID: %s, to: %s, text: %s", sessionID, to, text)


	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		m.logger.Errorf("DEBUG: Validation failed: %v", err)
		return nil, err
	}

	m.logger.Infof("DEBUG: Client found and JID parsed successfully: %s -> %s", to, jid.String())


	m.logger.Infof("DEBUG: Skipping IsConnected() check to avoid deadlock")

	m.logger.Infof("DEBUG: Calling client.SendTextMessage...")
	resp, err := client.SendTextMessage(ctx, jid, text, contextInfo)
	if err != nil {
		m.logger.Errorf("DEBUG: client.SendTextMessage failed: %v", err)
		return nil, err
	}

	m.logger.Infof("DEBUG: client.SendTextMessage completed successfully")
	return resp, nil
}


func (m *MeowServiceImpl) SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name, address string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendLocationMessage(ctx, jid, latitude, longitude, name, address)
}


func (m *MeowServiceImpl) SendContactMessage(ctx context.Context, sessionID, to, displayName, vcard string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendContactMessage(ctx, jid, displayName, vcard)
}




func (m *MeowServiceImpl) SetChatPresence(ctx context.Context, sessionID, chatJID, state, media string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}


	jid, err := JID.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid JID %s: %w", chatJID, err)
	}


	if err := m.validatePresenceParams(state, media); err != nil {
		return err
	}


	chatPresence := waTypes.ChatPresence(state)
	mediaType := waTypes.ChatPresenceMedia(media)


	err = client.client.SendChatPresence(jid, chatPresence, mediaType)
	if err != nil {
		return fmt.Errorf("failed to set chat presence: %w", err)
	}

	return nil
}


func (m *MeowServiceImpl) validatePresenceParams(state, media string) error {

	if state != "composing" && state != "paused" {
		return fmt.Errorf("invalid state: %s (valid: composing, paused)", state)
	}


	if media != "" && media != "audio" {
		return fmt.Errorf("invalid media: %s (valid: \"\", audio)", media)
	}

	return nil
}


func (m *MeowServiceImpl) MarkMessageRead(ctx context.Context, sessionID, chatJID string, messageIDs []string) error {
	m.logger.Infof("MarkMessageRead service called with sessionID: %s, chatJID: %s, messageIDs: %v", sessionID, chatJID, messageIDs)

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse chat JID using our utility function (handles phone numbers)
	jid, ok := parseJID(chatJID)
	if !ok {
		m.logger.Errorf("Failed to parse chat JID: %s", chatJID)
		return fmt.Errorf("invalid chat JID %s", chatJID)
	}

	m.logger.Infof("Parsed JID: %s -> %s", chatJID, jid.String())

	err := client.MarkMessageRead(ctx, jid, messageIDs)
	if err != nil {
		m.logger.Errorf("MarkMessageRead failed: %v", err)
		return err
	}

	m.logger.Infof("MarkMessageRead service completed successfully")
	return nil
}


func (m *MeowServiceImpl) CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*waTypes.GroupInfo, error) {
	m.logger.Infof("CreateGroup service called with sessionID: %s, name: %s, participants: %v", sessionID, name, participants)

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		m.logger.Errorf("Client not found for session: %s", sessionID)
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	m.logger.Infof("Skipping IsConnected() check to avoid deadlock")

	// Convert participant phone numbers to JIDs using our utility function
	participantJIDs := make([]waTypes.JID, len(participants))
	for i, participant := range participants {
		jid, ok := parseJID(participant)
		if !ok {
			m.logger.Errorf("Failed to parse participant JID: %s", participant)
			return nil, fmt.Errorf("invalid participant JID %s", participant)
		}
		m.logger.Infof("Parsed participant JID: %s -> %s", participant, jid.String())
		participantJIDs[i] = jid
	}

	m.logger.Infof("Calling client.CreateGroup with name: %s, participants: %v", name, participantJIDs)
	groupInfo, err := client.CreateGroup(ctx, name, participantJIDs)
	if err != nil {
		m.logger.Errorf("Failed to create group: %v", err)
		return nil, err
	}

	m.logger.Infof("Successfully created group: %s (JID: %s)", name, groupInfo.JID.String())
	m.logger.Infof("CreateGroup service completed successfully")
	return groupInfo, nil
}


func (m *MeowServiceImpl) GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.GetGroupInfo(ctx, jid)
}


func (m *MeowServiceImpl) JoinGroupWithLink(ctx context.Context, sessionID, inviteCode string) (*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	return client.JoinGroupWithLink(ctx, inviteCode)
}


func (m *MeowServiceImpl) LeaveGroup(ctx context.Context, sessionID, groupJID string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.LeaveGroup(ctx, jid)
}


func (m *MeowServiceImpl) GetGroupInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return "", fmt.Errorf("client not found for session %s", sessionID)
	}

	
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.GetGroupInviteLink(ctx, jid, reset)
}


func (m *MeowServiceImpl) UpdateGroupParticipants(ctx context.Context, sessionID, groupJID string, participants []string, action string) error {
	m.logger.Infof("UpdateGroupParticipants service called with sessionID: %s, groupJID: %s, participants: %v, action: %s", sessionID, groupJID, participants, action)

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		m.logger.Errorf("Client not found for session: %s", sessionID)
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse group JID
	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		m.logger.Errorf("Failed to parse group JID: %s, error: %v", groupJID, err)
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}
	m.logger.Infof("Parsed group JID: %s -> %s", groupJID, groupJIDParsed.String())

	// Parse participant JIDs using our utility function
	participantJIDs := make([]waTypes.JID, len(participants))
	for i, participant := range participants {
		pJid, ok := parseJID(participant)
		if !ok {
			m.logger.Errorf("Failed to parse participant JID: %s", participant)
			return fmt.Errorf("invalid participant JID %s", participant)
		}
		m.logger.Infof("Parsed participant JID: %s -> %s", participant, pJid.String())
		participantJIDs[i] = pJid
	}

	// Validate action
	validActions := []string{"add", "remove", "promote", "demote"}
	isValidAction := false
	for _, validAction := range validActions {
		if action == validAction {
			isValidAction = true
			break
		}
	}
	if !isValidAction {
		m.logger.Errorf("Invalid action: %s", action)
		return fmt.Errorf("invalid action: %s. Must be one of: add, remove, promote, demote", action)
	}

	m.logger.Infof("Calling client.UpdateGroupParticipants")
	err = client.UpdateGroupParticipants(ctx, groupJIDParsed, participantJIDs, action)
	if err != nil {
		m.logger.Errorf("Failed to update group participants: %v", err)
		return err
	}

	m.logger.Infof("UpdateGroupParticipants service completed successfully")
	return nil
}


func (m *MeowServiceImpl) SetGroupName(ctx context.Context, sessionID, groupJID, name string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.SetGroupName(ctx, jid, name)
}


func (m *MeowServiceImpl) SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.SetGroupTopic(ctx, jid, topic)
}


func (m *MeowServiceImpl) SendImageMessage(ctx context.Context, sessionID, to string, imageData []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendImageMessage(ctx, jid, imageData, caption, mimeType)
}


func (m *MeowServiceImpl) SendAudioMessage(ctx context.Context, sessionID, to string, audioData []byte, mimeType string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendAudioMessage(ctx, jid, audioData, mimeType)
}


func (m *MeowServiceImpl) SendDocumentMessage(ctx context.Context, sessionID, to string, documentData []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendDocumentMessage(ctx, jid, documentData, filename, caption, mimeType)
}


func (m *MeowServiceImpl) SendVideoMessage(ctx context.Context, sessionID, to string, videoData []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendVideoMessage(ctx, jid, videoData, caption, mimeType)
}


func (m *MeowServiceImpl) SendStickerMessage(ctx context.Context, sessionID, to string, stickerData []byte, mimeType string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendStickerMessage(ctx, jid, stickerData, mimeType)
}


func (m *MeowServiceImpl) SendButtonsMessage(ctx context.Context, sessionID, to, text string, buttons []types.Button, footer string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendButtonsMessage(ctx, jid, text, buttons, footer)
}


func (m *MeowServiceImpl) SendListMessage(ctx context.Context, sessionID, to, text, buttonText string, sections []types.Section, footer string) (*whatsmeow.SendResponse, error) {

	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendListMessage(ctx, jid, text, buttonText, sections, footer)
}






func (m *MeowServiceImpl) DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error {
	m.logger.Infof("DeleteMessage service called with sessionID: %s, chatJID: %s, messageID: %s, forEveryone: %v", sessionID, chatJID, messageID, forEveryone)


	if err := m.validateChatOperation(sessionID, chatJID, messageID); err != nil {
		return err
	}

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}


	m.logger.Infof("Skipping IsConnected() check to avoid deadlock")


	// Parse chat JID using our utility function (handles phone numbers)
	jid, ok := parseJID(chatJID)
	if !ok {
		m.logger.Errorf("Failed to parse chat JID: %s", chatJID)
		return fmt.Errorf("invalid chat JID %s", chatJID)
	}

	m.logger.Infof("Parsed JID: %s -> %s", chatJID, jid.String())

	// Build revoke message
	m.logger.Infof("Building revoke message for messageID: %s", messageID)
	revokeMsg := client.client.BuildRevoke(jid, waTypes.EmptyJID, waTypes.MessageID(messageID))

	// Send revoke message
	m.logger.Infof("Sending revoke message to WhatsApp API")
	_, err := client.client.SendMessage(ctx, jid, revokeMsg)
	if err != nil {
		m.logger.Errorf("Failed to delete message: %v", err)
		return fmt.Errorf("failed to delete message %s: %w", messageID, err)
	}

	m.logger.Infof("Successfully deleted message %s", messageID)
	m.logger.Infof("DeleteMessage service completed successfully")
	return nil
}


func (m *MeowServiceImpl) validateChatOperation(sessionID, chatJID, messageID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	if chatJID == "" {
		return fmt.Errorf("chat JID cannot be empty")
	}
	if messageID == "" {
		return fmt.Errorf("message ID cannot be empty")
	}


	actualMessageID := messageID
	if strings.HasPrefix(messageID, "me:") {
		actualMessageID = messageID[len("me:"):]
	}


	if len(actualMessageID) < 10 {
		return fmt.Errorf("invalid message ID format")
	}

	return nil
}


func (m *MeowServiceImpl) EditMessage(ctx context.Context, sessionID, chatJID, messageID, newText string) (*types.SendResponse, error) {
	m.logger.Infof("EditMessage service called with sessionID: %s, chatJID: %s, messageID: %s, newText: %s", sessionID, chatJID, messageID, newText)


	if err := m.validateChatOperation(sessionID, chatJID, messageID); err != nil {
		return nil, err
	}


	if newText == "" {
		return nil, fmt.Errorf("new text cannot be empty")
	}
	if len(newText) > 4096 {
		return nil, fmt.Errorf("new text too long (max 4096 characters)")
	}

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}


	m.logger.Infof("Skipping IsConnected() check to avoid deadlock")


	// Parse chat JID using our utility function (handles phone numbers)
	jid, ok := parseJID(chatJID)
	if !ok {
		m.logger.Errorf("Failed to parse chat JID: %s", chatJID)
		return nil, fmt.Errorf("invalid chat JID %s", chatJID)
	}

	m.logger.Infof("Parsed JID: %s -> %s", chatJID, jid.String())

	// Create text message
	m.logger.Infof("Creating text message with newText: %s", newText)
	textMsg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &newText,
		},
	}

	// Build edit message
	m.logger.Infof("Building edit message for messageID: %s", messageID)
	editMsg := client.client.BuildEdit(jid, waTypes.MessageID(messageID), textMsg)

	// Send edit message
	m.logger.Infof("Sending edit message to WhatsApp API")
	resp, err := client.client.SendMessage(ctx, jid, editMsg)
	if err != nil {
		m.logger.Errorf("Failed to edit message: %v", err)
		return nil, fmt.Errorf("failed to edit message %s: %w", messageID, err)
	}

	m.logger.Infof("Successfully edited message %s", messageID)

	// Convert response
	m.logger.Infof("Converting response from WhatsApp API")
	response := types.NewSendResponseFromWhatsmeow(&resp, messageID)
	m.logger.Infof("EditMessage service completed successfully")
	return &response, nil
}


func (m *MeowServiceImpl) DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error) {
	m.logger.Infof("Downloading media for message %s in session %s", messageID, sessionID)






	return nil, "", fmt.Errorf("media download requires message metadata storage - not implemented in this simplified version")
}


func (m *MeowServiceImpl) DownloadMediaWithInfo(ctx context.Context, sessionID string, mediaInfo MediaDownloadInfo) ([]byte, string, error) {
	m.logger.Infof("Downloading media with info for session %s", sessionID)

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, "", fmt.Errorf("client not found for session %s", sessionID)
	}

	var mediaData []byte
	var err error
	var mimeType string


	switch mediaInfo.MediaType {
	case "image":
		msg := &waE2E.Message{ImageMessage: &waE2E.ImageMessage{
			URL:           &mediaInfo.URL,
			DirectPath:    &mediaInfo.DirectPath,
			MediaKey:      mediaInfo.MediaKey,
			Mimetype:      &mediaInfo.MimeType,
			FileEncSHA256: mediaInfo.FileEncSHA256,
			FileSHA256:    mediaInfo.FileSHA256,
			FileLength:    &mediaInfo.FileLength,
		}}
		mediaData, err = client.client.Download(ctx, msg.GetImageMessage())
		mimeType = msg.GetImageMessage().GetMimetype()

	case "video":
		msg := &waE2E.Message{VideoMessage: &waE2E.VideoMessage{
			URL:           &mediaInfo.URL,
			DirectPath:    &mediaInfo.DirectPath,
			MediaKey:      mediaInfo.MediaKey,
			Mimetype:      &mediaInfo.MimeType,
			FileEncSHA256: mediaInfo.FileEncSHA256,
			FileSHA256:    mediaInfo.FileSHA256,
			FileLength:    &mediaInfo.FileLength,
		}}
		mediaData, err = client.client.Download(ctx, msg.GetVideoMessage())
		mimeType = msg.GetVideoMessage().GetMimetype()

	case "audio":
		msg := &waE2E.Message{AudioMessage: &waE2E.AudioMessage{
			URL:           &mediaInfo.URL,
			DirectPath:    &mediaInfo.DirectPath,
			MediaKey:      mediaInfo.MediaKey,
			Mimetype:      &mediaInfo.MimeType,
			FileEncSHA256: mediaInfo.FileEncSHA256,
			FileSHA256:    mediaInfo.FileSHA256,
			FileLength:    &mediaInfo.FileLength,
		}}
		mediaData, err = client.client.Download(ctx, msg.GetAudioMessage())
		mimeType = msg.GetAudioMessage().GetMimetype()

	case "document":
		msg := &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{
			URL:           &mediaInfo.URL,
			DirectPath:    &mediaInfo.DirectPath,
			MediaKey:      mediaInfo.MediaKey,
			Mimetype:      &mediaInfo.MimeType,
			FileEncSHA256: mediaInfo.FileEncSHA256,
			FileSHA256:    mediaInfo.FileSHA256,
			FileLength:    &mediaInfo.FileLength,
		}}
		mediaData, err = client.client.Download(ctx, msg.GetDocumentMessage())
		mimeType = msg.GetDocumentMessage().GetMimetype()

	default:
		return nil, "", fmt.Errorf("unsupported media type: %s", mediaInfo.MediaType)
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to download %s media: %w", mediaInfo.MediaType, err)
	}

	m.logger.Infof("Successfully downloaded %s media (%d bytes)", mediaInfo.MediaType, len(mediaData))
	return mediaData, mimeType, nil
}


func (m *MeowServiceImpl) ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string) error {
	m.logger.Infof("Reacting to message %s in chat %s for session %s with emoji: %s", messageID, chatJID, sessionID, emoji)


	if err := m.validateChatOperation(sessionID, chatJID, messageID); err != nil {
		return err
	}


	if emoji == "" {
		return fmt.Errorf("emoji cannot be empty")
	}
	if len(emoji) > 10 {
		return fmt.Errorf("emoji too long")
	}

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}


	m.logger.Infof("Skipping IsConnected() check to avoid deadlock")


	jid, ok := parseJID(chatJID)
	if !ok {
		return fmt.Errorf("invalid chat JID %s", chatJID)
	}


	actualMessageID := messageID
	fromMe := waTypes.EmptyJID
	if strings.HasPrefix(messageID, "me:") {
		actualMessageID = messageID[len("me:"):]

		fromMe = client.client.Store.ID.ToNonAD()
		m.logger.Infof("Detected 'me:' prefix, treating as own message. Original: %s, Actual: %s", messageID, actualMessageID)
	}


	reaction := emoji
	if reaction == "remove" {
		reaction = ""
		m.logger.Infof("Converting 'remove' to empty string for reaction removal")
	}


	reactionMsg := client.client.BuildReaction(jid, fromMe, waTypes.MessageID(actualMessageID), reaction)


	_, err := client.client.SendMessage(ctx, jid, reactionMsg)
	if err != nil {
		return fmt.Errorf("failed to react to message %s: %w", messageID, err)
	}

	m.logger.Infof("Successfully reacted to message %s with %s", messageID, emoji)
	return nil
}


func (m *MeowServiceImpl) SendPollMessage(ctx context.Context, sessionID, to, name string, options []string, selectableCount int) (*whatsmeow.SendResponse, error) {
	m.logger.Infof("DEBUG: SendPollMessage called - sessionID: %s, to: %s, name: %s", sessionID, to, name)


	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		m.logger.Errorf("DEBUG: Validation failed: %v", err)
		return nil, err
	}

	m.logger.Infof("DEBUG: Client found and JID parsed successfully: %s -> %s", to, jid.String())
	m.logger.Infof("DEBUG: Calling client.SendPollMessage...")

	resp, err := client.SendPollMessage(ctx, jid, name, options, selectableCount)
	if err != nil {
		m.logger.Errorf("DEBUG: client.SendPollMessage failed: %v", err)
		return nil, err
	}

	m.logger.Infof("DEBUG: SendPollMessage completed successfully")
	return resp, nil
}


func (m *MeowServiceImpl) ListGroups(ctx context.Context, sessionID string) ([]*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}


	groups, err := client.GetClient().GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	return groups, nil
}

// SetGlobalPresence sets the global presence status for the session
func (m *MeowServiceImpl) SetGlobalPresence(ctx context.Context, sessionID string, presence waTypes.Presence) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	err := client.client.SendPresence(presence)
	if err != nil {
		return fmt.Errorf("failed to set global presence: %w", err)
	}

	return nil
}

// CheckUsersOnWhatsApp checks if phone numbers are registered on WhatsApp
func (m *MeowServiceImpl) CheckUsersOnWhatsApp(ctx context.Context, sessionID string, phones []string) ([]waTypes.IsOnWhatsAppResponse, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	resp, err := client.client.IsOnWhatsApp(phones)
	if err != nil {
		return nil, fmt.Errorf("failed to check users on WhatsApp: %w", err)
	}

	return resp, nil
}

// GetUserInfo gets detailed information about WhatsApp users
func (m *MeowServiceImpl) GetUserInfo(ctx context.Context, sessionID string, phones []string) (map[waTypes.JID]waTypes.UserInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	// Convert phone numbers to JIDs
	var jids []waTypes.JID
	for _, phone := range phones {
		jid, ok := parseJID(phone)
		if !ok {
			return nil, fmt.Errorf("invalid phone number: %s", phone)
		}
		jids = append(jids, jid)
	}

	resp, err := client.client.GetUserInfo(jids)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return resp, nil
}

// AvatarInfo represents avatar information
type AvatarInfo struct {
	URL       string `json:"url,omitempty"`
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	DirectURL string `json:"directUrl,omitempty"`
}

// GetUserAvatar gets avatar information for a WhatsApp user
func (m *MeowServiceImpl) GetUserAvatar(ctx context.Context, sessionID, phone string, preview bool) (*AvatarInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse phone number to JID
	jid, ok := parseJID(phone)
	if !ok {
		return nil, fmt.Errorf("invalid phone number: %s", phone)
	}

	// Get avatar info
	avatarInfo, err := client.client.GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
		Preview: preview,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar info: %w", err)
	}

	return &AvatarInfo{
		URL:       avatarInfo.URL,
		ID:        avatarInfo.ID,
		Type:      avatarInfo.Type,
		DirectURL: avatarInfo.DirectPath,
	}, nil
}

// GetContacts gets all contacts from the WhatsApp session
func (m *MeowServiceImpl) GetContacts(ctx context.Context, sessionID string) (map[waTypes.JID]waTypes.ContactInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	contacts, err := client.client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	return contacts, nil
}

// GetSubscribedNewsletters gets all newsletters that the session is subscribed to
func (m *MeowServiceImpl) GetSubscribedNewsletters(ctx context.Context, sessionID string) ([]*waTypes.NewsletterMetadata, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	newsletters, err := client.client.GetSubscribedNewsletters()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribed newsletters: %w", err)
	}

	return newsletters, nil
}


func (m *MeowServiceImpl) GetClientJID(sessionID string) (string, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return "", fmt.Errorf("client not found for session %s", sessionID)
	}

	jid := client.GetJID()
	if jid.IsEmpty() {
		return "", fmt.Errorf("client is not authenticated for session %s", sessionID)
	}

	return jid.String(), nil
}


func (m *MeowServiceImpl) RestartClient(sessionID string) error {
	m.logger.Infof("Restarting client for session %s", sessionID)
	
	
	if err := m.StopClient(sessionID); err != nil {
		m.logger.Warnf("Failed to stop client during restart for session %s: %v", sessionID, err)
	}

	
	return m.StartClient(sessionID)
}



func (m *MeowServiceImpl) ConnectOnStartup(ctx context.Context) error {
	m.logger.Infof("Starting automatic reconnection for previously connected sessions...")

	


	query := `
		SELECT id, name, COALESCE(device_jid, '') as device_jid, status, COALESCE(proxy_url, '') as proxy_url
		FROM sessions
		WHERE device_jid IS NOT NULL AND device_jid != ''
		ORDER BY created_at ASC
	`

	rows, err := m.clientManager.db.QueryxContext(ctx, query)
	if err != nil {
		m.logger.Errorf("Failed to query sessions with credentials: %v", err)
		return fmt.Errorf("failed to query sessions with credentials: %w", err)
	}
	defer rows.Close()

	reconnectedCount := 0
	skippedCount := 0
	errorCount := 0
	for rows.Next() {
		var sessionID, name, deviceJID, status, proxyURL string
		err := rows.Scan(&sessionID, &name, &deviceJID, &status, &proxyURL)
		if err != nil {
			m.logger.Errorf("Failed to scan session row: %v", err)
			errorCount++
			continue
		}

		m.logger.Infof("Attempting to reconnect session %s (name: %s, device_jid: %s, current_status: %s)", sessionID, name, deviceJID, status)

		

		hasDevice, err := m.hasDeviceCredentials(ctx, sessionID, deviceJID)
		if err != nil {
			m.logger.Errorf("Failed to check device credentials for session %s: %v", sessionID, err)

			m.updateSessionStatus(ctx, sessionID, "error")
			errorCount++
			continue
		}

		if !hasDevice {
			m.logger.Warnf("No valid device credentials found for session %s, updating status to disconnected", sessionID)

			m.updateSessionStatus(ctx, sessionID, "disconnected")
			skippedCount++
			continue
		}


		if m.clientManager.IsClientConnected(sessionID) {
			m.logger.Infof("Session %s is already connected, skipping", sessionID)

			m.updateSessionStatus(ctx, sessionID, "connected")
			reconnectedCount++
			continue
		}


		m.logger.Infof("Starting client for session %s...", sessionID)
		err = m.StartClient(sessionID)
		if err != nil {
			m.logger.Errorf("Failed to start client for session %s: %v", sessionID, err)

			m.updateSessionStatus(ctx, sessionID, "error")
			errorCount++
			continue
		}

		reconnectedCount++
		m.logger.Infof("Successfully started reconnection for session %s", sessionID)
	}

	if err := rows.Err(); err != nil {
		m.logger.Errorf("Error iterating session rows: %v", err)
		return fmt.Errorf("error iterating session rows: %w", err)
	}

	m.logger.Infof("Automatic reconnection completed. Attempted to reconnect %d sessions", reconnectedCount)
	m.logger.Infof("Reconnection summary - Success: %d, Skipped: %d, Errors: %d", reconnectedCount, skippedCount, errorCount)


	if reconnectedCount > 0 {
		m.logger.Infof("Waiting for %d sessions to establish connections...", reconnectedCount)

		go m.waitForConnectionsToEstablish(ctx, reconnectedCount)
	}

	return nil
}


func (m *MeowServiceImpl) hasDeviceCredentials(ctx context.Context, sessionID, deviceJID string) (bool, error) {

	if deviceJID != "" {
		jid, err := waTypes.ParseJID(deviceJID)
		if err != nil {
			m.logger.Warnf("Invalid device JID %s for session %s: %v", deviceJID, sessionID, err)
		} else {

			device, err := m.clientManager.container.GetDevice(ctx, jid)
			if err == nil && device != nil && device.ID != nil && !device.ID.IsEmpty() {

				identityKeyPair := device.GetIdentityKeyPair()
				if identityKeyPair != nil && identityKeyPair.PrivateKey() != nil {
					m.logger.Debugf("Found valid device credentials for session %s with JID %s", sessionID, deviceJID)
					return true, nil
				} else {
					m.logger.Debugf("Device %s for session %s exists but has no identity key pair", deviceJID, sessionID)
				}
			} else if err != nil {
				m.logger.Debugf("Failed to get device %s for session %s: %v", deviceJID, sessionID, err)
			} else {
				m.logger.Debugf("Device %s for session %s exists but has invalid ID", deviceJID, sessionID)
			}
		}
	}



	devices, err := m.clientManager.container.GetAllDevices(ctx)
	if err != nil {
		m.logger.Errorf("Failed to get all devices for session %s: %v", sessionID, err)
		return false, fmt.Errorf("failed to get devices: %w", err)
	}


	validDeviceCount := 0
	for _, device := range devices {
		if device != nil && device.ID != nil && !device.ID.IsEmpty() {

			identityKeyPair := device.GetIdentityKeyPair()
			if identityKeyPair != nil && identityKeyPair.PrivateKey() != nil {
				validDeviceCount++
				m.logger.Debugf("Found valid device credentials (device %d) for potential use with session %s", validDeviceCount, sessionID)
			}
		}
	}

	if validDeviceCount > 0 {
		m.logger.Debugf("Found %d valid device(s) that could be used for session %s", validDeviceCount, sessionID)
		return true, nil
	}

	m.logger.Debugf("No valid device credentials found for session %s", sessionID)
	return false, nil
}


func (m *MeowServiceImpl) waitForConnectionsToEstablish(ctx context.Context, expectedConnections int) {

	maxWaitTime := 30
	checkInterval := 2

	for i := 0; i < maxWaitTime; i += checkInterval {
		time.Sleep(time.Duration(checkInterval) * time.Second)


		connectedCount := 0
		query := `SELECT id FROM sessions WHERE status = 'connecting' OR status = 'connected'`
		rows, err := m.clientManager.db.QueryxContext(ctx, query)
		if err != nil {
			m.logger.Errorf("Failed to check session statuses: %v", err)
			return
		}

		for rows.Next() {
			var sessionID string
			if err := rows.Scan(&sessionID); err != nil {
				continue
			}


			if m.clientManager.IsClientConnected(sessionID) {
				connectedCount++

				m.updateSessionStatus(ctx, sessionID, "connected")
			}
		}
		rows.Close()

		m.logger.Infof("Connection progress: %d/%d sessions connected after %d seconds", connectedCount, expectedConnections, i+checkInterval)


		if connectedCount >= expectedConnections {
			m.logger.Infof("All %d sessions successfully connected!", connectedCount)
			return
		}
	}

	m.logger.Warnf("Connection timeout reached. Some sessions may still be connecting.")
}

func (m *MeowServiceImpl) updateSessionStatus(ctx context.Context, sessionID, status string) {
	query := `UPDATE sessions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := m.clientManager.db.ExecContext(ctx, query, status, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to update session %s status to %s: %v", sessionID, status, err)
	} else {
		m.logger.Debugf("Updated session %s status to %s", sessionID, status)
	}
}


