package meow

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// MeowClient wraps a whatsmeow.Client with additional functionality
type MeowClient struct {
	sessionID      string
	client         *whatsmeow.Client
	eventHandler   *EventHandler
	manager        *ClientManager
	logger         logger.Logger
	waLogger       waLog.Logger

	// State management
	mu           sync.RWMutex
	status       types.Status
	lastActivity time.Time
	qrCode       string

	// Event handler ID for cleanup
	eventHandlerID uint32

	// Reconnection management (like wuzapi)
	killChannel chan bool
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewMeowClient creates a new MeowClient
func NewMeowClient(sessionID string, deviceStore *store.Device, waLogger waLog.Logger, manager *ClientManager) (*MeowClient, error) {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("meow-client").Sub(sessionID)

	// Create whatsmeow client
	waClient := whatsmeow.NewClient(deviceStore, waLogger)

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Create MeowClient
	meowClient := &MeowClient{
		sessionID:    sessionID,
		client:       waClient,
		manager:      manager,
		logger:       appLogger,
		waLogger:     waLogger,
		status:       types.StatusDisconnected,
		lastActivity: time.Now(),
		killChannel:  make(chan bool, 1),
		ctx:          ctx,
		cancel:       cancel,
	}

	// Create and set up event handler
	eventHandler := NewEventHandler(sessionID, waLogger, meowClient)
	meowClient.eventHandler = eventHandler

	// Register event handler with whatsmeow client
	meowClient.eventHandlerID = waClient.AddEventHandler(eventHandler.HandleEvent)

	return meowClient, nil
}

// Connect connects the client to WhatsApp with auto-reconnection (like wuzapi)
func (mc *MeowClient) Connect(ctx context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.client.IsConnected() {
		return nil // Already connected
	}

	mc.setStatus(types.StatusConnecting)
	mc.logger.Infof("Connecting client for session %s", mc.sessionID)

	// Start the connection process in a goroutine (like wuzapi)
	go mc.startClientLoop()

	return nil
}

// startClientLoop implements the main connection loop (like wuzapi)
func (mc *MeowClient) startClientLoop() {
	mc.logger.Infof("Starting client loop for session %s", mc.sessionID)

	// Handle QR code generation if device is not registered (BEFORE connecting)
	if mc.client.Store.ID == nil {
		mc.logger.Infof("Device not registered for session %s, waiting for QR code", mc.sessionID)

		// Get QR channel BEFORE connecting (like wuzapi)
		qrChan, err := mc.client.GetQRChannel(context.Background())
		if err != nil {
			mc.logger.Errorf("Failed to get QR channel for session %s: %v", mc.sessionID, err)
			return
		}

		// Connect after getting QR channel
		err = mc.client.Connect()
		if err != nil {
			mc.logger.Errorf("Failed to connect client for session %s: %v", mc.sessionID, err)
			mc.setStatus(types.StatusError)
			return
		}

		// Process QR codes in a loop (like wuzapi)
		go mc.handleQRLoop(qrChan)
	} else {
		mc.logger.Infof("Already logged in, just connecting for session %s", mc.sessionID)

		// Connect for already logged in devices
		err := mc.client.Connect()
		if err != nil {
			mc.logger.Errorf("Failed to connect client for session %s: %v", mc.sessionID, err)
			mc.setStatus(types.StatusError)
			return
		}

		mc.setStatus(types.StatusConnected)
	}

	// Keep client alive with reconnection loop (like wuzapi)
	mc.keepAliveLoop()
}

// handleQRLoop processes QR codes continuously (like wuzapi)
func (mc *MeowClient) handleQRLoop(qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			mc.logger.Infof("QR code received for session %s", mc.sessionID)

			// Display QR code in terminal (like wuzapi does) - EVERY TIME
			fmt.Printf("\n=== NEW QR CODE FOR SESSION %s ===\n", mc.sessionID)
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Printf("QR code text: %s\n", evt.Code)
			fmt.Printf("=== Scan with WhatsApp to connect ===\n\n")

			mc.setQRCode(evt.Code)
		case "timeout":
			mc.logger.Warnf("QR code timeout for session %s", mc.sessionID)
			fmt.Printf("\n=== QR CODE TIMEOUT FOR SESSION %s ===\n", mc.sessionID)
			fmt.Printf("=== Generating new QR code... ===\n\n")
			mc.clearQRCode()
			// Don't kill on timeout, let it generate a new QR code
		case "success":
			mc.logger.Infof("QR pairing successful for session %s", mc.sessionID)
			fmt.Printf("\n=== QR PAIRING SUCCESSFUL FOR SESSION %s ===\n", mc.sessionID)
			fmt.Printf("=== WhatsApp connected successfully! ===\n\n")
			mc.clearQRCode()
			mc.setStatus(types.StatusConnected)
		default:
			mc.logger.Infof("QR event for session %s: %s", mc.sessionID, evt.Event)
		}
	}
}

