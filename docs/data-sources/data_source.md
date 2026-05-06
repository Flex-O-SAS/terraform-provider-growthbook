---
title: "growthbook_data_source Data Source"
description: |-
  Provides a GrowthBook Data Source data source.
---

# growthbook_data_source (Data Source)

Retrieves information about a read-only GrowthBook data source by ID.

## Example Usage

```hcl
data "growthbook_data_source" "warehouse" {
  id = "ds_abc123"
}
```

## Argument Reference

- `id` (String, Required) – The unique ID of the data source to look up.

## Attributes Reference

- `name` (String) – The name of the data source.
- `type` (String) – The data source type.
- `description` (String) – The description of the data source.
- `project_ids` (List(String)) – Project IDs associated with the data source.
- `event_tracker` (String) – The data source event tracker.
- `date_created` (String) – The creation date of the data source.
- `date_updated` (String) – The last update date of the data source.
