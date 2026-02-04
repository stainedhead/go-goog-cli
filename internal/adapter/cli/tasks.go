package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
)

// Command flags for tasks actions.
var (
	tasksListID        string
	tasksShowCompleted bool
	tasksShowHidden    bool
	tasksDueMax        string
	tasksDueMin        string
	tasksMaxResults    int64
	tasksTitle         string
	tasksNotes         string
	tasksDue           string
	tasksParent        string
	tasksPrevious      string
	tasksDeleteConfirm bool
	tasksClearConfirm  bool
)

// tasksCmd represents the tasks command group.
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage Google Tasks",
	Long: `Manage Google Tasks.

The tasks commands allow you to list, create, update, and manage
task lists and tasks in your Google Tasks account.`,
}

// ================ Task List Commands ================

// tasksListsCmd lists all task lists.
var tasksListsCmd = &cobra.Command{
	Use:     "lists",
	Aliases: []string{"list-lists"},
	Short:   "List all task lists",
	Long: `List all task lists in your Google Tasks account.

Task lists are containers for tasks. You can have multiple
task lists to organize different types of tasks.`,
	Example: `  # List all task lists
  goog tasks lists

  # List with JSON output
  goog tasks lists --format json`,
	Args: cobra.NoArgs,
	RunE: runTasksLists,
}

// tasksCreateListCmd creates a new task list.
var tasksCreateListCmd = &cobra.Command{
	Use:   "create-list <title>",
	Short: "Create a new task list",
	Long: `Create a new task list with the specified title.

The new task list will be empty and ready to add tasks.`,
	Example: `  # Create a task list
  goog tasks create-list "Work Tasks"

  # Create a task list with a specific account
  goog tasks create-list "Personal" --account personal`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksCreateList,
}

