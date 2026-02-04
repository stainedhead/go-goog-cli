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

func TestACLCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "acl", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "acl") {
		t.Error("expected output to contain 'acl'")
	}
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "add") {
		t.Error("expected output to contain 'add'")
	}
	if !contains(output, "remove") {
		t.Error("expected output to contain 'remove'")
	}
}

func TestACLListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "acl", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "<calendar-id>") {
		t.Error("expected output to contain '<calendar-id>'")
	}
}

func TestACLAddCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "acl", "add", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "add") {
		t.Error("expected output to contain 'add'")
	}
	if !contains(output, "--email") {
		t.Error("expected output to contain '--email'")
	}
	if !contains(output, "--role") {
		t.Error("expected output to contain '--role'")
	}
}

func TestACLRemoveCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "acl", "remove", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "remove") {
		t.Error("expected output to contain 'remove'")
	}
	if !contains(output, "<calendar-id>") {
		t.Error("expected output to contain '<calendar-id>'")
	}
	if !contains(output, "<rule-id>") {
		t.Error("expected output to contain '<rule-id>'")
	}
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
}

func TestShareCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "share", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "share") {
		t.Error("expected output to contain 'share'")
	}
	if !contains(output, "--email") {
		t.Error("expected output to contain '--email'")
	}
	if !contains(output, "--role") {
		t.Error("expected output to contain '--role'")
	}
}

func TestUnshareCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "unshare", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "unshare") {
		t.Error("expected output to contain 'unshare'")
	}
	if !contains(output, "<calendar-id>") {
		t.Error("expected output to contain '<calendar-id>'")
	}
	if !contains(output, "<rule-id>") {
		t.Error("expected output to contain '<rule-id>'")
	}
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
}

func TestACLListCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if aclListCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestACLAddCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if aclAddCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestACLRemoveCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(2)
	if aclRemoveCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestShareCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if shareCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestUnshareCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(2)
	if unshareCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestACLAddCmd_HasEmailFlag(t *testing.T) {
	// Verify the command has an --email flag
	flag := aclAddCmd.Flag("email")
	if flag == nil {
		t.Error("expected --email flag to be set")
	}
}

func TestACLAddCmd_HasRoleFlag(t *testing.T) {
	// Verify the command has a --role flag
	flag := aclAddCmd.Flag("role")
	if flag == nil {
		t.Error("expected --role flag to be set")
	}
}

func TestACLRemoveCmd_HasConfirmFlag(t *testing.T) {
	// Verify the command has a --confirm flag
	flag := aclRemoveCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}

func TestShareCmd_HasEmailFlag(t *testing.T) {
	// Verify the share command has an --email flag
	flag := shareCmd.Flag("email")
	if flag == nil {
		t.Error("expected --email flag to be set")
	}
}

func TestShareCmd_HasRoleFlag(t *testing.T) {
	// Verify the share command has a --role flag
	flag := shareCmd.Flag("role")
	if flag == nil {
		t.Error("expected --role flag to be set")
	}
}

func TestUnshareCmd_HasConfirmFlag(t *testing.T) {
	// Verify the unshare command has a --confirm flag
	flag := unshareCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}

