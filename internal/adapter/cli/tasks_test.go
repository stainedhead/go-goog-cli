package cli

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/spf13/cobra"
	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// ============================================================================
// Help Tests
// ============================================================================

func TestTasksCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(tasksCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"tasks", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "tasks") {
		t.Error("expected output to contain 'tasks'")
	}
	if !contains(output, "lists") {
		t.Error("expected output to contain 'lists'")
	}
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
}

func TestTasksListsCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(tasksCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"tasks", "lists", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "lists") {
		t.Error("expected output to contain 'lists'")
	}
}

func TestTasksCreateListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(tasksCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"tasks", "create-list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create-list") {
		t.Error("expected output to contain 'create-list'")
	}
}

func TestTasksListCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(tasksCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"tasks", "list", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "list") {
		t.Error("expected output to contain 'list'")
	}
}

func TestTasksCreateCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(tasksCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"tasks", "create", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "create") {
		t.Error("expected output to contain 'create'")
	}
	if !contains(output, "--notes") {
		t.Error("expected output to contain '--notes'")
	}
	if !contains(output, "--due") {
		t.Error("expected output to contain '--due'")
	}
}

func TestTasksCompleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(tasksCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"tasks", "complete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "complete") {
		t.Error("expected output to contain 'complete'")
	}
}

// ============================================================================
// Args Validation Tests
// ============================================================================

