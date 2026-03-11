package credentials

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/inflights-engineering/inflights-cli/internal/config"
)

type Credentials struct {
	Token string `json:"token"`
}

func Load() (*Credentials, error) {
	path := config.CredentialsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("not logged in (run 'inflights login' first)")
	}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("corrupted credentials file: %w", err)
	}
	if creds.Token == "" {
		return nil, fmt.Errorf("empty token in credentials file")
	}
	return &creds, nil
}

func Save(token string) error {
	path := config.CredentialsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	creds := Credentials{Token: token}
	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func Delete() error {
	path := config.CredentialsPath()
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
