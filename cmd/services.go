package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List available services",
	RunE:  runServices,
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}

type service struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Subtitle     string   `json:"subtitle"`
	Description  string   `json:"description"`
	ProductType  string   `json:"product_type"`
	PriceMinimum float64  `json:"price_minimum"`
	Industries   []string `json:"industries"`
	SensorTypes  []string `json:"sensor_types"`
}

func runServices(cmd *cobra.Command, args []string) error {
	client := api.NewUnauthenticated()

	body, err := client.Get("/services")
	if err != nil {
		return fmt.Errorf("failed to fetch services: %w", err)
	}

	var services []service
	if err := json.Unmarshal(body, &services); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(services)
		return nil
	}

	if len(services) == 0 {
		fmt.Println("No services found.")
		return nil
	}

	rows := make([][]string, len(services))
	for i, s := range services {
		rows[i] = []string{
			fmt.Sprintf("%d", s.ID),
			s.Name,
			s.ProductType,
			fmt.Sprintf("%.0f€", s.PriceMinimum),
		}
	}
	output.Table([]string{"ID", "Name", "Type", "Min Price"}, rows)
	return nil
}
