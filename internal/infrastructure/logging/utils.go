package logging

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// TruncateID truncates long IDs/UUIDs for console display
// Example: "c6ceef12-3456-7890-abcd-1234567890ab" -> "c6ceef…7890ab"
func TruncateID(id string) string {
	if len(id) <= 12 {
		return id
	}

	// For UUIDs (with dashes)
	if strings.Contains(id, "-") {
		parts := strings.Split(id, "-")
		if len(parts) >= 2 {
			return fmt.Sprintf("%s…%s", parts[0][:6], parts[len(parts)-1][:6])
		}
	}

	// For other long strings
	if len(id) > 12 {
		return fmt.Sprintf("%s…%s", id[:6], id[len(id)-6:])
	}

	return id
}

// TruncateString truncates long strings with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// FormatContext creates a context string from key-value pairs
func FormatContext(pairs ...interface{}) string {
	if len(pairs)%2 != 0 {
		return ""
	}

	var parts []string
	for i := 0; i < len(pairs); i += 2 {
		key := fmt.Sprintf("%v", pairs[i])
		value := fmt.Sprintf("%v", pairs[i+1])
		
		// Truncate long values for readability
		if isIDField(key) && len(value) > 12 {
			value = TruncateID(value)
		} else if len(value) > 50 {
			value = TruncateString(value, 50)
		}
		
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	
	return strings.Join(parts, " ")
}

// SanitizeMessage cleans up log messages for consistency
func SanitizeMessage(msg string) string {
	// Remove excessive whitespace
	msg = strings.TrimSpace(msg)
	msg = strings.ReplaceAll(msg, "\n", " ")
	msg = strings.ReplaceAll(msg, "\t", " ")
	
	// Replace multiple spaces with single space
	for strings.Contains(msg, "  ") {
		msg = strings.ReplaceAll(msg, "  ", " ")
	}
	
	return msg
}

// GetShortMessage creates a short, objective message
func GetShortMessage(action, resource string, success bool) string {
	if success {
		return fmt.Sprintf("%s %s", action, resource)
	}
	return fmt.Sprintf("%s %s failed", action, resource)
}

// Common message templates for consistency
var MessageTemplates = struct {
	SessionConnected    string
	SessionDisconnected string
	SessionFailed       string
	MessageSent         string
	MessageReceived     string
	MessageFailed       string
	DatabaseConnected   string
	DatabaseError       string
	ServerStarted       string
	ServerStopped       string
	RequestProcessed    string
	RequestFailed       string
}{
	SessionConnected:    "Session connected",
	SessionDisconnected: "Session disconnected",
	SessionFailed:       "Session connection failed",
	MessageSent:         "Message sent",
	MessageReceived:     "Message received",
	MessageFailed:       "Message send failed",
	DatabaseConnected:   "Database connected",
	DatabaseError:       "Database error",
	ServerStarted:       "Server started",
	ServerStopped:       "Server stopped",
	RequestProcessed:    "Request processed",
	RequestFailed:       "Request failed",
}

// LogLevel helpers
func IsDebugEnabled(logger Logger) bool {
	// This is a simplified check - in a real implementation,
	// you'd check the actual log level
	return true
}

// Performance helpers
type ExecTimer struct {
	start time.Time
	name  string
}

// NewExecTimer creates a new timer for performance logging
func NewExecTimer(name string) *ExecTimer {
	return &ExecTimer{
		start: time.Now(),
		name:  name,
	}
}

// Stop stops the timer and returns the duration
func (t *ExecTimer) Stop() time.Duration {
	return time.Since(t.start)
}

// LogDuration logs the duration with context
func (t *ExecTimer) LogDuration(logger Logger, msg string) {
	duration := t.Stop()
	logger.With().
		Str("operation", t.name).
		Dur("duration", duration).
		Logger().Info(msg)
}

// Context builders for common scenarios
type LogContextBuilder struct {
	pairs []interface{}
}

// NewContext creates a new context builder
func NewContext() *LogContextBuilder {
	return &LogContextBuilder{
		pairs: make([]interface{}, 0),
	}
}

// Add adds a key-value pair to the context
func (c *LogContextBuilder) Add(key string, value interface{}) *LogContextBuilder {
	c.pairs = append(c.pairs, key, value)
	return c
}

// Session adds session context
func (c *LogContextBuilder) Session(sessionID string) *LogContextBuilder {
	return c.Add("session", sessionID)
}

// User adds user context
func (c *LogContextBuilder) User(userID string) *LogContextBuilder {
	return c.Add("user", userID)
}

// Request adds request context
func (c *LogContextBuilder) Request(method, path string) *LogContextBuilder {
	return c.Add("method", method).Add("path", path)
}

// Error adds error context
func (c *LogContextBuilder) Error(err error) *LogContextBuilder {
	if err != nil {
		return c.Add("error", err.Error())
	}
	return c
}

// Duration adds duration context
func (c *LogContextBuilder) Duration(d time.Duration) *LogContextBuilder {
	return c.Add("duration", FormatDuration(d))
}

// Build returns the context pairs
func (c *LogContextBuilder) Build() []interface{} {
	return c.pairs
}

// Apply applies the context to a logger
func (c *LogContextBuilder) Apply(logger Logger) LogContext {
	ctx := logger.With()
	for i := 0; i < len(c.pairs); i += 2 {
		key := fmt.Sprintf("%v", c.pairs[i])
		value := c.pairs[i+1]
		
		switch v := value.(type) {
		case string:
			ctx = ctx.Str(key, v)
		case int:
			ctx = ctx.Int(key, v)
		case bool:
			ctx = ctx.Bool(key, v)
		case time.Duration:
			ctx = ctx.Dur(key, v)
		case time.Time:
			ctx = ctx.Time(key, v)
		default:
			ctx = ctx.Str(key, fmt.Sprintf("%v", v))
		}
	}
	return ctx
}

// GenerateTraceID generates a unique trace ID for request tracking
func GenerateTraceID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateCorrelationID generates a unique correlation ID
func GenerateCorrelationID() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Helper functions for common logging patterns

// LogHTTPRequest logs an HTTP request with standard context
func LogHTTPRequest(logger Logger, method, path string, status int, duration time.Duration, clientIP string) {
	level := "info"
	if status >= 500 {
		level = "error"
	} else if status >= 400 {
		level = "warn"
	}

	ctx := NewContext().
		Request(method, path).
		Add("status", status).
		Duration(duration).
		Add("client_ip", clientIP).
		Apply(logger)

	switch level {
	case "error":
		ctx.Logger().Error("HTTP request")
	case "warn":
		ctx.Logger().Warn("HTTP request")
	default:
		ctx.Logger().Info("HTTP request")
	}
}

// LogSessionEvent logs a session-related event
func LogSessionEvent(logger Logger, sessionID, event string, success bool, err error) {
	ctx := NewContext().
		Session(sessionID).
		Add("event", event).
		Add("success", success)

	if err != nil {
		ctx = ctx.Error(err)
	}

	msg := GetShortMessage(event, "session", success)
	
	if success {
		ctx.Apply(logger).Logger().Info(msg)
	} else {
		ctx.Apply(logger).Logger().Error(msg)
	}
}

// LogMessageEvent logs a message-related event
func LogMessageEvent(logger Logger, sessionID, messageID, direction string, success bool, err error) {
	ctx := NewContext().
		Session(sessionID).
		Add("msg_id", messageID).
		Add("direction", direction).
		Add("success", success)

	if err != nil {
		ctx = ctx.Error(err)
	}

	msg := GetShortMessage(direction, "message", success)
	
	if success {
		ctx.Apply(logger).Logger().Info(msg)
	} else {
		ctx.Apply(logger).Logger().Error(msg)
	}
}
