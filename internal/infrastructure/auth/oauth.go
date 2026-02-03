// Package auth provides OAuth2/PKCE authentication for Google APIs.
// It handles the complete OAuth2 authorization code flow with PKCE extension
// for secure authentication in CLI applications.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google API Scopes for commonly used services.
const (
	// Gmail scopes
	ScopeGmailReadonly = "https://www.googleapis.com/auth/gmail.readonly"
	ScopeGmailSend     = "https://www.googleapis.com/auth/gmail.send"
	ScopeGmailModify   = "https://www.googleapis.com/auth/gmail.modify"
	ScopeGmailCompose  = "https://www.googleapis.com/auth/gmail.compose"
	ScopeGmailLabels   = "https://www.googleapis.com/auth/gmail.labels"

	// Calendar scopes
	ScopeCalendarReadonly = "https://www.googleapis.com/auth/calendar.readonly"
	ScopeCalendarEvents   = "https://www.googleapis.com/auth/calendar.events"
	ScopeCalendar         = "https://www.googleapis.com/auth/calendar"

	// Drive scopes
	ScopeDriveReadonly = "https://www.googleapis.com/auth/drive.readonly"
	ScopeDriveFile     = "https://www.googleapis.com/auth/drive.file"
	ScopeDrive         = "https://www.googleapis.com/auth/drive"

	// User info scopes
	ScopeUserInfoEmail   = "https://www.googleapis.com/auth/userinfo.email"
	ScopeUserInfoProfile = "https://www.googleapis.com/auth/userinfo.profile"

	// OpenID scopes
	ScopeOpenID = "openid"
)

// Environment variable names for OAuth configuration.
const (
	EnvClientID     = "GOOG_CLIENT_ID"
	EnvClientSecret = "GOOG_CLIENT_SECRET"
	EnvRedirectPort = "GOOG_REDIRECT_PORT"
)

// Default configuration values.
const (
	DefaultRedirectPort = 8085
	DefaultRedirectPath = "/callback"
)

// Errors returned by the auth package.
var (
	ErrMissingClientID     = errors.New("GOOG_CLIENT_ID environment variable is not set")
	ErrMissingClientSecret = errors.New("GOOG_CLIENT_SECRET environment variable is not set")
	ErrOAuthError          = errors.New("OAuth error")
	ErrNoAuthCode          = errors.New("no authorization code received")
	ErrCallbackTimeout     = errors.New("callback timeout")
)

// CallbackServer handles the OAuth callback on localhost.
type CallbackServer struct {
	server     *http.Server
	listener   net.Listener
	codeChan   chan string
	errChan    chan error
	once       sync.Once
	serverURL  string
	shutdownWG sync.WaitGroup
}

// GenerateCodeVerifier generates a cryptographically random code verifier for PKCE.
// The verifier is 32 bytes (256 bits) of random data, base64url encoded without padding,
// resulting in a 43-character string.
func GenerateCodeVerifier() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate random bytes: %v", err))
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// GenerateCodeChallenge generates a code challenge from the code verifier using SHA256.
// The challenge is the SHA256 hash of the verifier, base64url encoded without padding.
func GenerateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// NewOAuthConfig creates a new OAuth2 configuration for Google APIs.
// It reads client credentials from environment variables:
//   - GOOG_CLIENT_ID: OAuth2 client ID
//   - GOOG_CLIENT_SECRET: OAuth2 client secret
//   - GOOG_REDIRECT_PORT: Localhost port for callback (default: 8085)
func NewOAuthConfig(scopes []string) *oauth2.Config {
	clientID := os.Getenv(EnvClientID)
	clientSecret := os.Getenv(EnvClientSecret)

	port := DefaultRedirectPort
	if portStr := os.Getenv(EnvRedirectPort); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("http://localhost:%d%s", port, DefaultRedirectPath),
	}
}

// NewOAuthConfigWithCredentials creates a new OAuth2 configuration with explicit credentials.
// This is useful when credentials are loaded from a config file rather than environment variables.
func NewOAuthConfigWithCredentials(clientID, clientSecret string, scopes []string, port int) *oauth2.Config {
	if port == 0 {
		port = DefaultRedirectPort
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("http://localhost:%d%s", port, DefaultRedirectPath),
	}
}

