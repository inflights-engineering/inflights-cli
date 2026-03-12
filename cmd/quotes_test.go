package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inflights-engineering/inflights-cli/internal/output"
)

var sampleQuotes = []map[string]any{
	{
		"id":           1,
		"quote_number": "Q-2026-001",
		"status":       "pending",
		"amount":       "750.00",
		"vat_percent":  "21",
		"quote_date":   "2026-03-10",
		"due_date":     "2026-04-10",
		"created_at":   "2026-03-10T10:00:00Z",
		"type":         "quote",
	},
	{
		"type":              "estimate",
		"flight_id":         2,
		"flight_public_uid": "FL-002",
		"status":            "pending",
		"amount":            "400.00",
		"product":           "Inspection",
		"created_at":        "2026-03-01T08:00:00Z",
	},
}

var sampleQuoteDetail = map[string]any{
	"id":           1,
	"quote_number": "Q-2026-001",
	"status":       "accepted",
	"amount":       "750.00",
	"vat_percent":  "21",
	"quote_date":   "2026-03-10",
	"due_date":     "2026-04-10",
	"created_at":   "2026-03-10T10:00:00Z",
	"type":         "quote",
	"accepted_at":  "2026-03-12T14:00:00Z",
	"flights": []map[string]any{
		{"id": 1, "public_uid": "FL-001", "status": "flight_scheduled", "product": "Aerial Survey"},
	},
}

func TestQuotes_List(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/quotes" {
			t.Errorf("got path %q, want /quotes", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleQuotes)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runQuotes(quotesCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "Q-2026-001") {
		t.Errorf("output = %q, want it to contain 'Q-2026-001'", out)
	}
	if !strings.Contains(out, "FL-002 (estimate)") {
		t.Errorf("output = %q, want it to contain 'FL-002 (estimate)'", out)
	}
}

func TestQuotes_FilterByStatus(t *testing.T) {
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
	quotesCmd.Flags().Set("status", "accepted")
	defer quotesCmd.Flags().Set("status", "")
	runQuotes(quotesCmd, []string{})

	// Assert
	if gotStatus != "accepted" {
		t.Errorf("got status param %q, want %q", gotStatus, "accepted")
	}
}

func TestQuotes_EmptyList(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runQuotes(quotesCmd, []string{})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if !strings.Contains(out, "No quotes found") {
		t.Errorf("output = %q, want it to contain 'No quotes found'", out)
	}
}

func TestQuotes_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleQuotes)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runQuotes(quotesCmd, []string{})
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

func TestQuotes_NotLoggedIn(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)

	// Act
	err := runQuotes(quotesCmd, []string{})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error when not logged in")
	}
}

func TestQuoteShow(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/quotes/1" {
			t.Errorf("got path %q, want /quotes/1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sampleQuoteDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runQuoteShow(quoteShowCmd, []string{"1"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	checks := []string{"Q-2026-001", "accepted", "750.00", "FL-001", "Aerial Survey"}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestQuoteShow_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sampleQuoteDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runQuoteShow(quoteShowCmd, []string{"1"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["quote_number"] != "Q-2026-001" {
		t.Errorf("got quote_number %v, want Q-2026-001", parsed["quote_number"])
	}
}

