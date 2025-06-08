package internal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccEnvironmentConfig(name string) string {
	return `
resource "growthbook_project" "test" {
  name = "` + name + `-proj"
}
resource "growthbook_environment" "test" {
  name        = "` + name + `"
  description = "Acceptance test environment"
  projects    = [growthbook_project.test.id]
}
data "growthbook_environment" "by_id" {
  id = growthbook_environment.test.id
}
`
}

func testAccEnvironmentConfigUpdate(name string) string {
	return `
resource "growthbook_project" "test" {
  name = "` + name + `-proj"
}
resource "growthbook_environment" "test" {
  name        = "` + name + `"
  description = "new description"
  toggle_on_list = true
  default_state = true
  projects    = [growthbook_project.test.id]
}
data "growthbook_environment" "by_id" {
  id = growthbook_environment.test.id
}
`
}

func TestAccGrowthBookEnvironment_basic(t *testing.T) {
	envName := acctest.RandomWithPrefix("tf-acc-env-")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig(envName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_environment.test", "id", envName),
					resource.TestCheckResourceAttr("growthbook_environment.test", "description", "Acceptance test environment"),
					resource.TestCheckResourceAttr("growthbook_environment.test", "projects.#", "1"),
					testCheckResourceAttrPrefix("growthbook_environment.test", "projects.0", "prj"),
					resource.TestCheckResourceAttr("growthbook_environment.test", "toggle_on_list", "false"),
					resource.TestCheckResourceAttr("growthbook_environment.test", "default_state", "false"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "id", envName),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "description", "Acceptance test environment"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "projects.#", "1"),
					testCheckResourceAttrPrefix("data.growthbook_environment.by_id", "projects.0", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "toggle_on_list", "false"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "default_state", "false"),
				),
			},
			{
				Config: testAccEnvironmentConfigUpdate(envName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_environment.test", "id", envName),
					resource.TestCheckResourceAttr("growthbook_environment.test", "description", "new description"),
					resource.TestCheckResourceAttr("growthbook_environment.test", "projects.#", "1"),
					testCheckResourceAttrPrefix("growthbook_environment.test", "projects.0", "prj"),
					resource.TestCheckResourceAttr("growthbook_environment.test", "toggle_on_list", "true"),
					resource.TestCheckResourceAttr("growthbook_environment.test", "default_state", "true"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "id", envName),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "description", "new description"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "projects.#", "1"),
					testCheckResourceAttrPrefix("data.growthbook_environment.by_id", "projects.0", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "toggle_on_list", "true"),
					resource.TestCheckResourceAttr("data.growthbook_environment.by_id", "default_state", "true"),
				),
			},
		},
	})
}
