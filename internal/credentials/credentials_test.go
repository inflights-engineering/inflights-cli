package credentials

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/config"
)

// setupTempCreds points credential storage at a temp directory for the test.
// It returns the path to the temp credentials file.
func setupTempCreds(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials")
	config.CredentialsPathOverride = path
	t.Cleanup(func() { config.CredentialsPathOverride = "" })
	return path
}

func TestSave(t *testing.T) {
	t.Run("creates file with correct permissions", func(t *testing.T) {
		// Arrange
		path := setupTempCreds(t)

		// Act
		err := Save("my-token")

		// Assert
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("file not created: %v", err)
		}
		if perm := info.Mode().Perm(); perm != 0600 {
			t.Errorf("got permissions %o, want 0600", perm)
		}
	})
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(path string) // write file content before Load()
		wantToken string
		wantErr   bool
	}{
		{
			name:    "errors when no file exists",
			setup:   func(path string) {},
			wantErr: true,
		},
		{
			name: "errors on corrupted JSON",
			setup: func(path string) {
				os.WriteFile(path, []byte(`not json`), 0600)
			},
			wantErr: true,
		},
		{
			name: "errors on empty token",
			setup: func(path string) {
				os.WriteFile(path, []byte(`{"token":""}`), 0600)
			},
			wantErr: true,
		},
		{
			name: "returns token from valid file",
			setup: func(path string) {
				Save("valid-token-123")
			},
			wantToken: "valid-token-123",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			path := setupTempCreds(t)
			tt.setup(path)

			// Act
			creds, err := Load()

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Error("got nil error, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if creds.Token != tt.wantToken {
				t.Errorf("got token %q, want %q", creds.Token, tt.wantToken)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		saveFirst  bool
		wantErr    bool
		wantGone   bool
	}{
		{
			name:      "removes existing credentials file",
			saveFirst: true,
			wantErr:   false,
			wantGone:  true,
		},
		{
			name:      "succeeds when no file exists",
			saveFirst: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			path := setupTempCreds(t)
			if tt.saveFirst {
				Save("token-to-delete")
			}

			// Act
			err := Delete()

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Error("got nil error, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if tt.wantGone {
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Error("file still exists after Delete()")
				}
			}
		})
	}
}
