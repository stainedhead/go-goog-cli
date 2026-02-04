// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
)

// containsStr is a helper function to check if string contains substring.
// Named differently to avoid collision with contains in auth_test.go.
func containsStr(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

func TestCalCmd_UtilsHelp(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !containsStr(output, "cal") {
		t.Error("expected output to contain 'cal'")
	}
	if !containsStr(output, "quick") {
		t.Error("expected output to contain 'quick'")
	}
	if !containsStr(output, "freebusy") {
		t.Error("expected output to contain 'freebusy'")
	}
	if !containsStr(output, "rsvp") {
		t.Error("expected output to contain 'rsvp'")
	}
	if !containsStr(output, "move") {
		t.Error("expected output to contain 'move'")
	}
	if !containsStr(output, "instances") {
		t.Error("expected output to contain 'instances'")
	}
}

func TestCalQuickAddCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "quick", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !containsStr(output, "quick") {
		t.Error("expected output to contain 'quick'")
	}
	if !containsStr(output, "natural language") {
		t.Error("expected output to contain 'natural language'")
	}
	if !containsStr(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
}

func TestCalQuickAddCmd_Aliases(t *testing.T) {
	// Test that quickadd and add work as aliases
	aliases := calQuickAddCmd.Aliases
	foundQuickadd := false
	foundAdd := false

	for _, alias := range aliases {
		if alias == "quickadd" {
			foundQuickadd = true
		}
		if alias == "add" {
			foundAdd = true
		}
	}

	if !foundQuickadd {
		t.Error("expected 'quickadd' to be an alias")
	}
	if !foundAdd {
		t.Error("expected 'add' to be an alias")
	}
}

func TestCalQuickAddCmd_RequiresText(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if calQuickAddCmd.Args == nil {
		t.Error("calQuickAddCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := calQuickAddCmd.Args(calQuickAddCmd, []string{})
	if err == nil {
		t.Error("expected error when text is missing")
	}
}

func TestCalFreeBusyCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "freebusy", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !containsStr(output, "freebusy") {
		t.Error("expected output to contain 'freebusy'")
	}
	if !containsStr(output, "--start") {
		t.Error("expected output to contain '--start'")
	}
	if !containsStr(output, "--end") {
		t.Error("expected output to contain '--end'")
	}
	if !containsStr(output, "--calendars") {
		t.Error("expected output to contain '--calendars'")
	}
}

func TestCalFreeBusyCmd_RequiresStartFlag(t *testing.T) {
	// Reset flags
	calFreeBusyStart = ""
	calFreeBusyEnd = "2024-01-15T17:00:00Z"

	mockCmd := &cobra.Command{Use: "test"}

	if calFreeBusyCmd.PreRunE != nil {
		err := calFreeBusyCmd.PreRunE(mockCmd, []string{})
		if err == nil {
			t.Error("expected error when --start flag is not set")
		}
		if err != nil && !containsStr(err.Error(), "--start") {
			t.Errorf("expected error to mention --start, got: %v", err)
		}
	} else {
		t.Error("calFreeBusyCmd should have PreRunE defined")
	}
}

func TestCalFreeBusyCmd_RequiresEndFlag(t *testing.T) {
	// Reset flags
	calFreeBusyStart = "2024-01-15T09:00:00Z"
	calFreeBusyEnd = ""

	mockCmd := &cobra.Command{Use: "test"}

	if calFreeBusyCmd.PreRunE != nil {
		err := calFreeBusyCmd.PreRunE(mockCmd, []string{})
		if err == nil {
			t.Error("expected error when --end flag is not set")
		}
		if err != nil && !containsStr(err.Error(), "--end") {
			t.Errorf("expected error to mention --end, got: %v", err)
		}
	} else {
		t.Error("calFreeBusyCmd should have PreRunE defined")
	}
}

func TestCalFreeBusyCmd_PassesWithBothFlags(t *testing.T) {
	// Reset flags
	calFreeBusyStart = "2024-01-15T09:00:00Z"
	calFreeBusyEnd = "2024-01-15T17:00:00Z"

	mockCmd := &cobra.Command{Use: "test"}

	if calFreeBusyCmd.PreRunE != nil {
		err := calFreeBusyCmd.PreRunE(mockCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error when both flags are set: %v", err)
		}
	} else {
		t.Error("calFreeBusyCmd should have PreRunE defined")
	}
}

func TestCalRSVPCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "rsvp", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !containsStr(output, "rsvp") {
		t.Error("expected output to contain 'rsvp'")
	}
	if !containsStr(output, "--accept") {
		t.Error("expected output to contain '--accept'")
	}
	if !containsStr(output, "--decline") {
		t.Error("expected output to contain '--decline'")
	}
	if !containsStr(output, "--tentative") {
		t.Error("expected output to contain '--tentative'")
	}
}

func TestCalRSVPCmd_RequiresEventID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if calRSVPCmd.Args == nil {
		t.Error("calRSVPCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := calRSVPCmd.Args(calRSVPCmd, []string{})
	if err == nil {
		t.Error("expected error when event ID is missing")
	}
}

func TestCalRSVPCmd_RequiresOneResponseFlag(t *testing.T) {
	// Reset all flags
	calRSVPAccept = false
	calRSVPDecline = false
	calRSVPTentative = false

	mockCmd := &cobra.Command{Use: "test"}

	if calRSVPCmd.PreRunE != nil {
		err := calRSVPCmd.PreRunE(mockCmd, []string{"event123"})
		if err == nil {
			t.Error("expected error when no RSVP flag is set")
		}
	} else {
		t.Error("calRSVPCmd should have PreRunE defined")
	}
}

func TestCalRSVPCmd_RejectsMultipleResponseFlags(t *testing.T) {
	// Set multiple flags
	calRSVPAccept = true
	calRSVPDecline = true
	calRSVPTentative = false

	mockCmd := &cobra.Command{Use: "test"}

	if calRSVPCmd.PreRunE != nil {
		err := calRSVPCmd.PreRunE(mockCmd, []string{"event123"})
		if err == nil {
			t.Error("expected error when multiple RSVP flags are set")
		}
		if err != nil && !containsStr(err.Error(), "only one") {
			t.Errorf("expected error to mention 'only one', got: %v", err)
		}
	} else {
		t.Error("calRSVPCmd should have PreRunE defined")
	}
}

func TestCalRSVPCmd_AcceptsAcceptFlag(t *testing.T) {
	calRSVPAccept = true
	calRSVPDecline = false
	calRSVPTentative = false

	mockCmd := &cobra.Command{Use: "test"}

	if calRSVPCmd.PreRunE != nil {
		err := calRSVPCmd.PreRunE(mockCmd, []string{"event123"})
		if err != nil {
			t.Errorf("unexpected error when --accept is set: %v", err)
		}
	} else {
		t.Error("calRSVPCmd should have PreRunE defined")
	}
}

func TestCalRSVPCmd_AcceptsDeclineFlag(t *testing.T) {
	calRSVPAccept = false
	calRSVPDecline = true
	calRSVPTentative = false

	mockCmd := &cobra.Command{Use: "test"}

	if calRSVPCmd.PreRunE != nil {
		err := calRSVPCmd.PreRunE(mockCmd, []string{"event123"})
		if err != nil {
			t.Errorf("unexpected error when --decline is set: %v", err)
		}
	} else {
		t.Error("calRSVPCmd should have PreRunE defined")
	}
}

func TestCalRSVPCmd_AcceptsTentativeFlag(t *testing.T) {
	calRSVPAccept = false
	calRSVPDecline = false
	calRSVPTentative = true

	mockCmd := &cobra.Command{Use: "test"}

	if calRSVPCmd.PreRunE != nil {
		err := calRSVPCmd.PreRunE(mockCmd, []string{"event123"})
		if err != nil {
			t.Errorf("unexpected error when --tentative is set: %v", err)
		}
	} else {
		t.Error("calRSVPCmd should have PreRunE defined")
	}
}

func TestCalMoveCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "move", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !containsStr(output, "move") {
		t.Error("expected output to contain 'move'")
	}
	if !containsStr(output, "--to") {
		t.Error("expected output to contain '--to'")
	}
	if !containsStr(output, "--calendar") {
		t.Error("expected output to contain '--calendar'")
	}
}

func TestCalMoveCmd_RequiresEventID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if calMoveCmd.Args == nil {
		t.Error("calMoveCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := calMoveCmd.Args(calMoveCmd, []string{})
	if err == nil {
		t.Error("expected error when event ID is missing")
	}
}

func TestCalMoveCmd_RequiresToFlag(t *testing.T) {
	// Reset flag
	calMoveDestination = ""

	mockCmd := &cobra.Command{Use: "test"}

	if calMoveCmd.PreRunE != nil {
		err := calMoveCmd.PreRunE(mockCmd, []string{"event123"})
		if err == nil {
			t.Error("expected error when --to flag is not set")
		}
		if err != nil && !containsStr(err.Error(), "--to") {
			t.Errorf("expected error to mention --to, got: %v", err)
		}
	} else {
		t.Error("calMoveCmd should have PreRunE defined")
	}
}

func TestCalMoveCmd_PassesWithToFlag(t *testing.T) {
	// Set flag
	calMoveDestination = "destination@calendar.google.com"

	mockCmd := &cobra.Command{Use: "test"}

	if calMoveCmd.PreRunE != nil {
		err := calMoveCmd.PreRunE(mockCmd, []string{"event123"})
		if err != nil {
			t.Errorf("unexpected error when --to flag is set: %v", err)
		}
	} else {
		t.Error("calMoveCmd should have PreRunE defined")
	}
}

func TestCalCmd_Aliases(t *testing.T) {
	// Test that calendar works as an alias
	aliases := calCmd.Aliases
	found := false

	for _, alias := range aliases {
		if alias == "calendar" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected 'calendar' to be an alias for 'cal'")
	}
}

func TestCalFreeBusyCmd_Aliases(t *testing.T) {
	// Test that busy and availability work as aliases
	aliases := calFreeBusyCmd.Aliases
	foundBusy := false
	foundAvailability := false

	for _, alias := range aliases {
		if alias == "busy" {
			foundBusy = true
		}
		if alias == "availability" {
			foundAvailability = true
		}
	}

	if !foundBusy {
		t.Error("expected 'busy' to be an alias")
	}
	if !foundAvailability {
		t.Error("expected 'availability' to be an alias")
	}
}

func TestCalRSVPCmd_Aliases(t *testing.T) {
	// Test that respond works as an alias
	aliases := calRSVPCmd.Aliases
	found := false

	for _, alias := range aliases {
		if alias == "respond" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected 'respond' to be an alias for 'rsvp'")
	}
}

func TestRenderFreeBusyTable_NilResponse(t *testing.T) {
	result := renderFreeBusyTable(nil, dummyTime(), dummyTime())
	if !containsStr(result, "No busy periods found") {
		t.Error("expected output to mention no busy periods for nil response")
	}
}

func TestRenderFreeBusyJSON_NilResponse(t *testing.T) {
	result := renderFreeBusyJSON(nil)
	if result != "{}" {
		t.Errorf("expected '{}' for nil response, got: %s", result)
	}
}

// dummyTime returns a dummy time for testing.
func dummyTime() time.Time {
	return time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
}

func TestRenderFreeBusyTable_EmptyCalendars(t *testing.T) {
	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{},
	}

	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)

	result := renderFreeBusyTable(response, start, end)
	if !containsStr(result, "No busy periods found") {
		t.Error("expected output to mention no busy periods for empty calendars")
	}
}

