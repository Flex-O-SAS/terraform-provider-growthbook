---
title: "growthbook_project Data Source"
description: |-
  Provides a GrowthBook Project data source.
---

# growthbook_project (Data Source)

Retrieves information about a GrowthBook project by name.

## Example Usage

```hcl
data "growthbook_project" "by_name" {
  name = "my-project"
}
```

## Argument Reference

- `name` (String, Required) – The name of the project to look up.

## Attributes Reference

- `id` (String) – The unique ID of the project.
- `description` (String) – The description of the project.
- `date_created` (String) – The creation date of the project.
- `date_updated` (String) – The last update date of the project.
