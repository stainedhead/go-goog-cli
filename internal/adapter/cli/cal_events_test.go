// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
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
		name      string
		input     []string
		expected  []string
		expectErr bool
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
		{
			name:      "invalid email",
			input:     []string{"notanemail"},
			expectErr: true,
		},
		{
			name:      "one valid one invalid",
			input:     []string{"valid@example.com", "invalid"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAttendees(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %v, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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

func TestParseTimeOfDay_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectH   int
		expectM   int
		expectErr bool
	}{
		{
			name:    "1am",
			input:   "1am",
			expectH: 1,
			expectM: 0,
		},
		{
			name:    "11pm",
			input:   "11pm",
			expectH: 23,
			expectM: 0,
		},
		{
			name:    "9:45am",
			input:   "9:45am",
			expectH: 9,
			expectM: 45,
		},
		{
			name:    "11:59pm",
			input:   "11:59pm",
			expectH: 23,
			expectM: 59,
		},
		{
			name:    "0:00 24h format",
			input:   "0:00",
			expectH: 0,
			expectM: 0,
		},
		{
			name:    "23:59 24h format",
			input:   "23:59",
			expectH: 23,
			expectM: 59,
		},
		{
			name:    "single digit hour 24h",
			input:   "9:30",
			expectH: 9,
			expectM: 30,
		},
		{
			name:      "invalid hour 24h",
			input:     "25:00",
			expectErr: true,
		},
		{
			name:      "invalid minute 24h",
			input:     "10:60",
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:    "whitespace trimming",
			input:   "  3pm  ",
			expectH: 15,
			expectM: 0,
		},
		{
			name:    "uppercase AM",
			input:   "3PM",
			expectH: 15,
			expectM: 0,
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

func TestParseDateTime_AdditionalFormats(t *testing.T) {
	loc := time.Local

	tests := []struct {
		name        string
		input       string
		expectErr   bool
		checkResult func(t *testing.T, result time.Time)
	}{
		{
			name:      "RFC3339 with timezone offset",
			input:     "2024-06-15T10:30:00+02:00",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				if result.Hour() != 10 || result.Minute() != 30 {
					t.Errorf("expected 10:30, got %02d:%02d", result.Hour(), result.Minute())
				}
			},
		},
		{
			name:      "whitespace only",
			input:     "   ",
			expectErr: true,
		},
		{
			name:      "partial date",
			input:     "2024-01",
			expectErr: true,
		},
		{
			name:      "time only",
			input:     "14:00",
			expectErr: true,
		},
		{
			name:      "date with full time and seconds",
			input:     "2024-03-20 09:15:30",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				expected := time.Date(2024, 3, 20, 9, 15, 30, 0, loc)
				if !result.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:      "today with 24h time",
			input:     "today 14:00",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				today := time.Now()
				if result.Year() != today.Year() ||
					result.Month() != today.Month() ||
					result.Day() != today.Day() ||
					result.Hour() != 14 ||
					result.Minute() != 0 {
					t.Errorf("expected today at 14:00, got %v", result)
				}
			},
		},
		{
			name:      "tomorrow with minutes",
			input:     "tomorrow 3:30pm",
			expectErr: false,
			checkResult: func(t *testing.T, result time.Time) {
				tomorrow := time.Now().AddDate(0, 0, 1)
				if result.Year() != tomorrow.Year() ||
					result.Month() != tomorrow.Month() ||
					result.Day() != tomorrow.Day() ||
					result.Hour() != 15 ||
					result.Minute() != 30 {
					t.Errorf("expected tomorrow at 3:30pm, got %v", result)
				}
			},
		},
		{
			name:      "today with invalid time",
			input:     "today badtime",
			expectErr: true,
		},
		{
			name:      "tomorrow with invalid time",
			input:     "tomorrow xyz",
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

func TestParseAttendees_AdditionalCases(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		expected  []string
		expectErr bool
	}{
		{
			name:     "mixed whitespace and empty",
			input:    []string{"  ", "user@example.com", "", "   user2@example.com   "},
			expected: []string{"user@example.com", "user2@example.com"},
		},
		{
			name:     "all empty strings",
			input:    []string{"", "", ""},
			expected: []string{},
		},
		{
			name:     "tabs and newlines",
			input:    []string{"\tuser@example.com\t"},
			expected: []string{"user@example.com"},
		},
		{
			name:     "large list",
			input:    []string{"a@b.com", "c@d.com", "e@f.com", "g@h.com", "i@j.com"},
			expected: []string{"a@b.com", "c@d.com", "e@f.com", "g@h.com", "i@j.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAttendees(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %v, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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

func TestCalCreateCmd_ValidatesRequiredFlags(t *testing.T) {
	// Test that both title and start are validated
	tests := []struct {
		name      string
		title     string
		start     string
		expectErr bool
	}{
		{
			name:      "both missing",
			title:     "",
			start:     "",
			expectErr: true,
		},
		{
			name:      "title missing",
			title:     "",
			start:     "2024-01-15 14:00",
			expectErr: true,
		},
		{
			name:      "start missing",
			title:     "Test Meeting",
			start:     "",
			expectErr: true,
		},
		{
			name:      "both present",
			title:     "Test Meeting",
			start:     "2024-01-15 14:00",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origTitle := calCreateTitle
			origStart := calCreateStart

			// Set test values
			calCreateTitle = tt.title
			calCreateStart = tt.start

			mockCmd := &cobra.Command{Use: "test"}

			err := calCreateCmd.PreRunE(mockCmd, []string{})

			// Restore original values
			calCreateTitle = origTitle
			calCreateStart = origStart

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

func TestCalDeleteCmd_ConfirmValidation(t *testing.T) {
	tests := []struct {
		name      string
		confirm   bool
		expectErr bool
	}{
		{
			name:      "confirm true",
			confirm:   true,
			expectErr: false,
		},
		{
			name:      "confirm false",
			confirm:   false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origConfirm := calDeleteConfirm
			calDeleteConfirm = tt.confirm

			mockCmd := &cobra.Command{Use: "test"}
			mockCmd.SetOut(new(bytes.Buffer))
			mockCmd.SetErr(new(bytes.Buffer))

			err := calDeleteCmd.PreRunE(mockCmd, []string{"event123"})

			calDeleteConfirm = origConfirm

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

func TestCalUpdateCmd_ArgsValidation(t *testing.T) {
	// Test Args validation
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
			args:      []string{"event123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"event123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calUpdateCmd.Args(calUpdateCmd, tt.args)
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

func TestCalDeleteCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"event123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"event123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calDeleteCmd.Args(calDeleteCmd, tt.args)
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

// =============================================================================
// Edge Case and Error Path Tests for Calendar Events
// =============================================================================

func TestParseDateTime_InvalidFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid date format",
			input: "15-01-2024",
		},
		{
			name:  "time without date",
			input: "14:30",
		},
		{
			name:  "random text",
			input: "not a valid date",
		},
		{
			name:  "partial ISO date",
			input: "2024-01",
		},
		{
			name:  "invalid month",
			input: "2024-13-01",
		},
		{
			name:  "invalid day",
			input: "2024-01-32",
		},
		{
			name:  "whitespace only",
			input: "   ",
		},
		{
			name:  "tomorrow with extra words",
			input: "tomorrow at the meeting",
		},
		{
			name:  "today with invalid time format",
			input: "today 25:00",
		},
		{
			name:  "yesterday keyword (not supported)",
			input: "yesterday",
		},
		{
			name:  "next week (not supported)",
			input: "next week",
		},
		{
			name:  "date with slash separators",
			input: "01/15/2024",
		},
		{
			name:  "european date format",
			input: "15.01.2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseDateTime(tt.input)
			if err == nil {
				t.Errorf("expected error for invalid input %q, got nil", tt.input)
			}
		})
	}
}

func TestParseDateTime_ValidEdgeCases(t *testing.T) {
	loc := time.Local

	tests := []struct {
		name        string
		input       string
		checkResult func(t *testing.T, result time.Time)
	}{
		{
			name:  "midnight",
			input: "2024-01-15 00:00",
			checkResult: func(t *testing.T, result time.Time) {
				expected := time.Date(2024, 1, 15, 0, 0, 0, 0, loc)
				if !result.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:  "end of day",
			input: "2024-01-15 23:59:59",
			checkResult: func(t *testing.T, result time.Time) {
				expected := time.Date(2024, 1, 15, 23, 59, 59, 0, loc)
				if !result.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:  "leap year date",
			input: "2024-02-29",
			checkResult: func(t *testing.T, result time.Time) {
				if result.Month() != time.February || result.Day() != 29 {
					t.Errorf("expected Feb 29, got %v", result)
				}
			},
		},
		{
			name:  "new year",
			input: "2024-01-01 00:00",
			checkResult: func(t *testing.T, result time.Time) {
				if result.Month() != time.January || result.Day() != 1 {
					t.Errorf("expected Jan 1, got %v", result)
				}
			},
		},
		{
			name:  "year end",
			input: "2024-12-31 23:59",
			checkResult: func(t *testing.T, result time.Time) {
				if result.Month() != time.December || result.Day() != 31 {
					t.Errorf("expected Dec 31, got %v", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDateTime(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestParseTimeOfDay_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "negative hour",
			input: "-5:00",
		},
		{
			name:  "hour > 23",
			input: "25:00",
		},
		{
			name:  "minute > 59",
			input: "10:60",
		},
		{
			name:  "negative minute",
			input: "10:-30",
		},
		{
			name:  "13pm (invalid)",
			input: "13pm",
		},
		{
			name:  "single digit without am/pm",
			input: "3",
		},
		{
			name:  "no colon in 24h format",
			input: "1400",
		},
		{
			name:  "am/pm without hour",
			input: "am",
		},
		{
			name:  "special characters",
			input: "@#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseTimeOfDay(tt.input)
			if err == nil {
				t.Errorf("expected error for invalid input %q, got nil", tt.input)
			}
		})
	}
}

func TestParseTimeOfDay_ValidBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectH int
		expectM int
	}{
		{
			name:    "1am (boundary)",
			input:   "1am",
			expectH: 1,
			expectM: 0,
		},
		{
			name:    "12:01am (just past midnight)",
			input:   "12:01am",
			expectH: 0,
			expectM: 1,
		},
		{
			name:    "11:59am (just before noon)",
			input:   "11:59am",
			expectH: 11,
			expectM: 59,
		},
		{
			name:    "12:01pm (just past noon)",
			input:   "12:01pm",
			expectH: 12,
			expectM: 1,
		},
		{
			name:    "11:59pm (just before midnight)",
			input:   "11:59pm",
			expectH: 23,
			expectM: 59,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, m, err := parseTimeOfDay(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if h != tt.expectH || m != tt.expectM {
				t.Errorf("expected %02d:%02d, got %02d:%02d", tt.expectH, tt.expectM, h, m)
			}
		})
	}
}

func TestParseAttendees_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		expected  int
		expectErr bool
	}{
		{
			name:     "very long email list",
			input:    []string{"a1@b.com", "a2@b.com", "a3@b.com", "a4@b.com", "a5@b.com", "a6@b.com", "a7@b.com", "a8@b.com", "a9@b.com", "a10@b.com"},
			expected: 10,
		},
		{
			name:     "duplicate emails",
			input:    []string{"user@example.com", "user@example.com"},
			expected: 2, // parseAttendees doesn't deduplicate
		},
		{
			name:     "mixed case emails",
			input:    []string{"User@Example.Com", "test@TEST.COM"},
			expected: 2,
		},
		{
			name:      "email with invalid TLD",
			input:     []string{"user@example.c"},
			expectErr: true, // Short TLD, depends on validation
		},
		{
			name:     "email with numbers only in domain",
			input:    []string{"user@123.com"},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAttendees(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(result) != tt.expected {
					t.Errorf("expected %d attendees, got %d", tt.expected, len(result))
				}
			}
		})
	}
}

func TestParseRelativeDate_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		daysOffset  int
		expectErr   bool
		checkResult func(t *testing.T, result time.Time)
	}{
		{
			name:       "tomorrow with 12am",
			input:      "tomorrow 12am",
			daysOffset: 1,
			checkResult: func(t *testing.T, result time.Time) {
				if result.Hour() != 0 {
					t.Errorf("expected hour 0, got %d", result.Hour())
				}
			},
		},
		{
			name:       "today with 12pm",
			input:      "today 12pm",
			daysOffset: 0,
			checkResult: func(t *testing.T, result time.Time) {
				if result.Hour() != 12 {
					t.Errorf("expected hour 12, got %d", result.Hour())
				}
			},
		},
		{
			name:       "tomorrow with 11:59pm",
			input:      "tomorrow 11:59pm",
			daysOffset: 1,
			checkResult: func(t *testing.T, result time.Time) {
				if result.Hour() != 23 || result.Minute() != 59 {
					t.Errorf("expected 23:59, got %02d:%02d", result.Hour(), result.Minute())
				}
			},
		},
		{
			name:       "today with 1:30am",
			input:      "today 1:30am",
			daysOffset: 0,
			checkResult: func(t *testing.T, result time.Time) {
				if result.Hour() != 1 || result.Minute() != 30 {
					t.Errorf("expected 01:30, got %02d:%02d", result.Hour(), result.Minute())
				}
			},
		},
		{
			name:       "tomorrow only (no time)",
			input:      "tomorrow",
			daysOffset: 1,
			checkResult: func(t *testing.T, result time.Time) {
				if result.Hour() != 0 || result.Minute() != 0 {
					t.Errorf("expected 00:00, got %02d:%02d", result.Hour(), result.Minute())
				}
			},
		},
		{
			name:       "today only (no time)",
			input:      "today",
			daysOffset: 0,
			checkResult: func(t *testing.T, result time.Time) {
				if result.Hour() != 0 || result.Minute() != 0 {
					t.Errorf("expected 00:00, got %02d:%02d", result.Hour(), result.Minute())
				}
			},
		},
		{
			name:       "tomorrow with invalid time",
			input:      "tomorrow 25:00",
			daysOffset: 1,
			expectErr:  true,
		},
		{
			name:       "today with text instead of time",
			input:      "today afternoon",
			daysOffset: 0,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRelativeDate(tt.input, tt.daysOffset)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %q", tt.input)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}
		})
	}
}

func TestCalCreateCmd_AllFlagCombinations(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		start     string
		expectErr bool
	}{
		{
			name:      "both present with whitespace",
			title:     "  Meeting  ",
			start:     "  2024-01-15 14:00  ",
			expectErr: false,
		},
		{
			name:      "title with special characters",
			title:     "Meeting: [URGENT] #123!",
			start:     "2024-01-15 14:00",
			expectErr: false,
		},
		{
			name:      "very long title",
			title:     "This is a very long meeting title that goes on and on and on and on and on and should still be accepted",
			start:     "2024-01-15 14:00",
			expectErr: false,
		},
		{
			name:      "unicode in title",
			title:     "会議 - Meeting",
			start:     "2024-01-15 14:00",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTitle := calCreateTitle
			origStart := calCreateStart

			calCreateTitle = tt.title
			calCreateStart = tt.start

			mockCmd := &cobra.Command{Use: "test"}
			err := calCreateCmd.PreRunE(mockCmd, []string{})

			calCreateTitle = origTitle
			calCreateStart = origStart

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

func TestCalDeleteCmd_WithoutConfirm(t *testing.T) {
	// Test that error message is printed to stderr
	origConfirm := calDeleteConfirm
	calDeleteConfirm = false

	mockCmd := &cobra.Command{Use: "test"}
	errBuf := new(bytes.Buffer)
	mockCmd.SetErr(errBuf)

	err := calDeleteCmd.PreRunE(mockCmd, []string{"event123"})

	calDeleteConfirm = origConfirm

	if err == nil {
		t.Error("expected error without --confirm flag")
	}

	errOutput := errBuf.String()
	if !contains(errOutput, "deletion requires --confirm") {
		t.Errorf("expected error message about --confirm, got: %s", errOutput)
	}
}

// =============================================================================
// Additional Edge Case Tests for Helper Functions
// =============================================================================

func TestParseRelativeDate_NegativeOffset(t *testing.T) {
	// Test with negative offset (yesterday - not typically used, but should work)
	result, err := parseRelativeDate("yesterday", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	yesterday := time.Now().AddDate(0, 0, -1)
	if result.Year() != yesterday.Year() ||
		result.Month() != yesterday.Month() ||
		result.Day() != yesterday.Day() {
		t.Errorf("expected yesterday's date, got %v", result)
	}
}

func TestParseRelativeDate_LargeOffset(t *testing.T) {
	// Test with large offset (far future)
	result, err := parseRelativeDate("future", 365)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	future := time.Now().AddDate(0, 0, 365)
	if result.Year() != future.Year() ||
		result.Month() != future.Month() ||
		result.Day() != future.Day() {
		t.Errorf("expected date 365 days from now, got %v", result)
	}
}

func TestParseRelativeDate_WithExtraSpaces(t *testing.T) {
	// Test with extra spaces in input
	result, err := parseRelativeDate("tomorrow   3pm", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tomorrow := time.Now().AddDate(0, 0, 1)
	if result.Year() != tomorrow.Year() ||
		result.Month() != tomorrow.Month() ||
		result.Day() != tomorrow.Day() ||
		result.Hour() != 15 {
		t.Errorf("expected tomorrow at 3pm, got %v", result)
	}
}

func TestParseRelativeDate_CaseInsensitive(t *testing.T) {
	// Test that TODAY, ToMoRrOw, etc. work
	tests := []struct {
		name   string
		input  string
		offset int
	}{
		{"uppercase TODAY", "TODAY", 0},
		{"mixed case ToMoRrOw", "ToMoRrOw", 1},
		{"uppercase TOMORROW", "TOMORROW", 1},
		{"mixed case tOdAy", "tOdAy", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRelativeDate(tt.input, tt.offset)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := time.Now().AddDate(0, 0, tt.offset)
			if result.Year() != expected.Year() ||
				result.Month() != expected.Month() ||
				result.Day() != expected.Day() {
				t.Errorf("expected date %v, got %v", expected, result)
			}
		})
	}
}

func TestParseTimeOfDay_InvalidFormats(t *testing.T) {
	// Additional invalid format tests
	tests := []string{
		"25pm",
		// Note: "0pm" is actually valid (converts to 12pm)
		"-1am",
		"12:60am",
		"99:99",
		"abc:def",
		"12:30:45pm", // Seconds not supported in am/pm format
		":",
		"12:",
		":30",
		"1200",  // No colon
		"12.30", // Wrong separator
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, err := parseTimeOfDay(input)
			if err == nil {
				t.Errorf("expected error for invalid input %q, got nil", input)
			}
		})
	}
}

func TestParseTimeOfDay_BoundaryTimes(t *testing.T) {
	tests := []struct {
		input   string
		expectH int
		expectM int
	}{
		{"0:00", 0, 0},
		{"0:01", 0, 1},
		{"0:59", 0, 59},
		{"23:00", 23, 0},
		{"23:59", 23, 59},
		{"12:00am", 0, 0},
		{"12:00pm", 12, 0},
		{"12:59am", 0, 59},
		{"12:59pm", 12, 59},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			h, m, err := parseTimeOfDay(tt.input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			if h != tt.expectH || m != tt.expectM {
				t.Errorf("expected %02d:%02d, got %02d:%02d", tt.expectH, tt.expectM, h, m)
			}
		})
	}
}

func TestParseAttendees_SpecialEmailFormats(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		expectErr bool
		expected  int
	}{
		{
			name:     "email with multiple dots",
			input:    []string{"first.middle.last@example.com"},
			expected: 1,
		},
		{
			name:     "email with plus sign",
			input:    []string{"user+tag123@example.com"},
			expected: 1,
		},
		{
			name:     "email with numbers",
			input:    []string{"user123@example456.com"},
			expected: 1,
		},
		{
			name:     "email with hyphens",
			input:    []string{"first-last@my-company.com"},
			expected: 1,
		},
		{
			name:     "email with underscore",
			input:    []string{"user_name@example.com"},
			expected: 1,
		},
		{
			name:      "email without @ sign",
			input:     []string{"userexample.com"},
			expectErr: true,
		},
		{
			name:      "email with space",
			input:     []string{"user @example.com"},
			expectErr: true,
		},
		{
			name:      "email with @ at start",
			input:     []string{"@example.com"},
			expectErr: true,
		},
		{
			name:      "email with @ at end",
			input:     []string{"user@"},
			expectErr: true,
		},
		{
			name:     "very long email",
			input:    []string{"verylongemailaddresswithnumbers12345@verylongdomainnamewithsubdomain.example.com"},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAttendees(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for input %v, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(result) != tt.expected {
					t.Errorf("expected %d attendees, got %d", tt.expected, len(result))
				}
			}
		})
	}
}

func TestParseDateTime_TimezoneHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "UTC timezone",
			input: "2024-01-15T14:00:00Z",
		},
		{
			name:  "positive offset",
			input: "2024-01-15T14:00:00+05:30",
		},
		{
			name:  "negative offset",
			input: "2024-01-15T14:00:00-08:00",
		},
		{
			name:  "UTC with milliseconds",
			input: "2024-01-15T14:00:00.000Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDateTime(tt.input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			// Just verify it parsed successfully
			if result.IsZero() {
				t.Error("expected non-zero time")
			}
		})
	}
}

func TestParseDateTime_LeapYearAndEdgeDates(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "leap year Feb 29",
			input: "2024-02-29",
		},
		{
			name:  "new year",
			input: "2024-01-01",
		},
		{
			name:  "end of year",
			input: "2024-12-31",
		},
		{
			name:  "end of month",
			input: "2024-01-31",
		},
		{
			name:  "Feb 28 non-leap year",
			input: "2023-02-28",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDateTime(tt.input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			if result.IsZero() {
				t.Error("expected non-zero time")
			}
		})
	}
}

func TestParseDateTime_WithSecondsAndSubseconds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "with seconds",
			input:    "2024-01-15 14:30:45",
			expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.Local),
		},
		{
			name:     "without seconds",
			input:    "2024-01-15 14:30",
			expected: time.Date(2024, 1, 15, 14, 30, 0, 0, time.Local),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDateTime(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseRelativeDate_WithMultipleWords(t *testing.T) {
	// Test that only first two words are processed (relative date + time)
	// Additional words should cause error
	tests := []struct {
		name      string
		input     string
		offset    int
		expectErr bool
	}{
		{
			name:      "tomorrow at 3pm (three words)",
			input:     "tomorrow at 3pm",
			offset:    1,
			expectErr: true, // "at" would be treated as time component and fail
		},
		{
			name:      "today in morning",
			input:     "today in morning",
			offset:    0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseRelativeDate(tt.input, tt.offset)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error for multi-word input")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestIsValidEmail_ComprehensiveCoverage(t *testing.T) {
	// Additional email validation edge cases
	tests := []struct {
		email string
		valid bool
	}{
		// Valid complex cases
		{"user+tag@sub.domain.example.com", true},
		{"first.last123@example-domain.co.uk", true},
		{"123@456.com", true},
		{"a@b.co", true},

		// The regex is simple and allows some edge cases - which is OK for basic validation
		{"user@domain..com", true}, // Double dots allowed by simple regex
		{".user@domain.com", true}, // Leading dot allowed
		{"user.@domain.com", true}, // Trailing dot allowed

		// Invalid cases
		{"user@domain", false},      // No TLD
		{"@domain.com", false},      // No local part
		{"user@", false},            // No domain
		{"user", false},             // No @ or domain
		{"user@@domain.com", false}, // Double @
		{"user @domain.com", false}, // Space in local
		{"user@domain .com", false}, // Space in domain
		{"user@domain.c", false},    // TLD too short
		{"", false},                 // Empty
		{"   ", false},              // Whitespace only
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}
}

// =============================================================================
// Execution tests for event CRUD operations (runCalCreate, runCalUpdate, runCalDelete)
// =============================================================================

func TestRunCalCreate_Success(t *testing.T) {
	// Setup - use future date to avoid validation errors
	futureDate := time.Now().AddDate(0, 1, 0)
	mockEvent := &calendar.Event{
		ID:       "created-event-id",
		Title:    "Test Event",
		Start:    futureDate,
		End:      futureDate.Add(time.Hour),
		HTMLLink: "https://calendar.google.com/event?eid=xyz",
	}

	mockRepo := &MockEventRepository{
		CreateResult: mockEvent,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	// Set command flags
	origTitle := calCreateTitle
	origStart := calCreateStart
	origEnd := calCreateEnd
	origAllDay := calCreateAllDay
	origCalendar := calCreateCalendar
	origFormat := formatFlag
	origQuiet := quietFlag

	calCreateTitle = "Test Event"
	calCreateStart = futureDate.Format("2006-01-02 15:04")
	calCreateEnd = futureDate.Add(time.Hour).Format("2006-01-02 15:04")
	calCreateAllDay = false
	calCreateCalendar = "primary"
	formatFlag = "plain"
	quietFlag = false

	defer func() {
		calCreateTitle = origTitle
		calCreateStart = origStart
		calCreateEnd = origEnd
		calCreateAllDay = origAllDay
		calCreateCalendar = origCalendar
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute
	err := runCalCreate(cmd, []string{})

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Test Event") {
		t.Error("expected output to contain event title")
	}
	if !contains(output, "created-event-id") {
		t.Error("expected output to contain event ID")
	}
}

func TestRunCalCreate_AllDayEvent(t *testing.T) {
	// Use future date
	futureDate := time.Now().AddDate(0, 1, 0)

	mockRepo := &MockEventRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calCreateTitle
	origStart := calCreateStart
	origEnd := calCreateEnd
	origAllDay := calCreateAllDay
	origFormat := formatFlag
	origQuiet := quietFlag

	calCreateTitle = "All Day Event"
	calCreateStart = futureDate.Format("2006-01-02")
	calCreateEnd = ""
	calCreateAllDay = true
	formatFlag = "plain"
	quietFlag = true

	defer func() {
		calCreateTitle = origTitle
		calCreateStart = origStart
		calCreateEnd = origEnd
		calCreateAllDay = origAllDay
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalCreate(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunCalCreate_InvalidStartTime(t *testing.T) {
	mockRepo := &MockEventRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calCreateTitle
	origStart := calCreateStart

	calCreateTitle = "Test Event"
	calCreateStart = "invalid-date"

	defer func() {
		calCreateTitle = origTitle
		calCreateStart = origStart
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalCreate(cmd, []string{})

	if err == nil {
		t.Error("expected error for invalid start time")
	}
	if !contains(err.Error(), "invalid start time") {
		t.Errorf("expected 'invalid start time' error, got: %v", err)
	}
}

func TestRunCalCreate_StartAfterEnd(t *testing.T) {
	mockRepo := &MockEventRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calCreateTitle
	origStart := calCreateStart
	origEnd := calCreateEnd

	calCreateTitle = "Test Event"
	calCreateStart = "2024-01-15 15:00"
	calCreateEnd = "2024-01-15 14:00"

	defer func() {
		calCreateTitle = origTitle
		calCreateStart = origStart
		calCreateEnd = origEnd
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalCreate(cmd, []string{})

	if err == nil {
		t.Error("expected error for start time after end time")
	}
	if !contains(err.Error(), "start time must be before end time") {
		t.Errorf("expected time ordering error, got: %v", err)
	}
}

func TestRunCalCreate_RepositoryError(t *testing.T) {
	// Use future date
	futureDate := time.Now().AddDate(0, 1, 0)

	mockRepo := &MockEventRepository{
		CreateErr: fmt.Errorf("API error: calendar not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calCreateTitle
	origStart := calCreateStart
	origEnd := calCreateEnd

	calCreateTitle = "Test Event"
	calCreateStart = futureDate.Format("2006-01-02 15:04")
	calCreateEnd = futureDate.Add(time.Hour).Format("2006-01-02 15:04")

	defer func() {
		calCreateTitle = origTitle
		calCreateStart = origStart
		calCreateEnd = origEnd
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalCreate(cmd, []string{})

	if err == nil {
		t.Error("expected error from repository")
	}
	if !contains(err.Error(), "failed to create event") {
		t.Errorf("expected create error, got: %v", err)
	}
}

func TestRunCalUpdate_Success(t *testing.T) {
	existingEvent := &calendar.Event{
		ID:    "event-123",
		Title: "Original Title",
		Start: time.Now(),
		End:   time.Now().Add(time.Hour),
	}

	updatedEvent := &calendar.Event{
		ID:       "event-123",
		Title:    "Updated Title",
		Start:    time.Now(),
		End:      time.Now().Add(time.Hour),
		HTMLLink: "https://calendar.google.com/event?eid=xyz",
	}

	mockRepo := &MockEventRepository{
		Event:        existingEvent,
		UpdateResult: updatedEvent,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calUpdateTitle
	origCalendar := calUpdateCalendar
	origFormat := formatFlag
	origQuiet := quietFlag

	calUpdateTitle = "Updated Title"
	calUpdateCalendar = "primary"
	formatFlag = "plain"
	quietFlag = false

	defer func() {
		calUpdateTitle = origTitle
		calUpdateCalendar = origCalendar
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalUpdate(cmd, []string{"event-123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Updated Title") {
		t.Error("expected output to contain updated title")
	}
}

func TestRunCalUpdate_EventNotFound(t *testing.T) {
	mockRepo := &MockEventRepository{
		GetErr: fmt.Errorf("event not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origCalendar := calUpdateCalendar
	calUpdateCalendar = "primary"
	defer func() { calUpdateCalendar = origCalendar }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalUpdate(cmd, []string{"nonexistent"})

	if err == nil {
		t.Error("expected error for nonexistent event")
	}
	if !contains(err.Error(), "failed to get event") {
		t.Errorf("expected get error, got: %v", err)
	}
}

func TestRunCalUpdate_InvalidTimeRange(t *testing.T) {
	existingEvent := &calendar.Event{
		ID:    "event-123",
		Title: "Test Event",
		Start: time.Now(),
		End:   time.Now().Add(time.Hour),
	}

	mockRepo := &MockEventRepository{
		Event: existingEvent,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origStart := calUpdateStart
	origEnd := calUpdateEnd
	origCalendar := calUpdateCalendar

	calUpdateStart = "2024-01-15 15:00"
	calUpdateEnd = "2024-01-15 14:00"
	calUpdateCalendar = "primary"

	defer func() {
		calUpdateStart = origStart
		calUpdateEnd = origEnd
		calUpdateCalendar = origCalendar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalUpdate(cmd, []string{"event-123"})

	if err == nil {
		t.Error("expected error for invalid time range")
	}
	if !contains(err.Error(), "start time must be before end time") {
		t.Errorf("expected time ordering error, got: %v", err)
	}
}

func TestRunCalDelete_Success(t *testing.T) {
	mockRepo := &MockEventRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origCalendar := calDeleteCalendar
	origQuiet := quietFlag

	calDeleteCalendar = "primary"
	quietFlag = false

	defer func() {
		calDeleteCalendar = origCalendar
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalDelete(cmd, []string{"event-123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "deleted successfully") {
		t.Error("expected success message in output")
	}
}

func TestRunCalDelete_RepositoryError(t *testing.T) {
	mockRepo := &MockEventRepository{
		DeleteErr: fmt.Errorf("event not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origCalendar := calDeleteCalendar
	calDeleteCalendar = "primary"
	defer func() { calDeleteCalendar = origCalendar }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalDelete(cmd, []string{"nonexistent"})

	if err == nil {
		t.Error("expected error from repository")
	}
	if !contains(err.Error(), "failed to delete event") {
		t.Errorf("expected delete error, got: %v", err)
	}
}

func TestRunCalDelete_QuietMode(t *testing.T) {
	mockRepo := &MockEventRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origCalendar := calDeleteCalendar
	origQuiet := quietFlag

	calDeleteCalendar = "primary"
	quietFlag = true

	defer func() {
		calDeleteCalendar = origCalendar
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalDelete(cmd, []string{"event-123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("expected no output in quiet mode, got: %s", output)
	}
}
