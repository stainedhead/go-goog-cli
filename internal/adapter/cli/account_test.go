// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/cobra"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

func TestAccountCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "account") {
		t.Error("expected output to contain 'account'")
	}
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "add") {
		t.Error("expected output to contain 'add'")
	}
	if !contains(output, "remove") {
		t.Error("expected output to contain 'remove'")
	}
	if !contains(output, "switch") {
		t.Error("expected output to contain 'switch'")
	}
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "rename") {
		t.Error("expected output to contain 'rename'")
	}
}

func TestAccountListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
}

func TestAccountAddCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "add", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "add") {
		t.Error("expected output to contain 'add'")
	}
	if !contains(output, "--scopes") {
		t.Error("expected output to contain '--scopes'")
	}
}

func TestAccountRemoveCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "remove"}) // Missing required argument

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestAccountSwitchCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "switch"}) // Missing required argument

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestAccountRenameCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Test with no arguments
	cmd.SetArgs([]string{"account", "rename"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing arguments")
	}

	// Test with one argument
	cmd.SetArgs([]string{"account", "rename", "old"})
	err = cmd.Execute()
	if err == nil {
		t.Error("expected error for missing second argument")
	}
}

func TestAccountShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
}

func TestAccountCmd_Aliases(t *testing.T) {
	tests := []struct {
		alias   string
		command string
	}{
		{"ls", "list"},
		{"rm", "remove"},
		{"use", "switch"},
		{"current", "show"},
		{"mv", "rename"},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"account", tt.alias, "--help"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error for alias %s: %v", tt.alias, err)
			}
		})
	}
}

func TestAccountCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":   false,
		"add":    false,
		"remove": false,
		"switch": false,
		"show":   false,
		"rename": false,
	}

	for _, sub := range accountCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with accountCmd", name)
		}
	}
}

func TestAccountRemoveCmd_HasArgsRequirement(t *testing.T) {
	if accountRemoveCmd.Args == nil {
		t.Error("expected Args to be set on remove command")
	}
}

func TestAccountSwitchCmd_HasArgsRequirement(t *testing.T) {
	if accountSwitchCmd.Args == nil {
		t.Error("expected Args to be set on switch command")
	}
}

func TestAccountRenameCmd_HasArgsRequirement(t *testing.T) {
	if accountRenameCmd.Args == nil {
		t.Error("expected Args to be set on rename command")
	}
}

func TestAccountAddCmd_HasScopesFlag(t *testing.T) {
	flag := accountAddCmd.Flag("scopes")
	if flag == nil {
		t.Error("expected --scopes flag to be defined on add command")
	}
}

func TestAccountRemoveCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"alias"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"alias", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountRemoveCmd.Args(accountRemoveCmd, tt.args)
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

func TestAccountSwitchCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"alias"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"alias", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountSwitchCmd.Args(accountSwitchCmd, tt.args)
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

func TestAccountRenameCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"old"},
			expectErr: true,
		},
		{
			name:      "two args",
			args:      []string{"old", "new"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"old", "new", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountRenameCmd.Args(accountRenameCmd, tt.args)
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

func TestAccountAddCmd_ArgsValidation(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "no args - valid",
			args:      []string{},
			expectErr: false,
		},
		{
			name:      "one arg - valid",
			args:      []string{"alias"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"alias", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountAddCmd.Args(accountAddCmd, tt.args)
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

func TestOutputAccountsTable(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "default",
			Email:     "user@example.com",
			IsDefault: true,
			Added:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Alias:     "work",
			Email:     "work@company.com",
			IsDefault: false,
			Added:     time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC),
		},
	}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsTable(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Check for table headers
	if !contains(output, "ALIAS") {
		t.Error("expected output to contain 'ALIAS'")
	}
	if !contains(output, "EMAIL") {
		t.Error("expected output to contain 'EMAIL'")
	}
	if !contains(output, "DEFAULT") {
		t.Error("expected output to contain 'DEFAULT'")
	}
	// Check for account data
	if !contains(output, "default") {
		t.Error("expected output to contain 'default'")
	}
	if !contains(output, "user@example.com") {
		t.Error("expected output to contain 'user@example.com'")
	}
	if !contains(output, "work") {
		t.Error("expected output to contain 'work'")
	}
}

func TestOutputAccountsPlain(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "default",
			Email:     "user@example.com",
			IsDefault: true,
			Added:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Alias:     "work",
			Email:     "work@company.com",
			IsDefault: false,
			Added:     time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC),
		},
	}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsPlain(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Check for account data in plain format
	if !contains(output, "default:") {
		t.Error("expected output to contain 'default:'")
	}
	if !contains(output, "user@example.com") {
		t.Error("expected output to contain 'user@example.com'")
	}
	if !contains(output, "(default)") {
		t.Error("expected output to contain '(default)' marker")
	}
	if !contains(output, "work:") {
		t.Error("expected output to contain 'work:'")
	}
}

func TestOutputAccountsTable_EmptyList(t *testing.T) {
	accounts := []*accountuc.Account{}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsTable(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Should still have headers
	if !contains(output, "ALIAS") {
		t.Error("expected output to contain 'ALIAS' header even for empty list")
	}
}

func TestOutputAccountsPlain_EmptyList(t *testing.T) {
	accounts := []*accountuc.Account{}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsPlain(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Should be empty or minimal for no accounts
	if output != "" {
		t.Errorf("expected empty output for empty account list, got: %s", output)
	}
}
