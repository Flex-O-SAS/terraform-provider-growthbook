---
title: "growthbook_feature Resource"
description: |-
  Provides a GrowthBook Feature resource.
---

# growthbook_feature

Manages a GrowthBook feature flag.

## Example Usage

```hcl
resource "growthbook_feature" "example" {
  id            = "my-feature"
  description   = "A feature flag"
  owner         = "owner@example.com"
  project       = growthbook_project.example.id
  value_type    = "boolean"
  default_value = "false"
  tags          = ["beta", "internal"]
}
```

## Argument Reference

- `id` (String, Required) – The unique ID of the feature.
- `description` (String, Optional) – The description of the feature.
- `owner` (String, Optional) – The owner of the feature.
- `project` (String, Required) – The project ID this feature belongs to.
- `value_type` (String, Required) – The type of value for the feature (e.g., `boolean`, `string`).
- `default_value` (String, Required) – The default value for the feature.
- `tags` (List of String, Optional) – Tags associated with the feature.

## Attributes Reference

- `archived` (Boolean) – Whether the feature is archived.
- `environments` (Map of Object) – Map of environment configs for the feature.
- `prerequisites` (List of String) – List of prerequisite feature IDs.
- `date_created` (String) – The creation date of the feature.
- `date_updated` (String) – The last update date of the feature.

## Import

Features can be imported using the feature ID:

```sh
terraform import growthbook_feature.example <feature_id>
```
