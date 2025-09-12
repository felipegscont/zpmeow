package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBSslMode  string `env:"DB_SSLMODE"`
	GinMode    string `env:"GIN_MODE"`
	ServerPort string `env:"SERVER_PORT"`
	DBUrl      string

	LogLevel          string `env:"LOG_LEVEL"`
	LogFormat         string `env:"LOG_FORMAT"`
	LogConsoleColor   bool   `env:"LOG_CONSOLE_COLOR"`
	LogFileEnabled    bool   `env:"LOG_FILE_ENABLED"`
	LogFilePath       string `env:"LOG_FILE_PATH"`
	LogFileMaxSize    int    `env:"LOG_FILE_MAX_SIZE"`
	LogFileMaxBackups int    `env:"LOG_FILE_MAX_BACKUPS"`
	LogFileMaxAge     int    `env:"LOG_FILE_MAX_AGE"`
	LogFileCompress   bool   `env:"LOG_FILE_COMPRESS"`
	LogFileFormat     string `env:"LOG_FILE_FORMAT"`
}

func LoadConfig() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSslMode:  os.Getenv("DB_SSLMODE"),
		GinMode:    os.Getenv("GIN_MODE"),
		ServerPort: os.Getenv("SERVER_PORT"),

		LogLevel:          os.Getenv("LOG_LEVEL"),
		LogFormat:         os.Getenv("LOG_FORMAT"),
		LogConsoleColor:   getBoolEnv("LOG_CONSOLE_COLOR", true),
		LogFileEnabled:    getBoolEnv("LOG_FILE_ENABLED", true),
		LogFilePath:       os.Getenv("LOG_FILE_PATH"),
		LogFileMaxSize:    getIntEnv("LOG_FILE_MAX_SIZE", 100),
		LogFileMaxBackups: getIntEnv("LOG_FILE_MAX_BACKUPS", 3),
		LogFileMaxAge:     getIntEnv("LOG_FILE_MAX_AGE", 28),
		LogFileCompress:   getBoolEnv("LOG_FILE_COMPRESS", true),
		LogFileFormat:     os.Getenv("LOG_FILE_FORMAT"),
	}

	if cfg.GinMode == "" {
		cfg.GinMode = "debug"
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.LogFormat == "" {
		cfg.LogFormat = "console"
	}
	if cfg.LogFilePath == "" {
		cfg.LogFilePath = "log/app.log"
	}
	if cfg.LogFileFormat == "" {
		cfg.LogFileFormat = "json"
	}

	cfg.DBUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslMode)

	return cfg, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

type LoggerConfig struct {
	Level          string
	Format         string
	ConsoleColor   bool
	FileEnabled    bool
	FilePath       string
	FileMaxSize    int
	FileMaxBackups int
	FileMaxAge     int
	FileCompress   bool
	FileFormat     string
}

func (c *LoggerConfig) GetLevel() string {
	return c.Level
}

func (c *LoggerConfig) GetFormat() string {
	return c.Format
}

func (c *LoggerConfig) GetConsoleColor() bool {
	return c.ConsoleColor
}

func (c *LoggerConfig) GetFileEnabled() bool {
	return c.FileEnabled
}

func (c *LoggerConfig) GetFilePath() string {
	return c.FilePath
}

func (c *LoggerConfig) GetFileMaxSize() int {
	return c.FileMaxSize
}

func (c *LoggerConfig) GetFileMaxBackups() int {
	return c.FileMaxBackups
}

func (c *LoggerConfig) GetFileMaxAge() int {
	return c.FileMaxAge
}

func (c *LoggerConfig) GetFileCompress() bool {
	return c.FileCompress
}

func (c *LoggerConfig) GetFileFormat() string {
	return c.FileFormat
}

func (c *Config) GetLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:          c.LogLevel,
		Format:         c.LogFormat,
		ConsoleColor:   c.LogConsoleColor,
		FileEnabled:    c.LogFileEnabled,
		FilePath:       c.LogFilePath,
		FileMaxSize:    c.LogFileMaxSize,
		FileMaxBackups: c.LogFileMaxBackups,
		FileMaxAge:     c.LogFileMaxAge,
		FileCompress:   c.LogFileCompress,
		FileFormat:     c.LogFileFormat,
	}
}

func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:          "info",
		Format:         "console",
		ConsoleColor:   true,
		FileEnabled:    true,
		FilePath:       "log/app.log",
		FileMaxSize:    100,
		FileMaxBackups: 3,
		FileMaxAge:     28,
		FileCompress:   true,
		FileFormat:     "json",
	}
}
