---
title: "growthbook_feature Data Source"
description: |-
  Provides a GrowthBook Feature data source.
---

# growthbook_feature (Data Source)

Retrieves information about a GrowthBook feature by ID.

## Example Usage

```hcl
data "growthbook_feature" "by_id" {
  id = "my-feature"
}
```

## Argument Reference

- `id` (String, Required) – The ID of the feature to look up.

## Attributes Reference

- `description` (String) – The description of the feature.
- `owner` (String) – The owner of the feature.
- `project` (String) – The project ID this feature belongs to.
- `value_type` (String) – The type of value for the feature.
- `default_value` (String) – The default value for the feature.
- `tags` (List of String) – Tags associated with the feature.
- `archived` (Boolean) – Whether the feature is archived.
- `environments` (Map of Object) – Map of environment configs for the feature.
- `prerequisites` (List of String) – List of prerequisite feature IDs.
- `date_created` (String) – The creation date of the feature.
- `date_updated` (String) – The last update date of the feature.
