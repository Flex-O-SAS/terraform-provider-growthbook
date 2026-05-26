---
title: "growthbook_environments Data Source"
description: |-
  Provides a GrowthBook Environments data source.
---

# growthbook_environments (Data Source)

Retrieves the list of all GrowthBook environments.

## Example Usage

```hcl
data "growthbook_environments" "all" {}

output "environment_ids" {
  value = [for env in data.growthbook_environments.all.environments : env.id]
}
```

## Attributes Reference

- `environments` (List of Object) – The list of all environments. Each object has the following attributes:
  - `id` (String) – The ID of the environment.
  - `description` (String) – The description of the environment.
  - `toggle_on_list` (Boolean) – Whether the toggle is shown on the feature list.
  - `default_state` (Boolean) – The default state for new features in this environment.
  - `projects` (List of String) – List of project IDs associated with the environment.
