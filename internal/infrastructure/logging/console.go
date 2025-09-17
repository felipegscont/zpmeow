package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// ANSI color codes for console output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorWhite  = "\033[97m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

// Level colors for the new format
var levelColors = map[string]string{
	"debug": ColorGray,
	"info":  ColorGreen,
	"warn":  ColorYellow,
	"error": ColorRed,
	"fatal": ColorRed + ColorBold,
	"panic": ColorRed + ColorBold,
}

// Level abbreviations
var levelAbbrev = map[string]string{
	"debug": "DBG",
	"info":  "INF",
	"warn":  "WRN",
	"error": "ERR",
	"fatal": "FTL",
	"panic": "PNC",
}

// ConsoleWriter provides a clean, robust console output
type ConsoleWriter struct {
	Out     io.Writer
	NoColor bool
}

// Write implements io.Writer interface
func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}

	// Parse JSON log entry
	if err := parseJSON(p, &event); err != nil {
		// Fallback: write raw if parsing fails
		_, writeErr := w.Out.Write(p)
		if writeErr != nil {
			return 0, writeErr
		}
		return len(p), nil
	}

	// Format for console in the requested format
	formatted := w.formatEvent(event)

	// Write formatted output
	_, writeErr := w.Out.Write([]byte(formatted))
	if writeErr != nil {
		return 0, writeErr
	}

	return len(p), nil // Return original length to satisfy zerolog
}

// formatEvent formats a log event in the requested format:
// 12:40PM INF Mensagem enviada com sucesso module=MeowService action=SendMessage to=559981276953@s.meow.net type=text
func (w *ConsoleWriter) formatEvent(event map[string]interface{}) string {
	var parts []string

	// 1. Timestamp with full date and time: 17-09-2025 13:25:30
	if timestamp, ok := event["time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
			timeStr := t.Format("02-01-2006 15:04:05")
			parts = append(parts, timeStr)
		}
	}

	// 2. Level abbreviation with color
	level := "info"
	if l, ok := event["level"].(string); ok {
		level = l
	}

	levelStr := levelAbbrev[level]
	if levelStr == "" {
		levelStr = strings.ToUpper(level)[:3]
	}

	if !w.NoColor {
		color := levelColors[level]
		levelStr = color + levelStr + ColorReset
	}
	parts = append(parts, levelStr)

	// 3. Message
	message := ""
	if msg, ok := event["message"].(string); ok {
		message = msg
	}
	parts = append(parts, message)

	// 4. Context fields in key=value format
	context := w.formatContext(event)
	if context != "" {
		parts = append(parts, context)
	}

	return strings.Join(parts, " ") + "\n"
}

// formatContext formats context fields as key=value pairs
func (w *ConsoleWriter) formatContext(event map[string]interface{}) string {
	var contextParts []string

	// Skip standard fields
	skipFields := map[string]bool{
		"time":    true,
		"level":   true,
		"message": true,
		"caller":  true,
	}

	// Define order for common fields
	fieldOrder := []string{"module", "action", "event", "userId", "jid", "to", "from", "type", "traceId", "correlationId", "plan", "ts", "payload", "reason", "error"}

	// Add ordered fields first
	for _, field := range fieldOrder {
		if value, exists := event[field]; exists {
			valueStr := w.formatValue(field, value)
			contextParts = append(contextParts, fmt.Sprintf("%s=%s", field, valueStr))
		}
	}

	// Add remaining fields
	for key, value := range event {
		if skipFields[key] {
			continue
		}

		// Skip if already added in ordered fields
		found := false
		for _, field := range fieldOrder {
			if key == field {
				found = true
				break
			}
		}
		if found {
			continue
		}

		valueStr := w.formatValue(key, value)
		contextParts = append(contextParts, fmt.Sprintf("%s=%s", key, valueStr))
	}

	return strings.Join(contextParts, " ")
}

// formatValue formats a value based on its type and key
func (w *ConsoleWriter) formatValue(key string, value interface{}) string {
	switch v := value.(type) {
	case string:
		// For payload field, only show in debug and truncate if too long
		if key == "payload" && len(v) > 100 {
			return fmt.Sprintf(`"%s..."`, v[:97])
		}
		// Quote strings that contain spaces or special characters
		if strings.Contains(v, " ") || strings.Contains(v, "=") {
			return fmt.Sprintf(`"%s"`, v)
		}
		return v
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// isIDField checks if a field name suggests it contains an ID/UUID
func isIDField(key string) bool {
	key = strings.ToLower(key)
	return strings.Contains(key, "id") ||
		strings.Contains(key, "uuid") ||
		strings.Contains(key, "hash") ||
		strings.Contains(key, "token") ||
		strings.Contains(key, "session")
}

// parseJSON parses JSON log events using encoding/json
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// FormatDuration formats duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return d.String()
}

// FormatBytes formats byte count in human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
