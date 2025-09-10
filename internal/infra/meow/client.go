package meow

import (
	"context"
	"errors"
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


type MeowClient struct {
	sessionID    string
	client       *whatsmeow.Client
	eventHandler *EventHandler
	manager      *ClientManager
	logger       logger.Logger
	waLogger     waLog.Logger

	
	mu           sync.RWMutex
	status       types.Status
	lastActivity time.Time
	qrCode       string

	
	eventHandlerID uint32

	
	killChannel chan bool
	ctx         context.Context
	cancel      context.CancelFunc

	
	qrStopChannel chan bool

	
	qrLoopActive bool
	qrLoopCancel context.CancelFunc
}


func NewMeowClient(sessionID string, deviceStore *store.Device, waLogger waLog.Logger, manager *ClientManager) (*MeowClient, error) {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("meow-client").Sub(sessionID)

	
	waClient := whatsmeow.NewClient(deviceStore, waLogger)

	
	ctx, cancel := context.WithCancel(context.Background())

	
	meowClient := &MeowClient{
		sessionID:     sessionID,
		client:        waClient,
		manager:       manager,
		logger:        appLogger,
		waLogger:      waLogger,
		status:        types.StatusDisconnected,
		lastActivity:  time.Now(),
		killChannel:   make(chan bool, 1),
		qrStopChannel: make(chan bool, 1),
		ctx:           ctx,
		cancel:        cancel,
	}

	
	eventHandler := NewEventHandler(sessionID, waLogger, meowClient)
	meowClient.eventHandler = eventHandler

	
	meowClient.eventHandlerID = waClient.AddEventHandler(eventHandler.HandleEvent)

	return meowClient, nil
}


func (mc *MeowClient) Connect(ctx context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Early return if already connected
	if mc.client.IsConnected() {
		return nil
	}

	if err := Validation.ValidateSessionID(mc.sessionID); err != nil {
		return Error.WrapError(err, "connect validation failed")
	}

	mc.setStatus(types.StatusConnecting)
	mc.logger.Infof("Connecting client for session %s", mc.sessionID)

	// Start the client loop in a goroutine
	go mc.startClientLoop()

	return nil
}


func (mc *MeowClient) startClientLoop() {
	mc.logger.Infof("Starting client loop for session %s", mc.sessionID)

	
	if mc.client.Store.ID == nil {
		mc.logger.Infof("Device not registered for session %s, waiting for QR code", mc.sessionID)

		
		qrChan, err := mc.client.GetQRChannel(context.Background())
		if err != nil {
			mc.logger.Errorf("Failed to get QR channel for session %s: %v", mc.sessionID, err)
			return
		}

		
		err = mc.client.Connect()
		if err != nil {
			mc.logger.Errorf("Failed to connect client for session %s: %v", mc.sessionID, err)
			mc.setStatus(types.StatusError)
			return
		}

		
		go mc.handleQRLoop(qrChan)
	} else {
		mc.logger.Infof("Already logged in, just connecting for session %s", mc.sessionID)

		
		err := mc.client.Connect()
		if err != nil {
			mc.logger.Errorf("Failed to connect client for session %s: %v", mc.sessionID, err)
			mc.setStatus(types.StatusError)
			return
		}

		mc.setStatus(types.StatusConnected)
	}

	
	mc.keepAliveLoop()
}


func (mc *MeowClient) handleQRLoop(qrChan <-chan whatsmeow.QRChannelItem) {
	mc.logger.Infof("Starting QR loop for session %s", mc.sessionID)

	
	qrCtx, qrCancel := context.WithCancel(mc.ctx)
	mc.mu.Lock()
	mc.qrLoopActive = true
	mc.qrLoopCancel = qrCancel
	mc.mu.Unlock()

	defer func() {
		mc.mu.Lock()
		mc.qrLoopActive = false
		mc.qrLoopCancel = nil
		mc.mu.Unlock()
		qrCancel()
		mc.logger.Infof("QR loop finished for session %s", mc.sessionID)
	}()

	for {
		select {
		case evt, ok := <-qrChan:
			if !ok {
				mc.logger.Infof("QR channel closed for session %s", mc.sessionID)
				return
			}

			switch evt.Event {
			case "code":
				mc.logger.Infof("QR code received for session %s", mc.sessionID)

				
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

				
				if mc.manager != nil {
					go mc.manager.clearQRCodeInDatabase(mc.sessionID)
				}

			case "success":
				mc.logger.Infof("QR pairing successful for session %s", mc.sessionID)
				fmt.Printf("\n=== QR PAIRING SUCCESSFUL FOR SESSION %s ===\n", mc.sessionID)
				fmt.Printf("=== WhatsApp connected successfully! ===\n\n")

				
				mc.clearQRCode()
				mc.setStatus(types.StatusConnected)

				
				if mc.client.Store.ID != nil && mc.manager != nil {
					mc.manager.OnPairSuccess(mc.sessionID, mc.client.Store.ID.String())
				}

				
				mc.logger.Infof("QR pairing successful, stopping QR loop for session %s", mc.sessionID)
				return

			default:
				mc.logger.Infof("QR event for session %s: %s", mc.sessionID, evt.Event)
			}

		case <-mc.qrStopChannel:
			mc.logger.Infof("QR loop stop signal received for session %s", mc.sessionID)
			
			mc.clearQRCode()
			if mc.manager != nil {
				go mc.manager.clearQRCodeInDatabase(mc.sessionID)
			}
			return

		case <-qrCtx.Done():
			mc.logger.Infof("QR loop context cancelled for session %s", mc.sessionID)
			
			mc.clearQRCode()
			if mc.manager != nil {
				go mc.manager.clearQRCodeInDatabase(mc.sessionID)
			}
			return

		case <-mc.ctx.Done():
			mc.logger.Infof("Main context cancelled, stopping QR loop for session %s", mc.sessionID)
			return
		}
	}
}


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
			
		}
	}
}


