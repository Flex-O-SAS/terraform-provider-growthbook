---
title: "growthbook_dimension Data Source"
description: |-
  Provides a GrowthBook Dimension data source.
---

# growthbook_dimension (Data Source)

Retrieves information about a GrowthBook dimension by name.

## Example Usage

```hcl
data "growthbook_dimension" "country" {
  name = "Country"
}
```

## Argument Reference

- `name` (String, Required) – Name of the dimension to look up.

## Attributes Reference

- `id` (String) – The unique ID of the dimension.
- `datasource_id` (String) – ID of the data source this dimension belongs to.
- `identifier_type` (String) – Type of identifier, such as `user_id` or `anonymous_id`.
- `query` (String) – SQL query or equivalent for the dimension.
- `description` (String) – Description of the dimension.
- `owner` (String) – The user ID or email address of the owner.
- `managed_by` (String) – Where this dimension must be managed from. Valid values are ``, `api`, and `config`.
- `date_created` (String) – The creation date of the dimension.
- `date_updated` (String) – The last update date of the dimension.
