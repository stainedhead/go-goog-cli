// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestCalInstancesCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "instances", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !containsStr(output, "instances") {
		t.Error("expected output to contain 'instances'")
	}
	if !containsStr(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
	if !containsStr(output, "--start") {
		t.Error("expected output to contain '--start'")
	}
	if !containsStr(output, "--end") {
		t.Error("expected output to contain '--end'")
	}
	if !containsStr(output, "--max-results") {
		t.Error("expected output to contain '--max-results'")
	}
	if !containsStr(output, "<event-id>") {
		t.Error("expected output to contain '<event-id>'")
	}
}

func TestCalInstancesCmd_RequiresEventID(t *testing.T) {
	if calInstancesCmd.Args == nil {
		t.Error("calInstancesCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := calInstancesCmd.Args(calInstancesCmd, []string{})
	if err == nil {
		t.Error("expected error when event ID is missing")
	}
}

func TestCalInstancesCmd_AcceptsEventID(t *testing.T) {
	if calInstancesCmd.Args == nil {
		t.Error("calInstancesCmd should have Args validator defined")
		return
	}

	// Test with one arg - should pass
	err := calInstancesCmd.Args(calInstancesCmd, []string{"abc123"})
	if err != nil {
		t.Errorf("unexpected error when event ID is provided: %v", err)
	}
}

func TestCalInstancesCmd_HasCalendarFlag(t *testing.T) {
	flag := calInstancesCmd.Flag("calendar")
	if flag == nil {
		t.Error("expected --calendar flag to be set")
	}
	if flag.DefValue != "primary" {
		t.Errorf("expected --calendar default to be 'primary', got %s", flag.DefValue)
	}
}

func TestCalInstancesCmd_HasStartFlag(t *testing.T) {
	flag := calInstancesCmd.Flag("start")
	if flag == nil {
		t.Error("expected --start flag to be set")
	}
}

func TestCalInstancesCmd_HasEndFlag(t *testing.T) {
	flag := calInstancesCmd.Flag("end")
	if flag == nil {
		t.Error("expected --end flag to be set")
	}
}

func TestCalInstancesCmd_HasMaxResultsFlag(t *testing.T) {
	flag := calInstancesCmd.Flag("max-results")
	if flag == nil {
		t.Error("expected --max-results flag to be set")
	}
	if flag.DefValue != "25" {
		t.Errorf("expected --max-results default to be '25', got %s", flag.DefValue)
	}
}

func TestCalInstancesCmd_RegisteredWithCalCmd(t *testing.T) {
	found := false
	for _, sub := range calCmd.Commands() {
		if sub.Name() == "instances" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected 'instances' subcommand to be registered with calCmd")
	}
}

func TestCalInstancesCmd_HasRecurringAlias(t *testing.T) {
	aliases := calInstancesCmd.Aliases
	found := false

	for _, alias := range aliases {
		if alias == "recurring" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected 'recurring' to be an alias for 'instances'")
	}
}

func TestCalInstancesCmd_AliasWorks(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "recurring", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error for alias 'recurring': %v", err)
	}
}
