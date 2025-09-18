package application

// WhatsApp service types moved to interfaces.go to avoid duplication
// This file contains only the DTOs and response types

type MessageResponse struct {
	ID        string
	Timestamp int64
	Status    string
}

type GroupResponse struct {
	GroupJID string
	Name     string
	Members  []string
}

type UserInfo struct {
	JID           string
	Name          string
	ProfilePicURL string
	Status        string
	IsBlocked     bool
}

type ChatInfo struct {
	JID           string
	Name          string
	IsGroup       bool
	LastMessage   string
	LastTimestamp int64
	UnreadCount   int
	IsMuted       bool
	IsPinned      bool
	IsArchived    bool
}

type MessageInfo struct {
	ID        string
	FromJID   string
	ToJID     string
	Content   string
	Type      string
	Timestamp int64
	IsFromMe  bool
}
