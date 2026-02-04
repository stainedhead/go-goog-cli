// Package keyring provides secure credential storage using the system keyring.
// It supports macOS Keychain as the primary backend with an encrypted file
// fallback for environments where the system keyring is unavailable.
package keyring

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/99designs/keyring"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// ServiceName is the service identifier used in the system keyring.
	ServiceName = "go-goog-cli"

	// keyPrefix is the prefix for all keys stored by this application.
	keyPrefix = "goog"

	// pbkdf2Iterations is the number of iterations for PBKDF2 key derivation.
	// This provides protection against brute-force attacks.
	pbkdf2Iterations = 100000

	// saltSize is the size of the random salt in bytes.
	saltSize = 32
)

// ErrKeyNotFound is returned when a requested key does not exist in the store.
var ErrKeyNotFound = errors.New("key not found")

// Store defines the interface for secure credential storage.
type Store interface {
	// Set stores a value for the given account and key.
	Set(account, key string, value []byte) error

	// Get retrieves a value for the given account and key.
	// Returns ErrKeyNotFound if the key does not exist.
	Get(account, key string) ([]byte, error)

	// Delete removes a value for the given account and key.
	// Returns nil if the key does not exist (idempotent).
	Delete(account, key string) error

	// List returns all keys stored for the given account.
	List(account string) ([]string, error)
}

// KeyringStore implements Store using the system keyring.
type KeyringStore struct {
	ring keyring.Keyring
}

// FileStore implements Store using encrypted files as a fallback.
type FileStore struct {
	baseDir string
}

// NewStore creates a new Store using the appropriate backend for the platform.
// On macOS, it uses Keychain. If the system keyring is unavailable, it falls
// back to encrypted file storage at ~/.config/goog/tokens/.
func NewStore() (Store, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Try to open the system keyring
	ring, err := openKeyring(configDir)
	if err != nil {
		// Fall back to file-based storage
		return NewFileStore(configDir)
	}

	return &KeyringStore{ring: ring}, nil
}

// NewFileStore creates a file-based Store at the specified directory.
// This is used as a fallback when the system keyring is unavailable.
func NewFileStore(baseDir string) (*FileStore, error) {
	tokensDir := filepath.Join(baseDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create tokens directory: %w", err)
	}
	return &FileStore{baseDir: baseDir}, nil
}

// openKeyring attempts to open the system keyring with appropriate configuration.
func openKeyring(configDir string) (keyring.Keyring, error) {
	backends := []keyring.BackendType{}

	switch runtime.GOOS {
	case "darwin":
		backends = append(backends, keyring.KeychainBackend)
	case "linux":
		backends = append(backends, keyring.SecretServiceBackend)
	case "windows":
		backends = append(backends, keyring.WinCredBackend)
	}

	// Always add file backend as final fallback
	backends = append(backends, keyring.FileBackend)

	// Derive machine-specific password for file backend
	machinePassword := deriveMachinePassword()

	cfg := keyring.Config{
		ServiceName:                    ServiceName,
		AllowedBackends:                backends,
		FileDir:                        filepath.Join(configDir, "keyring"),
		FilePasswordFunc:               keyring.FixedStringPrompt(machinePassword),
		KeychainTrustApplication:       true,
		KeychainSynchronizable:         false,
		KeychainAccessibleWhenUnlocked: true,
	}

	return keyring.Open(cfg)
}

// deriveMachinePassword creates a machine-specific password by combining
// hostname and user information. This ensures the password is unique per machine
// and user, making stolen keyring files harder to decrypt on other machines.
func deriveMachinePassword() string {
	var components []string

	// Add hostname
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		components = append(components, hostname)
	}

	// Add current user info
	if currentUser, err := user.Current(); err == nil {
		if currentUser.Username != "" {
			components = append(components, currentUser.Username)
		}
		if currentUser.Uid != "" {
			components = append(components, currentUser.Uid)
		}
		if currentUser.HomeDir != "" {
			components = append(components, currentUser.HomeDir)
		}
	}

	// Add a constant prefix to make the intent clear
	components = append([]string{"go-goog-cli-keyring"}, components...)

	// Combine all components and hash them
	combined := strings.Join(components, ":")
	hash := sha256.Sum256([]byte(combined))

	// Return hex-encoded hash as the password
	return fmt.Sprintf("%x", hash)
}

// getConfigDir returns the configuration directory for the application.
func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "goog"), nil
}

// formatKey creates a namespaced key from account and key name.
// Format: "goog:<account>:<key>"
func formatKey(account, key string) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefix, account, key)
}

// parseKey extracts account and key name from a namespaced key.
// Returns (account, key, ok) where ok is false if the key format is invalid.
func parseKey(fullKey string) (account, key string, ok bool) {
	parts := strings.SplitN(fullKey, ":", 3)
	if len(parts) != 3 {
		return "", "", false
	}
	if parts[0] != keyPrefix {
		return "", "", false
	}
	return parts[1], parts[2], true
}

// Set stores a value in the system keyring.
func (s *KeyringStore) Set(account, key string, value []byte) error {
	fullKey := formatKey(account, key)
	item := keyring.Item{
		Key:  fullKey,
		Data: value,
	}
	return s.ring.Set(item)
}

// Get retrieves a value from the system keyring.
func (s *KeyringStore) Get(account, key string) ([]byte, error) {
	fullKey := formatKey(account, key)
	item, err := s.ring.Get(fullKey)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return item.Data, nil
}

// Delete removes a value from the system keyring.
func (s *KeyringStore) Delete(account, key string) error {
	fullKey := formatKey(account, key)
	err := s.ring.Remove(fullKey)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return nil // Idempotent delete
		}
		return err
	}
	return nil
}

