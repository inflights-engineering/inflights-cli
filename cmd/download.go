package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [flight-id or public-uid]",
	Short: "Download flight deliverables",
	Long: `Download all deliverables and images for a flight.

Files are saved to the current directory, or to the directory specified by --output.`,
	Args: cobra.ExactArgs(1),
	RunE: runDownload,
}

func init() {
	downloadCmd.Flags().StringP("output", "o", ".", "Output directory")
	rootCmd.AddCommand(downloadCmd)
}

type downloadResponse struct {
	Documents  []downloadDocument  `json:"documents"`
	PictureSet *downloadPictureSet `json:"picture_set"`
}

type downloadDocument struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Filename    string `json:"filename"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

type downloadPictureSet struct {
	ID          string `json:"id"`
	DownloadURL string `json:"download_url"`
}

func runDownload(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	body, err := client.Get(fmt.Sprintf("/flights/%s/downloads", args[0]))
	if err != nil {
		return fmt.Errorf("failed to fetch downloads: %w", err)
	}

	var dl downloadResponse
	if err := json.Unmarshal(body, &dl); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if output.JSONOutput {
		output.JSON(dl)
		return nil
	}

	if len(dl.Documents) == 0 && dl.PictureSet == nil {
		fmt.Println("No downloads available for this flight.")
		return nil
	}

	outDir, _ := cmd.Flags().GetString("output")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, doc := range dl.Documents {
		if doc.DownloadURL == "" {
			continue
		}
		filename := doc.Filename
		if filename == "" {
			filename = doc.Title
		}
		dest := filepath.Join(outDir, filename)
		fmt.Printf("Downloading %s...\n", filename)
		if err := downloadFile(dest, doc.DownloadURL); err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		fmt.Printf("  Saved to %s\n", dest)
	}

	if dl.PictureSet != nil && dl.PictureSet.DownloadURL != "" {
		dest := filepath.Join(outDir, "images.zip")
		fmt.Printf("Downloading images.zip...\n")
		if err := downloadFile(dest, dl.PictureSet.DownloadURL); err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Saved to %s\n", dest)
		}
	}

	return nil
}

func downloadFile(dest, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
