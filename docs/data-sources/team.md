---
title: "growthbook_team Data Source"
description: |-
  Look up a GrowthBook team by name.
---

# growthbook_team

Use this data source to look up a GrowthBook team by name.

## Example Usage

```hcl
data "growthbook_team" "engineering" {
  name = "Engineering"
}
```

## Argument Reference

- `name` (String, Required) – The name of the team to look up.

## Attributes Reference

- `id` (String) – The unique ID of the team.
- `description` (String) – The team description.
- `role` (String) – The global role for team members.
- `limit_access_by_environment` (Boolean) – Whether access is restricted by environment.
- `environments` (List of String) – Environments the team has access to.
- `project_roles` (List of Object) – Per-project role assignments.
- `members` (List of String) – Member IDs in this team.
- `managed_by_idp` (Boolean) – Whether the team is managed by an identity provider.
- `created_by` (String) – The user who created the team.
- `date_created` (String) – The creation date.
- `date_updated` (String) – The last update date.
