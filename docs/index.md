---
page_title: "Provider: Growthbook"
description: |-
    Provider to manage resources on Growthbook
---

# Growthbook Provider

This Terraform provider allows you to manage GrowthBook resources and retrieve data from
your GrowthBook instance using Terraform.

## Configuring the provider

```hcl
provider "growthbook" {
  api_key              = "<your_growthbook_api_key>" # or set GROWTHBOOK_API_KEY env var
  api_url              = "https://api.growthbook.io/api/v1" # optional, override for self-hosted
  http_timeout         = 60 # optional, HTTP timeout in seconds (default: 60)
  insecure_skip_verify = false # optional, skip SSL verification (not recommended)
}
```

## Schema

## Required

- `api_key`: (String) GrowthBook API key. Can also be set via `GROWTHBOOK_API_KEY` env variable.

### Optional

- `api_url`: (String)  GrowthBook API base URL. Defaults to `https://api.growthbook.io/api/v1`.
- `http_timeout`: (Integer) Timeout (in seconds) for HTTP requests. Defaults to `60`.
- `insecure_skip_verify`: (Boolean) If true, disables SSL certificate verification (not recommended for prod). Defaults to `false`.


## Example usage

```hcl
resource "growthbook_project" "example" {
  name        = "my-project"
  description = "A sample project"
}

resource "growthbook_environment" "example" {
  name           = "production"
  description    = "Production environment"
  toggle_on_list = true
  default_state  = true
  projects       = [growthbook_project.example.id]
}

resource "growthbook_feature" "example" {
  id           = "my-feature"
  description  = "A feature flag"
  owner        = "owner@example.com"
  project      = growthbook_project.example.id
  value_type   = "boolean"
  default_value = "false"
  tags         = ["beta", "internal"]
}

resource "growthbook_sdk_connection" "example" {
  name        = "SDK Connection"
  language    = "go"
  environment = growthbook_environment.example.id
  projects    = [growthbook_project.example.id]
}

data "growthbook_project" "by_name" {
  name = "my-project"
}

data "growthbook_feature" "by_id" {
  id = "my-feature"
}

data "growthbook_environment" "by_id" {
  id = "production"
}

data "growthbook_sdk_connection" "by_id" {
  id = growthbook_sdk_connection.example.id
}
```
