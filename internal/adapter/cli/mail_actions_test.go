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

func TestMailActionsCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "mail") {
		t.Error("expected output to contain 'mail'")
	}
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "archive") {
		t.Error("expected output to contain 'archive'")
	}
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
	if !contains(output, "mark") {
		t.Error("expected output to contain 'mark'")
	}
}

func TestMailTrashCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "trash", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "message") {
		t.Error("expected output to contain 'message'")
	}
}

func TestMailUntrashCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "untrash", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "untrash") {
		t.Error("expected output to contain 'untrash'")
	}
	if !contains(output, "restore") || !contains(output, "trash") {
		t.Error("expected output to contain 'restore' or 'trash'")
	}
}

func TestMailArchiveCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "archive", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "archive") {
		t.Error("expected output to contain 'archive'")
	}
	if !contains(output, "INBOX") {
		t.Error("expected output to contain 'INBOX'")
	}
}

func TestMailDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "delete", "--help"})

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

func TestMailModifyCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "modify", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
	if !contains(output, "--add-labels") {
		t.Error("expected output to contain '--add-labels'")
	}
	if !contains(output, "--remove-labels") {
		t.Error("expected output to contain '--remove-labels'")
	}
}

func TestMailMarkCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "mark", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "mark") {
		t.Error("expected output to contain 'mark'")
	}
	if !contains(output, "--read") {
		t.Error("expected output to contain '--read'")
	}
	if !contains(output, "--unread") {
		t.Error("expected output to contain '--unread'")
	}
	if !contains(output, "--star") {
		t.Error("expected output to contain '--star'")
	}
	if !contains(output, "--unstar") {
		t.Error("expected output to contain '--unstar'")
	}
}

func TestMailDeleteCmd_RequiresConfirmFlag(t *testing.T) {
	// Test that PreRunE validates the --confirm flag
	// We test the validation logic directly since Cobra flag parsing
	// behavior varies in test contexts
	mailDeleteConfirm = false

	mockCmd := &cobra.Command{Use: "test"}

	if mailDeleteCmd.PreRunE != nil {
		err := mailDeleteCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when --confirm flag is not set")
		}
	} else {
		t.Error("mailDeleteCmd should have PreRunE defined")
	}
}

func TestMailMarkCmd_RequiresAtLeastOneFlag(t *testing.T) {
	// Test that PreRunE validates at least one flag is set
	mailMarkRead = false
	mailMarkUnread = false
	mailMarkStar = false
	mailMarkUnstar = false

	mockCmd := &cobra.Command{Use: "test"}

	if mailMarkCmd.PreRunE != nil {
		err := mailMarkCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when no flags are set")
		}
	} else {
		t.Error("mailMarkCmd should have PreRunE defined")
	}
}

func TestMailModifyCmd_RequiresAtLeastOneFlag(t *testing.T) {
	// Test that PreRunE validates at least one flag is set
	mailModifyAddLabels = nil
	mailModifyRemoveLabels = nil

	mockCmd := &cobra.Command{Use: "test"}

	if mailModifyCmd.PreRunE != nil {
		err := mailModifyCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when no flags are set")
		}
	} else {
		t.Error("mailModifyCmd should have PreRunE defined")
	}
}

func TestMailMarkCmd_ConflictingReadFlags(t *testing.T) {
	// Test that --read and --unread cannot be used together
	mailMarkRead = true
	mailMarkUnread = true
	mailMarkStar = false
	mailMarkUnstar = false

	mockCmd := &cobra.Command{Use: "test"}

	if mailMarkCmd.PreRunE != nil {
		err := mailMarkCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when both --read and --unread are set")
		}
	} else {
		t.Error("mailMarkCmd should have PreRunE defined")
	}
}

func TestMailMarkCmd_ConflictingStarFlags(t *testing.T) {
	// Test that --star and --unstar cannot be used together
	mailMarkRead = false
	mailMarkUnread = false
	mailMarkStar = true
	mailMarkUnstar = true

	mockCmd := &cobra.Command{Use: "test"}

	if mailMarkCmd.PreRunE != nil {
		err := mailMarkCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when both --star and --unstar are set")
		}
	} else {
		t.Error("mailMarkCmd should have PreRunE defined")
	}
}

func TestMailTrashCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailTrashCmd.Args == nil {
		t.Error("mailTrashCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailTrashCmd.Args(mailTrashCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailUntrashCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailUntrashCmd.Args == nil {
		t.Error("mailUntrashCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailUntrashCmd.Args(mailUntrashCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailArchiveCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailArchiveCmd.Args == nil {
		t.Error("mailArchiveCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailArchiveCmd.Args(mailArchiveCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailDeleteCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailDeleteCmd.Args == nil {
		t.Error("mailDeleteCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailDeleteCmd.Args(mailDeleteCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailModifyCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailModifyCmd.Args == nil {
		t.Error("mailModifyCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailModifyCmd.Args(mailModifyCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailMarkCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailMarkCmd.Args == nil {
		t.Error("mailMarkCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailMarkCmd.Args(mailMarkCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailTrashCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailTrashCmd.Args(mailTrashCmd, tt.args)
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

func TestMailUntrashCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailUntrashCmd.Args(mailUntrashCmd, tt.args)
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

func TestMailArchiveCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailArchiveCmd.Args(mailArchiveCmd, tt.args)
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

func TestMailDeleteCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailDeleteCmd.Args(mailDeleteCmd, tt.args)
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

func TestMailModifyCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailModifyCmd.Args(mailModifyCmd, tt.args)
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

func TestMailMarkCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailMarkCmd.Args(mailMarkCmd, tt.args)
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

func TestMailDeleteCmd_ConfirmValidation(t *testing.T) {
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
			origConfirm := mailDeleteConfirm
			mailDeleteConfirm = tt.confirm

			mockCmd := &cobra.Command{Use: "test"}
			mockCmd.SetOut(new(bytes.Buffer))
			mockCmd.SetErr(new(bytes.Buffer))

			err := mailDeleteCmd.PreRunE(mockCmd, []string{"msg123"})

			mailDeleteConfirm = origConfirm

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

func TestMailMarkCmd_ValidSingleFlags(t *testing.T) {
	tests := []struct {
		name   string
		read   bool
		unread bool
		star   bool
		unstar bool
	}{
		{
			name: "only read",
			read: true,
		},
		{
			name:   "only unread",
			unread: true,
		},
		{
			name: "only star",
			star: true,
		},
		{
			name:   "only unstar",
			unstar: true,
		},
		{
			name: "read and star",
			read: true,
			star: true,
		},
		{
			name:   "unread and star",
			unread: true,
			star:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailMarkRead = tt.read
			mailMarkUnread = tt.unread
			mailMarkStar = tt.star
			mailMarkUnstar = tt.unstar

			mockCmd := &cobra.Command{Use: "test"}

			err := mailMarkCmd.PreRunE(mockCmd, []string{"msg123"})

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMailModifyCmd_ValidFlags(t *testing.T) {
	tests := []struct {
		name   string
		add    []string
		remove []string
	}{
		{
			name: "only add labels",
			add:  []string{"IMPORTANT"},
		},
		{
			name:   "only remove labels",
			remove: []string{"INBOX"},
		},
		{
			name:   "both add and remove",
			add:    []string{"IMPORTANT"},
			remove: []string{"INBOX"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailModifyAddLabels = tt.add
			mailModifyRemoveLabels = tt.remove

			mockCmd := &cobra.Command{Use: "test"}

			err := mailModifyCmd.PreRunE(mockCmd, []string{"msg123"})

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMailCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":    false,
		"read":    false,
		"search":  false,
		"trash":   false,
		"untrash": false,
		"archive": false,
		"delete":  false,
		"modify":  false,
		"mark":    false,
		"send":    false,
		"reply":   false,
		"forward": false,
	}

	for _, sub := range mailCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	// Note: not all subcommands may be registered yet
	registeredCount := 0
	for _, found := range subcommands {
		if found {
			registeredCount++
		}
	}

	if registeredCount == 0 {
		t.Error("expected at least some subcommands to be registered with mailCmd")
	}
}

func TestMailDeleteCmd_HasConfirmFlag(t *testing.T) {
	flag := mailDeleteCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be defined on delete command")
	}
}

func TestMailModifyCmd_HasFlags(t *testing.T) {
	flags := []string{"add-labels", "remove-labels"}

	for _, flagName := range flags {
		flag := mailModifyCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on modify command", flagName)
		}
	}
}

func TestMailMarkCmd_HasFlags(t *testing.T) {
	flags := []string{"read", "unread", "star", "unstar"}

	for _, flagName := range flags {
		flag := mailMarkCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on mark command", flagName)
		}
	}
}

// =============================================================================
// Tests using dependency injection with mocks
// =============================================================================

func TestRunMailList_WithMockDependencies(t *testing.T) {
	// Setup mock dependencies
	mockMessages := []*mail.Message{
		{ID: "msg1", Subject: "Test Subject 1", From: "sender1@example.com"},
		{ID: "msg2", Subject: "Test Subject 2", From: "sender2@example.com"},
	}

	mockRepo := &MockMessageRepository{
		Messages: mockMessages,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account: &accountuc.Account{
				Alias: "test",
				Email: "test@example.com",
			},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	// Save and restore global flags
	origFormat := formatFlag
	origMaxResults := mailListMaxResults
	origLabels := mailListLabels
	origUnreadOnly := mailListUnreadOnly
	formatFlag = "plain"
	mailListMaxResults = 10
	mailListLabels = []string{"INBOX"}
	mailListUnreadOnly = false
	defer func() {
		formatFlag = origFormat
		mailListMaxResults = origMaxResults
		mailListLabels = origLabels
		mailListUnreadOnly = origUnreadOnly
	}()

	// Create command and capture output
	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailList(cmd, []string{})
	if err != nil {
		t.Fatalf("runMailList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg1") || !contains(output, "msg2") {
		t.Errorf("expected output to contain message IDs, got: %s", output)
	}
}

func TestRunMailList_WithUnreadOnly(t *testing.T) {
	mockMessages := []*mail.Message{
		{ID: "unread1", Subject: "Unread Message", From: "sender@example.com"},
	}

	mockRepo := &MockMessageRepository{
		Messages: mockMessages,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origUnreadOnly := mailListUnreadOnly
	formatFlag = "plain"
	mailListUnreadOnly = true
	defer func() {
		formatFlag = origFormat
		mailListUnreadOnly = origUnreadOnly
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailList(cmd, []string{})
	if err != nil {
		t.Fatalf("runMailList failed: %v", err)
	}
}

func TestRunMailList_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ListErr: fmt.Errorf("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to list messages") {
		t.Errorf("expected error to contain 'failed to list messages', got: %v", err)
	}
}

func TestRunMailRead_WithMockDependencies(t *testing.T) {
	mockMessage := &mail.Message{
		ID:      "msg123",
		Subject: "Test Email Subject",
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Body:    "This is the test email body.",
	}

	mockRepo := &MockMessageRepository{
		Message: mockMessage,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailRead(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailRead failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Test Email Subject") {
		t.Errorf("expected output to contain subject, got: %s", output)
	}
	if !contains(output, "Message Body") {
		t.Errorf("expected output to contain 'Message Body', got: %s", output)
	}
}

func TestRunMailRead_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		GetErr: fmt.Errorf("message not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailRead(cmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to read message") {
		t.Errorf("expected error to contain 'failed to read message', got: %v", err)
	}
}

func TestRunMailSearch_WithMockDependencies(t *testing.T) {
	mockMessages := []*mail.Message{
		{ID: "search1", Subject: "Meeting Tomorrow", From: "boss@example.com"},
		{ID: "search2", Subject: "Meeting Update", From: "boss@example.com"},
	}

	mockRepo := &MockMessageRepository{
		SearchResult: &mail.ListResult[*mail.Message]{
			Items: mockMessages,
			Total: 2,
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailSearch(cmd, []string{"from:boss@example.com"})
	if err != nil {
		t.Fatalf("runMailSearch failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "search1") || !contains(output, "search2") {
		t.Errorf("expected output to contain search results, got: %s", output)
	}
}

func TestRunMailSearch_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		SearchErr: fmt.Errorf("search failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailSearch(cmd, []string{"invalid query"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to search messages") {
		t.Errorf("expected error to contain 'failed to search messages', got: %v", err)
	}
}

func TestRunMailTrash_WithMockDependencies(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailTrash(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailTrash failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "trash") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunMailTrash_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		TrashErr: fmt.Errorf("trash operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailTrash(cmd, []string{"msg123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to trash message") {
		t.Errorf("expected error to contain 'failed to trash message', got: %v", err)
	}
}

func TestRunMailUntrash_WithMockDependencies(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailUntrash(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailUntrash failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "restored") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunMailUntrash_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		UntrashErr: fmt.Errorf("untrash operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailUntrash(cmd, []string{"msg123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to restore message") {
		t.Errorf("expected error to contain 'failed to restore message', got: %v", err)
	}
}

func TestRunMailArchive_WithMockDependencies(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailArchive(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailArchive failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "archived") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunMailArchive_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ArchiveErr: fmt.Errorf("archive operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailArchive(cmd, []string{"msg123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to archive message") {
		t.Errorf("expected error to contain 'failed to archive message', got: %v", err)
	}
}

func TestRunMailDelete_WithMockDependencies(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailDelete(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailDelete failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "deleted") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunMailDelete_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		DeleteErr: fmt.Errorf("delete operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailDelete(cmd, []string{"msg123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to delete message") {
		t.Errorf("expected error to contain 'failed to delete message', got: %v", err)
	}
}

func TestRunMailModify_WithMockDependencies(t *testing.T) {
	mockMessage := &mail.Message{
		ID:     "msg123",
		Labels: []string{"INBOX", "IMPORTANT"},
	}
	mockRepo := &MockMessageRepository{
		ModifyResult: mockMessage,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origVerbose := verboseFlag
	origAddLabels := mailModifyAddLabels
	origRemoveLabels := mailModifyRemoveLabels
	quietFlag = false
	verboseFlag = true
	mailModifyAddLabels = []string{"IMPORTANT"}
	mailModifyRemoveLabels = []string{"INBOX"}
	defer func() {
		quietFlag = origQuiet
		verboseFlag = origVerbose
		mailModifyAddLabels = origAddLabels
		mailModifyRemoveLabels = origRemoveLabels
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailModify(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailModify failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "modified") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunMailModify_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ModifyErr: fmt.Errorf("modify operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origAddLabels := mailModifyAddLabels
	mailModifyAddLabels = []string{"IMPORTANT"}
	defer func() { mailModifyAddLabels = origAddLabels }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailModify(cmd, []string{"msg123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to modify message") {
		t.Errorf("expected error to contain 'failed to modify message', got: %v", err)
	}
}

func TestRunMailMark_MarkAsRead(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origRead := mailMarkRead
	origUnread := mailMarkUnread
	origStar := mailMarkStar
	origUnstar := mailMarkUnstar
	quietFlag = false
	mailMarkRead = true
	mailMarkUnread = false
	mailMarkStar = false
	mailMarkUnstar = false
	defer func() {
		quietFlag = origQuiet
		mailMarkRead = origRead
		mailMarkUnread = origUnread
		mailMarkStar = origStar
		mailMarkUnstar = origUnstar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMark(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMark failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "read") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunMailMark_StarMessage(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origRead := mailMarkRead
	origUnread := mailMarkUnread
	origStar := mailMarkStar
	origUnstar := mailMarkUnstar
	quietFlag = false
	mailMarkRead = false
	mailMarkUnread = false
	mailMarkStar = true
	mailMarkUnstar = false
	defer func() {
		quietFlag = origQuiet
		mailMarkRead = origRead
		mailMarkUnread = origUnread
		mailMarkStar = origStar
		mailMarkUnstar = origUnstar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMark(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMark failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "starred") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunMailMark_MultipleActions(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origRead := mailMarkRead
	origUnread := mailMarkUnread
	origStar := mailMarkStar
	origUnstar := mailMarkUnstar
	quietFlag = false
	mailMarkRead = true
	mailMarkUnread = false
	mailMarkStar = true
	mailMarkUnstar = false
	defer func() {
		quietFlag = origQuiet
		mailMarkRead = origRead
		mailMarkUnread = origUnread
		mailMarkStar = origStar
		mailMarkUnstar = origUnstar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMark(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMark failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
	if !contains(output, "read") || !contains(output, "starred") {
		t.Errorf("expected both read and starred in output, got: %s", output)
	}
}

func TestRunMailMark_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ModifyErr: fmt.Errorf("mark operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origRead := mailMarkRead
	mailMarkRead = true
	defer func() { mailMarkRead = origRead }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMark(cmd, []string{"msg123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to mark message") {
		t.Errorf("expected error to contain 'failed to mark message', got: %v", err)
	}
}

func TestRunMailList_QuietMode(t *testing.T) {
	mockMessages := []*mail.Message{
		{ID: "msg1", Subject: "Test Subject 1", From: "sender1@example.com"},
	}

	mockRepo := &MockMessageRepository{
		Messages: mockMessages,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailList(cmd, []string{})
	if err != nil {
		t.Fatalf("runMailList failed: %v", err)
	}
}

func TestRunMailSearch_EmptyResults(t *testing.T) {
	mockRepo := &MockMessageRepository{
		SearchResult: &mail.ListResult[*mail.Message]{
			Items: []*mail.Message{},
			Total: 0,
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailSearch(cmd, []string{"nonexistent query"})
	if err != nil {
		t.Fatalf("runMailSearch failed: %v", err)
	}
}

func TestRunMailSearch_WithPagination(t *testing.T) {
	mockMessages := []*mail.Message{
		{ID: "msg1", Subject: "Test 1", From: "sender@example.com"},
	}

	mockRepo := &MockMessageRepository{
		SearchResult: &mail.ListResult[*mail.Message]{
			Items: mockMessages,
			Total: 100, // Total is greater than returned
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailSearch(cmd, []string{"test query"})
	if err != nil {
		t.Fatalf("runMailSearch failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "showing first") {
		t.Errorf("expected pagination info in output, got: %s", output)
	}
}

func TestRunMailList_AccountResolveError(t *testing.T) {
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

	err := runMailList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "no account found") {
		t.Errorf("expected error to contain 'no account found', got: %v", err)
	}
}

func TestRunMailTrash_QuietMode(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailTrash(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailTrash failed: %v", err)
	}

	output := buf.String()
	// In quiet mode, should not show message
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunMailUntrash_QuietMode(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailUntrash(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailUntrash failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunMailArchive_QuietMode(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailArchive(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailArchive failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunMailDelete_QuietMode(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailDelete(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailDelete failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunMailModify_QuietMode(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origAddLabels := mailModifyAddLabels
	quietFlag = true
	mailModifyAddLabels = []string{"IMPORTANT"}
	defer func() {
		quietFlag = origQuiet
		mailModifyAddLabels = origAddLabels
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailModify(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailModify failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunMailMark_QuietMode(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origRead := mailMarkRead
	quietFlag = true
	mailMarkRead = true
	defer func() {
		quietFlag = origQuiet
		mailMarkRead = origRead
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMark(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMark failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunMailMark_MarkAsUnread(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origRead := mailMarkRead
	origUnread := mailMarkUnread
	origStar := mailMarkStar
	origUnstar := mailMarkUnstar
	quietFlag = false
	mailMarkRead = false
	mailMarkUnread = true
	mailMarkStar = false
	mailMarkUnstar = false
	defer func() {
		quietFlag = origQuiet
		mailMarkRead = origRead
		mailMarkUnread = origUnread
		mailMarkStar = origStar
		mailMarkUnstar = origUnstar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMark(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMark failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "unread") {
		t.Errorf("expected output to contain 'unread', got: %s", output)
	}
}

func TestRunMailMark_Unstar(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origRead := mailMarkRead
	origUnread := mailMarkUnread
	origStar := mailMarkStar
	origUnstar := mailMarkUnstar
	quietFlag = false
	mailMarkRead = false
	mailMarkUnread = false
	mailMarkStar = false
	mailMarkUnstar = true
	defer func() {
		quietFlag = origQuiet
		mailMarkRead = origRead
		mailMarkUnread = origUnread
		mailMarkStar = origStar
		mailMarkUnstar = origUnstar
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMark(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMark failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "unstarred") {
		t.Errorf("expected output to contain 'unstarred', got: %s", output)
	}
}

func TestRunMailRead_JSONFormat(t *testing.T) {
	mockMessage := &mail.Message{
		ID:      "msg123",
		Subject: "JSON Test Subject",
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
	}

	mockRepo := &MockMessageRepository{
		Message: mockMessage,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailRead(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailRead failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "JSON Test Subject") {
		t.Errorf("expected JSON output to contain subject, got: %s", output)
	}
}

func TestRunMailList_JSONFormat(t *testing.T) {
	mockMessages := []*mail.Message{
		{ID: "msg1", Subject: "JSON List Test", From: "sender@example.com"},
	}

	mockRepo := &MockMessageRepository{
		Messages: mockMessages,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailList(cmd, []string{})
	if err != nil {
		t.Fatalf("runMailList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "JSON List Test") {
		t.Errorf("expected JSON output to contain subject, got: %s", output)
	}
}

func TestRunMailSearch_JSONFormat(t *testing.T) {
	mockMessages := []*mail.Message{
		{ID: "msg1", Subject: "JSON Search Test", From: "sender@example.com"},
	}

	mockRepo := &MockMessageRepository{
		SearchResult: &mail.ListResult[*mail.Message]{
			Items: mockMessages,
			Total: 1,
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
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

	err := runMailSearch(cmd, []string{"test"})
	if err != nil {
		t.Fatalf("runMailSearch failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "JSON Search Test") {
		t.Errorf("expected JSON output to contain subject, got: %s", output)
	}
}

// =============================================================================
// Tests for mail move command
// =============================================================================

func TestMailMoveCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "move", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "move") {
		t.Error("expected output to contain 'move'")
	}
	if !contains(output, "--to") {
		t.Error("expected output to contain '--to'")
	}
	if !contains(output, "label") {
		t.Error("expected output to contain 'label'")
	}
}

func TestMailMoveCmd_HasArgsRequirement(t *testing.T) {
	if mailMoveCmd.Args == nil {
		t.Error("mailMoveCmd should have Args validator defined")
	}
}

func TestMailMoveCmd_HasToFlag(t *testing.T) {
	flag := mailMoveCmd.Flag("to")
	if flag == nil {
		t.Error("expected --to flag to be defined on move command")
	}
}

func TestMailMoveCmd_RequiresToFlag(t *testing.T) {
	origDestination := mailMoveDestination
	mailMoveDestination = ""
	defer func() { mailMoveDestination = origDestination }()

	mockCmd := &cobra.Command{Use: "test"}
	mockCmd.SetOut(new(bytes.Buffer))
	mockCmd.SetErr(new(bytes.Buffer))

	if mailMoveCmd.PreRunE != nil {
		err := mailMoveCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when --to flag is not set")
		}
	} else {
		t.Error("mailMoveCmd should have PreRunE defined")
	}
}

func TestMailMoveCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"msg123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"msg123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailMoveCmd.Args(mailMoveCmd, tt.args)
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

func TestMailMoveCmd_ToFlagValidation(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		expectErr   bool
	}{
		{
			name:        "with destination",
			destination: "IMPORTANT",
			expectErr:   false,
		},
		{
			name:        "empty destination",
			destination: "",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origDestination := mailMoveDestination
			mailMoveDestination = tt.destination
			defer func() { mailMoveDestination = origDestination }()

			mockCmd := &cobra.Command{Use: "test"}
			mockCmd.SetOut(new(bytes.Buffer))
			mockCmd.SetErr(new(bytes.Buffer))

			err := mailMoveCmd.PreRunE(mockCmd, []string{"msg123"})

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

func TestRunMailMove_WithMockDependencies(t *testing.T) {
	mockMessage := &mail.Message{
		ID:     "msg123",
		Labels: []string{"IMPORTANT"},
	}
	mockRepo := &MockMessageRepository{
		ModifyResult: mockMessage,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origDestination := mailMoveDestination
	quietFlag = false
	mailMoveDestination = "IMPORTANT"
	defer func() {
		quietFlag = origQuiet
		mailMoveDestination = origDestination
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMove(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMove failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "msg123") || !contains(output, "moved") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
	if !contains(output, "IMPORTANT") {
		t.Errorf("expected output to contain destination label, got: %s", output)
	}
}

func TestRunMailMove_Error(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ModifyErr: fmt.Errorf("move operation failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origDestination := mailMoveDestination
	mailMoveDestination = "IMPORTANT"
	defer func() { mailMoveDestination = origDestination }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMove(cmd, []string{"msg123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to move message") {
		t.Errorf("expected error to contain 'failed to move message', got: %v", err)
	}
}

func TestRunMailMove_QuietMode(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origDestination := mailMoveDestination
	quietFlag = true
	mailMoveDestination = "IMPORTANT"
	defer func() {
		quietFlag = origQuiet
		mailMoveDestination = origDestination
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMove(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMove failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunMailMove_VerboseMode(t *testing.T) {
	mockMessage := &mail.Message{
		ID:     "msg123",
		Labels: []string{"IMPORTANT", "STARRED"},
	}
	mockRepo := &MockMessageRepository{
		ModifyResult: mockMessage,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origQuiet := quietFlag
	origVerbose := verboseFlag
	origDestination := mailMoveDestination
	quietFlag = false
	verboseFlag = true
	mailMoveDestination = "IMPORTANT"
	defer func() {
		quietFlag = origQuiet
		verboseFlag = origVerbose
		mailMoveDestination = origDestination
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailMove(cmd, []string{"msg123"})
	if err != nil {
		t.Fatalf("runMailMove failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "labels") {
		t.Errorf("expected verbose output to contain label information, got: %s", output)
	}
}

// TestMailCmd_MoveSubcommandRegistered verifies the move subcommand is registered.
func TestMailCmd_MoveSubcommandRegistered(t *testing.T) {
	found := false
	for _, sub := range mailCmd.Commands() {
		if sub.Name() == "move" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected subcommand 'move' to be registered with mailCmd")
	}
}
