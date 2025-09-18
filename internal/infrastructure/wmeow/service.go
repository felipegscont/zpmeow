package wmeow

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"meow/internal/config"
	"meow/internal/domain/session"
	"meow/internal/infrastructure/logging"
	"meow/internal/infrastructure/webhooks"
	"meow/internal/shared/types"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waEvents "go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type ButtonData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type ListItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListSection struct {
	Title string     `json:"title"`
	Rows  []ListItem `json:"rows"`
}

type GroupInfo struct {
	JID          string   `json:"jid"`
	Name         string   `json:"name"`
	Topic        string   `json:"topic"`
	Participants []string `json:"participants"`
	Admins       []string `json:"admins"`
	Owner        string   `json:"owner"`
	CreatedAt    int64    `json:"createdAt"`
	Size         int      `json:"size"`
	Announce     bool     `json:"announce"`
	Locked       bool     `json:"locked"`
	Ephemeral    bool     `json:"ephemeral"`
}

type ChatInfo struct {
	JID         string    `json:"jid"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "contact", "group"
	LastMessage string    `json:"last_message,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	UnreadCount int       `json:"unread_count"`
	Pinned      bool      `json:"pinned"`
	Muted       bool      `json:"muted"`
	Archived    bool      `json:"archived"`
}

type Service interface {
	StartClient(sessionID string) error
	StopClient(sessionID string) error
	LogoutClient(sessionID string) error
	GetQRCode(sessionID string) (string, error)
	PairPhone(sessionID, phoneNumber string) (string, error)
	IsClientConnected(sessionID string) bool
	GetClientStatus(sessionID string) types.Status
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
	SetGroupMemberAddMode(ctx context.Context, sessionID, groupJID string, mode string) error
	GetGroupRequestParticipants(ctx context.Context, sessionID, groupJID string) ([]string, error)
	UpdateGroupRequestParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error
	LinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error
	UnlinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error
	GetSubGroups(ctx context.Context, sessionID, communityJID string) ([]string, error)
	GetLinkedGroupsParticipants(ctx context.Context, sessionID, communityJID string) ([]string, error)

	CheckUser(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error)
	GetUserInfo(ctx context.Context, sessionID string, phones []string) (map[string]UserInfoResult, error)
	GetAvatar(ctx context.Context, sessionID, phone string) (*AvatarResult, error)
	GetContacts(ctx context.Context, sessionID string) ([]ContactResult, error)
	SetUserPresence(ctx context.Context, sessionID, state string) error

	GetPrivacySettings(ctx context.Context, sessionID string) (*PrivacySettingsResult, error)
	TryFetchPrivacySettings(ctx context.Context, sessionID string, ignoreCache bool) (*PrivacySettingsResult, error)
	SetPrivacySetting(ctx context.Context, sessionID string, settingType, settingValue string) (*PrivacySettingsResult, error)
	GetBlocklist(ctx context.Context, sessionID string) ([]string, error)
	UpdateBlocklist(ctx context.Context, sessionID, jidStr, action string) ([]string, error)

	UpdateSessionWebhook(sessionID, webhookURL string) error
	UpdateSessionSubscriptions(sessionID string, events []string) error

	SetDisappearingTimer(ctx context.Context, sessionID, chatJID string, timer time.Duration) error
	ListChats(ctx context.Context, sessionID string, chatType string) ([]ChatInfo, error)
	GetChatInfo(ctx context.Context, sessionID, chatJID string) (*ChatInfo, error)
	PinChat(ctx context.Context, sessionID, chatJID string, pinned bool) error
	MuteChat(ctx context.Context, sessionID, chatJID string, muted bool, duration time.Duration) error
	ArchiveChat(ctx context.Context, sessionID, chatJID string, archived bool) error

	CreateNewsletter(ctx context.Context, sessionID string, params *whatsmeow.CreateNewsletterParams) (*waTypes.NewsletterMetadata, error)
	GetNewsletterInfo(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
	GetNewsletterInfoWithInvite(ctx context.Context, sessionID, inviteKey string) (*NewsletterInfo, error)
	FollowNewsletter(ctx context.Context, sessionID, newsletterJID string) error
	UnfollowNewsletter(ctx context.Context, sessionID, newsletterJID string) error
	GetSubscribedNewsletters(ctx context.Context, sessionID string) ([]NewsletterInfo, error)
	GetNewsletterMessages(ctx context.Context, sessionID, newsletterJID string, params *whatsmeow.GetNewsletterMessagesParams) ([]NewsletterMessage, error)
	GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterJID string, params *whatsmeow.GetNewsletterMessagesParams) ([]NewsletterMessage, error)
	NewsletterMarkViewed(ctx context.Context, sessionID, newsletterJID string, serverIDs []waTypes.MessageServerID) error
	NewsletterSendReaction(ctx context.Context, sessionID, newsletterJID string, serverID waTypes.MessageServerID, reaction string, messageID waTypes.MessageID) error
	NewsletterToggleMute(ctx context.Context, sessionID, newsletterJID string, mute bool) error
	NewsletterSubscribeLiveUpdates(ctx context.Context, sessionID, newsletterJID string) error
	UploadNewsletter(ctx context.Context, sessionID string, data []byte, mediaType whatsmeow.MediaType) (*whatsmeow.UploadResponse, error)
	SendNewsletterMessage(ctx context.Context, sessionID, newsletterJID string, message *waE2E.Message, mediaHandle string) (*whatsmeow.SendResponse, error)
	UploadNewsletterReader(ctx context.Context, sessionID string, data []byte, mediaType whatsmeow.MediaType) (*whatsmeow.UploadResponse, error)
}

type UserCheckResult struct {
	Query        string `json:"query"`
	IsInmeow     bool   `json:"is_in_meow"`
	JID          string `json:"jid"`
	VerifiedName string `json:"verified_name,omitempty"`
}

type UserInfoResult struct {
	JID          string `json:"jid"`
	DisplayName  string `json:"display_name,omitempty"`
	VerifiedName string `json:"verified_name,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Status       string `json:"status,omitempty"`
	PictureID    string `json:"picture_id,omitempty"`
	DeviceCount  int    `json:"device_count,omitempty"`
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

type PrivacySettingsResult struct {
	GroupAdd     string `json:"groupAdd"`     // Who can add to groups: "all", "contacts", "contact_blacklist", "none"
	LastSeen     string `json:"lastSeen"`     // Who can see last seen: "all", "contacts", "contact_blacklist", "none"
	Status       string `json:"status"`       // Who can see status: "all", "contacts", "contact_blacklist", "none"
	Profile      string `json:"profile"`      // Who can see profile photo: "all", "contacts", "contact_blacklist", "none"
	ReadReceipts string `json:"readReceipts"` // Read receipts: "all", "none"
	CallAdd      string `json:"callAdd"`      // Who can call: "all", "known"
	Online       string `json:"online"`       // Who can see online status: "all", "match_last_seen"
}

type NewsletterInfo struct {
	JID         string    `json:"jid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Picture     string    `json:"picture,omitempty"`
	Verified    bool      `json:"verified"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	OwnerJID    string    `json:"owner_jid,omitempty"`
	Subscribers int       `json:"subscribers,omitempty"`
	Muted       bool      `json:"muted,omitempty"`
}

type NewsletterMessage struct {
	ID        string         `json:"id"`
	ServerID  string         `json:"server_id"`
	Content   string         `json:"content"`
	MediaType string         `json:"media_type,omitempty"`
	MediaURL  string         `json:"media_url,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Views     int            `json:"views,omitempty"`
	Reactions map[string]int `json:"reactions,omitempty"`
}

type serviceImpl struct {
	clients        map[string]*Client
	sessions       session.Repository
	logger         logging.Logger
	container      *sqlstore.Container
	waLogger       waLog.Logger
	config         config.MeowConfigProvider
	webhookService *webhooks.Service
	mu             sync.RWMutex
}