// tasksDeleteListCmd deletes a task list.
var tasksDeleteListCmd = &cobra.Command{
	Use:   "delete-list <list-id>",
	Short: "Delete a task list",
	Long: `Delete a task list and all its tasks.

WARNING: This action is irreversible. All tasks in the list
will be permanently deleted.

The --confirm flag is required to prevent accidental deletion.`,
	Example: `  # Delete a task list (requires --confirm)
  goog tasks delete-list list123 --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !tasksDeleteConfirm {
			cmd.PrintErrln("Error: deletion requires --confirm flag")
			cmd.PrintErrln("Use --confirm to confirm this action")
			return fmt.Errorf("confirmation required")
		}
		return nil
	},
	RunE: runTasksDeleteList,
}

// tasksUpdateListCmd updates a task list.
var tasksUpdateListCmd = &cobra.Command{
	Use:   "update-list <list-id>",
	Short: "Update a task list title",
	Long: `Update the title of a task list.

Use the --title flag to specify the new title.`,
	Example: `  # Update a task list title
  goog tasks update-list list123 --title "New Title"`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if tasksTitle == "" {
			return fmt.Errorf("--title flag is required")
		}
		return nil
	},
	RunE: runTasksUpdateList,
}

// ================ Task Commands ================

// tasksListCmd lists tasks in a task list.
var tasksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks in a task list",
	Long: `List tasks in the specified task list.

Use --list to specify the task list ID (default: @default).
Use --show-completed to include completed tasks.`,
	Example: `  # List tasks in default list
  goog tasks list

  # List tasks with completed ones
  goog tasks list --show-completed

  # List tasks in a specific list
  goog tasks list --list "work-list-id"`,
	Args: cobra.NoArgs,
	RunE: runTasksList,
}

// tasksGetCmd gets a specific task.
var tasksGetCmd = &cobra.Command{
	Use:   "get <task-id>",
	Short: "Get details of a specific task",
	Long: `Get detailed information about a specific task.

Use --list to specify the task list ID (default: @default).`,
	Example: `  # Get a task
  goog tasks get task123

  # Get a task from a specific list
  goog tasks get task123 --list "work-list-id"`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksGet,
}

// tasksCreateCmd creates a new task.
var tasksCreateCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Long: `Create a new task with the specified title.

Use optional flags to set notes, due date, and parent task.`,
	Example: `  # Create a simple task
  goog tasks create "Buy groceries"

  # Create a task with notes and due date
  goog tasks create "Submit report" --notes "Q4 financial report" --due "2024-12-31"

  # Create a subtask
  goog tasks create "Review slides" --parent "task-parent-id"`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksCreate,
}

// tasksUpdateCmd updates a task.
var tasksUpdateCmd = &cobra.Command{
	Use:   "update <task-id>",
	Short: "Update a task",
	Long: `Update properties of an existing task.

Use optional flags to update title, notes, or due date.`,
	Example: `  # Update task title
  goog tasks update task123 --title "New title"

  # Update notes
  goog tasks update task123 --notes "Updated notes"

  # Update due date
  goog tasks update task123 --due "2024-12-31"`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if tasksTitle == "" && tasksNotes == "" && tasksDue == "" {
			return fmt.Errorf("at least one of --title, --notes, or --due must be specified")
		}
		return nil
	},
	RunE: runTasksUpdate,
}

// tasksCompleteCmd marks a task as completed.
var tasksCompleteCmd = &cobra.Command{
	Use:   "complete <task-id>",
	Short: "Mark a task as completed",
	Long: `Mark a task as completed.

The task status will be changed to "completed" and the
completion date will be set to now.`,
	Example: `  # Complete a task
  goog tasks complete task123`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksComplete,
}

// tasksReopenCmd reopens a completed task.
var tasksReopenCmd = &cobra.Command{
	Use:   "reopen <task-id>",
	Short: "Reopen a completed task",
	Long: `Reopen a completed task, changing its status to "needsAction".

The completion date will be cleared.`,
	Example: `  # Reopen a task
  goog tasks reopen task123`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksReopen,
}

// tasksDeleteCmd deletes a task.
var tasksDeleteCmd = &cobra.Command{
	Use:   "delete <task-id>",
	Short: "Delete a task",
	Long: `Delete a task permanently.

The --confirm flag is required to prevent accidental deletion.`,
	Example: `  # Delete a task (requires --confirm)
  goog tasks delete task123 --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !tasksDeleteConfirm {
			cmd.PrintErrln("Error: deletion requires --confirm flag")
			cmd.PrintErrln("Use --confirm to confirm this action")
			return fmt.Errorf("confirmation required")
		}
		return nil
	},
	RunE: runTasksDelete,
}

// tasksMoveCmd moves a task.
var tasksMoveCmd = &cobra.Command{
	Use:   "move <task-id>",
	Short: "Move a task to a different position",
	Long: `Move a task to a different position in the list or under a different parent.

Use --parent to move under a parent task (make it a subtask).
Use --previous to move after a specific task.`,
	Example: `  # Move task under a parent
  goog tasks move task123 --parent parent-task-id

  # Move task after another task
  goog tasks move task123 --previous previous-task-id`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksMove,
}

// tasksClearCmd clears completed tasks from a list.
var tasksClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear completed tasks from a list",
	Long: `Clear all completed tasks from the specified task list.

WARNING: This action is irreversible. All completed tasks
will be permanently deleted.

The --confirm flag is required to prevent accidental deletion.`,
	Example: `  # Clear completed tasks (requires --confirm)
  goog tasks clear --confirm

  # Clear from a specific list
  goog tasks clear --list "work-list-id" --confirm`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !tasksClearConfirm {
			cmd.PrintErrln("Error: clearing completed tasks requires --confirm flag")
			cmd.PrintErrln("Use --confirm to confirm this action")
			return fmt.Errorf("confirmation required")
		}
		return nil
	},
	RunE: runTasksClear,
}

// ================ Command Implementations ================

