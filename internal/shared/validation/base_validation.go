package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"zpmeow/internal/shared/errors"
)

type BaseValidator struct{}

func NewBaseValidator() *BaseValidator {
	return &BaseValidator{}
}

func (v *BaseValidator) Validate(value interface{}) error {
	return nil
}

var (
	sessionIDRegex   = regexp.MustCompile(`^[a-fA-F0-9]{32}$|^[a-zA-Z0-9-]{3,50}$`)
	phoneNumberRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	urlRegex         = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	messageIDRegex   = regexp.MustCompile(`^[A-Za-z0-9-]+$`)
)

func (v *BaseValidator) ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	if !IsValidPhoneNumber(phone) {
		return fmt.Errorf("invalid phone number format: %s", phone)
	}

	return nil
}

func (v *BaseValidator) ValidatePhoneNumbers(phones []string) error {
	if len(phones) == 0 {
		return fmt.Errorf("at least one phone number is required")
	}

	for i, phone := range phones {
		if err := v.ValidatePhoneNumber(phone); err != nil {
			return fmt.Errorf("phone number at index %d: %w", i, err)
		}
	}

	return nil
}

func (v *BaseValidator) ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	if len(url) > 2048 {
		return fmt.Errorf("URL too long (max 2048 characters)")
	}

	if !urlRegex.MatchString(url) {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

func (v *BaseValidator) ValidateSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	if !IsValidSessionID(sessionID) {
		return fmt.Errorf("invalid session ID format: %s", sessionID)
	}

	return nil
}

func (v *BaseValidator) ValidateMessageID(messageID string) error {
	if messageID == "" {
		return fmt.Errorf("message ID cannot be empty")
	}

	if len(messageID) < 10 || len(messageID) > 100 {
		return fmt.Errorf("message ID length must be between 10 and 100 characters")
	}

	if !messageIDRegex.MatchString(messageID) {
		return fmt.Errorf("message ID contains invalid characters")
	}

	return nil
}

func (v *BaseValidator) ValidateMessageIDs(messageIDs []string) error {
	if len(messageIDs) == 0 {
		return fmt.Errorf("at least one message ID is required")
	}

	for i, messageID := range messageIDs {
		if err := v.ValidateMessageID(messageID); err != nil {
			return fmt.Errorf("message ID at index %d: %w", i, err)
		}
	}

	return nil
}

func (v *BaseValidator) ValidateStringLength(value, fieldName string, minLength, maxLength int) error {
	if value == "" && minLength > 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	if len(value) < minLength {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, minLength)
	}

	if len(value) > maxLength {
		return fmt.Errorf("%s must be at most %d characters long", fieldName, maxLength)
	}

	return nil
}

func (v *BaseValidator) ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.NewValidationError(fmt.Sprintf("%s is required", fieldName))
	}
	return nil
}

func (v *BaseValidator) ValidateOptionalLength(value, fieldName string, maxLength int) error {
	if value == "" {
		return nil // Optional field
	}

	if len(value) > maxLength {
		return fmt.Errorf("%s too long (max %d characters)", fieldName, maxLength)
	}

	return nil
}

func (v *BaseValidator) ValidateArrayLength(arr []string, fieldName string, minLength, maxLength int) error {
	if len(arr) < minLength {
		return fmt.Errorf("%s must have at least %d items", fieldName, minLength)
	}

	if len(arr) > maxLength {
		return fmt.Errorf("%s cannot have more than %d items", fieldName, maxLength)
	}

	return nil
}

func (v *BaseValidator) ValidateNoEmptyStrings(arr []string, fieldName string) error {
	for i, item := range arr {
		if strings.TrimSpace(item) == "" {
			return fmt.Errorf("%s at index %d cannot be empty", fieldName, i)
		}
	}
	return nil
}

func (v *BaseValidator) ValidateRegex(value, fieldName, pattern string) error {
	if value == "" {
		return nil // Skip validation for empty values
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern for %s", fieldName)
	}

	if !regex.MatchString(value) {
		return fmt.Errorf("invalid %s format", fieldName)
	}

	return nil
}

func (v *BaseValidator) ValidateNumericRange(value float64, fieldName string, min, max float64) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %g and %g", fieldName, min, max)
	}
	return nil
}

func (v *BaseValidator) ValidatePositive(value int64, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", fieldName)
	}
	return nil
}

func (v *BaseValidator) ValidateNonNegative(value int, fieldName string) error {
	if value < 0 {
		return fmt.Errorf("%s cannot be negative", fieldName)
	}
	return nil
}

func IsValidSessionID(id string) bool {
	return sessionIDRegex.MatchString(strings.TrimSpace(id))
}

func IsValidPhoneNumber(phoneNumber string) bool {
	return phoneNumberRegex.MatchString(strings.TrimSpace(phoneNumber))
}

func IsValidURL(url string) bool {
	return urlRegex.MatchString(strings.TrimSpace(url))
}

func IsValidMessageID(messageID string) bool {
	return messageIDRegex.MatchString(strings.TrimSpace(messageID))
}

var DefaultBaseValidator = &BaseValidator{}

// Utility functions consolidated from utils/validation.go

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidHTTPURL(url string) bool {
	if url == "" {
		return false
	}
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func SanitizeString(input string) string {
	cleaned := regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}

func GenerateUUID() string {
	return uuid.New().String()
}