func (mc *MeowClient) Disconnect() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if !mc.client.IsConnected() {
		return
	}

	mc.logger.Infof("Disconnecting client for session %s", mc.sessionID)

	
	mc.stopQRLoop()

	
	select {
	case mc.killChannel <- true:
	default:
		
	}

	mc.client.Disconnect()
	mc.setStatus(types.StatusDisconnected)
}


func (mc *MeowClient) Stop() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.logger.Infof("Stopping client for session %s", mc.sessionID)

	
	mc.stopQRLoop()

	
	mc.cancel()

	
	select {
	case mc.killChannel <- true:
	default:
		
	}

	if mc.client.IsConnected() {
		mc.client.Disconnect()
	}

	mc.setStatus(types.StatusDisconnected)
}


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


func (mc *MeowClient) GetQRCode(ctx context.Context) (string, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Early return if already authenticated
	if mc.client.IsConnected() && mc.client.Store.ID != nil {
		return "", errors.New(ErrClientAlreadyAuth)
	}

	// Early return if session is already connected
	if Status.IsConnectedStatus(mc.status) {
		return "", errors.New(ErrSessionNotConnected)
	}

	// Return existing QR code if available
	if mc.qrCode != "" {
		return mc.qrCode, nil
	}

	// Check if client is connected
	if !mc.client.IsConnected() {
		return "", fmt.Errorf(ErrClientNotConnected + ", call Connect() first")
	}

	// QR code not yet available
	return "", errors.New(ErrQRNotAvailable)
}


func (mc *MeowClient) PairPhone(ctx context.Context, phoneNumber string) (string, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.client.IsConnected() && mc.client.Store.ID != nil {
		return "", fmt.Errorf("client is already authenticated")
	}

	if !mc.client.IsConnected() {
		return "", fmt.Errorf("client is not connected, call Connect() first")
	}

	
	code, err := mc.client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", fmt.Errorf("failed to request pairing code: %w", err)
	}

	mc.logger.Infof("Pairing code requested for session %s, phone %s", mc.sessionID, phoneNumber)
	return code, nil
}


func (mc *MeowClient) IsConnected() bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.client.IsConnected()
}


func (mc *MeowClient) GetStatus() types.Status {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.status
}


func (mc *MeowClient) GetLastActivity() time.Time {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.lastActivity
}


func (mc *MeowClient) GetJID() waTypes.JID {
	if mc.client.Store.ID != nil {
		return *mc.client.Store.ID
	}
	return waTypes.JID{}
}


func (mc *MeowClient) SendMessage(ctx context.Context, to waTypes.JID, message string) error {
	

	
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


func (mc *MeowClient) GetClient() *whatsmeow.Client {
	return mc.client
}


func (mc *MeowClient) SendTextMessage(ctx context.Context, to waTypes.JID, text string, contextInfo *waE2E.ContextInfo) (*whatsmeow.SendResponse, error) {
	mc.logger.Infof("DEBUG: MeowClient.SendTextMessage called - to: %s, text: %s", to.String(), text)

	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	mc.logger.Infof("DEBUG: Creating message...")
	msg := MsgBuilder.BuildTextMessage(text, contextInfo)

	mc.logger.Infof("DEBUG: Calling whatsmeow client.SendMessage...")
	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		mc.logger.Errorf("DEBUG: whatsmeow client.SendMessage failed: %v", err)
		return nil, Error.WrapError(err, "failed to send text message")
	}

	mc.logger.Infof("DEBUG: whatsmeow client.SendMessage succeeded")
	// Skip updateActivity() to avoid potential deadlock
	mc.logger.Infof("DEBUG: SendTextMessage completed successfully")
	return resp, nil
}


