package meow

import (
	"context"
	"fmt"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"

	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// MediaDownloadInfo contains the information needed to download media
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


// parseJID parses a phone number or JID string into a WhatsApp JID
// Deprecated: Use JID.ParseJID instead
func parseJID(arg string) (waTypes.JID, bool) {
	jid, err := JID.ParseJID(arg)
	return jid, err == nil
}

// validateAndGetClient validates session and returns client with parsed JID - eliminates code duplication
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


func NewMeowService(db *sqlx.DB, container *sqlstore.Container, waLogger waLog.Logger) session.WhatsAppService {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("meow-service")
	clientManager := NewClientManager(db, container, waLogger)

	service := &MeowServiceImpl{
		clientManager: clientManager,
		logger:        appLogger,
		waLogger:      waLogger,
	}

	return service
}


func (m *MeowServiceImpl) StartClient(sessionID string) error {
	// Validate session ID
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
	// Validate session ID
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

	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		m.logger.Errorf("DEBUG: Validation failed: %v", err)
		return nil, err
	}

	m.logger.Infof("DEBUG: Client found and JID parsed successfully: %s -> %s", to, jid.String())

	// Skip IsConnected() check to avoid deadlock
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
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendLocationMessage(ctx, jid, latitude, longitude, name, address)
}


func (m *MeowServiceImpl) SendContactMessage(ctx context.Context, sessionID, to, displayName, vcard string) (*whatsmeow.SendResponse, error) {
	// Use utility function to validate and get client - eliminates code duplication
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

	// Parse chat JID - use our JID utility for better parsing
	jid, err := JID.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid JID %s: %w", chatJID, err)
	}

	// Validate presence parameters
	if err := m.validatePresenceParams(state, media); err != nil {
		return err
	}

	// Convert to whatsmeow types
	chatPresence := waTypes.ChatPresence(state)
	mediaType := waTypes.ChatPresenceMedia(media)

	// Send chat presence using whatsmeow client
	err = client.client.SendChatPresence(jid, chatPresence, mediaType)
	if err != nil {
		return fmt.Errorf("failed to set chat presence: %w", err)
	}

	return nil
}

// validatePresenceParams validates chat presence parameters
func (m *MeowServiceImpl) validatePresenceParams(state, media string) error {
	// Validate state - only composing and paused are supported by whatsmeow
	if state != "composing" && state != "paused" {
		return fmt.Errorf("invalid state: %s (valid: composing, paused)", state)
	}

	// Validate media - empty string (ChatPresenceMediaText), "audio" are valid
	if media != "" && media != "audio" {
		return fmt.Errorf("invalid media: %s (valid: \"\", audio)", media)
	}

	return nil
}


func (m *MeowServiceImpl) MarkMessageRead(ctx context.Context, sessionID, chatJID string, messageIDs []string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid JID %s: %w", chatJID, err)
	}

	return client.MarkMessageRead(ctx, jid, messageIDs)
}


func (m *MeowServiceImpl) CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	
	participantJIDs := make([]waTypes.JID, len(participants))
	for i, participant := range participants {
		jid, err := waTypes.ParseJID(participant)
		if err != nil {
			return nil, fmt.Errorf("invalid participant JID %s: %w", participant, err)
		}
		participantJIDs[i] = jid
	}

	return client.CreateGroup(ctx, name, participantJIDs)
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
	_, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	
	_, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	
	participantJIDs := make([]waTypes.JID, len(participants))
	for i, participant := range participants {
		pJid, err := waTypes.ParseJID(participant)
		if err != nil {
			return fmt.Errorf("invalid participant JID %s: %w", participant, err)
		}
		participantJIDs[i] = pJid
	}

	
	
	return fmt.Errorf("group participant management not implemented for new whatsmeow version")
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
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendImageMessage(ctx, jid, imageData, caption, mimeType)
}


func (m *MeowServiceImpl) SendAudioMessage(ctx context.Context, sessionID, to string, audioData []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendAudioMessage(ctx, jid, audioData, mimeType)
}


func (m *MeowServiceImpl) SendDocumentMessage(ctx context.Context, sessionID, to string, documentData []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendDocumentMessage(ctx, jid, documentData, filename, caption, mimeType)
}


func (m *MeowServiceImpl) SendVideoMessage(ctx context.Context, sessionID, to string, videoData []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendVideoMessage(ctx, jid, videoData, caption, mimeType)
}


func (m *MeowServiceImpl) SendStickerMessage(ctx context.Context, sessionID, to string, stickerData []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendStickerMessage(ctx, jid, stickerData, mimeType)
}


func (m *MeowServiceImpl) SendButtonsMessage(ctx context.Context, sessionID, to, text string, buttons []types.Button, footer string) (*whatsmeow.SendResponse, error) {
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendButtonsMessage(ctx, jid, text, buttons, footer)
}


