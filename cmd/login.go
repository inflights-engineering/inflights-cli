package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/credentials"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

// Overridable for tests.
var (
	openBrowserFn = openBrowser
	pollInterval  = 2 * time.Second
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Inflights",
	RunE:  runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

type loginTokenResponse struct {
	LoginToken string `json:"login_token"`
	LoginURL   string `json:"login_url"`
}

type tokenExchangeResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	User   struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	} `json:"user"`
}

func runLogin(cmd *cobra.Command, args []string) error {
	client := api.NewUnauthenticated()

	// Step 1: Create login token
	body, err := client.Post("/auth/login_tokens", nil)
	if err != nil {
		return fmt.Errorf("failed to create login token: %w", err)
	}

	var loginResp loginTokenResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	// Step 2: Open browser
	fmt.Printf("Opening browser… %s\n", loginResp.LoginURL)
	openBrowserFn(loginResp.LoginURL)

	// Step 3: Poll for token exchange
	fmt.Println("Waiting for authentication…")
	for {
		time.Sleep(pollInterval)

		body, err := client.Post("/auth/token_exchange", map[string]string{
			"login_token": loginResp.LoginToken,
		})
		if err != nil {
			return fmt.Errorf("token exchange failed: %w", err)
		}

		var exchangeResp tokenExchangeResponse
		if err := json.Unmarshal(body, &exchangeResp); err != nil {
			return fmt.Errorf("failed to parse exchange response: %w", err)
		}

		if exchangeResp.Status == "pending" {
			continue
		}

		// Step 4: Store token
		if err := credentials.Save(exchangeResp.Token); err != nil {
			return fmt.Errorf("failed to save credentials: %w", err)
		}

		if output.JSONOutput {
			output.JSON(exchangeResp.User)
		} else {
			fmt.Printf("Authenticated as %s (%s)\n", exchangeResp.User.Email, exchangeResp.User.Role)
		}
		return nil
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
