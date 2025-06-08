---
title: "growthbook_sdk_connection Data Source"
description: |-
  Provides a GrowthBook SDK Connection data source.
---

# growthbook_sdk_connection (Data Source)

Retrieves information about a GrowthBook SDK connection by ID.

## Example Usage

```hcl
data "growthbook_sdk_connection" "by_id" {
  id = growthbook_sdk_connection.example.id
}
```

## Argument Reference

- `id` (String, Required) – The ID of the SDK connection to look up.

## Attributes Reference

- `name` (String) – The name of the SDK connection.
- `language` (String) – The language for the SDK connection.
- `environment` (String) – The environment ID for the SDK connection.
- `sdk_version` (String) – The SDK version.
- `projects` (List of String) – List of project IDs associated with the SDK connection.
- `encrypt_payload` (Boolean) – Whether to encrypt the payload.
- `include_visual_experiments` (Boolean) – Include visual experiments.
- `include_draft_experiments` (Boolean) – Include draft experiments.
- `include_experiment_names` (Boolean) – Include experiment names.
- `include_redirect_experiments` (Boolean) – Include redirect experiments.
- `include_rule_ids` (Boolean) – Include rule IDs.
- `proxy_enabled` (Boolean) – Whether proxy is enabled.
- `proxy_host` (String) – Proxy host.
- `hash_secure_attributes` (Boolean) – Hash secure attributes.
- `remote_eval_enabled` (Boolean) – Enable remote evaluation.
- `saved_group_references_enabled` (Boolean) – Enable saved group references.
- `organization` (String) – The organization ID.
- `key` (String, Sensitive) – The SDK key.
- `proxy_signing_key` (String, Sensitive) – The proxy signing key.
- `sse_enabled` (Boolean) – Whether SSE is enabled.
- `encryption_key` (String, Sensitive) – The encryption key.
- `date_created` (String) – The creation date of the SDK connection.
- `date_updated` (String) – The last update date of the SDK connection.
