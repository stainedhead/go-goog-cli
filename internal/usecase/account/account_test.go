// Package account provides application use cases for account management.
package account

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"golang.org/x/oauth2"
)

// mockStore implements auth.Store for testing.
type mockStore struct {
	data map[string]map[string][]byte
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[string]map[string][]byte),
	}
}

func (m *mockStore) Set(account, key string, value []byte) error {
	if m.data[account] == nil {
		m.data[account] = make(map[string][]byte)
	}
	m.data[account][key] = value
	return nil
}

func (m *mockStore) Get(account, key string) ([]byte, error) {
	if m.data[account] == nil {
		return nil, errors.New("key not found")
	}
	v, ok := m.data[account][key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return v, nil
}

func (m *mockStore) Delete(account, key string) error {
	if m.data[account] != nil {
		delete(m.data[account], key)
	}
	return nil
}

func (m *mockStore) List(account string) ([]string, error) {
	if m.data[account] == nil {
		return []string{}, nil
	}
	keys := make([]string, 0, len(m.data[account]))
	for k := range m.data[account] {
		keys = append(keys, k)
	}
	return keys, nil
}

// mockAuthFlow implements AuthFlow for testing.
type mockAuthFlow struct {
	email string
	token *oauth2.Token
	err   error
}

func (m *mockAuthFlow) Run(ctx context.Context, scopes []string) (string, *oauth2.Token, error) {
	if m.err != nil {
		return "", nil, m.err
	}
	return m.email, m.token, nil
}

// createTestConfig creates a config for testing.
func createTestConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.NewConfig()
	return cfg
}

func TestAccountService_Add(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		},
	}

	svc := NewService(cfg, store, authFlow)

	t.Run("add new account", func(t *testing.T) {
		acc, err := svc.Add(context.Background(), "work", []string{"https://www.googleapis.com/auth/gmail.readonly"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.Alias != "work" {
			t.Errorf("expected alias 'work', got '%s'", acc.Alias)
		}
		if acc.Email != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got '%s'", acc.Email)
		}
	})

	t.Run("add duplicate account", func(t *testing.T) {
		_, err := svc.Add(context.Background(), "work", []string{})
		if !errors.Is(err, account.ErrAccountExists) {
			t.Errorf("expected ErrAccountExists, got %v", err)
		}
	})
}

