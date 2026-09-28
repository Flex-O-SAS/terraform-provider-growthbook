package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

func TestRulesToAPI_CoverageOnlySentForRolloutRules(t *testing.T) {
	t.Parallel()

	rules := []featureRuleModel{
		{Type: types.StringValue("force"), Value: types.StringValue("true"), Coverage: types.Float64Value(1.0)},
		{Type: types.StringValue("rollout"), Value: types.StringValue("true"), Coverage: types.Float64Value(0.5)},
		{Type: types.StringValue("experiment-ref"), Coverage: types.Float64Value(1.0)},
	}

	out := rulesToAPI(rules)

	if out[0].Coverage != nil {
		t.Errorf("force rule: coverage must not be sent (GrowthBook rejects it), got %v", *out[0].Coverage)
	}
	if out[1].Coverage == nil || *out[1].Coverage != 0.5 {
		t.Errorf("rollout rule: expected coverage 0.5, got %v", out[1].Coverage)
	}
	if out[2].Coverage != nil {
		t.Errorf("experiment-ref rule: coverage must not be sent, got %v", *out[2].Coverage)
	}
}

func TestRulesFromAPI_NonRolloutCoverageMirrorsSchemaDefault(t *testing.T) {
	t.Parallel()

	half := 0.5
	rules := []growthbookapi.FeatureRule{
		{Type: "force", Value: "true"},
		{Type: "rollout", Value: "true", Coverage: &half},
	}

	out := rulesFromAPI(rules)

	// GrowthBook never returns coverage on force rules; null would contradict the planned default (1.0).
	if out[0].Coverage.IsNull() || out[0].Coverage.ValueFloat64() != 1.0 {
		t.Errorf("force rule: expected coverage 1.0 (schema default), got %v", out[0].Coverage)
	}
	if out[1].Coverage.ValueFloat64() != 0.5 {
		t.Errorf("rollout rule: expected coverage 0.5, got %v", out[1].Coverage)
	}
}
