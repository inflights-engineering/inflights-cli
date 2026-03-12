package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var gearCmd = &cobra.Command{
	Use:   "gear",
	Short: "Manage equipment",
}

var gearListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available equipment types",
	Long: `List available equipment types. Optionally filter by category.

Valid categories:
  drone              Drones
  payload            Payloads / sensors
  drone_and_payload  Combined drone + payload
  gnss_receiver      GNSS receivers`,
	RunE: runGearList,
}

var gearMineCmd = &cobra.Command{
	Use:   "mine",
	Short: "List your equipment",
	RunE:  runGearMine,
}

var gearAddCmd = &cobra.Command{
	Use:   "add [equipment-type-id]",
	Short: "Add equipment to your profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runGearAdd,
}

var gearRemoveCmd = &cobra.Command{
	Use:   "remove [equipment-id]",
	Short: "Remove equipment from your profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runGearRemove,
}

func init() {
	gearListCmd.Flags().String("category", "", "Filter by category (drone, payload, drone_and_payload, gnss_receiver)")
	gearCmd.AddCommand(gearListCmd)
	gearCmd.AddCommand(gearMineCmd)
	gearCmd.AddCommand(gearAddCmd)
	gearCmd.AddCommand(gearRemoveCmd)
	rootCmd.AddCommand(gearCmd)
}

type equipmentType struct {
	ID          string   `json:"id"`
	Brand       string   `json:"brand"`
	ProductName string   `json:"product_name"`
	Category    string   `json:"category"`
	SensorTypes []string `json:"sensor_types"`
	Resolution  *int     `json:"resolution"`
}

type equipment struct {
	ID            string         `json:"id"`
	EquipmentType *equipmentType `json:"equipment_type"`
	FullDayRate   json.Number    `json:"full_day_rate"`
	PriceMinimum  json.Number    `json:"price_minimum"`
	PricePerHa    json.Number    `json:"price_per_ha"`
	SurfaceMin    json.Number    `json:"surface_minimum"`
}

func runGearList(cmd *cobra.Command, args []string) error {
	client := api.NewUnauthenticated()

	params := url.Values{}
	if c, _ := cmd.Flags().GetString("category"); c != "" {
		params.Set("category", c)
	}

	path := "/equipment_types"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := client.Get(path)
	if err != nil {
		return fmt.Errorf("failed to fetch equipment types: %w", err)
	}

	var types []equipmentType
	if err := json.Unmarshal(body, &types); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(types)
		return nil
	}

	if len(types) == 0 {
		fmt.Println("No equipment types found.")
		return nil
	}

	rows := make([][]string, len(types))
	for i, t := range types {
		rows[i] = []string{
			t.ID,
			t.Brand,
			t.ProductName,
			t.Category,
		}
	}
	output.Table([]string{"ID", "Brand", "Product", "Category"}, rows)
	return nil
}

func runGearMine(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Get("/equipments")
	if err != nil {
		return fmt.Errorf("failed to fetch equipment: %w", err)
	}

	var equipments []equipment
	if err := json.Unmarshal(body, &equipments); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(equipments)
		return nil
	}

	if len(equipments) == 0 {
		fmt.Println("No equipment found.")
		return nil
	}

	rows := make([][]string, len(equipments))
	for i, e := range equipments {
		brand, product := "", ""
		if e.EquipmentType != nil {
			brand = e.EquipmentType.Brand
			product = e.EquipmentType.ProductName
		}
		rows[i] = []string{
			e.ID,
			brand,
			product,
			e.PricePerHa.String(),
			e.PriceMinimum.String(),
		}
	}
	output.Table([]string{"ID", "Brand", "Product", "Price/ha", "Min Price"}, rows)
	return nil
}

func runGearAdd(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Post("/equipments", map[string]string{
		"equipment_type_id": args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to add equipment: %w", err)
	}

	var e equipment
	if err := json.Unmarshal(body, &e); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(e)
	} else {
		brand, product := "", ""
		if e.EquipmentType != nil {
			brand = e.EquipmentType.Brand
			product = e.EquipmentType.ProductName
		}
		fmt.Printf("Added %s %s to your equipment.\n", brand, product)
	}
	return nil
}

func runGearRemove(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	_, err = client.Delete("/equipments/" + args[0])
	if err != nil {
		return fmt.Errorf("failed to remove equipment: %w", err)
	}

	if output.JSONOutput {
		output.JSON(map[string]string{"status": "removed", "id": args[0]})
	} else {
		fmt.Println("Equipment removed.")
	}
	return nil
}
