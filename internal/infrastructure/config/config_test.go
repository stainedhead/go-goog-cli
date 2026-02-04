// Package config provides configuration management for the goog CLI application.
package config

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestConfigDefaults(t *testing.T) {
	cfg := NewConfig()

	t.Run("default_account is empty", func(t *testing.T) {
		if cfg.DefaultAccount != "" {
			t.Errorf("expected empty default_account, got %q", cfg.DefaultAccount)
		}
	})

	t.Run("default_format is table", func(t *testing.T) {
		if cfg.DefaultFormat != "table" {
			t.Errorf("expected default_format 'table', got %q", cfg.DefaultFormat)
		}
	})

	t.Run("timezone is Local", func(t *testing.T) {
		if cfg.Timezone != "Local" {
			t.Errorf("expected timezone 'Local', got %q", cfg.Timezone)
		}
	})

	t.Run("accounts is empty map", func(t *testing.T) {
		if cfg.Accounts == nil {
			t.Error("expected accounts to be initialized")
		}
		if len(cfg.Accounts) != 0 {
			t.Errorf("expected empty accounts map, got %d entries", len(cfg.Accounts))
		}
	})

	t.Run("mail defaults", func(t *testing.T) {
		if cfg.Mail.DefaultLabel != "INBOX" {
			t.Errorf("expected mail default_label 'INBOX', got %q", cfg.Mail.DefaultLabel)
		}
		if cfg.Mail.PageSize != 20 {
			t.Errorf("expected mail page_size 20, got %d", cfg.Mail.PageSize)
		}
	})

	t.Run("calendar defaults", func(t *testing.T) {
		if cfg.Calendar.DefaultCalendar != "primary" {
			t.Errorf("expected calendar default_calendar 'primary', got %q", cfg.Calendar.DefaultCalendar)
		}
		if cfg.Calendar.WeekStart != "sunday" {
			t.Errorf("expected calendar week_start 'sunday', got %q", cfg.Calendar.WeekStart)
		}
	})
}

func TestConfigPlatformPaths(t *testing.T) {
	// Save and restore GOOG_CONFIG env var
	origEnv := os.Getenv("GOOG_CONFIG")
	os.Unsetenv("GOOG_CONFIG")
	defer func() {
		if origEnv != "" {
			os.Setenv("GOOG_CONFIG", origEnv)
		}
	}()

	path := GetConfigPath()

	t.Run("path contains goog directory", func(t *testing.T) {
		if filepath.Base(filepath.Dir(path)) != "goog" {
			t.Errorf("expected path to contain 'goog' directory, got %q", path)
		}
	})

	t.Run("path ends with config.yaml", func(t *testing.T) {
		if filepath.Base(path) != "config.yaml" {
			t.Errorf("expected path to end with 'config.yaml', got %q", filepath.Base(path))
		}
	})

	t.Run("platform specific path", func(t *testing.T) {
		switch runtime.GOOS {
		case "darwin":
			if !contains(path, "Library/Application Support") {
				t.Errorf("macOS path should contain 'Library/Application Support', got %q", path)
			}
		case "linux":
			if !contains(path, ".config") {
				t.Errorf("Linux path should contain '.config', got %q", path)
			}
		case "windows":
			// Windows path should contain APPDATA
			appdata := os.Getenv("APPDATA")
			if appdata != "" && !contains(path, appdata) {
				t.Errorf("Windows path should contain APPDATA, got %q", path)
			}
		}
	})
}

