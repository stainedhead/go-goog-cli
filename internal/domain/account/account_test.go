// Package account provides domain entities and interfaces for account management.
package account

import (
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name          string
		alias         string
		email         string
		wantAlias     string
		wantEmail     string
		wantIsDefault bool
	}{
		{
			name:          "creates account with valid alias and email",
			alias:         "work",
			email:         "user@example.com",
			wantAlias:     "work",
			wantEmail:     "user@example.com",
			wantIsDefault: false,
		},
		{
			name:          "creates account with gmail address",
			alias:         "personal",
			email:         "user@gmail.com",
			wantAlias:     "personal",
			wantEmail:     "user@gmail.com",
			wantIsDefault: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			account := NewAccount(tt.alias, tt.email)
			after := time.Now()

			if account == nil {
				t.Fatal("NewAccount returned nil")
			}

			if account.Alias != tt.wantAlias {
				t.Errorf("Alias = %q, want %q", account.Alias, tt.wantAlias)
			}

			if account.Email != tt.wantEmail {
				t.Errorf("Email = %q, want %q", account.Email, tt.wantEmail)
			}

			if account.IsDefault != tt.wantIsDefault {
				t.Errorf("IsDefault = %v, want %v", account.IsDefault, tt.wantIsDefault)
			}

			if account.Added.Before(before) || account.Added.After(after) {
				t.Errorf("Added time %v not in expected range [%v, %v]", account.Added, before, after)
			}

			if account.Scopes == nil {
				t.Error("Scopes should be initialized to empty slice, got nil")
			}

			if len(account.Scopes) != 0 {
				t.Errorf("Scopes length = %d, want 0", len(account.Scopes))
			}
		})
	}
}

func TestAccountValidation(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		email     string
		wantError error
	}{
		{
			name:      "valid account",
			alias:     "work",
			email:     "user@example.com",
			wantError: nil,
		},
		{
			name:      "valid account with subdomain",
			alias:     "corp",
			email:     "user@mail.corp.example.com",
			wantError: nil,
		},
		{
			name:      "valid account with plus addressing",
			alias:     "tagged",
			email:     "user+tag@example.com",
			wantError: nil,
		},
		{
			name:      "empty alias",
			alias:     "",
			email:     "user@example.com",
			wantError: ErrInvalidAlias,
		},
		{
			name:      "whitespace only alias",
			alias:     "   ",
			email:     "user@example.com",
			wantError: ErrInvalidAlias,
		},
		{
			name:      "empty email",
			alias:     "work",
			email:     "",
			wantError: ErrInvalidEmail,
		},
		{
			name:      "email without @",
			alias:     "work",
			email:     "userexample.com",
			wantError: ErrInvalidEmail,
		},
		{
			name:      "email without domain",
			alias:     "work",
			email:     "user@",
			wantError: ErrInvalidEmail,
		},
		{
			name:      "email without local part",
			alias:     "work",
			email:     "@example.com",
			wantError: ErrInvalidEmail,
		},
		{
			name:      "email with multiple @",
			alias:     "work",
			email:     "user@@example.com",
			wantError: ErrInvalidEmail,
		},
		{
			name:      "email with spaces",
			alias:     "work",
			email:     "user @example.com",
			wantError: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := NewAccount(tt.alias, tt.email)
			err := account.Validate()

			if tt.wantError == nil {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Errorf("Validate() error = nil, want %v", tt.wantError)
				return
			}

			if err != tt.wantError {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantError)
			}
		})
	}
}

