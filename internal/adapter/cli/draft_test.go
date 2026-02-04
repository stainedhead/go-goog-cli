// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

func TestDraftCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "draft") {
		t.Error("expected output to contain 'draft'")
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
	if !contains(output, "send") {
		t.Error("expected output to contain 'send'")
	}
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
}

func TestDraftListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "--limit") {
		t.Error("expected output to contain '--limit'")
	}
}

func TestDraftShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "show", "--help"})

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

func TestDraftCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "create", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
	if !contains(output, "--to") {
		t.Error("expected output to contain '--to'")
	}
	if !contains(output, "--subject") {
		t.Error("expected output to contain '--subject'")
	}
	if !contains(output, "--body") {
		t.Error("expected output to contain '--body'")
	}
}

func TestDraftUpdateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "update", "--help"})

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
}

func TestDraftSendCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "send", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "send") {
		t.Error("expected output to contain 'send'")
	}
	if !contains(output, "<id>") {
		t.Error("expected output to contain '<id>'")
	}
}

func TestDraftDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(draftCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"draft", "delete", "--help"})

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
}

func TestDraftShowCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftShowCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftUpdateCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftUpdateCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftSendCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftSendCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftDeleteCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if draftDeleteCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestDraftCmd_Aliases(t *testing.T) {
	tests := []struct {
		name    string
		command string
		alias   string
	}{
		{"list alias ls", "list", "ls"},
		{"show alias get", "show", "get"},
		{"show alias read", "show", "read"},
		{"delete alias rm", "delete", "rm"},
		{"delete alias remove", "delete", "remove"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find the subcommand
			var targetCmd *cobra.Command
			for _, sub := range draftCmd.Commands() {
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

func TestDraftShowCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"draft123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"draft123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := draftShowCmd.Args(draftShowCmd, tt.args)
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

func TestDraftUpdateCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"draft123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"draft123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := draftUpdateCmd.Args(draftUpdateCmd, tt.args)
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

func TestDraftSendCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"draft123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"draft123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := draftSendCmd.Args(draftSendCmd, tt.args)
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

func TestDraftDeleteCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"draft123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"draft123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := draftDeleteCmd.Args(draftDeleteCmd, tt.args)
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

func TestDraftCreateCmd_Validation(t *testing.T) {
	tests := []struct {
		name      string
		to        []string
		subject   string
		expectErr bool
	}{
		{
			name:      "empty to list",
			to:        []string{},
			subject:   "Test Subject",
			expectErr: true,
		},
		{
			name:      "nil to list",
			to:        nil,
			subject:   "Test Subject",
			expectErr: true,
		},
		{
			name:      "empty subject",
			to:        []string{"user@example.com"},
			subject:   "",
			expectErr: true,
		},
		{
			name:      "valid input",
			to:        []string{"user@example.com"},
			subject:   "Test Subject",
			expectErr: false,
		},
		{
			name:      "multiple recipients",
			to:        []string{"user1@example.com", "user2@example.com"},
			subject:   "Test Subject",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTo := draftTo
			origSubject := draftSubject

			draftTo = tt.to
			draftSubject = tt.subject

			mockCmd := &cobra.Command{Use: "test"}

			err := draftCreateCmd.PreRunE(mockCmd, []string{})

			draftTo = origTo
			draftSubject = origSubject

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

func TestDraftCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":   false,
		"show":   false,
		"create": false,
		"update": false,
		"send":   false,
		"delete": false,
	}

	for _, sub := range draftCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with draftCmd", name)
		}
	}
}

func TestDraftListCmd_HasFlags(t *testing.T) {
	flags := []string{"limit"}

	for _, flagName := range flags {
		flag := draftListCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on list command", flagName)
		}
	}
}

func TestDraftCreateCmd_HasFlags(t *testing.T) {
	flags := []string{"to", "subject", "body"}

	for _, flagName := range flags {
		flag := draftCreateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on create command", flagName)
		}
	}
}

func TestDraftUpdateCmd_HasFlags(t *testing.T) {
	flags := []string{"to", "subject", "body"}

	for _, flagName := range flags {
		flag := draftUpdateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on update command", flagName)
		}
	}
}

func TestDraftListCmd_DefaultLimit(t *testing.T) {
	flag := draftListCmd.Flag("limit")
	if flag == nil {
		t.Fatal("expected --limit flag to be set")
	}

	if flag.DefValue != "20" {
		t.Errorf("expected default limit to be '20', got '%s'", flag.DefValue)
	}
}

// =============================================================================
// Tests using dependency injection with mocks
// =============================================================================

