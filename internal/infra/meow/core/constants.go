package core

// WhatsApp Event Types
const (
	// Messages and Communication
	EventMessage              = "Message"
	EventUndecryptableMessage = "UndecryptableMessage"
	EventReceipt              = "Receipt"
	EventMediaRetry           = "MediaRetry"
	EventReadReceipt          = "ReadReceipt"

	// Groups and Contacts
	EventGroupInfo       = "GroupInfo"
	EventJoinedGroup     = "JoinedGroup"
	EventPicture         = "Picture"
	EventBlocklistChange = "BlocklistChange"
	EventBlocklist       = "Blocklist"

	// Connection and Session
	EventConnected                   = "Connected"
	EventDisconnected                = "Disconnected"
	EventConnectFailure              = "ConnectFailure"
	EventKeepAliveRestored           = "KeepAliveRestored"
	EventKeepAliveTimeout            = "KeepAliveTimeout"
	EventLoggedOut                   = "LoggedOut"
	EventClientOutdated              = "ClientOutdated"
	EventTemporaryBan                = "TemporaryBan"
	EventStreamError                 = "StreamError"
	EventStreamReplaced              = "StreamReplaced"
	EventPairSuccess                 = "PairSuccess"
	EventPairError                   = "PairError"
	EventQR                          = "QR"
	EventQRScannedWithoutMultidevice = "QRScannedWithoutMultidevice"

	// Privacy and Settings
	EventPrivacySettings = "PrivacySettings"
	EventPushNameSetting = "PushNameSetting"
	EventUserAbout       = "UserAbout"

	// Synchronization and State
	EventAppState             = "AppState"
	EventAppStateSyncComplete = "AppStateSyncComplete"
	EventHistorySync          = "HistorySync"
	EventOfflineSyncCompleted = "OfflineSyncCompleted"
	EventOfflineSyncPreview   = "OfflineSyncPreview"

	// Calls
	EventCallOffer        = "CallOffer"
	EventCallAccept       = "CallAccept"
	EventCallTerminate    = "CallTerminate"
	EventCallOfferNotice  = "CallOfferNotice"
	EventCallRelayLatency = "CallRelayLatency"

	// Presence and Activity
	EventPresence     = "Presence"
	EventChatPresence = "ChatPresence"

	// Identity
	EventIdentityChange = "IdentityChange"

	// Errors
	EventCATRefreshError = "CATRefreshError"

	// Newsletter (WhatsApp Channels)
	EventNewsletterJoin       = "NewsletterJoin"
	EventNewsletterLeave      = "NewsletterLeave"
	EventNewsletterMuteChange = "NewsletterMuteChange"
	EventNewsletterLiveUpdate = "NewsletterLiveUpdate"

	// Facebook/Meta Bridge
	EventFBMessage = "FBMessage"

	// Special - receives all events
	EventAll = "All"
)
