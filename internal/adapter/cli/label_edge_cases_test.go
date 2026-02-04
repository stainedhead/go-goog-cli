// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// TestRunLabelCreate_WithNilColor tests label creation output when no color is set
func TestRunLabelCreate_WithNilColor(t *testing.T) {
	mockRepo := &MockLabelRepository{
		CreateResult: &mail.Label{
			ID:    "label-123",
			Name:  "TestLabel",
			Type:  mail.LabelTypeUser,
			Color: nil,
		},
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
	if contains(output, "Background:") {
		t.Errorf("should not show color info when nil, got: %s", output)
	}
}

// TestRunLabelCreate_WithColorInResult tests label creation output when result includes color
func TestRunLabelCreate_WithColorInResult(t *testing.T) {
	mockRepo := &MockLabelRepository{
		CreateResult: &mail.Label{
			ID:   "label-123",
			Name: "ColoredLabel",
			Type: mail.LabelTypeUser,
			Color: &mail.LabelColor{
				Background: "#ff0000",
				Text:       "#ffffff",
			},
		},
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
	if !contains(output, "Background:") {
		t.Errorf("expected background color in output, got: %s", output)
	}
	if !contains(output, "#ff0000") {
		t.Errorf("expected background color value in output, got: %s", output)
	}
}

// TestRunLabelList_EmptyResults tests label list with no labels
func TestRunLabelList_EmptyResults(t *testing.T) {
	mockRepo := &MockLabelRepository{
		Labels: []*mail.Label{},
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
	// Should complete without error even with no labels
	if len(output) == 0 {
		t.Error("expected some output even with empty results")
	}
}

// TestRunLabelShow_SystemLabel tests showing a system label
func TestRunLabelShow_SystemLabel(t *testing.T) {
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

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelShow(cmd, []string{"INBOX"})
	if err != nil {
		t.Fatalf("runLabelShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "INBOX") {
		t.Errorf("expected output to contain label name, got: %s", output)
	}
}

// TestRunLabelCreate_JSONOutputWithColors tests JSON output for label creation with colors
func TestRunLabelCreate_JSONOutputWithColors(t *testing.T) {
	mockRepo := &MockLabelRepository{
		CreateResult: &mail.Label{
			ID:   "label-123",
			Name: "JSONLabel",
			Type: mail.LabelTypeUser,
			Color: &mail.LabelColor{
				Background: "#00ff00",
				Text:       "#000000",
			},
		},
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
	formatFlag = "json"
	labelBackgroundColor = "#00ff00"
	labelTextColor = "#000000"
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelCreate(cmd, []string{"JSONLabel"})
	if err != nil {
		t.Fatalf("runLabelCreate failed: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

// TestRunLabelUpdate_WithPartialColorUpdate tests updating only background or text color
func TestRunLabelUpdate_WithPartialColorUpdate(t *testing.T) {
	tests := []struct {
		name              string
		existingColor     *mail.LabelColor
		updateBg          string
		updateText        string
		expectedBg        string
		expectedText      string
		expectedToSucceed bool
	}{
		{
			name:              "update bg only with existing color",
			existingColor:     &mail.LabelColor{Background: "#000000", Text: "#ffffff"},
			updateBg:          "#ff0000",
			updateText:        "",
			expectedBg:        "#ff0000",
			expectedText:      "#ffffff",
			expectedToSucceed: true,
		},
		{
			name:              "update text only with existing color",
			existingColor:     &mail.LabelColor{Background: "#000000", Text: "#ffffff"},
			updateBg:          "",
			updateText:        "#00ff00",
			expectedBg:        "#000000",
			expectedText:      "#00ff00",
			expectedToSucceed: true,
		},
		{
			name:              "update bg only with no existing color",
			existingColor:     nil,
			updateBg:          "#ff0000",
			updateText:        "",
			expectedBg:        "#ff0000",
			expectedText:      "#ffffff",
			expectedToSucceed: true,
		},
		{
			name:              "update text only with no existing color",
			existingColor:     nil,
			updateBg:          "",
			updateText:        "#00ff00",
			expectedBg:        "#000000",
			expectedText:      "#00ff00",
			expectedToSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabel := &mail.Label{
				ID:    "label123",
				Name:  "TestLabel",
				Type:  mail.LabelTypeUser,
				Color: tt.existingColor,
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
			labelBackgroundColor = tt.updateBg
			labelTextColor = tt.updateText
			defer func() {
				formatFlag = origFormat
				labelBackgroundColor = origBg
				labelTextColor = origText
			}()

			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runLabelUpdate(cmd, []string{"TestLabel"})
			if tt.expectedToSucceed {
				if err != nil {
					t.Errorf("expected success, got error: %v", err)
				}
				output := buf.String()
				if !contains(output, "Label updated successfully") {
					t.Errorf("expected success message, got: %s", output)
				}
			} else {
				if err == nil {
					t.Error("expected error, got nil")
				}
			}
		})
	}
}

// TestRunLabelList_MixedLabelTypes tests listing both system and user labels
func TestRunLabelList_MixedLabelTypes(t *testing.T) {
	mockLabels := []*mail.Label{
		{ID: "INBOX", Name: "INBOX", Type: mail.LabelTypeSystem},
		{ID: "SENT", Name: "SENT", Type: mail.LabelTypeSystem},
		{ID: "label1", Name: "Work", Type: mail.LabelTypeUser},
		{ID: "label2", Name: "Personal", Type: mail.LabelTypeUser},
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
	if !contains(output, "INBOX") {
		t.Errorf("expected output to contain system label INBOX, got: %s", output)
	}
	if !contains(output, "Work") {
		t.Errorf("expected output to contain user label Work, got: %s", output)
	}
}

// TestRunLabelList_TableFormat tests label list with table output format
func TestRunLabelList_TableFormat(t *testing.T) {
	mockLabels := []*mail.Label{
		{ID: "label1", Name: "TableLabel1", Type: mail.LabelTypeUser},
		{ID: "label2", Name: "TableLabel2", Type: mail.LabelTypeUser},
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
	formatFlag = "table"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelList(cmd, []string{})
	if err != nil {
		t.Fatalf("runLabelList failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "TableLabel1") || !contains(output, "TableLabel2") {
		t.Errorf("expected output to contain label names, got: %s", output)
	}
}

// TestRunLabelShow_TableFormat tests label show with table output format
func TestRunLabelShow_TableFormat(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "TableShowLabel",
		Type: mail.LabelTypeUser,
		Color: &mail.LabelColor{
			Background: "#4285f4",
			Text:       "#ffffff",
		},
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
	formatFlag = "table"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelShow(cmd, []string{"TableShowLabel"})
	if err != nil {
		t.Fatalf("runLabelShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "TableShowLabel") {
		t.Errorf("expected output to contain label name, got: %s", output)
	}
}

// TestRunLabelCreate_SpecialCharactersInName tests label creation with special characters
func TestRunLabelCreate_SpecialCharactersInName(t *testing.T) {
	tests := []struct {
		name      string
		labelName string
	}{
		{
			name:      "label with spaces",
			labelName: "My Work Label",
		},
		{
			name:      "label with slash",
			labelName: "Work/Projects",
		},
		{
			name:      "label with dash",
			labelName: "High-Priority",
		},
		{
			name:      "label with underscore",
			labelName: "My_Label",
		},
		{
			name:      "label with numbers",
			labelName: "Label123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockLabelRepository{
				CreateResult: &mail.Label{
					ID:   "label-special",
					Name: tt.labelName,
					Type: mail.LabelTypeUser,
				},
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
			labelTextColor = ""
			defer func() {
				formatFlag = origFormat
				labelBackgroundColor = origBg
				labelTextColor = origText
			}()

			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runLabelCreate(cmd, []string{tt.labelName})
			if err != nil {
				t.Fatalf("runLabelCreate failed: %v", err)
			}

			output := buf.String()
			if !contains(output, "Label created successfully") {
				t.Errorf("expected success message, got: %s", output)
			}
		})
	}
}

// TestRunLabelUpdate_EmptyUpdate tests label update with no color flags set
func TestRunLabelUpdate_EmptyUpdate(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "TestLabel",
		Type: mail.LabelTypeUser,
		Color: &mail.LabelColor{
			Background: "#000000",
			Text:       "#ffffff",
		},
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
	labelTextColor = ""
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
		labelTextColor = origText
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"TestLabel"})
	if err != nil {
		t.Fatalf("runLabelUpdate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

// TestRunLabelShow_WithColorInfo tests showing a label with color information
func TestRunLabelShow_WithColorInfo(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "ColoredLabel",
		Type: mail.LabelTypeUser,
		Color: &mail.LabelColor{
			Background: "#4285f4",
			Text:       "#ffffff",
		},
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

	err := runLabelShow(cmd, []string{"ColoredLabel"})
	if err != nil {
		t.Fatalf("runLabelShow failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "ColoredLabel") {
		t.Errorf("expected output to contain label name, got: %s", output)
	}
}

// TestLabelDeleteCmd_WithoutConfirm tests that delete command requires --confirm flag
func TestLabelDeleteCmd_WithoutConfirm(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "TestLabel",
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

	origConfirm := labelConfirm
	labelConfirm = false
	defer func() { labelConfirm = origConfirm }()

	mockCmd := &cobra.Command{Use: "test"}

	err := labelDeleteCmd.PreRunE(mockCmd, []string{"TestLabel"})
	if err == nil {
		t.Error("expected error for missing --confirm flag, got nil")
	}
	if !contains(err.Error(), "confirm") {
		t.Errorf("expected error to mention confirm flag, got: %v", err)
	}
}

// TestRunLabelUpdate_TableFormat tests label update with table output format
func TestRunLabelUpdate_TableFormat(t *testing.T) {
	mockLabel := &mail.Label{
		ID:   "label123",
		Name: "TableUpdateLabel",
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
	formatFlag = "table"
	labelBackgroundColor = "#0000ff"
	defer func() {
		formatFlag = origFormat
		labelBackgroundColor = origBg
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runLabelUpdate(cmd, []string{"TableUpdateLabel"})
	if err != nil {
		t.Fatalf("runLabelUpdate failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "Label updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}

// TestRunLabelCreate_DefaultColorsBehavior tests default color assignment
func TestRunLabelCreate_DefaultColorsBehavior(t *testing.T) {
	tests := []struct {
		name       string
		setBg      string
		setText    string
		expectBoth bool
	}{
		{
			name:       "only background set, text should default",
			setBg:      "#ff0000",
			setText:    "",
			expectBoth: true,
		},
		{
			name:       "only text set, background should default",
			setBg:      "",
			setText:    "#00ff00",
			expectBoth: true,
		},
		{
			name:       "neither set, no color",
			setBg:      "",
			setText:    "",
			expectBoth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockLabelRepository{
				CreateResult: &mail.Label{
					ID:   "label-default",
					Name: "DefaultColorLabel",
					Type: mail.LabelTypeUser,
				},
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
			labelBackgroundColor = tt.setBg
			labelTextColor = tt.setText
			defer func() {
				formatFlag = origFormat
				labelBackgroundColor = origBg
				labelTextColor = origText
			}()

			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runLabelCreate(cmd, []string{"DefaultColorLabel"})
			if err != nil {
				t.Fatalf("runLabelCreate failed: %v", err)
			}

			output := buf.String()
			if !contains(output, "Label created successfully") {
				t.Errorf("expected success message, got: %s", output)
			}
		})
	}
}
