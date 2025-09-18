package dto

import (
	"encoding/base64"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)


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

type GetNewsletterInfoRequest struct {
	JID string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
}

type GetNewsletterInfoWithInviteRequest struct {
	InviteKey string `json:"invite_key" binding:"required" example:"invite_key_123"`
}


type CreateNewsletterRequest struct {
	Name        string `json:"name" binding:"required" example:"My Newsletter"`
	Description string `json:"description,omitempty" example:"Newsletter description"`
	Picture     string `json:"picture,omitempty"` // Base64 encoded image
}

type CreateNewsletterResponse struct {
	Success bool                    `json:"success"`
	Data    *CreateNewsletterResult `json:"data,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

type CreateNewsletterResult struct {
	JID         string    `json:"jid"`
	ServerID    string    `json:"server_id"`
	Timestamp   time.Time `json:"timestamp"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}


type GetNewsletterMessagesRequest struct {
	JID    string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	Count  int    `json:"count,omitempty" example:"50"`
	Before string `json:"before,omitempty" example:"server_message_id_123"`
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

type GetNewsletterMessageUpdatesRequest struct {
	JID    string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	Count  int    `json:"count,omitempty" example:"50"`
	Before string `json:"before,omitempty" example:"server_message_id_123"`
}

func (r *GetNewsletterMessageUpdatesRequest) ToWhatsmeowParams() (*whatsmeow.GetNewsletterMessagesParams, error) {
	params := &whatsmeow.GetNewsletterMessagesParams{
		Count: r.Count,
	}

	return params, nil
}


type FollowNewsletterRequest struct {
	JID string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
}

type UnfollowNewsletterRequest struct {
	JID string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
}


type NewsletterReactionRequest struct {
	JID       string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	ServerID  string `json:"server_id" binding:"required" example:"server_msg_123"`
	MessageID string `json:"message_id" binding:"required" example:"msg_123"`
	Reaction  string `json:"reaction" binding:"required" example:"👍"`
}

type SendReactionRequest struct {
	ServerID  string `json:"server_id" binding:"required" example:"server_msg_123"`
	MessageID string `json:"message_id" binding:"required" example:"msg_123"`
	Reaction  string `json:"reaction" binding:"required" example:"👍"`
}

type MarkViewedRequest struct {
	JID       string   `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	ServerIDs []string `json:"server_ids" binding:"required" example:"[\"server_msg_1\", \"server_msg_2\"]"`
}

type ToggleMuteRequest struct {
	JID  string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	Mute bool   `json:"mute" example:"true"`
}

type SendNewsletterMessageRequest struct {
	Content     string `json:"content,omitempty" example:"Hello newsletter!"`
	MediaHandle string `json:"media_handle,omitempty" example:"media_handle_from_upload"`
	MediaType   string `json:"media_type,omitempty" example:"image"`
	Caption     string `json:"caption,omitempty" example:"Image caption"`
}


type UploadNewsletterMediaRequest struct {
	Data      string `json:"data" binding:"required" example:"base64_encoded_data"` // Base64 encoded media
	MediaType string `json:"media_type" binding:"required" example:"image"`         // Media type: image, video, audio, document
	Filename  string `json:"filename,omitempty" example:"image.jpg"`                // Optional filename
}

type UploadNewsletterMediaResponse struct {
	Success bool                         `json:"success"`
	Data    *UploadNewsletterMediaResult `json:"data,omitempty"`
	Error   string                       `json:"error,omitempty"`
}

type UploadNewsletterMediaResult struct {
	URL        string `json:"url"`
	DirectPath string `json:"direct_path"`
	Handle     string `json:"handle"`
	ObjectID   string `json:"object_id"`
	MediaKey   string `json:"media_key"`
	FileHash   string `json:"file_hash"`
	FileLength int64  `json:"file_length"`
}


type NewsletterInfoResponse struct {
	Success bool            `json:"success"`
	Data    *NewsletterInfo `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type NewsletterListResponse struct {
	Success bool            `json:"success"`
	Data    *NewsletterList `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type NewsletterList struct {
	Newsletters []NewsletterInfo `json:"newsletters"`
	Count       int              `json:"count"`
	Total       int              `json:"total,omitempty"`
}

type NewsletterMessagesResponse struct {
	Success bool                      `json:"success"`
	Data    *NewsletterMessagesResult `json:"data,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

type NewsletterMessagesResult struct {
	Messages []NewsletterMessage `json:"messages"`
	Count    int                 `json:"count"`
	HasMore  bool                `json:"has_more"`
}

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type SendNewsletterMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	ServerID  string `json:"server_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Error     string `json:"error,omitempty"`
}


func (r *CreateNewsletterRequest) ToWhatsmeowParams() (*whatsmeow.CreateNewsletterParams, error) {
	params := &whatsmeow.CreateNewsletterParams{
		Name:        r.Name,
		Description: r.Description,
	}

	if r.Picture != "" {
		pictureBytes, err := base64.StdEncoding.DecodeString(r.Picture)
		if err != nil {
			return nil, fmt.Errorf("invalid picture format: %w", err)
		}
		params.Picture = pictureBytes
	}

	return params, nil
}

func (r *GetNewsletterInfoRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

func (r *FollowNewsletterRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

func (r *UnfollowNewsletterRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

func (r *GetNewsletterMessagesRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

func (r *GetNewsletterMessagesRequest) ToWhatsmeowParams() (*whatsmeow.GetNewsletterMessagesParams, error) {
	params := &whatsmeow.GetNewsletterMessagesParams{
		Count: r.Count,
	}

	return params, nil
}

func ConvertFromWhatsmeowNewsletterInfo(info *types.NewsletterMetadata) *NewsletterInfo {
	if info == nil {
		return nil
	}

	return &NewsletterInfo{
		JID:         info.ID.String(),
		Name:        info.ThreadMeta.Name.Text,
		Description: info.ThreadMeta.Description.Text,
		Picture:     "", // Picture extraction not available in current whatsmeow version
		Verified:    info.ThreadMeta.VerificationState == types.NewsletterVerificationStateVerified,
		CreatedAt:   time.Now(), // Using current time - jsontime conversion complex
		UpdatedAt:   time.Now(), // Using current time - jsontime conversion complex
		OwnerJID:    "",         // Not available in NewsletterMetadata
		Subscribers: info.ThreadMeta.SubscriberCount,
		Muted:       info.ViewerMeta != nil && info.ViewerMeta.Mute == types.NewsletterMuteOn,
	}
}

func ConvertFromWhatsmeowNewsletterInfoList(infos []*types.NewsletterMetadata) []NewsletterInfo {
	if infos == nil {
		return []NewsletterInfo{}
	}

	result := make([]NewsletterInfo, len(infos))
	for i, info := range infos {
		if converted := ConvertFromWhatsmeowNewsletterInfo(info); converted != nil {
			result[i] = *converted
		}
	}
	return result
}
