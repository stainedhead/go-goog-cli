// Package auth provides OAuth2/PKCE authentication for Google APIs.
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
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

// TestNewOAuthConfigWithCredentials tests OAuth config creation with explicit credentials.
func TestNewOAuthConfigWithCredentials(t *testing.T) {
	clientID := "explicit-client-id"
	clientSecret := "explicit-client-secret"
	scopes := []string{ScopeGmailReadonly, ScopeCalendarReadonly}

	t.Run("uses provided credentials", func(t *testing.T) {
		cfg := NewOAuthConfigWithCredentials(clientID, clientSecret, scopes, 9000)
		if cfg.ClientID != clientID {
			t.Errorf("expected client ID %q, got %q", clientID, cfg.ClientID)
		}
		if cfg.ClientSecret != clientSecret {
			t.Errorf("expected client secret %q, got %q", clientSecret, cfg.ClientSecret)
		}
	})

	t.Run("uses custom port", func(t *testing.T) {
		cfg := NewOAuthConfigWithCredentials(clientID, clientSecret, scopes, 9000)
		if !strings.Contains(cfg.RedirectURL, ":9000") {
			t.Errorf("expected redirect URL to contain port 9000, got %q", cfg.RedirectURL)
		}
	})

	t.Run("uses default port when 0", func(t *testing.T) {
		cfg := NewOAuthConfigWithCredentials(clientID, clientSecret, scopes, 0)
		expectedPort := fmt.Sprintf(":%d", DefaultRedirectPort)
		if !strings.Contains(cfg.RedirectURL, expectedPort) {
			t.Errorf("expected redirect URL to contain default port %s, got %q", expectedPort, cfg.RedirectURL)
		}
	})

	t.Run("scopes are correctly set", func(t *testing.T) {
		cfg := NewOAuthConfigWithCredentials(clientID, clientSecret, scopes, 0)
		if len(cfg.Scopes) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(cfg.Scopes))
		}
	})
}

// TestNewOAuthConfigWithCustomPort tests OAuth config with custom redirect port.
func TestNewOAuthConfigWithCustomPort(t *testing.T) {
	// Save and restore env vars
	origClientID := getEnvOrDefault("GOOG_CLIENT_ID", "")
	origClientSecret := getEnvOrDefault("GOOG_CLIENT_SECRET", "")
	origPort := getEnvOrDefault("GOOG_REDIRECT_PORT", "")
	defer func() {
		setEnvForTest("GOOG_CLIENT_ID", origClientID)
		setEnvForTest("GOOG_CLIENT_SECRET", origClientSecret)
		setEnvForTest("GOOG_REDIRECT_PORT", origPort)
	}()

	setEnvForTest("GOOG_CLIENT_ID", "test-client-id")
	setEnvForTest("GOOG_CLIENT_SECRET", "test-client-secret")

	t.Run("uses custom port from environment", func(t *testing.T) {
		setEnvForTest("GOOG_REDIRECT_PORT", "9999")
		cfg := NewOAuthConfig([]string{ScopeGmailReadonly})
		if !strings.Contains(cfg.RedirectURL, ":9999") {
			t.Errorf("expected redirect URL to contain port 9999, got %q", cfg.RedirectURL)
		}
	})

	t.Run("uses default port when env not set", func(t *testing.T) {
		setEnvForTest("GOOG_REDIRECT_PORT", "")
		cfg := NewOAuthConfig([]string{ScopeGmailReadonly})
		expectedPort := fmt.Sprintf(":%d", DefaultRedirectPort)
		if !strings.Contains(cfg.RedirectURL, expectedPort) {
			t.Errorf("expected redirect URL to contain default port, got %q", cfg.RedirectURL)
		}
	})
}

// TestValidateConfig tests OAuth configuration validation.
func TestValidateConfig(t *testing.T) {
	t.Run("returns error for missing client ID", func(t *testing.T) {
		cfg := &oauth2.Config{
			ClientID:     "",
			ClientSecret: "secret",
		}
		err := ValidateConfig(cfg)
		if err == nil {
			t.Error("expected error for missing client ID")
		}
		if err != ErrMissingClientID {
			t.Errorf("expected ErrMissingClientID, got %v", err)
		}
	})

	t.Run("returns error for missing client secret", func(t *testing.T) {
		cfg := &oauth2.Config{
			ClientID:     "client-id",
			ClientSecret: "",
		}
		err := ValidateConfig(cfg)
		if err == nil {
			t.Error("expected error for missing client secret")
		}
		if err != ErrMissingClientSecret {
			t.Errorf("expected ErrMissingClientSecret, got %v", err)
		}
	})

	t.Run("returns nil for valid config", func(t *testing.T) {
		cfg := &oauth2.Config{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		}
		err := ValidateConfig(cfg)
		if err != nil {
			t.Errorf("expected no error for valid config, got %v", err)
		}
	})
}

