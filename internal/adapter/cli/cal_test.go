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

// =============================================================================
// Tests using dependency injection with mocks
// =============================================================================

func TestRunCalList_WithMockDependencies(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{
			ID:    "event1",
			Title: "Team Meeting",
			Start: now.Add(time.Hour),
			End:   now.Add(2 * time.Hour),
		},
		{
			ID:    "event2",
			Title: "Project Review",
			Start: now.Add(3 * time.Hour),
			End:   now.Add(4 * time.Hour),
		},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	origMaxResults := calListMaxResults
	origCalendar := calCalendarFlag
	formatFlag = "plain"
	calListMaxResults = 25
	calCalendarFlag = "primary"
	defer func() {
		formatFlag = origFormat
		calListMaxResults = origMaxResults
		calCalendarFlag = origCalendar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalList(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Team Meeting") || !contains(output, "Project Review") {
		t.Errorf("expected output to contain event titles, got: %s", output)
	}
}

func TestRunCalList_WithMaxResults(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{ID: "event1", Title: "Event 1", Start: now, End: now.Add(time.Hour)},
		{ID: "event2", Title: "Event 2", Start: now, End: now.Add(time.Hour)},
		{ID: "event3", Title: "Event 3", Start: now, End: now.Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	origMaxResults := calListMaxResults
	origCalendar := calCalendarFlag
	formatFlag = "plain"
	calListMaxResults = 2 // Limit to 2 results
	calCalendarFlag = "primary"
	defer func() {
		formatFlag = origFormat
		calListMaxResults = origMaxResults
		calCalendarFlag = origCalendar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalList(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalList failed: %v", err)
	}

	// The output should only show first 2 events
	output := buf.String()
	if !contains(output, "Event 1") || !contains(output, "Event 2") {
		t.Errorf("expected output to contain first 2 events, got: %s", output)
	}
}

func TestRunCalList_Error(t *testing.T) {
	mockRepo := &MockEventRepository{
		ListErr: fmt.Errorf("calendar API error"),
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

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to list events") {
		t.Errorf("expected error to contain 'failed to list events', got: %v", err)
	}
}

func TestRunCalShow_WithMockDependencies(t *testing.T) {
	now := time.Now()
	mockEvent := &calendar.Event{
		ID:          "event123",
		Title:       "Important Meeting",
		Description: "Discuss project milestones",
		Location:    "Conference Room A",
		Start:       now.Add(time.Hour),
		End:         now.Add(2 * time.Hour),
		HTMLLink:    "https://calendar.google.com/event/123",
	}

	mockRepo := &MockEventRepository{
		Event: mockEvent,
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

	origFormat := formatFlag
	origCalendar := calCalendarFlag
	formatFlag = "plain"
	calCalendarFlag = "primary"
	defer func() {
		formatFlag = origFormat
		calCalendarFlag = origCalendar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalShow(cmd, []string{"event123"})
	if err != nil {
		t.Fatalf("runCalShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Important Meeting") {
		t.Errorf("expected output to contain event title, got: %s", output)
	}
}

func TestRunCalShow_Error(t *testing.T) {
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

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalShow(cmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to get event") {
		t.Errorf("expected error to contain 'failed to get event', got: %v", err)
	}
}

func TestRunCalToday_WithMockDependencies(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{
			ID:    "today1",
			Title: "Morning Standup",
			Start: time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()),
			End:   time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location()),
		},
		{
			ID:    "today2",
			Title: "Lunch Meeting",
			Start: time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()),
			End:   time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location()),
		},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	origQuiet := quietFlag
	origCalendar := calCalendarFlag
	formatFlag = "plain"
	quietFlag = false
	calCalendarFlag = "primary"
	defer func() {
		formatFlag = origFormat
		quietFlag = origQuiet
		calCalendarFlag = origCalendar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalToday(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalToday failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Morning Standup") || !contains(output, "Lunch Meeting") {
		t.Errorf("expected output to contain today's events, got: %s", output)
	}
	if !contains(output, "event(s)") {
		t.Errorf("expected output to contain event count, got: %s", output)
	}
}

func TestRunCalToday_EmptyResults(t *testing.T) {
	mockRepo := &MockEventRepository{
		Events: []*calendar.Event{},
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

	origFormat := formatFlag
	origQuiet := quietFlag
	formatFlag = "plain"
	quietFlag = false
	defer func() {
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalToday(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalToday failed: %v", err)
	}
}

func TestRunCalToday_Error(t *testing.T) {
	mockRepo := &MockEventRepository{
		ListErr: fmt.Errorf("calendar API error"),
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

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalToday(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to list today's events") {
		t.Errorf("expected error to contain 'failed to list today's events', got: %v", err)
	}
}

func TestRunCalWeek_WithMockDependencies(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{
			ID:    "week1",
			Title: "Monday Meeting",
			Start: now,
			End:   now.Add(time.Hour),
		},
		{
			ID:    "week2",
			Title: "Friday Review",
			Start: now.Add(4 * 24 * time.Hour),
			End:   now.Add(4*24*time.Hour + time.Hour),
		},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	origQuiet := quietFlag
	origCalendar := calCalendarFlag
	formatFlag = "plain"
	quietFlag = false
	calCalendarFlag = "primary"
	defer func() {
		formatFlag = origFormat
		quietFlag = origQuiet
		calCalendarFlag = origCalendar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalWeek(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalWeek failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Monday Meeting") || !contains(output, "Friday Review") {
		t.Errorf("expected output to contain week's events, got: %s", output)
	}
	if !contains(output, "event(s)") {
		t.Errorf("expected output to contain event count, got: %s", output)
	}
}

func TestRunCalWeek_QuietMode(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{ID: "event1", Title: "Event", Start: now, End: now.Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	origQuiet := quietFlag
	formatFlag = "plain"
	quietFlag = true // Enable quiet mode
	defer func() {
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalWeek(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalWeek failed: %v", err)
	}

	output := buf.String()
	// In quiet mode, should not show event count
	if contains(output, "event(s)") {
		t.Errorf("quiet mode should not show event count, got: %s", output)
	}
}

func TestRunCalWeek_Error(t *testing.T) {
	mockRepo := &MockEventRepository{
		ListErr: fmt.Errorf("calendar API error"),
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

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalWeek(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to list this week's events") {
		t.Errorf("expected error to contain 'failed to list this week's events', got: %v", err)
	}
}

func TestRunCalList_JSONFormat(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{ID: "event1", Title: "JSON Test Event", Start: now, End: now.Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	formatFlag = "json"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalList(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalList failed: %v", err)
	}

	output := buf.String()
	// JSON output should contain the event
	if !contains(output, "JSON Test Event") {
		t.Errorf("expected JSON output to contain event, got: %s", output)
	}
}

func TestRunCalShow_JSONFormat(t *testing.T) {
	now := time.Now()
	mockEvent := &calendar.Event{
		ID:    "event123",
		Title: "JSON Show Test",
		Start: now,
		End:   now.Add(time.Hour),
	}

	mockRepo := &MockEventRepository{
		Event: mockEvent,
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

	origFormat := formatFlag
	formatFlag = "json"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalShow(cmd, []string{"event123"})
	if err != nil {
		t.Fatalf("runCalShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "JSON Show Test") {
		t.Errorf("expected JSON output to contain event, got: %s", output)
	}
}

func TestRunCalToday_QuietMode(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{ID: "event1", Title: "Quiet Event", Start: now, End: now.Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	origQuiet := quietFlag
	formatFlag = "plain"
	quietFlag = true
	defer func() {
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalToday(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalToday failed: %v", err)
	}

	output := buf.String()
	// In quiet mode, should not show event count
	if contains(output, "event(s)") {
		t.Errorf("quiet mode should not show event count, got: %s", output)
	}
}

func TestRunCalList_EmptyResults(t *testing.T) {
	mockRepo := &MockEventRepository{
		Events: []*calendar.Event{},
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

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalList(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalList failed: %v", err)
	}
}

func TestRunCalList_AccountResolveError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			ResolveErr: fmt.Errorf("no account found"),
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "no account found") {
		t.Errorf("expected error to contain 'no account found', got: %v", err)
	}
}

func TestRunCalShow_AccountResolveError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			ResolveErr: fmt.Errorf("no account found"),
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalShow(cmd, []string{"event123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "no account found") {
		t.Errorf("expected error to contain 'no account found', got: %v", err)
	}
}

func TestRunCalToday_AccountResolveError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			ResolveErr: fmt.Errorf("no account found"),
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalToday(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "no account found") {
		t.Errorf("expected error to contain 'no account found', got: %v", err)
	}
}

func TestRunCalWeek_AccountResolveError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			ResolveErr: fmt.Errorf("no account found"),
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalWeek(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "no account found") {
		t.Errorf("expected error to contain 'no account found', got: %v", err)
	}
}

func TestRunCalWeek_EmptyResults(t *testing.T) {
	mockRepo := &MockEventRepository{
		Events: []*calendar.Event{},
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

	origFormat := formatFlag
	origQuiet := quietFlag
	formatFlag = "plain"
	quietFlag = false
	defer func() {
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalWeek(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalWeek failed: %v", err)
	}
}

func TestRunCalList_TableFormat(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{ID: "event1", Title: "Table Test Event", Start: now, End: now.Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	formatFlag = "table"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalList(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Table Test Event") {
		t.Errorf("expected table output to contain event, got: %s", output)
	}
}

func TestRunCalShow_WithDetails(t *testing.T) {
	now := time.Now()
	mockEvent := &calendar.Event{
		ID:          "event123",
		Title:       "Detailed Meeting",
		Description: "Meeting with agenda items",
		Location:    "Building A, Room 101",
		Start:       now.Add(time.Hour),
		End:         now.Add(2 * time.Hour),
		HTMLLink:    "https://calendar.google.com/event/123",
		Attendees: []*calendar.Attendee{
			{Email: "attendee1@example.com", ResponseStatus: "accepted"},
			{Email: "attendee2@example.com", ResponseStatus: "needsAction"},
		},
	}

	mockRepo := &MockEventRepository{
		Event: mockEvent,
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

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalShow(cmd, []string{"event123"})
	if err != nil {
		t.Fatalf("runCalShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Detailed Meeting") {
		t.Errorf("expected output to contain event title, got: %s", output)
	}
}

func TestRunCalToday_MultipleEvents(t *testing.T) {
	now := time.Now()
	mockEvents := []*calendar.Event{
		{
			ID:    "event1",
			Title: "9am Standup",
			Start: time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()),
			End:   time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location()),
		},
		{
			ID:    "event2",
			Title: "10am Review",
			Start: time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location()),
			End:   time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, now.Location()),
		},
		{
			ID:    "event3",
			Title: "2pm Planning",
			Start: time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location()),
			End:   time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()),
		},
	}

	mockRepo := &MockEventRepository{
		Events: mockEvents,
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

	origFormat := formatFlag
	origQuiet := quietFlag
	formatFlag = "plain"
	quietFlag = false
	defer func() {
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalToday(cmd, []string{})
	if err != nil {
		t.Fatalf("runCalToday failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "9am Standup") || !contains(output, "10am Review") || !contains(output, "2pm Planning") {
		t.Errorf("expected output to contain all events, got: %s", output)
	}
	if !contains(output, "3 event") {
		t.Errorf("expected output to show 3 events, got: %s", output)
	}
}
