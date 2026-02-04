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

// =============================================================================
// Execution tests for instances operations
// =============================================================================

func TestRunCalInstances_Success(t *testing.T) {
	now := time.Now()
	mockInstances := []*calendar.Event{
		{ID: "instance-1", Title: "Recurring Event", Start: now, End: now.Add(time.Hour)},
		{ID: "instance-2", Title: "Recurring Event", Start: now.AddDate(0, 0, 7), End: now.AddDate(0, 0, 7).Add(time.Hour)},
		{ID: "instance-3", Title: "Recurring Event", Start: now.AddDate(0, 0, 14), End: now.AddDate(0, 0, 14).Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockInstances,
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

	origCalendar := calInstancesCalendar
	origFormat := formatFlag
	origQuiet := quietFlag

	calInstancesCalendar = "primary"
	formatFlag = "plain"
	quietFlag = false

	defer func() {
		calInstancesCalendar = origCalendar
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalInstances(cmd, []string{"recurring-event-id"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "3 instance(s)") {
		t.Error("expected output to contain instance count")
	}
}

func TestRunCalInstances_WithTimeRange(t *testing.T) {
	now := time.Now()
	mockInstances := []*calendar.Event{
		{ID: "instance-1", Title: "Recurring Event", Start: now, End: now.Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockInstances,
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

	origCalendar := calInstancesCalendar
	origStart := calInstancesStart
	origEnd := calInstancesEnd
	origFormat := formatFlag
	origQuiet := quietFlag

	calInstancesCalendar = "primary"
	calInstancesStart = now.Format(time.RFC3339)
	calInstancesEnd = now.AddDate(0, 1, 0).Format(time.RFC3339)
	formatFlag = "plain"
	quietFlag = true

	defer func() {
		calInstancesCalendar = origCalendar
		calInstancesStart = origStart
		calInstancesEnd = origEnd
		formatFlag = origFormat
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalInstances(cmd, []string{"recurring-event-id"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunCalInstances_InvalidStartTime(t *testing.T) {
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

	origStart := calInstancesStart
	calInstancesStart = "invalid-date"
	defer func() { calInstancesStart = origStart }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalInstances(cmd, []string{"recurring-event-id"})

	if err == nil {
		t.Error("expected error for invalid start time")
	}
	if !contains(err.Error(), "invalid start time format") {
		t.Errorf("expected invalid time error, got: %v", err)
	}
}

func TestRunCalInstances_InvalidTimeRange(t *testing.T) {
	now := time.Now()
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

	origStart := calInstancesStart
	origEnd := calInstancesEnd

	calInstancesStart = now.AddDate(0, 1, 0).Format(time.RFC3339)
	calInstancesEnd = now.Format(time.RFC3339)

	defer func() {
		calInstancesStart = origStart
		calInstancesEnd = origEnd
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalInstances(cmd, []string{"recurring-event-id"})

	if err == nil {
		t.Error("expected error for invalid time range")
	}
	if !contains(err.Error(), "start time must be before end time") {
		t.Errorf("expected time ordering error, got: %v", err)
	}
}

func TestRunCalInstances_RepositoryError(t *testing.T) {
	mockRepo := &MockEventRepository{
		InstancesErr: fmt.Errorf("event not found"),
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

	origCalendar := calInstancesCalendar
	calInstancesCalendar = "primary"
	defer func() { calInstancesCalendar = origCalendar }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalInstances(cmd, []string{"nonexistent"})

	if err == nil {
		t.Error("expected error from repository")
	}
	if !contains(err.Error(), "failed to get event instances") {
		t.Errorf("expected instances error, got: %v", err)
	}
}

func TestRunCalInstances_MaxResults(t *testing.T) {
	now := time.Now()
	mockInstances := []*calendar.Event{
		{ID: "instance-1", Title: "Event", Start: now, End: now.Add(time.Hour)},
		{ID: "instance-2", Title: "Event", Start: now.AddDate(0, 0, 1), End: now.AddDate(0, 0, 1).Add(time.Hour)},
		{ID: "instance-3", Title: "Event", Start: now.AddDate(0, 0, 2), End: now.AddDate(0, 0, 2).Add(time.Hour)},
		{ID: "instance-4", Title: "Event", Start: now.AddDate(0, 0, 3), End: now.AddDate(0, 0, 3).Add(time.Hour)},
		{ID: "instance-5", Title: "Event", Start: now.AddDate(0, 0, 4), End: now.AddDate(0, 0, 4).Add(time.Hour)},
	}

	mockRepo := &MockEventRepository{
		Events: mockInstances,
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

	origCalendar := calInstancesCalendar
	origMaxResults := calInstancesMaxResults
	origQuiet := quietFlag

	calInstancesCalendar = "primary"
	calInstancesMaxResults = 3
	quietFlag = false

	defer func() {
		calInstancesCalendar = origCalendar
		calInstancesMaxResults = origMaxResults
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runCalInstances(cmd, []string{"recurring-event-id"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "3 instance(s)") {
		t.Error("expected output to show max 3 instances")
	}
}
