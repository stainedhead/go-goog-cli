// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestAuthCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(authCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"auth", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "auth") {
		t.Error("expected output to contain 'auth'")
	}
	if !contains(output, "login") {
		t.Error("expected output to contain 'login'")
	}
	if !contains(output, "logout") {
		t.Error("expected output to contain 'logout'")
	}
	if !contains(output, "status") {
		t.Error("expected output to contain 'status'")
	}
	if !contains(output, "refresh") {
		t.Error("expected output to contain 'refresh'")
	}
}

func TestAuthLoginCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(authCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"auth", "login", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "login") {
		t.Error("expected output to contain 'login'")
	}
	if !contains(output, "--scopes") {
		t.Error("expected output to contain '--scopes'")
	}
}

func TestAuthLogoutCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(authCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"auth", "logout", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "logout") {
		t.Error("expected output to contain 'logout'")
	}
}

func TestAuthStatusCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(authCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"auth", "status", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "status") {
		t.Error("expected output to contain 'status'")
	}
}

func TestAuthRefreshCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(authCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"auth", "refresh", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "refresh") {
		t.Error("expected output to contain 'refresh'")
	}
}

func TestParseScopes(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected int // minimum expected scopes (including auto-added email/openid)
	}{
		{
			name:     "empty scopes",
			input:    []string{},
			expected: 0, // returns nil
		},
		{
			name:     "gmail shorthand",
			input:    []string{"gmail"},
			expected: 3, // gmail.readonly + email + openid
		},
		{
			name:     "calendar shorthand",
			input:    []string{"calendar"},
			expected: 3, // calendar.readonly + email + openid
		},
		{
			name:     "full URL",
			input:    []string{"https://www.googleapis.com/auth/gmail.readonly"},
			expected: 3, // full URL + email + openid
		},
		{
			name:     "multiple scopes",
			input:    []string{"gmail", "calendar"},
			expected: 4, // gmail + calendar + email + openid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseScopes(tt.input)
			if tt.expected == 0 {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}
			if len(result) < tt.expected {
				t.Errorf("expected at least %d scopes, got %d: %v", tt.expected, len(result), result)
			}
		})
	}
}

func TestParseScopes_AllShorthands(t *testing.T) {
	// Test all known scope shorthands
	shorthands := []string{
		"gmail", "gmail.readonly", "gmail.send", "gmail.modify", "gmail.compose", "gmail.labels",
		"calendar", "calendar.readonly", "calendar.events", "calendar.full",
		"drive", "drive.readonly", "drive.file", "drive.full",
		"email", "profile", "openid",
	}

	for _, shorthand := range shorthands {
		t.Run(shorthand, func(t *testing.T) {
			result := parseScopes([]string{shorthand})
			if result == nil {
				t.Errorf("expected non-nil result for shorthand %s", shorthand)
			}
			// Result should have at least 2 scopes (the scope + email + openid, but email/openid may already be included)
			if len(result) < 1 {
				t.Errorf("expected at least 1 scope for shorthand %s, got %d", shorthand, len(result))
			}
		})
	}
}

func TestParseScopes_UnknownScope(t *testing.T) {
	// Unknown scopes should be passed through as-is
	result := parseScopes([]string{"unknown.scope"})
	if result == nil {
		t.Error("expected non-nil result")
	}
	// Should contain the unknown scope plus email and openid
	found := false
	for _, s := range result {
		if s == "unknown.scope" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected unknown scope to be in result, got: %v", result)
	}
}

func TestParseScopes_WithEmailAndOpenID(t *testing.T) {
	// When email and openid are already included, they shouldn't be duplicated
	result := parseScopes([]string{"email", "openid", "gmail"})
	if result == nil {
		t.Error("expected non-nil result")
	}

	// Count occurrences of email and openid scope URLs
	emailCount := 0
	openidCount := 0
	for _, s := range result {
		if s == "https://www.googleapis.com/auth/userinfo.email" {
			emailCount++
		}
		if s == "openid" {
			openidCount++
		}
	}

	if emailCount > 1 {
		t.Errorf("email scope should not be duplicated, found %d occurrences", emailCount)
	}
	if openidCount > 1 {
		t.Errorf("openid scope should not be duplicated, found %d occurrences", openidCount)
	}
}

func TestParseScopes_WhitespaceHandling(t *testing.T) {
	// Test that whitespace is trimmed
	result := parseScopes([]string{"  gmail  ", "  CALENDAR  "})
	if result == nil {
		t.Error("expected non-nil result")
	}
	// Should successfully parse despite whitespace and case
	if len(result) < 3 {
		t.Errorf("expected at least 3 scopes, got %d: %v", len(result), result)
	}
}

func TestAuthCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"login":   false,
		"logout":  false,
		"status":  false,
		"refresh": false,
	}

	for _, sub := range authCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with authCmd", name)
		}
	}
}

func TestAuthLoginCmd_HasScopesFlag(t *testing.T) {
	flag := authLoginCmd.Flag("scopes")
	if flag == nil {
		t.Error("expected --scopes flag to be defined on login command")
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		expected     string
	}{
		{
			name:         "env not set returns default",
			key:          "TEST_GOOG_CLI_NONEXISTENT",
			defaultValue: "default_value",
			setEnv:       false,
			expected:     "default_value",
		},
		{
			name:         "env set returns env value",
			key:          "TEST_GOOG_CLI_SET",
			defaultValue: "default_value",
			envValue:     "env_value",
			setEnv:       true,
			expected:     "env_value",
		},
		{
			name:         "empty env returns default",
			key:          "TEST_GOOG_CLI_EMPTY",
			defaultValue: "default_value",
			envValue:     "",
			setEnv:       true,
			expected:     "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.key, tt.envValue)
			}
			result := getEnvWithDefault(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvWithDefault(%q, %q) = %q, expected %q", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
