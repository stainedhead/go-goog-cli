// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
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