func (mc *MeowClient) SendLocationMessage(ctx context.Context, to waTypes.JID, latitude, longitude float64, name, address string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	msg := MsgBuilder.BuildLocationMessage(latitude, longitude, name, address)
	sender := NewMessageSender(mc.client)

	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send location message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendContactMessage(ctx context.Context, to waTypes.JID, displayName, vcard string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	msg := MsgBuilder.BuildContactMessage(displayName, vcard)
	sender := NewMessageSender(mc.client)

	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send contact message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) ReactToMessage(ctx context.Context, chatJID waTypes.JID, messageID string, emoji string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	
	
	return fmt.Errorf("message reactions not implemented for new whatsmeow version")
}


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


func (mc *MeowClient) MarkMessageRead(ctx context.Context, chatJID waTypes.JID, messageIDs []string) error {
	if !mc.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	
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


func (mc *MeowClient) JoinGroupWithLink(ctx context.Context, inviteCode string) (*waTypes.GroupInfo, error) {
	if !mc.IsConnected() {
		return nil, fmt.Errorf("client is not connected")
	}

	groupJID, err := mc.client.JoinGroupWithLink(inviteCode)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	
	groupInfo, err := mc.client.GetGroupInfo(groupJID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info after joining: %w", err)
	}

	mc.updateActivity()
	return groupInfo, nil
}


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



func (mc *MeowClient) UpdateGroupParticipants(ctx context.Context, groupJID waTypes.JID, participants []waTypes.JID, action string) error {
	return fmt.Errorf("group participant management not implemented for new whatsmeow version")
}


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


func (mc *MeowClient) Cleanup() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	
	if mc.eventHandlerID != 0 {
		mc.client.RemoveEventHandler(mc.eventHandlerID)
	}

	
	if mc.client.IsConnected() {
		mc.client.Disconnect()
	}

	mc.logger.Infof("Cleaned up client for session %s", mc.sessionID)
}


func (mc *MeowClient) setStatus(status types.Status) {
	mc.status = status
	mc.lastActivity = time.Now()

	
	if mc.manager != nil {
		go mc.manager.OnClientStatusChange(mc.sessionID, status)
	}
}


func (mc *MeowClient) updateActivity() {
	// Try to acquire lock with timeout to avoid deadlock
	done := make(chan bool, 1)
	go func() {
		mc.mu.Lock()
		mc.lastActivity = time.Now()
		mc.mu.Unlock()
		done <- true
	}()

	select {
	case <-done:
		// Successfully updated activity
		return
	case <-time.After(100 * time.Millisecond):
		// Timeout - log warning but don't block
		mc.logger.Warnf("updateActivity timeout for session %s - potential deadlock avoided", mc.sessionID)
		return
	}
}


func (mc *MeowClient) setQRCode(qrCode string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.qrCode = qrCode
	mc.logger.Infof("QR code updated for session %s", mc.sessionID)
}


func (mc *MeowClient) clearQRCode() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.qrCode = ""
}


func (mc *MeowClient) stopQRLoop() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Early return if already stopped
	if !mc.qrLoopActive {
		mc.logger.Debugf("QR loop already stopped for session %s", mc.sessionID)
		return
	}

	// Try to send stop signal with timeout
	if Channel.SafeChannelSend(mc.qrStopChannel, true, 1*time.Second) {
		mc.logger.Infof("QR loop stop signal sent for session %s", mc.sessionID)
		mc.qrLoopActive = false
		return
	}

	// Fallback to context cancellation
	mc.logger.Warnf("QR loop stop signal timeout for session %s, using context cancellation", mc.sessionID)
	if mc.qrLoopCancel != nil {
		mc.qrLoopCancel()
	}
	mc.qrLoopActive = false
}


func (mc *MeowClient) onConnected() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.setStatus(types.StatusConnected)
	mc.clearQRCode()

	
	if mc.manager != nil {
		go mc.manager.clearQRCodeInDatabase(mc.sessionID)
	}

	
	mc.stopQRLoop()

	mc.logger.Infof("Client connected for session %s", mc.sessionID)
}


func (mc *MeowClient) onDisconnected() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.setStatus(types.StatusDisconnected)
	mc.logger.Infof("Client disconnected for session %s", mc.sessionID)
}


func (mc *MeowClient) onError(err error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.setStatus(types.StatusError)
	mc.logger.Errorf("Client error for session %s: %v", mc.sessionID, err)
}


func (mc *MeowClient) GetSessionID() string {
	return mc.sessionID
}


