package session

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
	"meow/internal/domain/session"
)

// CreateSessionCommand represents the command to create a new session
type CreateSessionCommand struct {
	Name               string
	ProxyConfiguration string
	WebhookEndpoint    string
}

// Validate validates the create session command
func (c CreateSessionCommand) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return common.NewValidationError("name", c.Name, "session name is required")
	}

	if len(c.Name) < 3 {
		return common.NewValidationError("name", c.Name, "session name must be at least 3 characters")
	}

	if len(c.Name) > 100 {
		return common.NewValidationError("name", c.Name, "session name must not exceed 100 characters")
	}

	return nil
}

// CreateSessionResult represents the result of creating a session
type CreateSessionResult struct {
	SessionID string
	Name      string
	Status    string
	ApiKey    string
}

// CreateSessionUseCase handles the creation of new sessions
type CreateSessionUseCase struct {
	sessionRepo    ports.SessionRepository
	idGenerator    ports.IDGenerator
	eventPublisher ports.EventPublisher
	logger         ports.Logger
}

// NewCreateSessionUseCase creates a new CreateSessionUseCase
func NewCreateSessionUseCase(
	sessionRepo ports.SessionRepository,
	idGenerator ports.IDGenerator,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo:    sessionRepo,
		idGenerator:    idGenerator,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// Handle executes the create session use case
func (uc *CreateSessionUseCase) Handle(ctx context.Context, cmd CreateSessionCommand) (*CreateSessionResult, error) {
	// 1. Validate command
	if err := cmd.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid create session command", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Check if session with same name already exists
	exists, err := uc.sessionRepo.Exists(ctx, cmd.Name)
	if err != nil {
		uc.logger.Error(ctx, "Failed to check session existence", "name", cmd.Name, "error", err)
		return nil, fmt.Errorf("failed to check session existence: %w", err)
	}

	if exists {
		return nil, common.NewBusinessRuleError("unique_session_name", "session with this name already exists")
	}

	// 3. Generate unique identifiers
	sessionID := uc.idGenerator.GenerateSessionID()
	apiKey := uc.idGenerator.GenerateAPIKey()

	// 4. Create domain entity
	sessionEntity, err := session.NewSession(sessionID, cmd.Name)
	if err != nil {
		uc.logger.Error(ctx, "Failed to create session entity", "sessionID", sessionID, "name", cmd.Name, "error", err)
		return nil, fmt.Errorf("failed to create session entity: %w", err)
	}

	// 5. Configure session if additional parameters provided
	if cmd.ProxyConfiguration != "" {
		if err := sessionEntity.SetProxyConfiguration(cmd.ProxyConfiguration); err != nil {
			uc.logger.Warn(ctx, "Invalid proxy configuration", "proxy", cmd.ProxyConfiguration, "error", err)
			return nil, fmt.Errorf("invalid proxy configuration: %w", err)
		}
	}

	if cmd.WebhookEndpoint != "" {
		if err := sessionEntity.SetWebhookEndpoint(cmd.WebhookEndpoint); err != nil {
			uc.logger.Warn(ctx, "Invalid webhook endpoint", "webhook", cmd.WebhookEndpoint, "error", err)
			return nil, fmt.Errorf("invalid webhook endpoint: %w", err)
		}
	}

	// 6. Set the generated API key
	if err := sessionEntity.SetApiKey(apiKey); err != nil {
		uc.logger.Error(ctx, "Failed to set API key", "sessionID", sessionID, "error", err)
		return nil, fmt.Errorf("failed to set API key: %w", err)
	}

	// 7. Persist the session
	if err := uc.sessionRepo.Create(ctx, sessionEntity); err != nil {
		uc.logger.Error(ctx, "Failed to persist session", "sessionID", sessionID, "error", err)
		return nil, fmt.Errorf("failed to persist session: %w", err)
	}

	// 8. Publish domain events
	events := sessionEntity.GetEvents()
	if len(events) > 0 {
		if err := uc.eventPublisher.PublishBatch(ctx, events); err != nil {
			uc.logger.Warn(ctx, "Failed to publish domain events", "sessionID", sessionID, "error", err)
			// Don't fail the operation, just log the warning
		}
		sessionEntity.ClearEvents()
	}

	uc.logger.Info(ctx, "Session created successfully", "sessionID", sessionID, "name", cmd.Name)

	// 9. Return result
	return &CreateSessionResult{
		SessionID: sessionEntity.SessionID().Value(),
		Name:      sessionEntity.Name().Value(),
		Status:    sessionEntity.Status().String(),
		ApiKey:    sessionEntity.ApiKey().Value(),
	}, nil
}
