package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	goexif "github.com/rwcarlsen/goexif/exif"

	"github.com/inflights-engineering/inflights-cli/internal/api"
	"github.com/inflights-engineering/inflights-cli/internal/output"
	"github.com/spf13/cobra"
)

const defaultConcurrency = 5

var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".tif":  true,
	".tiff": true,
	".dng":  true,
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to a flight",
}

var uploadDataCmd = &cobra.Command{
	Use:   "data [flight-id or public-uid] [files...]",
	Short: "Upload data files (deliverables)",
	Long: `Upload processed deliverables to a flight.

Each file goes through a 3-step process: presign → upload to S3 → confirm.`,
	Args: cobra.MinimumNArgs(2),
	RunE: runUploadData,
}

var uploadImagesCmd = &cobra.Command{
	Use:   "images [flight-id or public-uid] [path...]",
	Short: "Upload images to a flight",
	Long: `Upload images to a flight.

Accepts individual files or directories. When a directory is given,
it is scanned for image files (.jpg, .jpeg, .png, .tif, .tiff, .dng).

Each image goes through presign → upload to S3 → confirm.
After all images are uploaded, the dataset is finalized.
EXIF metadata (GPS, altitude, camera) is automatically extracted and sent.`,
	Args: cobra.MinimumNArgs(2),
	RunE: runUploadImages,
}

func init() {
	uploadDataCmd.Flags().IntP("concurrency", "c", defaultConcurrency, "Number of parallel uploads")
	uploadDataCmd.Flags().Int("deliverable", 0, "Deliverable type ID (see 'inflights flight <uid>')")
	uploadImagesCmd.Flags().IntP("concurrency", "c", defaultConcurrency, "Number of parallel uploads")
	uploadCmd.AddCommand(uploadDataCmd)
	uploadCmd.AddCommand(uploadImagesCmd)
	rootCmd.AddCommand(uploadCmd)
}

type presignResponse struct {
	FileID      string      `json:"file_id"`
	PresignData presignData `json:"presign_data"`
}

type presignData struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

type confirmResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type finalizeResponse struct {
	DatasetID     string `json:"dataset_id"`
	DatasetStatus string `json:"dataset_status"`
	PictureCount  int    `json:"picture_count"`
}

type uploadResult struct {
	index    int
	response *confirmResponse
	err      error
	filePath string
}

// resolveImageFiles takes a mix of files and directories and returns all image file paths.
func resolveImageFiles(paths []string) ([]string, error) {
	var files []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("cannot access %s: %w", p, err)
		}
		if !info.IsDir() {
			files = append(files, p)
			continue
		}
		entries, err := os.ReadDir(p)
		if err != nil {
			return nil, fmt.Errorf("cannot read directory %s: %w", p, err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if imageExtensions[ext] {
				files = append(files, filepath.Join(p, e.Name()))
			}
		}
	}
	return files, nil
}

// extractEXIF reads EXIF metadata from an image file.
func extractEXIF(filePath string) map[string]any {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	x, err := goexif.Decode(f)
	if err != nil {
		return nil
	}

	exifData := map[string]any{}

	if lat, lon, err := x.LatLong(); err == nil {
		exifData["latitude"] = lat
		exifData["longitude"] = lon
	}

	if tag, err := x.Get(goexif.GPSAltitude); err == nil {
		if num, den, err := tag.Rat2(0); err == nil && den != 0 {
			exifData["altitude"] = float64(num) / float64(den)
		}
	}

	if tm, err := x.DateTime(); err == nil {
		exifData["date_time_original"] = tm.Format("2006-01-02T15:04:05Z")
	}

	if tag, err := x.Get(goexif.Make); err == nil {
		exifData["make"] = strings.TrimSpace(tag.String())
	}
	if tag, err := x.Get(goexif.Model); err == nil {
		exifData["model"] = strings.TrimSpace(tag.String())
	}

	if len(exifData) == 0 {
		return nil
	}
	return exifData
}

type uploadOptions struct {
	withEXIF         bool
	productElementID int
}

func uploadParallel(client *api.Client, flightID string, files []string, endpoint string, concurrency int, opts uploadOptions) []uploadResult {
	results := make([]uploadResult, len(files))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var completed int64
	total := int64(len(files))

	for i, filePath := range files {
		wg.Add(1)
		go func(i int, filePath string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			resp, err := uploadSingleFile(client, flightID, filePath, endpoint, opts)
			results[i] = uploadResult{
				index:    i,
				response: resp,
				err:      err,
				filePath: filePath,
			}

			n := atomic.AddInt64(&completed, 1)
			if !output.JSONOutput {
				if err != nil {
					fmt.Printf("[%d/%d] %s — failed: %v\n", n, total, filepath.Base(filePath), err)
				} else {
					fmt.Printf("[%d/%d] %s — done\n", n, total, filepath.Base(filePath))
				}
			}
		}(i, filePath)
	}

	wg.Wait()
	return results
}

