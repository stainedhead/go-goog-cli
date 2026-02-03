// Package auth provides OAuth2/PKCE authentication for Google APIs.
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TestGenerateCodeVerifier tests that the code verifier meets PKCE requirements.
func TestGenerateCodeVerifier(t *testing.T) {
	verifier := GenerateCodeVerifier()

	t.Run("length is 43 characters (256 bits base64url)", func(t *testing.T) {
		// 32 bytes = 256 bits, base64url encoded without padding = 43 chars
		if len(verifier) != 43 {
			t.Errorf("expected verifier length 43, got %d", len(verifier))
		}
	})

	t.Run("contains only base64url characters", func(t *testing.T) {
		for _, c := range verifier {
			if !isBase64URLChar(c) {
				t.Errorf("verifier contains invalid character: %c", c)
			}
		}
	})

	t.Run("generates unique values", func(t *testing.T) {
		verifier2 := GenerateCodeVerifier()
		if verifier == verifier2 {
			t.Error("expected unique verifiers, got identical values")
		}
	})
}

// TestGenerateCodeChallenge tests the PKCE code challenge generation.
func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	challenge := GenerateCodeChallenge(verifier)

	t.Run("produces SHA256 hash of verifier", func(t *testing.T) {
		// Manually compute expected challenge
		h := sha256.Sum256([]byte(verifier))
		expected := base64.RawURLEncoding.EncodeToString(h[:])

		if challenge != expected {
			t.Errorf("expected challenge %q, got %q", expected, challenge)
		}
	})

	t.Run("contains only base64url characters", func(t *testing.T) {
		for _, c := range challenge {
			if !isBase64URLChar(c) {
				t.Errorf("challenge contains invalid character: %c", c)
			}
		}
	})

	t.Run("length is 43 characters", func(t *testing.T) {
		// SHA256 = 32 bytes, base64url without padding = 43 chars
		if len(challenge) != 43 {
			t.Errorf("expected challenge length 43, got %d", len(challenge))
		}
	})
}

// TestNewOAuthConfig tests OAuth2 configuration creation.
func TestNewOAuthConfig(t *testing.T) {
	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	scopes := []string{ScopeGmailReadonly, ScopeCalendarReadonly}
	cfg := NewOAuthConfig(scopes)

	t.Run("client ID from environment", func(t *testing.T) {
		if cfg.ClientID != "test-client-id" {
			t.Errorf("expected client ID 'test-client-id', got %q", cfg.ClientID)
		}
	})

	t.Run("client secret from environment", func(t *testing.T) {
		if cfg.ClientSecret != "test-client-secret" {
			t.Errorf("expected client secret 'test-client-secret', got %q", cfg.ClientSecret)
		}
	})

	t.Run("scopes are set", func(t *testing.T) {
		if len(cfg.Scopes) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(cfg.Scopes))
		}
		if cfg.Scopes[0] != ScopeGmailReadonly {
			t.Errorf("expected first scope %q, got %q", ScopeGmailReadonly, cfg.Scopes[0])
		}
	})

	t.Run("uses Google endpoints", func(t *testing.T) {
		if !strings.Contains(cfg.Endpoint.AuthURL, "accounts.google.com") {
			t.Errorf("expected Google auth URL, got %q", cfg.Endpoint.AuthURL)
		}
		if !strings.Contains(cfg.Endpoint.TokenURL, "oauth2.googleapis.com") {
			t.Errorf("expected Google token URL, got %q", cfg.Endpoint.TokenURL)
		}
	})

	t.Run("redirect URI is localhost", func(t *testing.T) {
		if !strings.HasPrefix(cfg.RedirectURL, "http://localhost:") {
			t.Errorf("expected localhost redirect URL, got %q", cfg.RedirectURL)
		}
	})
}

