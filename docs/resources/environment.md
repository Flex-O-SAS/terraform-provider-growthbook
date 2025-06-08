---
title: "growthbook_environment Resource"
description: |-
  Provides a GrowthBook Environment resource.
---

# growthbook_environment

Manages a GrowthBook environment.

## Example Usage

```hcl
resource "growthbook_environment" "example" {
  name           = "production"
  description    = "Production environment"
  toggle_on_list = true
  default_state  = true
  projects       = [growthbook_project.example.id]
}
```

## Argument Reference

- `name` (String, Required) – The name of the environment.
- `description` (String, Optional) – The description of the environment.
- `toggle_on_list` (Boolean, Optional) – Whether the environment is toggled on in the list.
- `default_state` (Boolean, Optional) – The default state of the environment.
- `projects` (List of String, Optional) – List of project IDs associated with the environment.

## Attributes Reference

- `id` (String) – The unique ID of the environment.

## Import

Environments can be imported using the environment ID:

```sh
terraform import growthbook_environment.example <environment_id>
```
