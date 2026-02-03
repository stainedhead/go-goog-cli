// Package account provides application use cases for account management.
package account

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/infrastructure/auth"
	"golang.org/x/oauth2"
)

// DefaultOAuthFlow implements the OAuth2/PKCE flow for CLI authentication.
type DefaultOAuthFlow struct {
	openBrowser func(url string) error
}

// NewDefaultOAuthFlow creates a new DefaultOAuthFlow.
func NewDefaultOAuthFlow() *DefaultOAuthFlow {
	return &DefaultOAuthFlow{
		openBrowser: auth.OpenBrowser,
	}
}

// Run executes the OAuth flow and returns the email and token.
func (f *DefaultOAuthFlow) Run(ctx context.Context, scopes []string) (string, *oauth2.Token, error) {
	// Create OAuth config
	cfg := auth.NewOAuthConfig(scopes)

	// Validate config
	if err := auth.ValidateConfig(cfg); err != nil {
		return "", nil, fmt.Errorf("invalid OAuth config: %w", err)
	}

	// Generate PKCE parameters
	verifier := auth.GenerateCodeVerifier()
	challenge := auth.GenerateCodeChallenge(verifier)

	// Start callback server
	callbackServer, serverURL, err := auth.StartCallbackServer(ctx, 0)
	if err != nil {
		return "", nil, fmt.Errorf("failed to start callback server: %w", err)
	}

	// Update redirect URL to use the actual server port
	cfg.RedirectURL = serverURL + auth.DefaultRedirectPath

	// Generate state for CSRF protection
	state := auth.GenerateCodeVerifier() // Reuse verifier generation for state

	// Get authorization URL
	authURL := auth.GetAuthorizationURL(cfg, state, challenge)

	// Open browser
	fmt.Println("Opening browser for authentication...")
	fmt.Printf("If the browser doesn't open, visit this URL:\n%s\n", authURL)

	if err := f.openBrowser(authURL); err != nil {
		fmt.Printf("Warning: could not open browser: %v\n", err)
	}

	// Wait for callback with timeout
	callbackCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	code, err := auth.WaitForCallback(callbackCtx, callbackServer)
	if err != nil {
		return "", nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Exchange code for token
	token, err := auth.ExchangeCode(ctx, cfg, code, verifier)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user email from token
	email, err := f.getUserEmail(ctx, token)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user email: %w", err)
	}

	return email, token, nil
}

// getUserEmail retrieves the user's email from the Google userinfo endpoint.
func (f *DefaultOAuthFlow) getUserEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("userinfo request failed: %s - %s", resp.Status, string(body))
	}

	var userInfo struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	if userInfo.Email == "" {
		return "", fmt.Errorf("no email in userinfo response")
	}

	return userInfo.Email, nil
}
