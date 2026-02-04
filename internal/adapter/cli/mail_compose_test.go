// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"strings"
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
		name      string
		input     []string
		expected  []string
		expectErr bool
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
		{
			name:      "invalid email",
			input:     []string{"notanemail"},
			expectErr: true,
		},
		{
			name:      "one valid one invalid",
			input:     []string{"valid@example.com", "invalid"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEmailRecipients(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %v, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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

func TestParseEmailRecipients_AdditionalCases(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		expected  []string
		expectErr bool
	}{
		{
			name:     "filter empty strings",
			input:    []string{"user@example.com", "", "  "},
			expected: []string{"user@example.com"},
		},
		{
			name:     "mixed whitespace and valid",
			input:    []string{"  ", "user1@example.com", "", "   user2@example.com   "},
			expected: []string{"user1@example.com", "user2@example.com"},
		},
		{
			name:     "all empty strings",
			input:    []string{"", "", ""},
			expected: []string{},
		},
		{
			name:     "tabs and spaces",
			input:    []string{"\tuser@example.com\t"},
			expected: []string{"user@example.com"},
		},
		{
			name:     "large list",
			input:    []string{"a@b.com", "c@d.com", "e@f.com", "g@h.com", "i@j.com"},
			expected: []string{"a@b.com", "c@d.com", "e@f.com", "g@h.com", "i@j.com"},
		},
		{
			name:      "email missing @",
			input:     []string{"userexample.com"},
			expectErr: true,
		},
		{
			name:      "email missing domain",
			input:     []string{"user@"},
			expectErr: true,
		},
		{
			name:      "email missing local part",
			input:     []string{"@example.com"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEmailRecipients(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %v, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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

func TestBuildReplySubject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain subject",
			input:    "Hello World",
			expected: "Re: Hello World",
		},
		{
			name:     "empty subject",
			input:    "",
			expected: "Re: ",
		},
		{
			name:     "already has Re:",
			input:    "Re: Hello World",
			expected: "Re: Hello World",
		},
		{
			name:     "lowercase re:",
			input:    "re: Hello World",
			expected: "re: Hello World",
		},
		{
			name:     "RE: uppercase",
			input:    "RE: Hello World",
			expected: "RE: Hello World",
		},
		{
			name:     "Re: with extra spaces",
			input:    "Re:  Hello World",
			expected: "Re:  Hello World",
		},
		{
			name:     "subject starting with Re but not prefix",
			input:    "Regarding your request",
			expected: "Re: Regarding your request",
		},
		{
			name:     "subject with just 'Re'",
			input:    "Re",
			expected: "Re: Re",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildReplySubject(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMailSendCmd_Validation(t *testing.T) {
	tests := []struct {
		name      string
		to        []string
		expectErr bool
	}{
		{
			name:      "empty to list",
			to:        []string{},
			expectErr: true,
		},
		{
			name:      "nil to list",
			to:        nil,
			expectErr: true,
		},
		{
			name:      "with recipients",
			to:        []string{"user@example.com"},
			expectErr: false,
		},
		{
			name:      "multiple recipients",
			to:        []string{"user1@example.com", "user2@example.com"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTo := mailSendTo
			mailSendTo = tt.to

			mockCmd := &cobra.Command{Use: "test"}

			err := mailSendCmd.PreRunE(mockCmd, []string{})

			mailSendTo = origTo

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

func TestMailReplyCmd_Validation(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		expectErr bool
	}{
		{
			name:      "empty body",
			body:      "",
			expectErr: true,
		},
		{
			name:      "with body",
			body:      "Thanks for your message!",
			expectErr: false,
		},
		{
			name:      "whitespace only body - still valid",
			body:      "   ",
			expectErr: false, // whitespace is technically non-empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origBody := mailReplyBody
			mailReplyBody = tt.body

			mockCmd := &cobra.Command{Use: "test"}

			err := mailReplyCmd.PreRunE(mockCmd, []string{"msg123"})

			mailReplyBody = origBody

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

func TestMailForwardCmd_Validation(t *testing.T) {
	tests := []struct {
		name      string
		to        []string
		expectErr bool
	}{
		{
			name:      "empty to list",
			to:        []string{},
			expectErr: true,
		},
		{
			name:      "nil to list",
			to:        nil,
			expectErr: true,
		},
		{
			name:      "with recipients",
			to:        []string{"user@example.com"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTo := mailForwardTo
			mailForwardTo = tt.to

			mockCmd := &cobra.Command{Use: "test"}

			err := mailForwardCmd.PreRunE(mockCmd, []string{"msg123"})

			mailForwardTo = origTo

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

func TestMailReplyCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailReplyCmd.Args(mailReplyCmd, tt.args)
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

func TestMailForwardCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailForwardCmd.Args(mailForwardCmd, tt.args)
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

func TestMailComposeCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"send":    false,
		"reply":   false,
		"forward": false,
	}

	for _, sub := range mailCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with mailCmd", name)
		}
	}
}

func TestMailSendCmd_HasRequiredFlags(t *testing.T) {
	flags := []string{"to", "cc", "bcc", "subject", "body", "html"}

	for _, flagName := range flags {
		flag := mailSendCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on send command", flagName)
		}
	}
}

func TestMailReplyCmd_HasRequiredFlags(t *testing.T) {
	flags := []string{"body", "all"}

	for _, flagName := range flags {
		flag := mailReplyCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on reply command", flagName)
		}
	}
}

func TestMailForwardCmd_HasRequiredFlags(t *testing.T) {
	flags := []string{"to", "body"}

	for _, flagName := range flags {
		flag := mailForwardCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on forward command", flagName)
		}
	}
}

// =============================================================================
// Edge Case and Error Path Tests
// =============================================================================

func TestParseEmailRecipients_ErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "email with spaces",
			input:     []string{"user @example.com"},
			expectErr: true,
			errMsg:    "invalid email address",
		},
		{
			name:      "email with multiple @",
			input:     []string{"user@@example.com"},
			expectErr: true,
			errMsg:    "invalid email address",
		},
		{
			name:      "email without TLD",
			input:     []string{"user@example"},
			expectErr: true,
			errMsg:    "invalid email address",
		},
		{
			name:      "just @ symbol",
			input:     []string{"@"},
			expectErr: true,
			errMsg:    "invalid email address",
		},
		{
			name:      "special characters only",
			input:     []string{"!@#$%"},
			expectErr: true,
			errMsg:    "invalid email address",
		},
		{
			name:      "very long email",
			input:     []string{"verylongemailaddressthatexceedsnormallimitsbutshouldbetested@verylongdomainname.com"},
			expectErr: false, // Should be valid
		},
		{
			name:      "email with plus sign (valid)",
			input:     []string{"user+tag@example.com"},
			expectErr: false,
		},
		{
			name:      "email with dots (valid)",
			input:     []string{"first.last@example.com"},
			expectErr: false,
		},
		{
			name:      "email with numbers (valid)",
			input:     []string{"user123@example123.com"},
			expectErr: false,
		},
		{
			name:      "mixed valid and invalid with whitespace",
			input:     []string{"valid@example.com", "  ", "invalid@"},
			expectErr: true,
			errMsg:    "invalid email address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEmailRecipients(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %v, got nil", tt.input)
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) == 0 && len(tt.input) > 0 {
					// Only fail if input was non-empty
					hasNonEmpty := false
					for _, s := range tt.input {
						if strings.TrimSpace(s) != "" {
							hasNonEmpty = true
							break
						}
					}
					if hasNonEmpty {
						t.Error("expected non-empty result")
					}
				}
			}
		})
	}
}

func TestBuildReplySubject_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mixed case Re:",
			input:    "rE: Test",
			expected: "rE: Test",
		},
		{
			name:     "Re: with leading/trailing spaces",
			input:    "  Re: Test  ",
			expected: "Re:   Re: Test  ", // Doesn't trim, just checks prefix
		},
		{
			name:     "multiple Re: prefixes",
			input:    "Re: Re: Test",
			expected: "Re: Re: Test",
		},
		{
			name:     "RE: in middle",
			input:    "Test RE: something",
			expected: "Re: Test RE: something",
		},
		{
			name:     "very long subject",
			input:    "This is a very long subject line that goes on and on and on and on and on and on",
			expected: "Re: This is a very long subject line that goes on and on and on and on and on and on",
		},
		{
			name:     "unicode characters",
			input:    "Hello 世界",
			expected: "Re: Hello 世界",
		},
		{
			name:     "special characters",
			input:    "Test: [IMPORTANT] #123",
			expected: "Re: Test: [IMPORTANT] #123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildReplySubject(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMailSendCmd_EmptyBodies(t *testing.T) {
	// Test that send can work with empty body
	origTo := mailSendTo
	origSubject := mailSendSubject
	origBody := mailSendBody

	mailSendTo = []string{"user@example.com"}
	mailSendSubject = "Empty Body Test"
	mailSendBody = "" // Empty body should be allowed

	mockCmd := &cobra.Command{Use: "test"}
	err := mailSendCmd.PreRunE(mockCmd, []string{})

	mailSendTo = origTo
	mailSendSubject = origSubject
	mailSendBody = origBody

	if err != nil {
		t.Errorf("unexpected error with empty body: %v", err)
	}
}

func TestMailSendCmd_EmptySubject(t *testing.T) {
	// Test that send can work with empty subject
	origTo := mailSendTo
	origSubject := mailSendSubject
	origBody := mailSendBody

	mailSendTo = []string{"user@example.com"}
	mailSendSubject = "" // Empty subject should be allowed
	mailSendBody = "Test body"

	mockCmd := &cobra.Command{Use: "test"}
	err := mailSendCmd.PreRunE(mockCmd, []string{})

	mailSendTo = origTo
	mailSendSubject = origSubject
	mailSendBody = origBody

	if err != nil {
		t.Errorf("unexpected error with empty subject: %v", err)
	}
}

func TestMailSendCmd_MultipleRecipients(t *testing.T) {
	// Test validation with many recipients
	tests := []struct {
		name      string
		to        []string
		cc        []string
		bcc       []string
		expectErr bool
	}{
		{
			name:      "many to recipients",
			to:        []string{"user1@example.com", "user2@example.com", "user3@example.com", "user4@example.com", "user5@example.com"},
			expectErr: false,
		},
		{
			name:      "to, cc, and bcc",
			to:        []string{"user1@example.com"},
			cc:        []string{"user2@example.com", "user3@example.com"},
			bcc:       []string{"user4@example.com"},
			expectErr: false,
		},
		{
			name:      "invalid in to list",
			to:        []string{"valid@example.com", "invalid"},
			expectErr: false, // PreRunE only checks length, not validity
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTo := mailSendTo
			mailSendTo = tt.to

			mockCmd := &cobra.Command{Use: "test"}
			err := mailSendCmd.PreRunE(mockCmd, []string{})

			mailSendTo = origTo

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

func TestMailReplyCmd_EmptyMessageID(t *testing.T) {
	// Test Args validator with empty string (not zero args)
	err := mailReplyCmd.Args(mailReplyCmd, []string{""})
	// Should pass Args validation (it only checks count), but would fail in execution
	if err != nil {
		t.Errorf("Args validator should accept empty string: %v", err)
	}
}

func TestMailForwardCmd_EmptyMessageID(t *testing.T) {
	// Test Args validator with empty string (not zero args)
	err := mailForwardCmd.Args(mailForwardCmd, []string{""})
	// Should pass Args validation (it only checks count), but would fail in execution
	if err != nil {
		t.Errorf("Args validator should accept empty string: %v", err)
	}
}

func TestParseEmailRecipients_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		expectErr bool
	}{
		{
			name:      "single character local part",
			input:     []string{"a@example.com"},
			expectErr: false,
		},
		{
			name:      "email with subdomain",
			input:     []string{"user@mail.example.com"},
			expectErr: false,
		},
		{
			name:      "email with hyphen in domain",
			input:     []string{"user@ex-ample.com"},
			expectErr: false,
		},
		{
			name:      "email with underscore in local",
			input:     []string{"user_name@example.com"},
			expectErr: false,
		},
		{
			name:      "empty string in middle of list",
			input:     []string{"user1@example.com", "", "user2@example.com"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEmailRecipients(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Count non-empty inputs
				expectedCount := 0
				for _, s := range tt.input {
					if strings.TrimSpace(s) != "" {
						expectedCount++
					}
				}
				if len(result) != expectedCount {
					t.Errorf("expected %d results, got %d", expectedCount, len(result))
				}
			}
		})
	}
}
