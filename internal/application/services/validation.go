package services

import (
	"fmt"
	"regexp"
	"strings"

	"zpmeow/internal/shared/utils"
)

// ApplicationValidationService provides centralized validation for application layer
type ApplicationValidationService struct{}

// ValidatePhoneNumber validates phone numbers for application use
func (v *ApplicationValidationService) ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}
	
	if !utils.IsValidPhoneNumber(phone) {
		return fmt.Errorf("invalid phone number format: %s", phone)
	}
	
	return nil
}

// ValidatePhoneNumbers validates multiple phone numbers
func (v *ApplicationValidationService) ValidatePhoneNumbers(phones []string) error {
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

// ValidateMessageID validates WhatsApp message IDs
func (v *ApplicationValidationService) ValidateMessageID(messageID string) error {
	if messageID == "" {
		return fmt.Errorf("message ID cannot be empty")
	}
	
	// WhatsApp message IDs are typically alphanumeric with some special characters
	if len(messageID) < 10 || len(messageID) > 100 {
		return fmt.Errorf("message ID length must be between 10 and 100 characters")
	}
	
	// Basic format validation
	validMessageID := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if !validMessageID.MatchString(messageID) {
		return fmt.Errorf("message ID contains invalid characters")
	}
	
	return nil
}

// ValidateMessageIDs validates multiple message IDs
func (v *ApplicationValidationService) ValidateMessageIDs(messageIDs []string) error {
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

// ValidateWebhookURL validates webhook URLs
func (v *ApplicationValidationService) ValidateWebhookURL(url string) error {
	if url == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}
	
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("webhook URL must start with http:// or https://")
	}
	
	if len(url) > 2048 {
		return fmt.Errorf("webhook URL too long (max 2048 characters)")
	}
	
	return nil
}

// ValidateWebhookEvents validates webhook event types
func (v *ApplicationValidationService) ValidateWebhookEvents(events []string) error {
	if len(events) == 0 {
		return nil // Empty events list is valid (means all events)
	}
	
	// This would use the constants from meow/core/constants.go
	// For now, we'll do basic validation
	for i, event := range events {
		if event == "" {
			return fmt.Errorf("event at index %d cannot be empty", i)
		}
		
		if len(event) > 50 {
			return fmt.Errorf("event name at index %d too long (max 50 characters)", i)
		}
	}
	
	return nil
}

// ValidateTextMessage validates text message content
func (v *ApplicationValidationService) ValidateTextMessage(text string) error {
	if text == "" {
		return fmt.Errorf("message text cannot be empty")
	}
	
	// WhatsApp has a limit of ~65,536 characters for text messages
	if len(text) > 65536 {
		return fmt.Errorf("message text too long (max 65,536 characters)")
	}
	
	return nil
}

// ValidateCaption validates media captions
func (v *ApplicationValidationService) ValidateCaption(caption string) error {
	// Caption is optional, so empty is valid
	if caption == "" {
		return nil
	}
	
	// WhatsApp caption limit is typically around 1024 characters
	if len(caption) > 1024 {
		return fmt.Errorf("caption too long (max 1024 characters)")
	}
	
	return nil
}

// ValidateFilename validates file names
func (v *ApplicationValidationService) ValidateFilename(filename string) error {
	if filename == "" {
		return nil // Filename is optional
	}
	
	if len(filename) > 255 {
		return fmt.Errorf("filename too long (max 255 characters)")
	}
	
	// Check for invalid characters in filename
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	if invalidChars.MatchString(filename) {
		return fmt.Errorf("filename contains invalid characters")
	}
	
	return nil
}

// ValidateVCard validates vCard format
func (v *ApplicationValidationService) ValidateVCard(vcard string) error {
	if vcard == "" {
		return fmt.Errorf("vCard cannot be empty")
	}
	
	if !strings.HasPrefix(vcard, "BEGIN:VCARD") {
		return fmt.Errorf("vCard must start with BEGIN:VCARD")
	}
	
	if !strings.HasSuffix(vcard, "END:VCARD") {
		return fmt.Errorf("vCard must end with END:VCARD")
	}
	
	return nil
}

// ValidatePollOptions validates poll options
func (v *ApplicationValidationService) ValidatePollOptions(options []string) error {
	if len(options) < 2 {
		return fmt.Errorf("poll must have at least 2 options")
	}
	
	if len(options) > 12 {
		return fmt.Errorf("poll cannot have more than 12 options")
	}
	
	for i, option := range options {
		if option == "" {
			return fmt.Errorf("poll option at index %d cannot be empty", i)
		}
		
		if len(option) > 100 {
			return fmt.Errorf("poll option at index %d too long (max 100 characters)", i)
		}
	}
	
	return nil
}

// Global instance
var DefaultValidationService = &ApplicationValidationService{}
