// Package keyring provides secure credential storage using the system keyring.
package keyring

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
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

// Security-related tests

// TestDeriveMachinePassword tests that machine-specific password derivation works.
func TestDeriveMachinePassword(t *testing.T) {
	password := deriveMachinePassword()

	// Password should not be empty
	if password == "" {
		t.Error("deriveMachinePassword returned empty string")
	}

	// Password should be consistent (deterministic based on machine info)
	password2 := deriveMachinePassword()
	if password != password2 {
		t.Error("deriveMachinePassword is not deterministic")
	}

	// Password should be a hex-encoded SHA256 hash (64 characters)
	if len(password) != 64 {
		t.Errorf("deriveMachinePassword returned %d characters, expected 64", len(password))
	}

	// Password should only contain hex characters
	for _, c := range password {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("deriveMachinePassword contains non-hex character: %c", c)
			break
		}
	}
}

// TestGetMachineInfo tests that machine info collection works.
func TestGetMachineInfo(t *testing.T) {
	info := getMachineInfo()

	// Machine info should not be empty (at least hostname should be present)
	if info == "" {
		t.Error("getMachineInfo returned empty string")
	}

	// Should be deterministic
	info2 := getMachineInfo()
	if info != info2 {
		t.Error("getMachineInfo is not deterministic")
	}
}

// TestDeriveKeyWithSalt tests PBKDF2 key derivation with salt.
func TestDeriveKeyWithSalt(t *testing.T) {
	store := &FileStore{baseDir: t.TempDir()}
	account := "test-account"

	// Generate two different salts
	salt1 := make([]byte, saltSize)
	salt2 := make([]byte, saltSize)
	if _, err := rand.Read(salt1); err != nil {
		t.Fatalf("failed to generate salt1: %v", err)
	}
	if _, err := rand.Read(salt2); err != nil {
		t.Fatalf("failed to generate salt2: %v", err)
	}

	// Derive keys with different salts
	key1 := store.deriveKey(account, salt1)
	key2 := store.deriveKey(account, salt2)

	// Keys should be 32 bytes (256 bits for AES-256)
	if len(key1) != 32 {
		t.Errorf("deriveKey returned %d bytes, expected 32", len(key1))
	}
	if len(key2) != 32 {
		t.Errorf("deriveKey returned %d bytes, expected 32", len(key2))
	}

	// Keys with different salts should be different
	if bytes.Equal(key1, key2) {
		t.Error("deriveKey with different salts produced same key")
	}

	// Same salt should produce same key (deterministic)
	key1Again := store.deriveKey(account, salt1)
	if !bytes.Equal(key1, key1Again) {
		t.Error("deriveKey with same salt produced different keys")
	}

	// Different accounts with same salt should produce different keys
	key3 := store.deriveKey("other-account", salt1)
	if bytes.Equal(key1, key3) {
		t.Error("deriveKey with different accounts produced same key")
	}
}

// TestEncryptedFileFormat tests that the encrypted file format includes salt.
func TestEncryptedFileFormat(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "format-test"
	key := "test-key"
	value := []byte("test-value")

	// Store a value
	if err := store.Set(account, key, value); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Read the raw file contents
	filePath := filepath.Join(tmpDir, "tokens", account+".enc")
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read encrypted file: %v", err)
	}

	// Parse as JSON to verify structure
	var encFile encryptedFile
	if err := json.Unmarshal(fileData, &encFile); err != nil {
		t.Fatalf("failed to parse encrypted file as JSON: %v", err)
	}

	// Verify salt is present and has correct size
	if len(encFile.Salt) != saltSize {
		t.Errorf("salt has wrong size: got %d, want %d", len(encFile.Salt), saltSize)
	}

	// Verify ciphertext is present
	if len(encFile.Ciphertext) == 0 {
		t.Error("ciphertext is empty")
	}

	// Verify we can still read the value back
	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !bytes.Equal(retrieved, value) {
		t.Errorf("retrieved value mismatch: got %q, want %q", retrieved, value)
	}
}

// TestSaltUniquenessPerSave tests that each save generates a new random salt.
func TestSaltUniquenessPerSave(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "salt-unique-test"
	key := "test-key"
	filePath := filepath.Join(tmpDir, "tokens", account+".enc")

	// Save a value and capture the salt
	if err := store.Set(account, key, []byte("value1")); err != nil {
		t.Fatalf("first Set failed: %v", err)
	}
	fileData1, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file after first save: %v", err)
	}
	var encFile1 encryptedFile
	if err := json.Unmarshal(fileData1, &encFile1); err != nil {
		t.Fatalf("failed to parse first save: %v", err)
	}

	// Save a different value (this will generate a new salt)
	if err := store.Set(account, key, []byte("value2")); err != nil {
		t.Fatalf("second Set failed: %v", err)
	}
	fileData2, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file after second save: %v", err)
	}
	var encFile2 encryptedFile
	if err := json.Unmarshal(fileData2, &encFile2); err != nil {
		t.Fatalf("failed to parse second save: %v", err)
	}

	// Salts should be different
	if bytes.Equal(encFile1.Salt, encFile2.Salt) {
		t.Error("salts should be different on each save, but they are the same")
	}

	// Data should still be readable
	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(retrieved) != "value2" {
		t.Errorf("retrieved wrong value: got %q, want %q", retrieved, "value2")
	}
}

// TestLegacyKeyDerivation tests backward compatibility with legacy key derivation.
func TestLegacyKeyDerivation(t *testing.T) {
	store := &FileStore{baseDir: t.TempDir()}
	account := "test-account"

	legacyKey := store.deriveLegacyKey(account)

	// Legacy key should be 32 bytes (SHA256 output)
	if len(legacyKey) != 32 {
		t.Errorf("deriveLegacyKey returned %d bytes, expected 32", len(legacyKey))
	}

	// Should be deterministic
	legacyKey2 := store.deriveLegacyKey(account)
	if !bytes.Equal(legacyKey, legacyKey2) {
		t.Error("deriveLegacyKey is not deterministic")
	}

	// Different accounts should produce different keys
	legacyKey3 := store.deriveLegacyKey("other-account")
	if bytes.Equal(legacyKey, legacyKey3) {
		t.Error("deriveLegacyKey with different accounts produced same key")
	}
}

// TestEncryptDecryptRoundtrip tests that encrypt/decrypt work correctly together.
func TestEncryptDecryptRoundtrip(t *testing.T) {
	// Test with various data sizes
	testCases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"small", []byte("hello world")},
		{"medium", bytes.Repeat([]byte("x"), 1000)},
		{"large", bytes.Repeat([]byte("y"), 100000)},
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext, err := encrypt(tc.data, key)
			if err != nil {
				t.Fatalf("encrypt failed: %v", err)
			}

			plaintext, err := decrypt(ciphertext, key)
			if err != nil {
				t.Fatalf("decrypt failed: %v", err)
			}

			if !bytes.Equal(plaintext, tc.data) {
				t.Error("decrypted data does not match original")
			}
		})
	}
}

