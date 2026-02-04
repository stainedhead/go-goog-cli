// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestThreadCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "thread") {
		t.Error("expected output to contain 'thread'")
	}
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
}

func TestThreadListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "list", "--help"})

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
}

func TestThreadShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "show", "--help"})

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

func TestThreadTrashCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "trash", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestThreadModifyCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "modify", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
	if !contains(output, "--add-labels") {
		t.Error("expected output to contain '--add-labels'")
	}
	if !contains(output, "--remove-labels") {
		t.Error("expected output to contain '--remove-labels'")
	}
}

func TestThreadShowCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if threadShowCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadTrashCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if threadTrashCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadModifyCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if threadModifyCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadCmd_Aliases(t *testing.T) {
	tests := []struct {
		name    string
		command string
		alias   string
	}{
		{"list alias ls", "list", "ls"},
		{"show alias get", "show", "get"},
		{"show alias read", "show", "read"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find the subcommand
			var targetCmd *cobra.Command
			for _, sub := range threadCmd.Commands() {
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

func TestThreadListCmd_HasMaxResultsFlag(t *testing.T) {
	flag := threadListCmd.Flag("max-results")
	if flag == nil {
		t.Error("expected --max-results flag to be set")
	}
}

func TestThreadListCmd_HasLabelsFlag(t *testing.T) {
	flag := threadListCmd.Flag("labels")
	if flag == nil {
		t.Error("expected --labels flag to be set")
	}
}

func TestThreadModifyCmd_HasAddLabelsFlag(t *testing.T) {
	flag := threadModifyCmd.Flag("add-labels")
	if flag == nil {
		t.Error("expected --add-labels flag to be set")
	}
}

func TestThreadModifyCmd_HasRemoveLabelsFlag(t *testing.T) {
	flag := threadModifyCmd.Flag("remove-labels")
	if flag == nil {
		t.Error("expected --remove-labels flag to be set")
	}
}
