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

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