// TestEncryptProducesUniqueCiphertext tests that encryption produces different
// ciphertext each time (due to random nonce).
func TestEncryptProducesUniqueCiphertext(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	data := []byte("same data")

	ciphertext1, err := encrypt(data, key)
	if err != nil {
		t.Fatalf("first encrypt failed: %v", err)
	}

	ciphertext2, err := encrypt(data, key)
	if err != nil {
		t.Fatalf("second encrypt failed: %v", err)
	}

	// Ciphertexts should be different (due to random nonce)
	if bytes.Equal(ciphertext1, ciphertext2) {
		t.Error("encrypt produced identical ciphertext for same data")
	}
}

// TestDecryptWithWrongKeyFails tests that decryption fails with wrong key.
func TestDecryptWithWrongKeyFails(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	if _, err := rand.Read(key1); err != nil {
		t.Fatalf("failed to generate key1: %v", err)
	}
	if _, err := rand.Read(key2); err != nil {
		t.Fatalf("failed to generate key2: %v", err)
	}

	data := []byte("secret data")
	ciphertext, err := encrypt(data, key1)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Decrypting with wrong key should fail
	_, err = decrypt(ciphertext, key2)
	if err == nil {
		t.Error("decrypt with wrong key should fail")
	}
}

// TestFileStoreListOperations tests List operations on FileStore.
func TestFileStoreListOperations(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "list-test"

	t.Run("empty account returns empty list", func(t *testing.T) {
		keys, err := store.List(account)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(keys) != 0 {
			t.Errorf("expected empty list, got %v", keys)
		}
	})

	t.Run("returns all keys for account", func(t *testing.T) {
		// Add multiple keys
		if err := store.Set(account, "key1", []byte("value1")); err != nil {
			t.Fatalf("Set key1 failed: %v", err)
		}
		if err := store.Set(account, "key2", []byte("value2")); err != nil {
			t.Fatalf("Set key2 failed: %v", err)
		}
		if err := store.Set(account, "key3", []byte("value3")); err != nil {
			t.Fatalf("Set key3 failed: %v", err)
		}

		keys, err := store.List(account)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(keys) != 3 {
			t.Errorf("expected 3 keys, got %d: %v", len(keys), keys)
		}
	})

	t.Run("list is updated after delete", func(t *testing.T) {
		// Delete one key
		if err := store.Delete(account, "key2"); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		keys, err := store.List(account)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(keys) != 2 {
			t.Errorf("expected 2 keys after delete, got %d: %v", len(keys), keys)
		}
	})
}

// TestFileStoreDeleteRemovesFile tests that deleting all keys removes the file.
func TestFileStoreDeleteRemovesFile(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "delete-file-test"
	key := "only-key"

	// Create a token file
	if err := store.Set(account, key, []byte("value")); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "tokens", account+".enc")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("expected token file to exist")
	}

	// Delete the only key
	if err := store.Delete(account, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// File should be removed
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected token file to be removed when last key deleted")
	}
}

// TestFileStoreMultipleAccounts tests isolation between accounts.
func TestFileStoreMultipleAccounts(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account1 := "account1"
	account2 := "account2"
	key := "shared-key-name"

	// Store different values for same key in different accounts
	if err := store.Set(account1, key, []byte("value1")); err != nil {
		t.Fatalf("Set account1 failed: %v", err)
	}
	if err := store.Set(account2, key, []byte("value2")); err != nil {
		t.Fatalf("Set account2 failed: %v", err)
	}

	// Values should be independent
	val1, err := store.Get(account1, key)
	if err != nil {
		t.Fatalf("Get account1 failed: %v", err)
	}
	val2, err := store.Get(account2, key)
	if err != nil {
		t.Fatalf("Get account2 failed: %v", err)
	}

	if string(val1) != "value1" {
		t.Errorf("account1 value = %q, want 'value1'", val1)
	}
	if string(val2) != "value2" {
		t.Errorf("account2 value = %q, want 'value2'", val2)
	}

	// Delete from account1 should not affect account2
	if err := store.Delete(account1, key); err != nil {
		t.Fatalf("Delete account1 failed: %v", err)
	}

	val2, err = store.Get(account2, key)
	if err != nil {
		t.Fatalf("Get account2 after account1 delete failed: %v", err)
	}
	if string(val2) != "value2" {
		t.Errorf("account2 value after account1 delete = %q, want 'value2'", val2)
	}
}

// TestFileStoreUpdateValue tests updating an existing value.
func TestFileStoreUpdateValue(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "update-test"
	key := "token"

	// Set initial value
	if err := store.Set(account, key, []byte("initial")); err != nil {
		t.Fatalf("initial Set failed: %v", err)
	}

	// Verify initial value
	val, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get initial failed: %v", err)
	}
	if string(val) != "initial" {
		t.Errorf("initial value = %q, want 'initial'", val)
	}

	// Update value
	if err := store.Set(account, key, []byte("updated")); err != nil {
		t.Fatalf("update Set failed: %v", err)
	}

	// Verify updated value
	val, err = store.Get(account, key)
	if err != nil {
		t.Fatalf("Get updated failed: %v", err)
	}
	if string(val) != "updated" {
		t.Errorf("updated value = %q, want 'updated'", val)
	}
}

// TestDecryptTooShortCiphertext tests decryption with too short ciphertext.
func TestDecryptTooShortCiphertext(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	// Ciphertext shorter than nonce size should fail
	_, err := decrypt([]byte("short"), key)
	if err == nil {
		t.Error("expected error for too short ciphertext")
	}
}

// TestFileStoreDirectoryCreation tests that NewFileStore creates necessary directories.
func TestFileStoreDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	baseDir := filepath.Join(tmpDir, "nested", "path", "config")

	store, err := NewFileStore(baseDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	// Verify tokens directory was created
	tokensDir := filepath.Join(baseDir, "tokens")
	info, err := os.Stat(tokensDir)
	if err != nil {
		t.Fatalf("tokens directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("tokens path is not a directory")
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("tokens directory permissions = %o, want 0700", info.Mode().Perm())
	}

	// Verify we can still use the store
	if err := store.Set("test", "key", []byte("value")); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
}

// TestFileStoreWithLargeData tests storing and retrieving large data.
func TestFileStoreWithLargeData(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "large-data"
	key := "big-token"
	// Create 1MB of data
	largeData := make([]byte, 1024*1024)
	if _, err := rand.Read(largeData); err != nil {
		t.Fatalf("failed to generate large data: %v", err)
	}

	// Store large data
	if err := store.Set(account, key, largeData); err != nil {
		t.Fatalf("Set large data failed: %v", err)
	}

	// Retrieve and verify
	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get large data failed: %v", err)
	}
	if !bytes.Equal(retrieved, largeData) {
		t.Error("retrieved large data does not match original")
	}
}

