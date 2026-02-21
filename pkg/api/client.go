// Package api provides a client for the API‑Football service.
package api

import (
    "context"
    "encoding/json"

    "fmt"
    "io"
    "log"
    "net/http"
    "time"
)

const BaseURL = "https://v3.football.api-sports.io"

// APIError represents an error returned by the API.
type APIError struct {
    StatusCode int
    Message    string
    Body       []byte
}

func (e *APIError) Error() string { return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message) }

// Client is the API‑Football client.
// It holds the API key and an injectable http.Client.
type Client struct {
    ApiKey string
    // HttpClient is used for all requests. It can be overridden for testing.
    HttpClient *http.Client
}

// NewClient creates a new Client with the provided API key.
// The underlying http.Client has a default timeout of 10 seconds.
func NewClient(apiKey string) *Client {
    return &Client{ApiKey: apiKey, HttpClient: &http.Client{Timeout: 10 * time.Second}}
}

// WithHTTPClient allows injecting a custom http.Client (e.g., for tests).
func (c *Client) WithHTTPClient(client *http.Client) *Client {
    c.HttpClient = client
    return c
}

// sendRequest performs an HTTP GET request with context, retries, and structured logging.
// The target argument should be a pointer to the expected response struct.
func (c *Client) sendRequest(ctx context.Context, endpoint string, target interface{}) error {
    // Build request with context.
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", BaseURL, endpoint), nil)
    if err != nil {
        return err
    }
    req.Header.Set("x-apisports-key", c.ApiKey)
    req.Header.Set("Content-Type", "application/json")

    // Retry logic – up to 3 attempts with exponential backoff.
    var resp *http.Response
    for attempt := 1; attempt <= 3; attempt++ {
        log.Printf("API request attempt %d: %s", attempt, endpoint)
        resp, err = c.HttpClient.Do(req)
        if err == nil && resp.StatusCode < 500 {
            break // success or client error – no retry.
        }
        // On transient server error, wait before retrying.
        backoff := time.Duration(attempt*100) * time.Millisecond
        log.Printf("Retrying after %v due to error: %v", backoff, err)
        time.Sleep(backoff)
    }
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        // Read body for debugging.
        body, _ := io.ReadAll(resp.Body)
        return &APIError{StatusCode: resp.StatusCode, Message: resp.Status, Body: body}
    }

    // Decode the generic API response wrapper.
    var apiResp APIResponse
    apiResp.Response = target
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return err
    }
    return nil
}

// GetTeam searches for teams by name.
func (c *Client) GetTeam(ctx context.Context, teamName string) ([]TeamResponse, error) {
    var teams []TeamResponse
    endpoint := fmt.Sprintf("/teams?search=%s", teamName)
    err := c.sendRequest(ctx, endpoint, &teams)
    return teams, err
}

// GetLatestFixturesForTeam returns the latest fixtures for a team, limited by `limit`.
func (c *Client) GetLatestFixturesForTeam(ctx context.Context, teamID int, limit int) ([]FixtureResponse, error) {
    var fixtures []FixtureResponse
    endpoint := fmt.Sprintf("/fixtures?team=%d&last=%d", teamID, limit)
    err := c.sendRequest(ctx, endpoint, &fixtures)
    return fixtures, err
}

// GetFixtureDetails fetches detailed information for a fixture.
func (c *Client) GetFixtureDetails(ctx context.Context, fixtureID int) ([]FixtureResponse, error) {
    var details []FixtureResponse
    endpoint := fmt.Sprintf("/fixtures?id=%d", fixtureID)
    err := c.sendRequest(ctx, endpoint, &details)
    return details, err
}

// GetPlayerStatsForFixture returns player statistics for a fixture.
func (c *Client) GetPlayerStatsForFixture(ctx context.Context, fixtureID int) ([]PlayerStatsParent, error) {
    var stats []PlayerStatsParent
    endpoint := fmt.Sprintf("/fixtures/players?fixture=%d", fixtureID)
    err := c.sendRequest(ctx, endpoint, &stats)
    return stats, err
}

// NextPage extracts the next page number from the API response paging info.
func (c *Client) NextPage(paging Paging) (int, bool) {
    if paging.Current < paging.Total {
        return paging.Current + 1, true
    }
    return 0, false
}

// Helper to create a context with a per‑request timeout.
func (c *Client) contextWithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), d)
}

// Paging struct is defined in models.go and contains fields like Current, Total, etc.
