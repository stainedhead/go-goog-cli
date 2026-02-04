// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
)

func TestConfigCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "config") {
		t.Error("expected output to contain 'config'")
	}
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "set") {
		t.Error("expected output to contain 'set'")
	}
	if !contains(output, "get") {
		t.Error("expected output to contain 'get'")
	}
}

func TestConfigShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
}

func TestConfigSetCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Test with no arguments
	cmd.SetArgs([]string{"config", "set"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing arguments")
	}

	// Test with one argument
	cmd.SetArgs([]string{"config", "set", "key"})
	err = cmd.Execute()
	if err == nil {
		t.Error("expected error for missing value argument")
	}
}

func TestConfigGetCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Test with no arguments
	cmd.SetArgs([]string{"config", "get"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestConfigSetCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "set", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "set") {
		t.Error("expected output to contain 'set'")
	}
	if !contains(output, "default_account") {
		t.Error("expected output to contain 'default_account'")
	}
	if !contains(output, "default_format") {
		t.Error("expected output to contain 'default_format'")
	}
}

func TestConfigGetCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "get", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "get") {
		t.Error("expected output to contain 'get'")
	}
}

func TestConfigPathCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "path", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "path") {
		t.Error("expected output to contain 'path'")
	}
}

func TestConfigCmd_Aliases(t *testing.T) {
	tests := []struct {
		alias   string
		command string
	}{
		{"view", "show"},
		{"list", "show"},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(configCmd)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"config", tt.alias, "--help"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error for alias %s: %v", tt.alias, err)
			}
		})
	}
}

func TestConfigCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"show": false,
		"set":  false,
		"get":  false,
		"path": false,
	}

	for _, sub := range configCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with configCmd", name)
		}
	}
}

func TestConfigSetCmd_HasArgsRequirement(t *testing.T) {
	if configSetCmd.Args == nil {
		t.Error("expected Args to be set on set command")
	}
}

func TestConfigGetCmd_HasArgsRequirement(t *testing.T) {
	if configGetCmd.Args == nil {
		t.Error("expected Args to be set on get command")
	}
}