func NewService(container *sqlstore.Container, waLogger waLog.Logger, sessionRepo session.Repository, cfg config.MeowConfigProvider, webhookService *webhooks.Service) Service {
	return &serviceImpl{
		clients:        make(map[string]*Client),
		sessions:       sessionRepo,
		logger:         logging.GetLogger().Sub("meow"),
		container:      container,
		waLogger:       waLogger,
		config:         cfg,
		webhookService: webhookService,
	}
}

func (m *serviceImpl) StartClient(sessionID string) error {
	m.logger.Infof("Starting client for session %s", sessionID)
	client := m.getOrCreateClient(sessionID)
	return client.Connect()
}

func (m *serviceImpl) StopClient(sessionID string) error {
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

func (m *serviceImpl) LogoutClient(sessionID string) error {
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

func (m *serviceImpl) getClient(sessionID string) *Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[sessionID]
}

func (m *serviceImpl) getOrCreateClient(sessionID string) *Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[sessionID]; exists {
		return client
	}

	sessionEntity, err := m.sessions.GetByID(context.Background(), sessionID)
	var expectedDeviceJID string
	if err == nil && !sessionEntity.WaJID.IsEmpty() {
		expectedDeviceJID = sessionEntity.WaJID.Value()
		m.logger.Infof("Creating client for session %s with expected device JID: %s", sessionID, expectedDeviceJID)
	}

	eventProcessor := NewEventProcessor(sessionID, "", m.sessions, m.webhookService)

	client, err := NewClientWithDeviceJID(sessionID, expectedDeviceJID, m.container, m.waLogger, eventProcessor, m.sessions)
	if err != nil {
		m.logger.Errorf("Failed to create Client for session %s: %v", sessionID, err)
		return nil
	}

	m.clients[sessionID] = client
	return client
}

func (m *serviceImpl) validateAndGetClient(sessionID string) (*Client, error) {
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

func (m *serviceImpl) validateAndGetClientForSending(sessionID string) (*Client, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}
	return client, nil
}

func (m *serviceImpl) validateButtons(buttons []ButtonData) error {
	if len(buttons) == 0 {
		return fmt.Errorf("at least one button is required")
	}
	if len(buttons) > 3 {
		return fmt.Errorf("maximum 3 buttons allowed")
	}
	return nil
}

func (m *serviceImpl) buildmeowButtons(buttons []ButtonData) []*waProto.ButtonsMessage_Button {
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

func (m *serviceImpl) validateListSections(sections []ListSection) error {
	if len(sections) == 0 {
		return fmt.Errorf("at least one section is required")
	}
	return nil
}

func (m *serviceImpl) buildmeowListSections(sections []ListSection) []*waProto.ListMessage_Section {
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

func (m *serviceImpl) removeClient(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, sessionID)
}

func (m *serviceImpl) GetQRCode(sessionID string) (string, error) {
	client := m.getOrCreateClient(sessionID)
	if client == nil {
		return "", fmt.Errorf("failed to create client for session %s", sessionID)
	}
	return client.GetQRCode()
}

func (m *serviceImpl) PairPhone(sessionID, phoneNumber string) (string, error) {
	m.logger.Infof("Pairing phone %s for session %s", phoneNumber, sessionID)
	client := m.getOrCreateClient(sessionID)
	return client.PairPhone(phoneNumber)
}

func (m *serviceImpl) IsClientConnected(sessionID string) bool {
	client := m.getClient(sessionID)
	return client != nil && client.IsConnected()
}

func (m *serviceImpl) GetClientStatus(sessionID string) types.Status {
	client := m.getClient(sessionID)
	if client == nil {
		return types.StatusDisconnected
	}
	return client.GetStatus()
}

func (m *serviceImpl) ConnectOnStartup(ctx context.Context) error {
	m.logger.Infof("Connecting sessions with credentials on startup")

	sessions, err := m.sessions.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sessions with credentials: %w", err)
	}

	m.logger.Infof("Found %d sessions with credentials to reconnect", len(sessions))

	for _, sessionEntity := range sessions {
		if sessionEntity.WaJID.IsEmpty() {
			m.logger.Warnf("Session %s (%s) has no device_jid but was returned as active - fixing status",
				sessionEntity.ID, sessionEntity.Name)

			sessionEntity.Status = session.StatusDisconnected
			if err := m.sessions.Update(ctx, sessionEntity); err != nil {
				m.logger.Errorf("Failed to fix session %s status: %v", sessionEntity.ID, err)
			}
			continue
		}

		if !m.deviceExistsInDatabase(sessionEntity.WaJID.Value()) {
			m.logger.Warnf("Session %s (%s) has device_jid %s but device not found in whatsmeow_device table - marking as disconnected",
				sessionEntity.ID.Value(), sessionEntity.Name.Value(), sessionEntity.WaJID.Value())

			sessionEntity.Status = session.StatusDisconnected
			err := sessionEntity.SetWaJID("") // Clear invalid device_jid
			if err != nil {
				m.logger.Errorf("Failed to clear WaJID for session %s: %v", sessionEntity.ID.Value(), err)
			}
			if err := m.sessions.Update(ctx, sessionEntity); err != nil {
				m.logger.Errorf("Failed to fix session %s: %v", sessionEntity.ID, err)
			}
			continue
		}

		m.logger.Infof("Attempting to reconnect session %s (status: %s, wa_jid: %s)",
			sessionEntity.ID.Value(), sessionEntity.Status, sessionEntity.WaJID.Value())

		if err := m.StartClient(sessionEntity.ID.Value()); err != nil {
			m.logger.Errorf("Failed to start client for session %s: %v", sessionEntity.ID.Value(), err)
			continue
		}

		m.logger.Infof("Successfully initiated reconnection for session %s", sessionEntity.ID.Value())
	}

	return nil
}

func (m *serviceImpl) deviceExistsInDatabase(deviceJID string) bool {
	if deviceJID == "" {
		return false
	}

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

func (m *serviceImpl) DeleteMessage(ctx context.Context, sessionID, phone, messageID string, forEveryone bool) error {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	var revokeTarget waTypes.JID
	if forEveryone {
		revokeTarget = waTypes.EmptyJID // Delete for everyone
	} else {
		if client.GetClient().Store.ID == nil {
			return fmt.Errorf("unable to get client ID for delete operation")
		}
		revokeTarget = *client.GetClient().Store.ID // Delete for me only
	}

	revokeMsg := client.GetClient().BuildRevoke(recipient, revokeTarget, messageID)

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

func (m *serviceImpl) EditMessage(ctx context.Context, sessionID, phone, messageID, newText string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	editMsg := client.GetClient().BuildEdit(recipient, messageID, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: &newText,
		},
	})

	resp, err := client.GetClient().SendMessage(ctx, recipient, editMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to edit message: %w", err)
	}

	m.logger.Debugf("Message %s edited for phone %s in session %s",
		messageID, phone, sessionID)

	return &resp, nil
}

func (m *serviceImpl) DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error) {
	return nil, "", fmt.Errorf("media download functionality pending")
}

func (m *serviceImpl) ReactToMessage(ctx context.Context, sessionID, phone, messageID, emoji string) error {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return err
	}

	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	reaction := emoji
	if emoji == "remove" {
		reaction = ""
	}

	fromMe := false
	actualMessageID := messageID
	if strings.HasPrefix(messageID, "me:") {
		fromMe = true
		actualMessageID = messageID[len("me:"):]
	}

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

func (m *serviceImpl) SendTextMessage(ctx context.Context, sessionID, phone, text string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendTextMessage(client.GetClient(), phone, text)
}

func (m *serviceImpl) SendImageMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendImageMessage(client.GetClient(), phone, data, caption)
}

func (m *serviceImpl) SendAudioMessage(ctx context.Context, sessionID, phone string, data []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendAudioMessage(client.GetClient(), phone, data, mimeType)
}

func (m *serviceImpl) SendVideoMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendVideoMessage(client.GetClient(), phone, data, caption, mimeType)
}

