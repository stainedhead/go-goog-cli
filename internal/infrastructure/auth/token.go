// Package auth provides OAuth2/PKCE authentication for Google APIs.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
)

// Key names used for storing token data in the keyring.
const (
	KeyToken  = "oauth_token"
	KeyScopes = "oauth_scopes"
)

// Errors for token management.
var (
	ErrTokenNotFound = errors.New("token not found")
	ErrScopesNotSet  = errors.New("scopes not set for account")
	errKeyNotFound   = errors.New("key not found") // Internal error for mock store
)

// Store defines the interface for secure credential storage.
// This mirrors the keyring.Store interface to avoid circular dependencies.
type Store interface {
	Set(account, key string, value []byte) error
	Get(account, key string) ([]byte, error)
	Delete(account, key string) error
	List(account string) ([]string, error)
}

// TokenManager handles OAuth2 token storage and retrieval.
// It uses a keyring Store for secure token persistence.
type TokenManager struct {
	store Store
}

// NewTokenManager creates a new TokenManager with the given store.
func NewTokenManager(store Store) *TokenManager {
	return &TokenManager{
		store: store,
	}
}

// SaveToken stores an OAuth2 token for the given account.
// The token is serialized to JSON and stored securely in the keyring.
func (tm *TokenManager) SaveToken(account string, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := tm.store.Set(account, KeyToken, data); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	return nil
}

// LoadToken retrieves an OAuth2 token for the given account.
// Returns ErrTokenNotFound if no token exists for the account.
func (tm *TokenManager) LoadToken(account string) (*oauth2.Token, error) {
	data, err := tm.store.Get(account, KeyToken)
	if err != nil {
		if isKeyNotFoundError(err) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to retrieve token: %w", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// DeleteToken removes the OAuth2 token for the given account.
// This operation is idempotent - it does not error if the token doesn't exist.
func (tm *TokenManager) DeleteToken(account string) error {
	// Delete token
	if err := tm.store.Delete(account, KeyToken); err != nil && !isKeyNotFoundError(err) {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	// Also delete scopes
	if err := tm.store.Delete(account, KeyScopes); err != nil && !isKeyNotFoundError(err) {
		return fmt.Errorf("failed to delete scopes: %w", err)
	}

	return nil
}

// RefreshToken refreshes an expired OAuth2 token for the given account.
// It loads the existing token, refreshes it using the provided config, and saves the new token.
func (tm *TokenManager) RefreshToken(ctx context.Context, account string, cfg *oauth2.Config) (*oauth2.Token, error) {
	// Load existing token
	token, err := tm.LoadToken(account)
	if err != nil {
		return nil, fmt.Errorf("failed to load token for refresh: %w", err)
	}

	// Create a token source that will refresh the token
	ts := cfg.TokenSource(ctx, token)

	// Get a new token (this will refresh if expired)
	newToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Save the new token if it's different
	if newToken.AccessToken != token.AccessToken {
		if err := tm.SaveToken(account, newToken); err != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", err)
		}
	}

	return newToken, nil
}

// GetTokenSource returns an oauth2.TokenSource for the given account.
// The token source will automatically refresh the token when it expires.
func (tm *TokenManager) GetTokenSource(ctx context.Context, account string) (oauth2.TokenSource, error) {
	// Load the token
	token, err := tm.LoadToken(account)
	if err != nil {
		return nil, err
	}

	// Load scopes to create the config
	scopes, err := tm.GetGrantedScopes(account)
	if err != nil {
		// If scopes aren't stored, use an empty slice
		// The token should still work, just can't refresh
		scopes = []string{}
	}

	// Create OAuth config
	cfg := NewOAuthConfig(scopes)

	// Create a reusable token source that auto-refreshes
	ts := cfg.TokenSource(ctx, token)

	// Wrap in a ReuseTokenSource for efficiency
	return oauth2.ReuseTokenSource(token, ts), nil
}

// SaveScopes stores the granted OAuth scopes for the given account.
func (tm *TokenManager) SaveScopes(account string, scopes []string) error {
	data, err := json.Marshal(scopes)
	if err != nil {
		return fmt.Errorf("failed to marshal scopes: %w", err)
	}

	if err := tm.store.Set(account, KeyScopes, data); err != nil {
		return fmt.Errorf("failed to store scopes: %w", err)
	}

	return nil
}

// GetGrantedScopes retrieves the granted OAuth scopes for the given account.
func (tm *TokenManager) GetGrantedScopes(account string) ([]string, error) {
	data, err := tm.store.Get(account, KeyScopes)
	if err != nil {
		if isKeyNotFoundError(err) {
			return nil, ErrScopesNotSet
		}
		return nil, fmt.Errorf("failed to retrieve scopes: %w", err)
	}

	var scopes []string
	if err := json.Unmarshal(data, &scopes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scopes: %w", err)
	}

	return scopes, nil
}

// HasScope checks if the account has been granted the specified scope.
// Returns false if the account doesn't exist or scopes aren't stored.
func (tm *TokenManager) HasScope(account, scope string) bool {
	scopes, err := tm.GetGrantedScopes(account)
	if err != nil {
		return false
	}

	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// TokenInfo contains information about a stored token.
type TokenInfo struct {
	Account    string
	HasToken   bool
	IsExpired  bool
	Scopes     []string
	TokenType  string
	ExpiryTime string
}

// GetTokenInfo returns information about the token for the given account.
func (tm *TokenManager) GetTokenInfo(account string) (*TokenInfo, error) {
	info := &TokenInfo{
		Account:  account,
		HasToken: false,
	}

	// Try to load token
	token, err := tm.LoadToken(account)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return info, nil
		}
		return nil, err
	}

	info.HasToken = true
	info.TokenType = token.TokenType
	info.IsExpired = !token.Valid()

	if !token.Expiry.IsZero() {
		info.ExpiryTime = token.Expiry.String()
	}

	// Try to load scopes
	scopes, err := tm.GetGrantedScopes(account)
	if err == nil {
		info.Scopes = scopes
	}

	return info, nil
}

// isKeyNotFoundError checks if the error indicates a key was not found.
// This handles both the internal errKeyNotFound and external keyring errors.
func isKeyNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// Check for our internal error
	if errors.Is(err, errKeyNotFound) {
		return true
	}
	// Check for common "not found" error messages
	errStr := err.Error()
	return errStr == "key not found" || errStr == "secret not found in keyring"
}
