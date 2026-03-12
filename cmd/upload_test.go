package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/output"
)

func writeTestFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func makeUploadServer(t *testing.T, s3Server *httptest.Server) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/presign"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"file_id":      "abc123.tif",
				"presign_data": map[string]any{"url": s3Server.URL, "fields": map[string]string{"key": "cache/abc123.tif"}},
			})
		case strings.HasSuffix(r.URL.Path, "/confirm"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{"id": "doc-uuid-1", "filename": "test.tif", "size": 12})
		case strings.HasSuffix(r.URL.Path, "/finalize"):
			json.NewEncoder(w).Encode(map[string]any{"dataset_id": "ds-uuid-1", "dataset_status": "uploaded", "picture_count": 1})
		}
	}))
}

func TestUploadData(t *testing.T) {
	// Arrange
	var presignCalled, confirmCalled bool
	var confirmBody map[string]any

	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer s3Server.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/uploads/presign"):
			presignCalled = true
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"file_id": "abc123.tif",
				"presign_data": map[string]any{
					"url":    s3Server.URL,
					"fields": map[string]string{"key": "cache/abc123.tif"},
				},
			})
		case strings.HasSuffix(r.URL.Path, "/uploads/confirm"):
			confirmCalled = true
			json.NewDecoder(r.Body).Decode(&confirmBody)
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"id":       "doc-uuid-1",
				"filename": "test.tif",
				"size":     12,
			})
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	filePath := writeTestFile(t, "test.tif", "file-content")

	// Act
	out := captureOutput(t, func() {
		err := runUploadData(uploadDataCmd, []string{"42", filePath})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !presignCalled {
		t.Error("presign endpoint was not called")
	}
	if !confirmCalled {
		t.Error("confirm endpoint was not called")
	}
	if confirmBody["file_id"] != "abc123.tif" {
		t.Errorf("got file_id %v, want abc123.tif", confirmBody["file_id"])
	}
	if !strings.Contains(out, "1/1 files uploaded") {
		t.Errorf("output = %q, want it to contain '1/1 files uploaded'", out)
	}
}

func TestUploadData_JSONOutput(t *testing.T) {
	// Arrange
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer s3Server.Close()

	server := makeUploadServer(t, s3Server)
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	filePath := writeTestFile(t, "test.tif", "file-content")

	// Act
	out := captureOutput(t, func() {
		err := runUploadData(uploadDataCmd, []string{"42", filePath})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed []map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(parsed) != 1 {
		t.Errorf("got %d items, want 1", len(parsed))
	}
}

func TestUploadImages(t *testing.T) {
	// Arrange
	var finalizeCalled bool

	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer s3Server.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/images/presign"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"file_id":      "img123.jpg",
				"presign_data": map[string]any{"url": s3Server.URL, "fields": map[string]string{}},
			})
		case strings.HasSuffix(r.URL.Path, "/images/confirm"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{"id": "pic-uuid-1", "filename": "photo.jpg", "size": 1024})
		case strings.HasSuffix(r.URL.Path, "/images/finalize"):
			finalizeCalled = true
			json.NewEncoder(w).Encode(map[string]any{
				"dataset_id":     "ds-uuid-1",
				"dataset_status": "uploaded",
				"picture_count":  1,
			})
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	filePath := writeTestFile(t, "photo.jpg", "jpeg-content")

	// Act
	out := captureOutput(t, func() {
		err := runUploadImages(uploadImagesCmd, []string{"42", filePath})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !finalizeCalled {
		t.Error("finalize endpoint was not called")
	}
	if !strings.Contains(out, "1/1 images uploaded") {
		t.Errorf("output = %q, want it to contain '1/1 images uploaded'", out)
	}
	if !strings.Contains(out, "finalized") {
		t.Errorf("output = %q, want it to contain 'finalized'", out)
	}
}

func TestUploadImages_FolderScanning(t *testing.T) {
	// Arrange — create a directory with mixed files
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "img1.jpg"), []byte("jpg"), 0644)
	os.WriteFile(filepath.Join(dir, "img2.jpeg"), []byte("jpeg"), 0644)
	os.WriteFile(filepath.Join(dir, "img3.png"), []byte("png"), 0644)
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("text"), 0644)
	os.WriteFile(filepath.Join(dir, "data.csv"), []byte("csv"), 0644)

	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer s3Server.Close()

	var confirmCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/presign"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"file_id":      "abc.jpg",
				"presign_data": map[string]any{"url": s3Server.URL, "fields": map[string]string{}},
			})
		case strings.HasSuffix(r.URL.Path, "/confirm"):
			confirmCount++
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{"id": "pic-uuid", "filename": "img.jpg", "size": 3})
		case strings.HasSuffix(r.URL.Path, "/finalize"):
			json.NewEncoder(w).Encode(map[string]any{"dataset_id": "ds-uuid", "dataset_status": "uploaded", "picture_count": 3})
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runUploadImages(uploadImagesCmd, []string{"42", dir})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert — only 3 image files should be uploaded, not txt/csv
	if confirmCount != 3 {
		t.Errorf("got %d confirms, want 3 (should skip non-image files)", confirmCount)
	}
	if !strings.Contains(out, "Found 3 images") {
		t.Errorf("output = %q, want it to contain 'Found 3 images'", out)
	}
}

