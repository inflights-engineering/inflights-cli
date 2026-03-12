package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/output"
)

var sampleDownloads = map[string]any{
	"documents": []map[string]any{},
	"picture_set": map[string]any{
		"id":           "ps-uuid-1",
		"download_url": "",
	},
}

func TestDownload(t *testing.T) {
	// Arrange — serve download list + a file
	fileServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("file-content"))
	}))
	defer fileServer.Close()

	downloads := map[string]any{
		"documents": []map[string]any{
			{
				"id":           "doc-uuid-1",
				"title":        "deliverable.tif",
				"filename":     "deliverable.tif",
				"type":         "processed_deliverable",
				"download_url": fileServer.URL + "/deliverable.tif",
			},
		},
		"picture_set": nil,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(downloads)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	outDir := t.TempDir()
	downloadCmd.Flags().Set("output", outDir)
	defer downloadCmd.Flags().Set("output", ".")

	// Act
	out := captureOutput(t, func() {
		err := runDownload(downloadCmd, []string{"1"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "deliverable.tif") {
		t.Errorf("output = %q, want it to contain 'deliverable.tif'", out)
	}
	content, err := os.ReadFile(filepath.Join(outDir, "deliverable.tif"))
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if string(content) != "file-content" {
		t.Errorf("got file content %q, want 'file-content'", string(content))
	}
}

func TestDownload_Empty(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"documents":   []map[string]any{},
			"picture_set": nil,
		})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runDownload(downloadCmd, []string{"1"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "No downloads available") {
		t.Errorf("output = %q, want it to contain 'No downloads available'", out)
	}
}

func TestDownload_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleDownloads)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runDownload(downloadCmd, []string{"1"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["documents"] == nil {
		t.Error("got nil documents, want array")
	}
}