func (m *serviceImpl) SendDocumentMessage(ctx context.Context, sessionID, phone string, data []byte, filename, caption, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendDocumentMessage(client.GetClient(), phone, data, filename, caption, mimeType)
}

func (m *serviceImpl) SendStickerMessage(ctx context.Context, sessionID, phone string, data []byte, mimeType string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendStickerMessage(client.GetClient(), phone, data, mimeType)
}

func (m *serviceImpl) SendContactMessage(ctx context.Context, sessionID, phone, contactName, contactPhone string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendContactMessage(client.GetClient(), phone, contactName, contactPhone)
}

func (m *serviceImpl) SendLocationMessage(ctx context.Context, sessionID, phone string, latitude, longitude float64, name, address string) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClientForSending(sessionID)
	if err != nil {
		return nil, err
	}
	return SendLocationMessage(client.GetClient(), phone, latitude, longitude, name, address)
}

func (m *serviceImpl) MarkAsRead(ctx context.Context, sessionID, phone string, messageIDs []string) error {
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

func (m *serviceImpl) SendButtonMessage(ctx context.Context, sessionID, phone, title string, buttons []ButtonData) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	if err := m.validateButtons(buttons); err != nil {
		return nil, err
	}

	waButtons := m.buildmeowButtons(buttons)

	buttonMsg := &waProto.Message{
		ButtonsMessage: &waProto.ButtonsMessage{
			ContentText: &title,
			HeaderType:  waProto.ButtonsMessage_EMPTY.Enum(),
			Buttons:     waButtons,
		},
	}

	resp, err := client.GetClient().SendMessage(ctx, recipient, buttonMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send button message: %w", err)
	}

	m.logger.Debugf("Button message sent to %s from session %s", phone, sessionID)

	return &resp, nil
}

func (m *serviceImpl) SendListMessage(ctx context.Context, sessionID, phone, title, description, buttonText, footerText string, sections []ListSection) (*whatsmeow.SendResponse, error) {
	client, err := m.validateAndGetClient(sessionID)
	if err != nil {
		return nil, err
	}

	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	if err := m.validateListSections(sections); err != nil {
		return nil, err
	}

	waSections := m.buildmeowListSections(sections)

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

	resp, err := client.GetClient().SendMessage(ctx, recipient, listMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send list message: %w", err)
	}

	m.logger.Debugf("List message sent to %s from session %s", phone, sessionID)

	return &resp, nil
}

func (m *serviceImpl) SendPollMessage(ctx context.Context, sessionID, phone, name string, options []string, selectableCount int) (*whatsmeow.SendResponse, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	recipient, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

	if len(options) < 2 {
		return nil, fmt.Errorf("at least 2 options are required for a poll")
	}
	if len(options) > 12 {
		return nil, fmt.Errorf("maximum 12 options allowed for a poll")
	}

	if selectableCount <= 0 {
		selectableCount = 1 // Default to single select
	}
	if selectableCount > len(options) {
		selectableCount = len(options)
	}

	pollMsg := client.GetClient().BuildPollCreation(name, options, selectableCount)

	resp, err := client.GetClient().SendMessage(ctx, recipient, pollMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send poll message: %w", err)
	}

	m.logger.Debugf("Poll message '%s' sent to %s from session %s", name, phone, sessionID)

	return &resp, nil
}

func (m *serviceImpl) SetPresence(ctx context.Context, sessionID, phone, state, media string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	var chatJID waTypes.JID
	if phone != "" {
		jid, err := parsePhoneToJID(phone)
		if err != nil {
			return fmt.Errorf("invalid phone number: %w", err)
		}
		chatJID = jid
	}

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

func (m *serviceImpl) CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

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

	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participantJIDs,
	}

	groupInfo, err := client.GetClient().CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

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

	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Group '%s' created successfully: %s for session %s",
		name, groupInfo.JID.String(), sessionID)

	return result, nil
}

func (m *serviceImpl) ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	groups, err := client.GetClient().GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get joined groups: %w", err)
	}

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

		for _, participant := range group.Participants {
			result.Participants = append(result.Participants, participant.JID.String())
		}

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

func (m *serviceImpl) GetGroupInfo(ctx context.Context, sessionID, groupJID string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	groupInfo, err := client.GetClient().GetGroupInfo(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

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

	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Retrieved group info for %s in session %s", groupJID, sessionID)

	return result, nil
}

func (m *serviceImpl) JoinGroup(ctx context.Context, sessionID, inviteLink string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	groupJID, err := client.GetClient().JoinGroupWithLink(inviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	groupInfo, err := client.GetClient().GetGroupInfo(groupJID)
	if err != nil {
		result := &GroupInfo{
			JID: groupJID.String(),
		}
		m.logger.Debugf("Successfully joined group %s via invite link for session %s (basic info)",
			groupJID.String(), sessionID)
		return result, nil
	}

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

	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Successfully joined group %s via invite link for session %s",
		groupInfo.JID.String(), sessionID)

	return result, nil
}

func (m *serviceImpl) JoinGroupWithInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	inviterJID, err := waTypes.ParseJID(inviter)
	if err != nil {
		return nil, fmt.Errorf("invalid inviter JID %s: %w", inviter, err)
	}

	err = client.GetClient().JoinGroupWithInvite(groupJIDParsed, inviterJID, code, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to join group with invite: %w", err)
	}

	groupInfo, err := client.GetClient().GetGroupInfo(groupJIDParsed)
	if err != nil {
		result := &GroupInfo{
			JID: groupJIDParsed.String(),
		}
		m.logger.Debugf("Successfully joined group %s via specific invite for session %s (basic info)",
			groupJIDParsed.String(), sessionID)
		return result, nil
	}

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

	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Successfully joined group %s via specific invite for session %s",
		groupJIDParsed.String(), sessionID)

	return result, nil
}

func (m *serviceImpl) LeaveGroup(ctx context.Context, sessionID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	err = client.GetClient().LeaveGroup(jid)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	m.logger.Debugf("Successfully left group %s for session %s", groupJID, sessionID)

	return nil
}

func (m *serviceImpl) GetInviteLink(ctx context.Context, sessionID, groupJID string, reset bool) (string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return "", fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return "", fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	inviteLink, err := client.GetClient().GetGroupInviteLink(jid, reset)
	if err != nil {
		return "", fmt.Errorf("failed to get invite link: %w", err)
	}

	m.logger.Debugf("Retrieved invite link for group %s in session %s", groupJID, sessionID)

	return inviteLink, nil
}

func (m *serviceImpl) SetGroupEphemeral(ctx context.Context, sessionID, groupJID string, ephemeral bool, duration int) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	var expiration time.Duration
	if ephemeral && duration > 0 {
		expiration = time.Duration(duration) * time.Second
	} else if ephemeral {
		expiration = 24 * time.Hour // Default 24 hours
	} else {
		expiration = 0 // Disable ephemeral
	}

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

func (m *serviceImpl) GetGroupRequestParticipants(ctx context.Context, sessionID, groupJID string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	requestParticipants, err := client.GetClient().GetGroupRequestParticipants(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group request participants: %w", err)
	}

	var participants []string
	for _, participant := range requestParticipants {
		participants = append(participants, participant.JID.String())
	}

	m.logger.Debugf("Retrieved %d group request participants for group %s in session %s",
		len(participants), groupJID, sessionID)

	return participants, nil
}

func (m *serviceImpl) UpdateGroupRequestParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

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

	var requestChange whatsmeow.ParticipantRequestChange
	switch action {
	case "approve":
		requestChange = whatsmeow.ParticipantChangeApprove
	case "reject":
		requestChange = whatsmeow.ParticipantChangeReject
	default:
		return fmt.Errorf("invalid action %s. Valid actions: approve, reject", action)
	}

	_, err = client.GetClient().UpdateGroupRequestParticipants(groupJIDParsed, participantJIDs, requestChange)
	if err != nil {
		return fmt.Errorf("failed to %s group request participants: %w", action, err)
	}

	m.logger.Debugf("Successfully %sed %d group request participants for group %s in session %s",
		action, len(participantJIDs), groupJID, sessionID)

	return nil
}

