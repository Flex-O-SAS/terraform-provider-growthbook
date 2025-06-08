# GrowthBook Terraform Provider

The GrowthBook Terraform Provider enables you to manage GrowthBook resources (projects, features, environments, metrics, data sources, experiments, groups, segments, API keys, webhooks, integrations, SDK connections, members, and organizations) via Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- GrowthBook account and API key

## Provider Configuration

```hcl
provider "growthbook" {
  api_key              = "<your_growthbook_api_key>" # or set GROWTHBOOK_API_KEY env var
  api_url              = "https://app.growthbook.io/api/v1" # optional, override for self-hosted
  http_timeout         = 60 # optional, HTTP timeout in seconds (default: 60)
  insecure_skip_verify = false # optional, skip SSL verification (not recommended)
}
```

### Arguments

| Name                  | Type    | Required | Default                                 | Description                                                                 |
|-----------------------|---------|----------|-----------------------------------------|-----------------------------------------------------------------------------|
| `api_key`             | string  | yes      |                                         | GrowthBook API key. Can also be set via `GROWTHBOOK_API_KEY` env variable.  |
| `api_url`             | string  | no       | `https://app.growthbook.io/api/v1`      | GrowthBook API base URL.                                                    |
| `http_timeout`        | int     | no       | `60`                                    | Timeout (in seconds) for HTTP requests.                                     |
| `insecure_skip_verify`| bool    | no       | `false`                                 | If true, disables SSL certificate verification (not recommended for prod).   |

## Resources

- `growthbook_project`
- `growthbook_feature`
- `growthbook_environment`
- `growthbook_metric`
- `growthbook_data_source`
- `growthbook_experiment`
- `growthbook_group`
- `growthbook_segment`
- `growthbook_api_key`
- `growthbook_webhook`
- `growthbook_integration`
- `growthbook_sdk_connection`
- `growthbook_member`
- `growthbook_organization`

See the [docs](https://registry.terraform.io/providers/growthbook/growthbook/latest/docs/resources) for resource-specific arguments and import instructions.

## Import

All resources support import. For example:

```sh
terraform import growthbook_project.example <project_id>
```

## Example Usage

```hcl
provider "growthbook" {
  api_key = var.growthbook_api_key
}

resource "growthbook_project" "example" {
  name = "My Project"
}

resource "growthbook_feature" "example" {
  project_id = growthbook_project.example.id
  key        = "my-feature"
  name       = "My Feature"
}
```

## Logging

All API requests and responses are logged using Go's slog as structured JSON. API keys are redacted from logs. Logging output and level can be configured in future releases.

## Error Handling

All errors are reported with user-friendly messages and diagnostics. HTTP timeouts and SSL errors are surfaced clearly.


