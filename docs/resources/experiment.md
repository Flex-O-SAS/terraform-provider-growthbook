---
title: "growthbook_experiment Resource"
description: |-
  Provides a GrowthBook Experiment resource.
---

# growthbook_experiment

Manages a GrowthBook experiment. Experiments define variations, targeting, metrics, status, analysis settings, and rollout phases.

## Example Usage

`variations` and `phases` are JSON-encoded strings.

```hcl
resource "growthbook_experiment" "checkout_cta" {
  name         = "Checkout CTA"
  tracking_key = "checkout-cta"
  type         = "standard"
  project      = "project_abc123"
  hypothesis   = "Changing the CTA increases checkout completion"
  description  = "CTA copy experiment on checkout"
  tags         = ["checkout"]
  owner        = "owner@example.com"
  status       = "draft"

  datasource_id       = "ds_abc123"
  assignment_query_id = "assignments"
  hash_attribute      = "id"
  metrics             = ["met_signup"]
  secondary_metrics   = ["met_revenue"]
  guardrail_metrics   = ["met_refunds"]
  stats_engine        = "bayesian"

  variations = jsonencode([
    { id = "0", name = "Control", weight = 0.5 },
    { id = "1", name = "New CTA", weight = 0.5 }
  ])

  phases = jsonencode([
    {
      name      = "Main"
      coverage  = 1
      condition = {}
    }
  ])
}
```

## Argument Reference

- `name` (String, Required) – Name of the experiment.
- `tracking_key` (String, Required) – Tracking key for the experiment.
- `variations` (String, Required) – JSON-encoded array of experiment variations.
- `type` (String, Optional) – Experiment type. Valid values are `standard` and `multi-armed-bandit`.
- `project` (String, Optional) – Project ID which the experiment belongs to.
- `hypothesis` (String, Optional) – Hypothesis of the experiment.
- `description` (String, Optional) – Description of the experiment.
- `tags` (List(String), Optional) – Tags associated with the experiment.
- `owner` (String, Optional) – The user ID or email address of the owner.
- `archived` (Boolean, Optional) – Whether the experiment is archived.
- `status` (String, Optional) – Experiment status. Valid values are `draft`, `running`, and `stopped`.
- `auto_refresh` (Boolean, Optional) – Whether experiment results automatically refresh.
- `hash_attribute` (String, Optional) – Attribute used for hashing users into variations.
- `fallback_attribute` (String, Optional) – Fallback attribute used for hashing.
- `datasource_id` (String, Optional) – ID of the data source used for analysis.
- `assignment_query_id` (String, Optional) – ID of an assignment query associated with the data source.
- `segment_id` (String, Optional) – Only users in this segment are included.
- `metrics` (List(String), Optional) – Primary metric IDs.
- `secondary_metrics` (List(String), Optional) – Secondary metric IDs.
- `guardrail_metrics` (List(String), Optional) – Guardrail metric IDs.
- `stats_engine` (String, Optional) – Stats engine. Valid values are `bayesian` and `frequentist`.
- `phases` (String, Optional) – JSON-encoded array of experiment phases.

## Attributes Reference

- `id` (String) – The unique ID of the experiment.
- `date_created` (String) – The creation date of the experiment.
- `date_updated` (String) – The last update date of the experiment.

## Import

Experiments can be imported using the experiment ID:

```sh
terraform import growthbook_experiment.example <experiment_id>
```
