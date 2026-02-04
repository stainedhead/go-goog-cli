// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestCalCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "create", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
	if !contains(output, "--title") {
		t.Error("expected output to contain '--title'")
	}
	if !contains(output, "--start") {
		t.Error("expected output to contain '--start'")
	}
	if !contains(output, "--end") {
		t.Error("expected output to contain '--end'")
	}
	if !contains(output, "--location") {
		t.Error("expected output to contain '--location'")
	}
	if !contains(output, "--description") {
		t.Error("expected output to contain '--description'")
	}
	if !contains(output, "--attendees") {
		t.Error("expected output to contain '--attendees'")
	}
	if !contains(output, "--all-day") {
		t.Error("expected output to contain '--all-day'")
	}
	if !contains(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
}

func TestCalUpdateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "update", "--help"})

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
	if !contains(output, "--title") {
		t.Error("expected output to contain '--title'")
	}
	if !contains(output, "--start") {
		t.Error("expected output to contain '--start'")
	}
	if !contains(output, "--end") {
		t.Error("expected output to contain '--end'")
	}
	if !contains(output, "--location") {
		t.Error("expected output to contain '--location'")
	}
	if !contains(output, "--description") {
		t.Error("expected output to contain '--description'")
	}
	if !contains(output, "--attendees") {
		t.Error("expected output to contain '--attendees'")
	}
	if !contains(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
}

func TestCalDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "delete", "--help"})

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
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
	if !contains(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
}

func TestCalCreateCmd_RequiresTitleFlag(t *testing.T) {
	// Reset flags
	calCreateTitle = ""
	calCreateStart = "2024-01-15 14:00"

	mockCmd := &cobra.Command{Use: "test"}

	if calCreateCmd.PreRunE != nil {
		err := calCreateCmd.PreRunE(mockCmd, []string{})
		if err == nil {
			t.Error("expected error when --title flag is missing")
		}
	} else {
		t.Error("calCreateCmd should have PreRunE defined")
	}
}

func TestCalCreateCmd_RequiresStartFlag(t *testing.T) {
	// Reset flags
	calCreateTitle = "Test Meeting"
	calCreateStart = ""

	mockCmd := &cobra.Command{Use: "test"}

	if calCreateCmd.PreRunE != nil {
		err := calCreateCmd.PreRunE(mockCmd, []string{})
		if err == nil {
			t.Error("expected error when --start flag is missing")
		}
	} else {
		t.Error("calCreateCmd should have PreRunE defined")
	}
}

