---
title: "growthbook_namespace Data Source"
description: |-
  Provides a GrowthBook Namespace data source.
---

# growthbook_namespace (Data Source)

Retrieves information about a GrowthBook namespace by display name.

## Example Usage

```hcl
data "growthbook_namespace" "checkout" {
  display_name = "Checkout experiments"
}
```

## Argument Reference

- `display_name` (String, Required) – Human-readable display name of the namespace to look up.

## Attributes Reference

- `id` (String) – The unique internal identifier for the namespace.
- `description` (String) – Description of the namespace.
- `status` (String) – Namespace status. Valid values are `active` and `inactive`.
- `format` (String) – Namespace format. Valid values are `legacy` and `multiRange`.
- `hash_attribute` (String) – The user attribute used to assign users to namespace buckets.
- `seed` (String) – The seed used for bucket hashing.
