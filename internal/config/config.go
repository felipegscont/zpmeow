package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Auth     AuthConfig     `json:"auth"`
	Logging  LoggingConfig  `json:"logging"`
	CORS     CORSConfig     `json:"cors"`
	Webhook  WebhookConfig  `json:"webhook"`
	Meow     MeowConfig     `json:"meow"`
	Security SecurityConfig `json:"security"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	Name            string        `json:"name"`
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	URL             string        `json:"url"` // Computed field
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string        `json:"port"`
	Mode         string        `json:"mode"` // debug, release, test
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	GlobalAPIKey    string        `json:"global_api_key"`
	SessionTimeout  time.Duration `json:"session_timeout"`
	TokenExpiration time.Duration `json:"token_expiration"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level          string `json:"level"`
	Format         string `json:"format"`
	ConsoleColor   bool   `json:"console_color"`
	FileEnabled    bool   `json:"file_enabled"`
	FilePath       string `json:"file_path"`
	FileMaxSize    int    `json:"file_max_size"`
	FileMaxBackups int    `json:"file_max_backups"`
	FileMaxAge     int    `json:"file_max_age"`
	FileCompress   bool   `json:"file_compress"`
	FileFormat     string `json:"file_format"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowAllOrigins  bool     `json:"allow_all_origins"`
	AllowOrigins     []string `json:"allow_origins"`
	AllowMethods     []string `json:"allow_methods"`
	AllowHeaders     []string `json:"allow_headers"`
	ExposeHeaders    []string `json:"expose_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
}

// WebhookConfig holds webhook configuration
type WebhookConfig struct {
	Timeout           time.Duration `json:"timeout"`
	MaxRetries        int           `json:"max_retries"`
	InitialBackoff    time.Duration `json:"initial_backoff"`
	MaxBackoff        time.Duration `json:"max_backoff"`
	BackoffMultiplier float64       `json:"backoff_multiplier"`
}

// MeowConfig holds meow client configuration
type MeowConfig struct {
	MaxRetries        int           `json:"max_retries"`
	RetryInterval     time.Duration `json:"retry_interval"`
	ConnectionTimeout time.Duration `json:"connection_timeout"`
	QRCodeTimeout     time.Duration `json:"qr_code_timeout"`
	ReconnectDelay    time.Duration `json:"reconnect_delay"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	RateLimitEnabled bool          `json:"rate_limit_enabled"`
	RateLimitRPS     int           `json:"rate_limit_rps"`
	RequestTimeout   time.Duration `json:"request_timeout"`
	MaxRequestSize   int64         `json:"max_request_size"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Database: loadDatabaseConfig(),
		Server:   loadServerConfig(),
		Auth:     loadAuthConfig(),
		Logging:  loadLoggingConfig(),
		CORS:     loadCORSConfig(),
		Webhook:  loadWebhookConfig(),
		Meow:     loadMeowConfig(),
		Security: loadSecurityConfig(),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Auth.GlobalAPIKey == "" {
		return fmt.Errorf("global API key is required")
	}
	return nil
}

// GetDatabaseURL returns the computed database URL
func (c *Config) GetDatabaseURL() string {
	return c.Database.URL
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Mode == "release"
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Mode == "debug"
}

// IsTest returns true if running in test mode
func (c *Config) IsTest() bool {
	return c.Server.Mode == "test"
}

// loadDatabaseConfig loads database configuration from environment
func loadDatabaseConfig() DatabaseConfig {
	cfg := DatabaseConfig{
		Host:            getEnvOrDefault("DB_HOST", "localhost"),
		Port:            getEnvOrDefault("DB_PORT", "5432"),
		User:            getEnvOrDefault("DB_USER", "postgres"),
		Password:        getEnvOrDefault("DB_PASSWORD", ""),
		Name:            getEnvOrDefault("DB_NAME", "meow"),
		SSLMode:         getEnvOrDefault("DB_SSLMODE", "disable"),
		MaxOpenConns:    getIntEnvOrDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getIntEnvOrDefault("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getDurationEnvOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}

	// Compute database URL
	cfg.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	return cfg
}

// loadServerConfig loads server configuration from environment
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:         getEnvOrDefault("SERVER_PORT", "8080"),
		Mode:         getEnvOrDefault("GIN_MODE", "debug"),
		ReadTimeout:  getDurationEnvOrDefault("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:  getDurationEnvOrDefault("SERVER_IDLE_TIMEOUT", 120*time.Second),
	}
}

// loadAuthConfig loads authentication configuration from environment
func loadAuthConfig() AuthConfig {
	return AuthConfig{
		GlobalAPIKey:    getEnvOrDefault("GLOBAL_API_KEY", ""),
		SessionTimeout:  getDurationEnvOrDefault("AUTH_SESSION_TIMEOUT", 24*time.Hour),
		TokenExpiration: getDurationEnvOrDefault("AUTH_TOKEN_EXPIRATION", 1*time.Hour),
	}
}

// loadLoggingConfig loads logging configuration from environment
func loadLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:          getEnvOrDefault("LOG_LEVEL", "info"),
		Format:         getEnvOrDefault("LOG_FORMAT", "console"),
		ConsoleColor:   getBoolEnvOrDefault("LOG_CONSOLE_COLOR", true),
		FileEnabled:    getBoolEnvOrDefault("LOG_FILE_ENABLED", true),
		FilePath:       getEnvOrDefault("LOG_FILE_PATH", "log/app.log"),
		FileMaxSize:    getIntEnvOrDefault("LOG_FILE_MAX_SIZE", 100),
		FileMaxBackups: getIntEnvOrDefault("LOG_FILE_MAX_BACKUPS", 3),
		FileMaxAge:     getIntEnvOrDefault("LOG_FILE_MAX_AGE", 28),
		FileCompress:   getBoolEnvOrDefault("LOG_FILE_COMPRESS", true),
		FileFormat:     getEnvOrDefault("LOG_FILE_FORMAT", "json"),
	}
}