func TestAccountService_Remove(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add an account first
	_, err := svc.Add(context.Background(), "work", []string{})
	if err != nil {
		t.Fatalf("failed to add account: %v", err)
	}

	t.Run("remove existing account", func(t *testing.T) {
		err := svc.Remove("work")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify account is removed
		_, err = svc.cfg.GetAccount("work")
		if err == nil {
			t.Error("expected account to be removed")
		}
	})

	t.Run("remove non-existent account", func(t *testing.T) {
		err := svc.Remove("nonexistent")
		if !errors.Is(err, account.ErrAccountNotFound) {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})
}

func TestAccountService_List(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	t.Run("list empty accounts", func(t *testing.T) {
		accounts, err := svc.List()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(accounts) != 0 {
			t.Errorf("expected 0 accounts, got %d", len(accounts))
		}
	})

	// Add accounts
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	authFlow.email = "personal@example.com"
	_, _ = svc.Add(context.Background(), "personal", []string{})

	t.Run("list multiple accounts", func(t *testing.T) {
		accounts, err := svc.List()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(accounts) != 2 {
			t.Errorf("expected 2 accounts, got %d", len(accounts))
		}
	})
}

func TestAccountService_Switch(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add accounts
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	authFlow.email = "personal@example.com"
	_, _ = svc.Add(context.Background(), "personal", []string{})

	t.Run("switch to existing account", func(t *testing.T) {
		err := svc.Switch("personal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.DefaultAccount != "personal" {
			t.Errorf("expected default account 'personal', got '%s'", cfg.DefaultAccount)
		}
	})

	t.Run("switch to non-existent account", func(t *testing.T) {
		err := svc.Switch("nonexistent")
		if !errors.Is(err, account.ErrAccountNotFound) {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})
}

func TestAccountService_Show(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	t.Run("show with no accounts", func(t *testing.T) {
		_, err := svc.Show()
		if !errors.Is(err, account.ErrAccountNotFound) {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})

	// Add an account and set as default
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	cfg.DefaultAccount = "work"

	t.Run("show default account", func(t *testing.T) {
		acc, err := svc.Show()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.Alias != "work" {
			t.Errorf("expected alias 'work', got '%s'", acc.Alias)
		}
	})
}

func TestAccountService_Rename(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add an account
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})

	t.Run("rename existing account", func(t *testing.T) {
		err := svc.Rename("work", "office")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify old alias is gone
		_, err = svc.cfg.GetAccount("work")
		if err == nil {
			t.Error("expected old alias to be removed")
		}

		// Verify new alias exists
		acc, err := svc.cfg.GetAccount("office")
		if err != nil {
			t.Fatalf("failed to get renamed account: %v", err)
		}
		if acc.Email != "work@example.com" {
			t.Errorf("expected email 'work@example.com', got '%s'", acc.Email)
		}
	})

	t.Run("rename non-existent account", func(t *testing.T) {
		err := svc.Rename("nonexistent", "new")
		if !errors.Is(err, account.ErrAccountNotFound) {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})

	t.Run("rename to existing alias", func(t *testing.T) {
		// Add another account
		authFlow.email = "personal@example.com"
		_, _ = svc.Add(context.Background(), "personal", []string{})

		err := svc.Rename("personal", "office")
		if !errors.Is(err, account.ErrAccountExists) {
			t.Errorf("expected ErrAccountExists, got %v", err)
		}
	})
}

func TestAccountService_ResolveAccount(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add accounts
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	authFlow.email = "personal@example.com"
	_, _ = svc.Add(context.Background(), "personal", []string{})

	t.Run("resolve from flag value", func(t *testing.T) {
		acc, err := svc.ResolveAccount("work")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.Alias != "work" {
			t.Errorf("expected alias 'work', got '%s'", acc.Alias)
		}
	})

	t.Run("resolve from env var", func(t *testing.T) {
		os.Setenv("GOOG_ACCOUNT", "personal")
		defer os.Unsetenv("GOOG_ACCOUNT")

		acc, err := svc.ResolveAccount("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.Alias != "personal" {
			t.Errorf("expected alias 'personal', got '%s'", acc.Alias)
		}
	})

	t.Run("resolve from default account", func(t *testing.T) {
		cfg.DefaultAccount = "work"
		acc, err := svc.ResolveAccount("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.Alias != "work" {
			t.Errorf("expected alias 'work', got '%s'", acc.Alias)
		}
	})

	t.Run("resolve fallback to first account", func(t *testing.T) {
		cfg.DefaultAccount = ""
		acc, err := svc.ResolveAccount("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should return one of the accounts (order not guaranteed)
		if acc.Alias != "work" && acc.Alias != "personal" {
			t.Errorf("expected alias 'work' or 'personal', got '%s'", acc.Alias)
		}
	})

	t.Run("resolve with no accounts", func(t *testing.T) {
		emptyCfg := createTestConfig(t)
		emptySvc := NewService(emptyCfg, store, authFlow)

		_, err := emptySvc.ResolveAccount("")
		if !errors.Is(err, account.ErrAccountNotFound) {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})
}

func TestAccountService_Add_AuthFlowError(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		err: errors.New("authentication failed"),
	}

	svc := NewService(cfg, store, authFlow)

	_, err := svc.Add(context.Background(), "work", []string{"openid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "authentication failed: authentication failed" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAccountService_Add_DefaultScopes(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		},
	}

	svc := NewService(cfg, store, authFlow)

	// Add with empty scopes - should use defaults
	acc, err := svc.Add(context.Background(), "work", []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(acc.Scopes) == 0 {
		t.Error("expected default scopes to be set")
	}
}

func TestAccountService_Remove_DefaultAccountReassignment(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add two accounts
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	authFlow.email = "personal@example.com"
	_, _ = svc.Add(context.Background(), "personal", []string{})

	// Set work as default
	cfg.DefaultAccount = "work"

	// Remove work account
	err := svc.Remove("work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Default should be reassigned
	if cfg.DefaultAccount == "work" {
		t.Error("expected default account to be reassigned")
	}
}

func TestAccountService_Rename_DefaultAccountUpdate(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "work@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add account
	_, _ = svc.Add(context.Background(), "work", []string{})

	// Set as default
	cfg.DefaultAccount = "work"

	// Rename
	err := svc.Rename("work", "office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Default should be updated
	if cfg.DefaultAccount != "office" {
		t.Errorf("expected default account 'office', got '%s'", cfg.DefaultAccount)
	}
}

func TestAccountService_GetTokenManager(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{}

	svc := NewService(cfg, store, authFlow)

	tm := svc.GetTokenManager()
	if tm == nil {
		t.Error("expected non-nil token manager")
	}
}

func TestAccountService_GetConfig(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{}

	svc := NewService(cfg, store, authFlow)

	c := svc.GetConfig()
	if c != cfg {
		t.Error("expected config to match")
	}
}

func TestAccountService_Remove_DefaultAccountCleared(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add only one account
	_, err := svc.Add(context.Background(), "work", []string{})
	if err != nil {
		t.Fatalf("failed to add account: %v", err)
	}

	// Verify it's set as default (since it's the only one)
	if cfg.DefaultAccount != "work" {
		t.Errorf("expected default account 'work', got '%s'", cfg.DefaultAccount)
	}

	// Remove the account
	err = svc.Remove("work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Default should be cleared since there are no more accounts
	if cfg.DefaultAccount != "" {
		t.Errorf("expected default account to be cleared, got '%s'", cfg.DefaultAccount)
	}
}

func TestAccountService_Switch_SavesConfig(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add accounts
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	authFlow.email = "personal@example.com"
	_, _ = svc.Add(context.Background(), "personal", []string{})

	// Switch accounts
	err := svc.Switch("personal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify config was updated
	if cfg.DefaultAccount != "personal" {
		t.Errorf("expected default account 'personal', got '%s'", cfg.DefaultAccount)
	}
}

func TestAccountService_Rename_WithNoToken(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add account
	_, _ = svc.Add(context.Background(), "work", []string{})

	// Manually delete the token from store to simulate missing token
	store.Delete("work", "token")
	store.Delete("work", "scopes")

	// Rename should still work even if token is not found
	err := svc.Rename("work", "office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify rename happened
	_, err = svc.cfg.GetAccount("office")
	if err != nil {
		t.Error("expected new alias to exist")
	}
}

func TestAccountService_ResolveAccount_FlagTakesPrecedence(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add accounts
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	authFlow.email = "personal@example.com"
	_, _ = svc.Add(context.Background(), "personal", []string{})

	// Set environment variable
	os.Setenv("GOOG_ACCOUNT", "personal")
	defer os.Unsetenv("GOOG_ACCOUNT")

	// Set default
	cfg.DefaultAccount = "work"

	// Resolve with flag - should take precedence
	acc, err := svc.ResolveAccount("personal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc.Alias != "personal" {
		t.Errorf("expected alias 'personal' (from flag), got '%s'", acc.Alias)
	}
}

func TestAccountService_ResolveAccount_NonExistentFlagValue(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add an account
	_, _ = svc.Add(context.Background(), "work", []string{})

	// Resolve with non-existent flag value
	_, err := svc.ResolveAccount("nonexistent")
	if !errors.Is(err, account.ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestAccountService_Add_SavesScopes(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		},
	}

	svc := NewService(cfg, store, authFlow)

	scopes := []string{"https://www.googleapis.com/auth/gmail.readonly", "https://www.googleapis.com/auth/calendar.readonly"}
	acc, err := svc.Add(context.Background(), "work", scopes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify scopes are saved
	if len(acc.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(acc.Scopes))
	}
}

func TestAccountService_List_WithDefaults(t *testing.T) {
	store := newMockStore()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add accounts
	authFlow.email = "work@example.com"
	_, _ = svc.Add(context.Background(), "work", []string{})
	authFlow.email = "personal@example.com"
	_, _ = svc.Add(context.Background(), "personal", []string{})

	// Set default
	cfg.DefaultAccount = "work"

	accounts, err := svc.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that default is marked correctly
	for _, acc := range accounts {
		if acc.Alias == "work" && !acc.IsDefault {
			t.Error("expected work to be marked as default")
		}
		if acc.Alias == "personal" && acc.IsDefault {
			t.Error("expected personal to not be marked as default")
		}
	}
}

// mockStoreWithErrors implements Store with configurable errors.
type mockStoreWithErrors struct {
	*mockStore
	setError    error
	deleteError error
}

func newMockStoreWithErrors() *mockStoreWithErrors {
	return &mockStoreWithErrors{
		mockStore: newMockStore(),
	}
}

func (m *mockStoreWithErrors) Set(account, key string, value []byte) error {
	if m.setError != nil {
		return m.setError
	}
	return m.mockStore.Set(account, key, value)
}

func (m *mockStoreWithErrors) Delete(account, key string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	return m.mockStore.Delete(account, key)
}

func TestAccountService_Add_TokenSaveError(t *testing.T) {
	store := newMockStoreWithErrors()
	store.setError = errors.New("keyring error")
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	_, err := svc.Add(context.Background(), "work", []string{"openid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsString(err.Error(), "failed to save token") {
		t.Errorf("expected 'failed to save token' error, got %v", err)
	}
}

func TestAccountService_Remove_TokenDeleteError(t *testing.T) {
	store := newMockStoreWithErrors()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add an account first
	_, err := svc.Add(context.Background(), "work", []string{})
	if err != nil {
		t.Fatalf("failed to add account: %v", err)
	}

	// Now make delete fail
	store.deleteError = errors.New("keyring delete error")

	err = svc.Remove("work")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsString(err.Error(), "failed to delete tokens") {
		t.Errorf("expected 'failed to delete tokens' error, got %v", err)
	}
}

func TestAccountService_Rename_TokenSaveError(t *testing.T) {
	store := newMockStoreWithErrors()
	cfg := createTestConfig(t)
	authFlow := &mockAuthFlow{
		email: "test@example.com",
		token: &oauth2.Token{AccessToken: "test"},
	}

	svc := NewService(cfg, store, authFlow)

	// Add an account first
	_, err := svc.Add(context.Background(), "work", []string{})
	if err != nil {
		t.Fatalf("failed to add account: %v", err)
	}

	// Now make set fail for the rename
	store.setError = errors.New("keyring save error")

	err = svc.Rename("work", "office")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsString(err.Error(), "failed to save token") {
		t.Errorf("expected 'failed to save token' error, got %v", err)
	}
}

// helper function
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