func TestQuoteConfirm_Quote(t *testing.T) {
	// Arrange
	var lastPath, lastMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastPath = r.URL.Path
		lastMethod = r.Method
		switch r.URL.Path {
		case "/quotes":
			json.NewEncoder(w).Encode([]map[string]any{
				{"id": 1, "quote_number": "Q-2026-001", "status": "pending", "type": "quote", "amount": "750"},
			})
		case "/quotes/1/accept":
			json.NewEncoder(w).Encode(sampleQuoteDetail)
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runQuoteConfirm(quoteConfirmCmd, []string{"Q-2026-001"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if lastMethod != "POST" {
		t.Errorf("got method %q, want POST", lastMethod)
	}
	if lastPath != "/quotes/1/accept" {
		t.Errorf("got path %q, want /quotes/1/accept", lastPath)
	}
	if !strings.Contains(out, "Q-2026-001 accepted") {
		t.Errorf("output = %q, want it to contain 'Q-2026-001 accepted'", out)
	}
}

func TestQuoteConfirm_Estimate(t *testing.T) {
	// Arrange
	var lastPath string
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastPath = r.URL.Path
		switch r.URL.Path {
		case "/quotes":
			json.NewEncoder(w).Encode([]map[string]any{
				{"type": "estimate", "flight_id": 60, "flight_public_uid": "BT779", "status": "pending", "amount": "500"},
			})
		case "/quotes/accept_estimate":
			json.NewDecoder(r.Body).Decode(&gotBody)
			json.NewEncoder(w).Encode(map[string]any{"type": "estimate", "flight_id": 60, "status": "accepted"})
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runQuoteConfirm(quoteConfirmCmd, []string{"BT779"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if lastPath != "/quotes/accept_estimate" {
		t.Errorf("got path %q, want /quotes/accept_estimate", lastPath)
	}
	if gotBody["flight_id"] != float64(60) {
		t.Errorf("got flight_id %v, want 60", gotBody["flight_id"])
	}
	if !strings.Contains(out, "BT779 accepted") {
		t.Errorf("output = %q, want it to contain 'BT779 accepted'", out)
	}
}

func TestQuoteConfirm_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/quotes":
			json.NewEncoder(w).Encode([]map[string]any{
				{"id": 1, "quote_number": "Q-2026-001", "status": "pending", "type": "quote", "amount": "750"},
			})
		case "/quotes/1/accept":
			json.NewEncoder(w).Encode(sampleQuoteDetail)
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runQuoteConfirm(quoteConfirmCmd, []string{"Q-2026-001"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["quote_number"] != "Q-2026-001" {
		t.Errorf("got quote_number %v, want Q-2026-001", parsed["quote_number"])
	}
}

func TestQuoteConfirm_NotFound(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	err := runQuoteConfirm(quoteConfirmCmd, []string{"NONEXISTENT"})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error")
	}
	if !strings.Contains(err.Error(), "no quote or estimate found") {
		t.Errorf("got error %q, want it to contain 'no quote or estimate found'", err.Error())
	}
}

func TestQuoteReject_Quote(t *testing.T) {
	// Arrange
	var lastPath, lastMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastPath = r.URL.Path
		lastMethod = r.Method
		switch r.URL.Path {
		case "/quotes":
			json.NewEncoder(w).Encode([]map[string]any{
				{"id": 1, "quote_number": "Q-2026-001", "status": "pending", "type": "quote", "amount": "750"},
			})
		case "/quotes/1/reject":
			json.NewEncoder(w).Encode(sampleQuoteDetail)
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runQuoteReject(quoteRejectCmd, []string{"Q-2026-001"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if lastMethod != "POST" {
		t.Errorf("got method %q, want POST", lastMethod)
	}
	if lastPath != "/quotes/1/reject" {
		t.Errorf("got path %q, want /quotes/1/reject", lastPath)
	}
	if !strings.Contains(out, "Q-2026-001 rejected") {
		t.Errorf("output = %q, want it to contain 'Q-2026-001 rejected'", out)
	}
}

func TestQuoteReject_Estimate(t *testing.T) {
	// Arrange
	var lastPath string
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastPath = r.URL.Path
		switch r.URL.Path {
		case "/quotes":
			json.NewEncoder(w).Encode([]map[string]any{
				{"type": "estimate", "flight_id": 60, "flight_public_uid": "PC376", "status": "pending", "amount": "500"},
			})
		case "/quotes/reject_estimate":
			json.NewDecoder(r.Body).Decode(&gotBody)
			json.NewEncoder(w).Encode(map[string]any{"type": "estimate", "flight_id": 60, "status": "rejected"})
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	out := captureOutput(t, func() {
		err := runQuoteReject(quoteRejectCmd, []string{"PC376"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if lastPath != "/quotes/reject_estimate" {
		t.Errorf("got path %q, want /quotes/reject_estimate", lastPath)
	}
	if gotBody["flight_id"] != float64(60) {
		t.Errorf("got flight_id %v, want 60", gotBody["flight_id"])
	}
	if !strings.Contains(out, "PC376 rejected") {
		t.Errorf("output = %q, want it to contain 'PC376 rejected'", out)
	}
}

func TestQuoteReject_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/quotes":
			json.NewEncoder(w).Encode([]map[string]any{
				{"id": 1, "quote_number": "Q-2026-001", "status": "pending", "type": "quote", "amount": "750"},
			})
		case "/quotes/1/reject":
			json.NewEncoder(w).Encode(map[string]any{"id": 1, "quote_number": "Q-2026-001", "status": "rejected"})
		}
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	// Act
	out := captureOutput(t, func() {
		err := runQuoteReject(quoteRejectCmd, []string{"Q-2026-001"})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["status"] != "rejected" {
		t.Errorf("got status %v, want rejected", parsed["status"])
	}
}

func TestQuoteReject_NotFound(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	err := runQuoteReject(quoteRejectCmd, []string{"NONEXISTENT"})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error")
	}
	if !strings.Contains(err.Error(), "no quote or estimate found") {
		t.Errorf("got error %q, want it to contain 'no quote or estimate found'", err.Error())
	}
}

func TestQuoteShow_NotFound(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error":{"id":"not_found","message":"Quote not found"}}`))
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	// Act
	err := runQuoteShow(quoteShowCmd, []string{"999"})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error for not found")
	}
}
