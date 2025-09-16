package application

import (
	"context"
	"fmt"

	"zpmeow/internal/application/commands"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/shared/validation"
)

// NewsletterService implements newsletter use cases following Clean Architecture
type NewsletterService struct {
	wameowService interface {
		CreateNewsletter(ctx context.Context, sessionID, name, description string) (interface{}, error)
		GetNewsletter(ctx context.Context, sessionID, newsletterJID string) (interface{}, error)
		ListNewsletters(ctx context.Context, sessionID string) (interface{}, error)
		SubscribeToNewsletter(ctx context.Context, sessionID, newsletterJID string) error
		UnsubscribeFromNewsletter(ctx context.Context, sessionID, newsletterJID string) error
		SendNewsletterMessage(ctx context.Context, sessionID, newsletterJID, content string) (interface{}, error)
		GetNewsletterMessages(ctx context.Context, sessionID, newsletterJID string, limit int) (interface{}, error)
		GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterJID string) (interface{}, error)
		MarkNewsletterViewed(ctx context.Context, sessionID, newsletterJID, messageID string) error
		SendNewsletterReaction(ctx context.Context, sessionID, newsletterJID, messageID, reaction string) error
		ToggleNewsletterMute(ctx context.Context, sessionID, newsletterJID string, mute bool) error
		SubscribeLiveUpdates(ctx context.Context, sessionID, newsletterJID string) error
		UploadNewsletterMedia(ctx context.Context, sessionID string, mediaData []byte, mediaType string) (interface{}, error)
		GetNewsletterByInvite(ctx context.Context, sessionID, inviteKey string) (interface{}, error)
	}
	sessionRepo session.Repository
	validator   *validation.Validator
}

// NewNewsletterService creates a new NewsletterService instance
func NewNewsletterService(
	wameowService interface{},
	sessionRepo session.Repository,
	validator *validation.Validator,
) *NewsletterService {
	return &NewsletterService{
		wameowService: wameowService.(interface {
			CreateNewsletter(ctx context.Context, sessionID, name, description string) (interface{}, error)
			GetNewsletter(ctx context.Context, sessionID, newsletterJID string) (interface{}, error)
			ListNewsletters(ctx context.Context, sessionID string) (interface{}, error)
			SubscribeToNewsletter(ctx context.Context, sessionID, newsletterJID string) error
			UnsubscribeFromNewsletter(ctx context.Context, sessionID, newsletterJID string) error
			SendNewsletterMessage(ctx context.Context, sessionID, newsletterJID, content string) (interface{}, error)
			GetNewsletterMessages(ctx context.Context, sessionID, newsletterJID string, limit int) (interface{}, error)
			GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterJID string) (interface{}, error)
			MarkNewsletterViewed(ctx context.Context, sessionID, newsletterJID, messageID string) error
			SendNewsletterReaction(ctx context.Context, sessionID, newsletterJID, messageID, reaction string) error
			ToggleNewsletterMute(ctx context.Context, sessionID, newsletterJID string, mute bool) error
			SubscribeLiveUpdates(ctx context.Context, sessionID, newsletterJID string) error
			UploadNewsletterMedia(ctx context.Context, sessionID string, mediaData []byte, mediaType string) (interface{}, error)
			GetNewsletterByInvite(ctx context.Context, sessionID, inviteKey string) (interface{}, error)
		}),
		sessionRepo: sessionRepo,
		validator:   validator,
	}
}

// CreateNewsletter creates a new newsletter using command pattern
func (s *NewsletterService) CreateNewsletter(ctx context.Context, cmd *commands.CreateNewsletterCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.CreateNewsletter(ctx, cmd.SessionID, cmd.Name, cmd.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create newsletter: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// GetNewsletter gets newsletter information using command pattern
func (s *NewsletterService) GetNewsletter(ctx context.Context, cmd *commands.GetNewsletterCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.GetNewsletter(ctx, cmd.SessionID, cmd.NewsletterJID)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter: %w", err)
	}

	return result, nil
}

// ListNewsletters lists all newsletters
func (s *NewsletterService) ListNewsletters(ctx context.Context, sessionID string) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.ListNewsletters(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list newsletters: %w", err)
	}

	return result, nil
}

