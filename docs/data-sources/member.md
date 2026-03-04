---
title: "growthbook_member Data Source"
description: |-
  Look up a GrowthBook organization member by email.
---

# growthbook_member

Use this data source to look up a GrowthBook organization member by email.

## Example Usage

```hcl
data "growthbook_member" "joao" {
  email = "joao@example.com"
}
```

## Argument Reference

- `email` (String, Required) – The email of the member to look up.

## Attributes Reference

- `id` (String) – The unique member ID.
- `name` (String) – The member's display name.
- `role` (String) – The member's global role.
- `limit_access_by_environment` (Boolean) – Whether access is restricted by environment.
- `environments` (List of String) – Environments the member has access to.
- `project_roles` (List of Object) – Per-project role assignments.
- `teams` (List of String) – Team IDs the member belongs to.
- `managed_by_idp` (Boolean) – Whether managed by an identity provider.
- `last_login_date` (String) – The member's last login date.
- `date_created` (String) – The creation date.
- `date_updated` (String) – The last update date.
