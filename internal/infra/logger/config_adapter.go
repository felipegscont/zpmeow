package logger

// ConfigAdapter adapts external config to logger Config interface
type ConfigAdapter struct {
	Level           string
	Format          string
	ConsoleColor    bool
	FileEnabled     bool
	FilePath        string
	FileMaxSize     int
	FileMaxBackups  int
	FileMaxAge      int
	FileCompress    bool
	FileFormat      string
}

// GetLevel implements Config interface
func (c *ConfigAdapter) GetLevel() string {
	return c.Level
}

// GetFormat implements Config interface
func (c *ConfigAdapter) GetFormat() string {
	return c.Format
}

// GetConsoleColor implements Config interface
func (c *ConfigAdapter) GetConsoleColor() bool {
	return c.ConsoleColor
}

// GetFileEnabled implements Config interface
func (c *ConfigAdapter) GetFileEnabled() bool {
	return c.FileEnabled
}

// GetFilePath implements Config interface
func (c *ConfigAdapter) GetFilePath() string {
	return c.FilePath
}

// GetFileMaxSize implements Config interface
func (c *ConfigAdapter) GetFileMaxSize() int {
	return c.FileMaxSize
}

// GetFileMaxBackups implements Config interface
func (c *ConfigAdapter) GetFileMaxBackups() int {
	return c.FileMaxBackups
}

// GetFileMaxAge implements Config interface
func (c *ConfigAdapter) GetFileMaxAge() int {
	return c.FileMaxAge
}

// GetFileCompress implements Config interface
func (c *ConfigAdapter) GetFileCompress() bool {
	return c.FileCompress
}

// GetFileFormat implements Config interface
func (c *ConfigAdapter) GetFileFormat() string {
	return c.FileFormat
}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter(level, format, filePath, fileFormat string, consoleColor, fileEnabled, fileCompress bool, fileMaxSize, fileMaxBackups, fileMaxAge int) Config {
	return &ConfigAdapter{
		Level:           level,
		Format:          format,
		ConsoleColor:    consoleColor,
		FileEnabled:     fileEnabled,
		FilePath:        filePath,
		FileMaxSize:     fileMaxSize,
		FileMaxBackups:  fileMaxBackups,
		FileMaxAge:      fileMaxAge,
		FileCompress:    fileCompress,
		FileFormat:      fileFormat,
	}
}
