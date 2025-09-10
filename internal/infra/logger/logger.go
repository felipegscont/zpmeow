package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zpmeow/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Logger is our application logger interface
type Logger interface {
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
	With() LoggerContext
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	Sub(module string) Logger
}

// LoggerContext provides a fluent interface for adding context to logs
type LoggerContext interface {
	Str(key, val string) LoggerContext
	Int(key string, i int) LoggerContext
	Bool(key string, b bool) LoggerContext
	Err(err error) LoggerContext
	Dur(key string, d time.Duration) LoggerContext
	Time(key string, t time.Time) LoggerContext
	Interface(key string, i interface{}) LoggerContext
	Logger() Logger
}

// Config holds logger configuration (interface for external config)
type Config interface {
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

// zerologLogger implements our Logger interface using zerolog
type zerologLogger struct {
	logger zerolog.Logger
	module string
}

// zerologContext implements LoggerContext
type zerologContext struct {
	ctx zerolog.Context
}

// Initialize sets up the global logger with the given configuration
func Initialize(config Config) Logger {
	// Set global log level
	level, err := zerolog.ParseLevel(strings.ToLower(config.GetLevel()))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Create writers
	var writers []io.Writer

	// Console writer
	if config.GetFormat() == "console" {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
			NoColor:    !config.GetConsoleColor(),
		}

		// Custom level formatting with colors
		consoleWriter.FormatLevel = func(i interface{}) string {
			if i == nil {
				return ""
			}
			level := strings.ToUpper(i.(string))
			if !config.GetConsoleColor() {
				return level
			}
			
			switch level {
			case "DEBUG":
				return "\x1b[36m" + level + "\x1b[0m" // Cyan
			case "INFO":
				return "\x1b[32m" + level + "\x1b[0m" // Green
			case "WARN":
				return "\x1b[33m" + level + "\x1b[0m" // Yellow
			case "ERROR":
				return "\x1b[31m" + level + "\x1b[0m" // Red
			case "FATAL":
				return "\x1b[35m" + level + "\x1b[0m" // Magenta
			default:
				return level
			}
		}

		writers = append(writers, consoleWriter)
	} else {
		// JSON format to stdout
		writers = append(writers, os.Stdout)
	}

	// File writer
	if config.GetFileEnabled() {
		// Ensure log directory exists
		logDir := filepath.Dir(config.GetFilePath())
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
		} else {
			fileWriter := &lumberjack.Logger{
				Filename:   config.GetFilePath(),
				MaxSize:    config.GetFileMaxSize(),
				MaxBackups: config.GetFileMaxBackups(),
				MaxAge:     config.GetFileMaxAge(),
				Compress:   config.GetFileCompress(),
			}

			if config.GetFileFormat() == "console" {
				// Console format for file (without colors)
				consoleFileWriter := zerolog.ConsoleWriter{
					Out:        fileWriter,
					TimeFormat: "2006-01-02 15:04:05",
					NoColor:    true,
				}
				writers = append(writers, consoleFileWriter)
			} else {
				// JSON format for file
				writers = append(writers, fileWriter)
			}
		}
	}

	// Create multi-writer
	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = zerolog.MultiLevelWriter(writers...)
	}

	// Create logger
	logger := zerolog.New(writer).With().
		Timestamp().
		Caller().
		Logger()

	// Set as global logger
	log.Logger = logger

	return &zerologLogger{
		logger: logger,
		module: "app",
	}
}

// Debug logs a debug message
func (l *zerologLogger) Debug(msg string) {
	l.logger.Debug().Str("module", l.module).Msg(msg)
}

// Debugf logs a formatted debug message
func (l *zerologLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Str("module", l.module).Msgf(format, args...)
}

// Info logs an info message
func (l *zerologLogger) Info(msg string) {
	l.logger.Info().Str("module", l.module).Msg(msg)
}

// Infof logs a formatted info message
func (l *zerologLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Str("module", l.module).Msgf(format, args...)
}

// Warn logs a warning message
func (l *zerologLogger) Warn(msg string) {
	l.logger.Warn().Str("module", l.module).Msg(msg)
}

// Warnf logs a formatted warning message
func (l *zerologLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Str("module", l.module).Msgf(format, args...)
}

// Error logs an error message
func (l *zerologLogger) Error(msg string) {
	l.logger.Error().Str("module", l.module).Msg(msg)
}

