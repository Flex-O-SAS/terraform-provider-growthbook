package internal

import (
	"encoding/json"
)

// rawJSONToString converts json.RawMessage to a normalized JSON string.
// Returns empty string if msg is nil or empty.
func rawJSONToString(msg json.RawMessage) string {
	if len(msg) == 0 {
		return ""
	}

	var v any
	if err := json.Unmarshal(msg, &v); err != nil {
		return string(msg)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return string(msg)
	}
	return string(b)
}

// stringToRawJSON converts a JSON string to json.RawMessage.
// Returns nil if s is empty.
func stringToRawJSON(s string) json.RawMessage {
	if s == "" {
		return nil
	}
	return json.RawMessage(s)
}
