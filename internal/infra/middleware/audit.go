package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"meow/internal/infra/logging"
)

// AuditEvent represents an audit log entry
type AuditEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	RequestID   string                 `json:"request_id"`
	UserID      string                 `json:"user_id,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	StatusCode  int                    `json:"status_code"`
	ClientIP    string                 `json:"client_ip"`
	UserAgent   string                 `json:"user_agent"`
	RequestBody map[string]interface{} `json:"request_body,omitempty"`
	Duration    int64                  `json:"duration_ms"`
	Success     bool                   `json:"success"`
	Error       string                 `json:"error,omitempty"`
}

// AuditMiddleware logs important actions for security and compliance
func AuditMiddleware() gin.HandlerFunc {
	logger := logging.GetLogger().Sub("audit")

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Skip audit for health checks and static content
		if shouldSkipAudit(path) {
			c.Next()
			return
		}

		// Capture request body for important endpoints
		var requestBody map[string]interface{}
		if shouldCaptureBody(path, method) {
			requestBody = captureRequestBody(c)
		}

		// Process request
		c.Next()

		// Create audit event
		event := AuditEvent{
			Timestamp:   start,
			RequestID:   getRequestID(c),
			UserID:      getUserID(c),
			SessionID:   getSessionID(c),
			Action:      determineAction(method, path),
			Resource:    determineResource(path),
			Method:      method,
			Path:        path,
			StatusCode:  c.Writer.Status(),
			ClientIP:    c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			RequestBody: requestBody,
			Duration:    time.Since(start).Milliseconds(),
			Success:     c.Writer.Status() < 400,
		}

		// Add error information if request failed
		if !event.Success {
			if err, exists := c.Get("error"); exists {
				if errStr, ok := err.(string); ok {
					event.Error = errStr
				}
			}
		}

		// Log audit event
		logAuditEvent(logger, &event)
	}
}

// shouldSkipAudit determines if the request should be audited
func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/ping",
		"/health",
		"/metrics",
		"/favicon.ico",
		"/static/",
		"/swagger/",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

// shouldCaptureBody determines if request body should be captured
func shouldCaptureBody(path, method string) bool {
	// Only capture body for POST, PUT, PATCH requests
	if method != "POST" && method != "PUT" && method != "PATCH" {
		return false
	}

	// Capture for important endpoints
	importantPaths := []string{
		"/sessions",
		"/session/",
		"/api/",
	}

	for _, importantPath := range importantPaths {
		if strings.HasPrefix(path, importantPath) {
			return true
		}
	}

	return false
}

// captureRequestBody safely captures and parses request body
func captureRequestBody(c *gin.Context) map[string]interface{} {
	if c.Request.Body == nil {
		return nil
	}

	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}

	// Restore body for subsequent handlers
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Parse JSON
	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		return nil
	}

	// Remove sensitive fields
	sanitizeRequestBody(requestBody)

	return requestBody
}

// sanitizeRequestBody removes sensitive information from request body
func sanitizeRequestBody(body map[string]interface{}) {
	sensitiveFields := []string{
		"password",
		"token",
		"api_key",
		"secret",
		"private_key",
		"auth",
		"authorization",
	}

	for _, field := range sensitiveFields {
		if _, exists := body[field]; exists {
			body[field] = "[REDACTED]"
		}
	}
}

// getUserID extracts user ID from context
func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// getSessionID extracts session ID from context or URL
func getSessionID(c *gin.Context) string {
	// Try to get from context first
	if sessionID, exists := c.Get("session_id"); exists {
		if id, ok := sessionID.(string); ok {
			return id
		}
	}

	// Try to get from URL parameters
	if sessionID := c.Param("sessionId"); sessionID != "" {
		return sessionID
	}

	if sessionID := c.Param("id"); sessionID != "" {
		return sessionID
	}

	return ""
}

// determineAction determines the action being performed
func determineAction(method, path string) string {
	switch method {
	case "GET":
		if strings.Contains(path, "/list") {
			return "list"
		}
		return "read"
	case "POST":
		if strings.Contains(path, "/create") {
			return "create"
		}
		if strings.Contains(path, "/send") {
			return "send"
		}
		if strings.Contains(path, "/connect") {
			return "connect"
		}
		if strings.Contains(path, "/disconnect") {
			return "disconnect"
		}
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return method
	}
}

// determineResource determines the resource being accessed
func determineResource(path string) string {
	if strings.Contains(path, "/sessions") {
		return "session"
	}
	if strings.Contains(path, "/message") {
		return "message"
	}
	if strings.Contains(path, "/contacts") {
		return "contact"
	}
	if strings.Contains(path, "/groups") {
		return "group"
	}
	if strings.Contains(path, "/webhook") {
		return "webhook"
	}
	if strings.Contains(path, "/privacy") {
		return "privacy"
	}

	// Extract resource from path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return parts[0]
	}

	return "unknown"
}

// logAuditEvent logs the audit event with appropriate level
func logAuditEvent(logger logging.Logger, event *AuditEvent) {
	eventJSON, _ := json.Marshal(event)

	if event.Success {
		logger.Infof("AUDIT: %s", string(eventJSON))
	} else {
		logger.Warnf("AUDIT_FAILED: %s", string(eventJSON))
	}

	// Log critical actions at higher level
	if isCriticalAction(event.Action, event.Resource) {
		logger.Errorf("AUDIT_CRITICAL: %s", string(eventJSON))
	}
}

// isCriticalAction determines if an action is critical and needs special attention
func isCriticalAction(action, resource string) bool {
	criticalActions := map[string][]string{
		"session": {"create", "delete"},
		"webhook": {"update", "delete"},
		"privacy": {"update"},
	}

	if actions, exists := criticalActions[resource]; exists {
		for _, criticalAction := range actions {
			if action == criticalAction {
				return true
			}
		}
	}

	return false
}