// TestGetAuthorizationURL tests authorization URL generation.
func TestGetAuthorizationURL(t *testing.T) {
	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	scopes := []string{ScopeGmailReadonly}
	cfg := NewOAuthConfig(scopes)

	state := "test-state-123"
	codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	authURL := GetAuthorizationURL(cfg, state, codeChallenge)

	parsedURL, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("failed to parse auth URL: %v", err)
	}

	t.Run("URL points to Google", func(t *testing.T) {
		if parsedURL.Host != "accounts.google.com" {
			t.Errorf("expected host 'accounts.google.com', got %q", parsedURL.Host)
		}
	})

	t.Run("contains state parameter", func(t *testing.T) {
		gotState := parsedURL.Query().Get("state")
		if gotState != state {
			t.Errorf("expected state %q, got %q", state, gotState)
		}
	})

	t.Run("contains code_challenge", func(t *testing.T) {
		gotChallenge := parsedURL.Query().Get("code_challenge")
		if gotChallenge != codeChallenge {
			t.Errorf("expected code_challenge %q, got %q", codeChallenge, gotChallenge)
		}
	})

	t.Run("code_challenge_method is S256", func(t *testing.T) {
		method := parsedURL.Query().Get("code_challenge_method")
		if method != "S256" {
			t.Errorf("expected code_challenge_method 'S256', got %q", method)
		}
	})

	t.Run("response_type is code", func(t *testing.T) {
		responseType := parsedURL.Query().Get("response_type")
		if responseType != "code" {
			t.Errorf("expected response_type 'code', got %q", responseType)
		}
	})

	t.Run("contains redirect_uri", func(t *testing.T) {
		redirectURI := parsedURL.Query().Get("redirect_uri")
		if redirectURI == "" {
			t.Error("expected redirect_uri to be present")
		}
	})

	t.Run("contains scope", func(t *testing.T) {
		scope := parsedURL.Query().Get("scope")
		if scope == "" {
			t.Error("expected scope to be present")
		}
	})

	t.Run("access_type is offline", func(t *testing.T) {
		accessType := parsedURL.Query().Get("access_type")
		if accessType != "offline" {
			t.Errorf("expected access_type 'offline', got %q", accessType)
		}
	})
}

// TestCallbackServer tests the local callback server.
func TestCallbackServer(t *testing.T) {
	t.Run("extracts authorization code from callback", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Start the callback server
		codeChan := make(chan string, 1)
		errChan := make(chan error, 1)

		server, serverURL, err := StartCallbackServer(ctx, 0) // 0 = random port
		if err != nil {
			t.Fatalf("failed to start callback server: %v", err)
		}

		// Simulate the OAuth callback in a goroutine
		go func() {
			code, err := WaitForCallback(ctx, server)
			if err != nil {
				errChan <- err
				return
			}
			codeChan <- code
		}()

		// Give the server a moment to start listening
		time.Sleep(100 * time.Millisecond)

		// Simulate OAuth callback
		callbackURL := serverURL + "/callback?code=test-auth-code&state=test-state"
		resp, err := http.Get(callbackURL)
		if err != nil {
			t.Fatalf("failed to make callback request: %v", err)
		}
		resp.Body.Close()

		// Wait for result
		select {
		case code := <-codeChan:
			if code != "test-auth-code" {
				t.Errorf("expected code 'test-auth-code', got %q", code)
			}
		case err := <-errChan:
			t.Fatalf("callback server error: %v", err)
		case <-ctx.Done():
			t.Fatal("timeout waiting for callback")
		}
	})

	t.Run("returns error on OAuth error response", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server, serverURL, err := StartCallbackServer(ctx, 0)
		if err != nil {
			t.Fatalf("failed to start callback server: %v", err)
		}

		errChan := make(chan error, 1)
		go func() {
			_, err := WaitForCallback(ctx, server)
			errChan <- err
		}()

		time.Sleep(100 * time.Millisecond)

		// Simulate OAuth error callback
		callbackURL := serverURL + "/callback?error=access_denied&error_description=User+denied+access"
		resp, err := http.Get(callbackURL)
		if err != nil {
			t.Fatalf("failed to make callback request: %v", err)
		}
		resp.Body.Close()

		select {
		case err := <-errChan:
			if err == nil {
				t.Error("expected error for OAuth error response")
			}
			if !strings.Contains(err.Error(), "access_denied") {
				t.Errorf("expected error to contain 'access_denied', got %q", err.Error())
			}
		case <-ctx.Done():
			t.Fatal("timeout waiting for error")
		}
	})
}

