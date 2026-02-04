// Package config provides configuration management for the goog CLI application.
// It handles loading, saving, and managing application configuration with
// platform-specific paths and environment variable overrides.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	// DefaultAccount is the email of the default Google account to use.
	DefaultAccount string `yaml:"default_account" mapstructure:"default_account"`

	// DefaultFormat specifies the default output format (json|plain|table).
	DefaultFormat string `yaml:"default_format" mapstructure:"default_format"`

	// Timezone specifies the timezone for displaying dates and times.
	Timezone string `yaml:"timezone" mapstructure:"timezone"`

	// Accounts contains configuration for each authenticated account.
	Accounts map[string]AccountConfig `yaml:"accounts" mapstructure:"accounts"`

	// Mail contains mail-specific settings.
	Mail MailConfig `yaml:"mail" mapstructure:"mail"`

	// Calendar contains calendar-specific settings.
	Calendar CalendarConfig `yaml:"calendar" mapstructure:"calendar"`
}

// AccountConfig represents configuration for a single Google account.
type AccountConfig struct {
	// Email is the account's email address.
	Email string `yaml:"email" mapstructure:"email"`

	// Scopes lists the OAuth scopes granted to this account.
	Scopes []string `yaml:"scopes" mapstructure:"scopes"`

	// AddedAt is the timestamp when the account was added.
	AddedAt time.Time `yaml:"added_at" mapstructure:"added_at"`
}

// MailConfig contains mail-specific settings.
type MailConfig struct {
	// DefaultLabel is the default label to list messages from.
	DefaultLabel string `yaml:"default_label" mapstructure:"default_label"`

	// PageSize is the default number of messages to fetch per page.
	PageSize int `yaml:"page_size" mapstructure:"page_size"`
}

// CalendarConfig contains calendar-specific settings.
type CalendarConfig struct {
	// DefaultCalendar is the ID of the default calendar to use.
	DefaultCalendar string `yaml:"default_calendar" mapstructure:"default_calendar"`

	// WeekStart specifies the first day of the week (sunday|monday).
	WeekStart string `yaml:"week_start" mapstructure:"week_start"`
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		DefaultAccount: "",
		DefaultFormat:  "table",
		Timezone:       "Local",
		Accounts:       make(map[string]AccountConfig),
		Mail: MailConfig{
			DefaultLabel: "INBOX",
			PageSize:     20,
		},
		Calendar: CalendarConfig{
			DefaultCalendar: "primary",
			WeekStart:       "sunday",
		},
	}
}

// GetConfigPath returns the platform-specific configuration file path.
// The path can be overridden by setting the GOOG_CONFIG environment variable.
func GetConfigPath() string {
	// Check for environment variable override
	if envPath := os.Getenv("GOOG_CONFIG"); envPath != "" {
		return envPath
	}

	var configDir string

	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/goog/config.yaml
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		configDir = filepath.Join(home, "Library", "Application Support", "goog")

	case "windows":
		// Windows: %APPDATA%\goog\config.yaml
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				home = "."
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "goog")

	default:
		// Linux and others: ~/.config/goog/config.yaml
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				home = "."
			}
			configDir = filepath.Join(home, ".config")
		}
		configDir = filepath.Join(configDir, "goog")
	}

	return filepath.Join(configDir, "config.yaml")
}

// Load reads the configuration from the config file.
// If the file does not exist, it creates a default configuration.
// Environment variables can override specific settings:
//   - GOOG_ACCOUNT overrides default_account
//   - GOOG_FORMAT overrides default_format
//   - GOOG_CONFIG overrides the config file path
func Load() (*Config, error) {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	_, err := os.Stat(configPath)
	configExists := err == nil

	// Set up viper
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("default_account", "")
	v.SetDefault("default_format", "table")
	v.SetDefault("timezone", "Local")
	v.SetDefault("accounts", make(map[string]AccountConfig))
	v.SetDefault("mail.default_label", "INBOX")
	v.SetDefault("mail.page_size", 20)
	v.SetDefault("calendar.default_calendar", "primary")
	v.SetDefault("calendar.week_start", "sunday")

	// Read config file if it exists
	if configExists {
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Apply environment variable overrides
	if envAccount := os.Getenv("GOOG_ACCOUNT"); envAccount != "" {
		v.Set("default_account", envAccount)
	}
	if envFormat := os.Getenv("GOOG_FORMAT"); envFormat != "" {
		v.Set("default_format", envFormat)
	}

	// Unmarshal into config struct with custom decode hook for time.Time
	cfg := NewConfig()
	if err := v.Unmarshal(cfg, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			stringToTimeHookFunc(),
		),
	)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Ensure accounts map is initialized
	if cfg.Accounts == nil {
		cfg.Accounts = make(map[string]AccountConfig)
	}

	// If config didn't exist, save the default
	if !configExists {
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
	}

	return cfg, nil
}

