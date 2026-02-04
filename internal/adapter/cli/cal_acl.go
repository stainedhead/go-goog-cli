// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	"github.com/stainedhead/go-goog-cli/internal/adapter/repository"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/keyring"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// ACL command flags.
var (
	aclEmail   string
	aclRole    string
	aclConfirm bool
)

// aclCmd represents the acl command group.
var aclCmd = &cobra.Command{
	Use:   "acl",
	Short: "Manage calendar sharing rules",
	Long: `Manage access control list (ACL) rules for Google Calendars.

ACL rules control who can access a calendar and what permissions they have.
Available roles are: reader, writer, owner, freeBusyReader.`,
}

// aclListCmd lists ACL rules for a calendar.
var aclListCmd = &cobra.Command{
	Use:   "list <calendar-id>",
	Short: "List sharing rules for a calendar",
	Long: `List all access control rules for a specific calendar.

Shows who has access to the calendar and their permission level.`,
	Example: `  # List sharing rules for primary calendar
  goog cal acl list primary

  # List sharing rules for a specific calendar
  goog cal acl list "example@group.calendar.google.com"

  # List with JSON output
  goog cal acl list primary --format json`,
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	RunE:    runACLList,
}

// aclAddCmd adds a new ACL rule.
var aclAddCmd = &cobra.Command{
	Use:   "add <calendar-id>",
	Short: "Add a sharing rule",
	Long: `Add a new access control rule to share a calendar with a user.

You must specify the email address and role for the new rule.
Available roles: reader (view only), writer (edit events), owner (full control).`,
	Example: `  # Share calendar with a user as reader
  goog cal acl add primary --email user@example.com --role reader

  # Share calendar with a user as writer
  goog cal acl add "mywork@group.calendar.google.com" --email colleague@example.com --role writer

  # Share calendar with owner access
  goog cal acl add primary --email admin@example.com --role owner`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if aclEmail == "" {
			return fmt.Errorf("required flag \"email\" not set")
		}
		return nil
	},
	RunE: runACLAdd,
}

// aclRemoveCmd removes an ACL rule.
var aclRemoveCmd = &cobra.Command{
	Use:   "remove <calendar-id> <rule-id>",
	Short: "Remove a sharing rule",
	Long: `Remove an access control rule from a calendar.

This will revoke the user's access to the calendar.
Requires --confirm flag for safety.`,
	Example: `  # Remove a sharing rule (requires confirmation)
  goog cal acl remove primary "user:user@example.com" --confirm

  # Remove from a specific calendar
  goog cal acl remove "mywork@group.calendar.google.com" "user:colleague@example.com" --confirm`,
	Aliases: []string{"rm", "delete"},
	Args:    cobra.ExactArgs(2),
	RunE:    runACLRemove,
}

// shareCmd is a user-friendly alias for acl add.
var shareCmd = &cobra.Command{
	Use:   "share <calendar-id>",
	Short: "Share a calendar with a user",
	Long: `Share a calendar with a user by adding an access control rule.

This is a user-friendly alias for 'goog cal acl add'.
You must specify the email address and role for sharing.
Available roles: reader (view only), writer (edit events), owner (full control).`,
	Example: `  # Share calendar with a user as reader
  goog cal share primary --email user@example.com --role reader

  # Share calendar with a user as writer
  goog cal share "mywork@group.calendar.google.com" --email colleague@example.com --role writer`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if aclEmail == "" {
			return fmt.Errorf("required flag \"email\" not set")
		}
		return nil
	},
	RunE: runACLAdd,
}

// unshareCmd is a user-friendly alias for acl remove.
var unshareCmd = &cobra.Command{
	Use:   "unshare <calendar-id> <rule-id>",
	Short: "Unshare a calendar from a user",
	Long: `Remove a user's access to a calendar.

This is a user-friendly alias for 'goog cal acl remove'.
Requires --confirm flag for safety.`,
	Example: `  # Unshare calendar from a user (requires confirmation)
  goog cal unshare primary "user:user@example.com" --confirm

  # Unshare from a specific calendar
  goog cal unshare "mywork@group.calendar.google.com" "user:colleague@example.com" --confirm`,
	Args: cobra.ExactArgs(2),
	RunE: runACLRemove,
}

