package session

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Compiled regex for UUID validation (performance optimization)
	uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// IdentifierService handles session identifier resolution and validation
type IdentifierService interface {
	// ResolveIdentifier determines if input is UUID or name and validates format
	ResolveIdentifier(identifier string) (*IdentifierInfo, error)

	// IsValidUUID checks if string is a valid UUID format
	IsValidUUID(identifier string) bool

	// IsValidName checks if string is a valid session name format
	IsValidName(identifier string) bool

	// NormalizeIdentifier cleans and normalizes the identifier
	NormalizeIdentifier(identifier string) string
}

// IdentifierInfo contains information about a resolved identifier
type IdentifierInfo struct {
	Original   string
	Normalized string
	Type       IdentifierType
	IsValid    bool
}

// IdentifierType represents the type of identifier
type IdentifierType string

const (
	IdentifierTypeUUID    IdentifierType = "uuid"
	IdentifierTypeName    IdentifierType = "name"
	IdentifierTypeUnknown IdentifierType = "unknown"
)

// SessionIdentifierService implements the IdentifierService interface
type SessionIdentifierService struct{}

// NewIdentifierService creates a new session identifier service
func NewIdentifierService() IdentifierService {
	return &SessionIdentifierService{}
}

// ResolveIdentifier determines the type and validates the identifier
func (s *SessionIdentifierService) ResolveIdentifier(identifier string) (*IdentifierInfo, error) {
	if identifier == "" {
		return nil, fmt.Errorf("identifier cannot be empty")
	}

	normalized := s.NormalizeIdentifier(identifier)

	info := &IdentifierInfo{
		Original:   identifier,
		Normalized: normalized,
	}

	// Check if it's a UUID
	if s.IsValidUUID(normalized) {
		info.Type = IdentifierTypeUUID
		info.IsValid = true
		return info, nil
	}

	// Check if it's a valid name
	if s.IsValidName(normalized) {
		info.Type = IdentifierTypeName
		info.IsValid = true
		return info, nil
	}

	// Neither UUID nor valid name
	info.Type = IdentifierTypeUnknown
	info.IsValid = false
	return info, fmt.Errorf("invalid identifier format: must be UUID or valid session name")
}

// IsValidUUID checks if the identifier is a valid UUID
func (s *SessionIdentifierService) IsValidUUID(identifier string) bool {
	if identifier == "" {
		return false
	}

	// Quick format check: UUID should be 36 characters with hyphens
	if len(identifier) != 36 || !strings.Contains(identifier, "-") {
		return false
	}

	// Validate with compiled regex pattern (pure domain validation)
	return uuidRegex.MatchString(strings.ToLower(identifier))
}

// IsValidName checks if the identifier is a valid session name
func (s *SessionIdentifierService) IsValidName(identifier string) bool {
	if identifier == "" {
		return false
	}

	// Use the existing SessionName validation logic
	_, err := NewSessionName(identifier)
	return err == nil
}

// NormalizeIdentifier cleans and normalizes the identifier
func (s *SessionIdentifierService) NormalizeIdentifier(identifier string) string {
	// Trim whitespace
	normalized := strings.TrimSpace(identifier)

	// Convert to lowercase for UUIDs (UUIDs are case-insensitive)
	if s.looksLikeUUID(normalized) {
		normalized = strings.ToLower(normalized)
	}

	return normalized
}

// looksLikeUUID does a quick check if string looks like UUID format
func (s *SessionIdentifierService) looksLikeUUID(identifier string) bool {
	return len(identifier) == 36 && strings.Contains(identifier, "-")
}

// GetIdentifierType returns the type of identifier without full validation
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

// ValidateIdentifierFormat validates identifier format without checking existence
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

// SuggestCorrection suggests a corrected version of invalid identifier
func (s *SessionIdentifierService) SuggestCorrection(identifier string) string {
	normalized := s.NormalizeIdentifier(identifier)

	// If it looks like UUID but has issues, try to fix common problems
	if s.looksLikeUUID(normalized) {
		// Remove extra characters, fix case
		return strings.ToLower(normalized)
	}

	// For names, remove invalid characters and truncate if needed
	if len(normalized) > 100 {
		normalized = normalized[:100]
	}

	// Remove invalid characters for session names
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

// IsEmpty checks if identifier is empty or whitespace only
func (s *SessionIdentifierService) IsEmpty(identifier string) bool {
	return strings.TrimSpace(identifier) == ""
}

// GetValidationRules returns validation rules for identifiers
func (s *SessionIdentifierService) GetValidationRules() map[string]string {
	return map[string]string{
		"uuid_format":     "Must be valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
		"name_length":     "Must be between 3 and 100 characters",
		"name_characters": "Can only contain letters, numbers, hyphens, and underscores",
		"not_empty":       "Cannot be empty or whitespace only",
	}
}
