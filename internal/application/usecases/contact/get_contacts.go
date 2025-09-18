package contact

import (
	"context"
	"fmt"
	"strings"

	"meow/internal/application/common"
	"meow/internal/application/ports"
)

// GetContactsQuery represents the query to get contacts
type GetContactsQuery struct {
	SessionID string
	Limit     int
	Offset    int
}

// Validate validates the get contacts query
func (q GetContactsQuery) Validate() error {
	if strings.TrimSpace(q.SessionID) == "" {
		return common.NewValidationError("sessionID", q.SessionID, "session ID is required")
	}

	if q.Limit < 0 {
		return common.NewValidationError("limit", q.Limit, "limit cannot be negative")
	}

	if q.Limit > 1000 {
		return common.NewValidationError("limit", q.Limit, "limit cannot exceed 1000")
	}

	if q.Offset < 0 {
		return common.NewValidationError("offset", q.Offset, "offset cannot be negative")
	}

	return nil
}

// ContactView represents a contact view model
type ContactView struct {
	JID          string
	Name         string
	Notify       string
	PushName     string
	BusinessName string
	IsBlocked    bool
	IsMuted      bool
	IsContact    bool
	Avatar       string
}

// GetContactsResult represents the result of getting contacts
type GetContactsResult struct {
	SessionID string
	Contacts  []ContactView
	Total     int
	Limit     int
	Offset    int
}

