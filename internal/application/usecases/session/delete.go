package session

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// DeleteSessionCommand represents the command to delete a session
type DeleteSessionCommand struct {
	SessionID string
	Force     bool // Force deletion even if session is connected
}

// Validate validates the delete session command
func (c DeleteSessionCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	return nil
}

// DeleteSessionResult represents the result of deleting a session
type DeleteSessionResult struct {
	SessionID string
	Name      string
	Deleted   bool
}

// DeleteSessionUseCase handles the deletion of sessions
type DeleteSessionUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	eventPublisher  ports.EventPublisher
	logger          ports.Logger
}

// NewDeleteSessionUseCase creates a new DeleteSessionUseCase
func NewDeleteSessionUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
) *DeleteSessionUseCase {
	return &DeleteSessionUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		eventPublisher:  eventPublisher,
		logger:          logger,
	}
}

// Handle executes the delete session use case
func (uc *DeleteSessionUseCase) Handle(ctx context.Context, cmd DeleteSessionCommand) (*DeleteSessionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid delete session command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, cmd.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check business rules for deletion
	if sessionEntity.IsConnected() && !cmd.Force {
		return nil, common.NewBusinessRuleError(
			"session_deletion_not_allowed",
			"cannot delete connected session without force flag",
		)
	}

	// 4. Disconnect session if it's still connected
	if sessionEntity.IsConnected() || sessionEntity.IsConnecting() {
		uc.logger.Info(ctx, "Disconnecting session before deletion", "sessionID", cmd.SessionID)

		// Disconnect from WhatsApp service
		if err := uc.whatsappService.DisconnectSession(ctx, cmd.SessionID); err != nil {
			uc.logger.Warn(ctx, "Failed to disconnect session via WhatsApp service", "sessionID", cmd.SessionID, "error", err)
			// Continue with deletion even if service call fails
		}

		// Update domain state
		if err := sessionEntity.Disconnect("deletion requested"); err != nil {
			uc.logger.Warn(ctx, "Failed to disconnect session in domain", "sessionID", cmd.SessionID, "error", err)
			// Continue with deletion
		}
	}

	// 5. Mark session for deletion (domain operation)
	sessionEntity.Delete()

	// 6. Publish domain events before deletion
	events := sessionEntity.GetEvents()
	if len(events) > 0 {
		if err := uc.eventPublisher.PublishBatch(ctx, events); err != nil {
			uc.logger.Warn(ctx, "Failed to publish domain events", "sessionID", cmd.SessionID, "error", err)
			// Don't fail the operation, just log the warning
		}
	}

	// 7. Delete session from repository
	if err := uc.sessionRepo.Delete(ctx, cmd.SessionID); err != nil {
		uc.logger.Error(ctx, "Failed to delete session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to delete session: %w", err)
	}

	uc.logger.Info(ctx, "Session deleted successfully", "sessionID", cmd.SessionID, "name", sessionEntity.Name().Value())

	// 8. Return result
	return &DeleteSessionResult{
		SessionID: sessionEntity.SessionID().Value(),
		Name:      sessionEntity.Name().Value(),
		Deleted:   true,
	}, nil
}