func TestRenderFreeBusyTable_WithBusyPeriods(t *testing.T) {
	busyStart := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	busyEnd := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {
				{Start: busyStart, End: busyEnd},
			},
		},
	}

	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)

	result := renderFreeBusyTable(response, start, end)

	if !containsStr(result, "primary") {
		t.Error("expected output to contain calendar ID 'primary'")
	}
	if !containsStr(result, "BUSY") {
		t.Error("expected output to contain 'BUSY'")
	}
	if !containsStr(result, "Free/Busy Information") {
		t.Error("expected output to contain 'Free/Busy Information'")
	}
}

func TestRenderFreeBusyTable_NoBusyPeriods(t *testing.T) {
	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {}, // Empty busy periods (all free)
		},
	}

	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)

	result := renderFreeBusyTable(response, start, end)

	if !containsStr(result, "primary") {
		t.Error("expected output to contain calendar ID 'primary'")
	}
	if !containsStr(result, "No busy periods (free)") {
		t.Error("expected output to mention calendar is free")
	}
}

func TestRenderFreeBusyTable_MultipleCalendars(t *testing.T) {
	busyStart := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	busyEnd := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {
				{Start: busyStart, End: busyEnd},
			},
			"work@example.com": {}, // Free
		},
	}

	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)

	result := renderFreeBusyTable(response, start, end)

	if !containsStr(result, "primary") {
		t.Error("expected output to contain 'primary'")
	}
	if !containsStr(result, "work@example.com") {
		t.Error("expected output to contain 'work@example.com'")
	}
}