func init() {
	// Add acl subcommands
	aclCmd.AddCommand(aclListCmd)
	aclCmd.AddCommand(aclAddCmd)
	aclCmd.AddCommand(aclRemoveCmd)

	// Add flags to aclAddCmd
	aclAddCmd.Flags().StringVar(&aclEmail, "email", "", "email address of the user to share with (required)")
	aclAddCmd.Flags().StringVar(&aclRole, "role", "reader", "access role: reader, writer, owner, freeBusyReader")
	_ = aclAddCmd.MarkFlagRequired("email")

	// Add flags to aclRemoveCmd
	aclRemoveCmd.Flags().BoolVar(&aclConfirm, "confirm", false, "confirm removal of sharing rule")

	// Add flags to shareCmd (alias for aclAddCmd)
	shareCmd.Flags().StringVar(&aclEmail, "email", "", "email address of the user to share with (required)")
	shareCmd.Flags().StringVar(&aclRole, "role", "reader", "access role: reader, writer, owner, freeBusyReader")
	_ = shareCmd.MarkFlagRequired("email")

	// Add flags to unshareCmd (alias for aclRemoveCmd)
	unshareCmd.Flags().BoolVar(&aclConfirm, "confirm", false, "confirm removal of sharing rule")

	// Add to cal command
	calCmd.AddCommand(aclCmd)
	calCmd.AddCommand(shareCmd)
	calCmd.AddCommand(unshareCmd)
}

// getACLRepository creates an ACL repository for the current account.
func getACLRepository(ctx context.Context) (*repository.GCalACLRepository, error) {
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
		return nil, fmt.Errorf("no account found: %w (run 'goog auth login' to authenticate)", err)
	}

	// Get token source
	tokenMgr := svc.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w (run 'goog auth login' to authenticate)", err)
	}

	// Create GCal service
	gcalService, err := repository.NewGCalService(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return gcalService.ACL(), nil
}

// runACLList handles the acl list command.
func runACLList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	calendarID := args[0]

	repo, err := getACLRepository(ctx)
	if err != nil {
		return err
	}

	rules, err := repo.List(ctx, calendarID)
	if err != nil {
		return fmt.Errorf("failed to list ACL rules: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderACLRules(rules)
	cmd.Println(output)

	return nil
}

// runACLAdd handles the acl add and share commands.
func runACLAdd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	calendarID := args[0]

	// Validate role
	validRoles := map[string]bool{
		calendar.AccessRoleOwner:          true,
		calendar.AccessRoleWriter:         true,
		calendar.AccessRoleReader:         true,
		calendar.AccessRoleFreeBusyReader: true,
	}
	if !validRoles[aclRole] {
		return fmt.Errorf("invalid role %q: must be one of reader, writer, owner, freeBusyReader", aclRole)
	}

	repo, err := getACLRepository(ctx)
	if err != nil {
		return err
	}

	// Create new ACL rule
	scope := calendar.NewUserACLScope(aclEmail)
	rule := calendar.NewACLRule(scope, aclRole)

	created, err := repo.Insert(ctx, calendarID, rule)
	if err != nil {
		return fmt.Errorf("failed to add ACL rule: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderACLRule(created))
	} else {
		cmd.Printf("Sharing rule added successfully.\n")
		cmd.Printf("ID: %s\n", created.ID)
		cmd.Printf("Email: %s\n", aclEmail)
		cmd.Printf("Role: %s\n", created.Role)
	}

	return nil
}

// runACLRemove handles the acl remove and unshare commands.
func runACLRemove(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	calendarID := args[0]
	ruleID := args[1]

	// Require confirmation
	if !aclConfirm {
		return fmt.Errorf("removal requires --confirm flag")
	}

	repo, err := getACLRepository(ctx)
	if err != nil {
		return err
	}

	// Get the rule first to display info
	rule, err := repo.Get(ctx, calendarID, ruleID)
	if err != nil {
		return fmt.Errorf("ACL rule not found: %s", ruleID)
	}

	if err := repo.Delete(ctx, calendarID, ruleID); err != nil {
		return fmt.Errorf("failed to remove ACL rule: %w", err)
	}

	if !quietFlag {
		scopeValue := ""
		if rule.Scope != nil {
			scopeValue = rule.Scope.Value
		}
		cmd.Printf("Sharing rule '%s' removed successfully.\n", ruleID)
		if scopeValue != "" {
			cmd.Printf("User %s no longer has access.\n", scopeValue)
		}
	}

	return nil
}
