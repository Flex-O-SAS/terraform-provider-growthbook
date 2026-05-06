package internal_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccMetricConfig(name, datasourceID string) string {
	return `
resource "growthbook_metric" "test" {
  name          = "` + name + `"
  datasource_id = "` + datasourceID + `"
  type          = "count"
  description   = "Acceptance test metric"
  behavior      = jsonencode({ goal = "increase" })
  sql           = jsonencode({ numerator = "count(*)" })
}
data "growthbook_metric" "by_name" {
  name = growthbook_metric.test.name
}
`
}

func testAccMetricConfigUpdate(name, datasourceID string) string {
	return `
resource "growthbook_metric" "test" {
  name          = "` + name + `"
  datasource_id = "` + datasourceID + `"
  type          = "count"
  description   = "Updated metric"
  behavior      = jsonencode({ goal = "decrease" })
  sql           = jsonencode({ numerator = "sum(total)" })
  archived      = true
}
data "growthbook_metric" "by_name" {
  name = growthbook_metric.test.name
}
`
}

func TestAccMetric_basic(t *testing.T) {
	t.Parallel()

	datasourceID := os.Getenv("GROWTHBOOK_TEST_DATASOURCE_ID")
	if datasourceID == "" {
		t.Skip("GROWTHBOOK_TEST_DATASOURCE_ID not set")
	}

	name := acctest.RandomWithPrefix("tf-acc-metric-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricConfig(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_metric.test", "name", name),
					resource.TestCheckResourceAttr("growthbook_metric.test", "datasource_id", datasourceID),
					resource.TestCheckResourceAttr("data.growthbook_metric.by_name", "name", name),
				),
			},
			{
				Config: testAccMetricConfigUpdate(name, datasourceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_metric.test", "description", "Updated metric"),
					resource.TestCheckResourceAttr("growthbook_metric.test", "archived", "true"),
					resource.TestCheckResourceAttr("data.growthbook_metric.by_name", "archived", "true"),
				),
			},
		},
	})
}
