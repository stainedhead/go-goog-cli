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

// DefaultBrowserOpener implements BrowserOpener using the auth package.
type DefaultBrowserOpener struct{}

// Open opens the URL in the default browser.
func (d *DefaultBrowserOpener) Open(url string) error {
	return auth.OpenBrowser(url)
}

// DefaultPKCEGenerator implements PKCEGenerator using the auth package.
type DefaultPKCEGenerator struct{}

// GenerateVerifier generates a code verifier.
func (d *DefaultPKCEGenerator) GenerateVerifier() string {
	return auth.GenerateCodeVerifier()
}

// GenerateChallenge generates a code challenge from a verifier.
func (d *DefaultPKCEGenerator) GenerateChallenge(verifier string) string {
	return auth.GenerateCodeChallenge(verifier)
}

// DefaultOAuthProvider implements OAuthProvider using the auth package.
type DefaultOAuthProvider struct {
	config *oauth2.Config
}

// NewDefaultOAuthProvider creates a new DefaultOAuthProvider with the given scopes.
func NewDefaultOAuthProvider(scopes []string) *DefaultOAuthProvider {
	return &DefaultOAuthProvider{
		config: auth.NewOAuthConfig(scopes),
	}
}

// GetAuthURL returns the OAuth2 authorization URL with PKCE parameters.
func (d *DefaultOAuthProvider) GetAuthURL(state, codeChallenge string) string {
	return auth.GetAuthorizationURL(d.config, state, codeChallenge)
}

// Exchange exchanges an authorization code for an OAuth2 token.
func (d *DefaultOAuthProvider) Exchange(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	return auth.ExchangeCode(ctx, d.config, code, codeVerifier)
}

// TokenSource returns a token source that auto-refreshes the token.
func (d *DefaultOAuthProvider) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return d.config.TokenSource(ctx, token)
}

// GetRedirectURL returns the configured redirect URL.
func (d *DefaultOAuthProvider) GetRedirectURL() string {
	return d.config.RedirectURL
}

// SetRedirectURL sets the redirect URL for the OAuth flow.
func (d *DefaultOAuthProvider) SetRedirectURL(url string) {
	d.config.RedirectURL = url
}

// Validate checks that the OAuth configuration has required credentials.
func (d *DefaultOAuthProvider) Validate() error {
	return auth.ValidateConfig(d.config)
}

// DefaultCallbackServer implements CallbackServer using the auth package.
type DefaultCallbackServer struct {
	server *auth.CallbackServer
}

// Start starts the callback server and returns the server URL.
func (d *DefaultCallbackServer) Start(ctx context.Context) (string, error) {
	server, serverURL, err := auth.StartCallbackServer(ctx, 0)
	if err != nil {
		return "", err
	}
	d.server = server
	return serverURL, nil
}

// WaitForCallback waits for the OAuth callback and returns the authorization code.
func (d *DefaultCallbackServer) WaitForCallback(ctx context.Context) (string, error) {
	return auth.WaitForCallback(ctx, d.server)
}

// Stop stops the callback server.
func (d *DefaultCallbackServer) Stop() error {
	// The callback server is stopped automatically in WaitForCallback
	return nil
}

// GetServerURL returns the server URL.
func (d *DefaultCallbackServer) GetServerURL() string {
	if d.server == nil {
		return ""
	}
	return d.server.GetServerURL()
}

// DefaultUserInfoFetcher implements UserInfoFetcher using Google's userinfo API.
type DefaultUserInfoFetcher struct {
	client HTTPClient
}

// NewDefaultUserInfoFetcher creates a new DefaultUserInfoFetcher with the given HTTP client.
func NewDefaultUserInfoFetcher(client HTTPClient) *DefaultUserInfoFetcher {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &DefaultUserInfoFetcher{client: client}
}

// GetUserEmail retrieves the user's email from Google's userinfo endpoint.
func (d *DefaultUserInfoFetcher) GetUserEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := d.client.Do(req)
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

