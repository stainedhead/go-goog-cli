// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

func TestCalendarsCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendars", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "calendars") {
		t.Error("expected output to contain 'calendars'")
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
	if !contains(output, "clear") {
		t.Error("expected output to contain 'clear'")
	}
}

func TestCalendarsListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendars", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "calendars") {
		t.Error("expected output to contain 'calendars'")
	}
}

func TestCalendarsShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendars", "show", "--help"})

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

func TestCalendarsCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendars", "create", "--help"})

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
	if !contains(output, "--description") {
		t.Error("expected output to contain '--description'")
	}
	if !contains(output, "--timezone") {
		t.Error("expected output to contain '--timezone'")
	}
}

func TestCalendarsUpdateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendars", "update", "--help"})

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
	if !contains(output, "--description") {
		t.Error("expected output to contain '--description'")
	}
	if !contains(output, "--timezone") {
		t.Error("expected output to contain '--timezone'")
	}
}

func TestCalendarsDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendars", "delete", "--help"})

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
}

func TestCalendarsClearCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendars", "clear", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "clear") {
		t.Error("expected output to contain 'clear'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
}

func TestCalendarsShowCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if calendarsShowCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestCalendarsUpdateCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if calendarsUpdateCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestCalendarsDeleteCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if calendarsDeleteCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestCalendarsClearCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if calendarsClearCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestCalendarsCmd_Aliases(t *testing.T) {
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
			for _, sub := range calendarsCmd.Commands() {
				usePrefix := tt.command
				if len(sub.Use) >= len(usePrefix) && sub.Use[:len(usePrefix)] == usePrefix {
					targetCmd = sub
					break
				}
				if sub.Use == tt.command {
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

func TestCalendarsDeleteCmd_HasConfirmFlag(t *testing.T) {
	// Verify the command has a --confirm flag
	flag := calendarsDeleteCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}

func TestCalendarsClearCmd_HasConfirmFlag(t *testing.T) {
	// Verify the command has a --confirm flag
	flag := calendarsClearCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}

func TestCalendarsCreateCmd_HasTitleFlag(t *testing.T) {
	// Verify the command has a --title flag
	flag := calendarsCreateCmd.Flag("title")
	if flag == nil {
		t.Error("expected --title flag to be set")
	}
}

func TestCalendarsCreateCmd_TitleFlagRequired(t *testing.T) {
	// Verify the --title flag is marked as required via annotations
	flag := calendarsCreateCmd.Flag("title")
	if flag == nil {
		t.Fatal("expected --title flag to exist")
	}

	// Check that PreRunE validation exists
	if calendarsCreateCmd.PreRunE == nil {
		t.Error("expected PreRunE validation for title flag")
	}
}

func TestCalendarsCmd_CalendarAlias(t *testing.T) {
	// Test that 'calendar' is an alias for 'calendars'
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "calendar", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
}

func TestCalendarsCmd_ViaCalendarTopLevel(t *testing.T) {
	// Test that 'goog calendar calendars list' works since 'calendar' is an alias for 'cal'
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"calendar", "calendars", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
}

func TestCalendarsShowCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"primary"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"primary", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calendarsShowCmd.Args(calendarsShowCmd, tt.args)
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

func TestCalendarsUpdateCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"primary"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"primary", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calendarsUpdateCmd.Args(calendarsUpdateCmd, tt.args)
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

func TestCalendarsDeleteCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"primary"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"primary", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calendarsDeleteCmd.Args(calendarsDeleteCmd, tt.args)
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

func TestCalendarsClearCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"primary"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"primary", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calendarsClearCmd.Args(calendarsClearCmd, tt.args)
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

func TestCalendarsCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":   false,
		"show":   false,
		"create": false,
		"update": false,
		"delete": false,
		"clear":  false,
	}

	for _, sub := range calendarsCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with calendarsCmd", name)
		}
	}
}

func TestCalendarsCreateCmd_HasFlags(t *testing.T) {
	flags := []string{"title", "description", "timezone"}

	for _, flagName := range flags {
		flag := calendarsCreateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on create command", flagName)
		}
	}
}

func TestCalendarsUpdateCmd_HasFlags(t *testing.T) {
	flags := []string{"title", "description", "timezone"}

	for _, flagName := range flags {
		flag := calendarsUpdateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on update command", flagName)
		}
	}
}

// =============================================================================
// Execution tests for calendar management operations
// =============================================================================

