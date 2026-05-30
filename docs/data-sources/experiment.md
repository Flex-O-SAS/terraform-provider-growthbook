---
title: "growthbook_experiment Data Source"
description: |-
  Provides a GrowthBook Experiment data source.
---

# growthbook_experiment (Data Source)

Retrieves information about a GrowthBook experiment by name.

## Example Usage

```hcl
data "growthbook_experiment" "checkout_cta" {
  name = "Checkout CTA"
}
```

## Argument Reference

- `name` (String, Required) – Name of the experiment to look up.

## Attributes Reference

- `id` (String) – The unique ID of the experiment.
- `tracking_key` (String) – Tracking key for the experiment.
- `variations` (String) – JSON-encoded array of experiment variations.
- `type` (String) – Experiment type. Valid values are `standard` and `multi-armed-bandit`.
- `project` (String) – Project ID which the experiment belongs to.
- `hypothesis` (String) – Hypothesis of the experiment.
- `description` (String) – Description of the experiment.
- `tags` (List(String)) – Tags associated with the experiment.
- `owner` (String) – The user ID or email address of the owner.
- `archived` (Boolean) – Whether the experiment is archived.
- `status` (String) – Experiment status. Valid values are `draft`, `running`, and `stopped`.
- `auto_refresh` (Boolean) – Whether experiment results automatically refresh.
- `hash_attribute` (String) – Attribute used for hashing users into variations.
- `fallback_attribute` (String) – Fallback attribute used for hashing.
- `datasource_id` (String) – ID of the data source used for analysis.
- `assignment_query_id` (String) – ID of an assignment query associated with the data source.
- `segment_id` (String) – Only users in this segment are included.
- `metrics` (List(String)) – Primary metric IDs.
- `secondary_metrics` (List(String)) – Secondary metric IDs.
- `guardrail_metrics` (List(String)) – Guardrail metric IDs.
- `stats_engine` (String) – Stats engine. Valid values are `bayesian` and `frequentist`.
- `phases` (String) – JSON-encoded array of experiment phases.
- `date_created` (String) – The creation date of the experiment.
- `date_updated` (String) – The last update date of the experiment.
