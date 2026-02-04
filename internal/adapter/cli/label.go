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

// Label command flags.
var (
	labelBackgroundColor string
	labelTextColor       string
	labelConfirm         bool
)

// labelCmd represents the label command group.
var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage email labels",
	Long: `Manage email labels (folders/categories) in your Gmail account.

Labels help organize your email by categorizing messages.
The label commands allow you to list, view, create, update,
and delete labels.`,
}

// labelListCmd lists all labels.
var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all labels",
	Long: `List all email labels in your account.

Displays both system labels (INBOX, SENT, etc.) and
user-created labels.`,
	Aliases: []string{"ls"},
	Example: `  # List all labels
  goog label list

  # List labels with JSON output
  goog label list --format json`,
	RunE: runLabelList,
}

// labelShowCmd shows a single label.
var labelShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show label details",
	Long: `Show details of a specific label.

Displays the label ID, name, type, visibility settings,
and color if set.`,
	Aliases: []string{"get", "info"},
	Example: `  # Show label by name
  goog label show "Work"

  # Show system label
  goog label show "INBOX"

  # Show with JSON output
  goog label show "Work" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runLabelShow,
}

// labelCreateCmd creates a new label.
var labelCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new label",
	Long: `Create a new email label.

You can optionally set the label color using background
and text color flags.

Color values should be valid hex colors (e.g., #4285f4).`,
	Example: `  # Create a simple label
  goog label create "Work Projects"

  # Create a label with colors
  goog label create "Urgent" --background "#ff0000" --text "#ffffff"`,
	Args: cobra.ExactArgs(1),
	RunE: runLabelCreate,
}

// labelUpdateCmd updates an existing label.
var labelUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update a label",
	Long: `Update an existing label.

You can update the label's color using the background
and text color flags. Note that system labels cannot
be modified.`,
	Example: `  # Update label colors
  goog label update "Work" --background "#4285f4" --text "#ffffff"`,
	Args: cobra.ExactArgs(1),
	RunE: runLabelUpdate,
}

