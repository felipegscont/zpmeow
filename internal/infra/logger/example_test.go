package logger

import (
	"os"
	"testing"
	"time"
)

func TestLoggerBasicUsage(t *testing.T) {
	// Test with custom config
	config := NewConfigAdapter(
		"info", "console", "log/test.log", "json",
		true, false, true, 100, 3, 28,
	)

	log := Initialize(config)
	
	// Test basic logging
	log.Info("Test info message")
	log.Warn("Test warning message")
	log.Error("Test error message")
	
	// Test formatted logging
	log.Infof("Test formatted message: %s", "hello")
	log.Errorf("Test error with value: %d", 42)
}

func TestSubLoggers(t *testing.T) {
	config := NewConfigAdapter(
		"info", "console", "log/test.log", "json",
		true, false, true, 100, 3, 28,
	)

	log := Initialize(config)
	
	// Create sub-loggers
	httpLogger := log.Sub("http")
	dbLogger := log.Sub("database")
	
	httpLogger.Info("HTTP request received")
	dbLogger.Error("Database connection failed")
	
	// Nested sub-loggers
	userLogger := httpLogger.Sub("user")
	userLogger.Info("User authenticated")
}

func TestStructuredLogging(t *testing.T) {
	config := NewConfigAdapter(
		"info", "console", "log/test.log", "json",
		true, false, true, 100, 3, 28,
	)

	log := Initialize(config)
	
	// Test structured logging
	log.With().
		Str("user_id", "123").
		Int("status_code", 200).
		Dur("duration", time.Millisecond*150).
		Bool("success", true).
		Logger().
		Info("Request processed")
	
	// Test with fields
	fields := map[string]interface{}{
		"action":    "login",
		"user_id":   "456",
		"timestamp": time.Now(),
	}
	
	log.WithFields(fields).Info("User action")
}

func TestWALogAdapter(t *testing.T) {
	config := NewConfigAdapter(
		"info", "console", "log/test.log", "json",
		true, false, true, 100, 3, 28,
	)

	Initialize(config)

	// Test waLog adapter
	waLogger := GetWALogger("whatsapp")
	
	waLogger.Infof("WhatsApp client initialized")
	waLogger.Debugf("Debug message with value: %v", 42)
	waLogger.Warnf("Warning message")
	waLogger.Errorf("Error message")
	
	// Test sub-logger
	subLogger := waLogger.Sub("session")
	subLogger.Infof("Session created")
}

func TestFileLogging(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()

	config := NewConfigAdapter(
		"debug", "console", tempDir+"/test.log", "json",
		false, true, false, 1, 2, 1,
	)

	log := Initialize(config)
	
	// Write some logs
	log.Info("Test file logging")
	log.Error("Test error in file")
	
	// Check if file was created
	filePath := tempDir + "/test.log"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", filePath)
	}
}

func ExampleLogger() {
	// Initialize logger with default config
	log := Initialize(DefaultConfig())
	
	// Basic logging
	log.Info("Application started")
	log.Errorf("Failed to connect: %v", "connection refused")
	
	// Structured logging
	log.With().
		Str("user_id", "123").
		Int("attempts", 3).
		Logger().
		Warn("Login failed")
	
	// Sub-loggers
	httpLogger := log.Sub("http")
	httpLogger.Info("Server listening on :8080")
	
	// waLog adapter for whatsmeow
	waLogger := GetWALogger("whatsapp")
	waLogger.Infof("WhatsApp client ready")
}