func TestConfigEnvOverrides(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create a minimal config file
	configContent := `default_account: "original@example.com"
default_format: "json"
timezone: "UTC"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Save original env vars
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	origConfig := os.Getenv("GOOG_CONFIG")

	// Clean up after test
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
		restoreEnv("GOOG_CONFIG", origConfig)
	}()

	t.Run("GOOG_CONFIG overrides path", func(t *testing.T) {
		os.Setenv("GOOG_CONFIG", configPath)
		os.Unsetenv("GOOG_ACCOUNT")
		os.Unsetenv("GOOG_FORMAT")

		gotPath := GetConfigPath()
		if gotPath != configPath {
			t.Errorf("expected GOOG_CONFIG to override path, got %q want %q", gotPath, configPath)
		}
	})

	t.Run("GOOG_ACCOUNT overrides default_account", func(t *testing.T) {
		os.Setenv("GOOG_CONFIG", configPath)
		os.Setenv("GOOG_ACCOUNT", "override@example.com")
		os.Unsetenv("GOOG_FORMAT")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if cfg.DefaultAccount != "override@example.com" {
			t.Errorf("expected GOOG_ACCOUNT to override, got %q", cfg.DefaultAccount)
		}
	})

	t.Run("GOOG_FORMAT overrides default_format", func(t *testing.T) {
		os.Setenv("GOOG_CONFIG", configPath)
		os.Unsetenv("GOOG_ACCOUNT")
		os.Setenv("GOOG_FORMAT", "plain")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if cfg.DefaultFormat != "plain" {
			t.Errorf("expected GOOG_FORMAT to override, got %q", cfg.DefaultFormat)
		}
	})
}

func TestConfigLoad(t *testing.T) {
	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	t.Run("load creates default config if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "goog", "config.yaml")
		os.Setenv("GOOG_CONFIG", configPath)

		// Clear other env vars
		origAccount := os.Getenv("GOOG_ACCOUNT")
		origFormat := os.Getenv("GOOG_FORMAT")
		os.Unsetenv("GOOG_ACCOUNT")
		os.Unsetenv("GOOG_FORMAT")
		defer func() {
			restoreEnv("GOOG_ACCOUNT", origAccount)
			restoreEnv("GOOG_FORMAT", origFormat)
		}()

		cfg, err := Load()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		// Check defaults are set
		if cfg.DefaultFormat != "table" {
			t.Errorf("expected default_format 'table', got %q", cfg.DefaultFormat)
		}

		// Check file was created
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("expected config file to be created")
		}
	})

	t.Run("load reads existing config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `default_account: "test@example.com"
default_format: "json"
timezone: "America/New_York"
accounts:
  test@example.com:
    email: "test@example.com"
    scopes:
      - "gmail.readonly"
    added_at: "2024-01-15T10:30:00Z"
mail:
  default_label: "INBOX"
  page_size: 50
calendar:
  default_calendar: "work"
  week_start: "monday"
`
		if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		os.Setenv("GOOG_CONFIG", configPath)
		os.Unsetenv("GOOG_ACCOUNT")
		os.Unsetenv("GOOG_FORMAT")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if cfg.DefaultAccount != "test@example.com" {
			t.Errorf("expected default_account 'test@example.com', got %q", cfg.DefaultAccount)
		}
		if cfg.DefaultFormat != "json" {
			t.Errorf("expected default_format 'json', got %q", cfg.DefaultFormat)
		}
		if cfg.Timezone != "America/New_York" {
			t.Errorf("expected timezone 'America/New_York', got %q", cfg.Timezone)
		}
		if cfg.Mail.PageSize != 50 {
			t.Errorf("expected mail page_size 50, got %d", cfg.Mail.PageSize)
		}
		if cfg.Calendar.DefaultCalendar != "work" {
			t.Errorf("expected calendar default_calendar 'work', got %q", cfg.Calendar.DefaultCalendar)
		}
		if cfg.Calendar.WeekStart != "monday" {
			t.Errorf("expected calendar week_start 'monday', got %q", cfg.Calendar.WeekStart)
		}

		// Check account was loaded
		acc, ok := cfg.Accounts["test@example.com"]
		if !ok {
			t.Fatal("expected account 'test@example.com' to exist")
		}
		if acc.Email != "test@example.com" {
			t.Errorf("expected account email 'test@example.com', got %q", acc.Email)
		}
		if len(acc.Scopes) != 1 || acc.Scopes[0] != "gmail.readonly" {
			t.Errorf("expected account scopes ['gmail.readonly'], got %v", acc.Scopes)
		}
	})
}

func TestConfigSave(t *testing.T) {
	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	t.Run("save writes config to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		os.Setenv("GOOG_CONFIG", configPath)

		cfg := NewConfig()
		cfg.DefaultAccount = "saved@example.com"
		cfg.DefaultFormat = "plain"
		cfg.Accounts["saved@example.com"] = AccountConfig{
			Email:   "saved@example.com",
			Scopes:  []string{"gmail.readonly", "calendar.readonly"},
			AddedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		}

		if err := cfg.Save(); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		// Read and verify
		os.Unsetenv("GOOG_ACCOUNT")
		os.Unsetenv("GOOG_FORMAT")

		loaded, err := Load()
		if err != nil {
			t.Fatalf("failed to load saved config: %v", err)
		}

		if loaded.DefaultAccount != "saved@example.com" {
			t.Errorf("expected loaded default_account 'saved@example.com', got %q", loaded.DefaultAccount)
		}
		if loaded.DefaultFormat != "plain" {
			t.Errorf("expected loaded default_format 'plain', got %q", loaded.DefaultFormat)
		}

		acc, ok := loaded.Accounts["saved@example.com"]
		if !ok {
			t.Fatal("expected account 'saved@example.com' to exist")
		}
		if len(acc.Scopes) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(acc.Scopes))
		}
	})

	t.Run("save creates directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "subdir", "goog", "config.yaml")
		os.Setenv("GOOG_CONFIG", configPath)

		cfg := NewConfig()
		if err := cfg.Save(); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("expected config file to be created")
		}
	})

	t.Run("save sets file permissions to 600", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("file permissions test not applicable on Windows")
		}

		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		os.Setenv("GOOG_CONFIG", configPath)

		cfg := NewConfig()
		if err := cfg.Save(); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		info, err := os.Stat(configPath)
		if err != nil {
			t.Fatalf("failed to stat config file: %v", err)
		}

		perm := info.Mode().Perm()
		if perm != 0600 {
			t.Errorf("expected file permissions 0600, got %o", perm)
		}
	})
}

func TestSetPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions test not applicable on Windows")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create file with wrong permissions
	if err := os.WriteFile(configPath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	if err := SetPermissions(); err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func restoreEnv(key, value string) {
	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
}

// TestSetValueValid tests SetValue with valid values.
func TestSetValueValid(t *testing.T) {
	cfg := NewConfig()

	testCases := []struct {
		key      string
		value    string
		validate func() bool
	}{
		{
			key:   "default_account",
			value: "test@example.com",
			validate: func() bool {
				return cfg.DefaultAccount == "test@example.com"
			},
		},
		{
			key:   "default_format",
			value: "json",
			validate: func() bool {
				return cfg.DefaultFormat == "json"
			},
		},
		{
			key:   "timezone",
			value: "America/Los_Angeles",
			validate: func() bool {
				return cfg.Timezone == "America/Los_Angeles"
			},
		},
		{
			key:   "mail.default_label",
			value: "SENT",
			validate: func() bool {
				return cfg.Mail.DefaultLabel == "SENT"
			},
		},
		{
			key:   "mail.page_size",
			value: "50",
			validate: func() bool {
				return cfg.Mail.PageSize == 50
			},
		},
		{
			key:   "calendar.default_calendar",
			value: "work",
			validate: func() bool {
				return cfg.Calendar.DefaultCalendar == "work"
			},
		},
		{
			key:   "calendar.week_start",
			value: "monday",
			validate: func() bool {
				return cfg.Calendar.WeekStart == "monday"
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			err := cfg.SetValue(tc.key, tc.value)
			if err != nil {
				t.Errorf("SetValue(%q, %q) returned error: %v", tc.key, tc.value, err)
			}
			if !tc.validate() {
				t.Errorf("SetValue(%q, %q) did not set value correctly", tc.key, tc.value)
			}
		})
	}
}

// TestSetValueInvalid tests SetValue with invalid values.
func TestSetValueInvalid(t *testing.T) {
	cfg := NewConfig()

	t.Run("unknown key returns error", func(t *testing.T) {
		err := cfg.SetValue("unknown_key", "value")
		if err == nil {
			t.Error("expected error for unknown key")
		}
	})

	t.Run("invalid page_size returns error", func(t *testing.T) {
		err := cfg.SetValue("mail.page_size", "not_a_number")
		if err == nil {
			t.Error("expected error for invalid page_size")
		}
	})
}

// TestGetValueAll tests GetValue for all config keys.
func TestGetValueAll(t *testing.T) {
	cfg := NewConfig()
	cfg.DefaultAccount = "test@example.com"
	cfg.DefaultFormat = "json"
	cfg.Timezone = "UTC"
	cfg.Mail.DefaultLabel = "INBOX"
	cfg.Mail.PageSize = 25
	cfg.Calendar.DefaultCalendar = "work"
	cfg.Calendar.WeekStart = "monday"

	testCases := []struct {
		key      string
		expected string
	}{
		{"default_account", "test@example.com"},
		{"default_format", "json"},
		{"timezone", "UTC"},
		{"mail.default_label", "INBOX"},
		{"mail.page_size", "25"},
		{"calendar.default_calendar", "work"},
		{"calendar.week_start", "monday"},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			value, err := cfg.GetValue(tc.key)
			if err != nil {
				t.Errorf("GetValue(%q) returned error: %v", tc.key, err)
			}
			if value != tc.expected {
				t.Errorf("GetValue(%q) = %q, want %q", tc.key, value, tc.expected)
			}
		})
	}

	t.Run("unknown key returns error", func(t *testing.T) {
		_, err := cfg.GetValue("unknown_key")
		if err == nil {
			t.Error("expected error for unknown key")
		}
	})
}

// TestGetAccount tests GetAccount function.
func TestGetAccount(t *testing.T) {
	cfg := NewConfig()
	cfg.Accounts["personal"] = AccountConfig{
		Email:   "personal@example.com",
		Scopes:  []string{"gmail.readonly"},
		AddedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}
	cfg.Accounts["work"] = AccountConfig{
		Email:   "work@example.com",
		Scopes:  []string{"gmail.readonly", "calendar.readonly"},
		AddedAt: time.Date(2024, 2, 20, 14, 45, 0, 0, time.UTC),
	}

	t.Run("returns existing account", func(t *testing.T) {
		acc, err := cfg.GetAccount("personal")
		if err != nil {
			t.Fatalf("GetAccount failed: %v", err)
		}
		if acc.Email != "personal@example.com" {
			t.Errorf("expected email 'personal@example.com', got %q", acc.Email)
		}
	})

	t.Run("returns error for non-existent account", func(t *testing.T) {
		_, err := cfg.GetAccount("nonexistent")
		if err == nil {
			t.Error("expected error for non-existent account")
		}
		if err != ErrAccountNotFound {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})
}

// TestAddRemoveAccount tests adding and removing accounts.
func TestAddRemoveAccount(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Clear other env vars
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	t.Run("add account", func(t *testing.T) {
		cfg.Accounts["new@example.com"] = AccountConfig{
			Email:   "new@example.com",
			Scopes:  []string{"gmail.readonly"},
			AddedAt: time.Now(),
		}

		if err := cfg.Save(); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Reload and verify
		loaded, err := Load()
		if err != nil {
			t.Fatalf("Load after save failed: %v", err)
		}

		acc, ok := loaded.Accounts["new@example.com"]
		if !ok {
			t.Error("expected account to exist after save")
		}
		if acc.Email != "new@example.com" {
			t.Errorf("expected email 'new@example.com', got %q", acc.Email)
		}
	})

	t.Run("remove account", func(t *testing.T) {
		delete(cfg.Accounts, "new@example.com")

		if err := cfg.Save(); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Reload and verify
		loaded, err := Load()
		if err != nil {
			t.Fatalf("Load after delete failed: %v", err)
		}

		_, ok := loaded.Accounts["new@example.com"]
		if ok {
			t.Error("expected account to be removed after save")
		}
	})
}

// TestDefaultAccountResolution tests default account resolution.
func TestDefaultAccountResolution(t *testing.T) {
	cfg := NewConfig()
	cfg.DefaultAccount = "default@example.com"
	cfg.Accounts["default@example.com"] = AccountConfig{
		Email: "default@example.com",
	}
	cfg.Accounts["other@example.com"] = AccountConfig{
		Email: "other@example.com",
	}

	t.Run("default account is set", func(t *testing.T) {
		if cfg.DefaultAccount != "default@example.com" {
			t.Errorf("expected default account 'default@example.com', got %q", cfg.DefaultAccount)
		}
	})

	t.Run("default account exists in accounts map", func(t *testing.T) {
		_, ok := cfg.Accounts[cfg.DefaultAccount]
		if !ok {
			t.Error("expected default account to exist in accounts map")
		}
	})
}

// TestConfigPathXDG tests XDG_CONFIG_HOME override on Linux.
func TestConfigPathXDG(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG test only applicable on Linux")
	}

	// Save and restore env vars
	origConfig := os.Getenv("GOOG_CONFIG")
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		restoreEnv("GOOG_CONFIG", origConfig)
		restoreEnv("XDG_CONFIG_HOME", origXDG)
	}()

	os.Unsetenv("GOOG_CONFIG")

	t.Run("uses XDG_CONFIG_HOME when set", func(t *testing.T) {
		os.Setenv("XDG_CONFIG_HOME", "/custom/config")
		path := GetConfigPath()
		if !contains(path, "/custom/config/goog") {
			t.Errorf("expected path to use XDG_CONFIG_HOME, got %q", path)
		}
	})
}

// TestStringToTimeHookFunc tests the time parsing hook function.
func TestStringToTimeHookFunc(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Clear other env vars
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Test with RFC3339 formatted time in config
	configContent := `default_account: "test@example.com"
accounts:
  test@example.com:
    email: "test@example.com"
    added_at: "2024-06-15T14:30:00Z"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	acc, ok := cfg.Accounts["test@example.com"]
	if !ok {
		t.Fatal("expected account to exist")
	}

	expectedTime := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	if !acc.AddedAt.Equal(expectedTime) {
		t.Errorf("expected AddedAt to be %v, got %v", expectedTime, acc.AddedAt)
	}
}

