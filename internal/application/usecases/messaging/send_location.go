package messaging

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// SendLocationMessageCommand represents the command to send a location message
type SendLocationMessageCommand struct {
	SessionID string
	ChatJID   string
	Latitude  float64
	Longitude float64
	Name      string
	Address   string
}

// Validate validates the send location message command
func (c SendLocationMessageCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.ChatJID) == "" {
		return common.NewValidationError("chatJID", c.ChatJID, "chat JID is required")
	}

	if c.Latitude < -90 || c.Latitude > 90 {
		return common.NewValidationError("latitude", c.Latitude, "latitude must be between -90 and 90")
	}

	if c.Longitude < -180 || c.Longitude > 180 {
		return common.NewValidationError("longitude", c.Longitude, "longitude must be between -180 and 180")
	}

	if len(c.Name) > 100 {
		return common.NewValidationError("name", c.Name, "location name must not exceed 100 characters")
	}

	if len(c.Address) > 500 {
		return common.NewValidationError("address", c.Address, "address must not exceed 500 characters")
	}

	return nil
}

// SendLocationMessageResult represents the result of sending a location message
type SendLocationMessageResult struct {
	SessionID string
	ChatJID   string
	Latitude  float64
	Longitude float64
	MessageID string
	Sent      bool
}

// SendLocationMessageUseCase handles sending location messages via WhatsApp
type SendLocationMessageUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewSendLocationMessageUseCase creates a new SendLocationMessageUseCase
func NewSendLocationMessageUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *SendLocationMessageUseCase {
	return &SendLocationMessageUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the send location message use case
func (uc *SendLocationMessageUseCase) Handle(ctx context.Context, cmd SendLocationMessageCommand) (*SendLocationMessageResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid send location message command", "error", err)
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

	// 5. Send location message via WhatsApp service
	// Note: This would need to be added to the WhatsAppService interface
	// For now, we'll use a generic approach
	if err := uc.whatsappService.SendLocationMessage(ctx, cmd.SessionID, cmd.ChatJID, cmd.Latitude, cmd.Longitude, cmd.Name, cmd.Address); err != nil {
		uc.logger.Error(ctx, "Failed to send location message",
			"sessionID", cmd.SessionID,
			"chatJID", cmd.ChatJID,
			"latitude", cmd.Latitude,
			"longitude", cmd.Longitude,
			"error", err)
		return nil, fmt.Errorf("failed to send location message: %w", err)
	}

	uc.logger.Info(ctx, "Location message sent successfully",
		"sessionID", cmd.SessionID,
		"chatJID", cmd.ChatJID,
		"latitude", cmd.Latitude,
		"longitude", cmd.Longitude)

	// 6. Return result
	return &SendLocationMessageResult{
		SessionID: cmd.SessionID,
		ChatJID:   cmd.ChatJID,
		Latitude:  cmd.Latitude,
		Longitude: cmd.Longitude,
		MessageID: "", // Would be provided by WhatsApp service
		Sent:      true,
	}, nil
}
