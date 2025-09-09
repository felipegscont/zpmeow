package meow

import (
	"context"
	"fmt"
	"strings"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"

	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// parseJID parses a JID string and adds @s.whatsapp.net if needed (like wuzapi)
func parseJID(arg string) (waTypes.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return waTypes.NewJID(arg, waTypes.DefaultUserServer), true
	} else {
		recipient, err := waTypes.ParseJID(arg)
		if err != nil {
			return recipient, false
		} else if recipient.User == "" {
			return recipient, false
		}
		return recipient, true
	}
}

// MeowServiceImpl implements the WhatsAppService interface using the new architecture
type MeowServiceImpl struct {
	clientManager *ClientManager
	logger        logger.Logger
	waLogger      waLog.Logger
}

// NewMeowService creates a new Meow WhatsApp service implementation
func NewMeowService(db *sqlx.DB, container *sqlstore.Container, waLogger waLog.Logger) session.WhatsAppService {
	if waLogger == nil {
		waLogger = waLog.Noop
	}

	appLogger := logger.GetLogger().Sub("meow-service")
	clientManager := NewClientManager(db, container, waLogger)

	service := &MeowServiceImpl{
		clientManager: clientManager,
		logger:        appLogger,
		waLogger:      waLogger,
	}

	return service
}

// StartClient starts a WhatsApp client for the given session
func (m *MeowServiceImpl) StartClient(sessionID string) error {
	ctx := context.Background()
	
	m.logger.Infof("Starting client for session %s", sessionID)
	
	err := m.clientManager.StartClient(ctx, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to start client for session %s: %v", sessionID, err)
		return fmt.Errorf("failed to start client: %w", err)
	}

	m.logger.Infof("Client started successfully for session %s", sessionID)
	return nil
}

// StopClient stops a WhatsApp client for the given session
func (m *MeowServiceImpl) StopClient(sessionID string) error {
	m.logger.Infof("Stopping client for session %s", sessionID)
	
	err := m.clientManager.StopClient(sessionID)
	if err != nil {
		m.logger.Errorf("Failed to stop client for session %s: %v", sessionID, err)
		return fmt.Errorf("failed to stop client: %w", err)
	}

	m.logger.Infof("Client stopped successfully for session %s", sessionID)
	return nil
}

// LogoutClient logs out a WhatsApp client for the given session
func (m *MeowServiceImpl) LogoutClient(sessionID string) error {
	ctx := context.Background()
	
	m.logger.Infof("Logging out client for session %s", sessionID)
	
	err := m.clientManager.LogoutClient(ctx, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to logout client for session %s: %v", sessionID, err)
		return fmt.Errorf("failed to logout client: %w", err)
	}

	m.logger.Infof("Client logged out successfully for session %s", sessionID)
	return nil
}

// GetQRCode gets the QR code for a session
func (m *MeowServiceImpl) GetQRCode(sessionID string) (string, error) {
	ctx := context.Background()
	
	m.logger.Debugf("Getting QR code for session %s", sessionID)
	
	qrCode, err := m.clientManager.GetQRCode(ctx, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to get QR code for session %s: %v", sessionID, err)
		return "", fmt.Errorf("failed to get QR code: %w", err)
	}

	if qrCode == "" {
		return "", fmt.Errorf("QR code not available yet")
	}

	m.logger.Debugf("QR code retrieved for session %s", sessionID)
	return qrCode, nil
}

