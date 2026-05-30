---
title: "growthbook_dimension Resource"
description: |-
  Provides a GrowthBook Dimension resource.
---

# growthbook_dimension

Manages a GrowthBook dimension. Dimensions define SQL queries used to break experiment results down by user attributes.

## Example Usage

Creation requires a GrowthBook data source ID.

```hcl
resource "growthbook_dimension" "country" {
  name            = "Country"
  datasource_id   = "ds_abc123"
  identifier_type = "user_id"
  query           = "SELECT user_id, country FROM users"
  description     = "User country dimension"
  owner           = "owner@example.com"
  managed_by      = "api"
}
```

## Argument Reference

- `name` (String, Required) – Name of the dimension.
- `datasource_id` (String, Required) – ID of the data source this dimension belongs to.
- `identifier_type` (String, Required) – Type of identifier, such as `user_id` or `anonymous_id`.
- `query` (String, Required) – SQL query or equivalent for the dimension.
- `description` (String, Optional) – Description of the dimension.
- `owner` (String, Optional) – The user ID or email address of the owner.
- `managed_by` (String, Optional) – Where this dimension must be managed from. Valid values are ``, `api`, and `config`.

## Attributes Reference

- `id` (String) – The unique ID of the dimension.
- `date_created` (String) – The creation date of the dimension.
- `date_updated` (String) – The last update date of the dimension.

## Import

Dimensions can be imported using the dimension ID:

```sh
terraform import growthbook_dimension.example <dimension_id>
```