// TestCallbackServerGetServerURL tests the GetServerURL method.
func TestCallbackServerGetServerURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server, serverURL, err := StartCallbackServer(ctx, 0)
	if err != nil {
		t.Fatalf("failed to start callback server: %v", err)
	}

	t.Run("returns correct server URL", func(t *testing.T) {
		gotURL := server.GetServerURL()
		if gotURL != serverURL {
			t.Errorf("expected server URL %q, got %q", serverURL, gotURL)
		}
		if !strings.HasPrefix(gotURL, "http://localhost:") {
			t.Errorf("expected URL to start with http://localhost:, got %q", gotURL)
		}
	})

	// Clean up by triggering a callback
	go func() {
		time.Sleep(100 * time.Millisecond)
		http.Get(serverURL + "/callback?code=cleanup")
	}()
	WaitForCallback(ctx, server)
}

// TestCallbackServerNoAuthCode tests callback with no authorization code.
func TestCallbackServerNoAuthCode(t *testing.T) {
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

	// Simulate callback without code
	callbackURL := serverURL + "/callback"
	resp, err := http.Get(callbackURL)
	if err != nil {
		t.Fatalf("failed to make callback request: %v", err)
	}
	resp.Body.Close()

	select {
	case err := <-errChan:
		if err == nil {
			t.Error("expected error for missing auth code")
		}
		if err != ErrNoAuthCode {
			t.Errorf("expected ErrNoAuthCode, got %v", err)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for error")
	}
}

// TestCallbackServerTimeout tests callback server timeout behavior.
func TestCallbackServerTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	server, _, err := StartCallbackServer(ctx, 0)
	if err != nil {
		t.Fatalf("failed to start callback server: %v", err)
	}

	t.Run("returns error on context timeout", func(t *testing.T) {
		_, err := WaitForCallback(ctx, server)
		if err == nil {
			t.Error("expected timeout error")
		}
		if !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Errorf("expected context deadline error, got %v", err)
		}
	})
}

// TestStartCallbackServerWithSpecificPort tests starting server on a specific port.
func TestStartCallbackServerWithSpecificPort(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start server on a specific port
	server, serverURL, err := StartCallbackServer(ctx, 18765)
	if err != nil {
		t.Fatalf("failed to start callback server: %v", err)
	}

	t.Run("server uses requested port", func(t *testing.T) {
		if !strings.Contains(serverURL, ":18765") {
			t.Errorf("expected server URL to contain port 18765, got %q", serverURL)
		}
	})

	// Clean up
	go func() {
		time.Sleep(100 * time.Millisecond)
		http.Get(serverURL + "/callback?code=cleanup")
	}()
	WaitForCallback(ctx, server)
}

// TestGenerateCodeVerifierUniqueness tests that code verifiers are unique.
func TestGenerateCodeVerifierUniqueness(t *testing.T) {
	verifiers := make(map[string]bool)
	numIterations := 100

	for i := 0; i < numIterations; i++ {
		v := GenerateCodeVerifier()
		if verifiers[v] {
			t.Errorf("duplicate verifier generated on iteration %d", i)
		}
		verifiers[v] = true
	}

	if len(verifiers) != numIterations {
		t.Errorf("expected %d unique verifiers, got %d", numIterations, len(verifiers))
	}
}

// TestOpenBrowser tests the OpenBrowser function for various platforms.
// Note: This test doesn't actually open a browser; it just verifies the function exists
// and handles the current platform.
func TestOpenBrowser(t *testing.T) {
	// Skip this test in CI environments where browsers may not be available
	if os.Getenv("CI") == "true" {
		t.Skip("skipping browser test in CI environment")
	}

	// We can't fully test OpenBrowser without side effects, but we can
	// verify it doesn't panic and returns appropriate errors for invalid URLs
	t.Run("handles non-existent URL gracefully", func(t *testing.T) {
		// This will start a process but it will fail quickly
		// We're mainly testing that the function doesn't panic
		err := OpenBrowser("http://127.0.0.1:1") // Invalid URL that should fail
		// We don't check the error because the behavior is platform-dependent
		_ = err
	})
}
