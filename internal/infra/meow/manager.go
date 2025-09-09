package meow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"

	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// ClientManager manages multiple WhatsApp clients
type ClientManager struct {
	mu        sync.RWMutex
	clients   map[string]*MeowClient
	db        *sqlx.DB
	container *sqlstore.Container
	logger    logger.Logger
	waLogger  waLog.Logger
}

// NewClientManager creates a new client manager
func NewClientManager(db *sqlx.DB, container *sqlstore.Container, waLogger waLog.Logger) *ClientManager {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("client-manager")

	return &ClientManager{
		clients:   make(map[string]*MeowClient),
		db:        db,
		container: container,
		logger:    appLogger,
		waLogger:  waLogger,
	}
}

// GetClient retrieves a client by session ID
func (cm *ClientManager) GetClient(sessionID string) (*MeowClient, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	client, exists := cm.clients[sessionID]
	return client, exists
}

// CreateClient creates a new WhatsApp client for the session
func (cm *ClientManager) CreateClient(ctx context.Context, sessionID string) (*MeowClient, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Check if client already exists
	if client, exists := cm.clients[sessionID]; exists {
		return client, nil
	}

	// Get or create device store for this session
	deviceStore, err := cm.getOrCreateDeviceStore(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device store: %w", err)
	}

	// Create new MeowClient
	client, err := NewMeowClient(sessionID, deviceStore, cm.waLogger, cm)
	if err != nil {
		return nil, fmt.Errorf("failed to create meow client: %w", err)
	}

	// Store client
	cm.clients[sessionID] = client

	cm.logger.Infof("Created new client for session %s", sessionID)
	return client, nil
}

// RemoveClient removes a client from the manager
func (cm *ClientManager) RemoveClient(sessionID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, exists := cm.clients[sessionID]; exists {
		// Disconnect the client if it's connected
		if client.IsConnected() {
			client.Disconnect()
		}
		delete(cm.clients, sessionID)
		cm.logger.Infof("Removed client for session %s", sessionID)
	}
}

// GetAllClients returns all active clients
func (cm *ClientManager) GetAllClients() map[string]*MeowClient {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Create a copy to avoid race conditions
	clients := make(map[string]*MeowClient)
	for sessionID, client := range cm.clients {
		clients[sessionID] = client
	}
	return clients
}

// GetClientStatus returns the status of a client
func (cm *ClientManager) GetClientStatus(sessionID string) types.Status {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return types.StatusDisconnected
	}
	return client.GetStatus()
}

// IsClientConnected checks if a client is connected
func (cm *ClientManager) IsClientConnected(sessionID string) bool {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return false
	}
	return client.IsConnected()
}

// StartClient starts a WhatsApp client
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

// StopClient stops a WhatsApp client
func (cm *ClientManager) StopClient(sessionID string) error {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return nil // Client doesn't exist, nothing to stop
	}

	client.Disconnect()
	return nil
}

// LogoutClient logs out and removes a WhatsApp client
func (cm *ClientManager) LogoutClient(ctx context.Context, sessionID string) error {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		return nil // Client doesn't exist, nothing to logout
	}

	err := client.Logout(ctx)
	if err != nil {
		cm.logger.Errorf("Failed to logout client %s: %v", sessionID, err)
	}

	// Remove client from manager
	cm.RemoveClient(sessionID)
	return err
}

// GetQRCode gets the QR code for a session
func (cm *ClientManager) GetQRCode(ctx context.Context, sessionID string) (string, error) {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		// Create client if it doesn't exist
		var err error
		client, err = cm.CreateClient(ctx, sessionID)
		if err != nil {
			return "", fmt.Errorf("failed to create client: %w", err)
		}
	}

	return client.GetQRCode(ctx)
}

// PairPhone pairs a phone number with the session
func (cm *ClientManager) PairPhone(ctx context.Context, sessionID, phoneNumber string) (string, error) {
	client, exists := cm.GetClient(sessionID)
	if !exists {
		// Create client if it doesn't exist
		var err error
		client, err = cm.CreateClient(ctx, sessionID)
		if err != nil {
			return "", fmt.Errorf("failed to create client: %w", err)
		}
	}

	return client.PairPhone(ctx, phoneNumber)
}

// Cleanup removes inactive clients
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

// Shutdown gracefully shuts down all clients
func (cm *ClientManager) Shutdown(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.logger.Infof("Shutting down client manager...")

	for sessionID, client := range cm.clients {
		cm.logger.Infof("Disconnecting client %s", sessionID)
		client.Disconnect()
	}

	// Clear all clients
	cm.clients = make(map[string]*MeowClient)
	cm.logger.Infof("Client manager shutdown complete")
	return nil
}