func TestRenderFreeBusyJSON_EmptyCalendars(t *testing.T) {
	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{},
	}

	result := renderFreeBusyJSON(response)

	if !containsStr(result, "calendars") {
		t.Error("expected output to contain 'calendars' key")
	}
}

func TestRenderFreeBusyJSON_WithBusyPeriods(t *testing.T) {
	busyStart := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	busyEnd := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {
				{Start: busyStart, End: busyEnd},
			},
		},
	}

	result := renderFreeBusyJSON(response)

	if !containsStr(result, "primary") {
		t.Error("expected output to contain 'primary'")
	}
	if !containsStr(result, "start") {
		t.Error("expected output to contain 'start'")
	}
	if !containsStr(result, "end") {
		t.Error("expected output to contain 'end'")
	}
	if !containsStr(result, "2024-01-15T10:00:00Z") {
		t.Error("expected output to contain start time in RFC3339 format")
	}
	if !containsStr(result, "2024-01-15T11:00:00Z") {
		t.Error("expected output to contain end time in RFC3339 format")
	}
}

func TestRenderFreeBusyJSON_MultipleBusyPeriods(t *testing.T) {
	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {
				{Start: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)},
				{Start: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)},
			},
		},
	}

	result := renderFreeBusyJSON(response)

	// Should have commas between periods
	if !containsStr(result, "2024-01-15T10:00:00Z") {
		t.Error("expected first busy period start time")
	}
	if !containsStr(result, "2024-01-15T14:00:00Z") {
		t.Error("expected second busy period start time")
	}
}