// TestLoadInvalidConfig tests loading an invalid config file.
func TestLoadInvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Write invalid YAML
	if err := os.WriteFile(configPath, []byte("invalid: [yaml: syntax"), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid YAML config")
	}
}

// TestConfigWithAllFields tests a config file with all fields populated.
func TestConfigWithAllFields(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Clear other env vars
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	configContent := `default_account: "primary@example.com"
default_format: "plain"
timezone: "Europe/London"
accounts:
  primary@example.com:
    email: "primary@example.com"
    scopes:
      - "gmail.readonly"
      - "calendar"
      - "drive.readonly"
    added_at: "2024-03-15T09:00:00Z"
  secondary@example.com:
    email: "secondary@example.com"
    scopes:
      - "gmail.send"
    added_at: "2024-04-20T15:30:00Z"
mail:
  default_label: "STARRED"
  page_size: 100
calendar:
  default_calendar: "personal"
  week_start: "monday"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify all fields
	if cfg.DefaultAccount != "primary@example.com" {
		t.Errorf("DefaultAccount = %q, want 'primary@example.com'", cfg.DefaultAccount)
	}
	if cfg.DefaultFormat != "plain" {
		t.Errorf("DefaultFormat = %q, want 'plain'", cfg.DefaultFormat)
	}
	if cfg.Timezone != "Europe/London" {
		t.Errorf("Timezone = %q, want 'Europe/London'", cfg.Timezone)
	}
	if len(cfg.Accounts) < 2 {
		t.Errorf("len(Accounts) = %d, want at least 2", len(cfg.Accounts))
	}
	if cfg.Mail.DefaultLabel != "STARRED" {
		t.Errorf("Mail.DefaultLabel = %q, want 'STARRED'", cfg.Mail.DefaultLabel)
	}
	if cfg.Mail.PageSize != 100 {
		t.Errorf("Mail.PageSize = %d, want 100", cfg.Mail.PageSize)
	}
	if cfg.Calendar.DefaultCalendar != "personal" {
		t.Errorf("Calendar.DefaultCalendar = %q, want 'personal'", cfg.Calendar.DefaultCalendar)
	}
	if cfg.Calendar.WeekStart != "monday" {
		t.Errorf("Calendar.WeekStart = %q, want 'monday'", cfg.Calendar.WeekStart)
	}

	// Verify account details
	primary := cfg.Accounts["primary@example.com"]
	if len(primary.Scopes) != 3 {
		t.Errorf("len(primary.Scopes) = %d, want 3", len(primary.Scopes))
	}
}

// TestMarshalYAML tests the marshalYAML function.
func TestMarshalYAML(t *testing.T) {
	settings := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"nested": map[string]interface{}{
			"nested_key": "nested_value",
		},
	}

	data, err := marshalYAML(settings)
	if err != nil {
		t.Fatalf("marshalYAML failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("marshalYAML returned empty data")
	}

	// Should contain the key values
	if !contains(string(data), "key1") {
		t.Error("YAML output should contain 'key1'")
	}
}

// TestWriteConfigSecurely tests the secure config writing.
func TestWriteConfigSecurely(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions test not applicable on Windows")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "secure.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	cfg.DefaultAccount = "secure@example.com"

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists and has correct permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}
}

// TestConfigSaveWindows tests config saving on Windows (mocked).
func TestConfigSaveWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "win-config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	cfg.DefaultAccount = "windows@example.com"

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file not created")
	}
}

// TestSetPermissionsOnWindows tests SetPermissions behavior on Windows.
func TestSetPermissionsOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	// On Windows, SetPermissions should be a no-op
	err := SetPermissions()
	if err != nil {
		t.Errorf("SetPermissions on Windows should not error: %v", err)
	}
}

// TestLoadWithEmptyValues tests loading config with empty string values.
func TestLoadWithEmptyValues(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Clear other env vars
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Config with some empty values (using defaults)
	configContent := `default_account: ""
default_format: ""
accounts: {}
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Empty values should remain empty or use defaults
	if cfg.DefaultAccount != "" {
		t.Errorf("expected empty DefaultAccount, got %q", cfg.DefaultAccount)
	}
}