func runTasksLists(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskListRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// List task lists
	lists, err := repo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list task lists: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTaskLists(lists))

	return nil
}

func runTasksCreateList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	title := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskListRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Create task list
	taskList, err := domaintasks.NewTaskList(title)
	if err != nil {
		return fmt.Errorf("invalid task list: %w", err)
	}

	created, err := repo.Create(ctx, taskList)
	if err != nil {
		return fmt.Errorf("failed to create task list: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTaskList(created))

	return nil
}

func runTasksDeleteList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	listID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskListRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Delete task list
	err = repo.Delete(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to delete task list: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderSuccess(fmt.Sprintf("Task list '%s' deleted", listID)))

	return nil
}

func runTasksUpdateList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	listID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskListRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Get existing task list
	taskList, err := repo.Get(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to get task list: %w", err)
	}

	// Update title
	taskList.Title = tasksTitle
	taskList.Updated = time.Now()

	updated, err := repo.Update(ctx, taskList)
	if err != nil {
		return fmt.Errorf("failed to update task list: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTaskList(updated))

	return nil
}

func runTasksList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Build list options
	opts := domaintasks.ListOptions{
		ShowCompleted: tasksShowCompleted,
		ShowHidden:    tasksShowHidden,
		MaxResults:    tasksMaxResults,
	}

	// List tasks
	result, err := repo.List(ctx, tasksListID, opts)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTasks(result.Items))

	return nil
}

func runTasksGet(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	taskID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Get task
	task, err := repo.Get(ctx, tasksListID, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTask(task))

	return nil
}

func runTasksCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	title := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Create task
	task, err := domaintasks.NewTask(title, tasksListID)
	if err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}

	// Set optional fields
	if tasksNotes != "" {
		task.Notes = tasksNotes
	}
	if tasksDue != "" {
		dueDate, err := time.Parse("2006-01-02", tasksDue)
		if err != nil {
			return fmt.Errorf("invalid due date format (use YYYY-MM-DD): %w", err)
		}
		task.Due = &dueDate
	}
	if tasksParent != "" {
		err := task.SetParent(tasksParent)
		if err != nil {
			return fmt.Errorf("invalid parent: %w", err)
		}
	}

	created, err := repo.Create(ctx, tasksListID, task)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTask(created))

	return nil
}

func runTasksUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	taskID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Get existing task
	task, err := repo.Get(ctx, tasksListID, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Update fields
	if tasksTitle != "" {
		task.Title = tasksTitle
	}
	if tasksNotes != "" {
		task.Notes = tasksNotes
	}
	if tasksDue != "" {
		dueDate, err := time.Parse("2006-01-02", tasksDue)
		if err != nil {
			return fmt.Errorf("invalid due date format (use YYYY-MM-DD): %w", err)
		}
		task.Due = &dueDate
	}
	task.Updated = time.Now()

	updated, err := repo.Update(ctx, tasksListID, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTask(updated))

	return nil
}

func runTasksComplete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	taskID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Get existing task
	task, err := repo.Get(ctx, tasksListID, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Complete task
	err = task.Complete()
	if err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	updated, err := repo.Update(ctx, tasksListID, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTask(updated))

	return nil
}

func runTasksReopen(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	taskID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Get existing task
	task, err := repo.Get(ctx, tasksListID, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Reopen task
	err = task.Reopen()
	if err != nil {
		return fmt.Errorf("failed to reopen task: %w", err)
	}

	updated, err := repo.Update(ctx, tasksListID, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTask(updated))

	return nil
}

func runTasksDelete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	taskID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Delete task
	err = repo.Delete(ctx, tasksListID, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderSuccess(fmt.Sprintf("Task '%s' deleted", taskID)))

	return nil
}

