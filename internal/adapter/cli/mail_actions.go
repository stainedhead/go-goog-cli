// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	"github.com/stainedhead/go-goog-cli/internal/adapter/repository"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// Command flags for mail actions.
var (
	mailDeleteConfirm      bool
	mailModifyAddLabels    []string
	mailModifyRemoveLabels []string
	mailMarkRead           bool
	mailMarkUnread         bool
	mailMarkStar           bool
	mailMarkUnstar         bool
	mailListMaxResults     int
	mailListLabels         []string
	mailListUnreadOnly     bool
	mailSearchMaxResults   int
)

// mailCmd represents the mail command group.
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Manage Gmail messages",
	Long: `Manage Gmail messages.

The mail commands allow you to list, read, search, and manage
email messages in your Gmail account.`,
}

// mailTrashCmd moves a message to trash.
var mailTrashCmd = &cobra.Command{
	Use:   "trash <message-id>",
	Short: "Move a message to trash",
	Long: `Move a message to the trash folder.

The message can be restored using the 'untrash' command within
30 days. After that, it will be permanently deleted.`,
	Example: `  # Trash a message by ID
  goog mail trash msg123abc

  # Trash a message using a specific account
  goog mail trash msg123abc --account work`,
	Args: cobra.ExactArgs(1),
	RunE: runMailTrash,
}

// mailUntrashCmd restores a message from trash.
var mailUntrashCmd = &cobra.Command{
	Use:   "untrash <message-id>",
	Short: "Restore a message from trash",
	Long: `Restore a message from the trash folder.

This removes the message from trash and restores it to its
previous labels.`,
	Example: `  # Restore a message from trash
  goog mail untrash msg123abc`,
	Args: cobra.ExactArgs(1),
	RunE: runMailUntrash,
}

// mailArchiveCmd archives a message.
var mailArchiveCmd = &cobra.Command{
	Use:   "archive <message-id>",
	Short: "Archive a message",
	Long: `Archive a message by removing the INBOX label.

The message remains accessible via search and other labels
but will no longer appear in the inbox.`,
	Example: `  # Archive a message
  goog mail archive msg123abc`,
	Args: cobra.ExactArgs(1),
	RunE: runMailArchive,
}