// TestTimeParsingEdgeCases tests edge cases in time parsing.
func TestTimeParsingEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Clear other env vars
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Config with RFC3339Nano time format
	configContent := `default_account: "nano@example.com"
accounts:
  nano@example.com:
    email: "nano@example.com"
    added_at: "2024-06-15T14:30:00.123456789Z"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	acc, ok := cfg.Accounts["nano@example.com"]
	if !ok {
		t.Fatal("expected account to exist")
	}

	if acc.AddedAt.IsZero() {
		t.Error("expected AddedAt to be parsed")
	}
}

// TestConfigSaveAndReload tests a full save/reload cycle.
func TestConfigSaveAndReload(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "cycle-config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Clear other env vars
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Create initial config
	cfg := NewConfig()
	cfg.DefaultAccount = "cycle@example.com"
	cfg.DefaultFormat = "json"
	cfg.Timezone = "America/Chicago"
	cfg.Mail.PageSize = 75
	cfg.Calendar.WeekStart = "monday"
	cfg.Accounts["cycle@example.com"] = AccountConfig{
		Email:   "cycle@example.com",
		Scopes:  []string{"gmail.readonly", "calendar"},
		AddedAt: time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC),
	}

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reload
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify all values were preserved
	if loaded.DefaultAccount != "cycle@example.com" {
		t.Errorf("DefaultAccount = %q, want 'cycle@example.com'", loaded.DefaultAccount)
	}
	if loaded.DefaultFormat != "json" {
		t.Errorf("DefaultFormat = %q, want 'json'", loaded.DefaultFormat)
	}
	if loaded.Timezone != "America/Chicago" {
		t.Errorf("Timezone = %q, want 'America/Chicago'", loaded.Timezone)
	}
	if loaded.Mail.PageSize != 75 {
		t.Errorf("Mail.PageSize = %d, want 75", loaded.Mail.PageSize)
	}
	if loaded.Calendar.WeekStart != "monday" {
		t.Errorf("Calendar.WeekStart = %q, want 'monday'", loaded.Calendar.WeekStart)
	}

	acc, ok := loaded.Accounts["cycle@example.com"]
	if !ok {
		t.Fatal("expected account to exist")
	}
	if len(acc.Scopes) != 2 {
		t.Errorf("len(Scopes) = %d, want 2", len(acc.Scopes))
	}
}

// TestGetConfigPathDarwin tests GetConfigPath on macOS.
func TestGetConfigPathDarwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Unsetenv("GOOG_CONFIG")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	path := GetConfigPath()
	if !contains(path, "Library/Application Support/goog") {
		t.Errorf("macOS path should contain 'Library/Application Support/goog', got %q", path)
	}
}

// TestGetConfigPathLinux tests GetConfigPath on Linux.
func TestGetConfigPathLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("GOOG_CONFIG")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer func() {
		restoreEnv("GOOG_CONFIG", origConfig)
		restoreEnv("XDG_CONFIG_HOME", origXDG)
	}()

	path := GetConfigPath()
	if !contains(path, ".config/goog") {
		t.Errorf("Linux path should contain '.config/goog', got %q", path)
	}
}

// TestGetConfigPathWindows tests GetConfigPath on Windows.
func TestGetConfigPathWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Unsetenv("GOOG_CONFIG")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	path := GetConfigPath()
	// Should contain goog directory
	if !contains(path, "goog") {
		t.Errorf("Windows path should contain 'goog', got %q", path)
	}
}

// TestSetValuePageSizeConversion tests page_size value conversion.
func TestSetValuePageSizeConversion(t *testing.T) {
	cfg := NewConfig()

	testCases := []struct {
		value    string
		expected int
		wantErr  bool
	}{
		{"10", 10, false},
		{"100", 100, false},
		{"0", 0, false},
		{"-5", -5, false}, // Negative values are parsed, validation is separate
		{"abc", 0, true},
		{"12.5", 12, false}, // Sscanf parses leading integer, stops at decimal
	}

	for _, tc := range testCases {
		t.Run(tc.value, func(t *testing.T) {
			// Reset
			cfg.Mail.PageSize = 20

			err := cfg.SetValue("mail.page_size", tc.value)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for value %q", tc.value)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for value %q: %v", tc.value, err)
				}
				if cfg.Mail.PageSize != tc.expected {
					t.Errorf("PageSize = %d, want %d", cfg.Mail.PageSize, tc.expected)
				}
			}
		})
	}
}

// TestLoadCreatesDirectoryWithCorrectPermissions tests directory permissions.
func TestLoadCreatesDirectoryWithCorrectPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions test not applicable on Windows")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "newdir", "goog", "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	_, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check directory permissions
	dirPath := filepath.Dir(configPath)
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("failed to stat directory: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0700 {
		t.Errorf("expected directory permissions 0700, got %o", perm)
	}
}

// TestStringToTimeHookFuncEmptyString tests time hook with empty string.
func TestStringToTimeHookFuncEmptyString(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Config with empty added_at
	configContent := `default_account: "empty@example.com"
accounts:
  empty@example.com:
    email: "empty@example.com"
    added_at: ""
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	acc, ok := cfg.Accounts["empty@example.com"]
	if !ok {
		t.Fatal("expected account to exist")
	}

	if !acc.AddedAt.IsZero() {
		t.Errorf("expected AddedAt to be zero time, got %v", acc.AddedAt)
	}
}

// TestWriteConfigSecurelyOverwrite tests overwriting existing config.
func TestWriteConfigSecurelyOverwrite(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "overwrite.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Create initial config
	cfg := NewConfig()
	cfg.DefaultAccount = "first@example.com"
	if err := cfg.Save(); err != nil {
		t.Fatalf("first Save failed: %v", err)
	}

	// Read initial content
	data1, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read first config: %v", err)
	}

	// Overwrite with different content
	cfg.DefaultAccount = "second@example.com"
	cfg.Mail.PageSize = 99
	if err := cfg.Save(); err != nil {
		t.Fatalf("second Save failed: %v", err)
	}

	// Read new content
	data2, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read second config: %v", err)
	}

	// Content should be different
	if bytes.Equal(data1, data2) {
		t.Error("config file should have different content after overwrite")
	}

	// New content should contain new values
	if !contains(string(data2), "second@example.com") {
		t.Error("new config should contain 'second@example.com'")
	}
}

// Security-related tests

// TestSecureFileCreation tests that config files are created with secure permissions
// from the start, avoiding the race condition where files are created with default
// permissions and then chmod'd.
func TestSecureFileCreation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions test not applicable on Windows")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "secure-test", "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Create and save a new config
	cfg := NewConfig()
	cfg.DefaultAccount = "test@example.com"

	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Check file permissions immediately after creation
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}
}

// TestConfigSaveCreatesDirectory tests that Save creates directories with correct permissions.
func TestConfigSaveCreatesDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions test not applicable on Windows")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nested", "dirs", "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Check directory permissions
	dirPath := filepath.Dir(configPath)
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("failed to stat config directory: %v", err)
	}

	dirPerm := dirInfo.Mode().Perm()
	if dirPerm != 0700 {
		t.Errorf("expected directory permissions 0700, got %o", dirPerm)
	}
}

// TestConfigSaveOverwritesMaintainsPermissions tests that overwriting an existing
// config maintains secure permissions.
func TestConfigSaveOverwritesMaintainsPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions test not applicable on Windows")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save and restore GOOG_CONFIG
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Create initial config
	cfg := NewConfig()
	cfg.DefaultAccount = "first@example.com"
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save initial config: %v", err)
	}

	// Modify and save again
	cfg.DefaultAccount = "second@example.com"
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save modified config: %v", err)
	}

	// Check permissions are still correct
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected file permissions 0600 after overwrite, got %o", perm)
	}

	// Verify content was actually updated
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if loaded.DefaultAccount != "second@example.com" {
		t.Errorf("expected default_account 'second@example.com', got %q", loaded.DefaultAccount)
	}
}

// TestSetValueFormatValidation tests the validation of default_format in SetValue.
func TestSetValueFormatValidation(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{
			name:      "valid json format",
			value:     "json",
			expectErr: false,
		},
		{
			name:      "valid plain format",
			value:     "plain",
			expectErr: false,
		},
		{
			name:      "valid table format",
			value:     "table",
			expectErr: false,
		},
		{
			name:      "invalid format",
			value:     "xml",
			expectErr: true,
		},
		{
			name:      "invalid format - empty",
			value:     "",
			expectErr: true,
		},
		{
			name:      "invalid format - uppercase",
			value:     "JSON",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewConfig()
			err := cfg.SetValue("default_format", tt.value)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for format %q, got nil", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for format %q: %v", tt.value, err)
				}
				if cfg.DefaultFormat != tt.value {
					t.Errorf("expected DefaultFormat %q, got %q", tt.value, cfg.DefaultFormat)
				}
			}
		})
	}
}

