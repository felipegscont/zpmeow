package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ID represents a unique identifier in the domain
type ID struct {
	value string
}

// IDGenerator defines the interface for ID generation
type IDGenerator interface {
	Generate() string
}

// defaultIDGenerator is a simple counter-based generator for domain purity
var defaultGenerator IDGenerator = &counterIDGenerator{counter: 0}

// SetIDGenerator allows injection of ID generator (for testing or different implementations)
func SetIDGenerator(generator IDGenerator) {
	defaultGenerator = generator
}

// counterIDGenerator is a simple implementation for domain purity
type counterIDGenerator struct {
	counter int
}

func (g *counterIDGenerator) Generate() string {
	g.counter++
	return fmt.Sprintf("domain-id-%d", g.counter)
}

// NewID creates a new ID from a string value
func NewID(value string) (ID, error) {
	if value == "" {
		return ID{}, fmt.Errorf("ID cannot be empty")
	}

	// Basic validation - just check it's not empty and reasonable length
	if len(value) < 1 || len(value) > 100 {
		return ID{}, fmt.Errorf("invalid ID format: must be between 1 and 100 characters")
	}

	return ID{value: value}, nil
}

// GenerateID creates a new ID using the configured generator
func GenerateID() ID {
	return ID{value: defaultGenerator.Generate()}
}

// Value returns the string value of the ID
func (id ID) Value() string {
	return id.value
}

// String returns the string representation
func (id ID) String() string {
	return id.value
}

// IsEmpty checks if the ID is empty
func (id ID) IsEmpty() bool {
	return id.value == ""
}

// Equals compares two IDs
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}

// Timestamp represents a domain timestamp
type Timestamp struct {
	value time.Time
}

// NewTimestamp creates a new timestamp
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{value: t}
}

// Now creates a timestamp with current time
func Now() Timestamp {
	return Timestamp{value: time.Now()}
}

// Value returns the time value
func (ts Timestamp) Value() time.Time {
	return ts.value
}

// IsZero checks if timestamp is zero
func (ts Timestamp) IsZero() bool {
	return ts.value.IsZero()
}

// Before checks if this timestamp is before another
func (ts Timestamp) Before(other Timestamp) bool {
	return ts.value.Before(other.value)
}

// After checks if this timestamp is after another
func (ts Timestamp) After(other Timestamp) bool {
	return ts.value.After(other.value)
}

// Name represents a domain name with validation rules
type Name struct {
	value string
}

// NewName creates a new name with validation
func NewName(value string, minLength, maxLength int) (Name, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return Name{}, fmt.Errorf("name cannot be empty")
	}

	if len(trimmed) < minLength {
		return Name{}, fmt.Errorf("name must be at least %d characters long", minLength)
	}

	if len(trimmed) > maxLength {
		return Name{}, fmt.Errorf("name cannot exceed %d characters", maxLength)
	}

	// Basic validation - only alphanumeric, hyphens, underscores
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !nameRegex.MatchString(trimmed) {
		return Name{}, fmt.Errorf("name can only contain letters, numbers, hyphens, and underscores")
	}

	return Name{value: trimmed}, nil
}

// Value returns the string value
func (n Name) Value() string {
	return n.value
}

// String returns the string representation
func (n Name) String() string {
	return n.value
}

// IsEmpty checks if name is empty
func (n Name) IsEmpty() bool {
	return n.value == ""
}

// Equals compares two names
func (n Name) Equals(other Name) bool {
	return n.value == other.value
}