func TestACLCmd_Aliases(t *testing.T) {
	tests := []struct {
		name    string
		command string
		alias   string
	}{
		{"list alias ls", "list", "ls"},
		{"remove alias rm", "remove", "rm"},
		{"remove alias delete", "remove", "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find the subcommand
			var targetCmd *cobra.Command
			for _, sub := range aclCmd.Commands() {
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

func TestACLAddCmd_EmailFlagRequired(t *testing.T) {
	// Verify that PreRunE validation exists for email flag
	if aclAddCmd.PreRunE == nil {
		t.Error("expected PreRunE validation for email flag")
	}
}

func TestShareCmd_EmailFlagRequired(t *testing.T) {
	// Verify that PreRunE validation exists for email flag
	if shareCmd.PreRunE == nil {
		t.Error("expected PreRunE validation for email flag")
	}
}

func TestACLListCmd_ViaAlias(t *testing.T) {
	// Test that 'goog cal acl ls' works
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "acl", "ls", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") || !contains(output, "sharing rules") {
		t.Error("expected output to describe list command")
	}
}

func TestACLRemoveCmd_ViaAlias(t *testing.T) {
	// Test that 'goog cal acl rm' works
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(calCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"cal", "acl", "rm", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "remove") || !contains(output, "--confirm") {
		t.Error("expected output to describe remove command with confirm flag")
	}
}

func TestCalCmd_HasACLSubcommand(t *testing.T) {
	// Verify that cal command has acl subcommand
	found := false
	for _, sub := range calCmd.Commands() {
		if sub.Use == "acl" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected cal command to have acl subcommand")
	}
}

func TestCalCmd_HasShareSubcommand(t *testing.T) {
	// Verify that cal command has share subcommand
	found := false
	for _, sub := range calCmd.Commands() {
		if sub.Use == "share <calendar-id>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected cal command to have share subcommand")
	}
}

func TestCalCmd_HasUnshareSubcommand(t *testing.T) {
	// Verify that cal command has unshare subcommand
	found := false
	for _, sub := range calCmd.Commands() {
		if sub.Use == "unshare <calendar-id> <rule-id>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected cal command to have unshare subcommand")
	}
}

func TestACLListCmd_ArgsValidation(t *testing.T) {
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
			err := aclListCmd.Args(aclListCmd, tt.args)
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

func TestACLAddCmd_ArgsValidation(t *testing.T) {
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
			err := aclAddCmd.Args(aclAddCmd, tt.args)
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

func TestACLRemoveCmd_ArgsValidation(t *testing.T) {
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
			expectErr: true, // needs 2 args
		},
		{
			name:      "two args",
			args:      []string{"primary", "ruleId"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"primary", "ruleId", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aclRemoveCmd.Args(aclRemoveCmd, tt.args)
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

func TestShareCmd_ArgsValidation(t *testing.T) {
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
			err := shareCmd.Args(shareCmd, tt.args)
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

func TestUnshareCmd_ArgsValidation(t *testing.T) {
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
			expectErr: true, // needs 2 args
		},
		{
			name:      "two args",
			args:      []string{"primary", "ruleId"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"primary", "ruleId", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unshareCmd.Args(unshareCmd, tt.args)
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

func TestACLCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":   false,
		"add":    false,
		"remove": false,
	}

	for _, sub := range aclCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with aclCmd", name)
		}
	}
}

// =============================================================================
// Execution tests for ACL operations
// =============================================================================

func TestRunACLList_Success(t *testing.T) {
	mockRules := []*calendar.ACLRule{
		{ID: "user:test@example.com", Role: "owner", Scope: &calendar.ACLScope{Type: "user", Value: "test@example.com"}},
		{ID: "user:other@example.com", Role: "reader", Scope: &calendar.ACLScope{Type: "user", Value: "other@example.com"}},
	}

	mockRepo := &MockACLRepository{
		Rules: mockRules,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: mockRepo,
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

	err := runACLList(cmd, []string{"primary"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "test@example.com") || !contains(output, "other@example.com") {
		t.Error("expected output to contain email addresses")
	}
}

func TestRunACLList_Error(t *testing.T) {
	mockRepo := &MockACLRepository{
		ListErr: fmt.Errorf("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runACLList(cmd, []string{"primary"})

	if err == nil {
		t.Error("expected error from repository")
	}
	if !contains(err.Error(), "failed to list ACL rules") {
		t.Errorf("expected list error, got: %v", err)
	}
}

func TestRunACLAdd_Success(t *testing.T) {
	mockRule := &calendar.ACLRule{
		ID:    "user:user@example.com",
		Role:  "reader",
		Scope: &calendar.ACLScope{Type: "user", Value: "user@example.com"},
	}

	mockRepo := &MockACLRepository{
		InsertResult: mockRule,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origEmail := aclEmail
	origRole := aclRole
	origFormat := formatFlag

	aclEmail = "user@example.com"
	aclRole = "reader"
	formatFlag = "plain"

	defer func() {
		aclEmail = origEmail
		aclRole = origRole
		formatFlag = origFormat
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runACLAdd(cmd, []string{"primary"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "user@example.com") {
		t.Error("expected output to contain email address")
	}
	if !contains(output, "reader") {
		t.Error("expected output to contain role")
	}
}

func TestRunACLAdd_InvalidRole(t *testing.T) {
	mockRepo := &MockACLRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origEmail := aclEmail
	origRole := aclRole

	aclEmail = "user@example.com"
	aclRole = "invalid-role"

	defer func() {
		aclEmail = origEmail
		aclRole = origRole
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runACLAdd(cmd, []string{"primary"})

	if err == nil {
		t.Error("expected error for invalid role")
	}
	if !contains(err.Error(), "invalid role") {
		t.Errorf("expected invalid role error, got: %v", err)
	}
}

func TestRunACLRemove_Success(t *testing.T) {
	mockRule := &calendar.ACLRule{
		ID:    "user:user@example.com",
		Role:  "reader",
		Scope: &calendar.ACLScope{Type: "user", Value: "user@example.com"},
	}

	mockRepo := &MockACLRepository{
		Rule: mockRule,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origConfirm := aclConfirm
	origQuiet := quietFlag

	aclConfirm = true
	quietFlag = false

	defer func() {
		aclConfirm = origConfirm
		quietFlag = origQuiet
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runACLRemove(cmd, []string{"primary", "user:user@example.com"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "removed successfully") {
		t.Error("expected success message in output")
	}
}

func TestRunACLRemove_WithoutConfirm(t *testing.T) {
	mockRepo := &MockACLRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origConfirm := aclConfirm
	aclConfirm = false
	defer func() { aclConfirm = origConfirm }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runACLRemove(cmd, []string{"primary", "user:user@example.com"})

	if err == nil {
		t.Error("expected error for missing confirmation")
	}
	if !contains(err.Error(), "confirm") {
		t.Errorf("expected confirm error, got: %v", err)
	}
}

func TestRunACLRemove_RuleNotFound(t *testing.T) {
	mockRepo := &MockACLRepository{
		GetErr: fmt.Errorf("rule not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origConfirm := aclConfirm
	aclConfirm = true
	defer func() { aclConfirm = origConfirm }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runACLRemove(cmd, []string{"primary", "nonexistent"})

	if err == nil {
		t.Error("expected error for nonexistent rule")
	}
	if !contains(err.Error(), "ACL rule not found") {
		t.Errorf("expected rule not found error, got: %v", err)
	}
}