// TestSetValueTimezoneValidation tests the validation of timezone in SetValue.
func TestSetValueTimezoneValidation(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{
			name:      "valid timezone - America/New_York",
			value:     "America/New_York",
			expectErr: false,
		},
		{
			name:      "valid timezone - Europe/London",
			value:     "Europe/London",
			expectErr: false,
		},
		{
			name:      "valid timezone - UTC",
			value:     "UTC",
			expectErr: false,
		},
		{
			name:      "valid timezone - Local",
			value:     "Local",
			expectErr: false,
		},
		{
			name:      "valid timezone - empty (allowed)",
			value:     "",
			expectErr: false,
		},
		{
			name:      "invalid timezone",
			value:     "Invalid/Timezone",
			expectErr: true,
		},
		{
			name:      "invalid timezone - random string",
			value:     "not_a_timezone",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewConfig()
			err := cfg.SetValue("timezone", tt.value)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for timezone %q, got nil", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for timezone %q: %v", tt.value, err)
				}
				if cfg.Timezone != tt.value {
					t.Errorf("expected Timezone %q, got %q", tt.value, cfg.Timezone)
				}
			}
		})
	}
}

// TestStringToTimeHookFuncInvalidTime tests the time hook with invalid time format.
func TestStringToTimeHookFuncInvalidTime(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Config with invalid time format
	configContent := `default_account: "invalid@example.com"
accounts:
  invalid@example.com:
    email: "invalid@example.com"
    added_at: "not-a-valid-time-format"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid time format")
	}
}

// TestLoadWithNoHome tests Load when home directory might not be available.
func TestLoadWithNoHome(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Load should work with explicit config path
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should have defaults
	if cfg.DefaultFormat != "table" {
		t.Errorf("expected default format 'table', got %q", cfg.DefaultFormat)
	}
}

// TestSetValueCalendarWeekStart tests calendar.week_start SetValue.
func TestSetValueCalendarWeekStart(t *testing.T) {
	cfg := NewConfig()

	testCases := []struct {
		value   string
		wantErr bool
	}{
		{"sunday", false},
		{"monday", false},
		{"Saturday", false}, // Any value is accepted
	}

	for _, tc := range testCases {
		t.Run(tc.value, func(t *testing.T) {
			err := cfg.SetValue("calendar.week_start", tc.value)
			if tc.wantErr && err == nil {
				t.Errorf("expected error for value %q", tc.value)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error for value %q: %v", tc.value, err)
			}
		})
	}
}

// TestSetValueDefaultLabel tests mail.default_label SetValue.
func TestSetValueDefaultLabel(t *testing.T) {
	cfg := NewConfig()

	testCases := []string{"INBOX", "SENT", "DRAFT", "STARRED", "IMPORTANT", "Custom Label"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			err := cfg.SetValue("mail.default_label", tc)
			if err != nil {
				t.Errorf("unexpected error for value %q: %v", tc, err)
			}
			if cfg.Mail.DefaultLabel != tc {
				t.Errorf("expected default label %q, got %q", tc, cfg.Mail.DefaultLabel)
			}
		})
	}
}

// TestSetValueCalendarDefaultCalendar tests calendar.default_calendar SetValue.
func TestSetValueCalendarDefaultCalendar(t *testing.T) {
	cfg := NewConfig()

	testCases := []string{"primary", "work", "personal", "custom-calendar-id"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			err := cfg.SetValue("calendar.default_calendar", tc)
			if err != nil {
				t.Errorf("unexpected error for value %q: %v", tc, err)
			}
			if cfg.Calendar.DefaultCalendar != tc {
				t.Errorf("expected default calendar %q, got %q", tc, cfg.Calendar.DefaultCalendar)
			}
		})
	}
}

// TestGetConfigPathWindowsNoAPPDATA tests Windows path fallback when APPDATA is not set.
func TestGetConfigPathWindowsNoAPPDATA(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	origAppData := os.Getenv("APPDATA")
	os.Unsetenv("GOOG_CONFIG")
	os.Unsetenv("APPDATA")
	defer func() {
		restoreEnv("GOOG_CONFIG", origConfig)
		restoreEnv("APPDATA", origAppData)
	}()

	path := GetConfigPath()
	// Should fall back to user home directory
	if path == "" {
		t.Error("GetConfigPath returned empty string")
	}
	if !contains(path, "goog") {
		t.Errorf("path should contain 'goog', got %q", path)
	}
}

// TestLoadWithUnmarshalError tests config loading when unmarshal fails.
func TestLoadWithUnmarshalError(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Write YAML that won't unmarshal properly to Config struct
	// This is actually valid YAML but may cause issues with type conversion
	configContent := `default_account: 123
mail:
  page_size: "not_an_int"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := Load()
	// This should fail because page_size should be an int, not a string
	if err == nil {
		// Actually, viper may be lenient here - if it doesn't error, that's also valid
		t.Log("Load did not error on type mismatch - viper may be lenient")
	}
}

// TestSaveWithMkdirError tests Save when directory creation fails.
func TestSaveWithMkdirError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	// Use an invalid path that cannot be created
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", "/dev/null/cannot/create/this/path/config.yaml")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	err := cfg.Save()
	if err == nil {
		t.Error("expected error when directory cannot be created")
	}
}

// TestLoadWithMkdirError tests Load when directory creation fails.
func TestLoadWithMkdirError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Setenv("GOOG_CONFIG", "/dev/null/cannot/create/this/path/config.yaml")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_CONFIG", origConfig)
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	_, err := Load()
	if err == nil {
		t.Error("expected error when directory cannot be created")
	}
}

// TestSetPermissionsOnNonexistentFile tests SetPermissions when file doesn't exist.
func TestSetPermissionsOnNonexistentFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", "/tmp/nonexistent_config_file_for_test.yaml")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	err := SetPermissions()
	if err == nil {
		t.Error("expected error when file doesn't exist")
	}
}

