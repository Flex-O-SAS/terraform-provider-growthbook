package growthbookapi

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCheckStatuses_IncludesResponseBodyInError(t *testing.T) {
	t.Parallel()

	body := `{"message":"Request body: [environments.dev.rules.0] Unrecognized key: \"coverage\""}`
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}

	err := checkStatuses("POST", resp)

	if err == nil {
		t.Fatal("expected an error for status 400")
	}
	if !strings.Contains(err.Error(), `Unrecognized key: \"coverage\"`) {
		t.Errorf("error should carry the GrowthBook message, got: %s", err)
	}
}
