---
title: "growthbook_metric Resource"
description: |-
  Provides a GrowthBook Metric resource.
---

# growthbook_metric

Manages a GrowthBook legacy metric. Metrics define the SQL or event logic used to analyze experiments.

## Example Usage

Creation requires a GrowthBook data source ID. `behavior` and `sql` are JSON-encoded strings.

```hcl
resource "growthbook_metric" "signup" {
  name          = "Signup conversion"
  datasource_id = "ds_abc123"
  type          = "binomial"
  description   = "Users who completed signup"
  owner         = "owner@example.com"
  projects      = ["project_abc123"]
  tags          = ["activation"]
  archived      = false
  managed_by    = "api"

  behavior = jsonencode({
    goal = "increase"
  })

  sql = jsonencode({
    userIdType = "user_id"
    query      = "SELECT user_id, timestamp FROM signups"
  })
}
```

## Argument Reference

- `name` (String, Required) – Name of the metric.
- `datasource_id` (String, Required) – ID of the data source.
- `type` (String, Required) – Type of metric. Valid values are `binomial`, `count`, `duration`, and `revenue`.
- `description` (String, Optional) – Description of the metric.
- `owner` (String, Optional) – The user ID or email address of the owner.
- `tags` (List(String), Optional) – List of tags.
- `projects` (List(String), Optional) – List of project IDs for projects that can access this metric.
- `archived` (Boolean, Optional) – Whether the metric is archived.
- `behavior` (String, Optional) – JSON-encoded behavior settings for the metric.
- `sql` (String, Optional) – JSON-encoded SQL metric definition. The GrowthBook API allows only one of `sql`, `sqlBuilder`, or `mixpanel`; this provider exposes `sql`.
- `managed_by` (String, Optional) – Where this metric must be managed from. Valid values are ``, `api`, `config`, and `admin`.

## Attributes Reference

- `id` (String) – The unique ID of the metric.
- `date_created` (String) – The creation date of the metric.
- `date_updated` (String) – The last update date of the metric.

## Import

Metrics can be imported using the metric ID:

```sh
terraform import growthbook_metric.example <metric_id>
```
