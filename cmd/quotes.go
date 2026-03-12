package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var quotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "List quotes",
	Long: `List quotes and price estimates.

Valid statuses:
  pending    Quote sent, awaiting client action
  accepted   Quote accepted by client`,
	RunE: runQuotes,
}

var quoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "Manage a quote",
}

var quoteShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show quote details",
	Args:  cobra.ExactArgs(1),
	RunE:  runQuoteShow,
}

var quoteConfirmCmd = &cobra.Command{
	Use:   "confirm [quote-number or flight-uid]",
	Short: "Accept a quote or price estimate",
	Long: `Accept a quote or price estimate.

Pass the quote number (e.g. Q-2026-001) or flight UID (e.g. BT779) as shown
in the quotes list. The command auto-detects whether it's a quote or estimate.`,
	Args: cobra.ExactArgs(1),
	RunE: runQuoteConfirm,
}

var quoteRejectCmd = &cobra.Command{
	Use:   "reject [quote-number or flight-uid]",
	Short: "Reject a quote or cancel an estimate",
	Long: `Reject a quote or cancel a price estimate.

Pass the quote number (e.g. Q-2026-001) or flight UID (e.g. BT779) as shown
in the quotes list. The command auto-detects whether it's a quote or estimate.`,
	Args: cobra.ExactArgs(1),
	RunE: runQuoteReject,
}

func init() {
	quotesCmd.Flags().String("status", "", "Filter by status (pending, accepted)")
	quoteCmd.AddCommand(quoteShowCmd)
	quoteCmd.AddCommand(quoteConfirmCmd)
	quoteCmd.AddCommand(quoteRejectCmd)
	rootCmd.AddCommand(quotesCmd)
	rootCmd.AddCommand(quoteCmd)
}

type quote struct {
	ID          int         `json:"id"`
	QuoteNumber string      `json:"quote_number"`
	Status      string      `json:"status"`
	Amount      json.Number `json:"amount"`
	VATPercent  json.Number `json:"vat_percent"`
	QuoteDate   string      `json:"quote_date"`
	DueDate     string      `json:"due_date"`
	CreatedAt   string      `json:"created_at"`
	Type        string      `json:"type"`
	// Estimate-only fields
	FlightID        int    `json:"flight_id,omitempty"`
	FlightPublicUID string `json:"flight_public_uid,omitempty"`
	Product         string `json:"product,omitempty"`
}

type quoteDetail struct {
	quote
	AcceptedAt string        `json:"accepted_at"`
	Flights    []quoteFlight `json:"flights"`
}

type quoteFlight struct {
	ID        int    `json:"id"`
	PublicUID string `json:"public_uid"`
	Status    string `json:"status"`
	Product   string `json:"product"`
}

func runQuotes(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	params := url.Values{}
	if s, _ := cmd.Flags().GetString("status"); s != "" {
		params.Set("status", s)
	}

	path := "/quotes"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := client.Get(path)
	if err != nil {
		return fmt.Errorf("failed to fetch quotes: %w", err)
	}

	var quotes []quote
	if err := json.Unmarshal(body, &quotes); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(quotes)
		return nil
	}

	if len(quotes) == 0 {
		fmt.Println("No quotes found.")
		return nil
	}

	rows := make([][]string, len(quotes))
	for i, q := range quotes {
		id := q.QuoteNumber
		if q.Type == "estimate" {
			id = q.FlightPublicUID + " (estimate)"
		}
		rows[i] = []string{
			id,
			q.Type,
			q.Status,
			q.Amount.String(),
			q.CreatedAt,
		}
	}
	output.Table([]string{"Quote", "Type", "Status", "Amount", "Created"}, rows)
	return nil
}

func runQuoteConfirm(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	identifier := args[0]

	// Fetch all quotes to find the matching one
	listBody, err := client.Get("/quotes")
	if err != nil {
		return fmt.Errorf("failed to fetch quotes: %w", err)
	}

	var quotes []quote
	if err := json.Unmarshal(listBody, &quotes); err != nil {
		return fmt.Errorf("failed to parse quotes: %w", err)
	}

	// Find matching quote or estimate
	var match *quote
	for i, q := range quotes {
		if q.QuoteNumber == identifier || q.FlightPublicUID == identifier {
			match = &quotes[i]
			break
		}
	}
	if match == nil {
		return fmt.Errorf("no quote or estimate found matching %q", identifier)
	}

	var body []byte
	if match.Type == "estimate" {
		body, err = client.Post("/quotes/accept_estimate", map[string]int{
			"flight_id": match.FlightID,
		})
	} else {
		body, err = client.Post(fmt.Sprintf("/quotes/%d/accept", match.ID), nil)
	}
	if err != nil {
		return fmt.Errorf("failed to confirm quote: %w", err)
	}

	if output.JSONOutput {
		var raw json.RawMessage
		json.Unmarshal(body, &raw)
		output.JSON(raw)
	} else {
		if match.Type == "estimate" {
			fmt.Printf("Estimate for flight %s accepted.\n", match.FlightPublicUID)
		} else {
			fmt.Printf("Quote %s accepted.\n", match.QuoteNumber)
		}
	}
	return nil
}

func runQuoteReject(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	identifier := args[0]

	// Fetch all quotes to find the matching one
	listBody, err := client.Get("/quotes")
	if err != nil {
		return fmt.Errorf("failed to fetch quotes: %w", err)
	}

	var quotes []quote
	if err := json.Unmarshal(listBody, &quotes); err != nil {
		return fmt.Errorf("failed to parse quotes: %w", err)
	}

	var match *quote
	for i, q := range quotes {
		if q.QuoteNumber == identifier || q.FlightPublicUID == identifier {
			match = &quotes[i]
			break
		}
	}
	if match == nil {
		return fmt.Errorf("no quote or estimate found matching %q", identifier)
	}

	var body []byte
	if match.Type == "estimate" {
		body, err = client.Post("/quotes/reject_estimate", map[string]int{
			"flight_id": match.FlightID,
		})
	} else {
		body, err = client.Post(fmt.Sprintf("/quotes/%d/reject", match.ID), nil)
	}
	if err != nil {
		return fmt.Errorf("failed to reject: %w", err)
	}

	if output.JSONOutput {
		var raw json.RawMessage
		json.Unmarshal(body, &raw)
		output.JSON(raw)
	} else {
		if match.Type == "estimate" {
			fmt.Printf("Estimate for flight %s rejected.\n", match.FlightPublicUID)
		} else {
			fmt.Printf("Quote %s rejected.\n", match.QuoteNumber)
		}
	}
	return nil
}

func runQuoteShow(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Get("/quotes/" + args[0])
	if err != nil {
		return fmt.Errorf("failed to fetch quote: %w", err)
	}

	var q quoteDetail
	if err := json.Unmarshal(body, &q); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(q)
		return nil
	}

	output.Print("Quote:", q.QuoteNumber)
	output.Print("Status:", q.Status)
	output.Print("Amount:", q.Amount.String())
	output.Print("VAT:", q.VATPercent.String()+"%")
	output.Print("Quote date:", q.QuoteDate)
	output.Print("Due date:", q.DueDate)
	output.Print("Accepted at:", q.AcceptedAt)

	if len(q.Flights) > 0 {
		fmt.Println()
		rows := make([][]string, len(q.Flights))
		for i, f := range q.Flights {
			rows[i] = []string{
				fmt.Sprintf("%d", f.ID),
				f.PublicUID,
				f.Status,
				f.Product,
			}
		}
		output.Table([]string{"Flight ID", "UID", "Status", "Product"}, rows)
	}
	return nil
}
