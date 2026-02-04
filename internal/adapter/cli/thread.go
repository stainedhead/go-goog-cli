// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	"github.com/stainedhead/go-goog-cli/internal/adapter/repository"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// Thread command flags.
var (
	threadMaxResults   int
	threadLabels       []string
	threadAddLabels    []string
	threadRemoveLabels []string
)

// threadCmd represents the thread command group.
var threadCmd = &cobra.Command{
	Use:   "thread",
	Short: "Manage email threads",
	Long: `Manage email threads (conversations) in your Gmail account.

Threads group related messages together in a conversation.
The thread commands allow you to list, view, trash, and
modify labels on entire threads at once.`,
}

// threadListCmd lists threads.
var threadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List email threads",
	Long: `List email threads in your account.

Displays threads with their ID, snippet, and message count.
Use --labels to filter by specific labels.`,
	Aliases: []string{"ls"},
	Example: `  # List recent threads
  goog thread list

  # List threads with limit
  goog thread list --max-results 20

  # List threads with specific labels
  goog thread list --labels INBOX --labels UNREAD

  # List threads with JSON output
  goog thread list --format json`,
	RunE: runThreadList,
}

// threadShowCmd shows a single thread with all messages.
var threadShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show thread with all messages",
	Long: `Show the full content of a thread including all messages.

Displays thread details including all messages in the
conversation in chronological order.`,
	Aliases: []string{"get", "read"},
	Example: `  # Show thread by ID
  goog thread show abc123

  # Show thread with JSON output
  goog thread show abc123 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runThreadShow,
}

// threadTrashCmd moves a thread to trash.
var threadTrashCmd = &cobra.Command{
	Use:   "trash <id>",
	Short: "Move thread to trash",
	Long: `Move an entire thread to trash.

All messages in the thread will be moved to the Trash label.
Trashed threads can be recovered within 30 days.`,
	Example: `  # Trash a thread by ID
  goog thread trash abc123`,
	Args: cobra.ExactArgs(1),
	RunE: runThreadTrash,
}

// threadModifyCmd modifies labels on a thread.
var threadModifyCmd = &cobra.Command{
	Use:   "modify <id>",
	Short: "Modify thread labels",
	Long: `Modify labels on an entire thread.

Add or remove labels from all messages in the thread
using the --add-labels and --remove-labels flags.`,
	Example: `  # Add labels to a thread
  goog thread modify abc123 --add-labels "Work" --add-labels "Important"

  # Remove labels from a thread
  goog thread modify abc123 --remove-labels "INBOX"

  # Add and remove labels simultaneously
  goog thread modify abc123 --add-labels "Archive" --remove-labels "INBOX"`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(threadAddLabels) == 0 && len(threadRemoveLabels) == 0 {
			return fmt.Errorf("at least one of --add-labels or --remove-labels is required")
		}
		return nil
	},
	RunE: runThreadModify,
}

func init() {
	// Add thread subcommands
	threadCmd.AddCommand(threadListCmd)
	threadCmd.AddCommand(threadShowCmd)
	threadCmd.AddCommand(threadTrashCmd)
	threadCmd.AddCommand(threadModifyCmd)

	// List flags
	threadListCmd.Flags().IntVar(&threadMaxResults, "max-results", 20, "maximum number of threads to list")
	threadListCmd.Flags().StringSliceVar(&threadLabels, "labels", nil, "filter by label IDs")

	// Modify flags
	threadModifyCmd.Flags().StringSliceVar(&threadAddLabels, "add-labels", nil, "labels to add")
	threadModifyCmd.Flags().StringSliceVar(&threadRemoveLabels, "remove-labels", nil, "labels to remove")

	// Add to root
	rootCmd.AddCommand(threadCmd)
}

// getThreadRepository creates a thread repository for the current account.
// Deprecated: Use getThreadRepositoryFromDeps for testability.
func getThreadRepository(ctx context.Context) (*repository.GmailThreadRepository, error) {
	tokenSource, err := getTokenSource(ctx)
	if err != nil {
		return nil, err
	}

	// Create Gmail repository
	gmailRepo, err := repository.NewGmailRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail repository: %w", err)
	}

	return repository.NewGmailThreadRepository(gmailRepo), nil
}

// runThreadList handles the thread list command.
func runThreadList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repo, err := getThreadRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	opts := mail.ListOptions{
		MaxResults: threadMaxResults,
		LabelIDs:   threadLabels,
	}

	result, err := repo.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list threads: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderThreads(result.Items)
	cmd.Println(output)

	if !quietFlag && result.NextPageToken != "" {
		cmd.Println("\n(More threads available. Use --max-results to adjust.)")
	}

	return nil
}

// runThreadShow handles the thread show command.
func runThreadShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	threadID := args[0]

	repo, err := getThreadRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	thread, err := repo.Get(ctx, threadID)
	if err != nil {
		return fmt.Errorf("failed to get thread: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderThread(thread)
	cmd.Println(output)

	// Show messages content if not in quiet mode and not JSON
	if !quietFlag && formatFlag != "json" {
		for i, msg := range thread.Messages {
			cmd.Printf("\n--- Message %d of %d ---\n", i+1, len(thread.Messages))
			cmd.Printf("From: %s\n", msg.From)
			cmd.Printf("Date: %s\n", msg.Date.Format("Mon, 02 Jan 2006 15:04:05"))
			cmd.Printf("Subject: %s\n", msg.Subject)
			if msg.Body != "" {
				cmd.Println("\n" + msg.Body)
			}
		}
	}

	return nil
}

// runThreadTrash handles the thread trash command.
func runThreadTrash(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	threadID := args[0]

	repo, err := getThreadRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	if err := repo.Trash(ctx, threadID); err != nil {
		return fmt.Errorf("failed to trash thread: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Thread %s moved to trash.\n", threadID)
	}

	return nil
}

// runThreadModify handles the thread modify command.
func runThreadModify(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	threadID := args[0]

	repo, err := getThreadRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	req := mail.ModifyRequest{
		AddLabels:    threadAddLabels,
		RemoveLabels: threadRemoveLabels,
	}

	thread, err := repo.Modify(ctx, threadID, req)
	if err != nil {
		return fmt.Errorf("failed to modify thread: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderThread(thread))
	} else {
		cmd.Printf("Thread %s modified successfully.\n", threadID)
		if len(threadAddLabels) > 0 {
			cmd.Printf("Added labels: %s\n", strings.Join(threadAddLabels, ", "))
		}
		if len(threadRemoveLabels) > 0 {
			cmd.Printf("Removed labels: %s\n", strings.Join(threadRemoveLabels, ", "))
		}
	}

	return nil
}
