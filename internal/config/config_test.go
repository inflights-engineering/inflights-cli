package config

import (
	"os"
	"testing"
)

func TestBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "returns default when env is not set",
			envValue: "",
			want:     DefaultBaseURL,
		},
		{
			name:     "returns env value when set",
			envValue: "http://localhost:3000/api/v1",
			want:     "http://localhost:3000/api/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			os.Unsetenv("INFLIGHTS_API_URL")
			if tt.envValue != "" {
				os.Setenv("INFLIGHTS_API_URL", tt.envValue)
				defer os.Unsetenv("INFLIGHTS_API_URL")
			}

			// Act
			got := BaseURL()

			// Assert
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCredentialsPath(t *testing.T) {
	tests := []struct {
		name     string
		override string
		want     string
	}{
		{
			name:     "returns override when set",
			override: "/tmp/test-creds",
			want:     "/tmp/test-creds",
		},
		{
			name:     "returns non-empty default when no override",
			override: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			CredentialsPathOverride = tt.override
			defer func() { CredentialsPathOverride = "" }()

			// Act
			got := CredentialsPath()

			// Assert
			if tt.want != "" && got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if tt.want == "" && got == "" {
				t.Error("got empty string, want non-empty default path")
			}
		})
	}
}
