package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) GenerateSessionID() string {
	return uuid.New().String()
}

func (g *UUIDGenerator) GenerateAPIKey() string {
	bytes := make([]byte, 32) // 64 character hex string
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("api_%s", uuid.New().String())
	}
	return fmt.Sprintf("api_%s", hex.EncodeToString(bytes))
}