// PairPhone pairs a phone number with the session
func (m *MeowServiceImpl) PairPhone(sessionID, phoneNumber string) (string, error) {
	ctx := context.Background()
	
	m.logger.Infof("Pairing phone %s with session %s", phoneNumber, sessionID)
	
	pairingCode, err := m.clientManager.PairPhone(ctx, sessionID, phoneNumber)
	if err != nil {
		m.logger.Errorf("Failed to pair phone for session %s: %v", sessionID, err)
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	m.logger.Infof("Phone pairing initiated for session %s", sessionID)
	return pairingCode, nil
}

// IsClientConnected checks if a client is connected
func (m *MeowServiceImpl) IsClientConnected(sessionID string) bool {
	connected := m.clientManager.IsClientConnected(sessionID)
	m.logger.Debugf("Client connection status for session %s: %t", sessionID, connected)
	return connected
}

// GetClientStatus gets the status of a client
func (m *MeowServiceImpl) GetClientStatus(sessionID string) types.Status {
	status := m.clientManager.GetClientStatus(sessionID)
	m.logger.Debugf("Client status for session %s: %s", sessionID, status)
	return status
}

// CheckAndUpdateDeviceJID checks and updates device JID for a session
func (m *MeowServiceImpl) CheckAndUpdateDeviceJID(sessionID string) {
	m.logger.Infof("Checking and updating device JID for session %s", sessionID)
	m.clientManager.CheckAndUpdateDeviceJID(sessionID)
}

// GetClient returns the MeowClient for a session (for advanced operations)
func (m *MeowServiceImpl) GetClient(sessionID string) (*MeowClient, bool) {
	return m.clientManager.GetClient(sessionID)
}

// GetAllClients returns all active clients
func (m *MeowServiceImpl) GetAllClients() map[string]*MeowClient {
	return m.clientManager.GetAllClients()
}

// GetStats returns service statistics
func (m *MeowServiceImpl) GetStats() map[string]interface{} {
	stats := m.clientManager.GetStats()
	m.logger.Debugf("Service stats: %+v", stats)
	return stats
}

// Cleanup performs cleanup operations
func (m *MeowServiceImpl) Cleanup() {
	m.logger.Infof("Performing service cleanup")
	m.clientManager.Cleanup()
}

// Shutdown gracefully shuts down the service
func (m *MeowServiceImpl) Shutdown(ctx context.Context) error {
	m.logger.Infof("Shutting down Meow service")
	return m.clientManager.Shutdown(ctx)
}

// Health check methods

// IsHealthy checks if the service is healthy
func (m *MeowServiceImpl) IsHealthy() bool {
	// Basic health check - could be expanded
	return m.clientManager != nil
}

// GetHealthStatus returns detailed health status
func (m *MeowServiceImpl) GetHealthStatus() map[string]interface{} {
	stats := m.GetStats()
	
	return map[string]interface{}{
		"healthy":     m.IsHealthy(),
		"stats":       stats,
		"service":     "meow-whatsapp",
		"version":     "1.0.0",
	}
}

// Advanced client operations

// SendMessage sends a message through a specific client
func (m *MeowServiceImpl) SendMessage(ctx context.Context, sessionID, to, message string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendMessage(ctx, jid, message)
}

// SendTextMessage sends a text message with optional context info
func (m *MeowServiceImpl) SendTextMessage(ctx context.Context, sessionID, to, text string, contextInfo *waE2E.ContextInfo) error {
	m.logger.Infof("DEBUG: SendTextMessage called - sessionID: %s, to: %s, text: %s", sessionID, to, text)

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		m.logger.Errorf("DEBUG: Client not found for session %s", sessionID)
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	m.logger.Infof("DEBUG: Client found for session %s", sessionID)

	// CORREÇÃO: Remover verificação IsConnected() que pode causar deadlock
	// O wuzapi não faz essa verificação e funciona corretamente
	m.logger.Infof("DEBUG: Skipping IsConnected() check to avoid deadlock")

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		m.logger.Errorf("DEBUG: Failed to parse JID %s", to)
		return fmt.Errorf("invalid JID %s", to)
	}

	m.logger.Infof("DEBUG: JID parsed successfully: %s -> %s", to, jid.String())

	m.logger.Infof("DEBUG: Calling client.SendTextMessage...")
	err := client.SendTextMessage(ctx, jid, text, contextInfo)
	if err != nil {
		m.logger.Errorf("DEBUG: client.SendTextMessage failed: %v", err)
		return err
	}

	m.logger.Infof("DEBUG: client.SendTextMessage completed successfully")
	return nil
}

// SendLocationMessage sends a location message
func (m *MeowServiceImpl) SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name, address string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendLocationMessage(ctx, jid, latitude, longitude, name, address)
}

// SendContactMessage sends a contact message
func (m *MeowServiceImpl) SendContactMessage(ctx context.Context, sessionID, to, displayName, vcard string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendContactMessage(ctx, jid, displayName, vcard)
}

