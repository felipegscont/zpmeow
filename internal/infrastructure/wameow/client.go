package wameow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infrastructure/logging"
	"zpmeow/internal/shared/types"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// WameowClient - Unified WhatsApp client with all functionality
type WameowClient struct {
	sessionID    string
	client       *whatsmeow.Client
	eventHandler EventHandler
	sessionRepo  session.Repository
	logger       logging.Logger
	waLogger     waLog.Logger

	// Status and activity tracking
	mu           sync.RWMutex
	status       types.Status
	lastActivity time.Time

	// QR code management
	qrCode       string
	qrCodeBase64 string // Base64 encoded QR code image
	qrLoopActive bool
	qrLoopCancel context.CancelFunc

	// Event handling
	eventHandlerID uint32

	// Context and channels for lifecycle management
	ctx           context.Context
	cancel        context.CancelFunc
	killChannel   chan bool
	qrStopChannel chan bool

	// Retry configuration
	maxRetries    int
	retryCount    int
	retryInterval time.Duration

	// Helpers
	sessionHelper    *SessionHelper
	qrHelper         *QRCodeHelper
	connectionHelper *ConnectionHelper
}

type EventHandler interface {
	HandleEvent(interface{})
}

// NewWameowClient creates a new unified WhatsApp client
func NewWameowClient(sessionID string, container *sqlstore.Container, waLogger waLog.Logger, eventHandler EventHandler, sessionRepo session.Repository) (*WameowClient, error) {
	return NewWameowClientWithDeviceJID(sessionID, "", container, waLogger, eventHandler, sessionRepo)
}

// NewWameowClientWithDeviceJID creates a new unified WhatsApp client with expected device JID
func NewWameowClientWithDeviceJID(sessionID, expectedDeviceJID string, container *sqlstore.Container, waLogger waLog.Logger, eventHandler EventHandler, sessionRepo session.Repository) (*WameowClient, error) {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logging.GetLogger().Sub("wameow-client").Sub(sessionID)

	// Get or create device store with expected device JID
	deviceStore := GetDeviceStoreForSession(sessionID, expectedDeviceJID, container)
	if deviceStore == nil {
		return nil, fmt.Errorf("failed to create device store for session %s", sessionID)
	}

	// Configure device properties
	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_UNKNOWN.Enum()
	osName := "zpmeow"
	store.DeviceProps.Os = &osName

	// Create WhatsApp client
	waClient := whatsmeow.NewClient(deviceStore, waLogger)
	if waClient == nil {
		return nil, fmt.Errorf("failed to create WhatsApp client for session %s", sessionID)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create helpers
	sessionHelper := NewSessionHelper(sessionRepo, appLogger)
	qrHelper := NewQRCodeHelper(appLogger)
	connectionHelper := NewConnectionHelper(appLogger)

	client := &WameowClient{
		sessionID:        sessionID,
		client:           waClient,
		eventHandler:     eventHandler,
		sessionRepo:      sessionRepo,
		logger:           appLogger,
		waLogger:         waLogger,
		status:           types.StatusDisconnected,
		lastActivity:     time.Now(),
		ctx:              ctx,
		cancel:           cancel,
		killChannel:      make(chan bool, 1),
		qrStopChannel:    make(chan bool, 1),
		maxRetries:       5,
		retryCount:       0,
		retryInterval:    30 * time.Second,
		sessionHelper:    sessionHelper,
		qrHelper:         qrHelper,
		connectionHelper: connectionHelper,
	}

	// Register event handler if provided
	if eventHandler != nil {
		client.eventHandlerID = waClient.AddEventHandler(eventHandler.HandleEvent)
	}

	return client, nil
}

// WameowClient methods
func (c *WameowClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Validate client and store
	if err := ValidateClientAndStore(c.client, c.sessionID, c.logger); err != nil {
		return err
	}

	if c.client.IsConnected() {
		c.logger.Debugf("Client already connected for session %s", c.sessionID)
		return nil
	}

	c.setStatus(types.StatusConnecting)
	c.sessionHelper.UpdateSessionStatus(c.sessionID, session.StatusConnecting)

	// Start client loop in background
	go c.startClientLoop()

	return nil
}

func (c *WameowClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Infof("Disconnecting client for session %s", c.sessionID)

	// Stop QR loop if active
	c.stopQRLoop()

	// Disconnect client safely
	c.connectionHelper.SafeDisconnect(c.client, c.sessionID)

	// Cancel context
	if c.cancel != nil {
		c.cancel()
	}

	c.setStatus(types.StatusDisconnected)
	c.sessionHelper.UpdateSessionStatus(c.sessionID, session.StatusDisconnected)

	return nil
}

func (c *WameowClient) GetQRCode() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.qrCode == "" {
		return "", fmt.Errorf("no QR code available")
	}

	return c.qrCode, nil
}

