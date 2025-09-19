package dto

import "time"

type MessageResult struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type GroupResult struct {
	GroupJID string   `json:"groupJid"`
	Name     string   `json:"name"`
	Members  []string `json:"members"`
}

type EventData struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

type MediaUploadResult struct {
	URL      string `json:"url"`
	MediaID  string `json:"mediaId"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}
