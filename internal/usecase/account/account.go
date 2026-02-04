// Package account provides application use cases for account management.
package account

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/auth"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"golang.org/x/oauth2"
)

const (
	// EnvAccount is the environment variable name for specifying the default account.
	EnvAccount = "GOOG_ACCOUNT"
)

// Store defines the interface for secure credential storage.
type Store interface {
	Set(account, key string, value []byte) error
	Get(account, key string) ([]byte, error)
	Delete(account, key string) error
	List(account string) ([]string, error)
}

// AuthFlow defines the interface for OAuth authentication flow.
type AuthFlow interface {
	// Run executes the OAuth flow and returns the email and token.
	Run(ctx context.Context, scopes []string) (string, *oauth2.Token, error)
}

// Service provides account management operations.
type Service struct {
	cfg      *config.Config
	store    Store
	authFlow AuthFlow
	tokens   *auth.TokenManager
}

// NewService creates a new account service.
func NewService(cfg *config.Config, store Store, authFlow AuthFlow) *Service {
	return &Service{
		cfg:      cfg,
		store:    store,
		authFlow: authFlow,
		tokens:   auth.NewTokenManager(store),
	}
}

// Add adds a new account by running the OAuth flow.
func (s *Service) Add(ctx context.Context, alias string, scopes []string) (*account.Account, error) {
	// Check if account already exists
	_, err := s.cfg.GetAccount(alias)
	if err == nil {
		return nil, account.ErrAccountExists
	}

	// Set default scopes if none provided
	if len(scopes) == 0 {
		scopes = []string{
			auth.ScopeGmailReadonly,
			auth.ScopeCalendarReadonly,
			auth.ScopeUserInfoEmail,
			auth.ScopeOpenID,
		}
	}

	// Run OAuth flow
	email, token, err := s.authFlow.Run(ctx, scopes)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Save token to keyring
	if err := s.tokens.SaveToken(alias, token); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	// Save scopes
	if err := s.tokens.SaveScopes(alias, scopes); err != nil {
		return nil, fmt.Errorf("failed to save scopes: %w", err)
	}

	// Create account config
	accConfig := config.AccountConfig{
		Email:   email,
		Scopes:  scopes,
		AddedAt: time.Now(),
	}

	// Add to config
	s.cfg.Accounts[alias] = accConfig

	// If this is the first account, set as default
	if len(s.cfg.Accounts) == 1 {
		s.cfg.DefaultAccount = alias
	}

	// Save config
	if err := s.cfg.Save(); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	// Create domain account
	acc := account.NewAccount(alias, email)
	acc.Scopes = scopes
	acc.IsDefault = s.cfg.DefaultAccount == alias

	return acc, nil
}

// Remove removes an account and its tokens.
func (s *Service) Remove(alias string) error {
	// Check if account exists
	_, err := s.cfg.GetAccount(alias)
	if err != nil {
		return account.ErrAccountNotFound
	}

	// Delete tokens from keyring
	if err := s.tokens.DeleteToken(alias); err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}

	// Remove from config
	delete(s.cfg.Accounts, alias)

	// If this was the default account, clear it
	if s.cfg.DefaultAccount == alias {
		s.cfg.DefaultAccount = ""
		// Set a new default if accounts remain (use sorted order for deterministic behavior)
		if len(s.cfg.Accounts) > 0 {
			aliases := make([]string, 0, len(s.cfg.Accounts))
			for a := range s.cfg.Accounts {
				aliases = append(aliases, a)
			}
			sort.Strings(aliases)
			s.cfg.DefaultAccount = aliases[0]
		}
	}

	// Save config
	if err := s.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// List returns all configured accounts in sorted order by alias.