func TestConfigSetCmd_ArgsValidation(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "no args",
			args:      []string{},
			expectErr: true,
		},
		{
			name:      "one arg",
			args:      []string{"key"},
			expectErr: true,
		},
		{
			name:      "two args",
			args:      []string{"key", "value"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"key", "value", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configSetCmd.Args(configSetCmd, tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestConfigGetCmd_ArgsValidation(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "no args",
			args:      []string{},
			expectErr: true,
		},
		{
			name:      "one arg",
			args:      []string{"key"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"key", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configGetCmd.Args(configGetCmd, tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// =============================================================================
// Tests for RunE functions - these test the actual command execution
// =============================================================================

func TestRunConfigPath(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runConfigPath(cmd, []string{})
	if err != nil {
		t.Fatalf("runConfigPath failed: %v", err)
	}

	output := buf.String()
	// Should output a path
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRunConfigShow(t *testing.T) {
	// Setup: Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yaml"
	t.Setenv("GOOG_CONFIG", configPath)

	// Create a test config using config package
	// This will initialize a default config
	_, err := setupTestConfig(configPath)
	if err != nil {
		t.Fatalf("failed to setup test config: %v", err)
	}

	tests := []struct {
		name           string
		expectedFields []string
	}{
		{
			name: "show default config",
			expectedFields: []string{
				"default_account:",
				"default_format:",
				"timezone:",
				"mail:",
				"default_label:",
				"page_size:",
				"calendar:",
				"default_calendar:",
				"week_start:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runConfigShow(cmd, []string{})
			if err != nil {
				t.Fatalf("runConfigShow failed: %v", err)
			}

			output := buf.String()
			for _, field := range tt.expectedFields {
				if !contains(output, field) {
					t.Errorf("expected output to contain %q, got: %s", field, output)
				}
			}
		})
	}
}

func TestRunConfigShow_WithAccounts(t *testing.T) {
	// Setup: Create a temporary config with accounts
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yaml"
	t.Setenv("GOOG_CONFIG", configPath)

	cfg, err := setupTestConfig(configPath)
	if err != nil {
		t.Fatalf("failed to setup test config: %v", err)
	}

	// Add test accounts
	cfg.Accounts["work"] = accountConfigForTest("work@example.com", []string{"gmail.readonly"})
	cfg.Accounts["personal"] = accountConfigForTest("personal@example.com", []string{"gmail.modify"})
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err = runConfigShow(cmd, []string{})
	if err != nil {
		t.Fatalf("runConfigShow failed: %v", err)
	}

	output := buf.String()

	// Check for account fields
	expectedFields := []string{
		"accounts:",
		"personal:",
		"email: personal@example.com",
		"work:",
		"email: work@example.com",
		"scopes:",
	}

	for _, field := range expectedFields {
		if !contains(output, field) {
			t.Errorf("expected output to contain %q, got: %s", field, output)
		}
	}
}

func TestRunConfigGet(t *testing.T) {
	// Setup: Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yaml"
	t.Setenv("GOOG_CONFIG", configPath)

	cfg, err := setupTestConfig(configPath)
	if err != nil {
		t.Fatalf("failed to setup test config: %v", err)
	}

	// Set some test values
	cfg.DefaultFormat = "json"
	cfg.Timezone = "America/New_York"
	cfg.Mail.PageSize = 50
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	tests := []struct {
		name          string
		key           string
		expectedValue string
		expectError   bool
	}{
		{
			name:          "get default_format",
			key:           "default_format",
			expectedValue: "json",
			expectError:   false,
		},
		{
			name:          "get timezone",
			key:           "timezone",
			expectedValue: "America/New_York",
			expectError:   false,
		},
		{
			name:          "get mail.page_size",
			key:           "mail.page_size",
			expectedValue: "50",
			expectError:   false,
		},
		{
			name:          "get mail.default_label",
			key:           "mail.default_label",
			expectedValue: "INBOX",
			expectError:   false,
		},
		{
			name:          "get calendar.default_calendar",
			key:           "calendar.default_calendar",
			expectedValue: "primary",
			expectError:   false,
		},
		{
			name:          "get calendar.week_start",
			key:           "calendar.week_start",
			expectedValue: "sunday",
			expectError:   false,
		},
		{
			name:        "get invalid key",
			key:         "invalid.key",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runConfigGet(cmd, []string{tt.key})

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				output := buf.String()
				if !contains(output, tt.expectedValue) {
					t.Errorf("expected output to contain %q, got: %q", tt.expectedValue, output)
				}
			}
		})
	}
}

func TestRunConfigSet(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "set default_account",
			key:         "default_account",
			value:       "test@example.com",
			expectError: false,
		},
		{
			name:        "set default_format to json",
			key:         "default_format",
			value:       "json",
			expectError: false,
		},
		{
			name:        "set default_format to plain",
			key:         "default_format",
			value:       "plain",
			expectError: false,
		},
		{
			name:        "set default_format to table",
			key:         "default_format",
			value:       "table",
			expectError: false,
		},
		{
			name:        "set default_format to invalid value",
			key:         "default_format",
			value:       "xml",
			expectError: true,
			errorMsg:    "invalid format",
		},
		{
			name:        "set timezone",
			key:         "timezone",
			value:       "America/New_York",
			expectError: false,
		},
		{
			name:        "set timezone to Local",
			key:         "timezone",
			value:       "Local",
			expectError: false,
		},
		{
			name:        "set timezone to invalid value",
			key:         "timezone",
			value:       "Invalid/Timezone",
			expectError: true,
			errorMsg:    "invalid timezone",
		},
		{
			name:        "set mail.default_label",
			key:         "mail.default_label",
			value:       "SENT",
			expectError: false,
		},
		{
			name:        "set mail.page_size",
			key:         "mail.page_size",
			value:       "100",
			expectError: false,
		},
		{
			name:        "set mail.page_size to invalid value",
			key:         "mail.page_size",
			value:       "not-a-number",
			expectError: true,
			errorMsg:    "invalid page_size",
		},
		{
			name:        "set calendar.default_calendar",
			key:         "calendar.default_calendar",
			value:       "work@example.com",
			expectError: false,
		},
		{
			name:        "set calendar.week_start",
			key:         "calendar.week_start",
			value:       "monday",
			expectError: false,
		},
		{
			name:        "set unknown key",
			key:         "unknown.key",
			value:       "value",
			expectError: true,
			errorMsg:    "unknown config key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: Create a fresh temporary config for each test
			tempDir := t.TempDir()
			configPath := tempDir + "/config.yaml"
			t.Setenv("GOOG_CONFIG", configPath)

			_, err := setupTestConfig(configPath)
			if err != nil {
				t.Fatalf("failed to setup test config: %v", err)
			}

			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err = runConfigSet(cmd, []string{tt.key, tt.value})

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain %q, got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				output := buf.String()
				if !contains(output, "Set") {
					t.Errorf("expected success message in output, got: %q", output)
				}
				if !contains(output, tt.key) {
					t.Errorf("expected output to contain key %q, got: %q", tt.key, output)
				}
			}
		})
	}
}

func TestRunConfigSet_PersistsValue(t *testing.T) {
	// Setup: Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yaml"
	t.Setenv("GOOG_CONFIG", configPath)

	_, err := setupTestConfig(configPath)
	if err != nil {
		t.Fatalf("failed to setup test config: %v", err)
	}

	// Set a value
	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err = runConfigSet(cmd, []string{"default_format", "json"})
	if err != nil {
		t.Fatalf("runConfigSet failed: %v", err)
	}

	// Verify the value was persisted by reading it back
	buf.Reset()
	err = runConfigGet(cmd, []string{"default_format"})
	if err != nil {
		t.Fatalf("runConfigGet failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "json") {
		t.Errorf("expected to read back 'json', got: %q", output)
	}
}

func TestRunConfigSet_MultipleValues(t *testing.T) {
	// Setup: Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yaml"
	t.Setenv("GOOG_CONFIG", configPath)

	_, err := setupTestConfig(configPath)
	if err != nil {
		t.Fatalf("failed to setup test config: %v", err)
	}

	cmd := &cobra.Command{Use: "test"}

	// Set multiple values
	values := map[string]string{
		"default_format":            "json",
		"timezone":                  "America/Los_Angeles",
		"mail.page_size":            "75",
		"calendar.default_calendar": "work@example.com",
	}

	for key, value := range values {
		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err = runConfigSet(cmd, []string{key, value})
		if err != nil {
			t.Fatalf("runConfigSet failed for %s: %v", key, err)
		}
	}

	// Verify all values were set correctly
	for key, expectedValue := range values {
		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err = runConfigGet(cmd, []string{key})
		if err != nil {
			t.Fatalf("runConfigGet failed for %s: %v", key, err)
		}

		output := buf.String()
		if !contains(output, expectedValue) {
			t.Errorf("expected %s to be %q, got: %q", key, expectedValue, output)
		}
	}
}

// Helper functions

func setupTestConfig(configPath string) (*config.Config, error) {
	// Load config (will create default if it doesn't exist)
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Helper to create AccountConfig for tests
func accountConfigForTest(email string, scopes []string) config.AccountConfig {
	return config.AccountConfig{
		Email:  email,
		Scopes: scopes,
	}
}
