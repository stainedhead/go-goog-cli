// Package account provides domain entities and interfaces for account management.
package account

import (
	"strings"
	"time"
)

// Account represents a Google account configured for use with the CLI.
type Account struct {
	Alias     string
	Email     string
	Scopes    []string
	Added     time.Time
	LastUsed  time.Time
	IsDefault bool
}

// NewAccount creates a new Account with the given alias and email.
// The Added timestamp is set to the current time.
func NewAccount(alias, email string) *Account {
	return &Account{
		Alias:     alias,
		Email:     email,
		Scopes:    make([]string, 0),
		Added:     time.Now(),
		IsDefault: false,
	}
}

// Validate checks that the account has valid alias and email values.
// Returns ErrInvalidAlias if alias is empty or whitespace-only.
// Returns ErrInvalidEmail if email format is invalid.
func (a *Account) Validate() error {
	if strings.TrimSpace(a.Alias) == "" {
		return ErrInvalidAlias
	}

	if err := validateEmail(a.Email); err != nil {
		return err
	}

	return nil
}

// validateEmail checks if the email address has a valid format.
func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}

	// Check for spaces
	if strings.Contains(email, " ") {
		return ErrInvalidEmail
	}

	// Split by @ to validate structure
	parts := strings.Split(email, "@")

	// Must have exactly one @
	if len(parts) != 2 {
		return ErrInvalidEmail
	}

	local := parts[0]
	domain := parts[1]

	// Both parts must be non-empty
	if local == "" || domain == "" {
		return ErrInvalidEmail
	}

	return nil
}

// HasScope returns true if the account has the specified scope.
func (a *Account) HasScope(scope string) bool {
	for _, s := range a.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// AddScope adds a scope to the account if it doesn't already exist.
func (a *Account) AddScope(scope string) {
	if !a.HasScope(scope) {
		a.Scopes = append(a.Scopes, scope)
	}
}

// RemoveScope removes a scope from the account if it exists.
func (a *Account) RemoveScope(scope string) {
	for i, s := range a.Scopes {
		if s == scope {
			// Remove by swapping with last element and truncating
			a.Scopes[i] = a.Scopes[len(a.Scopes)-1]
			a.Scopes = a.Scopes[:len(a.Scopes)-1]
			return
		}
	}
}