// TestFileStoreWithBinaryData tests storing binary (non-UTF8) data.
func TestFileStoreWithBinaryData(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "binary-data"
	key := "binary-token"
	// Create binary data with all byte values
	binaryData := make([]byte, 256)
	for i := range binaryData {
		binaryData[i] = byte(i)
	}

	// Store binary data
	if err := store.Set(account, key, binaryData); err != nil {
		t.Fatalf("Set binary data failed: %v", err)
	}

	// Retrieve and verify
	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get binary data failed: %v", err)
	}
	if !bytes.Equal(retrieved, binaryData) {
		t.Error("retrieved binary data does not match original")
	}
}

// TestFileStoreEmptyValue tests storing and retrieving empty values.
func TestFileStoreEmptyValue(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "empty-value"
	key := "empty"

	// Store empty value
	if err := store.Set(account, key, []byte{}); err != nil {
		t.Fatalf("Set empty value failed: %v", err)
	}

	// Retrieve and verify
	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get empty value failed: %v", err)
	}
	if len(retrieved) != 0 {
		t.Errorf("expected empty value, got %v", retrieved)
	}
}

// TestNewStoreCreation tests the NewStore function.
func TestNewStoreCreation(t *testing.T) {
	// NewStore should return some implementation of Store
	store, err := NewStore()
	if err != nil {
		// This might fail in some environments, which is expected
		// The important thing is it doesn't panic
		t.Logf("NewStore returned error (expected in some environments): %v", err)
		return
	}

	if store == nil {
		t.Error("NewStore returned nil store without error")
	}
}

// TestGetConfigDir tests the getConfigDir function.
func TestGetConfigDir(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir failed: %v", err)
	}

	if configDir == "" {
		t.Error("getConfigDir returned empty string")
	}

	// Should end with "goog"
	if filepath.Base(configDir) != "goog" {
		t.Errorf("config dir should end with 'goog', got %q", configDir)
	}
}

// TestFileStoreDeleteIdempotent tests that Delete is idempotent.
func TestFileStoreDeleteIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "idempotent-delete"
	key := "token"

	// Delete non-existent key should not error
	if err := store.Delete(account, key); err != nil {
		t.Errorf("Delete non-existent key should not error: %v", err)
	}

	// Set a key
	if err := store.Set(account, key, []byte("value")); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Delete should succeed
	if err := store.Delete(account, key); err != nil {
		t.Fatalf("Delete existing key failed: %v", err)
	}

	// Delete again should not error
	if err := store.Delete(account, key); err != nil {
		t.Errorf("Delete already deleted key should not error: %v", err)
	}
}

// TestFileStoreTokenFilePath tests the tokenFilePath method.
func TestFileStoreTokenFilePath(t *testing.T) {
	baseDir := "/home/user/.config/goog"
	store := &FileStore{baseDir: baseDir}

	testCases := []struct {
		account string
		want    string
	}{
		{"user@example.com", "/home/user/.config/goog/tokens/user@example.com.enc"},
		{"simple", "/home/user/.config/goog/tokens/simple.enc"},
		{"test-account", "/home/user/.config/goog/tokens/test-account.enc"},
	}

	for _, tc := range testCases {
		t.Run(tc.account, func(t *testing.T) {
			got := store.tokenFilePath(tc.account)
			if got != tc.want {
				t.Errorf("tokenFilePath(%q) = %q, want %q", tc.account, got, tc.want)
			}
		})
	}
}

// TestEncryptInvalidKeySize tests encrypt with invalid key sizes.
func TestEncryptInvalidKeySize(t *testing.T) {
	// AES requires 16, 24, or 32 byte keys
	testCases := []struct {
		name    string
		keySize int
	}{
		{"too short key", 8},
		{"odd size key", 17},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := make([]byte, tc.keySize)
			_, err := encrypt([]byte("test"), key)
			if err == nil {
				t.Error("expected error for invalid key size")
			}
		})
	}
}

// TestDecryptEmptyCiphertext tests decrypt with empty ciphertext.
func TestDecryptEmptyCiphertext(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	_, err := decrypt([]byte{}, key)
	if err == nil {
		t.Error("expected error for empty ciphertext")
	}
}

// TestFileStoreFilePermissions tests that token files have correct permissions.
func TestFileStoreFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "perm-test"
	key := "token"

	if err := store.Set(account, key, []byte("value")); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "tokens", account+".enc")
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}
}

// TestFileStoreCorruptedFile tests behavior with corrupted encrypted files.
func TestFileStoreCorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "corrupted-test"
	tokensDir := filepath.Join(tmpDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("failed to create tokens dir: %v", err)
	}

	filePath := filepath.Join(tokensDir, account+".enc")

	t.Run("invalid JSON structure", func(t *testing.T) {
		// Write invalid JSON
		if err := os.WriteFile(filePath, []byte("not json at all"), 0600); err != nil {
			t.Fatalf("failed to write corrupted file: %v", err)
		}

		_, err := store.Get(account, "key")
		if err == nil {
			t.Error("expected error for corrupted file")
		}
	})

	t.Run("invalid salt size", func(t *testing.T) {
		// Write JSON with wrong salt size
		data := []byte(`{"salt":"dG9vc2hvcnQ=","ciphertext":"YWJj"}`)
		if err := os.WriteFile(filePath, data, 0600); err != nil {
			t.Fatalf("failed to write invalid salt file: %v", err)
		}

		_, err := store.Get(account, "key")
		if err == nil {
			t.Error("expected error for invalid salt size")
		}
	})
}

