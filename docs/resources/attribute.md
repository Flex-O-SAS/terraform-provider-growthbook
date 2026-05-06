---
title: "growthbook_attribute Resource"
description: |-
  Provides a GrowthBook Attribute resource.
---

# growthbook_attribute

Manages a GrowthBook attribute. Attributes define typed user properties that can be used for targeting, hashing, and analysis.

## Example Usage

```hcl
resource "growthbook_attribute" "plan" {
  property    = "plan"
  datatype    = "enum"
  format      = ""
  enum_values = "free,pro,enterprise"
  description = "Customer subscription plan"
  projects    = ["project_abc123"]
  archived    = false
}
```

## Argument Reference

- `property` (String, Required) – The attribute property key.
- `datatype` (String, Required) – The attribute datatype. Valid values are `boolean`, `string`, `number`, `secureString`, `enum`, `string[]`, `number[]`, and `secureString[]`.
- `format` (String, Optional) – The attribute format. Valid values are ``, `version`, `date`, and `isoCountryCode`.
- `enum_values` (String, Optional) – Comma-separated enum values for attributes with `datatype = "enum"`.
- `projects` (List(String), Optional) – Array of project IDs that can use this attribute.
- `archived` (Boolean, Optional) – Whether the attribute is archived.
- `description` (String, Optional) – The description of the attribute.

## Attributes Reference

All arguments above are also exported. GrowthBook attributes are identified by `property` and do not have a separate Terraform `id` attribute.

## Import

Attributes can be imported using the attribute property:

```sh
terraform import growthbook_attribute.example <property>
```
