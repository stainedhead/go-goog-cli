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

func TestLabelCmd_Help(t *testing.T) {
	// Create a new root command for testing
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "--help"})

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check output contains expected content
	output := buf.String()
	if !contains(output, "label") {
		t.Error("expected output to contain 'label'")
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
}

func TestLabelListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
	if !contains(output, "system labels") {
		t.Error("expected output to contain 'system labels'")
	}
}

func TestLabelShowCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "show", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "show") {
		t.Error("expected output to contain 'show'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
}

func TestLabelCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "create", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
	if !contains(output, "--background") {
		t.Error("expected output to contain '--background'")
	}
	if !contains(output, "--text") {
		t.Error("expected output to contain '--text'")
	}
}

func TestLabelUpdateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "update", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "update") {
		t.Error("expected output to contain 'update'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
	if !contains(output, "--background") {
		t.Error("expected output to contain '--background'")
	}
}

func TestLabelDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(labelCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"label", "delete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
	if !contains(output, "<name>") {
		t.Error("expected output to contain '<name>'")
	}
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
}

func TestLabelShowCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelShowCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelCreateCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelCreateCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelUpdateCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelUpdateCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelDeleteCmd_HasArgsRequirement(t *testing.T) {
	// Verify the command has Args set to ExactArgs(1)
	if labelDeleteCmd.Args == nil {
		t.Error("expected Args to be set")
	}
}

func TestLabelCmd_Aliases(t *testing.T) {
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
			for _, sub := range labelCmd.Commands() {
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

func TestLabelDeleteCmd_HasConfirmFlag(t *testing.T) {
	// Verify the command has a --confirm flag
	flag := labelDeleteCmd.Flag("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}

func TestLabelShowCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"MyLabel"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"MyLabel", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := labelShowCmd.Args(labelShowCmd, tt.args)
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

func TestLabelCreateCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"NewLabel"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"NewLabel", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := labelCreateCmd.Args(labelCreateCmd, tt.args)
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

func TestLabelUpdateCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"MyLabel"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"MyLabel", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := labelUpdateCmd.Args(labelUpdateCmd, tt.args)
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

func TestLabelDeleteCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"MyLabel"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"MyLabel", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := labelDeleteCmd.Args(labelDeleteCmd, tt.args)
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

func TestLabelCmd_SubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":   false,
		"show":   false,
		"create": false,
		"update": false,
		"delete": false,
	}

	for _, sub := range labelCmd.Commands() {
		if _, ok := subcommands[sub.Name()]; ok {
			subcommands[sub.Name()] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %s to be registered with labelCmd", name)
		}
	}
}

func TestLabelCreateCmd_HasFlags(t *testing.T) {
	flags := []string{"background", "text"}

	for _, flagName := range flags {
		flag := labelCreateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on create command", flagName)
		}
	}
}

func TestLabelUpdateCmd_HasFlags(t *testing.T) {
	flags := []string{"background", "text"}

	for _, flagName := range flags {
		flag := labelUpdateCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("expected --%s flag to be defined on update command", flagName)
		}
	}
}

func TestLabelDeleteCmd_ConfirmValidation(t *testing.T) {
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
			origConfirm := labelConfirm
			labelConfirm = tt.confirm

			mockCmd := &cobra.Command{Use: "test"}

			err := labelDeleteCmd.PreRunE(mockCmd, []string{"MyLabel"})

			labelConfirm = origConfirm

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

func TestRunLabelList_WithMockDependencies(t *testing.T) {
	mockLabels := []*mail.Label{
		{ID: "label1", Name: "Work", Type: mail.LabelTypeUser},
		{ID: "label2", Name: "Personal", Type: mail.LabelTypeUser},
		{ID: "INBOX", Name: "INBOX", Type: mail.LabelTypeSystem},
	}

	mockRepo := &MockLabelRepository{
		Labels: mockLabels,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelList(cmd, []string{})
	if err != nil {
		t.Fatalf("runLabelList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Work") || !contains(output, "Personal") {
		t.Errorf("expected output to contain label names, got: %s", output)
	}
}

func TestRunLabelList_Error(t *testing.T) {
	mockRepo := &MockLabelRepository{
		ListErr: fmt.Errorf("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to list labels") {
		t.Errorf("expected error to contain 'failed to list labels', got: %v", err)
	}
}

func TestRunLabelShow_WithMockDependencies(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelShow(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Work") {
		t.Errorf("expected output to contain label name, got: %s", output)
	}
}

func TestRunLabelShow_NotFound(t *testing.T) {
	mockRepo := &MockLabelRepository{
		GetByNameErr: fmt.Errorf("not found"),
		GetErr:       fmt.Errorf("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelShow(cmd, []string{"NonExistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "label not found") {
		t.Errorf("expected error to contain 'label not found', got: %v", err)
	}
}

func TestRunLabelCreate_WithMockDependencies(t *testing.T) {
	mockRepo := &MockLabelRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	origText := labelTextColor
	formatFlag = "plain"
	labelBackgroundColor = ""
	labelTextColor = ""
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelCreate(cmd, []string{"NewLabel"})
	if err != nil {
		t.Fatalf("runLabelCreate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label created successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunLabelCreate_WithColors(t *testing.T) {
	mockRepo := &MockLabelRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	origText := labelTextColor
	formatFlag = "plain"
	labelBackgroundColor = "#ff0000"
	labelTextColor = "#ffffff"
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelCreate(cmd, []string{"ColoredLabel"})
	if err != nil {
		t.Fatalf("runLabelCreate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label created successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunLabelCreate_Error(t *testing.T) {
	mockRepo := &MockLabelRepository{
		CreateErr: fmt.Errorf("create failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelCreate(cmd, []string{"NewLabel"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to create label") {
		t.Errorf("expected error to contain 'failed to create label', got: %v", err)
	}
}

func TestRunLabelUpdate_WithMockDependencies(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	origText := labelTextColor
	formatFlag = "plain"
	labelBackgroundColor = "#0000ff"
	labelTextColor = "#ffffff"
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelUpdate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunLabelUpdate_SystemLabel(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "INBOX",
		Name: "INBOX",
		Type: mail.LabelTypeSystem,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origBg := labelBackgroundColor
	labelBackgroundColor = "#ff0000"
	defer func() { labelBackgroundColor = origBg }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"INBOX"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "cannot modify system label") {
		t.Errorf("expected error to contain 'cannot modify system label', got: %v", err)
	}
}

func TestRunLabelDelete_WithMockDependencies(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelDelete(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelDelete failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Work") || !contains(output, "deleted") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunLabelDelete_SystemLabel(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "INBOX",
		Name: "INBOX",
		Type: mail.LabelTypeSystem,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelDelete(cmd, []string{"INBOX"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "cannot delete system label") {
		t.Errorf("expected error to contain 'cannot delete system label', got: %v", err)
	}
}

func TestRunLabelDelete_QuietMode(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelDelete(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelDelete failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("quiet mode should not produce output, got: %s", output)
	}
}

func TestRunLabelList_JSONFormat(t *testing.T) {
	mockLabels := []*mail.Label{
		{ID: "label1", Name: "Work", Type: mail.LabelTypeUser},
	}

	mockRepo := &MockLabelRepository{
		Labels: mockLabels,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelList(cmd, []string{})
	if err != nil {
		t.Fatalf("runLabelList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Work") {
		t.Errorf("expected JSON output to contain label name, got: %s", output)
	}
}

func TestRunLabelShow_JSONFormat(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelShow(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Work") {
		t.Errorf("expected JSON output to contain label name, got: %s", output)
	}
}

func TestRunLabelCreate_JSONFormat(t *testing.T) {
	mockRepo := &MockLabelRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	origText := labelTextColor
	formatFlag = "json"
	labelBackgroundColor = ""
	labelTextColor = ""
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelCreate(cmd, []string{"NewLabel"})
	if err != nil {
		t.Fatalf("runLabelCreate failed: %v", err)
	}

	output := buf.String()
	// JSON output should contain label data
	if len(output) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestRunLabelUpdate_JSONFormat(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	formatFlag = "json"
	labelBackgroundColor = "#0000ff"
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelUpdate failed: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestRunLabelUpdate_UpdateError(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label:     mockLabel,
		UpdateErr: fmt.Errorf("update failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origBg := labelBackgroundColor
	labelBackgroundColor = "#ff0000"
	defer func() { labelBackgroundColor = origBg }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"Work"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to update label") {
		t.Errorf("expected error to contain 'failed to update label', got: %v", err)
	}
}

func TestRunLabelUpdate_NotFound(t *testing.T) {
	mockRepo := &MockLabelRepository{
		GetByNameErr: fmt.Errorf("not found"),
		GetErr:       fmt.Errorf("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origBg := labelBackgroundColor
	labelBackgroundColor = "#ff0000"
	defer func() { labelBackgroundColor = origBg }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"NonExistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "label not found") {
		t.Errorf("expected error to contain 'label not found', got: %v", err)
	}
}

func TestRunLabelDelete_NotFound(t *testing.T) {
	mockRepo := &MockLabelRepository{
		GetByNameErr: fmt.Errorf("not found"),
		GetErr:       fmt.Errorf("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelDelete(cmd, []string{"NonExistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "label not found") {
		t.Errorf("expected error to contain 'label not found', got: %v", err)
	}
}

func TestRunLabelDelete_DeleteError(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		Label:     mockLabel,
		DeleteErr: fmt.Errorf("delete failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelDelete(cmd, []string{"Work"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to delete label") {
		t.Errorf("expected error to contain 'failed to delete label', got: %v", err)
	}
}

func TestRunLabelCreate_WithPartialColors(t *testing.T) {
	mockRepo := &MockLabelRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	t.Run("background only", func(t *testing.T) {
		origFormat := formatFlag
		origBg := labelBackgroundColor
		origText := labelTextColor
		formatFlag = "plain"
		labelBackgroundColor = "#ff0000"
		labelTextColor = ""
		defer func() {
			formatFlag = origFormat
			labelBackgroundColor = origBg
			labelTextColor = origText
		}()

		cmd := &cobra.Command{Use: "test"}
		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := runLabelCreate(cmd, []string{"TestLabel"})
		if err != nil {
			t.Fatalf("runLabelCreate failed: %v", err)
		}

		output := buf.String()
		if !contains(output, "Label created successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})

	t.Run("text only", func(t *testing.T) {
		origFormat := formatFlag
		origBg := labelBackgroundColor
		origText := labelTextColor
		formatFlag = "plain"
		labelBackgroundColor = ""
		labelTextColor = "#ffffff"
		defer func() {
			formatFlag = origFormat
			labelBackgroundColor = origBg
			labelTextColor = origText
		}()

		cmd := &cobra.Command{Use: "test"}
		var buf bytes.Buffer
		cmd.SetOut(&buf)

		err := runLabelCreate(cmd, []string{"TestLabel2"})
		if err != nil {
			t.Fatalf("runLabelCreate failed: %v", err)
		}

		output := buf.String()
		if !contains(output, "Label created successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})
}

func TestRunLabelUpdate_WithColorPreservation(t *testing.T) {
	mockLabel := &mail.Label{
		ID:    "label123",
		Name:  "Work",
		Type:  mail.LabelTypeUser,
		Color: &mail.LabelColor{Background: "#0000ff", Text: "#ffffff"},
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	origText := labelTextColor
	formatFlag = "plain"
	labelBackgroundColor = "#ff0000"
	labelTextColor = ""
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelUpdate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunLabelShow_ByID(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		GetByNameErr: fmt.Errorf("not found by name"),
		Label:        mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelShow(cmd, []string{"label123"})
	if err != nil {
		t.Fatalf("runLabelShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Work") {
		t.Errorf("expected output to contain label name, got: %s", output)
	}
}

func TestRunLabelDelete_ByID(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		GetByNameErr: fmt.Errorf("not found by name"),
		Label:        mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
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

	err := runLabelDelete(cmd, []string{"label123"})
	if err != nil {
		t.Fatalf("runLabelDelete failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "deleted") {
		t.Errorf("expected confirmation message, got: %s", output)
	}
}

func TestRunLabelUpdate_ByID(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "Work",
		Type: mail.LabelTypeUser,
	}

	mockRepo := &MockLabelRepository{
		GetByNameErr: fmt.Errorf("not found by name"),
		Label:        mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	formatFlag = "plain"
	labelBackgroundColor = "#0000ff"
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"label123"})
	if err != nil {
		t.Fatalf("runLabelUpdate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRunLabelList_RepositoryError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelErr: fmt.Errorf("failed to create repository"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelList(cmd, []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunLabelShow_RepositoryError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelErr: fmt.Errorf("failed to create repository"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelShow(cmd, []string{"Work"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunLabelCreate_RepositoryError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelErr: fmt.Errorf("failed to create repository"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelCreate(cmd, []string{"NewLabel"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunLabelUpdate_RepositoryError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelErr: fmt.Errorf("failed to create repository"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"Work"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunLabelDelete_RepositoryError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelErr: fmt.Errorf("failed to create repository"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelDelete(cmd, []string{"Work"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunLabelUpdate_WithNoColorAndTextOnly(t *testing.T) {
	mockLabel := &mail.Label{
		ID:    "label123",
		Name:  "Work",
		Type:  mail.LabelTypeUser,
		Color: nil,
	}

	mockRepo := &MockLabelRepository{
		Label: mockLabel,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	origBg := labelBackgroundColor
	origText := labelTextColor
	formatFlag = "plain"
	labelBackgroundColor = ""
	labelTextColor = "#ffffff"
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"Work"})
	if err != nil {
		t.Fatalf("runLabelUpdate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}
