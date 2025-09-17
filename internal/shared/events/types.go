package events

// EventType represents a WhatsApp event type (shared across layers)
type EventType string

const (
	// Message events
	EventTypeMessage              EventType = "Message"
	EventTypeUndecryptableMessage EventType = "UndecryptableMessage"
	EventTypeReceipt              EventType = "Receipt"
	EventTypeMediaRetry           EventType = "MediaRetry"
	EventTypeMediaRetryError      EventType = "MediaRetryError"
	EventTypeReadReceipt          EventType = "ReadReceipt"

	// Connection events
	EventTypeConnected         EventType = "Connected"
	EventTypeDisconnected      EventType = "Disconnected"
	EventTypeConnectFailure    EventType = "ConnectFailure"
	EventTypeKeepAliveRestored EventType = "KeepAliveRestored"
	EventTypeKeepAliveTimeout  EventType = "KeepAliveTimeout"
	EventTypeLoggedOut         EventType = "LoggedOut"
	EventTypeClientOutdated    EventType = "ClientOutdated"
	EventTypeTemporaryBan      EventType = "TemporaryBan"
	EventTypeStreamError       EventType = "StreamError"
	EventTypeStreamReplaced    EventType = "StreamReplaced"

	// Authentication events
	EventTypePairSuccess                 EventType = "PairSuccess"
	EventTypePairError                   EventType = "PairError"
	EventTypeQR                          EventType = "QR"
	EventTypeQRScannedWithoutMultidevice EventType = "QRScannedWithoutMultidevice"
	EventTypeManualLoginReconnect        EventType = "ManualLoginReconnect"

	// Groups and Contacts
	EventTypeGroupInfo       EventType = "GroupInfo"
	EventTypeJoinedGroup     EventType = "JoinedGroup"
	EventTypePicture         EventType = "Picture"
	EventTypeBlocklistChange EventType = "BlocklistChange"
	EventTypeBlocklist       EventType = "Blocklist"
	EventTypeContact         EventType = "Contact"

	// Privacy and Settings
	EventTypePrivacySettings EventType = "PrivacySettings"
	EventTypePushNameSetting EventType = "PushNameSetting"
	EventTypePushName        EventType = "PushName"
	EventTypeUserAbout       EventType = "UserAbout"
	EventTypeBusinessName    EventType = "BusinessName"

	// Presence
	EventTypePresence     EventType = "Presence"
	EventTypeChatPresence EventType = "ChatPresence"

	// Special
	EventTypeAll EventType = "All" // Receives all events
)

// IsValid checks if the event type is valid
func (et EventType) IsValid() bool {
	switch et {
	case EventTypeMessage, EventTypeUndecryptableMessage, EventTypeReceipt,
		EventTypeMediaRetry, EventTypeMediaRetryError, EventTypeReadReceipt,
		EventTypeConnected, EventTypeDisconnected, EventTypeConnectFailure,
		EventTypeKeepAliveRestored, EventTypeKeepAliveTimeout, EventTypeLoggedOut,
		EventTypeClientOutdated, EventTypeTemporaryBan, EventTypeStreamError,
		EventTypeStreamReplaced, EventTypePairSuccess, EventTypePairError,
		EventTypeQR, EventTypeQRScannedWithoutMultidevice, EventTypeManualLoginReconnect,
		EventTypeGroupInfo, EventTypeJoinedGroup, EventTypePicture,
		EventTypeBlocklistChange, EventTypeBlocklist, EventTypeContact,
		EventTypePrivacySettings, EventTypePushNameSetting, EventTypePushName,
		EventTypeUserAbout, EventTypeBusinessName, EventTypePresence,
		EventTypeChatPresence, EventTypeAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the event type
func (et EventType) String() string {
	return string(et)
}

// GetAllEventTypes returns all valid event types
func GetAllEventTypes() []EventType {
	return []EventType{
		EventTypeMessage, EventTypeUndecryptableMessage, EventTypeReceipt,
		EventTypeMediaRetry, EventTypeMediaRetryError, EventTypeReadReceipt,
		EventTypeConnected, EventTypeDisconnected, EventTypeConnectFailure,
		EventTypeKeepAliveRestored, EventTypeKeepAliveTimeout, EventTypeLoggedOut,
		EventTypeClientOutdated, EventTypeTemporaryBan, EventTypeStreamError,
		EventTypeStreamReplaced, EventTypePairSuccess, EventTypePairError,
		EventTypeQR, EventTypeQRScannedWithoutMultidevice, EventTypeManualLoginReconnect,
		EventTypeGroupInfo, EventTypeJoinedGroup, EventTypePicture,
		EventTypeBlocklistChange, EventTypeBlocklist, EventTypeContact,
		EventTypePrivacySettings, EventTypePushNameSetting, EventTypePushName,
		EventTypeUserAbout, EventTypeBusinessName, EventTypePresence,
		EventTypeChatPresence, EventTypeAll,
	}
}

// GetEventTypeNames returns all valid event type names as strings
func GetEventTypeNames() []string {
	eventTypes := GetAllEventTypes()
	names := make([]string, len(eventTypes))
	for i, eventType := range eventTypes {
		names[i] = string(eventType)
	}
	return names
}