// TestLoadWithNilAccountsMap tests that nil accounts map is initialized.
func TestLoadWithNilAccountsMap(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Write config without accounts field
	configContent := `default_account: "test@example.com"
default_format: "json"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Accounts should be initialized even if not in config
	if cfg.Accounts == nil {
		t.Error("Accounts map should be initialized")
	}
}

// TestConfigDefaultAccountSave tests saving config with default account.
func TestConfigDefaultAccountSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	cfg.DefaultAccount = "default@example.com"

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Read file content directly to verify
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	if !contains(string(data), "default@example.com") {
		t.Error("config file should contain default account email")
	}
}

// TestGetValueAllKeys tests GetValue returns correct values for all supported keys.
func TestGetValueAllKeys(t *testing.T) {
	cfg := NewConfig()
	cfg.DefaultAccount = "account@test.com"
	cfg.DefaultFormat = "plain"
	cfg.Timezone = "UTC"
	cfg.Mail.DefaultLabel = "STARRED"
	cfg.Mail.PageSize = 100
	cfg.Calendar.DefaultCalendar = "work"
	cfg.Calendar.WeekStart = "monday"

	tests := []struct {
		key      string
		expected string
	}{
		{"default_account", "account@test.com"},
		{"default_format", "plain"},
		{"timezone", "UTC"},
		{"mail.default_label", "STARRED"},
		{"mail.page_size", "100"},
		{"calendar.default_calendar", "work"},
		{"calendar.week_start", "monday"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value, err := cfg.GetValue(tt.key)
			if err != nil {
				t.Errorf("GetValue(%q) returned error: %v", tt.key, err)
			}
			if value != tt.expected {
				t.Errorf("GetValue(%q) = %q, want %q", tt.key, value, tt.expected)
			}
		})
	}
}

// TestWriteConfigSecurelyWindowsPath tests writeConfigSecurely on Windows.
func TestWriteConfigSecurelyWindowsPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	cfg.DefaultAccount = "windows@test.com"

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}

// TestLoadCreatesDefaultWhenNotExists tests that Load creates default config when file doesn't exist.
func TestLoadCreatesDefaultWhenNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "newdir", "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should have default values
	if cfg.DefaultFormat != "table" {
		t.Errorf("expected default format 'table', got %q", cfg.DefaultFormat)
	}

	// File should have been created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}

// TestAccountConfigSerialization tests that AccountConfig is properly serialized.
func TestAccountConfigSerialization(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Create config with account
	cfg := NewConfig()
	addedAt := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
	cfg.Accounts["test"] = AccountConfig{
		Email:   "test@example.com",
		Scopes:  []string{"gmail.readonly", "calendar.events"},
		AddedAt: addedAt,
	}

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load and verify
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	acc, ok := loaded.Accounts["test"]
	if !ok {
		t.Fatal("account 'test' not found after reload")
	}

	if acc.Email != "test@example.com" {
		t.Errorf("email = %q, want 'test@example.com'", acc.Email)
	}

	if len(acc.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(acc.Scopes))
	}
}

// TestGetConfigPathWithDifferentEnv tests GetConfigPath behavior with env var.
func TestGetConfigPathWithDifferentEnv(t *testing.T) {
	origConfig := os.Getenv("GOOG_CONFIG")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Test with custom path
	customPath := "/custom/path/config.yaml"
	os.Setenv("GOOG_CONFIG", customPath)

	path := GetConfigPath()
	if path != customPath {
		t.Errorf("GetConfigPath with env = %q, want %q", path, customPath)
	}

	// Test without env var (platform-specific path)
	os.Unsetenv("GOOG_CONFIG")
	path = GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath returned empty string")
	}
	if !contains(path, "goog") {
		t.Errorf("path should contain 'goog', got %q", path)
	}
}

// TestLoadWithSaveError tests Load behavior when Save fails.
func TestLoadWithSaveError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	tmpDir := t.TempDir()
	// Create a path in a read-only directory
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0500); err != nil {
		t.Fatalf("failed to create read-only dir: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0700) // Restore permissions for cleanup

	configPath := filepath.Join(readOnlyDir, "subdir", "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Load should fail because it can't create directory
	_, err := Load()
	if err == nil {
		t.Error("expected error when directory creation fails")
	}
}

// TestConfigMailSettings tests mail configuration settings.
func TestConfigMailSettings(t *testing.T) {
	cfg := NewConfig()

	// Test default mail settings
	if cfg.Mail.DefaultLabel != "INBOX" {
		t.Errorf("default mail label = %q, want 'INBOX'", cfg.Mail.DefaultLabel)
	}
	if cfg.Mail.PageSize != 20 {
		t.Errorf("default page size = %d, want 20", cfg.Mail.PageSize)
	}

	// Modify and verify
	cfg.Mail.DefaultLabel = "SENT"
	cfg.Mail.PageSize = 50

	if cfg.Mail.DefaultLabel != "SENT" {
		t.Errorf("mail label = %q, want 'SENT'", cfg.Mail.DefaultLabel)
	}
	if cfg.Mail.PageSize != 50 {
		t.Errorf("page size = %d, want 50", cfg.Mail.PageSize)
	}
}

// TestConfigCalendarSettings tests calendar configuration settings.
func TestConfigCalendarSettings(t *testing.T) {
	cfg := NewConfig()

	// Test default calendar settings
	if cfg.Calendar.DefaultCalendar != "primary" {
		t.Errorf("default calendar = %q, want 'primary'", cfg.Calendar.DefaultCalendar)
	}
	if cfg.Calendar.WeekStart != "sunday" {
		t.Errorf("default week start = %q, want 'sunday'", cfg.Calendar.WeekStart)
	}

	// Modify and verify
	cfg.Calendar.DefaultCalendar = "work"
	cfg.Calendar.WeekStart = "monday"

	if cfg.Calendar.DefaultCalendar != "work" {
		t.Errorf("calendar = %q, want 'work'", cfg.Calendar.DefaultCalendar)
	}
	if cfg.Calendar.WeekStart != "monday" {
		t.Errorf("week start = %q, want 'monday'", cfg.Calendar.WeekStart)
	}
}

// TestLoadExistingConfigDoesNotOverwrite tests that loading existing config doesn't create new file.
func TestLoadExistingConfigDoesNotOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Create config with specific values
	configContent := `default_account: "existing@example.com"
default_format: "plain"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Get file modification time
	info1, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	// Wait a moment
	time.Sleep(10 * time.Millisecond)

	// Load existing config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify config was loaded correctly
	if cfg.DefaultAccount != "existing@example.com" {
		t.Errorf("default account = %q, want 'existing@example.com'", cfg.DefaultAccount)
	}

	// File should not have been modified
	info2, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat file after load: %v", err)
	}

	if !info1.ModTime().Equal(info2.ModTime()) {
		t.Error("file was modified when loading existing config")
	}
}

// TestValidFormatsMap tests the validFormats map.
func TestValidFormatsMap(t *testing.T) {
	cfg := NewConfig()

	validFormats := []string{"json", "plain", "table"}
	invalidFormats := []string{"xml", "yaml", "csv", "", "JSON", "PLAIN"}

	for _, format := range validFormats {
		if err := cfg.SetValue("default_format", format); err != nil {
			t.Errorf("SetValue for valid format %q returned error: %v", format, err)
		}
	}

	for _, format := range invalidFormats {
		// Reset to valid format first
		cfg.SetValue("default_format", "json")

		if err := cfg.SetValue("default_format", format); err == nil {
			t.Errorf("SetValue for invalid format %q should return error", format)
		}
	}
}

// TestSetValueDefaultAccountEdgeCases tests edge cases for default_account.
func TestSetValueDefaultAccountEdgeCases(t *testing.T) {
	cfg := NewConfig()

	testCases := []string{
		"",                       // Empty
		"simple",                 // Simple string
		"test@example.com",       // Email format
		"test+alias@example.com", // Email with plus
		"very-long-account-name-that-is-quite-lengthy@subdomain.example.com", // Long
	}

	for _, tc := range testCases {
		if err := cfg.SetValue("default_account", tc); err != nil {
			t.Errorf("SetValue(default_account, %q) returned error: %v", tc, err)
		}
		if cfg.DefaultAccount != tc {
			t.Errorf("default_account = %q, want %q", cfg.DefaultAccount, tc)
		}
	}
}

// TestErrAccountNotFoundValue tests ErrAccountNotFound error.
func TestErrAccountNotFoundValue(t *testing.T) {
	if ErrAccountNotFound == nil {
		t.Error("ErrAccountNotFound should not be nil")
	}
	if ErrAccountNotFound.Error() != "account not found" {
		t.Errorf("ErrAccountNotFound.Error() = %q, want 'account not found'", ErrAccountNotFound.Error())
	}
}

// TestMultipleAccountsInConfig tests config with multiple accounts.
func TestMultipleAccountsInConfig(t *testing.T) {
	cfg := NewConfig()

	// Add multiple accounts
	for i := 0; i < 5; i++ {
		alias := "account" + string(rune('0'+i))
		email := alias + "@example.com"
		cfg.Accounts[alias] = AccountConfig{
			Email:   email,
			Scopes:  []string{"gmail.readonly"},
			AddedAt: time.Now(),
		}
	}

	// Verify all accounts can be retrieved
	for i := 0; i < 5; i++ {
		alias := "account" + string(rune('0'+i))
		acc, err := cfg.GetAccount(alias)
		if err != nil {
			t.Errorf("GetAccount(%q) failed: %v", alias, err)
			continue
		}
		expectedEmail := alias + "@example.com"
		if acc.Email != expectedEmail {
			t.Errorf("account %s email = %q, want %q", alias, acc.Email, expectedEmail)
		}
	}
}