// TestExchangeCode tests the code exchange functionality.
func TestExchangeCode(t *testing.T) {
	// Create a mock OAuth server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}

		// Verify PKCE code_verifier is present
		codeVerifier := r.FormValue("code_verifier")
		if codeVerifier == "" {
			t.Error("expected code_verifier in request")
		}

		// Return mock token response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "mock-access-token",
			"token_type": "Bearer",
			"refresh_token": "mock-refresh-token",
			"expires_in": 3600
		}`))
	}))
	defer mockServer.Close()

	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	cfg := NewOAuthConfig([]string{ScopeGmailReadonly})
	// Override the token URL to use mock server
	cfg.Endpoint.TokenURL = mockServer.URL

	ctx := context.Background()
	token, err := ExchangeCode(ctx, cfg, "test-auth-code", "test-code-verifier")
	if err != nil {
		t.Fatalf("ExchangeCode failed: %v", err)
	}

	t.Run("returns access token", func(t *testing.T) {
		if token.AccessToken != "mock-access-token" {
			t.Errorf("expected access token 'mock-access-token', got %q", token.AccessToken)
		}
	})

	t.Run("returns refresh token", func(t *testing.T) {
		if token.RefreshToken != "mock-refresh-token" {
			t.Errorf("expected refresh token 'mock-refresh-token', got %q", token.RefreshToken)
		}
	})

	t.Run("returns token type", func(t *testing.T) {
		if token.TokenType != "Bearer" {
			t.Errorf("expected token type 'Bearer', got %q", token.TokenType)
		}
	})

	t.Run("token has valid expiry", func(t *testing.T) {
		if token.Expiry.IsZero() {
			t.Error("expected token to have valid expiry time")
		}
		if token.Expiry.Before(time.Now()) {
			t.Error("expected token expiry to be in the future")
		}
	})
}

// TestGoogleScopes tests that scope constants are correctly defined.
func TestGoogleScopes(t *testing.T) {
	testCases := []struct {
		name     string
		scope    string
		contains string
	}{
		{"Gmail readonly", ScopeGmailReadonly, "gmail.readonly"},
		{"Gmail send", ScopeGmailSend, "gmail.send"},
		{"Gmail modify", ScopeGmailModify, "gmail.modify"},
		{"Calendar readonly", ScopeCalendarReadonly, "calendar.readonly"},
		{"Calendar events", ScopeCalendarEvents, "calendar.events"},
		{"Drive readonly", ScopeDriveReadonly, "drive.readonly"},
		{"Drive files", ScopeDriveFile, "drive.file"},
		{"User info email", ScopeUserInfoEmail, "userinfo.email"},
		{"User info profile", ScopeUserInfoProfile, "userinfo.profile"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !strings.Contains(tc.scope, tc.contains) {
				t.Errorf("expected scope to contain %q, got %q", tc.contains, tc.scope)
			}
			if !strings.HasPrefix(tc.scope, "https://www.googleapis.com/auth/") {
				t.Errorf("expected scope to start with Google API prefix, got %q", tc.scope)
			}
		})
	}
}

// Helper functions

func isBase64URLChar(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_'
}

func getEnvOrDefault(key, defaultValue string) string {
	if v, ok := lookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func setEnvForTest(key, value string) {
	if value == "" {
		unsetEnv(key)
	} else {
		setEnv(key, value)
	}
}