// ReactToMessage sends a reaction to a message
func (m *MeowServiceImpl) ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(chatJID)
	if !ok {
		return fmt.Errorf("invalid JID %s", chatJID)
	}

	return client.ReactToMessage(ctx, jid, messageID, emoji)
}

// SetChatPresence sets the presence in a chat
func (m *MeowServiceImpl) SetChatPresence(ctx context.Context, sessionID, chatJID, presence string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse JID
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid JID %s: %w", chatJID, err)
	}

	// Convert string to presence type
	var presenceType waTypes.Presence
	switch presence {
	case "typing":
		presenceType = waTypes.PresenceAvailable // Usar available como fallback
	case "recording":
		presenceType = waTypes.PresenceAvailable // Recording não existe mais, usar available
	case "paused":
		presenceType = waTypes.PresenceUnavailable
	default:
		return fmt.Errorf("invalid presence type: %s", presence)
	}

	return client.SetChatPresence(ctx, jid, presenceType)
}

// MarkMessageRead marks messages as read
func (m *MeowServiceImpl) MarkMessageRead(ctx context.Context, sessionID, chatJID string, messageIDs []string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse JID
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid JID %s: %w", chatJID, err)
	}

	return client.MarkMessageRead(ctx, jid, messageIDs)
}

// CreateGroup creates a new WhatsApp group
func (m *MeowServiceImpl) CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse participant JIDs
	participantJIDs := make([]waTypes.JID, len(participants))
	for i, participant := range participants {
		jid, err := waTypes.ParseJID(participant)
		if err != nil {
			return nil, fmt.Errorf("invalid participant JID %s: %w", participant, err)
		}
		participantJIDs[i] = jid
	}

	return client.CreateGroup(ctx, name, participantJIDs)
}

// GetGroupInfo retrieves information about a group
func (m *MeowServiceImpl) GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.GetGroupInfo(ctx, jid)
}

// JoinGroupWithLink joins a group using an invite link
func (m *MeowServiceImpl) JoinGroupWithLink(ctx context.Context, sessionID, inviteCode string) (*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	return client.JoinGroupWithLink(ctx, inviteCode)
}

// LeaveGroup leaves a group
func (m *MeowServiceImpl) LeaveGroup(ctx context.Context, sessionID, groupJID string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.LeaveGroup(ctx, jid)
}

// GetGroupInviteLink gets the invite link for a group
func (m *MeowServiceImpl) GetGroupInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return "", fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.GetGroupInviteLink(ctx, jid, reset)
}

// UpdateGroupParticipants updates group participants (add, remove, promote, demote)
func (m *MeowServiceImpl) UpdateGroupParticipants(ctx context.Context, sessionID, groupJID string, participants []string, action string) error {
	_, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse group JID (not used in current implementation)
	_, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Parse participant JIDs
	participantJIDs := make([]waTypes.JID, len(participants))
	for i, participant := range participants {
		pJid, err := waTypes.ParseJID(participant)
		if err != nil {
			return fmt.Errorf("invalid participant JID %s: %w", participant, err)
		}
		participantJIDs[i] = pJid
	}

	// TODO: ParticipantChange types were removed in new whatsmeow version
	// Need to implement using new API methods
	return fmt.Errorf("group participant management not implemented for new whatsmeow version")
}

// SetGroupName sets the name of a group
func (m *MeowServiceImpl) SetGroupName(ctx context.Context, sessionID, groupJID, name string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.SetGroupName(ctx, jid, name)
}

// SetGroupTopic sets the topic/description of a group
func (m *MeowServiceImpl) SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	return client.SetGroupTopic(ctx, jid, topic)
}

// SendImageMessage sends an image message
func (m *MeowServiceImpl) SendImageMessage(ctx context.Context, sessionID, to string, imageData []byte, caption, mimeType string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendImageMessage(ctx, jid, imageData, caption, mimeType)
}

// SendAudioMessage sends an audio message
func (m *MeowServiceImpl) SendAudioMessage(ctx context.Context, sessionID, to string, audioData []byte, mimeType string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendAudioMessage(ctx, jid, audioData, mimeType)
}

