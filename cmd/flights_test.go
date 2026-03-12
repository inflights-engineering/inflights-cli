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

var sampleGeoJSON = `{
	"type": "Feature",
	"geometry": {
		"type": "GeometryCollection",
		"geometries": [{
			"type": "Polygon",
			"coordinates": [[[2.35, 48.85], [2.36, 48.85], [2.36, 48.86], [2.35, 48.86], [2.35, 48.85]]]
		}]
	},
	"properties": {}
}`

func writeGeoJSON(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "area.geojson")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFlightOrder(t *testing.T) {
	// Arrange
	var gotPath, gotMethod string
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleFlightDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	geojsonPath := writeGeoJSON(t, sampleGeoJSON)

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")
	flightOrderCmd.Flags().Set("description", "Test flight")
	defer flightOrderCmd.Flags().Set("description", "")

	// Act
	out := captureOutput(t, func() {
		err := runFlightOrder(flightOrderCmd, []string{geojsonPath})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	if gotMethod != "POST" {
		t.Errorf("got method %q, want POST", gotMethod)
	}
	if gotPath != "/flights" {
		t.Errorf("got path %q, want /flights", gotPath)
	}
	if gotBody["product_id"] != float64(3) {
		t.Errorf("got product_id %v, want 3", gotBody["product_id"])
	}
	if gotBody["description_user"] != "Test flight" {
		t.Errorf("got description_user %v, want 'Test flight'", gotBody["description_user"])
	}
	if gotBody["areas"] == nil {
		t.Error("got nil areas, want GeoJSON")
	}
	if !strings.Contains(out, "FL-001") {
		t.Errorf("output = %q, want it to contain 'FL-001'", out)
	}
	if !strings.Contains(out, "created") {
		t.Errorf("output = %q, want it to contain 'created'", out)
	}
}

func TestFlightOrder_JSONOutput(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleFlightDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")
	output.JSONOutput = true

	geojsonPath := writeGeoJSON(t, sampleGeoJSON)

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")

	// Act
	out := captureOutput(t, func() {
		err := runFlightOrder(flightOrderCmd, []string{geojsonPath})
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

func TestFlightOrder_InvalidFile(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")

	// Act
	err := runFlightOrder(flightOrderCmd, []string{"/nonexistent/file.geojson"})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error for missing file")
	}
	if !strings.Contains(err.Error(), "failed to read GeoJSON") {
		t.Errorf("got error %q, want it to contain 'failed to read GeoJSON'", err.Error())
	}
}

func TestFlightOrder_InvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	badPath := writeGeoJSON(t, "not json at all")

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")

	// Act
	err := runFlightOrder(flightOrderCmd, []string{badPath})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("got error %q, want it to contain 'invalid JSON'", err.Error())
	}
}

func TestFlightOrder_NormalizesPlainPolygon(t *testing.T) {
	// Arrange — bare Polygon, not wrapped in Feature > GeometryCollection
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleFlightDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	plainPolygon := `{
		"type": "Polygon",
		"coordinates": [[[-3.70, 40.41], [-3.69, 40.41], [-3.69, 40.42], [-3.70, 40.42], [-3.70, 40.41]]]
	}`
	geojsonPath := writeGeoJSON(t, plainPolygon)

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")

	// Act
	captureOutput(t, func() {
		err := runFlightOrder(flightOrderCmd, []string{geojsonPath})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert — should be wrapped into Feature > GeometryCollection > [Polygon]
	areas, _ := gotBody["areas"].(map[string]any)
	if areas["type"] != "Feature" {
		t.Errorf("got areas.type %v, want Feature", areas["type"])
	}
	geom, _ := areas["geometry"].(map[string]any)
	if geom["type"] != "GeometryCollection" {
		t.Errorf("got geometry.type %v, want GeometryCollection", geom["type"])
	}
	geometries, _ := geom["geometries"].([]any)
	if len(geometries) != 1 {
		t.Errorf("got %d geometries, want 1", len(geometries))
	}
}

func TestFlightOrder_NormalizesFeatureWithPolygon(t *testing.T) {
	// Arrange — Feature with Polygon (standard GeoJSON from most tools)
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleFlightDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	featurePolygon := `{
		"type": "Feature",
		"geometry": {
			"type": "Polygon",
			"coordinates": [[[-3.70, 40.41], [-3.69, 40.41], [-3.69, 40.42], [-3.70, 40.42], [-3.70, 40.41]]]
		},
		"properties": {}
	}`
	geojsonPath := writeGeoJSON(t, featurePolygon)

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")

	// Act
	captureOutput(t, func() {
		err := runFlightOrder(flightOrderCmd, []string{geojsonPath})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert
	areas, _ := gotBody["areas"].(map[string]any)
	geom, _ := areas["geometry"].(map[string]any)
	if geom["type"] != "GeometryCollection" {
		t.Errorf("got geometry.type %v, want GeometryCollection", geom["type"])
	}
}

func TestFlightOrder_NormalizesFeatureCollection(t *testing.T) {
	// Arrange — FeatureCollection with two polygon features
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleFlightDetail)
	}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	fc := `{
		"type": "FeatureCollection",
		"features": [
			{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-3.70,40.41],[-3.69,40.41],[-3.69,40.42],[-3.70,40.42],[-3.70,40.41]]]},"properties":{}},
			{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-3.68,40.41],[-3.67,40.41],[-3.67,40.42],[-3.68,40.42],[-3.68,40.41]]]},"properties":{}}
		]
	}`
	geojsonPath := writeGeoJSON(t, fc)

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")

	// Act
	captureOutput(t, func() {
		err := runFlightOrder(flightOrderCmd, []string{geojsonPath})
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	})

	// Assert — both polygons should be merged into one GeometryCollection
	areas, _ := gotBody["areas"].(map[string]any)
	geom, _ := areas["geometry"].(map[string]any)
	geometries, _ := geom["geometries"].([]any)
	if len(geometries) != 2 {
		t.Errorf("got %d geometries, want 2", len(geometries))
	}
}

func TestFlightOrder_UnsupportedType(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	setupTestEnv(t, server)
	saveTestCredentials(t, "test-token")

	point := `{"type": "Point", "coordinates": [-3.70, 40.41]}`
	geojsonPath := writeGeoJSON(t, point)

	flightOrderCmd.Flags().Set("service", "3")
	defer flightOrderCmd.Flags().Set("service", "0")

	// Act
	err := runFlightOrder(flightOrderCmd, []string{geojsonPath})

	// Assert
	if err == nil {
		t.Fatal("got nil error, want error for unsupported type")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("got error %q, want it to contain 'unsupported'", err.Error())
	}
}
