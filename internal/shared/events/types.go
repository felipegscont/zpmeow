package events

type EventType string

const (
	EventTypeMessage              EventType = "Message"
	EventTypeUndecryptableMessage EventType = "UndecryptableMessage"
	EventTypeReceipt              EventType = "Receipt"
	EventTypeMediaRetry           EventType = "MediaRetry"
	EventTypeMediaRetryError      EventType = "MediaRetryError"
	EventTypeReadReceipt          EventType = "ReadReceipt"

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

	EventTypePairSuccess                 EventType = "PairSuccess"
	EventTypePairError                   EventType = "PairError"
	EventTypeQR                          EventType = "QR"
	EventTypeQRScannedWithoutMultidevice EventType = "QRScannedWithoutMultidevice"
	EventTypeManualLoginReconnect        EventType = "ManualLoginReconnect"

	EventTypeGroupInfo       EventType = "GroupInfo"
	EventTypeJoinedGroup     EventType = "JoinedGroup"
	EventTypePicture         EventType = "Picture"
	EventTypeBlocklistChange EventType = "BlocklistChange"
	EventTypeBlocklist       EventType = "Blocklist"
	EventTypeContact         EventType = "Contact"

	EventTypePrivacySettings EventType = "PrivacySettings"
	EventTypePushNameSetting EventType = "PushNameSetting"
	EventTypePushName        EventType = "PushName"
	EventTypeUserAbout       EventType = "UserAbout"
	EventTypeBusinessName    EventType = "BusinessName"

	EventTypePresence     EventType = "Presence"
	EventTypeChatPresence EventType = "ChatPresence"

	EventTypeAll EventType = "All" // Receives all events
)

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

func (et EventType) String() string {
	return string(et)
}

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

func GetEventTypeNames() []string {
	eventTypes := GetAllEventTypes()
	names := make([]string, len(eventTypes))
	for i, eventType := range eventTypes {
		names[i] = string(eventType)
	}
	return names
}
