// Package account provides domain entities and interfaces for account management.
package account

import "errors"

// Domain errors for account operations.
var (
	// ErrAccountNotFound is returned when an account with the specified alias does not exist.
	ErrAccountNotFound = errors.New("account not found")

	// ErrAccountExists is returned when attempting to create an account with an alias that already exists.
	ErrAccountExists = errors.New("account already exists")

	// ErrInvalidAlias is returned when the account alias is empty or invalid.
	ErrInvalidAlias = errors.New("invalid alias: alias cannot be empty")

	// ErrInvalidEmail is returned when the email address format is invalid.
	ErrInvalidEmail = errors.New("invalid email: must be a valid email address")
)

// Repository defines the interface for account persistence operations.
type Repository interface {
	// Save persists an account. Returns ErrAccountExists if an account with
	// the same alias already exists.
	Save(account *Account) error

	// Get retrieves an account by alias. Returns ErrAccountNotFound if the
	// account does not exist.
	Get(alias string) (*Account, error)

	// List returns all configured accounts.
	List() ([]*Account, error)

	// Delete removes an account by alias. Returns ErrAccountNotFound if the
	// account does not exist.
	Delete(alias string) error

	// SetDefault marks the account with the given alias as the default.
	// Returns ErrAccountNotFound if the account does not exist.
	SetDefault(alias string) error

	// GetDefault returns the default account. Returns ErrAccountNotFound if
	// no default account is configured.
	GetDefault() (*Account, error)
}