func (m *serviceImpl) LinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	err = client.GetClient().LinkGroup(communityJIDParsed, groupJIDParsed)
	if err != nil {
		return fmt.Errorf("failed to link group to community: %w", err)
	}

	m.logger.Debugf("Successfully linked group %s to community %s in session %s",
		groupJID, communityJID, sessionID)

	return nil
}

func (m *serviceImpl) UnlinkGroup(ctx context.Context, sessionID, communityJID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	err = client.GetClient().UnlinkGroup(communityJIDParsed, groupJIDParsed)
	if err != nil {
		return fmt.Errorf("failed to unlink group from community: %w", err)
	}

	m.logger.Debugf("Successfully unlinked group %s from community %s in session %s",
		groupJID, communityJID, sessionID)

	return nil
}

func (m *serviceImpl) GetSubGroups(ctx context.Context, sessionID, communityJID string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return nil, fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	subGroups, err := client.GetClient().GetSubGroups(communityJIDParsed)
	if err != nil {
		return nil, fmt.Errorf("failed to get subgroups: %w", err)
	}

	var groups []string
	for _, group := range subGroups {
		groups = append(groups, group.JID.String())
	}

	m.logger.Debugf("Retrieved %d subgroups for community %s in session %s",
		len(groups), communityJID, sessionID)

	return groups, nil
}

func (m *serviceImpl) GetLinkedGroupsParticipants(ctx context.Context, sessionID, communityJID string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	communityJIDParsed, err := waTypes.ParseJID(communityJID)
	if err != nil {
		return nil, fmt.Errorf("invalid community JID %s: %w", communityJID, err)
	}

	participants, err := client.GetClient().GetLinkedGroupsParticipants(communityJIDParsed)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked groups participants: %w", err)
	}

	var participantList []string
	for _, participant := range participants {
		participantList = append(participantList, participant.String())
	}

	m.logger.Debugf("Retrieved %d linked groups participants for community %s in session %s",
		len(participantList), communityJID, sessionID)

	return participantList, nil
}

func (m *serviceImpl) CheckUser(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

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

	resp, err := client.GetClient().IsOnWhatsApp(validPhones)
	if err != nil {
		return nil, fmt.Errorf("failed to check users on WhatsApp: %w", err)
	}

	var results []UserCheckResult
	for _, item := range resp {
		verifiedName := ""
		if item.VerifiedName != nil {
			verifiedName = item.VerifiedName.Details.GetVerifiedName()
		}

		result := UserCheckResult{
			Query:        item.Query,
			IsInmeow:     item.IsIn,
			JID:          item.JID.String(),
			VerifiedName: verifiedName,
		}
		results = append(results, result)
	}

	m.logger.Debugf("Checked %d users for session %s, found %d on meow",
		len(phones), sessionID, len(results))

	return results, nil
}

func (m *serviceImpl) GetUserInfo(ctx context.Context, sessionID string, phones []string) (map[string]UserInfoResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

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

	resp, err := client.GetClient().GetUserInfo(jids)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	results := make(map[string]UserInfoResult)
	for jid, info := range resp {
		verifiedName := ""
		if info.VerifiedName != nil {
			verifiedName = info.VerifiedName.Details.GetVerifiedName()
		}

		result := UserInfoResult{
			JID:          jid.String(),
			DisplayName:  "", // Will be populated from contact info if available
			VerifiedName: verifiedName,
			Avatar:       "", // Will be populated if available
			Status:       info.Status,
			PictureID:    info.PictureID,
			DeviceCount:  len(info.Devices),
		}
		results[jid.String()] = result
	}

	m.logger.Debugf("Retrieved info for %d users for session %s",
		len(results), sessionID)

	return results, nil
}

func (m *serviceImpl) GetAvatar(ctx context.Context, sessionID, phone string) (*AvatarResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := parsePhoneToJID(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number %s: %w", phone, err)
	}

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

func (m *serviceImpl) GetContacts(ctx context.Context, sessionID string) ([]ContactResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	contacts, err := client.GetClient().Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	var results []ContactResult
	for jid, contact := range contacts {
		result := ContactResult{
			JID:          jid.String(),
			Name:         contact.FullName,
			Notify:       contact.PushName,
			PushName:     contact.PushName,
			BusinessName: contact.BusinessName,
			IsBlocked:    false, // meow doesn't expose this easily
			IsMuted:      false, // meow doesn't expose this easily
		}
		results = append(results, result)
	}

	m.logger.Debugf("Retrieved %d contacts for session %s", len(results), sessionID)

	return results, nil
}

func (m *serviceImpl) SetUserPresence(ctx context.Context, sessionID, state string) error {
	return m.SetPresence(ctx, sessionID, "", state, "")
}

func (m *serviceImpl) GetInviteInfo(ctx context.Context, sessionID, inviteLink string) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	groupInfo, err := client.GetClient().GetGroupInfoFromLink(inviteLink)
	if err != nil {
		return nil, fmt.Errorf("failed to get invite info: %w", err)
	}

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

	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Retrieved invite info for group %s from link for session %s",
		groupInfo.JID.String(), sessionID)

	return result, nil
}

func (m *serviceImpl) GetGroupInfoFromInvite(ctx context.Context, sessionID, groupJID, inviter, code string, expiration int64) (*GroupInfo, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	inviterJID, err := waTypes.ParseJID(inviter)
	if err != nil {
		return nil, fmt.Errorf("invalid inviter JID %s: %w", inviter, err)
	}

	groupInfo, err := client.GetClient().GetGroupInfoFromInvite(groupJIDParsed, inviterJID, code, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info from invite: %w", err)
	}

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

	for _, participant := range groupInfo.Participants {
		result.Participants = append(result.Participants, participant.JID.String())
	}

	for _, participant := range groupInfo.Participants {
		if participant.IsAdmin {
			result.Admins = append(result.Admins, participant.JID.String())
		}
	}

	m.logger.Debugf("Retrieved group info from specific invite for group %s in session %s",
		groupInfo.JID.String(), sessionID)

	return result, nil
}

func (m *serviceImpl) UpdateParticipants(ctx context.Context, sessionID, groupJID, action string, participants []string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	groupJIDParsed, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

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

	_, err = client.GetClient().UpdateGroupParticipants(groupJIDParsed, participantJIDs, participantChange)
	if err != nil {
		return fmt.Errorf("failed to %s participants: %w", action, err)
	}

	m.logger.Debugf("Successfully %s %d participants in group %s for session %s",
		action, len(participantJIDs), groupJID, sessionID)

	return nil
}

func (m *serviceImpl) SetGroupName(ctx context.Context, sessionID, groupJID, name string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	err = client.GetClient().SetGroupName(jid, name)
	if err != nil {
		return fmt.Errorf("failed to set group name: %w", err)
	}

	m.logger.Debugf("Successfully set group name to '%s' for group %s in session %s",
		name, groupJID, sessionID)

	return nil
}

func (m *serviceImpl) SetGroupTopic(ctx context.Context, sessionID, groupJID, topic string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	err = client.GetClient().SetGroupTopic(jid, "", "", topic)
	if err != nil {
		return fmt.Errorf("failed to set group topic: %w", err)
	}

	m.logger.Debugf("Successfully set group topic to '%s' for group %s in session %s",
		topic, groupJID, sessionID)

	return nil
}

func (m *serviceImpl) SetGroupPhoto(ctx context.Context, sessionID, groupJID string, photo []byte) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	if len(photo) == 0 {
		return fmt.Errorf("photo data cannot be empty")
	}

	pictureID, err := client.GetClient().SetGroupPhoto(jid, photo)
	if err != nil {
		return fmt.Errorf("failed to set group photo: %w", err)
	}

	m.logger.Debugf("Successfully set group photo (ID: %s) for group %s in session %s",
		pictureID, groupJID, sessionID)

	return nil
}

