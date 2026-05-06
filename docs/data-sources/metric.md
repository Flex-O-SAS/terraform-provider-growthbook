---
title: "growthbook_metric Data Source"
description: |-
  Provides a GrowthBook Metric data source.
---

# growthbook_metric (Data Source)

Retrieves information about a GrowthBook legacy metric by name.

## Example Usage

```hcl
data "growthbook_metric" "signup" {
  name = "Signup conversion"
}
```

## Argument Reference

- `name` (String, Required) – Name of the metric to look up.

## Attributes Reference

- `id` (String) – The unique ID of the metric.
- `datasource_id` (String) – ID of the data source.
- `type` (String) – Type of metric. Valid values are `binomial`, `count`, `duration`, and `revenue`.
- `description` (String) – Description of the metric.
- `owner` (String) – The user ID or email address of the owner.
- `tags` (List(String)) – List of tags.
- `projects` (List(String)) – List of project IDs for projects that can access this metric.
- `archived` (Boolean) – Whether the metric is archived.
- `behavior` (String) – JSON-encoded behavior settings for the metric.
- `sql` (String) – JSON-encoded SQL metric definition.
- `managed_by` (String) – Where this metric must be managed from. Valid values are ``, `api`, `config`, and `admin`.
- `date_created` (String) – The creation date of the metric.
- `date_updated` (String) – The last update date of the metric.