// SubscribeToNewsletter subscribes to a newsletter using command pattern
func (s *NewsletterService) SubscribeToNewsletter(ctx context.Context, cmd *commands.SubscribeToNewsletterCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	err := s.wameowService.SubscribeToNewsletter(ctx, cmd.SessionID, cmd.NewsletterJID)
	if err != nil {
		return fmt.Errorf("failed to subscribe to newsletter: %w", err)
	}

	// whatsmeow will handle events automatically
	return nil
}

// UnsubscribeFromNewsletter unsubscribes from a newsletter using command pattern
func (s *NewsletterService) UnsubscribeFromNewsletter(ctx context.Context, cmd *commands.UnsubscribeFromNewsletterCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	err := s.wameowService.UnsubscribeFromNewsletter(ctx, cmd.SessionID, cmd.NewsletterJID)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from newsletter: %w", err)
	}

	// whatsmeow will handle events automatically
	return nil
}

// SendNewsletterMessage sends a newsletter message using command pattern
func (s *NewsletterService) SendNewsletterMessage(ctx context.Context, cmd *commands.SendNewsletterMessageCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.SendNewsletterMessage(ctx, cmd.SessionID, cmd.NewsletterJID, cmd.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to send newsletter message: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// GetNewsletterMessages gets newsletter messages using command pattern
func (s *NewsletterService) GetNewsletterMessages(ctx context.Context, cmd *commands.GetNewsletterMessagesCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.GetNewsletterMessages(ctx, cmd.SessionID, cmd.NewsletterJID, cmd.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter messages: %w", err)
	}

	return result, nil
}

// GetNewsletterMessageUpdates gets newsletter message updates using command pattern
func (s *NewsletterService) GetNewsletterMessageUpdates(ctx context.Context, cmd *commands.GetNewsletterMessageUpdatesCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.GetNewsletterMessageUpdates(ctx, cmd.SessionID, cmd.NewsletterJID)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter message updates: %w", err)
	}

	return result, nil
}

// MarkNewsletterViewed marks newsletter as viewed using command pattern
func (s *NewsletterService) MarkNewsletterViewed(ctx context.Context, cmd *commands.MarkNewsletterViewedCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	return s.wameowService.MarkNewsletterViewed(ctx, cmd.SessionID, cmd.NewsletterJID, cmd.MessageID)
}

// SendNewsletterReaction sends newsletter reaction using command pattern
func (s *NewsletterService) SendNewsletterReaction(ctx context.Context, cmd *commands.SendNewsletterReactionCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	return s.wameowService.SendNewsletterReaction(ctx, cmd.SessionID, cmd.NewsletterJID, cmd.MessageID, cmd.Reaction)
}

// ToggleNewsletterMute toggles newsletter mute using command pattern
func (s *NewsletterService) ToggleNewsletterMute(ctx context.Context, cmd *commands.ToggleNewsletterMuteCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	return s.wameowService.ToggleNewsletterMute(ctx, cmd.SessionID, cmd.NewsletterJID, cmd.Mute)
}

// SubscribeLiveUpdates subscribes to live updates using command pattern
func (s *NewsletterService) SubscribeLiveUpdates(ctx context.Context, cmd *commands.SubscribeLiveUpdatesCommand) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return err
	}

	return s.wameowService.SubscribeLiveUpdates(ctx, cmd.SessionID, cmd.NewsletterJID)
}

// UploadNewsletterMedia uploads newsletter media using command pattern
func (s *NewsletterService) UploadNewsletterMedia(ctx context.Context, cmd *commands.UploadNewsletterMediaCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.UploadNewsletterMedia(ctx, cmd.SessionID, cmd.MediaData, cmd.MediaType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload newsletter media: %w", err)
	}

	return result, nil
}

// GetNewsletterByInvite gets newsletter by invite using command pattern
func (s *NewsletterService) GetNewsletterByInvite(ctx context.Context, cmd *commands.GetNewsletterByInviteCommand) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	if err := s.validateSession(ctx, cmd.SessionID); err != nil {
		return nil, err
	}

	result, err := s.wameowService.GetNewsletterByInvite(ctx, cmd.SessionID, cmd.InviteKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter by invite: %w", err)
	}

	return result, nil
}

// Helper methods

func (s *NewsletterService) validateSession(ctx context.Context, sessionID string) error {
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return fmt.Errorf("session is not connected")
	}

	return nil
}
