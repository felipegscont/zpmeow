package session

import (
	"context"
	"fmt"
	"strings"

	"zpmeow/internal/application/common"
	"zpmeow/internal/application/ports"
)

// PairPhoneCommand represents the command to pair a phone with a session
type PairPhoneCommand struct {
	SessionID   string
	PhoneNumber string
}

// Validate validates the pair phone command
func (c PairPhoneCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	if strings.TrimSpace(c.PhoneNumber) == "" {
		return common.NewValidationError("phoneNumber", c.PhoneNumber, "phone number is required")
	}

	// Basic phone number validation
	phoneNumber := strings.TrimSpace(c.PhoneNumber)
	if len(phoneNumber) < 10 || len(phoneNumber) > 15 {
		return common.NewValidationError("phoneNumber", c.PhoneNumber, "phone number must be between 10 and 15 digits")
	}

	return nil
}

// PairPhoneResult represents the result of pairing a phone
type PairPhoneResult struct {
	SessionID   string
	PhoneNumber string
	Success     bool
	Message     string
}

// PairPhoneUseCase handles pairing a phone number with a session
type PairPhoneUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	eventPublisher  ports.EventPublisher
	logger          ports.Logger
}

// NewPairPhoneUseCase creates a new PairPhoneUseCase
func NewPairPhoneUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
) *PairPhoneUseCase {
	return &PairPhoneUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		eventPublisher:  eventPublisher,
		logger:          logger,
	}
}

// Handle executes the pair phone use case
func (uc *PairPhoneUseCase) Handle(ctx context.Context, cmd PairPhoneCommand) (*PairPhoneResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid pair phone command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, cmd.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session can be paired (business rule)
	if !sessionEntity.CanConnect() {
		return nil, common.NewBusinessRuleError(
			"session_pairing_not_allowed",
			fmt.Sprintf("session cannot be paired from current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Pair phone via WhatsApp service
	if err := uc.whatsappService.PairWithPhone(ctx, cmd.SessionID, cmd.PhoneNumber); err != nil {
		uc.logger.Error(ctx, "Failed to pair phone with session",
			"sessionID", cmd.SessionID,
			"phoneNumber", cmd.PhoneNumber,
			"error", err)

		// Set session to error state
		sessionEntity.SetError("phone pairing failed: " + err.Error())

		// Try to update session state
		if updateErr := uc.sessionRepo.Update(ctx, sessionEntity); updateErr != nil {
			uc.logger.Error(ctx, "Failed to update session after pairing error", "sessionID", cmd.SessionID, "error", updateErr)
		}

		return nil, fmt.Errorf("failed to pair phone with session: %w", err)
	}

	// 5. Update session state if needed
	if err := uc.sessionRepo.Update(ctx, sessionEntity); err != nil {
		uc.logger.Error(ctx, "Failed to update session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// 6. Publish domain events
	events := sessionEntity.GetEvents()
	if len(events) > 0 {
		if err := uc.eventPublisher.PublishBatch(ctx, events); err != nil {
			uc.logger.Warn(ctx, "Failed to publish domain events", "sessionID", cmd.SessionID, "error", err)
			// Don't fail the operation, just log the warning
		}
		sessionEntity.ClearEvents()
	}

	uc.logger.Info(ctx, "Phone paired successfully",
		"sessionID", cmd.SessionID,
		"phoneNumber", cmd.PhoneNumber)

	// 7. Return result
	return &PairPhoneResult{
		SessionID:   cmd.SessionID,
		PhoneNumber: cmd.PhoneNumber,
		Success:     true,
		Message:     "Phone paired successfully",
	}, nil
}
