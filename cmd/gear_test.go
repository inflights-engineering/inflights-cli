package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/output"
)

var sampleEquipmentTypes = []map[string]any{
	{
		"id":           "et-uuid-1",
		"brand":        "DJI",
		"product_name": "Mavic 3 Enterprise",
		"category":     "drone",
		"sensor_types": []string{"rgb"},
		"resolution":   20,
	},
	{
		"id":           "et-uuid-2",
		"brand":        "Leica",
		"product_name": "BLK360",
		"category":     "payload",
		"sensor_types": []string{"lidar"},
		"resolution":   nil,
	},
}

var sampleEquipments = []map[string]any{
	{
		"id": "eq-uuid-1",
		"equipment_type": map[string]any{
			"id":           "et-uuid-1",
			"brand":        "DJI",
			"product_name": "Mavic 3 Enterprise",
			"category":     "drone",
		},
		"full_day_rate":   "500.00",
		"price_minimum":   "200.00",
		"price_per_ha":    "15.00",
		"surface_minimum": "5.0",
	},
}

func TestGearList(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/equipment_types" {
			t.Errorf("got path %q, want /equipment_types", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleEquipmentTypes)
	}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	out := captureOutput(t, func() {
		err := runGearList(gearListCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "Mavic 3 Enterprise") {
		t.Errorf("output = %q, want it to contain 'Mavic 3 Enterprise'", out)
	}
	if !strings.Contains(out, "BLK360") {
		t.Errorf("output = %q, want it to contain 'BLK360'", out)
	}
}

func TestGearList_FilterByCategory(t *testing.T) {
	// Arrange
	var gotCategory string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCategory = r.URL.Query().Get("category")
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	gearListCmd.Flags().Set("category", "drone")
	defer gearListCmd.Flags().Set("category", "")
	runGearList(gearListCmd, []string{})

	// Assert
	if gotCategory != "drone" {
		t.Errorf("got category param %q, want %q", gotCategory, "drone")
	}
}

func TestGearList_EmptyList(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	out := captureOutput(t, func() {
		err := runGearList(gearListCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "No equipment types found") {
		t.Errorf("output = %q, want it to contain 'No equipment types found'", out)
	}
}

func TestGearList_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleEquipmentTypes)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runGearList(gearListCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed []map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(parsed) != 2 {
		t.Errorf("got %d items, want 2", len(parsed))
	}
}

func TestGearMine(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/equipments" {
			t.Errorf("got path %q, want /equipments", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleEquipments)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runGearMine(gearMineCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "DJI") {
		t.Errorf("output = %q, want it to contain 'DJI'", out)
	}
	if !strings.Contains(out, "15.00") {
		t.Errorf("output = %q, want it to contain '15.00'", out)
	}
}

func TestGearMine_EmptyList(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runGearMine(gearMineCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "No equipment found") {
		t.Errorf("output = %q, want it to contain 'No equipment found'", out)
	}
}

func TestGearMine_NotLoggedIn(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	err := runGearMine(gearMineCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error when not logged in")
	}
}
