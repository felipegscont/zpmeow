package dto

import (
	"encoding/base64"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// ============================================================================
// NEWSLETTER INFO DTOs
// ============================================================================

// NewsletterInfo represents newsletter information
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

// GetNewsletterInfoRequest represents the request to get newsletter info
type GetNewsletterInfoRequest struct {
	JID string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
}

// GetNewsletterInfoWithInviteRequest represents the request to get newsletter info via invite
type GetNewsletterInfoWithInviteRequest struct {
	InviteKey string `json:"invite_key" binding:"required" example:"invite_key_123"`
}

// ============================================================================
// CREATE NEWSLETTER DTOs
// ============================================================================

// CreateNewsletterRequest represents the request to create a newsletter
type CreateNewsletterRequest struct {
	Name        string `json:"name" binding:"required" example:"My Newsletter"`
	Description string `json:"description,omitempty" example:"Newsletter description"`
	Picture     string `json:"picture,omitempty"` // Base64 encoded image
}

// CreateNewsletterResponse represents the response for creating a newsletter
type CreateNewsletterResponse struct {
	Success bool                    `json:"success"`
	Data    *CreateNewsletterResult `json:"data,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

// CreateNewsletterResult represents the result of creating a newsletter
type CreateNewsletterResult struct {
	JID         string    `json:"jid"`
	ServerID    string    `json:"server_id"`
	Timestamp   time.Time `json:"timestamp"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

// ============================================================================
// NEWSLETTER MESSAGES DTOs
// ============================================================================

// GetNewsletterMessagesRequest represents the request to get newsletter messages
type GetNewsletterMessagesRequest struct {
	JID    string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	Count  int    `json:"count,omitempty" example:"50"`
	Before string `json:"before,omitempty" example:"server_message_id_123"`
}

// NewsletterMessage represents a newsletter message
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

// GetNewsletterMessageUpdatesRequest represents the request to get message updates
type GetNewsletterMessageUpdatesRequest struct {
	JID    string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	Count  int    `json:"count,omitempty" example:"50"`
	Before string `json:"before,omitempty" example:"server_message_id_123"`
}

// ToWhatsmeowParams converts GetNewsletterMessageUpdatesRequest to whatsmeow params
func (r *GetNewsletterMessageUpdatesRequest) ToWhatsmeowParams() (*whatsmeow.GetNewsletterMessagesParams, error) {
	params := &whatsmeow.GetNewsletterMessagesParams{
		Count: r.Count,
	}

	return params, nil
}

// ============================================================================
// FOLLOW/UNFOLLOW DTOs
// ============================================================================

// FollowNewsletterRequest represents the request to follow a newsletter
type FollowNewsletterRequest struct {
	JID string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
}

// UnfollowNewsletterRequest represents the request to unfollow a newsletter
type UnfollowNewsletterRequest struct {
	JID string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
}

// ============================================================================
// NEWSLETTER INTERACTIONS DTOs
// ============================================================================

// NewsletterReactionRequest represents the request to react to a newsletter message
type NewsletterReactionRequest struct {
	JID       string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	ServerID  string `json:"server_id" binding:"required" example:"server_msg_123"`
	MessageID string `json:"message_id" binding:"required" example:"msg_123"`
	Reaction  string `json:"reaction" binding:"required" example:"👍"`
}

// SendReactionRequest represents the request to send a reaction (simplified version)
type SendReactionRequest struct {
	ServerID  string `json:"server_id" binding:"required" example:"server_msg_123"`
	MessageID string `json:"message_id" binding:"required" example:"msg_123"`
	Reaction  string `json:"reaction" binding:"required" example:"👍"`
}

// MarkViewedRequest represents the request to mark newsletter messages as viewed
type MarkViewedRequest struct {
	JID       string   `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	ServerIDs []string `json:"server_ids" binding:"required" example:"[\"server_msg_1\", \"server_msg_2\"]"`
}

// ToggleMuteRequest represents the request to mute/unmute a newsletter
type ToggleMuteRequest struct {
	JID  string `json:"jid" binding:"required" example:"120363123456789012@newsletter"`
	Mute bool   `json:"mute" example:"true"`
}

// SendNewsletterMessageRequest represents the request to send a message to a newsletter
type SendNewsletterMessageRequest struct {
	Content     string `json:"content,omitempty" example:"Hello newsletter!"`
	MediaHandle string `json:"media_handle,omitempty" example:"media_handle_from_upload"`
	MediaType   string `json:"media_type,omitempty" example:"image"`
	Caption     string `json:"caption,omitempty" example:"Image caption"`
}

// ============================================================================
// NEWSLETTER UPLOAD DTOs
// ============================================================================

// UploadNewsletterMediaRequest represents the request to upload media to newsletter
type UploadNewsletterMediaRequest struct {
	Data      string `json:"data" binding:"required" example:"base64_encoded_data"` // Base64 encoded media
	MediaType string `json:"media_type" binding:"required" example:"image"`         // Media type: image, video, audio, document
	Filename  string `json:"filename,omitempty" example:"image.jpg"`                // Optional filename
}

// UploadNewsletterMediaResponse represents the response for media upload
type UploadNewsletterMediaResponse struct {
	Success bool                         `json:"success"`
	Data    *UploadNewsletterMediaResult `json:"data,omitempty"`
	Error   string                       `json:"error,omitempty"`
}

// UploadNewsletterMediaResult represents the result of media upload
type UploadNewsletterMediaResult struct {
	URL        string `json:"url"`
	DirectPath string `json:"direct_path"`
	Handle     string `json:"handle"`
	ObjectID   string `json:"object_id"`
	MediaKey   string `json:"media_key"`
	FileHash   string `json:"file_hash"`
	FileLength int64  `json:"file_length"`
}

// ============================================================================
// RESPONSE DTOs
// ============================================================================

// NewsletterInfoResponse represents the response for newsletter info requests
type NewsletterInfoResponse struct {
	Success bool            `json:"success"`
	Data    *NewsletterInfo `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// NewsletterListResponse represents the response for listing newsletters
type NewsletterListResponse struct {
	Success bool            `json:"success"`
	Data    *NewsletterList `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// NewsletterList represents a list of newsletters with metadata
type NewsletterList struct {
	Newsletters []NewsletterInfo `json:"newsletters"`
	Count       int              `json:"count"`
	Total       int              `json:"total,omitempty"`
}

// NewsletterMessagesResponse represents the response for newsletter messages
type NewsletterMessagesResponse struct {
	Success bool                      `json:"success"`
	Data    *NewsletterMessagesResult `json:"data,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

// NewsletterMessagesResult represents the result of getting newsletter messages
type NewsletterMessagesResult struct {
	Messages []NewsletterMessage `json:"messages"`
	Count    int                 `json:"count"`
	HasMore  bool                `json:"has_more"`
}

// StandardResponse represents a standard success/error response
type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SendNewsletterMessageResponse represents the response for sending a newsletter message
type SendNewsletterMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	ServerID  string `json:"server_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ============================================================================
// CONVERSION FUNCTIONS
// ============================================================================

// ToWhatsmeowParams converts CreateNewsletterRequest to whatsmeow CreateNewsletterParams
func (r *CreateNewsletterRequest) ToWhatsmeowParams() (*whatsmeow.CreateNewsletterParams, error) {
	params := &whatsmeow.CreateNewsletterParams{
		Name:        r.Name,
		Description: r.Description,
	}

	// Convert base64 picture to bytes if provided
	if r.Picture != "" {
		pictureBytes, err := base64.StdEncoding.DecodeString(r.Picture)
		if err != nil {
			return nil, fmt.Errorf("invalid picture format: %w", err)
		}
		params.Picture = pictureBytes
	}

	return params, nil
}

// ToJID converts string JID to types.JID
func (r *GetNewsletterInfoRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

// ToJID converts string JID to types.JID for follow request
func (r *FollowNewsletterRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

// ToJID converts string JID to types.JID for unfollow request
func (r *UnfollowNewsletterRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

// ToJID converts string JID to types.JID for messages request
func (r *GetNewsletterMessagesRequest) ToJID() (types.JID, error) {
	return types.ParseJID(r.JID)
}

// ToWhatsmeowParams converts GetNewsletterMessagesRequest to whatsmeow params
func (r *GetNewsletterMessagesRequest) ToWhatsmeowParams() (*whatsmeow.GetNewsletterMessagesParams, error) {
	params := &whatsmeow.GetNewsletterMessagesParams{
		Count: r.Count,
	}

	return params, nil
}

// ConvertFromWhatsmeowNewsletterInfo converts whatsmeow types to our DTOs
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

// ConvertFromWhatsmeowNewsletterInfoList converts a slice of whatsmeow newsletter info
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
