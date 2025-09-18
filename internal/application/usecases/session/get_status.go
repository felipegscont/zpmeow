package session

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// GetSessionStatusQuery represents the query to get session status
type GetSessionStatusQuery struct {
	SessionID string
}

// Validate validates the get session status query
func (q GetSessionStatusQuery) Validate() error {
	if strings.TrimSpace(q.SessionID) == "" {
		return common.NewValidationError("sessionID", q.SessionID, "session ID is required")
	}

	return nil
}

// SessionStatusView represents the session status view model
type SessionStatusView struct {
	SessionID       string
	Status          string
	IsConnected     bool
	IsAuthenticated bool
	WaJID           string
	QRCode          string
	LastSeen        string
	ConnectionInfo  map[string]interface{}
}

// GetSessionStatusUseCase handles getting session status information
type GetSessionStatusUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewGetSessionStatusUseCase creates a new GetSessionStatusUseCase
func NewGetSessionStatusUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *GetSessionStatusUseCase {
	return &GetSessionStatusUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the get session status use case
func (uc *GetSessionStatusUseCase) Handle(ctx context.Context, query GetSessionStatusQuery) (*SessionStatusView, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid get session status query", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, query.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", query.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Get real-time status from WhatsApp service
	whatsappStatus, err := uc.whatsappService.GetSessionStatus(ctx, query.SessionID)
	if err != nil {
		uc.logger.Warn(ctx, "Failed to get WhatsApp status", "sessionID", query.SessionID, "error", err)
		// Continue with domain status if WhatsApp service fails
		whatsappStatus = sessionEntity.Status().String()
	}

	// 4. Get QR code if session is not authenticated
	qrCode := ""
	if !sessionEntity.IsAuthenticated() && (sessionEntity.IsConnecting() || sessionEntity.IsDisconnected()) {
		qrCode, err = uc.whatsappService.GetQRCode(ctx, query.SessionID)
		if err != nil {
			uc.logger.Warn(ctx, "Failed to get QR code", "sessionID", query.SessionID, "error", err)
			// Don't fail, QR code might not be available
		}
	}

	// 5. Build connection info
	connectionInfo := map[string]interface{}{
		"domain_status":   sessionEntity.Status().String(),
		"whatsapp_status": whatsappStatus,
		"has_proxy":       sessionEntity.HasProxy(),
		"has_webhook":     sessionEntity.HasWebhook(),
		"created_at":      sessionEntity.CreatedAt().Value().Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":      sessionEntity.UpdatedAt().Value().Format("2006-01-02T15:04:05Z07:00"),
	}

	// 6. Create status view
	statusView := &SessionStatusView{
		SessionID:       sessionEntity.SessionID().Value(),
		Status:          whatsappStatus, // Use real-time status
		IsConnected:     sessionEntity.IsConnected(),
		IsAuthenticated: sessionEntity.IsAuthenticated(),
		WaJID:           sessionEntity.WaJID().Value(),
		QRCode:          qrCode,
		LastSeen:        sessionEntity.UpdatedAt().Value().Format("2006-01-02T15:04:05Z07:00"),
		ConnectionInfo:  connectionInfo,
	}

	uc.logger.Debug(ctx, "Session status retrieved successfully",
		"sessionID", query.SessionID,
		"status", whatsappStatus,
		"isConnected", sessionEntity.IsConnected())

	return statusView, nil
}