func (s *Service) List() ([]*account.Account, error) {
	// Get sorted list of aliases for deterministic ordering
	aliases := make([]string, 0, len(s.cfg.Accounts))
	for alias := range s.cfg.Accounts {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	accounts := make([]*account.Account, 0, len(s.cfg.Accounts))
	for _, alias := range aliases {
		accCfg := s.cfg.Accounts[alias]
		acc := account.NewAccount(alias, accCfg.Email)
		acc.Scopes = accCfg.Scopes
		acc.Added = accCfg.AddedAt
		acc.IsDefault = s.cfg.DefaultAccount == alias
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

// Switch sets the default account.
func (s *Service) Switch(alias string) error {
	// Check if account exists
	_, err := s.cfg.GetAccount(alias)
	if err != nil {
		return account.ErrAccountNotFound
	}

	s.cfg.DefaultAccount = alias

	// Save config
	if err := s.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Show returns the current (default) account.
func (s *Service) Show() (*account.Account, error) {
	// Try to resolve the current account
	return s.ResolveAccount("")
}

// Rename changes an account's alias.
func (s *Service) Rename(oldAlias, newAlias string) error {
	// Check if old alias exists
	accCfg, err := s.cfg.GetAccount(oldAlias)
	if err != nil {
		return account.ErrAccountNotFound
	}

	// Check if new alias already exists
	_, err = s.cfg.GetAccount(newAlias)
	if err == nil {
		return account.ErrAccountExists
	}

	// Load token with old alias
	token, err := s.tokens.LoadToken(oldAlias)
	if err != nil && err != auth.ErrTokenNotFound {
		return fmt.Errorf("failed to load token: %w", err)
	}

	// Load scopes with old alias
	// Note: We ignore the error here because scopes are optional metadata.
	// If scopes don't exist (e.g., old account without scopes saved),
	// we'll just skip copying them. The token itself is the critical data.
	scopes, _ := s.tokens.GetGrantedScopes(oldAlias)

	// Save token with new alias (if token exists)
	if token != nil {
		if err := s.tokens.SaveToken(newAlias, token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
		if scopes != nil {
			if err := s.tokens.SaveScopes(newAlias, scopes); err != nil {
				return fmt.Errorf("failed to save scopes: %w", err)
			}
		}
		// Delete old token - best effort, don't fail the rename if this errors
		// The new token is already saved, so the rename is functionally complete.
		if err := s.tokens.DeleteToken(oldAlias); err != nil {
			// Log or ignore - the old token remains but won't be used
			// since the config now points to newAlias
			_ = err
		}
	}

	// Add with new alias
	s.cfg.Accounts[newAlias] = *accCfg

	// Remove old alias
	delete(s.cfg.Accounts, oldAlias)

	// Update default if needed
	if s.cfg.DefaultAccount == oldAlias {
		s.cfg.DefaultAccount = newAlias
	}

	// Save config
	if err := s.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// ResolveAccount resolves the account to use based on:
// 1. Flag value (if provided)
// 2. GOOG_ACCOUNT environment variable
// 3. Default account in config
// 4. First account (fallback)
func (s *Service) ResolveAccount(flagValue string) (*account.Account, error) {
	var alias string

	// 1. Flag value
	if flagValue != "" {
		alias = flagValue
	}

	// 2. Environment variable
	if alias == "" {
		alias = os.Getenv(EnvAccount)
	}

	// 3. Default account
	if alias == "" {
		alias = s.cfg.DefaultAccount
	}

	// 4. First account fallback (use sorted order for deterministic behavior)
	if alias == "" && len(s.cfg.Accounts) > 0 {
		aliases := make([]string, 0, len(s.cfg.Accounts))
		for a := range s.cfg.Accounts {
			aliases = append(aliases, a)
		}
		sort.Strings(aliases)
		alias = aliases[0]
	}

	// Still no alias - no accounts configured
	if alias == "" {
		return nil, account.ErrAccountNotFound
	}

	// Get account config
	accCfg, err := s.cfg.GetAccount(alias)
	if err != nil {
		return nil, account.ErrAccountNotFound
	}

	// Create domain account
	acc := account.NewAccount(alias, accCfg.Email)
	acc.Scopes = accCfg.Scopes
	acc.Added = accCfg.AddedAt
	acc.IsDefault = s.cfg.DefaultAccount == alias

	return acc, nil
}

// GetTokenManager returns the token manager for auth operations.
func (s *Service) GetTokenManager() *auth.TokenManager {
	return s.tokens
}

// GetConfig returns the config for auth operations.
func (s *Service) GetConfig() *config.Config {
	return s.cfg
}

// Account is a type alias for the domain account for external use.
type Account = account.Account
