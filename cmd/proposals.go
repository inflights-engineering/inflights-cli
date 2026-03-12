package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var proposalListCmd = &cobra.Command{
	Use:   "proposals",
	Short: "List proposals",
	RunE:  runProposalList,
}

var proposalCmd = &cobra.Command{
	Use:   "proposal",
	Short: "Manage a proposal",
}

var proposalShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show proposal details",
	Args:  cobra.ExactArgs(1),
	RunE:  runProposalShow,
}

var proposalAcceptCmd = &cobra.Command{
	Use:   "accept [id]",
	Short: "Accept a proposal",
	Args:  cobra.ExactArgs(1),
	RunE:  runProposalAccept,
}

var proposalRejectCmd = &cobra.Command{
	Use:   "reject [id]",
	Short: "Reject a proposal",
	Args:  cobra.ExactArgs(1),
	RunE:  runProposalReject,
}

func init() {
	proposalListCmd.Flags().String("status", "", "Filter by status")
	proposalRejectCmd.Flags().String("reason", "", "Reason for rejection")
	proposalCmd.AddCommand(proposalShowCmd)
	proposalCmd.AddCommand(proposalAcceptCmd)
	proposalCmd.AddCommand(proposalRejectCmd)
	rootCmd.AddCommand(proposalListCmd)
	rootCmd.AddCommand(proposalCmd)
}

type proposal struct {
	ID              string      `json:"id"`
	Status          string      `json:"status"`
	FlightID        int         `json:"flight_id"`
	FlightPublicUID string      `json:"flight_public_uid"`
	ScheduledDate   string      `json:"scheduled_date"`
	PricePilot      json.Number `json:"price_pilot"`
	CreatedAt       string      `json:"created_at"`
}

type proposalDetail struct {
	proposal
	BackupScheduledDate  string                `json:"backup_scheduled_date"`
	ReasonForRejection   string                `json:"reason_for_rejection"`
	EquipmentType        *proposalEquipment    `json:"equipment_type"`
	Flight               *proposalFlightDetail `json:"flight"`
}

type proposalEquipment struct {
	ID          string `json:"id"`
	Brand       string `json:"brand"`
	ProductName string `json:"product_name"`
}

type proposalFlightDetail struct {
	ID              int         `json:"id"`
	PublicUID       string      `json:"public_uid"`
	Status          string      `json:"status"`
	Product         string      `json:"product"`
	AreaInHa        json.Number `json:"area_in_ha"`
	DescriptionUser string      `json:"description_user"`
}

func runProposalList(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	params := url.Values{}
	if s, _ := cmd.Flags().GetString("status"); s != "" {
		params.Set("status", s)
	}

	path := "/proposals"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := client.Get(path)
	if err != nil {
		return fmt.Errorf("failed to fetch proposals: %w", err)
	}

	var proposals []proposal
	if err := json.Unmarshal(body, &proposals); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(proposals)
		return nil
	}

	if len(proposals) == 0 {
		fmt.Println("No proposals found.")
		return nil
	}

	rows := make([][]string, len(proposals))
	for i, p := range proposals {
		rows[i] = []string{
			p.ID,
			p.Status,
			p.FlightPublicUID,
			p.PricePilot.String(),
		}
	}
	output.Table([]string{"ID", "Status", "Flight", "Price"}, rows)
	return nil
}

func runProposalAccept(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Post("/proposals/"+args[0]+"/accept", nil)
	if err != nil {
		return fmt.Errorf("failed to accept proposal: %w", err)
	}

	var p proposalDetail
	if err := json.Unmarshal(body, &p); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(p)
	} else {
		fmt.Printf("Proposal %s accepted.\n", p.ID)
	}
	return nil
}

func runProposalReject(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	reqBody := map[string]string{}
	if reason, _ := cmd.Flags().GetString("reason"); reason != "" {
		reqBody["reason"] = reason
	}

	body, err := client.Post("/proposals/"+args[0]+"/reject", reqBody)
	if err != nil {
		return fmt.Errorf("failed to reject proposal: %w", err)
	}

	var p proposalDetail
	if err := json.Unmarshal(body, &p); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(p)
	} else {
		fmt.Printf("Proposal %s rejected.\n", p.ID)
	}
	return nil
}

func runProposalShow(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Get("/proposals/" + args[0])
	if err != nil {
		return fmt.Errorf("failed to fetch proposal: %w", err)
	}

	var p proposalDetail
	if err := json.Unmarshal(body, &p); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(p)
		return nil
	}

	output.Print("ID:", p.ID)
	output.Print("Status:", p.Status)
	output.Print("Price:", p.PricePilot.String())
	output.Print("Scheduled:", p.ScheduledDate)
	output.Print("Backup date:", p.BackupScheduledDate)
	if p.ReasonForRejection != "" {
		output.Print("Rejected:", p.ReasonForRejection)
	}
	if p.EquipmentType != nil {
		output.Print("Equipment:", fmt.Sprintf("%s %s", p.EquipmentType.Brand, p.EquipmentType.ProductName))
	}
	if p.Flight != nil {
		fmt.Println()
		output.Print("Flight ID:", fmt.Sprintf("%d", p.Flight.ID))
		output.Print("Flight UID:", p.Flight.PublicUID)
		output.Print("Flight status:", p.Flight.Status)
		output.Print("Product:", p.Flight.Product)
		output.Print("Area (ha):", p.Flight.AreaInHa.String())
		output.Print("Description:", p.Flight.DescriptionUser)
	}
	return nil
}