func (c *WameowClient) GetQRCodeBase64() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.qrCodeBase64 == "" {
		return "", fmt.Errorf("no QR code image available")
	}

	return c.qrCodeBase64, nil
}

func (c *WameowClient) PairPhone(phoneNumber string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Infof("Pairing phone %s for session %s", phoneNumber, c.sessionID)

	if phoneNumber == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}

	code, err := c.client.PairPhone(context.Background(), phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		c.logger.Errorf("Failed to pair phone for session %s: %v", c.sessionID, err)
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	c.logger.Infof("Pairing code generated for session %s", c.sessionID)
	return code, nil
}

func (c *WameowClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client.IsConnected()
}

func (c *WameowClient) GetStatus() types.Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

func (c *WameowClient) GetClient() *whatsmeow.Client {
	return c.client
}

func (c *WameowClient) GetJID() waTypes.JID {
	if c.client.Store.ID == nil {
		return waTypes.EmptyJID
	}
	return *c.client.Store.ID
}

func (c *WameowClient) Logout() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Infof("Logging out session %s", c.sessionID)

	err := c.client.Logout(context.Background())
	if err != nil {
		c.logger.Errorf("Failed to logout session %s: %v", c.sessionID, err)
		return fmt.Errorf("failed to logout: %w", err)
	}

	if c.client.IsConnected() {
		c.client.Disconnect()
	}

	c.setStatus(types.StatusDisconnected)
	c.logger.Infof("Successfully logged out session %s", c.sessionID)
	return nil
}

func (c *WameowClient) Reconnect(ctx context.Context) error {
	c.logger.Infof("Attempting to reconnect session %s", c.sessionID)

	if err := c.Disconnect(); err != nil {
		c.logger.Warnf("Error during disconnect before reconnect for session %s: %v", c.sessionID, err)
	}

	time.Sleep(2 * time.Second)
	return c.Connect()
}

func (c *WameowClient) GetLastActivity() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastActivity
}

func (c *WameowClient) UpdateActivity() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastActivity = time.Now()
}

func (c *WameowClient) IsLoggedIn() bool {
	return c.client.Store.ID != nil
}

func (c *WameowClient) GetSessionID() string {
	return c.sessionID
}

// setStatus safely sets the status with logging
func (c *WameowClient) setStatus(status types.Status) {
	c.status = status
	c.lastActivity = time.Now()
	if status == types.StatusConnected || status == types.StatusDisconnected ||
		status == types.StatusConnecting || status == types.StatusError {
		c.logger.Infof("Session %s status: %s", c.sessionID, status)
	}
}


// startClientLoop starts the main client loop
func (c *WameowClient) startClientLoop() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Errorf("Client loop panic for session %s: %v", c.sessionID, r)
		}
	}()

	// Check if device is registered
	if !IsDeviceRegistered(c.client) {
		c.logger.Infof("Device not registered for session %s, starting QR code process", c.sessionID)
		c.handleNewDeviceRegistration()
	} else {
		c.logger.Infof("Device already registered for session %s, connecting directly", c.sessionID)
		c.handleExistingDeviceConnection()
	}
}

// handleNewDeviceRegistration handles QR code generation for new devices
func (c *WameowClient) handleNewDeviceRegistration() {
	qrChan, err := c.client.GetQRChannel(context.Background())
	if err != nil {
		c.logger.Errorf("Failed to get QR channel for session %s: %v", c.sessionID, err)
		c.setStatus(types.StatusDisconnected)
		return
	}

	err = c.client.Connect()
	if err != nil {
		c.logger.Errorf("Failed to connect client for session %s: %v", c.sessionID, err)
		c.setStatus(types.StatusDisconnected)
		return
	}

	c.handleQRLoop(qrChan)
}