// loadCORSConfig loads CORS configuration from environment
func loadCORSConfig() CORSConfig {
	return CORSConfig{
		AllowAllOrigins:  getBoolEnvOrDefault("CORS_ALLOW_ALL_ORIGINS", true),
		AllowOrigins:     getStringSliceEnvOrDefault("CORS_ALLOW_ORIGINS", []string{}),
		AllowMethods:     getStringSliceEnvOrDefault("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		AllowHeaders:     getStringSliceEnvOrDefault("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-API-Key"}),
		ExposeHeaders:    getStringSliceEnvOrDefault("CORS_EXPOSE_HEADERS", []string{}),
		AllowCredentials: getBoolEnvOrDefault("CORS_ALLOW_CREDENTIALS", false),
		MaxAge:           getIntEnvOrDefault("CORS_MAX_AGE", 86400),
	}
}

// loadWebhookConfig loads webhook configuration from environment
func loadWebhookConfig() WebhookConfig {
	return WebhookConfig{
		Timeout:           getDurationEnvOrDefault("WEBHOOK_TIMEOUT", 30*time.Second),
		MaxRetries:        getIntEnvOrDefault("WEBHOOK_MAX_RETRIES", 3),
		InitialBackoff:    getDurationEnvOrDefault("WEBHOOK_INITIAL_BACKOFF", 1*time.Second),
		MaxBackoff:        getDurationEnvOrDefault("WEBHOOK_MAX_BACKOFF", 30*time.Second),
		BackoffMultiplier: getFloat64EnvOrDefault("WEBHOOK_BACKOFF_MULTIPLIER", 2.0),
	}
}

// loadMeowConfig loads meow configuration from environment
func loadMeowConfig() MeowConfig {
	return MeowConfig{
		MaxRetries:        getIntEnvOrDefault("meow_MAX_RETRIES", 3),
		RetryInterval:     getDurationEnvOrDefault("meow_RETRY_INTERVAL", 5*time.Second),
		ConnectionTimeout: getDurationEnvOrDefault("meow_CONNECTION_TIMEOUT", 30*time.Second),
		QRCodeTimeout:     getDurationEnvOrDefault("meow_QR_CODE_TIMEOUT", 60*time.Second),
		ReconnectDelay:    getDurationEnvOrDefault("meow_RECONNECT_DELAY", 10*time.Second),
	}
}

// loadSecurityConfig loads security configuration from environment
func loadSecurityConfig() SecurityConfig {
	return SecurityConfig{
		RateLimitEnabled: getBoolEnvOrDefault("SECURITY_RATE_LIMIT_ENABLED", false),
		RateLimitRPS:     getIntEnvOrDefault("SECURITY_RATE_LIMIT_RPS", 100),
		RequestTimeout:   getDurationEnvOrDefault("SECURITY_REQUEST_TIMEOUT", 30*time.Second),
		MaxRequestSize:   getInt64EnvOrDefault("SECURITY_MAX_REQUEST_SIZE", 10*1024*1024), // 10MB
	}
}

// Helper functions for environment variable parsing

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnvOrDefault returns environment variable as bool or default
func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getIntEnvOrDefault returns environment variable as int or default
func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getInt64EnvOrDefault returns environment variable as int64 or default
func getInt64EnvOrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getFloat64EnvOrDefault returns environment variable as float64 or default
func getFloat64EnvOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getDurationEnvOrDefault returns environment variable as duration or default
func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getStringSliceEnvOrDefault returns environment variable as string slice or default
func getStringSliceEnvOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