func TestTasksCreateListCmd_ArgsValidation(t *testing.T) {
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
			args:      []string{"Work Tasks"},
			expectErr: false,
		},
		{
			name:      "too many args",
			args:      []string{"Work", "Tasks"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tasksCreateListCmd.Args(tasksCreateListCmd, tt.args)
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

func TestTasksDeleteListCmd_RequiresConfirm(t *testing.T) {
	// Save original flag value
	origConfirm := tasksDeleteConfirm
	defer func() { tasksDeleteConfirm = origConfirm }()

	tasksDeleteConfirm = false

	err := tasksDeleteListCmd.PreRunE(tasksDeleteListCmd, []string{"list123"})
	if err == nil {
		t.Error("expected error when --confirm is not set")
	}

	tasksDeleteConfirm = true
	err = tasksDeleteListCmd.PreRunE(tasksDeleteListCmd, []string{"list123"})
	if err != nil {
		t.Errorf("unexpected error with --confirm set: %v", err)
	}
}

func TestTasksUpdateListCmd_RequiresTitle(t *testing.T) {
	// Save original flag value
	origTitle := tasksTitle
	defer func() { tasksTitle = origTitle }()

	tasksTitle = ""
	err := tasksUpdateListCmd.PreRunE(tasksUpdateListCmd, []string{"list123"})
	if err == nil {
		t.Error("expected error when --title is not set")
	}

	tasksTitle = "New Title"
	err = tasksUpdateListCmd.PreRunE(tasksUpdateListCmd, []string{"list123"})
	if err != nil {
		t.Errorf("unexpected error with --title set: %v", err)
	}
}

func TestTasksUpdateCmd_RequiresAtLeastOneFlag(t *testing.T) {
	// Save original flag values
	origTitle := tasksTitle
	origNotes := tasksNotes
	origDue := tasksDue
	defer func() {
		tasksTitle = origTitle
		tasksNotes = origNotes
		tasksDue = origDue
	}()

	tasksTitle = ""
	tasksNotes = ""
	tasksDue = ""
	err := tasksUpdateCmd.PreRunE(tasksUpdateCmd, []string{"task123"})
	if err == nil {
		t.Error("expected error when no flags are set")
	}

	tasksTitle = "New Title"
	err = tasksUpdateCmd.PreRunE(tasksUpdateCmd, []string{"task123"})
	if err != nil {
		t.Errorf("unexpected error with --title set: %v", err)
	}
}

// ============================================================================
// Flag Tests
// ============================================================================

func TestTasksCmd_HasListFlag(t *testing.T) {
	flag := tasksCmd.PersistentFlags().Lookup("list")
	if flag == nil {
		t.Error("expected --list flag to be set")
	}
	if flag.DefValue != "@default" {
		t.Errorf("expected default value '@default', got %s", flag.DefValue)
	}
}

func TestTasksListCmd_HasShowCompletedFlag(t *testing.T) {
	flag := tasksListCmd.Flags().Lookup("show-completed")
	if flag == nil {
		t.Error("expected --show-completed flag to be set")
	}
}

func TestTasksCreateCmd_HasNotesFlag(t *testing.T) {
	flag := tasksCreateCmd.Flags().Lookup("notes")
	if flag == nil {
		t.Error("expected --notes flag to be set")
	}
}

func TestTasksCreateCmd_HasDueFlag(t *testing.T) {
	flag := tasksCreateCmd.Flags().Lookup("due")
	if flag == nil {
		t.Error("expected --due flag to be set")
	}
}

func TestTasksCreateCmd_HasParentFlag(t *testing.T) {
	flag := tasksCreateCmd.Flags().Lookup("parent")
	if flag == nil {
		t.Error("expected --parent flag to be set")
	}
}

func TestTasksDeleteCmd_HasConfirmFlag(t *testing.T) {
	flag := tasksDeleteCmd.Flags().Lookup("confirm")
	if flag == nil {
		t.Error("expected --confirm flag to be set")
	}
}

// ============================================================================
// Command Execution Tests with Mocks
// ============================================================================

func TestRunTasksLists_Success(t *testing.T) {
	mockLists := []*domaintasks.TaskList{
		{ID: "list1", Title: "My Tasks", Updated: time.Now()},
		{ID: "list2", Title: "Work Tasks", Updated: time.Now()},
	}

	mockRepo := &MockTaskListRepository{
		Lists: mockLists,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
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

	err := runTasksLists(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "My Tasks") {
		t.Error("expected output to contain 'My Tasks'")
	}
	if !contains(output, "Work Tasks") {
		t.Error("expected output to contain 'Work Tasks'")
	}
}

func TestRunTasksLists_Error(t *testing.T) {
	mockRepo := &MockTaskListRepository{
		ListErr: errors.New("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksLists(cmd, []string{})

	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestRunTasksCreateList_Success(t *testing.T) {
	mockList := &domaintasks.TaskList{
		ID:      "newlist",
		Title:   "New List",
		Updated: time.Now(),
	}

	mockRepo := &MockTaskListRepository{
		TaskList: mockList,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
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

	err := runTasksCreateList(cmd, []string{"New List"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "New List") {
		t.Error("expected output to contain 'New List'")
	}
}

func TestRunTasksList_Success(t *testing.T) {
	mockTasks := &domaintasks.ListResult[*domaintasks.Task]{
		Items: []*domaintasks.Task{
			{ID: "task1", Title: "Buy groceries", Status: "needsAction"},
			{ID: "task2", Title: "Write report", Status: "completed"},
		},
	}

	mockRepo := &MockTaskRepository{
		Tasks: mockTasks,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
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

	err := runTasksList(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Buy groceries") {
		t.Error("expected output to contain 'Buy groceries'")
	}
}

func TestRunTasksGet_Success(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:     "task1",
		Title:  "Test Task",
		Status: "needsAction",
		Notes:  "Some notes",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
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

	err := runTasksGet(cmd, []string{"task1"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "Test Task") {
		t.Error("expected output to contain 'Test Task'")
	}
}

func TestRunTasksCreate_Success(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:     "newtask",
		Title:  "New Task",
		Status: "needsAction",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag values
	origNotes := tasksNotes
	origDue := tasksDue
	origParent := tasksParent
	defer func() {
		tasksNotes = origNotes
		tasksDue = origDue
		tasksParent = origParent
	}()

	tasksNotes = ""
	tasksDue = ""
	tasksParent = ""

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksCreate(cmd, []string{"New Task"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "New Task") {
		t.Error("expected output to contain 'New Task'")
	}
}

func TestRunTasksComplete_Success(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:         "task1",
		Title:      "Test Task",
		Status:     "needsAction",
		TaskListID: "@default",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
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

	err := runTasksComplete(cmd, []string{"task1"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTasksReopen_Success(t *testing.T) {
	completed := time.Now()
	mockTask := &domaintasks.Task{
		ID:         "task1",
		Title:      "Test Task",
		Status:     "completed",
		Completed:  &completed,
		TaskListID: "@default",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
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

	err := runTasksReopen(cmd, []string{"task1"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTasksDelete_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
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

	err := runTasksDelete(cmd, []string{"task1"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "deleted") {
		t.Error("expected output to contain 'deleted'")
	}
}

func TestRunTasksMove_Success(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:     "task1",
		Title:  "Test Task",
		Status: "needsAction",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag values
	origParent := tasksParent
	origPrevious := tasksPrevious
	defer func() {
		tasksParent = origParent
		tasksPrevious = origPrevious
	}()

	tasksParent = "parent123"
	tasksPrevious = ""

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksMove(cmd, []string{"task1"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTasksClear_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
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

	err := runTasksClear(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "cleared") {
		t.Error("expected output to contain 'cleared'")
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestRunTasksCreate_InvalidDueDate(t *testing.T) {
	mockRepo := &MockTaskRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	// Save and restore flag values
	origDue := tasksDue
	defer func() { tasksDue = origDue }()

	tasksDue = "invalid-date"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksCreate(cmd, []string{"New Task"})

	if err == nil {
		t.Error("expected error for invalid due date")
	}
}

func TestRunTasksUpdate_InvalidDueDate(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:         "task1",
		Title:      "Test Task",
		TaskListID: "@default",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	// Save and restore flag values
	origDue := tasksDue
	defer func() { tasksDue = origDue }()

	tasksDue = "invalid-date"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdate(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error for invalid due date")
	}
}

func TestRunTasksDeleteList_Success(t *testing.T) {
	mockRepo := &MockTaskListRepository{}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
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

	err := runTasksDeleteList(cmd, []string{"list123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "deleted") {
		t.Error("expected output to contain 'deleted'")
	}
}

func TestRunTasksUpdateList_Success(t *testing.T) {
	mockList := &domaintasks.TaskList{
		ID:      "list123",
		Title:   "Old Title",
		Updated: time.Now(),
	}

	mockRepo := &MockTaskListRepository{
		TaskList: mockList,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag value
	origTitle := tasksTitle
	defer func() { tasksTitle = origTitle }()

	tasksTitle = "New Title"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdateList(cmd, []string{"list123"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTasksUpdate_Success(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:         "task1",
		Title:      "Old Title",
		TaskListID: "@default",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag values
	origTitle := tasksTitle
	origNotes := tasksNotes
	origDue := tasksDue
	defer func() {
		tasksTitle = origTitle
		tasksNotes = origNotes
		tasksDue = origDue
	}()

	tasksTitle = "New Title"
	tasksNotes = "New notes"
	tasksDue = ""

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdate(cmd, []string{"task1"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "New Title") {
		t.Error("expected output to contain 'New Title'")
	}
}

func TestRunTasksCreate_WithDueDate(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:     "newtask",
		Title:  "New Task",
		Status: "needsAction",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag values
	origNotes := tasksNotes
	origDue := tasksDue
	origParent := tasksParent
	defer func() {
		tasksNotes = origNotes
		tasksDue = origDue
		tasksParent = origParent
	}()

	tasksNotes = "Task notes"
	tasksDue = "2024-12-31"
	tasksParent = ""

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksCreate(cmd, []string{"New Task"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTasksCreate_WithParent(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:     "newtask",
		Title:  "Subtask",
		Status: "needsAction",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag values
	origNotes := tasksNotes
	origDue := tasksDue
	origParent := tasksParent
	defer func() {
		tasksNotes = origNotes
		tasksDue = origDue
		tasksParent = origParent
	}()

	tasksNotes = ""
	tasksDue = ""
	tasksParent = "parent123"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksCreate(cmd, []string{"Subtask"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTasksList_WithFilters(t *testing.T) {
	mockTasks := &domaintasks.ListResult[*domaintasks.Task]{
		Items: []*domaintasks.Task{
			{ID: "task1", Title: "Task 1", Status: "needsAction"},
			{ID: "task2", Title: "Task 2", Status: "completed"},
		},
	}

	mockRepo := &MockTaskRepository{
		Tasks: mockTasks,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag values
	origShowCompleted := tasksShowCompleted
	origShowHidden := tasksShowHidden
	origMaxResults := tasksMaxResults
	defer func() {
		tasksShowCompleted = origShowCompleted
		tasksShowHidden = origShowHidden
		tasksMaxResults = origMaxResults
	}()

	tasksShowCompleted = true
	tasksShowHidden = false
	tasksMaxResults = 50

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksList(cmd, []string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTasksCreate_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		CreateErr: errors.New("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	// Save and restore flag values
	origNotes := tasksNotes
	origDue := tasksDue
	origParent := tasksParent
	defer func() {
		tasksNotes = origNotes
		tasksDue = origDue
		tasksParent = origParent
	}()

	tasksNotes = ""
	tasksDue = ""
	tasksParent = ""

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksCreate(cmd, []string{"New Task"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksUpdate_WithDueDate(t *testing.T) {
	mockTask := &domaintasks.Task{
		ID:         "task1",
		Title:      "Test Task",
		TaskListID: "@default",
	}

	mockRepo := &MockTaskRepository{
		Task: mockTask,
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origFormat := formatFlag
	formatFlag = "plain"
	defer func() { formatFlag = origFormat }()

	// Save and restore flag values
	origTitle := tasksTitle
	origNotes := tasksNotes
	origDue := tasksDue
	defer func() {
		tasksTitle = origTitle
		tasksNotes = origNotes
		tasksDue = origDue
	}()

	tasksTitle = ""
	tasksNotes = ""
	tasksDue = "2024-12-31"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdate(cmd, []string{"task1"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// ============================================================================
// Additional Coverage Tests
// ============================================================================

func TestRunTasksCreateList_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskListRepository{
		CreateErr: errors.New("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksCreateList(cmd, []string{"New List"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksDeleteList_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskListRepository{
		DeleteErr: errors.New("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksDeleteList(cmd, []string{"list123"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksUpdateList_GetError(t *testing.T) {
	mockRepo := &MockTaskListRepository{
		GetErr: errors.New("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := tasksTitle
	defer func() { tasksTitle = origTitle }()
	tasksTitle = "New Title"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdateList(cmd, []string{"list123"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksUpdateList_UpdateError(t *testing.T) {
	mockRepo := &MockTaskListRepository{
		TaskList:  &domaintasks.TaskList{ID: "list123", Title: "Old Title"},
		UpdateErr: errors.New("update failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskListRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := tasksTitle
	defer func() { tasksTitle = origTitle }()
	tasksTitle = "New Title"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdateList(cmd, []string{"list123"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksList_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		ListErr: errors.New("API error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksList(cmd, []string{})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksGet_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetErr: errors.New("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksGet(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksUpdate_GetError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetErr: errors.New("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := tasksTitle
	defer func() { tasksTitle = origTitle }()
	tasksTitle = "New Title"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdate(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksUpdate_UpdateError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		Task:      &domaintasks.Task{ID: "task1", Title: "Old", TaskListID: "@default"},
		UpdateErr: errors.New("update failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTitle := tasksTitle
	defer func() { tasksTitle = origTitle }()
	tasksTitle = "New Title"

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksUpdate(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksComplete_GetError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetErr: errors.New("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksComplete(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksComplete_UpdateError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		Task:      &domaintasks.Task{ID: "task1", Status: "needsAction", TaskListID: "@default"},
		UpdateErr: errors.New("update failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksComplete(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksReopen_GetError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetErr: errors.New("not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksReopen(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksReopen_UpdateError(t *testing.T) {
	completed := time.Now()
	mockRepo := &MockTaskRepository{
		Task:      &domaintasks.Task{ID: "task1", Status: "completed", Completed: &completed, TaskListID: "@default"},
		UpdateErr: errors.New("update failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksReopen(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksDelete_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteErr: errors.New("delete failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksDelete(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksMove_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		MoveErr: errors.New("move failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origParent := tasksParent
	origPrevious := tasksPrevious
	defer func() {
		tasksParent = origParent
		tasksPrevious = origPrevious
	}()
	tasksParent = "parent123"
	tasksPrevious = ""

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksMove(cmd, []string{"task1"})

	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestRunTasksClear_RepositoryError(t *testing.T) {
	mockRepo := &MockTaskRepository{
		ClearErr: errors.New("clear failed"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			TaskRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runTasksClear(cmd, []string{})

	if err == nil {
		t.Error("expected error from repository")
	}
}
