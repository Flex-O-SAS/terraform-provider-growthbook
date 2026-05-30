---
title: "growthbook_segment Resource"
description: |-
  Provides a GrowthBook Segment resource.
---

# growthbook_segment

Manages a GrowthBook segment. Segments define reusable user cohorts from a data source using SQL or fact table filters.

## Example Usage

Creation requires a GrowthBook data source ID.

```hcl
resource "growthbook_segment" "active_customers" {
  name            = "Active customers"
  type            = "SQL"
  datasource_id   = "ds_abc123"
  identifier_type = "user_id"
  query           = "SELECT user_id FROM users WHERE active = true"
  owner           = "owner@example.com"
  description     = "Users with an active account"
  projects        = ["project_abc123"]
  managed_by      = "api"
}
```

## Argument Reference

- `name` (String, Required) – Name of the segment.
- `type` (String, Required) – GrowthBook segment type. Valid values are `SQL` and `FACT`.
- `datasource_id` (String, Required) – ID of the data source this segment belongs to.
- `identifier_type` (String, Required) – Type of identifier, such as `user_id` or `anonymous_id`.
- `owner` (String, Optional) – The user ID or email address of the owner.
- `description` (String, Optional) – Description of the segment.
- `query` (String, Optional) – SQL query that defines the segment. Required by GrowthBook for `SQL` segments.
- `fact_table_id` (String, Optional) – ID of the fact table this segment belongs to. Required by GrowthBook for `FACT` segments.
- `projects` (List(String), Optional) – List of project IDs for projects that can access this segment.
- `managed_by` (String, Optional) – Where this segment must be managed from. Valid values are ``, `api`, and `config`.

## Attributes Reference

- `id` (String) – The unique ID of the segment.
- `date_created` (String) – The creation date of the segment.
- `date_updated` (String) – The last update date of the segment.

## Import

Segments can be imported using the segment ID:

```sh
terraform import growthbook_segment.example <segment_id>
```
