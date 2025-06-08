package internal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGrowthBookProject_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "growthbook_project" "test" {
  name = "tf-acc-ds-proj"
  description = "super description"
}
data "growthbook_project" "by_name" {
  name = growthbook_project.test.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.growthbook_project.by_name", "name", "tf-acc-ds-proj"),
					testCheckResourceAttrPrefix("data.growthbook_project.by_name", "id", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_project.by_name", "description", "super description"),
				),
			},
		},
	})
}
