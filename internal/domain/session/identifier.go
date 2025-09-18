package session

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type IdentifierService interface {
	ResolveIdentifier(identifier string) (*IdentifierInfo, error)

	IsValidUUID(identifier string) bool

	IsValidName(identifier string) bool

	NormalizeIdentifier(identifier string) string
}

type IdentifierInfo struct {
	Original   string
	Normalized string
	Type       IdentifierType
	IsValid    bool
}

type IdentifierType string

const (
	IdentifierTypeUUID    IdentifierType = "uuid"
	IdentifierTypeName    IdentifierType = "name"
	IdentifierTypeUnknown IdentifierType = "unknown"
)

type SessionIdentifierService struct{}

func NewIdentifierService() IdentifierService {
	return &SessionIdentifierService{}
}

func (s *SessionIdentifierService) ResolveIdentifier(identifier string) (*IdentifierInfo, error) {
	if identifier == "" {
		return nil, fmt.Errorf("identifier cannot be empty")
	}

	normalized := s.NormalizeIdentifier(identifier)

	info := &IdentifierInfo{
		Original:   identifier,
		Normalized: normalized,
	}

	if s.IsValidUUID(normalized) {
		info.Type = IdentifierTypeUUID
		info.IsValid = true
		return info, nil
	}

	if s.IsValidName(normalized) {
		info.Type = IdentifierTypeName
		info.IsValid = true
		return info, nil
	}

	info.Type = IdentifierTypeUnknown
	info.IsValid = false
	return info, fmt.Errorf("invalid identifier format: must be UUID or valid session name")
}

func (s *SessionIdentifierService) IsValidUUID(identifier string) bool {
	if identifier == "" {
		return false
	}

	if len(identifier) != 36 || !strings.Contains(identifier, "-") {
		return false
	}

	return uuidRegex.MatchString(strings.ToLower(identifier))
}

func (s *SessionIdentifierService) IsValidName(identifier string) bool {
	if identifier == "" {
		return false
	}

	_, err := NewSessionName(identifier)
	return err == nil
}

func (s *SessionIdentifierService) NormalizeIdentifier(identifier string) string {
	normalized := strings.TrimSpace(identifier)

	if s.looksLikeUUID(normalized) {
		normalized = strings.ToLower(normalized)
	}

	return normalized
}

func (s *SessionIdentifierService) looksLikeUUID(identifier string) bool {
	return len(identifier) == 36 && strings.Contains(identifier, "-")
}

func (s *SessionIdentifierService) GetIdentifierType(identifier string) IdentifierType {
	normalized := s.NormalizeIdentifier(identifier)

	if s.looksLikeUUID(normalized) {
		return IdentifierTypeUUID
	}

	if len(normalized) >= 3 && len(normalized) <= 100 {
		return IdentifierTypeName
	}

	return IdentifierTypeUnknown
}

func (s *SessionIdentifierService) ValidateIdentifierFormat(identifier string) error {
	info, err := s.ResolveIdentifier(identifier)
	if err != nil {
		return err
	}

	if !info.IsValid {
		return fmt.Errorf("invalid identifier format")
	}

	return nil
}

func (s *SessionIdentifierService) SuggestCorrection(identifier string) string {
	normalized := s.NormalizeIdentifier(identifier)

	if s.looksLikeUUID(normalized) {
		return strings.ToLower(normalized)
	}

	if len(normalized) > 100 {
		normalized = normalized[:100]
	}

	var result strings.Builder
	for _, r := range normalized {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}

	suggestion := result.String()
	if len(suggestion) < 3 {
		return "session-" + suggestion
	}

	return suggestion
}

func (s *SessionIdentifierService) IsEmpty(identifier string) bool {
	return strings.TrimSpace(identifier) == ""
}

func (s *SessionIdentifierService) GetValidationRules() map[string]string {
	return map[string]string{
		"uuid_format":     "Must be valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
		"name_length":     "Must be between 3 and 100 characters",
		"name_characters": "Can only contain letters, numbers, hyphens, and underscores",
		"not_empty":       "Cannot be empty or whitespace only",
	}
}
