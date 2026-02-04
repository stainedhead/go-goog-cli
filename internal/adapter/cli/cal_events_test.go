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