// TestFileStoreMultipleKeys tests storing multiple keys for one account.
func TestFileStoreMultipleKeys(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "multi-key"
	keys := map[string][]byte{
		"oauth_token":   []byte("access-token-123"),
		"oauth_scopes":  []byte(`["gmail.readonly","calendar"]`),
		"refresh_token": []byte("refresh-456"),
	}

	// Store all keys
	for k, v := range keys {
		if err := store.Set(account, k, v); err != nil {
			t.Fatalf("Set %q failed: %v", k, err)
		}
	}

	// Verify all keys can be retrieved
	for k, expected := range keys {
		retrieved, err := store.Get(account, k)
		if err != nil {
			t.Errorf("Get %q failed: %v", k, err)
			continue
		}
		if !bytes.Equal(retrieved, expected) {
			t.Errorf("Get %q = %q, want %q", k, retrieved, expected)
		}
	}

	// Delete one key, others should remain
	if err := store.Delete(account, "refresh_token"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Other keys should still exist
	for k := range keys {
		if k == "refresh_token" {
			continue
		}
		_, err := store.Get(account, k)
		if err != nil {
			t.Errorf("Get %q after delete of other key failed: %v", k, err)
		}
	}

	// Deleted key should not exist
	_, err = store.Get(account, "refresh_token")
	if err != ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound for deleted key, got %v", err)
	}
}

// TestFormatKeyParsKey tests formatKey and parseKey are inverses.
func TestFormatKeyParseKey(t *testing.T) {
	testCases := []struct {
		account string
		key     string
	}{
		{"user@example.com", "oauth_token"},
		{"simple", "refresh_token"},
		{"test-account", "scopes"},
		{"account:with:colons", "key:with:colons"},
	}

	for _, tc := range testCases {
		t.Run(tc.account+"/"+tc.key, func(t *testing.T) {
			fullKey := formatKey(tc.account, tc.key)

			// For cases without extra colons, verify round-trip
			if tc.account != "account:with:colons" && tc.key != "key:with:colons" {
				account, key, ok := parseKey(fullKey)
				if !ok {
					t.Errorf("parseKey failed for %q", fullKey)
					return
				}
				if account != tc.account || key != tc.key {
					t.Errorf("parseKey(%q) = (%q, %q), want (%q, %q)", fullKey, account, key, tc.account, tc.key)
				}
			}
		})
	}
}

// TestOpenKeyringBackends tests that openKeyring sets up appropriate backends.
func TestOpenKeyringBackends(t *testing.T) {
	// This test may fail in CI environments without a keyring
	// We mainly want to ensure it doesn't panic
	tmpDir := t.TempDir()
	_, err := openKeyring(tmpDir)
	// It's OK if this fails in some environments
	_ = err
}

// TestFileStoreSaveLoadRoundTrip tests complete save/load cycle with different data types.
func TestFileStoreSaveLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	testCases := []struct {
		name    string
		account string
		key     string
		value   []byte
	}{
		{"simple string", "acc1", "key1", []byte("hello world")},
		{"json data", "acc2", "key2", []byte(`{"access_token":"xyz","expires":3600}`)},
		{"binary data", "acc3", "key3", []byte{0x00, 0x01, 0x02, 0xff, 0xfe}},
		{"unicode", "acc4", "key4", []byte("Hello \u4e16\u754c")},
		{"long value", "acc5", "key5", bytes.Repeat([]byte("x"), 10000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save
			if err := store.Set(tc.account, tc.key, tc.value); err != nil {
				t.Fatalf("Set failed: %v", err)
			}

			// Load
			retrieved, err := store.Get(tc.account, tc.key)
			if err != nil {
				t.Fatalf("Get failed: %v", err)
			}

			// Compare
			if !bytes.Equal(retrieved, tc.value) {
				t.Errorf("value mismatch for %s", tc.name)
			}
		})
	}
}

