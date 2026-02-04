// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"

	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/keyring"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
	"golang.org/x/oauth2"
)

// getAccountService creates an account service with config and keyring store.
// It returns the service and any error encountered during initialization.
func getAccountService() (*accountuc.Service, error) {
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

	return svc, nil
}

// getResolvedAccount creates an account service and resolves the current account.
// It uses accountFlag to determine which account to use.
// Returns the service, resolved account, and any error encountered.
func getResolvedAccount() (*accountuc.Service, *accountuc.Account, error) {
	svc, err := getAccountService()
	if err != nil {
		return nil, nil, err
	}

	// Resolve account
	acc, err := svc.ResolveAccount(accountFlag)
	if err != nil {
		return nil, nil, fmt.Errorf("no account found: %w (run 'goog auth login' to authenticate)", err)
	}

	return svc, acc, nil
}

// getTokenSource creates an account service, resolves the account, and returns a token source.
// This is the most common operation needed by repository factory functions.
func getTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	svc, acc, err := getResolvedAccount()
	if err != nil {
		return nil, err
	}

	// Get token source
	tokenMgr := svc.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w (run 'goog auth login' to authenticate)", err)
	}

	return tokenSource, nil
}

// getTokenSourceWithEmail creates an account service, resolves the account, and returns
// both a token source and the account's email address. This is useful for operations
// that need to know the sender's email.
func getTokenSourceWithEmail(ctx context.Context) (oauth2.TokenSource, string, error) {
	svc, acc, err := getResolvedAccount()
	if err != nil {
		return nil, "", err
	}

	// Get token source
	tokenMgr := svc.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get token: %w (run 'goog auth login' to authenticate)", err)
	}

	return tokenSource, acc.Email, nil
}