// getOrCreateDeviceStore gets or creates a device store for the session
func (cm *ClientManager) getOrCreateDeviceStore(ctx context.Context, sessionID string) (*store.Device, error) {
	// Try to get existing device first
	devices, err := cm.container.GetAllDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	// Look for existing device with matching session ID in the device name or similar
	for _, device := range devices {
		// This is a simple approach - in production you might want to store
		// session ID mapping in a separate table
		if device != nil {
			return device, nil
		}
	}

	// If no existing device found, create a new one
	deviceStore := cm.container.NewDevice()
	cm.logger.Infof("Created new device store for session %s", sessionID)
	return deviceStore, nil
}

// OnClientStatusChange is called when a client's status changes
func (cm *ClientManager) OnClientStatusChange(sessionID string, status types.Status) {
	cm.logger.Infof("Client %s status changed to %s", sessionID, status)

	// Update database when status changes
	go cm.updateSessionInDatabase(sessionID, status)

	// If status is connected, check and update device JID if missing
	if status == types.StatusConnected {
		go cm.CheckAndUpdateDeviceJID(sessionID)
	}
}

// updateSessionInDatabase updates the session status in the database
func (cm *ClientManager) updateSessionInDatabase(sessionID string, status types.Status) {
	ctx := context.Background()

	// Update only status (device_jid is updated separately on PairSuccess)
	query := `UPDATE sessions SET status = $1, updated_at = NOW() WHERE id = $2`
	args := []interface{}{string(status), sessionID}

	_, err := cm.db.ExecContext(ctx, query, args...)
	if err != nil {
		cm.logger.Errorf("Failed to update session %s status to %s: %v", sessionID, status, err)
	} else {
		cm.logger.Infof("Successfully updated session %s status to %s", sessionID, status)
	}
}

// updateSessionDeviceJID updates the device_jid for a specific session after successful pairing
func (cm *ClientManager) updateSessionDeviceJID(sessionID string, deviceJID string) {
	ctx := context.Background()

	cm.logger.Infof("Attempting to update device_jid for session %s with JID %s", sessionID, deviceJID)

	query := `UPDATE sessions SET device_jid = $1, status = $2, updated_at = NOW() WHERE id = $3`
	args := []interface{}{deviceJID, string(types.StatusConnected), sessionID}

	result, err := cm.db.ExecContext(ctx, query, args...)
	if err != nil {
		cm.logger.Errorf("Failed to update device_jid for session %s: %v", sessionID, err)
	} else {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			cm.logger.Warnf("No rows affected when updating device_jid for session %s - session may not exist", sessionID)
		} else {
			cm.logger.Infof("Successfully registered device JID %s for session %s", deviceJID, sessionID)
		}
	}
}

// OnPairSuccess is called when a session successfully pairs with WhatsApp
func (cm *ClientManager) OnPairSuccess(sessionID string, deviceJID string) {
	cm.logger.Infof("OnPairSuccess called for session %s with device JID %s", sessionID, deviceJID)

	if sessionID == "" {
		cm.logger.Errorf("OnPairSuccess called with empty session ID")
		return
	}

	if deviceJID == "" {
		cm.logger.Errorf("OnPairSuccess called with empty device JID for session %s", sessionID)
		return
	}

	// Update the database with the device JID
	go cm.updateSessionDeviceJID(sessionID, deviceJID)
}

// CheckAndUpdateDeviceJID checks if a connected session has a missing device_jid and updates it
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

	// Check if device_jid is missing in database
	ctx := context.Background()
	var deviceJID string
	query := `SELECT COALESCE(device_jid, '') FROM sessions WHERE id = $1`
	err := cm.db.QueryRowContext(ctx, query, sessionID).Scan(&deviceJID)
	if err != nil {
		cm.logger.Errorf("Failed to check device_jid for session %s: %v", sessionID, err)
		return
	}

	if deviceJID == "" {
		// Device JID is missing, update it
		newDeviceJID := client.client.Store.ID.String()
		cm.logger.Infof("Session %s: Found missing device_jid, updating to %s", sessionID, newDeviceJID)
		go cm.updateSessionDeviceJID(sessionID, newDeviceJID)
	}
}

// GetStats returns statistics about the client manager
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
		"total_clients":       len(cm.clients),
		"connected_clients":   connected,
		"disconnected_clients": disconnected,
	}
}
