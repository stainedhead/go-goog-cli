// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

func TestAccountCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "account") {
		t.Error("expected output to contain 'account'")
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
	if !contains(output, "switch") {
		t.Error("expected output to contain 'switch'")
	}
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "rename") {
		t.Error("expected output to contain 'rename'")
	}
}

func TestAccountListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
}

func TestAccountAddCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "add", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "add") {
		t.Error("expected output to contain 'add'")
	}
	if !contains(output, "--scopes") {
		t.Error("expected output to contain '--scopes'")
	}
}

func TestAccountRemoveCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "remove"}) // Missing required argument

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestAccountSwitchCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "switch"}) // Missing required argument

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestAccountRenameCmd_Args(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Test with no arguments
	cmd.SetArgs([]string{"account", "rename"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing arguments")
	}

	// Test with one argument
	cmd.SetArgs([]string{"account", "rename", "old"})
	err = cmd.Execute()
	if err == nil {
		t.Error("expected error for missing second argument")
	}
}

func TestAccountShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(accountCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"account", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
}

func TestAccountCmd_Aliases(t *testing.T) {
	tests := []struct {
		alias   string
		command string
	}{
		{"ls", "list"},
		{"rm", "remove"},
		{"use", "switch"},
		{"current", "show"},
		{"mv", "rename"},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"account", tt.alias, "--help"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error for alias %s: %v", tt.alias, err)
			}
		})
	}
}

func TestAccountCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":   false,
		"add":    false,
		"remove": false,
		"switch": false,
		"show":   false,
		"rename": false,
	}

	for _, sub := range accountCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with accountCmd", name)
		}
	}
}

func TestAccountRemoveCmd_HasArgsRequirement(t *testing.T) {
	if accountRemoveCmd.Args == nil {
		t.Error("expected Args to be set on remove command")
	}
}

func TestAccountSwitchCmd_HasArgsRequirement(t *testing.T) {
	if accountSwitchCmd.Args == nil {
		t.Error("expected Args to be set on switch command")
	}
}

func TestAccountRenameCmd_HasArgsRequirement(t *testing.T) {
	if accountRenameCmd.Args == nil {
		t.Error("expected Args to be set on rename command")
	}
}

func TestAccountAddCmd_HasScopesFlag(t *testing.T) {
	flag := accountAddCmd.Flag("scopes")
	if flag == nil {
		t.Error("expected --scopes flag to be defined on add command")
	}
}

func TestAccountRemoveCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"alias"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"alias", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountRemoveCmd.Args(accountRemoveCmd, tt.args)
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

func TestAccountSwitchCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"alias"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"alias", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountSwitchCmd.Args(accountSwitchCmd, tt.args)
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

func TestAccountRenameCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"old"},
			expectErr: true,
		},
		{
			name:      "two args",
			args:      []string{"old", "new"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"old", "new", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountRenameCmd.Args(accountRenameCmd, tt.args)
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

func TestAccountAddCmd_ArgsValidation(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "no args - valid",
			args:      []string{},
			expectErr: false,
		},
		{
			name:      "one arg - valid",
			args:      []string{"alias"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"alias", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accountAddCmd.Args(accountAddCmd, tt.args)
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

func TestOutputAccountsTable(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "default",
			Email:     "user@example.com",
			IsDefault: true,
			Added:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Alias:     "work",
			Email:     "work@company.com",
			IsDefault: false,
			Added:     time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC),
		},
	}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsTable(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Check for table headers
	if !contains(output, "ALIAS") {
		t.Error("expected output to contain 'ALIAS'")
	}
	if !contains(output, "EMAIL") {
		t.Error("expected output to contain 'EMAIL'")
	}
	if !contains(output, "DEFAULT") {
		t.Error("expected output to contain 'DEFAULT'")
	}
	// Check for account data
	if !contains(output, "default") {
		t.Error("expected output to contain 'default'")
	}
	if !contains(output, "user@example.com") {
		t.Error("expected output to contain 'user@example.com'")
	}
	if !contains(output, "work") {
		t.Error("expected output to contain 'work'")
	}
}

func TestOutputAccountsPlain(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "default",
			Email:     "user@example.com",
			IsDefault: true,
			Added:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Alias:     "work",
			Email:     "work@company.com",
			IsDefault: false,
			Added:     time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC),
		},
	}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsPlain(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Check for account data in plain format
	if !contains(output, "default:") {
		t.Error("expected output to contain 'default:'")
	}
	if !contains(output, "user@example.com") {
		t.Error("expected output to contain 'user@example.com'")
	}
	if !contains(output, "(default)") {
		t.Error("expected output to contain '(default)' marker")
	}
	if !contains(output, "work:") {
		t.Error("expected output to contain 'work:'")
	}
}

