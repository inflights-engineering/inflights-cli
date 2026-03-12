package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

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

var flightOrderCmd = &cobra.Command{
	Use:   "order [geojson-file]",
	Short: "Order a new flight",
	Long: `Order a new flight by providing a GeoJSON file defining the area.

The GeoJSON file must be a Feature with a GeometryCollection of Polygons.
Use --service to specify the service (get IDs from 'inflights services').

Example:
  inflights order area.geojson --service 3
  inflights order area.geojson --service 3 --description "Roof inspection"`,
	Args: cobra.ExactArgs(1),
	RunE: runFlightOrder,
}

func init() {
	flightsCmd.Flags().String("status", "", "Filter by status")
	flightsCmd.Flags().String("public-uid", "", "Filter by public UID")
	flightOrderCmd.Flags().Int("service", 0, "Service ID (required, see 'inflights services')")
	flightOrderCmd.Flags().String("description", "", "Description or notes for the flight")
	flightOrderCmd.MarkFlagRequired("service")
	rootCmd.AddCommand(flightsCmd)
	rootCmd.AddCommand(flightShowCmd)
	rootCmd.AddCommand(flightOrderCmd)
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

func runFlightOrder(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	// Read GeoJSON file
	geojsonBytes, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read GeoJSON file: %w", err)
	}

	var geojson map[string]any
	if err := json.Unmarshal(geojsonBytes, &geojson); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", args[0], err)
	}

	// Normalize GeoJSON into the expected Feature > GeometryCollection > [Polygon] structure
	areas, err := normalizeGeoJSON(geojson)
	if err != nil {
		return err
	}

	productID, _ := cmd.Flags().GetInt("service")
	description, _ := cmd.Flags().GetString("description")

	payload := map[string]any{
		"product_id":  productID,
		"areas":       map[string]any(areas),
		"skip_obtain": false,
	}
	if description != "" {
		payload["description_user"] = description
	}

	body, err := client.Post("/flights", payload)
	if err != nil {
		return fmt.Errorf("failed to create flight: %w", err)
	}

	var f flightDetail
	if err := json.Unmarshal(body, &f); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(f)
	} else {
		fmt.Printf("Flight %s created.\n", f.PublicUID)
		output.Print("ID:", fmt.Sprintf("%d", f.ID))
		output.Print("UID:", f.PublicUID)
		output.Print("Status:", f.Status)
		output.Print("Product:", f.Product)
	}
	return nil
}

// normalizeGeoJSON converts common GeoJSON formats into the
// Feature > GeometryCollection > [Polygon] structure the API expects.
//
// Accepted inputs:
//   - Feature with GeometryCollection of Polygons (already correct)
//   - Feature with a Polygon geometry (wrapped into GeometryCollection)
//   - Feature with a MultiPolygon geometry (each polygon becomes a geometry)
//   - FeatureCollection with polygon features (merged into one Feature)
//   - Bare Polygon geometry (wrapped into Feature > GeometryCollection)
//   - Bare MultiPolygon geometry (wrapped into Feature > GeometryCollection)
func normalizeGeoJSON(input map[string]any) (map[string]any, error) {
	typ, _ := input["type"].(string)

	switch typ {
	case "Feature":
		return normalizeFeature(input)
	case "FeatureCollection":
		return normalizeFeatureCollection(input)
	case "Polygon", "MultiPolygon":
		return normalizeFeature(map[string]any{
			"type":       "Feature",
			"geometry":   input,
			"properties": map[string]any{},
		})
	case "GeometryCollection":
		return normalizeFeature(map[string]any{
			"type":       "Feature",
			"geometry":   input,
			"properties": map[string]any{},
		})
	default:
		return nil, fmt.Errorf("unsupported GeoJSON type %q, expected Feature, FeatureCollection, Polygon, or MultiPolygon", typ)
	}
}

func normalizeFeature(feature map[string]any) (map[string]any, error) {
	geom, _ := feature["geometry"].(map[string]any)
	if geom == nil {
		return nil, fmt.Errorf("feature has no geometry")
	}

	geomType, _ := geom["type"].(string)

	switch geomType {
	case "GeometryCollection":
		// Already in the right format
		return feature, nil
	case "Polygon":
		feature["geometry"] = map[string]any{
			"type":       "GeometryCollection",
			"geometries": []any{geom},
		}
		return feature, nil
	case "MultiPolygon":
		coords, _ := geom["coordinates"].([]any)
		geometries := make([]any, len(coords))
		for i, ring := range coords {
			geometries[i] = map[string]any{
				"type":        "Polygon",
				"coordinates": ring,
			}
		}
		feature["geometry"] = map[string]any{
			"type":       "GeometryCollection",
			"geometries": geometries,
		}
		return feature, nil
	default:
		return nil, fmt.Errorf("unsupported geometry type %q, expected Polygon or MultiPolygon", geomType)
	}
}

func normalizeFeatureCollection(fc map[string]any) (map[string]any, error) {
	features, _ := fc["features"].([]any)
	if len(features) == 0 {
		return nil, fmt.Errorf("FeatureCollection has no features")
	}

	var geometries []any
	for _, f := range features {
		feat, _ := f.(map[string]any)
		geom, _ := feat["geometry"].(map[string]any)
		if geom == nil {
			continue
		}
		geomType, _ := geom["type"].(string)
		switch geomType {
		case "Polygon":
			geometries = append(geometries, geom)
		case "MultiPolygon":
			coords, _ := geom["coordinates"].([]any)
			for _, ring := range coords {
				geometries = append(geometries, map[string]any{
					"type":        "Polygon",
					"coordinates": ring,
				})
			}
		case "GeometryCollection":
			nested, _ := geom["geometries"].([]any)
			geometries = append(geometries, nested...)
		}
	}

	if len(geometries) == 0 {
		return nil, fmt.Errorf("no polygons found in FeatureCollection")
	}

	return map[string]any{
		"type": "Feature",
		"geometry": map[string]any{
			"type":       "GeometryCollection",
			"geometries": geometries,
		},
		"properties": map[string]any{},
	}, nil
}
