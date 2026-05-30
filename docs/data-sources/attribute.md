---
title: "growthbook_attribute Data Source"
description: |-
  Provides a GrowthBook Attribute data source.
---

# growthbook_attribute (Data Source)

Retrieves information about a GrowthBook attribute by property.

## Example Usage

```hcl
data "growthbook_attribute" "plan" {
  property = "plan"
}
```

## Argument Reference

- `property` (String, Required) – The attribute property key to look up.

## Attributes Reference

- `datatype` (String) – The attribute datatype. Valid values are `boolean`, `string`, `number`, `secureString`, `enum`, `string[]`, `number[]`, and `secureString[]`.
- `format` (String) – The attribute format. Valid values are ``, `version`, `date`, and `isoCountryCode`.
- `enum_values` (String) – Comma-separated enum values for enum attributes.
- `projects` (List(String)) – Array of project IDs that can use this attribute.
- `archived` (Boolean) – Whether the attribute is archived.
- `description` (String) – The description of the attribute.
