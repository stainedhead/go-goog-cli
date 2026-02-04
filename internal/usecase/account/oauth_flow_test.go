// Package account provides application use cases for account management.
package account

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestDefaultOAuthFlow_Run_Success(t *testing.T) {
	expectedToken := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		OAuthProvider: &MockOAuthProvider{
			AuthURL: "https://accounts.google.com/auth?test=1",
			Token:   expectedToken,
		},
		BrowserOpener: &MockBrowserOpener{},
		CallbackServer: &MockCallbackServer{
			ServerURL: "http://localhost:8085",
			Code:      "test-auth-code",
		},
		UserInfoFetcher: &MockUserInfoFetcher{
			Email: "test@example.com",
		},
		PKCEGenerator: &MockPKCEGenerator{
			Verifier:  "test-verifier",
			Challenge: "test-challenge",
		},
	})

	email, token, err := flow.Run(context.Background(), []string{"openid"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", email)
	}
	if token.AccessToken != expectedToken.AccessToken {
		t.Errorf("expected access token '%s', got '%s'", expectedToken.AccessToken, token.AccessToken)
	}
}

func TestDefaultOAuthFlow_Run_ValidationError(t *testing.T) {
	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		OAuthProvider: &MockOAuthProvider{
			ValidateErr: errors.New("missing client ID"),
		},
		BrowserOpener:  &MockBrowserOpener{},
		CallbackServer: &MockCallbackServer{},
		PKCEGenerator:  &MockPKCEGenerator{},
	})

	_, _, err := flow.Run(context.Background(), []string{"openid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, err) || err.Error() != "invalid OAuth config: missing client ID" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDefaultOAuthFlow_Run_CallbackServerStartError(t *testing.T) {
	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		OAuthProvider: &MockOAuthProvider{},
		BrowserOpener: &MockBrowserOpener{},
		CallbackServer: &MockCallbackServer{
			StartErr: errors.New("port already in use"),
		},
		PKCEGenerator: &MockPKCEGenerator{},
	})

	_, _, err := flow.Run(context.Background(), []string{"openid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to start callback server: port already in use" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDefaultOAuthFlow_Run_BrowserOpenError(t *testing.T) {
	// Browser open error should not fail the flow (it's just a warning)
	expectedToken := &oauth2.Token{
		AccessToken: "test-access-token",
	}

	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		OAuthProvider: &MockOAuthProvider{
			AuthURL: "https://accounts.google.com/auth",
			Token:   expectedToken,
		},
		BrowserOpener: &MockBrowserOpener{
			Err: errors.New("no browser available"),
		},
		CallbackServer: &MockCallbackServer{
			ServerURL: "http://localhost:8085",
			Code:      "test-auth-code",
		},
		UserInfoFetcher: &MockUserInfoFetcher{
			Email: "test@example.com",
		},
		PKCEGenerator: &MockPKCEGenerator{
			Verifier:  "test-verifier",
			Challenge: "test-challenge",
		},
	})

	email, token, err := flow.Run(context.Background(), []string{"openid"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", email)
	}
	if token.AccessToken != expectedToken.AccessToken {
		t.Errorf("expected access token '%s', got '%s'", expectedToken.AccessToken, token.AccessToken)
	}
}

func TestDefaultOAuthFlow_Run_CallbackError(t *testing.T) {
	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		OAuthProvider: &MockOAuthProvider{
			AuthURL: "https://accounts.google.com/auth",
		},
		BrowserOpener: &MockBrowserOpener{},
		CallbackServer: &MockCallbackServer{
			ServerURL:   "http://localhost:8085",
			CallbackErr: errors.New("user denied access"),
		},
		PKCEGenerator: &MockPKCEGenerator{},
	})

	_, _, err := flow.Run(context.Background(), []string{"openid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "authentication failed: user denied access" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDefaultOAuthFlow_Run_ExchangeError(t *testing.T) {
	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		OAuthProvider: &MockOAuthProvider{
			AuthURL:     "https://accounts.google.com/auth",
			ExchangeErr: errors.New("invalid authorization code"),
		},
		BrowserOpener: &MockBrowserOpener{},
		CallbackServer: &MockCallbackServer{
			ServerURL: "http://localhost:8085",
			Code:      "invalid-code",
		},
		PKCEGenerator: &MockPKCEGenerator{},
	})

	_, _, err := flow.Run(context.Background(), []string{"openid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to exchange code: invalid authorization code" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDefaultOAuthFlow_Run_UserInfoError(t *testing.T) {
	expectedToken := &oauth2.Token{
		AccessToken: "test-access-token",
	}

	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		OAuthProvider: &MockOAuthProvider{
			AuthURL: "https://accounts.google.com/auth",
			Token:   expectedToken,
		},
		BrowserOpener: &MockBrowserOpener{},
		CallbackServer: &MockCallbackServer{
			ServerURL: "http://localhost:8085",
			Code:      "test-auth-code",
		},
		UserInfoFetcher: &MockUserInfoFetcher{
			Err: errors.New("API error"),
		},
		PKCEGenerator: &MockPKCEGenerator{},
	})

	_, _, err := flow.Run(context.Background(), []string{"openid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to get user email: API error" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewDefaultOAuthFlow(t *testing.T) {
	flow := NewDefaultOAuthFlow()
	if flow == nil {
		t.Fatal("expected non-nil flow")
	}
	if flow.browserOpener == nil {
		t.Error("expected browserOpener to be set")
	}
	if flow.callbackServer == nil {
		t.Error("expected callbackServer to be set")
	}
	if flow.userInfoFetcher == nil {
		t.Error("expected userInfoFetcher to be set")
	}
	if flow.pkceGenerator == nil {
		t.Error("expected pkceGenerator to be set")
	}
}

func TestNewDefaultOAuthFlowWithConfig_Defaults(t *testing.T) {
	// Test with all nil config - should set defaults
	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{})
	if flow == nil {
		t.Fatal("expected non-nil flow")
	}
	if flow.browserOpener == nil {
		t.Error("expected browserOpener to be set to default")
	}
	if flow.callbackServer == nil {
		t.Error("expected callbackServer to be set to default")
	}
	if flow.userInfoFetcher == nil {
		t.Error("expected userInfoFetcher to be set to default")
	}
	if flow.pkceGenerator == nil {
		t.Error("expected pkceGenerator to be set to default")
	}
}

func TestDefaultUserInfoFetcher_GetUserEmail_Success(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: NewMockHTTPResponse(http.StatusOK, `{"email": "test@example.com"}`),
	}

	fetcher := NewDefaultUserInfoFetcher(mockClient)
	token := &oauth2.Token{AccessToken: "test-token"}

	email, err := fetcher.GetUserEmail(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", email)
	}
}

func TestDefaultUserInfoFetcher_GetUserEmail_HTTPError(t *testing.T) {
	mockClient := &MockHTTPClient{
		Err: errors.New("network error"),
	}

	fetcher := NewDefaultUserInfoFetcher(mockClient)
	token := &oauth2.Token{AccessToken: "test-token"}

	_, err := fetcher.GetUserEmail(context.Background(), token)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "network error" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDefaultUserInfoFetcher_GetUserEmail_NonOKStatus(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: NewMockHTTPResponse(http.StatusUnauthorized, "unauthorized"),
	}

	fetcher := NewDefaultUserInfoFetcher(mockClient)
	token := &oauth2.Token{AccessToken: "test-token"}

	_, err := fetcher.GetUserEmail(context.Background(), token)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultUserInfoFetcher_GetUserEmail_InvalidJSON(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: NewMockHTTPResponse(http.StatusOK, "not json"),
	}

	fetcher := NewDefaultUserInfoFetcher(mockClient)
	token := &oauth2.Token{AccessToken: "test-token"}

	_, err := fetcher.GetUserEmail(context.Background(), token)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultUserInfoFetcher_GetUserEmail_EmptyEmail(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: NewMockHTTPResponse(http.StatusOK, `{"email": ""}`),
	}

	fetcher := NewDefaultUserInfoFetcher(mockClient)
	token := &oauth2.Token{AccessToken: "test-token"}

	_, err := fetcher.GetUserEmail(context.Background(), token)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "no email in userinfo response" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDefaultUserInfoFetcher_NilClient(t *testing.T) {
	// Should use default http.Client when nil is passed
	fetcher := NewDefaultUserInfoFetcher(nil)
	if fetcher.client == nil {
		t.Error("expected client to be set to default")
	}
}

func TestDefaultCallbackServer_GetServerURL_NilServer(t *testing.T) {
	server := &DefaultCallbackServer{}
	url := server.GetServerURL()
	if url != "" {
		t.Errorf("expected empty string, got '%s'", url)
	}
}

func TestDefaultCallbackServer_Stop(t *testing.T) {
	server := &DefaultCallbackServer{}
	err := server.Stop()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMockOAuthProvider(t *testing.T) {
	token := &oauth2.Token{AccessToken: "test"}
	mock := &MockOAuthProvider{
		AuthURL:     "https://auth.example.com",
		Token:       token,
		RedirectURL: "http://localhost/callback",
	}

	// Test GetAuthURL
	url := mock.GetAuthURL("state", "challenge")
	if url != "https://auth.example.com" {
		t.Errorf("expected auth URL, got '%s'", url)
	}

	// Test Exchange
	result, err := mock.Exchange(context.Background(), "code", "verifier")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AccessToken != token.AccessToken {
		t.Errorf("expected token, got '%v'", result)
	}

	// Test GetRedirectURL
	if mock.GetRedirectURL() != "http://localhost/callback" {
		t.Errorf("expected redirect URL")
	}

	// Test SetRedirectURL
	mock.SetRedirectURL("http://new/callback")
	if mock.GetRedirectURL() != "http://new/callback" {
		t.Errorf("expected updated redirect URL")
	}

	// Test Validate
	if mock.Validate() != nil {
		t.Error("expected nil validate error")
	}

	// Test with error
	mock.ValidateErr = errors.New("invalid")
	if mock.Validate() == nil {
		t.Error("expected validate error")
	}

	// Test Exchange with error
	mock.ExchangeErr = errors.New("exchange failed")
	_, err = mock.Exchange(context.Background(), "code", "verifier")
	if err == nil {
		t.Error("expected exchange error")
	}
}

func TestMockBrowserOpener(t *testing.T) {
	mock := &MockBrowserOpener{}

	err := mock.Open("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.OpenedURL != "https://example.com" {
		t.Errorf("expected URL to be recorded")
	}

	// Test with error
	mock.Err = errors.New("no browser")
	err = mock.Open("https://example.com")
	if err == nil {
		t.Error("expected error")
	}
}

func TestMockCallbackServer(t *testing.T) {
	mock := &MockCallbackServer{
		ServerURL: "http://localhost:8085",
		Code:      "test-code",
	}

	// Test Start
	url, err := mock.Start(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != "http://localhost:8085" {
		t.Errorf("expected server URL")
	}

	// Test WaitForCallback
	code, err := mock.WaitForCallback(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "test-code" {
		t.Errorf("expected code")
	}

	// Test GetServerURL
	if mock.GetServerURL() != "http://localhost:8085" {
		t.Errorf("expected server URL")
	}

	// Test Stop
	if mock.Stop() != nil {
		t.Error("expected nil stop error")
	}

	// Test with errors
	mock.StartErr = errors.New("start failed")
	_, err = mock.Start(context.Background())
	if err == nil {
		t.Error("expected start error")
	}

	mock.CallbackErr = errors.New("callback failed")
	_, err = mock.WaitForCallback(context.Background())
	if err == nil {
		t.Error("expected callback error")
	}
}

func TestMockUserInfoFetcher(t *testing.T) {
	mock := &MockUserInfoFetcher{
		Email: "test@example.com",
	}

	email, err := mock.GetUserEmail(context.Background(), &oauth2.Token{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email != "test@example.com" {
		t.Errorf("expected email")
	}

	// Test with error
	mock.Err = errors.New("fetch failed")
	_, err = mock.GetUserEmail(context.Background(), &oauth2.Token{})
	if err == nil {
		t.Error("expected error")
	}
}

func TestMockPKCEGenerator(t *testing.T) {
	mock := &MockPKCEGenerator{
		Verifier:  "test-verifier",
		Challenge: "test-challenge",
	}

	if mock.GenerateVerifier() != "test-verifier" {
		t.Error("expected verifier")
	}
	if mock.GenerateChallenge("test") != "test-challenge" {
		t.Error("expected challenge")
	}
}

func TestMockHTTPClient(t *testing.T) {
	response := NewMockHTTPResponse(http.StatusOK, "test body")
	mock := &MockHTTPClient{
		Response: response,
	}

	req, _ := http.NewRequest("GET", "https://example.com", nil)
	resp, err := mock.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK")
	}

	// Test with error
	mock.Err = errors.New("network error")
	_, err = mock.Do(req)
	if err == nil {
		t.Error("expected error")
	}
}

func TestMockTokenSource(t *testing.T) {
	token := &oauth2.Token{AccessToken: "test"}
	mock := &MockTokenSource{
		MockToken: token,
	}

	result, err := mock.Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AccessToken != token.AccessToken {
		t.Errorf("expected token")
	}

	// Test with error
	mock.Err = errors.New("token error")
	_, err = mock.Token()
	if err == nil {
		t.Error("expected error")
	}
}

// TestDefaultBrowserOpener_Open tests the DefaultBrowserOpener (cannot actually open browser in test).
func TestDefaultBrowserOpener_Open(t *testing.T) {
	// We can't test actually opening a browser in automated tests
	// but we can verify the type implements the interface
	opener := &DefaultBrowserOpener{}
	var _ BrowserOpener = opener
}

// TestDefaultPKCEGenerator_GenerateVerifier tests the PKCE verifier generation.
func TestDefaultPKCEGenerator_GenerateVerifier(t *testing.T) {
	generator := &DefaultPKCEGenerator{}
	verifier := generator.GenerateVerifier()

	// Verifier should be a non-empty string
	if verifier == "" {
		t.Error("expected non-empty verifier")
	}

	// Verifier should be 43-128 characters (standard PKCE requirement)
	if len(verifier) < 43 || len(verifier) > 128 {
		t.Errorf("verifier length %d outside valid range [43, 128]", len(verifier))
	}
}

// TestDefaultPKCEGenerator_GenerateChallenge tests the PKCE challenge generation.
func TestDefaultPKCEGenerator_GenerateChallenge(t *testing.T) {
	generator := &DefaultPKCEGenerator{}
	verifier := generator.GenerateVerifier()
	challenge := generator.GenerateChallenge(verifier)

	// Challenge should be a non-empty string
	if challenge == "" {
		t.Error("expected non-empty challenge")
	}

	// Challenge should be URL-safe base64 encoded SHA256 (43 chars)
	if len(challenge) != 43 {
		t.Errorf("challenge length %d, expected 43", len(challenge))
	}
}

// TestDefaultPKCEGenerator_ConsistentChallenge tests that the same verifier produces the same challenge.
func TestDefaultPKCEGenerator_ConsistentChallenge(t *testing.T) {
	generator := &DefaultPKCEGenerator{}
	verifier := "test-verifier-12345678901234567890123456789012"

	challenge1 := generator.GenerateChallenge(verifier)
	challenge2 := generator.GenerateChallenge(verifier)

	if challenge1 != challenge2 {
		t.Errorf("same verifier should produce same challenge: %q != %q", challenge1, challenge2)
	}
}

// TestDefaultOAuthProvider_NewDefaultOAuthProvider tests creating a new OAuth provider.
func TestDefaultOAuthProvider_NewDefaultOAuthProvider(t *testing.T) {
	scopes := []string{"openid", "email"}
	provider := NewDefaultOAuthProvider(scopes)

	if provider == nil {
		t.Fatal("expected non-nil provider")
	}
	if provider.config == nil {
		t.Error("expected non-nil config")
	}
}

// TestDefaultOAuthProvider_GetRedirectURL tests getting the redirect URL.
func TestDefaultOAuthProvider_GetRedirectURL(t *testing.T) {
	provider := NewDefaultOAuthProvider([]string{"openid"})
	url := provider.GetRedirectURL()

	// Default config should have a redirect URL
	if url == "" {
		t.Error("expected non-empty redirect URL")
	}
}

// TestDefaultOAuthProvider_SetRedirectURL tests setting the redirect URL.
func TestDefaultOAuthProvider_SetRedirectURL(t *testing.T) {
	provider := NewDefaultOAuthProvider([]string{"openid"})
	newURL := "http://localhost:12345/callback"

	provider.SetRedirectURL(newURL)
	result := provider.GetRedirectURL()

	if result != newURL {
		t.Errorf("expected redirect URL %q, got %q", newURL, result)
	}
}

// TestDefaultOAuthProvider_Validate tests validation (will fail without credentials).
func TestDefaultOAuthProvider_Validate(t *testing.T) {
	provider := NewDefaultOAuthProvider([]string{"openid"})

	// Validate may return an error if client ID/secret are not configured
	// We just verify it doesn't panic
	_ = provider.Validate()
}

// TestDefaultOAuthProvider_GetAuthURL tests generating an auth URL.
func TestDefaultOAuthProvider_GetAuthURL(t *testing.T) {
	provider := NewDefaultOAuthProvider([]string{"openid"})
	url := provider.GetAuthURL("test-state", "test-challenge")

	if url == "" {
		t.Error("expected non-empty auth URL")
	}

	// URL should contain state parameter
	if !containsStr(url, "state=test-state") {
		t.Errorf("auth URL should contain state parameter: %s", url)
	}

	// URL should contain code challenge
	if !containsStr(url, "code_challenge=test-challenge") {
		t.Errorf("auth URL should contain code challenge: %s", url)
	}
}

// TestDefaultOAuthProvider_TokenSource tests getting a token source.
func TestDefaultOAuthProvider_TokenSource(t *testing.T) {
	provider := NewDefaultOAuthProvider([]string{"openid"})
	token := &oauth2.Token{
		AccessToken:  "test-token",
		RefreshToken: "test-refresh",
		Expiry:       time.Now().Add(time.Hour),
	}

	ts := provider.TokenSource(context.Background(), token)
	if ts == nil {
		t.Error("expected non-nil token source")
	}
}

// TestDefaultCallbackServer_NilServerGetServerURL tests GetServerURL when server is nil.
func TestDefaultCallbackServer_NilServerGetServerURL(t *testing.T) {
	server := &DefaultCallbackServer{}
	url := server.GetServerURL()

	if url != "" {
		t.Errorf("expected empty URL for nil server, got %q", url)
	}
}

// TestDefaultCallbackServer_Stop_Idempotent tests that Stop can be called multiple times.
func TestDefaultCallbackServer_Stop_Idempotent(t *testing.T) {
	server := &DefaultCallbackServer{}

	// Should not error even when server wasn't started
	err := server.Stop()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Should not error on second call
	err = server.Stop()
	if err != nil {
		t.Errorf("expected nil error on second stop, got %v", err)
	}
}

// TestDefaultUserInfoFetcher_GetUserEmail_WithNilClient tests that nil client gets default.
func TestDefaultUserInfoFetcher_GetUserEmail_WithNilClient(t *testing.T) {
	fetcher := NewDefaultUserInfoFetcher(nil)

	// Verify client is set to default
	if fetcher.client == nil {
		t.Error("expected default client to be set")
	}
}

// TestDefaultOAuthFlow_Run_WithNilProvider tests that nil provider gets default.
func TestDefaultOAuthFlow_Run_WithNilProvider(t *testing.T) {
	// Create flow without provider - it should create one internally
	flow := NewDefaultOAuthFlowWithConfig(OAuthFlowConfig{
		BrowserOpener: &MockBrowserOpener{},
		CallbackServer: &MockCallbackServer{
			StartErr: errors.New("test error"), // Stop early
		},
		PKCEGenerator: &MockPKCEGenerator{
			Verifier:  "test-verifier",
			Challenge: "test-challenge",
		},
	})

	// The flow will fail when trying to start the callback server
	// but it should have created a default provider
	_, _, err := flow.Run(context.Background(), []string{"openid"})
	if err == nil {
		t.Error("expected error from callback server start")
	}
}

// Helper function
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