// SendDocumentMessage sends a document message
func (m *MeowServiceImpl) SendDocumentMessage(ctx context.Context, sessionID, to string, documentData []byte, filename, caption, mimeType string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendDocumentMessage(ctx, jid, documentData, filename, caption, mimeType)
}

// SendVideoMessage sends a video message
func (m *MeowServiceImpl) SendVideoMessage(ctx context.Context, sessionID, to string, videoData []byte, caption, mimeType string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendVideoMessage(ctx, jid, videoData, caption, mimeType)
}

// SendStickerMessage sends a sticker message
func (m *MeowServiceImpl) SendStickerMessage(ctx context.Context, sessionID, to string, stickerData []byte, mimeType string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendStickerMessage(ctx, jid, stickerData, mimeType)
}

// SendButtonsMessage sends a buttons message
func (m *MeowServiceImpl) SendButtonsMessage(ctx context.Context, sessionID, to, text string, buttons []types.Button, footer string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendButtonsMessage(ctx, jid, text, buttons, footer)
}

// SendListMessage sends a list message
func (m *MeowServiceImpl) SendListMessage(ctx context.Context, sessionID, to, text, buttonText string, sections []types.Section, footer string) error {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		return fmt.Errorf("invalid JID %s", to)
	}

	return client.SendListMessage(ctx, jid, text, buttonText, sections, footer)
}

// SendPollMessage sends a poll message
func (m *MeowServiceImpl) SendPollMessage(ctx context.Context, sessionID, to, name string, options []string, selectableCount int) error {
	m.logger.Infof("DEBUG: SendPollMessage called - sessionID: %s, to: %s, name: %s", sessionID, to, name)

	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		m.logger.Errorf("DEBUG: Client not found for session %s", sessionID)
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	m.logger.Infof("DEBUG: Client found for session %s", sessionID)

	// CORREÇÃO: Remover verificação IsConnected() seguindo padrão do wuzapi

	// Parse JID using wuzapi-style parsing
	jid, ok := parseJID(to)
	if !ok {
		m.logger.Errorf("DEBUG: Failed to parse JID %s", to)
		return fmt.Errorf("invalid JID %s", to)
	}

	m.logger.Infof("DEBUG: JID parsed successfully: %s", jid.String())
	m.logger.Infof("DEBUG: Calling client.SendPollMessage...")

	err := client.SendPollMessage(ctx, jid, name, options, selectableCount)
	if err != nil {
		m.logger.Errorf("DEBUG: client.SendPollMessage failed: %v", err)
		return err
	}

	m.logger.Infof("DEBUG: SendPollMessage completed successfully")
	return nil
}

// ListGroups lists all groups for a session
func (m *MeowServiceImpl) ListGroups(ctx context.Context, sessionID string) ([]*waTypes.GroupInfo, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	// Get all groups from the client
	groups, err := client.GetClient().GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	// Convert to GroupInfo slice
	var groupInfos []*waTypes.GroupInfo
	for _, group := range groups {
		// group is already a *GroupInfo, no need to call GetGroupInfo again
		groupInfos = append(groupInfos, group)
	}

	return groupInfos, nil
}

// GetClientJID returns the JID of a connected client
func (m *MeowServiceImpl) GetClientJID(sessionID string) (string, error) {
	client, exists := m.clientManager.GetClient(sessionID)
	if !exists {
		return "", fmt.Errorf("client not found for session %s", sessionID)
	}

	jid := client.GetJID()
	if jid.IsEmpty() {
		return "", fmt.Errorf("client is not authenticated for session %s", sessionID)
	}

	return jid.String(), nil
}

// RestartClient restarts a client (stop and start)
func (m *MeowServiceImpl) RestartClient(sessionID string) error {
	m.logger.Infof("Restarting client for session %s", sessionID)
	
	// Stop the client first
	if err := m.StopClient(sessionID); err != nil {
		m.logger.Warnf("Failed to stop client during restart for session %s: %v", sessionID, err)
	}

	// Start the client
	return m.StartClient(sessionID)
}

