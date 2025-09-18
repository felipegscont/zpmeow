package wmeow

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"meow/internal/domain/session"
	"meow/internal/infrastructure/logging"
	"meow/internal/shared/types"

	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
)

type Helper struct {
	sessionRepo session.Repository
	logger      logging.Logger
}

func NewHelper(sessionRepo session.Repository, logger logging.Logger) *Helper {
	return &Helper{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (h *Helper) UpdateSessionStatus(sessionID string, status session.Status) {
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

func (h *Helper) UpdateSessionQRCode(sessionID string, qrCode string) {
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

func ValidateClientAndStore(client *whatsmeow.Client, sessionID string, logger logging.Logger) error {
	if client == nil {
		return fmt.Errorf("meow client is nil for session %s", sessionID)
	}

	if client.Store == nil {
		return fmt.Errorf("meow client store is nil for session %s", sessionID)
	}

	return nil
}

func GetOrCreateDeviceStore(sessionID string, container *sqlstore.Container) *store.Device {
	return GetDeviceStoreForSession(sessionID, "", container)
}

func GetDeviceStoreForSession(sessionID, expectedDeviceJID string, container *sqlstore.Container) *store.Device {
	var deviceStore *store.Device
	var err error

	if expectedDeviceJID != "" {
		jid, ok := parseJID(expectedDeviceJID)
		if ok {
			deviceStore, err = container.GetDevice(context.Background(), jid)
			if err != nil {
				fmt.Printf("Failed to get device for JID %s: %v\n", expectedDeviceJID, err)
				deviceStore = container.NewDevice()
			} else {
				fmt.Printf("Successfully retrieved existing device for JID %s\n", expectedDeviceJID)
			}
		} else {
			fmt.Printf("Failed to parse JID %s, creating new device\n", expectedDeviceJID)
			deviceStore = container.NewDevice()
		}
	} else {
		fmt.Printf("No device JID provided for session %s, creating new device\n", sessionID)
		deviceStore = container.NewDevice()
	}

	if deviceStore == nil {
		fmt.Printf("Device store is nil, creating fallback device\n")
		deviceStore = container.NewDevice()
	}

	return deviceStore
}

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

type QRHelper struct {
	logger logging.Logger
}

func NewQRHelper(logger logging.Logger) *QRHelper {
	return &QRHelper{
		logger: logger,
	}
}

func (h *QRHelper) GenerateQRCodeImage(qrText string) string {
	qrPNG, err := qrcode.Encode(qrText, qrcode.Medium, 256)
	if err != nil {
		h.logger.Errorf("Failed to generate QR code image: %v", err)
		return ""
	}

	base64Str := base64.StdEncoding.EncodeToString(qrPNG)
	return "data:image/png;base64," + base64Str
}

func (h *QRHelper) DisplayQRCodeInTerminal(qrCode, sessionID string) {
	fmt.Printf("\n=== QR Code for Session %s ===\n", sessionID)
	qrterminal.GenerateHalfBlock(qrCode, qrterminal.L, os.Stdout)
	fmt.Printf("QR Code String: %s\n", qrCode)
	fmt.Printf("=== End QR Code ===\n\n")
}

type ConnHelper struct {
	logger logging.Logger
}

func NewConnHelper(logger logging.Logger) *ConnHelper {
	return &ConnHelper{
		logger: logger,
	}
}

func (h *ConnHelper) SafeConnect(client *whatsmeow.Client, sessionID string) error {
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

func (h *ConnHelper) SafeDisconnect(client *whatsmeow.Client, sessionID string) {
	if client != nil && client.IsConnected() {
		h.logger.Infof("Disconnecting client for session %s", sessionID)
		client.Disconnect()
	}
}

func IsDeviceRegistered(client *whatsmeow.Client) bool {
	return client != nil && client.Store != nil && client.Store.ID != nil
}

type StatusHelper struct {
	mu     *sync.RWMutex
	status types.Status
	logger logging.Logger
}

func NewStatusHelper(logger logging.Logger) *StatusHelper {
	return &StatusHelper{
		mu:     &sync.RWMutex{},
		status: types.StatusDisconnected,
		logger: logger,
	}
}

func (h *StatusHelper) SetStatus(status types.Status, sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.status = status
	if status == types.StatusConnected || status == types.StatusDisconnected ||
		status == types.StatusConnecting || status == types.StatusError {
		h.logger.Infof("Session %-36s status: %-10s", sessionID, status)
	}
}

func (h *StatusHelper) GetStatus() types.Status {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status
}