// OAuthFlowConfig contains the dependencies for the OAuth flow.
type OAuthFlowConfig struct {
	OAuthProvider   OAuthProvider
	BrowserOpener   BrowserOpener
	CallbackServer  CallbackServer
	UserInfoFetcher UserInfoFetcher
	PKCEGenerator   PKCEGenerator
}

// DefaultOAuthFlow implements the OAuth2/PKCE flow for CLI authentication.
type DefaultOAuthFlow struct {
	provider        OAuthProvider
	browserOpener   BrowserOpener
	callbackServer  CallbackServer
	userInfoFetcher UserInfoFetcher
	pkceGenerator   PKCEGenerator
}

// NewDefaultOAuthFlow creates a new DefaultOAuthFlow with default implementations.
func NewDefaultOAuthFlow() *DefaultOAuthFlow {
	return &DefaultOAuthFlow{
		browserOpener:   &DefaultBrowserOpener{},
		callbackServer:  &DefaultCallbackServer{},
		userInfoFetcher: NewDefaultUserInfoFetcher(nil),
		pkceGenerator:   &DefaultPKCEGenerator{},
	}
}

// NewDefaultOAuthFlowWithConfig creates a new DefaultOAuthFlow with the provided configuration.
func NewDefaultOAuthFlowWithConfig(cfg OAuthFlowConfig) *DefaultOAuthFlow {
	flow := &DefaultOAuthFlow{
		provider:        cfg.OAuthProvider,
		browserOpener:   cfg.BrowserOpener,
		callbackServer:  cfg.CallbackServer,
		userInfoFetcher: cfg.UserInfoFetcher,
		pkceGenerator:   cfg.PKCEGenerator,
	}

	// Set defaults for nil dependencies
	if flow.browserOpener == nil {
		flow.browserOpener = &DefaultBrowserOpener{}
	}
	if flow.callbackServer == nil {
		flow.callbackServer = &DefaultCallbackServer{}
	}
	if flow.userInfoFetcher == nil {
		flow.userInfoFetcher = NewDefaultUserInfoFetcher(nil)
	}
	if flow.pkceGenerator == nil {
		flow.pkceGenerator = &DefaultPKCEGenerator{}
	}

	return flow
}

// Run executes the OAuth flow and returns the email and token.
func (f *DefaultOAuthFlow) Run(ctx context.Context, scopes []string) (string, *oauth2.Token, error) {
	// Create OAuth provider if not injected
	provider := f.provider
	if provider == nil {
		provider = NewDefaultOAuthProvider(scopes)
	}

	// Validate config
	if err := provider.Validate(); err != nil {
		return "", nil, fmt.Errorf("invalid OAuth config: %w", err)
	}

	// Generate PKCE parameters
	verifier := f.pkceGenerator.GenerateVerifier()
	challenge := f.pkceGenerator.GenerateChallenge(verifier)

	// Start callback server
	serverURL, err := f.callbackServer.Start(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to start callback server: %w", err)
	}

	// Update redirect URL to use the actual server port
	provider.SetRedirectURL(serverURL + auth.DefaultRedirectPath)

	// Generate state for CSRF protection
	state := f.pkceGenerator.GenerateVerifier()

	// Get authorization URL
	authURL := provider.GetAuthURL(state, challenge)

	// Open browser
	fmt.Println("Opening browser for authentication...")
	fmt.Printf("If the browser doesn't open, visit this URL:\n%s\n", authURL)

	if err := f.browserOpener.Open(authURL); err != nil {
		fmt.Printf("Warning: could not open browser: %v\n", err)
	}

	// Wait for callback with timeout
	callbackCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	code, err := f.callbackServer.WaitForCallback(callbackCtx)
	if err != nil {
		return "", nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Exchange code for token
	token, err := provider.Exchange(ctx, code, verifier)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user email from token
	email, err := f.userInfoFetcher.GetUserEmail(ctx, token)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user email: %w", err)
	}

	return email, token, nil
}