func (m *serviceImpl) RemoveGroupPhoto(ctx context.Context, sessionID, groupJID string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	_, err = client.GetClient().SetGroupPhoto(jid, nil)
	if err != nil {
		return fmt.Errorf("failed to remove group photo: %w", err)
	}

	m.logger.Debugf("Successfully removed group photo for group %s in session %s",
		groupJID, sessionID)

	return nil
}

func (m *serviceImpl) SetGroupAnnounce(ctx context.Context, sessionID, groupJID string, announceOnly bool) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

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

func (m *serviceImpl) SetGroupLocked(ctx context.Context, sessionID, groupJID string, locked bool) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

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

func (m *serviceImpl) SetGroupJoinApprovalMode(ctx context.Context, sessionID, groupJID string, requireApproval bool) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

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

func (m *serviceImpl) SetGroupMemberAddMode(ctx context.Context, sessionID, groupJID string, mode string) error {
	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID %s: %w", groupJID, err)
	}

	var memberAddMode waTypes.GroupMemberAddMode
	switch strings.ToLower(mode) {
	case "all", "everyone":
		memberAddMode = waTypes.GroupMemberAddModeAllMember
	case "admin", "admins", "admin_only":
		memberAddMode = waTypes.GroupMemberAddModeAdmin
	default:
		return fmt.Errorf("invalid member add mode: %s. Valid options: all, admin", mode)
	}

	err = client.GetClient().SetGroupMemberAddMode(jid, memberAddMode)
	if err != nil {
		return fmt.Errorf("failed to set group member add mode: %w", err)
	}

	m.logger.Debugf("Successfully set group member add mode to %s for group %s in session %s",
		mode, groupJID, sessionID)

	return nil
}

func (m *serviceImpl) UpdateSessionWebhook(sessionID, webhookURL string) error {
	m.logger.Infof("Updated webhook URL for session %s", sessionID)
	return nil
}

func (m *serviceImpl) UpdateSessionSubscriptions(sessionID string, events []string) error {
	m.logger.Infof("Updated event subscriptions for session %s", sessionID)
	return nil
}

func (m *serviceImpl) GetPrivacySettings(ctx context.Context, sessionID string) (*PrivacySettingsResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	settings := client.GetClient().GetPrivacySettings(ctx)

	result := &PrivacySettingsResult{
		GroupAdd:     string(settings.GroupAdd),
		LastSeen:     string(settings.LastSeen),
		Status:       string(settings.Status),
		Profile:      string(settings.Profile),
		ReadReceipts: string(settings.ReadReceipts),
		CallAdd:      string(settings.CallAdd),
		Online:       string(settings.Online),
	}

	m.logger.Debugf("Successfully retrieved privacy settings for session %s", sessionID)
	return result, nil
}

func (m *serviceImpl) TryFetchPrivacySettings(ctx context.Context, sessionID string, ignoreCache bool) (*PrivacySettingsResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	settings, err := client.GetClient().TryFetchPrivacySettings(ctx, ignoreCache)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch privacy settings from server: %w", err)
	}

	result := &PrivacySettingsResult{
		GroupAdd:     string(settings.GroupAdd),
		LastSeen:     string(settings.LastSeen),
		Status:       string(settings.Status),
		Profile:      string(settings.Profile),
		ReadReceipts: string(settings.ReadReceipts),
		CallAdd:      string(settings.CallAdd),
		Online:       string(settings.Online),
	}

	m.logger.Debugf("Successfully fetched privacy settings from server for session %s (ignoreCache: %v)", sessionID, ignoreCache)
	return result, nil
}

func (m *serviceImpl) SetPrivacySetting(ctx context.Context, sessionID string, settingType, settingValue string) (*PrivacySettingsResult, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	privacySettingType := waTypes.PrivacySettingType(settingType)
	privacySetting := waTypes.PrivacySetting(settingValue)

	m.logger.Debugf("🔧 Setting privacy: type='%s' value='%s' for session %s", settingType, settingValue, sessionID)
	m.logger.Debugf("🔧 Converted types: privacySettingType='%s' privacySetting='%s'", string(privacySettingType), string(privacySetting))

	settings, err := client.GetClient().SetPrivacySetting(ctx, privacySettingType, privacySetting)
	if err != nil {
		m.logger.Errorf("❌ Failed to set privacy setting %s to %s: %v", settingType, settingValue, err)
		return nil, fmt.Errorf("failed to set privacy setting %s to %s: %w", settingType, settingValue, err)
	}

	result := &PrivacySettingsResult{
		GroupAdd:     string(settings.GroupAdd),
		LastSeen:     string(settings.LastSeen),
		Status:       string(settings.Status),
		Profile:      string(settings.Profile),
		ReadReceipts: string(settings.ReadReceipts),
		CallAdd:      string(settings.CallAdd),
		Online:       string(settings.Online),
	}

	m.logger.Debugf("✅ Successfully set privacy setting %s to %s for session %s", settingType, settingValue, sessionID)
	m.logger.Debugf("📊 Updated settings: %+v", result)

	return result, nil
}

func (m *serviceImpl) GetBlocklist(ctx context.Context, sessionID string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	blocklist, err := client.GetClient().GetBlocklist()
	if err != nil {
		m.logger.Errorf("❌ Failed to get blocklist for session %s: %v", sessionID, err)
		return nil, fmt.Errorf("failed to get blocklist: %w", err)
	}

	jidList := make([]string, len(blocklist.JIDs))
	for i, jid := range blocklist.JIDs {
		jidList[i] = jid.String()
	}

	m.logger.Debugf("✅ Successfully retrieved blocklist for session %s: %d entries", sessionID, len(jidList))
	return jidList, nil
}

func (m *serviceImpl) UpdateBlocklist(ctx context.Context, sessionID, jidStr, action string) ([]string, error) {
	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(jidStr)
	if err != nil {
		m.logger.Errorf("❌ Failed to parse JID %s for session %s: %v", jidStr, sessionID, err)
		return nil, fmt.Errorf("invalid JID format: %w", err)
	}

	if action != "block" && action != "unblock" {
		return nil, fmt.Errorf("invalid action '%s', must be 'block' or 'unblock'", action)
	}

	var blocklistAction waEvents.BlocklistChangeAction
	switch action {
	case "block":
		blocklistAction = waEvents.BlocklistChangeActionBlock
	case "unblock":
		blocklistAction = waEvents.BlocklistChangeActionUnblock
	default:
		return nil, fmt.Errorf("invalid action '%s', must be 'block' or 'unblock'", action)
	}

	updatedBlocklist, err := client.GetClient().UpdateBlocklist(jid, blocklistAction)
	if err != nil {
		m.logger.Errorf("❌ Failed to %s user %s for session %s: %v", action, jidStr, sessionID, err)
		return nil, fmt.Errorf("failed to %s user: %w", action, err)
	}

	jidList := make([]string, len(updatedBlocklist.JIDs))
	for i, jid := range updatedBlocklist.JIDs {
		jidList[i] = jid.String()
	}

	m.logger.Debugf("✅ Successfully %sed user %s for session %s. New blocklist has %d entries", action, jidStr, sessionID, len(jidList))
	return jidList, nil
}

func (m *serviceImpl) SetDisappearingTimer(ctx context.Context, sessionID, chatJID string, timer time.Duration) error {
	m.logger.Debugf("🕐 Setting disappearing timer for chat %s in session %s to %v", chatJID, sessionID, timer)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	err = client.GetClient().SetDisappearingTimer(jid, timer, time.Now())
	if err != nil {
		m.logger.Errorf("❌ Failed to set disappearing timer for chat %s in session %s: %v", chatJID, sessionID, err)
		return fmt.Errorf("failed to set disappearing timer: %w", err)
	}

	timerStr := "off"
	if timer > 0 {
		timerStr = timer.String()
	}
	m.logger.Debugf("✅ Successfully set disappearing timer for chat %s in session %s to %s", chatJID, sessionID, timerStr)
	return nil
}