// Save writes the configuration to the config file.
// It creates the config directory if it doesn't exist and
// creates the file with secure permissions (0600) from the start to avoid
// race conditions where the file could be read before permissions are set.
func (c *Config) Save() error {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set up viper for writing
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set values from config struct
	v.Set("default_account", c.DefaultAccount)
	v.Set("default_format", c.DefaultFormat)
	v.Set("timezone", c.Timezone)
	v.Set("accounts", c.Accounts)
	v.Set("mail", c.Mail)
	v.Set("calendar", c.Calendar)

	// Write config securely to avoid race condition
	if err := writeConfigSecurely(configPath, v); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// writeConfigSecurely writes the viper configuration to a file with secure
// permissions (0600) from the start. This avoids the race condition where
// the file is created with default permissions and then chmod'd.
func writeConfigSecurely(configPath string, v *viper.Viper) error {
	// On Windows, just use viper's default behavior
	if runtime.GOOS == "windows" {
		if err := v.WriteConfig(); err != nil {
			if os.IsNotExist(err) {
				return v.SafeWriteConfig()
			}
			return err
		}
		return nil
	}

	// For Unix-like systems, create file with secure permissions first
	// Use O_CREATE|O_WRONLY|O_TRUNC to create or truncate the file
	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Get the configuration as YAML using viper's AllSettings
	settings := v.AllSettings()
	yamlData, err := marshalYAML(settings)
	if err != nil {
		f.Close()
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to the file
	if _, err := f.Write(yamlData); err != nil {
		f.Close()
		return fmt.Errorf("failed to write config: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close config file: %w", err)
	}

	return nil
}

// marshalYAML marshals the settings map to YAML format.
func marshalYAML(settings map[string]interface{}) ([]byte, error) {
	// Use gopkg.in/yaml.v3 for marshaling (viper uses this internally)
	return yaml.Marshal(settings)
}

// SetPermissions sets the config file permissions to 0600 (owner read/write only).
// This is a no-op on Windows where file permissions work differently.
func SetPermissions() error {
	if runtime.GOOS == "windows" {
		return nil
	}

	configPath := GetConfigPath()
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", configPath, err)
	}

	return nil
}

// ErrAccountNotFound is returned when the requested account is not found.
var ErrAccountNotFound = fmt.Errorf("account not found")

// GetAccount retrieves an account configuration by alias.
func (c *Config) GetAccount(alias string) (*AccountConfig, error) {
	acc, ok := c.Accounts[alias]
	if !ok {
		return nil, ErrAccountNotFound
	}
	return &acc, nil
}

// validFormats lists the valid output format options.
var validFormats = map[string]bool{
	"json":  true,
	"plain": true,
	"table": true,
}

// SetValue sets a configuration value by key path (e.g., "mail.page_size").
func (c *Config) SetValue(key, value string) error {
	switch key {
	case "default_account":
		c.DefaultAccount = value
	case "default_format":
		if !validFormats[value] {
			return fmt.Errorf("invalid format %q: must be one of json, plain, table", value)
		}
		c.DefaultFormat = value
	case "timezone":
		// Validate timezone using time.LoadLocation
		if value != "" && value != "Local" {
			if _, err := time.LoadLocation(value); err != nil {
				return fmt.Errorf("invalid timezone %q: %w", value, err)
			}
		}
		c.Timezone = value
	case "mail.default_label":
		c.Mail.DefaultLabel = value
	case "mail.page_size":
		var pageSize int
		if _, err := fmt.Sscanf(value, "%d", &pageSize); err != nil {
			return fmt.Errorf("invalid page_size: %w", err)
		}
		c.Mail.PageSize = pageSize
	case "calendar.default_calendar":
		c.Calendar.DefaultCalendar = value
	case "calendar.week_start":
		c.Calendar.WeekStart = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

// GetValue retrieves a configuration value by key path.
func (c *Config) GetValue(key string) (string, error) {
	switch key {
	case "default_account":
		return c.DefaultAccount, nil
	case "default_format":
		return c.DefaultFormat, nil
	case "timezone":
		return c.Timezone, nil
	case "mail.default_label":
		return c.Mail.DefaultLabel, nil
	case "mail.page_size":
		return fmt.Sprintf("%d", c.Mail.PageSize), nil
	case "calendar.default_calendar":
		return c.Calendar.DefaultCalendar, nil
	case "calendar.week_start":
		return c.Calendar.WeekStart, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// stringToTimeHookFunc returns a mapstructure decode hook that converts
// strings to time.Time using RFC3339 format.
func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		str := data.(string)
		if str == "" {
			return time.Time{}, nil
		}

		// Try RFC3339 first (standard format)
		parsed, err := time.Parse(time.RFC3339, str)
		if err == nil {
			return parsed, nil
		}

		// Try RFC3339Nano
		parsed, err = time.Parse(time.RFC3339Nano, str)
		if err == nil {
			return parsed, nil
		}

		return nil, fmt.Errorf("unable to parse time: %s", str)
	}
}
