---
title: "growthbook_fact_table Data Source"
description: |-
  Provides a GrowthBook Fact Table data source.
---

# growthbook_fact_table (Data Source)

Retrieves information about a GrowthBook fact table by name.

## Example Usage

```hcl
data "growthbook_fact_table" "orders" {
  name = "Orders"
}
```

## Argument Reference

- `name` (String, Required) – Name of the fact table to look up.

## Attributes Reference

- `id` (String) – The unique ID of the fact table.
- `datasource` (String) – The data source ID.
- `user_id_types` (List(String)) – List of identifier columns in this table, such as `id` or `anonymous_id`.
- `sql` (String) – The SQL query for this fact table.
- `description` (String) – Description of the fact table.
- `owner` (String) – The user ID or email address of the owner.
- `projects` (List(String)) – List of associated project IDs.
- `tags` (List(String)) – List of associated tags.
- `event_name` (String) – The event name used in SQL template variables.
- `managed_by` (String) – Where this fact table must be managed from. Valid values are ``, `api`, and `admin`.
- `archived` (Boolean) – Whether the fact table is archived.
- `date_created` (String) – The creation date of the fact table.
- `date_updated` (String) – The last update date of the fact table.