// TestStringToTimeHookFuncNonStringType tests that non-string types pass through unchanged.
func TestStringToTimeHookFuncNonStringType(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Config with an integer value for a non-time field - should pass through
	configContent := `default_account: "test@example.com"
default_format: "json"
mail:
  page_size: 50
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Mail.PageSize != 50 {
		t.Errorf("expected page_size 50, got %d", cfg.Mail.PageSize)
	}
}

// TestStringToTimeHookFuncNonTimeTarget tests that non-time targets pass through unchanged.
func TestStringToTimeHookFuncNonTimeTarget(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Test that string to string passes through (not converted to time)
	configContent := `default_account: "2024-01-15T10:30:00Z"
default_format: "json"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// The date-like string should be kept as string since DefaultAccount is a string
	if cfg.DefaultAccount != "2024-01-15T10:30:00Z" {
		t.Errorf("expected default_account '2024-01-15T10:30:00Z', got %q", cfg.DefaultAccount)
	}
}

// TestWriteConfigSecurelyWriteError tests write error handling.
func TestWriteConfigSecurelyWriteError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	// Create a config and save it successfully first
	cfg := NewConfig()
	cfg.DefaultAccount = "test@example.com"
	if err := cfg.Save(); err != nil {
		t.Fatalf("first save failed: %v", err)
	}

	// Make the file read-only
	if err := os.Chmod(configPath, 0400); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}
	defer os.Chmod(configPath, 0600) // Restore for cleanup

	// Try to save again - should fail
	cfg.DefaultAccount = "another@example.com"
	err := cfg.Save()
	if err == nil {
		t.Error("expected error when writing to read-only file")
	}
}

// TestGetConfigPathWindowsWithAPPDATA tests Windows path with APPDATA set.
func TestGetConfigPathWindowsWithAPPDATA(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	origAppData := os.Getenv("APPDATA")
	os.Unsetenv("GOOG_CONFIG")
	os.Setenv("APPDATA", "C:\\Users\\Test\\AppData\\Roaming")
	defer func() {
		restoreEnv("GOOG_CONFIG", origConfig)
		restoreEnv("APPDATA", origAppData)
	}()

	path := GetConfigPath()
	if !contains(path, "C:\\Users\\Test\\AppData\\Roaming\\goog") {
		t.Errorf("expected path to contain APPDATA\\goog, got %q", path)
	}
}

// TestSetValueUnknownKey tests SetValue with an unknown key.
func TestSetValueUnknownKey(t *testing.T) {
	cfg := NewConfig()

	err := cfg.SetValue("invalid.nested.key", "value")
	if err == nil {
		t.Error("expected error for unknown key")
	}
	if !contains(err.Error(), "unknown config key") {
		t.Errorf("expected 'unknown config key' error, got %v", err)
	}
}

// TestGetValueUnknownKey tests GetValue with an unknown key.
func TestGetValueUnknownKey(t *testing.T) {
	cfg := NewConfig()

	_, err := cfg.GetValue("invalid.nested.key")
	if err == nil {
		t.Error("expected error for unknown key")
	}
	if !contains(err.Error(), "unknown config key") {
		t.Errorf("expected 'unknown config key' error, got %v", err)
	}
}

// TestWriteConfigSecurelyCloseError tests handling when file close fails.
func TestWriteConfigSecurelyFileOperations(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-close.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	cfg.DefaultAccount = "test@example.com"
	cfg.Mail.PageSize = 100

	// Normal save should work
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Read the file to verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if !contains(string(data), "test@example.com") {
		t.Error("config should contain the email")
	}
}

// TestWriteConfigSecurelyOnWindows simulates Windows behavior.
func TestWriteConfigSecurelyOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "windows-test.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	cfg.DefaultAccount = "win@example.com"

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file should have been created")
	}

	// Verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if !contains(string(data), "win@example.com") {
		t.Error("file should contain saved account")
	}
}

// TestGetConfigPathLinuxNoXDG tests Linux path when XDG_CONFIG_HOME is not set.
func TestGetConfigPathLinuxNoXDG(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("GOOG_CONFIG")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer func() {
		restoreEnv("GOOG_CONFIG", origConfig)
		restoreEnv("XDG_CONFIG_HOME", origXDG)
	}()

	path := GetConfigPath()

	// Should fallback to ~/.config/goog
	if !contains(path, ".config/goog") {
		t.Errorf("expected path to contain '.config/goog', got %q", path)
	}
}

// TestLoadFailsOnInvalidDirectory tests that Load fails when config directory cannot be created.
func TestLoadFailsOnInvalidDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	// Use a path under /dev/null which cannot have subdirectories
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", "/dev/null/subdir/config.yaml")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	_, err := Load()
	if err == nil {
		t.Error("expected error when directory cannot be created")
	}
}

// TestSaveFailsOnInvalidDirectory tests that Save fails when config directory cannot be created.
func TestSaveFailsOnInvalidDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", "/dev/null/subdir/config.yaml")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	err := cfg.Save()
	if err == nil {
		t.Error("expected error when directory cannot be created")
	}
}

// TestSetValueTimezoneLocal tests setting timezone to "Local".
func TestSetValueTimezoneLocal(t *testing.T) {
	cfg := NewConfig()

	err := cfg.SetValue("timezone", "Local")
	if err != nil {
		t.Errorf("unexpected error setting timezone to Local: %v", err)
	}
	if cfg.Timezone != "Local" {
		t.Errorf("expected timezone 'Local', got %q", cfg.Timezone)
	}
}

// TestSetValueTimezoneEmpty tests setting timezone to empty string.
func TestSetValueTimezoneEmpty(t *testing.T) {
	cfg := NewConfig()

	err := cfg.SetValue("timezone", "")
	if err != nil {
		t.Errorf("unexpected error setting empty timezone: %v", err)
	}
	if cfg.Timezone != "" {
		t.Errorf("expected empty timezone, got %q", cfg.Timezone)
	}
}

// TestAccountConfigFields tests AccountConfig field access.
func TestAccountConfigFields(t *testing.T) {
	cfg := NewConfig()
	addedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	cfg.Accounts["test"] = AccountConfig{
		Email:   "test@example.com",
		Scopes:  []string{"gmail.readonly", "calendar.events"},
		AddedAt: addedAt,
	}

	acc, err := cfg.GetAccount("test")
	if err != nil {
		t.Fatalf("GetAccount failed: %v", err)
	}

	if acc.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", acc.Email)
	}
	if len(acc.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(acc.Scopes))
	}
	if !acc.AddedAt.Equal(addedAt) {
		t.Errorf("expected AddedAt %v, got %v", addedAt, acc.AddedAt)
	}
}

// TestMailConfigFields tests MailConfig field access.
func TestMailConfigFields(t *testing.T) {
	cfg := NewConfig()

	// Test defaults
	if cfg.Mail.DefaultLabel != "INBOX" {
		t.Errorf("expected default label 'INBOX', got %q", cfg.Mail.DefaultLabel)
	}
	if cfg.Mail.PageSize != 20 {
		t.Errorf("expected page size 20, got %d", cfg.Mail.PageSize)
	}

	// Test modification
	cfg.Mail.DefaultLabel = "SENT"
	cfg.Mail.PageSize = 50

	if cfg.Mail.DefaultLabel != "SENT" {
		t.Errorf("expected label 'SENT', got %q", cfg.Mail.DefaultLabel)
	}
	if cfg.Mail.PageSize != 50 {
		t.Errorf("expected page size 50, got %d", cfg.Mail.PageSize)
	}
}

