package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"zpmeow/internal/infra/logger"
	"zpmeow/internal/shared/types"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// MeowClient represents a WhatsApp client for a specific session
type MeowClient struct {
	sessionID    string
	client       *whatsmeow.Client
	eventHandler EventHandler
	logger       logger.Logger
	waLogger     waLog.Logger

	// Status management
	mu           sync.RWMutex
	status       types.Status
	lastActivity time.Time
	qrCode       string

	// Event handling
	eventHandlerID uint32

	// Lifecycle management
	killChannel chan bool
	ctx         context.Context
	cancel      context.CancelFunc

	// QR code management
	qrStopChannel chan bool
	qrLoopActive  bool
	qrLoopCancel  context.CancelFunc
}

// EventHandler interface for handling WhatsApp events
type EventHandler interface {
	HandleEvent(interface{})
}

// NewMeowClient creates a new WhatsApp client for a session
func NewMeowClient(sessionID string, deviceStore *store.Device, waLogger waLog.Logger, eventHandler EventHandler) (*MeowClient, error) {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("meow-client").Sub(sessionID)

	// Create WhatsApp client
	waClient := whatsmeow.NewClient(deviceStore, waLogger)

	// Create context for lifecycle management
	ctx, cancel := context.WithCancel(context.Background())

	// Create MeowClient instance
	meowClient := &MeowClient{
		sessionID:     sessionID,
		client:        waClient,
		eventHandler:  eventHandler,
		logger:        appLogger,
		waLogger:      waLogger,
		status:        types.StatusDisconnected,
		lastActivity:  time.Now(),
		killChannel:   make(chan bool, 1),
		qrStopChannel: make(chan bool, 1),
		ctx:           ctx,
		cancel:        cancel,
	}

	// Register event handler
	if eventHandler != nil {
		meowClient.eventHandlerID = waClient.AddEventHandler(eventHandler.HandleEvent)
	}

	return meowClient, nil
}

// Connect establishes connection to WhatsApp
func (mc *MeowClient) Connect(ctx context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Check if already connected
	if mc.client.IsConnected() {
		return nil
	}

	mc.setStatus(types.StatusConnecting)
	mc.logger.Infof("Connecting client for session %s", mc.sessionID)

	// Start client connection in background
	go mc.startClientLoop()

	return nil
}

// Disconnect closes the WhatsApp connection
func (mc *MeowClient) Disconnect() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.logger.Infof("Disconnecting client for session %s", mc.sessionID)

	// Stop QR loop if active
	mc.stopQRLoop()

	// Disconnect WhatsApp client
	if mc.client.IsConnected() {
		mc.client.Disconnect()
	}

	// Cancel context
	if mc.cancel != nil {
		mc.cancel()
	}

	mc.setStatus(types.StatusDisconnected)
	return nil
}

// IsConnected returns the connection status
func (mc *MeowClient) IsConnected() bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.client.IsConnected()
}

// GetStatus returns the current client status
func (mc *MeowClient) GetStatus() types.Status {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.status
}

// GetQRCode returns the current QR code for pairing
func (mc *MeowClient) GetQRCode() (string, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.qrCode == "" {
		return "", fmt.Errorf("no QR code available")
	}

	return mc.qrCode, nil
}

// PairPhone pairs the client with a phone number
func (mc *MeowClient) PairPhone(phoneNumber string) (string, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.logger.Infof("Pairing phone %s for session %s", phoneNumber, mc.sessionID)

	// Validate phone number format
	if phoneNumber == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}

	// Request pairing code
	code, err := mc.client.PairPhone(context.Background(), phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		mc.logger.Errorf("Failed to pair phone for session %s: %v", mc.sessionID, err)
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	mc.logger.Infof("Pairing code generated for session %s", mc.sessionID)
	return code, nil
}

