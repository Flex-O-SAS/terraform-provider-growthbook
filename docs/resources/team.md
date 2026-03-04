---
title: "growthbook_team Resource"
description: |-
  Provides a GrowthBook Team resource.
---

# growthbook_team

Manages a GrowthBook team. Teams allow you to assign roles and permissions to groups of members.

## Example Usage

```hcl
resource "growthbook_team" "engineering" {
  name        = "Engineering"
  description = "Backend engineering team"
  role        = "engineer"

  environments = ["production", "staging"]
  limit_access_by_environment = true

  members = ["user_abc123", "user_def456"]

  project_roles {
    project                     = "prj_main"
    role                        = "admin"
    limit_access_by_environment = false
    environments                = []
  }
}
```

## Argument Reference

- `name` (String, Required) – The team name.
- `description` (String, Optional) – The team description.
- `role` (String, Required) – The global role for team members (e.g. `readonly`, `collaborator`, `engineer`, `analyst`, `experimenter`, `admin`).
- `limit_access_by_environment` (Boolean, Optional) – Whether to restrict access to specific environments. Defaults to `false`.
- `environments` (List of String, Optional) – Environments the team has access to. Empty means all.
- `members` (Set of String, Optional) – Member IDs in this team.
- `project_roles` (List of Object, Optional) – Per-project role overrides. Each object supports:
  - `project` (String, Required) – The project ID.
  - `role` (String, Required) – The role for this project.
  - `limit_access_by_environment` (Boolean, Optional) – Defaults to `false`.
  - `environments` (List of String, Optional) – Environments for this project role.

## Attributes Reference

- `id` (String) – The unique ID of the team.
- `managed_by_idp` (Boolean) – Whether the team is managed by an identity provider.
- `created_by` (String) – The user who created the team.
- `date_created` (String) – The creation date.
- `date_updated` (String) – The last update date.

## Import

Teams can be imported using the team ID:

```sh
terraform import growthbook_team.example <team_id>
```