func (m *MeowServiceImpl) SendListMessage(ctx context.Context, sessionID, to, text, buttonText string, sections []types.Section, footer string) (*whatsmeow.SendResponse, error) {
	// Use utility function to validate and get client - eliminates code duplication
	client, jid, err := m.validateAndGetClient(sessionID, to)
	if err != nil {
		return nil, err
	}

	return client.SendListMessage(ctx, jid, text, buttonText, sections, footer)
}

// ============================================================================
// Chat Operations Implementation
// ============================================================================

// DeleteMessage deletes a message from a chat
func (m *MeowServiceImpl) DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error {
	m.logger.Infof("Deleting message %s in chat %s for session %s (forEveryone: %v)", messageID, chatJID, sessionID, forEveryone)

	// Validate inputs
	if err := m.validateChatOperation(sessionID, chatJID, messageID); err != nil {
		return err
	}

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Check if client is connected
	if !client.IsConnected() {
		return fmt.Errorf("client for session %s is not connected", sessionID)
	}

	// Parse chat JID
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID %s: %w", chatJID, err)
	}

	// Use whatsmeow's BuildRevoke to create revoke message
	revokeMsg := client.client.BuildRevoke(jid, waTypes.EmptyJID, waTypes.MessageID(messageID))

	// Send the revoke message with timeout
	_, err = client.client.SendMessage(ctx, jid, revokeMsg)
	if err != nil {
		return fmt.Errorf("failed to delete message %s: %w", messageID, err)
	}

	m.logger.Infof("Successfully deleted message %s", messageID)
	return nil
}

// validateChatOperation validates common parameters for chat operations
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

	// Validate message ID format (basic check)
	if len(messageID) < 10 {
		return fmt.Errorf("invalid message ID format")
	}

	return nil
}

// EditMessage edits a text message in a chat
func (m *MeowServiceImpl) EditMessage(ctx context.Context, sessionID, chatJID, messageID, newText string) (*types.SendResponse, error) {
	m.logger.Infof("Editing message %s in chat %s for session %s", messageID, chatJID, sessionID)

	// Validate inputs
	if err := m.validateChatOperation(sessionID, chatJID, messageID); err != nil {
		return nil, err
	}

	// Validate new text
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

	// Check if client is connected
	if !client.IsConnected() {
		return nil, fmt.Errorf("client for session %s is not connected", sessionID)
	}

	// Parse chat JID
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat JID %s: %w", chatJID, err)
	}

	// Create the new text message
	textMsg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &newText,
		},
	}

	// Use whatsmeow's BuildEdit to create edit message
	editMsg := client.client.BuildEdit(jid, waTypes.MessageID(messageID), textMsg)

	// Send the edit message
	resp, err := client.client.SendMessage(ctx, jid, editMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to edit message %s: %w", messageID, err)
	}

	m.logger.Infof("Successfully edited message %s", messageID)

	// Convert whatsmeow response to our types.SendResponse
	response := types.NewSendResponseFromWhatsmeow(&resp, messageID)
	return &response, nil
}

// DownloadMedia downloads media from a message and returns the data and mime type
func (m *MeowServiceImpl) DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error) {
	m.logger.Infof("Downloading media for message %s in session %s", messageID, sessionID)

	// Note: This is a simplified implementation. In a real-world scenario,
	// you would need to store message metadata (URL, DirectPath, MediaKey, etc.)
	// when messages are received and retrieve them here using the messageID.
	// For now, this returns an error indicating the limitation.

	return nil, "", fmt.Errorf("media download requires message metadata storage - not implemented in this simplified version")
}

// DownloadMediaWithInfo downloads media using provided media information
func (m *MeowServiceImpl) DownloadMediaWithInfo(ctx context.Context, sessionID string, mediaInfo MediaDownloadInfo) ([]byte, string, error) {
	m.logger.Infof("Downloading media with info for session %s", sessionID)

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, "", fmt.Errorf("client not found for session %s", sessionID)
	}

	var mediaData []byte
	var err error
	var mimeType string

	// Create appropriate message based on media type
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

// ReactToMessage sends a reaction to a message
func (m *MeowServiceImpl) ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string) error {
	m.logger.Infof("Reacting to message %s in chat %s for session %s with emoji: %s", messageID, chatJID, sessionID, emoji)

	// Validate inputs
	if err := m.validateChatOperation(sessionID, chatJID, messageID); err != nil {
		return err
	}

	// Validate emoji (basic check)
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

	// Check if client is connected
	if !client.IsConnected() {
		return fmt.Errorf("client for session %s is not connected", sessionID)
	}

	// Parse chat JID
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID %s: %w", chatJID, err)
	}

	// Use whatsmeow's BuildReaction to create reaction message
	reactionMsg := client.client.BuildReaction(jid, waTypes.EmptyJID, waTypes.MessageID(messageID), emoji)

	// Send the reaction message
	_, err = client.client.SendMessage(ctx, jid, reactionMsg)
	if err != nil {
		return fmt.Errorf("failed to react to message %s: %w", messageID, err)
	}

	m.logger.Infof("Successfully reacted to message %s with %s", messageID, emoji)
	return nil
}