// labelDeleteCmd deletes a label.
var labelDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a label",
	Long: `Delete a label permanently.

This action cannot be undone. Messages with this label
will not be deleted, but they will no longer have this
label attached.

System labels cannot be deleted.

Requires --confirm flag for safety.`,
	Aliases: []string{"rm", "remove"},
	Example: `  # Delete a label (requires confirmation)
  goog label delete "Old Projects" --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !labelConfirm {
			return fmt.Errorf("deletion requires --confirm flag")
		}
		return nil
	},
	RunE: runLabelDelete,
}

func init() {
	// Add label subcommands
	labelCmd.AddCommand(labelListCmd)
	labelCmd.AddCommand(labelShowCmd)
	labelCmd.AddCommand(labelCreateCmd)
	labelCmd.AddCommand(labelUpdateCmd)
	labelCmd.AddCommand(labelDeleteCmd)

	// Create flags
	labelCreateCmd.Flags().StringVar(&labelBackgroundColor, "background", "", "background color (hex, e.g., #4285f4)")
	labelCreateCmd.Flags().StringVar(&labelTextColor, "text", "", "text color (hex, e.g., #ffffff)")

	// Update flags
	labelUpdateCmd.Flags().StringVar(&labelBackgroundColor, "background", "", "background color (hex, e.g., #4285f4)")
	labelUpdateCmd.Flags().StringVar(&labelTextColor, "text", "", "text color (hex, e.g., #ffffff)")

	// Delete flags
	labelDeleteCmd.Flags().BoolVar(&labelConfirm, "confirm", false, "confirm deletion")

	// Add to root
	rootCmd.AddCommand(labelCmd)
}

// getLabelRepository creates a label repository for the current account.
func getLabelRepository(ctx context.Context) (*repository.GmailLabelRepository, error) {
	tokenSource, err := getTokenSource(ctx)
	if err != nil {
		return nil, err
	}

	// Create Gmail repository
	gmailRepo, err := repository.NewGmailRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail repository: %w", err)
	}

	return repository.NewGmailLabelRepository(gmailRepo), nil
}

// runLabelList handles the label list command.
func runLabelList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repo, err := getLabelRepository(ctx)
	if err != nil {
		return err
	}

	labels, err := repo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list labels: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderLabels(labels)
	cmd.Println(output)

	return nil
}

// runLabelShow handles the label show command.
func runLabelShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	labelName := args[0]

	repo, err := getLabelRepository(ctx)
	if err != nil {
		return err
	}

	// Try to find label by name first
	label, err := repo.GetByName(ctx, labelName)
	if err != nil {
		// If not found by name, try by ID
		label, err = repo.Get(ctx, labelName)
		if err != nil {
			return fmt.Errorf("label not found: %s", labelName)
		}
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderLabel(label)
	cmd.Println(output)

	return nil
}

// runLabelCreate handles the label create command.
func runLabelCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	labelName := args[0]

	repo, err := getLabelRepository(ctx)
	if err != nil {
		return err
	}

	// Create new label
	label := mail.NewLabel("", labelName)

	// Set colors if provided
	if labelBackgroundColor != "" || labelTextColor != "" {
		bg := labelBackgroundColor
		if bg == "" {
			bg = "#000000"
		}
		text := labelTextColor
		if text == "" {
			text = "#ffffff"
		}
		label.SetColor(bg, text)
	}

	created, err := repo.Create(ctx, label)
	if err != nil {
		return fmt.Errorf("failed to create label: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderLabel(created))
	} else {
		cmd.Printf("Label created successfully.\n")
		cmd.Printf("ID: %s\n", created.ID)
		cmd.Printf("Name: %s\n", created.Name)
		if created.Color != nil {
			cmd.Printf("Background: %s\n", created.Color.Background)
			cmd.Printf("Text: %s\n", created.Color.Text)
		}
	}

	return nil
}

// runLabelUpdate handles the label update command.
func runLabelUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	labelName := args[0]

	repo, err := getLabelRepository(ctx)
	if err != nil {
		return err
	}

	// Get existing label
	label, err := repo.GetByName(ctx, labelName)
	if err != nil {
		// Try by ID if not found by name
		label, err = repo.Get(ctx, labelName)
		if err != nil {
			return fmt.Errorf("label not found: %s", labelName)
		}
	}

	// Check if system label
	if label.IsSystemLabel() {
		return fmt.Errorf("cannot modify system label: %s", label.Name)
	}

	// Update colors if provided
	if labelBackgroundColor != "" || labelTextColor != "" {
		bg := labelBackgroundColor
		text := labelTextColor
		if label.Color != nil {
			if bg == "" {
				bg = label.Color.Background
			}
			if text == "" {
				text = label.Color.Text
			}
		} else {
			if bg == "" {
				bg = "#000000"
			}
			if text == "" {
				text = "#ffffff"
			}
		}
		label.SetColor(bg, text)
	}

	updated, err := repo.Update(ctx, label)
	if err != nil {
		return fmt.Errorf("failed to update label: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderLabel(updated))
	} else {
		cmd.Printf("Label updated successfully.\n")
		cmd.Printf("ID: %s\n", updated.ID)
		cmd.Printf("Name: %s\n", updated.Name)
	}

	return nil
}

// runLabelDelete handles the label delete command.
func runLabelDelete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	labelName := args[0]

	repo, err := getLabelRepository(ctx)
	if err != nil {
		return err
	}

	// Get label to verify it exists and get its ID
	label, err := repo.GetByName(ctx, labelName)
	if err != nil {
		// Try by ID if not found by name
		label, err = repo.Get(ctx, labelName)
		if err != nil {
			return fmt.Errorf("label not found: %s", labelName)
		}
	}

	// Check if system label
	if label.IsSystemLabel() {
		return fmt.Errorf("cannot delete system label: %s", label.Name)
	}

	if err := repo.Delete(ctx, label.ID); err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Label '%s' deleted successfully.\n", label.Name)
	}

	return nil
}
