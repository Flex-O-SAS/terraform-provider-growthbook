---
title: "growthbook_saved_group Data Source"
description: |-
  Provides a GrowthBook Saved Group data source.
---

# growthbook_saved_group (Data Source)

Retrieves information about a GrowthBook saved group by name.

## Example Usage

```hcl
data "growthbook_saved_group" "employees" {
  name = "Employees"
}
```

## Argument Reference

- `name` (String, Required) – The display name of the saved group to look up.

## Attributes Reference

- `id` (String) – The unique ID of the saved group.
- `type` (String) – The type of saved group. Valid values are `condition` and `list`.
- `condition` (String) – When `type = "condition"`, the JSON-encoded condition for the group.
- `attribute_key` (String) – When `type = "list"`, the attribute key the group is based on.
- `values` (List(String)) – When `type = "list"`, the list of values for the attribute key.
- `owner` (String) – The user ID or email address of the owner.
- `projects` (List(String)) – Array of project IDs that can use this saved group.
- `description` (String) – The description of the saved group.
- `date_created` (String) – The creation date of the saved group.
- `date_updated` (String) – The last update date of the saved group.
