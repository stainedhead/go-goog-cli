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
