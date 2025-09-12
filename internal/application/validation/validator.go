package validation

import (
	"fmt"
	"regexp"
	"strings"
	"zpmeow/internal/application/dto/request"
	"zpmeow/internal/shared/errors"
)

// Validator provides validation for application DTOs
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// Validation patterns
var (
	phoneNumberRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	sessionIDRegex   = regexp.MustCompile(`^[a-fA-F0-9]{32}$|^[a-zA-Z0-9_-]{3,50}$`)
	urlRegex         = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
)

// ValidateCreateSessionRequest validates create session request
func (v *Validator) ValidateCreateSessionRequest(req *request.CreateSessionRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.NewValidationError("session name is required")
	}

	name := strings.TrimSpace(req.Name)
	if len(name) < 3 {
		return errors.NewValidationError("session name must be at least 3 characters long")
	}

	if len(name) > 50 {
		return errors.NewValidationError("session name must be at most 50 characters long")
	}

	// Check for valid characters
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
		return errors.NewValidationError("session name can only contain letters, numbers, underscores, and hyphens")
	}

	// Check for reserved names
	reservedNames := []string{"admin", "system", "default", "test", "api", "www"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(name, reserved) {
			return errors.NewValidationError(fmt.Sprintf("'%s' is a reserved session name", reserved))
		}
	}

	return nil
}

// ValidateGetSessionRequest validates get session request
func (v *Validator) ValidateGetSessionRequest(req *request.GetSessionRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.IDOrName) == "" {
		return errors.NewValidationError("session id or name is required")
	}

	return nil
}

// ValidateDeleteSessionRequest validates delete session request
func (v *Validator) ValidateDeleteSessionRequest(req *request.DeleteSessionRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ID) == "" {
		return errors.NewValidationError("session id is required")
	}

	return nil
}

// ValidateConnectSessionRequest validates connect session request
func (v *Validator) ValidateConnectSessionRequest(req *request.ConnectSessionRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ID) == "" {
		return errors.NewValidationError("session id is required")
	}

	return nil
}

// ValidateDisconnectSessionRequest validates disconnect session request
func (v *Validator) ValidateDisconnectSessionRequest(req *request.DisconnectSessionRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ID) == "" {
		return errors.NewValidationError("session id is required")
	}

	return nil
}

// ValidateGetQRCodeRequest validates get QR code request
func (v *Validator) ValidateGetQRCodeRequest(req *request.GetQRCodeRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ID) == "" {
		return errors.NewValidationError("session id is required")
	}

	return nil
}

// ValidatePairSessionRequest validates pair session request
func (v *Validator) ValidatePairSessionRequest(req *request.PairSessionRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.SessionID) == "" {
		return errors.NewValidationError("session id is required")
	}

	if strings.TrimSpace(req.PhoneNumber) == "" {
		return errors.NewValidationError("phone number is required")
	}

	// Validate phone number format
	phoneNumber := strings.TrimSpace(req.PhoneNumber)
	if !phoneNumberRegex.MatchString(phoneNumber) {
		return errors.NewValidationError("invalid phone number format")
	}

	return nil
}

// ValidateSetProxyRequest validates set proxy request
func (v *Validator) ValidateSetProxyRequest(req *request.SetProxyRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ID) == "" {
		return errors.NewValidationError("session id is required")
	}

	if strings.TrimSpace(req.ProxyURL) == "" {
		return errors.NewValidationError("proxy url is required")
	}

	// Validate URL format
	proxyURL := strings.TrimSpace(req.ProxyURL)
	if !urlRegex.MatchString(proxyURL) {
		return errors.NewValidationError("invalid proxy url format")
	}

	return nil
}

// ValidateClearProxyRequest validates clear proxy request
func (v *Validator) ValidateClearProxyRequest(req *request.ClearProxyRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ID) == "" {
		return errors.NewValidationError("session id is required")
	}

	return nil
}

// ValidateGetSessionStatusRequest validates get session status request
func (v *Validator) ValidateGetSessionStatusRequest(req *request.GetSessionStatusRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ID) == "" {
		return errors.NewValidationError("session id is required")
	}

	return nil
}

// ValidateProxyRequest validates proxy request
func (v *Validator) ValidateProxyRequest(req *request.ProxyRequest) error {
	if req == nil {
		return errors.NewValidationError("request cannot be nil")
	}

	if strings.TrimSpace(req.ProxyURL) == "" {
		return errors.NewValidationError("proxy url is required")
	}

	// Validate URL format
	proxyURL := strings.TrimSpace(req.ProxyURL)
	if !urlRegex.MatchString(proxyURL) {
		return errors.NewValidationError("invalid proxy url format")
	}

	return nil
}

// Helper functions for common validations

// IsValidSessionID checks if a session ID is valid
func IsValidSessionID(id string) bool {
	return sessionIDRegex.MatchString(strings.TrimSpace(id))
}

// IsValidPhoneNumber checks if a phone number is valid
func IsValidPhoneNumber(phoneNumber string) bool {
	return phoneNumberRegex.MatchString(strings.TrimSpace(phoneNumber))
}

// IsValidURL checks if a URL is valid
func IsValidURL(url string) bool {
	return urlRegex.MatchString(strings.TrimSpace(url))
}