// ConnectOnStartup connects to WhatsApp on server startup if sessions were previously connected
// This is similar to wuzapi's connectOnStartup functionality
func (m *MeowServiceImpl) ConnectOnStartup(ctx context.Context) error {
	m.logger.Infof("Starting automatic reconnection for previously connected sessions...")

	// Query sessions that were connected when the server was last shut down
	query := `
		SELECT id, name, COALESCE(device_jid, '') as device_jid, status, COALESCE(proxy_url, '') as proxy_url
		FROM sessions
		WHERE status = 'connected'
	`

	rows, err := m.clientManager.db.QueryxContext(ctx, query)
	if err != nil {
		m.logger.Errorf("Failed to query connected sessions: %v", err)
		return fmt.Errorf("failed to query connected sessions: %w", err)
	}
	defer rows.Close()

	reconnectedCount := 0
	for rows.Next() {
		var sessionID, name, deviceJID, status, proxyURL string
		err := rows.Scan(&sessionID, &name, &deviceJID, &status, &proxyURL)
		if err != nil {
			m.logger.Errorf("Failed to scan session row: %v", err)
			continue
		}

		m.logger.Infof("Attempting to reconnect session %s (name: %s, device_jid: %s)", sessionID, name, deviceJID)

		// Check if device store exists for this session
		hasDevice, err := m.hasDeviceCredentials(ctx, sessionID, deviceJID)
		if err != nil {
			m.logger.Errorf("Failed to check device credentials for session %s: %v", sessionID, err)
			continue
		}

		if !hasDevice {
			m.logger.Warnf("No device credentials found for session %s, skipping reconnection", sessionID)
			// Update status to disconnected since we can't reconnect without credentials
			m.updateSessionStatus(ctx, sessionID, "disconnected")
			continue
		}

		// Attempt to start the client
		err = m.StartClient(sessionID)
		if err != nil {
			m.logger.Errorf("Failed to start client for session %s: %v", sessionID, err)
			// Update status to error
			m.updateSessionStatus(ctx, sessionID, "error")
			continue
		}

		reconnectedCount++
		m.logger.Infof("Successfully started reconnection for session %s", sessionID)
	}

	if err := rows.Err(); err != nil {
		m.logger.Errorf("Error iterating session rows: %v", err)
		return fmt.Errorf("error iterating session rows: %w", err)
	}

	m.logger.Infof("Automatic reconnection completed. Attempted to reconnect %d sessions", reconnectedCount)
	return nil
}

// hasDeviceCredentials checks if device credentials exist for the session
func (m *MeowServiceImpl) hasDeviceCredentials(ctx context.Context, sessionID, deviceJID string) (bool, error) {
	// If we have a device JID, try to get the specific device (like wuzapi does)
	if deviceJID != "" {
		// Parse the JID to get the proper format
		jid, err := waTypes.ParseJID(deviceJID)
		if err != nil {
			m.logger.Warnf("Invalid device JID %s for session %s: %v", deviceJID, sessionID, err)
		} else {
			// Try to get the specific device
			device, err := m.clientManager.container.GetDevice(ctx, jid)
			if err == nil && device != nil && device.ID != nil {
				m.logger.Debugf("Found device credentials for session %s with JID %s", sessionID, deviceJID)
				return true, nil
			}
		}
	}

	// Fallback: check if we have any devices that could potentially be used
	// This handles cases where the JID mapping might not be perfect or the session
	// was created before proper JID tracking
	devices, err := m.clientManager.container.GetAllDevices(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get devices: %w", err)
	}

	// If we have any valid devices, we can potentially reconnect
	// The client manager will handle device selection during connection
	for _, device := range devices {
		if device != nil && device.ID != nil {
			m.logger.Debugf("Found potential device credentials for session %s", sessionID)
			return true, nil
		}
	}

	m.logger.Debugf("No device credentials found for session %s", sessionID)
	return false, nil
}

// updateSessionStatus updates the session status in the database
func (m *MeowServiceImpl) updateSessionStatus(ctx context.Context, sessionID, status string) {
	query := `UPDATE sessions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := m.clientManager.db.ExecContext(ctx, query, status, sessionID)
	if err != nil {
		m.logger.Errorf("Failed to update session %s status to %s: %v", sessionID, status, err)
	} else {
		m.logger.Debugf("Updated session %s status to %s", sessionID, status)
	}
}

// Helper functions are now using whatsmeow's built-in JID parsing
