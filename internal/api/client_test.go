package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestClient creates a Client pointing at a local httptest server.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := &Client{
		BaseURL:    server.URL,
		HTTPClient: &http.Client{},
	}
	return client, server
}

func TestHTTPMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(c *Client) ([]byte, error)
	}{
		{
			name:   "GET request",
			method: "GET",
			call:   func(c *Client) ([]byte, error) { return c.Get("/test") },
		},
		{
			name:   "POST request",
			method: "POST",
			call:   func(c *Client) ([]byte, error) { return c.Post("/test", nil) },
		},
		{
			name:   "PATCH request",
			method: "PATCH",
			call:   func(c *Client) ([]byte, error) { return c.Patch("/test", nil) },
		},
		{
			name:   "DELETE request",
			method: "DELETE",
			call:   func(c *Client) ([]byte, error) { return c.Delete("/test") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var gotMethod string
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				w.Write([]byte(`{}`))
			})

			// Act
			_, err := tt.call(client)

			// Assert
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if gotMethod != tt.method {
				t.Errorf("got method %q, want %q", gotMethod, tt.method)
			}
		})
	}
}

func TestPost_SendsJSONBody(t *testing.T) {
	// Arrange
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]string
		json.Unmarshal(body, &payload)
		if payload["key"] != "value" {
			t.Errorf("got payload key %q, want %q", payload["key"], "value")
		}
		w.Write([]byte(`{"created":true}`))
	})

	// Act
	body, err := client.Post("/create", map[string]string{"key": "value"})

	// Assert
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if !strings.Contains(string(body), "created") {
		t.Errorf("got body %q, want it to contain 'created'", string(body))
	}
}

func TestAuthorizationHeader(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantHeader string
	}{
		{
			name:       "sends Bearer token when set",
			token:      "my-secret-token",
			wantHeader: "Bearer my-secret-token",
		},
		{
			name:       "sends no header when token is empty",
			token:      "",
			wantHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var gotHeader string
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				gotHeader = r.Header.Get("Authorization")
				w.Write([]byte(`{}`))
			})
			client.Token = tt.token

			// Act
			_, err := client.Get("/test")

			// Assert
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if gotHeader != tt.wantHeader {
				t.Errorf("got Authorization %q, want %q", gotHeader, tt.wantHeader)
			}
		})
	}
}

func TestRequestHeaders(t *testing.T) {
	// Arrange
	var gotContentType, gotAccept string
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		gotAccept = r.Header.Get("Accept")
		w.Write([]byte(`{}`))
	})

	// Act
	client.Get("/test")

	// Assert
	if gotContentType != "application/json" {
		t.Errorf("got Content-Type %q, want %q", gotContentType, "application/json")
	}
	if gotAccept != "application/json" {
		t.Errorf("got Accept %q, want %q", gotAccept, "application/json")
	}
}

func TestAPIErrors(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		wantSubstr string
	}{
		{
			name:       "parses structured error response",
			status:     404,
			body:       `{"error":{"id":"not_found","message":"Record not found"}}`,
			wantSubstr: "Record not found",
		},
		{
			name:       "falls back to raw body for unstructured errors",
			status:     500,
			body:       `Internal Server Error`,
			wantSubstr: "Internal Server Error",
		},
		{
			name:       "includes status code in error",
			status:     422,
			body:       `{"error":{"id":"validation_error","message":"Invalid input"}}`,
			wantSubstr: "422",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			})

			// Act
			_, err := client.Get("/error")

			// Assert
			if err == nil {
				t.Fatal("got nil error, want error")
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Errorf("got error %q, want it to contain %q", err.Error(), tt.wantSubstr)
			}
		})
	}
}

func TestSuccessStatusCodes(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{name: "200 OK", status: 200},
		{name: "201 Created", status: 201},
		{name: "202 Accepted", status: 202},
		{name: "204 No Content", status: 204},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			})

			// Act
			_, err := client.Get("/ok")

			// Assert
			if err != nil {
				t.Errorf("got error %v, want nil for status %d", err, tt.status)
			}
		})
	}
}
