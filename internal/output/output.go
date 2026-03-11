package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

var JSONOutput bool

// JSON prints data as formatted JSON.
func JSON(data any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}

// Table prints data as a formatted table.
func Table(headers []string, rows [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headerRow := table.Row{}
	for _, h := range headers {
		headerRow = append(headerRow, h)
	}
	t.AppendHeader(headerRow)

	for _, row := range rows {
		tableRow := table.Row{}
		for _, cell := range row {
			tableRow = append(tableRow, cell)
		}
		t.AppendRow(tableRow)
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}

// Print prints a simple key-value pair.
func Print(label, value string) {
	fmt.Printf("%-15s %s\n", label, value)
}

// Error prints an error message to stderr.
func Error(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
}
