---
title: "growthbook_project Resource"
description: |-
  Provides a GrowthBook Project resource.
---

# growthbook_project

Manages a GrowthBook project.

## Example Usage

```hcl
resource "growthbook_project" "example" {
  name        = "my-project"
  description = "A sample project"
}
```

## Argument Reference

- `name` (String, Required) – The name of the project.
- `description` (String, Optional) – The description of the project.

## Attributes Reference

- `id` (String) – The unique ID of the project.
- `date_created` (String) – The creation date of the project.
- `date_updated` (String) – The last update date of the project.

## Import

Projects can be imported using the project ID:

```sh
terraform import growthbook_project.example <project_id>
```
