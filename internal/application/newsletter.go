package application

import (
	"context"
	"fmt"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

// NewsletterApp implements newsletter use cases following Clean Architecture
type NewsletterApp struct {
	meowService interface {
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

// NewNewsletterApp creates a new NewsletterApp instance
func NewNewsletterApp(
	meowService interface{},
	sessionRepo session.Repository,
	validator *validation.Validator,
) *NewsletterApp {
	return &NewsletterApp{
		meowService: meowService.(interface {
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

// CreateNewsletter creates a new newsletter using DTO
func (s *NewsletterApp) CreateNewsletter(ctx context.Context, sessionID string, req *dto.CreateNewsletterRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.meowService.CreateNewsletter(ctx, sessionID, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create newsletter: %w", err)
	}

	// whatsmeow will handle events automatically
	return result, nil
}

// GetNewsletter gets newsletter information using DTO
func (s *NewsletterApp) GetNewsletter(ctx context.Context, sessionID string, req *dto.GetNewsletterInfoRequest) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.meowService.GetNewsletter(ctx, sessionID, req.JID)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter: %w", err)
	}

	return result, nil
}

// ListNewsletters lists all newsletters
func (s *NewsletterApp) ListNewsletters(ctx context.Context, sessionID string) (interface{}, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.meowService.ListNewsletters(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list newsletters: %w", err)
	}

	return result, nil
}

// SubscribeToNewsletter subscribes to a newsletter using DTO
func (s *NewsletterApp) SubscribeToNewsletter(ctx context.Context, sessionID string, newsletterJID string) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	err := s.meowService.SubscribeToNewsletter(ctx, sessionID, newsletterJID)
	if err != nil {
		return fmt.Errorf("failed to subscribe to newsletter: %w", err)
	}

	// whatsmeow will handle events automatically
	return nil
}

// Helper methods

// Helper methods

func (s *NewsletterApp) validateSession(ctx context.Context, sessionID string) error {
	sessionEntity, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if sessionEntity.Status != session.StatusConnected {
		return fmt.Errorf("session is not connected")
	}

	return nil
}
