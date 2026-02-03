// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/keyring"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// accountCmd represents the account command group.
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage Google accounts",
	Long: `Manage multiple Google accounts.

The account commands allow you to add, remove, list, and switch
between different Google accounts.`,
}

// accountListCmd lists all accounts.
var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured accounts",
	Long: `List all configured Google accounts.

Displays a table showing the alias, email, default status,
and when each account was added.`,
	Aliases: []string{"ls"},
	RunE:    runAccountList,
}

// accountAddCmd adds a new account.
var accountAddCmd = &cobra.Command{
	Use:   "add [alias]",
	Short: "Add a new Google account",
	Long: `Add a new Google account.

This runs the OAuth authentication flow and stores the
credentials under the specified alias. If no alias is
provided, the email address will be used.`,
	Example: `  # Add account with alias
  goog account add work

  # Add account with specific scopes
  goog account add personal --scopes gmail.modify,calendar`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAccountAdd,
}

// accountRemoveCmd removes an account.
var accountRemoveCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove a Google account",
	Long: `Remove a Google account.

This removes the account configuration and deletes the
stored credentials from the keyring.`,
	Example: `  goog account remove work`,
	Aliases: []string{"rm", "delete"},
	Args:    cobra.ExactArgs(1),
	RunE:    runAccountRemove,
}

// accountSwitchCmd switches the default account.
var accountSwitchCmd = &cobra.Command{
	Use:   "switch <alias>",
	Short: "Set the default account",
	Long: `Set the default Google account.

The default account is used when no --account flag is specified.`,
	Example: `  goog account switch personal`,
	Aliases: []string{"use", "default"},
	Args:    cobra.ExactArgs(1),
	RunE:    runAccountSwitch,
}

// accountShowCmd shows the current account.
var accountShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current account",
	Long: `Show details of the current account.

Displays information about the currently active account,
including alias, email, scopes, and token status.`,
	Aliases: []string{"current", "info"},
	RunE:    runAccountShow,
}

// accountRenameCmd renames an account.
var accountRenameCmd = &cobra.Command{
	Use:   "rename <old-alias> <new-alias>",
	Short: "Rename an account",
	Long: `Rename an account alias.

This changes the alias used to refer to the account while
preserving all credentials and configuration.`,
	Example: `  goog account rename work office`,
	Aliases: []string{"mv"},
	Args:    cobra.ExactArgs(2),
	RunE:    runAccountRename,
}

var accountAddScopes []string

func init() {
	// Add account subcommands
	accountCmd.AddCommand(accountListCmd)
	accountCmd.AddCommand(accountAddCmd)
	accountCmd.AddCommand(accountRemoveCmd)
	accountCmd.AddCommand(accountSwitchCmd)
	accountCmd.AddCommand(accountShowCmd)
	accountCmd.AddCommand(accountRenameCmd)

	// Add flags
	accountAddCmd.Flags().StringSliceVar(&accountAddScopes, "scopes", nil, "OAuth scopes to request")

	// Add to root
	rootCmd.AddCommand(accountCmd)
}

// runAccountList handles the account list command.
func runAccountList(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create account service
	svc := accountuc.NewService(cfg, store, nil)

	// List accounts
	accounts, err := svc.List()
	if err != nil {
		return fmt.Errorf("failed to list accounts: %w", err)
	}

	if len(accounts) == 0 {
		cmd.Println("No accounts configured.")
		cmd.Println("Run 'goog auth login' or 'goog account add' to add an account.")
		return nil
	}

	// Output based on format
	switch formatFlag {
	case "json":
		return outputAccountsJSON(cmd, accounts)
	case "plain":
		return outputAccountsPlain(cmd, accounts)
	default:
		return outputAccountsTable(cmd, accounts)
	}
}

// runAccountAdd handles the account add command.
func runAccountAdd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create OAuth flow
	flow := accountuc.NewDefaultOAuthFlow()

	// Create account service
	svc := accountuc.NewService(cfg, store, flow)

	// Get alias
	alias := "default"
	if len(args) > 0 {
		alias = args[0]
	}

	// Parse scopes
	scopes := parseScopes(accountAddScopes)

	// Add account
	acc, err := svc.Add(ctx, alias, scopes)
	if err != nil {
		return fmt.Errorf("failed to add account: %w", err)
	}

	cmd.Printf("Successfully added account '%s' (%s)\n", acc.Alias, acc.Email)
	if acc.IsDefault {
		cmd.Println("This account is set as the default.")
	}

	return nil
}

// runAccountRemove handles the account remove command.
func runAccountRemove(cmd *cobra.Command, args []string) error {
	alias := args[0]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create account service
	svc := accountuc.NewService(cfg, store, nil)

	// Get account info before removal
	acc, err := svc.ResolveAccount(alias)
	if err != nil {
		return fmt.Errorf("account not found: %s", alias)
	}

	// Remove account
	if err := svc.Remove(alias); err != nil {
		return fmt.Errorf("failed to remove account: %w", err)
	}

	cmd.Printf("Successfully removed account '%s' (%s)\n", acc.Alias, acc.Email)
	return nil
}

