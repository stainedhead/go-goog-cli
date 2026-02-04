// Package account provides application use cases for account management.
package account

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// MockOAuthProvider is a mock implementation of OAuthProvider for testing.
type MockOAuthProvider struct {
	AuthURL     string
	Token       *oauth2.Token
	ExchangeErr error
	ValidateErr error
	RedirectURL string
	TokenSrc    oauth2.TokenSource
}

// GetAuthURL returns the mock auth URL.
func (m *MockOAuthProvider) GetAuthURL(state, codeChallenge string) string {
	return m.AuthURL
}

// Exchange returns the mock token or error.
func (m *MockOAuthProvider) Exchange(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	if m.ExchangeErr != nil {
		return nil, m.ExchangeErr
	}
	return m.Token, nil
}

// TokenSource returns the mock token source.
func (m *MockOAuthProvider) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return m.TokenSrc
}

// GetRedirectURL returns the mock redirect URL.
func (m *MockOAuthProvider) GetRedirectURL() string {
	return m.RedirectURL
}

// SetRedirectURL sets the mock redirect URL.
func (m *MockOAuthProvider) SetRedirectURL(url string) {
	m.RedirectURL = url
}

// Validate returns the mock validation error.
func (m *MockOAuthProvider) Validate() error {
	return m.ValidateErr
}

// MockBrowserOpener is a mock implementation of BrowserOpener for testing.
type MockBrowserOpener struct {
	OpenedURL string
	Err       error
}

// Open records the URL and returns the mock error.
func (m *MockBrowserOpener) Open(url string) error {
	m.OpenedURL = url
	return m.Err
}

// MockCallbackServer is a mock implementation of CallbackServer for testing.
type MockCallbackServer struct {
	ServerURL   string
	StartErr    error
	Code        string
	CallbackErr error
	StopErr     error
}

// Start returns the mock server URL and error.
func (m *MockCallbackServer) Start(ctx context.Context) (string, error) {
	if m.StartErr != nil {
		return "", m.StartErr
	}
	return m.ServerURL, nil
}

// WaitForCallback returns the mock code and error.
func (m *MockCallbackServer) WaitForCallback(ctx context.Context) (string, error) {
	if m.CallbackErr != nil {
		return "", m.CallbackErr
	}
	return m.Code, nil
}

// Stop returns the mock stop error.
func (m *MockCallbackServer) Stop() error {
	return m.StopErr
}

// GetServerURL returns the mock server URL.
func (m *MockCallbackServer) GetServerURL() string {
	return m.ServerURL
}

// MockUserInfoFetcher is a mock implementation of UserInfoFetcher for testing.
type MockUserInfoFetcher struct {
	Email string
	Err   error
}

// GetUserEmail returns the mock email and error.
func (m *MockUserInfoFetcher) GetUserEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.Email, nil
}

// MockPKCEGenerator is a mock implementation of PKCEGenerator for testing.
type MockPKCEGenerator struct {
	Verifier  string
	Challenge string
}

// GenerateVerifier returns the mock verifier.
func (m *MockPKCEGenerator) GenerateVerifier() string {
	return m.Verifier
}

// GenerateChallenge returns the mock challenge.
func (m *MockPKCEGenerator) GenerateChallenge(verifier string) string {
	return m.Challenge
}

// MockHTTPClient is a mock implementation of HTTPClient for testing.
type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

// Do returns the mock response and error.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Response, nil
}

// NewMockHTTPResponse creates a mock HTTP response with the given status code and body.
func NewMockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
		Status:     http.StatusText(statusCode),
	}
}

// MockTokenSource is a mock implementation of oauth2.TokenSource for testing.
type MockTokenSource struct {
	MockToken *oauth2.Token
	Err       error
}

// Token returns the mock token and error.
func (m *MockTokenSource) Token() (*oauth2.Token, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.MockToken, nil
}
