package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultBaseURL     = "https://inflights.com/api/v1"
	CredentialsDirName = ".inflights"
	CredentialsFile    = "credentials"
)

func BaseURL() string {
	if url := os.Getenv("INFLIGHTS_API_URL"); url != "" {
		return url
	}
	return DefaultBaseURL
}

// CredentialsPathOverride allows tests to redirect credential storage.
var CredentialsPathOverride string

func CredentialsPath() string {
	if CredentialsPathOverride != "" {
		return CredentialsPathOverride
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, CredentialsDirName, CredentialsFile)
}
