---
title: "growthbook_segment Data Source"
description: |-
  Provides a GrowthBook Segment data source.
---

# growthbook_segment (Data Source)

Retrieves information about a GrowthBook segment by name.

## Example Usage

```hcl
data "growthbook_segment" "active_customers" {
  name = "Active customers"
}
```

## Argument Reference

- `name` (String, Required) – Name of the segment to look up.

## Attributes Reference

- `id` (String) – The unique ID of the segment.
- `type` (String) – GrowthBook segment type. Valid values are `SQL` and `FACT`.
- `datasource_id` (String) – ID of the data source this segment belongs to.
- `identifier_type` (String) – Type of identifier, such as `user_id` or `anonymous_id`.
- `owner` (String) – The user ID or email address of the owner.
- `description` (String) – Description of the segment.
- `query` (String) – SQL query that defines the segment.
- `fact_table_id` (String) – ID of the fact table for `FACT` segments.
- `projects` (List(String)) – List of project IDs for projects that can access this segment.
- `managed_by` (String) – Where this segment must be managed from. Valid values are ``, `api`, and `config`.
- `date_created` (String) – The creation date of the segment.
- `date_updated` (String) – The last update date of the segment.
