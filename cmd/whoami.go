package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show the current authenticated user",
	RunE:  runWhoami,
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

type meResponse struct {
	User struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	} `json:"user"`
}

func runWhoami(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Get("/auth/me")
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	var resp meResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(resp.User)
	} else {
		output.Print("Email:", resp.User.Email)
		output.Print("Name:", fmt.Sprintf("%s %s", resp.User.FirstName, resp.User.LastName))
		output.Print("Role:", resp.User.Role)
	}
	return nil
}
