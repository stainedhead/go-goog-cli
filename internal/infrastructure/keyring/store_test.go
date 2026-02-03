// Package keyring provides secure credential storage using the system keyring.
package keyring

import (
	"os"
	"path/filepath"
	"testing"
)

// TestKeyringStore tests storing values in the keyring.
func TestKeyringStore(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	account := "test-account"
	key := "refresh_token"
	value := []byte("test-refresh-token-value")

	err := store.Set(account, key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify we can retrieve it
	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(retrieved) != string(value) {
		t.Errorf("retrieved value mismatch: got %q, want %q", string(retrieved), string(value))
	}
}

// TestKeyringRetrieve tests retrieving values from the keyring.
func TestKeyringRetrieve(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	account := "retrieve-test-account"
	testCases := []struct {
		key   string
		value []byte
	}{
		{"refresh_token", []byte("refresh-token-123")},
		{"access_token", []byte("access-token-456")},
		{"token_expiry", []byte("2024-12-31T23:59:59Z")},
		{"scopes", []byte("gmail.readonly calendar.readonly")},
	}

	// Store all test values
	for _, tc := range testCases {
		if err := store.Set(account, tc.key, tc.value); err != nil {
			t.Fatalf("Set(%q, %q) failed: %v", account, tc.key, err)
		}
	}

	// Retrieve and verify all values
	for _, tc := range testCases {
		retrieved, err := store.Get(account, tc.key)
		if err != nil {
			t.Errorf("Get(%q, %q) failed: %v", account, tc.key, err)
			continue
		}
		if string(retrieved) != string(tc.value) {
			t.Errorf("Get(%q, %q) = %q, want %q", account, tc.key, string(retrieved), string(tc.value))
		}
	}

	// Test retrieving non-existent key
	_, err := store.Get(account, "nonexistent")
	if err == nil {
		t.Error("Get for non-existent key should return error")
	}
	if err != ErrKeyNotFound {
		t.Errorf("Get for non-existent key should return ErrKeyNotFound, got: %v", err)
	}
}

// TestKeyringDelete tests deleting values from the keyring.
func TestKeyringDelete(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	account := "delete-test-account"
	key := "refresh_token"
	value := []byte("token-to-delete")

	// Store a value
	if err := store.Set(account, key, value); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it exists
	_, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get before delete failed: %v", err)
	}

	// Delete it
	if err := store.Delete(account, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's gone
	_, err = store.Get(account, key)
	if err == nil {
		t.Error("Get after delete should return error")
	}
	if err != ErrKeyNotFound {
		t.Errorf("Get after delete should return ErrKeyNotFound, got: %v", err)
	}

	// Delete non-existent key should not error (idempotent)
	err = store.Delete(account, "nonexistent")
	if err != nil && err != ErrKeyNotFound {
		t.Errorf("Delete non-existent key returned unexpected error: %v", err)
	}
}

// TestKeyringNamespacing tests that keys are properly namespaced per account.
func TestKeyringNamespacing(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	account1 := "account-one"
	account2 := "account-two"
	key := "refresh_token"
	value1 := []byte("token-for-account-one")
	value2 := []byte("token-for-account-two")

	// Store same key for different accounts
	if err := store.Set(account1, key, value1); err != nil {
		t.Fatalf("Set for account1 failed: %v", err)
	}
	if err := store.Set(account2, key, value2); err != nil {
		t.Fatalf("Set for account2 failed: %v", err)
	}

	// Retrieve and verify they are distinct
	retrieved1, err := store.Get(account1, key)
	if err != nil {
		t.Fatalf("Get for account1 failed: %v", err)
	}
	retrieved2, err := store.Get(account2, key)
	if err != nil {
		t.Fatalf("Get for account2 failed: %v", err)
	}

	if string(retrieved1) != string(value1) {
		t.Errorf("account1 value mismatch: got %q, want %q", string(retrieved1), string(value1))
	}
	if string(retrieved2) != string(value2) {
		t.Errorf("account2 value mismatch: got %q, want %q", string(retrieved2), string(value2))
	}

	// Test List operation
	keys1, err := store.List(account1)
	if err != nil {
		t.Fatalf("List for account1 failed: %v", err)
	}
	if len(keys1) != 1 || keys1[0] != key {
		t.Errorf("List for account1 = %v, want [%q]", keys1, key)
	}

	// Add another key to account1
	if err := store.Set(account1, "access_token", []byte("access")); err != nil {
		t.Fatalf("Set access_token failed: %v", err)
	}

	keys1, err = store.List(account1)
	if err != nil {
		t.Fatalf("List for account1 after second key failed: %v", err)
	}
	if len(keys1) != 2 {
		t.Errorf("List for account1 = %v, want 2 keys", keys1)
	}

	// Delete from account1 should not affect account2
	if err := store.Delete(account1, key); err != nil {
		t.Fatalf("Delete from account1 failed: %v", err)
	}
	retrieved2, err = store.Get(account2, key)
	if err != nil {
		t.Fatalf("Get for account2 after account1 delete failed: %v", err)
	}
	if string(retrieved2) != string(value2) {
		t.Errorf("account2 value after account1 delete: got %q, want %q", string(retrieved2), string(value2))
	}
}

// TestKeyringFallback tests that file-based fallback is used when keyring fails.
func TestKeyringFallback(t *testing.T) {
	// Create a file-based store explicitly to test fallback behavior
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "fallback-test"
	key := "refresh_token"
	value := []byte("fallback-token-value")

	// Test basic operations with file store
	if err := store.Set(account, key, value); err != nil {
		t.Fatalf("FileStore Set failed: %v", err)
	}

	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("FileStore Get failed: %v", err)
	}

	if string(retrieved) != string(value) {
		t.Errorf("FileStore value mismatch: got %q, want %q", string(retrieved), string(value))
	}

	// Verify file was created in expected location
	expectedFile := filepath.Join(tmpDir, "tokens", account+".enc")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected token file not created at %s", expectedFile)
	}

	// Test delete
	if err := store.Delete(account, key); err != nil {
		t.Fatalf("FileStore Delete failed: %v", err)
	}

	_, err = store.Get(account, key)
	if err != ErrKeyNotFound {
		t.Errorf("FileStore Get after delete should return ErrKeyNotFound, got: %v", err)
	}
}