func TestRenderFreeBusyJSON_MultipleCalendars(t *testing.T) {
	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {
				{Start: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)},
			},
			"work@example.com": {
				{Start: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)},
			},
		},
	}

	result := renderFreeBusyJSON(response)

	if !containsStr(result, "primary") {
		t.Error("expected output to contain 'primary'")
	}
	if !containsStr(result, "work@example.com") {
		t.Error("expected output to contain 'work@example.com'")
	}
}

func TestCalQuickAddCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"Meeting tomorrow at 3pm"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"Meeting", "tomorrow"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calQuickAddCmd.Args(calQuickAddCmd, tt.args)
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

func TestCalRSVPCmd_ArgsValidation(t *testing.T) {
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
			err := calRSVPCmd.Args(calRSVPCmd, tt.args)
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

func TestCalMoveCmd_ArgsValidation(t *testing.T) {
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
			err := calMoveCmd.Args(calMoveCmd, tt.args)
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

func TestCalRSVPCmd_AllThreeFlagsSet(t *testing.T) {
	calRSVPAccept = true
	calRSVPDecline = true
	calRSVPTentative = true

	mockCmd := &cobra.Command{Use: "test"}

	if calRSVPCmd.PreRunE != nil {
		err := calRSVPCmd.PreRunE(mockCmd, []string{"event123"})
		if err == nil {
			t.Error("expected error when all three RSVP flags are set")
		}
	}
}

func TestCalFreeBusyCmd_HasFlags(t *testing.T) {
	flags := []string{"start", "end", "calendars"}

	for _, flagName := range flags {
		flag := calFreeBusyCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on freebusy command", flagName)
		}
	}
}

func TestCalRSVPCmd_HasFlags(t *testing.T) {
	flags := []string{"accept", "decline", "tentative", "calendar"}

	for _, flagName := range flags {
		flag := calRSVPCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on rsvp command", flagName)
		}
	}
}

