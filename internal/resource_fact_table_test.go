package internal_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFactTableConfig(name, datasourceID string) string {
	return `
resource "growthbook_fact_table" "test" {
  name          = "` + name + `"
  datasource    = "` + datasourceID + `"
  user_id_types = ["user"]
  sql           = "select user_id from events"
  description   = "Acceptance test fact table"
}
data "growthbook_fact_table" "by_name" {
  name = growthbook_fact_table.test.name
}
`
}

func testAccFactTableConfigUpdate(name, datasourceID string) string {
	return `
resource "growthbook_fact_table" "test" {
  name          = "` + name + `"
  datasource    = "` + datasourceID + `"
  user_id_types = ["user", "anonymous"]
  sql           = "select user_id, anonymous_id from events"
  description   = "Updated fact table"
}
data "growthbook_fact_table" "by_name" {
  name = growthbook_fact_table.test.name
}
`
}

func TestAccFactTable_basic(t *testing.T) {
	t.Parallel()

	datasourceID := os.Getenv("GROWTHBOOK_TEST_DATASOURCE_ID")
	if datasourceID == "" {
		t.Skip("GROWTHBOOK_TEST_DATASOURCE_ID not set")
	}

	name := acctest.RandomWithPrefix("tf-acc-fact-table-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFactTableConfig(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_fact_table.test", "name", name),
					resource.TestCheckResourceAttr("data.growthbook_fact_table.by_name", "datasource", datasourceID),
				),
			},
			{
				Config: testAccFactTableConfigUpdate(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_fact_table.test", "description", "Updated fact table"),
					resource.TestCheckResourceAttr("data.growthbook_fact_table.by_name", "user_id_types.#", "2"),
				),
			},
		},
	})
}
