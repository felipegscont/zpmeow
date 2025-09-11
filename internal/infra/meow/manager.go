package meow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/infra/webhook"
	"zpmeow/internal/types"

	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)


type ClientManager struct {
	mu             sync.RWMutex
	clients        map[string]*MeowClient
	db             *sqlx.DB
	container      *sqlstore.Container
	logger         logger.Logger
	waLogger       waLog.Logger
	webhookService *webhook.WebhookService
	sessionService session.SessionService
}


func NewClientManager(db *sqlx.DB, container *sqlstore.Container, waLogger waLog.Logger, webhookService *webhook.WebhookService, sessionService session.SessionService) *ClientManager {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("client-manager")

	return &ClientManager{
		clients:        make(map[string]*MeowClient),
		db:             db,
		container:      container,
		logger:         appLogger,
		waLogger:       waLogger,
		webhookService: webhookService,
		sessionService: sessionService,
	}
}


func (cm *ClientManager) GetClient(sessionID string) (*MeowClient, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	client, exists := cm.clients[sessionID]
	return client, exists
}


func (cm *ClientManager) CreateClient(ctx context.Context, sessionID string) (*MeowClient, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()


	if err := Validation.ValidateSessionID(sessionID); err != nil {
		return nil, Error.WrapError(err, "create client validation failed")
	}


	if client, exists := cm.clients[sessionID]; exists {
		return client, nil
	}


	deviceStore, err := cm.getOrCreateDeviceStore(ctx, sessionID)
	if err != nil {
		return nil, Error.WrapError(err, "failed to get device store")
	}


	client, err := NewMeowClient(sessionID, deviceStore, cm.waLogger, cm, cm.webhookService, cm.sessionService)
	if err != nil {
		return nil, Error.WrapError(err, "failed to create meow client")
	}


	cm.clients[sessionID] = client

	cm.logger.Infof("Created new client for session %s", sessionID)
	return client, nil
}


func (cm *ClientManager) RemoveClient(sessionID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, exists := cm.clients[sessionID]; exists {
		
		if client.IsConnected() {
			client.Disconnect()
		}
		delete(cm.clients, sessionID)
		cm.logger.Infof("Removed client for session %s", sessionID)
	}
}


func (cm *ClientManager) GetAllClients() map[string]*MeowClient {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	
	clients := make(map[string]*MeowClient)
	for sessionID, client := range cm.clients {
		clients[sessionID] = client
	}
	return clients
}


func (cm *ClientManager) GetClientStatus(sessionID string) types.Status {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return types.StatusDisconnected
	}
	return client.GetStatus()
}


func (cm *ClientManager) IsClientConnected(sessionID string) bool {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return false
	}
	return client.IsConnected()
}


func (cm *ClientManager) StartClient(ctx context.Context, sessionID string) error {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		var err error
		client, err = cm.CreateClient(ctx, sessionID)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
	}

	return client.Connect(ctx)
}


func (cm *ClientManager) StopClient(sessionID string) error {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return nil 
	}

	client.Disconnect()
	return nil
}


func (cm *ClientManager) LogoutClient(ctx context.Context, sessionID string) error {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return nil 
	}

	err := client.Logout(ctx)
	if err != nil {
		cm.logger.Errorf("Failed to logout client %s: %v", sessionID, err)
	}

	
	cm.RemoveClient(sessionID)
	return err
}


func (cm *ClientManager) GetQRCode(ctx context.Context, sessionID string) (string, error) {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		
		var err error
		client, err = cm.CreateClient(ctx, sessionID)
		if err != nil {
			return "", fmt.Errorf("failed to create client: %w", err)
		}
	}

	return client.GetQRCode(ctx)
}


func (cm *ClientManager) PairPhone(ctx context.Context, sessionID, phoneNumber string) (string, error) {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		
		var err error
		client, err = cm.CreateClient(ctx, sessionID)
		if err != nil {
			return "", fmt.Errorf("failed to create client: %w", err)
		}
	}

	return client.PairPhone(ctx, phoneNumber)
}


func (cm *ClientManager) Cleanup() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for sessionID, client := range cm.clients {
		if !client.IsConnected() && time.Since(client.GetLastActivity()) > 30*time.Minute {
			delete(cm.clients, sessionID)
			cm.logger.Infof("Cleaned up inactive client for session %s", sessionID)
		}
	}
}


func (cm *ClientManager) Shutdown(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.logger.Infof("Shutting down client manager...")

	for sessionID, client := range cm.clients {
		cm.logger.Infof("Disconnecting client %s", sessionID)
		client.Disconnect()
	}

	
	cm.clients = make(map[string]*MeowClient)
	cm.logger.Infof("Client manager shutdown complete")
	return nil
}