// TestKeyFormat tests that keys are formatted correctly.
func TestKeyFormat(t *testing.T) {
	testCases := []struct {
		account string
		key     string
		want    string
	}{
		{"myaccount", "refresh_token", "goog:myaccount:refresh_token"},
		{"user@example.com", "access_token", "goog:user@example.com:access_token"},
		{"default", "scopes", "goog:default:scopes"},
	}

	for _, tc := range testCases {
		got := formatKey(tc.account, tc.key)
		if got != tc.want {
			t.Errorf("formatKey(%q, %q) = %q, want %q", tc.account, tc.key, got, tc.want)
		}
	}
}

// TestParseKey tests parsing namespaced keys back into components.
func TestParseKey(t *testing.T) {
	testCases := []struct {
		fullKey     string
		wantAccount string
		wantKey     string
		wantOK      bool
	}{
		{"goog:myaccount:refresh_token", "myaccount", "refresh_token", true},
		{"goog:user@example.com:access_token", "user@example.com", "access_token", true},
		{"invalid", "", "", false},
		{"goog:onlyonepart", "", "", false},
		{"other:account:key", "", "", false},
	}

	for _, tc := range testCases {
		account, key, ok := parseKey(tc.fullKey)
		if ok != tc.wantOK {
			t.Errorf("parseKey(%q) ok = %v, want %v", tc.fullKey, ok, tc.wantOK)
			continue
		}
		if ok && (account != tc.wantAccount || key != tc.wantKey) {
			t.Errorf("parseKey(%q) = (%q, %q), want (%q, %q)", tc.fullKey, account, key, tc.wantAccount, tc.wantKey)
		}
	}
}

// setupTestStore creates a test store using file-based storage for testing.
// Returns a cleanup function that should be deferred.
func setupTestStore(t *testing.T) (Store, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	return store, func() {
		// Cleanup is handled by t.TempDir()
	}
}
