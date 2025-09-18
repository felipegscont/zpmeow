package messaging

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// ContactInfo represents contact information to be sent
type ContactInfo = ports.ContactInfo

// validateContactInfo validates the contact info
func validateContactInfo(ci ContactInfo) error {
	if strings.TrimSpace(ci.Name) == "" {
		return common.NewValidationError("name", ci.Name, "contact name is required")
	}

	if strings.TrimSpace(ci.Phone) == "" {
		return common.NewValidationError("phone", ci.Phone, "contact phone is required")
	}

	if len(ci.Name) > 100 {
		return common.NewValidationError("name", ci.Name, "contact name must not exceed 100 characters")
	}

	if len(ci.Organization) > 100 {
		return common.NewValidationError("organization", ci.Organization, "organization must not exceed 100 characters")
	}

	if len(ci.Email) > 100 {
		return common.NewValidationError("email", ci.Email, "email must not exceed 100 characters")
	}

	return nil
}

// SendContactMessageCommand represents the command to send a contact message
type SendContactMessageCommand struct {
	SessionID string
	ChatJID   string
	Contacts  []ContactInfo
}

// Validate validates the send contact message command
func (c SendContactMessageCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if len(c.Contacts) == 0 {
		return common.NewValidationError("contacts", "", "at least one contact is required")
	}

	if len(c.Contacts) > 10 {
		return common.NewValidationError("contacts", "", "maximum 10 contacts allowed")
	}

	for i, contact := range c.Contacts {
		if err := validateContactInfo(contact); err != nil {
			return fmt.Errorf("contact %d validation failed: %w", i, err)
		}
	}

	return nil
}

// SendContactMessageResult represents the result of sending a contact message
type SendContactMessageResult struct {
	SessionID    string
	ChatJID      string
	ContactCount int
	MessageID    string
	Sent         bool
}

// SendContactMessageUseCase handles sending contact messages via WhatsApp
type SendContactMessageUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewSendContactMessageUseCase creates a new SendContactMessageUseCase
func NewSendContactMessageUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *SendContactMessageUseCase {
	return &SendContactMessageUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the send contact message use case
func (uc *SendContactMessageUseCase) Handle(ctx context.Context, cmd SendContactMessageCommand) (*SendContactMessageResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid send contact message command", "error", err)
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
			fmt.Sprintf("session must be connected to send messages, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Check if session is authenticated
	if !sessionEntity.IsAuthenticated() {
		return nil, common.NewBusinessRuleError(
			"session_not_authenticated",
			"session must be authenticated to send messages",
		)
	}

	// 5. Send contact message via WhatsApp service
	// Note: This would need to be added to the WhatsAppService interface
	if err := uc.whatsappService.SendContactMessage(ctx, cmd.SessionID, cmd.ChatJID, cmd.Contacts); err != nil {
		uc.logger.Error(ctx, "Failed to send contact message",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"contactCount", len(cmd.Contacts),
			"error", err)
		return nil, fmt.Errorf("failed to send contact message: %w", err)
	}

	uc.logger.Info(ctx, "Contact message sent successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"contactCount", len(cmd.Contacts))

	// 6. Return result
	return &SendContactMessageResult{
		SessionID:    cmd.SessionID,
		ChatJID:      cmd.ChatJID,
		ContactCount: len(cmd.Contacts),
		MessageID:    "", // Would be provided by WhatsApp service
		Sent:         true,
	}, nil
}
