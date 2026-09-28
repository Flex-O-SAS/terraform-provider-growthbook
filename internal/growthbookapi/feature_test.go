package growthbookapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateFeatureAppliesSavedGroupTargetingWithFollowUpUpdate(t *testing.T) {
	t.Parallel()

	var paths []string
	var updateBody Feature
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.Method + " " + r.URL.Path {
		case "POST /features":
			_ = json.NewEncoder(w).Encode(map[string]Feature{
				"feature": {
					ID: "feature-id",
					Environments: map[string]FeatureEnvironmentConfig{
						"production": {
							Enabled: true,
							Rules: []FeatureRule{
								{ID: "rule-id", Type: "force", Enabled: true},
							},
						},
					},
				},
			})
		case "POST /features/feature-id":
			if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
				t.Fatalf("decoding update body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]Feature{"feature": updateBody})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	feature := &Feature{
		ID: "feature-id",
		Environments: map[string]FeatureEnvironmentConfig{
			"production": {
				Enabled: true,
				Rules: []FeatureRule{
					{
						ID:      "rule-id",
						Type:    "force",
						Enabled: true,
						SavedGroupTargeting: []FeatureSavedGroupTargeting{
							{MatchType: "all", SavedGroups: []string{"saved-group-id"}},
						},
					},
				},
			},
		},
	}

	created, err := client.CreateFeature(context.Background(), feature)
	if err != nil {
		t.Fatalf("CreateFeature returned error: %v", err)
	}

	if got, want := len(paths), 2; got != want {
		t.Fatalf("expected %d requests, got %d: %v", want, got, paths)
	}
	if got, want := paths[0], "POST /features"; got != want {
		t.Fatalf("expected first request %q, got %q", want, got)
	}
	if got, want := paths[1], "POST /features/feature-id"; got != want {
		t.Fatalf("expected second request %q, got %q", want, got)
	}
	if got := updateBody.Environments["production"].Rules[0].SavedGroupTargeting; len(got) != 1 {
		t.Fatalf("expected update body to include saved group targeting, got %#v", got)
	}
	if got := created.Environments["production"].Rules[0].SavedGroupTargeting; len(got) != 1 {
		t.Fatalf("expected returned feature to include saved group targeting, got %#v", got)
	}
}

func TestCreateFeatureAppliesScheduleRulesWithFollowUpUpdate(t *testing.T) {
	t.Parallel()

	var paths []string
	var updateBody Feature
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.Method + " " + r.URL.Path {
		case "POST /features":
			_ = json.NewEncoder(w).Encode(map[string]Feature{
				"feature": {
					ID: "feature-id",
					Environments: map[string]FeatureEnvironmentConfig{
						"production": {
							Enabled: true,
							Rules: []FeatureRule{
								{ID: "rule-id", Type: "force", Enabled: true},
							},
						},
					},
				},
			})
		case "POST /features/feature-id":
			if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
				t.Fatalf("decoding update body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]Feature{"feature": updateBody})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	timestamp := "2026-06-01T00:00:00Z"
	feature := &Feature{
		ID: "feature-id",
		Environments: map[string]FeatureEnvironmentConfig{
			"production": {
				Enabled: true,
				Rules: []FeatureRule{
					{
						ID:      "rule-id",
						Type:    "force",
						Enabled: true,
						ScheduleRules: []FeatureScheduleRule{
							{Timestamp: &timestamp, Enabled: true},
						},
					},
				},
			},
		},
	}

	created, err := client.CreateFeature(context.Background(), feature)
	if err != nil {
		t.Fatalf("CreateFeature returned error: %v", err)
	}

	if got, want := len(paths), 2; got != want {
		t.Fatalf("expected %d requests, got %d: %v", want, got, paths)
	}
	if got := updateBody.Environments["production"].Rules[0].ScheduleRules; len(got) != 1 {
		t.Fatalf("expected update body to include schedule rules, got %#v", got)
	}
	if got := created.Environments["production"].Rules[0].ScheduleRules; len(got) != 1 {
		t.Fatalf("expected returned feature to include schedule rules, got %#v", got)
	}
}

func TestCreateFeatureSkipsFollowUpUpdateWithoutUpdateOnlyFields(t *testing.T) {
	t.Parallel()

	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" || r.URL.Path != "/features" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]Feature{
			"feature": {
				ID: "feature-id",
				Environments: map[string]FeatureEnvironmentConfig{
					"production": {
						Enabled: true,
						Rules: []FeatureRule{
							{ID: "rule-id", Type: "force", Enabled: true},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	_, err := client.CreateFeature(context.Background(), &Feature{
		ID: "feature-id",
		Environments: map[string]FeatureEnvironmentConfig{
			"production": {
				Enabled: true,
				Rules: []FeatureRule{
					{ID: "rule-id", Type: "force", Enabled: true},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateFeature returned error: %v", err)
	}
	if got, want := requestCount, 1; got != want {
		t.Fatalf("expected %d request, got %d", want, got)
	}
}