func (m *serviceImpl) ListChats(ctx context.Context, sessionID string, chatType string) ([]ChatInfo, error) {
	m.logger.Debugf("📋 Listing chats for session %s (type: %s)", sessionID, chatType)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	var chats []ChatInfo

	if chatType == "all" || chatType == "groups" {
		groups, err := client.GetClient().GetJoinedGroups()
		if err != nil {
			m.logger.Errorf("❌ Failed to get joined groups for session %s: %v", sessionID, err)
			return nil, fmt.Errorf("failed to get joined groups: %w", err)
		}

		for _, group := range groups {
			chatInfo := ChatInfo{
				JID:         group.JID.String(),
				Name:        group.Name,
				Type:        "group",
				UnreadCount: 0,     // Future: Implement unread count from app state
				Pinned:      false, // Future: Get from app state
				Muted:       false, // Future: Get from app state
				Archived:    false, // Future: Get from app state
			}
			chats = append(chats, chatInfo)
		}
	}

	if chatType == "all" || chatType == "contacts" {
		contacts, err := m.GetContacts(ctx, sessionID)
		if err != nil {
			m.logger.Errorf("❌ Failed to get contacts for session %s: %v", sessionID, err)
		} else {
			for _, contact := range contacts {
				chatInfo := ChatInfo{
					JID:         contact.JID,
					Name:        contact.Name,
					Type:        "contact",
					UnreadCount: 0,     // Future: Implement unread count from app state
					Pinned:      false, // Future: Get from app state
					Muted:       contact.IsMuted,
					Archived:    false, // Future: Get from app state
				}
				chats = append(chats, chatInfo)
			}
		}
	}

	m.logger.Debugf("✅ Successfully listed %d chats for session %s", len(chats), sessionID)
	return chats, nil
}

func (m *serviceImpl) GetChatInfo(ctx context.Context, sessionID, chatJID string) (*ChatInfo, error) {
	m.logger.Debugf("ℹ️ Getting chat info for %s in session %s", chatJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat JID: %w", err)
	}

	var chatInfo ChatInfo

	if jid.Server == waTypes.GroupServer {
		groupInfo, err := client.GetClient().GetGroupInfo(jid)
		if err != nil {
			m.logger.Errorf("❌ Failed to get group info for %s in session %s: %v", chatJID, sessionID, err)
			return nil, fmt.Errorf("failed to get group info: %w", err)
		}

		chatInfo = ChatInfo{
			JID:         groupInfo.JID.String(),
			Name:        groupInfo.Name,
			Type:        "group",
			UnreadCount: 0,     // Future: Implement unread count from app state
			Pinned:      false, // Future: Get from app state
			Muted:       false, // Future: Get from app state
			Archived:    false, // Future: Get from app state
		}
	} else {
		userInfo, err := m.GetUserInfo(ctx, sessionID, []string{chatJID})
		if err != nil {
			m.logger.Errorf("❌ Failed to get user info for %s in session %s: %v", chatJID, sessionID, err)
			return nil, fmt.Errorf("failed to get user info: %w", err)
		}

		if info, exists := userInfo[chatJID]; exists {
			chatInfo = ChatInfo{
				JID:         info.JID,
				Name:        info.DisplayName,
				Type:        "contact",
				UnreadCount: 0,     // Future: Implement unread count from app state
				Pinned:      false, // Future: Get from app state
				Muted:       false, // Future: Get from app state
				Archived:    false, // Future: Get from app state
			}
		} else {
			return nil, fmt.Errorf("chat not found: %s", chatJID)
		}
	}

	m.logger.Debugf("✅ Successfully got chat info for %s in session %s", chatJID, sessionID)
	return &chatInfo, nil
}

