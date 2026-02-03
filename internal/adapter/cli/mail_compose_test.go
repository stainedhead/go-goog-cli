// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestMailComposeCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "mail") {
		t.Error("expected output to contain 'mail'")
	}
	if !contains(output, "send") {
		t.Error("expected output to contain 'send'")
	}
	if !contains(output, "reply") {
		t.Error("expected output to contain 'reply'")
	}
	if !contains(output, "forward") {
		t.Error("expected output to contain 'forward'")
	}
}

func TestMailSendCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "send", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "send") {
		t.Error("expected output to contain 'send'")
	}
	if !contains(output, "--to") {
		t.Error("expected output to contain '--to'")
	}
	if !contains(output, "--cc") {
		t.Error("expected output to contain '--cc'")
	}
	if !contains(output, "--bcc") {
		t.Error("expected output to contain '--bcc'")
	}
	if !contains(output, "--subject") {
		t.Error("expected output to contain '--subject'")
	}
	if !contains(output, "--body") {
		t.Error("expected output to contain '--body'")
	}
	if !contains(output, "--html") {
		t.Error("expected output to contain '--html'")
	}
}

func TestMailReplyCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "reply", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "reply") {
		t.Error("expected output to contain 'reply'")
	}
	if !contains(output, "--body") {
		t.Error("expected output to contain '--body'")
	}
	if !contains(output, "--all") {
		t.Error("expected output to contain '--all'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestMailForwardCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "forward", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "forward") {
		t.Error("expected output to contain 'forward'")
	}
	if !contains(output, "--to") {
		t.Error("expected output to contain '--to'")
	}
	if !contains(output, "--body") {
		t.Error("expected output to contain '--body'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestMailSendCmd_RequiresToFlag(t *testing.T) {
	// Test that PreRunE validates the --to flag
	// We test the validation logic directly since Cobra flag parsing
	// behavior varies in test contexts

	// Reset flags
	mailSendTo = []string{}

	// Create a mock command to test the PreRunE function
	mockCmd := &cobra.Command{Use: "test"}

	// Invoke the PreRunE function directly
	if mailSendCmd.PreRunE != nil {
		err := mailSendCmd.PreRunE(mockCmd, []string{})
		if err == nil {
			t.Error("expected error when --to flag is missing")
		}
	} else {
		t.Error("mailSendCmd should have PreRunE defined")
	}
}

func TestMailReplyCmd_RequiresIDArg(t *testing.T) {
	// Test that Args validator requires exactly one argument
	// The reply command should require exactly 1 argument (message ID)
	if mailReplyCmd.Args == nil {
		t.Error("mailReplyCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailReplyCmd.Args(mailReplyCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailReplyCmd_RequiresBodyFlag(t *testing.T) {
	// Test that PreRunE validates the --body flag
	mailReplyBody = ""

	mockCmd := &cobra.Command{Use: "test"}

	if mailReplyCmd.PreRunE != nil {
		err := mailReplyCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when --body flag is missing")
		}
	} else {
		t.Error("mailReplyCmd should have PreRunE defined")
	}
}

func TestMailForwardCmd_RequiresIDArg(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailForwardCmd.Args == nil {
		t.Error("mailForwardCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailForwardCmd.Args(mailForwardCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailForwardCmd_RequiresToFlag(t *testing.T) {
	// Test that PreRunE validates the --to flag
	mailForwardTo = []string{}

	mockCmd := &cobra.Command{Use: "test"}

	if mailForwardCmd.PreRunE != nil {
		err := mailForwardCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when --to flag is missing")
		}
	} else {
		t.Error("mailForwardCmd should have PreRunE defined")
	}
}

func TestParseEmailRecipients(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "single email",
			input:    []string{"user@example.com"},
			expected: []string{"user@example.com"},
		},
		{
			name:     "multiple emails",
			input:    []string{"user1@example.com", "user2@example.com"},
			expected: []string{"user1@example.com", "user2@example.com"},
		},
		{
			name:     "with whitespace",
			input:    []string{"  user@example.com  "},
			expected: []string{"user@example.com"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseEmailRecipients(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d recipients, got %d", len(tt.expected), len(result))
				return
			}
			for i, email := range result {
				if email != tt.expected[i] {
					t.Errorf("expected %q at index %d, got %q", tt.expected[i], i, email)
				}
			}
		})
	}
}
