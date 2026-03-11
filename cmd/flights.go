package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var flightsCmd = &cobra.Command{
	Use:   "flights",
	Short: "List flights",
	Long: `List flights. Optionally filter by status or public UID.

Valid statuses:
  needs_flight_proposal   Needs pilot proposals
  proposal_pending        Proposals awaiting action
  price_not_final         Awaiting quote from Inflights
  quote_sent              Quote sent to client
  pilot_found             Pilot scheduling flight
  flight_scheduled        Flight date set
  flight_flown            Flight completed, awaiting upload
  raw_data_uploaded       Data uploaded, awaiting processing
  insights_generated      Processing complete
  done                    Invoice created`,
	RunE: runFlights,
}

var flightShowCmd = &cobra.Command{
	Use:   "flight [id]",
	Short: "Show flight details",
	Args:  cobra.ExactArgs(1),
	RunE:  runFlightShow,
}

func init() {
	flightsCmd.Flags().String("status", "", "Filter by status")
	flightsCmd.Flags().String("public-uid", "", "Filter by public UID")
	rootCmd.AddCommand(flightsCmd)
	rootCmd.AddCommand(flightShowCmd)
}

type flight struct {
	ID            int             `json:"id"`
	PublicUID     string          `json:"public_uid"`
	Status        string          `json:"status"`
	Product       string          `json:"product"`
	ScheduledDate string          `json:"scheduled_date"`
	AreaInHa      json.Number     `json:"area_in_ha"`
	PriceClient   json.Number     `json:"price_client"`
	CreatedAt     string          `json:"created_at"`
}

type flightDetail struct {
	flight
	FlownAt         string       `json:"flown_at"`
	CompletedAt     string       `json:"completed_at"`
	CancelledAt     string       `json:"cancelled_at"`
	DescriptionUser string       `json:"description_user"`
	Reference       string       `json:"reference"`
	Pilot           *flightActor `json:"pilot"`
	Customer        *flightActor `json:"customer"`
}

type flightActor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func runFlights(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	params := url.Values{}
	if s, _ := cmd.Flags().GetString("status"); s != "" {
		params.Set("status", s)
	}
	if uid, _ := cmd.Flags().GetString("public-uid"); uid != "" {
		params.Set("public_uid", uid)
	}

	path := "/flights"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := client.Get(path)
	if err != nil {
		return fmt.Errorf("failed to fetch flights: %w", err)
	}

	var flights []flight
	if err := json.Unmarshal(body, &flights); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(flights)
		return nil
	}

	if len(flights) == 0 {
		fmt.Println("No flights found.")
		return nil
	}

	rows := make([][]string, len(flights))
	for i, f := range flights {
		rows[i] = []string{
			fmt.Sprintf("%d", f.ID),
			f.PublicUID,
			f.Status,
			f.Product,
			f.ScheduledDate,
		}
	}
	output.Table([]string{"ID", "UID", "Status", "Product", "Scheduled"}, rows)
	return nil
}

func runFlightShow(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Get("/flights/" + args[0])
	if err != nil {
		return fmt.Errorf("failed to fetch flight: %w", err)
	}

	var f flightDetail
	if err := json.Unmarshal(body, &f); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(f)
		return nil
	}

	output.Print("ID:", fmt.Sprintf("%d", f.ID))
	output.Print("UID:", f.PublicUID)
	output.Print("Status:", f.Status)
	output.Print("Product:", f.Product)
	output.Print("Scheduled:", f.ScheduledDate)
	output.Print("Area (ha):", f.AreaInHa.String())
	output.Print("Price:", f.PriceClient.String()+"€")
	output.Print("Description:", f.DescriptionUser)
	output.Print("Reference:", f.Reference)
	if f.Pilot != nil {
		output.Print("Pilot:", f.Pilot.Name)
	}
	if f.Customer != nil {
		output.Print("Customer:", f.Customer.Name)
	}
	return nil
}
