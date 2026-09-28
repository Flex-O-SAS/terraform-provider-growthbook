---
title: "growthbook_saved_group Resource"
description: |-
  Provides a GrowthBook Saved Group resource.
---

# growthbook_saved_group

Manages a GrowthBook saved group. Saved groups define reusable targeting groups based on a JSON condition or a list of attribute values.

## Example Usage

```hcl
resource "growthbook_saved_group" "employees" {
  name          = "Employees"
  type          = "list"
  attribute_key = "email"
  values        = ["alice@example.com", "bob@example.com"]
  owner         = "owner@example.com"
  projects      = ["project_abc123"]
  description   = "Internal employee allow list"
}
```

For condition-based saved groups, `condition` is a JSON-encoded string:

```hcl
resource "growthbook_saved_group" "beta_users" {
  name      = "Beta users"
  type      = "condition"
  condition = jsonencode({ beta = true })
}
```

## Argument Reference

- `name` (String, Required) – The display name of the saved group.
- `type` (String, Optional) – The type of saved group. Valid values are `condition` and `list`.
- `condition` (String, Optional) – When `type = "condition"`, the JSON-encoded condition for the group.
- `attribute_key` (String, Optional) – When `type = "list"`, the attribute key the group is based on.
- `values` (List(String), Optional) – When `type = "list"`, the list of values for the attribute key.
- `owner` (String, Optional) – The user ID or email address of the owner.
- `projects` (List(String), Optional) – Array of project IDs that can use this saved group.
- `description` (String, Optional) – The description of the saved group.

## Attributes Reference

- `id` (String) – The unique ID of the saved group.
- `date_created` (String) – The creation date of the saved group.
- `date_updated` (String) – The last update date of the saved group.

## Import

Saved groups can be imported using the saved group ID:

```sh
terraform import growthbook_saved_group.example <saved_group_id>
```
