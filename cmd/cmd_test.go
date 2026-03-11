package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/inflights-engineering/inflights-cli/internal/config"
	"github.com/inflights-engineering/inflights-cli/internal/output"
)

// setupTestEnv points the API and credentials at a temp dir and test server.
// Returns a function to capture stdout.
func setupTestEnv(t *testing.T, server *httptest.Server) {
	t.Helper()

	// Point API client at test server
	os.Setenv("INFLIGHTS_API_URL", server.URL)
	t.Cleanup(func() { os.Unsetenv("INFLIGHTS_API_URL") })

	// Point credentials at temp dir
	dir := t.TempDir()
	config.CredentialsPathOverride = filepath.Join(dir, "credentials")
	t.Cleanup(func() { config.CredentialsPathOverride = "" })

	// Reset --json flag
	output.JSONOutput = false
}

// captureOutput captures stdout during fn execution.
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// saveTestCredentials writes a token to the temp credentials file.
func saveTestCredentials(t *testing.T, token string) {
	t.Helper()
	path := config.CredentialsPath()
	os.MkdirAll(filepath.Dir(path), 0700)
	data, _ := json.Marshal(map[string]string{"token": token})
	os.WriteFile(path, data, 0600)
}

// --- Root command tests ---

func TestRoot_HasExpectedSubcommands(t *testing.T) {
	tests := []struct {
		name    string
		command string
	}{
		{name: "login is registered", command: "login"},
		{name: "logout is registered", command: "logout"},
		{name: "whoami is registered", command: "whoami"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange & Act
			found := false
			for _, cmd := range rootCmd.Commands() {
				if cmd.Name() == tt.command {
					found = true
					break
				}
			}

			// Assert
			if !found {
				t.Errorf("subcommand %q not found on root", tt.command)
			}
		})
	}
}

func TestRoot_JSONFlagExists(t *testing.T) {
	// Arrange & Act
	flag := rootCmd.PersistentFlags().Lookup("json")

	// Assert
	if flag == nil {
		t.Fatal("--json flag not found on root command")
	}
	if flag.DefValue != "false" {
		t.Errorf("--json default = %q, want %q", flag.DefValue, "false")
	}
}

// --- Login command tests ---

func TestLogin_FullFlow(t *testing.T) {
	// Arrange
	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/login_tokens":
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]string{
				"login_token": "test-uuid-123",
				"login_url":   "http://localhost/cli/login/test-uuid-123",
			})
		case "/auth/token_exchange":
			pollCount++
			if pollCount < 2 {
				// First poll: pending
				w.WriteHeader(202)
				json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
			} else {
				// Second poll: success
				json.NewEncoder(w).Encode(map[string]any{
					"token": "bearer-token-abc",
					"user": map[string]string{
						"id":         "user-uuid",
						"email":      "test@inflights.com",
						"first_name": "Test",
						"last_name":  "User",
						"role":       "customer",
					},
				})
			}
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)

	// Disable browser opening and speed up polling
	origBrowser := openBrowserFn
	openBrowserFn = func(url string) {}
	defer func() { openBrowserFn = origBrowser }()

	origInterval := pollInterval
	pollInterval = 10 * time.Millisecond
	defer func() { pollInterval = origInterval }()

	// Act
	out := captureOutput(t, func() {
		err := runLogin(loginCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if pollCount < 2 {
		t.Errorf("got %d polls, want at least 2", pollCount)
	}

	// Check credentials were saved
	data, err := os.ReadFile(config.CredentialsPath())
	if err != nil {
		t.Fatalf("credentials file not created: %v", err)
	}
	var creds map[string]string
	json.Unmarshal(data, &creds)
	if creds["token"] != "bearer-token-abc" {
		t.Errorf("got saved token %q, want %q", creds["token"], "bearer-token-abc")
	}

	// Check output
	if !bytes.Contains([]byte(out), []byte("test@inflights.com")) {
		t.Errorf("output = %q, want it to contain email", out)
	}
}

func TestLogin_APIError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"error":{"id":"server_error","message":"Something broke"}}`))
	}))
	defer server.Close()
	setupTestEnv(t, server)

	openBrowserFn = func(url string) {}
	defer func() { openBrowserFn = openBrowser }()

	// Act
	err := runLogin(loginCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error")
	}
}

// --- Logout command tests ---

func TestLogout(t *testing.T) {
	tests := []struct {
		name          string
		saveFirst     bool
		wantFileGone  bool
	}{
		{
			name:         "removes credentials file",
			saveFirst:    true,
			wantFileGone: true,
		},
		{
			name:      "succeeds when not logged in",
			saveFirst: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			defer server.Close()
			setupTestEnv(t, server)

			if tt.saveFirst {
				saveTestCredentials(t, "token-to-remove")
			}

			// Act
			out := captureOutput(t, func() {
				err := runLogout(logoutCmd, []string{})
				if err != nil {
					t.Fatalf("got error %v, want nil", err)
				}
			})

			// Assert
			if tt.wantFileGone {
				if _, err := os.Stat(config.CredentialsPath()); !os.IsNotExist(err) {
					t.Error("credentials file still exists after logout")
				}
			}
			if !bytes.Contains([]byte(out), []byte("Logged out")) {
				t.Errorf("output = %q, want it to contain 'Logged out'", out)
			}
		})
	}
}

// --- Whoami command tests ---

func TestWhoami(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header is sent
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token-xyz" {
			w.WriteHeader(401)
			w.Write([]byte(`{"error":{"id":"unauthorized","message":"Invalid token"}}`))
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"user": map[string]string{
				"id":         "user-uuid",
				"email":      "elias@inflights.com",
				"first_name": "Elias",
				"last_name":  "Music",
				"role":       "admin",
			},
		})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token-xyz")

	// Act
	out := captureOutput(t, func() {
		err := runWhoami(whoamiCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !bytes.Contains([]byte(out), []byte("elias@inflights.com")) {
		t.Errorf("output = %q, want it to contain email", out)
	}
	if !bytes.Contains([]byte(out), []byte("admin")) {
		t.Errorf("output = %q, want it to contain role", out)
	}
}

func TestWhoami_NotLoggedIn(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)
	// Don't save credentials

	// Act
	err := runWhoami(whoamiCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error when not logged in")
	}
}

func TestWhoami_InvalidToken(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":{"id":"unauthorized","message":"Invalid token"}}`))
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "expired-token")

	// Act
	err := runWhoami(whoamiCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error for invalid token")
	}
}
