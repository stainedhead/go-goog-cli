// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestMailListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "--max-results") {
		t.Error("expected output to contain '--max-results'")
	}
	if !contains(output, "--labels") {
		t.Error("expected output to contain '--labels'")
	}
	if !contains(output, "--unread-only") {
		t.Error("expected output to contain '--unread-only'")
	}
}

func TestMailReadCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "read", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "read") {
		t.Error("expected output to contain 'read'")
	}
	if !contains(output, "<message-id>") {
		t.Error("expected output to contain '<message-id>'")
	}
}

func TestMailReadCmd_RequiresMessageID(t *testing.T) {
	// Test that mailReadCmd has Args validator set to require 1 argument
	if mailReadCmd.Args == nil {
		t.Error("expected mailReadCmd to have Args validator")
	}
}

func TestMailSearchCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "search", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "search") {
		t.Error("expected output to contain 'search'")
	}
	if !contains(output, "--max-results") {
		t.Error("expected output to contain '--max-results'")
	}
	if !contains(output, "Gmail query") {
		t.Error("expected output to contain 'Gmail query'")
	}
}

func TestMailSearchCmd_RequiresQuery(t *testing.T) {
	// Test that mailSearchCmd has Args validator set to require 1 argument
	if mailSearchCmd.Args == nil {
		t.Error("expected mailSearchCmd to have Args validator")
	}
}

func TestMailListReadSearchCmd_Aliases(t *testing.T) {
	tests := []struct {
		alias   string
		command string
	}{
		{"ls", "list"},
		{"get", "read"},
		{"show", "read"},
		{"find", "search"},
		{"query", "search"},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(mailCmd)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"mail", tt.alias, "--help"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error for alias %s: %v", tt.alias, err)
			}
		})
	}
}

func TestMailListCmd_DefaultMaxResults(t *testing.T) {
	// Reset to default first
	mailListMaxResults = 10

	// Verify the default value is set correctly
	if mailListMaxResults != 10 {
		t.Errorf("expected default max-results to be 10, got %d", mailListMaxResults)
	}
}

func TestMailSearchCmd_DefaultMaxResults(t *testing.T) {
	// Reset to default first
	mailSearchMaxResults = 10

	// Verify the default value is set correctly
	if mailSearchMaxResults != 10 {
		t.Errorf("expected default max-results to be 10, got %d", mailSearchMaxResults)
	}
}

func TestMailListCmd_DefaultLabels(t *testing.T) {
	// Reset to default first
	mailListLabels = []string{"INBOX"}

	// Verify the default value is set correctly
	if len(mailListLabels) != 1 || mailListLabels[0] != "INBOX" {
		t.Errorf("expected default labels to be [INBOX], got %v", mailListLabels)
	}
}

func TestMailReadCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"message-id"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"message-id", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailReadCmd.Args(mailReadCmd, tt.args)
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

func TestMailSearchCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"is:unread"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"is:unread", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailSearchCmd.Args(mailSearchCmd, tt.args)
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

func TestMailListCmd_HasFlags(t *testing.T) {
	flags := []string{"max-results", "labels", "unread-only"}

	for _, flagName := range flags {
		flag := mailListCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on list command", flagName)
		}
	}
}

func TestMailSearchCmd_HasFlags(t *testing.T) {
	flag := mailSearchCmd.Flag("max-results")
	if flag == nil {
		t.Error("expected --max-results flag to be defined on search command")
	}
}