func TestCalUpdateCmd_RequiresIDArg(t *testing.T) {
	if calUpdateCmd.Args == nil {
		t.Error("calUpdateCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := calUpdateCmd.Args(calUpdateCmd, []string{})
	if err == nil {
		t.Error("expected error when event ID is missing")
	}
}

func TestCalDeleteCmd_RequiresIDArg(t *testing.T) {
	if calDeleteCmd.Args == nil {
		t.Error("calDeleteCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := calDeleteCmd.Args(calDeleteCmd, []string{})
	if err == nil {
		t.Error("expected error when event ID is missing")
	}
}

func TestCalDeleteCmd_RequiresConfirmFlag(t *testing.T) {
	// Reset flag
	calDeleteConfirm = false

	mockCmd := &cobra.Command{Use: "test"}
	mockCmd.SetOut(new(bytes.Buffer))
	mockCmd.SetErr(new(bytes.Buffer))

	if calDeleteCmd.PreRunE != nil {
		err := calDeleteCmd.PreRunE(mockCmd, []string{"event123"})
		if err == nil {
			t.Error("expected error when --confirm flag is missing")
		}
	} else {
		t.Error("calDeleteCmd should have PreRunE defined")
	}
}

func TestParseDateTime(t *testing.T) {
	// Get current location for testing
	loc := time.Local

	tests := []struct {
		name        string
		input       string
		expectErr   bool
		checkResult func(t *testing.T, result time.Time)
	}{
		{
			name:      "ISO format with time",
			input:     "2024-01-15 14:00",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				expected := time.Date(2024, 1, 15, 14, 0, 0, 0, loc)
				if !result.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:      "ISO format date only",
			input:     "2024-01-15",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				expected := time.Date(2024, 1, 15, 0, 0, 0, 0, loc)
				if !result.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:      "ISO format with seconds",
			input:     "2024-01-15 14:30:45",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				expected := time.Date(2024, 1, 15, 14, 30, 45, 0, loc)
				if !result.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:      "RFC3339 format",
			input:     "2024-01-15T14:00:00Z",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				expected := time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)
				if !result.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:      "tomorrow keyword",
			input:     "tomorrow",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				tomorrow := time.Now().AddDate(0, 0, 1)
				// Check year, month, day match
				if result.Year() != tomorrow.Year() ||
					result.Month() != tomorrow.Month() ||
					result.Day() != tomorrow.Day() {
					t.Errorf("expected tomorrow's date, got %v", result)
				}
			},
		},
		{
			name:      "tomorrow with time",
			input:     "tomorrow 3pm",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				tomorrow := time.Now().AddDate(0, 0, 1)
				if result.Year() != tomorrow.Year() ||
					result.Month() != tomorrow.Month() ||
					result.Day() != tomorrow.Day() ||
					result.Hour() != 15 {
					t.Errorf("expected tomorrow at 3pm, got %v", result)
				}
			},
		},
		{
			name:      "today keyword",
			input:     "today",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				today := time.Now()
				if result.Year() != today.Year() ||
					result.Month() != today.Month() ||
					result.Day() != today.Day() {
					t.Errorf("expected today's date, got %v", result)
				}
			},
		},
		{
			name:      "today with time",
			input:     "today 2pm",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				today := time.Now()
				if result.Year() != today.Year() ||
					result.Month() != today.Month() ||
					result.Day() != today.Day() ||
					result.Hour() != 14 {
					t.Errorf("expected today at 2pm, got %v", result)
				}
			},
		},
		{
			name:      "invalid format",
			input:     "not a date",
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDateTime(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestParseAttendees(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "single attendee",
			input:    []string{"user@example.com"},
			expected: []string{"user@example.com"},
		},
		{
			name:     "multiple attendees",
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
			name:     "filter empty strings",
			input:    []string{"user@example.com", "", "  "},
			expected: []string{"user@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAttendees(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d attendees, got %d", len(tt.expected), len(result))
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

func TestParseTimeOfDay(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectH   int
		expectM   int
		expectErr bool
	}{
		{
			name:    "3pm",
			input:   "3pm",
			expectH: 15,
			expectM: 0,
		},
		{
			name:    "3:30pm",
			input:   "3:30pm",
			expectH: 15,
			expectM: 30,
		},
		{
			name:    "10am",
			input:   "10am",
			expectH: 10,
			expectM: 0,
		},
		{
			name:    "12pm (noon)",
			input:   "12pm",
			expectH: 12,
			expectM: 0,
		},
		{
			name:    "12am (midnight)",
			input:   "12am",
			expectH: 0,
			expectM: 0,
		},
		{
			name:    "24-hour format",
			input:   "14:00",
			expectH: 14,
			expectM: 0,
		},
		{
			name:    "24-hour format with minutes",
			input:   "14:30",
			expectH: 14,
			expectM: 30,
		},
		{
			name:      "invalid format",
			input:     "invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, m, err := parseTimeOfDay(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if h != tt.expectH || m != tt.expectM {
				t.Errorf("expected %02d:%02d, got %02d:%02d", tt.expectH, tt.expectM, h, m)
			}
		})
	}
}

func TestCalEventCmdSubcommands_Registered(t *testing.T) {
	// Verify subcommands are registered with calCmd
	subcommands := map[string]bool{
		"create": false,
		"update": false,
		"delete": false,
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

func TestCalCreateCmd_HasRequiredFlags(t *testing.T) {
	flags := []string{"title", "start", "end", "location", "description", "attendees", "all-day", "calendar"}

	for _, flagName := range flags {
		flag := calCreateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on create command", flagName)
		}
	}
}

func TestCalUpdateCmd_HasRequiredFlags(t *testing.T) {
	flags := []string{"title", "start", "end", "location", "description", "attendees", "calendar"}

	for _, flagName := range flags {
		flag := calUpdateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on update command", flagName)
		}
	}
}

func TestCalDeleteCmd_HasRequiredFlags(t *testing.T) {
	flags := []string{"confirm", "calendar"}

	for _, flagName := range flags {
		flag := calDeleteCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on delete command", flagName)
		}
	}
}
