// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "config") {
		t.Error("expected output to contain 'config'")
	}
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "set") {
		t.Error("expected output to contain 'set'")
	}
	if !contains(output, "get") {
		t.Error("expected output to contain 'get'")
	}
}

func TestConfigShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
}

func TestConfigSetCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Test with no arguments
	cmd.SetArgs([]string{"config", "set"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing arguments")
	}

	// Test with one argument
	cmd.SetArgs([]string{"config", "set", "key"})
	err = cmd.Execute()
	if err == nil {
		t.Error("expected error for missing value argument")
	}
}

func TestConfigGetCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Test with no arguments
	cmd.SetArgs([]string{"config", "get"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestConfigSetCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "set", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "set") {
		t.Error("expected output to contain 'set'")
	}
	if !contains(output, "default_account") {
		t.Error("expected output to contain 'default_account'")
	}
	if !contains(output, "default_format") {
		t.Error("expected output to contain 'default_format'")
	}
}

func TestConfigGetCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "get", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "get") {
		t.Error("expected output to contain 'get'")
	}
}

func TestConfigPathCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(configCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "path", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "path") {
		t.Error("expected output to contain 'path'")
	}
}

func TestConfigCmd_Aliases(t *testing.T) {
	tests := []struct {
		alias   string
		command string
	}{
		{"view", "show"},
		{"list", "show"},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(configCmd)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"config", tt.alias, "--help"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error for alias %s: %v", tt.alias, err)
			}
		})
	}
}
