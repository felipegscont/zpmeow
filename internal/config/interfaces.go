package config

import "time"

// ConfigProvider provides access to application configuration
type ConfigProvider interface {
	GetDatabase() DatabaseConfigProvider
	GetServer() ServerConfigProvider
	GetAuth() AuthConfigProvider
	GetLogging() LoggingConfigProvider
	GetCORS() CORSConfigProvider
	GetWebhook() WebhookConfigProvider
	GetMeow() MeowConfigProvider
	GetSecurity() SecurityConfigProvider
}

// DatabaseConfigProvider provides database configuration
type DatabaseConfigProvider interface {
	GetHost() string
	GetPort() string
	GetUser() string
	GetPassword() string
	GetName() string
	GetSSLMode() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxLifetime() time.Duration
	GetURL() string
}

// ServerConfigProvider provides server configuration
type ServerConfigProvider interface {
	GetPort() string
	GetMode() string
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetIdleTimeout() time.Duration
}

// AuthConfigProvider provides authentication configuration
type AuthConfigProvider interface {
	GetGlobalAPIKey() string
	GetSessionTimeout() time.Duration
	GetTokenExpiration() time.Duration
}

// LoggingConfigProvider provides logging configuration
type LoggingConfigProvider interface {
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

// CORSConfigProvider provides CORS configuration
type CORSConfigProvider interface {
	GetAllowAllOrigins() bool
	GetAllowOrigins() []string
	GetAllowMethods() []string
	GetAllowHeaders() []string
	GetExposeHeaders() []string
	GetAllowCredentials() bool
	GetMaxAge() int
}

// WebhookConfigProvider provides webhook configuration
type WebhookConfigProvider interface {
	GetTimeout() time.Duration
	GetMaxRetries() int
	GetInitialBackoff() time.Duration
	GetMaxBackoff() time.Duration
	GetBackoffMultiplier() float64
}

// MeowConfigProvider provides meow configuration (exported)
type MeowConfigProvider interface {
	GetMaxRetries() int
	GetRetryInterval() time.Duration
	GetConnectionTimeout() time.Duration
	GetQRCodeTimeout() time.Duration
	GetReconnectDelay() time.Duration
}

// SecurityConfigProvider provides security configuration
type SecurityConfigProvider interface {
	GetRateLimitEnabled() bool
	GetRateLimitRPS() int
	GetRequestTimeout() time.Duration
	GetMaxRequestSize() int64
}

// Implementation of interfaces

// GetDatabase returns database configuration provider
func (c *Config) GetDatabase() DatabaseConfigProvider {
	return &c.Database
}

// GetServer returns server configuration provider
func (c *Config) GetServer() ServerConfigProvider {
	return &c.Server
}

// GetAuth returns auth configuration provider
func (c *Config) GetAuth() AuthConfigProvider {
	return &c.Auth
}

// GetLogging returns logging configuration provider
func (c *Config) GetLogging() LoggingConfigProvider {
	return &c.Logging
}

// GetCORS returns CORS configuration provider
func (c *Config) GetCORS() CORSConfigProvider {
	return &c.CORS
}

// GetWebhook returns webhook configuration provider
func (c *Config) GetWebhook() WebhookConfigProvider {
	return &c.Webhook
}

// GetMeow returns meow configuration provider
func (c *Config) GetMeow() MeowConfigProvider {
	return &c.Meow
}

// GetSecurity returns security configuration provider
func (c *Config) GetSecurity() SecurityConfigProvider {
	return &c.Security
}

// DatabaseConfig interface implementations

func (d *DatabaseConfig) GetHost() string                   { return d.Host }
func (d *DatabaseConfig) GetPort() string                   { return d.Port }
func (d *DatabaseConfig) GetUser() string                   { return d.User }
func (d *DatabaseConfig) GetPassword() string               { return d.Password }
func (d *DatabaseConfig) GetName() string                   { return d.Name }
func (d *DatabaseConfig) GetSSLMode() string                { return d.SSLMode }
func (d *DatabaseConfig) GetMaxOpenConns() int              { return d.MaxOpenConns }
func (d *DatabaseConfig) GetMaxIdleConns() int              { return d.MaxIdleConns }
func (d *DatabaseConfig) GetConnMaxLifetime() time.Duration { return d.ConnMaxLifetime }
func (d *DatabaseConfig) GetURL() string                    { return d.URL }

// ServerConfig interface implementations

func (s *ServerConfig) GetPort() string                { return s.Port }
func (s *ServerConfig) GetMode() string                { return s.Mode }
func (s *ServerConfig) GetReadTimeout() time.Duration  { return s.ReadTimeout }
func (s *ServerConfig) GetWriteTimeout() time.Duration { return s.WriteTimeout }
func (s *ServerConfig) GetIdleTimeout() time.Duration  { return s.IdleTimeout }

// AuthConfig interface implementations

func (a *AuthConfig) GetGlobalAPIKey() string           { return a.GlobalAPIKey }
func (a *AuthConfig) GetSessionTimeout() time.Duration  { return a.SessionTimeout }
func (a *AuthConfig) GetTokenExpiration() time.Duration { return a.TokenExpiration }

// LoggingConfig interface implementations

func (l *LoggingConfig) GetLevel() string       { return l.Level }
func (l *LoggingConfig) GetFormat() string      { return l.Format }
func (l *LoggingConfig) GetConsoleColor() bool  { return l.ConsoleColor }
func (l *LoggingConfig) GetFileEnabled() bool   { return l.FileEnabled }
func (l *LoggingConfig) GetFilePath() string    { return l.FilePath }
func (l *LoggingConfig) GetFileMaxSize() int    { return l.FileMaxSize }
func (l *LoggingConfig) GetFileMaxBackups() int { return l.FileMaxBackups }
func (l *LoggingConfig) GetFileMaxAge() int     { return l.FileMaxAge }
func (l *LoggingConfig) GetFileCompress() bool  { return l.FileCompress }
func (l *LoggingConfig) GetFileFormat() string  { return l.FileFormat }

// CORSConfig interface implementations

func (c *CORSConfig) GetAllowAllOrigins() bool   { return c.AllowAllOrigins }
func (c *CORSConfig) GetAllowOrigins() []string  { return c.AllowOrigins }
func (c *CORSConfig) GetAllowMethods() []string  { return c.AllowMethods }
func (c *CORSConfig) GetAllowHeaders() []string  { return c.AllowHeaders }
func (c *CORSConfig) GetExposeHeaders() []string { return c.ExposeHeaders }
func (c *CORSConfig) GetAllowCredentials() bool  { return c.AllowCredentials }
func (c *CORSConfig) GetMaxAge() int             { return c.MaxAge }

// WebhookConfig interface implementations

func (w *WebhookConfig) GetTimeout() time.Duration        { return w.Timeout }
func (w *WebhookConfig) GetMaxRetries() int               { return w.MaxRetries }
func (w *WebhookConfig) GetInitialBackoff() time.Duration { return w.InitialBackoff }
func (w *WebhookConfig) GetMaxBackoff() time.Duration     { return w.MaxBackoff }
func (w *WebhookConfig) GetBackoffMultiplier() float64    { return w.BackoffMultiplier }

// MeowConfig interface implementations

func (w *MeowConfig) GetMaxRetries() int                  { return w.MaxRetries }
func (w *MeowConfig) GetRetryInterval() time.Duration     { return w.RetryInterval }
func (w *MeowConfig) GetConnectionTimeout() time.Duration { return w.ConnectionTimeout }
func (w *MeowConfig) GetQRCodeTimeout() time.Duration     { return w.QRCodeTimeout }
func (w *MeowConfig) GetReconnectDelay() time.Duration    { return w.ReconnectDelay }

// SecurityConfig interface implementations

func (s *SecurityConfig) GetRateLimitEnabled() bool        { return s.RateLimitEnabled }
func (s *SecurityConfig) GetRateLimitRPS() int             { return s.RateLimitRPS }
func (s *SecurityConfig) GetRequestTimeout() time.Duration { return s.RequestTimeout }
func (s *SecurityConfig) GetMaxRequestSize() int64         { return s.MaxRequestSize }
