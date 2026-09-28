package growthbookapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteFeature_ArchivesBeforeDeleting(t *testing.T) {
	t.Parallel()

	var calls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch r.Method {
		case http.MethodPost:
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["archived"] != true {
				t.Errorf("archive request must set archived=true, got %v", body)
			}
			_, _ = w.Write([]byte(`{"feature":{"id":"f1","archived":true}}`))
		case http.MethodDelete:
			_, _ = w.Write([]byte(`{"deletedId":"f1"}`))
		}
	}))
	defer srv.Close()

	c, ok := NewClient(srv.URL, "key").(*Client)
	if !ok {
		t.Fatal("NewClient must return *Client")
	}

	if err := c.DeleteFeature(context.Background(), "f1"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := []string{"POST /features/f1", "DELETE /features/f1"}
	if len(calls) != len(want) || calls[0] != want[0] || calls[1] != want[1] {
		t.Errorf("expected calls %v, got %v", want, calls)
	}
}
