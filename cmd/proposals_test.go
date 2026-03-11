package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/output"
)

var sampleProposals = []map[string]any{
	{
		"id":                "proposal-uuid-10",
		"status":            "pending",
		"flight_id":         1,
		"flight_public_uid": "FL-001",
		"scheduled_date":    "2026-04-01",
		"price_pilot":    "250.00",
		"created_at":     "2026-03-10T10:00:00Z",
	},
	{
		"id":                "proposal-uuid-11",
		"status":            "accepted",
		"flight_id":         2,
		"flight_public_uid": "FL-002",
		"scheduled_date":    "2026-04-05",
		"price_pilot":    "300.00",
		"created_at":     "2026-03-08T08:00:00Z",
	},
}

var sampleProposalDetail = map[string]any{
	"id":                    "proposal-uuid-10",
	"status":                "pending",
	"flight_id":             1,
	"flight_public_uid":     "FL-001",
	"scheduled_date":        "2026-04-01",
	"price_pilot":           "250.00",
	"created_at":            "2026-03-10T10:00:00Z",
	"backup_scheduled_date": "2026-04-03",
	"reason_for_rejection":  nil,
	"equipment_type": map[string]any{
		"id":           "equip-uuid-5",
		"brand":        "DJI",
		"product_name": "Mavic 3 Enterprise",
	},
	"flight": map[string]any{
		"id":               1,
		"public_uid":       "FL-001",
		"status":           "pilot_found",
		"product":          "Aerial Survey",
		"area_in_ha":       "12.5",
		"description_user": "Map the construction site",
	},
}

func TestProposals_List(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proposals" {
			t.Errorf("got path %q, want /proposals", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleProposals)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runProposalList(proposalListCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "pending") {
		t.Errorf("output = %q, want it to contain 'pending'", out)
	}
	if !strings.Contains(out, "accepted") {
		t.Errorf("output = %q, want it to contain 'accepted'", out)
	}
}

func TestProposals_FilterByStatus(t *testing.T) {
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
	proposalListCmd.Flags().Set("status", "pending")
	defer proposalListCmd.Flags().Set("status", "")
	runProposalList(proposalListCmd, []string{})

	// Assert
	if gotStatus != "pending" {
		t.Errorf("got status param %q, want %q", gotStatus, "pending")
	}
}

func TestProposals_EmptyList(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runProposalList(proposalListCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "No proposals found") {
		t.Errorf("output = %q, want it to contain 'No proposals found'", out)
	}
}

func TestProposals_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleProposals)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runProposalList(proposalListCmd, []string{})
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

func TestProposals_NotLoggedIn(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	err := runProposalList(proposalListCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error when not logged in")
	}
}

func TestProposalShow(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proposals/proposal-uuid-10" {
			t.Errorf("got path %q, want /proposals/10", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleProposalDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runProposalShow(proposalShowCmd, []string{"proposal-uuid-10"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	checks := []string{"pending", "250.00", "DJI", "Mavic 3 Enterprise", "FL-001", "Map the construction site"}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestProposalShow_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleProposalDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runProposalShow(proposalShowCmd, []string{"proposal-uuid-10"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestProposalShow_NotFound(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error":{"id":"not_found","message":"Proposal not found"}}`))
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	err := runProposalShow(proposalShowCmd, []string{"999"})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error for not found")
	}
}
