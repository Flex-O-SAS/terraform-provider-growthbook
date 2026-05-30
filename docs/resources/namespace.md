---
title: "growthbook_namespace Resource"
description: |-
  Provides a GrowthBook Namespace resource.
---

# growthbook_namespace

Manages a GrowthBook namespace. Namespaces reserve mutually exclusive traffic ranges for experiments.

## Example Usage

```hcl
resource "growthbook_namespace" "checkout" {
  display_name   = "Checkout experiments"
  description    = "Mutually exclusive checkout experiments"
  status         = "active"
  format         = "multiRange"
  hash_attribute = "id"
}
```

## Argument Reference

- `display_name` (String, Required) – Human-readable display name. Must be unique within the organization.
- `description` (String, Optional) – Description of the namespace.
- `status` (String, Optional) – Namespace status. Valid values are `active` and `inactive`.
- `format` (String, Optional) – Namespace format. Valid values are `legacy` and `multiRange`. `multiRange` supports multiple ranges per experiment and a configurable hash attribute.
- `hash_attribute` (String, Optional) – Required by GrowthBook when `format = "multiRange"`. The user attribute used to assign users to namespace buckets.

## Attributes Reference

- `id` (String) – The unique internal identifier for the namespace.
- `seed` (String) – The seed used for bucket hashing.

## Import

Namespaces can be imported using the namespace ID:

```sh
terraform import growthbook_namespace.example <namespace_id>
```