func runTasksMove(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()
	taskID := args[0]

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Move task
	moved, err := repo.Move(ctx, tasksListID, taskID, tasksParent, tasksPrevious)
	if err != nil {
		return fmt.Errorf("failed to move task: %w", err)
	}

	// Render output
	p := presenter.New(formatFlag)
	cmd.Println(p.RenderTask(moved))

	return nil
}

func runTasksClear(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	deps := GetDependencies()

	// Resolve account
	account, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("account resolution failed: %w", err)
	}

	// Get token source
	tokenSource, err := deps.AccountService.GetTokenManager().GetTokenSource(ctx, account.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token source: %w", err)
	}

	// Create repository
	repo, err := deps.RepoFactory.NewTaskRepository(ctx, tokenSource)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Clear completed tasks
	err = repo.Clear(ctx, tasksListID)
	if err != nil {
		return fmt.Errorf("failed to clear completed tasks: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderSuccess("Completed tasks cleared"))

	return nil
}

// init registers all tasks commands.
func init() {
	// Register tasks command group
	rootCmd.AddCommand(tasksCmd)

	// Task list commands
	tasksCmd.AddCommand(tasksListsCmd)
	tasksCmd.AddCommand(tasksCreateListCmd)
	tasksCmd.AddCommand(tasksDeleteListCmd)
	tasksCmd.AddCommand(tasksUpdateListCmd)

	// Task commands
	tasksCmd.AddCommand(tasksListCmd)
	tasksCmd.AddCommand(tasksGetCmd)
	tasksCmd.AddCommand(tasksCreateCmd)
	tasksCmd.AddCommand(tasksUpdateCmd)
	tasksCmd.AddCommand(tasksCompleteCmd)
	tasksCmd.AddCommand(tasksReopenCmd)
	tasksCmd.AddCommand(tasksDeleteCmd)
	tasksCmd.AddCommand(tasksMoveCmd)
	tasksCmd.AddCommand(tasksClearCmd)

	// Persistent flags for all tasks commands
	tasksCmd.PersistentFlags().StringVar(&tasksListID, "list", "@default", "task list ID")

	// Flags for tasks list
	tasksListCmd.Flags().BoolVar(&tasksShowCompleted, "show-completed", false, "include completed tasks")
	tasksListCmd.Flags().BoolVar(&tasksShowHidden, "show-hidden", false, "include hidden tasks")
	tasksListCmd.Flags().Int64Var(&tasksMaxResults, "max-results", 100, "maximum number of tasks to return")

	// Flags for tasks create
	tasksCreateCmd.Flags().StringVar(&tasksNotes, "notes", "", "task notes")
	tasksCreateCmd.Flags().StringVar(&tasksDue, "due", "", "due date (YYYY-MM-DD)")
	tasksCreateCmd.Flags().StringVar(&tasksParent, "parent", "", "parent task ID (for subtasks)")

	// Flags for tasks update
	tasksUpdateCmd.Flags().StringVar(&tasksTitle, "title", "", "new task title")
	tasksUpdateCmd.Flags().StringVar(&tasksNotes, "notes", "", "new task notes")
	tasksUpdateCmd.Flags().StringVar(&tasksDue, "due", "", "new due date (YYYY-MM-DD)")

	// Flags for update-list
	tasksUpdateListCmd.Flags().StringVar(&tasksTitle, "title", "", "new task list title")

	// Flags for delete operations
	tasksDeleteListCmd.Flags().BoolVar(&tasksDeleteConfirm, "confirm", false, "confirm deletion")
	tasksDeleteCmd.Flags().BoolVar(&tasksDeleteConfirm, "confirm", false, "confirm deletion")
	tasksClearCmd.Flags().BoolVar(&tasksClearConfirm, "confirm", false, "confirm clearing completed tasks")

	// Flags for move
	tasksMoveCmd.Flags().StringVar(&tasksParent, "parent", "", "parent task ID")
	tasksMoveCmd.Flags().StringVar(&tasksPrevious, "previous", "", "previous sibling task ID")
}
