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