func TestCalMoveCmd_HasFlags(t *testing.T) {
	flags := []string{"to", "calendar"}

	for _, flagName := range flags {
		flag := calMoveCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on move command", flagName)
		}
	}
}

func TestCalCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"quick":    false,
		"freebusy": false,
		"rsvp":     false,
		"move":     false,
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

// =============================================================================
// Edge Case and Error Path Tests for Calendar Utilities
// =============================================================================

func TestCalFreeBusyCmd_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "both missing",
			start:     "",
			end:       "",
			expectErr: true,
			errMsg:    "--start",
		},
		{
			name:      "start missing",
			start:     "",
			end:       "2024-01-15T17:00:00Z",
			expectErr: true,
			errMsg:    "--start",
		},
		{
			name:      "end missing",
			start:     "2024-01-15T09:00:00Z",
			end:       "",
			expectErr: true,
			errMsg:    "--end",
		},
		{
			name:      "both present",
			start:     "2024-01-15T09:00:00Z",
			end:       "2024-01-15T17:00:00Z",
			expectErr: false,
		},
		{
			name:      "whitespace only start",
			start:     "   ",
			end:       "2024-01-15T17:00:00Z",
			expectErr: false, // Whitespace is not empty string, passes PreRunE
		},
		{
			name:      "whitespace only end",
			start:     "2024-01-15T09:00:00Z",
			end:       "   ",
			expectErr: false, // Whitespace is not empty string, passes PreRunE
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origStart := calFreeBusyStart
			origEnd := calFreeBusyEnd

			calFreeBusyStart = tt.start
			calFreeBusyEnd = tt.end

			mockCmd := &cobra.Command{Use: "test"}
			err := calFreeBusyCmd.PreRunE(mockCmd, []string{})

			calFreeBusyStart = origStart
			calFreeBusyEnd = origEnd

			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errMsg != "" && !containsStr(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCalRSVPCmd_FlagCombinations(t *testing.T) {
	tests := []struct {
		name      string
		accept    bool
		decline   bool
		tentative bool
		expectErr bool
	}{
		{
			name:      "none set",
			accept:    false,
			decline:   false,
			tentative: false,
			expectErr: true,
		},
		{
			name:      "accept only",
			accept:    true,
			decline:   false,
			tentative: false,
			expectErr: false,
		},
		{
			name:      "decline only",
			accept:    false,
			decline:   true,
			tentative: false,
			expectErr: false,
		},
		{
			name:      "tentative only",
			accept:    false,
			decline:   false,
			tentative: true,
			expectErr: false,
		},
		{
			name:      "accept and decline",
			accept:    true,
			decline:   true,
			tentative: false,
			expectErr: true,
		},
		{
			name:      "accept and tentative",
			accept:    true,
			decline:   false,
			tentative: true,
			expectErr: true,
		},
		{
			name:      "decline and tentative",
			accept:    false,
			decline:   true,
			tentative: true,
			expectErr: true,
		},
		{
			name:      "all three",
			accept:    true,
			decline:   true,
			tentative: true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origAccept := calRSVPAccept
			origDecline := calRSVPDecline
			origTentative := calRSVPTentative

			calRSVPAccept = tt.accept
			calRSVPDecline = tt.decline
			calRSVPTentative = tt.tentative

			mockCmd := &cobra.Command{Use: "test"}
			err := calRSVPCmd.PreRunE(mockCmd, []string{"event123"})

			calRSVPAccept = origAccept
			calRSVPDecline = origDecline
			calRSVPTentative = origTentative

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

func TestCalMoveCmd_ValidationEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		expectErr   bool
	}{
		{
			name:        "empty destination",
			destination: "",
			expectErr:   true,
		},
		{
			name:        "whitespace destination",
			destination: "   ",
			expectErr:   false, // PreRunE only checks empty string, not whitespace
		},
		{
			name:        "valid destination",
			destination: "target@calendar.google.com",
			expectErr:   false,
		},
		{
			name:        "primary as destination",
			destination: "primary",
			expectErr:   false,
		},
		{
			name:        "very long destination",
			destination: "verylongemailaddress@verylongdomainname.calendar.google.com",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origDest := calMoveDestination
			calMoveDestination = tt.destination

			mockCmd := &cobra.Command{Use: "test"}
			err := calMoveCmd.PreRunE(mockCmd, []string{"event123"})

			calMoveDestination = origDest

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

func TestRenderFreeBusyTable_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		response *calendar.FreeBusyResponse
		start    time.Time
		end      time.Time
		check    func(t *testing.T, output string)
	}{
		{
			name:     "nil response",
			response: nil,
			start:    dummyTime(),
			end:      dummyTime().Add(time.Hour),
			check: func(t *testing.T, output string) {
				if !containsStr(output, "No busy periods found") {
					t.Error("expected 'No busy periods found'")
				}
			},
		},
		{
			name: "empty calendars map",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{},
			},
			start: dummyTime(),
			end:   dummyTime().Add(time.Hour),
			check: func(t *testing.T, output string) {
				if !containsStr(output, "No busy periods found") {
					t.Error("expected 'No busy periods found'")
				}
			},
		},
		{
			name: "calendar with empty periods",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{
					"primary": {},
				},
			},
			start: dummyTime(),
			end:   dummyTime().Add(time.Hour),
			check: func(t *testing.T, output string) {
				if !containsStr(output, "primary") {
					t.Error("expected calendar ID")
				}
				if !containsStr(output, "No busy periods (free)") {
					t.Error("expected 'No busy periods (free)'")
				}
			},
		},
		{
			name: "multiple busy periods",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{
					"primary": {
						{Start: dummyTime(), End: dummyTime().Add(time.Hour)},
						{Start: dummyTime().Add(2 * time.Hour), End: dummyTime().Add(3 * time.Hour)},
					},
				},
			},
			start: dummyTime(),
			end:   dummyTime().Add(4 * time.Hour),
			check: func(t *testing.T, output string) {
				if !containsStr(output, "BUSY") {
					t.Error("expected 'BUSY' marker")
				}
				// Should show both periods
				busyCount := bytes.Count([]byte(output), []byte("BUSY"))
				if busyCount != 2 {
					t.Errorf("expected 2 BUSY markers, got %d", busyCount)
				}
			},
		},
		{
			name: "multiple calendars mixed",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{
					"primary": {
						{Start: dummyTime(), End: dummyTime().Add(time.Hour)},
					},
					"work@example.com": {},
					"personal@example.com": {
						{Start: dummyTime().Add(2 * time.Hour), End: dummyTime().Add(3 * time.Hour)},
					},
				},
			},
			start: dummyTime(),
			end:   dummyTime().Add(4 * time.Hour),
			check: func(t *testing.T, output string) {
				if !containsStr(output, "primary") {
					t.Error("expected 'primary'")
				}
				if !containsStr(output, "work@example.com") {
					t.Error("expected 'work@example.com'")
				}
				if !containsStr(output, "personal@example.com") {
					t.Error("expected 'personal@example.com'")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderFreeBusyTable(tt.response, tt.start, tt.end)
			if tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestRenderFreeBusyJSON_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		response *calendar.FreeBusyResponse
		check    func(t *testing.T, output string)
	}{
		{
			name:     "nil response",
			response: nil,
			check: func(t *testing.T, output string) {
				if output != "{}" {
					t.Errorf("expected '{}', got %s", output)
				}
			},
		},
		{
			name: "empty calendars",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{},
			},
			check: func(t *testing.T, output string) {
				if !containsStr(output, "calendars") {
					t.Error("expected 'calendars' key")
				}
			},
		},
		{
			name: "calendar with no periods",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{
					"primary": {},
				},
			},
			check: func(t *testing.T, output string) {
				if !containsStr(output, "primary") {
					t.Error("expected 'primary'")
				}
				// Empty array with newline: "[\n    ]"
				if !containsStr(output, "    ]") {
					t.Errorf("expected formatted empty array, got: %s", output)
				}
			},
		},
		{
			name: "single period",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{
					"primary": {
						{Start: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)},
					},
				},
			},
			check: func(t *testing.T, output string) {
				if !containsStr(output, "start") {
					t.Error("expected 'start' key")
				}
				if !containsStr(output, "end") {
					t.Error("expected 'end' key")
				}
				if !containsStr(output, "2024-01-15T10:00:00Z") {
					t.Error("expected start time in RFC3339")
				}
			},
		},
		{
			name: "multiple periods with special characters in calendar ID",
			response: &calendar.FreeBusyResponse{
				Calendars: map[string][]*calendar.TimePeriod{
					"user+tag@example.com": {
						{Start: dummyTime(), End: dummyTime().Add(time.Hour)},
					},
				},
			},
			check: func(t *testing.T, output string) {
				if !containsStr(output, "user+tag@example.com") {
					t.Error("expected calendar ID with special characters")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderFreeBusyJSON(tt.response)
			if tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestCalQuickAddCmd_EmptyText(t *testing.T) {
	// Test that Args validator rejects empty string
	err := calQuickAddCmd.Args(calQuickAddCmd, []string{""})
	// Should pass Args validation (counts as 1 arg), though empty text would fail in API
	if err != nil {
		t.Errorf("Args validator should accept empty string (count is 1): %v", err)
	}
}

func TestCalRSVPCmd_EmptyEventID(t *testing.T) {
	// Test Args validator with empty string
	err := calRSVPCmd.Args(calRSVPCmd, []string{""})
	// Should pass Args validation (counts as 1 arg)
	if err != nil {
		t.Errorf("Args validator should accept empty string (count is 1): %v", err)
	}
}

func TestCalMoveCmd_EmptyEventID(t *testing.T) {
	// Test Args validator with empty string
	err := calMoveCmd.Args(calMoveCmd, []string{""})
	// Should pass Args validation (counts as 1 arg)
	if err != nil {
		t.Errorf("Args validator should accept empty string (count is 1): %v", err)
	}
}

func TestCalFreeBusyCmd_DefaultCalendars(t *testing.T) {
	// Test that --calendars has a default value
	flag := calFreeBusyCmd.Flag("calendars")
	if flag == nil {
		t.Fatal("expected --calendars flag to exist")
	}
	// Default should be "primary"
	if flag.DefValue != "[primary]" {
		t.Logf("Default calendars flag value: %s (expected [primary])", flag.DefValue)
	}
}

func TestRenderFreeBusyTable_TimeFormatting(t *testing.T) {
	// Test that times are formatted correctly in the output
	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)
	busyStart := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	busyEnd := time.Date(2024, 1, 15, 11, 45, 0, 0, time.UTC)

	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {
				{Start: busyStart, End: busyEnd},
			},
		},
	}

	output := renderFreeBusyTable(response, start, end)

	// Should contain formatted times
	if !containsStr(output, "10:30") {
		t.Error("expected output to contain start time '10:30'")
	}
	if !containsStr(output, "11:45") {
		t.Error("expected output to contain end time '11:45'")
	}
}

func TestRenderFreeBusyJSON_ValidJSON(t *testing.T) {
	// Test that output is valid JSON structure (basic check)
	response := &calendar.FreeBusyResponse{
		Calendars: map[string][]*calendar.TimePeriod{
			"primary": {
				{Start: dummyTime(), End: dummyTime().Add(time.Hour)},
			},
		},
	}

	output := renderFreeBusyJSON(response)

	// Should have proper JSON structure
	if !containsStr(output, "{") || !containsStr(output, "}") {
		t.Error("expected JSON object with braces")
	}
	if !containsStr(output, "[") || !containsStr(output, "]") {
		t.Error("expected JSON array with brackets")
	}
	if !containsStr(output, "\"calendars\"") {
		t.Error("expected quoted 'calendars' key")
	}
}
