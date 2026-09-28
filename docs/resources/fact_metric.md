---
title: "growthbook_fact_metric Resource"
description: |-
  Provides a GrowthBook Fact Metric resource.
---

# growthbook_fact_metric

Manages a GrowthBook fact metric. Fact metrics define experiment metrics based on fact tables, including numerator, denominator, capping, window, prior, and regression-adjustment settings.

## Example Usage

Creation requires a GrowthBook data source ID in `datasource`. JSON object fields are provided as JSON-encoded strings.

```hcl
resource "growthbook_fact_metric" "revenue_per_user" {
  name        = "Revenue per user"
  metric_type = "mean"
  datasource  = "ds_abc123"
  owner       = "owner@example.com"
  projects    = ["project_abc123"]
  tags        = ["revenue"]
  inverse     = false
  managed_by  = "api"
  archived    = false

  numerator = jsonencode({
    factTableId = "fact_orders"
    column      = "revenue"
  })

  capping_settings = jsonencode({
    type = "none"
  })

  window_settings = jsonencode({
    type = "conversion"
    delayHours = 0
    windowValue = 72
    windowUnit = "hours"
  })
}
```

## Argument Reference

- `name` (String, Required) – Name of the fact metric.
- `metric_type` (String, Required) – Metric type. Valid values are `proportion`, `retention`, `mean`, `quantile`, `ratio`, and `dailyParticipation`.
- `numerator` (String, Required) – JSON-encoded numerator definition.
- `datasource` (String, Optional) – The data source ID.
- `description` (String, Optional) – Description of the fact metric.
- `owner` (String, Optional) – The user ID or email address of the owner.
- `projects` (List(String), Optional) – List of associated project IDs.
- `tags` (List(String), Optional) – List of associated tags.
- `denominator` (String, Optional) – JSON-encoded denominator definition. Used when `metric_type = "ratio"`.
- `inverse` (Boolean, Optional) – Set to true for things like bounce rate where you want the metric to decrease.
- `capping_settings` (String, Optional) – JSON-encoded settings that control how outliers are handled.
- `window_settings` (String, Optional) – JSON-encoded settings that control the conversion window for the metric.
- `prior_settings` (String, Optional) – JSON-encoded Bayesian prior settings. If omitted, organization defaults are used.
- `regression_adjustment_settings` (String, Optional) – JSON-encoded regression adjustment (CUPED) settings.
- `risk_threshold_success` (Number, Optional) – Deprecated GrowthBook risk threshold for low risk, as a proportion.
- `risk_threshold_danger` (Number, Optional) – Deprecated GrowthBook risk threshold for high risk, as a proportion.
- `min_percent_change` (Number, Optional) – Minimum percent change to consider uplift significant, as a proportion.
- `max_percent_change` (Number, Optional) – Maximum percent change to consider uplift significant, as a proportion.
- `min_sample_size` (Number, Optional) – Minimum sample size for the metric.
- `target_mde` (Number, Optional) – Percentage change to reliably detect before ending an experiment, as a proportion.
- `managed_by` (String, Optional) – Where this fact metric must be managed from. Valid values are ``, `api`, and `admin`.
- `archived` (Boolean, Optional) – Whether the fact metric is archived.

## Attributes Reference

- `id` (String) – The unique ID of the fact metric.
- `date_created` (String) – The creation date of the fact metric.
- `date_updated` (String) – The last update date of the fact metric.

## Import

Fact metrics can be imported using the fact metric ID:

```sh
terraform import growthbook_fact_metric.example <fact_metric_id>
```
