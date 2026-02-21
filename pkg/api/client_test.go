import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

// testRoundTripper rewrites request URLs to the test server.
type testRoundTripper struct {
    serverURL string
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    // Replace the scheme and host with the test server's.
    req.URL.Scheme = "http"
    req.URL.Host = t.serverURL
    return http.DefaultTransport.RoundTrip(req)
}

func TestGetTeamSuccess(t *testing.T) {
    // Mock API response.
    mockResp := APIResponse{
        Get: "/teams",
        Response: []TeamResponse{{Team: Team{ID: 1, Name: "Test FC"}}},
    }
    body, _ := json.Marshal(mockResp)

    // Create test server.
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(body)
    }))
    defer ts.Close()

    // Create client with custom transport to redirect to test server.
    c := NewClient("dummy-key")
    c.HttpClient = &http.Client{Transport: &testRoundTripper{serverURL: ts.Listener.Addr().String()}, Timeout: 5 * time.Second}

    teams, err := c.GetTeam(context.Background(), "Test")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if len(teams) != 1 || teams[0].Team.Name != "Test FC" {
        t.Fatalf("unexpected team data: %+v", teams)
    }
}
