package session

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// DisconnectSessionCommand represents the command to disconnect a session
type DisconnectSessionCommand struct {
	SessionID string
	Reason    string
}

// Validate validates the disconnect session command
func (c DisconnectSessionCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	// Reason is optional, but if provided should not be too long
	if len(c.Reason) > 500 {
		return common.NewValidationError("reason", c.Reason, "reason must not exceed 500 characters")
	}

	return nil
}

// DisconnectSessionResult represents the result of disconnecting a session
type DisconnectSessionResult struct {
	SessionID string
	Status    string
	Reason    string
}

// DisconnectSessionUseCase handles the disconnection of sessions from WhatsApp
type DisconnectSessionUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	eventPublisher  ports.EventPublisher
	logger          ports.Logger
}

// NewDisconnectSessionUseCase creates a new DisconnectSessionUseCase
func NewDisconnectSessionUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
) *DisconnectSessionUseCase {
	return &DisconnectSessionUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		eventPublisher:  eventPublisher,
		logger:          logger,
	}
}

// Handle executes the disconnect session use case
func (uc *DisconnectSessionUseCase) Handle(ctx context.Context, cmd DisconnectSessionCommand) (*DisconnectSessionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid disconnect session command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, cmd.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is already disconnected
	if sessionEntity.IsDisconnected() {
		uc.logger.Info(ctx, "Session already disconnected", "sessionID", cmd.SessionID)
		return &DisconnectSessionResult{
			SessionID: sessionEntity.SessionID().Value(),
			Status:    sessionEntity.Status().String(),
			Reason:    "already disconnected",
		}, nil
	}

	// 4. Set default reason if not provided
	reason := cmd.Reason
	if reason == "" {
		reason = "user requested"
	}

	// 5. Disconnect from WhatsApp service
	if err := uc.whatsappService.DisconnectSession(ctx, cmd.SessionID); err != nil {
		uc.logger.Warn(ctx, "Failed to disconnect session via WhatsApp service", "sessionID", cmd.SessionID, "error", err)
		// Continue with domain disconnection even if service call fails
	}

	// 6. Update session to disconnected state (domain operation)
	if err := sessionEntity.Disconnect(reason); err != nil {
		uc.logger.Error(ctx, "Failed to disconnect session in domain", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to disconnect session: %w", err)
	}

	// 7. Persist updated session
	if err := uc.sessionRepo.Update(ctx, sessionEntity); err != nil {
		uc.logger.Error(ctx, "Failed to update session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// 8. Publish domain events
	events := sessionEntity.GetEvents()
	if len(events) > 0 {
		if err := uc.eventPublisher.PublishBatch(ctx, events); err != nil {
			uc.logger.Warn(ctx, "Failed to publish domain events", "sessionID", cmd.SessionID, "error", err)
			// Don't fail the operation, just log the warning
		}
		sessionEntity.ClearEvents()
	}

	uc.logger.Info(ctx, "Session disconnected successfully", "sessionID", cmd.SessionID, "reason", reason)

	// 9. Return result
	return &DisconnectSessionResult{
		SessionID: sessionEntity.SessionID().Value(),
		Status:    sessionEntity.Status().String(),
		Reason:    reason,
	}, nil
}
