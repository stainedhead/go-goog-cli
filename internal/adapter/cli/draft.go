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
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/keyring"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// Draft command flags.
var (
	draftTo      []string
	draftSubject string
	draftBody    string
	draftLimit   int
)

// draftCmd represents the draft command group.
var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Manage email drafts",
	Long: `Manage email drafts in your Gmail account.

The draft commands allow you to list, view, create, update,
send, and delete email drafts.`,
}

// draftListCmd lists all drafts.
var draftListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all drafts",
	Long: `List all email drafts in your account.

Displays a table showing draft ID, subject, recipients, and
last updated date.`,
	Aliases: []string{"ls"},
	Example: `  # List all drafts
  goog draft list

  # List drafts with JSON output
  goog draft list --format json

  # List drafts with limit
  goog draft list --limit 10`,
	RunE: runDraftList,
}

// draftShowCmd shows a single draft.
var draftShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show draft content",
	Long: `Show the full content of a draft.

Displays the draft details including recipient, subject,
body, and metadata.`,
	Aliases: []string{"get", "read"},
	Example: `  # Show draft by ID
  goog draft show abc123

  # Show draft with JSON output
  goog draft show abc123 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runDraftShow,
}

// draftCreateCmd creates a new draft.
var draftCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new draft",
	Long: `Create a new email draft.

The draft will be saved but not sent. You can later
edit it or send it using the send command.`,
	Example: `  # Create a draft with all fields
  goog draft create --to user@example.com --subject "Hello" --body "Message content"

  # Create a draft with multiple recipients
  goog draft create --to user1@example.com --to user2@example.com --subject "Group message"`,
	RunE: runDraftCreate,
}

// draftUpdateCmd updates an existing draft.
var draftUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a draft",
	Long: `Update an existing email draft.

Specify which fields to update using flags. Fields not
specified will retain their current values.`,
	Example: `  # Update subject
  goog draft update abc123 --subject "New Subject"

  # Update body
  goog draft update abc123 --body "New content"

  # Update recipient and subject
  goog draft update abc123 --to new@example.com --subject "Updated"`,
	Args: cobra.ExactArgs(1),
	RunE: runDraftUpdate,
}

// draftSendCmd sends a draft.
var draftSendCmd = &cobra.Command{
	Use:   "send <id>",
	Short: "Send a draft",
	Long: `Send an existing draft as an email.

The draft will be deleted after sending and the message
will appear in your Sent folder.`,
	Example: `  # Send a draft by ID
  goog draft send abc123`,
	Args: cobra.ExactArgs(1),
	RunE: runDraftSend,
}

// draftDeleteCmd deletes a draft.
var draftDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a draft",
	Long: `Delete a draft permanently.

