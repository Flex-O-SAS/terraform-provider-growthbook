---
title: "growthbook_member_role Resource"
description: |-
  Manages the role and permissions of an existing GrowthBook organization member.
---

# growthbook_member_role

Manages the role and permissions of an existing GrowthBook organization member.

~> **Note:** This resource does not create or invite members. Members must already exist in the organization (invited via the GrowthBook UI). This resource adopts an existing member by email and manages their role configuration.

~> **Warning:** Destroying this resource will **remove the member from the organization**.

## Example Usage

```hcl
resource "growthbook_member_role" "joao" {
  email = "joao@example.com"
  role  = "engineer"

  limit_access_by_environment = true
  environments                = ["production", "staging"]

  project_roles {
    project                     = "prj_main"
    role                        = "admin"
    limit_access_by_environment = false
    environments                = []
  }
}
```

## Argument Reference

- `email` (String, Required, ForceNew) – The member's email address. Used to look up the existing member.
- `role` (String, Required) – The member's global role (e.g. `readonly`, `collaborator`, `engineer`, `analyst`, `experimenter`, `admin`).
- `limit_access_by_environment` (Boolean, Optional) – Whether to restrict access to specific environments. Defaults to `false`.
- `environments` (List of String, Optional) – Environments the member has access to. Empty means all.
- `project_roles` (List of Object, Optional) – Per-project role overrides. Each object supports:
  - `project` (String, Required) – The project ID.
  - `role` (String, Required) – The role for this project.
  - `limit_access_by_environment` (Boolean, Optional) – Defaults to `false`.
  - `environments` (List of String, Optional) – Environments for this project role.

## Attributes Reference

- `id` (String) – The unique member ID.
- `name` (String) – The member's display name.
- `teams` (List of String) – Team IDs the member belongs to.
- `managed_by_idp` (Boolean) – Whether the member is managed by an identity provider.
- `last_login_date` (String) – The member's last login date.
- `date_created` (String) – The creation date.
- `date_updated` (String) – The last update date.

## Import

Member roles can be imported using the member ID:

```sh
terraform import growthbook_member_role.example <member_id>
```
