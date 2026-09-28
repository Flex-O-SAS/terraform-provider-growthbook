package internal_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSegmentConfig(name, datasourceID string) string {
	return `
resource "growthbook_segment" "test" {
  name            = "` + name + `"
  type            = "SQL"
  datasource_id   = "` + datasourceID + `"
  identifier_type = "user"
  query           = "select user_id from users"
  description     = "Acceptance test segment"
}
data "growthbook_segment" "by_name" {
  name = growthbook_segment.test.name
}
`
}

func testAccSegmentConfigUpdate(name, datasourceID string) string {
	return `
resource "growthbook_segment" "test" {
  name            = "` + name + `"
  type            = "SQL"
  datasource_id   = "` + datasourceID + `"
  identifier_type = "user"
  query           = "select distinct user_id from users"
  description     = "Updated segment"
}
data "growthbook_segment" "by_name" {
  name = growthbook_segment.test.name
}
`
}

func TestAccSegment_basic(t *testing.T) {
	t.Parallel()

	datasourceID := os.Getenv("GROWTHBOOK_TEST_DATASOURCE_ID")
	if datasourceID == "" {
		t.Skip("GROWTHBOOK_TEST_DATASOURCE_ID not set")
	}

	name := acctest.RandomWithPrefix("tf-acc-segment-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentConfig(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_segment.test", "name", name),
					resource.TestCheckResourceAttr("growthbook_segment.test", "datasource_id", datasourceID),
					resource.TestCheckResourceAttr("data.growthbook_segment.by_name", "name", name),
				),
			},
			{
				Config: testAccSegmentConfigUpdate(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_segment.test", "description", "Updated segment"),
					resource.TestCheckResourceAttr("data.growthbook_segment.by_name", "query", "select distinct user_id from users"),
				),
			},
		},
	})
}
