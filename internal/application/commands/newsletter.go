package commands

import (
	"fmt"
	"strings"
)

// CreateNewsletterCommand represents a command to create a newsletter
type CreateNewsletterCommand struct {
	SessionID   string `json:"session_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Picture     string `json:"picture,omitempty"` // base64 or URL
}

// SendNewsletterMessageCommand represents a command to send a newsletter message
type SendNewsletterMessageCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
	Content       string `json:"content,omitempty"`
	MediaHandle   string `json:"media_handle,omitempty"`
	MediaType     string `json:"media_type,omitempty"`
	Caption       string `json:"caption,omitempty"`
}

// SubscribeNewsletterCommand represents a command to subscribe to a newsletter
type SubscribeNewsletterCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
}

// UnsubscribeNewsletterCommand represents a command to unsubscribe from a newsletter
type UnsubscribeNewsletterCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
}

// GetNewsletterMessagesCommand represents a command to get newsletter messages
type GetNewsletterMessagesCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
	Count         int    `json:"count,omitempty"`
	Before        string `json:"before,omitempty"`
}

// GetNewsletterMessageUpdatesCommand represents a command to get newsletter message updates
type GetNewsletterMessageUpdatesCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
	Count         int    `json:"count,omitempty"`
	Before        string `json:"before,omitempty"`
}

// MarkNewsletterMessagesViewedCommand represents a command to mark newsletter messages as viewed
type MarkNewsletterMessagesViewedCommand struct {
	SessionID     string   `json:"session_id" validate:"required"`
	NewsletterJID string   `json:"newsletter_jid" validate:"required"`
	ServerIDs     []string `json:"server_ids" validate:"required"`
}

// SendNewsletterReactionCommand represents a command to send a reaction to a newsletter message
type SendNewsletterReactionCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
	MessageID     string `json:"message_id" validate:"required"`
	ServerID      string `json:"server_id" validate:"required"`
	Reaction      string `json:"reaction" validate:"required"`
}

// ToggleNewsletterMuteCommand represents a command to toggle newsletter mute status
type ToggleNewsletterMuteCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
	Mute          bool   `json:"mute"`
}

// UploadNewsletterMediaCommand represents a command to upload media for newsletter
type UploadNewsletterMediaCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	MediaType string `json:"media_type" validate:"required"`
	MediaData []byte `json:"media_data" validate:"required"`
}

// GetNewsletterInfoCommand represents a command to get newsletter information
type GetNewsletterInfoCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
}

// GetNewsletterInfoWithInviteCommand represents a command to get newsletter info with invite
type GetNewsletterInfoWithInviteCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	InviteKey string `json:"invite_key" validate:"required"`
}

// Validation methods

func (c *CreateNewsletterCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("newsletter name is required")
	}
	return nil
}

func (c *SendNewsletterMessageCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}

	// Either content or media handle is required
	if strings.TrimSpace(c.Content) == "" && c.MediaHandle == "" {
		return fmt.Errorf("either content or media handle is required")
	}

	// If media handle is provided, media type is required
	if c.MediaHandle != "" && c.MediaType == "" {
		return fmt.Errorf("media type is required when media handle is provided")
	}

	return nil
}

func (c *SubscribeNewsletterCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

func (c *UnsubscribeNewsletterCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

func (c *GetNewsletterMessagesCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	if c.Count < 0 {
		return fmt.Errorf("count cannot be negative")
	}
	return nil
}

func (c *GetNewsletterMessageUpdatesCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	if c.Count < 0 {
		return fmt.Errorf("count cannot be negative")
	}
	return nil
}

func (c *MarkNewsletterMessagesViewedCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	if len(c.ServerIDs) == 0 {
		return fmt.Errorf("at least one server ID is required")
	}
	return nil
}

func (c *SendNewsletterReactionCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	if c.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	if c.ServerID == "" {
		return fmt.Errorf("server ID is required")
	}
	if strings.TrimSpace(c.Reaction) == "" {
		return fmt.Errorf("reaction is required")
	}
	return nil
}

func (c *ToggleNewsletterMuteCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

func (c *UploadNewsletterMediaCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.MediaType == "" {
		return fmt.Errorf("media type is required")
	}
	if len(c.MediaData) == 0 {
		return fmt.Errorf("media data is required")
	}
	return nil
}

func (c *GetNewsletterInfoCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

// GetNewsletterCommand represents a command to get newsletter information
type GetNewsletterCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
}

func (c *GetNewsletterCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

// SubscribeToNewsletterCommand represents a command to subscribe to a newsletter
type SubscribeToNewsletterCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
}

func (c *SubscribeToNewsletterCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

// UnsubscribeFromNewsletterCommand represents a command to unsubscribe from a newsletter
type UnsubscribeFromNewsletterCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
}

func (c *UnsubscribeFromNewsletterCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

// MarkNewsletterViewedCommand represents a command to mark newsletter as viewed
type MarkNewsletterViewedCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
	MessageID     string `json:"message_id" validate:"required"`
}

func (c *MarkNewsletterViewedCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	if c.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	return nil
}

// SubscribeLiveUpdatesCommand represents a command to subscribe to live updates
type SubscribeLiveUpdatesCommand struct {
	SessionID     string `json:"session_id" validate:"required"`
	NewsletterJID string `json:"newsletter_jid" validate:"required"`
}

func (c *SubscribeLiveUpdatesCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.NewsletterJID == "" {
		return fmt.Errorf("newsletter JID is required")
	}
	return nil
}

// GetNewsletterByInviteCommand represents a command to get newsletter by invite
type GetNewsletterByInviteCommand struct {
	SessionID string `json:"session_id" validate:"required"`
	InviteKey string `json:"invite_key" validate:"required"`
}

func (c *GetNewsletterByInviteCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.InviteKey == "" {
		return fmt.Errorf("invite key is required")
	}
	return nil
}

func (c *GetNewsletterInfoWithInviteCommand) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if c.InviteKey == "" {
		return fmt.Errorf("invite key is required")
	}
	return nil
}
