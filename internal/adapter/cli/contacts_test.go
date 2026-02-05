package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	domaincontacts "github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// ============================================================================
// Help Tests
// ============================================================================

func TestContactsCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "contacts") {
		t.Error("expected output to contain 'contacts'")
	}
}

func TestContactsListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
}

func TestContactsGetCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "get", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "get") {
		t.Error("expected output to contain 'get'")
	}
}

func TestContactsCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "create", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
}

func TestContactsUpdateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "update", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "update") {
		t.Error("expected output to contain 'update'")
	}
}

func TestContactsDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "delete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
}

func TestContactsSearchCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "search", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "search") {
		t.Error("expected output to contain 'search'")
	}
}

// ============================================================================
// Args Validation Tests
// ============================================================================

func TestContactsDeleteCmd_RequiresConfirm(t *testing.T) {
	origConfirm := contactsDeleteConfirm
	defer func() { contactsDeleteConfirm = origConfirm }()

	contactsDeleteConfirm = false
	err := contactsDeleteCmd.PreRunE(contactsDeleteCmd, []string{"people/c123"})
	if err == nil {
		t.Error("expected error when --confirm is not set")
	}

	contactsDeleteConfirm = true
	err = contactsDeleteCmd.PreRunE(contactsDeleteCmd, []string{"people/c123"})
	if err != nil {
		t.Errorf("unexpected error with --confirm set: %v", err)
	}
}

func TestContactsGetCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"people/c123"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"people/c123", "extra"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := contactsGetCmd.Args(contactsGetCmd, tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// ============================================================================
// Flag Tests
// ============================================================================

func TestContactsListCmd_HasMaxResultsFlag(t *testing.T) {
	flag := contactsListCmd.Flags().Lookup("max-results")
	if flag == nil {
		t.Error("expected --max-results flag to be set")
	}
}

func TestContactsCreateCmd_HasNameFlags(t *testing.T) {
	flag := contactsCreateCmd.Flags().Lookup("given-name")
	if flag == nil {
		t.Error("expected --given-name flag to be set")
	}

	flag = contactsCreateCmd.Flags().Lookup("family-name")
	if flag == nil {
		t.Error("expected --family-name flag to be set")
	}
}

func TestContactsCreateCmd_HasEmailFlag(t *testing.T) {
	flag := contactsCreateCmd.Flags().Lookup("email")
	if flag == nil {
		t.Error("expected --email flag to be set")
	}
}

func TestContactsDeleteCmd_HasConfirmFlag(t *testing.T) {
	flag := contactsDeleteCmd.Flags().Lookup("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}

// ============================================================================
// Command Execution Tests with Mocks
// ============================================================================

func TestRunContactsList_Success(t *testing.T) {
	mockContacts := &domaincontacts.ListResult[*domaincontacts.Contact]{
		Items: []*domaincontacts.Contact{
			{
				ResourceName: "people/c1",
				Names:        []domaincontacts.Name{{DisplayName: "John Doe"}},
			},
			{
				ResourceName: "people/c2",
				Names:        []domaincontacts.Name{{DisplayName: "Jane Smith"}},
			},
		},
	}

	mockRepo := &MockContactRepository{
		Contacts: mockContacts,
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsList(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "John Doe") {
		t.Error("expected output to contain 'John Doe'")
	}
	if !contains(output, "Jane Smith") {
		t.Error("expected output to contain 'Jane Smith'")
	}
}

func TestRunContactsGet_Success(t *testing.T) {
	mockContact := &domaincontacts.Contact{
		ResourceName: "people/c123",
		Names:        []domaincontacts.Name{{DisplayName: "John Doe", GivenName: "John", FamilyName: "Doe"}},
		EmailAddresses: []domaincontacts.Email{
			{Value: "john@example.com", Type: "work"},
		},
	}

	mockRepo := &MockContactRepository{
		Contact: mockContact,
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsGet(cmd, []string{"people/c123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "John Doe") {
		t.Error("expected output to contain 'John Doe'")
	}
}

func TestRunContactsCreate_Success(t *testing.T) {
	mockContact := &domaincontacts.Contact{
		ResourceName: "people/c123",
		Names:        []domaincontacts.Name{{GivenName: "John", FamilyName: "Doe"}},
	}

	mockRepo := &MockContactRepository{
		Contact: mockContact,
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	origGivenName := contactsGivenName
	origFamilyName := contactsFamilyName
	defer func() {
		contactsGivenName = origGivenName
		contactsFamilyName = origFamilyName
	}()

	contactsGivenName = "John"
	contactsFamilyName = "Doe"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsCreate(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunContactsDelete_Success(t *testing.T) {
	mockRepo := &MockContactRepository{}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsDelete(cmd, []string{"people/c123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "deleted") {
		t.Error("expected output to contain 'deleted'")
	}
}

func TestRunContactsSearch_Success(t *testing.T) {
	mockContacts := &domaincontacts.ListResult[*domaincontacts.Contact]{
		Items: []*domaincontacts.Contact{
			{
				ResourceName: "people/c1",
				Names:        []domaincontacts.Name{{DisplayName: "John Doe"}},
			},
		},
	}

	mockRepo := &MockContactRepository{
		SearchResult: mockContacts,
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsSearch(cmd, []string{"John"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "John Doe") {
		t.Error("expected output to contain 'John Doe'")
	}
}

// ============================================================================
// Group Command Tests
// ============================================================================

func TestContactsGroupsCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(contactsCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"contacts", "groups", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "groups") {
		t.Error("expected output to contain 'groups'")
	}
}

func TestRunContactsGroups_Success(t *testing.T) {
	mockGroups := []*domaincontacts.ContactGroup{
		{ResourceName: "contactGroups/g1", Name: "Family"},
		{ResourceName: "contactGroups/g2", Name: "Friends"},
	}

	mockGroupRepo := &MockContactGroupRepository{
		Groups: mockGroups,
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactGroupRepo = mockGroupRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsGroups(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Family") {
		t.Error("expected output to contain 'Family'")
	}
}

func TestRunContactsGroupCreate_Success(t *testing.T) {
	mockGroup := &domaincontacts.ContactGroup{
		ResourceName: "contactGroups/g123",
		Name:         "Work Contacts",
	}

	mockGroupRepo := &MockContactGroupRepository{
		Group: mockGroup,
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactGroupRepo = mockGroupRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsGroupCreate(cmd, []string{"Work Contacts"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Work Contacts") {
		t.Error("expected output to contain 'Work Contacts'")
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestRunContactsList_Error(t *testing.T) {
	mockRepo := &MockContactRepository{
		ListErr: errors.New("API error"),
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsList(cmd, []string{})

	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestRunContactsGet_Error(t *testing.T) {
	mockRepo := &MockContactRepository{
		GetErr: errors.New("not found"),
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsGet(cmd, []string{"people/c123"})

	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestRunContactsCreate_Error(t *testing.T) {
	mockRepo := &MockContactRepository{
		CreateErr: errors.New("create failed"),
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origGivenName := contactsGivenName
	defer func() { contactsGivenName = origGivenName }()
	contactsGivenName = "John"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsCreate(cmd, []string{})

	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestRunContactsDelete_Error(t *testing.T) {
	mockRepo := &MockContactRepository{
		DeleteErr: errors.New("delete failed"),
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsDelete(cmd, []string{"people/c123"})

	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestRunContactsSearch_Error(t *testing.T) {
	mockRepo := &MockContactRepository{
		SearchErr: errors.New("search failed"),
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactRepo = mockRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsSearch(cmd, []string{"query"})

	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestRunContactsGroups_Error(t *testing.T) {
	mockGroupRepo := &MockContactGroupRepository{
		ListErr: errors.New("API error"),
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactGroupRepo = mockGroupRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsGroups(cmd, []string{})

	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestRunContactsGroupCreate_Error(t *testing.T) {
	mockGroupRepo := &MockContactGroupRepository{
		CreateErr: errors.New("create failed"),
	}

	mockFactory := &MockRepositoryFactory{}
	mockFactory.ContactGroupRepo = mockGroupRepo

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: mockFactory,
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runContactsGroupCreate(cmd, []string{"Work"})

	if err == nil {
		t.Error("expected error but got none")
	}
}