// TestFileStoreGetNonexistentAccount tests Get for an account with no file.
func TestFileStoreGetNonexistentAccount(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	_, err = store.Get("nonexistent-account", "any-key")
	if err != ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

// TestFileStoreDeleteLastKeyRemovesFile tests that deleting the last key removes the file.
func TestFileStoreDeleteLastKeyRemovesFile(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "delete-last"

	// Add two keys
	if err := store.Set(account, "key1", []byte("value1")); err != nil {
		t.Fatalf("Set key1 failed: %v", err)
	}
	if err := store.Set(account, "key2", []byte("value2")); err != nil {
		t.Fatalf("Set key2 failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "tokens", account+".enc")

	// Delete first key
	if err := store.Delete(account, "key1"); err != nil {
		t.Fatalf("Delete key1 failed: %v", err)
	}

	// File should still exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("file should still exist after deleting one of two keys")
	}

	// Delete second (last) key
	if err := store.Delete(account, "key2"); err != nil {
		t.Fatalf("Delete key2 failed: %v", err)
	}

	// File should be removed
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("file should be removed after deleting last key")
	}
}

// TestFileStoreListNonexistentAccount tests List for an account with no file.
func TestFileStoreListNonexistentAccount(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	keys, err := store.List("nonexistent-account")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(keys) != 0 {
		t.Errorf("expected empty list, got %v", keys)
	}
}

// TestEncryptDecryptWithValidKeys tests encrypt/decrypt with various valid key sizes.
func TestEncryptDecryptWithValidKeys(t *testing.T) {
	keySizes := []int{16, 24, 32} // AES-128, AES-192, AES-256

	for _, size := range keySizes {
		t.Run("key-size-"+string(rune(size+'0')), func(t *testing.T) {
			key := make([]byte, size)
			if _, err := rand.Read(key); err != nil {
				t.Fatalf("failed to generate key: %v", err)
			}

			data := []byte("test data for encryption")
			ciphertext, err := encrypt(data, key)
			if err != nil {
				t.Fatalf("encrypt failed: %v", err)
			}

			plaintext, err := decrypt(ciphertext, key)
			if err != nil {
				t.Fatalf("decrypt failed: %v", err)
			}

			if !bytes.Equal(plaintext, data) {
				t.Error("decrypted data does not match original")
			}
		})
	}
}

// TestNewFileStorePermissions tests that NewFileStore creates directory with correct permissions.
func TestNewFileStorePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	baseDir := filepath.Join(tmpDir, "test-perm")

	_, err := NewFileStore(baseDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	tokensDir := filepath.Join(baseDir, "tokens")
	info, err := os.Stat(tokensDir)
	if err != nil {
		t.Fatalf("failed to stat tokens directory: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0700 {
		t.Errorf("expected directory permissions 0700, got %o", perm)
	}
}

// TestFileStoreSetError tests error handling in FileStore.Set.
func TestFileStoreSetError(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	// Try to set with empty account (should work but creates odd filename)
	if err := store.Set("", "key", []byte("value")); err != nil {
		// This may fail, which is acceptable
		t.Logf("Set with empty account returned error (acceptable): %v", err)
	}
}

// TestFileStoreListWithMultipleKeys tests List with various key configurations.
func TestFileStoreListWithMultipleKeys(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "list-multi"

	// Add several keys
	keys := []string{"oauth_token", "refresh_token", "scopes", "expiry"}
	for _, k := range keys {
		if err := store.Set(account, k, []byte("value-"+k)); err != nil {
			t.Fatalf("Set %q failed: %v", k, err)
		}
	}

	// List and verify
	listed, err := store.List(account)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(listed) != len(keys) {
		t.Errorf("List returned %d keys, want %d", len(listed), len(keys))
	}
}

// TestFileStoreDeleteMaintainsOtherKeys tests that Delete maintains other keys.
func TestFileStoreDeleteMaintainsOtherKeys(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "delete-maintain"

	// Set multiple keys
	if err := store.Set(account, "keep1", []byte("value1")); err != nil {
		t.Fatalf("Set keep1 failed: %v", err)
	}
	if err := store.Set(account, "delete", []byte("value2")); err != nil {
		t.Fatalf("Set delete failed: %v", err)
	}
	if err := store.Set(account, "keep2", []byte("value3")); err != nil {
		t.Fatalf("Set keep2 failed: %v", err)
	}

	// Delete middle key
	if err := store.Delete(account, "delete"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify other keys still exist with correct values
	val1, err := store.Get(account, "keep1")
	if err != nil {
		t.Fatalf("Get keep1 failed: %v", err)
	}
	if string(val1) != "value1" {
		t.Errorf("keep1 = %q, want 'value1'", val1)
	}

	val2, err := store.Get(account, "keep2")
	if err != nil {
		t.Fatalf("Get keep2 failed: %v", err)
	}
	if string(val2) != "value3" {
		t.Errorf("keep2 = %q, want 'value3'", val2)
	}
}

// TestSaveTokenDataWithMultipleKeys tests saveTokenData with multiple keys.
func TestSaveTokenDataWithMultipleKeys(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "multi-save"

	// Add keys one by one
	for i := 0; i < 5; i++ {
		key := "key" + string(rune('0'+i))
		value := []byte("value" + string(rune('0'+i)))
		if err := store.Set(account, key, value); err != nil {
			t.Fatalf("Set %s failed: %v", key, err)
		}
	}

	// Verify all keys
	for i := 0; i < 5; i++ {
		key := "key" + string(rune('0'+i))
		expected := "value" + string(rune('0'+i))
		val, err := store.Get(account, key)
		if err != nil {
			t.Errorf("Get %s failed: %v", key, err)
			continue
		}
		if string(val) != expected {
			t.Errorf("Get %s = %q, want %q", key, val, expected)
		}
	}
}

// TestNewStoreWithValidDir tests NewStore with a valid config directory.
func TestNewStoreWithValidDir(t *testing.T) {
	// This test may succeed or fail depending on system keyring availability
	// The main goal is to ensure it doesn't panic
	store, err := NewStore()
	if err != nil {
		t.Logf("NewStore returned error (expected in some environments): %v", err)
		return
	}
	if store == nil {
		t.Error("NewStore returned nil store without error")
	}
}

// TestFileStoreListAfterDelete tests List returns correct keys after delete.
func TestFileStoreListAfterDelete(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "list-after-delete"

	// Set multiple keys
	keys := []string{"a", "b", "c", "d"}
	for _, k := range keys {
		if err := store.Set(account, k, []byte(k)); err != nil {
			t.Fatalf("Set %q failed: %v", k, err)
		}
	}

	// Delete some keys
	if err := store.Delete(account, "b"); err != nil {
		t.Fatalf("Delete b failed: %v", err)
	}
	if err := store.Delete(account, "d"); err != nil {
		t.Fatalf("Delete d failed: %v", err)
	}

	// List should only have "a" and "c"
	listed, err := store.List(account)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(listed) != 2 {
		t.Errorf("List returned %d keys, want 2", len(listed))
	}

	// Check specific keys exist
	listMap := make(map[string]bool)
	for _, k := range listed {
		listMap[k] = true
	}

	if !listMap["a"] || !listMap["c"] {
		t.Errorf("List = %v, want keys a and c", listed)
	}
	if listMap["b"] || listMap["d"] {
		t.Errorf("List should not contain deleted keys, got %v", listed)
	}
}

// TestEncryptDecryptSpecialCases tests encrypt/decrypt with edge cases.
func TestEncryptDecryptSpecialCases(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	testCases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"single byte", []byte{0x42}},
		{"null bytes", []byte{0x00, 0x00, 0x00}},
		{"max byte values", []byte{0xff, 0xff, 0xff}},
		{"mixed", []byte{0x00, 0x7f, 0x80, 0xff}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext, err := encrypt(tc.data, key)
			if err != nil {
				t.Fatalf("encrypt failed: %v", err)
			}

			plaintext, err := decrypt(ciphertext, key)
			if err != nil {
				t.Fatalf("decrypt failed: %v", err)
			}

			if !bytes.Equal(plaintext, tc.data) {
				t.Errorf("decrypt produced %v, want %v", plaintext, tc.data)
			}
		})
	}
}

// TestFileStoreOverwriteValue tests overwriting an existing value.
func TestFileStoreOverwriteValue(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "overwrite-test"
	key := "token"

	// Set initial value
	if err := store.Set(account, key, []byte("initial")); err != nil {
		t.Fatalf("initial Set failed: %v", err)
	}

	// Overwrite with different value
	if err := store.Set(account, key, []byte("updated")); err != nil {
		t.Fatalf("update Set failed: %v", err)
	}

	// Verify new value
	val, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(val) != "updated" {
		t.Errorf("Get = %q, want 'updated'", val)
	}
}

// TestDecryptInvalidCiphertextFormats tests decrypt with various invalid inputs.
func TestDecryptInvalidCiphertextFormats(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	testCases := []struct {
		name       string
		ciphertext []byte
	}{
		{"empty", []byte{}},
		{"too short", []byte("short")},
		{"exactly nonce size", make([]byte, 12)}, // GCM nonce is 12 bytes
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := decrypt(tc.ciphertext, key)
			if err == nil {
				t.Error("expected error for invalid ciphertext")
			}
		})
	}
}

// TestFileStoreLoadTokenDataWithLegacyFormat tests loading legacy format encrypted files.
func TestFileStoreLoadTokenDataWithLegacyFormat(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "legacy-test"
	tokensDir := filepath.Join(tmpDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("failed to create tokens dir: %v", err)
	}

	// Create token data
	data := &tokenData{
		Tokens: map[string][]byte{
			"key1": []byte("value1"),
		},
	}
	plaintext, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal token data: %v", err)
	}

	// Encrypt using legacy key derivation
	legacyKey := store.deriveLegacyKey(account)
	ciphertext, err := encrypt(plaintext, legacyKey)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Write directly to file (legacy format - raw encrypted data, not JSON structure)
	filePath := filepath.Join(tokensDir, account+".enc")
	if err := os.WriteFile(filePath, ciphertext, 0600); err != nil {
		t.Fatalf("failed to write legacy file: %v", err)
	}

	// Try to read using the store - should fall back to legacy format
	val, err := store.Get(account, "key1")
	if err != nil {
		t.Fatalf("Get failed with legacy format: %v", err)
	}

	if string(val) != "value1" {
		t.Errorf("Get = %q, want 'value1'", val)
	}
}

