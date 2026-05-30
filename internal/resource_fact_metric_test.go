package internal_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFactMetricConfig(name, datasourceID, factTableID string) string {
	return `
resource "growthbook_fact_metric" "test" {
  name        = "` + name + `"
  metric_type = "proportion"
  datasource  = "` + datasourceID + `"
  numerator   = jsonencode({ factTableId = "` + factTableID + `", column = "" })
  description = "Acceptance test fact metric"
}
data "growthbook_fact_metric" "by_name" {
  name = growthbook_fact_metric.test.name
}
`
}

func testAccFactMetricConfigUpdate(name, datasourceID, factTableID string) string {
	return `
resource "growthbook_fact_metric" "test" {
  name        = "` + name + `"
  metric_type = "proportion"
  datasource  = "` + datasourceID + `"
  numerator   = jsonencode({ factTableId = "` + factTableID + `", column = "" })
  denominator = jsonencode({ factTableId = "` + factTableID + `", column = "$$distinctUsers" })
  inverse     = true
  description = "Updated fact metric"
}
data "growthbook_fact_metric" "by_name" {
  name = growthbook_fact_metric.test.name
}
`
}

func TestAccFactMetric_basic(t *testing.T) {
	t.Parallel()

	datasourceID := os.Getenv("GROWTHBOOK_TEST_DATASOURCE_ID")
	if datasourceID == "" {
		t.Skip("GROWTHBOOK_TEST_DATASOURCE_ID not set")
	}

	factTableID := os.Getenv("GROWTHBOOK_TEST_FACT_TABLE_ID")
	if factTableID == "" {
		t.Skip("GROWTHBOOK_TEST_FACT_TABLE_ID not set")
	}

	name := acctest.RandomWithPrefix("tf-acc-fact-metric-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFactMetricConfig(name, datasourceID, factTableID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_fact_metric.test", "name", name),
					resource.TestCheckResourceAttr("growthbook_fact_metric.test", "metric_type", "proportion"),
					resource.TestCheckResourceAttr(
						"data.growthbook_fact_metric.by_name", "description", "Acceptance test fact metric"),
				),
			},
			{
				Config: testAccFactMetricConfigUpdate(name, datasourceID, factTableID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_fact_metric.test", "inverse", "true"),
					resource.TestCheckResourceAttr("data.growthbook_fact_metric.by_name", "description", "Updated fact metric"),
				),
			},
		},
	})
}
