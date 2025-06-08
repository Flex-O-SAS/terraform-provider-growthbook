---
title: "growthbook_environment Data Source"
description: |-
  Provides a GrowthBook Environment data source.
---

# growthbook_environment (Data Source)

Retrieves information about a GrowthBook environment by ID.

## Example Usage

```hcl
data "growthbook_environment" "by_id" {
  id = "production"
}
```

## Argument Reference

- `id` (String, Required) – The ID of the environment to look up.

## Attributes Reference

- `name` (String) – The name of the environment.
- `description` (String) – The description of the environment.
- `toggle_on_list` (Boolean) – Whether the environment is toggled on in the list.
- `default_state` (Boolean) – The default state of the environment.
- `projects` (List of String) – List of project IDs associated with the environment.
