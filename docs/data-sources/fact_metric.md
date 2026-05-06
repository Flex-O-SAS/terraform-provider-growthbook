---
title: "growthbook_fact_metric Data Source"
description: |-
  Provides a GrowthBook Fact Metric data source.
---

# growthbook_fact_metric (Data Source)

Retrieves information about a GrowthBook fact metric by name.

## Example Usage

```hcl
data "growthbook_fact_metric" "revenue_per_user" {
  name = "Revenue per user"
}
```

## Argument Reference

- `name` (String, Required) – Name of the fact metric to look up.

## Attributes Reference

- `id` (String) – The unique ID of the fact metric.
- `metric_type` (String) – Metric type. Valid values are `proportion`, `retention`, `mean`, `quantile`, `ratio`, and `dailyParticipation`.
- `numerator` (String) – JSON-encoded numerator definition.
- `datasource` (String) – The data source ID.
- `description` (String) – Description of the fact metric.
- `owner` (String) – The user ID or email address of the owner.
- `projects` (List(String)) – List of associated project IDs.
- `tags` (List(String)) – List of associated tags.
- `denominator` (String) – JSON-encoded denominator definition.
- `inverse` (Boolean) – Whether lower values are better, such as bounce rate.
- `capping_settings` (String) – JSON-encoded settings that control how outliers are handled.
- `window_settings` (String) – JSON-encoded settings that control the conversion window for the metric.
- `prior_settings` (String) – JSON-encoded Bayesian prior settings.
- `regression_adjustment_settings` (String) – JSON-encoded regression adjustment (CUPED) settings.
- `risk_threshold_success` (Number) – Deprecated GrowthBook risk threshold for low risk, as a proportion.
- `risk_threshold_danger` (Number) – Deprecated GrowthBook risk threshold for high risk, as a proportion.
- `min_percent_change` (Number) – Minimum percent change to consider uplift significant, as a proportion.
- `max_percent_change` (Number) – Maximum percent change to consider uplift significant, as a proportion.
- `min_sample_size` (Number) – Minimum sample size for the metric.
- `target_mde` (Number) – Percentage change to reliably detect before ending an experiment, as a proportion.
- `managed_by` (String) – Where this fact metric must be managed from. Valid values are ``, `api`, and `admin`.
- `archived` (Boolean) – Whether the fact metric is archived.
- `date_created` (String) – The creation date of the fact metric.
- `date_updated` (String) – The last update date of the fact metric.