// Errorf logs a formatted error message
func (l *zerologLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Str("module", l.module).Msgf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *zerologLogger) Fatal(msg string) {
	l.logger.Fatal().Str("module", l.module).Msg(msg)
}

// Fatalf logs a formatted fatal message and exits
func (l *zerologLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Str("module", l.module).Msgf(format, args...)
}

// With returns a logger context for adding fields
func (l *zerologLogger) With() LoggerContext {
	return &zerologContext{
		ctx: l.logger.With().Str("module", l.module),
	}
}

// WithField returns a logger with a single field
func (l *zerologLogger) WithField(key string, value interface{}) Logger {
	return &zerologLogger{
		logger: l.logger.With().Str("module", l.module).Interface(key, value).Logger(),
		module: l.module,
	}
}

// WithFields returns a logger with multiple fields
func (l *zerologLogger) WithFields(fields map[string]interface{}) Logger {
	ctx := l.logger.With().Str("module", l.module)
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &zerologLogger{
		logger: ctx.Logger(),
		module: l.module,
	}
}

// Sub creates a sub-logger with a module name
func (l *zerologLogger) Sub(module string) Logger {
	var fullModule string
	if l.module != "" {
		fullModule = fmt.Sprintf("%s/%s", l.module, module)
	} else {
		fullModule = module
	}
	
	return &zerologLogger{
		logger: l.logger,
		module: fullModule,
	}
}

// LoggerContext implementation
func (c *zerologContext) Str(key, val string) LoggerContext {
	return &zerologContext{ctx: c.ctx.Str(key, val)}
}

func (c *zerologContext) Int(key string, i int) LoggerContext {
	return &zerologContext{ctx: c.ctx.Int(key, i)}
}

func (c *zerologContext) Bool(key string, b bool) LoggerContext {
	return &zerologContext{ctx: c.ctx.Bool(key, b)}
}

func (c *zerologContext) Err(err error) LoggerContext {
	return &zerologContext{ctx: c.ctx.Err(err)}
}

func (c *zerologContext) Dur(key string, d time.Duration) LoggerContext {
	return &zerologContext{ctx: c.ctx.Dur(key, d)}
}

func (c *zerologContext) Time(key string, t time.Time) LoggerContext {
	return &zerologContext{ctx: c.ctx.Time(key, t)}
}

func (c *zerologContext) Interface(key string, i interface{}) LoggerContext {
	return &zerologContext{ctx: c.ctx.Interface(key, i)}
}

func (c *zerologContext) Logger() Logger {
	return &zerologLogger{
		logger: c.ctx.Logger(),
		module: "",
	}
}

// Global logger instance
var globalLogger Logger

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		globalLogger = Initialize(config.DefaultLoggerConfig())
	}
	return globalLogger
}



// SetLogger sets the global logger instance
func SetLogger(logger Logger) {
	globalLogger = logger
}

// waLogAdapter adapts our Logger interface to implement waLog.Logger
type waLogAdapter struct {
	logger Logger
}

// NewWALogAdapter creates a new adapter that implements waLog.Logger interface
func NewWALogAdapter(logger Logger) waLog.Logger {
	return &waLogAdapter{
		logger: logger,
	}
}

// Warnf implements waLog.Logger interface
func (w *waLogAdapter) Warnf(msg string, args ...interface{}) {
	w.logger.Warnf(msg, args...)
}

// Errorf implements waLog.Logger interface
func (w *waLogAdapter) Errorf(msg string, args ...interface{}) {
	w.logger.Errorf(msg, args...)
}

// Infof implements waLog.Logger interface
func (w *waLogAdapter) Infof(msg string, args ...interface{}) {
	w.logger.Infof(msg, args...)
}

// Debugf implements waLog.Logger interface
func (w *waLogAdapter) Debugf(msg string, args ...interface{}) {
	w.logger.Debugf(msg, args...)
}

// Sub implements waLog.Logger interface
func (w *waLogAdapter) Sub(module string) waLog.Logger {
	return &waLogAdapter{
		logger: w.logger.Sub(module),
	}
}

// GetWALogger returns a waLog.Logger that uses our logger system
func GetWALogger(module string) waLog.Logger {
	return NewWALogAdapter(GetLogger().Sub(module))
}
