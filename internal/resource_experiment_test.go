package internal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccExperimentConfig(name string) string {
	return `
resource "growthbook_experiment" "test" {
  name         = "` + name + `"
  tracking_key = "` + name + `-tracking"
  description  = "Acceptance test experiment"
  variations   = jsonencode([
    { key = "0", name = "Control" },
    { key = "1", name = "Variant" }
  ])
}
data "growthbook_experiment" "by_name" {
  name = growthbook_experiment.test.name
}
`
}

func testAccExperimentConfigUpdate(name string) string {
	return `
resource "growthbook_experiment" "test" {
  name         = "` + name + `"
  tracking_key = "` + name + `-tracking"
  description  = "Updated experiment"
  hypothesis   = "Variant increases conversions"
  variations   = jsonencode([
    { key = "0", name = "Control" },
    { key = "1", name = "Variant" }
  ])
  phases = jsonencode([{ name = "Phase 1" }])
}
data "growthbook_experiment" "by_name" {
  name = growthbook_experiment.test.name
}
`
}

func TestAccExperiment_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-acc-experiment-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExperimentConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_experiment.test", "name", name),
					resource.TestCheckResourceAttr("data.growthbook_experiment.by_name", "tracking_key", name+"-tracking"),
				),
			},
			{
				Config: testAccExperimentConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_experiment.test", "description", "Updated experiment"),
					resource.TestCheckResourceAttr(
						"data.growthbook_experiment.by_name", "hypothesis", "Variant increases conversions"),
				),
			},
		},
	})
}
