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
// Deprecated: Use GetDependencies().AccountService instead for testability.
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
// Deprecated: Use GetDependencies().AccountService.ResolveAccount() instead for testability.
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
// Deprecated: Use getTokenSourceFromDeps() instead for testability.
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
// Deprecated: Use getTokenSourceWithEmailFromDeps() instead for testability.
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

// =============================================================================
// Dependency Injection-based Factory Functions
// =============================================================================

// getTokenSourceFromDeps resolves the account and returns a token source using injected dependencies.
// This function supports dependency injection for testing.
func getTokenSourceFromDeps(ctx context.Context) (oauth2.TokenSource, error) {
	deps := GetDependencies()

	acc, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return nil, fmt.Errorf("no account found: %w (run 'goog auth login' to authenticate)", err)
	}

	tokenMgr := deps.AccountService.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w (run 'goog auth login' to authenticate)", err)
	}

	return tokenSource, nil
}

// getTokenSourceWithEmailFromDeps resolves the account and returns a token source
// along with the account's email address using injected dependencies.
func getTokenSourceWithEmailFromDeps(ctx context.Context) (oauth2.TokenSource, string, error) {
	deps := GetDependencies()

	acc, err := deps.AccountService.ResolveAccount(accountFlag)
	if err != nil {
		return nil, "", fmt.Errorf("no account found: %w (run 'goog auth login' to authenticate)", err)
	}

	tokenMgr := deps.AccountService.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get token: %w (run 'goog auth login' to authenticate)", err)
	}

	return tokenSource, acc.Email, nil
}

// getMessageRepositoryFromDeps creates a message repository using injected dependencies.
func getMessageRepositoryFromDeps(ctx context.Context) (MessageRepository, string, error) {
	tokenSource, email, err := getTokenSourceWithEmailFromDeps(ctx)
	if err != nil {
		return nil, "", err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewMessageRepository(ctx, tokenSource)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create message repository: %w", err)
	}

	return repo, email, nil
}

// getDraftRepositoryFromDeps creates a draft repository using injected dependencies.
func getDraftRepositoryFromDeps(ctx context.Context) (DraftRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewDraftRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create draft repository: %w", err)
	}

	return repo, nil
}

// getThreadRepositoryFromDeps creates a thread repository using injected dependencies.
func getThreadRepositoryFromDeps(ctx context.Context) (ThreadRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewThreadRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create thread repository: %w", err)
	}

	return repo, nil
}

// getLabelRepositoryFromDeps creates a label repository using injected dependencies.
func getLabelRepositoryFromDeps(ctx context.Context) (LabelRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewLabelRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create label repository: %w", err)
	}

	return repo, nil
}

// getEventRepositoryFromDeps creates an event repository using injected dependencies.
func getEventRepositoryFromDeps(ctx context.Context) (EventRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewEventRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create event repository: %w", err)
	}

	return repo, nil
}

// getCalendarRepositoryFromDeps creates a calendar repository using injected dependencies.
func getCalendarRepositoryFromDeps(ctx context.Context) (CalendarRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewCalendarRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar repository: %w", err)
	}

	return repo, nil
}

// getACLRepositoryFromDeps creates an ACL repository using injected dependencies.
func getACLRepositoryFromDeps(ctx context.Context) (ACLRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewACLRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACL repository: %w", err)
	}

	return repo, nil
}

// getFreeBusyRepositoryFromDeps creates a free/busy repository using injected dependencies.
func getFreeBusyRepositoryFromDeps(ctx context.Context) (FreeBusyRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewFreeBusyRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create free/busy repository: %w", err)
	}

	return repo, nil
}

// getContactRepositoryFromDeps creates a contact repository using injected dependencies.
func getContactRepositoryFromDeps(ctx context.Context) (ContactRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewContactRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact repository: %w", err)
	}

	return repo, nil
}

// getContactGroupRepositoryFromDeps creates a contact group repository using injected dependencies.
func getContactGroupRepositoryFromDeps(ctx context.Context) (ContactGroupRepository, error) {
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		return nil, err
	}

	deps := GetDependencies()
	repo, err := deps.RepoFactory.NewContactGroupRepository(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact group repository: %w", err)
	}

	return repo, nil
}

// getAccountServiceFromDeps returns the account service using injected dependencies.
// This function supports dependency injection for testing.
func getAccountServiceFromDeps() AccountService {
	deps := GetDependencies()
	return deps.AccountService
}
