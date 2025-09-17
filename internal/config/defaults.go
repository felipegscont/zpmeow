package config

import "time"

// DefaultConfig returns a configuration with all default values
func DefaultConfig() *Config {
	return &Config{
		Database: DefaultDatabaseConfig(),
		Server:   DefaultServerConfig(),
		Auth:     DefaultAuthConfig(),
		Logging:  DefaultLoggingConfig(),
		CORS:     DefaultCORSConfig(),
		Webhook:  DefaultWebhookConfig(),
		Meow:     DefaultMeowConfig(),
		Security: DefaultSecurityConfig(),
	}
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            "localhost",
		Port:            "5432",
		User:            "postgres",
		Password:        "",
		Name:            "meow",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:         "8080",
		Mode:         "debug",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// DefaultAuthConfig returns default authentication configuration
func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		GlobalAPIKey:    "",
		SessionTimeout:  24 * time.Hour,
		TokenExpiration: 1 * time.Hour,
	}
}

// DefaultLoggingConfig returns default logging configuration
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:          "info",
		Format:         "console",
		ConsoleColor:   true,
		FileEnabled:    false, // Disabled by default, enable via LOG_FILE_ENABLED=true
		FilePath:       "log/app.log",
		FileMaxSize:    100,
		FileMaxBackups: 3,
		FileMaxAge:     28,
		FileCompress:   true,
		FileFormat:     "json",
	}
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowAllOrigins: true,
		AllowOrigins:    []string{},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "Authorization", "X-API-Key",
		},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// DefaultWebhookConfig returns default webhook configuration
func DefaultWebhookConfig() WebhookConfig {
	return WebhookConfig{
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		InitialBackoff:    1 * time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// DefaultMeowConfig returns default meow configuration
func DefaultMeowConfig() MeowConfig {
	return MeowConfig{
		MaxRetries:        3,
		RetryInterval:     5 * time.Second,
		ConnectionTimeout: 30 * time.Second,
		QRCodeTimeout:     60 * time.Second,
		ReconnectDelay:    10 * time.Second,
	}
}

// DefaultSecurityConfig returns default security configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		RateLimitEnabled: false,
		RateLimitRPS:     100,
		RequestTimeout:   30 * time.Second,
		MaxRequestSize:   10 * 1024 * 1024, // 10MB
	}
}

// ProductionConfig returns a configuration optimized for production
func ProductionConfig() *Config {
	cfg := DefaultConfig()

	// Production server settings
	cfg.Server.Mode = "release"
	cfg.Server.ReadTimeout = 15 * time.Second
	cfg.Server.WriteTimeout = 15 * time.Second
	cfg.Server.IdleTimeout = 60 * time.Second

	// Production logging settings
	cfg.Logging.Level = "warn"
	cfg.Logging.Format = "json"
	cfg.Logging.ConsoleColor = false
	cfg.Logging.FileEnabled = true
	cfg.Logging.FileFormat = "json"

	// Production CORS settings
	cfg.CORS.AllowAllOrigins = false
	cfg.CORS.AllowCredentials = true

	// Production security settings
	cfg.Security.RateLimitEnabled = true
	cfg.Security.RateLimitRPS = 50
	cfg.Security.RequestTimeout = 15 * time.Second
	cfg.Security.MaxRequestSize = 5 * 1024 * 1024 // 5MB

	// Production database settings
	cfg.Database.SSLMode = "require"
	cfg.Database.MaxOpenConns = 50
	cfg.Database.MaxIdleConns = 10
	cfg.Database.ConnMaxLifetime = 10 * time.Minute

	return cfg
}

// TestConfig returns a configuration optimized for testing
func TestConfig() *Config {
	cfg := DefaultConfig()

	// Test server settings
	cfg.Server.Mode = "test"
	cfg.Server.Port = "0" // Random port

	// Test logging settings
	cfg.Logging.Level = "debug"
	cfg.Logging.Format = "console"
	cfg.Logging.FileEnabled = false

	// Test database settings
	cfg.Database.Name = "meow_test"
	cfg.Database.MaxOpenConns = 5
	cfg.Database.MaxIdleConns = 2

	// Test webhook settings
	cfg.Webhook.Timeout = 5 * time.Second
	cfg.Webhook.MaxRetries = 1

	// Test meow settings
	cfg.Meow.ConnectionTimeout = 5 * time.Second
	cfg.Meow.QRCodeTimeout = 10 * time.Second
	cfg.Meow.MaxRetries = 1

	return cfg
}

// DevelopmentConfig returns a configuration optimized for development
func DevelopmentConfig() *Config {
	cfg := DefaultConfig()

	// Development server settings
	cfg.Server.Mode = "debug"

	// Development logging settings
	cfg.Logging.Level = "debug"
	cfg.Logging.Format = "console"
	cfg.Logging.ConsoleColor = true
	cfg.Logging.FileEnabled = false // Keep simple for development

	// Development CORS settings (more permissive)
	cfg.CORS.AllowAllOrigins = true
	cfg.CORS.AllowCredentials = false

	// Development security settings (less restrictive)
	cfg.Security.RateLimitEnabled = false
	cfg.Security.RequestTimeout = 60 * time.Second
	cfg.Security.MaxRequestSize = 50 * 1024 * 1024 // 50MB for development

	return cfg
}