func TestRunDraftList_WithMockDependencies(t *testing.T) {
	mockDrafts := []*mail.Draft{
		{ID: "draft1", Message: &mail.Message{Subject: "Test Draft 1", To: []string{"user@example.com"}}},
		{ID: "draft2", Message: &mail.Message{Subject: "Test Draft 2", To: []string{"user2@example.com"}}},
	}

	mockRepo := &MockDraftRepository{
		Drafts: mockDrafts,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origLimit := draftLimit
	formatFlag = "plain"
	draftLimit = 10
	defer func() {
		formatFlag = origFormat
		draftLimit = origLimit
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftList(cmd, []string{})
	if err != nil {
		t.Fatalf("runDraftList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "draft1") || !contains(output, "draft2") {
		t.Errorf("expected output to contain draft IDs, got: %s", output)
	}
}

func TestRunDraftList_Error(t *testing.T) {
	mockRepo := &MockDraftRepository{
		ListErr: fmt.Errorf("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to list drafts") {
		t.Errorf("expected error to contain 'failed to list drafts', got: %v", err)
	}
}

func TestRunDraftShow_WithMockDependencies(t *testing.T) {
	mockDraft := &mail.Draft{
		ID: "draft123",
		Message: &mail.Message{
			Subject: "Test Draft Subject",
			To:      []string{"recipient@example.com"},
			Body:    "This is the draft body.",
		},
	}

	mockRepo := &MockDraftRepository{
		Draft: mockDraft,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
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

	err := runDraftShow(cmd, []string{"draft123"})
	if err != nil {
		t.Fatalf("runDraftShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Test Draft Subject") {
		t.Errorf("expected output to contain subject, got: %s", output)
	}
}

func TestRunDraftShow_Error(t *testing.T) {
	mockRepo := &MockDraftRepository{
		GetErr: fmt.Errorf("draft not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftShow(cmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to get draft") {
		t.Errorf("expected error to contain 'failed to get draft', got: %v", err)
	}
}

func TestRunDraftCreate_WithMockDependencies(t *testing.T) {
	mockRepo := &MockDraftRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origTo := draftTo
	origSubject := draftSubject
	origBody := draftBody
	formatFlag = "plain"
	draftTo = []string{"user@example.com"}
	draftSubject = "Test Subject"
	draftBody = "Test body"
	defer func() {
		formatFlag = origFormat
		draftTo = origTo
		draftSubject = origSubject
		draftBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftCreate(cmd, []string{})
	if err != nil {
		t.Fatalf("runDraftCreate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Draft created successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunDraftCreate_Error(t *testing.T) {
	mockRepo := &MockDraftRepository{
		CreateErr: fmt.Errorf("create failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	draftTo = []string{"user@example.com"}
	draftSubject = "Test"
	defer func() {
		draftTo = nil
		draftSubject = ""
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftCreate(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to create draft") {
		t.Errorf("expected error to contain 'failed to create draft', got: %v", err)
	}
}

func TestRunDraftUpdate_WithMockDependencies(t *testing.T) {
	mockDraft := &mail.Draft{
		ID: "draft123",
		Message: &mail.Message{
			Subject: "Original Subject",
			To:      []string{"original@example.com"},
		},
	}

	mockRepo := &MockDraftRepository{
		Draft: mockDraft,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origSubject := draftSubject
	formatFlag = "plain"
	draftSubject = "Updated Subject"
	defer func() {
		formatFlag = origFormat
		draftSubject = origSubject
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftUpdate(cmd, []string{"draft123"})
	if err != nil {
		t.Fatalf("runDraftUpdate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Draft updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunDraftUpdate_Error(t *testing.T) {
	mockRepo := &MockDraftRepository{
		GetErr: fmt.Errorf("draft not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftUpdate(cmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunDraftSend_WithMockDependencies(t *testing.T) {
	mockRepo := &MockDraftRepository{
		SendResult: &mail.Message{
			ID:      "sent-msg-id",
			Subject: "Sent Subject",
			To:      []string{"recipient@example.com"},
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
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

	err := runDraftSend(cmd, []string{"draft123"})
	if err != nil {
		t.Fatalf("runDraftSend failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Draft sent successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunDraftSend_Error(t *testing.T) {
	mockRepo := &MockDraftRepository{
		SendErr: fmt.Errorf("send failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftSend(cmd, []string{"draft123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to send draft") {
		t.Errorf("expected error to contain 'failed to send draft', got: %v", err)
	}
}

func TestRunDraftDelete_WithMockDependencies(t *testing.T) {
	mockRepo := &MockDraftRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
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

	err := runDraftDelete(cmd, []string{"draft123"})
	if err != nil {
		t.Fatalf("runDraftDelete failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "draft123") || !contains(output, "deleted") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunDraftDelete_Error(t *testing.T) {
	mockRepo := &MockDraftRepository{
		DeleteErr: fmt.Errorf("delete failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftDelete(cmd, []string{"draft123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to delete draft") {
		t.Errorf("expected error to contain 'failed to delete draft', got: %v", err)
	}
}

func TestRunDraftDelete_QuietMode(t *testing.T) {
	mockRepo := &MockDraftRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
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

	err := runDraftDelete(cmd, []string{"draft123"})
	if err != nil {
		t.Fatalf("runDraftDelete failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunDraftList_JSONFormat(t *testing.T) {
	mockDrafts := []*mail.Draft{
		{ID: "draft1", Message: &mail.Message{Subject: "JSON Draft"}},
	}

	mockRepo := &MockDraftRepository{
		Drafts: mockDrafts,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
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

	err := runDraftList(cmd, []string{})
	if err != nil {
		t.Fatalf("runDraftList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "draft1") {
		t.Errorf("expected JSON output to contain draft ID, got: %s", output)
	}
}

func TestRunDraftCreate_JSONFormat(t *testing.T) {
	mockRepo := &MockDraftRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origTo := draftTo
	origSubject := draftSubject
	formatFlag = "json"
	draftTo = []string{"user@example.com"}
	draftSubject = "Test Subject"
	defer func() {
		formatFlag = origFormat
		draftTo = origTo
		draftSubject = origSubject
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runDraftCreate(cmd, []string{})
	if err != nil {
		t.Fatalf("runDraftCreate failed: %v", err)
	}
}

func TestRunDraftSend_JSONFormat(t *testing.T) {
	mockRepo := &MockDraftRepository{
		SendResult: &mail.Message{
			ID:      "sent-msg-id",
			Subject: "Sent Subject",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: mockRepo,
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

	err := runDraftSend(cmd, []string{"draft123"})
	if err != nil {
		t.Fatalf("runDraftSend failed: %v", err)
	}
}
