package internal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFeatureConfig(id string) string {
	return `
resource "growthbook_project" "test" {
  name = "` + id + `-proj"
}
resource "growthbook_feature" "test" {
  name         = "` + id + `"
  description  = "Acceptance test feature"
  owner        = "owner@example.com"
  project      = growthbook_project.test.id
  value_type   = "boolean"
  default_value = "false"
  tags         = ["acc"]
}
data "growthbook_feature" "by_id" {
  id = growthbook_feature.test.id
}
`
}

func testAccFeatureConfigUpdate(id string) string {
	return `
resource "growthbook_project" "test" {
  name = "` + id + `-proj"
}
resource "growthbook_feature" "test" {
  name         = "` + id + `"
  description  = "new description"
  owner        = "other@example.com"
  project      = growthbook_project.test.id
  value_type   = "json"
  default_value = "{\"key\": \"value\"}"
  tags         = ["new"]
}
data "growthbook_feature" "by_id" {
  id = growthbook_feature.test.id
}
`
}

func TestAccGrowthBookFeature_basic(t *testing.T) {
	t.Parallel()

	featureID := acctest.RandomWithPrefix("tf-acc-feature-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig(featureID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_feature.test", "id", featureID),
					resource.TestCheckResourceAttr("growthbook_feature.test", "description", "Acceptance test feature"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "owner", "owner@example.com"),
					testCheckResourceAttrPrefix("growthbook_feature.test", "project", "prj"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "value_type", "boolean"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "default_value", "false"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "tags.0", "acc"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "id", featureID),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "description", "Acceptance test feature"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "owner", "owner@example.com"),
					testCheckResourceAttrPrefix("data.growthbook_feature.by_id", "project", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "value_type", "boolean"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "default_value", "false"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "tags.#", "1"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "tags.0", "acc"),
				),
			},
			{
				Config: testAccFeatureConfigUpdate(featureID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_feature.test", "id", featureID),
					resource.TestCheckResourceAttr("growthbook_feature.test", "description", "new description"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "owner", "other@example.com"),
					testCheckResourceAttrPrefix("growthbook_feature.test", "project", "prj"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "value_type", "json"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "default_value", "{\"key\": \"value\"}"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("growthbook_feature.test", "tags.0", "new"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "id", featureID),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "description", "new description"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "owner", "other@example.com"),
					testCheckResourceAttrPrefix("data.growthbook_feature.by_id", "project", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "value_type", "json"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "default_value", "{\"key\": \"value\"}"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "tags.#", "1"),
					resource.TestCheckResourceAttr("data.growthbook_feature.by_id", "tags.0", "new"),
				),
			},
		},
	})
}

func testAccFeatureWithRulesConfig(id, forceValue string) string {
	return `
resource "growthbook_feature" "test" {
  name          = "` + id + `"
  owner         = "owner@example.com"
  value_type    = "boolean"
  default_value = "false"
  environments = {
    production = {
      enabled = true
      rules = [
        {
          type      = "force"
          enabled   = true
          condition = "{\"id\":\"someone@example.com\"}"
          value     = "` + forceValue + `"
        },
        {
          type           = "rollout"
          enabled        = true
          value          = "true"
          coverage       = 0.5
          hash_attribute = "id"
        },
      ]
    }
  }
}
`
}

func TestAccGrowthBookFeature_forceAndRolloutRules(t *testing.T) {
	t.Parallel()

	featureID := acctest.RandomWithPrefix("tf-acc-feature-rules-")
	const rules = "environments.production.rules"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureWithRulesConfig(featureID, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".#", "2"),
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".0.type", "force"),
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".0.value", "true"),
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".1.type", "rollout"),
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".1.coverage", "0.5"),
				),
			},
			{
				Config: testAccFeatureWithRulesConfig(featureID, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".0.type", "force"),
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".0.value", "false"),
					resource.TestCheckResourceAttr("growthbook_feature.test", rules+".1.coverage", "0.5"),
				),
			},
		},
	})
}
