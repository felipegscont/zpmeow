package logging

import (
	"context"
	"log"

	"meow/internal/application/ports"
)

// LoggerAdapter adapts the existing logging to implement application ports
type LoggerAdapter struct {
	// We can wrap the existing logger here if needed
}

// NewLoggerAdapter creates a new logger adapter
func NewLoggerAdapter() ports.Logger {
	return &LoggerAdapter{}
}

// Debug implements ports.Logger
func (l *LoggerAdapter) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	// For now, use standard log - can be improved later
	log.Printf("[DEBUG] %s %v", msg, keysAndValues)
}

// Info implements ports.Logger
func (l *LoggerAdapter) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s %v", msg, keysAndValues)
}

// Warn implements ports.Logger
func (l *LoggerAdapter) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.Printf("[WARN] %s %v", msg, keysAndValues)
}

// Error implements ports.Logger
func (l *LoggerAdapter) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, keysAndValues)
}