func TestUploadImages_PartialFailure_NoFinalize(t *testing.T) {
	// Arrange — S3 fails for every other file
	callCount := 0
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount%2 == 0 {
			w.WriteHeader(403)
			w.Write([]byte("Access Denied"))
			return
		}
		w.WriteHeader(204)
	}))
	defer s3Server.Close()

	var finalizeCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/presign"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"file_id":      "abc.jpg",
				"presign_data": map[string]any{"url": s3Server.URL, "fields": map[string]string{}},
			})
		case strings.HasSuffix(r.URL.Path, "/confirm"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{"id": "pic-uuid", "filename": "img.jpg", "size": 3})
		case strings.HasSuffix(r.URL.Path, "/finalize"):
			finalizeCalled = true
			json.NewEncoder(w).Encode(map[string]any{"dataset_id": "ds-uuid", "dataset_status": "uploaded", "picture_count": 1})
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	file1 := writeTestFile(t, "img1.jpg", "jpg1")
	file2 := writeTestFile(t, "img2.jpg", "jpg2")

	// Act
	out := captureOutput(t, func() {
		runUploadImages(uploadImagesCmd, []string{"42", file1, file2})
	})

	// Assert — should NOT finalize when there are failures
	if finalizeCalled {
		t.Error("finalize was called despite partial failure")
	}
	if !strings.Contains(out, "Skipping finalize") {
		t.Errorf("output = %q, want it to contain 'Skipping finalize'", out)
	}
}

func TestUploadData_S3Failure(t *testing.T) {
	// Arrange — S3 returns error
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte("Access Denied"))
	}))
	defer s3Server.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{
			"file_id":      "abc123.tif",
			"presign_data": map[string]any{"url": s3Server.URL, "fields": map[string]string{}},
		})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	filePath := writeTestFile(t, "test.tif", "file-content")

	// Act
	out := captureOutput(t, func() {
		runUploadData(uploadDataCmd, []string{"42", filePath})
	})

	// Assert — should report error but not crash
	if !strings.Contains(out, "failed") {
		t.Errorf("output = %q, want it to contain 'failed'", out)
	}
}

func TestUploadData_FileNotFound(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{
			"file_id":      "abc123.tif",
			"presign_data": map[string]any{"url": "http://localhost", "fields": map[string]string{}},
		})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		runUploadData(uploadDataCmd, []string{"42", "/nonexistent/file.tif"})
	})

	// Assert
	if !strings.Contains(out, "failed") {
		t.Errorf("output = %q, want it to contain 'failed'", out)
	}
}

// Verify S3 receives the file content
func TestUploadData_S3ReceivesFile(t *testing.T) {
	// Arrange
	var s3Body []byte
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s3Body, _ = io.ReadAll(r.Body)
		w.WriteHeader(204)
	}))
	defer s3Server.Close()

	server := makeUploadServer(t, s3Server)
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	filePath := writeTestFile(t, "test.tif", "my-file-data")

	// Act
	captureOutput(t, func() {
		runUploadData(uploadDataCmd, []string{"42", filePath})
	})

	// Assert — S3 should have received the file content in multipart body
	if !strings.Contains(string(s3Body), "my-file-data") {
		t.Error("S3 did not receive the file content")
	}
	if !strings.Contains(string(s3Body), "cache/abc123.tif") {
		t.Error("S3 did not receive the presign key field")
	}
}

func TestResolveImageFiles(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.jpg"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "b.PNG"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "c.tiff"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "d.txt"), []byte(""), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)

	// Act
	files, err := resolveImageFiles([]string{dir})
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(files) != 3 {
		t.Errorf("got %d files, want 3", len(files))
	}
}

func TestResolveImageFiles_MixedArgs(t *testing.T) {
	// Arrange — one explicit file + one directory
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.jpg"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "b.jpg"), []byte(""), 0644)
	explicit := writeTestFile(t, "explicit.dng", "")

	// Act
	files, err := resolveImageFiles([]string{explicit, dir})
	if err != nil {
		t.Fatal(err)
	}

	// Assert — 1 explicit + 2 from directory
	if len(files) != 3 {
		t.Errorf("got %d files, want 3", len(files))
	}
}