// keepAliveLoop maintains the connection (like wuzapi)
func (mc *MeowClient) keepAliveLoop() {
	for {
		select {
		case <-mc.killChannel:
			mc.logger.Infof("Received kill signal for session %s", mc.sessionID)
			mc.client.Disconnect()
			mc.setStatus(types.StatusDisconnected)
			return
		case <-mc.ctx.Done():
			mc.logger.Infof("Context cancelled for session %s", mc.sessionID)
			mc.client.Disconnect()
			mc.setStatus(types.StatusDisconnected)
			return
		default:
			time.Sleep(1000 * time.Millisecond)
			// Keep the loop running to maintain connection
		}
	}
}

// Disconnect disconnects the client from WhatsApp
func (mc *MeowClient) Disconnect() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if !mc.client.IsConnected() {
		return
	}

	mc.logger.Infof("Disconnecting client for session %s", mc.sessionID)

	// Send kill signal to stop the keep-alive loop
	select {
	case mc.killChannel <- true:
	default:
		// Channel might be full, that's ok
	}

	mc.client.Disconnect()
	mc.setStatus(types.StatusDisconnected)
}

// Stop stops the client and cancels all operations
func (mc *MeowClient) Stop() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.logger.Infof("Stopping client for session %s", mc.sessionID)

	// Cancel context to stop all goroutines
	mc.cancel()

	// Send kill signal
	select {
	case mc.killChannel <- true:
	default:
		// Channel might be full, that's ok
	}

	if mc.client.IsConnected() {
		mc.client.Disconnect()
	}

	mc.setStatus(types.StatusDisconnected)
}

// Logout logs out the client and clears device data
func (mc *MeowClient) Logout(ctx context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.logger.Infof("Logging out client for session %s", mc.sessionID)

	if mc.client.IsConnected() {
		err := mc.client.Logout(ctx)
		if err != nil {
			mc.logger.Errorf("Failed to logout client %s: %v", mc.sessionID, err)
			return err
		}
	}

	mc.setStatus(types.StatusDisconnected)
	return nil
}

// GetQRCode gets the QR code for authentication
func (mc *MeowClient) GetQRCode(ctx context.Context) (string, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.client.IsConnected() && mc.client.Store.ID != nil {
		return "", fmt.Errorf("client is already authenticated")
	}

	// If we already have a QR code, return it
	if mc.qrCode != "" {
		return mc.qrCode, nil
	}

	// If not connected, we need to connect first to get QR code
	if !mc.client.IsConnected() {
		return "", fmt.Errorf("client is not connected, call Connect() first")
	}

	// QR code will be set by the event handler
	return "", fmt.Errorf("QR code not yet available, please wait")
}

// PairPhone initiates phone pairing
func (mc *MeowClient) PairPhone(ctx context.Context, phoneNumber string) (string, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.client.IsConnected() && mc.client.Store.ID != nil {
		return "", fmt.Errorf("client is already authenticated")
	}

	if !mc.client.IsConnected() {
		return "", fmt.Errorf("client is not connected, call Connect() first")
	}

	// Request pairing code
	code, err := mc.client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", fmt.Errorf("failed to request pairing code: %w", err)
	}

	mc.logger.Infof("Pairing code requested for session %s, phone %s", mc.sessionID, phoneNumber)
	return code, nil
}

// IsConnected checks if the client is connected
func (mc *MeowClient) IsConnected() bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.client.IsConnected()
}

// GetStatus returns the current status
func (mc *MeowClient) GetStatus() types.Status {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.status
}

// GetLastActivity returns the last activity time
func (mc *MeowClient) GetLastActivity() time.Time {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.lastActivity
}

// GetJID returns the client's JID if authenticated
func (mc *MeowClient) GetJID() waTypes.JID {
	if mc.client.Store.ID != nil {
		return *mc.client.Store.ID
	}
	return waTypes.JID{}
}

// SendMessage sends a text message
func (mc *MeowClient) SendMessage(ctx context.Context, to waTypes.JID, message string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Create a simple text message
	msg := &waE2E.Message{
		Conversation: &message,
	}

	_, err := mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// GetClient returns the underlying whatsmeow client (use with caution)
func (mc *MeowClient) GetClient() *whatsmeow.Client {
	return mc.client
}

// SendTextMessage sends a text message with optional context info
func (mc *MeowClient) SendTextMessage(ctx context.Context, to waTypes.JID, text string, contextInfo *waE2E.ContextInfo) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &text,
		},
	}

	if contextInfo != nil {
		msg.ExtendedTextMessage.ContextInfo = contextInfo
	}

	_, err := mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send text message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// SendLocationMessage sends a location message
