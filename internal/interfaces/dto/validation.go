package dto

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Validator wraps the go-playground validator with custom validation methods
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance with custom validation rules
func NewValidator() *Validator {
	validate := validator.New()

	// Register custom validation functions
	validate.RegisterValidation("session_id", validateSessionID)
	validate.RegisterValidation("phone_number", validatePhoneNumber)
	validate.RegisterValidation("message_id", validateMessageID)
	validate.RegisterValidation("webhook_url", validateWebhookURL)
	validate.RegisterValidation("jid", validateJID)

	return &Validator{
		validate: validate,
	}
}

// Validate validates a struct using struct tags
func (v *Validator) Validate(value interface{}) error {
	return v.validate.Struct(value)
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// Custom validation functions
func validateSessionID(fl validator.FieldLevel) bool {
	return IsValidSessionID(fl.Field().String())
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	return IsValidPhoneNumber(fl.Field().String())
}

func validateMessageID(fl validator.FieldLevel) bool {
	return IsValidMessageID(fl.Field().String())
}

func validateWebhookURL(fl validator.FieldLevel) bool {
	return IsValidHTTPURL(fl.Field().String())
}

func validateJID(fl validator.FieldLevel) bool {
	jid := fl.Field().String()
	return strings.Contains(jid, "@") && len(jid) > 5
}

var (
	sessionIDRegex   = regexp.MustCompile(`^[a-fA-F0-9]{32}$|^[a-zA-Z0-9-]{3,50}$`)
	phoneNumberRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	// urlRegex and messageIDRegex are defined in types.go to avoid duplication
)

func (v *Validator) ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	if !IsValidPhoneNumber(phone) {
		return fmt.Errorf("invalid phone number format: %s", phone)
	}

	return nil
}

func (v *Validator) ValidatePhoneNumbers(phones []string) error {
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

func (v *Validator) ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	if len(url) > 2048 {
		return fmt.Errorf("URL too long (max 2048 characters)")
	}

	// Use URL type validation from types.go
	if _, err := NewURL(url); err != nil {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

func (v *Validator) ValidateSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	if !IsValidSessionID(sessionID) {
		return fmt.Errorf("invalid session ID format: %s", sessionID)
	}

	return nil
}

func (v *Validator) ValidateMessageID(messageID string) error {
	if messageID == "" {
		return fmt.Errorf("message ID cannot be empty")
	}

	if len(messageID) < 10 || len(messageID) > 100 {
		return fmt.Errorf("message ID length must be between 10 and 100 characters")
	}

	// Use MessageID type validation from types.go
	if _, err := NewMessageID(messageID); err != nil {
		return fmt.Errorf("message ID contains invalid characters")
	}

	return nil
}

func (v *Validator) ValidateMessageIDs(messageIDs []string) error {
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

func (v *Validator) ValidateStringLength(value, fieldName string, minLength, maxLength int) error {
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

func (v *Validator) ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func (v *Validator) ValidateOptionalLength(value, fieldName string, maxLength int) error {
	if value == "" {
		return nil // Optional field
	}

	if len(value) > maxLength {
		return fmt.Errorf("%s too long (max %d characters)", fieldName, maxLength)
	}

	return nil
}

func (v *Validator) ValidateArrayLength(arr []string, fieldName string, minLength, maxLength int) error {
	if len(arr) < minLength {
		return fmt.Errorf("%s must have at least %d items", fieldName, minLength)
	}

	if len(arr) > maxLength {
		return fmt.Errorf("%s cannot have more than %d items", fieldName, maxLength)
	}

	return nil
}

func (v *Validator) ValidateNoEmptyStrings(arr []string, fieldName string) error {
	for i, item := range arr {
		if strings.TrimSpace(item) == "" {
			return fmt.Errorf("%s at index %d cannot be empty", fieldName, i)
		}
	}
	return nil
}

func (v *Validator) ValidateRegex(value, fieldName, pattern string) error {
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

func (v *Validator) ValidateNumericRange(value float64, fieldName string, min, max float64) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %g and %g", fieldName, min, max)
	}
	return nil
}

func (v *Validator) ValidatePositive(value int64, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", fieldName)
	}
	return nil
}

func (v *Validator) ValidateNonNegative(value int, fieldName string) error {
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
	_, err := NewURL(strings.TrimSpace(url))
	return err == nil
}

func IsValidMessageID(messageID string) bool {
	_, err := NewMessageID(strings.TrimSpace(messageID))
	return err == nil
}

var DefaultValidator = &Validator{}

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
