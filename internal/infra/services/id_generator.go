package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

// UUIDGenerator implements the IDGenerator interface using UUIDs
type UUIDGenerator struct{}

// NewUUIDGenerator creates a new UUID-based ID generator
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

// GenerateSessionID generates a unique session identifier using UUID
func (g *UUIDGenerator) GenerateSessionID() string {
	return uuid.New().String()
}

// GenerateAPIKey generates a unique API key using random hex
func (g *UUIDGenerator) GenerateAPIKey() string {
	bytes := make([]byte, 32) // 64 character hex string
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to UUID if random fails
		return fmt.Sprintf("api_%s", uuid.New().String())
	}
	return fmt.Sprintf("api_%s", hex.EncodeToString(bytes))
}
