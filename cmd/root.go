package cmd

import (
	"fmt"
	"os"

	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "inflights",
	Short: "Inflights CLI",
	Long:  "The command-line interface for Inflights.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&output.JSONOutput, "json", false, "Output in JSON format")
}