// handleExistingDeviceConnection handles connection for already registered devices
func (c *WameowClient) handleExistingDeviceConnection() {
	c.logger.Infof("Connecting existing device for session %s", c.sessionID)

	err := c.client.Connect()
	if err != nil {
		c.logger.Errorf("Failed to connect client for session %s: %v", c.sessionID, err)
		c.setStatus(types.StatusDisconnected)
		c.sessionHelper.UpdateSessionStatus(c.sessionID, session.StatusDisconnected)
		return
	}

	// Wait a bit for connection to establish
	time.Sleep(2 * time.Second)

	if c.client.IsConnected() {
		c.logger.Infof("Successfully connected session %s", c.sessionID)
		c.setStatus(types.StatusConnected)
		c.sessionHelper.UpdateSessionStatus(c.sessionID, session.StatusConnected)
	} else {
		c.logger.Warnf("Connection attempt completed but client not connected for session %s", c.sessionID)
		c.setStatus(types.StatusDisconnected)
		c.sessionHelper.UpdateSessionStatus(c.sessionID, session.StatusDisconnected)
	}
}

// handleQRLoop handles QR code events
func (c *WameowClient) handleQRLoop(qrChan <-chan whatsmeow.QRChannelItem) {
	c.mu.Lock()
	c.qrLoopActive = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.qrLoopActive = false
		c.mu.Unlock()
	}()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("QR loop cancelled for session %s", c.sessionID)
			return

		case <-c.qrStopChannel:
			c.logger.Infof("QR loop stopped for session %s", c.sessionID)
			return

		case evt, ok := <-qrChan:
			if !ok {
				c.logger.Infof("QR channel closed for session %s", c.sessionID)
				c.setStatus(types.StatusDisconnected)
				// Clear QR code in database
				c.sessionHelper.UpdateSessionQRCode(c.sessionID, "")
				return
			}

			switch evt.Event {
			case "code":
				c.mu.Lock()
				c.qrCode = evt.Code
				c.qrCodeBase64 = c.qrHelper.GenerateQRCodeImage(evt.Code)
				c.mu.Unlock()

				c.qrHelper.DisplayQRCodeInTerminal(evt.Code, c.sessionID)
				c.logger.Infof("QR code generated for session %s", c.sessionID)
				c.setStatus(types.StatusConnecting)

				// Update QR code in database
				c.sessionHelper.UpdateSessionQRCode(c.sessionID, evt.Code)

			case "success":
				c.logger.Infof("QR code scanned successfully for session %s", c.sessionID)
				c.setStatus(types.StatusConnected)

				go c.persistQRSuccess()
				return

			case "timeout":
				c.logger.Warnf("QR code timeout for session %s", c.sessionID)
				c.mu.Lock()
				c.qrCode = ""
				c.qrCodeBase64 = ""
				c.mu.Unlock()

				c.setStatus(types.StatusDisconnected)

				// Update QR code in database (clear it)
				c.sessionHelper.UpdateSessionQRCode(c.sessionID, "")
				return

			default:
				c.logger.Infof("QR event: %s for session %s", evt.Event, c.sessionID)
			}
		}
	}
}

// stopQRLoop stops the QR code loop
func (c *WameowClient) stopQRLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.qrLoopActive {
		select {
		case c.qrStopChannel <- true:
		default:
		}
	}
}

// persistQRSuccess persists successful QR scan to database
func (c *WameowClient) persistQRSuccess() {
	if c.sessionRepo == nil {
		c.logger.Warnf("No session repository available for session %s", c.sessionID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionEntity, err := c.sessionRepo.GetByID(ctx, c.sessionID)
	if err != nil {
		c.logger.Errorf("Failed to get session %s from database: %v", c.sessionID, err)
		return
	}

	var deviceJID string
	if c.client != nil && c.client.Store.ID != nil {
		deviceJID = c.client.Store.ID.String()
	}

	if deviceJID != "" {
		// Validate device uniqueness before assigning
		if err := c.sessionRepo.ValidateDeviceUniqueness(ctx, c.sessionID, deviceJID); err != nil {
			c.logger.Errorf("Device uniqueness validation failed for session %s: %v", c.sessionID, err)
			// Don't update the session if device is already in use
			return
		}

		c.logger.Infof("Assigning device JID %s to session %s", deviceJID, c.sessionID)
		sessionEntity.WaJID = deviceJID
	}

	sessionEntity.SetStatus(session.StatusConnected)
	sessionEntity.SetQRCode("") // Clear QR code after successful pairing

	if err := c.sessionRepo.Update(ctx, sessionEntity); err != nil {
		c.logger.Errorf("Failed to update session %s in database after QR scan: %v", c.sessionID, err)
		return
	}

	c.logger.Infof("Successfully updated session %s in database after QR scan: JID=%s, Status=%s", c.sessionID, deviceJID, session.StatusConnected)
}

