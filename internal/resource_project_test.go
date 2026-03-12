package internal_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGrowthBookProject_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-acc-proj-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "growthbook_project" "test" {
  name        = "` + name + `"
  description = "super description"
}
data "growthbook_project" "by_name" {
  name = growthbook_project.test.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.growthbook_project.by_name", "name", name),
					testCheckResourceAttrPrefix("data.growthbook_project.by_name", "id", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_project.by_name", "description", "super description"),
				),
			},
		},
	})
}

// generateManyProjectsHCL generates HCL for N projects with unique names and descriptions.
func generateManyProjectsHCL(prefix string, n int) string {
	var b strings.Builder
	b.Grow(n * 150)

	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b,
			`resource "growthbook_project" "proj_%[1]d" {
				name        = "%[2]s%[1]d"
				description = "Acceptance test project %[1]d"
			}
		`, i, prefix)
	}

	return b.String()
}

func TestAccDataSourceGrowthBookProject_manyProjects(t *testing.T) {
	t.Parallel()

	namePrefix := acctest.RandomWithPrefix("tf-acc-proj-")
	config := generateManyProjectsHCL(namePrefix, 500) + `
# Pick one project to read as a data source
data "growthbook_project" "by_name1" {
  name = growthbook_project.proj_1.name
}
data "growthbook_project" "by_name2" {
  name = growthbook_project.proj_499.name
}
  data "growthbook_project" "by_name3" {
  name = growthbook_project.proj_200.name
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.growthbook_project.by_name1", "name", namePrefix+"1"),
					resource.TestCheckResourceAttr("data.growthbook_project.by_name1", "description", "Acceptance test project 1"),
					testCheckResourceAttrPrefix("data.growthbook_project.by_name1", "id", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_project.by_name2", "name", namePrefix+"499"),
					resource.TestCheckResourceAttr("data.growthbook_project.by_name2", "description", "Acceptance test project 499"),
					testCheckResourceAttrPrefix("data.growthbook_project.by_name2", "id", "prj"),
					resource.TestCheckResourceAttr("data.growthbook_project.by_name3", "name", namePrefix+"200"),
					resource.TestCheckResourceAttr("data.growthbook_project.by_name3", "description", "Acceptance test project 200"),
					testCheckResourceAttrPrefix("data.growthbook_project.by_name3", "id", "prj"),
				),
			},
		},
	})
}
