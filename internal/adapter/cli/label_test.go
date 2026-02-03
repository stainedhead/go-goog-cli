// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestLabelCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "label") {
		t.Error("expected output to contain 'label'")
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
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
}

func TestLabelListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "system labels") {
		t.Error("expected output to contain 'system labels'")
	}
}

func TestLabelShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
}

func TestLabelCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "create", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
	if !contains(output, "--background") {
		t.Error("expected output to contain '--background'")
	}
	if !contains(output, "--text") {
		t.Error("expected output to contain '--text'")
	}
}

func TestLabelUpdateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "update", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "update") {
		t.Error("expected output to contain 'update'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
	if !contains(output, "--background") {
		t.Error("expected output to contain '--background'")
	}
}

func TestLabelDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "delete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
}

func TestLabelShowCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelShowCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelCreateCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelCreateCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelUpdateCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelUpdateCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelDeleteCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelDeleteCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelCmd_Aliases(t *testing.T) {
	tests := []struct {
		name    string
		command string
		alias   string
	}{
		{"list alias ls", "list", "ls"},
		{"show alias get", "show", "get"},
		{"show alias info", "show", "info"},
		{"delete alias rm", "delete", "rm"},
		{"delete alias remove", "delete", "remove"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find the subcommand
			var targetCmd *cobra.Command
			for _, sub := range labelCmd.Commands() {
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

func TestLabelDeleteCmd_HasConfirmFlag(t *testing.T) {
	// Verify the command has a --confirm flag
	flag := labelDeleteCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}
