package core

import (
	"context"
	"fmt"
	"sync"

	"zpmeow/internal/infra/logger"

	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// ClientManager manages WhatsApp clients for multiple sessions
type ClientManager struct {
	mu        sync.RWMutex
	clients   map[string]*MeowClient
	container *sqlstore.Container
	logger    logger.Logger
	waLogger  waLog.Logger
}

// NewClientManager creates a new client manager
func NewClientManager(container *sqlstore.Container, waLogger waLog.Logger) *ClientManager {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("client-manager")

	return &ClientManager{
		clients:   make(map[string]*MeowClient),
		container: container,
		logger:    appLogger,
		waLogger:  waLogger,
	}
}

// GetClient gets a client by session ID
func (m *ClientManager) GetClient(sessionID string) *MeowClient {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if client, exists := m.clients[sessionID]; exists {
		return client
	}
	return nil
}

// CreateClient creates a new client for a session
func (m *ClientManager) CreateClient(sessionID string, eventHandler EventHandler) (*MeowClient, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if client already exists
	if existingClient, exists := m.clients[sessionID]; exists {
		m.logger.Warnf("Client already exists for session %s", sessionID)
		return existingClient, nil
	}

	// Get device store for this session
	// Create a JID for the session (using session ID as user part)
	sessionJID := types.NewJID(sessionID, types.DefaultUserServer)
	deviceStore, err := m.container.GetDevice(context.Background(), sessionJID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device store for session %s: %w", sessionID, err)
	}

	// Create new MeowClient
	client, err := NewMeowClient(sessionID, deviceStore, m.waLogger, eventHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for session %s: %w", sessionID, err)
	}

	// Store client
	m.clients[sessionID] = client
	m.logger.Infof("Created client for session %s", sessionID)

	return client, nil
}

// RemoveClient removes a client and cleans up resources
func (m *ClientManager) RemoveClient(sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[sessionID]
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Cleanup client resources
	client.Cleanup()

	// Remove from map
	delete(m.clients, sessionID)
	m.logger.Infof("Removed client for session %s", sessionID)

	return nil
}

// GetAllClients returns a copy of all clients
func (m *ClientManager) GetAllClients() map[string]*MeowClient {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	clientsCopy := make(map[string]*MeowClient)
	for sessionID, client := range m.clients {
		clientsCopy[sessionID] = client
	}

	return clientsCopy
}

// HasClient checks if a client exists for the given session
func (m *ClientManager) HasClient(sessionID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.clients[sessionID]
	return exists
}

// GetClientCount returns the number of active clients
func (m *ClientManager) GetClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.clients)
}

// ConnectClient connects a client for the given session
func (m *ClientManager) ConnectClient(ctx context.Context, sessionID string, eventHandler EventHandler) error {
	// Get or create client
	client := m.GetClient(sessionID)
	if client == nil {
		var err error
		client, err = m.CreateClient(sessionID, eventHandler)
		if err != nil {
			return fmt.Errorf("failed to create client for session %s: %w", sessionID, err)
		}
	}

	// Connect the client
	return client.Connect(ctx)
}

// DisconnectClient disconnects a client for the given session
func (m *ClientManager) DisconnectClient(sessionID string) error {
	client := m.GetClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	return client.Disconnect()
}

// Shutdown gracefully shuts down all clients
func (m *ClientManager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Infof("Shutting down client manager with %d clients", len(m.clients))

	// Disconnect all clients
	for sessionID, client := range m.clients {
		m.logger.Infof("Shutting down client for session %s", sessionID)
		if err := client.Disconnect(); err != nil {
			m.logger.Errorf("Failed to disconnect client for session %s: %v", sessionID, err)
		}
		client.Cleanup()
	}

	// Clear clients map
	m.clients = make(map[string]*MeowClient)

	m.logger.Infof("Client manager shutdown complete")
	return nil
}
