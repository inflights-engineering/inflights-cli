package cmd

import (
	"fmt"

	"github.com/inflights-engineering/inflights-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and clear stored credentials",
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := credentials.Delete(); err != nil {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}
	fmt.Println("Logged out.")
	return nil
}
