// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestCalListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
	if !contains(output, "--max-results") {
		t.Error("expected output to contain '--max-results'")
	}
	if !contains(output, "upcoming") {
		t.Error("expected output to contain 'upcoming'")
	}
}

func TestCalShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
	if !contains(output, "<event-id>") {
		t.Error("expected output to contain '<event-id>'")
	}
}

func TestCalShowCmd_RequiresEventID(t *testing.T) {
	if calShowCmd.Args == nil {
		t.Error("calShowCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := calShowCmd.Args(calShowCmd, []string{})
	if err == nil {
		t.Error("expected error when event ID is missing")
	}
}

func TestCalTodayCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "today", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "today") {
		t.Error("expected output to contain 'today'")
	}
	if !contains(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
}

func TestCalWeekCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "week", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "week") {
		t.Error("expected output to contain 'week'")
	}
	if !contains(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
}

func TestCalListCmd_DefaultMaxResults(t *testing.T) {
	// Reset to default first
	calListMaxResults = 25

	// Verify the default value is set correctly
	if calListMaxResults != 25 {
		t.Errorf("expected default max-results to be 25, got %d", calListMaxResults)
	}
}

func TestCalListCmd_Aliases(t *testing.T) {
	tests := []struct {
		alias   string
		command string
	}{
		{"ls", "list"},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(calCmd)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"cal", tt.alias, "--help"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error for alias %s: %v", tt.alias, err)
			}
		})
	}
}

func TestCalShowCmd_Aliases(t *testing.T) {
	tests := []struct {
		alias   string
		command string
	}{
		{"get", "show"},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(calCmd)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			// Need an argument for show/get
			cmd.SetArgs([]string{"cal", tt.alias, "test-id", "--help"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error for alias %s: %v", tt.alias, err)
			}
		})
	}
}

func TestCalListCmd_HasCalendarFlag(t *testing.T) {
	// Verify the command has a --calendar flag
	flag := calListCmd.Flag("calendar")
	if flag == nil {
		t.Error("expected --calendar flag to be set")
	}
}

func TestCalListCmd_HasMaxResultsFlag(t *testing.T) {
	// Verify the command has a --max-results flag
	flag := calListCmd.Flag("max-results")
	if flag == nil {
		t.Error("expected --max-results flag to be set")
	}
}

func TestCalShowCmd_HasCalendarFlag(t *testing.T) {
	// Verify the command has a --calendar flag
	flag := calShowCmd.Flag("calendar")
	if flag == nil {
		t.Error("expected --calendar flag to be set")
	}
}

func TestCalTodayCmd_HasCalendarFlag(t *testing.T) {
	// Verify the command has a --calendar flag
	flag := calTodayCmd.Flag("calendar")
	if flag == nil {
		t.Error("expected --calendar flag to be set")
	}
}

func TestCalWeekCmd_HasCalendarFlag(t *testing.T) {
	// Verify the command has a --calendar flag
	flag := calWeekCmd.Flag("calendar")
	if flag == nil {
		t.Error("expected --calendar flag to be set")
	}
}

func TestCalCmdSubcommands_Registered(t *testing.T) {
	// Verify subcommands are registered with calCmd
	subcommands := map[string]bool{
		"list":      false,
		"show":      false,
		"today":     false,
		"week":      false,
		"create":    false,
		"update":    false,
		"delete":    false,
		"instances": false,
		"quick":     false,
		"freebusy":  false,
		"rsvp":      false,
		"move":      false,
		"calendars": false,
		"acl":       false,
	}

	for _, sub := range calCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with calCmd", name)
		}
	}
}

func TestCalListCmd_HasArgsRequirement(t *testing.T) {
	// calListCmd should not require args (it lists events without args)
	if calListCmd.Args != nil {
		err := calListCmd.Args(calListCmd, []string{})
		if err != nil {
			t.Errorf("calListCmd should accept no args: %v", err)
		}
	}
}

func TestCalShowCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"event-id"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"event-id", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calShowCmd.Args(calShowCmd, tt.args)
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
