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

// Types moved to interfaces.go to avoid duplication

type MessageInfo struct {
	ID        string
	FromJID   string
	ToJID     string
	Content   string
	Type      string
	Timestamp int64
	IsFromMe  bool
}
