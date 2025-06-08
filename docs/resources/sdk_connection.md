---
title: "growthbook_sdk_connection Resource"
description: |-
  Provides a GrowthBook SDK Connection resource.
---

# growthbook_sdk_connection

Manages a GrowthBook SDK Connection.

## Example Usage

```hcl
resource "growthbook_sdk_connection" "example" {
  name        = "SDK Connection"
  language    = "go"
  environment = growthbook_environment.example.id
  projects    = [growthbook_project.example.id]
}
```

## Argument Reference

- `name` (String, Required) – The name of the SDK connection.
- `language` (String, Required) – The language for the SDK connection.
- `environment` (String, Required) – The environment ID for the SDK connection.
- `sdk_version` (String, Optional) – The SDK version.
- `projects` (List of String, Optional) – List of project IDs associated with the SDK connection.
- `encrypt_payload` (Boolean, Optional) – Whether to encrypt the payload.
- `include_visual_experiments` (Boolean, Optional) – Include visual experiments.
- `include_draft_experiments` (Boolean, Optional) – Include draft experiments.
- `include_experiment_names` (Boolean, Optional) – Include experiment names.
- `include_redirect_experiments` (Boolean, Optional) – Include redirect experiments.
- `include_rule_ids` (Boolean, Optional) – Include rule IDs.
- `proxy_enabled` (Boolean, Optional) – Whether proxy is enabled.
- `proxy_host` (String, Optional) – Proxy host.
- `hash_secure_attributes` (Boolean, Optional) – Hash secure attributes.
- `remote_eval_enabled` (Boolean, Optional) – Enable remote evaluation.
- `saved_group_references_enabled` (Boolean, Optional) – Enable saved group references.

## Attributes Reference

- `id` (String) – The unique ID of the SDK connection.
- `organization` (String) – The organization ID.
- `key` (String, Sensitive) – The SDK key.
- `proxy_signing_key` (String, Sensitive) – The proxy signing key.
- `sse_enabled` (Boolean) – Whether SSE is enabled.
- `encryption_key` (String, Sensitive) – The encryption key.
- `date_created` (String) – The creation date of the SDK connection.
- `date_updated` (String) – The last update date of the SDK connection.

## Import

SDK Connections can be imported using the SDK connection ID:

```sh
terraform import growthbook_sdk_connection.example <sdk_connection_id>
```