// List returns all keys stored for the given account.
func (s *KeyringStore) List(account string) ([]string, error) {
	keys, err := s.ring.Keys()
	if err != nil {
		return nil, err
	}

	prefix := formatKey(account, "")
	var result []string
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			_, keyName, ok := parseKey(k)
			if ok {
				result = append(result, keyName)
			}
		}
	}
	return result, nil
}

// tokenData represents the structure of encrypted token files.
type tokenData struct {
	Tokens map[string][]byte `json:"tokens"`
}

// encryptedFile represents the structure of the encrypted file on disk,
// including the salt used for key derivation.
type encryptedFile struct {
	Salt       []byte `json:"salt"`       // Random salt for PBKDF2
	Ciphertext []byte `json:"ciphertext"` // Encrypted token data
}

// Set stores a value in an encrypted file.
func (s *FileStore) Set(account, key string, value []byte) error {
	data, err := s.loadTokenData(account)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to load token data: %w", err)
	}
	if data == nil {
		data = &tokenData{Tokens: make(map[string][]byte)}
	}

	data.Tokens[key] = value
	return s.saveTokenData(account, data)
}

// Get retrieves a value from an encrypted file.
func (s *FileStore) Get(account, key string) ([]byte, error) {
	data, err := s.loadTokenData(account)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("failed to load token data: %w", err)
	}

	value, ok := data.Tokens[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

// Delete removes a value from an encrypted file.
func (s *FileStore) Delete(account, key string) error {
	data, err := s.loadTokenData(account)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // Idempotent delete
		}
		return fmt.Errorf("failed to load token data: %w", err)
	}

	delete(data.Tokens, key)

	if len(data.Tokens) == 0 {
		// Remove the file if no tokens remain
		return os.Remove(s.tokenFilePath(account))
	}

	return s.saveTokenData(account, data)
}

// List returns all keys stored for the given account.
func (s *FileStore) List(account string) ([]string, error) {
	data, err := s.loadTokenData(account)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to load token data: %w", err)
	}

	keys := make([]string, 0, len(data.Tokens))
	for k := range data.Tokens {
		keys = append(keys, k)
	}
	return keys, nil
}

// tokenFilePath returns the path to the token file for the given account.
func (s *FileStore) tokenFilePath(account string) string {
	return filepath.Join(s.baseDir, "tokens", account+".enc")
}

// loadTokenData loads and decrypts token data from a file.
func (s *FileStore) loadTokenData(account string) (*tokenData, error) {
	filePath := s.tokenFilePath(account)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the encrypted file structure
	var encFile encryptedFile
	if err := json.Unmarshal(fileData, &encFile); err != nil {
		// Try legacy format (raw encrypted data without salt)
		return s.loadLegacyTokenData(fileData, account)
	}

	// Validate salt
	if len(encFile.Salt) != saltSize {
		return nil, errors.New("invalid salt in encrypted file")
	}

	// Derive key using PBKDF2 with the stored salt
	key := s.deriveKey(account, encFile.Salt)

	plaintext, err := decrypt(encFile.Ciphertext, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token data: %w", err)
	}

	var data tokenData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &data, nil
}

// loadLegacyTokenData handles loading token data encrypted with the old format
// (simple SHA256 key derivation without salt). This provides backward compatibility.
func (s *FileStore) loadLegacyTokenData(encryptedData []byte, account string) (*tokenData, error) {
	// Use legacy key derivation
	key := s.deriveLegacyKey(account)

	plaintext, err := decrypt(encryptedData, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token data: %w", err)
	}

	var data tokenData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &data, nil
}

// saveTokenData encrypts and saves token data to a file.
func (s *FileStore) saveTokenData(account string, data *tokenData) error {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Generate a new random salt for each save
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using PBKDF2 with the new salt
	key := s.deriveKey(account, salt)

	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data: %w", err)
	}

	// Create the encrypted file structure
	encFile := encryptedFile{
		Salt:       salt,
		Ciphertext: ciphertext,
	}

	fileData, err := json.Marshal(encFile)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted file: %w", err)
	}

	filePath := s.tokenFilePath(account)
	return os.WriteFile(filePath, fileData, 0600)
}

// deriveKey derives an encryption key using PBKDF2 with machine-specific info.
// The salt is stored alongside the encrypted data to allow decryption.
// Uses 100,000 iterations of PBKDF2-HMAC-SHA256 for key stretching.
func (s *FileStore) deriveKey(account string, salt []byte) []byte {
	// Combine account with machine-specific information
	machineInfo := getMachineInfo()
	input := fmt.Sprintf("go-goog-cli-file-store:%s:%s", account, machineInfo)

	// Use PBKDF2 with SHA256 for key stretching
	// 32 bytes = 256 bits for AES-256
	return pbkdf2.Key([]byte(input), salt, pbkdf2Iterations, 32, sha256.New)
}

// deriveLegacyKey provides backward compatibility with the old key derivation.
// This uses simple SHA256 hashing without salt or iterations.
func (s *FileStore) deriveLegacyKey(account string) []byte {
	input := fmt.Sprintf("go-goog-cli-file-store:%s", account)
	hash := sha256.Sum256([]byte(input))
	return hash[:]
}

// getMachineInfo returns a string combining machine-specific identifiers.
// This makes encrypted files harder to use if copied to another machine.
func getMachineInfo() string {
	var components []string

	// Add hostname
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		components = append(components, hostname)
	}

	// Add current user info
	if currentUser, err := user.Current(); err == nil {
		if currentUser.Username != "" {
			components = append(components, currentUser.Username)
		}
		if currentUser.Uid != "" {
			components = append(components, currentUser.Uid)
		}
	}

	return strings.Join(components, ":")
}

// encrypt encrypts plaintext using AES-GCM.
func encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// decrypt decrypts ciphertext using AES-GCM.
func decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