// GetJID returns the WhatsApp JID for this client
func (mc *MeowClient) GetJID() waTypes.JID {
	if mc.client.Store.ID == nil {
		return waTypes.EmptyJID
	}
	return *mc.client.Store.ID
}

// Internal methods

// setStatus updates the client status (internal use)
func (mc *MeowClient) setStatus(status types.Status) {
	mc.status = status
	mc.lastActivity = time.Now()
	mc.logger.Debugf("Status changed to %s for session %s", status, mc.sessionID)
}

// startClientLoop handles the client connection lifecycle
func (mc *MeowClient) startClientLoop() {
	mc.logger.Infof("Starting client loop for session %s", mc.sessionID)

	// Check if device is registered
	if mc.client.Store.ID == nil {
		mc.logger.Infof("Device not registered for session %s, waiting for QR code", mc.sessionID)

		// Get QR channel
		qrChan, err := mc.client.GetQRChannel(context.Background())
		if err != nil {
			mc.logger.Errorf("Failed to get QR channel for session %s: %v", mc.sessionID, err)
			mc.setStatus(types.StatusError)
			return
		}

		// Connect client
		err = mc.client.Connect()
		if err != nil {
			mc.logger.Errorf("Failed to connect client for session %s: %v", mc.sessionID, err)
			mc.setStatus(types.StatusError)
			return
		}

		// Handle QR code loop
		go mc.handleQRLoop(qrChan)
	} else {
		mc.logger.Infof("Already logged in, just connecting for session %s", mc.sessionID)

		// Connect directly
		err := mc.client.Connect()
		if err != nil {
			mc.logger.Errorf("Failed to connect client for session %s: %v", mc.sessionID, err)
			mc.setStatus(types.StatusError)
			return
		}

		mc.setStatus(types.StatusConnected)
	}
}

// handleQRLoop manages QR code generation and display
func (mc *MeowClient) handleQRLoop(qrChan <-chan whatsmeow.QRChannelItem) {
	mc.mu.Lock()
	mc.qrLoopActive = true
	ctx, cancel := context.WithCancel(mc.ctx)
	mc.qrLoopCancel = cancel
	mc.mu.Unlock()

	defer func() {
		mc.mu.Lock()
		mc.qrLoopActive = false
		mc.qrLoopCancel = nil
		mc.mu.Unlock()
	}()

	for {
		select {
		case evt, ok := <-qrChan:
			if !ok {
				mc.logger.Infof("QR channel closed for session %s", mc.sessionID)
				return
			}

			if evt.Event == "code" {
				mc.mu.Lock()
				mc.qrCode = evt.Code
				mc.mu.Unlock()

				mc.logger.Infof("QR code generated for session %s", mc.sessionID)
			} else {
				mc.logger.Infof("QR event: %s for session %s", evt.Event, mc.sessionID)
			}

		case <-ctx.Done():
			mc.logger.Infof("QR loop cancelled for session %s", mc.sessionID)
			return

		case <-mc.qrStopChannel:
			mc.logger.Infof("QR loop stopped for session %s", mc.sessionID)
			return
		}
	}
}

// stopQRLoop stops the QR code generation loop
func (mc *MeowClient) stopQRLoop() {
	if mc.qrLoopActive && mc.qrLoopCancel != nil {
		mc.qrLoopCancel()
	}

	// Send stop signal (non-blocking)
	select {
	case mc.qrStopChannel <- true:
	default:
	}
}

// Cleanup performs cleanup when the client is destroyed
func (mc *MeowClient) Cleanup() {
	mc.logger.Infof("Cleaning up client for session %s", mc.sessionID)

	// Remove event handler
	if mc.eventHandlerID != 0 {
		mc.client.RemoveEventHandler(mc.eventHandlerID)
	}

	// Disconnect if connected
	if mc.client.IsConnected() {
		mc.client.Disconnect()
	}

	// Cancel context
	if mc.cancel != nil {
		mc.cancel()
	}

	// Stop QR loop
	mc.stopQRLoop()
}