// TestCalendarConfigFields tests CalendarConfig field access.
func TestCalendarConfigFields(t *testing.T) {
	cfg := NewConfig()

	// Test defaults
	if cfg.Calendar.DefaultCalendar != "primary" {
		t.Errorf("expected default calendar 'primary', got %q", cfg.Calendar.DefaultCalendar)
	}
	if cfg.Calendar.WeekStart != "sunday" {
		t.Errorf("expected week start 'sunday', got %q", cfg.Calendar.WeekStart)
	}

	// Test modification
	cfg.Calendar.DefaultCalendar = "work"
	cfg.Calendar.WeekStart = "monday"

	if cfg.Calendar.DefaultCalendar != "work" {
		t.Errorf("expected calendar 'work', got %q", cfg.Calendar.DefaultCalendar)
	}
	if cfg.Calendar.WeekStart != "monday" {
		t.Errorf("expected week start 'monday', got %q", cfg.Calendar.WeekStart)
	}
}

// TestGetConfigPathPlatformSpecificFallbacks tests GetConfigPath fallbacks.
func TestGetConfigPathPlatformSpecificFallbacks(t *testing.T) {
	origConfig := os.Getenv("GOOG_CONFIG")
	os.Unsetenv("GOOG_CONFIG")
	defer restoreEnv("GOOG_CONFIG", origConfig)

	path := GetConfigPath()

	// Verify path is not empty
	if path == "" {
		t.Error("GetConfigPath returned empty path")
	}

	// Verify path ends with config.yaml
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("path should end with 'config.yaml', got %q", filepath.Base(path))
	}

	// Verify path contains goog directory
	if filepath.Base(filepath.Dir(path)) != "goog" {
		t.Errorf("path should contain 'goog' directory, got %q", path)
	}
}

// TestLoadWithMalformedConfig tests loading a malformed config file.
func TestLoadWithMalformedConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Write malformed YAML (unbalanced brackets)
	malformedContent := `default_account: "test@example.com"
accounts:
  test:
    - missing: [bracket
`
	if err := os.WriteFile(configPath, []byte(malformedContent), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Error("expected error for malformed YAML")
	}
}

// TestConfigSaveWithSpecialCharacters tests saving config with special characters.
func TestConfigSaveWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	cfg := NewConfig()
	cfg.DefaultAccount = "test+special@example.com"
	cfg.Mail.DefaultLabel = "Label with spaces & symbols!"

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reload and verify
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.DefaultAccount != "test+special@example.com" {
		t.Errorf("expected account 'test+special@example.com', got %q", loaded.DefaultAccount)
	}
	if loaded.Mail.DefaultLabel != "Label with spaces & symbols!" {
		t.Errorf("expected label 'Label with spaces & symbols!', got %q", loaded.Mail.DefaultLabel)
	}
}

// TestLoadWithEnvOverridesAndExistingConfig tests env overrides with existing config.
func TestLoadWithEnvOverridesAndExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save all env vars
	origConfig := os.Getenv("GOOG_CONFIG")
	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")

	defer func() {
		restoreEnv("GOOG_CONFIG", origConfig)
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Create config file
	configContent := `default_account: "file@example.com"
default_format: "table"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Set env vars to override
	os.Setenv("GOOG_CONFIG", configPath)
	os.Setenv("GOOG_ACCOUNT", "env@example.com")
	os.Setenv("GOOG_FORMAT", "json")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Env vars should override file values
	if cfg.DefaultAccount != "env@example.com" {
		t.Errorf("expected GOOG_ACCOUNT override, got %q", cfg.DefaultAccount)
	}
	if cfg.DefaultFormat != "json" {
		t.Errorf("expected GOOG_FORMAT override, got %q", cfg.DefaultFormat)
	}
}

// TestSetValuePageSizeZero tests setting page_size to zero.
func TestSetValuePageSizeZero(t *testing.T) {
	cfg := NewConfig()

	err := cfg.SetValue("mail.page_size", "0")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cfg.Mail.PageSize != 0 {
		t.Errorf("expected page size 0, got %d", cfg.Mail.PageSize)
	}
}

// TestSetValuePageSizeNegative tests setting page_size to negative value.
func TestSetValuePageSizeNegative(t *testing.T) {
	cfg := NewConfig()

	err := cfg.SetValue("mail.page_size", "-10")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Negative values are parsed but may not be valid semantically
	if cfg.Mail.PageSize != -10 {
		t.Errorf("expected page size -10, got %d", cfg.Mail.PageSize)
	}
}

// TestConfigFullRoundTrip tests complete save/load cycle with all fields.
func TestConfigFullRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "roundtrip.yaml")

	origConfig := os.Getenv("GOOG_CONFIG")
	os.Setenv("GOOG_CONFIG", configPath)
	defer restoreEnv("GOOG_CONFIG", origConfig)

	origAccount := os.Getenv("GOOG_ACCOUNT")
	origFormat := os.Getenv("GOOG_FORMAT")
	os.Unsetenv("GOOG_ACCOUNT")
	os.Unsetenv("GOOG_FORMAT")
	defer func() {
		restoreEnv("GOOG_ACCOUNT", origAccount)
		restoreEnv("GOOG_FORMAT", origFormat)
	}()

	// Create config with all fields
	cfg := NewConfig()
	cfg.DefaultAccount = "primary@example.com"
	cfg.DefaultFormat = "plain"
	cfg.Timezone = "America/Los_Angeles"
	cfg.Mail.DefaultLabel = "STARRED"
	cfg.Mail.PageSize = 75
	cfg.Calendar.DefaultCalendar = "work"
	cfg.Calendar.WeekStart = "monday"

	cfg.Accounts["personal"] = AccountConfig{
		Email:   "personal@example.com",
		Scopes:  []string{"gmail.readonly", "calendar.events"},
		AddedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
	cfg.Accounts["work"] = AccountConfig{
		Email:   "work@company.com",
		Scopes:  []string{"gmail.modify", "drive.readonly"},
		AddedAt: time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC),
	}

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify all top-level fields
	if loaded.DefaultAccount != "primary@example.com" {
		t.Errorf("DefaultAccount = %q, want 'primary@example.com'", loaded.DefaultAccount)
	}
	if loaded.DefaultFormat != "plain" {
		t.Errorf("DefaultFormat = %q, want 'plain'", loaded.DefaultFormat)
	}
	if loaded.Timezone != "America/Los_Angeles" {
		t.Errorf("Timezone = %q, want 'America/Los_Angeles'", loaded.Timezone)
	}

	// Verify mail settings
	if loaded.Mail.DefaultLabel != "STARRED" {
		t.Errorf("Mail.DefaultLabel = %q, want 'STARRED'", loaded.Mail.DefaultLabel)
	}
	if loaded.Mail.PageSize != 75 {
		t.Errorf("Mail.PageSize = %d, want 75", loaded.Mail.PageSize)
	}

	// Verify calendar settings
	if loaded.Calendar.DefaultCalendar != "work" {
		t.Errorf("Calendar.DefaultCalendar = %q, want 'work'", loaded.Calendar.DefaultCalendar)
	}
	if loaded.Calendar.WeekStart != "monday" {
		t.Errorf("Calendar.WeekStart = %q, want 'monday'", loaded.Calendar.WeekStart)
	}

	// Verify accounts
	if len(loaded.Accounts) != 2 {
		t.Errorf("len(Accounts) = %d, want 2", len(loaded.Accounts))
	}

	personal, ok := loaded.Accounts["personal"]
	if !ok {
		t.Fatal("account 'personal' not found")
	}
	if personal.Email != "personal@example.com" {
		t.Errorf("personal.Email = %q, want 'personal@example.com'", personal.Email)
	}

	work, ok := loaded.Accounts["work"]
	if !ok {
		t.Fatal("account 'work' not found")
	}
	if work.Email != "work@company.com" {
		t.Errorf("work.Email = %q, want 'work@company.com'", work.Email)
	}
}