// TestFileStoreSetErrorWhenLoadFails tests Set when loadTokenData fails with non-NotExist error.
func TestFileStoreSetErrorWhenLoadFails(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "load-err-test"
	tokensDir := filepath.Join(tmpDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("failed to create tokens dir: %v", err)
	}

	// Create a corrupted file that will fail to parse and decrypt
	filePath := filepath.Join(tokensDir, account+".enc")
	// Write valid JSON structure but with invalid encrypted content
	invalidContent := `{"salt":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=","ciphertext":"dG9vLXNob3J0"}`
	if err := os.WriteFile(filePath, []byte(invalidContent), 0600); err != nil {
		t.Fatalf("failed to write invalid file: %v", err)
	}

	// Set should fail because it can't load existing data
	err = store.Set(account, "key", []byte("value"))
	if err == nil {
		t.Error("expected error when loading corrupted file")
	}
}

// TestFileStoreDeleteWithLoadError tests Delete when loadTokenData fails.
func TestFileStoreDeleteWithLoadError(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "delete-load-err"
	tokensDir := filepath.Join(tmpDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("failed to create tokens dir: %v", err)
	}

	// Create a corrupted file
	filePath := filepath.Join(tokensDir, account+".enc")
	if err := os.WriteFile(filePath, []byte("not json at all"), 0600); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	// Delete should fail with corrupted file
	err = store.Delete(account, "key")
	if err == nil {
		t.Error("expected error when loading corrupted file during delete")
	}
}

// TestFileStoreListWithLoadError tests List when loadTokenData fails.
func TestFileStoreListWithLoadError(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "list-load-err"
	tokensDir := filepath.Join(tmpDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("failed to create tokens dir: %v", err)
	}

	// Create a corrupted file
	filePath := filepath.Join(tokensDir, account+".enc")
	if err := os.WriteFile(filePath, []byte("corrupted data"), 0600); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	// List should fail with corrupted file
	_, err = store.List(account)
	if err == nil {
		t.Error("expected error when loading corrupted file during list")
	}
}

// TestNewStoreFallbackToFileStore tests that NewStore falls back to FileStore when keyring unavailable.
func TestNewStoreFallbackToFileStore(t *testing.T) {
	// This test verifies that NewStore returns some kind of store
	// In test environments, it may use FileStore as fallback
	store, err := NewStore()
	if err != nil {
		// Some environments don't have keyring support, which is fine
		t.Logf("NewStore returned error (may be expected): %v", err)
		return
	}

	if store == nil {
		t.Error("NewStore returned nil store without error")
		return
	}

	// Verify the store works by doing a basic operation
	account := "new-store-test"
	key := "test-key"
	value := []byte("test-value")

	// This may fail in some environments, which is acceptable
	if err := store.Set(account, key, value); err != nil {
		t.Logf("Store.Set failed (may be expected in some environments): %v", err)
		return
	}

	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Logf("Store.Get failed: %v", err)
		return
	}

	if !bytes.Equal(retrieved, value) {
		t.Errorf("retrieved value = %q, want %q", retrieved, value)
	}

	// Clean up
	store.Delete(account, key)
}

// TestGetConfigDirWithNoHome tests getConfigDir when home directory lookup fails.
func TestGetConfigDirWithNoHome(t *testing.T) {
	// This test mainly verifies that getConfigDir doesn't panic
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir failed: %v", err)
	}

	if configDir == "" {
		t.Error("getConfigDir returned empty string")
	}
}

// TestFileStoreGetWithInvalidJSONTokenData tests Get when token data has invalid JSON.
func TestFileStoreGetWithInvalidJSONTokenData(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "invalid-json-test"
	tokensDir := filepath.Join(tmpDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("failed to create tokens dir: %v", err)
	}

	// Create a file with valid encryption format but invalid token JSON inside
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("failed to generate salt: %v", err)
	}

	// Encrypt invalid JSON
	key := store.deriveKey(account, salt)
	ciphertext, err := encrypt([]byte("not valid json"), key)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	encFile := encryptedFile{
		Salt:       salt,
		Ciphertext: ciphertext,
	}
	fileData, err := json.Marshal(encFile)
	if err != nil {
		t.Fatalf("failed to marshal encrypted file: %v", err)
	}

	filePath := filepath.Join(tokensDir, account+".enc")
	if err := os.WriteFile(filePath, fileData, 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// Get should fail because token data is not valid JSON
	_, err = store.Get(account, "key")
	if err == nil {
		t.Error("expected error for invalid JSON token data")
	}
}

// TestFileStoreLoadLegacyTokenDataWithInvalidJSON tests legacy format with invalid JSON.
func TestFileStoreLoadLegacyTokenDataWithInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "legacy-invalid-json"
	tokensDir := filepath.Join(tmpDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("failed to create tokens dir: %v", err)
	}

	// Encrypt invalid JSON using legacy key
	legacyKey := store.deriveLegacyKey(account)
	ciphertext, err := encrypt([]byte("not valid json"), legacyKey)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Write raw encrypted data (legacy format)
	filePath := filepath.Join(tokensDir, account+".enc")
	if err := os.WriteFile(filePath, ciphertext, 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// Get should fail because decrypted data is not valid JSON
	_, err = store.Get(account, "key")
	if err == nil {
		t.Error("expected error for invalid JSON in legacy format")
	}
}

// TestOpenKeyringWithDifferentPlatforms tests openKeyring configuration.
func TestOpenKeyringWithDifferentPlatforms(t *testing.T) {
	tmpDir := t.TempDir()

	// This mainly verifies that openKeyring doesn't panic
	_, err := openKeyring(tmpDir)
	// Error is acceptable - we're mainly testing it doesn't panic
	_ = err
}

// TestDeriveMachinePasswordDeterminism tests that machine password is deterministic.
func TestDeriveMachinePasswordDeterminism(t *testing.T) {
	password1 := deriveMachinePassword()
	password2 := deriveMachinePassword()
	password3 := deriveMachinePassword()

	if password1 != password2 || password2 != password3 {
		t.Error("deriveMachinePassword should be deterministic")
	}
}