func TestOutputAccountsTable_EmptyList(t *testing.T) {
	accounts := []*accountuc.Account{}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsTable(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Should still have headers
	if !contains(output, "ALIAS") {
		t.Error("expected output to contain 'ALIAS' header even for empty list")
	}
}

func TestOutputAccountsPlain_EmptyList(t *testing.T) {
	accounts := []*accountuc.Account{}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsPlain(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Should be empty or minimal for no accounts
	if output != "" {
		t.Errorf("expected empty output for empty account list, got: %s", output)
	}
}

// =============================================================================
// Execution Tests with Mocks
// =============================================================================

// MockAccountServiceFull extends MockAccountService with all needed methods.
type MockAccountServiceFull struct {
	MockAccountService
	AddFunc    func(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error)
	RemoveFunc func(alias string) error
	SwitchFunc func(alias string) error
	RenameFunc func(oldAlias, newAlias string) error
	AddResult  *accountuc.Account
	AddErr     error
	RemoveErr  error
	SwitchErr  error
	RenameErr  error
}

func (m *MockAccountServiceFull) Add(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error) {
	if m.AddFunc != nil {
		return m.AddFunc(ctx, alias, scopes)
	}
	if m.AddErr != nil {
		return nil, m.AddErr
	}
	if m.AddResult != nil {
		return m.AddResult, nil
	}
	return &accountuc.Account{
		Alias:     alias,
		Email:     "test@example.com",
		IsDefault: true,
	}, nil
}

func (m *MockAccountServiceFull) Remove(alias string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(alias)
	}
	return m.RemoveErr
}

func (m *MockAccountServiceFull) Switch(alias string) error {
	if m.SwitchFunc != nil {
		return m.SwitchFunc(alias)
	}
	return m.SwitchErr
}

func (m *MockAccountServiceFull) Rename(oldAlias, newAlias string) error {
	if m.RenameFunc != nil {
		return m.RenameFunc(oldAlias, newAlias)
	}
	return m.RenameErr
}

// TestAccountList_Execution tests account list command execution with mocks.
func TestAccountList_Execution(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		setupMock     func(*MockAccountServiceFull)
		expectError   bool
		expectOutputs []string
	}{
		{
			name: "list with accounts - table format",
			args: []string{"account", "list"},
			setupMock: func(m *MockAccountServiceFull) {
				m.Accounts = []*accountuc.Account{
					{
						Alias:     "default",
						Email:     "user@example.com",
						IsDefault: true,
						Added:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					},
					{
						Alias:     "work",
						Email:     "work@company.com",
						IsDefault: false,
						Added:     time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC),
					},
				}
			},
			expectError: false,
			expectOutputs: []string{
				"ALIAS",
				"EMAIL",
				"default",
				"user@example.com",
				"work",
				"work@company.com",
			},
		},
		{
			name: "list with no accounts",
			args: []string{"account", "list"},
			setupMock: func(m *MockAccountServiceFull) {
				m.Accounts = []*accountuc.Account{}
			},
			expectError: false,
			expectOutputs: []string{
				"No accounts configured",
			},
		},
		{
			name: "list with format flag - json",
			args: []string{"account", "list", "--format", "json"},
			setupMock: func(m *MockAccountServiceFull) {
				m.Accounts = []*accountuc.Account{
					{
						Alias:     "test",
						Email:     "test@example.com",
						IsDefault: true,
						Added:     time.Now(),
					},
				}
			},
			expectError: false,
			expectOutputs: []string{
				"\"alias\": \"test\"",
				"\"email\": \"test@example.com\"",
			},
		},
		{
			name: "list with format flag - plain",
			args: []string{"account", "list", "--format", "plain"},
			setupMock: func(m *MockAccountServiceFull) {
				m.Accounts = []*accountuc.Account{
					{
						Alias:     "test",
						Email:     "test@example.com",
						IsDefault: true,
						Added:     time.Now(),
					},
				}
			},
			expectError: false,
			expectOutputs: []string{
				"test: test@example.com (default)",
			},
		},
		{
			name: "list with error",
			args: []string{"account", "list"},
			setupMock: func(m *MockAccountServiceFull) {
				m.ListErr = fmt.Errorf("failed to list")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAccountServiceFull{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Set dependencies
			deps := &Dependencies{
				AccountService: mockSvc,
				RepoFactory:    &MockRepositoryFactory{},
			}
			SetDependencies(deps)
			defer ResetDependencies()

			// Create command
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Reset global flags
			formatFlag = "table"
			accountFlag = ""

			// Execute
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check outputs
			if len(tt.expectOutputs) > 0 {
				output := buf.String()
				for _, expected := range tt.expectOutputs {
					if !contains(output, expected) {
						t.Errorf("expected output to contain %q, got: %s", expected, output)
					}
				}
			}
		})
	}
}

