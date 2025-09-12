package shared

// Re-export shared types and utilities to avoid aliases

import (
	"zpmeow/internal/shared/errors"
	"zpmeow/internal/shared/patterns"
	"zpmeow/internal/shared/types"
	"zpmeow/internal/shared/utils"
)

// Common types
type ID = types.ID
type Timestamp = types.Timestamp
type Status = types.Status
type SendResponse = types.SendResponse
type HTTPLogEntry = types.HTTPLogEntry
type LogLevel = types.LogLevel

// Status constants
const (
	StatusDisconnected = types.StatusDisconnected
	StatusConnecting   = types.StatusConnecting
	StatusConnected    = types.StatusConnected
	StatusError        = types.StatusError
	StatusDeleted      = types.StatusDeleted
)

// Log level constants
const (
	LogLevelInfo  = types.LogLevelInfo
	LogLevelWarn  = types.LogLevelWarn
	LogLevelError = types.LogLevelError
)

// Utilities

// Utility instances
var (
	JID = utils.JID
)

// Errors

// Error utilities
var (
	Error = errors.Error
	MapDomainError = errors.MapDomainError
	IsValidationError = errors.IsValidationError
	IsConflictError = errors.IsConflictError
	IsNotFoundError = errors.IsNotFoundError
)

// Error mappings
var (
	DomainErrorMappings = errors.DomainErrorMappings
)

// Patterns

// Pattern types
type MediaStrategy = patterns.MediaStrategy
type MediaSender = patterns.MediaSender
type MediaStrategyFactory = patterns.MediaStrategyFactory
type Converter[T, U any] = patterns.Converter[T, U]
type BatchConverter[T, U any] = patterns.BatchConverter[T, U]
type BidirectionalConverter[T, U any] = patterns.BidirectionalConverter[T, U]
type EntityToDTOConverter[Entity, DTO any] = patterns.EntityToDTOConverter[Entity, DTO]
type DTOToEntityConverter[DTO, Entity any] = patterns.DTOToEntityConverter[DTO, Entity]
