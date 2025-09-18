package application

import (
	"context"
	"fmt"

	"meow/internal/domain/session"
	"meow/internal/interfaces/dto"
	"meow/internal/shared/validation"
)

type NewsletterApp struct {
	newsletterService interface {
		CreateNewsletter(ctx context.Context, sessionID, name, description string) (*NewsletterInfo, error)
		GetNewsletter(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
		ListNewsletters(ctx context.Context, sessionID string) (*NewsletterInfo, error)
		SubscribeToNewsletter(ctx context.Context, sessionID, newsletterJID string) error
		UnsubscribeFromNewsletter(ctx context.Context, sessionID, newsletterJID string) error
		SendNewsletterMessage(ctx context.Context, sessionID, newsletterJID, content string) (*NewsletterInfo, error)
		GetNewsletterMessages(ctx context.Context, sessionID, newsletterJID string, limit int) (*NewsletterInfo, error)
		GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
		MarkNewsletterViewed(ctx context.Context, sessionID, newsletterJID, messageID string) error
		SendNewsletterReaction(ctx context.Context, sessionID, newsletterJID, messageID, reaction string) error
		ToggleNewsletterMute(ctx context.Context, sessionID, newsletterJID string, mute bool) error
		SubscribeLiveUpdates(ctx context.Context, sessionID, newsletterJID string) error
		UploadNewsletterMedia(ctx context.Context, sessionID string, mediaData []byte, mediaType string) (*NewsletterInfo, error)
		GetNewsletterByInvite(ctx context.Context, sessionID, inviteKey string) (*NewsletterInfo, error)
	}
	sessionRepo session.Repository
	validator   *validation.Validator
}

func NewNewsletterApp(
	newsletterService interface{},
	sessionRepo session.Repository,
	validator *validation.Validator,
) *NewsletterApp {
	return &NewsletterApp{
		newsletterService: newsletterService.(interface {
			CreateNewsletter(ctx context.Context, sessionID, name, description string) (*NewsletterInfo, error)
			GetNewsletter(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
			ListNewsletters(ctx context.Context, sessionID string) (*NewsletterInfo, error)
			SubscribeToNewsletter(ctx context.Context, sessionID, newsletterJID string) error
			UnsubscribeFromNewsletter(ctx context.Context, sessionID, newsletterJID string) error
			SendNewsletterMessage(ctx context.Context, sessionID, newsletterJID, content string) (*NewsletterInfo, error)
			GetNewsletterMessages(ctx context.Context, sessionID, newsletterJID string, limit int) (*NewsletterInfo, error)
			GetNewsletterMessageUpdates(ctx context.Context, sessionID, newsletterJID string) (*NewsletterInfo, error)
			MarkNewsletterViewed(ctx context.Context, sessionID, newsletterJID, messageID string) error
			SendNewsletterReaction(ctx context.Context, sessionID, newsletterJID, messageID, reaction string) error
			ToggleNewsletterMute(ctx context.Context, sessionID, newsletterJID string, mute bool) error
			SubscribeLiveUpdates(ctx context.Context, sessionID, newsletterJID string) error
			UploadNewsletterMedia(ctx context.Context, sessionID string, mediaData []byte, mediaType string) (*NewsletterInfo, error)
			GetNewsletterByInvite(ctx context.Context, sessionID, inviteKey string) (*NewsletterInfo, error)
		}),
		sessionRepo: sessionRepo,
		validator:   validator,
	}
}

func (s *NewsletterApp) CreateNewsletter(ctx context.Context, sessionID string, req *dto.CreateNewsletterRequest) (*NewsletterInfo, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.newsletterService.CreateNewsletter(ctx, sessionID, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create newsletter: %w", err)
	}

	return result, nil
}

func (s *NewsletterApp) GetNewsletter(ctx context.Context, sessionID string, req *dto.GetNewsletterInfoRequest) (*NewsletterInfo, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.newsletterService.GetNewsletter(ctx, sessionID, req.JID)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter: %w", err)
	}

	return result, nil
}

func (s *NewsletterApp) ListNewsletters(ctx context.Context, sessionID string) (*NewsletterInfo, error) {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	result, err := s.newsletterService.ListNewsletters(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list newsletters: %w", err)
	}

	return result, nil
}

func (s *NewsletterApp) SubscribeToNewsletter(ctx context.Context, sessionID string, newsletterJID string) error {
	if err := s.validateSession(ctx, sessionID); err != nil {
		return err
	}

	err := s.newsletterService.SubscribeToNewsletter(ctx, sessionID, newsletterJID)
	if err != nil {
		return fmt.Errorf("failed to subscribe to newsletter: %w", err)
	}

	return nil
}

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
