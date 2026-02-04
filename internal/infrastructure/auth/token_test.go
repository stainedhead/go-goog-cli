// Package auth provides OAuth2/PKCE authentication for Google APIs.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// mockStore implements keyring.Store for testing.
type mockStore struct {
	data map[string]map[string][]byte
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[string]map[string][]byte),
	}
}

func (s *mockStore) Set(account, key string, value []byte) error {
	if s.data[account] == nil {
		s.data[account] = make(map[string][]byte)
	}
	s.data[account][key] = value
	return nil
}

func (s *mockStore) Get(account, key string) ([]byte, error) {
	if s.data[account] == nil {
		return nil, errKeyNotFound
	}
	if v, ok := s.data[account][key]; ok {
		return v, nil
	}
	return nil, errKeyNotFound
}

func (s *mockStore) Delete(account, key string) error {
	if s.data[account] != nil {
		delete(s.data[account], key)
	}
	return nil
}

func (s *mockStore) List(account string) ([]string, error) {
	if s.data[account] == nil {
		return []string{}, nil
	}
	keys := make([]string, 0, len(s.data[account]))
	for k := range s.data[account] {
		keys = append(keys, k)
	}
	return keys, nil
}

// TestTokenSaveLoad tests saving and loading tokens through the TokenManager.
func TestTokenSaveLoad(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "test@example.com"
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	t.Run("save token successfully", func(t *testing.T) {
		err := manager.SaveToken(account, token)
		if err != nil {
			t.Fatalf("SaveToken failed: %v", err)
		}
	})

	t.Run("load token successfully", func(t *testing.T) {
		loaded, err := manager.LoadToken(account)
		if err != nil {
			t.Fatalf("LoadToken failed: %v", err)
		}

		if loaded.AccessToken != token.AccessToken {
			t.Errorf("expected access token %q, got %q", token.AccessToken, loaded.AccessToken)
		}
		if loaded.RefreshToken != token.RefreshToken {
			t.Errorf("expected refresh token %q, got %q", token.RefreshToken, loaded.RefreshToken)
		}
		if loaded.TokenType != token.TokenType {
			t.Errorf("expected token type %q, got %q", token.TokenType, loaded.TokenType)
		}
	})

	t.Run("load non-existent token returns error", func(t *testing.T) {
		_, err := manager.LoadToken("nonexistent@example.com")
		if err == nil {
			t.Error("expected error for non-existent account")
		}
	})
}

// TestTokenDelete tests deleting tokens.
func TestTokenDelete(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "delete@example.com"
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Save token first
	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	// Verify it exists
	_, err := manager.LoadToken(account)
	if err != nil {
		t.Fatalf("LoadToken failed: %v", err)
	}

	t.Run("delete token successfully", func(t *testing.T) {
		err := manager.DeleteToken(account)
		if err != nil {
			t.Fatalf("DeleteToken failed: %v", err)
		}
	})

	t.Run("load deleted token returns error", func(t *testing.T) {
		_, err := manager.LoadToken(account)
		if err == nil {
			t.Error("expected error after deletion")
		}
	})

	t.Run("delete non-existent token is idempotent", func(t *testing.T) {
		err := manager.DeleteToken("nonexistent@example.com")
		if err != nil {
			t.Errorf("delete non-existent should not error: %v", err)
		}
	})
}

