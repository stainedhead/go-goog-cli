// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog", Short: rootCmd.Short, Long: rootCmd.Long}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "goog") {
		t.Error("expected output to contain 'goog'")
	}
	// Check for parts of the description
	if !contains(output, "Google") {
		t.Error("expected output to contain 'Google'")
	}
}

func TestRootCmd_HasGlobalFlags(t *testing.T) {
	flags := []string{"account", "format", "quiet", "verbose", "config"}

	for _, flagName := range flags {
		flag := rootCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("expected global flag --%s to be defined", flagName)
		}
	}
}

func TestRootCmd_HasVersionSubcommand(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected version subcommand to be registered")
	}
}

func TestVersionCmd_Execute(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(versionCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"version"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "goog") {
		t.Error("expected output to contain 'goog'")
	}
}

func TestRootCmd_SubcommandsRegistered(t *testing.T) {
	// Check that all expected top-level commands are registered
	expectedCommands := []string{
		"version",
		"auth",
		"account",
		"config",
		"mail",
		"cal",
		"label",
		"draft",
		"thread",
	}

	registeredCommands := make(map[string]bool)
	for _, sub := range rootCmd.Commands() {
		registeredCommands[sub.Name()] = true
	}

	for _, cmd := range expectedCommands {
		if !registeredCommands[cmd] {
			t.Errorf("expected subcommand %s to be registered with rootCmd", cmd)
		}
	}
}

func TestRootCmd_FormatFlagDefaultValue(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("format")
	if flag == nil {
		t.Fatal("expected format flag to exist")
	}
	if flag.DefValue != "table" {
		t.Errorf("expected format flag default to be 'table', got %s", flag.DefValue)
	}
}

func TestRootCmd_QuietFlagDefaultValue(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("quiet")
	if flag == nil {
		t.Fatal("expected quiet flag to exist")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected quiet flag default to be 'false', got %s", flag.DefValue)
	}
}

func TestRootCmd_VerboseFlagDefaultValue(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("verbose")
	if flag == nil {
		t.Fatal("expected verbose flag to exist")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected verbose flag default to be 'false', got %s", flag.DefValue)
	}
}
