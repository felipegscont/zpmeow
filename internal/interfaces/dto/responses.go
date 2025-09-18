package dto

import "time"

// Legacy response types - consider migrating to common.go structures

// MessageResult represents the result of a message operation
type MessageResult struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// GroupResult represents the result of a group operation
type GroupResult struct {
	GroupJID string   `json:"groupJid"`
	Name     string   `json:"name"`
	Members  []string `json:"members"`
}

// EventData represents webhook event data
type EventData struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// MediaUploadResult represents the result of a media upload operation
type MediaUploadResult struct {
	URL      string `json:"url"`
	MediaID  string `json:"mediaId"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}
