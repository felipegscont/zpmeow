package infra

// Re-export infrastructure types to avoid aliases
import (
	"zpmeow/internal/infra/database/models"
	"zpmeow/internal/infra/database/repositories"
	"zpmeow/internal/infra/webhook"
	"zpmeow/internal/infra/meow/core"
	"zpmeow/internal/infra/meow/service"
	"zpmeow/internal/infra/meow/adapter"
)

// Database types
type SessionModel = models.SessionModel
type PostgresSessionRepository = repositories.PostgresSessionRepository

// Database constructors
var (
	NewPostgresSessionRepository = repositories.NewPostgresSessionRepository
)

// HTTP utilities - will be populated when utils are properly structured

// Webhook types
type WebhookService = webhook.WebhookService
type WebhookHTTPClient = webhook.WebhookHTTPClient
type HTTPClient = webhook.HTTPClient

// Webhook constructors
var (
	NewWebhookService = webhook.NewWebhookService
	NewWebhookHTTPClient = webhook.NewWebhookHTTPClient
)

// Meow types
type MeowServiceImpl = service.MeowServiceImpl
type MeowClient = adapter.MeowClient

// Meow constructors
var (
	NewMeowService = service.NewMeowService
)

// Meow constants
const (
	// Messages and Communication
	EventMessage                    = core.EventMessage
	EventUndecryptableMessage      = core.EventUndecryptableMessage
	EventReceipt                   = core.EventReceipt
	EventMediaRetry                = core.EventMediaRetry
	EventReadReceipt               = core.EventReadReceipt

	// Groups and Contacts
	EventGroupInfo                 = core.EventGroupInfo
	EventJoinedGroup               = core.EventJoinedGroup
	EventPicture                   = core.EventPicture
	EventBlocklistChange           = core.EventBlocklistChange
	EventBlocklist                 = core.EventBlocklist

	// Connection and Session
	EventConnected                 = core.EventConnected
	EventDisconnected              = core.EventDisconnected
	EventConnectFailure            = core.EventConnectFailure
	EventKeepAliveRestored         = core.EventKeepAliveRestored
	EventKeepAliveTimeout          = core.EventKeepAliveTimeout
	EventLoggedOut                 = core.EventLoggedOut
	EventClientOutdated            = core.EventClientOutdated
	EventTemporaryBan              = core.EventTemporaryBan
	EventStreamError               = core.EventStreamError
	EventStreamReplaced            = core.EventStreamReplaced
	EventPairSuccess               = core.EventPairSuccess
	EventPairError                 = core.EventPairError
	EventQR                        = core.EventQR
	EventQRScannedWithoutMultidevice = core.EventQRScannedWithoutMultidevice

	// Privacy and Settings
	EventPrivacySettings           = core.EventPrivacySettings
	EventPushNameSetting           = core.EventPushNameSetting
	EventUserAbout                 = core.EventUserAbout

	// Synchronization and State
	EventAppState                  = core.EventAppState
	EventAppStateSyncComplete      = core.EventAppStateSyncComplete
	EventHistorySync               = core.EventHistorySync
	EventOfflineSyncCompleted      = core.EventOfflineSyncCompleted
	EventOfflineSyncPreview        = core.EventOfflineSyncPreview

	// Calls
	EventCallOffer                 = core.EventCallOffer
	EventCallAccept                = core.EventCallAccept
	EventCallTerminate             = core.EventCallTerminate
	EventCallOfferNotice           = core.EventCallOfferNotice
	EventCallRelayLatency          = core.EventCallRelayLatency

	// Presence and Activity
	EventPresence                  = core.EventPresence
	EventChatPresence              = core.EventChatPresence

	// Identity
	EventIdentityChange            = core.EventIdentityChange

	// Errors
	EventCATRefreshError           = core.EventCATRefreshError

	// Newsletter (WhatsApp Channels)
	EventNewsletterJoin            = core.EventNewsletterJoin
	EventNewsletterLeave           = core.EventNewsletterLeave
	EventNewsletterMuteChange      = core.EventNewsletterMuteChange
	EventNewsletterLiveUpdate      = core.EventNewsletterLiveUpdate

	// Facebook/Meta Bridge
	EventFBMessage                 = core.EventFBMessage

	// Special - receives all events
	EventAll                       = core.EventAll
)
