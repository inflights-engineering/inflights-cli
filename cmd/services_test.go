package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/output"
)

func TestServices_ListsServices(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services" {
			t.Errorf("got path %q, want /services", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("got method %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":            1,
				"name":          "Aerial Survey",
				"subtitle":      "High-res mapping",
				"product_type":  "survey",
				"price_minimum": 500,
				"industries":    []string{"construction"},
				"sensor_types":  []string{"lidar"},
			},
			{
				"id":            2,
				"name":          "Inspection",
				"subtitle":      "Asset inspection",
				"product_type":  "inspection",
				"price_minimum": 300,
				"industries":    []string{"energy"},
				"sensor_types":  []string{"rgb"},
			},
		})
	}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	out := captureOutput(t, func() {
		err := runServices(servicesCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "Aerial Survey") {
		t.Errorf("output = %q, want it to contain 'Aerial Survey'", out)
	}
	if !strings.Contains(out, "Inspection") {
		t.Errorf("output = %q, want it to contain 'Inspection'", out)
	}
}

func TestServices_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Aerial Survey", "product_type": "survey", "price_minimum": 500},
		})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runServices(servicesCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert — output should be valid JSON
	var parsed []map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, out)
	}
	if len(parsed) != 1 {
		t.Errorf("got %d items, want 1", len(parsed))
	}
}

func TestServices_EmptyList(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	err := runServices(servicesCmd, []string{})

	// Assert
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
}

func TestServices_APIError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"error":{"id":"server_error","message":"Something broke"}}`))
	}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	err := runServices(servicesCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error")
	}
}
