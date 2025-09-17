package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	waLog "go.mau.fi/whatsmeow/util/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger interface for structured logging
type Logger interface {
	// Core logging methods
	Debug(msg string)
	Debugf(format string, args ...interface{})
	Info(msg string)
	Infof(format string, args ...interface{})
	Warn(msg string)
	Warnf(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
	Fatal(msg string)
	Fatalf(format string, args ...interface{})

	// Structured logging with context
	With() LogContext
	Sub(module string) Logger
}

// LogContext interface for building structured log entries
type LogContext interface {
	Str(key, val string) LogContext
	Int(key string, val int) LogContext
	Bool(key string, val bool) LogContext
	Dur(key string, val time.Duration) LogContext
	Time(key string, val time.Time) LogContext
	Err(err error) LogContext
	Logger() Logger
}

// ZerologLogger implements Logger interface using zerolog
type ZerologLogger struct {
	logger zerolog.Logger
	module string
}

// ZerologContext implements LogContext interface
type ZerologContext struct {
	context zerolog.Context
	base    *ZerologLogger
}

// Global logger instance
var globalLogger Logger

// LoggerConfig holds configuration for the logger
type LoggerConfig interface {
	GetLevel() string
	GetFormat() string
	GetConsoleColor() bool
	GetFileEnabled() bool
	GetFilePath() string
	GetFileMaxSize() int
	GetFileMaxBackups() int
	GetFileMaxAge() int
	GetFileCompress() bool
	GetFileFormat() string
}

// Initialize creates and configures the global logger
func Initialize(config LoggerConfig) Logger {
	// Set global log level
	level := parseLogLevel(config.GetLevel())
	zerolog.SetGlobalLevel(level)

	// Configure time format
	zerolog.TimeFieldFormat = time.RFC3339

	var writers []io.Writer

	// Console output with TTY detection
	if config.GetFormat() == "console" || config.GetFormat() == "" {
		// Detect TTY and configure output writer
		var out io.Writer = os.Stdout

		// Windows compatibility
		if runtime.GOOS == "windows" {
			out = colorable.NewColorableStdout()
		}

		// Detect if we should use colors
		useColor := shouldUseColor(out, config.GetConsoleColor())

		consoleWriter := zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: "02-01-2006 15:04:05",
			NoColor:    !useColor,
		}
		writers = append(writers, consoleWriter)
	}

	// File output with JSON format
	if config.GetFileEnabled() {
		// Ensure log directory exists
		logDir := filepath.Dir(config.GetFilePath())
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
		}

		fileWriter := &lumberjack.Logger{
			Filename:   config.GetFilePath(),
			MaxSize:    config.GetFileMaxSize(),
			MaxBackups: config.GetFileMaxBackups(),
			MaxAge:     config.GetFileMaxAge(),
			Compress:   config.GetFileCompress(),
		}
		writers = append(writers, fileWriter)
	}

	// Create multi-writer for dual output
	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else {
		writer = os.Stdout
	}

	// Create logger with context
	logger := zerolog.New(writer).With().
		Timestamp().
		Logger()

	globalLogger = &ZerologLogger{
		logger: logger,
		module: "main",
	}

	return globalLogger
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Fallback to default logger
		globalLogger = &ZerologLogger{
			logger: log.Logger,
			module: "default",
		}
	}
	return globalLogger
}

// SetLogger sets the global logger instance
func SetLogger(logger Logger) {
	globalLogger = logger
}

// GetWALogger creates a meow-compatible logger
func GetWALogger(module string) waLog.Logger {
	return NewWALogger(module)
}

// shouldUseColor determines if colors should be used based on TTY detection and FORCE_COLOR
func shouldUseColor(out io.Writer, configColor bool) bool {
	// Check FORCE_COLOR environment variable first
	if forceColor := os.Getenv("FORCE_COLOR"); forceColor != "" {
		return forceColor != "0" && forceColor != "false"
	}

	// If config explicitly disables color, respect it
	if !configColor {
		return false
	}

	// Check if output is a TTY
	if f, ok := out.(*os.File); ok {
		return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}

	// For colorable writer (Windows), check the underlying file
	if runtime.GOOS == "windows" {
		// Try to get the underlying file descriptor
		return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	}

	return false
}

// parseLogLevel converts string level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// ZerologLogger implementation
func (l *ZerologLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *ZerologLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *ZerologLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *ZerologLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l *ZerologLogger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l *ZerologLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l *ZerologLogger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

func (l *ZerologLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

func (l *ZerologLogger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

func (l *ZerologLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

func (l *ZerologLogger) With() LogContext {
	ctx := l.logger.With()
	return &ZerologContext{
		context: ctx,
		base:    l,
	}
}

func (l *ZerologLogger) Sub(module string) Logger {
	fullModule := l.module
	if module != "" {
		if l.module != "" {
			fullModule = l.module + "." + module
		} else {
			fullModule = module
		}
	}
	return &ZerologLogger{
		logger: l.logger,
		module: fullModule,
	}
}

// ZerologContext implementation
func (c *ZerologContext) Str(key, val string) LogContext {
	// Truncate long values for console readability
	if len(val) > 50 && (strings.Contains(key, "id") || strings.Contains(key, "uuid") || strings.Contains(key, "hash")) {
		val = TruncateID(val)
	}
	c.context = c.context.Str(key, val)
	return c
}

func (c *ZerologContext) Int(key string, val int) LogContext {
	c.context = c.context.Int(key, val)
	return c
}

func (c *ZerologContext) Bool(key string, val bool) LogContext {
	c.context = c.context.Bool(key, val)
	return c
}

func (c *ZerologContext) Dur(key string, val time.Duration) LogContext {
	c.context = c.context.Dur(key, val)
	return c
}

func (c *ZerologContext) Time(key string, val time.Time) LogContext {
	c.context = c.context.Time(key, val)
	return c
}

func (c *ZerologContext) Err(err error) LogContext {
	c.context = c.context.Err(err)
	return c
}

func (c *ZerologContext) Logger() Logger {
	return &ZerologLogger{
		logger: c.context.Logger(),
		module: c.base.module,
	}
}