// runAccountSwitch handles the account switch command.
func runAccountSwitch(cmd *cobra.Command, args []string) error {
	alias := args[0]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create account service
	svc := accountuc.NewService(cfg, store, nil)

	// Switch account
	if err := svc.Switch(alias); err != nil {
		return fmt.Errorf("failed to switch account: %w", err)
	}

	// Get account info
	acc, _ := svc.ResolveAccount(alias)

	cmd.Printf("Switched to account '%s'", alias)
	if acc != nil {
		cmd.Printf(" (%s)", acc.Email)
	}
	cmd.Println()

	return nil
}

// runAccountShow handles the account show command.
func runAccountShow(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create account service
	svc := accountuc.NewService(cfg, store, nil)

	// Resolve current account
	acc, err := svc.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("no account found: %w", err)
	}

	// Get token info
	tokenMgr := svc.GetTokenManager()
	tokenInfo, _ := tokenMgr.GetTokenInfo(acc.Alias)

	// Display account info
	cmd.Printf("Alias:       %s\n", acc.Alias)
	cmd.Printf("Email:       %s\n", acc.Email)
	cmd.Printf("Default:     %v\n", acc.IsDefault)
	cmd.Printf("Added:       %s\n", acc.Added.Format(time.RFC3339))

	if tokenInfo != nil && tokenInfo.HasToken {
		cmd.Printf("Token:       Valid\n")
		if tokenInfo.IsExpired {
			cmd.Printf("Status:      EXPIRED\n")
		} else {
			cmd.Printf("Status:      ACTIVE\n")
		}
		if tokenInfo.ExpiryTime != "" {
			cmd.Printf("Expires:     %s\n", tokenInfo.ExpiryTime)
		}
	} else {
		cmd.Printf("Token:       Not found\n")
	}

	if len(acc.Scopes) > 0 {
		cmd.Println("Scopes:")
		for _, scope := range acc.Scopes {
			cmd.Printf("  - %s\n", scope)
		}
	}

	return nil
}

// runAccountRename handles the account rename command.
func runAccountRename(cmd *cobra.Command, args []string) error {
	oldAlias := args[0]
	newAlias := args[1]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create account service
	svc := accountuc.NewService(cfg, store, nil)

	// Rename account
	if err := svc.Rename(oldAlias, newAlias); err != nil {
		return fmt.Errorf("failed to rename account: %w", err)
	}

	cmd.Printf("Successfully renamed account '%s' to '%s'\n", oldAlias, newAlias)
	return nil
}

// outputAccountsTable outputs accounts in table format.
func outputAccountsTable(cmd *cobra.Command, accounts []*accountuc.Account) error {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ALIAS\tEMAIL\tDEFAULT\tADDED")
	for _, acc := range accounts {
		defaultStr := ""
		if acc.IsDefault {
			defaultStr = "*"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			acc.Alias,
			acc.Email,
			defaultStr,
			acc.Added.Format("2006-01-02"),
		)
	}
	return w.Flush()
}

// outputAccountsPlain outputs accounts in plain format.
func outputAccountsPlain(cmd *cobra.Command, accounts []*accountuc.Account) error {
	for _, acc := range accounts {
		defaultStr := ""
		if acc.IsDefault {
			defaultStr = " (default)"
		}
		cmd.Printf("%s: %s%s\n", acc.Alias, acc.Email, defaultStr)
	}
	return nil
}

// outputAccountsJSON outputs accounts in JSON format.
func outputAccountsJSON(cmd *cobra.Command, accounts []*accountuc.Account) error {
	type accountJSON struct {
		Alias     string   `json:"alias"`
		Email     string   `json:"email"`
		IsDefault bool     `json:"is_default"`
		Added     string   `json:"added"`
		Scopes    []string `json:"scopes,omitempty"`
	}

	result := make([]accountJSON, len(accounts))
	for i, acc := range accounts {
		result[i] = accountJSON{
			Alias:     acc.Alias,
			Email:     acc.Email,
			IsDefault: acc.IsDefault,
			Added:     acc.Added.Format(time.RFC3339),
			Scopes:    acc.Scopes,
		}
	}

	// Use standard library for JSON output
	encoder := os.Stdout
	fmt.Fprintf(encoder, "[\n")
	for i, a := range result {
		fmt.Fprintf(encoder, "  {\n")
		fmt.Fprintf(encoder, "    \"alias\": %q,\n", a.Alias)
		fmt.Fprintf(encoder, "    \"email\": %q,\n", a.Email)
		fmt.Fprintf(encoder, "    \"is_default\": %v,\n", a.IsDefault)
		fmt.Fprintf(encoder, "    \"added\": %q", a.Added)
		if len(a.Scopes) > 0 {
			fmt.Fprintf(encoder, ",\n    \"scopes\": [")
			for j, s := range a.Scopes {
				if j > 0 {
					fmt.Fprintf(encoder, ", ")
				}
				fmt.Fprintf(encoder, "%q", s)
			}
			fmt.Fprintf(encoder, "]")
		}
		fmt.Fprintf(encoder, "\n  }")
		if i < len(result)-1 {
			fmt.Fprintf(encoder, ",")
		}
		fmt.Fprintf(encoder, "\n")
	}
	fmt.Fprintf(encoder, "]\n")

	return nil
}

// Account type alias for domain account to avoid import conflicts in output functions.
type Account = accountuc.Account
