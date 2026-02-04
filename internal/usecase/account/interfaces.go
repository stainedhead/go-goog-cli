// Package account provides application use cases for account management.
package account

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// OAuthProvider handles OAuth2 authentication configuration and token exchange.
type OAuthProvider interface {
	// GetAuthURL returns the OAuth2 authorization URL with PKCE parameters.
	GetAuthURL(state, codeChallenge string) string
	// Exchange exchanges an authorization code for an OAuth2 token.
	Exchange(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error)
	// TokenSource returns a token source that auto-refreshes the token.
	TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource
	// GetRedirectURL returns the configured redirect URL.
	GetRedirectURL() string
	// SetRedirectURL sets the redirect URL for the OAuth flow.
	SetRedirectURL(url string)
	// Validate checks that the OAuth configuration has required credentials.
	Validate() error
}

// BrowserOpener opens URLs in the default browser.
type BrowserOpener interface {
	// Open opens the specified URL in the default browser.
	Open(url string) error
}

// CallbackServer handles OAuth callbacks on localhost.
type CallbackServer interface {
	// Start starts the callback server and returns the server URL.
	Start(ctx context.Context) (serverURL string, err error)
	// WaitForCallback waits for the OAuth callback and returns the authorization code.
	WaitForCallback(ctx context.Context) (code string, err error)
	// Stop stops the callback server.
	Stop() error
	// GetServerURL returns the server URL.
	GetServerURL() string
}

// TokenStore stores and retrieves OAuth tokens securely.
type TokenStore interface {
	// Save stores an OAuth2 token for the given alias.
	Save(alias string, token *oauth2.Token) error
	// Load retrieves an OAuth2 token for the given alias.
	Load(alias string) (*oauth2.Token, error)
	// Delete removes the OAuth2 token for the given alias.
	Delete(alias string) error
}

// UserInfoFetcher retrieves user information from an OAuth provider.
type UserInfoFetcher interface {
	// GetUserEmail retrieves the user's email from the provider.
	GetUserEmail(ctx context.Context, token *oauth2.Token) (string, error)
}

// PKCEGenerator generates PKCE parameters for OAuth2 flows.
type PKCEGenerator interface {
	// GenerateVerifier generates a code verifier.
	GenerateVerifier() string
	// GenerateChallenge generates a code challenge from a verifier.
	GenerateChallenge(verifier string) string
}

// HTTPClient defines the interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