func (mc *MeowClient) SendLocationMessage(ctx context.Context, to waTypes.JID, latitude, longitude float64, name, address string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  &latitude,
			DegreesLongitude: &longitude,
		},
	}

	if name != "" {
		msg.LocationMessage.Name = &name
	}
	if address != "" {
		msg.LocationMessage.Address = &address
	}

	_, err := mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send location message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// SendContactMessage sends a contact message
func (mc *MeowClient) SendContactMessage(ctx context.Context, to waTypes.JID, displayName, vcard string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	msg := &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: &displayName,
			Vcard:       &vcard,
		},
	}

	_, err := mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send contact message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// ReactToMessage sends a reaction to a message
func (mc *MeowClient) ReactToMessage(ctx context.Context, chatJID waTypes.JID, messageID string, emoji string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// TODO: MessageKey was removed in new whatsmeow version
	// Need to implement using new reaction API
	return fmt.Errorf("message reactions not implemented for new whatsmeow version")
}

// SetChatPresence sets the presence in a chat
func (mc *MeowClient) SetChatPresence(ctx context.Context, chatJID waTypes.JID, presence waTypes.Presence) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	err := mc.client.SendPresence(presence)
	if err != nil {
		return fmt.Errorf("failed to set presence: %w", err)
	}

	mc.updateActivity()
	return nil
}

// MarkMessageRead marks messages as read
func (mc *MeowClient) MarkMessageRead(ctx context.Context, chatJID waTypes.JID, messageIDs []string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Convert message IDs to MessageIDs
	msgIDList := make([]waTypes.MessageID, len(messageIDs))
	for i, msgID := range messageIDs {
		msgIDList[i] = waTypes.MessageID(msgID)
	}

	err := mc.client.MarkRead(msgIDList, time.Now(), chatJID, chatJID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	mc.updateActivity()
	return nil
}

// CreateGroup creates a new WhatsApp group
func (mc *MeowClient) CreateGroup(ctx context.Context, name string, participants []waTypes.JID) (*waTypes.GroupInfo, error) {
	if !mc.IsConnected() {
		return nil, fmt.Errorf("client is not connected")
	}

	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participants,
	}

	groupInfo, err := mc.client.CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	mc.updateActivity()
	return groupInfo, nil
}

// GetGroupInfo retrieves information about a group
func (mc *MeowClient) GetGroupInfo(ctx context.Context, groupJID waTypes.JID) (*waTypes.GroupInfo, error) {
	if !mc.IsConnected() {
		return nil, fmt.Errorf("client is not connected")
	}

	groupInfo, err := mc.client.GetGroupInfo(groupJID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

	mc.updateActivity()
	return groupInfo, nil
}

// JoinGroupWithLink joins a group using an invite link
func (mc *MeowClient) JoinGroupWithLink(ctx context.Context, inviteCode string) (*waTypes.GroupInfo, error) {
	if !mc.IsConnected() {
		return nil, fmt.Errorf("client is not connected")
	}

	groupJID, err := mc.client.JoinGroupWithLink(inviteCode)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	// Get group info after joining
	groupInfo, err := mc.client.GetGroupInfo(groupJID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info after joining: %w", err)
	}

	mc.updateActivity()
	return groupInfo, nil
}

// LeaveGroup leaves a group
func (mc *MeowClient) LeaveGroup(ctx context.Context, groupJID waTypes.JID) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	err := mc.client.LeaveGroup(groupJID)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	mc.updateActivity()
	return nil
}

// GetGroupInviteLink gets the invite link for a group
func (mc *MeowClient) GetGroupInviteLink(ctx context.Context, groupJID waTypes.JID, reset bool) (string, error) {
	if !mc.IsConnected() {
		return "", fmt.Errorf("client is not connected")
	}

	link, err := mc.client.GetGroupInviteLink(groupJID, reset)
	if err != nil {
		return "", fmt.Errorf("failed to get group invite link: %w", err)
	}

	mc.updateActivity()
	return link, nil
}

// UpdateGroupParticipants updates group participants (add, remove, promote, demote)
// TODO: ParticipantChange type was removed in new whatsmeow version
func (mc *MeowClient) UpdateGroupParticipants(ctx context.Context, groupJID waTypes.JID, participants []waTypes.JID, action string) error {
	return fmt.Errorf("group participant management not implemented for new whatsmeow version")
}

// SetGroupName sets the name of a group
func (mc *MeowClient) SetGroupName(ctx context.Context, groupJID waTypes.JID, name string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	err := mc.client.SetGroupName(groupJID, name)
	if err != nil {
		return fmt.Errorf("failed to set group name: %w", err)
	}

	mc.updateActivity()
	return nil
}

