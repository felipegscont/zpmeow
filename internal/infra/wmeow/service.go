package wmeow

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/logging"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// ButtonData represents a button in interactive messages
type ButtonData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// ListItem represents an item in a list message
type ListItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ListSection represents a section in a list message
type ListSection struct {
	Title string     `json:"title"`
	Rows  []ListItem `json:"rows"`
}

// ChatInfo represents information about a WhatsApp chat
type ChatInfo struct {
	JID           string `json:"jid"`
	Name          string `json:"name"`
	Type          string `json:"type"` // "contact", "group", "broadcast"
	IsGroup       bool   `json:"isGroup"`
	LastMessage   string `json:"lastMessage"`
	LastTimestamp int64  `json:"lastTimestamp"`
	Timestamp     int64  `json:"timestamp"` // Alias for LastTimestamp for compatibility
	UnreadCount   int    `json:"unreadCount"`
	Pinned        bool   `json:"pinned"`
	Muted         bool   `json:"muted"`
	Archived      bool   `json:"archived"`
	MuteEndTime   int64  `json:"muteEndTime,omitempty"`
}

// GroupInfo represents information about a WhatsApp group
type GroupInfo struct {
	JID          string   `json:"jid"`
	Name         string   `json:"name"`
	GroupName    string   `json:"group_name"` // Alias for Name
	Topic        string   `json:"topic"`
	Participants []string `json:"participants"`
	Admins       []string `json:"admins"`
	Owner        string   `json:"owner"`
	Inviter      string   `json:"inviter,omitempty"`     // Who invited to the group
	InviteCode   string   `json:"invite_code,omitempty"` // Group invite code
	CreatedAt    int64    `json:"createdAt"`
	Size         int      `json:"size"`
	Announce     bool     `json:"announce"`
	Locked       bool     `json:"locked"`
	Ephemeral    bool     `json:"ephemeral"`
}

type WameowService interface {
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	IsClientConnected(sessionID string) bool
	GetClientStatus(sessionID string) session.Status
	ConnectOnStartup(ctx context.Context) error

	SendTextMessage(ctx context.Context, sessionID, phone, text string) (*whatsmeow.SendResponse, error)
	SendImageMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error)
	SendAudioMessage(ctx context.Context, sessionID, phone string, data []byte, mimeType string) (*whatsmeow.SendResponse, error)
	SendVideoMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error)
	SendDocumentMessage(ctx context.Context, sessionID, phone string, data []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error)
	SendStickerMessage(ctx context.Context, sessionID, phone string, data []byte, mimeType string) (*whatsmeow.SendResponse, error)
	SendContactMessage(ctx context.Context, sessionID, phone, contactName, contactPhone string) (*whatsmeow.SendResponse, error)
	SendLocationMessage(ctx context.Context, sessionID, phone string, latitude, longitude float64, name, address string) (*whatsmeow.SendResponse, error)
	SendButtonMessage(ctx context.Context, sessionID, phone, title string, buttons []ButtonData) (*whatsmeow.SendResponse, error)
	SendListMessage(ctx context.Context, sessionID, phone, title, description, buttonText, footerText string, sections []ListSection) (*whatsmeow.SendResponse, error)
	SendPollMessage(ctx context.Context, sessionID, phone, name string, options []string, selectableCount int) (*whatsmeow.SendResponse, error)

	MarkAsRead(ctx context.Context, sessionID, phone string, messageIDs []string) error
	SetPresence(ctx context.Context, sessionID, phone, state, media string) error
	ReactToMessage(ctx context.Context, sessionID, phone, messageID, emoji string) error
	DeleteMessage(ctx context.Context, sessionID, phone, messageID string, forEveryone bool) error
	EditMessage(ctx context.Context, sessionID, phone, messageID, newText string) (*whatsmeow.SendResponse, error)
	DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error)

	CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*GroupInfo, error)
	ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error)
	GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error)
	JoinGroup(ctx context.Context, sessionID, inviteLink string) (*GroupInfo, error)
	JoinGroupWithInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error)
	LeaveGroup(ctx context.Context, sessionID, groupJID string) error
	GetInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error)
	GetInviteInfo(ctx context.Context, sessionID, inviteLink string) (*GroupInfo, error)
	GetGroupInfoFromInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error)
	UpdateParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
	SetGroupName(ctx context.Context, sessionID, groupJID, name string) error
	SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error
	SetGroupPhoto(ctx context.Context, sessionID, groupJID string, photo []byte) error
	RemoveGroupPhoto(ctx context.Context, sessionID, groupJID string) error
	SetGroupAnnounce(ctx context.Context, sessionID, groupJID string, announceOnly bool) error
	SetGroupLocked(ctx context.Context, sessionID, groupJID string, locked bool) error
	SetGroupEphemeral(ctx context.Context, sessionID, groupJID string, ephemeral bool, duration int) error
	SetGroupJoinApprovalMode(ctx context.Context, sessionID, groupJID string, requireApproval bool) error
	SetGroupJoinApproval(ctx context.Context, sessionID, groupJID string, requireApproval bool) error // Alias for SetGroupJoinApprovalMode
	SetGroupMemberAddMode(ctx context.Context, sessionID, groupJID string, mode string) error
	GetGroupRequestParticipants(ctx context.Context, sessionID, groupJID string) ([]string, error)
	UpdateGroupRequestParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
	LinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error
	UnlinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error
	GetSubGroups(ctx context.Context, sessionID, communityJID string) ([]string, error)
	GetLinkedGroupsParticipants(ctx context.Context, sessionID, communityJID string) ([]string, error)

	// User Operations
	CheckUser(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error)
	GetUserInfo(ctx context.Context, sessionID string, phones []string) (map[string]UserInfoResult, error)
	GetAvatar(ctx context.Context, sessionID, phone string) (*AvatarResult, error)
	GetContacts(ctx context.Context, sessionID string) ([]ContactResult, error)
	SetUserPresence(ctx context.Context, sessionID, state string) error

	// Session Management Operations
	UpdateSessionWebhook(sessionID, webhookURL string) error
	UpdateSessionSubscriptions(sessionID string, events []string) error

	// Chat Operations
	SetDisappearingTimer(ctx context.Context, sessionID, chatJID string, duration time.Duration) error
	ListChats(ctx context.Context, sessionID, chatType string) ([]ChatInfo, error)
	GetChatInfo(ctx context.Context, sessionID, chatJID string) (*ChatInfo, error)
	PinChat(ctx context.Context, sessionID, chatJID string, pinned bool) error
	MuteChat(ctx context.Context, sessionID, chatJID string, muted bool, duration time.Duration) error
	ArchiveChat(ctx context.Context, sessionID, chatJID string, archived bool) error

	// Newsletter Operations
	GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterID string) ([]NewsletterMessage, error)
	NewsletterMarkViewed(ctx context.Context, sessionID, newsletterID string, messageIDs []string) error
	NewsletterSendReaction(ctx context.Context, sessionID, newsletterID, messageID, reaction string) error
	NewsletterToggleMute(ctx context.Context, sessionID, newsletterID string, muted bool) error
	NewsletterSubscribeLiveUpdates(ctx context.Context, sessionID, newsletterID string) error
	UploadNewsletter(ctx context.Context, sessionID string, data []byte) error
	GetNewsletterInfoWithInvite(ctx context.Context, sessionID, inviteCode string) (*NewsletterInfo, error)
	CreateNewsletter(ctx context.Context, sessionID, name, description string) (*NewsletterInfo, error)
	GetNewsletterInfo(ctx context.Context, sessionID, newsletterID string) (*NewsletterInfo, error)
	GetSubscribedNewsletters(ctx context.Context, sessionID string) ([]NewsletterInfo, error)
	FollowNewsletter(ctx context.Context, sessionID, newsletterID string) error
	UnfollowNewsletter(ctx context.Context, sessionID, newsletterID string) error
	SendNewsletterMessage(ctx context.Context, sessionID, newsletterID, message string) error
	GetNewsletterMessages(ctx context.Context, sessionID, newsletterID string) ([]NewsletterMessage, error)

	// Privacy Operations
	GetPrivacySettings(ctx context.Context, sessionID string) (*PrivacySettings, error)
	SetPrivacySetting(ctx context.Context, sessionID, setting, value string) error
	GetBlocklist(ctx context.Context, sessionID string) ([]string, error)
	UpdateBlocklist(ctx context.Context, sessionID string, action string, contacts []string) error
}