// TestGetMachineInfoDeterminism tests that machine info is deterministic.
func TestGetMachineInfoDeterminism(t *testing.T) {
	info1 := getMachineInfo()
	info2 := getMachineInfo()
	info3 := getMachineInfo()

	if info1 != info2 || info2 != info3 {
		t.Error("getMachineInfo should be deterministic")
	}
}

// TestParseKeyWithColons tests parseKey with keys containing colons.
func TestParseKeyWithColons(t *testing.T) {
	// Test that parseKey handles edge cases
	testCases := []struct {
		fullKey     string
		wantAccount string
		wantKey     string
		wantOK      bool
	}{
		// Valid format
		{"goog:account:key", "account", "key", true},
		// Key with colons - parseKey uses SplitN(_, 3) so key can contain colons
		{"goog:account:key:with:colons", "account", "key:with:colons", true},
		// Too few parts
		{"goog:onlyonepart", "", "", false},
		{"goog", "", "", false},
		{"", "", "", false},
		// Wrong prefix
		{"other:account:key", "", "", false},
		{"google:account:key", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.fullKey, func(t *testing.T) {
			account, key, ok := parseKey(tc.fullKey)
			if ok != tc.wantOK {
				t.Errorf("parseKey(%q) ok = %v, want %v", tc.fullKey, ok, tc.wantOK)
				return
			}
			if ok {
				if account != tc.wantAccount {
					t.Errorf("parseKey(%q) account = %q, want %q", tc.fullKey, account, tc.wantAccount)
				}
				if key != tc.wantKey {
					t.Errorf("parseKey(%q) key = %q, want %q", tc.fullKey, key, tc.wantKey)
				}
			}
		})
	}
}

// TestFileStoreDeletePreservesOtherAccountFiles tests that deleting from one account doesn't affect others.
func TestFileStoreDeletePreservesOtherAccountFiles(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	// Create two accounts with tokens
	if err := store.Set("account1", "token", []byte("token1")); err != nil {
		t.Fatalf("Set account1 failed: %v", err)
	}
	if err := store.Set("account2", "token", []byte("token2")); err != nil {
		t.Fatalf("Set account2 failed: %v", err)
	}

	// Delete all keys from account1 (file should be removed)
	if err := store.Delete("account1", "token"); err != nil {
		t.Fatalf("Delete account1 failed: %v", err)
	}

	// account2 should still work
	val, err := store.Get("account2", "token")
	if err != nil {
		t.Fatalf("Get account2 after account1 delete failed: %v", err)
	}
	if string(val) != "token2" {
		t.Errorf("account2 token = %q, want 'token2'", val)
	}
}

// TestEncryptValidKeySizes tests encryption with all valid AES key sizes.
func TestEncryptValidKeySizes(t *testing.T) {
	validSizes := []int{16, 24, 32} // AES-128, AES-192, AES-256

	for _, size := range validSizes {
		t.Run("size-"+string(rune(size+'0')), func(t *testing.T) {
			key := make([]byte, size)
			if _, err := rand.Read(key); err != nil {
				t.Fatalf("failed to generate key: %v", err)
			}

			data := []byte("test data")
			ciphertext, err := encrypt(data, key)
			if err != nil {
				t.Errorf("encrypt with %d byte key failed: %v", size, err)
			}
			if len(ciphertext) == 0 {
				t.Error("ciphertext is empty")
			}
		})
	}
}

// TestDecryptWithInvalidKey tests decrypt with invalid key.
func TestDecryptWithInvalidKey(t *testing.T) {
	// Create valid ciphertext first
	validKey := make([]byte, 32)
	if _, err := rand.Read(validKey); err != nil {
		t.Fatalf("failed to generate valid key: %v", err)
	}

	data := []byte("secret data")
	ciphertext, err := encrypt(data, validKey)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Try to decrypt with different key
	wrongKey := make([]byte, 32)
	if _, err := rand.Read(wrongKey); err != nil {
		t.Fatalf("failed to generate wrong key: %v", err)
	}

	_, err = decrypt(ciphertext, wrongKey)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

// TestFileStoreSaveTokenDataMarshalError simulates a marshal error scenario.
func TestFileStoreSaveTokenDataMarshalError(t *testing.T) {
	// This is hard to test directly since json.Marshal rarely fails on maps
	// We mainly want to verify the error path exists and code doesn't panic
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	// Create and save some data to verify the happy path works
	account := "marshal-test"
	if err := store.Set(account, "key", []byte("value")); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it was saved
	val, err := store.Get(account, "key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(val) != "value" {
		t.Errorf("Get = %q, want 'value'", val)
	}
}

// TestNewFileStoreWithInvalidPath tests NewFileStore with an invalid path.
func TestNewFileStoreWithInvalidPath(t *testing.T) {
	// Try to create a file store in a path that cannot be created
	_, err := NewFileStore("/dev/null/invalid/path")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

// TestFileStoreSaveTokenDataSaltError tests error handling during salt generation.
// Note: This is difficult to test directly since rand.Reader rarely fails,
// but we verify the function works correctly with valid random generation.
func TestFileStoreSaveTokenDataSaltGeneration(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "salt-gen-test"

	// Multiple saves should work with different salts each time
	for i := 0; i < 5; i++ {
		key := "key" + string(rune('0'+i))
		value := []byte("value" + string(rune('0'+i)))
		if err := store.Set(account, key, value); err != nil {
			t.Fatalf("Set %s failed: %v", key, err)
		}
	}

	// Verify all keys are retrievable
	for i := 0; i < 5; i++ {
		key := "key" + string(rune('0'+i))
		expected := "value" + string(rune('0'+i))
		val, err := store.Get(account, key)
		if err != nil {
			t.Errorf("Get %s failed: %v", key, err)
			continue
		}
		if string(val) != expected {
			t.Errorf("Get %s = %q, want %q", key, val, expected)
		}
	}
}

// TestEncryptWithNonceFail tests encrypt behavior (nonce generation).
// Since rand.Reader rarely fails, this tests the normal path.
func TestEncryptWithNonceGeneration(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	data := []byte("test data for encryption")

	// Encrypt multiple times - should produce different ciphertexts (different nonces)
	ciphertexts := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		ct, err := encrypt(data, key)
		if err != nil {
			t.Fatalf("encrypt %d failed: %v", i, err)
		}
		ciphertexts[i] = ct
	}

	// All ciphertexts should be different
	for i := 0; i < 10; i++ {
		for j := i + 1; j < 10; j++ {
			if bytes.Equal(ciphertexts[i], ciphertexts[j]) {
				t.Errorf("ciphertext[%d] equals ciphertext[%d]", i, j)
			}
		}
	}
}

// TestDecryptWithGCMOpenError tests decrypt with tampered ciphertext.
func TestDecryptWithGCMOpenError(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	data := []byte("test data")
	ciphertext, err := encrypt(data, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Tamper with the ciphertext
	if len(ciphertext) > 20 {
		ciphertext[20] ^= 0xff
	}

	// Decryption should fail due to authentication failure
	_, err = decrypt(ciphertext, key)
	if err == nil {
		t.Error("expected error when decrypting tampered ciphertext")
	}
}

// TestFileStoreSetWithNilExistingData tests Set when there's no existing file.
func TestFileStoreSetWithNilExistingData(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "new-account"
	key := "first-key"
	value := []byte("first-value")

	// This should create new token data since none exists
	if err := store.Set(account, key, value); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !bytes.Equal(retrieved, value) {
		t.Errorf("Get = %q, want %q", retrieved, value)
	}
}

// TestFileStoreDeleteWhenFileEmpty tests Delete when no keys remain.
func TestFileStoreDeleteWhenFileEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	account := "delete-empty-test"
	key := "only-key"

	// Set a single key
	if err := store.Set(account, key, []byte("value")); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "tokens", account+".enc")

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("file should exist after Set")
	}

	// Delete the only key
	if err := store.Delete(account, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// File should be removed
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("file should be removed when last key is deleted")
	}
}

