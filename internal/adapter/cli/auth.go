// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/auth"
)

var (
	authScopes []string
)

// authCmd represents the auth command group.
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication with Google",
	Long: `Manage authentication with Google APIs.

The auth commands handle OAuth2 authentication, including login,
logout, and checking authentication status.`,
}

// authLoginCmd handles OAuth login.
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Google",
	Long: `Start the OAuth2 authentication flow.

This opens a browser window for Google authentication. After
authenticating, the tokens are stored securely in the system keyring.

By default, the following scopes are requested:
  - Gmail (readonly)
  - Calendar (readonly)
  - User email info

Use --scopes to request specific scopes.`,
	Example: `  # Login with default scopes
  goog auth login

  # Login with specific scopes
  goog auth login --scopes gmail.modify,calendar

  # Login and add as a named account
  goog auth login --account work`,
	RunE: runAuthLogin,
}

// authLogoutCmd handles logout.
var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove credentials for the current account",
	Long: `Remove stored credentials for the current account.

This deletes the OAuth tokens from the keyring but does not
revoke the tokens on Google's servers.`,
	Example: `  # Logout from current account
  goog auth logout

  # Logout from specific account
  goog auth logout --account work`,
	RunE: runAuthLogout,
}

// authStatusCmd shows authentication status.
var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long: `Display the current authentication status.

Shows the current account, email, granted scopes, and token
expiry information.`,
	Example: `  # Show status for current account
  goog auth status

  # Show status for specific account
  goog auth status --account work`,
	RunE: runAuthStatus,
}

// authRefreshCmd forces token refresh.
var authRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Force refresh the access token",
	Long: `Force a refresh of the OAuth access token.

This can be useful if you're experiencing authentication issues
or want to ensure you have a fresh token.`,
	Example: `  # Refresh current account token
  goog auth refresh

  # Refresh specific account token
  goog auth refresh --account work`,
	RunE: runAuthRefresh,
}

func init() {
	// Add auth subcommands
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authRefreshCmd)

	// Login flags
	authLoginCmd.Flags().StringSliceVar(&authScopes, "scopes", nil, "OAuth scopes to request (comma-separated)")

	// Add to root
	rootCmd.AddCommand(authCmd)
}

// runAuthLogin handles the auth login command.
func runAuthLogin(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get account service using dependency injection
	svc := getAccountServiceFromDeps()

	// Determine alias
	alias := accountFlag
	if alias == "" {
		alias = "default"
	}

	// Parse scopes
	scopes := parseScopes(authScopes)

	// Add account
	acc, err := svc.Add(ctx, alias, scopes)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	cmd.Printf("Successfully logged in as %s\n", acc.Email)
	cmd.Printf("Account alias: %s\n", acc.Alias)
	if acc.IsDefault {
		cmd.Println("This account is set as the default.")
	}

	return nil
}

// runAuthLogout handles the auth logout command.
func runAuthLogout(cmd *cobra.Command, args []string) error {
	// Get account service using dependency injection
	svc := getAccountServiceFromDeps()

	// Resolve account
	acc, err := svc.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("no account found: %w", err)
	}

	// Remove account
	if err := svc.Remove(acc.Alias); err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	cmd.Printf("Successfully logged out from %s (%s)\n", acc.Alias, acc.Email)
	return nil
}

// runAuthStatus handles the auth status command.
func runAuthStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get account service using dependency injection
	svc := getAccountServiceFromDeps()

	// Resolve account
	acc, err := svc.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("no account found: %w", err)
	}

	// Display status
	cmd.Printf("Account:     %s\n", acc.Alias)
	cmd.Printf("Email:       %s\n", acc.Email)
	cmd.Printf("Default:     %v\n", acc.IsDefault)
	cmd.Printf("Added:       %s\n", acc.Added.Format(time.RFC3339))

	// Try to get token source to verify token status
	tokenMgr := svc.GetTokenManager()
	_, err = tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		cmd.Printf("Token:       Not found\n")
		cmd.Println("Status:      NOT AUTHENTICATED")
	} else {
		cmd.Printf("Token:       Valid\n")
		cmd.Println("Status:      ACTIVE")
	}

	if len(acc.Scopes) > 0 {
		cmd.Println("Scopes:")
		for _, scope := range acc.Scopes {
			cmd.Printf("  - %s\n", scope)
		}
	}

	return nil
}

// runAuthRefresh handles the auth refresh command.
func runAuthRefresh(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get account service using dependency injection
	svc := getAccountServiceFromDeps()

	// Resolve account
	acc, err := svc.ResolveAccount(accountFlag)
	if err != nil {
		return fmt.Errorf("no account found: %w", err)
	}

	// Get token manager and force a token refresh by getting a new token source
	tokenMgr := svc.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Force a token fetch which will refresh if expired
	token, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	cmd.Printf("Successfully refreshed token for %s\n", acc.Alias)
	if !token.Expiry.IsZero() {
		cmd.Printf("New expiry: %s\n", token.Expiry.Format(time.RFC3339))
	}

	return nil
}

// parseScopes converts scope shorthand to full scope URLs.
func parseScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}

	scopeMap := map[string]string{
		"gmail":             auth.ScopeGmailReadonly,
		"gmail.readonly":    auth.ScopeGmailReadonly,
		"gmail.send":        auth.ScopeGmailSend,
		"gmail.modify":      auth.ScopeGmailModify,
		"gmail.compose":     auth.ScopeGmailCompose,
		"gmail.labels":      auth.ScopeGmailLabels,
		"calendar":          auth.ScopeCalendarReadonly,
		"calendar.readonly": auth.ScopeCalendarReadonly,
		"calendar.events":   auth.ScopeCalendarEvents,
		"calendar.full":     auth.ScopeCalendar,
		"drive":             auth.ScopeDriveReadonly,
		"drive.readonly":    auth.ScopeDriveReadonly,
		"drive.file":        auth.ScopeDriveFile,
		"drive.full":        auth.ScopeDrive,
		"email":             auth.ScopeUserInfoEmail,
		"profile":           auth.ScopeUserInfoProfile,
		"openid":            auth.ScopeOpenID,
	}

	result := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(strings.ToLower(scope))
		if fullScope, ok := scopeMap[scope]; ok {
			result = append(result, fullScope)
		} else if strings.HasPrefix(scope, "https://") {
			// Already a full URL
			result = append(result, scope)
		} else {
			// Treat as-is, let Google validate
			result = append(result, scope)
		}
	}

	// Always include email and openid for user identification
	hasEmail := false
	hasOpenID := false
	for _, s := range result {
		if s == auth.ScopeUserInfoEmail {
			hasEmail = true
		}
		if s == auth.ScopeOpenID {
			hasOpenID = true
		}
	}
	if !hasEmail {
		result = append(result, auth.ScopeUserInfoEmail)
	}
	if !hasOpenID {
		result = append(result, auth.ScopeOpenID)
	}

	return result
}

// getEnvWithDefault returns the environment variable value or default.
func getEnvWithDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