func TestAccountScopes(t *testing.T) {
	t.Run("AddScope adds new scope", func(t *testing.T) {
		account := NewAccount("work", "user@example.com")
		scope := "https://www.googleapis.com/auth/gmail.readonly"

		account.AddScope(scope)

		if len(account.Scopes) != 1 {
			t.Fatalf("Scopes length = %d, want 1", len(account.Scopes))
		}

		if account.Scopes[0] != scope {
			t.Errorf("Scopes[0] = %q, want %q", account.Scopes[0], scope)
		}
	})

	t.Run("AddScope does not add duplicate", func(t *testing.T) {
		account := NewAccount("work", "user@example.com")
		scope := "https://www.googleapis.com/auth/gmail.readonly"

		account.AddScope(scope)
		account.AddScope(scope)

		if len(account.Scopes) != 1 {
			t.Errorf("Scopes length = %d, want 1 (no duplicate)", len(account.Scopes))
		}
	})

	t.Run("AddScope adds multiple different scopes", func(t *testing.T) {
		account := NewAccount("work", "user@example.com")
		scope1 := "https://www.googleapis.com/auth/gmail.readonly"
		scope2 := "https://www.googleapis.com/auth/calendar.readonly"

		account.AddScope(scope1)
		account.AddScope(scope2)

		if len(account.Scopes) != 2 {
			t.Fatalf("Scopes length = %d, want 2", len(account.Scopes))
		}
	})

	t.Run("RemoveScope removes existing scope", func(t *testing.T) {
		account := NewAccount("work", "user@example.com")
		scope := "https://www.googleapis.com/auth/gmail.readonly"

		account.AddScope(scope)
		account.RemoveScope(scope)

		if len(account.Scopes) != 0 {
			t.Errorf("Scopes length = %d, want 0", len(account.Scopes))
		}
	})

	t.Run("RemoveScope does nothing for non-existent scope", func(t *testing.T) {
		account := NewAccount("work", "user@example.com")
		scope1 := "https://www.googleapis.com/auth/gmail.readonly"
		scope2 := "https://www.googleapis.com/auth/calendar.readonly"

		account.AddScope(scope1)
		account.RemoveScope(scope2)

		if len(account.Scopes) != 1 {
			t.Errorf("Scopes length = %d, want 1", len(account.Scopes))
		}

		if account.Scopes[0] != scope1 {
			t.Errorf("Scopes[0] = %q, want %q", account.Scopes[0], scope1)
		}
	})

	t.Run("RemoveScope removes correct scope from multiple", func(t *testing.T) {
		account := NewAccount("work", "user@example.com")
		scope1 := "https://www.googleapis.com/auth/gmail.readonly"
		scope2 := "https://www.googleapis.com/auth/calendar.readonly"
		scope3 := "https://www.googleapis.com/auth/gmail.send"

		account.AddScope(scope1)
		account.AddScope(scope2)
		account.AddScope(scope3)
		account.RemoveScope(scope2)

		if len(account.Scopes) != 2 {
			t.Fatalf("Scopes length = %d, want 2", len(account.Scopes))
		}

		for _, s := range account.Scopes {
			if s == scope2 {
				t.Errorf("scope2 should have been removed, but found in Scopes")
			}
		}
	})
}

func TestAccountHasScope(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []string
		checkFor string
		wantHas  bool
	}{
		{
			name:     "has scope when present",
			scopes:   []string{"https://www.googleapis.com/auth/gmail.readonly"},
			checkFor: "https://www.googleapis.com/auth/gmail.readonly",
			wantHas:  true,
		},
		{
			name:     "does not have scope when absent",
			scopes:   []string{"https://www.googleapis.com/auth/gmail.readonly"},
			checkFor: "https://www.googleapis.com/auth/calendar.readonly",
			wantHas:  false,
		},
		{
			name:     "does not have scope when empty",
			scopes:   []string{},
			checkFor: "https://www.googleapis.com/auth/gmail.readonly",
			wantHas:  false,
		},
		{
			name: "has scope among multiple",
			scopes: []string{
				"https://www.googleapis.com/auth/gmail.readonly",
				"https://www.googleapis.com/auth/calendar.readonly",
				"https://www.googleapis.com/auth/gmail.send",
			},
			checkFor: "https://www.googleapis.com/auth/calendar.readonly",
			wantHas:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := NewAccount("work", "user@example.com")
			for _, s := range tt.scopes {
				account.AddScope(s)
			}

			got := account.HasScope(tt.checkFor)
			if got != tt.wantHas {
				t.Errorf("HasScope(%q) = %v, want %v", tt.checkFor, got, tt.wantHas)
			}
		})
	}
}
