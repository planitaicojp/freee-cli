package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"

	"github.com/planitaicojp/freee-cli/internal/config"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
)

const (
	authURL             = "https://accounts.secure.freee.co.jp/public_api/authorize"
	tokenURL            = "https://accounts.secure.freee.co.jp/public_api/token"
	defaultCallbackPort = 8080
	defaultScope        = "read write"
)

// OAuthConfig creates an OAuth2 config for freee.
func OAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: redirectURL,
		Scopes:      []string{"read", "write"},
	}
}

// generatePKCE generates a PKCE code verifier and challenge (S256).
func generatePKCE() (codeVerifier, codeChallenge string) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	codeVerifier = base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(codeVerifier))
	codeChallenge = base64.RawURLEncoding.EncodeToString(h[:])
	return
}

// LoginResult holds the result of an OAuth2 login flow.
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// Login performs the OAuth2 Authorization Code flow with a local callback server.
// The redirect URI must match exactly what is registered in the freee developer console.
// Default: http://localhost:8080/callback
func Login(clientID, clientSecret string) (*LoginResult, error) {
	port := defaultCallbackPort
	redirectURL := fmt.Sprintf("http://localhost:%d/callback", port)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, fmt.Errorf("starting callback server on port %d (is it already in use?): %w", port, err)
	}

	oauthCfg := OAuthConfig(clientID, clientSecret, redirectURL)

	// Generate state parameter
	state, err := randomString(16)
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("generating state: %w", err)
	}

	// Channel to receive the authorization code
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			errCh <- fmt.Errorf("state mismatch")
			http.Error(w, "State mismatch", http.StatusBadRequest)
			return
		}
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			errCh <- fmt.Errorf("authorization error: %s", errMsg)
			fmt.Fprintf(w, "<html><body><h1>Authorization Failed</h1><p>%s</p></body></html>", errMsg)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			http.Error(w, "No code received", http.StatusBadRequest)
			return
		}
		codeCh <- code
		fmt.Fprint(w, "<html><body><h1>Authorization Successful</h1><p>You can close this window.</p></body></html>")
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	defer server.Shutdown(context.Background())

	// PKCE (Proof Key for Code Exchange)
	codeVerifier, codeChallenge := generatePKCE()

	// Open browser
	authURLStr := oauthCfg.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	fmt.Fprintf(os.Stderr, "Opening browser for authorization...\n")
	fmt.Fprintf(os.Stderr, "If the browser doesn't open, visit:\n%s\n\n", authURLStr)
	openBrowser(authURLStr)

	// Wait for callback
	fmt.Fprintf(os.Stderr, "Waiting for authorization...\n")
	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authorization timeout (5 minutes)")
	}

	// Exchange code for token (with PKCE code_verifier)
	token, err := oauthCfg.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		return nil, &cerrors.AuthError{Message: fmt.Sprintf("token exchange failed: %v", err)}
	}

	return &LoginResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

// EnsureToken returns a valid access token, refreshing if necessary.
func EnsureToken(profile string, cred config.OAuthCredentials, cfg *config.Config) (string, error) {
	// Priority 1: environment variable
	if t := os.Getenv(config.EnvToken); t != "" {
		return t, nil
	}

	// Priority 2: valid cached token
	if time.Until(cred.ExpiresAt) > 5*time.Minute {
		return cred.AccessToken, nil
	}

	// Priority 3: refresh token
	if cred.RefreshToken == "" {
		return "", &cerrors.AuthError{Message: "no valid token, run 'freee auth login'"}
	}

	fmt.Fprintf(os.Stderr, "Refreshing token...\n")
	oauthCfg := OAuthConfig(cred.ClientID, cred.ClientSecret, "")
	token, err := oauthCfg.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: cred.RefreshToken,
	}).Token()
	if err != nil {
		return "", &cerrors.AuthError{Message: fmt.Sprintf("token refresh failed: %v", err)}
	}

	// Update stored credentials
	creds, err := config.LoadCredentials()
	if err != nil {
		return "", err
	}
	creds.Set(profile, config.OAuthCredentials{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
		ClientID:     cred.ClientID,
		ClientSecret: cred.ClientSecret,
	})
	if err := creds.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save refreshed token: %v\n", err)
	}

	return token.AccessToken, nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	if cmd != nil {
		_ = cmd.Start()
	}
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