func (mc *MeowClient) SendImageMessage(ctx context.Context, to waTypes.JID, imageData []byte, caption string, mimeType string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	// Upload image
	uploader := NewMediaUploader(mc.client)
	uploaded, err := uploader.UploadMedia(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return nil, Error.WrapError(err, "failed to upload image")
	}

	// Build and send message
	params := MediaMessageParams{
		UploadResponse: uploaded,
		Caption:        caption,
		MimeType:       mimeType,
	}
	msg := MsgBuilder.BuildImageMessage(params)

	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send image message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendAudioMessage(ctx context.Context, to waTypes.JID, audioData []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	// Upload audio using MediaUploader utility
	uploader := NewMediaUploader(mc.client)
	uploaded, err := uploader.UploadMedia(ctx, audioData, whatsmeow.MediaAudio)
	if err != nil {
		return nil, Error.WrapError(err, "failed to upload audio")
	}

	// Build and send message using MessageBuilder
	params := MediaMessageParams{
		UploadResponse: uploaded,
		MimeType:       func() string {
			if mimeType != "" {
				return mimeType
			}
			return "audio/ogg; codecs=opus"
		}(),
	}

	msg := MsgBuilder.BuildAudioMessage(params, true) // PTT = true
	sender := NewMessageSender(mc.client)

	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send audio message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendDocumentMessage(ctx context.Context, to waTypes.JID, documentData []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	// Upload document using MediaUploader utility
	uploader := NewMediaUploader(mc.client)
	uploaded, err := uploader.UploadMedia(ctx, documentData, whatsmeow.MediaDocument)
	if err != nil {
		return nil, Error.WrapError(err, "failed to upload document")
	}

	// Build and send message using MessageBuilder
	params := MediaMessageParams{
		UploadResponse: uploaded,
		Caption:        caption,
		MimeType:       mimeType,
		FileName:       filename,
	}
	msg := MsgBuilder.BuildDocumentMessage(params)

	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send document message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendVideoMessage(ctx context.Context, to waTypes.JID, videoData []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	// Upload video using MediaUploader utility
	uploader := NewMediaUploader(mc.client)
	uploaded, err := uploader.UploadMedia(ctx, videoData, whatsmeow.MediaVideo)
	if err != nil {
		return nil, Error.WrapError(err, "failed to upload video")
	}

	// Build and send message using MessageBuilder
	params := MediaMessageParams{
		UploadResponse: uploaded,
		Caption:        caption,
		MimeType:       mimeType,
	}
	msg := MsgBuilder.BuildVideoMessage(params)

	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send video message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendStickerMessage(ctx context.Context, to waTypes.JID, stickerData []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	// Upload sticker using MediaUploader utility
	uploader := NewMediaUploader(mc.client)
	uploaded, err := uploader.UploadMedia(ctx, stickerData, whatsmeow.MediaImage)
	if err != nil {
		return nil, Error.WrapError(err, "failed to upload sticker")
	}

	// Build and send message using MessageBuilder
	params := MediaMessageParams{
		UploadResponse: uploaded,
		MimeType:       mimeType,
	}
	msg := MsgBuilder.BuildStickerMessage(params)

	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send sticker message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendButtonsMessage(ctx context.Context, to waTypes.JID, text string, buttons []types.Button, footer string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	// Build buttons message using MessageBuilder
	msg := MsgBuilder.BuildButtonsMessage(text, buttons, footer)

	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send buttons message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendListMessage(ctx context.Context, to waTypes.JID, text, buttonText string, sections []types.Section, footer string) (*whatsmeow.SendResponse, error) {
	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	// Build list message using MessageBuilder
	msg := MsgBuilder.BuildListMessage(text, buttonText, sections, footer)

	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		return nil, Error.WrapError(err, "failed to send list message")
	}

	// Skip updateActivity() to avoid potential deadlock
	return resp, nil
}


func (mc *MeowClient) SendPollMessage(ctx context.Context, to waTypes.JID, name string, options []string, selectableCount int) (*whatsmeow.SendResponse, error) {
	mc.logger.Infof("DEBUG: MeowClient.SendPollMessage called - to: %s, name: %s", to.String(), name)

	// Skip connection validation to avoid deadlock - validation is done at service level
	mc.logger.Infof("DEBUG: Skipping connection validation to avoid deadlock")

	mc.logger.Infof("DEBUG: Building poll message...")
	msg := MsgBuilder.BuildPollMessage(name, options, selectableCount)

	mc.logger.Infof("DEBUG: Calling whatsmeow client.SendMessage...")
	sender := NewMessageSender(mc.client)
	resp, err := sender.SendMessage(ctx, to, msg)
	if err != nil {
		mc.logger.Errorf("DEBUG: whatsmeow client.SendMessage failed: %v", err)
		return nil, Error.WrapError(err, "failed to send poll message")
	}

	mc.logger.Infof("DEBUG: whatsmeow client.SendMessage succeeded")
	// Skip updateActivity() to avoid potential deadlock
	mc.logger.Infof("DEBUG: SendPollMessage completed successfully")
	return resp, nil
}