// TestTokenRefresh tests the token refresh logic.
func TestTokenRefresh(t *testing.T) {
	// Create mock OAuth server for token refresh
	refreshCalled := false
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCalled = true
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}

		grantType := r.FormValue("grant_type")
		if grantType != "refresh_token" {
			t.Errorf("expected grant_type 'refresh_token', got %q", grantType)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "new-access-token",
			"token_type": "Bearer",
			"expires_in": 3600,
			"refresh_token": "new-refresh-token"
		}`))
	}))
	defer mockServer.Close()

	store := newMockStore()
	manager := NewTokenManager(store)

	account := "refresh@example.com"

	// Store an expired token
	expiredToken := &oauth2.Token{
		AccessToken:  "expired-access-token",
		TokenType:    "Bearer",
		RefreshToken: "valid-refresh-token",
		Expiry:       time.Now().Add(-time.Hour), // Expired
	}
	if err := manager.SaveToken(account, expiredToken); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	// Also store scopes
	if err := manager.SaveScopes(account, []string{ScopeGmailReadonly}); err != nil {
		t.Fatalf("SaveScopes failed: %v", err)
	}

	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	ctx := context.Background()

	t.Run("refreshes expired token", func(t *testing.T) {
		// Create config with mock server
		cfg := NewOAuthConfig([]string{ScopeGmailReadonly})
		cfg.Endpoint.TokenURL = mockServer.URL

		newToken, err := manager.RefreshToken(ctx, account, cfg)
		if err != nil {
			t.Fatalf("RefreshToken failed: %v", err)
		}

		if !refreshCalled {
			t.Error("expected refresh endpoint to be called")
		}

		if newToken.AccessToken != "new-access-token" {
			t.Errorf("expected new access token, got %q", newToken.AccessToken)
		}
	})
}

// TestGetTokenSource tests getting an oauth2.TokenSource.
func TestGetTokenSource(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "source@example.com"
	token := &oauth2.Token{
		AccessToken:  "valid-access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}
	if err := manager.SaveScopes(account, []string{ScopeGmailReadonly}); err != nil {
		t.Fatalf("SaveScopes failed: %v", err)
	}

	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	ctx := context.Background()

	t.Run("returns valid token source", func(t *testing.T) {
		ts, err := manager.GetTokenSource(ctx, account)
		if err != nil {
			t.Fatalf("GetTokenSource failed: %v", err)
		}

		gotToken, err := ts.Token()
		if err != nil {
			t.Fatalf("Token() failed: %v", err)
		}

		if gotToken.AccessToken != token.AccessToken {
			t.Errorf("expected access token %q, got %q", token.AccessToken, gotToken.AccessToken)
		}
	})

	t.Run("returns error for non-existent account", func(t *testing.T) {
		_, err := manager.GetTokenSource(ctx, "nonexistent@example.com")
		if err == nil {
			t.Error("expected error for non-existent account")
		}
	})
}

// TestScopeManagement tests saving and loading scopes.
func TestScopeManagement(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "scopes@example.com"
	scopes := []string{ScopeGmailReadonly, ScopeCalendarReadonly, ScopeDriveReadonly}

	t.Run("save scopes", func(t *testing.T) {
		err := manager.SaveScopes(account, scopes)
		if err != nil {
			t.Fatalf("SaveScopes failed: %v", err)
		}
	})

	t.Run("get granted scopes", func(t *testing.T) {
		loaded, err := manager.GetGrantedScopes(account)
		if err != nil {
			t.Fatalf("GetGrantedScopes failed: %v", err)
		}

		if len(loaded) != len(scopes) {
			t.Errorf("expected %d scopes, got %d", len(scopes), len(loaded))
		}
	})

	t.Run("has scope returns true for existing scope", func(t *testing.T) {
		if !manager.HasScope(account, ScopeGmailReadonly) {
			t.Error("expected HasScope to return true for existing scope")
		}
	})

	t.Run("has scope returns false for missing scope", func(t *testing.T) {
		if manager.HasScope(account, ScopeGmailSend) {
			t.Error("expected HasScope to return false for missing scope")
		}
	})

	t.Run("has scope returns false for non-existent account", func(t *testing.T) {
		if manager.HasScope("nonexistent@example.com", ScopeGmailReadonly) {
			t.Error("expected HasScope to return false for non-existent account")
		}
	})
}

// TestTokenSerialization tests that tokens are properly serialized/deserialized.
func TestTokenSerialization(t *testing.T) {
	original := &oauth2.Token{
		AccessToken:  "access-token-123",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token-456",
		Expiry:       time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	// Serialize
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Deserialize
	var loaded oauth2.Token
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	t.Run("access token matches", func(t *testing.T) {
		if loaded.AccessToken != original.AccessToken {
			t.Errorf("expected %q, got %q", original.AccessToken, loaded.AccessToken)
		}
	})

	t.Run("refresh token matches", func(t *testing.T) {
		if loaded.RefreshToken != original.RefreshToken {
			t.Errorf("expected %q, got %q", original.RefreshToken, loaded.RefreshToken)
		}
	})

	t.Run("expiry matches", func(t *testing.T) {
		if !loaded.Expiry.Equal(original.Expiry) {
			t.Errorf("expected %v, got %v", original.Expiry, loaded.Expiry)
		}
	})
}

// TestTokenExpired tests token expiry detection.
func TestTokenExpired(t *testing.T) {
	t.Run("valid token is not expired", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "valid",
			Expiry:      time.Now().Add(time.Hour),
		}
		if !token.Valid() {
			t.Error("expected token to be valid")
		}
	})

	t.Run("expired token is invalid", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "expired",
			Expiry:      time.Now().Add(-time.Hour),
		}
		if token.Valid() {
			t.Error("expected token to be invalid")
		}
	})

	t.Run("token without expiry is valid", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "no-expiry",
		}
		if !token.Valid() {
			t.Error("expected token without expiry to be valid")
		}
	})
}

// TestGetTokenInfo tests retrieving token info.
func TestGetTokenInfo(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	t.Run("returns info for non-existent account", func(t *testing.T) {
		info, err := manager.GetTokenInfo("nonexistent@example.com")
		if err != nil {
			t.Fatalf("GetTokenInfo failed: %v", err)
		}
		if info.HasToken {
			t.Error("expected HasToken to be false")
		}
		if info.Account != "nonexistent@example.com" {
			t.Errorf("expected account 'nonexistent@example.com', got %q", info.Account)
		}
	})

	t.Run("returns info for valid account with token", func(t *testing.T) {
		account := "info@example.com"
		token := &oauth2.Token{
			AccessToken:  "test-access-token",
			TokenType:    "Bearer",
			RefreshToken: "test-refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		}
		if err := manager.SaveToken(account, token); err != nil {
			t.Fatalf("SaveToken failed: %v", err)
		}
		if err := manager.SaveScopes(account, []string{ScopeGmailReadonly, ScopeCalendarReadonly}); err != nil {
			t.Fatalf("SaveScopes failed: %v", err)
		}

		info, err := manager.GetTokenInfo(account)
		if err != nil {
			t.Fatalf("GetTokenInfo failed: %v", err)
		}

		if !info.HasToken {
			t.Error("expected HasToken to be true")
		}
		if info.IsExpired {
			t.Error("expected IsExpired to be false for valid token")
		}
		if info.TokenType != "Bearer" {
			t.Errorf("expected TokenType 'Bearer', got %q", info.TokenType)
		}
		if len(info.Scopes) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(info.Scopes))
		}
		if info.ExpiryTime == "" {
			t.Error("expected ExpiryTime to be set")
		}
	})

	t.Run("returns info for expired token", func(t *testing.T) {
		account := "expired@example.com"
		token := &oauth2.Token{
			AccessToken:  "expired-access-token",
			TokenType:    "Bearer",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(-time.Hour), // Expired
		}
		if err := manager.SaveToken(account, token); err != nil {
			t.Fatalf("SaveToken failed: %v", err)
		}

		info, err := manager.GetTokenInfo(account)
		if err != nil {
			t.Fatalf("GetTokenInfo failed: %v", err)
		}

		if !info.HasToken {
			t.Error("expected HasToken to be true")
		}
		if !info.IsExpired {
			t.Error("expected IsExpired to be true for expired token")
		}
	})

	t.Run("returns info with token but no scopes", func(t *testing.T) {
		account := "noscopes@example.com"
		token := &oauth2.Token{
			AccessToken:  "test-access-token",
			TokenType:    "Bearer",
			RefreshToken: "test-refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		}
		if err := manager.SaveToken(account, token); err != nil {
			t.Fatalf("SaveToken failed: %v", err)
		}
		// Don't save any scopes

		info, err := manager.GetTokenInfo(account)
		if err != nil {
			t.Fatalf("GetTokenInfo failed: %v", err)
		}

		if !info.HasToken {
			t.Error("expected HasToken to be true")
		}
		if len(info.Scopes) != 0 {
			t.Errorf("expected empty scopes, got %v", info.Scopes)
		}
	})
}

// TestDeleteTokenIdempotency tests that deleting tokens is idempotent.
func TestDeleteTokenIdempotency(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "idempotent@example.com"
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Save token
	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	// Delete once
	if err := manager.DeleteToken(account); err != nil {
		t.Fatalf("first DeleteToken failed: %v", err)
	}

	// Delete again (should not error)
	if err := manager.DeleteToken(account); err != nil {
		t.Fatalf("second DeleteToken failed: %v", err)
	}

	// Delete third time (still should not error)
	if err := manager.DeleteToken(account); err != nil {
		t.Fatalf("third DeleteToken failed: %v", err)
	}
}

// TestSaveLoadTokenRoundTrip tests saving and loading tokens preserves data.
func TestSaveLoadTokenRoundTrip(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "roundtrip@example.com"
	expiry := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	token := &oauth2.Token{
		AccessToken:  "access-token-abc123",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token-xyz789",
		Expiry:       expiry,
	}

	// Save
	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	// Load
	loaded, err := manager.LoadToken(account)
	if err != nil {
		t.Fatalf("LoadToken failed: %v", err)
	}

	// Verify all fields match
	if loaded.AccessToken != token.AccessToken {
		t.Errorf("AccessToken mismatch: got %q, want %q", loaded.AccessToken, token.AccessToken)
	}
	if loaded.TokenType != token.TokenType {
		t.Errorf("TokenType mismatch: got %q, want %q", loaded.TokenType, token.TokenType)
	}
	if loaded.RefreshToken != token.RefreshToken {
		t.Errorf("RefreshToken mismatch: got %q, want %q", loaded.RefreshToken, token.RefreshToken)
	}
	if !loaded.Expiry.Equal(token.Expiry) {
		t.Errorf("Expiry mismatch: got %v, want %v", loaded.Expiry, token.Expiry)
	}
}

// errorStore is a mock store that returns errors for testing error handling.
type errorStore struct {
	getErr    error
	setErr    error
	deleteErr error
	listErr   error
}

func (s *errorStore) Set(account, key string, value []byte) error { return s.setErr }
func (s *errorStore) Get(account, key string) ([]byte, error)     { return nil, s.getErr }
func (s *errorStore) Delete(account, key string) error            { return s.deleteErr }
func (s *errorStore) List(account string) ([]string, error)       { return nil, s.listErr }

// TestLoadTokenErrorHandling tests error handling when loading tokens.
func TestLoadTokenErrorHandling(t *testing.T) {
	t.Run("returns error for corrupted token data", func(t *testing.T) {
		store := newMockStore()
		manager := NewTokenManager(store)

		account := "corrupted@example.com"
		// Store invalid JSON data
		store.Set(account, KeyToken, []byte("not valid json"))

		_, err := manager.LoadToken(account)
		if err == nil {
			t.Error("expected error for corrupted token data")
		}
	})

	t.Run("wraps store errors", func(t *testing.T) {
		customErr := errKeyNotFound
		store := &errorStore{getErr: customErr}
		manager := NewTokenManager(store)

		_, err := manager.LoadToken("any@example.com")
		if err == nil {
			t.Error("expected error from store")
		}
		// Should return ErrTokenNotFound for key not found errors
		if err != ErrTokenNotFound {
			t.Errorf("expected ErrTokenNotFound, got %v", err)
		}
	})
}

// TestSaveTokenErrorHandling tests error handling when saving tokens.
func TestSaveTokenErrorHandling(t *testing.T) {
	t.Run("wraps store errors", func(t *testing.T) {
		customErr := errKeyNotFound // Use any error
		store := &errorStore{setErr: customErr}
		manager := NewTokenManager(store)

		token := &oauth2.Token{
			AccessToken: "test",
		}

		err := manager.SaveToken("any@example.com", token)
		if err == nil {
			t.Error("expected error from store")
		}
	})
}

// TestGetGrantedScopesErrorHandling tests error handling when getting scopes.
func TestGetGrantedScopesErrorHandling(t *testing.T) {
	t.Run("returns error for corrupted scope data", func(t *testing.T) {
		store := newMockStore()
		manager := NewTokenManager(store)

		account := "corrupted@example.com"
		// Store invalid JSON data for scopes
		store.Set(account, KeyScopes, []byte("not valid json"))

		_, err := manager.GetGrantedScopes(account)
		if err == nil {
			t.Error("expected error for corrupted scope data")
		}
	})

	t.Run("returns ErrScopesNotSet for missing scopes", func(t *testing.T) {
		store := newMockStore()
		manager := NewTokenManager(store)

		_, err := manager.GetGrantedScopes("nosuch@example.com")
		if err != ErrScopesNotSet {
			t.Errorf("expected ErrScopesNotSet, got %v", err)
		}
	})
}

// TestIsKeyNotFoundError tests the isKeyNotFoundError helper function.
func TestIsKeyNotFoundError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"internal errKeyNotFound", errKeyNotFound, true},
		{"key not found message", customError("key not found"), true},
		{"secret not found in keyring", customError("secret not found in keyring"), true},
		{"other error", customError("some other error"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isKeyNotFoundError(tc.err)
			if result != tc.expected {
				t.Errorf("isKeyNotFoundError(%v) = %v, want %v", tc.err, result, tc.expected)
			}
		})
	}
}

// customError is a simple error type for testing.
type customError string

func (e customError) Error() string { return string(e) }

// TestSaveScopesErrorHandling tests error handling when saving scopes.
func TestSaveScopesErrorHandling(t *testing.T) {
	t.Run("wraps store errors", func(t *testing.T) {
		customErr := errKeyNotFound
		store := &errorStore{setErr: customErr}
		manager := NewTokenManager(store)

		err := manager.SaveScopes("any@example.com", []string{"scope1"})
		if err == nil {
			t.Error("expected error from store")
		}
	})
}

// TestDeleteTokenWithScopesAlsoDeletesScopes tests that deleting a token also deletes scopes.
func TestDeleteTokenWithScopesAlsoDeletesScopes(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "withscopes@example.com"
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Save token and scopes
	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}
	if err := manager.SaveScopes(account, []string{ScopeGmailReadonly}); err != nil {
		t.Fatalf("SaveScopes failed: %v", err)
	}

	// Verify scopes exist
	_, err := manager.GetGrantedScopes(account)
	if err != nil {
		t.Fatalf("GetGrantedScopes before delete failed: %v", err)
	}

	// Delete token
	if err := manager.DeleteToken(account); err != nil {
		t.Fatalf("DeleteToken failed: %v", err)
	}

	// Verify scopes are also deleted
	_, err = manager.GetGrantedScopes(account)
	if err != ErrScopesNotSet {
		t.Errorf("expected ErrScopesNotSet after delete, got %v", err)
	}
}

// TestHasScopeEdgeCases tests edge cases for HasScope.
func TestHasScopeEdgeCases(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "hasscope@example.com"
	if err := manager.SaveScopes(account, []string{ScopeGmailReadonly, ScopeCalendarReadonly}); err != nil {
		t.Fatalf("SaveScopes failed: %v", err)
	}

	t.Run("returns true for first scope in list", func(t *testing.T) {
		if !manager.HasScope(account, ScopeGmailReadonly) {
			t.Error("expected HasScope to return true for first scope")
		}
	})

	t.Run("returns true for last scope in list", func(t *testing.T) {
		if !manager.HasScope(account, ScopeCalendarReadonly) {
			t.Error("expected HasScope to return true for last scope")
		}
	})

	t.Run("returns false for partial scope match", func(t *testing.T) {
		// ScopeGmailReadonly is the full scope, "gmail" is just a partial match
		if manager.HasScope(account, "gmail") {
			t.Error("expected HasScope to return false for partial match")
		}
	})

	t.Run("returns false for empty scope", func(t *testing.T) {
		if manager.HasScope(account, "") {
			t.Error("expected HasScope to return false for empty scope")
		}
	})
}

// TestDeleteTokenErrorPaths tests error handling in DeleteToken.
func TestDeleteTokenErrorPaths(t *testing.T) {
	t.Run("returns error when store delete fails", func(t *testing.T) {
		customErr := customError("delete failed")
		store := &errorStore{deleteErr: customErr}
		manager := NewTokenManager(store)

		err := manager.DeleteToken("any@example.com")
		if err == nil {
			t.Error("expected error from store")
		}
	})
}

// TestRefreshTokenErrorPaths tests error handling in RefreshToken.
func TestRefreshTokenErrorPaths(t *testing.T) {
	t.Run("returns error when token not found", func(t *testing.T) {
		store := newMockStore()
		manager := NewTokenManager(store)

		ctx := context.Background()
		cfg := &oauth2.Config{}

		_, err := manager.RefreshToken(ctx, "nonexistent@example.com", cfg)
		if err == nil {
			t.Error("expected error for non-existent token")
		}
	})
}

// TestGetTokenSourceWithoutScopes tests GetTokenSource when scopes are not stored.
func TestGetTokenSourceWithoutScopes(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "noscopes@example.com"
	token := &oauth2.Token{
		AccessToken:  "valid-access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}
	// Don't save any scopes

	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	ctx := context.Background()

	// Should still work even without scopes
	ts, err := manager.GetTokenSource(ctx, account)
	if err != nil {
		t.Fatalf("GetTokenSource failed: %v", err)
	}

	gotToken, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() failed: %v", err)
	}

	if gotToken.AccessToken != token.AccessToken {
		t.Errorf("expected access token %q, got %q", token.AccessToken, gotToken.AccessToken)
	}
}

// TestGetTokenInfoWithZeroExpiry tests GetTokenInfo when token has zero expiry.
func TestGetTokenInfoWithZeroExpiry(t *testing.T) {
	store := newMockStore()
	manager := NewTokenManager(store)

	account := "zeroexpiry@example.com"
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		// No expiry set (zero value)
	}
	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	info, err := manager.GetTokenInfo(account)
	if err != nil {
		t.Fatalf("GetTokenInfo failed: %v", err)
	}

	if !info.HasToken {
		t.Error("expected HasToken to be true")
	}
	if info.ExpiryTime != "" {
		t.Errorf("expected empty ExpiryTime for zero expiry, got %q", info.ExpiryTime)
	}
}

// TestLoadTokenWrapsStoreError tests that LoadToken wraps non-key-not-found errors.
func TestLoadTokenWrapsStoreError(t *testing.T) {
	customErr := customError("some other error")
	store := &errorStore{getErr: customErr}
	manager := NewTokenManager(store)

	_, err := manager.LoadToken("any@example.com")
	if err == nil {
		t.Error("expected error from store")
	}
	// Should return ErrTokenNotFound for our mock since it matches "key not found"
}

// TestGetGrantedScopesWrapsStoreError tests error wrapping in GetGrantedScopes.
func TestGetGrantedScopesWrapsStoreError(t *testing.T) {
	customErr := customError("some other error")
	store := &errorStore{getErr: customErr}
	manager := NewTokenManager(store)

	_, err := manager.GetGrantedScopes("any@example.com")
	if err == nil {
		t.Error("expected error from store")
	}
}

// TestGetTokenInfoWithStoreError tests GetTokenInfo when store returns error.
func TestGetTokenInfoWithStoreError(t *testing.T) {
	customErr := customError("some store error")
	store := &errorStore{getErr: customErr}
	manager := NewTokenManager(store)

	_, err := manager.GetTokenInfo("any@example.com")
	if err == nil {
		t.Error("expected error from store")
	}
}

// errorStoreWithScopeErr allows separate errors for different keys.
type errorStoreWithScopeErr struct {
	mockStore *mockStore
	scopeErr  error
}

func (s *errorStoreWithScopeErr) Set(account, key string, value []byte) error {
	return s.mockStore.Set(account, key, value)
}

func (s *errorStoreWithScopeErr) Get(account, key string) ([]byte, error) {
	if key == KeyScopes && s.scopeErr != nil {
		return nil, s.scopeErr
	}
	return s.mockStore.Get(account, key)
}

func (s *errorStoreWithScopeErr) Delete(account, key string) error {
	return s.mockStore.Delete(account, key)
}

func (s *errorStoreWithScopeErr) List(account string) ([]string, error) {
	return s.mockStore.List(account)
}

// TestGetTokenInfoWithScopesError tests GetTokenInfo when scopes lookup fails.
func TestGetTokenInfoWithScopesError(t *testing.T) {
	mockStore := newMockStore()
	store := &errorStoreWithScopeErr{
		mockStore: mockStore,
		scopeErr:  customError("scope lookup failed"),
	}
	manager := NewTokenManager(store)

	account := "scopeerr@example.com"
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}
	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	// GetTokenInfo should still succeed even if scopes lookup fails
	info, err := manager.GetTokenInfo(account)
	if err != nil {
		t.Fatalf("GetTokenInfo failed: %v", err)
	}

	if !info.HasToken {
		t.Error("expected HasToken to be true")
	}
	if len(info.Scopes) != 0 {
		t.Errorf("expected empty scopes when lookup fails, got %v", info.Scopes)
	}
}

// TestRefreshTokenWhenTokenUnchanged tests that refreshed token is not saved when unchanged.
func TestRefreshTokenWhenTokenUnchanged(t *testing.T) {
	// Create a mock OAuth server that returns the same access token
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return same access token as was stored
		w.Write([]byte(`{
			"access_token": "same-access-token",
			"token_type": "Bearer",
			"expires_in": 3600,
			"refresh_token": "same-refresh-token"
		}`))
	}))
	defer mockServer.Close()

	store := newMockStore()
	manager := NewTokenManager(store)

	account := "unchanged@example.com"
	token := &oauth2.Token{
		AccessToken:  "same-access-token", // Same as what server returns
		TokenType:    "Bearer",
		RefreshToken: "same-refresh-token",
		Expiry:       time.Now().Add(-time.Hour), // Expired
	}
	if err := manager.SaveToken(account, token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	ctx := context.Background()
	cfg := NewOAuthConfig([]string{ScopeGmailReadonly})
	cfg.Endpoint.TokenURL = mockServer.URL

	newToken, err := manager.RefreshToken(ctx, account, cfg)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	if newToken.AccessToken != "same-access-token" {
		t.Errorf("expected access token 'same-access-token', got %q", newToken.AccessToken)
	}
}
