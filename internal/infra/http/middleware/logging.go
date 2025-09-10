package middleware

import (
	"strings"

	"zpmeow/internal/infra/logger"

	"github.com/gin-gonic/gin"
)

// LogLevel represents the log level for HTTP requests
type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// HTTPLogEntry represents a structured HTTP log entry
type HTTPLogEntry struct {
	Method    string
	Path      string
	Status    int
	Latency   string
	ClientIP  string
	UserAgent string
	Error     string
	Level     LogLevel
}

// Logger returns a gin.HandlerFunc (middleware) that logs requests using our logger
func Logger() gin.HandlerFunc {
	httpLogger := logger.GetLogger().Sub("http")

	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		// Skip logging for certain paths to reduce noise
		if shouldSkipLogging(params.Path) {
			return ""
		}

		// Create and log the HTTP entry
		entry := createHTTPLogEntry(params)
		logHTTPRequest(httpLogger, entry)

		// Return empty string since we're handling logging ourselves
		return ""
	})
}

// createHTTPLogEntry creates a structured log entry from gin parameters
func createHTTPLogEntry(params gin.LogFormatterParams) HTTPLogEntry {
	entry := HTTPLogEntry{
		Method:   params.Method,
		Path:     params.Path,
		Status:   params.StatusCode,
		Latency:  params.Latency.String(),
		ClientIP: params.ClientIP,
		Level:    determineLogLevel(params.StatusCode),
	}

	// Add error information if present
	if params.ErrorMessage != "" {
		entry.Error = params.ErrorMessage
	}

	// Add user agent for non-static requests
	if !isStaticResource(params.Path) && params.Request != nil {
		entry.UserAgent = params.Request.UserAgent()
	}

	return entry
}

// logHTTPRequest logs the HTTP request with the appropriate level
func logHTTPRequest(httpLogger logger.Logger, entry HTTPLogEntry) {
	logEntry := httpLogger.With().
		Str("method", entry.Method).
		Str("path", entry.Path).
		Int("status", entry.Status).
		Str("latency", entry.Latency).
		Str("client_ip", entry.ClientIP)

	// Add optional fields
	if entry.Error != "" {
		logEntry = logEntry.Str("error", entry.Error)
	}
	if entry.UserAgent != "" {
		logEntry = logEntry.Str("user_agent", entry.UserAgent)
	}

	// Log with appropriate level
	switch entry.Level {
	case LogLevelError:
		logEntry.Logger().Error("HTTP Request")
	case LogLevelWarn:
		logEntry.Logger().Warn("HTTP Request")
	default:
		logEntry.Logger().Info("HTTP Request")
	}
}

// Paths to skip logging for noise reduction
var skipLogPaths = map[string]bool{
	"/ping":        true,
	"/health":      true,
	"/favicon.ico": true,
}

// shouldSkipLogging determines if we should skip logging for certain paths
func shouldSkipLogging(path string) bool {
	return skipLogPaths[path]
}

// determineLogLevel determines the appropriate log level based on HTTP status code
func determineLogLevel(statusCode int) LogLevel {
	switch {
	case statusCode >= 500:
		return LogLevelError
	case statusCode >= 400:
		return LogLevelWarn
	default:
		return LogLevelInfo
	}
}

// Static resource configuration
var (
	staticPrefixes = []string{
		"/swagger/",
		"/static/",
		"/assets/",
	}

	staticExtensions = map[string]bool{
		".css":   true,
		".js":    true,
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".ico":   true,
		".svg":   true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
	}
)

// isStaticResource checks if the path is for a static resource
func isStaticResource(path string) bool {
	// Check prefixes first (most common case)
	for _, prefix := range staticPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Check extensions using map lookup for O(1) performance
	for ext := range staticExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
