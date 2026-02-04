// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

func TestThreadCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "thread") {
		t.Error("expected output to contain 'thread'")
	}
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
}

func TestThreadListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "--max-results") {
		t.Error("expected output to contain '--max-results'")
	}
	if !contains(output, "--labels") {
		t.Error("expected output to contain '--labels'")
	}
}

func TestThreadShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "show", "--help"})

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

func TestThreadTrashCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "trash", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestThreadModifyCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "modify", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
	if !contains(output, "--add-labels") {
		t.Error("expected output to contain '--add-labels'")
	}
	if !contains(output, "--remove-labels") {
		t.Error("expected output to contain '--remove-labels'")
	}
}

func TestThreadShowCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if threadShowCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadTrashCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if threadTrashCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadModifyCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if threadModifyCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadCmd_Aliases(t *testing.T) {
	tests := []struct {
		name    string
		command string
		alias   string
	}{
		{"list alias ls", "list", "ls"},
		{"show alias get", "show", "get"},
		{"show alias read", "show", "read"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find the subcommand
			var targetCmd *cobra.Command
			for _, sub := range threadCmd.Commands() {
				if sub.Use[:len(tt.command)] == tt.command || sub.Use == tt.command {
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

func TestThreadListCmd_HasMaxResultsFlag(t *testing.T) {
	flag := threadListCmd.Flag("max-results")
	if flag == nil {
		t.Error("expected --max-results flag to be set")
	}
}

func TestThreadListCmd_HasLabelsFlag(t *testing.T) {
	flag := threadListCmd.Flag("labels")
	if flag == nil {
		t.Error("expected --labels flag to be set")
	}
}

func TestThreadModifyCmd_HasAddLabelsFlag(t *testing.T) {
	flag := threadModifyCmd.Flag("add-labels")
	if flag == nil {
		t.Error("expected --add-labels flag to be set")
	}
}

func TestThreadModifyCmd_HasRemoveLabelsFlag(t *testing.T) {
	flag := threadModifyCmd.Flag("remove-labels")
	if flag == nil {
		t.Error("expected --remove-labels flag to be set")
	}
}

func TestThreadShowCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"thread123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"thread123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := threadShowCmd.Args(threadShowCmd, tt.args)
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

func TestThreadTrashCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"thread123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"thread123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := threadTrashCmd.Args(threadTrashCmd, tt.args)
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

func TestThreadModifyCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"thread123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"thread123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := threadModifyCmd.Args(threadModifyCmd, tt.args)
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

func TestThreadCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":   false,
		"show":   false,
		"trash":  false,
		"modify": false,
	}

	for _, sub := range threadCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with threadCmd", name)
		}
	}
}

func TestThreadListCmd_HasAllFlags(t *testing.T) {
	flags := []string{"max-results", "labels"}

	for _, flagName := range flags {
		flag := threadListCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on list command", flagName)
		}
	}
}

func TestThreadModifyCmd_HasAllFlags(t *testing.T) {
	flags := []string{"add-labels", "remove-labels"}

	for _, flagName := range flags {
		flag := threadModifyCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on modify command", flagName)
		}
	}
}

func TestThreadShowCmd_HasIdArg(t *testing.T) {
	// Check that Use string contains <id>
	if !contains(threadShowCmd.Use, "<id>") {
		t.Error("expected Use to contain '<id>'")
	}
}

func TestThreadTrashCmd_HasIdArg(t *testing.T) {
	// Check that Use string contains <id>
	if !contains(threadTrashCmd.Use, "<id>") {
		t.Error("expected Use to contain '<id>'")
	}
}

func TestThreadModifyCmd_HasIdArg(t *testing.T) {
	// Check that Use string contains <id>
	if !contains(threadModifyCmd.Use, "<id>") {
		t.Error("expected Use to contain '<id>'")
	}
}

func TestThreadListCmd_DefaultMaxResults(t *testing.T) {
	flag := threadListCmd.Flag("max-results")
	if flag == nil {
		t.Fatal("expected --max-results flag to be set")
	}

	// Check default value
	if flag.DefValue != "20" {
		t.Errorf("expected default max-results to be '20', got '%s'", flag.DefValue)
	}
}

