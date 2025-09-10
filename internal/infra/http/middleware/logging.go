package middleware

import (
	"strings"

	"zpmeow/internal/infra/logger"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin.HandlerFunc (middleware) that logs requests using our logger
func Logger() gin.HandlerFunc {
	log := logger.GetLogger().Sub("http")

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Skip logging for certain paths to reduce noise
		if shouldSkipLogging(param.Path) {
			return ""
		}

		// Determine log level based on status code
		logLevel := getLogLevel(param.StatusCode)

		// Create structured log entry
		logEntry := log.With().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Dur("latency", param.Latency).
			Str("client_ip", param.ClientIP).
			Time("timestamp", param.TimeStamp)

		// Add error information if present
		if param.ErrorMessage != "" {
			logEntry = logEntry.Str("error", param.ErrorMessage)
		}

		// Add user agent for non-static requests
		if !isStaticResource(param.Path) && param.Request != nil {
			logEntry = logEntry.Str("user_agent", param.Request.UserAgent())
		}

		// Log with appropriate level
		switch logLevel {
		case "error":
			logEntry.Logger().Error("HTTP Request")
		case "warn":
			logEntry.Logger().Warn("HTTP Request")
		default:
			logEntry.Logger().Info("HTTP Request")
		}

		// Return empty string since we're handling logging ourselves
		return ""
	})
}

// shouldSkipLogging determines if we should skip logging for certain paths
func shouldSkipLogging(path string) bool {
	skipPaths := []string{
		"/ping",
		"/health",
		"/favicon.ico",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}

	return false
}

// getLogLevel determines the appropriate log level based on HTTP status code
func getLogLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warn"
	default:
		return "info"
	}
}

// isStaticResource checks if the path is for a static resource
func isStaticResource(path string) bool {
	staticPrefixes := []string{
		"/swagger/",
		"/static/",
		"/assets/",
	}

	staticExtensions := []string{
		".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg",
		".woff", ".woff2", ".ttf", ".eot",
	}

	// Check prefixes
	for _, prefix := range staticPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Check extensions
	for _, ext := range staticExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