func TestRunCalendarsList_Success(t *testing.T) {
	mockCalendars := []*calendar.Calendar{
		{ID: "primary", Title: "Personal Calendar", Primary: true},
		{ID: "work@group.calendar.google.com", Title: "Work Calendar"},
	}

	mockRepo := &MockCalendarRepository{
		Calendars: mockCalendars,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsList(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Personal Calendar") {
		t.Error("expected output to contain 'Personal Calendar'")
	}
	if !contains(output, "Work Calendar") {
		t.Error("expected output to contain 'Work Calendar'")
	}
}

func TestRunCalendarsShow_Success(t *testing.T) {
	mockCal := &calendar.Calendar{
		ID:       "primary",
		Title:    "Personal Calendar",
		Primary:  true,
		TimeZone: "America/New_York",
	}

	mockRepo := &MockCalendarRepository{
		Calendar: mockCal,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsShow(cmd, []string{"primary"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Personal Calendar") {
		t.Error("expected output to contain calendar title")
	}
}

func TestRunCalendarsCreate_Success(t *testing.T) {
	mockCal := &calendar.Calendar{
		ID:    "new-cal-id",
		Title: "New Calendar",
	}

	mockRepo := &MockCalendarRepository{
		CreateResult: mockCal,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calendarsTitle
	origFormat := formatFlag

	calendarsTitle = "New Calendar"
	formatFlag = "plain"

	defer func() {
		calendarsTitle = origTitle
		formatFlag = origFormat
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsCreate(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "New Calendar") {
		t.Error("expected output to contain calendar title")
	}
	if !contains(output, "new-cal-id") {
		t.Error("expected output to contain calendar ID")
	}
}

func TestRunCalendarsCreate_Error(t *testing.T) {
	mockRepo := &MockCalendarRepository{
		CreateErr: fmt.Errorf("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calendarsTitle
	calendarsTitle = "New Calendar"
	defer func() { calendarsTitle = origTitle }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsCreate(cmd, []string{})

	if err == nil {
		t.Error("expected error from repository")
	}
	if !contains(err.Error(), "failed to create calendar") {
		t.Errorf("expected create error, got: %v", err)
	}
}

func TestRunCalendarsUpdate_Success(t *testing.T) {
	existingCal := &calendar.Calendar{
		ID:         "cal-123",
		Title:      "Old Title",
		AccessRole: "owner",
	}

	updatedCal := &calendar.Calendar{
		ID:    "cal-123",
		Title: "New Title",
	}

	mockRepo := &MockCalendarRepository{
		Calendar:     existingCal,
		UpdateResult: updatedCal,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calendarsTitle
	origFormat := formatFlag

	calendarsTitle = "New Title"
	formatFlag = "plain"

	defer func() {
		calendarsTitle = origTitle
		formatFlag = origFormat
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsUpdate(cmd, []string{"cal-123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "New Title") {
		t.Error("expected output to contain new title")
	}
}

func TestRunCalendarsUpdate_NotOwner(t *testing.T) {
	existingCal := &calendar.Calendar{
		ID:         "cal-123",
		Title:      "Calendar",
		AccessRole: "reader",
	}

	mockRepo := &MockCalendarRepository{
		Calendar: existingCal,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := calendarsTitle
	calendarsTitle = "New Title"
	defer func() { calendarsTitle = origTitle }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsUpdate(cmd, []string{"cal-123"})

	if err == nil {
		t.Error("expected error for insufficient permissions")
	}
	if !contains(err.Error(), "cannot modify calendar") {
		t.Errorf("expected permission error, got: %v", err)
	}
}

func TestRunCalendarsDelete_Success(t *testing.T) {
	existingCal := &calendar.Calendar{
		ID:         "cal-123",
		Title:      "Calendar to Delete",
		Primary:    false,
		AccessRole: "owner",
	}

	mockRepo := &MockCalendarRepository{
		Calendar: existingCal,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	quietFlag = false
	defer func() { quietFlag = origQuiet }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsDelete(cmd, []string{"cal-123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "deleted successfully") {
		t.Error("expected success message in output")
	}
}

func TestRunCalendarsDelete_PrimaryCalendar(t *testing.T) {
	existingCal := &calendar.Calendar{
		ID:         "primary",
		Title:      "Personal Calendar",
		Primary:    true,
		AccessRole: "owner",
	}

	mockRepo := &MockCalendarRepository{
		Calendar: existingCal,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsDelete(cmd, []string{"primary"})

	if err == nil {
		t.Error("expected error for deleting primary calendar")
	}
	if !contains(err.Error(), "cannot delete primary calendar") {
		t.Errorf("expected primary calendar error, got: %v", err)
	}
}

func TestRunCalendarsClear_Success(t *testing.T) {
	existingCal := &calendar.Calendar{
		ID:    "primary",
		Title: "Personal Calendar",
	}

	mockRepo := &MockCalendarRepository{
		Calendar: existingCal,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	quietFlag = false
	defer func() { quietFlag = origQuiet }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsClear(cmd, []string{"primary"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "cleared") {
		t.Error("expected success message in output")
	}
}

func TestRunCalendarsClear_Error(t *testing.T) {
	existingCal := &calendar.Calendar{
		ID:    "primary",
		Title: "Personal Calendar",
	}

	mockRepo := &MockCalendarRepository{
		Calendar: existingCal,
		ClearErr: fmt.Errorf("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalendarsClear(cmd, []string{"primary"})

	if err == nil {
		t.Error("expected error from repository")
	}
	if !contains(err.Error(), "failed to clear calendar") {
		t.Errorf("expected clear error, got: %v", err)
	}
}
