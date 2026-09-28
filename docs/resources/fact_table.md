---
title: "growthbook_fact_table Resource"
description: |-
  Provides a GrowthBook Fact Table resource.
---

# growthbook_fact_table

Manages a GrowthBook fact table. Fact tables define reusable SQL event tables used by fact metrics and fact segments.

## Example Usage

Creation requires a GrowthBook data source ID in `datasource`.

```hcl
resource "growthbook_fact_table" "orders" {
  name          = "Orders"
  datasource    = "ds_abc123"
  user_id_types = ["user_id"]
  sql           = "SELECT user_id, timestamp, revenue FROM orders"
  description   = "Order events"
  owner         = "owner@example.com"
  projects      = ["project_abc123"]
  tags          = ["revenue"]
  event_name    = "Order Completed"
  managed_by    = "api"
  archived      = false
}
```

## Argument Reference

- `name` (String, Required) – Name of the fact table.
- `datasource` (String, Required) – The data source ID.
- `user_id_types` (List(String), Required) – List of identifier columns in this table. For example, `id` or `anonymous_id`.
- `sql` (String, Required) – The SQL query for this fact table.
- `description` (String, Optional) – Description of the fact table.
- `owner` (String, Optional) – The user ID or email address of the owner.
- `projects` (List(String), Optional) – List of associated project IDs.
- `tags` (List(String), Optional) – List of associated tags.
- `event_name` (String, Optional) – The event name used in SQL template variables.
- `managed_by` (String, Optional) – Where this fact table must be managed from. Valid values are ``, `api`, and `admin`.
- `archived` (Boolean, Optional) – Whether the fact table is archived.

## Attributes Reference

- `id` (String) – The unique ID of the fact table.
- `date_created` (String) – The creation date of the fact table.
- `date_updated` (String) – The last update date of the fact table.

## Import

Fact tables can be imported using the fact table ID:

```sh
terraform import growthbook_fact_table.example <fact_table_id>
```