func (m *serviceImpl) PinChat(ctx context.Context, sessionID, chatJID string, pinned bool) error {
	m.logger.Debugf("📌 %s chat %s in session %s", map[bool]string{true: "Pinning", false: "Unpinning"}[pinned], chatJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	patch := appstate.BuildPin(jid, pinned)

	err = client.GetClient().SendAppState(ctx, patch)
	if err != nil {
		m.logger.Errorf("❌ Failed to %s chat %s in session %s: %v", map[bool]string{true: "pin", false: "unpin"}[pinned], chatJID, sessionID, err)
		return fmt.Errorf("failed to %s chat: %w", map[bool]string{true: "pin", false: "unpin"}[pinned], err)
	}

	m.logger.Debugf("✅ Successfully %s chat %s in session %s", map[bool]string{true: "pinned", false: "unpinned"}[pinned], chatJID, sessionID)
	return nil
}

func (m *serviceImpl) MuteChat(ctx context.Context, sessionID, chatJID string, muted bool, duration time.Duration) error {
	m.logger.Debugf("🔇 %s chat %s in session %s for %v", map[bool]string{true: "Muting", false: "Unmuting"}[muted], chatJID, sessionID, duration)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	patch := appstate.BuildMute(jid, muted, duration)

	err = client.GetClient().SendAppState(ctx, patch)
	if err != nil {
		m.logger.Errorf("❌ Failed to %s chat %s in session %s: %v", map[bool]string{true: "mute", false: "unmute"}[muted], chatJID, sessionID, err)
		return fmt.Errorf("failed to %s chat: %w", map[bool]string{true: "mute", false: "unmute"}[muted], err)
	}

	durationStr := "permanently"
	if duration > 0 {
		durationStr = "for " + duration.String()
	}
	m.logger.Debugf("✅ Successfully %s chat %s in session %s %s", map[bool]string{true: "muted", false: "unmuted"}[muted], chatJID, sessionID, durationStr)
	return nil
}

func (m *serviceImpl) ArchiveChat(ctx context.Context, sessionID, chatJID string, archived bool) error {
	m.logger.Debugf("📦 %s chat %s in session %s", map[bool]string{true: "Archiving", false: "Unarchiving"}[archived], chatJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("invalid chat JID: %w", err)
	}

	patch := appstate.BuildArchive(jid, archived, time.Time{}, nil)

	err = client.GetClient().SendAppState(ctx, patch)
	if err != nil {
		m.logger.Errorf("❌ Failed to %s chat %s in session %s: %v", map[bool]string{true: "archive", false: "unarchive"}[archived], chatJID, sessionID, err)
		return fmt.Errorf("failed to %s chat: %w", map[bool]string{true: "archive", false: "unarchive"}[archived], err)
	}

	m.logger.Debugf("✅ Successfully %s chat %s in session %s", map[bool]string{true: "archived", false: "unarchived"}[archived], chatJID, sessionID)
	return nil
}

func (m *serviceImpl) CreateNewsletter(ctx context.Context, sessionID string, params *whatsmeow.CreateNewsletterParams) (*waTypes.NewsletterMetadata, error) {
	m.logger.Debugf("📰 Creating newsletter '%s' in session %s", params.Name, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	if params.Name == "" {
		return nil, fmt.Errorf("newsletter name is required")
	}

	resp, err := client.GetClient().CreateNewsletter(*params)
	if err != nil {
		m.logger.Errorf("❌ Failed to create newsletter '%s' in session %s: %v", params.Name, sessionID, err)
		return nil, fmt.Errorf("failed to create newsletter: %w", err)
	}

	m.logger.Debugf("✅ Successfully created newsletter '%s' with JID %s in session %s", params.Name, resp.ID, sessionID)
	return resp, nil
}

func (m *serviceImpl) UploadNewsletterReader(ctx context.Context, sessionID string, data []byte, mediaType whatsmeow.MediaType) (*whatsmeow.UploadResponse, error) {
	m.logger.Debugf("📰 Uploading newsletter media (reader) for session %s", sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	reader := bytes.NewReader(data)

	resp, err := client.GetClient().UploadNewsletterReader(ctx, reader, mediaType)
	if err != nil {
		m.logger.Errorf("❌ Failed to upload newsletter media (reader) for session %s: %v", sessionID, err)
		return nil, fmt.Errorf("failed to upload newsletter media: %w", err)
	}

	m.logger.Debugf("✅ Successfully uploaded newsletter media (reader) for session %s", sessionID)
	return &resp, nil
}

func (m *serviceImpl) GetNewsletterInfo(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error) {
	m.logger.Debugf("📰 Getting newsletter info for %s in session %s", newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return nil, fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	info, err := client.GetClient().GetNewsletterInfo(jid)
	if err != nil {
		m.logger.Errorf("❌ Failed to get newsletter info for %s in session %s: %v", newsletterJID, sessionID, err)
		return nil, fmt.Errorf("failed to get newsletter info: %w", err)
	}

	result := &NewsletterInfo{
		JID:         info.ID.String(),
		Name:        info.ThreadMeta.Name.Text,
		Description: info.ThreadMeta.Description.Text,
		Picture:     "", // Picture extraction not available in current whatsmeow version
		Verified:    info.ThreadMeta.VerificationState == waTypes.NewsletterVerificationStateVerified,
		CreatedAt:   time.Now(), // Using current time - jsontime conversion complex
		UpdatedAt:   time.Now(), // Using current time - jsontime conversion complex
		OwnerJID:    "",         // Not available in NewsletterMetadata
		Subscribers: info.ThreadMeta.SubscriberCount,
		Muted:       info.ViewerMeta != nil && info.ViewerMeta.Mute == waTypes.NewsletterMuteOn,
	}

	m.logger.Debugf("✅ Successfully got newsletter info for %s in session %s", newsletterJID, sessionID)
	return result, nil
}

func (m *serviceImpl) GetNewsletterInfoWithInvite(ctx context.Context, sessionID, inviteKey string) (*NewsletterInfo, error) {
	m.logger.Debugf("📰 Getting newsletter info with invite key in session %s", sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	if inviteKey == "" {
		return nil, fmt.Errorf("invite key is required")
	}

	info, err := client.GetClient().GetNewsletterInfoWithInvite(inviteKey)
	if err != nil {
		m.logger.Errorf("❌ Failed to get newsletter info with invite in session %s: %v", sessionID, err)
		return nil, fmt.Errorf("failed to get newsletter info with invite: %w", err)
	}

	result := &NewsletterInfo{
		JID:         info.ID.String(),
		Name:        info.ThreadMeta.Name.Text,
		Description: info.ThreadMeta.Description.Text,
		Picture:     "", // Picture extraction not available in current whatsmeow version
		Verified:    info.ThreadMeta.VerificationState == waTypes.NewsletterVerificationStateVerified,
		CreatedAt:   time.Now(), // Using current time - jsontime conversion complex
		UpdatedAt:   time.Now(), // Using current time - jsontime conversion complex
		OwnerJID:    "",         // Not available in NewsletterMetadata
		Subscribers: info.ThreadMeta.SubscriberCount,
		Muted:       info.ViewerMeta != nil && info.ViewerMeta.Mute == waTypes.NewsletterMuteOn,
	}

	m.logger.Debugf("✅ Successfully got newsletter info with invite in session %s", sessionID)
	return result, nil
}

func (m *serviceImpl) FollowNewsletter(ctx context.Context, sessionID, newsletterJID string) error {
	m.logger.Debugf("📰 Following newsletter %s in session %s", newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	err = client.GetClient().FollowNewsletter(jid)
	if err != nil {
		m.logger.Errorf("❌ Failed to follow newsletter %s in session %s: %v", newsletterJID, sessionID, err)
		return fmt.Errorf("failed to follow newsletter: %w", err)
	}

	m.logger.Debugf("✅ Successfully followed newsletter %s in session %s", newsletterJID, sessionID)
	return nil
}

func (m *serviceImpl) UnfollowNewsletter(ctx context.Context, sessionID, newsletterJID string) error {
	m.logger.Debugf("📰 Unfollowing newsletter %s in session %s", newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	err = client.GetClient().UnfollowNewsletter(jid)
	if err != nil {
		m.logger.Errorf("❌ Failed to unfollow newsletter %s in session %s: %v", newsletterJID, sessionID, err)
		return fmt.Errorf("failed to unfollow newsletter: %w", err)
	}

	m.logger.Debugf("✅ Successfully unfollowed newsletter %s in session %s", newsletterJID, sessionID)
	return nil
}

func (m *serviceImpl) GetSubscribedNewsletters(ctx context.Context, sessionID string) ([]NewsletterInfo, error) {
	m.logger.Debugf("📰 Getting subscribed newsletters in session %s", sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	newsletters, err := client.GetClient().GetSubscribedNewsletters()
	if err != nil {
		m.logger.Errorf("❌ Failed to get subscribed newsletters in session %s: %v", sessionID, err)
		return nil, fmt.Errorf("failed to get subscribed newsletters: %w", err)
	}

	result := make([]NewsletterInfo, len(newsletters))
	for i, newsletter := range newsletters {
		result[i] = NewsletterInfo{
			JID:         newsletter.ID.String(),
			Name:        newsletter.ThreadMeta.Name.Text,
			Description: newsletter.ThreadMeta.Description.Text,
			Picture:     "", // Picture extraction not available in current whatsmeow version
			Verified:    newsletter.ThreadMeta.VerificationState == waTypes.NewsletterVerificationStateVerified,
			CreatedAt:   time.Now(), // Using current time - jsontime conversion complex
			UpdatedAt:   time.Now(), // Using current time - jsontime conversion complex
			OwnerJID:    "",         // Not available in NewsletterMetadata
			Subscribers: newsletter.ThreadMeta.SubscriberCount,
			Muted:       newsletter.ViewerMeta != nil && newsletter.ViewerMeta.Mute == waTypes.NewsletterMuteOn,
		}
	}

	m.logger.Debugf("✅ Successfully got %d subscribed newsletters in session %s", len(result), sessionID)
	return result, nil
}

func (m *serviceImpl) GetNewsletterMessages(ctx context.Context, sessionID, newsletterJID string, params *whatsmeow.GetNewsletterMessagesParams) ([]NewsletterMessage, error) {
	m.logger.Debugf("📰 Getting newsletter messages for %s in session %s", newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return nil, fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	messages, err := client.GetClient().GetNewsletterMessages(jid, params)
	if err != nil {
		m.logger.Errorf("❌ Failed to get newsletter messages for %s in session %s: %v", newsletterJID, sessionID, err)
		return nil, fmt.Errorf("failed to get newsletter messages: %w", err)
	}

	result := make([]NewsletterMessage, len(messages))
	for i, msg := range messages {
		result[i] = NewsletterMessage{
			ID:        string(msg.MessageID),
			ServerID:  fmt.Sprintf("%d", msg.MessageServerID),
			Content:   extractMessageContent(msg),
			MediaType: extractMediaType(msg),
			MediaURL:  extractMediaURL(msg),
			Timestamp: msg.Timestamp,
			Views:     msg.ViewsCount,
			Reactions: msg.ReactionCounts,
		}
	}

	m.logger.Debugf("✅ Successfully got %d newsletter messages for %s in session %s", len(result), newsletterJID, sessionID)
	return result, nil
}

func (m *serviceImpl) GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterJID string, params *whatsmeow.GetNewsletterMessagesParams) ([]NewsletterMessage, error) {
	m.logger.Debugf("📰 Getting newsletter message updates for %s in session %s", newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return nil, fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	updates, err := client.GetClient().GetNewsletterMessages(jid, params)
	if err != nil {
		m.logger.Errorf("❌ Failed to get newsletter message updates for %s in session %s: %v", newsletterJID, sessionID, err)
		return nil, fmt.Errorf("failed to get newsletter message updates: %w", err)
	}

	result := make([]NewsletterMessage, len(updates))
	for i, update := range updates {
		result[i] = NewsletterMessage{
			ID:        string(update.MessageID),
			ServerID:  fmt.Sprintf("%d", update.MessageServerID),
			Content:   extractMessageContent(update),
			MediaType: extractMediaType(update),
			MediaURL:  extractMediaURL(update),
			Timestamp: update.Timestamp,
			Views:     update.ViewsCount,
			Reactions: update.ReactionCounts,
		}
	}

	m.logger.Debugf("✅ Successfully got %d newsletter message updates for %s in session %s", len(result), newsletterJID, sessionID)
	return result, nil
}

func extractMessageContent(msg *waTypes.NewsletterMessage) string {
	if msg.Message == nil {
		return ""
	}

	if msg.Message.Conversation != nil {
		return *msg.Message.Conversation
	}

	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.Text != nil {
		return *msg.Message.ExtendedTextMessage.Text
	}

	if msg.Message.ImageMessage != nil && msg.Message.ImageMessage.Caption != nil {
		return *msg.Message.ImageMessage.Caption
	}

	if msg.Message.VideoMessage != nil && msg.Message.VideoMessage.Caption != nil {
		return *msg.Message.VideoMessage.Caption
	}

	if msg.Message.DocumentMessage != nil && msg.Message.DocumentMessage.Caption != nil {
		return *msg.Message.DocumentMessage.Caption
	}

	return ""
}

func extractMediaType(msg *waTypes.NewsletterMessage) string {
	if msg.Message == nil {
		return ""
	}

	if msg.Message.ImageMessage != nil {
		return "image"
	}

	if msg.Message.VideoMessage != nil {
		return "video"
	}

	if msg.Message.AudioMessage != nil {
		return "audio"
	}

	if msg.Message.DocumentMessage != nil {
		return "document"
	}

	if msg.Message.StickerMessage != nil {
		return "sticker"
	}

	return "text"
}

func extractMediaURL(msg *waTypes.NewsletterMessage) string {
	if msg.Message == nil {
		return ""
	}

	if msg.Message.ImageMessage != nil && msg.Message.ImageMessage.URL != nil {
		return *msg.Message.ImageMessage.URL
	}

	if msg.Message.VideoMessage != nil && msg.Message.VideoMessage.URL != nil {
		return *msg.Message.VideoMessage.URL
	}

	if msg.Message.AudioMessage != nil && msg.Message.AudioMessage.URL != nil {
		return *msg.Message.AudioMessage.URL
	}

	if msg.Message.DocumentMessage != nil && msg.Message.DocumentMessage.URL != nil {
		return *msg.Message.DocumentMessage.URL
	}

	return ""
}

func (m *serviceImpl) NewsletterMarkViewed(ctx context.Context, sessionID, newsletterJID string, serverIDs []waTypes.MessageServerID) error {
	m.logger.Debugf("📰 Marking %d newsletter messages as viewed for %s in session %s", len(serverIDs), newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	err = client.GetClient().NewsletterMarkViewed(jid, serverIDs)
	if err != nil {
		m.logger.Errorf("❌ Failed to mark newsletter messages as viewed for %s in session %s: %v", newsletterJID, sessionID, err)
		return fmt.Errorf("failed to mark newsletter messages as viewed: %w", err)
	}

	m.logger.Debugf("✅ Successfully marked %d newsletter messages as viewed for %s in session %s", len(serverIDs), newsletterJID, sessionID)
	return nil
}

func (m *serviceImpl) NewsletterSendReaction(ctx context.Context, sessionID, newsletterJID string, serverID waTypes.MessageServerID, reaction string, messageID waTypes.MessageID) error {
	m.logger.Debugf("📰 Sending reaction '%s' to newsletter message %s in %s for session %s", reaction, serverID, newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	if reaction == "" {
		return fmt.Errorf("reaction cannot be empty")
	}

	err = client.GetClient().NewsletterSendReaction(jid, serverID, reaction, messageID)
	if err != nil {
		m.logger.Errorf("❌ Failed to send reaction to newsletter message %s in %s for session %s: %v", serverID, newsletterJID, sessionID, err)
		return fmt.Errorf("failed to send newsletter reaction: %w", err)
	}

	m.logger.Debugf("✅ Successfully sent reaction '%s' to newsletter message %s in %s for session %s", reaction, serverID, newsletterJID, sessionID)
	return nil
}

func (m *serviceImpl) NewsletterToggleMute(ctx context.Context, sessionID, newsletterJID string, mute bool) error {
	m.logger.Debugf("📰 %s newsletter %s in session %s", map[bool]string{true: "Muting", false: "Unmuting"}[mute], newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	err = client.GetClient().NewsletterToggleMute(jid, mute)
	if err != nil {
		m.logger.Errorf("❌ Failed to %s newsletter %s in session %s: %v", map[bool]string{true: "mute", false: "unmute"}[mute], newsletterJID, sessionID, err)
		return fmt.Errorf("failed to %s newsletter: %w", map[bool]string{true: "mute", false: "unmute"}[mute], err)
	}

	m.logger.Debugf("✅ Successfully %s newsletter %s in session %s", map[bool]string{true: "muted", false: "unmuted"}[mute], newsletterJID, sessionID)
	return nil
}

func (m *serviceImpl) NewsletterSubscribeLiveUpdates(ctx context.Context, sessionID, newsletterJID string) error {
	m.logger.Debugf("📰 Subscribing to live updates for newsletter %s in session %s", newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	_, err = client.GetClient().NewsletterSubscribeLiveUpdates(ctx, jid)
	if err != nil {
		m.logger.Errorf("❌ Failed to subscribe to live updates for newsletter %s in session %s: %v", newsletterJID, sessionID, err)
		return fmt.Errorf("failed to subscribe to newsletter live updates: %w", err)
	}

	m.logger.Debugf("✅ Successfully subscribed to live updates for newsletter %s in session %s", newsletterJID, sessionID)
	return nil
}

func (m *serviceImpl) UploadNewsletter(ctx context.Context, sessionID string, data []byte, mediaType whatsmeow.MediaType) (*whatsmeow.UploadResponse, error) {
	m.logger.Debugf("📰 Uploading newsletter media (%d bytes, type: %s) in session %s", len(data), mediaType, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("media data cannot be empty")
	}

	resp, err := client.GetClient().UploadNewsletter(ctx, data, mediaType)
	if err != nil {
		m.logger.Errorf("❌ Failed to upload newsletter media in session %s: %v", sessionID, err)
		return nil, fmt.Errorf("failed to upload newsletter media: %w", err)
	}

	m.logger.Debugf("✅ Successfully uploaded newsletter media (%d bytes) in session %s", len(data), sessionID)
	return &resp, nil
}

func (m *serviceImpl) SendNewsletterMessage(ctx context.Context, sessionID, newsletterJID string, message *waE2E.Message, mediaHandle string) (*whatsmeow.SendResponse, error) {
	m.logger.Debugf("📰 Sending message to newsletter %s in session %s", newsletterJID, sessionID)

	client := m.getClient(sessionID)
	if client == nil {
		return nil, fmt.Errorf("client not found for session %s", sessionID)
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("client not connected for session %s", sessionID)
	}

	jid, err := waTypes.ParseJID(newsletterJID)
	if err != nil {
		return nil, fmt.Errorf("invalid newsletter JID %s: %w", newsletterJID, err)
	}

	if message == nil {
		return nil, fmt.Errorf("message cannot be nil")
	}

	extra := whatsmeow.SendRequestExtra{}

	isMediaMessage := message.ImageMessage != nil || message.VideoMessage != nil ||
		message.AudioMessage != nil || message.DocumentMessage != nil ||
		message.StickerMessage != nil

	if isMediaMessage {
		if mediaHandle == "" {
			return nil, fmt.Errorf("MediaHandle is required for media messages in newsletters")
		}
		extra.MediaHandle = mediaHandle
		m.logger.Debugf("📰 Using MediaHandle for newsletter media message: %s", mediaHandle)
	}

	resp, err := client.GetClient().SendMessage(ctx, jid, message, extra)
	if err != nil {
		m.logger.Errorf("❌ Failed to send message to newsletter %s in session %s: %v", newsletterJID, sessionID, err)
		return nil, fmt.Errorf("failed to send newsletter message: %w", err)
	}

	m.logger.Debugf("✅ Successfully sent message to newsletter %s in session %s", newsletterJID, sessionID)
	return &resp, nil
}