// mailDeleteCmd permanently deletes a message.
var mailDeleteCmd = &cobra.Command{
	Use:   "delete <message-id>",
	Short: "Permanently delete a message",
	Long: `Permanently delete a message.

WARNING: This action is irreversible. The message will be
permanently deleted and cannot be recovered.

The --confirm flag is required to prevent accidental deletion.`,
	Example: `  # Permanently delete a message (requires --confirm)
  goog mail delete msg123abc --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !mailDeleteConfirm {
			cmd.PrintErrln("Error: permanent deletion requires --confirm flag")
			cmd.PrintErrln("This action is irreversible. Use 'goog mail trash' for recoverable deletion.")
			return fmt.Errorf("--confirm flag required for permanent deletion")
		}
		return nil
	},
	RunE: runMailDelete,
}

// mailModifyCmd modifies message labels.
var mailModifyCmd = &cobra.Command{
	Use:   "modify <message-id>",
	Short: "Modify message labels",
	Long: `Modify the labels on a message.

Use --add-labels to add labels and --remove-labels to remove labels.
At least one of these flags must be specified.

Common system labels:
  INBOX, STARRED, IMPORTANT, UNREAD, SPAM, TRASH, CATEGORY_PERSONAL,
  CATEGORY_SOCIAL, CATEGORY_PROMOTIONS, CATEGORY_UPDATES, CATEGORY_FORUMS`,
	Example: `  # Add a label to a message
  goog mail modify msg123abc --add-labels STARRED

  # Remove a label from a message
  goog mail modify msg123abc --remove-labels INBOX

  # Add and remove labels at the same time
  goog mail modify msg123abc --add-labels STARRED,IMPORTANT --remove-labels INBOX`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(mailModifyAddLabels) == 0 && len(mailModifyRemoveLabels) == 0 {
			cmd.PrintErrln("Error: at least one of --add-labels or --remove-labels is required")
			return fmt.Errorf("no labels specified to add or remove")
		}
		return nil
	},
	RunE: runMailModify,
}

// mailMarkCmd marks a message as read/unread/starred.
var mailMarkCmd = &cobra.Command{
	Use:   "mark <message-id>",
	Short: "Mark a message as read/unread/starred",
	Long: `Mark a message as read, unread, starred, or unstarred.

At least one flag must be specified. You can combine --read/--unread
with --star/--unstar, but cannot use conflicting flags together.`,
	Example: `  # Mark a message as read
  goog mail mark msg123abc --read

  # Mark a message as unread
  goog mail mark msg123abc --unread

  # Star a message
  goog mail mark msg123abc --star

  # Unstar a message
  goog mail mark msg123abc --unstar

  # Mark as read and star in one command
  goog mail mark msg123abc --read --star`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !mailMarkRead && !mailMarkUnread && !mailMarkStar && !mailMarkUnstar {
			cmd.PrintErrln("Error: at least one of --read, --unread, --star, or --unstar is required")
			return fmt.Errorf("no mark action specified")
		}
		if mailMarkRead && mailMarkUnread {
			cmd.PrintErrln("Error: cannot use both --read and --unread")
			return fmt.Errorf("conflicting flags: --read and --unread")
		}
		if mailMarkStar && mailMarkUnstar {
			cmd.PrintErrln("Error: cannot use both --star and --unstar")
			return fmt.Errorf("conflicting flags: --star and --unstar")
		}
		return nil
	},
	RunE: runMailMark,
}

// mailListCmd lists messages in the inbox.
var mailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List messages",
	Long: `List messages from your Gmail inbox.

By default, lists messages from the INBOX label. Use --labels
to filter by specific labels and --unread-only to show only
unread messages.`,
	Example: `  # List recent inbox messages
  goog mail list

  # List messages with specific labels
  goog mail list --labels INBOX,IMPORTANT

  # List only unread messages
  goog mail list --unread-only

  # List with JSON output
  goog mail list --format json

  # List more messages
  goog mail list --max-results 50`,
	Aliases: []string{"ls"},
	RunE:    runMailList,
}

// mailReadCmd reads a single message.
var mailReadCmd = &cobra.Command{
	Use:   "read <message-id>",
	Short: "Read a single message",
	Long: `Read and display a single email message.

Retrieves the full content of the specified message including
headers, body, and metadata.`,
	Example: `  # Read a message by ID
  goog mail read 18abc123def456

  # Read with JSON output
  goog mail read 18abc123def456 --format json

  # Read with plain text output
  goog mail read 18abc123def456 --format plain`,
	Aliases: []string{"get", "show"},
	Args:    cobra.ExactArgs(1),
	RunE:    runMailRead,
}

// mailSearchCmd searches for messages.
var mailSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for messages",
	Long: `Search for messages using Gmail query syntax.

Uses the same query format as the Gmail web interface.
Common operators include:
  - from:sender@example.com
  - to:recipient@example.com
  - subject:keyword
  - is:unread, is:starred, is:important
  - has:attachment
  - after:YYYY/MM/DD, before:YYYY/MM/DD
  - label:labelname`,
	Example: `  # Search for unread messages
  goog mail search "is:unread"

  # Search from a specific sender
  goog mail search "from:example@gmail.com"

  # Search with subject keyword
  goog mail search "subject:meeting"

  # Combine search terms
  goog mail search "from:boss@company.com is:unread after:2024/01/01"

  # Search with JSON output
  goog mail search "has:attachment" --format json`,
	Aliases: []string{"find", "query"},
	Args:    cobra.ExactArgs(1),
	RunE:    runMailSearch,
}

func init() {
	// Add mail subcommands
	mailCmd.AddCommand(mailListCmd)
	mailCmd.AddCommand(mailReadCmd)
	mailCmd.AddCommand(mailSearchCmd)
	mailCmd.AddCommand(mailTrashCmd)
	mailCmd.AddCommand(mailUntrashCmd)
	mailCmd.AddCommand(mailArchiveCmd)
	mailCmd.AddCommand(mailDeleteCmd)
	mailCmd.AddCommand(mailModifyCmd)
	mailCmd.AddCommand(mailMarkCmd)

	// List command flags
	mailListCmd.Flags().IntVar(&mailListMaxResults, "max-results", 10, "maximum number of messages to return")
	mailListCmd.Flags().StringSliceVar(&mailListLabels, "labels", []string{"INBOX"}, "filter by labels")
	mailListCmd.Flags().BoolVar(&mailListUnreadOnly, "unread-only", false, "show only unread messages")

	// Search command flags
	mailSearchCmd.Flags().IntVar(&mailSearchMaxResults, "max-results", 10, "maximum number of messages to return")

	// Delete flags
	mailDeleteCmd.Flags().BoolVar(&mailDeleteConfirm, "confirm", false, "confirm permanent deletion")

	// Modify flags
	mailModifyCmd.Flags().StringSliceVar(&mailModifyAddLabels, "add-labels", nil, "labels to add (comma-separated)")
	mailModifyCmd.Flags().StringSliceVar(&mailModifyRemoveLabels, "remove-labels", nil, "labels to remove (comma-separated)")

	// Mark flags
	mailMarkCmd.Flags().BoolVar(&mailMarkRead, "read", false, "mark as read")
	mailMarkCmd.Flags().BoolVar(&mailMarkUnread, "unread", false, "mark as unread")
	mailMarkCmd.Flags().BoolVar(&mailMarkStar, "star", false, "add star")
	mailMarkCmd.Flags().BoolVar(&mailMarkUnstar, "unstar", false, "remove star")

	// Add mail command to root
	rootCmd.AddCommand(mailCmd)
}

// getGmailRepository creates a GmailRepository using the current account's credentials.
// Returns the repository and the sender's email address.
func getGmailRepository(ctx context.Context) (*repository.GmailRepository, string, error) {
	tokenSource, email, err := getTokenSourceWithEmail(ctx)
	if err != nil {
		return nil, "", err
	}

	// Create Gmail repository
	repo, err := repository.NewGmailRepository(ctx, tokenSource)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Gmail client: %w", err)
	}

	return repo, email, nil
}

// runMailList handles the mail list command.
func runMailList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get Gmail repository
	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	// Build list options
	opts := mail.ListOptions{
		MaxResults: mailListMaxResults,
		LabelIDs:   mailListLabels,
	}

	// Add unread filter if requested
	if mailListUnreadOnly {
		opts.Query = "is:unread"
	}

	// List messages
	result, err := repo.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list messages: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderMessages(result.Items)
	cmd.Println(output)

	return nil
}

// runMailRead handles the mail read command.
func runMailRead(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	messageID := args[0]

	// Get Gmail repository
	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	// Get the message
	msg, err := repo.Get(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderMessage(msg)
	cmd.Println(output)

	// For plain format, also show the body content
	if formatFlag == "plain" && msg.Body != "" {
		cmd.Println("\n--- Message Body ---")
		cmd.Println(msg.Body)
	}

	return nil
}

// runMailSearch handles the mail search command.
func runMailSearch(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	query := args[0]

	// Get Gmail repository
	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	// Build search options
	opts := mail.ListOptions{
		MaxResults: mailSearchMaxResults,
	}

	// Search messages
	result, err := repo.Search(ctx, query, opts)
	if err != nil {
		return fmt.Errorf("failed to search messages: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderMessages(result.Items)
	cmd.Println(output)

	// Show result count if not empty
	if len(result.Items) > 0 && !quietFlag {
		cmd.Printf("\nFound %d message(s)", len(result.Items))
		if result.Total > len(result.Items) {
			cmd.Printf(" (showing first %d of ~%d)", len(result.Items), result.Total)
		}
		cmd.Println()
	}

	return nil
}

// runMailTrash handles the mail trash command.
func runMailTrash(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	messageID := args[0]

	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	if err := repo.Trash(ctx, messageID); err != nil {
		return fmt.Errorf("failed to trash message: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Message %s moved to trash\n", messageID)
	}
	return nil
}

// runMailUntrash handles the mail untrash command.
func runMailUntrash(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	messageID := args[0]

	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	if err := repo.Untrash(ctx, messageID); err != nil {
		return fmt.Errorf("failed to restore message from trash: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Message %s restored from trash\n", messageID)
	}
	return nil
}

// runMailArchive handles the mail archive command.
func runMailArchive(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	messageID := args[0]

	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	if err := repo.Archive(ctx, messageID); err != nil {
		return fmt.Errorf("failed to archive message: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Message %s archived (removed from INBOX)\n", messageID)
	}
	return nil
}

// runMailDelete handles the mail delete command.
func runMailDelete(cmd *cobra.Command, args []string) error {
	messageID := args[0]
	ctx := context.Background()

	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	if err := repo.Delete(ctx, messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Message %s permanently deleted\n", messageID)
	}
	return nil
}

// runMailModify handles the mail modify command.
func runMailModify(cmd *cobra.Command, args []string) error {
	messageID := args[0]
	ctx := context.Background()

	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	req := mail.ModifyRequest{
		AddLabels:    mailModifyAddLabels,
		RemoveLabels: mailModifyRemoveLabels,
	}

	msg, err := repo.Modify(ctx, messageID, req)
	if err != nil {
		return fmt.Errorf("failed to modify message labels: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Message %s labels modified\n", messageID)
		if verboseFlag && msg != nil {
			cmd.Printf("Current labels: %v\n", msg.Labels)
		}
	}
	return nil
}

// runMailMark handles the mail mark command.
func runMailMark(cmd *cobra.Command, args []string) error {
	messageID := args[0]
	ctx := context.Background()

	repo, _, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	// Build modify request based on flags
	req := mail.ModifyRequest{
		AddLabels:    []string{},
		RemoveLabels: []string{},
	}

	var actions []string

	if mailMarkRead {
		req.RemoveLabels = append(req.RemoveLabels, "UNREAD")
		actions = append(actions, "marked as read")
	}
	if mailMarkUnread {
		req.AddLabels = append(req.AddLabels, "UNREAD")
		actions = append(actions, "marked as unread")
	}
	if mailMarkStar {
		req.AddLabels = append(req.AddLabels, "STARRED")
		actions = append(actions, "starred")
	}
	if mailMarkUnstar {
		req.RemoveLabels = append(req.RemoveLabels, "STARRED")
		actions = append(actions, "unstarred")
	}

	_, err = repo.Modify(ctx, messageID, req)
	if err != nil {
		return fmt.Errorf("failed to mark message: %w", err)
	}

	if !quietFlag {
		actionStr := ""
		for i, a := range actions {
			if i > 0 {
				actionStr += " and "
			}
			actionStr += a
		}
		cmd.Printf("Message %s %s\n", messageID, actionStr)
	}
	return nil
}
