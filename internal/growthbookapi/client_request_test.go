package growthbookapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// stubAPI answers with the given statuses in order, repeating the last one, and
// records the body of every request it serves.
type stubAPI struct {
	statuses []int
	bodies   []string
	next     int
}

func (s *stubAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()
	s.bodies = append(s.bodies, string(body))

	status := s.statuses[len(s.statuses)-1]
	if s.next < len(s.statuses) {
		status = s.statuses[s.next]
	}
	s.next++

	if status == http.StatusTooManyRequests {
		w.Header().Set("Retry-After", "0")
	}
	w.WriteHeader(status)
	if status < 300 {
		_, _ = w.Write([]byte(`{"feature":{"id":"feat"}}`))
	}
}

func newTestClient(t *testing.T, statuses []int) (*Client, *stubAPI) {
	t.Helper()

	stub := &stubAPI{statuses: statuses}
	server := httptest.NewServer(stub)
	t.Cleanup(server.Close)

	return &Client{
		BaseURL:    server.URL,
		APIKey:     "api_key",
		HTTPClient: server.Client(),
		Backoff: BackoffConfig{
			MaxRetries:      3,
			InitialInterval: time.Millisecond,
			Multiplier:      2.0,
			MaxInterval:     5 * time.Millisecond,
		},
	}, stub
}

func TestGetFeatureRetriesRateLimitedRequest(t *testing.T) {
	client, stub := newTestClient(t, []int{http.StatusTooManyRequests, http.StatusOK})

	feature, err := client.GetFeature(context.Background(), "feat")
	if err != nil {
		t.Fatalf("GetFeature: %v", err)
	}
	if feature.ID != "feat" {
		t.Fatalf("got id %q, want %q", feature.ID, "feat")
	}
	if len(stub.bodies) != 2 {
		t.Fatalf("served %d requests, want 2", len(stub.bodies))
	}
}

func TestCreateFeatureReplaysBodyOnRetry(t *testing.T) {
	client, stub := newTestClient(t, []int{http.StatusTooManyRequests, http.StatusOK})

	if _, err := client.CreateFeature(context.Background(), &Feature{ID: "feat"}); err != nil {
		t.Fatalf("CreateFeature: %v", err)
	}
	if len(stub.bodies) != 2 {
		t.Fatalf("served %d requests, want 2", len(stub.bodies))
	}
	if stub.bodies[0] == "" || stub.bodies[0] != stub.bodies[1] {
		t.Fatalf("retry sent %q, want a replay of the original %q", stub.bodies[1], stub.bodies[0])
	}
	if !strings.Contains(stub.bodies[1], `"id":"feat"`) {
		t.Fatalf("retry body %q does not carry the payload", stub.bodies[1])
	}
}

func TestGetFeatureGivesUpAfterMaxRetries(t *testing.T) {
	client, stub := newTestClient(t, []int{http.StatusTooManyRequests})

	if _, err := client.GetFeature(context.Background(), "feat"); err == nil {
		t.Fatal("want an error once the retry budget is exhausted")
	}
	if len(stub.bodies) != 1+client.Backoff.MaxRetries {
		t.Fatalf("served %d requests, want %d", len(stub.bodies), 1+client.Backoff.MaxRetries)
	}
}

func TestGetFeatureRetriesServerError(t *testing.T) {
	client, stub := newTestClient(t, []int{http.StatusInternalServerError, http.StatusOK})

	if _, err := client.GetFeature(context.Background(), "feat"); err != nil {
		t.Fatalf("GetFeature: %v", err)
	}
	if len(stub.bodies) != 2 {
		t.Fatalf("served %d requests, want 2", len(stub.bodies))
	}
}
