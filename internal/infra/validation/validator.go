package validation

import (
	"fmt"
	"regexp"
	"strings"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateSessionName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("session name cannot be empty")
	}

	if len(name) < 3 {
		return fmt.Errorf("session name must be at least 3 characters")
	}

	if len(name) > 50 {
		return fmt.Errorf("session name cannot exceed 50 characters")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	if !matched {
		return fmt.Errorf("session name can only contain letters, numbers, dash and underscore")
	}

	return nil
}

func (v *Validator) ValidatePhoneNumber(phone string) error {
	if strings.TrimSpace(phone) == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	matched, _ := regexp.MatchString(`^\d+$`, cleaned)
	if !matched {
		return fmt.Errorf("phone number must contain only digits")
	}

	if len(cleaned) < 10 || len(cleaned) > 15 {
		return fmt.Errorf("phone number must be between 10 and 15 digits")
	}

	return nil
}

func (v *Validator) ValidateJID(jid string) error {
	if strings.TrimSpace(jid) == "" {
		return fmt.Errorf("JID cannot be empty")
	}

	if !strings.Contains(jid, "@") {
		return fmt.Errorf("JID must contain @ symbol")
	}

	parts := strings.Split(jid, "@")
	if len(parts) != 2 {
		return fmt.Errorf("JID must have exactly one @ symbol")
	}

	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("JID parts cannot be empty")
	}

	return nil
}

func (v *Validator) ValidateURL(url string) error {
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	return nil
}

func (v *Validator) ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func (v *Validator) Validate(data interface{}) error {
	return nil
}