func collectResults(results []uploadResult) (confirmed []confirmResponse, failed []string) {
	for _, r := range results {
		if r.err != nil {
			failed = append(failed, filepath.Base(r.filePath))
			continue
		}
		confirmed = append(confirmed, *r.response)
	}
	return
}

func runUploadData(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	flightID := args[0]
	files := args[1:]
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	deliverable, _ := cmd.Flags().GetInt("deliverable")

	if !output.JSONOutput {
		fmt.Printf("Uploading %d files...\n", len(files))
	}

	results := uploadParallel(client, flightID, files, "uploads", concurrency, uploadOptions{
		productElementID: deliverable,
	})
	confirmed, failed := collectResults(results)

	if output.JSONOutput {
		output.JSON(confirmed)
	} else {
		fmt.Printf("%d/%d files uploaded.\n", len(confirmed), len(files))
		if len(failed) > 0 {
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		}
	}
	return nil
}

func runUploadImages(cmd *cobra.Command, args []string) error {
	client, err := api.NewAuthenticated()
	if err != nil {
		return err
	}

	flightID := args[0]
	concurrency, _ := cmd.Flags().GetInt("concurrency")

	// Resolve directories into image file lists
	files, err := resolveImageFiles(args[1:])
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no image files found")
	}

	if !output.JSONOutput {
		fmt.Printf("Found %d images. Uploading with %d workers...\n", len(files), concurrency)
	}

	results := uploadParallel(client, flightID, files, "images", concurrency, uploadOptions{
		withEXIF: true,
	})
	confirmed, failed := collectResults(results)

	if len(failed) > 0 && !output.JSONOutput {
		fmt.Printf("\nFailed (%d):\n", len(failed))
		for _, name := range failed {
			fmt.Printf("  - %s\n", name)
		}
	}

	if len(confirmed) == 0 {
		return fmt.Errorf("no images were uploaded successfully")
	}

	// Don't finalize if there were failures
	if len(failed) > 0 {
		if output.JSONOutput {
			output.JSON(map[string]any{
				"images": confirmed,
				"failed": failed,
			})
		} else {
			fmt.Printf("\n%d/%d images uploaded. Skipping finalize due to failures.\n", len(confirmed), len(files))
			fmt.Println("Fix the issues and re-run to retry.")
		}
		return nil
	}

	// Finalize the dataset
	finalizeBody, err := client.Post(fmt.Sprintf("/flights/%s/images/finalize", flightID), nil)
	if err != nil {
		return fmt.Errorf("failed to finalize dataset: %w", err)
	}

	var fin finalizeResponse
	json.Unmarshal(finalizeBody, &fin)

	if output.JSONOutput {
		output.JSON(map[string]any{
			"images":   confirmed,
			"finalize": fin,
		})
	} else {
		fmt.Printf("\n%d/%d images uploaded.\n", len(confirmed), len(files))
		fmt.Printf("Dataset finalized (%d pictures).\n", fin.PictureCount)
	}
	return nil
}

func uploadSingleFile(client *api.Client, flightID, filePath, endpoint string, opts uploadOptions) (*confirmResponse, error) {
	filename := filepath.Base(filePath)

	// Step 1: Presign
	presignBody, err := client.Post(
		fmt.Sprintf("/flights/%s/%s/presign", flightID, endpoint),
		map[string]string{"filename": filename},
	)
	if err != nil {
		return nil, fmt.Errorf("presign failed: %w", err)
	}

	var presign presignResponse
	if err := json.Unmarshal(presignBody, &presign); err != nil {
		return nil, fmt.Errorf("failed to parse presign response: %w", err)
	}

	// Step 2: Upload to S3 via multipart form POST
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if err := uploadToS3(filePath, presign.PresignData); err != nil {
		return nil, fmt.Errorf("S3 upload failed: %w", err)
	}

	// Step 3: Confirm
	confirmPayload := map[string]any{
		"filename": filename,
		"file_id":  presign.FileID,
		"size":     fileInfo.Size(),
	}

	if opts.withEXIF {
		if exifData := extractEXIF(filePath); exifData != nil {
			confirmPayload["exif"] = exifData
		}
	}
	if opts.productElementID != 0 {
		confirmPayload["product_element_id"] = opts.productElementID
	}

	confirmBody, err := client.Post(
		fmt.Sprintf("/flights/%s/%s/confirm", flightID, endpoint),
		confirmPayload,
	)
	if err != nil {
		return nil, fmt.Errorf("confirm failed: %w", err)
	}

	var result confirmResponse
	if err := json.Unmarshal(confirmBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse confirm response: %w", err)
	}

	return &result, nil
}

func uploadToS3(filePath string, presign presignData) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add presign fields first (must come before file)
	for key, val := range presign.Fields {
		writer.WriteField(key, val)
	}

	// Add the file
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", presign.URL, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("S3 returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