func (cm *ClientManager) getOrCreateDeviceStore(ctx context.Context, sessionID string) (*store.Device, error) {
	
	var existingDeviceJID string
	query := `SELECT COALESCE(device_jid, '') FROM sessions WHERE id = $1`
	err := cm.db.QueryRowContext(ctx, query, sessionID).Scan(&existingDeviceJID)
	if err != nil {
		cm.logger.Warnf("Failed to check existing device_jid for session %s: %v", sessionID, err)
	}

	
	if existingDeviceJID != "" {
		cm.logger.Infof("Found existing device_jid %s for session %s, attempting to retrieve device", existingDeviceJID, sessionID)

		
		if jid, err := waTypes.ParseJID(existingDeviceJID); err == nil {
			if device, err := cm.container.GetDevice(ctx, jid); err == nil && device != nil {
				cm.logger.Infof("Successfully retrieved existing device for session %s with JID %s", sessionID, existingDeviceJID)
				return device, nil
			} else {
				cm.logger.Warnf("Failed to retrieve device for session %s with JID %s: %v", sessionID, existingDeviceJID, err)
			}
		} else {
			cm.logger.Warnf("Failed to parse existing device JID %s for session %s: %v", existingDeviceJID, sessionID, err)
		}
	}

	
	deviceStore := cm.container.NewDevice()
	cm.logger.Infof("Created new unique device store for session %s", sessionID)
	return deviceStore, nil
}


func (cm *ClientManager) OnClientStatusChange(sessionID string, status types.Status) {
	cm.logger.Infof("Client %s status changed to %s", sessionID, status)

	
	go cm.updateSessionInDatabase(sessionID, status)

	
	if status == types.StatusConnected {
		go cm.CheckAndUpdateDeviceJID(sessionID)
	}
}


func (cm *ClientManager) updateSessionInDatabase(sessionID string, status types.Status) {
	ctx := context.Background()

	
	query := `UPDATE sessions SET status = $1, updated_at = NOW() WHERE id = $2`
	args := []interface{}{string(status), sessionID}

	_, err := cm.db.ExecContext(ctx, query, args...)
	if err != nil {
		cm.logger.Errorf("Failed to update session %s status to %s: %v", sessionID, status, err)
	} else {
		cm.logger.Infof("Successfully updated session %s status to %s", sessionID, status)
	}
}


func (cm *ClientManager) updateSessionDeviceJID(sessionID string, deviceJID string) {
	ctx := context.Background()

	cm.logger.Infof("Attempting to update device_jid for session %s with JID %s", sessionID, deviceJID)

	
	query := `UPDATE sessions SET device_jid = $1, status = $2, qr_code = '', updated_at = NOW() WHERE id = $3`
	args := []interface{}{deviceJID, string(types.StatusConnected), sessionID}

	result, err := cm.db.ExecContext(ctx, query, args...)
	if err != nil {
		cm.logger.Errorf("Failed to update device_jid for session %s: %v", sessionID, err)
	} else {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			cm.logger.Warnf("No rows affected when updating device_jid for session %s - session may not exist", sessionID)
		} else {
			cm.logger.Infof("Successfully registered device JID %s for session %s and cleared QR code", deviceJID, sessionID)
		}
	}
}


func (cm *ClientManager) clearQRCodeInDatabase(sessionID string) {
	ctx := context.Background()

	cm.logger.Infof("Clearing QR code in database for session %s", sessionID)

	query := `UPDATE sessions SET qr_code = '', updated_at = NOW() WHERE id = $1`
	_, err := cm.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		cm.logger.Errorf("Failed to clear QR code for session %s: %v", sessionID, err)
	} else {
		cm.logger.Infof("Successfully cleared QR code for session %s", sessionID)
	}
}


func (cm *ClientManager) OnPairSuccess(sessionID string, deviceJID string) {
	cm.logger.Infof("OnPairSuccess called for session %s with device JID %s", sessionID, deviceJID)


	if err := Validation.ValidateSessionID(sessionID); err != nil {
		cm.logger.Errorf("OnPairSuccess validation failed: %v", err)
		return
	}


	if err := Validation.ValidateDeviceJID(deviceJID); err != nil {
		cm.logger.Errorf("OnPairSuccess validation failed for session %s: %v", sessionID, err)
		return
	}


	go cm.updateSessionDeviceJID(sessionID, deviceJID)
}


func (cm *ClientManager) CheckAndUpdateDeviceJID(sessionID string) {
	client, exists := cm.clients[sessionID]
	if !exists {
		cm.logger.Debugf("Session %s not found in client manager", sessionID)
		return
	}

	if client.client == nil || client.client.Store.ID == nil {
		cm.logger.Debugf("Session %s: Client or Store.ID is nil", sessionID)
		return
	}

	if client.status != types.StatusConnected {
		cm.logger.Debugf("Session %s: Not connected, skipping device JID check", sessionID)
		return
	}

	
	ctx := context.Background()
	var deviceJID string
	query := `SELECT COALESCE(device_jid, '') FROM sessions WHERE id = $1`
	err := cm.db.QueryRowContext(ctx, query, sessionID).Scan(&deviceJID)
	if err != nil {
		cm.logger.Errorf("Failed to check device_jid for session %s: %v", sessionID, err)
		return
	}

	if deviceJID == "" {
		
		newDeviceJID := client.client.Store.ID.String()
		cm.logger.Infof("Session %s: Found missing device_jid, updating to %s", sessionID, newDeviceJID)
		go cm.updateSessionDeviceJID(sessionID, newDeviceJID)
	}
}


func (cm *ClientManager) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	connected := 0
	disconnected := 0

	for _, client := range cm.clients {
		if client.IsConnected() {
			connected++
		} else {
			disconnected++
		}
	}

	return map[string]interface{}{
		"total_clients":        len(cm.clients),
		"connected_clients":    connected,
		"disconnected_clients": disconnected,
	}
}
