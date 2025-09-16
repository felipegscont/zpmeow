package wameow

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infrastructure/logging"
	"zpmeow/internal/shared/types"

	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
)

// SessionHelper contains common session operations
type SessionHelper struct {
	sessionRepo session.Repository
	logger      logging.Logger
}

// NewSessionHelper creates a new session helper
func NewSessionHelper(sessionRepo session.Repository, logger logging.Logger) *SessionHelper {
	return &SessionHelper{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

// UpdateSessionStatus updates the session status in the database
func (h *SessionHelper) UpdateSessionStatus(sessionID string, status session.Status) {
	if h.sessionRepo == nil {
		h.logger.Warnf("No session repository available for session %s", sessionID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionEntity, err := h.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		h.logger.Errorf("Failed to get session %s from database: %v", sessionID, err)
		return
	}

	sessionEntity.SetStatus(status)

	if err := h.sessionRepo.Update(ctx, sessionEntity); err != nil {
		h.logger.Errorf("Failed to update session %s status to %s in database: %v", sessionID, status, err)
		return
	}

	h.logger.Infof("Successfully updated session %s status to %s in database", sessionID, status)
}

// UpdateSessionQRCode updates the QR code in the database
func (h *SessionHelper) UpdateSessionQRCode(sessionID string, qrCode string) {
	if h.sessionRepo == nil {
		h.logger.Warnf("No session repository available for session %s", sessionID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionEntity, err := h.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		h.logger.Errorf("Failed to get session %s from database: %v", sessionID, err)
		return
	}

	sessionEntity.SetQRCode(qrCode)

	if err := h.sessionRepo.Update(ctx, sessionEntity); err != nil {
		h.logger.Errorf("Failed to update session %s QR code in database: %v", sessionID, err)
		return
	}

	h.logger.Infof("Successfully updated session %s QR code in database", sessionID)
}

// ValidateClientAndStore validates WhatsApp client and store
func ValidateClientAndStore(client *whatsmeow.Client, sessionID string, logger logging.Logger) error {
	if client == nil {
		return fmt.Errorf("WhatsApp client is nil for session %s", sessionID)
	}

	if client.Store == nil {
		return fmt.Errorf("WhatsApp client store is nil for session %s", sessionID)
	}

	return nil
}

// GetOrCreateDeviceStore gets or creates a device store for a session
// This function properly handles existing devices based on session's device_jid
func GetOrCreateDeviceStore(sessionID string, container *sqlstore.Container) *store.Device {
	return GetDeviceStoreForSession(sessionID, "", container)
}

// GetDeviceStoreForSession gets or creates a device store for a session with optional device JID
// This follows the reference implementation logic from wmiau.go.bak
func GetDeviceStoreForSession(sessionID, expectedDeviceJID string, container *sqlstore.Container) *store.Device {
	var deviceStore *store.Device
	var err error

	// If we have an expected device JID, try to get that specific device (like reference)
	if expectedDeviceJID != "" {
		// Parse the JID like in the reference
		jid, ok := parseJID(expectedDeviceJID)
		if ok {
			// Use container.GetDevice like in the reference implementation
			deviceStore, err = container.GetDevice(context.Background(), jid)
			if err != nil {
				fmt.Printf("Failed to get device for JID %s: %v\n", expectedDeviceJID, err)
				// Fallback to creating new device
				deviceStore = container.NewDevice()
			} else {
				fmt.Printf("Successfully retrieved existing device for JID %s\n", expectedDeviceJID)
			}
		} else {
			fmt.Printf("Failed to parse JID %s, creating new device\n", expectedDeviceJID)
			deviceStore = container.NewDevice()
		}
	} else {
		// No JID provided, create new device
		fmt.Printf("No device JID provided for session %s, creating new device\n", sessionID)
		deviceStore = container.NewDevice()
	}

	// Final fallback
	if deviceStore == nil {
		fmt.Printf("Device store is nil, creating fallback device\n")
		deviceStore = container.NewDevice()
	}

	return deviceStore
}

// parseJID parses a JID string into waTypes.JID (from reference implementation)
func parseJID(arg string) (waTypes.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return waTypes.NewJID(arg, waTypes.DefaultUserServer), true
	} else {
		recipient, err := waTypes.ParseJID(arg)
		if err != nil {
			fmt.Printf("Invalid JID: %v\n", err)
			return recipient, false
		} else if recipient.User == "" {
			fmt.Printf("Invalid JID no server specified\n")
			return recipient, false
		}
		return recipient, true
	}
}

// QRCodeHelper contains QR code related utilities
type QRCodeHelper struct {
	logger logging.Logger
}

// NewQRCodeHelper creates a new QR code helper
func NewQRCodeHelper(logger logging.Logger) *QRCodeHelper {
	return &QRCodeHelper{
		logger: logger,
	}
}

// GenerateQRCodeImage generates a base64 encoded QR code image
func (h *QRCodeHelper) GenerateQRCodeImage(qrText string) string {
	qrPNG, err := qrcode.Encode(qrText, qrcode.Medium, 256)
	if err != nil {
		h.logger.Errorf("Failed to generate QR code image: %v", err)
		return ""
	}

	base64Str := base64.StdEncoding.EncodeToString(qrPNG)
	return "data:image/png;base64," + base64Str
}

// DisplayQRCodeInTerminal displays QR code in terminal
func (h *QRCodeHelper) DisplayQRCodeInTerminal(qrCode, sessionID string) {
	fmt.Printf("\n=== QR Code for Session %s ===\n", sessionID)
	qrterminal.GenerateHalfBlock(qrCode, qrterminal.L, os.Stdout)
	fmt.Printf("QR Code String: %s\n", qrCode)
	fmt.Printf("=== End QR Code ===\n\n")
}

// ConnectionHelper contains connection related utilities
type ConnectionHelper struct {
	logger logging.Logger
}

// NewConnectionHelper creates a new connection helper
func NewConnectionHelper(logger logging.Logger) *ConnectionHelper {
	return &ConnectionHelper{
		logger: logger,
	}
}

// SafeConnect safely connects a WhatsApp client with proper error handling
func (h *ConnectionHelper) SafeConnect(client *whatsmeow.Client, sessionID string) error {
	if err := ValidateClientAndStore(client, sessionID, h.logger); err != nil {
		return err
	}

	if client.IsConnected() {
		h.logger.Debugf("Client already connected for session %s", sessionID)
		return nil
	}

	h.logger.Infof("Connecting client for session %s", sessionID)
	return client.Connect()
}

// SafeDisconnect safely disconnects a WhatsApp client
func (h *ConnectionHelper) SafeDisconnect(client *whatsmeow.Client, sessionID string) {
	if client != nil && client.IsConnected() {
		h.logger.Infof("Disconnecting client for session %s", sessionID)
		client.Disconnect()
	}
}

// IsDeviceRegistered checks if device is registered (has ID)
func IsDeviceRegistered(client *whatsmeow.Client) bool {
	return client != nil && client.Store != nil && client.Store.ID != nil
}

// StatusHelper contains status management utilities
type StatusHelper struct {
	mu     *sync.RWMutex
	status types.Status
	logger logging.Logger
}

// NewStatusHelper creates a new status helper
func NewStatusHelper(logger logging.Logger) *StatusHelper {
	return &StatusHelper{
		mu:     &sync.RWMutex{},
		status: types.StatusDisconnected,
		logger: logger,
	}
}

// SetStatus safely sets the status with logging
func (h *StatusHelper) SetStatus(status types.Status, sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.status = status
	if status == types.StatusConnected || status == types.StatusDisconnected ||
		status == types.StatusConnecting || status == types.StatusError {
		h.logger.Infof("Session %s status: %s", sessionID, status)
	}
}

// GetStatus safely gets the current status
func (h *StatusHelper) GetStatus() types.Status {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status
}
