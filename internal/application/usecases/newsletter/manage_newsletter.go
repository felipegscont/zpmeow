package newsletter

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// CreateNewsletterCommand represents the command to create a newsletter
type CreateNewsletterCommand struct {
	SessionID   string
	Name        string
	Description string
}

// Validate validates the create newsletter command
func (c CreateNewsletterCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.Name) == "" {
		return common.NewValidationError("name", c.Name, "newsletter name is required")
	}

	if len(c.Name) > 100 {
		return common.NewValidationError("name", c.Name, "newsletter name must not exceed 100 characters")
	}

	if len(c.Description) > 500 {
		return common.NewValidationError("description", c.Description, "newsletter description must not exceed 500 characters")
	}

	return nil
}

// SubscribeNewsletterCommand represents the command to subscribe to a newsletter
type SubscribeNewsletterCommand struct {
	SessionID     string
	NewsletterJID string
}

// Validate validates the subscribe newsletter command
func (c SubscribeNewsletterCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.NewsletterJID) == "" {
		return common.NewValidationError("newsletterJID", c.NewsletterJID, "newsletter JID is required")
	}

	return nil
}

// UnsubscribeNewsletterCommand represents the command to unsubscribe from a newsletter
type UnsubscribeNewsletterCommand struct {
	SessionID     string
	NewsletterJID string
}

// Validate validates the unsubscribe newsletter command
func (c UnsubscribeNewsletterCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.NewsletterJID) == "" {
		return common.NewValidationError("newsletterJID", c.NewsletterJID, "newsletter JID is required")
	}

	return nil
}

// GetNewsletterInfoQuery represents the query to get newsletter information
type GetNewsletterInfoQuery struct {
	SessionID     string
	NewsletterJID string
}

// Validate validates the get newsletter info query
func (q GetNewsletterInfoQuery) Validate() error {
	if strings.TrimSpace(q.SessionID) == "" {
		return common.NewValidationError("sessionID", q.SessionID, "session ID is required")
	}

	if strings.TrimSpace(q.NewsletterJID) == "" {
		return common.NewValidationError("newsletterJID", q.NewsletterJID, "newsletter JID is required")
	}

	return nil
}

// NewsletterView represents a newsletter view model
type NewsletterView struct {
	JID             string
	Name            string
	Description     string
	SubscriberCount int
	IsSubscribed    bool
	CreatedAt       string
	UpdatedAt       string
}

// NewsletterResult represents the result of newsletter operations
type NewsletterResult struct {
	SessionID     string
	NewsletterJID string
	Action        string
	Success       bool
	Message       string
	Newsletter    *NewsletterView
}

// CreateNewsletterUseCase handles creating newsletters
type CreateNewsletterUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewCreateNewsletterUseCase creates a new CreateNewsletterUseCase
func NewCreateNewsletterUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *CreateNewsletterUseCase {
	return &CreateNewsletterUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the create newsletter use case
func (uc *CreateNewsletterUseCase) Handle(ctx context.Context, cmd CreateNewsletterCommand) (*NewsletterResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid create newsletter command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, cmd.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is connected (business rule)
	if !sessionEntity.IsConnected() {
		return nil, common.NewBusinessRuleError(
			"session_not_connected",
			fmt.Sprintf("session must be connected to create newsletters, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Create newsletter via WhatsApp service
	newsletterJID, err := uc.whatsappService.CreateNewsletter(ctx, cmd.SessionID, cmd.Name, cmd.Description)
	if err != nil {
		uc.logger.Error(ctx, "Failed to create newsletter",
			"sessionID", cmd.SessionID,
			"name", cmd.Name,
			"error", err)
		return nil, fmt.Errorf("failed to create newsletter: %w", err)
	}

	// 5. Get newsletter info to return complete data
	newsletterInfo, err := uc.whatsappService.GetNewsletterInfo(ctx, cmd.SessionID, newsletterJID)
	if err != nil {
		uc.logger.Warn(ctx, "Failed to get newsletter info after creation",
			"sessionID", cmd.SessionID,
			"newsletterJID", newsletterJID,
			"error", err)
		// Don't fail the operation, just return basic info
		newsletterInfo = &ports.NewsletterInfo{
			JID:         newsletterJID,
			Name:        cmd.Name,
			Description: cmd.Description,
		}
	}

	// 6. Convert to view model
	newsletterView := &NewsletterView{
		JID:             newsletterInfo.JID,
		Name:            newsletterInfo.Name,
		Description:     newsletterInfo.Description,
		SubscriberCount: newsletterInfo.SubscriberCount,
		IsSubscribed:    newsletterInfo.IsSubscribed,
		CreatedAt:       newsletterInfo.CreatedAt,
		UpdatedAt:       newsletterInfo.UpdatedAt,
	}

	uc.logger.Info(ctx, "Newsletter created successfully",
		"sessionID", cmd.SessionID,
		"newsletterJID", newsletterJID,
		"name", cmd.Name)

	// 7. Return result
	return &NewsletterResult{
		SessionID:     cmd.SessionID,
		NewsletterJID: newsletterJID,
		Action:        "create",
		Success:       true,
		Message:       "Newsletter created successfully",
		Newsletter:    newsletterView,
	}, nil
}

// SubscribeNewsletterUseCase handles subscribing to newsletters
type SubscribeNewsletterUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewSubscribeNewsletterUseCase creates a new SubscribeNewsletterUseCase
func NewSubscribeNewsletterUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *SubscribeNewsletterUseCase {
	return &SubscribeNewsletterUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the subscribe newsletter use case
func (uc *SubscribeNewsletterUseCase) Handle(ctx context.Context, cmd SubscribeNewsletterCommand) (*NewsletterResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid subscribe newsletter command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, cmd.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is connected (business rule)
	if !sessionEntity.IsConnected() {
		return nil, common.NewBusinessRuleError(
			"session_not_connected",
			fmt.Sprintf("session must be connected to subscribe to newsletters, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Subscribe to newsletter via WhatsApp service
	if err := uc.whatsappService.SubscribeNewsletter(ctx, cmd.SessionID, cmd.NewsletterJID); err != nil {
		uc.logger.Error(ctx, "Failed to subscribe to newsletter",
			"sessionID", cmd.SessionID,
			"newsletterJID", cmd.NewsletterJID,
			"error", err)
		return nil, fmt.Errorf("failed to subscribe to newsletter: %w", err)
	}

	uc.logger.Info(ctx, "Successfully subscribed to newsletter",
		"sessionID", cmd.SessionID,
		"newsletterJID", cmd.NewsletterJID)

	// 5. Return result
	return &NewsletterResult{
		SessionID:     cmd.SessionID,
		NewsletterJID: cmd.NewsletterJID,
		Action:        "subscribe",
		Success:       true,
		Message:       "Successfully subscribed to newsletter",
	}, nil
}