func TestThreadModifyCmd_Validation(t *testing.T) {
	tests := []struct {
		name         string
		addLabels    []string
		removeLabels []string
		expectErr    bool
	}{
		{
			name:         "no labels",
			addLabels:    nil,
			removeLabels: nil,
			expectErr:    true,
		},
		{
			name:         "only add labels",
			addLabels:    []string{"IMPORTANT"},
			removeLabels: nil,
			expectErr:    false,
		},
		{
			name:         "only remove labels",
			addLabels:    nil,
			removeLabels: []string{"INBOX"},
			expectErr:    false,
		},
		{
			name:         "both add and remove",
			addLabels:    []string{"IMPORTANT"},
			removeLabels: []string{"INBOX"},
			expectErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origAdd := threadAddLabels
			origRemove := threadRemoveLabels
			threadAddLabels = tt.addLabels
			threadRemoveLabels = tt.removeLabels

			mockCmd := &cobra.Command{Use: "test"}

			err := threadModifyCmd.PreRunE(mockCmd, []string{"thread123"})

			threadAddLabels = origAdd
			threadRemoveLabels = origRemove

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

func TestRunThreadList_WithMockDependencies(t *testing.T) {
	mockThreads := []*mail.Thread{
		{ID: "thread1", Snippet: "Test thread snippet 1"},
		{ID: "thread2", Snippet: "Test thread snippet 2"},
	}

	mockRepo := &MockThreadRepository{
		Threads: mockThreads,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origMaxResults := threadMaxResults
	formatFlag = "plain"
	threadMaxResults = 20
	defer func() {
		formatFlag = origFormat
		threadMaxResults = origMaxResults
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadList(cmd, []string{})
	if err != nil {
		t.Fatalf("runThreadList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "thread1") || !contains(output, "thread2") {
		t.Errorf("expected output to contain thread IDs, got: %s", output)
	}
}

func TestRunThreadList_Error(t *testing.T) {
	mockRepo := &MockThreadRepository{
		ListErr: fmt.Errorf("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to list threads") {
		t.Errorf("expected error to contain 'failed to list threads', got: %v", err)
	}
}

func TestRunThreadShow_WithMockDependencies(t *testing.T) {
	mockThread := &mail.Thread{
		ID:      "thread123",
		Snippet: "Test thread snippet",
		Messages: []*mail.Message{
			{
				ID:      "msg1",
				Subject: "Test Subject",
				From:    "sender@example.com",
				Body:    "Message body",
				Date:    time.Now(),
			},
		},
	}

	mockRepo := &MockThreadRepository{
		Thread: mockThread,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
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

	err := runThreadShow(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Test Subject") {
		t.Errorf("expected output to contain subject, got: %s", output)
	}
}

func TestRunThreadShow_Error(t *testing.T) {
	mockRepo := &MockThreadRepository{
		GetErr: fmt.Errorf("thread not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadShow(cmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to get thread") {
		t.Errorf("expected error to contain 'failed to get thread', got: %v", err)
	}
}

func TestRunThreadTrash_WithMockDependencies(t *testing.T) {
	mockRepo := &MockThreadRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
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

	err := runThreadTrash(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadTrash failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "thread123") || !contains(output, "trash") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunThreadTrash_Error(t *testing.T) {
	mockRepo := &MockThreadRepository{
		TrashErr: fmt.Errorf("trash operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadTrash(cmd, []string{"thread123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to trash thread") {
		t.Errorf("expected error to contain 'failed to trash thread', got: %v", err)
	}
}

func TestRunThreadTrash_QuietMode(t *testing.T) {
	mockRepo := &MockThreadRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	quietFlag = true
	defer func() { quietFlag = origQuiet }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadTrash(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadTrash failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunThreadModify_WithMockDependencies(t *testing.T) {
	mockThread := &mail.Thread{
		ID:      "thread123",
		Snippet: "Test thread",
	}

	mockRepo := &MockThreadRepository{
		Thread: mockThread,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origAdd := threadAddLabels
	origRemove := threadRemoveLabels
	formatFlag = "plain"
	threadAddLabels = []string{"IMPORTANT"}
	threadRemoveLabels = []string{"INBOX"}
	defer func() {
		formatFlag = origFormat
		threadAddLabels = origAdd
		threadRemoveLabels = origRemove
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadModify(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadModify failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "thread123") || !contains(output, "modified") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunThreadModify_Error(t *testing.T) {
	mockRepo := &MockThreadRepository{
		ModifyErr: fmt.Errorf("modify operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origAdd := threadAddLabels
	threadAddLabels = []string{"IMPORTANT"}
	defer func() { threadAddLabels = origAdd }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadModify(cmd, []string{"thread123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to modify thread") {
		t.Errorf("expected error to contain 'failed to modify thread', got: %v", err)
	}
}

func TestRunThreadModify_JSONFormat(t *testing.T) {
	mockThread := &mail.Thread{
		ID:      "thread123",
		Snippet: "Test thread",
	}

	mockRepo := &MockThreadRepository{
		Thread: mockThread,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origAdd := threadAddLabels
	formatFlag = "json"
	threadAddLabels = []string{"IMPORTANT"}
	defer func() {
		formatFlag = origFormat
		threadAddLabels = origAdd
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadModify(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadModify failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "thread123") {
		t.Errorf("expected JSON output to contain thread ID, got: %s", output)
	}
}

func TestRunThreadList_WithLabelsFilter(t *testing.T) {
	mockThreads := []*mail.Thread{
		{ID: "thread1", Snippet: "Filtered thread"},
	}

	mockRepo := &MockThreadRepository{
		Threads: mockThreads,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origLabels := threadLabels
	formatFlag = "plain"
	threadLabels = []string{"INBOX", "UNREAD"}
	defer func() {
		formatFlag = origFormat
		threadLabels = origLabels
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadList(cmd, []string{})
	if err != nil {
		t.Fatalf("runThreadList failed: %v", err)
	}
}

func TestRunThreadShow_JSONFormat(t *testing.T) {
	mockThread := &mail.Thread{
		ID:      "thread123",
		Snippet: "JSON thread",
	}

	mockRepo := &MockThreadRepository{
		Thread: mockThread,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
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

	err := runThreadShow(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "thread123") {
		t.Errorf("expected JSON output to contain thread ID, got: %s", output)
	}
}

// =============================================================================
// Tests for thread untrash command
// =============================================================================

func TestThreadUntrashCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "untrash", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "untrash") {
		t.Error("expected output to contain 'untrash'")
	}
	if !contains(output, "Restore") {
		t.Error("expected output to contain 'Restore'")
	}
}

func TestThreadUntrashCmd_HasArgsRequirement(t *testing.T) {
	if threadUntrashCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadUntrashCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"thread123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"thread123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := threadUntrashCmd.Args(threadUntrashCmd, tt.args)
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

func TestRunThreadUntrash_WithMockDependencies(t *testing.T) {
	mockRepo := &MockThreadRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
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

	err := runThreadUntrash(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadUntrash failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "thread123") || !contains(output, "restored") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunThreadUntrash_Error(t *testing.T) {
	mockRepo := &MockThreadRepository{
		UntrashErr: fmt.Errorf("untrash operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadUntrash(cmd, []string{"thread123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to restore thread from trash") {
		t.Errorf("expected error to contain 'failed to restore thread from trash', got: %v", err)
	}
}

func TestRunThreadUntrash_QuietMode(t *testing.T) {
	mockRepo := &MockThreadRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	quietFlag = true
	defer func() { quietFlag = origQuiet }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadUntrash(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadUntrash failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

// =============================================================================
// Tests for thread delete command
// =============================================================================

func TestThreadDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(threadCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"thread", "delete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
	if !contains(output, "permanent") {
		t.Error("expected output to contain 'permanent'")
	}
}

func TestThreadDeleteCmd_HasArgsRequirement(t *testing.T) {
	if threadDeleteCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestThreadDeleteCmd_HasConfirmFlag(t *testing.T) {
	flag := threadDeleteCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be defined on delete command")
	}
}

func TestThreadDeleteCmd_RequiresConfirmFlag(t *testing.T) {
	origConfirm := threadDeleteConfirm
	threadDeleteConfirm = false
	defer func() { threadDeleteConfirm = origConfirm }()

	mockCmd := &cobra.Command{Use: "test"}
	mockCmd.SetOut(new(bytes.Buffer))
	mockCmd.SetErr(new(bytes.Buffer))

	if threadDeleteCmd.PreRunE != nil {
		err := threadDeleteCmd.PreRunE(mockCmd, []string{"thread123"})
		if err == nil {
			t.Error("expected error when --confirm flag is not set")
		}
	} else {
		t.Error("threadDeleteCmd should have PreRunE defined")
	}
}

func TestThreadDeleteCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"thread123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"thread123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := threadDeleteCmd.Args(threadDeleteCmd, tt.args)
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

func TestThreadDeleteCmd_ConfirmValidation(t *testing.T) {
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
			origConfirm := threadDeleteConfirm
			threadDeleteConfirm = tt.confirm
			defer func() { threadDeleteConfirm = origConfirm }()

			mockCmd := &cobra.Command{Use: "test"}
			mockCmd.SetOut(new(bytes.Buffer))
			mockCmd.SetErr(new(bytes.Buffer))

			err := threadDeleteCmd.PreRunE(mockCmd, []string{"thread123"})

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

func TestRunThreadDelete_WithMockDependencies(t *testing.T) {
	mockRepo := &MockThreadRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
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

	err := runThreadDelete(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadDelete failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "thread123") || !contains(output, "deleted") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunThreadDelete_Error(t *testing.T) {
	mockRepo := &MockThreadRepository{
		DeleteErr: fmt.Errorf("delete operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadDelete(cmd, []string{"thread123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to delete thread") {
		t.Errorf("expected error to contain 'failed to delete thread', got: %v", err)
	}
}

func TestRunThreadDelete_QuietMode(t *testing.T) {
	mockRepo := &MockThreadRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	quietFlag = true
	defer func() { quietFlag = origQuiet }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runThreadDelete(cmd, []string{"thread123"})
	if err != nil {
		t.Fatalf("runThreadDelete failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

// TestThreadCmd_NewSubcommandsRegistered verifies new subcommands are registered.
func TestThreadCmd_NewSubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":    false,
		"show":    false,
		"trash":   false,
		"untrash": false,
		"delete":  false,
		"modify":  false,
	}

	for _, sub := range threadCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with threadCmd", name)
		}
	}
}
