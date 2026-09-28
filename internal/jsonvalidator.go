package internal

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// jsonStringValidator validates that a string attribute contains valid JSON.
type jsonStringValidator struct{}

// JSONString returns a validator that ensures the string is valid JSON.
func JSONString() validator.String {
	return jsonStringValidator{}
}

func (v jsonStringValidator) Description(_ context.Context) string {
	return "value must be a valid JSON string"
}

func (v jsonStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonStringValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	s := req.ConfigValue.ValueString()
	if s == "" {
		return
	}
	if !json.Valid([]byte(s)) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON",
			"The value must be a valid JSON string. Use `jsonencode({...})` to construct it.",
		)
	}
}
