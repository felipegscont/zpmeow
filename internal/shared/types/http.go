package types

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

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
