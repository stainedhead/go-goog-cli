// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// TestRunDraftCreate_WithNilMessage tests draft creation output when created draft has nil message
func TestRunDraftCreate_WithNilMessage(t *testing.T) {
	mockRepo := &MockDraftRepository{
		CreateResult: &mail.Draft{
			ID:      "draft-123",
			Message: nil,
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
	if !contains(output, "draft-123") {
		t.Errorf("expected draft ID in output, got: %s", output)
	}
}

// TestRunDraftCreate_EmptyBody tests draft creation with empty body
func TestRunDraftCreate_EmptyBody(t *testing.T) {
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
	draftBody = ""
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

// TestRunDraftShow_WithNilMessage tests draft show when draft has nil message
func TestRunDraftShow_WithNilMessage(t *testing.T) {
	mockDraft := &mail.Draft{
		ID:      "draft123",
		Message: nil,
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
	if contains(output, "--- Body ---") {
		t.Errorf("should not show body separator for nil message, got: %s", output)
	}
}

// TestRunDraftList_EmptyResults tests draft list with no drafts
func TestRunDraftList_EmptyResults(t *testing.T) {
	mockRepo := &MockDraftRepository{
		Drafts: []*mail.Draft{},
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
	// Should complete without error even with no drafts
	if len(output) == 0 {
		t.Error("expected some output even with empty results")
	}
}

// TestRunDraftSend_EmptyRecipients tests draft send with no recipients
func TestRunDraftSend_EmptyRecipients(t *testing.T) {
	mockRepo := &MockDraftRepository{
		SendResult: &mail.Message{
			ID:      "sent-msg-id",
			Subject: "Sent Subject",
			To:      []string{},
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

// TestRunDraftUpdate_EmptyFlags tests draft update with no flags set
func TestRunDraftUpdate_EmptyFlags(t *testing.T) {
	mockDraft := &mail.Draft{
		ID: "draft123",
		Message: &mail.Message{
			Subject: "Original Subject",
			To:      []string{"original@example.com"},
			Body:    "Original body",
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
	origTo := draftTo
	origSubject := draftSubject
	origBody := draftBody
	formatFlag = "plain"
	draftTo = []string{}
	draftSubject = ""
	draftBody = ""
	defer func() {
		formatFlag = origFormat
		draftTo = origTo
		draftSubject = origSubject
		draftBody = origBody
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

// TestRunDraftList_WithLargeLimit tests draft list with very large limit
func TestRunDraftList_WithLargeLimit(t *testing.T) {
	mockDrafts := []*mail.Draft{
		{ID: "draft1", Message: &mail.Message{Subject: "Draft 1"}},
		{ID: "draft2", Message: &mail.Message{Subject: "Draft 2"}},
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
	draftLimit = 1000
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

// TestRunDraftCreate_LongSubject tests draft creation with very long subject
func TestRunDraftCreate_LongSubject(t *testing.T) {
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

	longSubject := "This is a very long subject line that exceeds typical email subject length limits and tests how the system handles extremely long subject lines in draft creation and display operations"

	origFormat := formatFlag
	origTo := draftTo
	origSubject := draftSubject
	origBody := draftBody
	formatFlag = "plain"
	draftTo = []string{"user@example.com"}
	draftSubject = longSubject
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

// TestRunDraftUpdate_PartialFieldUpdate tests updating only one field at a time
func TestRunDraftUpdate_PartialFieldUpdate(t *testing.T) {
	tests := []struct {
		name            string
		updateTo        []string
		updateSubject   string
		updateBody      string
		expectedSuccess bool
	}{
		{
			name:            "update to only",
			updateTo:        []string{"new@example.com"},
			updateSubject:   "",
			updateBody:      "",
			expectedSuccess: true,
		},
		{
			name:            "update subject only",
			updateTo:        []string{},
			updateSubject:   "New Subject Only",
			updateBody:      "",
			expectedSuccess: true,
		},
		{
			name:            "update body only",
			updateTo:        []string{},
			updateSubject:   "",
			updateBody:      "New body only",
			expectedSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDraft := &mail.Draft{
				ID: "draft123",
				Message: &mail.Message{
					Subject: "Original Subject",
					To:      []string{"original@example.com"},
					Body:    "Original body",
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
			origTo := draftTo
			origSubject := draftSubject
			origBody := draftBody
			formatFlag = "plain"
			draftTo = tt.updateTo
			draftSubject = tt.updateSubject
			draftBody = tt.updateBody
			defer func() {
				formatFlag = origFormat
				draftTo = origTo
				draftSubject = origSubject
				draftBody = origBody
			}()

			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runDraftUpdate(cmd, []string{"draft123"})
			if tt.expectedSuccess {
				if err != nil {
					t.Errorf("expected success, got error: %v", err)
				}
				output := buf.String()
				if !contains(output, "Draft updated successfully") {
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

// TestRunDraftShow_JSONFormat tests draft show with JSON output
func TestRunDraftShow_JSONFormat(t *testing.T) {
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
	formatFlag = "json"
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
	// JSON format should still show body in separate section if not quiet
	if !contains(output, "--- Body ---") {
		t.Errorf("expected body separator even in JSON format when not quiet, got: %s", output)
	}
}

// TestDraftCreateCmd_ValidationEdgeCases tests PreRunE validation edge cases
func TestDraftCreateCmd_ValidationEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		to        []string
		subject   string
		expectErr string
	}{
		{
			name:      "whitespace-only subject",
			to:        []string{"user@example.com"},
			subject:   "   ",
			expectErr: "",
		},
		{
			name:      "single recipient",
			to:        []string{"single@example.com"},
			subject:   "Subject",
			expectErr: "",
		},
		{
			name:      "many recipients",
			to:        []string{"user1@example.com", "user2@example.com", "user3@example.com", "user4@example.com"},
			subject:   "Subject",
			expectErr: "",
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

			if tt.expectErr != "" {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !contains(err.Error(), tt.expectErr) {
					t.Errorf("expected error containing '%s', got: %v", tt.expectErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestRunDraftList_TableFormat tests draft list with table output format
func TestRunDraftList_TableFormat(t *testing.T) {
	mockDrafts := []*mail.Draft{
		{ID: "draft1", Message: &mail.Message{Subject: "Table Test Draft 1", To: []string{"user@example.com"}}},
		{ID: "draft2", Message: &mail.Message{Subject: "Table Test Draft 2", To: []string{"user2@example.com"}}},
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
	formatFlag = "table"
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

// TestRunDraftCreate_WithCreatedDraftHavingColor tests draft creation when result includes color info
func TestRunDraftCreate_WithCreatedDraftHavingColor(t *testing.T) {
	mockRepo := &MockDraftRepository{
		CreateResult: &mail.Draft{
			ID: "draft-with-color",
			Message: &mail.Message{
				To:      []string{"test@example.com"},
				Subject: "Colored Draft",
				Body:    "Test body",
			},
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
	origTo := draftTo
	origSubject := draftSubject
	origBody := draftBody
	formatFlag = "plain"
	draftTo = []string{"test@example.com"}
	draftSubject = "Test"
	draftBody = "Body"
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
	if !contains(output, "Colored Draft") {
		t.Errorf("expected subject in output, got: %s", output)
	}
}

// TestRunDraftUpdate_JSONOutputFormat tests draft update with JSON output
func TestRunDraftUpdate_JSONOutputFormat(t *testing.T) {
	mockDraft := &mail.Draft{
		ID: "draft123",
		Message: &mail.Message{
			Subject: "Original Subject",
			To:      []string{"original@example.com"},
		},
	}

	mockRepo := &MockDraftRepository{
		Draft: mockDraft,
		UpdateResult: &mail.Draft{
			ID: "draft123",
			Message: &mail.Message{
				Subject: "Updated Subject",
				To:      []string{"original@example.com"},
			},
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
	origSubject := draftSubject
	formatFlag = "json"
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
	if len(output) == 0 {
		t.Error("expected non-empty JSON output")
	}
}
