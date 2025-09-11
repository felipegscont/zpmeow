package logger

import (
	"testing"

	"zpmeow/internal/config"
)

func TestLoggerWithConfig(t *testing.T) {

	defaultConfig := config.DefaultLoggerConfig()
	log := Initialize(defaultConfig)
	
	if log == nil {
		t.Fatal("Logger should not be nil")
	}
	

	log.Info("Test info message")
	log.Debug("Test debug message")
	log.Warn("Test warning message")
	log.Error("Test error message")
}

func TestLoggerWithCustomConfig(t *testing.T) {

	customConfig := &config.LoggerConfig{
		Level:           "debug",
		Format:          "console",
		ConsoleColor:    false,
		FileEnabled:     false,
		FilePath:        "test.log",
		FileMaxSize:     50,
		FileMaxBackups:  2,
		FileMaxAge:      14,
		FileCompress:    false,
		FileFormat:      "json",
	}
	
	log := Initialize(customConfig)
	
	if log == nil {
		t.Fatal("Logger should not be nil")
	}
	

	log.With().
		Str("test_key", "test_value").
		Int("test_number", 42).
		Logger().
		Info("Test structured message")
}

func TestSubLoggers(t *testing.T) {
	config := config.DefaultLoggerConfig()
	log := Initialize(config)
	

	httpLogger := log.Sub("http")
	dbLogger := log.Sub("database")
	
	httpLogger.Info("HTTP request received")
	dbLogger.Error("Database connection failed")
	

	userLogger := httpLogger.Sub("user")
	userLogger.Info("User authenticated")
}

func TestWALogAdapter(t *testing.T) {
	config := config.DefaultLoggerConfig()
	Initialize(config)


	waLogger := GetWALogger("whatsapp")
	
	waLogger.Infof("WhatsApp client initialized")
	waLogger.Debugf("Debug message with value: %v", 42)
	waLogger.Warnf("Warning message")
	waLogger.Errorf("Error message")
	

	subLogger := waLogger.Sub("session")
	subLogger.Infof("Session created")
}

func TestConfigInterface(t *testing.T) {

	var _ Config = (*config.LoggerConfig)(nil)
	

	cfg := config.DefaultLoggerConfig()
	
	if cfg.GetLevel() != "info" {
		t.Errorf("Expected level 'info', got '%s'", cfg.GetLevel())
	}
	
	if cfg.GetFormat() != "console" {
		t.Errorf("Expected format 'console', got '%s'", cfg.GetFormat())
	}
	
	if !cfg.GetConsoleColor() {
		t.Error("Expected console color to be true")
	}
	
	if !cfg.GetFileEnabled() {
		t.Error("Expected file enabled to be true")
	}
}
