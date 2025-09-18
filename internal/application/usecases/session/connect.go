package session

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// ConnectSessionCommand represents the command to connect a session
type ConnectSessionCommand struct {
	SessionID string
}

// Validate validates the connect session command
func (c ConnectSessionCommand) Validate() error {
	if strings.TrimSpace(c.SessionID) == "" {
		return common.NewValidationError("sessionID", c.SessionID, "session ID is required")
	}

	return nil
}

// ConnectSessionResult represents the result of connecting a session
type ConnectSessionResult struct {
	SessionID string
	Status    string
	QRCode    string
}

// ConnectSessionUseCase handles the connection of sessions to WhatsApp
type ConnectSessionUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	eventPublisher  ports.EventPublisher
	logger          ports.Logger
}

// NewConnectSessionUseCase creates a new ConnectSessionUseCase
func NewConnectSessionUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
) *ConnectSessionUseCase {
	return &ConnectSessionUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		eventPublisher:  eventPublisher,
		logger:          logger,
	}
}

// Handle executes the connect session use case
func (uc *ConnectSessionUseCase) Handle(ctx context.Context, cmd ConnectSessionCommand) (*ConnectSessionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid connect session command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, cmd.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session can be connected (domain business rule)
	if !sessionEntity.CanConnect() {
		return nil, common.NewBusinessRuleError(
			"session_connection_not_allowed",
			fmt.Sprintf("session cannot be connected from current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Initiate connection through WhatsApp service
	if err := uc.whatsappService.ConnectSession(ctx, cmd.SessionID); err != nil {
		uc.logger.Error(ctx, "Failed to connect session via WhatsApp service", "sessionID", cmd.SessionID, "error", err)

		// Set session to error state
		sessionEntity.SetError("connection failed: " + err.Error())

		// Try to update session state
		if updateErr := uc.sessionRepo.Update(ctx, sessionEntity); updateErr != nil {
			uc.logger.Error(ctx, "Failed to update session after connection error", "sessionID", cmd.SessionID, "error", updateErr)
		}

		return nil, fmt.Errorf("failed to connect session: %w", err)
	}

	// 5. Update session to connecting state
	if err := sessionEntity.Connect(); err != nil {
		uc.logger.Error(ctx, "Failed to set session to connecting state", "sessionID", cmd.SessionID, "error", err)
		return nil, fmt.Errorf("failed to update session state: %w", err)
	}

	// 6. Get QR code if available
	qrCode := ""
	if !sessionEntity.IsAuthenticated() {
		qrCode, err = uc.whatsappService.GetQRCode(ctx, cmd.SessionID)
		if err != nil {
			uc.logger.Warn(ctx, "Failed to get QR code", "sessionID", cmd.SessionID, "error", err)
			// Don't fail the operation, QR code might not be available yet
		} else if qrCode != "" {
			// Update session with QR code
			if err := sessionEntity.SetQRCode(qrCode); err != nil {
				uc.logger.Warn(ctx, "Failed to set QR code in session", "sessionID", cmd.SessionID, "error", err)
			}
		}
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

	uc.logger.Info(ctx, "Session connection initiated", "sessionID", cmd.SessionID, "status", sessionEntity.Status())

	// 9. Return result
	return &ConnectSessionResult{
		SessionID: sessionEntity.SessionID().Value(),
		Status:    sessionEntity.Status().String(),
		QRCode:    qrCode,
	}, nil
}