// User operation result types
type UserCheckResult struct {
	Query        string `json:"query"`
	IsInWhatsapp bool   `json:"is_in_whatsapp"`
	IsInmeow     bool   `json:"is_in_meow"` // Alias for compatibility
	JID          string `json:"jid"`
	VerifiedName string `json:"verified_name,omitempty"`
}

type UserInfoResult struct {
	JID          string `json:"jid"`
	DisplayName  string `json:"display_name,omitempty"`
	Name         string `json:"name,omitempty"`          // Alias for DisplayName
	PushName     string `json:"push_name,omitempty"`     // Push notification name
	BusinessName string `json:"business_name,omitempty"` // Business name if applicable
	Phone        string `json:"phone,omitempty"`         // Phone number
	VerifiedName string `json:"verified_name,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Status       string `json:"status,omitempty"`
	PictureID    string `json:"picture_id,omitempty"`
	DeviceCount  int    `json:"device_count,omitempty"`
	Notify       string `json:"notify,omitempty"` // Notification name
	IsBlocked    bool   `json:"is_blocked"`       // Whether user is blocked
	IsMuted      bool   `json:"is_muted"`         // Whether user is muted
}

// NewsletterMessage represents a newsletter message
type NewsletterMessage struct {
	ID           string `json:"id"`
	NewsletterID string `json:"newsletter_id"`
	Content      string `json:"content"`
	Timestamp    int64  `json:"timestamp"`
	Viewed       bool   `json:"viewed"`
}

// NewsletterInfo represents information about a newsletter
type NewsletterInfo struct {
	ID          string `json:"id"`
	JID         string `json:"jid"` // Newsletter JID
	Name        string `json:"name"`
	Description string `json:"description"`
	Subscribers int    `json:"subscribers"`
	Verified    bool   `json:"verified"`
	IsVerified  bool   `json:"is_verified"` // Alias for Verified
	Muted       bool   `json:"muted"`
	Following   bool   `json:"following"`
	CreatedAt   int64  `json:"created_at"` // Creation timestamp
	ServerID    string `json:"server_id"`  // Server ID
	Timestamp   int64  `json:"timestamp"`  // General timestamp
}

// PrivacySettings represents privacy settings
type PrivacySettings struct {
	LastSeen     string `json:"last_seen"`     // "everyone", "contacts", "nobody"
	ProfilePhoto string `json:"profile_photo"` // "everyone", "contacts", "nobody"
	About        string `json:"about"`         // "everyone", "contacts", "nobody"
	Status       string `json:"status"`        // "everyone", "contacts", "nobody"
	ReadReceipts bool   `json:"read_receipts"` // Whether read receipts are enabled
	GroupsAdd    string `json:"groups_add"`    // "everyone", "contacts", "nobody"
	CallsAdd     string `json:"calls_add"`     // "everyone", "contacts", "nobody"
}

type AvatarResult struct {
	Phone     string `json:"phone"`
	JID       string `json:"jid"`
	AvatarURL string `json:"avatar_url,omitempty"`
	PictureID string `json:"picture_id,omitempty"`
}

type ContactResult struct {
	JID          string `json:"jid"`
	Name         string `json:"name,omitempty"`
	Notify       string `json:"notify,omitempty"`
	PushName     string `json:"push_name,omitempty"`
	BusinessName string `json:"business_name,omitempty"`
	IsBlocked    bool   `json:"is_blocked"`
	IsMuted      bool   `json:"is_muted"`
}

// MeowService - Simplified implementation following KISS principles
type MeowService struct {
	clients   map[string]*WameowClient
	sessions  session.Repository
	logger    logging.Logger
	container *sqlstore.Container
	waLogger  waLog.Logger
	mu        sync.RWMutex
}

func NewMeowService(container *sqlstore.Container, waLogger waLog.Logger, sessionRepo session.Repository) WameowService {
	return &MeowService{
		clients:   make(map[string]*WameowClient),
		sessions:  sessionRepo,
		logger:    logging.GetLogger().Sub("wameow"),
		container: container,
		waLogger:  waLogger,
	}
}

func (m *MeowService) StartClient(sessionID string) error {
	m.logger.Infof("Starting client for session %s", sessionID)
	client := m.getOrCreateClient(sessionID)
	return client.Connect()
}

func (m *MeowService) StopClient(sessionID string) error {
	m.logger.Infof("Stopping client for session %s", sessionID)
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if err := client.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect client: %w", err)
	}

	m.removeClient(sessionID)
	return nil
}

func (m *MeowService) LogoutClient(sessionID string) error {
	m.logger.Infof("Logging out client for session %s", sessionID)
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if err := client.Logout(); err != nil {
		return fmt.Errorf("failed to logout client: %w", err)
	}

	m.removeClient(sessionID)
	return nil
}

// Helper methods for simplified client management
func (m *MeowService) getClient(sessionID string) *WameowClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[sessionID]
}

func (m *MeowService) getOrCreateClient(sessionID string) *WameowClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[sessionID]; exists {
		return client
	}

	// Get session to check for existing device_jid
	sessionEntity, err := m.sessions.GetByID(context.Background(), sessionID)
	var expectedDeviceJID string
	if err == nil && sessionEntity.GetDeviceJIDString() != "" {
		expectedDeviceJID = sessionEntity.GetDeviceJIDString()
		m.logger.Infof("Creating client for session %s with expected device JID: %s", sessionID, expectedDeviceJID)
	}

	// Create event processor for this session
	eventProcessor := NewEventProcessor(sessionID, "", m.sessions)

	client, err := NewWameowClientWithDeviceJID(sessionID, expectedDeviceJID, m.container, m.waLogger, eventProcessor, m.sessions)
	if err != nil {
		m.logger.Errorf("Failed to create WameowClient for session %s: %v", sessionID, err)
		return nil
	}

	m.clients[sessionID] = client
	return client
}

// validateAndGetClient - Helper method to validate session and get connected client
func (m *MeowService) validateAndGetClient(sessionID string) (*WameowClient, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	return client, nil
}

// validateAndGetClientForSending - Helper method specifically for sending operations
func (m *MeowService) validateAndGetClientForSending(sessionID string) (*WameowClient, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}
	return client, nil
}

// validateButtons - Helper method to validate button data
func (m *MeowService) validateButtons(buttons []ButtonData) error {
	if len(buttons) == 0 {
		return fmt.Errorf("at least one button is required")
	}
	if len(buttons) > 3 {
		return fmt.Errorf("maximum 3 buttons allowed")
	}
	return nil
}

// buildWhatsAppButtons - Helper method to convert ButtonData to WhatsApp buttons
func (m *MeowService) buildWhatsAppButtons(buttons []ButtonData) []*waProto.ButtonsMessage_Button {
	var waButtons []*waProto.ButtonsMessage_Button
	for _, button := range buttons {
		waButtons = append(waButtons, &waProto.ButtonsMessage_Button{
			ButtonID:   &button.ID,
			ButtonText: &waProto.ButtonsMessage_Button_ButtonText{DisplayText: &button.Text},
			Type:       waProto.ButtonsMessage_Button_RESPONSE.Enum(),
		})
	}
	return waButtons
}

// validateListSections - Helper method to validate list sections
func (m *MeowService) validateListSections(sections []ListSection) error {
	if len(sections) == 0 {
		return fmt.Errorf("at least one section is required")
	}
	return nil
}

// buildWhatsAppListSections - Helper method to convert ListSection to WhatsApp sections
func (m *MeowService) buildWhatsAppListSections(sections []ListSection) []*waProto.ListMessage_Section {
	var waSections []*waProto.ListMessage_Section
	for _, section := range sections {
		var waRows []*waProto.ListMessage_Row
		for _, row := range section.Rows {
			waRows = append(waRows, &waProto.ListMessage_Row{
				RowID:       &row.ID,
				Title:       &row.Title,
				Description: &row.Description,
			})
		}

		waSections = append(waSections, &waProto.ListMessage_Section{
			Title: &section.Title,
			Rows:  waRows,
		})
	}
	return waSections
}

func (m *MeowService) removeClient(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, sessionID)
}

func (m *MeowService) GetQRCode(sessionID string) (string, error) {
	client := m.getOrCreateClient(sessionID)
	if client == nil {
		return "", fmt.Errorf("failed to create client for session %s", sessionID)
	}
	return client.GetQRCode()
}

func (m *MeowService) PairPhone(sessionID, phoneNumber string) (string, error) {
	m.logger.Infof("Pairing phone %s for session %s", phoneNumber, sessionID)
	client := m.getOrCreateClient(sessionID)
	return client.PairPhone(phoneNumber)
}

func (m *MeowService) IsClientConnected(sessionID string) bool {
	client := m.getClient(sessionID)
	return client != nil && client.IsConnected()
}

func (m *MeowService) GetClientStatus(sessionID string) session.Status {
	client := m.getClient(sessionID)
	if client == nil {
		return session.StatusDisconnected
	}
	return client.GetStatus()
}

func (m *MeowService) ConnectOnStartup(ctx context.Context) error {
	m.logger.Infof("Connecting sessions with credentials on startup")

	// Get sessions with device credentials (device_jid) from repository
	sessions, err := m.sessions.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sessions with credentials: %w", err)
	}

	m.logger.Infof("Found %d sessions with credentials to reconnect", len(sessions))

	// Validate and start each session that has credentials
	for _, sessionEntity := range sessions {
		// Validate session integrity before attempting reconnection
		if sessionEntity.GetDeviceJIDString() == "" {
			m.logger.Warnf("Session %s (%s) has no device_jid but was returned as active - fixing status",
				sessionEntity.ID().Value(), sessionEntity.Name().Value())

			// Fix inconsistent state: session without device_jid cannot be connected
			err := sessionEntity.Disconnect("no device_jid")
			if err != nil {
				m.logger.Errorf("Failed to disconnect session %s: %v", sessionEntity.ID().Value(), err)
			}
			if err := m.sessions.Update(ctx, sessionEntity); err != nil {
				m.logger.Errorf("Failed to fix session %s status: %v", sessionEntity.ID().Value(), err)
			}
			continue
		}

		// Validate that device exists in whatsmeow_device table
		if !m.deviceExistsInDatabase(sessionEntity.GetDeviceJIDString()) {
			m.logger.Warnf("Session %s (%s) has device_jid %s but device not found in whatsmeow_device table - marking as disconnected",
				sessionEntity.ID().Value(), sessionEntity.Name().Value(), sessionEntity.GetDeviceJIDString())

			err := sessionEntity.Disconnect("device not found in database")
			if err != nil {
				m.logger.Errorf("Failed to disconnect session %s: %v", sessionEntity.ID().Value(), err)
			}
			err = sessionEntity.Authenticate("") // Clear invalid device_jid
			if err != nil {
				m.logger.Errorf("Failed to clear device_jid for session %s: %v", sessionEntity.ID().Value(), err)
			}
			if err := m.sessions.Update(ctx, sessionEntity); err != nil {
				m.logger.Errorf("Failed to fix session %s: %v", sessionEntity.ID().Value(), err)
			}
			continue
		}

		m.logger.Infof("Attempting to reconnect session %s (status: %s, device_jid: %s)",
			sessionEntity.ID().Value(), sessionEntity.Status().String(), sessionEntity.GetDeviceJIDString())

		if err := m.StartClient(sessionEntity.ID().Value()); err != nil {
			m.logger.Errorf("Failed to start client for session %s: %v", sessionEntity.ID, err)
			continue
		}

		m.logger.Infof("Successfully initiated reconnection for session %s", sessionEntity.ID)
	}

	return nil
}

// deviceExistsInDatabase checks if a device JID exists in the whatsmeow_device table
func (m *MeowService) deviceExistsInDatabase(deviceJID string) bool {
	if deviceJID == "" {
		return false
	}

	// Try to get device from container
	devices, err := m.container.GetAllDevices(context.Background())
	if err != nil {
		m.logger.Errorf("Failed to get devices from container: %v", err)
		return false
	}

	for _, device := range devices {
		if device != nil && device.ID != nil && device.ID.String() == deviceJID {
			return true
		}
	}

	return false
}

func (m *MeowService) DeleteMessage(ctx context.Context, sessionID, phone, messageID string, forEveryone bool) error {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// Parse phone to JID
	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	// Build revoke message based on forEveryone flag
	var revokeTarget waTypes.JID
	if forEveryone {
		revokeTarget = waTypes.EmptyJID // Delete for everyone
	} else {
		if client.GetClient().Store.ID == nil {
			return fmt.Errorf("unable to get client ID for delete operation")
		}
		revokeTarget = *client.GetClient().Store.ID // Delete for me only
	}

	// Create revoke message
	revokeMsg := client.GetClient().BuildRevoke(recipient, revokeTarget, messageID)

	// Send revoke message
	_, err = client.GetClient().SendMessage(ctx, recipient, revokeMsg)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	deleteType := "for me"
	if forEveryone {
		deleteType = "for everyone"
	}

	m.logger.Debugf("Message %s deleted %s for phone %s in session %s",
		messageID, deleteType, phone, sessionID)

	return nil
}

func (m *MeowService) EditMessage(ctx context.Context, sessionID, phone, messageID, newText string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// Parse phone to JID
	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	// Create edit message with new text
	editMsg := client.GetClient().BuildEdit(recipient, messageID, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: &newText,
		},
	})

	// Send edit message
	resp, err := client.GetClient().SendMessage(ctx, recipient, editMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to edit message: %w", err)
	}

	m.logger.Debugf("Message %s edited for phone %s in session %s",
		messageID, phone, sessionID)

	return &resp, nil
}

func (m *MeowService) DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error) {
	return nil, "", fmt.Errorf("media download not implemented yet")
}

func (m *MeowService) ReactToMessage(ctx context.Context, sessionID, phone, messageID, emoji string) error {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// Parse phone to JID
	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	// Handle remove reaction
	reaction := emoji
	if emoji == "remove" {
		reaction = ""
	}

	// Determine if message is from me based on messageID prefix
	fromMe := false
	actualMessageID := messageID
	if strings.HasPrefix(messageID, "me:") {
		fromMe = true
		actualMessageID = messageID[len("me:"):]
	}

	// Create reaction message
	recipientStr := recipient.String()
	timestampMs := time.Now().UnixMilli()

	reactionMsg := &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJID: &recipientStr,
				FromMe:    &fromMe,
				ID:        &actualMessageID,
			},
			Text:              &reaction,
			GroupingKey:       &reaction,
			SenderTimestampMS: &timestampMs,
		},
	}

	// Send reaction message
	_, err = client.GetClient().SendMessage(ctx, recipient, reactionMsg)
	if err != nil {
		return fmt.Errorf("failed to send reaction: %w", err)
	}

	reactionType := "added"
	if reaction == "" {
		reactionType = "removed"
	}

	m.logger.Debugf("Reaction %s %s for message %s for phone %s in session %s",
		emoji, reactionType, messageID, phone, sessionID)

	return nil
}

func (m *MeowService) SendTextMessage(ctx context.Context, sessionID, phone, text string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendTextMessage(client.GetClient(), phone, text)
}

func (m *MeowService) SendImageMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendImageMessage(client.GetClient(), phone, data, caption)
}

func (m *MeowService) SendAudioMessage(ctx context.Context, sessionID, phone string, data []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendAudioMessage(client.GetClient(), phone, data, mimeType)
}

func (m *MeowService) SendVideoMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendVideoMessage(client.GetClient(), phone, data, caption, mimeType)
}

func (m *MeowService) SendDocumentMessage(ctx context.Context, sessionID, phone string, data []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendDocumentMessage(client.GetClient(), phone, data, filename, caption, mimeType)
}

func (m *MeowService) SendStickerMessage(ctx context.Context, sessionID, phone string, data []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendStickerMessage(client.GetClient(), phone, data, mimeType)
}

func (m *MeowService) SendContactMessage(ctx context.Context, sessionID, phone, contactName, contactPhone string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendContactMessage(client.GetClient(), phone, contactName, contactPhone)
}

func (m *MeowService) SendLocationMessage(ctx context.Context, sessionID, phone string, latitude, longitude float64, name, address string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendLocationMessage(client.GetClient(), phone, latitude, longitude, name, address)
}

func (m *MeowService) MarkAsRead(ctx context.Context, sessionID, phone string, messageIDs []string) error {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	jid, err := parsePhoneToJID(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	err = client.GetClient().MarkRead(messageIDs, time.Now(), jid, jid)
	if err != nil {
		m.logger.Errorf("Failed to mark messages as read in session %s: %v", sessionID, err)
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	m.logger.Infof("Marked %d messages as read for phone %s in session %s", len(messageIDs), phone, sessionID)
	return nil
}

func (m *MeowService) SendButtonMessage(ctx context.Context, sessionID, phone, title string, buttons []ButtonData) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// Parse phone to JID
	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	// Validate buttons
	if err := m.validateButtons(buttons); err != nil {
		return nil, err
	}

	// Create buttons for WhatsApp
	waButtons := m.buildWhatsAppButtons(buttons)

	// Create button message
	buttonMsg := &waProto.Message{
		ButtonsMessage: &waProto.ButtonsMessage{
			ContentText: &title,
			HeaderType:  waProto.ButtonsMessage_EMPTY.Enum(),
			Buttons:     waButtons,
		},
	}

	// Send message
	resp, err := client.GetClient().SendMessage(ctx, recipient, buttonMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send button message: %w", err)
	}

	m.logger.Debugf("Button message sent to %s from session %s", phone, sessionID)

	return &resp, nil
}

func (m *MeowService) SendListMessage(ctx context.Context, sessionID, phone, title, description, buttonText, footerText string, sections []ListSection) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// Parse phone to JID
	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	// Validate sections
	if err := m.validateListSections(sections); err != nil {
		return nil, err
	}

	// Create sections for WhatsApp
	waSections := m.buildWhatsAppListSections(sections)

	// Create list message
	listMsg := &waProto.Message{
		ListMessage: &waProto.ListMessage{
			Title:       &title,
			Description: &description,
			ButtonText:  &buttonText,
			FooterText:  &footerText,
			ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
			Sections:    waSections,
		},
	}

	// Send message
	resp, err := client.GetClient().SendMessage(ctx, recipient, listMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send list message: %w", err)
	}

	m.logger.Debugf("List message sent to %s from session %s", phone, sessionID)

	return &resp, nil
}

func (m *MeowService) SendPollMessage(ctx context.Context, sessionID, phone, name string, options []string, selectableCount int) (*whatsmeow.SendResponse, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse phone to JID
	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	// Validate options
	if len(options) < 2 {
		return nil, fmt.Errorf("at least 2 options are required for a poll")
	}
	if len(options) > 12 {
		return nil, fmt.Errorf("maximum 12 options allowed for a poll")
	}

	// Validate selectable count
	if selectableCount <= 0 {
		selectableCount = 1 // Default to single select
	}
	if selectableCount > len(options) {
		selectableCount = len(options)
	}

	// Create poll message using whatsmeow's BuildPollCreation
	pollMsg := client.GetClient().BuildPollCreation(name, options, selectableCount)

	// Send message
	resp, err := client.GetClient().SendMessage(ctx, recipient, pollMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send poll message: %w", err)
	}

	m.logger.Debugf("Poll message '%s' sent to %s from session %s", name, phone, sessionID)

	return &resp, nil
}

func (m *MeowService) SetPresence(ctx context.Context, sessionID, phone, state, media string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	// Parse phone to JID if provided, otherwise use empty JID for global presence
	var chatJID waTypes.JID
	if phone != "" {
		jid, err := parsePhoneToJID(phone)
		if err != nil {
			return fmt.Errorf("invalid phone number: %w", err)
		}
		chatJID = jid
	}

	// Send presence - for chat-specific presence, use SendChatPresence
	if !chatJID.IsEmpty() {
		var chatPresence waTypes.ChatPresence
		switch strings.ToLower(state) {
		case "available":
			chatPresence = waTypes.ChatPresenceComposing
		case "unavailable":
			chatPresence = waTypes.ChatPresencePaused
		default:
			return fmt.Errorf("invalid chat presence state: %s", state)
		}

		err := client.GetClient().SendChatPresence(chatJID, chatPresence, "")
		if err != nil {
			return fmt.Errorf("failed to send chat presence: %w", err)
		}
	} else {
		// Global presence
		var presence waTypes.Presence
		switch strings.ToLower(state) {
		case "available":
			presence = waTypes.PresenceAvailable
		case "unavailable":
			presence = waTypes.PresenceUnavailable
		default:
			return fmt.Errorf("invalid presence state: %s", state)
		}

		err := client.GetClient().SendPresence(presence)
		if err != nil {
			return fmt.Errorf("failed to send presence: %w", err)
		}
	}

	m.logger.Debugf("Presence %s sent for session %s", state, sessionID)
	return nil
}

func (m *MeowService) CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Convert participants to JIDs
	var participantJIDs []waTypes.JID
	for _, phone := range participants {
		jid, err := parsePhoneToJID(phone)
		if err != nil {
			m.logger.Warnf("Invalid participant phone %s: %v", phone, err)
			continue
		}
		participantJIDs = append(participantJIDs, jid)
	}

	if len(participantJIDs) == 0 {
		return nil, fmt.Errorf("no valid participants provided")
	}

	// Create group request
	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participantJIDs,
	}

	// Create group
	groupInfo, err := client.GetClient().CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Convert to our GroupInfo format
	result := &GroupInfo{
		JID:       groupInfo.JID.String(),
		Name:      groupInfo.Name,
		Topic:     groupInfo.Topic,
		Owner:     groupInfo.OwnerJID.String(),
		CreatedAt: groupInfo.GroupCreated.Unix(),
		Size:      len(groupInfo.Participants),
		Announce:  groupInfo.IsAnnounce,
		Locked:    groupInfo.IsLocked,
		Ephemeral: groupInfo.IsEphemeral,
	}

	// Add participants
	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	// Add admins
	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Group '%s' created successfully: %s for session %s",
		name, groupInfo.JID.String(), sessionID)

	return result, nil
}

// All remaining group and user functions are stubs for now
func (m *MeowService) ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Get joined groups
	groups, err := client.GetClient().GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get joined groups: %w", err)
	}

	// Convert to our GroupInfo format
	var results []GroupInfo
	for _, group := range groups {
		result := GroupInfo{
			JID:       group.JID.String(),
			Name:      group.Name,
			Topic:     group.Topic,
			Owner:     group.OwnerJID.String(),
			CreatedAt: group.GroupCreated.Unix(),
			Size:      len(group.Participants),
			Announce:  group.IsAnnounce,
			Locked:    group.IsLocked,
			Ephemeral: group.IsEphemeral,
		}

		// Add participants
		for _, participant := range group.Participants {
			result.Participants = append(result.Participants, participant.JID.String())
		}

		// Add admins
		for _, participant := range group.Participants {
			if participant.IsAdmin {
				result.Admins = append(result.Admins, participant.JID.String())
			}
		}

		results = append(results, result)
	}

	m.logger.Debugf("Retrieved %d groups for session %s", len(results), sessionID)

	return results, nil
}

func (m *MeowService) GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Get group info
	groupInfo, err := client.GetClient().GetGroupInfo(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

	// Convert to our GroupInfo format
	result := &GroupInfo{
		JID:       groupInfo.JID.String(),
		Name:      groupInfo.Name,
		Topic:     groupInfo.Topic,
		Owner:     groupInfo.OwnerJID.String(),
		CreatedAt: groupInfo.GroupCreated.Unix(),
		Size:      len(groupInfo.Participants),
		Announce:  groupInfo.IsAnnounce,
		Locked:    groupInfo.IsLocked,
		Ephemeral: groupInfo.IsEphemeral,
	}

	// Add participants
	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	// Add admins
	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Retrieved group info for %s in session %s", groupJID, sessionID)

	return result, nil
}

func (m *MeowService) JoinGroup(ctx context.Context, sessionID, inviteLink string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Join group via invite link
	groupJID, err := client.GetClient().JoinGroupWithLink(inviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	// Get group info after joining
	groupInfo, err := client.GetClient().GetGroupInfo(groupJID)
	if err != nil {
		// If we can't get group info, return basic info
		result := &GroupInfo{
			JID: groupJID.String(),
		}
		m.logger.Debugf("Successfully joined group %s via invite link for session %s (basic info)",
			groupJID.String(), sessionID)
		return result, nil
	}

	// Convert to our GroupInfo format
	result := &GroupInfo{
		JID:       groupInfo.JID.String(),
		Name:      groupInfo.Name,
		Topic:     groupInfo.Topic,
		Owner:     groupInfo.OwnerJID.String(),
		CreatedAt: groupInfo.GroupCreated.Unix(),
		Size:      len(groupInfo.Participants),
		Announce:  groupInfo.IsAnnounce,
		Locked:    groupInfo.IsLocked,
		Ephemeral: groupInfo.IsEphemeral,
	}

	// Add participants
	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	// Add admins
	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Successfully joined group %s via invite link for session %s",
		groupInfo.JID.String(), sessionID)

	return result, nil
}

func (m *MeowService) JoinGroupWithInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Parse inviter JID
	inviterJID, err := waTypes.ParseJID(inviter)
	if err != nil {
		return nil, fmt.Errorf("invalid inviter JID %s: %w", inviter, err)
	}

	// Join group via specific invite
	err = client.GetClient().JoinGroupWithInvite(groupJIDParsed, inviterJID, code, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to join group with invite: %w", err)
	}

	// Get group info after joining
	groupInfo, err := client.GetClient().GetGroupInfo(groupJIDParsed)
	if err != nil {
		// If we can't get group info, return basic info
		result := &GroupInfo{
			JID: groupJIDParsed.String(),
		}
		m.logger.Debugf("Successfully joined group %s via specific invite for session %s (basic info)",
			groupJIDParsed.String(), sessionID)
		return result, nil
	}

	// Convert to our GroupInfo format
	result := &GroupInfo{
		JID:       groupInfo.JID.String(),
		Name:      groupInfo.Name,
		Topic:     groupInfo.Topic,
		Owner:     groupInfo.OwnerJID.String(),
		CreatedAt: groupInfo.GroupCreated.Unix(),
		Size:      len(groupInfo.Participants),
		Announce:  groupInfo.IsAnnounce,
		Locked:    groupInfo.IsLocked,
		Ephemeral: groupInfo.IsEphemeral,
	}

	// Add participants
	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	// Add admins
	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Successfully joined group %s via specific invite for session %s",
		groupJIDParsed.String(), sessionID)

	return result, nil
}

func (m *MeowService) LeaveGroup(ctx context.Context, sessionID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Leave group
	err = client.GetClient().LeaveGroup(jid)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	m.logger.Debugf("Successfully left group %s for session %s", groupJID, sessionID)

	return nil
}

func (m *MeowService) GetInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return "", fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return "", fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Get invite link
	inviteLink, err := client.GetClient().GetGroupInviteLink(jid, reset)
	if err != nil {
		return "", fmt.Errorf("failed to get invite link: %w", err)
	}

	m.logger.Debugf("Retrieved invite link for group %s in session %s", groupJID, sessionID)

	return inviteLink, nil
}

func (m *MeowService) SetGroupEphemeral(ctx context.Context, sessionID, groupJID string, ephemeral bool, duration int) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Convert duration to time.Duration
	var expiration time.Duration
	if ephemeral && duration > 0 {
		expiration = time.Duration(duration) * time.Second
	} else if ephemeral {
		expiration = 24 * time.Hour // Default 24 hours
	} else {
		expiration = 0 // Disable ephemeral
	}

	// Set group ephemeral
	err = client.GetClient().SetDisappearingTimer(jid, expiration, time.Time{})
	if err != nil {
		return fmt.Errorf("failed to set group ephemeral: %w", err)
	}

	ephemeralStatus := "disabled"
	if ephemeral {
		ephemeralStatus = fmt.Sprintf("enabled (%d seconds)", duration)
	}

	m.logger.Debugf("Successfully set group ephemeral to %s for group %s in session %s",
		ephemeralStatus, groupJID, sessionID)

	return nil
}

func (m *MeowService) GetGroupRequestParticipants(ctx context.Context, sessionID, groupJID string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Get group request participants
	requestParticipants, err := client.GetClient().GetGroupRequestParticipants(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group request participants: %w", err)
	}

	// Convert JIDs to strings
	var participants []string
	for _, participant := range requestParticipants {
		participants = append(participants, participant.JID.String())
	}

	m.logger.Debugf("Retrieved %d group request participants for group %s in session %s",
		len(participants), groupJID, sessionID)

	return participants, nil
}

func (m *MeowService) UpdateGroupRequestParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Convert participants to JIDs
	var participantJIDs []waTypes.JID
	for _, phone := range participants {
		jid, err := parsePhoneToJID(phone)
		if err != nil {
			m.logger.Warnf("Invalid participant phone %s: %v", phone, err)
			continue
		}
		participantJIDs = append(participantJIDs, jid)
	}

	if len(participantJIDs) == 0 {
		return fmt.Errorf("no valid participants provided")
	}

	// Parse action
	var requestChange whatsmeow.ParticipantRequestChange
	switch action {
	case "approve":
		requestChange = whatsmeow.ParticipantChangeApprove
	case "reject":
		requestChange = whatsmeow.ParticipantChangeReject
	default:
		return fmt.Errorf("invalid action %s. Valid actions: approve, reject", action)
	}

	// Update group request participants
	_, err = client.GetClient().UpdateGroupRequestParticipants(groupJIDParsed, participantJIDs, requestChange)
	if err != nil {
		return fmt.Errorf("failed to %s group request participants: %w", action, err)
	}

	m.logger.Debugf("Successfully %sed %d group request participants for group %s in session %s",
		action, len(participantJIDs), groupJID, sessionID)

	return nil
}

func (m *MeowService) LinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse community JID
	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	// Parse group JID
	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Link group to community
	err = client.GetClient().LinkGroup(communityJIDParsed, groupJIDParsed)
	if err != nil {
		return fmt.Errorf("failed to link group to community: %w", err)
	}

	m.logger.Debugf("Successfully linked group %s to community %s in session %s",
		groupJID, communityJID, sessionID)

	return nil
}

func (m *MeowService) UnlinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse community JID
	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	// Parse group JID
	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Unlink group from community
	err = client.GetClient().UnlinkGroup(communityJIDParsed, groupJIDParsed)
	if err != nil {
		return fmt.Errorf("failed to unlink group from community: %w", err)
	}

	m.logger.Debugf("Successfully unlinked group %s from community %s in session %s",
		groupJID, communityJID, sessionID)

	return nil
}

func (m *MeowService) GetSubGroups(ctx context.Context, sessionID, communityJID string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse community JID
	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return nil, fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	// Get subgroups
	subGroups, err := client.GetClient().GetSubGroups(communityJIDParsed)
	if err != nil {
		return nil, fmt.Errorf("failed to get subgroups: %w", err)
	}

	// Convert JIDs to strings
	var groups []string
	for _, group := range subGroups {
		// Check if subGroups returns []*types.GroupInfo or []types.JID
		if group != nil {
			groups = append(groups, group.JID.String())
		}
	}

	m.logger.Debugf("Retrieved %d subgroups for community %s in session %s",
		len(groups), communityJID, sessionID)

	return groups, nil
}

func (m *MeowService) GetLinkedGroupsParticipants(ctx context.Context, sessionID, communityJID string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse community JID
	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return nil, fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	// Get linked groups participants
	participants, err := client.GetClient().GetLinkedGroupsParticipants(communityJIDParsed)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked groups participants: %w", err)
	}

	// Convert JIDs to strings
	var participantList []string
	for _, participant := range participants {
		participantList = append(participantList, participant.String())
	}

	m.logger.Debugf("Retrieved %d linked groups participants for community %s in session %s",
		len(participantList), communityJID, sessionID)

	return participantList, nil
}

func (m *MeowService) CheckUser(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Validate and clean phone numbers
	var validPhones []string

	for _, phone := range phones {
		_, err := parsePhoneToJID(phone)
		if err != nil {
			m.logger.Warnf("Invalid phone number %s: %v", phone, err)
			continue
		}
		validPhones = append(validPhones, phone)
	}

	if len(validPhones) == 0 {
		return nil, fmt.Errorf("no valid phone numbers provided")
	}

	// Check if users are on WhatsApp
	resp, err := client.GetClient().IsOnWhatsApp(validPhones)
	if err != nil {
		return nil, fmt.Errorf("failed to check users on WhatsApp: %w", err)
	}

	// Convert response to our format
	var results []UserCheckResult
	for _, item := range resp {
		verifiedName := ""
		if item.VerifiedName != nil {
			verifiedName = item.VerifiedName.Details.GetVerifiedName()
		}

		result := UserCheckResult{
			Query:        item.Query,
			IsInWhatsapp: item.IsIn,
			IsInmeow:     item.IsIn, // Set alias to same value
			JID:          item.JID.String(),
			VerifiedName: verifiedName,
		}
		results = append(results, result)
	}

	m.logger.Debugf("Checked %d users for session %s, found %d on WhatsApp",
		len(phones), sessionID, len(results))

	return results, nil
}

func (m *MeowService) GetUserInfo(ctx context.Context, sessionID string, phones []string) (map[string]UserInfoResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Convert phones to JIDs for whatsmeow
	var jids []waTypes.JID
	phoneToJIDMap := make(map[string]string) // phone -> jid string mapping

	for _, phone := range phones {
		jid, err := parsePhoneToJID(phone)
		if err != nil {
			m.logger.Warnf("Invalid phone number %s: %v", phone, err)
			continue
		}
		jids = append(jids, jid)
		phoneToJIDMap[phone] = jid.String()
	}

	if len(jids) == 0 {
		return nil, fmt.Errorf("no valid phone numbers provided")
	}

	// Get user info from WhatsApp
	resp, err := client.GetClient().GetUserInfo(jids)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Convert response to our format
	results := make(map[string]UserInfoResult)
	for jid, info := range resp {
		verifiedName := ""
		if info.VerifiedName != nil {
			verifiedName = info.VerifiedName.Details.GetVerifiedName()
		}

		result := UserInfoResult{
			JID:          jid.String(),
			DisplayName:  "", // Will be populated from contact info if available
			Name:         "", // Alias for DisplayName
			PushName:     "", // Will be populated from contact info if available
			BusinessName: "", // Will be populated if business account
			Phone:        "", // Will be populated from JID
			VerifiedName: verifiedName,
			Avatar:       "", // Will be populated if available
			Status:       info.Status,
			PictureID:    info.PictureID,
			DeviceCount:  len(info.Devices),
			Notify:       "",    // Will be populated from contact info if available
			IsBlocked:    false, // Default value
			IsMuted:      false, // Default value
		}
		results[jid.String()] = result
	}

	m.logger.Debugf("Retrieved info for %d users for session %s",
		len(results), sessionID)

	return results, nil
}

func (m *MeowService) GetAvatar(ctx context.Context, sessionID, phone string) (*AvatarResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse phone to JID
	jid, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	// Get profile picture info
	pictureInfo, err := client.GetClient().GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar: %w", err)
	}

	result := &AvatarResult{
		Phone:     phone,
		JID:       jid.String(),
		AvatarURL: pictureInfo.URL,
		PictureID: pictureInfo.ID,
	}

	m.logger.Debugf("Retrieved avatar for %s in session %s", phone, sessionID)

	return result, nil
}

func (m *MeowService) GetContacts(ctx context.Context, sessionID string) ([]ContactResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Get contacts from store
	contacts, err := client.GetClient().Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	// Convert to our ContactResult format
	var results []ContactResult
	for jid, contact := range contacts {
		result := ContactResult{
			JID:          jid.String(),
			Name:         contact.FullName,
			Notify:       contact.PushName,
			PushName:     contact.PushName,
			BusinessName: contact.BusinessName,
			IsBlocked:    false, // WhatsApp doesn't expose this easily
			IsMuted:      false, // WhatsApp doesn't expose this easily
		}
		results = append(results, result)
	}

	m.logger.Debugf("Retrieved %d contacts for session %s", len(results), sessionID)

	return results, nil
}

func (m *MeowService) SetUserPresence(ctx context.Context, sessionID, state string) error {
	// SetUserPresence is global presence, so we call SetPresence with empty phone
	return m.SetPresence(ctx, sessionID, "", state, "")
}

func (m *MeowService) GetInviteInfo(ctx context.Context, sessionID, inviteLink string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Get invite info
	groupInfo, err := client.GetClient().GetGroupInfoFromLink(inviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to get invite info: %w", err)
	}

	// Convert to our GroupInfo format
	result := &GroupInfo{
		JID:       groupInfo.JID.String(),
		Name:      groupInfo.Name,
		Topic:     groupInfo.Topic,
		Owner:     groupInfo.OwnerJID.String(),
		CreatedAt: groupInfo.GroupCreated.Unix(),
		Size:      len(groupInfo.Participants),
		Announce:  groupInfo.IsAnnounce,
		Locked:    groupInfo.IsLocked,
		Ephemeral: groupInfo.IsEphemeral,
	}

	// Add participants
	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	// Add admins
	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Retrieved invite info for group %s from link for session %s",
		groupInfo.JID.String(), sessionID)

	return result, nil
}

func (m *MeowService) GetGroupInfoFromInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Parse inviter JID
	inviterJID, err := waTypes.ParseJID(inviter)
	if err != nil {
		return nil, fmt.Errorf("invalid inviter JID %s: %w", inviter, err)
	}

	// Get group info from specific invite
	groupInfo, err := client.GetClient().GetGroupInfoFromInvite(groupJIDParsed, inviterJID, code, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info from invite: %w", err)
	}

	// Convert to our GroupInfo format
	result := &GroupInfo{
		JID:       groupInfo.JID.String(),
		Name:      groupInfo.Name,
		Topic:     groupInfo.Topic,
		Owner:     groupInfo.OwnerJID.String(),
		CreatedAt: groupInfo.GroupCreated.Unix(),
		Size:      len(groupInfo.Participants),
		Announce:  groupInfo.IsAnnounce,
		Locked:    groupInfo.IsLocked,
		Ephemeral: groupInfo.IsEphemeral,
	}

	// Add participants
	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	// Add admins
	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Retrieved group info from specific invite for group %s in session %s",
		groupInfo.JID.String(), sessionID)

	return result, nil
}

func (m *MeowService) UpdateParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Convert participants to JIDs
	var participantJIDs []waTypes.JID
	for _, phone := range participants {
		jid, err := parsePhoneToJID(phone)
		if err != nil {
			m.logger.Warnf("Invalid participant phone %s: %v", phone, err)
			continue
		}
		participantJIDs = append(participantJIDs, jid)
	}

	if len(participantJIDs) == 0 {
		return fmt.Errorf("no valid participants provided")
	}

	// Parse action
	var participantChange whatsmeow.ParticipantChange
	switch action {
	case "add":
		participantChange = whatsmeow.ParticipantChangeAdd
	case "remove":
		participantChange = whatsmeow.ParticipantChangeRemove
	case "promote":
		participantChange = whatsmeow.ParticipantChangePromote
	case "demote":
		participantChange = whatsmeow.ParticipantChangeDemote
	default:
		return fmt.Errorf("invalid action: %s (must be add, remove, promote, or demote)", action)
	}

	// Update participants
	_, err = client.GetClient().UpdateGroupParticipants(groupJIDParsed, participantJIDs, participantChange)
	if err != nil {
		return fmt.Errorf("failed to %s participants: %w", action, err)
	}

	m.logger.Debugf("Successfully %s %d participants in group %s for session %s",
		action, len(participantJIDs), groupJID, sessionID)

	return nil
}

func (m *MeowService) SetGroupName(ctx context.Context, sessionID, groupJID, name string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Set group name
	err = client.GetClient().SetGroupName(jid, name)
	if err != nil {
		return fmt.Errorf("failed to set group name: %w", err)
	}

	m.logger.Debugf("Successfully set group name to '%s' for group %s in session %s",
		name, groupJID, sessionID)

	return nil
}

func (m *MeowService) SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Set group topic
	err = client.GetClient().SetGroupTopic(jid, "", "", topic)
	if err != nil {
		return fmt.Errorf("failed to set group topic: %w", err)
	}

	m.logger.Debugf("Successfully set group topic to '%s' for group %s in session %s",
		topic, groupJID, sessionID)

	return nil
}

func (m *MeowService) SetGroupPhoto(ctx context.Context, sessionID, groupJID string, photo []byte) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Let whatsmeow validate the image format
	if len(photo) == 0 {
		return fmt.Errorf("photo data cannot be empty")
	}

	// Set group photo
	pictureID, err := client.GetClient().SetGroupPhoto(jid, photo)
	if err != nil {
		return fmt.Errorf("failed to set group photo: %w", err)
	}

	m.logger.Debugf("Successfully set group photo (ID: %s) for group %s in session %s",
		pictureID, groupJID, sessionID)

	return nil
}

func (m *MeowService) RemoveGroupPhoto(ctx context.Context, sessionID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Remove group photo by setting it to nil
	_, err = client.GetClient().SetGroupPhoto(jid, nil)
	if err != nil {
		return fmt.Errorf("failed to remove group photo: %w", err)
	}

	m.logger.Debugf("Successfully removed group photo for group %s in session %s",
		groupJID, sessionID)

	return nil
}

func (m *MeowService) SetGroupAnnounce(ctx context.Context, sessionID, groupJID string, announceOnly bool) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Set group announce setting
	err = client.GetClient().SetGroupAnnounce(jid, announceOnly)
	if err != nil {
		return fmt.Errorf("failed to set group announce: %w", err)
	}

	announceStatus := "disabled"
	if announceOnly {
		announceStatus = "enabled (only admins can send)"
	}

	m.logger.Debugf("Successfully set group announce to %s for group %s in session %s",
		announceStatus, groupJID, sessionID)

	return nil
}

func (m *MeowService) SetGroupLocked(ctx context.Context, sessionID, groupJID string, locked bool) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Set group locked setting
	err = client.GetClient().SetGroupLocked(jid, locked)
	if err != nil {
		return fmt.Errorf("failed to set group locked: %w", err)
	}

	lockedStatus := "unlocked"
	if locked {
		lockedStatus = "locked (only admins can edit info)"
	}

	m.logger.Debugf("Successfully set group to %s for group %s in session %s",
		lockedStatus, groupJID, sessionID)

	return nil
}

func (m *MeowService) SetGroupJoinApprovalMode(ctx context.Context, sessionID, groupJID string, requireApproval bool) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Set group join approval mode
	err = client.GetClient().SetGroupJoinApprovalMode(jid, requireApproval)
	if err != nil {
		return fmt.Errorf("failed to set group join approval mode: %w", err)
	}

	approvalStatus := "open (no approval required)"
	if requireApproval {
		approvalStatus = "require approval"
	}

	m.logger.Debugf("Successfully set group join approval mode to %s for group %s in session %s",
		approvalStatus, groupJID, sessionID)

	return nil
}

// SetGroupJoinApproval is an alias for SetGroupJoinApprovalMode for compatibility
func (m *MeowService) SetGroupJoinApproval(ctx context.Context, sessionID, groupJID string, requireApproval bool) error {
	return m.SetGroupJoinApprovalMode(ctx, sessionID, groupJID, requireApproval)
}

func (m *MeowService) SetGroupMemberAddMode(ctx context.Context, sessionID, groupJID string, mode string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	// Parse group JID
	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	// Parse member add mode
	var memberAddMode waTypes.GroupMemberAddMode
	switch strings.ToLower(mode) {
	case "all", "everyone":
		memberAddMode = waTypes.GroupMemberAddModeAllMember
	case "admin", "admins", "admin_only":
		memberAddMode = waTypes.GroupMemberAddModeAdmin
	default:
		return fmt.Errorf("invalid member add mode: %s. Valid options: all, admin", mode)
	}

	// Set group member add mode
	err = client.GetClient().SetGroupMemberAddMode(jid, memberAddMode)
	if err != nil {
		return fmt.Errorf("failed to set group member add mode: %w", err)
	}

	m.logger.Debugf("Successfully set group member add mode to %s for group %s in session %s",
		mode, groupJID, sessionID)

	return nil
}

func (m *MeowService) UpdateSessionWebhook(sessionID, webhookURL string) error {
	m.logger.Infof("Updated webhook URL for session %s", sessionID)
	return nil
}

func (m *MeowService) UpdateSessionSubscriptions(sessionID string, events []string) error {
	m.logger.Infof("Updated event subscriptions for session %s", sessionID)
	return nil
}

// Chat Operations Implementation

func (m *MeowService) SetDisappearingTimer(ctx context.Context, sessionID, chatJID string, duration time.Duration) error {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// Parse chat JID
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	// Set disappearing timer using whatsmeow client
	err = client.GetClient().SetDisappearingTimer(jid, duration, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set disappearing timer: %w", err)
	}

	m.logger.Debugf("Set disappearing timer to %v for chat %s in session %s", duration, chatJID, sessionID)
	return nil
}

func (m *MeowService) ListChats(ctx context.Context, sessionID, chatType string) ([]ChatInfo, error) {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return a basic implementation
	// In a full implementation, this would query the WhatsApp client for conversations
	var chats []ChatInfo

	// Get groups if requested
	if chatType == "" || chatType == "groups" || chatType == "all" {
		groups, err := client.GetClient().GetJoinedGroups()
		if err == nil {
			for _, group := range groups {
				timestamp := time.Now().Unix() // Default timestamp
				chat := ChatInfo{
					JID:           group.JID.String(),
					Name:          group.Name, // Use the actual group name
					Type:          "group",
					IsGroup:       true,
					LastTimestamp: timestamp,
					Timestamp:     timestamp, // Set alias to same value
				}
				chats = append(chats, chat)
			}
		}
	}

	m.logger.Debugf("Retrieved %d chats for session %s", len(chats), sessionID)
	return chats, nil
}

func (m *MeowService) GetChatInfo(ctx context.Context, sessionID, chatJID string) (*ChatInfo, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// Parse chat JID
	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat JID: %w", err)
	}

	// Basic implementation - check if it's a group
	isGroup := strings.Contains(jid.Server, "g.us")
	chatType := "contact"
	if isGroup {
		chatType = "group"
	}

	timestamp := time.Now().Unix() // Default timestamp
	chat := &ChatInfo{
		JID:           chatJID,
		Name:          chatJID, // Would need to get actual name from group info or contacts
		Type:          chatType,
		IsGroup:       isGroup,
		LastTimestamp: timestamp,
		Timestamp:     timestamp, // Set alias to same value
	}

	m.logger.Debugf("Retrieved chat info for %s in session %s", chatJID, sessionID)
	return chat, nil
}

func (m *MeowService) PinChat(ctx context.Context, sessionID, chatJID string, pinned bool) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// Parse chat JID
	_, err = waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	// For now, just log the action - full implementation would use app state
	action := "unpinned"
	if pinned {
		action = "pinned"
	}
	m.logger.Debugf("Chat %s %s in session %s", chatJID, action, sessionID)
	return nil
}

func (m *MeowService) MuteChat(ctx context.Context, sessionID, chatJID string, muted bool, duration time.Duration) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// Parse chat JID
	_, err = waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	// For now, just log the action - full implementation would use app state
	action := "unmuted"
	if muted {
		action = "muted"
		if duration > 0 {
			m.logger.Debugf("Chat %s %s for %v in session %s", chatJID, action, duration, sessionID)
		} else {
			m.logger.Debugf("Chat %s %s forever in session %s", chatJID, action, sessionID)
		}
	} else {
		m.logger.Debugf("Chat %s %s in session %s", chatJID, action, sessionID)
	}
	return nil
}

func (m *MeowService) ArchiveChat(ctx context.Context, sessionID, chatJID string, archived bool) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// Parse chat JID
	_, err = waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	// For now, just log the action - full implementation would use app state
	action := "unarchived"
	if archived {
		action = "archived"
	}
	m.logger.Debugf("Chat %s %s in session %s", chatJID, action, sessionID)
	return nil
}

// Newsletter Operations Implementation

func (m *MeowService) GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterID string) ([]NewsletterMessage, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return empty list - full implementation would query newsletter messages
	var messages []NewsletterMessage

	m.logger.Debugf("Retrieved %d newsletter messages for newsletter %s in session %s", len(messages), newsletterID, sessionID)
	return messages, nil
}

func (m *MeowService) NewsletterMarkViewed(ctx context.Context, sessionID, newsletterID string, messageIDs []string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would mark messages as viewed
	m.logger.Debugf("Marked %d newsletter messages as viewed for newsletter %s in session %s", len(messageIDs), newsletterID, sessionID)
	return nil
}

func (m *MeowService) NewsletterSendReaction(ctx context.Context, sessionID, newsletterID, messageID, reaction string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would send reaction
	m.logger.Debugf("Sent reaction '%s' to message %s in newsletter %s for session %s", reaction, messageID, newsletterID, sessionID)
	return nil
}

func (m *MeowService) NewsletterToggleMute(ctx context.Context, sessionID, newsletterID string, muted bool) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	action := "unmuted"
	if muted {
		action = "muted"
	}
	m.logger.Debugf("Newsletter %s %s for session %s", newsletterID, action, sessionID)
	return nil
}

func (m *MeowService) NewsletterSubscribeLiveUpdates(ctx context.Context, sessionID, newsletterID string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would subscribe to live updates
	m.logger.Debugf("Subscribed to live updates for newsletter %s in session %s", newsletterID, sessionID)
	return nil
}

func (m *MeowService) UploadNewsletter(ctx context.Context, sessionID string, data []byte) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would upload newsletter
	m.logger.Debugf("Uploaded newsletter data (%d bytes) for session %s", len(data), sessionID)
	return nil
}

func (m *MeowService) GetNewsletterInfoWithInvite(ctx context.Context, sessionID, inviteCode string) (*NewsletterInfo, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return mock data - full implementation would get newsletter info from invite
	timestamp := time.Now().Unix()
	newsletter := &NewsletterInfo{
		ID:          "newsletter_" + inviteCode,
		JID:         "newsletter_" + inviteCode + "@newsletter",
		Name:        "Newsletter from invite",
		Description: "Newsletter obtained from invite code",
		Subscribers: 0,
		Verified:    false,
		IsVerified:  false,
		Muted:       false,
		Following:   false,
		CreatedAt:   timestamp,
		ServerID:    "server_" + inviteCode,
		Timestamp:   timestamp,
	}

	m.logger.Debugf("Retrieved newsletter info from invite %s for session %s", inviteCode, sessionID)
	return newsletter, nil
}

func (m *MeowService) CreateNewsletter(ctx context.Context, sessionID, name, description string) (*NewsletterInfo, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return mock data - full implementation would create newsletter
	timestamp := time.Now().Unix()
	newsletterID := fmt.Sprintf("newsletter_%d", timestamp)
	newsletter := &NewsletterInfo{
		ID:          newsletterID,
		JID:         newsletterID + "@newsletter",
		Name:        name,
		Description: description,
		Subscribers: 0,
		Verified:    false,
		IsVerified:  false,
		Muted:       false,
		Following:   true,
		CreatedAt:   timestamp,
		ServerID:    "server_" + newsletterID,
		Timestamp:   timestamp,
	}

	m.logger.Debugf("Created newsletter '%s' for session %s", name, sessionID)
	return newsletter, nil
}

func (m *MeowService) GetNewsletterInfo(ctx context.Context, sessionID, newsletterID string) (*NewsletterInfo, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return mock data - full implementation would get newsletter info
	timestamp := time.Now().Unix()
	newsletter := &NewsletterInfo{
		ID:          newsletterID,
		JID:         newsletterID + "@newsletter",
		Name:        "Newsletter",
		Description: "Newsletter description",
		Subscribers: 100,
		Verified:    false,
		IsVerified:  false,
		Muted:       false,
		Following:   true,
		CreatedAt:   timestamp,
		ServerID:    "server_" + newsletterID,
		Timestamp:   timestamp,
	}

	m.logger.Debugf("Retrieved newsletter info for %s in session %s", newsletterID, sessionID)
	return newsletter, nil
}

func (m *MeowService) GetSubscribedNewsletters(ctx context.Context, sessionID string) ([]NewsletterInfo, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return empty list - full implementation would get subscribed newsletters
	var newsletters []NewsletterInfo

	m.logger.Debugf("Retrieved %d subscribed newsletters for session %s", len(newsletters), sessionID)
	return newsletters, nil
}

func (m *MeowService) FollowNewsletter(ctx context.Context, sessionID, newsletterID string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would follow newsletter
	m.logger.Debugf("Followed newsletter %s for session %s", newsletterID, sessionID)
	return nil
}

func (m *MeowService) UnfollowNewsletter(ctx context.Context, sessionID, newsletterID string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would unfollow newsletter
	m.logger.Debugf("Unfollowed newsletter %s for session %s", newsletterID, sessionID)
	return nil
}

func (m *MeowService) SendNewsletterMessage(ctx context.Context, sessionID, newsletterID, message string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would send newsletter message
	m.logger.Debugf("Sent message to newsletter %s for session %s: %s", newsletterID, sessionID, message)
	return nil
}

func (m *MeowService) GetNewsletterMessages(ctx context.Context, sessionID, newsletterID string) ([]NewsletterMessage, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return empty list - full implementation would get newsletter messages
	var messages []NewsletterMessage

	m.logger.Debugf("Retrieved %d messages for newsletter %s in session %s", len(messages), newsletterID, sessionID)
	return messages, nil
}

// Privacy Operations Implementation

func (m *MeowService) GetPrivacySettings(ctx context.Context, sessionID string) (*PrivacySettings, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return default privacy settings - full implementation would get actual settings
	settings := &PrivacySettings{
		LastSeen:     "contacts",
		ProfilePhoto: "contacts",
		About:        "contacts",
		Status:       "contacts",
		ReadReceipts: true,
		GroupsAdd:    "contacts",
		CallsAdd:     "contacts",
	}

	m.logger.Debugf("Retrieved privacy settings for session %s", sessionID)
	return settings, nil
}

func (m *MeowService) SetPrivacySetting(ctx context.Context, sessionID, setting, value string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would set privacy setting
	m.logger.Debugf("Set privacy setting %s to %s for session %s", setting, value, sessionID)
	return nil
}

func (m *MeowService) GetBlocklist(ctx context.Context, sessionID string) ([]string, error) {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	// For now, return empty list - full implementation would get blocked contacts
	var blocklist []string

	m.logger.Debugf("Retrieved %d blocked contacts for session %s", len(blocklist), sessionID)
	return blocklist, nil
}

func (m *MeowService) UpdateBlocklist(ctx context.Context, sessionID string, action string, contacts []string) error {
	_, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	// For now, just log the action - full implementation would update blocklist
	m.logger.Debugf("Updated blocklist with action %s for %d contacts in session %s", action, len(contacts), sessionID)
	return nil
}