// GetAuthorizationURL generates the OAuth2 authorization URL with PKCE parameters.
// It includes the state parameter for CSRF protection and code_challenge for PKCE.
func GetAuthorizationURL(cfg *oauth2.Config, state, codeChallenge string) string {
	return cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
}

// StartCallbackServer starts a local HTTP server to handle the OAuth callback.
// If port is 0, a random available port will be used.
// Returns the server instance, the server URL, and any error.
func StartCallbackServer(ctx context.Context, port int) (*CallbackServer, string, error) {
	if port == 0 {
		port = DefaultRedirectPort
	}

	// Try to listen on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		// If the default port is in use, try a random port
		listener, err = net.Listen("tcp", "localhost:0")
		if err != nil {
			return nil, "", fmt.Errorf("failed to start callback server: %w", err)
		}
	}

	addr := listener.Addr().(*net.TCPAddr)
	serverURL := fmt.Sprintf("http://localhost:%d", addr.Port)

	cs := &CallbackServer{
		listener:  listener,
		codeChan:  make(chan string, 1),
		errChan:   make(chan error, 1),
		serverURL: serverURL,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(DefaultRedirectPath, cs.handleCallback)

	cs.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	cs.shutdownWG.Add(1)
	go func() {
		defer cs.shutdownWG.Done()
		if err := cs.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			cs.errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	return cs, serverURL, nil
}

// handleCallback processes the OAuth callback request.
func (cs *CallbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	cs.once.Do(func() {
		// Check for error response
		if errCode := r.URL.Query().Get("error"); errCode != "" {
			errDesc := r.URL.Query().Get("error_description")
			cs.errChan <- fmt.Errorf("%w: %s - %s", ErrOAuthError, errCode, errDesc)

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Authentication Failed</title></head>
<body>
<h1>Authentication Failed</h1>
<p>Error: %s</p>
<p>%s</p>
<p>You can close this window.</p>
</body>
</html>`, errCode, errDesc)
			return
		}

		// Extract authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			cs.errChan <- ErrNoAuthCode
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Authentication Failed</title></head>
<body>
<h1>Authentication Failed</h1>
<p>No authorization code received.</p>
<p>You can close this window.</p>
</body>
</html>`)
			return
		}

		cs.codeChan <- code

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Authentication Successful</title></head>
<body>
<h1>Authentication Successful!</h1>
<p>You have successfully authenticated with Google.</p>
<p>You can close this window and return to the terminal.</p>
</body>
</html>`)
	})
}

// WaitForCallback waits for the OAuth callback and returns the authorization code.
// It blocks until a callback is received, an error occurs, or the context is cancelled.
func WaitForCallback(ctx context.Context, cs *CallbackServer) (string, error) {
	defer func() {
		// Gracefully shutdown the server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cs.server.Shutdown(shutdownCtx)
		cs.shutdownWG.Wait()
	}()

	select {
	case code := <-cs.codeChan:
		return code, nil
	case err := <-cs.errChan:
		return "", err
	case <-ctx.Done():
		return "", fmt.Errorf("%w: %v", ErrCallbackTimeout, ctx.Err())
	}
}

// GetServerURL returns the server URL for the callback server.
func (cs *CallbackServer) GetServerURL() string {
	return cs.serverURL
}

// ExchangeCode exchanges the authorization code for an OAuth2 token.
// It includes the PKCE code_verifier for verification.
func ExchangeCode(ctx context.Context, cfg *oauth2.Config, code, codeVerifier string) (*oauth2.Token, error) {
	return cfg.Exchange(
		ctx,
		code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
}

// OpenBrowser opens the specified URL in the default browser.
// It uses platform-specific commands:
//   - macOS: open
//   - Linux: xdg-open
//   - Windows: cmd /c start
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// ValidateConfig checks that the OAuth2 configuration has required credentials.
func ValidateConfig(cfg *oauth2.Config) error {
	if cfg.ClientID == "" {
		return ErrMissingClientID
	}
	if cfg.ClientSecret == "" {
		return ErrMissingClientSecret
	}
	return nil
}

// Helper functions for testing (these wrap os functions to allow mocking).

func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func setEnv(key, value string) error {
	return os.Setenv(key, value)
}

func unsetEnv(key string) error {
	return os.Unsetenv(key)
}
