terraform {
  required_providers {
    growthbook = {
      source = "local/growthbook/growthbook"
      version = "0.1.0"
    }
  }
}

variable "growthbook_api_key" {
  description = "GrowthBook API Key"
  type        = string
  default     = "secret_admin_yYEDHjQODFrwRjBxpsD9aEVd3dl3iviYWkgdF3f9Bk"
}

variable "growthbook_api_url" {
  description = "GrowthBook API URL"
  type        = string
  default     = "https://api.growthbook.io/api/v1"
}

# Configuration du provider
provider "growthbook" {
  api_key = var.growthbook_api_key
  api_url = var.growthbook_api_url
}

resource "growthbook_attribute" "test" {
  property = "test"
  datatype = "enum"
  enum_values = "test1,test2,test3"
  description = "simple desc"
}

data "growthbook_attribute" "test_data" {
    property = growthbook_attribute.test.property
}

output "growthbook_attribute_test" {
    value = data.growthbook_attribute.test_data
}