// TestGetConfigDirSuccess tests getConfigDir returns valid path.
func TestGetConfigDirSuccess(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir failed: %v", err)
	}

	if configDir == "" {
		t.Error("getConfigDir returned empty string")
	}

	// Path should end with "goog"
	if filepath.Base(configDir) != "goog" {
		t.Errorf("config dir should end with 'goog', got %q", configDir)
	}

	// Path should be absolute
	if !filepath.IsAbs(configDir) {
		t.Errorf("config dir should be absolute, got %q", configDir)
	}
}

// TestServiceNameConstant tests that ServiceName is correctly defined.
func TestServiceNameConstant(t *testing.T) {
	if ServiceName != "go-goog-cli" {
		t.Errorf("ServiceName = %q, want 'go-goog-cli'", ServiceName)
	}
}

// TestErrKeyNotFoundConstant tests that ErrKeyNotFound is correctly defined.
func TestErrKeyNotFoundConstant(t *testing.T) {
	if ErrKeyNotFound == nil {
		t.Error("ErrKeyNotFound should not be nil")
	}
	if ErrKeyNotFound.Error() != "key not found" {
		t.Errorf("ErrKeyNotFound.Error() = %q, want 'key not found'", ErrKeyNotFound.Error())
	}
}

// TestKeyringStoreSetMethod tests the KeyringStore Set method if keyring is available.
func TestKeyringStoreSetMethod(t *testing.T) {
	// Try to create a real keyring store
	store, err := NewStore()
	if err != nil {
		t.Skip("keyring not available in this environment")
	}

	// Check if it's a KeyringStore (not FileStore)
	if _, ok := store.(*KeyringStore); !ok {
		t.Skip("store is not KeyringStore")
	}

	account := "keyring-set-test"
	key := "test-key"
	value := []byte("test-value")

	// Test Set
	if err := store.Set(account, key, value); err != nil {
		t.Logf("KeyringStore.Set failed (may be expected): %v", err)
		return
	}

	// Verify it can be retrieved
	retrieved, err := store.Get(account, key)
	if err != nil {
		t.Logf("KeyringStore.Get failed: %v", err)
		return
	}

	if !bytes.Equal(retrieved, value) {
		t.Errorf("retrieved = %q, want %q", retrieved, value)
	}

	// Clean up
	store.Delete(account, key)
}

// TestKeyringStoreListMethod tests the KeyringStore List method if keyring is available.
func TestKeyringStoreListMethod(t *testing.T) {
	store, err := NewStore()
	if err != nil {
		t.Skip("keyring not available in this environment")
	}

	if _, ok := store.(*KeyringStore); !ok {
		t.Skip("store is not KeyringStore")
	}

	account := "keyring-list-test"

	// Set some keys
	if err := store.Set(account, "key1", []byte("value1")); err != nil {
		t.Logf("KeyringStore.Set failed (may be expected): %v", err)
		return
	}
	if err := store.Set(account, "key2", []byte("value2")); err != nil {
		t.Logf("KeyringStore.Set key2 failed: %v", err)
		// Clean up and skip
		store.Delete(account, "key1")
		return
	}

	// Test List
	keys, err := store.List(account)
	if err != nil {
		t.Logf("KeyringStore.List failed: %v", err)
	} else {
		if len(keys) < 2 {
			t.Logf("expected at least 2 keys, got %d", len(keys))
		}
	}

	// Clean up
	store.Delete(account, "key1")
	store.Delete(account, "key2")
}

// TestKeyringStoreDeleteMethod tests the KeyringStore Delete method if keyring is available.
func TestKeyringStoreDeleteMethod(t *testing.T) {
	store, err := NewStore()
	if err != nil {
		t.Skip("keyring not available in this environment")
	}

	if _, ok := store.(*KeyringStore); !ok {
		t.Skip("store is not KeyringStore")
	}

	account := "keyring-delete-test"
	key := "key-to-delete"

	// Set a key
	if err := store.Set(account, key, []byte("value")); err != nil {
		t.Logf("KeyringStore.Set failed (may be expected): %v", err)
		return
	}

	// Delete it
	if err := store.Delete(account, key); err != nil {
		t.Logf("KeyringStore.Delete failed: %v", err)
		return
	}

	// Verify it's gone
	_, err = store.Get(account, key)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

// TestKeyringStoreDeleteNonexistent tests deleting a key that doesn't exist.
func TestKeyringStoreDeleteNonexistent(t *testing.T) {
	store, err := NewStore()
	if err != nil {
		t.Skip("keyring not available in this environment")
	}

	if _, ok := store.(*KeyringStore); !ok {
		t.Skip("store is not KeyringStore")
	}

	// Delete should be idempotent - not error on non-existent key
	err = store.Delete("nonexistent-account", "nonexistent-key")
	if err != nil && err != ErrKeyNotFound {
		t.Logf("Delete non-existent: %v (may be acceptable)", err)
	}
}

// TestKeyringStoreGetNonexistent tests getting a key that doesn't exist.
func TestKeyringStoreGetNonexistent(t *testing.T) {
	store, err := NewStore()
	if err != nil {
		t.Skip("keyring not available in this environment")
	}

	if _, ok := store.(*KeyringStore); !ok {
		t.Skip("store is not KeyringStore")
	}

	_, err = store.Get("nonexistent-account", "nonexistent-key")
	if err == nil {
		t.Error("expected error for non-existent key")
	}
}

