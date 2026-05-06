package internal_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccDimensionConfig(name, datasourceID string) string {
	return `
resource "growthbook_dimension" "test" {
  name            = "` + name + `"
  datasource_id   = "` + datasourceID + `"
  identifier_type = "user"
  query           = "select user_id, country from users"
  description     = "Acceptance test dimension"
}
data "growthbook_dimension" "by_name" {
  name = growthbook_dimension.test.name
}
`
}

func testAccDimensionConfigUpdate(name, datasourceID string) string {
	return `
resource "growthbook_dimension" "test" {
  name            = "` + name + `"
  datasource_id   = "` + datasourceID + `"
  identifier_type = "user"
  query           = "select user_id, region from users"
  description     = "Updated dimension"
}
data "growthbook_dimension" "by_name" {
  name = growthbook_dimension.test.name
}
`
}

func TestAccDimension_basic(t *testing.T) {
	t.Parallel()

	datasourceID := os.Getenv("GROWTHBOOK_TEST_DATASOURCE_ID")
	if datasourceID == "" {
		t.Skip("GROWTHBOOK_TEST_DATASOURCE_ID not set")
	}

	name := acctest.RandomWithPrefix("tf-acc-dimension-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDimensionConfig(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_dimension.test", "name", name),
					resource.TestCheckResourceAttr("data.growthbook_dimension.by_name", "query", "select user_id, country from users"),
				),
			},
			{
				Config: testAccDimensionConfigUpdate(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_dimension.test", "description", "Updated dimension"),
					resource.TestCheckResourceAttr("data.growthbook_dimension.by_name", "query", "select user_id, region from users"),
				),
			},
		},
	})
}