// GetContactsUseCase handles getting contacts for a session
type GetContactsUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewGetContactsUseCase creates a new GetContactsUseCase
func NewGetContactsUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *GetContactsUseCase {
	return &GetContactsUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the get contacts use case
func (uc *GetContactsUseCase) Handle(ctx context.Context, query GetContactsQuery) (*GetContactsResult, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid get contacts query", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, query.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", query.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is connected (business rule)
	if !sessionEntity.IsConnected() {
		return nil, common.NewBusinessRuleError(
			"session_not_connected",
			fmt.Sprintf("session must be connected to get contacts, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Get contacts via WhatsApp service
	contacts, err := uc.whatsappService.GetContacts(ctx, query.SessionID, query.Limit, query.Offset)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get contacts",
			"sessionID", query.SessionID,
			"error", err)
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	// 5. Convert to view models
	contactViews := make([]ContactView, len(contacts))
	for i, contact := range contacts {
		contactViews[i] = ContactView{
			JID:          contact.JID,
			Name:         contact.Name,
			Notify:       contact.Notify,
			PushName:     contact.PushName,
			BusinessName: contact.BusinessName,
			IsBlocked:    contact.IsBlocked,
			IsMuted:      contact.IsMuted,
			IsContact:    contact.IsContact,
			Avatar:       contact.Avatar,
		}
	}

	uc.logger.Debug(ctx, "Contacts retrieved successfully",
		"sessionID", query.SessionID,
		"count", len(contactViews))

	// 6. Return result
	return &GetContactsResult{
		SessionID: query.SessionID,
		Contacts:  contactViews,
		Total:     len(contactViews), // In a real implementation, this would be the total count
		Limit:     query.Limit,
		Offset:    query.Offset,
	}, nil
}

// CheckContactQuery represents the query to check if a number is on WhatsApp
type CheckContactQuery struct {
	SessionID string
	Phone     string
}

// Validate validates the check contact query
func (q CheckContactQuery) Validate() error {
	if strings.TrimSpace(q.SessionID) == "" {
		return common.NewValidationError("sessionID", q.SessionID, "session ID is required")
	}

	if strings.TrimSpace(q.Phone) == "" {
		return common.NewValidationError("phone", q.Phone, "phone number is required")
	}

	return nil
}

// CheckContactResult represents the result of checking a contact
type CheckContactResult struct {
	SessionID    string
	Phone        string
	IsOnWhatsApp bool
	JID          string
}

// CheckContactUseCase handles checking if a phone number is on WhatsApp
type CheckContactUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewCheckContactUseCase creates a new CheckContactUseCase
func NewCheckContactUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *CheckContactUseCase {
	return &CheckContactUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the check contact use case
func (uc *CheckContactUseCase) Handle(ctx context.Context, query CheckContactQuery) (*CheckContactResult, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid check contact query", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, query.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", query.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is connected (business rule)
	if !sessionEntity.IsConnected() {
		return nil, common.NewBusinessRuleError(
			"session_not_connected",
			fmt.Sprintf("session must be connected to check contacts, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Check contact via WhatsApp service
	isOnWhatsApp, jid, err := uc.whatsappService.CheckContact(ctx, query.SessionID, query.Phone)
	if err != nil {
		uc.logger.Error(ctx, "Failed to check contact",
			"sessionID", query.SessionID,
			"phone", query.Phone,
			"error", err)
		return nil, fmt.Errorf("failed to check contact: %w", err)
	}

	uc.logger.Debug(ctx, "Contact checked successfully",
		"sessionID", query.SessionID,
		"phone", query.Phone,
		"isOnWhatsApp", isOnWhatsApp)

	// 5. Return result
	return &CheckContactResult{
		SessionID:    query.SessionID,
		Phone:        query.Phone,
		IsOnWhatsApp: isOnWhatsApp,
		JID:          jid,
	}, nil
}

// GetUserInfoQuery represents the query to get user information
type GetUserInfoQuery struct {
	SessionID string
	UserJID   string
}

// Validate validates the get user info query
func (q GetUserInfoQuery) Validate() error {
	if strings.TrimSpace(q.SessionID) == "" {
		return common.NewValidationError("sessionID", q.SessionID, "session ID is required")
	}

	if strings.TrimSpace(q.UserJID) == "" {
		return common.NewValidationError("userJID", q.UserJID, "user JID is required")
	}

	return nil
}

// UserInfoView represents detailed user information
type UserInfoView struct {
	JID          string
	Name         string
	Notify       string
	PushName     string
	BusinessName string
	Phone        string
	Status       string
	Avatar       string
	IsBlocked    bool
	IsMuted      bool
	IsContact    bool
	LastSeen     string
}

// GetUserInfoUseCase handles getting detailed user information
type GetUserInfoUseCase struct {
	sessionRepo     ports.SessionRepository
	whatsappService ports.WhatsAppService
	logger          ports.Logger
}

// NewGetUserInfoUseCase creates a new GetUserInfoUseCase
func NewGetUserInfoUseCase(
	sessionRepo ports.SessionRepository,
	whatsappService ports.WhatsAppService,
	logger ports.Logger,
) *GetUserInfoUseCase {
	return &GetUserInfoUseCase{
		sessionRepo:     sessionRepo,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// Handle executes the get user info use case
func (uc *GetUserInfoUseCase) Handle(ctx context.Context, query GetUserInfoQuery) (*UserInfoView, error) {
	// 1. Validate query
	if err := query.Validate(); err != nil {
		uc.logger.Warn(ctx, "Invalid get user info query", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get session from repository
	sessionEntity, err := uc.sessionRepo.GetByID(ctx, query.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get session", "sessionID", query.SessionID, "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 3. Check if session is connected (business rule)
	if !sessionEntity.IsConnected() {
		return nil, common.NewBusinessRuleError(
			"session_not_connected",
			fmt.Sprintf("session must be connected to get user info, current status: %s", sessionEntity.Status()),
		)
	}

	// 4. Get user info via WhatsApp service
	userInfo, err := uc.whatsappService.GetUserInfo(ctx, query.SessionID, query.UserJID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get user info",
			"sessionID", query.SessionID,
			"userJID", query.UserJID,
			"error", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// 5. Convert to view model
	userInfoView := &UserInfoView{
		JID:          userInfo.JID,
		Name:         userInfo.Name,
		Notify:       userInfo.Notify,
		PushName:     userInfo.PushName,
		BusinessName: userInfo.BusinessName,
		Phone:        userInfo.Phone,
		Status:       userInfo.Status,
		Avatar:       userInfo.Avatar,
		IsBlocked:    userInfo.IsBlocked,
		IsMuted:      userInfo.IsMuted,
		IsContact:    userInfo.IsContact,
		LastSeen:     userInfo.LastSeen,
	}

	uc.logger.Debug(ctx, "User info retrieved successfully",
		"sessionID", query.SessionID,
		"userJID", query.UserJID,
		"userName", userInfo.Name)

	return userInfoView, nil
}