// TestAccountAdd_Execution tests account add command execution with mocks.
func TestAccountAdd_Execution(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		setupMock    func(*MockAccountServiceFull)
		expectError  bool
		expectOutput string
	}{
		{
			name: "add account with alias",
			args: []string{"account", "add", "work"},
			setupMock: func(m *MockAccountServiceFull) {
				m.AddResult = &accountuc.Account{
					Alias:     "work",
					Email:     "work@company.com",
					IsDefault: false,
				}
			},
			expectError:  false,
			expectOutput: "Successfully added account 'work'",
		},
		{
			name: "add account without alias defaults to 'default'",
			args: []string{"account", "add"},
			setupMock: func(m *MockAccountServiceFull) {
				m.AddFunc = func(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error) {
					if alias != "default" {
						t.Errorf("expected alias 'default', got %q", alias)
					}
					return &accountuc.Account{
						Alias:     alias,
						Email:     "user@example.com",
						IsDefault: true,
					}, nil
				}
			},
			expectError:  false,
			expectOutput: "Successfully added account 'default'",
		},
		{
			name: "add account with scopes",
			args: []string{"account", "add", "work", "--scopes", "gmail,calendar"},
			setupMock: func(m *MockAccountServiceFull) {
				m.AddFunc = func(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error) {
					if len(scopes) < 2 {
						t.Errorf("expected at least 2 scopes, got %d", len(scopes))
					}
					return &accountuc.Account{
						Alias:     alias,
						Email:     "work@company.com",
						IsDefault: false,
					}, nil
				}
			},
			expectError:  false,
			expectOutput: "Successfully added account",
		},
		{
			name: "add account with error",
			args: []string{"account", "add", "work"},
			setupMock: func(m *MockAccountServiceFull) {
				m.AddErr = fmt.Errorf("authentication failed")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAccountServiceFull{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Set dependencies
			deps := &Dependencies{
				AccountService: mockSvc,
				RepoFactory:    &MockRepositoryFactory{},
			}
			SetDependencies(deps)
			defer ResetDependencies()

			// Create command
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Reset global flags
			accountFlag = ""
			accountAddScopes = nil

			// Execute
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output
			if tt.expectOutput != "" {
				output := buf.String()
				if !contains(output, tt.expectOutput) {
					t.Errorf("expected output to contain %q, got: %s", tt.expectOutput, output)
				}
			}
		})
	}
}

// TestAccountRemove_Execution tests account remove command execution with mocks.
func TestAccountRemove_Execution(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		setupMock    func(*MockAccountServiceFull)
		expectError  bool
		expectOutput string
	}{
		{
			name: "remove existing account",
			args: []string{"account", "remove", "work"},
			setupMock: func(m *MockAccountServiceFull) {
				m.Account = &accountuc.Account{
					Alias: "work",
					Email: "work@company.com",
				}
			},
			expectError:  false,
			expectOutput: "Successfully removed account 'work'",
		},
		{
			name: "remove non-existent account",
			args: []string{"account", "remove", "nonexistent"},
			setupMock: func(m *MockAccountServiceFull) {
				m.ResolveErr = fmt.Errorf("account not found")
			},
			expectError: true,
		},
		{
			name: "remove with error",
			args: []string{"account", "remove", "work"},
			setupMock: func(m *MockAccountServiceFull) {
				m.Account = &accountuc.Account{
					Alias: "work",
					Email: "work@company.com",
				}
				m.RemoveErr = fmt.Errorf("failed to remove")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAccountServiceFull{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Set dependencies
			deps := &Dependencies{
				AccountService: mockSvc,
				RepoFactory:    &MockRepositoryFactory{},
			}
			SetDependencies(deps)
			defer ResetDependencies()

			// Create command
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Execute
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output
			if tt.expectOutput != "" {
				output := buf.String()
				if !contains(output, tt.expectOutput) {
					t.Errorf("expected output to contain %q, got: %s", tt.expectOutput, output)
				}
			}
		})
	}
}