// SetGroupTopic sets the topic/description of a group
func (mc *MeowClient) SetGroupTopic(ctx context.Context, groupJID waTypes.JID, topic string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	err := mc.client.SetGroupTopic(groupJID, "", topic, "")
	if err != nil {
		return fmt.Errorf("failed to set group topic: %w", err)
	}

	mc.updateActivity()
	return nil
}

// Cleanup cleans up resources
func (mc *MeowClient) Cleanup() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Remove event handler
	if mc.eventHandlerID != 0 {
		mc.client.RemoveEventHandler(mc.eventHandlerID)
	}

	// Disconnect if connected
	if mc.client.IsConnected() {
		mc.client.Disconnect()
	}

	mc.logger.Infof("Cleaned up client for session %s", mc.sessionID)
}

// setStatus sets the status and notifies the manager (internal method)
func (mc *MeowClient) setStatus(status types.Status) {
	mc.status = status
	mc.lastActivity = time.Now()
	
	// Notify manager of status change
	if mc.manager != nil {
		go mc.manager.OnClientStatusChange(mc.sessionID, status)
	}
}

// updateActivity updates the last activity time
func (mc *MeowClient) updateActivity() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.lastActivity = time.Now()
}

// setQRCode sets the QR code (called by event handler)
func (mc *MeowClient) setQRCode(qrCode string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.qrCode = qrCode
	mc.logger.Infof("QR code updated for session %s", mc.sessionID)
}

// clearQRCode clears the QR code (called by event handler)
func (mc *MeowClient) clearQRCode() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.qrCode = ""
}

// onConnected is called when the client connects successfully
func (mc *MeowClient) onConnected() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.setStatus(types.StatusConnected)
	mc.clearQRCode()
	mc.logger.Infof("Client connected for session %s", mc.sessionID)
}

// onDisconnected is called when the client disconnects
func (mc *MeowClient) onDisconnected() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.setStatus(types.StatusDisconnected)
	mc.logger.Infof("Client disconnected for session %s", mc.sessionID)
}

// onError is called when an error occurs
func (mc *MeowClient) onError(err error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.setStatus(types.StatusError)
	mc.logger.Errorf("Client error for session %s: %v", mc.sessionID, err)
}

// GetSessionID returns the session ID
func (mc *MeowClient) GetSessionID() string {
	return mc.sessionID
}

// SendImageMessage sends an image message
func (mc *MeowClient) SendImageMessage(ctx context.Context, to waTypes.JID, imageData []byte, caption string, mimeType string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Upload image
	uploaded, err := mc.client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
		},
	}

	if caption != "" {
		msg.ImageMessage.Caption = &caption
	}

	_, err = mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send image message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// SendAudioMessage sends an audio message
func (mc *MeowClient) SendAudioMessage(ctx context.Context, to waTypes.JID, audioData []byte, mimeType string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Upload audio
	uploaded, err := mc.client.Upload(ctx, audioData, whatsmeow.MediaAudio)
	if err != nil {
		return fmt.Errorf("failed to upload audio: %w", err)
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
		},
	}

	_, err = mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send audio message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// SendDocumentMessage sends a document message
func (mc *MeowClient) SendDocumentMessage(ctx context.Context, to waTypes.JID, documentData []byte, filename, caption, mimeType string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Upload document
	uploaded, err := mc.client.Upload(ctx, documentData, whatsmeow.MediaDocument)
	if err != nil {
		return fmt.Errorf("failed to upload document: %w", err)
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
			FileName:      &filename,
		},
	}

	if caption != "" {
		msg.DocumentMessage.Caption = &caption
	}

	_, err = mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send document message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// SendVideoMessage sends a video message
func (mc *MeowClient) SendVideoMessage(ctx context.Context, to waTypes.JID, videoData []byte, caption, mimeType string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Upload video
	uploaded, err := mc.client.Upload(ctx, videoData, whatsmeow.MediaVideo)
	if err != nil {
		return fmt.Errorf("failed to upload video: %w", err)
	}

	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
		},
	}

	if caption != "" {
		msg.VideoMessage.Caption = &caption
	}

	_, err = mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send video message: %w", err)
	}

	mc.updateActivity()
	return nil
}

// SendStickerMessage sends a sticker message
func (mc *MeowClient) SendStickerMessage(ctx context.Context, to waTypes.JID, stickerData []byte, mimeType string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Upload sticker
	uploaded, err := mc.client.Upload(ctx, stickerData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("failed to upload sticker: %w", err)
	}

	msg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
		},
	}

	_, err = mc.client.SendMessage(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send sticker message: %w", err)
	}

	mc.updateActivity()
	return nil
}
