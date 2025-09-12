package core

// ClientManager manages WhatsApp clients
type ClientManager struct {
	clients map[string]*MeowClient
}

// NewClientManager creates a new client manager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*MeowClient),
	}
}

// GetClient gets a client by session ID
func (m *ClientManager) GetClient(sessionID string) *MeowClient {
	if client, exists := m.clients[sessionID]; exists {
		return client
	}
	return nil
}

// CreateClient creates a new client
func (m *ClientManager) CreateClient(sessionID string) *MeowClient {
	client := NewMeowClient(sessionID)
	m.clients[sessionID] = client
	return client
}

// RemoveClient removes a client
func (m *ClientManager) RemoveClient(sessionID string) {
	delete(m.clients, sessionID)
}

// GetAllClients returns all clients
func (m *ClientManager) GetAllClients() map[string]*MeowClient {
	return m.clients
}
