package types

// LogLevel represents the severity level of HTTP logs
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
