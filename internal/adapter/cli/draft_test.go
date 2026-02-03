// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestDraftCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "draft") {
		t.Error("expected output to contain 'draft'")
	}
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
	if !contains(output, "update") {
		t.Error("expected output to contain 'update'")
	}
	if !contains(output, "send") {
		t.Error("expected output to contain 'send'")
	}
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
}

func TestDraftListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "--limit") {
		t.Error("expected output to contain '--limit'")
	}
}

func TestDraftShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestDraftCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "create", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
	if !contains(output, "--to") {
		t.Error("expected output to contain '--to'")
	}
	if !contains(output, "--subject") {
		t.Error("expected output to contain '--subject'")
	}
	if !contains(output, "--body") {
		t.Error("expected output to contain '--body'")
	}
}

func TestDraftUpdateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "update", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "update") {
		t.Error("expected output to contain 'update'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestDraftSendCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "send", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "send") {
		t.Error("expected output to contain 'send'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestDraftDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "delete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestDraftShowCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftShowCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftUpdateCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftUpdateCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftSendCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftSendCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftDeleteCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftDeleteCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftCmd_Aliases(t *testing.T) {
	tests := []struct {
		name    string
		command string
		alias   string
	}{
		{"list alias ls", "list", "ls"},
		{"show alias get", "show", "get"},
		{"show alias read", "show", "read"},
		{"delete alias rm", "delete", "rm"},
		{"delete alias remove", "delete", "remove"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find the subcommand
			var targetCmd *cobra.Command
			for _, sub := range draftCmd.Commands() {
				if sub.Use[:len(tt.command)] == tt.command || sub.Use == tt.command {
					targetCmd = sub
					break
				}
			}

			if targetCmd == nil {
				t.Fatalf("command %s not found", tt.command)
			}

			// Check alias exists
			found := false
			for _, alias := range targetCmd.Aliases {
				if alias == tt.alias {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected alias %s for command %s, got aliases: %v",
					tt.alias, tt.command, targetCmd.Aliases)
			}
		})
	}
}
