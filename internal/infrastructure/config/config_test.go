// Package config provides configuration management for the goog CLI application.
package config

import (
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
