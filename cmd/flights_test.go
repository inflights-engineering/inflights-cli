package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/output"
)

var sampleFlights = []map[string]any{
	{
		"id":             1,
		"public_uid":     "FL-001",
		"status":         "scheduled",
		"product":        "Aerial Survey",
		"scheduled_date": "2026-04-01",
		"area_in_ha":     12.5,
		"price_client":   750.0,
		"created_at":     "2026-03-10T10:00:00Z",
	},
	{
		"id":             2,
		"public_uid":     "FL-002",
		"status":         "completed",
		"product":        "Inspection",
		"scheduled_date": "2026-03-05",
		"area_in_ha":     3.2,
		"price_client":   400.0,
		"created_at":     "2026-03-01T08:00:00Z",
	},
}

var sampleFlightDetail = map[string]any{
	"id":               1,
	"public_uid":       "FL-001",
	"status":           "scheduled",
	"product":          "Aerial Survey",
	"scheduled_date":   "2026-04-01",
	"area_in_ha":       12.5,
	"price_client":     750.0,
	"created_at":       "2026-03-10T10:00:00Z",
	"flown_at":         nil,
	"completed_at":     nil,
	"cancelled_at":     nil,
	"description_user": "Map the construction site",
	"reference":        "REF-123",
	"pilot":            map[string]string{"id": "pilot-uuid", "name": "Jane Pilot"},
	"customer":         map[string]string{"id": "cust-uuid", "name": "Acme Corp"},
}

func TestFlights_List(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/flights" {
			t.Errorf("got path %q, want /flights", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleFlights)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runFlights(flightsCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "FL-001") {
		t.Errorf("output = %q, want it to contain 'FL-001'", out)
	}
	if !strings.Contains(out, "FL-002") {
		t.Errorf("output = %q, want it to contain 'FL-002'", out)
	}
}

func TestFlights_FilterByStatus(t *testing.T) {
	// Arrange
	var gotStatus string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotStatus = r.URL.Query().Get("status")
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	flightsCmd.Flags().Set("status", "scheduled")
	defer flightsCmd.Flags().Set("status", "")
	runFlights(flightsCmd, []string{})

	// Assert
	if gotStatus != "scheduled" {
		t.Errorf("got status param %q, want %q", gotStatus, "scheduled")
	}
}

func TestFlights_FilterByPublicUID(t *testing.T) {
	// Arrange
	var gotUID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUID = r.URL.Query().Get("public_uid")
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	flightsCmd.Flags().Set("public-uid", "FL-001")
	defer flightsCmd.Flags().Set("public-uid", "")
	runFlights(flightsCmd, []string{})

	// Assert
	if gotUID != "FL-001" {
		t.Errorf("got public_uid param %q, want %q", gotUID, "FL-001")
	}
}

func TestFlights_EmptyList(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runFlights(flightsCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "No flights found") {
		t.Errorf("output = %q, want it to contain 'No flights found'", out)
	}
}

func TestFlights_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleFlights)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runFlights(flightsCmd, []string{})
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

func TestFlights_NotLoggedIn(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	err := runFlights(flightsCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error when not logged in")
	}
}

func TestFlightShow(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/flights/1" {
			t.Errorf("got path %q, want /flights/1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleFlightDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runFlightShow(flightShowCmd, []string{"1"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	checks := []string{"FL-001", "scheduled", "Jane Pilot", "Acme Corp", "Map the construction site"}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestFlightShow_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleFlightDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runFlightShow(flightShowCmd, []string{"1"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["public_uid"] != "FL-001" {
		t.Errorf("got public_uid %v, want FL-001", parsed["public_uid"])
	}
}

func TestFlightShow_NotFound(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error":{"id":"not_found","message":"Flight not found"}}`))
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	err := runFlightShow(flightShowCmd, []string{"999"})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error for not found")
	}
	if !strings.Contains(err.Error(), "Flight not found") {
		t.Errorf("got error %q, want it to contain 'Flight not found'", err.Error())
	}
}