// TestAccountSwitch_Execution tests account switch command execution with mocks.
func TestAccountSwitch_Execution(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		setupMock    func(*MockAccountServiceFull)
		expectError  bool
		expectOutput string
	}{
		{
			name: "switch to existing account",
			args: []string{"account", "switch", "work"},
			setupMock: func(m *MockAccountServiceFull) {
				m.Account = &accountuc.Account{
					Alias: "work",
					Email: "work@company.com",
				}
			},
			expectError:  false,
			expectOutput: "Switched to account 'work'",
		},
		{
			name: "switch to non-existent account",
			args: []string{"account", "switch", "nonexistent"},
			setupMock: func(m *MockAccountServiceFull) {
				m.SwitchErr = fmt.Errorf("account not found")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAccountServiceFull{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Set dependencies
			deps := &Dependencies{
				AccountService: mockSvc,
				RepoFactory:    &MockRepositoryFactory{},
			}
			SetDependencies(deps)
			defer ResetDependencies()

			// Create command
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Execute
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output
			if tt.expectOutput != "" {
				output := buf.String()
				if !contains(output, tt.expectOutput) {
					t.Errorf("expected output to contain %q, got: %s", tt.expectOutput, output)
				}
			}
		})
	}
}

// TestAccountShow_Execution tests account show command execution with mocks.
func TestAccountShow_Execution(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		setupMock     func(*MockAccountServiceFull, *MockTokenManager)
		expectError   bool
		expectOutputs []string
	}{
		{
			name: "show current account with token",
			args: []string{"account", "show"},
			setupMock: func(m *MockAccountServiceFull, tm *MockTokenManager) {
				m.Account = &accountuc.Account{
					Alias:     "default",
					Email:     "user@example.com",
					IsDefault: true,
					Added:     time.Now(),
					Scopes:    []string{"email", "openid"},
				}
				m.TokenManager = tm
			},
			expectError: false,
			expectOutputs: []string{
				"Alias:",
				"default",
				"Email:",
				"user@example.com",
				"Default:",
				"true",
			},
		},
		{
			name: "show with no account",
			args: []string{"account", "show"},
			setupMock: func(m *MockAccountServiceFull, tm *MockTokenManager) {
				m.ResolveErr = fmt.Errorf("no account found")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockTM := &MockTokenManager{}
			mockSvc := &MockAccountServiceFull{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc, mockTM)
			}

			// Set dependencies
			deps := &Dependencies{
				AccountService: mockSvc,
				RepoFactory:    &MockRepositoryFactory{},
			}
			SetDependencies(deps)
			defer ResetDependencies()

			// Create command
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Reset global flags
			accountFlag = ""

			// Execute
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check outputs
			if len(tt.expectOutputs) > 0 {
				output := buf.String()
				for _, expected := range tt.expectOutputs {
					if !contains(output, expected) {
						t.Errorf("expected output to contain %q, got: %s", expected, output)
					}
				}
			}
		})
	}
}

// TestAccountRename_Execution tests account rename command execution with mocks.
func TestAccountRename_Execution(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		setupMock    func(*MockAccountServiceFull)
		expectError  bool
		expectOutput string
	}{
		{
			name: "rename existing account",
			args: []string{"account", "rename", "old", "new"},
			setupMock: func(m *MockAccountServiceFull) {
				// No error - successful rename
			},
			expectError:  false,
			expectOutput: "Successfully renamed account 'old' to 'new'",
		},
		{
			name: "rename non-existent account",
			args: []string{"account", "rename", "nonexistent", "new"},
			setupMock: func(m *MockAccountServiceFull) {
				m.RenameErr = fmt.Errorf("account not found")
			},
			expectError: true,
		},
		{
			name: "rename with duplicate name",
			args: []string{"account", "rename", "old", "existing"},
			setupMock: func(m *MockAccountServiceFull) {
				m.RenameErr = fmt.Errorf("account already exists")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAccountServiceFull{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Set dependencies
			deps := &Dependencies{
				AccountService: mockSvc,
				RepoFactory:    &MockRepositoryFactory{},
			}
			SetDependencies(deps)
			defer ResetDependencies()

			// Create command
			cmd := &cobra.Command{Use: "goog"}
			cmd.AddCommand(accountCmd)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Execute
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output
			if tt.expectOutput != "" {
				output := buf.String()
				if !contains(output, tt.expectOutput) {
					t.Errorf("expected output to contain %q, got: %s", tt.expectOutput, output)
				}
			}
		})
	}
}
