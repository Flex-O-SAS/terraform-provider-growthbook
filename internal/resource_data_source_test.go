package internal_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccDataSourceConfig(id string) string {
	return `
data "growthbook_data_source" "by_id" {
  id = "` + id + `"
}
`
}

func TestAccDataSource_basic(t *testing.T) {
	t.Parallel()

	id := os.Getenv("GROWTHBOOK_TEST_DATASOURCE_ID")
	if id == "" {
		t.Skip("GROWTHBOOK_TEST_DATASOURCE_ID not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{{
			Config: testAccDataSourceConfig(id),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.growthbook_data_source.by_id", "id", id),
			),
		}},
	})
}
