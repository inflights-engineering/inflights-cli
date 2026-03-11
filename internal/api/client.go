package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/inflights-engineering/inflights-cli/internal/config"
	"github.com/inflights-engineering/inflights-cli/internal/credentials"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type APIError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

// NewUnauthenticated creates a client without auth (for login flow).
func NewUnauthenticated() *Client {
	return &Client{
		BaseURL:    config.BaseURL(),
		HTTPClient: &http.Client{},
	}
}

// NewAuthenticated creates a client using stored credentials.
func NewAuthenticated() (*Client, error) {
	creds, err := credentials.Load()
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:    config.BaseURL(),
		Token:      creds.Token,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.do("GET", path, nil)
}

func (c *Client) Post(path string, body any) ([]byte, error) {
	return c.do("POST", path, body)
}

func (c *Client) Patch(path string, body any) ([]byte, error) {
	return c.do("PATCH", path, body)
}

func (c *Client) Delete(path string) ([]byte, error) {
	return c.do("DELETE", path, nil)
}

func (c *Client) do(method, path string, body any) ([]byte, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error.Message)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