func (m *MeowServiceImpl) SendPollMessage(ctx context.Context, sessionID, to, name string, options []string, selectableCount int) (*whatsmeow.SendResponse, error) {
	m.logger.Infof("DEBUG: SendPollMessage called - sessionID: %s, to: %s, name: %s", sessionID, to, name)

	// Use utility function to validate and get client - eliminates code duplication
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

	
	// Query ALL sessions that have device_jid (indicating they were previously authenticated)
	// We don't filter by status because sessions might be in various states after server restart
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

		
		// Check if session has valid device credentials in the store
		hasDevice, err := m.hasDeviceCredentials(ctx, sessionID, deviceJID)
		if err != nil {
			m.logger.Errorf("Failed to check device credentials for session %s: %v", sessionID, err)
			// Update status to error
			m.updateSessionStatus(ctx, sessionID, "error")
			errorCount++
			continue
		}

		if !hasDevice {
			m.logger.Warnf("No valid device credentials found for session %s, updating status to disconnected", sessionID)
			// Update status to disconnected since no valid credentials
			m.updateSessionStatus(ctx, sessionID, "disconnected")
			skippedCount++
			continue
		}

		// Check if client is already running (avoid duplicate connections)
		if m.clientManager.IsClientConnected(sessionID) {
			m.logger.Infof("Session %s is already connected, skipping", sessionID)
			// Ensure status is correct
			m.updateSessionStatus(ctx, sessionID, "connected")
			reconnectedCount++
			continue
		}

		// Start the client with proper error handling
		m.logger.Infof("Starting client for session %s...", sessionID)
		err = m.StartClient(sessionID)
		if err != nil {
			m.logger.Errorf("Failed to start client for session %s: %v", sessionID, err)
			// Update status to error
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

	// Wait a bit for connections to establish
	if reconnectedCount > 0 {
		m.logger.Infof("Waiting for %d sessions to establish connections...", reconnectedCount)
		// Give sessions time to connect (non-blocking)
		go m.waitForConnectionsToEstablish(ctx, reconnectedCount)
	}

	return nil
}


func (m *MeowServiceImpl) hasDeviceCredentials(ctx context.Context, sessionID, deviceJID string) (bool, error) {
	// First, try to get device by specific JID if provided
	if deviceJID != "" {
		jid, err := waTypes.ParseJID(deviceJID)
		if err != nil {
			m.logger.Warnf("Invalid device JID %s for session %s: %v", deviceJID, sessionID, err)
		} else {
			// Try to get the specific device
			device, err := m.clientManager.container.GetDevice(ctx, jid)
			if err == nil && device != nil && device.ID != nil && !device.ID.IsEmpty() {
				// Check if device has identity key pair (indicates it's authenticated)
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

	// If specific JID didn't work, check if there are any valid devices in the store
	// This is a fallback for cases where device_jid might be outdated
	devices, err := m.clientManager.container.GetAllDevices(ctx)
	if err != nil {
		m.logger.Errorf("Failed to get all devices for session %s: %v", sessionID, err)
		return false, fmt.Errorf("failed to get devices: %w", err)
	}

	// Look for any device with valid credentials
	validDeviceCount := 0
	for _, device := range devices {
		if device != nil && device.ID != nil && !device.ID.IsEmpty() {
			// Check if device has identity key pair (indicates it's authenticated)
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

// waitForConnectionsToEstablish waits for sessions to establish connections and updates their status
func (m *MeowServiceImpl) waitForConnectionsToEstablish(ctx context.Context, expectedConnections int) {
	// Wait up to 30 seconds for connections to establish
	maxWaitTime := 30
	checkInterval := 2

	for i := 0; i < maxWaitTime; i += checkInterval {
		time.Sleep(time.Duration(checkInterval) * time.Second)

		// Check how many sessions are actually connected
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

			// Check if this session is actually connected
			if m.clientManager.IsClientConnected(sessionID) {
				connectedCount++
				// Update status to connected if it's not already
				m.updateSessionStatus(ctx, sessionID, "connected")
			}
		}
		rows.Close()

		m.logger.Infof("Connection progress: %d/%d sessions connected after %d seconds", connectedCount, expectedConnections, i+checkInterval)

		// If all expected connections are established, we're done
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