This action cannot be undone.`,
	Aliases: []string{"rm", "remove"},
	Example: `  # Delete a draft by ID
  goog draft delete abc123`,
	Args: cobra.ExactArgs(1),
	RunE: runDraftDelete,
}

func init() {
	// Add draft subcommands
	draftCmd.AddCommand(draftListCmd)
	draftCmd.AddCommand(draftShowCmd)
	draftCmd.AddCommand(draftCreateCmd)
	draftCmd.AddCommand(draftUpdateCmd)
	draftCmd.AddCommand(draftSendCmd)
	draftCmd.AddCommand(draftDeleteCmd)

	// List flags
	draftListCmd.Flags().IntVar(&draftLimit, "limit", 20, "maximum number of drafts to list")

	// Create flags
	draftCreateCmd.Flags().StringSliceVar(&draftTo, "to", nil, "recipient email addresses")
	draftCreateCmd.Flags().StringVar(&draftSubject, "subject", "", "email subject")
	draftCreateCmd.Flags().StringVar(&draftBody, "body", "", "email body")

	// Update flags
	draftUpdateCmd.Flags().StringSliceVar(&draftTo, "to", nil, "recipient email addresses")
	draftUpdateCmd.Flags().StringVar(&draftSubject, "subject", "", "email subject")
	draftUpdateCmd.Flags().StringVar(&draftBody, "body", "", "email body")

	// Add to root
	rootCmd.AddCommand(draftCmd)
}

// getDraftRepository creates a draft repository for the current account.
func getDraftRepository(ctx context.Context) (*repository.GmailDraftRepository, error) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create account service
	svc := accountuc.NewService(cfg, store, nil)

	// Resolve account
	acc, err := svc.ResolveAccount(accountFlag)
	if err != nil {
		return nil, fmt.Errorf("no account found: %w", err)
	}

	// Get token source
	tokenMgr := svc.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create Gmail repository
	gmailRepo, err := repository.NewGmailRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail repository: %w", err)
	}

	return repository.NewGmailDraftRepository(gmailRepo), nil
}

// runDraftList handles the draft list command.
func runDraftList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repo, err := getDraftRepository(ctx)
	if err != nil {
		return err
	}

	opts := mail.ListOptions{
		MaxResults: draftLimit,
	}

	result, err := repo.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list drafts: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderDrafts(result.Items)
	cmd.Println(output)

	if !quietFlag && result.NextPageToken != "" {
		cmd.Println("\n(More drafts available. Use --limit to adjust.)")
	}

	return nil
}

// runDraftShow handles the draft show command.
func runDraftShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	draftID := args[0]

	repo, err := getDraftRepository(ctx)
	if err != nil {
		return err
	}

	draft, err := repo.Get(ctx, draftID)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderDraft(draft)
	cmd.Println(output)

	// Show body content if available and not in quiet mode
	if !quietFlag && draft.Message != nil && draft.Message.Body != "" {
		cmd.Println("\n--- Body ---")
		cmd.Println(draft.Message.Body)
	}

	return nil
}

// runDraftCreate handles the draft create command.
func runDraftCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Validate required fields
	if len(draftTo) == 0 {
		return fmt.Errorf("at least one recipient is required (--to)")
	}
	if draftSubject == "" {
		return fmt.Errorf("subject is required (--subject)")
	}

	repo, err := getDraftRepository(ctx)
	if err != nil {
		return err
	}

	// Create draft message
	msg := &mail.Message{
		To:      draftTo,
		Subject: draftSubject,
		Body:    draftBody,
	}

	draft := &mail.Draft{
		Message: msg,
	}

	created, err := repo.Create(ctx, draft)
	if err != nil {
		return fmt.Errorf("failed to create draft: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderDraft(created))
	} else {
		cmd.Printf("Draft created successfully.\n")
		cmd.Printf("ID: %s\n", created.ID)
		if created.Message != nil {
			cmd.Printf("To: %s\n", strings.Join(created.Message.To, ", "))
			cmd.Printf("Subject: %s\n", created.Message.Subject)
		}
	}

	return nil
}

// runDraftUpdate handles the draft update command.
func runDraftUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	draftID := args[0]

	repo, err := getDraftRepository(ctx)
	if err != nil {
		return err
	}

	// Get existing draft
	existing, err := repo.Get(ctx, draftID)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	// Update fields if provided
	if existing.Message == nil {
		existing.Message = &mail.Message{}
	}

	if len(draftTo) > 0 {
		existing.Message.To = draftTo
	}
	if draftSubject != "" {
		existing.Message.Subject = draftSubject
	}
	if draftBody != "" {
		existing.Message.Body = draftBody
	}

	updated, err := repo.Update(ctx, existing)
	if err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderDraft(updated))
	} else {
		cmd.Printf("Draft updated successfully.\n")
		cmd.Printf("ID: %s\n", updated.ID)
	}

	return nil
}

// runDraftSend handles the draft send command.
func runDraftSend(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	draftID := args[0]

	repo, err := getDraftRepository(ctx)
	if err != nil {
		return err
	}

	sent, err := repo.Send(ctx, draftID)
	if err != nil {
		return fmt.Errorf("failed to send draft: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderMessage(sent))
	} else {
		cmd.Printf("Draft sent successfully.\n")
		cmd.Printf("Message ID: %s\n", sent.ID)
		if len(sent.To) > 0 {
			cmd.Printf("To: %s\n", strings.Join(sent.To, ", "))
		}
		cmd.Printf("Subject: %s\n", sent.Subject)
	}

	return nil
}

// runDraftDelete handles the draft delete command.
func runDraftDelete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	draftID := args[0]

	repo, err := getDraftRepository(ctx)
	if err != nil {
		return err
	}

	if err := repo.Delete(ctx, draftID); err != nil {
		return fmt.Errorf("failed to delete draft: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Draft %s deleted successfully.\n", draftID)
	}

	return nil
}
